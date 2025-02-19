package config

import (
	"os"
	"testing"
)

func TestLoadBootstrapConfiguration(t *testing.T) {
	// Create a temporary configuration file
	configContent := `
config_elasticsearch:
  hosts:
    - https://elasticsearch.hedgehog.internal:9200
  index: service_configuration
  certificate_hash: "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  auth:
    username: hedgehog_config_reader
    password: ""

app_settings:
  log_level: "info"
  config_refresh_minutes: 5
  concurrency:
    max_concurrent_scrapers: 5
    max_concurrent_writers: 3
  retry:
    max_attempts: 3
    delay_seconds: 5
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Test successful configuration loading
	configuration, err := LoadBootstrapConfiguration(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load valid configuration: %v", err)
	}

	// Verify configuration values
	if len(configuration.ElasticsearchSettings.Hosts) != 1 {
		t.Errorf("Expected 1 host, got %d", len(configuration.ElasticsearchSettings.Hosts))
	}

	if configuration.ElasticsearchSettings.Hosts[0] != "https://elasticsearch.hedgehog.internal:9200" {
		t.Errorf("Unexpected host value: %s", configuration.ElasticsearchSettings.Hosts[0])
	}

	if configuration.ApplicationSettings.RefreshMinutes != 5 {
		t.Errorf("Expected refresh minutes to be 5, got %d", configuration.ApplicationSettings.RefreshMinutes)
	}

	// Test invalid configurations
	invalidConfigs := []struct {
		name    string
		content string
	}{
		{
			name: "missing hosts",
			content: `
config_elasticsearch:
  index: service_configuration
  auth:
    username: hedgehog_config_reader
app_settings:
  log_level: "info"
  config_refresh_minutes: 5
  concurrency:
    max_concurrent_scrapers: 5
    max_concurrent_writers: 3
  retry:
    max_attempts: 3
    delay_seconds: 5
`,
		},
		{
			name: "invalid refresh minutes",
			content: `
config_elasticsearch:
  hosts:
    - https://elasticsearch.hedgehog.internal:9200
  index: service_configuration
  auth:
    username: hedgehog_config_reader
app_settings:
  log_level: "info"
  config_refresh_minutes: 0
  concurrency:
    max_concurrent_scrapers: 5
    max_concurrent_writers: 3
  retry:
    max_attempts: 3
    delay_seconds: 5
`,
		},
	}

	for _, tc := range invalidConfigs {
		t.Run(tc.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "invalid-config-*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tc.content)); err != nil {
				t.Fatalf("Failed to write configuration: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temporary file: %v", err)
			}

			if _, err := LoadBootstrapConfiguration(tmpfile.Name()); err == nil {
				t.Error("Expected error for invalid configuration, but got none")
			}
		})
	}
}
