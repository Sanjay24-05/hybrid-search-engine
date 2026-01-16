// handlers/document_handler.go
package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sanjay24-05/hybrid-search-engine/models"
	"github.com/Sanjay24-05/hybrid-search-engine/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	uploadDir string
	maxSize   int64 // in bytes
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(uploadDir string, maxFileSizeMB int) *DocumentHandler {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("Warning: could not create upload directory: %v\n", err)
	}

	return &DocumentHandler{
		uploadDir: uploadDir,
		maxSize:   int64(maxFileSizeMB) * 1024 * 1024,
	}
}

// HandleUploadDocument handles document upload
// POST /api/documents
func (h *DocumentHandler) HandleUploadDocument(c *gin.Context) {
	// Parse multipart form with max size
	if err := c.Request.ParseMultipartForm(h.maxSize); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "Failed to parse form or file too large")
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "No file provided or file read error")
		return
	}
	defer file.Close()

	// Validate filename
	if file.Filename == "" {
		utils.RespondError(c, http.StatusBadRequest,
			utils.ErrQueryEmpty, "Filename cannot be empty")
		return
	}

	// Check file size
	if file.Size > h.maxSize {
		utils.RespondError(c, http.StatusRequestEntityTooLarge,
			fmt.Errorf("file size %d exceeds limit %d", file.Size, h.maxSize),
			"File is too large")
		return
	}

	// Generate unique ID and stored filename
	docID := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	storedFilename := docID + ext

	// Create file path
	filePath := filepath.Join(h.uploadDir, storedFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "Failed to create file")
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		// Clean up on error
		os.Remove(filePath)
		utils.RespondError(c, http.StatusInternalServerError, err, "Failed to save file")
		return
	}

	// Create document record
	doc := models.Document{
		ID:               docID,
		OriginalFilename: file.Filename,
		StoredFilename:   storedFilename,
		FileSize:         file.Size,
		MimeType:         file.Header.Get("Content-Type"),
		UploadedAt:       time.Now(),
	}

	// Return success response
	utils.RespondSuccess(c, http.StatusCreated, gin.H{
		"message":  "Document uploaded successfully",
		"document": doc,
	})
}

// HandleDownloadDocument handles document download
// GET /api/documents/:id
func (h *DocumentHandler) HandleDownloadDocument(c *gin.Context) {
	docID := c.Param("id")

	// Validate ID format (basic UUID validation)
	if docID == "" {
		utils.RespondError(c, http.StatusBadRequest,
			utils.ErrQueryEmpty, "Document ID is required")
		return
	}

	// In a real implementation, you would query the database
	// to get the stored filename and verify ownership.
	// For now, we'll use the ID directly (simplified).

	// Note: In production, retrieve document metadata from database
	// and verify user has access to it.

	utils.RespondError(c, http.StatusNotImplemented,
		fmt.Errorf("feature requires database integration"),
		"Document download not yet fully implemented")
}

// HandleDeleteDocument handles document deletion
// DELETE /api/documents/:id
func (h *DocumentHandler) HandleDeleteDocument(c *gin.Context) {
	docID := c.Param("id")

	// Validate ID
	if docID == "" {
		utils.RespondError(c, http.StatusBadRequest,
			utils.ErrQueryEmpty, "Document ID is required")
		return
	}

	// In a real implementation with database:
	// 1. Query database to get document metadata
	// 2. Verify user ownership
	// 3. Delete file from storage
	// 4. Delete database record

	utils.RespondError(c, http.StatusNotImplemented,
		fmt.Errorf("feature requires database integration"),
		"Document deletion not yet fully implemented")
}

// HandleListDocuments handles listing documents
// GET /api/documents
func (h *DocumentHandler) HandleListDocuments(c *gin.Context) {
	// In a real implementation with database:
	// 1. Query database for user's documents (with pagination)
	// 2. Return paginated list

	// For now, return empty list with structure
	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"message":   "Document listing requires authentication and database",
		"documents": []models.Document{},
		"total":     0,
		"page":      1,
		"limit":     10,
	})
}

// HandleSearchDocuments handles searching within documents
// POST /api/documents/search
func (h *DocumentHandler) HandleSearchDocuments(c *gin.Context) {
	var req struct {
		Query      string `json:"query" binding:"required"`
		MaxResults int    `json:"max_results,omitempty"`
	}

	// Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Validate query
	cleanQuery, err := utils.ValidateSearchQuery(req.Query)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "Invalid query")
		return
	}

	// Validate max results
	if req.MaxResults <= 0 || req.MaxResults > 100 {
		req.MaxResults = 20 // default
	}

	// In a real implementation:
	// 1. Search through user's documents
	// 2. Return relevant chunks/sections matching the query

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = ctx // Use context for timeout management

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"message": "Document search requires chunking and indexing",
		"query":   cleanQuery,
		"results": []gin.H{},
		"count":   0,
	})
}

// HandleDocumentStats handles retrieving document statistics
// GET /api/documents/:id/stats
func (h *DocumentHandler) HandleDocumentStats(c *gin.Context) {
	docID := c.Param("id")

	if docID == "" {
		utils.RespondError(c, http.StatusBadRequest,
			utils.ErrQueryEmpty, "Document ID is required")
		return
	}

	// In a real implementation:
	// Query database for document stats (chunk count, indexed status, etc.)

	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"document_id": docID,
		"message":     "Document statistics require database integration",
		"chunk_count": 0,
		"indexed":     false,
	})
}
