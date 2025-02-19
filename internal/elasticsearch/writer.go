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
func NewWriter(esclient *elasticsearch.Client, cfg WriterConfig) (*ESWriter, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("validating configuration: %w", err)
	}

	info, err := esclient.Info()
	if err != nil {
		return nil, fmt.Errorf("getting cluster info: %w", err)
	}
	defer info.Body.Close()

	var clusterInfo struct {
		Version struct {
			Number string `json:"number"`
		} `json:"version"`
	}

	if err := json.NewDecoder(info.Body).Decode(&clusterInfo); err != nil {
		return nil, fmt.Errorf("decoding cluster info: %w", err)
	}

	bi, err := createBulkIndexer(esclient, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create bulk indexer: %w", err)
	}

	return &ESWriter{
		client:        esclient,
		bulkIndexer:   bi,
		config:        cfg,
		indexPrefix:   cfg.IndexPrefix,
		isInitialized: true,
	}, nil
}

func validateConfig(cfg *WriterConfig) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	if cfg.IndexPrefix == "" {
		return fmt.Errorf("index prefix is required")
	}

	if cfg.BatchSize <= 0 {
		return fmt.Errorf("batch size must be positive")
	}

	if cfg.FlushInterval <= 0 {
		return fmt.Errorf("flush interval must be positive")
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

	if err := validateDocument(doc); err != nil {
		return fmt.Errorf("document validation failed: %w", err)
	}

	indexName := fmt.Sprintf("%s-%s", w.indexPrefix, doc.Timestamp.Format("2006.01.02"))

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

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
		return fmt.Errorf("timestamp is required")
	}
	if doc.DeviceID == "" {
		return fmt.Errorf("device ID is required")
	}
	if doc.MetricName == "" {
		return fmt.Errorf("metric name is required")
	}
	return nil
}

func generateDocumentID(doc MetricDocument) string {
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
