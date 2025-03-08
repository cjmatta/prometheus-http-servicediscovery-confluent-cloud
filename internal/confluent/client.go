package confluent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL                   = "https://api.confluent.cloud"
	environmentsPath          = "/org/v2/environments"
	kafkaClustersPath         = "/cmk/v2/clusters?environment=%s"
	schemaRegistryPath        = "/srcm/v2/clusters?environment=%s"
	ksqlPath                  = "/ksqldbcm/v2/clusters?environment=%s"
	computePoolsPath          = "/fcpm/v2/compute-pools?environment=%s"
	connectorsBasePath        = "/connect/v1/environments/%s/clusters/%s/connectors"
	defaultTimeout            = 30 * time.Second
)

// Client represents a Confluent Cloud API client
type Client struct {
	httpClient  *http.Client
	apiKey      string
	apiSecret   string
}

// Environment represents a Confluent Cloud environment
type Environment struct {
	ID   string `json:"id"`
	Name string `json:"display_name"`
}

// EnvironmentsResponse represents the response from the environments API
type EnvironmentsResponse struct {
	Data []Environment `json:"data"`
}

// Resource represents a Confluent Cloud resource with metadata
type Resource struct {
	ID           string            `json:"id"`
	ResourceType string            `json:"resource_type"`
	Labels       map[string]string `json:"labels"`
}

// KafkaClusterSpec represents the specification of a Kafka cluster
type KafkaClusterSpec struct {
	DisplayName  string `json:"display_name"`
	Availability string `json:"availability"`
	Cloud        string `json:"cloud"`
	Region       string `json:"region"`
}

// KafkaCluster represents a Kafka cluster
type KafkaCluster struct {
	ID         string          `json:"id"`
	Spec       KafkaClusterSpec `json:"spec"`
	Environment struct {
		ID string `json:"id"`
	} `json:"environment"`
}

// KafkaClustersResponse represents the response from the Kafka clusters API
type KafkaClustersResponse struct {
	Data []KafkaCluster `json:"data"`
}

// SchemaRegistrySpec represents the specification of a Schema Registry instance
type SchemaRegistrySpec struct {
	DisplayName string `json:"display_name"`
	Cloud       string `json:"cloud"`
	Region      string `json:"region"`
}

// SchemaRegistry represents a Schema Registry instance
type SchemaRegistry struct {
	ID         string            `json:"id"`
	Spec       SchemaRegistrySpec `json:"spec"`
	Environment struct {
		ID string `json:"id"`
	} `json:"environment"`
}

// SchemaRegistryResponse represents the response from the Schema Registry API
type SchemaRegistryResponse struct {
	Data []SchemaRegistry `json:"data"`
}

// KsqlDBSpec represents the specification of a KSQL database
type KsqlDBSpec struct {
	DisplayName string `json:"display_name"`
	Cloud       string `json:"cloud"`
	Region      string `json:"region"`
}

// KsqlDB represents a KSQL database
type KsqlDB struct {
	ID         string     `json:"id"`
	Spec       KsqlDBSpec `json:"spec"`
	Environment struct {
		ID string `json:"id"`
	} `json:"environment"`
}

// KsqlDBResponse represents the response from the KSQL API
type KsqlDBResponse struct {
	Data []KsqlDB `json:"data"`
}

// ComputePoolSpec represents the specification of a compute pool
type ComputePoolSpec struct {
	DisplayName string `json:"display_name"`
	Cloud       string `json:"cloud"`
	Region      string `json:"region"`
}

// ComputePool represents a compute pool
type ComputePool struct {
	ID         string        `json:"id"`
	Spec       ComputePoolSpec `json:"spec"`
	Environment struct {
		ID string `json:"id"`
	} `json:"environment"`
}

// ComputePoolsResponse represents the response from the compute pools API
type ComputePoolsResponse struct {
	Data []ComputePool `json:"data"`
}

// Connector represents a connector
type Connector struct {
	ID          string `json:"name"`
	ClusterID   string `json:"cluster_id"`
	Environment string `json:"environment_id"`
}

