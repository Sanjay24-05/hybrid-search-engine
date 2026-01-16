package routes

import (
	"encoding/json"
	"hybrid-search-engine/go-api/services"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	sources := map[string]bool{
		"serpapi":    r.URL.Query().Get("serpapi") == "1",
		"duckduckgo": r.URL.Query().Get("duckduckgo") == "1",
		"wikipedia":  r.URL.Query().Get("wikipedia") == "1",
	}

	results := services.Aggregate(query, sources)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"web":       results,
		"documents": []interface{}{}, // placeholder
	})
}
