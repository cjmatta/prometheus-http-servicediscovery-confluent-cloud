# Project Information

This project is a Prometheus HTTP service discovery mechanism for Confluent Cloud, written in golang.

## Key Documents
- [Project Specification](Spec.md) - Defines requirements, endpoints, and data structures
- [Implementation Plan](Implementation.md) - Provides iterative implementation steps

## Project Structure Overview
- **Language:** Go
- **Main Endpoints:** 
  - `/discovery` - Returns Prometheus-compatible service discovery data
  - `/health` - Health check endpoint for Kubernetes
- **Authentication:** Bearer Token (API key via Kubernetes ENV)
- **Caching:** In-memory, 30-minute duration

## Implementation Plan
1. Initialize Project & Setup
2. Implement Authentication and Configuration
3. Implement API clients for Confluent Cloud endpoints
4. Implement Caching Layer
5. Create HTTP Endpoint Handlers
6. Implement Label Prefix Handling and Validation
7. Format the Response for Prometheus
8. Add Comprehensive Error Handling
9. Testing & Integration
10. Deployment and Integration

## Commands
<!-- Add common build, test, lint commands here as you use them -->

## Code Style Preferences
<!-- Add code style preferences here as you establish them -->

## Important Codebase Information
### Confluent Cloud API Endpoints
- `/org/v2/environments` - Fetches environment information
- `/cmk/v2/clusters?environment={env_id}` - Fetches Kafka clusters for an environment
- `/srcm/v2/clusters?environment={env_id}` - Fetches Schema Registry instances for an environment
- `/ksqldbcm/v2/clusters?environment={env_id}` - Fetches KSQL databases for an environment
- `/fcpm/v2/compute-pools?environment={env_id}` - Fetches compute pools for an environment
- `/connect/v1/environments/{env_id}/clusters/{cluster_id}/connectors` - Fetches connectors

### API Dependencies
- The service first fetches all environments
- Then for each environment, it fetches resources (clusters, schema registries, etc.)
- For connectors, it first needs both environment ID and cluster ID

### Resource Types and Metadata
| Resource Type | Parameter | Metadata |
|---------------|-----------|----------|
| Kafka Clusters | `resource.kafka.id` | cloud_provider, environment_name, cluster_name |
| Schema Registry | `resource.schema_registry.id` | cloud_provider, environment_name |
| KSQL | `resource.ksql.id` | cloud_provider, environment_name, name |
| Compute Pools | `resource.compute_pool.id` | cloud_provider, environment_name |
| Connectors | `resource.connector.id` | cloud_provider, environment_name, connector_name, cluster_id |

### Configuration Environment Variables
- `CONFLUENT_API_KEY` - API key for Confluent Cloud authentication
- `CONFLUENT_API_SECRET` - API secret for Confluent Cloud authentication  
- `CACHE_DURATION` - Cache duration in minutes (default: 30)

### Response Format Structure
```json
[
  {
    "targets": ["metrics.confluent.cloud:443"],
    "labels": {
      "<prefix>cloud_provider": "aws",
      "<prefix>environment_name": "prod",
      "<prefix>cluster_name": "main-cluster"
    },
    "params": {
      "resource.kafka.id": ["lkc-abc123"]
    }
  }
]
```

### Authentication
- This service uses the API key as the Bearer token for client authentication
- The same API key and secret are used to authenticate with Confluent Cloud APIs
- In your `/discovery` endpoint requests, use: `Authorization: Bearer YOUR_API_KEY`