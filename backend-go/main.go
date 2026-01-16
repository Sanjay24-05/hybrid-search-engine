// main.go
package main

import (
	"fmt"
	"log"

	"github.com/Sanjay24-05/hybrid-search-engine/config"
	"github.com/Sanjay24-05/hybrid-search-engine/handlers"
	"github.com/Sanjay24-05/hybrid-search-engine/middleware"
	"github.com/Sanjay24-05/hybrid-search-engine/services/search"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Apply middleware
	router.Use(gin.Recovery()) // Recover from panics
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestSizeLimitMiddleware(cfg.MaxRequestSizeMB))

	// Initialize search aggregator
	aggregator := setupSearchAggregator(cfg)

	// Initialize handlers
	searchHandler := handlers.NewSearchHandler(aggregator)

	// Routes
	api := router.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.POST("/search", searchHandler.HandleSearch)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("CORS allowed origins: %s", cfg.AllowedOrigins)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupSearchAggregator(cfg *config.Config) *search.Aggregator {
	aggregator := search.NewAggregator()

	// Register Brave provider (if API key available)
	if cfg.BraveAPIKey != "" {
		brave := search.NewBraveProvider(cfg.BraveAPIKey, cfg.BraveDailyQuota)
		aggregator.RegisterProvider(brave)
		log.Println("✓ Brave Search enabled")
	} else {
		log.Println("⚠ Brave Search disabled (no API key)")
	}

	// Register DuckDuckGo provider
	ddg := search.NewDuckDuckGoProvider()
	aggregator.RegisterProvider(ddg)
	log.Println("✓ DuckDuckGo Search enabled")

	// Register Wikipedia provider
	wiki := search.NewWikipediaProvider()
	aggregator.RegisterProvider(wiki)
	log.Println("✓ Wikipedia Search enabled")

	return aggregator
}
