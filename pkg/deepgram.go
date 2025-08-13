package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	itfs "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest/interfaces"

	"github.com/sj0n/heepno/pkg/config"
	"github.com/sj0n/heepno/pkg/interfaces"
	"github.com/sj0n/heepno/pkg/shared"
	"github.com/spf13/cobra"
)

var (
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
	dg := interfaces.NewDeepgramProvider()

	shared.PrintTranscriptionStatus("Deepgram", config.Global.DeepgramModel, config.Global.Language, "Transcribing...")

	start := time.Now()
	result, err := dg.Transcribe(ctx, file)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}

	shared.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	transcript := result.(*itfs.PreRecordedResponse)
	var text string
	if len(transcript.Results.Channels) > 0 &&
		len(transcript.Results.Channels[0].Alternatives) > 0 &&
		transcript.Results.Channels[0].Alternatives[0].Paragraphs != nil {
		text = strings.TrimSpace(transcript.Results.Channels[0].Alternatives[0].Paragraphs.Transcript)
	}

	if text == "" {
		shared.UpdateTranscriptionStatus("", fmt.Errorf("The model failed to transcribe text from the audio. Try using a different service instead."))
		return nil
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
	RootCmd.AddCommand(deepgramCmd)

	deepgramCmd.Flags().StringVarP(&config.Global.Language, "language", "l", "", "Language to transcribe")
	deepgramCmd.Flags().StringVarP(&config.Global.DeepgramModel, "model", "m", "nova-2", "Model to use. See https://developers.deepgram.com/docs/models-languages-overview for more details.")
	deepgramCmd.Flags().StringVarP(&config.Global.Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	deepgramCmd.Flags().StringVarP(&config.Global.Format, "format", "f", "json", "Transcribe format. <json|text>")
}
