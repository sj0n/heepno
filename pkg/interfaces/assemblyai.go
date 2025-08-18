package interfaces

import (
	"context"
	"fmt"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/sj0n/heepno/pkg/config"
)

type AssemblyAIProvider struct {
	*assemblyai.Client
}

func NewAssemblAIProvider() *AssemblyAIProvider {
	return &AssemblyAIProvider{
		assemblyai.NewClientWithOptions(assemblyai.WithBaseURL("https://api.assemblyai.com/v2/transcript")),
	}
}

func (p *AssemblyAIProvider) Transcribe(ctx context.Context, file string) (any, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("File Error: %w", err)
	}
	defer f.Close()

	transcript, err := p.Transcripts.TranscribeFromReader(ctx, f, &assemblyai.TranscriptOptionalParams{
		LanguageCode: assemblyai.TranscriptLanguageCode(config.Global.Language),
		FormatText:   assemblyai.Bool(true),
		SpeechModel:  assemblyai.SpeechModel(config.Global.AaiModel),
	})
	if err != nil {
		return nil, fmt.Errorf("Transcription Error: %w", err)
	}

	return transcript, nil
}

// TODO: Implement translation
func (p *AssemblyAIProvider) Translate(ctx context.Context, file string) (any, error) {
	return nil, fmt.Errorf("Method Error: Not Implemented.")
}
