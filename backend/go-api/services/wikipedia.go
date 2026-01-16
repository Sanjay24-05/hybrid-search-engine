package services

import (
	"hybrid-search-engine/go-api/models"
)

func WikipediaSearch(query string) models.SearchResult {
	return models.SearchResult{
		Source: "wikipedia",
		Title:  "Wikipedia Result",
		URL:    "https://wikipedia.org",
	}
}