// NewClient creates a new Confluent Cloud API client
func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// GetEnvironments retrieves environments from Confluent Cloud
func (c *Client) GetEnvironments() ([]Environment, error) {
	log.Println("Fetching environments from Confluent Cloud API")
	req, err := http.NewRequest(http.MethodGet, baseURL+environmentsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// For debugging - print response body
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Environments API response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var envResp EnvironmentsResponse
	if err := json.Unmarshal(body, &envResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Found %d environments", len(envResp.Data))
	return envResp.Data, nil
}

// GetKafkaClusters retrieves Kafka clusters for a specific environment
func (c *Client) GetKafkaClusters(environmentID string) ([]KafkaCluster, error) {
	log.Printf("Fetching Kafka clusters for environment %s", environmentID)
	path := fmt.Sprintf(kafkaClustersPath, environmentID)
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// For debugging - print response body
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Kafka Clusters API response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var clustersResp KafkaClustersResponse
	if err := json.Unmarshal(body, &clustersResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Found %d Kafka clusters for environment %s", len(clustersResp.Data), environmentID)
	return clustersResp.Data, nil
}

// GetSchemaRegistries retrieves Schema Registry instances for a specific environment
func (c *Client) GetSchemaRegistries(environmentID string) ([]SchemaRegistry, error) {
	log.Printf("Fetching Schema Registry instances for environment %s", environmentID)
	path := fmt.Sprintf(schemaRegistryPath, environmentID)
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// For debugging - print response body
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Schema Registry API response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var srResp SchemaRegistryResponse
	if err := json.Unmarshal(body, &srResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Found %d Schema Registry instances for environment %s", len(srResp.Data), environmentID)
	return srResp.Data, nil
}

// GetKsqlDBs retrieves KSQL databases for a specific environment
func (c *Client) GetKsqlDBs(environmentID string) ([]KsqlDB, error) {
	log.Printf("Fetching KSQL databases for environment %s", environmentID)
	path := fmt.Sprintf(ksqlPath, environmentID)
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// For debugging - print response body
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("KSQL API response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var ksqlResp KsqlDBResponse
	if err := json.Unmarshal(body, &ksqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Found %d KSQL databases for environment %s", len(ksqlResp.Data), environmentID)
	return ksqlResp.Data, nil
}

// GetComputePools retrieves compute pools for a specific environment
func (c *Client) GetComputePools(environmentID string) ([]ComputePool, error) {
	log.Printf("Fetching compute pools for environment %s", environmentID)
	path := fmt.Sprintf(computePoolsPath, environmentID)
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// For debugging - print response body
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Compute Pools API response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var poolsResp ComputePoolsResponse
	if err := json.Unmarshal(body, &poolsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Found %d compute pools for environment %s", len(poolsResp.Data), environmentID)
	return poolsResp.Data, nil
}

// GetConnectors retrieves connectors for a specific environment and cluster
func (c *Client) GetConnectors(environmentID, clusterID string) ([]Connector, error) {
	log.Printf("Fetching connectors for environment %s, cluster %s", environmentID, clusterID)
	path := fmt.Sprintf(connectorsBasePath, environmentID, clusterID)
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// For debugging - print response body
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Connectors API response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var connectorNames []string
	if err := json.Unmarshal(body, &connectorNames); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert connector names to connector objects
	connectors := make([]Connector, len(connectorNames))
	for i, name := range connectorNames {
		connectors[i] = Connector{
			ID:          name,
			ClusterID:   clusterID,
			Environment: environmentID,
		}
	}

	log.Printf("Found %d connectors for environment %s, cluster %s", len(connectors), environmentID, clusterID)
	return connectors, nil
}

// GetAllResources fetches all resources and formats them with consistent metadata
func (c *Client) GetAllResources() ([]Resource, error) {
	var resources []Resource

	// Fetch environments
	environments, err := c.GetEnvironments()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch environments: %w", err)
	}

	// Create a map of environment IDs to names for easier lookup
	envMap := make(map[string]string)
	for _, env := range environments {
		envMap[env.ID] = env.Name
	}

	// Process each environment separately
	for _, env := range environments {
		log.Printf("Processing environment: %s (%s)", env.Name, env.ID)
		
		// Fetch Kafka clusters for this environment
		kafkaClusters, err := c.GetKafkaClusters(env.ID)
		if err != nil {
			log.Printf("Warning: failed to fetch Kafka clusters for environment %s: %v", env.ID, err)
		} else {
			for _, cluster := range kafkaClusters {
				// Map cloud provider from cloud field
				cloudProvider := cluster.Spec.Cloud
				if cloudProvider == "" {
					cloudProvider = "unknown"
				}
				
				resources = append(resources, Resource{
					ID:           cluster.ID,
					ResourceType: "kafka",
					Labels: map[string]string{
						"cloud_provider":   cloudProvider,
						"environment_name": env.Name,
						"cluster_name":     cluster.Spec.DisplayName,
						"region":           cluster.Spec.Region,
					},
				})
				
				// Fetch connectors for this Kafka cluster
				connectors, err := c.GetConnectors(env.ID, cluster.ID)
				if err != nil {
					log.Printf("Warning: failed to fetch connectors for environment %s, cluster %s: %v", 
						env.ID, cluster.ID, err)
				} else {
					for _, connector := range connectors {
						resources = append(resources, Resource{
							ID:           connector.ID,
							ResourceType: "connector",
							Labels: map[string]string{
								"cloud_provider":   cloudProvider, // Use cluster's provider
								"environment_name": env.Name,
								"connector_name":   connector.ID,
								"cluster_id":       connector.ClusterID,
								"region":           cluster.Spec.Region,
							},
						})
					}
				}
			}
		}

		// Fetch Schema Registry instances for this environment
		schemaRegistries, err := c.GetSchemaRegistries(env.ID)
		if err != nil {
			log.Printf("Warning: failed to fetch Schema Registry instances for environment %s: %v", env.ID, err)
		} else {
			for _, sr := range schemaRegistries {
				// Map cloud provider from cloud field
				cloudProvider := sr.Spec.Cloud
				if cloudProvider == "" {
					cloudProvider = "unknown"
				}
				
				resources = append(resources, Resource{
					ID:           sr.ID,
					ResourceType: "schema_registry",
					Labels: map[string]string{
						"cloud_provider":   cloudProvider,
						"environment_name": env.Name,
						"name":             sr.Spec.DisplayName,
						"region":           sr.Spec.Region,
					},
				})
			}
		}

		// Fetch KSQL databases for this environment
		ksqlDBs, err := c.GetKsqlDBs(env.ID)
		if err != nil {
			log.Printf("Warning: failed to fetch KSQL databases for environment %s: %v", env.ID, err)
		} else {
			for _, ksql := range ksqlDBs {
				// Map cloud provider from cloud field
				cloudProvider := ksql.Spec.Cloud
				if cloudProvider == "" {
					cloudProvider = "unknown"
				}
				
				resources = append(resources, Resource{
					ID:           ksql.ID,
					ResourceType: "ksql",
					Labels: map[string]string{
						"cloud_provider":   cloudProvider,
						"environment_name": env.Name,
						"name":             ksql.Spec.DisplayName,
						"region":           ksql.Spec.Region,
					},
				})
			}
		}

		// Fetch compute pools for this environment
		computePools, err := c.GetComputePools(env.ID)
		if err != nil {
			log.Printf("Warning: failed to fetch compute pools for environment %s: %v", env.ID, err)
		} else {
			for _, pool := range computePools {
				// Map cloud provider from cloud field
				cloudProvider := pool.Spec.Cloud
				if cloudProvider == "" {
					cloudProvider = "unknown"
				}
				
				resources = append(resources, Resource{
					ID:           pool.ID,
					ResourceType: "compute_pool",
					Labels: map[string]string{
						"cloud_provider":   cloudProvider,
						"environment_name": env.Name,
						"name":             pool.Spec.DisplayName,
						"region":           pool.Spec.Region,
					},
				})
			}
		}
	}
	
	log.Printf("Found %d total resources across %d environments", len(resources), len(environments))
	return resources, nil
}