// services/search/duckduckgo_provider.go
package search

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DuckDuckGoProvider struct {
	lastRequestTime time.Time
	minInterval     time.Duration
	mu              sync.Mutex
	client          *http.Client
}

func NewDuckDuckGoProvider() *DuckDuckGoProvider {
	return &DuckDuckGoProvider{
		minInterval: 1 * time.Second, // Be respectful: 1 req/sec
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (d *DuckDuckGoProvider) Name() string {
	return "duckduckgo"
}

func (d *DuckDuckGoProvider) Priority() int {
	return 2 // Fallback after Brave
}

func (d *DuckDuckGoProvider) IsAvailable() bool {
	return true // Always available (no API key needed)
}

func (d *DuckDuckGoProvider) Search(ctx context.Context, query string, maxResults int) ([]Result, error) {
	// Rate limiting
	d.mu.Lock()
	elapsed := time.Since(d.lastRequestTime)
	if elapsed < d.minInterval {
		time.Sleep(d.minInterval - elapsed)
	}
	d.lastRequestTime = time.Now()
	d.mu.Unlock()

	// Build URL
	ddgURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", ddgURL, nil)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: failed to create request: %w", err)
	}

	// IMPORTANT: DuckDuckGo requires User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Make request
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("duckduckgo: failed to parse HTML: %w", err)
	}

	// Extract results
	var results []Result
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		if i >= maxResults {
			return
		}

		// Extract title and URL
		titleElem := s.Find(".result__a")
		title := strings.TrimSpace(titleElem.Text())
		rawURL, exists := titleElem.Attr("href")

		if !exists || title == "" {
			return
		}

		// Extract description
		description := strings.TrimSpace(s.Find(".result__snippet").Text())

		// Fix DuckDuckGo URL format
		fixedURL := d.fixDDGURL(rawURL)

		results = append(results, Result{
			Title:       title,
			URL:         fixedURL,
			Description: description,
			Relevance:   0.7, // Lower than Brave
		})
	})

	return results, nil
}

// fixDDGURL extracts the real URL from DuckDuckGo's redirect format
func (d *DuckDuckGoProvider) fixDDGURL(rawURL string) string {
	// DDG URLs are like: /l/?uddg=https%3A%2F%2Fexample.com
	if strings.HasPrefix(rawURL, "/l/?uddg=") {
		encodedURL := strings.TrimPrefix(rawURL, "/l/?uddg=")
		decoded, err := url.QueryUnescape(encodedURL)
		if err == nil {
			return decoded
		}
	}
	return rawURL
}
