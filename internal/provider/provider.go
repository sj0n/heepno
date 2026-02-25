package provider

import (
	"context"

	"github.com/sj0n/heepno/internal/config"
)

type Provider interface {
	Transcribe(ctx context.Context, file string, cfg config.Config) (any, error)
	Translate(ctx context.Context, file string, cfg config.Config) (any, error)
}
