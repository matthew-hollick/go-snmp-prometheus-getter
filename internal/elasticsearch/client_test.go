package elasticsearch

import (
	"context"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestClient(t *testing.T) *Client {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	client, err := NewClient(cfg, "service_configuration")
	require.NoError(t, err)
	return client
}

func TestClient(t *testing.T) {
	client := setupTestClient(t)
	ctx := context.Background()

	// Test configuration
	config := &Config{
		ID:          "test_device",
		Name:        "Test Device",
		Description: "Test device for unit tests",
		SNMPSettings: struct {
			Host      string `json:"host"`
			Port      int    `json:"port"`
			Version   string `json:"version"`
			Community string `json:"community"`
			OIDs      []struct {
				OID         string `json:"oid"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Type        string `json:"type"`
			} `json:"oids"`
		}{
			Host:      "test.hedgehog.internal",
			Port:      161,
			Version:   "2c",
			Community: "public",
			OIDs: []struct {
				OID         string `json:"oid"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Type        string `json:"type"`
			}{
				{
					OID:         "1.3.6.1.2.1.1.3.0",
					Name:        "sysUpTime",
					Description: "System uptime",
					Type:        "gauge",
				},
			},
		},
		PrometheusSettings: struct {
			MetricPrefix string            `json:"metric_prefix"`
			Labels      map[string]string `json:"labels"`
		}{
			MetricPrefix: "test_device_",
			Labels: map[string]string{
				"environment": "test",
			},
		},
		Enabled: true,
	}

	// Test Create
	t.Run("Create", func(t *testing.T) {
		err := client.CreateConfig(ctx, config)
		require.NoError(t, err)

		// Wait for indexing
		time.Sleep(1 * time.Second)
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		retrieved, err := client.GetConfig(ctx, config.ID)
		require.NoError(t, err)
		assert.Equal(t, config.Name, retrieved.Name)
		assert.Equal(t, config.SNMPSettings.Host, retrieved.SNMPSettings.Host)
		assert.True(t, retrieved.CreatedAt.Before(time.Now()))
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		configs, err := client.ListConfigs(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, configs)
		found := false
		for _, c := range configs {
			if c.ID == config.ID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		config.Description = "Updated description"
		err := client.UpdateConfig(ctx, config)
		require.NoError(t, err)

		// Wait for indexing
		time.Sleep(1 * time.Second)

		updated, err := client.GetConfig(ctx, config.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", updated.Description)
		assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := client.DeleteConfig(ctx, config.ID)
		require.NoError(t, err)

		// Wait for indexing
		time.Sleep(1 * time.Second)

		_, err = client.GetConfig(ctx, config.ID)
		assert.Error(t, err)
	})
}
