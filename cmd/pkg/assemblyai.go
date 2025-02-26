package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/spf13/cobra"
)

var (
	aaiOutput string
	aaiModel  string
	aaiFormat string
	aaiCmd    = &cobra.Command{
		Use:   "aai <file>",
		Short: "Transcribe an audio file using AssemblyAI model.",
		Long:  "Transcribe an audio file using AssemblyAI model.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if os.Getenv("ASSEMBLYAI_API_KEY") == "" {
				fmt.Println("AssemblyAI Error: AssemblyAI API key is not set")
				os.Exit(1)
			}

			ctx := context.Background()
			client := aai.NewClient(os.Getenv("ASSEMBLYAI_API_KEY"))

			audioFile, err := os.Open(args[0])

			if err != nil {
				fmt.Println("File Error:", err)
				os.Exit(1)
			}
			defer audioFile.Close()

			fmt.Println("Model:", aaiModel)
			fmt.Println("Language:", Language)
			fmt.Println("Transcribing...")

			start := time.Now()
			transcript, err := client.Transcripts.TranscribeFromReader(ctx, audioFile, &aai.TranscriptOptionalParams{
				LanguageCode: aai.TranscriptLanguageCode(Language),
				SpeechModel:  aai.SpeechModel(aaiModel),
				FormatText:   aai.Bool(true),
			})
			elapsed := time.Since(start)

			if err != nil {
				fmt.Println("AssemblyAI Error:", err)
				os.Exit(1)
			}

			fmt.Printf("Transcribed in %s\n", elapsed)

			if aaiOutput != "" {
				cwd, err := os.Getwd()

				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}

				if aaiFormat == "json" {
					data, err := json.MarshalIndent(transcript, "", "  ")

					if err != nil {
						fmt.Println("JSON Error:", err)
						os.Exit(1)
					}

					fileName, err := writeToFile(aaiOutput, data, "json")

					if err != nil {
						fmt.Println("File Error:", err)
						os.Exit(1)
					}

					fmt.Printf("Transcription saved to %s\\%s\n", cwd, fileName)
				} else {
					fileName, err := writeToFile(aaiOutput, *transcript.Text, "text")

					if err != nil {
						fmt.Println("File Error:", err)
						os.Exit(1)
					}

					fmt.Printf("Transcription saved to %s\\%s\n", cwd, fileName)
				}
			} else {
				if aaiFormat == "json" {
					data, err := json.MarshalIndent(transcript, "", "  ")
					if err != nil {
						fmt.Println("OpenAI Error:", err)
						os.Exit(1)
					}
					fmt.Println(string(data))
				} else {
					fmt.Println(*transcript.Text)
				}
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(aaiCmd)

	aaiCmd.Flags().StringVarP(&Language, "language", "l", "", "Language to transcribe. See https://www.assemblyai.com/docs/getting-started/supported-languages for more details.")
	aaiCmd.Flags().StringVarP(&aaiModel, "model", "m", "best", "Model to use. <best|nano>")
	aaiCmd.Flags().StringVarP(&aaiFormat, "format", "f", "json", "Transcribe format. <json|text>")
	aaiCmd.Flags().StringVarP(&output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
}
