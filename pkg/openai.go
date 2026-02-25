package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/console"
	"github.com/sj0n/heepno/internal/output"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

var (
	translate     bool
	openaiCfg     config.Config
	openaiCmd     = &cobra.Command{
		Use:   "openai <file>",
		Short: "Transcribe an audio file using OpenAI model.",
		Long:  "Transcribe an audio file using OpenAI model.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openAI(args[0], openaiCfg)
		},
	}
)

func openAI(file string, cfg config.Config) error {
	if os.Getenv("OPENAI_API_KEY") == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client := provider.NewOpenAIProvider()

	var response any
	var err error

	start := time.Now()

	if translate {
		console.PrintTranscriptionStatus("OpenAI", cfg.OpenaiModel, cfg.Language, "Translating...")
		response, err = client.Translate(ctx, file, cfg)
	} else {
		console.PrintTranscriptionStatus("OpenAI", cfg.OpenaiModel, cfg.Language, "Transcribing...")
		response, err = client.Transcribe(ctx, file, cfg)
	}

	elapsed := time.Since(start)

	if err != nil {
		return err
	}

	transcript := response.(openai.AudioResponse)

	console.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	if cfg.Output != "" {
		if err := output.Save(transcript, transcript.Text, cfg.Format, cfg.Output); err != nil {
			return fmt.Errorf("file error: %w", err)
		}
	} else {
		if err := output.Print(transcript, transcript.Text, cfg.Format); err != nil {
			return fmt.Errorf("print error: %w", err)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(openaiCmd)

	openaiCmd.Flags().BoolVarP(&translate, "translate", "t", false, "Translate the audio file. Not setting this flag will transcribe the audio file.")
	openaiCmd.Flags().StringVarP(&openaiCfg.Language, "language", "l", "", "Language of the source audio. Setting this helps in accuracy and velocity.")
	openaiCmd.Flags().StringVarP(&openaiCfg.OpenaiModel, "model", "m", "whisper-1", "Model to use.")
	openaiCmd.Flags().StringVarP(&openaiCfg.Format, "format", "f", "json", "Format to use. json, text, srt, verbose_json, vtt")
	openaiCmd.Flags().StringVarP(&openaiCfg.Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
}
