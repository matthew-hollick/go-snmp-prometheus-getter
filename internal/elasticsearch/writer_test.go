package elasticsearch

import (
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     WriterConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: WriterConfig{
				Addresses:     []string{"http://localhost:9200"},
				IndexPrefix:  "metrics",
				BatchSize:    100,
				FlushInterval: time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing addresses",
			cfg: WriterConfig{
				IndexPrefix:  "metrics",
				BatchSize:    100,
				FlushInterval: time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing index prefix",
			cfg: WriterConfig{
				Addresses:     []string{"http://localhost:9200"},
				BatchSize:    100,
				FlushInterval: time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid batch size",
			cfg: WriterConfig{
				Addresses:     []string{"http://localhost:9200"},
				IndexPrefix:  "metrics",
				BatchSize:    0,
				FlushInterval: time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid flush interval",
			cfg: WriterConfig{
				Addresses:     []string{"http://localhost:9200"},
				IndexPrefix:  "metrics",
				BatchSize:    100,
				FlushInterval: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     MetricDocument
		wantErr bool
	}{
		{
			name: "valid document",
			doc: MetricDocument{
				Timestamp:  time.Now(),
				DeviceID:   "device1",
				MetricName: "cpu_usage",
				Value:      42.5,
			},
			wantErr: false,
		},
		{
			name: "missing timestamp",
			doc: MetricDocument{
				DeviceID:   "device1",
				MetricName: "cpu_usage",
				Value:      42.5,
			},
			wantErr: true,
		},
		{
			name: "missing device ID",
			doc: MetricDocument{
				Timestamp:  time.Now(),
				MetricName: "cpu_usage",
				Value:      42.5,
			},
			wantErr: true,
		},
		{
			name: "missing metric name",
			doc: MetricDocument{
				Timestamp: time.Now(),
				DeviceID:  "device1",
				Value:     42.5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDocument(tt.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateDocumentID(t *testing.T) {
	timestamp1 := time.Date(2025, 2, 18, 16, 0, 0, 0, time.UTC)
	timestamp2 := time.Date(2025, 2, 18, 16, 1, 0, 0, time.UTC)

	doc1 := MetricDocument{
		Timestamp:  timestamp1,
		DeviceID:   "device1",
		MetricName: "cpu_usage",
		Value:      42.5,
	}

	doc2 := MetricDocument{
		Timestamp:  timestamp1,
		DeviceID:   "device1",
		MetricName: "cpu_usage",
		Value:      42.5,
	}

	doc3 := MetricDocument{
		Timestamp:  timestamp2,
		DeviceID:   "device1",
		MetricName: "cpu_usage",
		Value:      42.5,
	}

	// Same document should generate same ID
	id1 := generateDocumentID(doc1)
	id2 := generateDocumentID(doc2)
	if id1 != id2 {
		t.Errorf("generateDocumentID() generated different IDs for same document: %v != %v", id1, id2)
	}

	// Different timestamps should generate different IDs
	id3 := generateDocumentID(doc3)
	if id1 == id3 {
		t.Errorf("generateDocumentID() generated same ID for different documents: %v == %v", id1, id3)
	}
}
