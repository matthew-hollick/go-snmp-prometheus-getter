package elasticsearch

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	
	// Regular expressions for validation
	hostnameRegex    = regexp.MustCompile(`^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.hedgehog\.internal$`)
	versionRegex     = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	durationRegex    = regexp.MustCompile(`^[0-9]+(ms|s|m|h)$`)
)

func init() {
	validate = validator.New()

	// Register custom validation functions
	if err := validate.RegisterValidation("hostname", validateHostname); err != nil {
		panic(fmt.Sprintf("failed to register hostname validator: %v", err))
	}
	if err := validate.RegisterValidation("duration", validateDuration); err != nil {
		panic(fmt.Sprintf("failed to register duration validator: %v", err))
	}
	if err := validate.RegisterValidation("version", validateVersion); err != nil {
		panic(fmt.Sprintf("failed to register version validator: %v", err))
	}
}

// validateHostname checks if a hostname matches our required pattern
func validateHostname(fl validator.FieldLevel) bool {
	return hostnameRegex.MatchString(fl.Field().String())
}

// validateDuration checks if a duration string is valid
func validateDuration(fl validator.FieldLevel) bool {
	if !durationRegex.MatchString(fl.Field().String()) {
		return false
	}
	
	// Try parsing the duration to ensure it's valid
	_, err := time.ParseDuration(fl.Field().String())
	return err == nil
}

// validateVersion checks if a version string matches semantic versioning
func validateVersion(fl validator.FieldLevel) bool {
	return versionRegex.MatchString(fl.Field().String())
}

// validateSNMPSettings performs additional validation on SNMP settings
func validateSNMPSettings(settings SNMPSettings) error {
	if settings.Port < 1 || settings.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if !strings.Contains(settings.Host, ".hedgehog.internal") {
		return fmt.Errorf("host must be in the hedgehog.internal domain")
	}

	validVersions := map[string]bool{"1": true, "2c": true, "3": true}
	if !validVersions[settings.Version] {
		return fmt.Errorf("invalid SNMP version: %s", settings.Version)
	}

	return nil
}

// validateCollectorSettings performs additional validation on collector settings
func validateCollectorSettings(settings CollectorSettings) error {
	if !strings.Contains(settings.Hostname, ".hedgehog.internal") {
		return fmt.Errorf("hostname must be in the hedgehog.internal domain")
	}

	if len(settings.Modules) == 0 {
		return fmt.Errorf("at least one module must be specified")
	}

	if len(settings.Metrics.Include) == 0 {
		return fmt.Errorf("at least one metric must be included")
	}

	// Check for duplicate metrics
	included := make(map[string]bool)
	for _, metric := range settings.Metrics.Include {
		if included[metric] {
			return fmt.Errorf("duplicate metric in include list: %s", metric)
		}
		included[metric] = true
	}

	excluded := make(map[string]bool)
	for _, metric := range settings.Metrics.Exclude {
		if excluded[metric] {
			return fmt.Errorf("duplicate metric in exclude list: %s", metric)
		}
		if included[metric] {
			return fmt.Errorf("metric cannot be both included and excluded: %s", metric)
		}
		excluded[metric] = true
	}

	return nil
}

// validateTags performs additional validation on tags
func validateTags(tags Tags) error {
	validEnvironments := map[string]bool{
		"development": true,
		"staging":     true,
		"production": true,
	}
	if !validEnvironments[tags.Environment] {
		return fmt.Errorf("invalid environment: %s", tags.Environment)
	}

	validRoles := map[string]bool{
		"network-switch": true,
		"router":        true,
		"firewall":      true,
		"server":        true,
	}
	if !validRoles[tags.Role] {
		return fmt.Errorf("invalid role: %s", tags.Role)
	}

	if tags.Location == "" {
		return fmt.Errorf("location is required")
	}

	return nil
}
