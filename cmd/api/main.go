package main

import (
	"log/slog"
	"net/http"
	"os"

	httpdelivery "github.com/lxmwaniky/file-storage-service/internal/delivery/http"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router := httpdelivery.NewRouter()

	port := ":8080"
	slog.Info("Starting server", "port", port)

	if err := http.ListenAndServe(port, router); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}