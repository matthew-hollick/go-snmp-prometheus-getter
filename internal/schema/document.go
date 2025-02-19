package schema

import (
	"time"
)

// Document represents a metrics document
type Document struct {
	Timestamp time.Time    `json:"@timestamp"`
	Event     EventInfo    `json:"event"`
	Host      HostInfo     `json:"host"`
	Observer  ObserverInfo `json:"observer"`
	Metrics   MetricsInfo  `json:"metrics"`
}

// EventInfo contains event metadata
type EventInfo struct {
	Kind     string    `json:"kind"`
	Category string    `json:"category"`
	Type     string    `json:"type"`
	Created  time.Time `json:"created"`
}

// HostInfo contains information about the monitored host
type HostInfo struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
}

// ObserverInfo contains information about the collector
type ObserverInfo struct {
	Type     string `json:"type"`
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

// SNMPMetrics contains SNMP-specific metrics
type SNMPMetrics struct {
	SysInfo    map[string]interface{} `json:"system"`
	Interfaces map[string]interface{} `json:"interfaces"`
	Resources  map[string]interface{} `json:"resources"`
}

// MetricsInfo contains all metrics for a device
type MetricsInfo struct {
	SNMP SNMPMetrics `json:"snmp"`
}
