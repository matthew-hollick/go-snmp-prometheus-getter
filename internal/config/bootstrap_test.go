package config

import (
	"os"
	"testing"
)

func TestLoadBootstrapConfiguration(t *testing.T) {
	configContent := `
[instance]
name = "test-instance"
log_level = "info"

[elasticsearch]
hosts = ["https://elasticsearch.hedgehog.internal:9200"]
index = "service_configuration"
certificate_hash = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

[elasticsearch.auth]
username = "hedgehog_config_reader"
password = ""

[concurrency]
max_scrapers = 5
max_writers = 3

[timing]
config_reload_interval = "5m"
scrape_timeout = "30s"
write_timeout = "10s"

[backoff]
initial_interval = "1s"
max_interval = "1m"
max_retries = 3
multiplier = 2.0

[metrics]
port = 9100
path = "/metrics"
`

	tmpfile, err := os.CreateTemp("", "config-*.toml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(configContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	configuration, err := LoadBootstrapConfiguration(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load valid configuration: %v", err)
	}

	// Verify configuration values
	expectedHost := "https://elasticsearch.hedgehog.internal:9200"
	if len(configuration.Elasticsearch.Hosts) != 1 || configuration.Elasticsearch.Hosts[0] != expectedHost {
		t.Errorf("Expected host to be %s, got %s", expectedHost, configuration.Elasticsearch.Hosts[0])
	}

	if configuration.Instance.LogLevel != "info" {
		t.Errorf("Expected log level to be info, got %s", configuration.Instance.LogLevel)
	}

	if configuration.Concurrency.MaxScrapers != 5 {
		t.Errorf("Expected max scrapers to be 5, got %d", configuration.Concurrency.MaxScrapers)
	}

	// Test invalid configuration file
	_, err = LoadBootstrapConfiguration("nonexistent.toml")
	if err == nil {
		t.Error("Expected error when loading nonexistent configuration file")
	}
}
