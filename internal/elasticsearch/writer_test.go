package elasticsearch

import (
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *WriterConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &WriterConfig{
				IndexPrefix:   "metrics",
				BatchSize:    100,
				FlushInterval: time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing index prefix",
			cfg: &WriterConfig{
				BatchSize:    100,
				FlushInterval: time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid batch size",
			cfg: &WriterConfig{
				IndexPrefix:   "metrics",
				BatchSize:    0,
				FlushInterval: time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid flush interval",
			cfg: &WriterConfig{
				IndexPrefix:   "metrics",
				BatchSize:    100,
				FlushInterval: 0,
			},
			wantErr: true,
		},
		{
			name:    "nil config",
			cfg:     nil,
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

func TestGenerateDocumentID(t *testing.T) {
	doc1 := &MetricsDocument{
		DeviceID:   "device1",
		MetricName: "metric1",
		Timestamp:  time.Now(),
	}

	doc2 := &MetricsDocument{
		DeviceID:   "device1",
		MetricName: "metric1",
		Timestamp:  doc1.Timestamp,
	}

	id1 := generateDocumentID(doc1)
	id2 := generateDocumentID(doc2)

	if id1 != id2 {
		t.Errorf("Expected identical IDs for same document content, got %s and %s", id1, id2)
	}
}
