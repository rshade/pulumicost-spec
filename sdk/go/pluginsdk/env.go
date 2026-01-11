// Package pluginsdk provides environment variable handling for PulumiCost plugins.
package pluginsdk

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

// Environment variable constants for PulumiCost plugins.
// These constants define the canonical names for all supported environment variables.
const (
	// EnvPort is the canonical environment variable for plugin gRPC port.
	// Plugins MUST use PULUMICOST_PLUGIN_PORT (no fallback to PORT).
	EnvPort = "PULUMICOST_PLUGIN_PORT"

	// EnvLogLevel is the canonical environment variable for log verbosity.
	// Supported values: debug, info, warn, error (plugin-specific validation).
	EnvLogLevel = "PULUMICOST_LOG_LEVEL"

	// EnvLogLevelFallback is the legacy environment variable for log verbosity.
	// GetLogLevel() checks PULUMICOST_LOG_LEVEL first, then falls back to LOG_LEVEL.
	EnvLogLevelFallback = "LOG_LEVEL"

	// EnvLogFormat is the environment variable for log output format.
	// Supported values: json, text (plugin-specific validation).
	EnvLogFormat = "PULUMICOST_LOG_FORMAT"

	// EnvLogFile is the environment variable for log file path.
	// Empty string means log to stderr.
	EnvLogFile = "PULUMICOST_LOG_FILE"

	// EnvTraceID is the environment variable for distributed tracing.
	// When set, this ID should be included in all related logs and responses.
	EnvTraceID = "PULUMICOST_TRACE_ID"

	// EnvTestMode is the environment variable for enabling test mode.
	// Only "true" enables test mode; all other values disable it.
	EnvTestMode = "PULUMICOST_TEST_MODE"
)

// GetPort returns the plugin port from PULUMICOST_PLUGIN_PORT environment variable.
// Returns 0 if not set or invalid. There is no fallback to PORT.
// Callers should treat 0 as "port not configured" and handle accordingly.
func GetPort() int {
	v := os.Getenv(EnvPort)
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
// Checks PULUMICOST_LOG_LEVEL first, then falls back to LOG_LEVEL.
// Returns empty string if neither is set.
func GetLogLevel() string {
	if v := os.Getenv(EnvLogLevel); v != "" {
		return v
	}
	return os.Getenv(EnvLogLevelFallback)
}

// GetLogFormat returns the log format from PULUMICOST_LOG_FORMAT.
// Returns empty string if not set.
func GetLogFormat() string {
	return os.Getenv(EnvLogFormat)
}

// GetLogFile returns the log file path from PULUMICOST_LOG_FILE.
// Returns empty string if not set (meaning stdout).
func GetLogFile() string {
	return os.Getenv(EnvLogFile)
}

// GetTraceID returns the trace ID from PULUMICOST_TRACE_ID.
// Returns empty string if not set.
func GetTraceID() string {
	return os.Getenv(EnvTraceID)
}

// GetTestMode returns true if PULUMICOST_TEST_MODE is set to "true".
// Logs a warning if the value is set but not "true" or "false".
// Returns false for any value other than "true".
func GetTestMode() bool {
	v := os.Getenv(EnvTestMode)
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

// IsTestMode returns true if PULUMICOST_TEST_MODE is set to "true".
// Unlike GetTestMode, this function does not log warnings for invalid values.
// Use this for repeated checks to avoid log spam.
func IsTestMode() bool {
	return os.Getenv(EnvTestMode) == "true"
}
