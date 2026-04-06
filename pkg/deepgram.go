package pkg

import (
	"context"

	"github.com/sj0n/heepno/internal/cli"
	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

var (
	deepgramCfg config.Config
	deepgramCmd = &cobra.Command{
		Use:   "dg <file>",
		Short: "Transcribe audio using Deepgram.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.ValidateFile(args[0]); err != nil {
				return err
			}

			if err := cli.RequireAPIKey("DEEPGRAM_API_KEY", "Deepgram"); err != nil {
				return err
			}

			client := provider.NewDeepgramProvider()
			return cli.Run(context.Background(), "Deepgram", deepgramCfg.Model, deepgramCfg.Language, deepgramCfg.Format, deepgramCfg.Output,
				func(ctx context.Context) (*provider.Result, error) {
					return client.Transcribe(ctx, args[0], deepgramCfg)
				})
		},
	}
)

func init() {
	RootCmd.AddCommand(deepgramCmd)

	deepgramCmd.Flags().StringVarP(&deepgramCfg.Language, "language", "l", "", "Language code")
	deepgramCmd.Flags().StringVarP(&deepgramCfg.Model, "model", "m", "nova-2", "Model to use")
	deepgramCmd.Flags().StringVarP(&deepgramCfg.Format, "format", "f", "json", "Output format: json, text")
	deepgramCmd.Flags().StringVarP(&deepgramCfg.Output, "output", "o", "", "Output file")
}
