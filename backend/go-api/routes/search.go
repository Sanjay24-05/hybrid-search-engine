package routes

import (
	"encoding/json"
	"hybrid-search-engine/go-api/services"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	userID := r.Header.Get("X-User-ID")

	webResults := services.Aggregate(query, map[string]bool{
		"serpapi":    r.URL.Query().Get("serpapi") == "1",
		"duckduckgo": r.URL.Query().Get("duckduckgo") == "1",
		"wikipedia":  r.URL.Query().Get("wikipedia") == "1",
	})

	docResults := []services.DocResult{}
	if r.URL.Query().Get("docs") == "1" {
		docs, err := services.SearchDocuments(userID, query)
		if err == nil {
			docResults = docs
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"documents": docResults,
		"web":       webResults,
	})
}
