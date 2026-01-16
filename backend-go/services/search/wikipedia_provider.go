// services/search/wikipedia_provider.go
package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type WikipediaProvider struct {
	baseURL string
	client  *http.Client
}

func NewWikipediaProvider() *WikipediaProvider {
	return &WikipediaProvider{
		baseURL: "https://en.wikipedia.org/w/api.php",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *WikipediaProvider) Name() string {
	return "wikipedia"
}

func (w *WikipediaProvider) Priority() int {
	return 3 // Lower priority
}

func (w *WikipediaProvider) IsAvailable() bool {
	return true // Always available
}

func (w *WikipediaProvider) Search(ctx context.Context, query string, maxResults int) ([]Result, error) {
	// Build search URL
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srlimit", fmt.Sprintf("%d", maxResults))
	params.Set("format", "json")

	searchURL := fmt.Sprintf("%s?%s", w.baseURL, params.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: failed to create request: %w", err)
	}

	// Make request
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wikipedia: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var wikiResp wikiSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&wikiResp); err != nil {
		return nil, fmt.Errorf("wikipedia: failed to decode response: %w", err)
	}

	// Convert to standard results
	results := make([]Result, 0, len(wikiResp.Query.Search))
	for _, r := range wikiResp.Query.Search {
		// Build Wikipedia URL
		pageURL := fmt.Sprintf("https://en.wikipedia.org/wiki/%s",
			url.PathEscape(strings.ReplaceAll(r.Title, " ", "_")))

		// Clean HTML from snippet
		cleanDesc := w.cleanHTML(r.Snippet)

		results = append(results, Result{
			Title:       r.Title,
			URL:         pageURL,
			Description: cleanDesc,
			Relevance:   0.8, // Good quality
		})
	}

	return results, nil
}

// cleanHTML removes HTML tags from Wikipedia snippets
func (w *WikipediaProvider) cleanHTML(html string) string {
	// Remove <span class="searchmatch"> tags
	re := regexp.MustCompile(`<[^>]*>`)
	clean := re.ReplaceAllString(html, "")

	// Decode HTML entities
	clean = strings.ReplaceAll(clean, "&quot;", "\"")
	clean = strings.ReplaceAll(clean, "&amp;", "&")
	clean = strings.ReplaceAll(clean, "&lt;", "<")
	clean = strings.ReplaceAll(clean, "&gt;", ">")

	return strings.TrimSpace(clean)
}

// Response structure
type wikiSearchResponse struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}
