# Detailed Project Implementation Plan and Step Breakdown

## High-Level Plan

Here’s a structured, step-by-step overview of the implementation plan for building a Golang-based HTTP Service Discovery endpoint integrating Confluent Cloud with Prometheus:

1.  Initialize Project & Setup
2.  Implement Authentication and Configuration
3.  Implement API clients for Confluent Cloud endpoints
4.  Implement Caching Layer
5.  Create HTTP Endpoint Handlers (/discovery and /health)
6.  Implement Label Prefix Handling and Validation
7.  Format the Response for Prometheus
8.  Add Comprehensive Error Handling
9.  Testing & Integration
10. Deployment and Integration

---

## Iterative Plan:

### Iteration 1: Setup and Configuration

*   Set up a basic Go project structure.
*   Handle .env configuration for API credentials and cache duration.
*   Verify the initial setup by printing configuration details.

#### Iterative Chunk Breakdown:

*   Load environment variables.
*   Initialize basic logging.

---

### Iterative Prompt (1)

Create a new Golang project structure with the following components:

*   `cmd/main.go` to bootstrap the application.
*   `internal/config/config.go` to load configurations from environment variables (`CONFLUENT_API_KEY`, `CONFLUENT_API_SECRET`, `CACHE_DURATION`).
*   Initialize logging using the standard `log` library.

Ensure the service can start, load these variables, and print them out clearly to validate successful initialization.

---

### Iterative Prompt (2)

Extend the project by implementing authentication handling:

*   Write a middleware function in `internal/middleware/auth.go` to validate the `Authorization` header containing the Bearer token.
*   Integrate middleware into the HTTP server setup in `cmd/main.go`.
*   Return HTTP 401 Unauthorized if token is missing or invalid.

---

### Iterative Prompt (3)

Implement an HTTP handler in `internal/http/health.go` to handle the `/health` endpoint:

*   It should simply return HTTP 200 OK to indicate service health.
*   Integrate this handler into the HTTP server.
*   Test this endpoint to ensure it returns HTTP 200.

---

### Iterative Prompt (4)

Implement Confluent Cloud API client:

*   Use Golang's `net/http` client to create API clients under `internal/confluent/client.go`.
*   Implement methods to fetch environments from `/org/v2/environments`.
*   Write a basic test to verify you can successfully authenticate and fetch environments.

---

### Iterative Prompt (5)

Enhance the API client by adding methods to:

*   Fetch resources from `/v2/metrics/cloud/descriptors/resources` for each environment retrieved.
*   Ensure you handle API pagination if applicable.
*   Log any errors during data fetching clearly.
*   Verify correctness by printing a sample response to console.

---

### Iterative Prompt (6)

Implement an in-memory caching layer:

*   Use a library like `go-cache` for simple TTL caching (`internal/cache/cache.go`).
*   Add methods for setting and retrieving cached data.
*   Cache resource responses for 30 minutes.
*   Test caching functionality separately with a simple Go test function in `internal/cache/cache_test.go`.

---

### Iterative Prompt (7)

Create the `/discovery` HTTP handler:

*   Implement in `internal/handlers/discovery.go`.
*   Fetch targets from query parameter (`targets`), validate they are comma-separated.
*   If cache is empty or expired, fetch data from Confluent API and populate cache.
*   If Confluent API call fails, respond with HTTP 500 error and descriptive message.
*   Format basic JSON output without labels initially; verify by responding with basic JSON containing just "targets".

---

### Iterative Prompt (8)

Extend discovery handler to format the JSON output as required:

*   For each resource returned from Confluent Cloud, generate a separate JSON entry.
*   Populate labels with metadata according to the resource type mapping:
    *   `resource.kafka.id`: cloud\_provider, environment\_name, cluster\_name
    *   `resource.schema_registry.id`: cloud\_provider, environment\_name
    *   `resource.ksql.id`: cloud\_provider, environment\_name, name
    *   `resource.compute_pool.id`: cloud\_provider, environment\_name
    *   `resource.connector.id`: cloud\_provider, environment\_name, connector\_name, cluster\_id
*   Populate `params` with the resource IDs, formatted correctly as described.
*   Verify correctness by manually invoking endpoint and inspecting response structure.

---

### Iterative Prompt (9)

Add optional query parameter handling for label prefixes in `/discovery`:

*   Parse optional `prefix` query parameter.
*   Validate that the prefix contains only alphanumeric characters and underscores. If invalid, respond with HTTP 400 Bad Request.
*   Prepend this prefix to every label key in the response if provided.
*   Verify this functionality with example requests.

---

### Iterative Prompt (10)

Ensure consistent JSON response formatting:

*   Verify the endpoint returns `application/json` content-type header explicitly.
*   Ensure HTTP response status codes (200, 400, 401, 500) are correct and clearly documented.
*   Test end-to-end using `curl` or Postman.

---

### Iterative Prompt (11)

Implement `/health` endpoint handling readiness and liveness probes:

*   Confirm the `/health` endpoint returns HTTP 200 without performing additional checks.
*   Integrate this endpoint into Kubernetes manifests (e.g., deployment YAML) as readiness and liveness probe.

---

## 9. Comprehensive Test Plan

### Unit Tests

*   Test configuration and environment variable loading.
*   Test middleware authentication handling.
*   Test caching functionality and expiry.
*   Test response formatting and label prefix handling.

### Integration Tests

*   End-to-end HTTP call simulation, validating JSON structure.
*   Failure scenarios (invalid authentication, Confluent API errors).

### Load Testing

*   Test caching performance under concurrent access.

---

## Final Review Checklist

*   ✅ Clear Golang project structure.
*   ✅ Configurable via environment variables.
*   ✅ Proper authentication handling.
*   ✅ Robust HTTP endpoint implementations (/discovery, /health).
*   ✅ API clients for Confluent Cloud integration.
*   ✅ Efficient in-memory caching (30 minutes TTL).
*   ✅ JSON formatted correctly for Prometheus.
*   ✅ Comprehensive error handling.
*   ✅ Optional prefix parameter with validation.
*   ✅ Basic logging for troubleshooting.
*   ✅ Kubernetes-friendly (/health probe).
*   ✅ Defined testing and validation strategy.

---

This provides a clear, actionable, iterative approach for safe implementation by your developer. Does this look ready to hand off, or would you like any final adjustments?