package pluginsdk_test

import (
	"os"
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

// setEnv sets an environment variable and restores it on cleanup.
func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

// unsetEnv clears an environment variable and restores it on cleanup.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	original, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		if existed {
			t.Setenv(key, original)
		}
	})
}

// ============================================================================
// User Story 1: Port Configuration Tests
// ============================================================================

func TestGetPort_Set(t *testing.T) {
	setEnv(t, pluginsdk.EnvPort, "8080")
	got := pluginsdk.GetPort()
	if got != 8080 {
		t.Errorf("GetPort() = %d, want 8080", got)
	}
}

func TestGetPort_NotSet_ReturnsZero(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvPort)
	got := pluginsdk.GetPort()
	if got != 0 {
		t.Errorf("GetPort() = %d, want 0 when not set", got)
	}
}

func TestGetPort_InvalidValue_ReturnsZero(t *testing.T) {
	setEnv(t, pluginsdk.EnvPort, "abc")
	got := pluginsdk.GetPort()
	if got != 0 {
		t.Errorf("GetPort() = %d, want 0 for invalid value", got)
	}
}

func TestGetPort_NonPositive_ReturnsZero(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"zero", "0"},
		{"negative", "-1"},
		{"negative large", "-8080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(t, pluginsdk.EnvPort, tt.value)
			got := pluginsdk.GetPort()
			if got != 0 {
				t.Errorf("GetPort() = %d, want 0 for %s", got, tt.name)
			}
		})
	}
}

func TestGetPort_Fallback(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvPort)
	setEnv(t, pluginsdk.EnvPortFallback, "9090")
	got := pluginsdk.GetPort()
	if got != 9090 {
		t.Errorf("GetPort() = %d, want 9090 (fallback should work)", got)
	}
}

// ============================================================================
// User Story 2: Logging Configuration Tests
// ============================================================================

func TestGetLogLevel_CanonicalVariable(t *testing.T) {
	setEnv(t, pluginsdk.EnvLogLevel, "debug")
	unsetEnv(t, pluginsdk.EnvLogLevelFallback)
	got := pluginsdk.GetLogLevel()
	if got != "debug" {
		t.Errorf("GetLogLevel() = %q, want %q", got, "debug")
	}
}

func TestGetLogLevel_FallbackVariable(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvLogLevel)
	setEnv(t, pluginsdk.EnvLogLevelFallback, "info")
	got := pluginsdk.GetLogLevel()
	if got != "info" {
		t.Errorf("GetLogLevel() = %q, want %q", got, "info")
	}
}

func TestGetLogLevel_CanonicalTakesPrecedence(t *testing.T) {
	setEnv(t, pluginsdk.EnvLogLevel, "debug")
	setEnv(t, pluginsdk.EnvLogLevelFallback, "info")
	got := pluginsdk.GetLogLevel()
	if got != "debug" {
		t.Errorf("GetLogLevel() = %q, want %q (canonical should take precedence)", got, "debug")
	}
}

func TestGetLogLevel_NeitherSet_ReturnsEmpty(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvLogLevel)
	unsetEnv(t, pluginsdk.EnvLogLevelFallback)
	got := pluginsdk.GetLogLevel()
	if got != "" {
		t.Errorf("GetLogLevel() = %q, want empty string", got)
	}
}

