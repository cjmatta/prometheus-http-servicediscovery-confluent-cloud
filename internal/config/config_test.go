package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("CONFLUENT_API_KEY", "test-key")
	os.Setenv("CONFLUENT_API_SECRET", "test-secret")
	os.Setenv("CACHE_DURATION", "15")

	// Load configuration
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Verify configuration values
	if cfg.ConfluentAPIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", cfg.ConfluentAPIKey)
	}

	if cfg.ConfluentAPISecret != "test-secret" {
		t.Errorf("Expected API secret 'test-secret', got '%s'", cfg.ConfluentAPISecret)
	}

	expectedDuration := 15 * time.Minute
	if cfg.CacheDuration != expectedDuration {
		t.Errorf("Expected cache duration %v, got %v", expectedDuration, cfg.CacheDuration)
	}

	// Test with invalid cache duration
	os.Setenv("CACHE_DURATION", "invalid")
	cfg, err = Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Should default to 30 minutes
	expectedDuration = 30 * time.Minute
	if cfg.CacheDuration != expectedDuration {
		t.Errorf("Expected default cache duration %v, got %v", expectedDuration, cfg.CacheDuration)
	}

	// Clean up
	os.Unsetenv("CONFLUENT_API_KEY")
	os.Unsetenv("CONFLUENT_API_SECRET")
	os.Unsetenv("CACHE_DURATION")
}

func TestLoadDefaultCacheDuration(t *testing.T) {
	// Set environment variables without cache duration
	os.Setenv("CONFLUENT_API_KEY", "test-key")
	os.Setenv("CONFLUENT_API_SECRET", "test-secret")
	os.Unsetenv("CACHE_DURATION")

	// Load configuration
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Verify default cache duration
	expectedDuration := 30 * time.Minute
	if cfg.CacheDuration != expectedDuration {
		t.Errorf("Expected default cache duration %v, got %v", expectedDuration, cfg.CacheDuration)
	}

	// Clean up
	os.Unsetenv("CONFLUENT_API_KEY")
	os.Unsetenv("CONFLUENT_API_SECRET")
}