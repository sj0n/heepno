package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sj0n/heepno/internal/config"
)

type OpenAIProvider struct {
	*openai.Client
}

func createOpenAIOptions(file string, cfg config.Config) openai.AudioRequest {
	return openai.AudioRequest{
		FilePath: file,
		Model:    cfg.OpenaiModel,
		Language: cfg.Language,
		Format:   getAudioRequestFormat(cfg.Format),
	}
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (p *OpenAIProvider) Transcribe(ctx context.Context, file string, cfg config.Config) (any, error) {
	options := createOpenAIOptions(file, cfg)
	response, err := p.CreateTranscription(ctx, options)

	if err != nil {
		return nil, fmt.Errorf("transcription error: %w", err)
	}

	return response, nil
}

func (p *OpenAIProvider) Translate(ctx context.Context, file string, cfg config.Config) (any, error) {
	options := createOpenAIOptions(file, cfg)
	response, err := p.CreateTranslation(ctx, options)

	if err != nil {
		return nil, fmt.Errorf("translation error: %w", err)
	}

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
