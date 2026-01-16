package main

import (
	"net/http"

	"hybrid-search-engine/go-api/middleware"
	"hybrid-search-engine/go-api/routes"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", routes.Login)
	mux.Handle("/search",
		middleware.Auth(
			http.HandlerFunc(routes.Search),
		),
	)

	handler := middleware.RateLimit(
		middleware.SecurityHeaders(
			middleware.CORS(mux),
		),
	)

	http.ListenAndServe(":8080", handler)
}
