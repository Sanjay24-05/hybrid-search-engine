package services

import (
	"hybrid-search-engine/go-api/models"
)

func DuckDuckGoSearch(query string) models.SearchResult {
	return models.SearchResult{
		Source: "duckduckgo",
		Title:  "DuckDuckGo Result",
		URL:    "https://duckduckgo.com",
	}
}
