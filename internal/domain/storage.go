package domain

import (
	"context"
	"io"
)

type StorageService interface {
	Save(ctx context.Context, key string, src io.Reader) error
	Delete(ctx context.Context, key string) error
	// Get(ctx context.Context, key string) (io.ReadCloser, error)
}
