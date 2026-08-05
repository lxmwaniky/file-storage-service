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

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Starting server", "port", port)
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
