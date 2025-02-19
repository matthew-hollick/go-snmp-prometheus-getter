package exporter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a client for the SNMP exporter.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Config represents the configuration for the SNMP exporter client.
type Config struct {
	// BaseURL is the base URL of the SNMP exporter (e.g., "http://localhost:9116").
	BaseURL string
	// Timeout is the timeout for HTTP requests.
	Timeout time.Duration
}

// NewClient creates a new SNMP exporter client.
func NewClient(cfg Config) (*Client, error) {
	// Validate and normalize the base URL.
	if !strings.HasPrefix(cfg.BaseURL, "http://") && !strings.HasPrefix(cfg.BaseURL, "https://") {
		return nil, fmt.Errorf("base URL must start with http:// or https://")
	}

	// Create HTTP client with timeout.
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	return &Client{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		httpClient: httpClient,
	}, nil
}

// QueryParams represents the parameters for querying the SNMP exporter.
type QueryParams struct {
	// Target is the SNMP device to query (required).
	Target string
	// Module is the SNMP module to use (optional, can be multiple comma-separated values).
	Module []string
	// Auth is the authentication configuration to use (optional).
	Auth string
	// Context is the SNMP context to use (optional).
	Context string
	// Transport is the transport protocol to use (optional, e.g., "udp" or "tcp").
	Transport string
	// Port is the port to use (optional, defaults to 161).
	Port int
	// Version is the SNMP version to use (optional, defaults to 2c).
	Version string
	// Timeout is the timeout for SNMP requests (optional).
	Timeout string
	// Retries is the number of retries for failed requests (optional).
	Retries int
}

// buildTargetString builds the target string with optional transport and port.
func (p *QueryParams) buildTargetString() string {
	var target string
	if p.Transport != "" {
		target = fmt.Sprintf("%s://%s", p.Transport, p.Target)
	} else {
		target = p.Target
	}

	if p.Port > 0 {
		target = fmt.Sprintf("%s:%d", target, p.Port)
	}

	return target
}

// GetMetrics queries the SNMP exporter for metrics.
func (c *Client) GetMetrics(ctx context.Context, params *QueryParams) ([]byte, error) {
	// Build query parameters.
	query := url.Values{}

	// Add target parameter (URL encoded).
	query.Set("target", params.buildTargetString())

	// Add module parameter if specified.
	if len(params.Module) > 0 {
		query.Set("module", strings.Join(params.Module, ","))
	}

	// Add auth parameter if specified.
	if params.Auth != "" {
		query.Set("auth", params.Auth)
	}

	// Add context parameter if specified.
	if params.Context != "" {
		query.Set("snmp_context", params.Context)
	}

	// Add version parameter if specified.
	if params.Version != "" {
		query.Set("version", params.Version)
	}

	// Add timeout parameter if specified.
	if params.Timeout != "" {
		query.Set("timeout", params.Timeout)
	}

	// Add retries parameter if specified.
	if params.Retries > 0 {
		query.Set("retries", fmt.Sprintf("%d", params.Retries))
	}

	// Parse base URL.
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Add query parameters to URL.
	baseURL.RawQuery = query.Encode()

	// Create request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var closeErr error
	// Send request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			closeErr = fmt.Errorf("failed to close response body: %w", cerr)
		}
	}()

	// Read response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check response status.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	if closeErr != nil {
		return nil, closeErr
	}

	return body, nil
}
