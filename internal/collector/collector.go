// Package collector provides functionality for collecting SNMP metrics.
package collector

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/matthew-hollick/go-snmp-prometheus-getter/internal/elasticsearch"
	"github.com/matthew-hollick/go-snmp-prometheus-getter/internal/exporter"
)

// Collector manages the collection of SNMP metrics.
type Collector struct {
	exporterClient *exporter.Client
	logger         *slog.Logger
	rateLimiter    *time.Ticker
	workerPool     chan struct{}
	wg             sync.WaitGroup
}

// Config represents the collector configuration.
type Config struct {
	MaxConcurrentScrapers int
	ScrapeIntervalSeconds int
}

// New creates a new collector.
func New(exporterClient *exporter.Client, cfg Config, logger *slog.Logger) *Collector {
	return &Collector{
		exporterClient: exporterClient,
		logger:         logger,
		rateLimiter:    time.NewTicker(time.Duration(cfg.ScrapeIntervalSeconds) * time.Second),
		workerPool:     make(chan struct{}, cfg.MaxConcurrentScrapers),
	}
}

// CollectMetrics collects metrics for a given configuration.
func (c *Collector) CollectMetrics(ctx context.Context, cfg elasticsearch.Config) error {
	// Wait for rate limiter.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.rateLimiter.C:
	}

	// Acquire worker from pool.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.workerPool <- struct{}{}:
	}

	// Ensure worker is released when done.
	defer func() {
		<-c.workerPool
	}()

	// Prepare query parameters.
	params := exporter.QueryParams{
		Target:  cfg.SNMPSettings.Host,
		Port:    cfg.SNMPSettings.Port,
		Module:  []string{}, // TODO: Determine modules based on OIDs.
		Auth:    cfg.SNMPSettings.Community,
		Version: cfg.SNMPSettings.Version,
		Timeout: cfg.SNMPSettings.Timeout,
		Retries: cfg.SNMPSettings.Retries,
	}

	// Query the exporter.
	metrics, err := c.exporterClient.GetMetrics(ctx, &params)
	if err != nil {
		return fmt.Errorf("failed to query exporter: %w", err)
	}

	// Validate response.
	if err := c.validateResponse(metrics); err != nil {
		return fmt.Errorf("invalid response from exporter: %w", err)
	}

	return nil
}

// validateResponse checks if the response contains valid metrics.
func (c *Collector) validateResponse(metrics []byte) error {
	if len(metrics) == 0 {
		return fmt.Errorf("empty response")
	}

	// TODO: Add more validation:
	// 1. Check if response is valid Prometheus format.
	// 2. Verify expected metrics are present.
	// 3. Check value ranges if specified in configuration.

	return nil
}

// Start begins collecting metrics for all configurations.
func (c *Collector) Start(ctx context.Context, configs []elasticsearch.Config) error {
	for _, cfg := range configs {
		if !cfg.Enabled {
			continue
		}

		c.wg.Add(1)

		go func(cfg elasticsearch.Config) {
			defer c.wg.Done()

			ticker := time.NewTicker(time.Duration(cfg.SNMPSettings.PollIntervalSeconds) * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := c.CollectMetrics(ctx, cfg); err != nil {
						c.logger.Error("Failed to collect metrics",
							"error", err,
							"target", cfg.SNMPSettings.Host,
							"id", cfg.ID)
					}
				}
			}
		}(cfg)
	}

	return nil
}

// Stop gracefully stops the collector.
func (c *Collector) Stop() {
	c.rateLimiter.Stop()
	c.wg.Wait()
}
