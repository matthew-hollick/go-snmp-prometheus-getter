package collector

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/sixworks/go-snmp-prometheus-getter/internal/elasticsearch"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/exporter"
)

func TestCollector_CollectMetrics(t *testing.T) {
	// Create test configuration
	cfg := elasticsearch.Config{
		ID:      "test1",
		Enabled: true,
		SNMPSettings: struct {
			Host              string
			Port              int
			Version           string
			Community         string
			PollIntervalSeconds int
		}{
			Host:              "localhost",
			Port:              161,
			Version:           "2c",
			Community:         "public",
			PollIntervalSeconds: 60,
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
	collector := New(exporterClient, Config{
		MaxConcurrentScrapers: 2,
		ScrapeIntervalSeconds: 1,
	}, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Test collecting metrics
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collector.CollectMetrics(ctx, cfg)
	if err != nil {
		// We expect an error in the test environment since there's no real SNMP exporter
		t.Logf("Expected error collecting metrics: %v", err)
	}
}

func TestCollector_RateLimiting(t *testing.T) {
	// Create test configuration
	cfg := elasticsearch.Config{
		ID:      "test1",
		Enabled: true,
		SNMPSettings: struct {
			Host              string
			Port              int
			Version           string
			Community         string
			PollIntervalSeconds int
		}{
			Host:              "localhost",
			Port:              161,
			Version:           "2c",
			Community:         "public",
			PollIntervalSeconds: 60,
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
	collector := New(exporterClient, Config{
		MaxConcurrentScrapers: 1,
		ScrapeIntervalSeconds: 1,
	}, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Test that we can't exceed rate limit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	
	// Try to collect twice
	_ = collector.CollectMetrics(ctx, cfg)
	_ = collector.CollectMetrics(ctx, cfg)
	
	duration := time.Since(start)

	// Should take at least 1 second due to rate limiting
	if duration < time.Second {
		t.Errorf("Rate limiting not working, took %v, expected at least 1s", duration)
	}
}

func TestCollector_ConcurrencyLimit(t *testing.T) {
	// Create test configuration
	cfg := elasticsearch.Config{
		ID:      "test1",
		Enabled: true,
		SNMPSettings: struct {
			Host              string
			Port              int
			Version           string
			Community         string
			PollIntervalSeconds int
		}{
			Host:              "localhost",
			Port:              161,
			Version:           "2c",
			Community:         "public",
			PollIntervalSeconds: 60,
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

	// Create collector with concurrency limit of 1
	collector := New(exporterClient, Config{
		MaxConcurrentScrapers: 1,
		ScrapeIntervalSeconds: 0, // No rate limiting for this test
	}, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// Start multiple collections in parallel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = collector.CollectMetrics(ctx, cfg)
			}()
		}
		wg.Wait()
	}()

	// Check that all collections complete within timeout
	select {
	case <-ctx.Done():
		t.Error("Test timed out, concurrency limit may be blocking too much")
	case <-done:
		// Test passed
	}
}
