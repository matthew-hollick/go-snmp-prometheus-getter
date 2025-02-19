package config

// Logging levels.
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Default values for application settings.
const (
	DefaultMaximumDataCollectors   = 5
	DefaultMaximumDataWriters      = 3
	DefaultMaximumRetryAttempts    = 3
	DefaultRetryWaitSeconds        = 5
)

// Minimum values for configuration validation.
const (
	MinimumDataCollectors   = 1
	MinimumDataWriters      = 1
	MinimumRetryAttempts    = 1
	MinimumRetryWaitSeconds = 1
)

// Maximum values for configuration validation.
const (
	MaximumDataCollectors   = 100
	MaximumDataWriters      = 50
	MaximumRetryAttempts    = 10
	MaximumRetryWaitSeconds = 60
)
