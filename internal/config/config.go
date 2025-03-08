package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	ConfluentAPIKey    string
	ConfluentAPISecret string
	CacheDuration      time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	apiKey := os.Getenv("CONFLUENT_API_KEY")
	apiSecret := os.Getenv("CONFLUENT_API_SECRET")

	cacheDurationStr := os.Getenv("CACHE_DURATION")
	cacheDuration := 30 * time.Minute // Default cache duration: 30 minutes

	if cacheDurationStr != "" {
		durationMinutes, err := strconv.Atoi(cacheDurationStr)
		if err != nil {
			log.Printf("Invalid CACHE_DURATION value: %s, using default of 30 minutes", cacheDurationStr)
		} else {
			cacheDuration = time.Duration(durationMinutes) * time.Minute
		}
	}

	return &Config{
		ConfluentAPIKey:    apiKey,
		ConfluentAPISecret: apiSecret,
		CacheDuration:      cacheDuration,
	}, nil
}