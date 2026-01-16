// models/search.go
package models

type SearchRequest struct {
	Query      string        `json:"query" binding:"required"`
	MaxResults int           `json:"max_results"`
	Sources    SearchSources `json:"sources"`
}

type SearchSources struct {
	Personal   bool `json:"personal"`   // Auto-enabled if docs exist
	Brave      bool `json:"brave"`      // Default enabled
	DuckDuckGo bool `json:"duckduckgo"` // Default disabled
	Wikipedia  bool `json:"wikipedia"`  // Default disabled
}

type SearchResult struct {
	Title       string                 `json:"title"`
	URL         string                 `json:"url,omitempty"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"` // "brave", "duckduckgo", "wikipedia", "personal"
	Relevance   float64                `json:"relevance"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type SearchResponse struct {
	Query       string         `json:"query"`
	Results     []SearchResult `json:"results"`
	Count       int            `json:"count"`
	SourcesUsed []string       `json:"sources_used"`
}

// Helper methods
func (s *SearchSources) HasAnyEnabled() bool {
	return s.Personal || s.Brave || s.DuckDuckGo || s.Wikipedia
}

func (s *SearchSources) EnabledSources() []string {
	var enabled []string
	if s.Personal {
		enabled = append(enabled, "personal")
	}
	if s.Brave {
		enabled = append(enabled, "brave")
	}
	if s.DuckDuckGo {
		enabled = append(enabled, "duckduckgo")
	}
	if s.Wikipedia {
		enabled = append(enabled, "wikipedia")
	}
	return enabled
}

// Apply defaults: Brave=true, rest=false
func (s *SearchSources) ApplyDefaults(hasDocuments bool) {
	// If personal not explicitly set, enable if user has documents
	if hasDocuments && !s.Personal {
		s.Personal = true
	}

	// If no sources specified, enable Brave by default
	if !s.HasAnyEnabled() {
		s.Brave = true
	}
}
