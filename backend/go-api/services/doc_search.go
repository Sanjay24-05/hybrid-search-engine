package services

import (
	"encoding/json"
	"net/http"
)

type DocResult struct {
	File    string `json:"file"`
	Snippet string `json:"snippet"`
}

func SearchDocuments(userID, query string) ([]DocResult, error) {
	req, _ := http.NewRequest(
		"GET",
		"http://localhost:8000/search-docs?q="+query,
		nil,
	)
	req.Header.Set("user-id", userID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []DocResult
	json.NewDecoder(resp.Body).Decode(&results)

	return results, nil
}
