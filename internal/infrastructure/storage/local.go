package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lxmwaniky/file-storage-service/internal/domain"
)

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) domain.StorageService {
	return &LocalStorage{uploadDir: uploadDir}
}

func (s *LocalStorage) resolvePath(key string) (string, error) {
	dstPath := filepath.Join(s.uploadDir, key)

	uploadDirAbs, err := filepath.Abs(s.uploadDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve upload dir: %w", err)
	}
	dstPathAbs, err := filepath.Abs(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve destination path: %w", err)
	}

	if dstPathAbs != uploadDirAbs && !strings.HasPrefix(dstPathAbs, uploadDirAbs+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid key %q: resolves outside upload directory", key)
	}

	return dstPathAbs, nil
}

func (s *LocalStorage) Save(ctx context.Context, key string, src io.Reader) error {
	if err := os.MkdirAll(s.uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local upload dir: %w", err)
	}

	dstPath, err := s.resolvePath(key)
	if err != nil {
		return err
	}

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

func (s *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	srcPath, err := s.resolvePath(key)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrObjectNotFound
		}
		return nil, fmt.Errorf("failed to open file on disk: %w", err)
	}

	return f, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	dstPath, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(dstPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete local file: %w", err)
	}
	return nil
}

func (s *LocalStorage) GeneratePresignedUploadURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return "", fmt.Errorf("presigned urls are not supported by local storage provider")
}
