package httpdelivery

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/lxmwaniky/file-storage-service/internal/domain"
	"github.com/lxmwaniky/file-storage-service/internal/infrastructure/repository"
)

func HandleListFiles(w http.ResponseWriter, r *http.Request, fileRepo *repository.FileRepository) {
	files, err := fileRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if files == nil {
		files = []domain.File{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

func HandleDownloadFile(w http.ResponseWriter, r *http.Request, fileRepo *repository.FileRepository, storageService domain.StorageService) {
	id := r.PathValue("id")
	if id == "" {
		slog.Warn("Download attempt with missing file ID")
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	fileMeta, err := fileRepo.FindByID(r.Context(), id)
	if err != nil {
		slog.Error("Database error during file lookup for download", "file_id", id, "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if fileMeta == nil {
		slog.Warn("Download requested for non-existent file ID in database", "file_id", id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"file not found"}`))
		return
	}

	targetFile, err := storageService.Get(r.Context(), fileMeta.StoredFileName)
	if err != nil {
		if errors.Is(err, domain.ErrObjectNotFound) {
			slog.Error("CRITICAL: Database record exists, but object is missing from storage!",
				"file_id", id,
				"stored_filename", fileMeta.StoredFileName,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"file missing from storage provider"}`))
			return
		}
		slog.Error("Failed to retrieve object from storage", "file_id", id, "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	w.Header().Set("Content-Type", fileMeta.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileMeta.OriginalName+"\"")

	slog.Info("Streaming file download to client",
		"file_id", fileMeta.ID,
		"original_name", fileMeta.OriginalName,
		"size_bytes", fileMeta.FileSize,
	)

	_, _ = io.Copy(w, targetFile)
}