func TestGetLogFormat(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"json format", "json", "json"},
		{"text format", "text", "text"},
		{"not set", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				setEnv(t, pluginsdk.EnvLogFormat, tt.value)
			} else {
				unsetEnv(t, pluginsdk.EnvLogFormat)
			}
			got := pluginsdk.GetLogFormat()
			if got != tt.expected {
				t.Errorf("GetLogFormat() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetLogFile(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"file path", "/var/log/plugin.log", "/var/log/plugin.log"},
		{"not set", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				setEnv(t, pluginsdk.EnvLogFile, tt.value)
			} else {
				unsetEnv(t, pluginsdk.EnvLogFile)
			}
			got := pluginsdk.GetLogFile()
			if got != tt.expected {
				t.Errorf("GetLogFile() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// User Story 3: Trace ID Configuration Tests
// ============================================================================

func TestGetTraceID_Set(t *testing.T) {
	setEnv(t, pluginsdk.EnvTraceID, "abc123")
	got := pluginsdk.GetTraceID()
	if got != "abc123" {
		t.Errorf("GetTraceID() = %q, want %q", got, "abc123")
	}
}

func TestGetTraceID_NotSet_ReturnsEmpty(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvTraceID)
	got := pluginsdk.GetTraceID()
	if got != "" {
		t.Errorf("GetTraceID() = %q, want empty string", got)
	}
}

// ============================================================================
// User Story 4: Test Mode Configuration Tests
// ============================================================================

func TestGetTestMode_True(t *testing.T) {
	setEnv(t, pluginsdk.EnvTestMode, "true")
	got := pluginsdk.GetTestMode()
	if !got {
		t.Error("GetTestMode() = false, want true")
	}
}

func TestGetTestMode_False(t *testing.T) {
	setEnv(t, pluginsdk.EnvTestMode, "false")
	got := pluginsdk.GetTestMode()
	if got {
		t.Error("GetTestMode() = true, want false")
	}
}

func TestGetTestMode_NotSet_ReturnsFalse(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvTestMode)
	got := pluginsdk.GetTestMode()
	if got {
		t.Error("GetTestMode() = true, want false when not set")
	}
}

func TestGetTestMode_InvalidValue_ReturnsFalse(t *testing.T) {
	tests := []string{"yes", "no", "1", "0", "TRUE", "FALSE", "True"}
	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			setEnv(t, pluginsdk.EnvTestMode, value)
			got := pluginsdk.GetTestMode()
			if got {
				t.Errorf("GetTestMode() = true for %q, want false", value)
			}
		})
	}
}

func TestIsTestMode_NoWarning(t *testing.T) {
	// IsTestMode should return the same result as GetTestMode but without logging.
	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"", false},
		{"yes", false},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if tt.value != "" {
				setEnv(t, pluginsdk.EnvTestMode, tt.value)
			} else {
				unsetEnv(t, pluginsdk.EnvTestMode)
			}
			got := pluginsdk.IsTestMode()
			if got != tt.expected {
				t.Errorf("IsTestMode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============================================================================
// User Story 5: Migration Support Tests
// ============================================================================

func TestGetLogLevel_FallbackWorks(t *testing.T) {
	// This test verifies that LOG_LEVEL fallback works for migration support.
	unsetEnv(t, pluginsdk.EnvLogLevel)
	unsetEnv(t, pluginsdk.EnvLogLevelPulumiCost)
	setEnv(t, pluginsdk.EnvLogLevelFallback, "warn")
	got := pluginsdk.GetLogLevel()
	if got != "warn" {
		t.Errorf("GetLogLevel() = %q, want %q (fallback should work)", got, "warn")
	}
}

func TestGetLogLevel_PulumiCostFallback(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvLogLevel)
	setEnv(t, pluginsdk.EnvLogLevelPulumiCost, "error")
	setEnv(t, pluginsdk.EnvLogLevelFallback, "warn")
	got := pluginsdk.GetLogLevel()
	if got != "error" {
		t.Errorf(
			"GetLogLevel() = %q, want %q (PULUMICOST_LOG_LEVEL should take precedence over LOG_LEVEL)",
			got, "error",
		)
	}
}

func TestGetLogFormat_Fallback(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvLogFormat)
	setEnv(t, pluginsdk.EnvLogFormatFallback, "text")
	got := pluginsdk.GetLogFormat()
	if got != "text" {
		t.Errorf("GetLogFormat() = %q, want %q (fallback should work)", got, "text")
	}
}

func TestGetLogFile_Fallback(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvLogFile)
	setEnv(t, pluginsdk.EnvLogFileFallback, "/tmp/legacy.log")
	got := pluginsdk.GetLogFile()
	if got != "/tmp/legacy.log" {
		t.Errorf("GetLogFile() = %q, want %q (fallback should work)", got, "/tmp/legacy.log")
	}
}

func TestGetTraceID_Fallback(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvTraceID)
	setEnv(t, pluginsdk.EnvTraceIDFallback, "legacy-trace-123")
	got := pluginsdk.GetTraceID()
	if got != "legacy-trace-123" {
		t.Errorf("GetTraceID() = %q, want %q (fallback should work)", got, "legacy-trace-123")
	}
}

func TestGetTestMode_Fallback(t *testing.T) {
	unsetEnv(t, pluginsdk.EnvTestMode)
	setEnv(t, pluginsdk.EnvTestModeFallback, "true")
	got := pluginsdk.GetTestMode()
	if !got {
		t.Error("GetTestMode() = false, want true (fallback should work)")
	}
}
