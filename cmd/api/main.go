package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/lxmwaniky/file-storage-service/internal/config"
	httpdelivery "github.com/lxmwaniky/file-storage-service/internal/delivery/http"
	"github.com/lxmwaniky/file-storage-service/internal/infrastructure/repository"
	"github.com/lxmwaniky/file-storage-service/internal/infrastructure/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	pgRepo, err := repository.NewPostgresRepo(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	fileRepo := repository.NewFileRepository(pgRepo.DB)

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		slog.Error("Failed to load AWS SDK configuration", "error", err)
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	storageProvider := storage.NewS3Storage(s3Client, cfg.AWSS3BucketName)
	slog.Info("Initialized S3 storage provider", "bucket", cfg.AWSS3BucketName, "region", cfg.AWSRegion)

	router := httpdelivery.NewRouter(fileRepo, storageProvider)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdownChan
	slog.Warn("Received shutdown signal, initiating graceful shutdown...", "signal", sig.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown due to timeout/error", "error", err)
	} else {
		slog.Info("HTTP server stopped gracefully")
	}

	if err := pgRepo.DB.Close(); err != nil {
		slog.Error("Failed to close database connection safely", "error", err)
	} else {
		slog.Info("Database connection pool closed successfully")
	}

	slog.Info("Service exited cleanly")
}
