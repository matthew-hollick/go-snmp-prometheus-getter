package service

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	esapi "github.com/elastic/go-elasticsearch/v8"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/cache"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/config"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/elasticsearch"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/exporter"
	"github.com/sixworks/go-snmp-prometheus-getter/internal/schema"
)

// Service handles the main application logic
type Service struct {
	cfg            *config.BootstrapConfiguration
	esClient       *elasticsearch.Client
	transformer    *schema.Transformer
	logger         *slog.Logger
	configCache    *cache.ConfigCache
	configRefresh  *time.Ticker
	workerPool     chan struct{}
	writerPool     chan struct{}
	wg             sync.WaitGroup
}

// NewService creates a new service instance
func NewService(cfg *config.BootstrapConfiguration, logger *slog.Logger) (*Service, error) {
	// Create Elasticsearch client configuration
	escfg := esapi.Config{
		Addresses: cfg.Elasticsearch.Hosts,
	}

	// Configure TLS if certificate hash is provided
	if cfg.Elasticsearch.CertificateHash != "" {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				VerifyConnection: func(cs tls.ConnectionState) error {
					// Get the certificate hash
					certHash := sha256.Sum256(cs.PeerCertificates[0].Raw)
					certHashHex := hex.EncodeToString(certHash[:])

					// Compare with configured hash
					if certHashHex != cfg.Elasticsearch.CertificateHash {
						return fmt.Errorf("certificate hash mismatch: got %s, want %s",
							certHashHex, cfg.Elasticsearch.CertificateHash)
					}
					return nil
				},
			},
		}
		escfg.Transport = transport
	}

	// Configure authentication if provided
	if cfg.Elasticsearch.Auth.Username != "" {
		escfg.Username = cfg.Elasticsearch.Auth.Username
		escfg.Password = cfg.Elasticsearch.Auth.Password
	}

	// Create Elasticsearch API client
	esclient, err := esapi.NewClient(escfg)
	if err != nil {
		return nil, fmt.Errorf("creating elasticsearch client: %w", err)
	}

	// Create our Elasticsearch client wrapper
	esWrapper := elasticsearch.NewClient(esclient, cfg.Elasticsearch.Index)

	// Create service components
	transformer := schema.NewTransformer(cfg.Instance.Name, "1.0.0")
	configCache := cache.NewConfigCache(cfg.Timing.ConfigReloadInterval.Duration)
	configRefresh := time.NewTicker(cfg.Timing.ConfigReloadInterval.Duration)
	workerPool := make(chan struct{}, cfg.Concurrency.MaxScrapers)
	writerPool := make(chan struct{}, cfg.Concurrency.MaxWriters)

	return &Service{
		cfg:           cfg,
		esClient:      esWrapper,
		transformer:   transformer,
		logger:        logger,
		configCache:   configCache,
		configRefresh: configRefresh,
		workerPool:    workerPool,
		writerPool:    writerPool,
	}, nil
}

// Start begins the service operation
func (s *Service) Start(ctx context.Context) error {
	s.logger.Info("starting service",
		"instance", s.cfg.Instance.Name,
		"elasticsearch_hosts", s.cfg.Elasticsearch.Hosts,
		"max_scrapers", s.cfg.Concurrency.MaxScrapers,
		"max_writers", s.cfg.Concurrency.MaxWriters,
	)

	// Initial configuration load
	if err := s.refreshConfigurations(ctx); err != nil {
		return fmt.Errorf("initial configuration load failed: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("shutting down service")
			s.configRefresh.Stop()
			s.wg.Wait()
			return nil

		case <-s.configRefresh.C:
			if err := s.refreshConfigurations(ctx); err != nil {
				s.logger.Error("refreshing configurations", "error", err)
			}
		}
	}
}

// refreshConfigurations fetches and processes device configurations
func (s *Service) refreshConfigurations(ctx context.Context) error {
	// Check if cache is still valid
	if !s.configCache.IsExpired() {
		s.logger.Debug("using cached configurations",
			"count", s.configCache.Count(),
			"last_updated", s.configCache.LastUpdated(),
		)
		return nil
	}

	// Fetch configurations from Elasticsearch
	configs, err := s.esClient.ListConfigs(ctx)
	if err != nil {
		return fmt.Errorf("listing configurations: %w", err)
	}

	s.logger.Info("fetched configurations",
		"count", len(configs),
	)

	// Update cache
	s.configCache.SetAll(configs)

	// Process each configuration
	for _, cfg := range configs {
		if !cfg.Enabled {
			s.logger.Debug("skipping disabled device",
				"id", cfg.ID,
				"name", cfg.Name,
			)
			continue
		}

		if err := s.processConfiguration(ctx, &cfg); err != nil {
			s.logger.Error("processing configuration",
				"id", cfg.ID,
				"name", cfg.Name,
				"error", err,
			)
		}
	}

	return nil
}

