// services/search/aggregator.go
package search

import (
	"context"
	"sort"
	"sync"
)

type Aggregator struct {
	providers []SearchProvider
	mu        sync.RWMutex
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		providers: make([]SearchProvider, 0),
	}
}

func (a *Aggregator) RegisterProvider(provider SearchProvider) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.providers = append(a.providers, provider)

	// Sort by priority
	sort.Slice(a.providers, func(i, j int) bool {
		return a.providers[i].Priority() < a.providers[j].Priority()
	})
}

func (a *Aggregator) Search(ctx context.Context, query string, maxResults int, enabledSources []string) ([]Result, error) {
	a.mu.RLock()
	providers := a.providers
	a.mu.RUnlock()

	// Filter providers based on enabled sources
	activeProviders := a.filterProviders(providers, enabledSources)

	if len(activeProviders) == 0 {
		return []Result{}, nil
	}

	// Search concurrently
	type providerResult struct {
		results  []Result
		err      error
		provider string
	}

	resultChan := make(chan providerResult, len(activeProviders))

	// Launch goroutines
	for _, provider := range activeProviders {
		go func(p SearchProvider) {
			if !p.IsAvailable() {
				resultChan <- providerResult{provider: p.Name()}
				return
			}

			results, err := p.Search(ctx, query, maxResults)
			resultChan <- providerResult{
				results:  results,
				err:      err,
				provider: p.Name(),
			}
		}(provider)
	}

	// Collect results
	allResults := []Result{}
	for i := 0; i < len(activeProviders); i++ {
		result := <-resultChan

		if result.err != nil {
			// Log error but continue
			continue
		}

		allResults = append(allResults, result.results...)
	}

	// Deduplicate
	deduplicated := a.deduplicateResults(allResults)

	// Rank
	ranked := a.rankResults(deduplicated)

	// Limit results
	if len(ranked) > maxResults {
		ranked = ranked[:maxResults]
	}

	return ranked, nil
}

func (a *Aggregator) filterProviders(providers []SearchProvider, enabledSources []string) []SearchProvider {
	if len(enabledSources) == 0 {
		return providers
	}

	enabled := make(map[string]bool)
	for _, source := range enabledSources {
		enabled[source] = true
	}

	filtered := []SearchProvider{}
	for _, provider := range providers {
		if enabled[provider.Name()] {
			filtered = append(filtered, provider)
		}
	}

	return filtered
}

func (a *Aggregator) deduplicateResults(results []Result) []Result {
	seen := make(map[string]Result)

	for _, result := range results {
		existing, exists := seen[result.URL]

		if !exists {
			seen[result.URL] = result
		} else if result.Relevance > existing.Relevance {
			// Keep higher relevance
			seen[result.URL] = result
		}
	}

	// Convert map back to slice
	unique := make([]Result, 0, len(seen))
	for _, result := range seen {
		unique = append(unique, result)
	}

	return unique
}

func (a *Aggregator) rankResults(results []Result) []Result {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})
	return results
}
