package httpdelivery

import (
	"net/http"

	"github.com/lxmwaniky/file-storage-service/internal/infrastructure/repository"
)

func NewRouter(fileRepo *repository.FileRepository) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.HandleFunc("POST /api/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		HandleUpload(w, r, fileRepo)
	})

	return mux
}
