package cache

import (
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	cache := New()

	// Set a value in the cache
	cache.Set("key1", "value1", 1*time.Minute)

	// Retrieve the value
	value, found := cache.Get("key1")
	if !found {
		t.Fatal("Expected to find key1 in cache, but it was not found")
	}

	// Verify the value
	if value != "value1" {
		t.Errorf("Expected value 'value1', got '%v'", value)
	}

	// Test non-existent key
	_, found = cache.Get("non-existent")
	if found {
		t.Error("Expected non-existent key to not be found, but it was found")
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := New()

	// Set a value with a very short expiration
	cache.Set("key1", "value1", 1*time.Millisecond)

	// Wait for expiration
	time.Sleep(5 * time.Millisecond)

	// Retrieve should not find the value
	_, found := cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired, but it was found")
	}
}

func TestCacheDelete(t *testing.T) {
	cache := New()

	// Set a value
	cache.Set("key1", "value1", 1*time.Minute)

	// Delete the value
	cache.Delete("key1")

	// Retrieve should not find the value
	_, found := cache.Get("key1")
	if found {
		t.Error("Expected key1 to be deleted, but it was found")
	}

	// Delete a non-existent key should not cause issues
	cache.Delete("non-existent")
}