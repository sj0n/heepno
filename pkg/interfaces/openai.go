package interfaces

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sj0n/heepno/pkg/config"
	"github.com/sj0n/heepno/pkg/shared"
)

type OpenAIProvider struct {
	*openai.Client
}

var OpenAIOptions openai.AudioRequest = openai.AudioRequest{
	FilePath: "",
	Model:    config.Global.OpenaiModel,
	Language: config.Global.Language,
	Format:   getAudioRequestFormat(config.Global.Format),
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (p *OpenAIProvider) Transcribe(ctx context.Context, file string) (any, error) {
	shared.PrintTranscriptionStatus("OpenAI", config.Global.OpenaiModel, config.Global.Language, "Transcribing...")

	start := time.Now()

	OpenAIOptions.FilePath = file
	response, err := p.CreateTranscription(ctx, OpenAIOptions)

	if err != nil {
		return nil, fmt.Errorf("Transcription Error: %w", err)
	}

	elapsed := time.Since(start)
	shared.UpdateTranscriptionStatus(fmt.Sprintf("Transcribed in %s", elapsed.Round(time.Second)), nil)

	return response, nil
}

func (p *OpenAIProvider) Translate(ctx context.Context, file string) (any, error) {
	shared.PrintTranscriptionStatus("OpenAI", config.Global.OpenaiModel, config.Global.Language, "Translating...")

	start := time.Now()

	OpenAIOptions.FilePath = file
	response, err := p.CreateTranslation(ctx, OpenAIOptions)

	if err != nil {
		return nil, fmt.Errorf("Translation Error: %w", err)
	}

	elapsed := time.Since(start)
	shared.UpdateTranscriptionStatus(fmt.Sprintf("Translated in %s", elapsed.Round(time.Second)), nil)

	return response, nil
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
