package elasticsearch

import (
	"fmt"
	"regexp"
)

var (
	// Compile regular expressions for validation
	durationRegex = regexp.MustCompile(`^\d+(ms|s|m|h)$`)
	versionRegex  = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
)

// ValidateConfig checks if a configuration is valid.
func ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if err := validateSNMPSettings(&config.SNMPSettings); err != nil {
		return fmt.Errorf("validating SNMP settings: %w", err)
	}

	if err := validateCollectorSettings(&config.CollectorSettings); err != nil {
		return fmt.Errorf("validating collector settings: %w", err)
	}

	if err := validateTags(&config.Tags); err != nil {
		return fmt.Errorf("validating tags: %w", err)
	}

	return nil
}

// validateSNMPSettings validates SNMP configuration settings.
func validateSNMPSettings(settings *SNMPSettings) error {
	if settings == nil {
		return fmt.Errorf("SNMP settings cannot be nil")
	}

	if settings.Host == "" {
		return fmt.Errorf("SNMP host is required")
	}

	if settings.Port < 1 || settings.Port > 65535 {
		return fmt.Errorf("invalid SNMP port: %d", settings.Port)
	}

	if settings.Version == "" {
		return fmt.Errorf("SNMP version is required")
	}

	if settings.Community == "" {
		return fmt.Errorf("SNMP community string is required")
	}

	if !durationRegex.MatchString(settings.Timeout) {
		return fmt.Errorf("invalid timeout format: %s", settings.Timeout)
	}

	return nil
}

// validateCollectorSettings validates collector configuration settings.
func validateCollectorSettings(settings *CollectorSettings) error {
	if settings == nil {
		return fmt.Errorf("collector settings cannot be nil")
	}

	if settings.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	if settings.Version == "" {
		return fmt.Errorf("version is required")
	}

	if !versionRegex.MatchString(settings.Version) {
		return fmt.Errorf("invalid version format: %s", settings.Version)
	}

	if len(settings.Modules) == 0 {
		return fmt.Errorf("at least one module must be specified")
	}

	if !durationRegex.MatchString(settings.CollectionInterval) {
		return fmt.Errorf("invalid collection interval format: %s", settings.CollectionInterval)
	}

	return validateMetrics(&settings.Metrics)
}

// validateMetrics validates metrics configuration settings.
func validateMetrics(metrics *MetricsSettings) error {
	if metrics == nil {
		return fmt.Errorf("metrics settings cannot be nil")
	}

	if len(metrics.Include) == 0 {
		return fmt.Errorf("at least one metric must be included")
	}

	included := make(map[string]bool)
	for _, metric := range metrics.Include {
		if included[metric] {
			return fmt.Errorf("duplicate metric in include list: %s", metric)
		}
		included[metric] = true
	}

	excluded := make(map[string]bool)
	for _, metric := range metrics.Exclude {
		if included[metric] {
			return fmt.Errorf("metric cannot be both included and excluded: %s", metric)
		}
		excluded[metric] = true
	}

	return nil
}

// validateTags validates tag configuration settings.
func validateTags(tags *Tags) error {
	if tags == nil {
		return fmt.Errorf("tags cannot be nil")
	}

	validEnvironments := map[string]bool{
		"development": true,
		"staging":     true,
		"production": true,
	}

	if !validEnvironments[tags.Environment] {
		return fmt.Errorf("invalid environment: %s", tags.Environment)
	}

	if tags.Location == "" {
		return fmt.Errorf("location is required")
	}

	if tags.Role == "" {
		return fmt.Errorf("role is required")
	}

	return nil
}
