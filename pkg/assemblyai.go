package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/console"
	"github.com/sj0n/heepno/internal/output"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

var (
	aaiCfg     config.Config
	aaiCmd     = &cobra.Command{
		Use:   "aai <file>",
		Short: "Transcribe an audio file using AssemblyAI model.",
		Long:  "Transcribe an audio file using AssemblyAI model.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return assemblyAI(args[0], aaiCfg)
		},
	}
)

func assemblyAI(file string, cfg config.Config) error {
	if os.Getenv("ASSEMBLYAI_API_KEY") == "" {
		return fmt.Errorf("ASSEMBLYAI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client := provider.NewAssemblyAIProvider()

	console.PrintTranscriptionStatus("AssemblyAI", cfg.AaiModel, cfg.Language, "Transcribing...")

	start := time.Now()
	result, err := client.Transcribe(ctx, file, cfg)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}

	transcript := result.(assemblyai.Transcript)

	console.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	var text string
	if transcript.Text != nil {
		text = *transcript.Text
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
	RootCmd.AddCommand(aaiCmd)

	aaiCmd.Flags().StringVarP(&aaiCfg.Language, "language", "l", "", "Language to transcribe. See https://www.assemblyai.com/docs/getting-started/supported-languages for more details.")
	aaiCmd.Flags().StringVarP(&aaiCfg.Format, "format", "f", "json", "Transcribe format. <json|text>")
	aaiCmd.Flags().StringVarP(&aaiCfg.Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	aaiCmd.Flags().StringVarP(&aaiCfg.AaiModel, "model", "m", "universal", "Model to use. <universal|slam-1(only support English.)>")
}
