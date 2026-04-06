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

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{openai.NewClient(os.Getenv("OPENAI_API_KEY"))}
}

func (p *OpenAIProvider) Transcribe(ctx context.Context, file string, cfg config.Config) (*Result, error) {
	resp, err := p.CreateTranscription(ctx, openai.AudioRequest{
		Model:    cfg.Model,
		FilePath: file,
		Language: cfg.Language,
		Format:   audioFormat(cfg.Format),
	})
	if err != nil {
		return nil, fmt.Errorf("transcription error: %w", err)
	}
	return &Result{Text: resp.Text, Raw: resp}, nil
}

func (p *OpenAIProvider) Translate(ctx context.Context, file string, cfg config.Config) (*Result, error) {
	resp, err := p.CreateTranslation(ctx, openai.AudioRequest{
		FilePath: file,
		Model:    cfg.Model,
		Language: cfg.Language,
		Format:   audioFormat(cfg.Format),
	})

	if err != nil {
		return nil, fmt.Errorf("translation error: %w", err)
	}
	return &Result{Text: resp.Text, Raw: resp}, nil
}

func audioFormat(format string) openai.AudioResponseFormat {
	switch format {
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
