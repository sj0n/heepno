package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/sj0n/heepno/internal/config"
	"github.com/sj0n/heepno/internal/provider"
	"github.com/spf13/cobra"
)

type ProviderConfig struct {
	Name             string
	APIKeyEnv        string
	ModelDefault     string
	SupportedFormats []string
	SupportTranslate bool
}

var providers = map[string]ProviderConfig{
	"OpenAI": {
		Name:             "OpenAI",
		APIKeyEnv:        "OPENAI_API_KEY",
		ModelDefault:     "whisper-1",
		SupportedFormats: []string{"json", "json_verbose", "text", "srt", "vtt"},
		SupportTranslate: true,
	},
	"AssemblyAI": {
		Name:             "AssemblyAI",
		APIKeyEnv:        "ASSEMBLYAI_API_KEY",
		ModelDefault:     "universal",
		SupportedFormats: []string{"json", "text"},
		SupportTranslate: false,
	},
	"Deepgram": {
		Name:             "Deepgram",
		APIKeyEnv:        "DEEPGRAM_API_KEY",
		ModelDefault:     "nova-3",
		SupportedFormats: []string{"json", "text"},
		SupportTranslate: false,
	},
}

func createCommand(shortCommand string, providerCfg ProviderConfig) *cobra.Command {
	var cliCfg config.Config

	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s <file>", shortCommand),
		Short: fmt.Sprintf("Transcribe audio using %s.", providerCfg.Name),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateFile(args[0]); err != nil {
				return err
			}

			if err := RequireAPIKey(providerCfg.APIKeyEnv, providerCfg.Name); err != nil {
				return err
			}

			var client provider.Provider
			switch providerCfg.Name {
			case "OpenAI":
				client = provider.NewOpenAIProvider()
			case "Deepgram":
				client = provider.NewDeepgramProvider()
			case "AssemblyAI":
				client = provider.NewAssemblyAIProvider()
			}

			return Run(context.Background(), providerCfg.Name, cliCfg.Model, cliCfg.Language, cliCfg.Format, cliCfg.Output, providerCfg.SupportedFormats,
				func(ctx context.Context) (*provider.Result, error) {
					if providerCfg.SupportTranslate {
						return client.Translate(ctx, args[0], cliCfg)
					}
					return client.Transcribe(ctx, args[0], cliCfg)
				})
		},
	}

	cmd.Flags().StringVarP(&cliCfg.Language, "language", "l", "", "Language code")
	cmd.Flags().StringVarP(&cliCfg.Model, "model", "m", providerCfg.ModelDefault, "Model to use")
	cmd.Flags().StringVarP(&cliCfg.Format, "format", "f", "json", fmt.Sprintf("Output format: %s", strings.Join(providerCfg.SupportedFormats, ", ")))
	cmd.Flags().StringVarP(&cliCfg.Output, "output", "o", "", "Output file")

	if providerCfg.SupportTranslate {
		cmd.Flags().BoolVarP(&cliCfg.Translate, "translate", "t", false, "Translate audio to English")
	}

	return cmd
}
