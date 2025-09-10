package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/sj0n/heepno/pkg/config"
	"github.com/sj0n/heepno/pkg/interfaces"
	"github.com/sj0n/heepno/pkg/shared"
	"github.com/spf13/cobra"
)

var (
	aaiCmd = &cobra.Command{
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
	client := interfaces.NewAssemblAIProvider()

	shared.PrintTranscriptionStatus("AssemblyAI", config.Global.AaiModel, config.Global.Language, "Transcribing...")

	start := time.Now()
	result, err := client.Transcribe(ctx, file)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}

	transcript := result.(assemblyai.Transcript)

	shared.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	var text string
	if transcript.Text != nil {
		text = *transcript.Text
	}

	if config.Global.Output != "" {
		if err := shared.Save(transcript, text, config.Global.Format, config.Global.Output); err != nil {
			return fmt.Errorf("File Error: %w", err)
		}
	} else {
		if err := shared.Print(transcript, text, config.Global.Format); err != nil {
			return fmt.Errorf("Print Error: %w", err)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(aaiCmd)

	aaiCmd.Flags().StringVarP(&config.Global.Language, "language", "l", "", "Language to transcribe. See https://www.assemblyai.com/docs/getting-started/supported-languages for more details.")
	aaiCmd.Flags().StringVarP(&config.Global.Format, "format", "f", "json", "Transcribe format. <json|text>")
	aaiCmd.Flags().StringVarP(&config.Global.Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	aaiCmd.Flags().StringVarP(&config.Global.AaiModel, "model", "m", "universal", "Model to use. <universal|slam-1(only support English.)>")
}
