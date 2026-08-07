package domain

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrObjectNotFound = errors.New("object not found in storage")

type StorageService interface {
	Save(ctx context.Context, key string, src io.Reader) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GeneratePresignedUploadURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}
