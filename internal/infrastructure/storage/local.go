package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	return &LocalStorage{uploadDir: uploadDir}
}

func (s *LocalStorage) Save(ctx context.Context, key string, src io.Reader) error {
	if err := os.MkdirAll(s.uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local upload dir: %w", err)
	}

	dstPath := filepath.Join(s.uploadDir, key)
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file stream to disk: %w", err)
	}

	return nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	dstPath := filepath.Join(s.uploadDir, key)
	if err := os.Remove(dstPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete local file: %w", err)
	}
	return nil
}
