package models

type SearchResult struct {
	Source  string `json:"source"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
	Score   float64
}
