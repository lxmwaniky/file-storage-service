package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(databaseURL string) (*PostgresRepo, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL Database")

	return &PostgresRepo{DB: db}, nil
}