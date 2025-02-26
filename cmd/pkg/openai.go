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
	userFormat  string
	oaiModel    string
	translate   bool

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
				Model:    oaiModel,
				Language: Language,
				Format:   getAudioRequestFormat(userFormat),
			}

			fmt.Println("Model:", oaiModel)

			if translate {
				fmt.Println("Translating...")
				start := time.Now()
				response, err := client.CreateTranslation(ctx, options)

				if err != nil {
					fmt.Println("OpenAI Error:", err)
					os.Exit(1)
				}

				elapsed := time.Since(start)
				fmt.Printf("Finished in: %s\n", elapsed)

				if output != "" {
					fmt.Println("Saving to file...")
					cwd, err := os.Getwd()

					if err != nil {
						fmt.Println("Error:", err)
					}

					switch userFormat {
					case "json", "verbose_json":
						file, err := os.Create(output + ".json")

						if err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						defer file.Close()

						data, err := json.MarshalIndent(response, "", "  ")

						if err != nil {
							fmt.Println("JSON Error:", err)
							os.Exit(1)
						}

						if _, err := file.Write(data); err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						fmt.Printf("Transcription saved to %s\\%s\n", cwd, file.Name())
						os.Exit(0)
					case "text":
						file, err := os.Create(output + ".txt")

						if err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						defer file.Close()

						if _, err := file.WriteString(response.Text); err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						fmt.Printf("Transcription saved to %s\\%s\n", cwd, file.Name())
						os.Exit(0)
					case "srt":
						file, err := os.Create(output + ".srt")

						if err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						defer file.Close()

						if _, err := file.WriteString(response.Text); err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						fmt.Printf("Transcription saved to %s\\%s\n", cwd, file.Name())
						os.Exit(0)
					case "vtt":
						file, err := os.Create(output + ".vtt")
						if err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						defer file.Close()

						if _, err := file.WriteString(response.Text); err != nil {
							fmt.Println("File Error:", err)
							os.Exit(1)
						}
						fmt.Printf("Transcription saved to %s\\%s\n", cwd, file.Name())
						os.Exit(0)
					}
				} else {
					fmt.Println(response)
					os.Exit(0)
				}
			}

			fmt.Println("Language:", Language)
			fmt.Println("Transcribing...")
			start := time.Now()
			response, err := client.CreateTranscription(ctx, options)
			if err != nil {
				fmt.Println("OpenAI Error:", err)
				os.Exit(1)
			}

			elapsed := time.Since(start)
			fmt.Printf("Finished in: %s\n", elapsed)

			if output != "" {
				cwd, err := os.Getwd()

				if err != nil {
					fmt.Println("Error:", err)
				}

				switch userFormat {
				case "json", "verbose_json":
					data, err := json.MarshalIndent(response, "", "  ")

					if err != nil {
						fmt.Println("JSON Error:", err)
						os.Exit(1)
					}

					fileName, err := writeToFile(output, data, "json")

					if err != nil {
						fmt.Println("File Error:", err)
						os.Exit(1)
					}

					fmt.Printf("Transcription saved in %s\\%s\n", cwd, fileName)
				default:
					fileName, err := writeToFile(output, response.Text, userFormat)

					if err != nil {
						fmt.Println("File Error:", err)
						os.Exit(1)
					}

					fmt.Printf("Transcription saved in %s\\%s\n", cwd, fileName)
				}
			} else {
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
			}
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

	openaiCmd.Flags().BoolVarP(&translate, "translate", "t", false, "Translate the audio file. Not setting this flag will transcribe the audio file.")
	openaiCmd.Flags().StringVarP(&Language, "language", "l", "", "Language of the source audio. Setting this helps in accuracy and velocity.")
	openaiCmd.Flags().StringVarP(&oaiModel, "model", "m", "whisper-1", "Model to use.")
	openaiCmd.Flags().StringVarP(&userFormat, "format", "f", "json", "Format to use. json, text, srt, verbose_json, vtt")
	openaiCmd.Flags().StringVarP(&output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
}
