// Package pluginsdk provides environment variable handling for FinFocus plugins.
package pluginsdk

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// Environment variable constants for FinFocus plugins.
// These constants define the canonical names for all supported environment variables.
const (
	// EnvPort is the canonical environment variable for plugin gRPC port.
	// Fallback chain: FINFOCUS_PLUGIN_PORT -> PULUMICOST_PLUGIN_PORT.
	EnvPort         = "FINFOCUS_PLUGIN_PORT"
	EnvPortFallback = "PULUMICOST_PLUGIN_PORT"

	// EnvLogLevel is the canonical environment variable for log verbosity.
	// Supported values: debug, info, warn, error (plugin-specific validation).
	// Fallback chain: FINFOCUS_LOG_LEVEL -> PULUMICOST_LOG_LEVEL -> LOG_LEVEL.
	EnvLogLevel           = "FINFOCUS_LOG_LEVEL"
	EnvLogLevelPulumiCost = "PULUMICOST_LOG_LEVEL"
	EnvLogLevelFallback   = "LOG_LEVEL"

	// EnvLogFormat is the environment variable for log output format.
	// Supported values: json, text (plugin-specific validation).
	// Fallback chain: FINFOCUS_LOG_FORMAT -> PULUMICOST_LOG_FORMAT.
	EnvLogFormat         = "FINFOCUS_LOG_FORMAT"
	EnvLogFormatFallback = "PULUMICOST_LOG_FORMAT"

	// EnvLogFile is the environment variable for log file path.
	// Empty string means log to stderr (not stdout).
	// Fallback chain: FINFOCUS_LOG_FILE -> PULUMICOST_LOG_FILE.
	EnvLogFile         = "FINFOCUS_LOG_FILE"
	EnvLogFileFallback = "PULUMICOST_LOG_FILE"

	// EnvTraceID is the environment variable for distributed tracing.
	// When set, this ID should be included in all related logs and responses.
	// Fallback chain: FINFOCUS_TRACE_ID -> PULUMICOST_TRACE_ID.
	EnvTraceID         = "FINFOCUS_TRACE_ID"
	EnvTraceIDFallback = "PULUMICOST_TRACE_ID"

	// EnvTestMode is the environment variable for enabling test mode.
	// Only "true" enables test mode; all other values disable it.
	// Fallback chain: FINFOCUS_TEST_MODE -> PULUMICOST_TEST_MODE.
	EnvTestMode         = "FINFOCUS_TEST_MODE"
	EnvTestModeFallback = "PULUMICOST_TEST_MODE"
)

// GetPort returns the plugin port from environment variables.
// Checks FINFOCUS_PLUGIN_PORT first, then falls back to PULUMICOST_PLUGIN_PORT.
// Returns 0 if not set or invalid.
// Callers should treat 0 as "port not configured" and handle accordingly.
func GetPort() int {
	v := os.Getenv(EnvPort)
	if v == "" {
		v = os.Getenv(EnvPortFallback)
	}
	if v == "" {
		return 0
	}
	port, err := strconv.Atoi(v)
	if err != nil || port <= 0 {
		return 0
	}
	return port
}

// GetLogLevel returns the log level from environment variables.
// Checks in order: FINFOCUS_LOG_LEVEL -> PULUMICOST_LOG_LEVEL -> LOG_LEVEL.
// Returns empty string if none are set.
func GetLogLevel() string {
	if v := os.Getenv(EnvLogLevel); v != "" {
		return v
	}
	if v := os.Getenv(EnvLogLevelPulumiCost); v != "" {
		return v
	}
	return os.Getenv(EnvLogLevelFallback)
}

// GetLogFormat returns the log format from environment variables.
// Checks FINFOCUS_LOG_FORMAT first, then falls back to PULUMICOST_LOG_FORMAT.
// Returns empty string if not set.
func GetLogFormat() string {
	if v := os.Getenv(EnvLogFormat); v != "" {
		return v
	}
	return os.Getenv(EnvLogFormatFallback)
}

// GetLogFile returns the log file path from environment variables.
// Checks FINFOCUS_LOG_FILE first, then falls back to PULUMICOST_LOG_FILE.
// Returns an empty string if unset or empty, indicating that logging will use os.Stderr.
func GetLogFile() string {
	if v := os.Getenv(EnvLogFile); v != "" {
		return v
	}
	return os.Getenv(EnvLogFileFallback)
}

// GetTraceID returns the trace ID from environment variables.
// Checks FINFOCUS_TRACE_ID first, then falls back to PULUMICOST_TRACE_ID.
// Returns empty string if not set.
func GetTraceID() string {
	if v := os.Getenv(EnvTraceID); v != "" {
		return v
	}
	return os.Getenv(EnvTraceIDFallback)
}

// GetTestMode returns true if test mode is enabled via environment variables.
// Checks FINFOCUS_TEST_MODE first, then falls back to PULUMICOST_TEST_MODE.
// Logs a warning if the value is set but not "true" or "false".
// Returns false for any value other than "true".
func GetTestMode() bool {
	v := os.Getenv(EnvTestMode)
	if v == "" {
		v = os.Getenv(EnvTestModeFallback)
	}
	if v == "" {
		return false
	}
	if v == "true" {
		return true
	}
	if v != "false" {
		log.Warn().
			Str("variable", EnvTestMode).
			Str("value", v).
			Msg("invalid test mode value, expected 'true' or 'false', defaulting to disabled")
	}
	return false
}

// IsTestMode returns true if test mode is enabled via environment variables.
// Checks FINFOCUS_TEST_MODE first, then falls back to PULUMICOST_TEST_MODE.
// Unlike GetTestMode, this function does not log warnings for invalid values.
// Use this for repeated checks to avoid log spam.
func IsTestMode() bool {
	v := os.Getenv(EnvTestMode)
	if v == "" {
		v = os.Getenv(EnvTestModeFallback)
	}
	return v == "true"
}
