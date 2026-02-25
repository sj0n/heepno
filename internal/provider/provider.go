package provider

import (
	"context"

	"github.com/sj0n/heepno/internal/config"
)

type Transcriber interface {
	Transcribe(ctx context.Context, file string, cfg config.Config) (any, error)
}

type Translator interface {
	Translate(ctx context.Context, file string, cfg config.Config) (any, error)
}
