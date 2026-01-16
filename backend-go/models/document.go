// models/document.go
package models

import "time"

type Document struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id,omitempty"` // For later when we add auth
	OriginalFilename string    `json:"original_filename"`
	StoredFilename   string    `json:"stored_filename"`
	FileSize         int64     `json:"file_size"`
	MimeType         string    `json:"mime_type"`
	UploadedAt       time.Time `json:"uploaded_at"`
	ChunkCount       int       `json:"chunk_count,omitempty"`
}
