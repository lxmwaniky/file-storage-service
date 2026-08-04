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

func (r *FileRepository) FindAll(ctx context.Context) ([]domain.File, error) {
	query := `SELECT id, original_name, stored_filename, file_size, mime_type, created_at FROM files ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %w", err)
	}
	defer rows.Close()

	var files []domain.File

	for rows.Next() {
		var f domain.File
		if err := rows.Scan(&f.ID, &f.OriginalName, &f.StoredFileName, &f.FileSize, &f.MimeType, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}
		files = append(files, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating file rows: %w", err)
	}

	return files, nil
}

func (r *FileRepository) FindByID(ctx context.Context, id string) (*domain.File, error) {
	query := `SELECT id, original_name, stored_filename, file_size, mime_type, created_at FROM files WHERE id = $1`

	var f domain.File
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.OriginalName,
		&f.StoredFileName,
		&f.FileSize,
		&f.MimeType,
		&f.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query file by id: %w", err)
	}

	return &f, nil
}