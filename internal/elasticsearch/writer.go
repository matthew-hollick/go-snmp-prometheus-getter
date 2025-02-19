package elasticsearch

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

// ESWriter implements the Writer interface for Elasticsearch
type ESWriter struct {
	client        *elasticsearch.Client
	bulkIndexer   esutil.BulkIndexer
	config        WriterConfig
	indexPrefix   string
	mu            sync.RWMutex
	isInitialized bool
}

// NewWriter creates a new Elasticsearch writer
func NewWriter(cfg WriterConfig) (*ESWriter, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	// Verify connection and certificate if hash provided
	if err := verifyConnection(client, cfg.CertificateHash); err != nil {
		return nil, fmt.Errorf("failed to verify elasticsearch connection: %w", err)
	}

	bi, err := createBulkIndexer(client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create bulk indexer: %w", err)
	}

	return &ESWriter{
		client:        client,
		bulkIndexer:   bi,
		config:        cfg,
		indexPrefix:   cfg.IndexPrefix,
		isInitialized: true,
	}, nil
}

func validateConfig(cfg WriterConfig) error {
	if len(cfg.Addresses) == 0 {
		return &ValidationError{Field: "addresses", Message: "at least one address is required"}
	}
	if cfg.IndexPrefix == "" {
		return &ValidationError{Field: "index_prefix", Message: "index prefix is required"}
	}
	if cfg.BatchSize <= 0 {
		return &ValidationError{Field: "batch_size", Message: "batch size must be positive"}
	}
	if cfg.FlushInterval <= 0 {
		return &ValidationError{Field: "flush_interval", Message: "flush interval must be positive"}
	}
	return nil
}

func verifyConnection(client *elasticsearch.Client, certHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to get cluster info: %w", err)
	}
	defer info.Body.Close()

	if certHash != "" {
		// TODO: Implement certificate hash verification
		// This would involve getting the server's certificate and comparing its hash
	}

	return nil
}

func createBulkIndexer(client *elasticsearch.Client, cfg WriterConfig) (esutil.BulkIndexer, error) {
	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        client,
		NumWorkers:    2, // Number of workers can be made configurable
		FlushBytes:    5e6,
		FlushInterval: cfg.FlushInterval,
		OnError: func(ctx context.Context, err error) {
			// TODO: Add proper error handling and logging
			fmt.Printf("bulk indexer error: %v\n", err)
		},
	})
}

// Write implements the Writer interface
func (w *ESWriter) Write(ctx context.Context, docs []MetricDocument) error {
	w.mu.RLock()
	if !w.isInitialized {
		w.mu.RUnlock()
		return errors.New("writer is not initialized")
	}
	w.mu.RUnlock()

	for _, doc := range docs {
		if err := w.WriteOne(ctx, doc); err != nil {
			return fmt.Errorf("failed to write document: %w", err)
		}
	}

	return nil
}

// WriteOne implements the Writer interface
func (w *ESWriter) WriteOne(ctx context.Context, doc MetricDocument) error {
	w.mu.RLock()
	if !w.isInitialized {
		w.mu.RUnlock()
		return errors.New("writer is not initialized")
	}
	w.mu.RUnlock()

	// Validate document
	if err := validateDocument(doc); err != nil {
		return fmt.Errorf("document validation failed: %w", err)
	}

	// Create index name with date suffix
	indexName := fmt.Sprintf("%s-%s", w.indexPrefix, doc.Timestamp.Format("2006.01.02"))

	// Convert document to JSON bytes
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Add document to bulk indexer
	err = w.bulkIndexer.Add(ctx, esutil.BulkIndexerItem{
		Action:     "index",
		Index:      indexName,
		DocumentID: generateDocumentID(doc),
		Body:       bytes.NewReader(docJSON),
	})

	if err != nil {
		return fmt.Errorf("failed to add document to bulk indexer: %w", err)
	}

	return nil
}

func validateDocument(doc MetricDocument) error {
	if doc.Timestamp.IsZero() {
		return &ValidationError{Field: "timestamp", Message: "timestamp is required"}
	}
	if doc.DeviceID == "" {
		return &ValidationError{Field: "device_id", Message: "device ID is required"}
	}
	if doc.MetricName == "" {
		return &ValidationError{Field: "metric_name", Message: "metric name is required"}
	}
	return nil
}

func generateDocumentID(doc MetricDocument) string {
	// Create a unique ID based on device ID, metric name, and timestamp
	idStr := fmt.Sprintf("%s-%s-%d", doc.DeviceID, doc.MetricName, doc.Timestamp.UnixNano())
	hash := sha256.Sum256([]byte(idStr))
	return hex.EncodeToString(hash[:])
}

// Close implements the Writer interface
func (w *ESWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isInitialized {
		return nil
	}

	if err := w.bulkIndexer.Close(context.Background()); err != nil {
		return fmt.Errorf("failed to close bulk indexer: %w", err)
	}

	w.isInitialized = false
	return nil
}
