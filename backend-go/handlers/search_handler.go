// handlers/search_handler.go
package handlers

import (
	"context"
	"net/http"

	"github.com/Sanjay24-05/hybrid-search-engine/models"
	"github.com/Sanjay24-05/hybrid-search-engine/services/search"
	"github.com/Sanjay24-05/hybrid-search-engine/utils"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	aggregator   *search.Aggregator
	hasDocuments bool // Will be dynamic later with auth
}

func NewSearchHandler(aggregator *search.Aggregator) *SearchHandler {
	return &SearchHandler{
		aggregator:   aggregator,
		hasDocuments: false, // TODO: Check per user when auth added
	}
}

func (h *SearchHandler) HandleSearch(c *gin.Context) {
	var req models.SearchRequest

	// Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Validate and sanitize query
	cleanQuery, err := utils.ValidateSearchQuery(req.Query)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err, "Invalid query")
		return
	}
	req.Query = cleanQuery

	// Apply defaults
	req.Sources.ApplyDefaults(h.hasDocuments)

	// Validate max results
	req.MaxResults = utils.ValidateMaxResults(req.MaxResults)

	// Check if any source is enabled
	if !req.Sources.HasAnyEnabled() {
		utils.RespondError(c, http.StatusBadRequest,
			utils.ErrQueryEmpty,
			"At least one search source must be enabled")
		return
	}

	// Perform search
	results, err := h.aggregator.Search(
		context.Background(),
		req.Query,
		req.MaxResults,
		req.Sources.EnabledSources(),
	)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err, "Search failed")
		return
	}

	// Convert to response format
	searchResults := make([]models.SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = models.SearchResult{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Source:      inferSource(r.URL), // Helper function
			Relevance:   r.Relevance,
		}
	}

	// Build response
	response := models.SearchResponse{
		Query:       req.Query,
		Results:     searchResults,
		Count:       len(searchResults),
		SourcesUsed: req.Sources.EnabledSources(),
	}

	utils.RespondSuccess(c, http.StatusOK, response)
}

// Helper to infer source from URL
func inferSource(url string) string {
	if url == "" {
		return "personal"
	}
	// Simple heuristic - can be improved
	if len(url) > 0 {
		if url[0:4] == "http" {
			return "web"
		}
	}
	return "unknown"
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
