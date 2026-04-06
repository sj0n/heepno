package pkg

import (
	"context"

	"github.com/sj0n/heepno/internal/cli"
	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

var (
	translate bool
	openaiCfg config.Config
	openaiCmd = &cobra.Command{
		Use:   "openai <file>",
		Short: "Transcribe audio using OpenAI Whisper.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.ValidateFile(args[0]); err != nil {
				return err
			}

			if err := cli.RequireAPIKey("OPENAI_API_KEY", "OpenAI"); err != nil {
				return err
			}

			client := provider.NewOpenAIProvider()
			return cli.Run(context.Background(), "OpenAI", openaiCfg.Model, openaiCfg.Language, openaiCfg.Format, openaiCfg.Output,
				func(ctx context.Context) (*provider.Result, error) {
					if translate {
						return client.Translate(ctx, args[0], openaiCfg)
					}
					return client.Transcribe(ctx, args[0], openaiCfg)
				})
		},
	}
)

func init() {
	RootCmd.AddCommand(openaiCmd)

	openaiCmd.Flags().BoolVarP(&translate, "translate", "t", false, "Translate audio to English")
	openaiCmd.Flags().StringVarP(&openaiCfg.Language, "language", "l", "", "Source language")
	openaiCmd.Flags().StringVarP(&openaiCfg.Model, "model", "m", "whisper-1", "Model to use")
	openaiCmd.Flags().StringVarP(&openaiCfg.Format, "format", "f", "json", "Output format: json, text, srt, verbose_json, vtt")
	openaiCmd.Flags().StringVarP(&openaiCfg.Output, "output", "o", "", "Output file")
}
