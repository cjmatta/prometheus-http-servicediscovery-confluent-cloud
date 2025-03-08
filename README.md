**Note:** This was an experiment with Claude Code, I don't know go, I don't write go, I just prompted.

# Prometheus HTTP Service Discovery for Confluent Cloud

This service provides an HTTP-based service discovery implementation for Prometheus to scrape metrics from Confluent Cloud.

## Overview

The service integrates with Confluent Cloud APIs to fetch resource information and presents it in a format compatible with Prometheus HTTP-based service discovery. It includes in-memory caching with a configurable TTL to minimize API calls to Confluent Cloud.

## Features

- HTTP-based service discovery for Prometheus
- In-memory caching (default: 30 minutes)
- Authentication via Bearer token
- Fetches and organizes metadata for different resource types:
  - Kafka clusters
  - Schema Registry instances
  - KSQL databases
  - Compute pools
  - Connectors
- Health endpoint for Kubernetes liveness/readiness probes

## Resource Types and Metadata

The service collects the following resource types and metadata:

| Resource Type | Parameter | Metadata |
|---------------|-----------|----------|
| Kafka Clusters | `resource.kafka.id` | cloud_provider, environment_name, cluster_name |
| Schema Registry | `resource.schema_registry.id` | cloud_provider, environment_name |
| KSQL | `resource.ksql.id` | cloud_provider, environment_name, name |
| Compute Pools | `resource.compute_pool.id` | cloud_provider, environment_name |
| Connectors | `resource.connector.id` | cloud_provider, environment_name, connector_name, cluster_id |

## Endpoints

### `/discovery`

- Method: `GET`
- Authentication: Bearer token (Confluent API key)
- Query Parameters:
  - Required: `targets` (comma-separated list)
  - Optional: `prefix` (label prefix)
- Response: JSON conforming to [Prometheus HTTP service discovery](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config) format

### `/health`

- Method: `GET`
- Response: HTTP 200 (OK)

## Configuration

The following environment variables are used for configuration:

- `CONFLUENT_API_KEY`: Confluent Cloud API key
- `CONFLUENT_API_SECRET`: Confluent Cloud API secret
- `CACHE_DURATION`: Cache duration in minutes (default: 30)

## Deployment

### Docker

Build and run the Docker image:

```shell
docker build -t prometheus-http-servicediscovery-confluent-cloud .
docker run -p 8080:8080 \
  -e CONFLUENT_API_KEY=your_api_key \
  -e CONFLUENT_API_SECRET=your_api_secret \
  prometheus-http-servicediscovery-confluent-cloud
```

### Kubernetes

1. Update the API key and secret in `kubernetes/deployment.yaml`
2. Apply the Kubernetes manifests:

```shell
kubectl apply -f kubernetes/deployment.yaml
```

## Prometheus Configuration

Add the following configuration to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'confluent-cloud'
    http_sd_configs:
      - url: 'http://prometheus-http-servicediscovery-confluent-cloud:8080/discovery?targets=metrics.confluent.cloud:443&prefix=confluent_'
        refresh_interval: 30m
        authorization:
          type: Bearer
          credentials: 'your_api_key'
```

## Response Format Example

```json
[
  {
    "targets": ["metrics.confluent.cloud:443"],
    "labels": {
      "confluent_cloud_provider": "aws",
      "confluent_environment_name": "prod",
      "confluent_cluster_name": "main-cluster"
    },
    "params": {
      "resource.kafka.id": ["lkc-abc123"]
    }
  },
  {
    "targets": ["metrics.confluent.cloud:443"],
    "labels": {
      "confluent_cloud_provider": "gcp",
      "confluent_environment_name": "staging",
      "confluent_name": "ksql-app"
    },
    "params": {
      "resource.ksql.id": ["lksqlc-def456"]
    }
  }
]
```

## Development

Build and run the service locally:

```shell
export CONFLUENT_API_KEY=your_api_key
export CONFLUENT_API_SECRET=your_api_secret
go run cmd/main.go
```
