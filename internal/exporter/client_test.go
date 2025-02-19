package exporter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			cfg: Config{
				BaseURL: "http://localhost:9116",
				Timeout: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid base URL",
			cfg: Config{
				BaseURL: "localhost:9116",
				Timeout: 5 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryParams_buildTargetString(t *testing.T) {
	tests := []struct {
		name   string
		params QueryParams
		want   string
	}{
		{
			name: "target only",
			params: QueryParams{
				Target: "192.0.0.8",
			},
			want: "192.0.0.8",
		},
		{
			name: "target with transport",
			params: QueryParams{
				Target:    "192.0.0.8",
				Transport: "tcp",
			},
			want: "tcp://192.0.0.8",
		},
		{
			name: "target with port",
			params: QueryParams{
				Target: "192.0.0.8",
				Port:   1161,
			},
			want: "192.0.0.8:1161",
		},
		{
			name: "target with transport and port",
			params: QueryParams{
				Target:    "192.0.0.8",
				Transport: "tcp",
				Port:      1161,
			},
			want: "tcp://192.0.0.8:1161",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.buildTargetString(); got != tt.want {
				t.Errorf("QueryParams.buildTargetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetMetrics(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/snmp" {
			t.Errorf("Expected path /snmp, got %s", r.URL.Path)
		}

		// Check query parameters
		query := r.URL.Query()
		if target := query.Get("target"); target != "tcp://192.0.0.8:1161" {
			t.Errorf("Expected target tcp://192.0.0.8:1161, got %s", target)
		}
		if module := query.Get("module"); module != "if_mib,system" {
			t.Errorf("Expected module if_mib,system, got %s", module)
		}
		if auth := query.Get("auth"); auth != "public_v2" {
			t.Errorf("Expected auth public_v2, got %s", auth)
		}

		// Send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test metrics"))
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetMetrics
	params := QueryParams{
		Target:    "192.0.0.8",
		Module:    []string{"if_mib", "system"},
		Auth:      "public_v2",
		Transport: "tcp",
		Port:      1161,
	}

	metrics, err := client.GetMetrics(context.Background(), params)
	if err != nil {
		t.Fatalf("GetMetrics() error = %v", err)
	}

	if string(metrics) != "test metrics" {
		t.Errorf("GetMetrics() = %v, want %v", string(metrics), "test metrics")
	}
}
