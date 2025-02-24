package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	userFormat string

	openaiCmd = &cobra.Command{
		Use:   "openai <file>",
		Short: "Transcribe an audio file using OpenAI model.",
		Long:  "Transcribe an audio file using OpenAI model.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if os.Getenv("OPENAI_API_KEY") == "" {
				fmt.Println("OpenAI Error: OpenAI API key is not set")
				os.Exit(1)
			}

			ctx := context.Background()
			client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

			options := openai.AudioRequest{
				FilePath: args[0],
				Model:    OpenAIModel,
				Language: Language,
				Format:   getAudioRequestFormat(userFormat),
			}

			if Translate {
				fmt.Println("Translating...")
				start := time.Now()
				response, err := client.CreateTranslation(ctx, options)
				if err != nil {
					fmt.Println("OpenAI Error:", err)
					os.Exit(1)
				}
				elapsed := time.Since(start)
				fmt.Println(response)
				fmt.Printf("Finished in: %s\n", elapsed)
				os.Exit(0)
			}

			fmt.Println("Transcribing...")
			start := time.Now()
			response, err := client.CreateTranscription(ctx, options)
			if err != nil {
				fmt.Println("OpenAI Error:", err)
				os.Exit(1)
			}

			elapsed := time.Since(start)
			
			if userFormat == "json" || userFormat == "verbose_json" {
				data, err := json.MarshalIndent(response, "", "  ")
				if err != nil {
					fmt.Println("OpenAI Error:", err)
					os.Exit(1)
				}
				fmt.Println(string(data))
			} else {
				fmt.Println(response.Text)
			}

			fmt.Printf("Transcribed in: %s\n", elapsed)
			os.Exit(0)
		},
	}
)

func getAudioRequestFormat(userFormat string) openai.AudioResponseFormat {
	switch userFormat {
	case "json":
		return openai.AudioResponseFormatJSON
	case "text":
		return openai.AudioResponseFormatText
	case "srt":
		return openai.AudioResponseFormatSRT
	case "verbose_json":
		return openai.AudioResponseFormatVerboseJSON
	case "vtt":
		return openai.AudioResponseFormatVTT
	default:
		return openai.AudioResponseFormatJSON
	}
}

func init() {
	RootCmd.AddCommand(openaiCmd)

	openaiCmd.Flags().BoolVarP(&Translate, "translate", "t", false, "Translate the audio file. Not setting this flag will transcribe the audio file.")
	openaiCmd.Flags().StringVarP(&Language, "language", "l", "", "Language of the source audio. Setting this helps in accuracy and velocity.")
	openaiCmd.Flags().StringVarP(&OpenAIModel, "model", "m", "whisper-1", "Model to use.")
	openaiCmd.Flags().StringVarP(&userFormat, "format", "f", "json", "Format to use. json, text, srt, verbose_json, vtt")
}
