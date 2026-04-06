package pkg

import (
	"context"

	"github.com/sj0n/heepno/internal/cli"
	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

var (
	aaiCfg config.Config
	aaiCmd = &cobra.Command{
		Use:   "aai <file>",
		Short: "Transcribe audio using AssemblyAI.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.ValidateFile(args[0]); err != nil {
				return err
			}

			if err := cli.RequireAPIKey("ASSEMBLYAI_API_KEY", "AssemblyAI"); err != nil {
				return err
			}

			client := provider.NewAssemblyAIProvider()
			return cli.Run(context.Background(), "AssemblyAI", aaiCfg.Model, aaiCfg.Language, aaiCfg.Format, aaiCfg.Output,
				func(ctx context.Context) (*provider.Result, error) {
					return client.Transcribe(ctx, args[0], aaiCfg)
				})
		},
	}
)

func init() {
	RootCmd.AddCommand(aaiCmd)

	aaiCmd.Flags().StringVarP(&aaiCfg.Language, "language", "l", "", "Language code")
	aaiCmd.Flags().StringVarP(&aaiCfg.Model, "model", "m", "universal", "Model: universal, slam-1 (English only)")
	aaiCmd.Flags().StringVarP(&aaiCfg.Format, "format", "f", "json", "Output format: json, text")
	aaiCmd.Flags().StringVarP(&aaiCfg.Output, "output", "o", "", "Output file")
}
