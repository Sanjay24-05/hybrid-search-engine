// config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Environment      string
	BraveAPIKey      string
	BraveDailyQuota  int
	AllowedOrigins   string
	MaxRequestSizeMB int
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Parse quota
	quota, err := strconv.Atoi(getEnv("BRAVE_DAILY_QUOTA", "66"))
	if err != nil {
		quota = 66
	}

	maxSize, err := strconv.Atoi(getEnv("MAX_REQUEST_SIZE_MB", "10"))
	if err != nil {
		maxSize = 10
	}

	config := &Config{
		Port:             getEnv("PORT", "8080"),
		Environment:      getEnv("ENV", "development"),
		BraveAPIKey:      getEnv("BRAVE_API_KEY", ""),
		BraveDailyQuota:  quota,
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		MaxRequestSizeMB: maxSize,
	}

	// Validate critical configs
	if config.BraveAPIKey == "" {
		fmt.Println("Warning: BRAVE_API_KEY not set, Brave search will be disabled")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