// processConfiguration handles a single device configuration
func (s *Service) processConfiguration(ctx context.Context, cfg *elasticsearch.Config) error {
	// Create exporter client for this device's configuration
	hostname := cfg.CollectorSettings.Hostname
	if !strings.HasPrefix(hostname, "http://") && !strings.HasPrefix(hostname, "https://") {
		hostname = "http://" + hostname
	}
	exporterClient, err := exporter.NewClient(exporter.Config{
		BaseURL: hostname,
		Timeout: s.cfg.Timing.ScrapeTimeout.Duration,
	})
	if err != nil {
		return fmt.Errorf("creating exporter client: %w", err)
	}

	// Create backoff configuration
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = s.cfg.Backoff.InitialInterval.Duration
	b.MaxInterval = s.cfg.Backoff.MaxInterval.Duration
	b.Multiplier = s.cfg.Backoff.Multiplier
	b.MaxElapsedTime = s.cfg.Timing.ScrapeTimeout.Duration

	// Acquire worker from pool
	select {
	case s.workerPool <- struct{}{}:
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			defer func() { <-s.workerPool }()

			operation := func() error {
				return s.collectMetrics(ctx, cfg, exporterClient)
			}

			if err := backoff.Retry(operation, b); err != nil {
				s.logger.Error("collecting metrics after retries",
					"device", cfg.Name,
					"error", err,
				)
			}
		}()
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// collectMetrics collects and processes metrics for a device
func (s *Service) collectMetrics(ctx context.Context, cfg *elasticsearch.Config, exporterClient *exporter.Client) error {
	params := exporter.QueryParams{
		Target:    cfg.SNMPSettings.Host,
		Port:      cfg.SNMPSettings.Port,
		Transport: "udp",
		Module:    cfg.CollectorSettings.Modules,
		Auth:      cfg.SNMPSettings.AuthName,
	}

	metrics, err := exporterClient.GetMetrics(ctx, params)
	if err != nil {
		return fmt.Errorf("getting metrics: %w", err)
	}

	doc, err := s.transformer.TransformMetrics(cfg.SNMPSettings.Host, metrics)
	if err != nil {
		return fmt.Errorf("transforming metrics: %w", err)
	}

	// Acquire writer from pool for document processing
	select {
	case s.writerPool <- struct{}{}:
		defer func() { <-s.writerPool }()

		// Create metrics document
		metricsDoc := elasticsearch.MetricsDocument{
			Timestamp:   doc.Timestamp,
			DeviceID:    cfg.ID,
			DeviceName:  cfg.Name,
			DeviceType:  cfg.Type,
			Host:        cfg.SNMPSettings.Host,
			Environment: cfg.Tags.Environment,
			Location:    cfg.Tags.Location,
			Role:        cfg.Tags.Role,
			Metrics: map[string]interface{}{
				"system":     doc.Metrics.SNMP.SysInfo,
				"interfaces": doc.Metrics.SNMP.Interfaces,
				"resources":  doc.Metrics.SNMP.Resources,
			},
		}

		// Store metrics in Elasticsearch
		if err := s.esClient.StoreMetrics(ctx, metricsDoc); err != nil {
			s.logger.Error("storing metrics",
				"device", cfg.Name,
				"error", err,
			)
			return fmt.Errorf("storing metrics: %w", err)
		}

		if s.logger.Enabled(ctx, slog.LevelDebug) {
			jsonDoc, err := json.MarshalIndent(metricsDoc, "", "  ")
			if err != nil {
				s.logger.Warn("failed to marshal metrics document", "error", err)
			} else {
				s.logger.Debug("stored metrics document",
					"device", cfg.Name,
					"json", string(jsonDoc),
				)
			}
		} else {
			s.logger.Info("stored metrics",
				"device", cfg.Name,
				"host", cfg.SNMPSettings.Host,
				"modules", cfg.CollectorSettings.Modules,
				"timestamp", metricsDoc.Timestamp,
			)
		}

		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
