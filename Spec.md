# Confluent Cloud Prometheus Service Discovery Endpoint Specification

## Overview
An HTTP endpoint service implemented in Go to integrate Confluent Cloud with Prometheus HTTP-based service discovery.

## Requirements
- **Language:** Golang
- **Endpoints:** `/discovery`, `/health`
- **Authentication:** Bearer Token (API key via Kubernetes ENV)
- **Caching:** In-memory, 30-minute duration
- **Response format:** JSON, conforming to Prometheus HTTP service discovery specification

## Architecture
- Service will run in Kubernetes
- Fetches data from Confluent Cloud API
- Caches data in-memory (30-minute TTL)

## Configuration
- ENV variables:
  - `CONFLUENT_API_KEY`
  - `CONFLUENT_API_SECRET`
  - `CACHE_DURATION` (default: 30 mins)

## Endpoint Details
### `/discovery`
- Method: `GET`
- Query Params:
  - Required: `targets` (comma-separated list)
  - Optional: `prefix` (label prefix)
- Response format:
```json
[
  {
    "targets": ["host1", "host2"],
    "labels": {
      "<prefix>cloud_provider": "aws",
      "<prefix>environment_name": "prod",
      "<prefix>cluster_name": "main-cluster"
    },
    "params": {
      "resource.kafka.id": ["kafka_id1", "kafka_id2"]
    }
  }
]
```

### `/health`
- Method: `GET`
- Response: HTTP 200

## Data Fetching Logic
1. Call Confluent API:
   - `/org/v2/environments`
   - `/v2/metrics/cloud/descriptors/resources`
2. Cache results for 30 minutes.
3. Return cached results or HTTP 500 on error.

## Metadata Included
- **resource.kafka.id:** `cloud_provider`, `environment_name`, `cluster_name`
- **resource.schema_registry.id:** `cloud_provider`, `environment_name`
- **resource.ksql.id:** `cloud_provider`, `environment_name`, `name`
- **resource.compute_pool.id:** `cloud_provider`, `environment_name`
- **resource.connector.id:** `cloud_provider`, `environment_name`, `connector_name`, `cluster_id`

## `/health` Endpoint
- Respond with HTTP 200.
- Used for Kubernetes readiness/liveness probes.

## Error Handling
- Authentication failure: `401 Unauthorized`
- Confluent API fetch failure: `500 Internal Server Error` (with detailed messages)
- Invalid `prefix` parameter: `400 Bad Request`

## Logging
- Basic logging: HTTP status codes, request times

## Detailed Project Implementation Plan

### Iteration 1: Setup and Configuration
- Initialize Golang project structure.
- Implement environment variable configuration (`cmd/main.go`, `internal/config/config.go`).
- Initialize logging.

### Iteration 2: Authentication Middleware
- Bearer token validation and middleware implementation (`internal/middleware/auth.go`).
- Integrate into HTTP server.

### Iteration 3: Health Endpoint
- Simple `/health` endpoint (`internal/http/health.go`).

### Iteration 4: Confluent API Client
- API client implementation (`internal/confluent/client.go`).
- Fetch environments and resources.

### Iteration 5: Caching Layer
- Implement in-memory caching (`internal/cache/cache.go`).
- Verify caching behavior and expiration.

### Iteration 6: Discovery Endpoint
- Implement `/discovery` handler (`internal/handlers/discovery.go`).
- Handle query parameters and basic JSON formatting.

### Iteration 7: Label Prefix Handling
- Implement validation and application of optional prefix parameter.

### Iteration 8: Comprehensive Error Handling
- Detailed error responses and logging for clarity.

### Iteration 9: Testing and Validation
- Unit tests (configuration, middleware, caching).
- Integration tests (end-to-end, response validation).
- Load testing for caching efficiency under concurrent requests.

### Iteration 10: Deployment
- Prepare Kubernetes deployment YAML.
- Implement readiness/liveness probes (`/health`).

This detailed implementation plan ensures incremental, manageable, and safe development progress.

