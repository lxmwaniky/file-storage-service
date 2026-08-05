package httpdelivery

import (
	"net/http"

	"github.com/lxmwaniky/file-storage-service/internal/domain"
	"github.com/lxmwaniky/file-storage-service/internal/infrastructure/repository"
)

func NewRouter(fileRepo *repository.FileRepository, storageService domain.StorageService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.HandleFunc("POST /api/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		HandleUpload(w, r, fileRepo, storageService)
	})

	mux.HandleFunc("GET /api/v1/files", func(w http.ResponseWriter, r *http.Request) {
		HandleListFiles(w, r, fileRepo)
	})

	mux.HandleFunc("GET /api/v1/files/{id}/download", func(w http.ResponseWriter, r *http.Request) {
		HandleDownloadFile(w, r, fileRepo)
	})

	return mux
}
