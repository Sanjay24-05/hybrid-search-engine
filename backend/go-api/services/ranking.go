package services

import (
	"hybrid-search-engine/go-api/models"
)

func Rank(results []models.SearchResult) []models.SearchResult {
	for i := range results {
		results[i].Score = float64(len(results[i].Title))
	}
	return results
}
