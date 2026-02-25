package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/sj0n/heepno/internal/config"
)

type AssemblyAIProvider struct {
	*assemblyai.Client
}

func NewAssemblyAIProvider() *AssemblyAIProvider {
	return &AssemblyAIProvider{
		assemblyai.NewClientWithOptions(assemblyai.WithBaseURL("https://api.assemblyai.com/v2/transcript")),
	}
}

func (p *AssemblyAIProvider) Transcribe(ctx context.Context, file string, cfg config.Config) (any, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("file error: %w", err)
	}
	defer f.Close()

	transcript, err := p.Transcripts.TranscribeFromReader(ctx, f, &assemblyai.TranscriptOptionalParams{
		LanguageCode: assemblyai.TranscriptLanguageCode(cfg.Language),
		FormatText:   assemblyai.Bool(true),
		SpeechModel:  assemblyai.SpeechModel(cfg.AaiModel),
	})
	if err != nil {
		return nil, fmt.Errorf("transcription error: %w", err)
	}

	return transcript, nil
}

// TODO: Implement translation
func (p *AssemblyAIProvider) Translate(ctx context.Context, file string, cfg config.Config) (any, error) {
	return nil, fmt.Errorf("method error: not implemented")
}
