package domain

import "time"

type File struct {
	ID             string    `json:"id"`
	OriginalName   string    `json:"original_name"`
	StoredFileName string    `json:"stored_filename"`
	FileSize       int64     `json:"file_size"`
	MimeType       string    `json:"mime_type"`
	CreatedAt      time.Time `json:"created_at"`
}
