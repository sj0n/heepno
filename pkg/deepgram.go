package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"

	"github.com/sj0n/heepno/pkg/shared"
	"github.com/spf13/cobra"
)

var (
	dgModel     string
	deepgramCmd = &cobra.Command{
		Use:   "dg <file>",
		Short: "Transcribe an audio file using Deepgram models.",
		Long:  "Transcribe an audio file using Deepgram models.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deepgram(args[0])
		},
	}
)

func deepgram(file string) error {
	if os.Getenv("DEEPGRAM_API_KEY") == "" {
		return fmt.Errorf("DEEPGRAM_API_KEY environment variable is not set")
	}

	ctx := context.Background()

	options := interfaces.PreRecordedTranscriptionOptions{
		Model:       dgModel,
		Language:    Language,
		SmartFormat: true,
	}

	c := client.NewWithDefaults()
	dg := prerecorded.New(c)

	fmt.Println("+----------------+----------------------+")
	fmt.Printf("| %-14s | %-20s |\n", "Model", dgModel)
	fmt.Printf("| %-14s | %-20s |\n", "Language", Language)
	fmt.Println("+----------------+----------------------+")
	fmt.Println("| Transcribing...|                      |")
	fmt.Println("+----------------+----------------------+")

	start := time.Now()
	response, err := dg.FromFile(ctx, file, &options)
	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("Transcription Error: %w", err)
	}

	fmt.Printf("| Transcribed in | %-20s |\n", elapsed)

	var text string
	if len(response.Results.Channels) > 0 &&
		len(response.Results.Channels[0].Alternatives) > 0 &&
		response.Results.Channels[0].Alternatives[0].Paragraphs != nil {
		text = strings.TrimSpace(response.Results.Channels[0].Alternatives[0].Paragraphs.Transcript)
	}

	if text == "" {
		fmt.Println("The model failed to transcribe text from the audio. Try using a different service instead.")
		return nil
	}

	if output != "" {
		if err := shared.Save(response, text, format, output); err != nil {
			return fmt.Errorf("File Error: %w", err)
		}
	} else {
		if err := shared.Print(response, text, format); err != nil {
			return fmt.Errorf("Print Error: %w", err)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deepgramCmd)

	deepgramCmd.Flags().StringVarP(&Language, "language", "l", "", "Language to transcribe")
	deepgramCmd.Flags().StringVarP(&dgModel, "model", "m", "nova-2", "Model to use. See https://developers.deepgram.com/docs/models-languages-overview for more details.")
	deepgramCmd.Flags().StringVarP(&output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	deepgramCmd.Flags().StringVarP(&format, "format", "f", "json", "Transcribe format. <json|text>")
}
