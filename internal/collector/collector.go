package collector

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/sixworks/go-snmp-prometheus-getter/internal/elasticsearch"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/exporter"
)

// Collector manages the collection of SNMP metrics
type Collector struct {
	exporterClient *exporter.Client
	logger        *slog.Logger
	rateLimiter   *time.Ticker
	workerPool    chan struct{}
	wg            sync.WaitGroup
}

// Config represents the collector configuration
type Config struct {
	MaxConcurrentScrapers int
	ScrapeIntervalSeconds int
}

// New creates a new collector
func New(exporterClient *exporter.Client, cfg Config, logger *slog.Logger) *Collector {
	return &Collector{
		exporterClient: exporterClient,
		logger:        logger,
		rateLimiter:   time.NewTicker(time.Duration(cfg.ScrapeIntervalSeconds) * time.Second),
		workerPool:    make(chan struct{}, cfg.MaxConcurrentScrapers),
	}
}

// CollectMetrics collects metrics for a given configuration
func (c *Collector) CollectMetrics(ctx context.Context, cfg elasticsearch.Config) error {
	// Wait for rate limiter
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.rateLimiter.C:
	}

	// Acquire worker from pool
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.workerPool <- struct{}{}:
	}
	defer func() { <-c.workerPool }()

	// Prepare query parameters
	params := exporter.QueryParams{
		Target:    cfg.SNMPSettings.Host,
		Port:      cfg.SNMPSettings.Port,
		Module:    []string{}, // TODO: Determine modules based on OIDs
		Auth:      cfg.SNMPSettings.Community,
	}

	// Log collection attempt
	c.logger.Debug("Collecting metrics",
		"target", params.Target,
		"port", params.Port,
		"auth", params.Auth)

	// Collect metrics with retries
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		metrics, err := c.exporterClient.GetMetrics(ctx, params)
		if err == nil {
			if err := c.validateResponse(metrics); err != nil {
				lastErr = fmt.Errorf("invalid response: %w", err)
				continue
			}

			// Log collected metrics in debug mode
			c.logger.Debug("Collected metrics",
				"target", params.Target,
				"metrics", string(metrics),
				"config_id", cfg.ID)

			return nil
		}
		lastErr = err
		
		c.logger.Debug("Collection attempt failed",
			"target", params.Target,
			"attempt", attempt,
			"error", err)

		// Check if we should retry
		if attempt < 3 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
				continue
			}
		}
	}

	return fmt.Errorf("failed after 3 attempts: %w", lastErr)
}

// validateResponse checks if the response contains valid metrics
func (c *Collector) validateResponse(metrics []byte) error {
	if len(metrics) == 0 {
		return fmt.Errorf("empty response")
	}

	// TODO: Add more validation:
	// 1. Check if response is valid Prometheus format
	// 2. Verify expected metrics are present
	// 3. Check value ranges if specified in configuration

	return nil
}

// Start begins collecting metrics for all configurations
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

// Stop gracefully stops the collector
func (c *Collector) Stop() {
	c.rateLimiter.Stop()
	c.wg.Wait()
}
