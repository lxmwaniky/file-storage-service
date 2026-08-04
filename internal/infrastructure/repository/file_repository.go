package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lxmwaniky/file-storage-service/internal/domain"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Save(ctx context.Context, file *domain.File) error {
	query := `
		INSERT INTO files (id, original_name, stored_filename, file_size, mime_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		file.ID,
		file.OriginalName,
		file.StoredFileName,
		file.FileSize,
		file.MimeType,
		file.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert file metadata: %w", err)
	}

	return nil
}
