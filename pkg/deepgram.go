package pkg

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	itfs "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest/interfaces"

	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/console"
	"github.com/sj0n/heepno/internal/output"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

var (
	deepgramCfg     config.Config
	deepgramCmd     = &cobra.Command{
		Use:   "dg <file>",
		Short: "Transcribe an audio file using Deepgram models.",
		Long:  "Transcribe an audio file using Deepgram models.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deepgram(args[0], deepgramCfg)
		},
	}
)

func deepgram(file string, cfg config.Config) error {
	if os.Getenv("DEEPGRAM_API_KEY") == "" {
		return fmt.Errorf("DEEPGRAM_API_KEY environment variable is not set")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	dg := provider.NewDeepgramProvider()

	console.PrintTranscriptionStatus("Deepgram", cfg.DeepgramModel, cfg.Language, "Transcribing...")

	start := time.Now()
	result, err := dg.Transcribe(ctx, file, cfg)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}

	console.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	transcript := result.(*itfs.PreRecordedResponse)
	var text string
	if len(transcript.Results.Channels) > 0 &&
		len(transcript.Results.Channels[0].Alternatives) > 0 &&
		transcript.Results.Channels[0].Alternatives[0].Paragraphs != nil {
		text = strings.TrimSpace(transcript.Results.Channels[0].Alternatives[0].Paragraphs.Transcript)
	}

	if text == "" {
		console.UpdateTranscriptionStatus("no speech detected", nil)
		return nil
	}

	if cfg.Output != "" {
		if err := output.Save(transcript, text, cfg.Format, cfg.Output); err != nil {
			return fmt.Errorf("file error: %w", err)
		}
	} else {
		if err := output.Print(transcript, text, cfg.Format); err != nil {
			return fmt.Errorf("print error: %w", err)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deepgramCmd)

	deepgramCmd.Flags().StringVarP(&deepgramCfg.Language, "language", "l", "", "Language to transcribe")
	deepgramCmd.Flags().StringVarP(&deepgramCfg.DeepgramModel, "model", "m", "nova-2", "Model to use. See https://developers.deepgram.com/docs/models-languages-overview for more details.")
	deepgramCmd.Flags().StringVarP(&deepgramCfg.Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	deepgramCmd.Flags().StringVarP(&deepgramCfg.Format, "format", "f", "json", "Transcribe format. <json|text>")
}
