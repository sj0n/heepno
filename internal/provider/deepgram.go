package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"

	"github.com/sj0n/heepno/internal/config"
)

type DeepgramProvider struct {
	client *api.Client
}

var deepgramInitOnce sync.Once

func NewDeepgramProvider() *DeepgramProvider {
	deepgramInitOnce.Do(func() {
		client.InitWithDefault()
	})

	return &DeepgramProvider{
		api.New(client.NewRESTWithDefaults()),
	}
}

func (p *DeepgramProvider) Transcribe(ctx context.Context, file string, cfg config.Config) (*Result, error) {
	resp, err := p.client.FromFile(ctx, file, &interfacesv1.PreRecordedTranscriptionOptions{
		Model:       cfg.Model,
		Language:    cfg.Language,
		SmartFormat: true,
	})
	if err != nil {
		return nil, fmt.Errorf("transcription error: %w", err)
	}

	if len(resp.Results.Channels) == 0 ||
		len(resp.Results.Channels[0].Alternatives) == 0 ||
		resp.Results.Channels[0].Alternatives[0].Transcript == "" {
		return &Result{Text: "", Raw: resp}, nil
	}

	text := strings.TrimSpace(resp.Results.Channels[0].Alternatives[0].Transcript)
	return &Result{Text: text, Raw: resp}, nil
}
