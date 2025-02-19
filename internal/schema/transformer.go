package schema

import (
	"fmt"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// Transformer converts Prometheus metrics to ECS documents
type Transformer struct {
	observerHostname string
	observerVersion  string
}

// NewTransformer creates a new transformer instance
func NewTransformer(observerHostname, observerVersion string) *Transformer {
	return &Transformer{
		observerHostname: observerHostname,
		observerVersion:  observerVersion,
	}
}

// TransformMetrics converts raw metrics data to a document
func (t *Transformer) TransformMetrics(target string, metricsData []byte) (*Document, error) {
	// Parse the metrics data
	parser := expfmt.TextParser{}
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
		Metrics: MetricsInfo{
			SNMP: SNMPMetrics{
				SysInfo:    make(map[string]interface{}),
				Interfaces: make(map[string]interface{}),
				Resources:  make(map[string]interface{}),
			},
		},
	}

	// Process each metric family
	for name, family := range metrics {
		// Skip empty metrics
		if len(family.Metric) == 0 {
			continue
		}

		// Get the first metric (we expect single values for system metrics)
		metric := family.Metric[0]

		// Process system metrics
		if strings.HasPrefix(name, "snmp_system_") {
			key := strings.TrimPrefix(name, "snmp_system_")
			switch family.GetType() {
			case dto.MetricType_COUNTER:
				doc.Metrics.SNMP.SysInfo[key] = metric.Counter.GetValue()
			case dto.MetricType_GAUGE:
				doc.Metrics.SNMP.SysInfo[key] = metric.Gauge.GetValue()
			default:
				doc.Metrics.SNMP.SysInfo[key] = metric.GetUntyped().GetValue()
			}
			continue
		}

		// Process interface metrics
		if strings.HasPrefix(name, "snmp_interface_") {
			key := strings.TrimPrefix(name, "snmp_interface_")
			switch family.GetType() {
			case dto.MetricType_COUNTER:
				doc.Metrics.SNMP.Interfaces[key] = metric.Counter.GetValue()
			case dto.MetricType_GAUGE:
				doc.Metrics.SNMP.Interfaces[key] = metric.Gauge.GetValue()
			default:
				doc.Metrics.SNMP.Interfaces[key] = metric.GetUntyped().GetValue()
			}
			continue
		}

		// Process resource metrics
		if strings.HasPrefix(name, "snmp_resource_") {
			key := strings.TrimPrefix(name, "snmp_resource_")
			switch family.GetType() {
			case dto.MetricType_COUNTER:
				doc.Metrics.SNMP.Resources[key] = metric.Counter.GetValue()
			case dto.MetricType_GAUGE:
				doc.Metrics.SNMP.Resources[key] = metric.Gauge.GetValue()
			default:
				doc.Metrics.SNMP.Resources[key] = metric.GetUntyped().GetValue()
			}
			continue
		}
	}

	return doc, nil
}

// ValidateDocument checks if a document meets our schema requirements
func (t *Transformer) ValidateDocument(doc *Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}

	if doc.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	if doc.Host.Hostname == "" {
		return fmt.Errorf("host.hostname is required")
	}

	if doc.Host.Type == "" {
		return fmt.Errorf("host.type is required")
	}

	return nil
}
