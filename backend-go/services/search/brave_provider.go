// services/search/brave_provider.go
package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BraveProvider struct {
	apiKey         string
	quotaLimit     int
	quotaRemaining int
	quotaResetTime time.Time
	mu             sync.Mutex
	client         *http.Client
}

func NewBraveProvider(apiKey string, dailyQuota int) *BraveProvider {
	return &BraveProvider{
		apiKey:         apiKey,
		quotaLimit:     dailyQuota,
		quotaRemaining: dailyQuota,
		quotaResetTime: getNextMidnight(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (b *BraveProvider) Name() string {
	return "brave"
}

func (b *BraveProvider) Priority() int {
	return 1 // Highest priority
}

func (b *BraveProvider) IsAvailable() bool {
	if b.apiKey == "" {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Reset quota if it's a new day
	if time.Now().After(b.quotaResetTime) {
		b.quotaRemaining = b.quotaLimit
		b.quotaResetTime = getNextMidnight()
	}

	return b.quotaRemaining > 0
}

func (b *BraveProvider) Search(ctx context.Context, query string, maxResults int) ([]Result, error) {
	if !b.IsAvailable() {
		return nil, errors.New("brave provider: quota exceeded or API key missing")
	}

	// Build URL
	baseURL := "https://api.search.brave.com/res/v1/web/search"
	params := url.Values{}
	params.Set("q", query)
	params.Set("count", fmt.Sprintf("%d", maxResults))

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("brave: failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Subscription-Token", b.apiKey)

	// Make request
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("brave: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brave: unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var braveResp braveResponse
	if err := json.NewDecoder(resp.Body).Decode(&braveResp); err != nil {
		return nil, fmt.Errorf("brave: failed to decode response: %w", err)
	}

	// Convert to standard results
	results := make([]Result, 0, len(braveResp.Web.Results))
	for _, r := range braveResp.Web.Results {
		results = append(results, Result{
			Title:       r.Title,
			URL:         r.URL,
			Description: r.Description,
			Relevance:   0.9, // Brave results are high quality
		})
	}

	// Consume quota
	b.consumeQuota()

	return results, nil
}

func (b *BraveProvider) consumeQuota() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.quotaRemaining > 0 {
		b.quotaRemaining--
	}
}

// Helper types for JSON parsing
type braveResponse struct {
	Web struct {
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"results"`
	} `json:"web"`
}

// Helper function
func getNextMidnight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}
