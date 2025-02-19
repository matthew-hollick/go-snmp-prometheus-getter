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
	esclient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	require.NoError(t, err)
	return NewClient(esclient, "service_configuration")
}

func TestClient(t *testing.T) {
	client := setupTestClient(t)
	ctx := context.Background()

	// Test configuration
	config := Config{
		ID:      "test_device",
		Name:    "Test Device",
		Type:    "network_device",
		Enabled: true,
		SNMPSettings: SNMPSettings{
			Host:                "test.hedgehog.internal",
			Port:                161,
			Version:             "2c",
			Community:           "public",
			AuthName:            "test_auth",
			Timeout:             "5s",
			Retries:             3,
			PollIntervalSeconds: 60,
		},
		CollectorSettings: CollectorSettings{
			Hostname:           "test.hedgehog.internal",
			Version:            "1.0.0",
			Modules:            []string{"if_mib", "system"},
			CollectionInterval: "60s",
			Metrics: MetricsSettings{
				Include: []string{"ifInOctets", "ifOutOctets"},
			},
		},
		Tags: Tags{
			Environment: "test",
			Location:    "test_datacenter",
			Role:        "test_device",
		},
	}

	// Test Save
	t.Run("Save", func(t *testing.T) {
		err := client.SaveConfig(ctx, config)
		require.NoError(t, err)

		// Wait for indexing
		time.Sleep(1 * time.Second)
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		configs, err := client.ListConfigs(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, configs)
		
		var found bool
		for _, c := range configs {
			if c.ID == config.ID {
				found = true
				assert.Equal(t, config.Name, c.Name)
				assert.Equal(t, config.SNMPSettings.Host, c.SNMPSettings.Host)
				assert.True(t, c.CreatedAt.Before(time.Now()))
				break
			}
		}
		assert.True(t, found, "Created config not found in list")
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := client.DeleteConfig(ctx, config.ID)
		require.NoError(t, err)

		// Wait for deletion
		time.Sleep(1 * time.Second)

		configs, err := client.ListConfigs(ctx)
		require.NoError(t, err)
		
		for _, c := range configs {
			assert.NotEqual(t, config.ID, c.ID, "Config should have been deleted")
		}
	})

	// Test metrics storage
	t.Run("StoreMetrics", func(t *testing.T) {
		doc := MetricsDocument{
			Timestamp:   time.Now(),
			DeviceID:   config.ID,
			DeviceName: config.Name,
			DeviceType: config.Type,
			Host:       config.SNMPSettings.Host,
			Environment: config.Tags.Environment,
			Location:    config.Tags.Location,
			Role:        config.Tags.Role,
			Metrics: map[string]interface{}{
				"ifInOctets":  float64(1000),
				"ifOutOctets": float64(2000),
			},
		}

		err := client.StoreMetrics(ctx, doc)
		require.NoError(t, err)
	})
}
