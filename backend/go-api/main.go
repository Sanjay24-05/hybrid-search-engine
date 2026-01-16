package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"hybrid-search-engine/go-api/middleware"
	"hybrid-search-engine/go-api/routes"
)

func main() {
	// Resolve absolute path to project root .env
	exePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envPath := filepath.Join(exePath, "..", "..", ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", routes.Login)
	mux.Handle(
		"/search",
		middleware.Auth(http.HandlerFunc(routes.Search)),
	)
	mux.Handle(
		"/upload",
		middleware.Auth(http.HandlerFunc(routes.Upload)),
	)

	handler := middleware.RateLimit(
		middleware.SecurityHeaders(
			middleware.CORS(mux),
		),
	)

	log.Printf("🚀 Go API running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
