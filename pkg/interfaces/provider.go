package interfaces

import (
	"context"
)

type Provider interface {
	Transcribe(ctx context.Context, file string) (any, error)
	Translate(ctx context.Context, file string) (any, error)
}
