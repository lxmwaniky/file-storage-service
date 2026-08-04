package main

import (
	"log/slog"
	"net/http"
	"os"

	httpdelivery "github.com/lxmwaniky/file-storage-service/internal/delivery/http"
	"github.com/lxmwaniky/file-storage-service/internal/infrastructure/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbConnString := "postgres://alex:wantam@localhost:5432/file_storage?sslmode=disable"

	pgRepo, err := repository.NewPostgresRepo(dbConnString)

	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pgRepo.DB.Close()

	fileRepo := repository.NewFileRepository(pgRepo.DB)

	router := httpdelivery.NewRouter(fileRepo)

	port := ":8080"
	slog.Info("Starting server", "port", port)

	if err := http.ListenAndServe(port, router); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}