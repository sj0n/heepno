package provider

import (
	"context"
	"fmt"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"

	"github.com/sj0n/heepno/internal/config"
)

type DeepgramProvider struct {
	*api.Client
}

func NewDeepgramProvider() *DeepgramProvider {
	client.InitWithDefault()
	c := client.NewRESTWithDefaults()
	return &DeepgramProvider{
		api.New(c),
	}
}

func (p *DeepgramProvider) Transcribe(ctx context.Context, file string, cfg config.Config) (any, error) {
	response, err := p.FromFile(ctx, file, &interfacesv1.PreRecordedTranscriptionOptions{
		Model:       cfg.DeepgramModel,
		Language:    cfg.Language,
		SmartFormat: true,
	})
	if err != nil {
		return nil, fmt.Errorf("transcription error: %w", err)
	}

	return response, nil
}

func (p *DeepgramProvider) Translate(ctx context.Context, file string, cfg config.Config) (any, error) {
	return nil, fmt.Errorf("method error: not implemented")
}
