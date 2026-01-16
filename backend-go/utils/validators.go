// utils/validators.go
package utils

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrQueryEmpty   = errors.New("query cannot be empty")
	ErrQueryTooLong = errors.New("query too long (max 500 characters)")
	ErrInvalidChars = errors.New("query contains invalid characters")
)

// ValidateSearchQuery validates and sanitizes search queries
func ValidateSearchQuery(query string) (string, error) {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// Check if empty
	if len(query) == 0 {
		return "", ErrQueryEmpty
	}

	// Check length
	if len(query) > 500 {
		return "", ErrQueryTooLong
	}

	// Sanitize - remove potentially dangerous characters
	sanitized := SanitizeInput(query)

	return sanitized, nil
}

// SanitizeInput removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove HTML tags
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")

	// Remove script-related characters
	input = strings.ReplaceAll(input, "<script>", "")
	input = strings.ReplaceAll(input, "</script>", "")

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim spaces again
	input = strings.TrimSpace(input)

	return input
}

// ValidateMaxResults validates max results parameter
func ValidateMaxResults(maxResults int) int {
	if maxResults <= 0 {
		return 10 // default
	}
	if maxResults > 50 {
		return 50 // cap at 50
	}
	return maxResults
}
