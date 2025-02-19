package config

// Logging levels
const (
	LoggingLevelDebug = "debug"
	LoggingLevelInfo  = "info"
	LoggingLevelWarn  = "warn"
	LoggingLevelError = "error"
)

// Default values for application settings
const (
	DefaultConfigurationRefreshMinutes = 5
	DefaultMaximumDataCollectors      = 5
	DefaultMaximumDataWriters         = 3
	DefaultMaximumRecoveryAttempts    = 3
	DefaultRecoveryWaitSeconds        = 5
)

// Minimum values for configuration validation
const (
	MinConfigRefreshMinutes    = 1
	MinConcurrentOperations    = 1
	MinRetryAttempts          = 1
	MinRetryDelaySeconds      = 1
)

// Maximum values for configuration validation
const (
	MaxRetryDelaySeconds      = 300
	MaxConcurrentOperations   = 100
)
