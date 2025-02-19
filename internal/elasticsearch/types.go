package elasticsearch

import (
	"context"
	"time"
)

// MetricDocument represents a single metric data point in Elasticsearch
type MetricDocument struct {
	Timestamp   time.Time              `json:"@timestamp"`
	DeviceID    string                 `json:"device_id"`
	DeviceName  string                 `json:"device_name"`
	MetricName  string                 `json:"metric_name"`
	Value       float64                `json:"value"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Environment string                 `json:"environment"`
}

// Writer defines the interface for writing metrics to Elasticsearch
type Writer interface {
	// Write writes a batch of metric documents to Elasticsearch
	Write(ctx context.Context, docs []MetricDocument) error

	// WriteOne writes a single metric document to Elasticsearch
	WriteOne(ctx context.Context, doc MetricDocument) error

	// Close closes the writer and releases any resources
	Close() error
}

// WriterConfig holds configuration for the Elasticsearch writer
type WriterConfig struct {
	Addresses        []string      `json:"addresses" yaml:"addresses" toml:"addresses"`
	Username         string        `json:"username" yaml:"username" toml:"username"`
	Password         string        `json:"password" yaml:"password" toml:"password"`
	IndexPrefix      string        `json:"index_prefix" yaml:"index_prefix" toml:"index_prefix"`
	BatchSize        int           `json:"batch_size" yaml:"batch_size" toml:"batch_size"`
	FlushInterval    time.Duration `json:"flush_interval" yaml:"flush_interval" toml:"flush_interval"`
	CertificateHash  string        `json:"certificate_hash" yaml:"certificate_hash" toml:"certificate_hash"`
	RetryMaxAttempts int           `json:"retry_max_attempts" yaml:"retry_max_attempts" toml:"retry_max_attempts"`
	RetryWaitTime    time.Duration `json:"retry_wait_time" yaml:"retry_wait_time" toml:"retry_wait_time"`
}

// ValidationError represents an error during document validation
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
