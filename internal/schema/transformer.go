package schema

import (
	"fmt"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// Transformer converts Prometheus metrics to ECS documents.
type Transformer struct {
	observerHostname string
	observerVersion  string
}

// NewTransformer creates a new transformer instance.
func NewTransformer(observerHostname, observerVersion string) *Transformer {
	return &Transformer{
		observerHostname: observerHostname,
		observerVersion:  observerVersion,
	}
}

// TransformMetrics converts raw metrics data to a document.
func (t *Transformer) TransformMetrics(target string, metricsData []byte) (*Document, error) {
	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(strings.NewReader(string(metricsData)))

	if err != nil {
		return nil, fmt.Errorf("parsing metrics: %w", err)
	}

	now := time.Now().UTC()

	doc := &Document{
		Timestamp: now,
		Event: EventInfo{
			Kind:     "metric",
			Category: "network",
			Type:     "info",
			Created:  now,
		},
		Host: HostInfo{
			Hostname: target,
			Type:     "network-device",
		},
		Observer: ObserverInfo{
			Type:     "snmp-collector",
			Version:  t.observerVersion,
			Hostname: t.observerHostname,
		},
		SNMP: SNMPMetrics{
			SysInfo:    make(map[string]interface{}),
			Interfaces: make(map[string]interface{}),
			Metrics:    make(map[string]interface{}),
			Resources:  make([]Resource, 0),
		},
	}

	// Process metrics
	for name, family := range metrics {
		for _, metric := range family.Metric {
			resource := Resource{
				Name:   name,
				Type:   family.GetType().String(),
				Values: make(map[string]interface{}),
			}

			switch family.GetType() {
			case dto.MetricType_COUNTER:
				resource.Values["value"] = metric.Counter.GetValue()
			case dto.MetricType_GAUGE:
				resource.Values["value"] = metric.Gauge.GetValue()
			case dto.MetricType_UNTYPED:
				resource.Values["value"] = metric.GetUntyped().GetValue()
			}

			// Add labels as values
			for _, label := range metric.Label {
				resource.Values[label.GetName()] = label.GetValue()
			}

			doc.SNMP.Resources = append(doc.SNMP.Resources, resource)
		}
	}

	return doc, nil
}

// ValidateDocument checks if a document meets our schema requirements.
func (t *Transformer) ValidateDocument(doc *Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}

	if doc.Host.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	if doc.Observer.Type == "" {
		return fmt.Errorf("observer type is required")
	}

	if len(doc.SNMP.Resources) == 0 {
		return fmt.Errorf("at least one metric resource is required")
	}

	return nil
}
