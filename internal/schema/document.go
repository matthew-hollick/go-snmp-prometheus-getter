package schema

import (
	"time"
)

// Document represents a metrics document.
type Document struct {
	Event     EventInfo    `json:"event"`
	Host      HostInfo     `json:"host"`
	Observer  ObserverInfo `json:"observer"`
	SNMP      SNMPMetrics  `json:"snmp"`
	Metrics   MetricsInfo  `json:"metrics"`
	Timestamp time.Time    `json:"@timestamp"`
}

// EventInfo contains event metadata.
type EventInfo struct {
	Created  time.Time `json:"created"`
	Kind     string    `json:"kind"`
	Category string    `json:"category"`
	Type     string    `json:"type"`
	Outcome  string    `json:"outcome"`
	Dataset  string    `json:"dataset"`
	Provider string    `json:"provider"`
}

// HostInfo contains information about the monitored host.
type HostInfo struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
}

// ObserverInfo contains information about the monitoring agent.
type ObserverInfo struct {
	Type     string `json:"type"`
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

// SNMPMetrics contains SNMP-specific metrics.
type SNMPMetrics struct {
	SysInfo    map[string]interface{} `json:"sys_info"`
	Interfaces map[string]interface{} `json:"interfaces"`
	Metrics    map[string]interface{} `json:"metrics"`
	Resources  []Resource            `json:"resources"`
}

// Resource represents a monitored SNMP resource.
type Resource struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Values      map[string]interface{} `json:"values"`
}

// MetricsInfo contains general metrics information.
type MetricsInfo struct {
	Name      string                 `json:"name"`
	Labels    map[string]string      `json:"labels"`
	Value     float64                `json:"value"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}
