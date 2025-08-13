package pkg

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sj0n/heepno/pkg/config"
	"github.com/sj0n/heepno/pkg/interfaces"
	"github.com/sj0n/heepno/pkg/shared"
	"github.com/spf13/cobra"
)

var (
	oaiModel  string
	translate bool
	openaiCmd = &cobra.Command{
		Use:   "openai <file>",
		Short: "Transcribe an audio file using OpenAI model.",
		Long:  "Transcribe an audio file using OpenAI model.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return openAI(args[0])
		},
	}
)

func openAI(file string) error {
	if os.Getenv("OPENAI_API_KEY") == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client := interfaces.NewOpenAIProvider()

	var response any
	var err error

	if translate {
		response, err = client.Translate(ctx, file)
	} else {
		response, err = client.Transcribe(ctx, file)
	}

	if err != nil {
		return err
	}

	transcript := response.(*openai.AudioResponse)

	if config.Global.Output != "" {
		if err := shared.Save(transcript, transcript.Text, config.Global.Format, config.Global.Output); err != nil {
			return fmt.Errorf("File Error: %w", err)
		}
	} else {
		if err := shared.Print(transcript, transcript.Text, config.Global.Format); err != nil {
			return fmt.Errorf("Print Error: %w", err)
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(openaiCmd)

	openaiCmd.Flags().BoolVarP(&translate, "translate", "t", false, "Translate the audio file. Not setting this flag will transcribe the audio file.")
	openaiCmd.Flags().StringVarP(&config.Global.Language, "language", "l", "", "Language of the source audio. Setting this helps in accuracy and velocity.")
	openaiCmd.Flags().StringVarP(&config.Global.OpenaiModel, "model", "m", "whisper-1", "Model to use.")
	openaiCmd.Flags().StringVarP(&config.Global.Format, "format", "f", "json", "Format to use. json, text, srt, verbose_json, vtt")
	openaiCmd.Flags().StringVarP(&config.Global.Output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
}
