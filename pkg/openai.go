package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
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
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	options := openai.AudioRequest{
		FilePath: file,
		Model:    oaiModel,
		Language: Language,
		Format:   getAudioRequestFormat(format),
	}

	fmt.Println("Model:", oaiModel)
	if translate {
		return handleTranslation(ctx, client, options)
	}
	return handleTranscription(ctx, client, options)
}

func handleTranslation(ctx context.Context, client *openai.Client, options openai.AudioRequest) error {
	fmt.Println("+----------------+----------------------+")
	fmt.Printf("| %-14s | %-20s |\n", "Language", Language)
	fmt.Println("+----------------+----------------------+")
	fmt.Println("| Translating...|                       |")
	fmt.Println("+----------------+----------------------+")
	
	start := time.Now()

	response, err := client.CreateTranslation(ctx, options)
	if err != nil {
		return fmt.Errorf("Translation Error: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("| Translated in | %-20s |\n", elapsed)

	if output != "" {
		if err := shared.Save(response, response.Text, format, output); err != nil {
			return fmt.Errorf("File Error: %w", err)
		}
	} else {
		if err := shared.Print(response, response.Text, format); err != nil {
			return fmt.Errorf("Print Error: %w", err)
		}
	}

	return nil
}

func handleTranscription(ctx context.Context, client *openai.Client, options openai.AudioRequest) error {
	fmt.Println("+----------------+----------------------+")
	fmt.Printf("| %-14s | %-20s |\n", "Language", Language)
	fmt.Println("+----------------+----------------------+")
	fmt.Println("| Transcribing...|                      |")
	fmt.Println("+----------------+----------------------+")
	
	start := time.Now()

	response, err := client.CreateTranscription(ctx, options)
	if err != nil {
		return fmt.Errorf("Transcription Error: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("| Transcribed in | %-20s |\n", elapsed)

	if output != "" {
		if err := shared.Save(response, response.Text, format, output); err != nil {
			return fmt.Errorf("File Error: %w", err)
		}
	} else {
		if err := shared.Print(response, response.Text, format); err != nil {
			return fmt.Errorf("Print Error: %w", err)
		}
	}

	return nil
}

func getAudioRequestFormat(format string) openai.AudioResponseFormat {
	switch format {
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
	openaiCmd.Flags().StringVarP(&format, "format", "f", "json", "Format to use. json, text, srt, verbose_json, vtt")
	openaiCmd.Flags().StringVarP(&output, "output", "o", "", "The name of the output file. If not specified, the output will be printed to the console.")
}
