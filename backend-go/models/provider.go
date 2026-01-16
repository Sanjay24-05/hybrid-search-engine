// models/provider.go
package models

import "time"

// ProviderConfig holds configuration for a search provider
type ProviderConfig struct {
	Name       string
	APIKey     string
	Enabled    bool
	Priority   int       // Lower number = higher priority
	Timeout    int       // seconds
	MaxResults int       // maximum results per query
	RateLimit  RateLimit // rate limiting settings
}

// RateLimit defines rate limiting parameters
type RateLimit struct {
	RequestsPerMinute int
	RequestsPerDay    int
	BurstSize         int
}

// ProviderStatus tracks the status and statistics of a provider
type ProviderStatus struct {
	Name               string    `json:"name"`
	Available          bool      `json:"available"`
	LastChecked        time.Time `json:"last_checked"`
	LastError          string    `json:"last_error,omitempty"`
	TotalRequests      int64     `json:"total_requests"`
	SuccessfulRequests int64     `json:"successful_requests"`
	FailedRequests     int64     `json:"failed_requests"`
	AvgResponseTime    float64   `json:"avg_response_time_ms"`
	CurrentQuotaUsage  int64     `json:"current_quota_usage,omitempty"`
	QuotaLimit         int64     `json:"quota_limit,omitempty"`
	NextQuotaReset     time.Time `json:"next_quota_reset,omitempty"`
}

// ProviderMetadata contains metadata about available providers
type ProviderMetadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Website     string `json:"website,omitempty"`
	RequiresKey bool   `json:"requires_key"`
	Enabled     bool   `json:"enabled"`
	Priority    int    `json:"priority"`
}

// ProviderRegistry holds all available providers
type ProviderRegistry struct {
	Providers map[string]*ProviderConfig
	Status    map[string]*ProviderStatus
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		Providers: make(map[string]*ProviderConfig),
		Status:    make(map[string]*ProviderStatus),
	}
}

// RegisterProvider adds a provider to the registry
func (pr *ProviderRegistry) RegisterProvider(config *ProviderConfig) {
	pr.Providers[config.Name] = config

	// Initialize status
	pr.Status[config.Name] = &ProviderStatus{
		Name:        config.Name,
		Available:   config.Enabled,
		LastChecked: time.Now(),
	}
}

// GetProvider returns a provider by name
func (pr *ProviderRegistry) GetProvider(name string) (*ProviderConfig, bool) {
	provider, exists := pr.Providers[name]
	return provider, exists
}

// GetStatus returns the status of a provider
func (pr *ProviderRegistry) GetStatus(name string) (*ProviderStatus, bool) {
	status, exists := pr.Status[name]
	return status, exists
}

// UpdateStatus updates provider status
func (pr *ProviderRegistry) UpdateStatus(name string, status *ProviderStatus) {
	status.LastChecked = time.Now()
	pr.Status[name] = status
}

// GetEnabledProviders returns all enabled providers sorted by priority
func (pr *ProviderRegistry) GetEnabledProviders() []*ProviderConfig {
	var enabled []*ProviderConfig
	for _, provider := range pr.Providers {
		if provider.Enabled {
			enabled = append(enabled, provider)
		}
	}

	// Sort by priority (bubble sort for simplicity, can use sort.Slice in production)
	for i := 0; i < len(enabled); i++ {
		for j := i + 1; j < len(enabled); j++ {
			if enabled[j].Priority < enabled[i].Priority {
				enabled[i], enabled[j] = enabled[j], enabled[i]
			}
		}
	}

	return enabled
}

// GetMetadata returns metadata about a provider
func (pr *ProviderRegistry) GetMetadata(name string) *ProviderMetadata {
	provider, exists := pr.GetProvider(name)
	if !exists {
		return nil
	}

	metadata := &ProviderMetadata{
		Name:        provider.Name,
		Enabled:     provider.Enabled,
		Priority:    provider.Priority,
		RequiresKey: provider.APIKey != "",
	}

	// Set provider-specific metadata
	switch name {
	case "brave":
		metadata.DisplayName = "Brave Search"
		metadata.Description = "Fast, privacy-focused search engine"
		metadata.Website = "https://search.brave.com"
		metadata.RequiresKey = true

	case "duckduckgo":
		metadata.DisplayName = "DuckDuckGo"
		metadata.Description = "Privacy-focused search engine"
		metadata.Website = "https://duckduckgo.com"
		metadata.RequiresKey = false

	case "wikipedia":
		metadata.DisplayName = "Wikipedia"
		metadata.Description = "Free encyclopedia"
		metadata.Website = "https://wikipedia.org"
		metadata.RequiresKey = false

	case "personal":
		metadata.DisplayName = "Personal Documents"
		metadata.Description = "Your uploaded documents and files"
		metadata.RequiresKey = false
	}

	return metadata
}

// GetAllMetadata returns metadata for all providers
func (pr *ProviderRegistry) GetAllMetadata() []*ProviderMetadata {
	var metadata []*ProviderMetadata
	for providerName := range pr.Providers {
		if meta := pr.GetMetadata(providerName); meta != nil {
			metadata = append(metadata, meta)
		}
	}
	return metadata
}

// DefaultProviderConfigs returns default configurations for all providers
func DefaultProviderConfigs() map[string]*ProviderConfig {
	return map[string]*ProviderConfig{
		"brave": {
			Name:       "brave",
			Enabled:    true, // Enabled by default
			Priority:   1,    // Highest priority
			Timeout:    10,   // 10 seconds
			MaxResults: 20,
			RateLimit: RateLimit{
				RequestsPerMinute: 60,
				RequestsPerDay:    66,
				BurstSize:         5,
			},
		},
		"duckduckgo": {
			Name:       "duckduckgo",
			Enabled:    false, // Disabled by default
			Priority:   2,
			Timeout:    10,
			MaxResults: 20,
			RateLimit: RateLimit{
				RequestsPerMinute: 100,
				RequestsPerDay:    10000,
				BurstSize:         10,
			},
		},
		"wikipedia": {
			Name:       "wikipedia",
			Enabled:    false, // Disabled by default
			Priority:   3,
			Timeout:    10,
			MaxResults: 10,
			RateLimit: RateLimit{
				RequestsPerMinute: 60,
				RequestsPerDay:    10000,
				BurstSize:         5,
			},
		},
		"personal": {
			Name:       "personal",
			Enabled:    false, // Enabled only if user has documents
			Priority:   0,     // Highest priority when enabled
			Timeout:    5,
			MaxResults: 50,
			RateLimit: RateLimit{
				RequestsPerMinute: 0, // Unlimited for local searches
				RequestsPerDay:    0,
				BurstSize:         0,
			},
		},
	}
}

// HealthCheckResult represents the result of a provider health check
type HealthCheckResult struct {
	Name         string           `json:"name"`
	Healthy      bool             `json:"healthy"`
	ResponseTime int64            `json:"response_time_ms"`
	Error        string           `json:"error,omitempty"`
	CheckedAt    time.Time        `json:"checked_at"`
	Metadata     ProviderMetadata `json:"metadata"`
}
