package httpdelivery

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

const MaxUploadSize = 10 << 20

func HandleUpload(w http.ResponseWriter, r *http.Request) {

	// Validation
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		slog.Warn("Failed to parse multipart form or file too large", "error", err)
		http.Error(w, "File is too large or invalid multipart form. Max size is 10MB", http.StatusRequestEntityTooLarge)
		return
	}

	// Retrieval
	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Warn("Missing 'file' field in request", "error", err)
		http.Error(w, "Field File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	slog.Info("Processing file upload",
		"filename", header.Filename,
		"size_bytes", header.Size)

	// Transformation
	fileBytes := make([]byte, 16)
	rand.Read(fileBytes)
	fileUUID := hex.EncodeToString(fileBytes)

	ext := filepath.Ext(header.Filename)
	safeFilename := fileUUID + ext

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		slog.Error("Failed to create upload directory", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join(uploadDir, safeFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		slog.Error("Failed to create destination file on disk", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Streaming
	if _, err := io.Copy(dst, file); err != nil {
		slog.Error("Failed to save contents", "error", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	slog.Info("File uploaded successfully", "saved_as", safeFilename)

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"status":"success", "message": "file uploaded successfully", "filename":"%s"}`, safeFilename)
}
