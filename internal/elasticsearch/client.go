package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

// Client wraps the Elasticsearch client for our specific use case
type Client struct {
	es    *elasticsearch.Client
	index string
}

// SNMPSettings contains SNMP protocol configuration for the device
type SNMPSettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Version  string `json:"version"`
	AuthName string `json:"auth_name"`
	Timeout  string `json:"timeout"`
	Retries  int    `json:"retries"`
}

// CollectorSettings contains settings for the SNMP metrics collector
type CollectorSettings struct {
	Hostname           string          `json:"hostname"`
	Version            string          `json:"version"`
	Modules            []string        `json:"modules"`
	CollectionInterval string          `json:"collection_interval"`
	Metrics           MetricsSettings `json:"metrics"`
}

// MetricsSettings defines which metrics to collect
type MetricsSettings struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude,omitempty"`
}

// Tags contains metadata tags for the device
type Tags struct {
	Environment string `json:"environment"`
	Location    string `json:"location"`
	Role        string `json:"role"`
}

// Config represents a device configuration document in Elasticsearch
type Config struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Enabled           bool              `json:"enabled"`
	SNMPSettings      SNMPSettings      `json:"snmp_settings"`
	CollectorSettings CollectorSettings `json:"collector_settings"`
	Tags              Tags              `json:"tags"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// MetricsDocument represents a metrics document in Elasticsearch
type MetricsDocument struct {
	Timestamp   time.Time              `json:"@timestamp"`
	DeviceID    string                 `json:"device_id"`
	DeviceName  string                 `json:"device_name"`
	DeviceType  string                 `json:"device_type"`
	Host        string                 `json:"host"`
	Environment string                 `json:"environment"`
	Location    string                 `json:"location"`
	Role        string                 `json:"role"`
	Metrics     map[string]interface{} `json:"metrics"`
}

// NewClient creates a new Elasticsearch client wrapper
func NewClient(esclient *elasticsearch.Client, index string, opts ...func(*Client)) *Client {
	client := &Client{
		es:    esclient,
		index: index,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// ListConfigs retrieves all device configurations from Elasticsearch
func (c *Client) ListConfigs(ctx context.Context) ([]Config, error) {
	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(c.index),
		c.es.Search.WithSize(1000),
	)
	if err != nil {
		return nil, fmt.Errorf("searching configs: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("search response error: %s", body)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source Config `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	configs := make([]Config, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		configs[i] = hit.Source
	}

	return configs, nil
}

// SaveConfig saves a device configuration to Elasticsearch
func (c *Client) SaveConfig(ctx context.Context, config Config) error {
	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = config.UpdatedAt
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(config); err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	res, err := c.es.Index(
		c.index,
		&buf,
		c.es.Index.WithContext(ctx),
		c.es.Index.WithDocumentID(config.ID),
	)
	if err != nil {
		return fmt.Errorf("indexing config: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("index response error: %s", body)
	}

	return nil
}

// DeleteConfig removes a device configuration from Elasticsearch
func (c *Client) DeleteConfig(ctx context.Context, id string) error {
	res, err := c.es.Delete(
		c.index,
		id,
		c.es.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("deleting config: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("delete response error: %s", body)
	}

	return nil
}

// StoreMetrics stores a metrics document in Elasticsearch
func (c *Client) StoreMetrics(ctx context.Context, doc MetricsDocument) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("encoding metrics document: %w", err)
	}

	res, err := c.es.Index(
		c.index,
		&buf,
		c.es.Index.WithContext(ctx),
		c.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("indexing metrics: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("index response error: %s", body)
	}

	return nil
}
