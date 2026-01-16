package services

import (
	"hybrid-search-engine/go-api/models"
)

func SerpAPISearch(query string) models.SearchResult {
	return models.SearchResult{
		Source: "serpapi",
		Title:  "SerpAPI Result",
		URL:    "https://google.com",
	}
}
