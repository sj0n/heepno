package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/spf13/cobra"
)

var aaiCmd = &cobra.Command{
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

		fmt.Println("Transcribing...")

		start := time.Now()
		transcript, err := client.Transcripts.TranscribeFromReader(ctx, audioFile, &aai.TranscriptOptionalParams{
			LanguageCode: aai.TranscriptLanguageCode(Language),
			SpeechModel:  aai.SpeechModel(AssemblyAIModel),
			FormatText:   aai.Bool(true),
		})

		if err != nil {
			fmt.Println("AssemblyAI Error:", err)
			os.Exit(1)
		}
		
		elapsed := time.Since(start)
		fmt.Println(*transcript.Text)
		fmt.Printf("Transcribed in %s\n", elapsed)
	},
}

func init() {
	RootCmd.AddCommand(aaiCmd)

	aaiCmd.Flags().StringVarP(&Language, "language", "l", "", "Language to transcribe. See https://www.assemblyai.com/docs/getting-started/supported-languages for more details.")
	aaiCmd.Flags().StringVarP(&AssemblyAIModel, "model", "m", "best", "Model to use. <best|nano>")
}
