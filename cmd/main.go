package main

import (
	"log"
	"net/http"

	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/cache"
	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/config"
	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/confluent"
	httpHandler "github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/http"
	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/handlers"
	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate required configuration
	if cfg.ConfluentAPIKey == "" || cfg.ConfluentAPISecret == "" {
		log.Fatalf("CONFLUENT_API_KEY and CONFLUENT_API_SECRET must be set")
	}

	// Log configuration (excluding sensitive information)
	log.Printf("Configuration loaded successfully")
	log.Printf("Cache duration set to %v", cfg.CacheDuration)

	// Initialize Confluent API client
	client := confluent.NewClient(cfg.ConfluentAPIKey, cfg.ConfluentAPISecret)

	// Initialize cache
	cacheInstance := cache.New()

	// Create router
	mux := http.NewServeMux()

	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(cfg.ConfluentAPIKey)

	// Register handlers
	mux.Handle("/health", httpHandler.HealthHandler())
	mux.Handle("/discovery", authMiddleware(handlers.DiscoveryHandler(client, cacheInstance, cfg.CacheDuration)))

	// Start the server
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}