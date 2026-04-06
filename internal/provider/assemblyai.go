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

func (p *AssemblyAIProvider) Transcribe(ctx context.Context, file string, cfg config.Config) (*Result, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("file error: %w", err)
	}
	defer f.Close()

	transcript, err := p.Transcripts.TranscribeFromReader(ctx, f, &assemblyai.TranscriptOptionalParams{
		LanguageCode: assemblyai.TranscriptLanguageCode(cfg.Language),
		FormatText:   assemblyai.Bool(true),
		SpeechModel:  assemblyai.SpeechModel(cfg.Model),
	})
	if err != nil {
		return nil, fmt.Errorf("transcription error: %w", err)
	}

	var text string
	if transcript.Text != nil {
		text = *transcript.Text
	}
	return &Result{Text: text, Raw: transcript}, nil
}

// Unused - satisfies interface, implement if needed
func (p *AssemblyAIProvider) Translate(ctx context.Context, file string, cfg config.Config) (*Result, error) {
	return nil, fmt.Errorf("translation not supported for AssemblyAI")
}
