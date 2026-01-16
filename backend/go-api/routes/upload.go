package routes

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB per request (Python enforces 50MB total)
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Extract user_id from header (set by auth middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusUnauthorized)
		return
	}

	// Prepare multipart request to Python service
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		http.Error(w, "Failed to create form file", http.StatusInternalServerError)
		return
	}

	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8000/upload",
		&buf,
	)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("user-id", userID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Python service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
