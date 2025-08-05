package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"

	"github.com/spf13/cobra"
)

var (
	dgModel     string
	deepgramCmd = &cobra.Command{
		Use:   "dg <file>",
		Short: "Transcribe an audio file using Deepgram models.",
		Long:  "Transcribe an audio file using Deepgram models.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if os.Getenv("DEEPGRAM_API_KEY") == "" {
				fmt.Println("Deepgram Error: Deepgram API key is not set")
				os.Exit(1)
			}

			ctx := context.Background()

			options := interfaces.PreRecordedTranscriptionOptions{
				Model:       dgModel,
				Language:    Language,
				SmartFormat: true,
			}

			c := client.NewWithDefaults()
			dg := prerecorded.New(c)

			fmt.Println("Model:", dgModel)
			fmt.Println("Language:", Language)
			fmt.Println("Transcribing...")

			start := time.Now()
			response, err := dg.FromFile(ctx, args[0], &options)
			elapsed := time.Since(start)

			if err != nil {
				fmt.Println("Deepgram Error: ", err)
				os.Exit(1)
			}
			fmt.Println("Finished in ", elapsed)

			if output != "" {
				cwd, err := os.Getwd()

				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}

				if format == "json" {
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

					fmt.Printf("Transcription saved to %s\\%s\n", cwd, fileName)
				} else {
					fileName, err := writeToFile(output, response.Results.Channels[0].Alternatives[0].Paragraphs.Transcript, "text")

					if err != nil {
						fmt.Println("File Error:", err)
						os.Exit(1)
					}

					fmt.Printf("Transcription saved to %s\\%s\n", cwd, fileName)
				}
			} else {
				if format == "json" {
					data, err := json.MarshalIndent(response, "", " ")

					if err != nil {
						fmt.Println("OpenAI Error:", err)
						os.Exit(1)
					}
					fmt.Println(string(data))
				} else {
					fmt.Println(response.Results.Channels[0].Alternatives[0].Paragraphs.Transcript)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(deepgramCmd)

	deepgramCmd.Flags().StringVarP(&Language, "language", "l", "", "Language to transcribe")
	deepgramCmd.Flags().StringVarP(&dgModel, "model", "m", "nova-2", "Model to use. See https://developers.deepgram.com/docs/models-languages-overview for more details.")
	deepgramCmd.Flags().StringVarP(&output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
	deepgramCmd.Flags().StringVarP(&format, "format", "f", "json", "Transcribe format. <json|text>")
}
