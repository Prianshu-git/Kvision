package config

import (
	"log"
	"os"
	"time"
)

// Config holds all configuration values for the agent
type Config struct {
	PromURL        string
	BackendURL     string
	ScrapeInterval time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	promURL := getEnv("PROM_URL", "")
	backendURL := getEnv("BACKEND_URL", "")
	scrapeIntervalStr := getEnv("SCRAPE_INTERVAL", "30s")

	if promURL == "" {
		log.Fatal("PROM_URL is required but missing")
	}

	if backendURL == "" {
		log.Fatal("BACKEND_URL is required but missing")
	}

	interval, err := time.ParseDuration(scrapeIntervalStr)
	if err != nil {
		log.Fatalf("Invalid SCRAPE_INTERVAL: %v", err)
	}

	return &Config{
		PromURL:        promURL,
		BackendURL:     backendURL,
		ScrapeInterval: interval,
	}
}

// helper to read environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
