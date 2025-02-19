package config

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// BootstrapConfiguration represents the initial configuration needed to start the service.
type BootstrapConfiguration struct {
	Instance      InstanceSettings      `toml:"instance"`
	Elasticsearch ElasticsearchSettings `toml:"elasticsearch"`
	Concurrency   ConcurrencySettings   `toml:"concurrency"`
	Timing        TimingSettings        `toml:"timing"`
	Backoff       BackoffSettings       `toml:"backoff"`
	Metrics       MetricsSettings       `toml:"metrics"`
}

// InstanceSettings contains instance identification and basic settings.
type InstanceSettings struct {
	Name     string `toml:"name"`
	LogLevel string `toml:"log_level"`
}

// ElasticsearchSettings contains the connection details for Elasticsearch.
type ElasticsearchSettings struct {
	Hosts           []string     `toml:"hosts"`
	Index           string       `toml:"index"`
	CertificateHash string       `toml:"certificate_hash"`
	Auth            AuthSettings `toml:"auth"`
}

// AuthSettings contains authentication details.
type AuthSettings struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// ConcurrencySettings controls parallel operations.
type ConcurrencySettings struct {
	MaxScrapers int `toml:"max_scrapers"`
	MaxWriters  int `toml:"max_writers"`
}

// TimingSettings controls various timeouts and intervals.
type TimingSettings struct {
	ConfigReloadInterval Duration `toml:"config_reload_interval"`
	ScrapeTimeout       Duration `toml:"scrape_timeout"`
	WriteTimeout        Duration `toml:"write_timeout"`
}

// BackoffSettings controls retry behaviour.
type BackoffSettings struct {
	InitialInterval Duration `toml:"initial_interval"`
	MaxInterval     Duration `toml:"max_interval"`
	MaxRetries      int      `toml:"max_retries"`
	Multiplier      float64  `toml:"multiplier"`
}

// MetricsSettings controls the collector's own metrics endpoint.
type MetricsSettings struct {
	Port int    `toml:"port"`
	Path string `toml:"path"`
}

// Duration is a wrapper around time.Duration for TOML parsing.
type Duration struct {
	time.Duration
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))

	return err
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.String()), nil
}

// LoadBootstrapConfiguration loads the bootstrap configuration from a TOML file.
func LoadBootstrapConfiguration(path string) (*BootstrapConfiguration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading configuration file: %w", err)
	}

	var config BootstrapConfiguration
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing configuration: %w", err)
	}

	if err := validateConfiguration(&config); err != nil {
		return nil, fmt.Errorf("validating configuration: %w", err)
	}

	return &config, nil
}

// validateConfiguration performs basic validation of the configuration.
func validateConfiguration(cfg *BootstrapConfiguration) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if len(cfg.Elasticsearch.Hosts) == 0 {
		return fmt.Errorf("at least one Elasticsearch host must be specified")
	}

	if cfg.Elasticsearch.Index == "" {
		return fmt.Errorf("Elasticsearch index must be specified")
	}

	if cfg.Concurrency.MaxScrapers < 1 {
		return fmt.Errorf("max scrapers must be at least 1")
	}

	if cfg.Concurrency.MaxWriters < 1 {
		return fmt.Errorf("max writers must be at least 1")
	}

	if cfg.Timing.ConfigReloadInterval.Duration < time.Second {
		return fmt.Errorf("config reload interval must be at least 1 second")
	}

	if cfg.Timing.ScrapeTimeout.Duration < time.Second {
		return fmt.Errorf("scrape timeout must be at least 1 second")
	}

	if cfg.Timing.WriteTimeout.Duration < time.Second {
		return fmt.Errorf("write timeout must be at least 1 second")
	}

	return nil
}
