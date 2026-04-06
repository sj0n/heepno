package provider

import (
	"context"

	"github.com/sj0n/heepno/internal/config"
)

type Result struct {
	Text string
	Raw  any
}

type Transcriber interface {
	Transcribe(ctx context.Context, file string, cfg config.Config) (*Result, error)
}

type Translator interface {
	Translate(ctx context.Context, file string, cfg config.Config) (*Result, error)
}
