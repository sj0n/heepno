package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/sj0n/heepno/pkg/shared"
	"github.com/spf13/cobra"
)

var (
	aaiModel string
	aaiCmd   = &cobra.Command{
		Use:   "aai <file>",
		Short: "Transcribe an audio file using AssemblyAI model.",
		Long:  "Transcribe an audio file using AssemblyAI model.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return assemblyAI(args[0])
		},
	}
)

func assemblyAI(file string) error {
	if os.Getenv("ASSEMBLYAI_API_KEY") == "" {
		return fmt.Errorf("ASSEMBLYAI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client := aai.NewClientWithOptions(aai.WithBaseURL("https://api.assemblyai.com/v2/transcript"))

	audioFile, err := os.Open(file)

	if err != nil {
		return fmt.Errorf("File Error: %w", err)
	}
	defer audioFile.Close()

	shared.PrintTranscriptionStatus("AssemblyAI", aaiModel, Language, "Transcribing...")

	start := time.Now()
	transcript, err := aaiTranscribe(ctx, client, audioFile)
	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	shared.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	var text string
	if transcript.Text != nil {
		text = *transcript.Text
	}

	if Output != "" {
		if err := shared.Save(transcript, text, Format, Output); err != nil {
			return fmt.Errorf("File Error: %w", err)
		}
	} else {
		if err := shared.Print(transcript, text, Format); err != nil {
			return fmt.Errorf("Print Error: %w", err)
		}
	}

	return nil
}

func aaiTranscribe(ctx context.Context, client *aai.Client, audioFile *os.File) (*aai.Transcript, error) {
	transcript, err := client.Transcripts.TranscribeFromReader(ctx, audioFile, &aai.TranscriptOptionalParams{
		LanguageCode: aai.TranscriptLanguageCode(Language),
		SpeechModel:  aai.SpeechModel(aaiModel),
		FormatText:   aai.Bool(true),
	})

	if err != nil {
		return nil, fmt.Errorf("Transcription Error: %w", err)
	}

	return &transcript, nil
}

func init() {
	RootCmd.AddCommand(aaiCmd)

	aaiCmd.Flags().StringVarP(&Language, "language", "l", "", "Language to transcribe. See https://www.assemblyai.com/docs/getting-started/supported-languages for more details.")
	aaiCmd.Flags().StringVarP(&Format, "format", "f", "json", "Transcribe format. <json|text>")
	aaiCmd.Flags().StringVarP(&Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	aaiCmd.Flags().StringVarP(&aaiModel, "model", "m", "universal", "Model to use. <universal|slam-1(only support English.)>")
}
