package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/cache"
	"github.com/cjmatta/prometheus-http-servicediscovery-confluent-cloud/internal/confluent"
)

const (
	cacheKey = "confluent_resources"
)

var (
	// validPrefixPattern is used to validate the prefix parameter
	validPrefixPattern = regexp.MustCompile(`^[a-zA-Z0-9_]*$`)
)

// Target represents a target for Prometheus to scrape
type Target struct {
	Targets []string            `json:"targets"`
	Labels  map[string]string   `json:"labels"`
	Params  map[string][]string `json:"params"`
}

// DiscoveryHandler handles the /discovery endpoint
func DiscoveryHandler(client *confluent.Client, cache *cache.Cache, cacheDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if we have cached data first, before potentially making API calls
		cachedData, found := cache.Get(cacheKey)
		var resourcesNeedFetching = !found
		
		// Parse query parameters
		targetsParam := r.URL.Query().Get("targets")
		if targetsParam == "" {
			http.Error(w, "Missing required 'targets' parameter", http.StatusBadRequest)
			return
		}

		// Split targets by comma
		targetsList := strings.Split(targetsParam, ",")

		// Get optional prefix
		prefix := r.URL.Query().Get("prefix")
		if prefix != "" && !validPrefixPattern.MatchString(prefix) {
			http.Error(w, "Invalid 'prefix' parameter. Must contain only alphanumeric characters and underscores", http.StatusBadRequest)
			return
		}

		// Add trailing underscore to prefix if it's not empty
		if prefix != "" && !strings.HasSuffix(prefix, "_") {
			prefix = prefix + "_"
		}

		// After validating parameters, fetch data if needed
		var resources []confluent.Resource

		if resourcesNeedFetching {
			// Fetch data from Confluent API since parameters are valid
			log.Println("Cache miss. Fetching data from Confluent API...")
			
			var err error
			resources, err = client.GetAllResources()
			if err != nil {
				log.Printf("Failed to fetch resources: %v", err)
				http.Error(w, "Failed to fetch resources from Confluent API", http.StatusInternalServerError)
				return
			}

			// Cache the results
			cache.Set(cacheKey, resources, cacheDuration)
		} else {
			// Use cached data
			log.Println("Using cached data")
			resources = cachedData.([]confluent.Resource)
		}

		// Format response for Prometheus
		response := formatResponse(resources, targetsList, prefix)

		// Set content type and return JSON response
		w.Header().Set("Content-Type", "application/json")
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		
		log.Printf("Returned %d resources to Prometheus", len(response))
	}
}

// formatResponse formats the response for Prometheus
func formatResponse(resources []confluent.Resource, targets []string, prefix string) []Target {
	var response []Target

	for _, resource := range resources {
		// Create a new target with the requested targets
		target := Target{
			Targets: targets,
			Labels:  make(map[string]string),
			Params:  make(map[string][]string),
		}

		// Add labels with optional prefix
		for k, v := range resource.Labels {
			target.Labels[prefix+k] = v
		}

		// Add resource ID to params based on resource type
		switch resource.ResourceType {
		case "kafka":
			target.Params["resource.kafka.id"] = []string{resource.ID}
		case "schema_registry":
			target.Params["resource.schema_registry.id"] = []string{resource.ID}
		case "ksql":
			target.Params["resource.ksql.id"] = []string{resource.ID}
		case "compute_pool":
			target.Params["resource.compute_pool.id"] = []string{resource.ID}
		case "connector":
			target.Params["resource.connector.id"] = []string{resource.ID}
		}

		response = append(response, target)
	}

	return response
}