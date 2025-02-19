package collector

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/matthew-hollick/go-snmp-prometheus-getter/internal/elasticsearch"
	"github.com/matthew-hollick/go-snmp-prometheus-getter/internal/exporter"
)

func TestCollector_CollectMetrics(t *testing.T) {
	// Create test configuration.
	cfg := elasticsearch.Config{
		ID:      "test_device",
		Enabled: true,
		SNMPSettings: elasticsearch.SNMPSettings{
			Host:                "test.hedgehog.internal",
			Port:                161,
			Version:             "2c",
			Community:           "public",
			PollIntervalSeconds: 60,
			AuthName:            "test_auth",
			Timeout:             "5s",
			Retries:             3,
		},
	}

	// Create test exporter client
	exporterClient, err := exporter.NewClient(exporter.Config{
		BaseURL: "http://localhost:9116",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create exporter client: %v", err)
	}

	// Create collector
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	collector := New(exporterClient, Config{
		MaxConcurrentScrapers: 1,
		ScrapeIntervalSeconds: 1,
	}, logger)

	// Test collecting metrics
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collector.CollectMetrics(ctx, cfg)
	if err != nil {
		t.Errorf("CollectMetrics() error = %v", err)
	}
}

func TestCollector_RateLimiting(t *testing.T) {
	// Create test configuration
	cfg := elasticsearch.Config{
		ID:      "test_device",
		Enabled: true,
		SNMPSettings: elasticsearch.SNMPSettings{
			Host:                "test.hedgehog.internal",
			Port:                161,
			Version:             "2c",
			Community:           "public",
			PollIntervalSeconds: 60,
			AuthName:            "test_auth",
			Timeout:             "5s",
			Retries:             3,
		},
	}

	// Create test exporter client
	exporterClient, err := exporter.NewClient(exporter.Config{
		BaseURL: "http://localhost:9116",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create exporter client: %v", err)
	}

	// Create collector with 1 second rate limit
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	collector := New(exporterClient, Config{
		MaxConcurrentScrapers: 1,
		ScrapeIntervalSeconds: 1,
	}, logger)

	// Test that we can't exceed rate limit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()

	// Try to collect twice
	err = collector.CollectMetrics(ctx, cfg)
	if err != nil {
		t.Errorf("CollectMetrics() error = %v", err)
	}

	err = collector.CollectMetrics(ctx, cfg)
	if err != nil {
		t.Errorf("CollectMetrics() error = %v", err)
	}

	elapsed := time.Since(start)

	// Should take at least 1 second due to rate limiting
	if elapsed < time.Second {
		t.Errorf("Rate limiting not working, elapsed time = %v", elapsed)
	}
}

func TestCollector_ConcurrencyLimit(t *testing.T) {
	// Create test configuration
	cfg := elasticsearch.Config{
		ID:      "test_device",
		Enabled: true,
		SNMPSettings: elasticsearch.SNMPSettings{
			Host:                "test.hedgehog.internal",
			Port:                161,
			Version:             "2c",
			Community:           "public",
			PollIntervalSeconds: 60,
			AuthName:            "test_auth",
			Timeout:             "5s",
			Retries:             3,
		},
	}

	// Create test collector with concurrency limit of 2
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	exporterClient := &exporter.Client{}
	collector := New(exporterClient, Config{
		MaxConcurrentScrapers: 2,
		ScrapeIntervalSeconds: 0,
	}, logger)

	// Test concurrent collections
	ctx := context.Background()
	start := time.Now()

	var wg sync.WaitGroup

	// Start 3 concurrent collections
	for i := 0; i < 3; i++ {
		wg.Add(1)
		
		go func() {
			defer wg.Done()
			
			err := collector.CollectMetrics(ctx, cfg)
			if err != nil {
				t.Errorf("CollectMetrics() error = %v", err)
			}
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)
	if elapsed < time.Second {
		t.Errorf("Concurrency limiting not working, elapsed time = %v", elapsed)
	}
}
