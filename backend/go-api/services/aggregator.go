package services

import (
	"sync"

	"hybrid-search-engine/go-api/models"
)

func Aggregate(query string, sources map[string]bool) []models.SearchResult {
	results := make(chan models.SearchResult)
	var wg sync.WaitGroup

	for src, enabled := range sources {
		if !enabled {
			continue
		}

		wg.Add(1)
		go func(source string) {
			defer wg.Done()
			switch source {
			case "duckduckgo":
				results <- DuckDuckGoSearch(query)
			case "wikipedia":
				results <- WikipediaSearch(query)
			case "serpapi":
				results <- SerpAPISearch(query)
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var collected []models.SearchResult
	for r := range results {
		collected = append(collected, r)
	}

	return Rank(collected)
}
