package pkg

import (
	"context"
	// "encoding/json"
	"fmt"
	"os"
	"time"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"

	"github.com/spf13/cobra"
)

var (
	deepgramCmd = &cobra.Command{
		Use:   "deepgram <file>",
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
				Model:       DeepgramModel,
				Language:    Language,
				SmartFormat: true,
			}

			c := client.NewWithDefaults()
			dg := prerecorded.New(c)

			fmt.Println("Transcribing...")
			fmt.Println(DeepgramModel)
			start := time.Now()
			response, err := dg.FromFile(ctx, args[0], &options)

			if err != nil {
				fmt.Println("Deepgram Error: ", err)
				os.Exit(1)
			}

			// data, err := json.MarshalIndent(response, "", "  ")

			// if err != nil {
			// 	fmt.Println("Marshal Error: ", err)
			// 	os.Exit(1)
			// }

			elapsed := time.Since(start)
			fmt.Println(*&response.Results.Channels[0].Alternatives[0].Paragraphs)
			fmt.Println("Finished in ", elapsed)

		},
	}
)

func init() {
	RootCmd.AddCommand(deepgramCmd)

	deepgramCmd.Flags().StringVarP(&Language, "language", "l", "", "Language to transcribe")
	deepgramCmd.Flags().StringVarP(&DeepgramModel, "model", "m", "nova-2", "Model to use. See https://developers.deepgram.com/docs/models-languages-overview for more details.")
}
