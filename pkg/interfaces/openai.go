package interfaces

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/sj0n/heepno/pkg/config"
)

type OpenAIProvider struct {
	*openai.Client
}

func createOpenAIOptions(file string) openai.AudioRequest {
	return openai.AudioRequest{
		FilePath: file,
		Model:    config.Global.OpenaiModel,
		Language: config.Global.Language,
		Format:   getAudioRequestFormat(config.Global.Format),
	}
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (p *OpenAIProvider) Transcribe(ctx context.Context, file string) (any, error) {
	options := createOpenAIOptions(file)
	response, err := p.CreateTranscription(ctx, options)

	if err != nil {
		return nil, fmt.Errorf("error transcribing: %w", err)
	}

	return response, nil
}

func (p *OpenAIProvider) Translate(ctx context.Context, file string) (any, error) {
	options := createOpenAIOptions(file)
	response, err := p.CreateTranslation(ctx, options)

	if err != nil {
		return nil, fmt.Errorf("error translating: %w", err)
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
