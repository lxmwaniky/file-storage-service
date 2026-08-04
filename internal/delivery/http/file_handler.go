package httpdelivery

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

func HandleDownloadFile(w http.ResponseWriter, r *http.Request, fileRepo *repository.FileRepository) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing file ID", http.StatusBadRequest)
		return
	}

	fileMeta, err := fileRepo.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if fileMeta == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := filepath.Join("./uploads", fileMeta.StoredFileName)
	targetFile, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File missing from storage provider", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	w.Header().Set("Content-Type", fileMeta.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileMeta.OriginalName+"\"")

	_, _ = io.Copy(w, targetFile)
}