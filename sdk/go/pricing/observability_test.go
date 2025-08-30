package pricing_test

import (
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pricing"
)

func TestValidateMetricName(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		expectError bool
	}{
		{"valid metric name", "pulumicost_requests_total", false},
		{"valid with underscore", "http_request_duration_seconds", false},
		{"valid with colon", "http:request_duration_seconds", false},
		{"empty name", "", true},
		{"starts with digit", "1invalid_metric", true},
		{"valid starting with letter", "valid_metric_123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateMetricName(tt.metricName)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for metric name '%s' but got none", tt.metricName)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for metric name '%s': %v", tt.metricName, err)
			}
		})
	}
}

func TestValidateMetricNameStrict(t *testing.T) {
	tests := []struct {
		name        string
		metricName  string
		expectError bool
	}{
		{"valid metric name", "pulumicost_requests_total", false},
		{"empty name", "", true},
		{"reserved prefix __", "__internal_metric", true},
		{"reserved prefix prometheus_", "prometheus_metric", true},
		{"too long", string(make([]byte, 201)), true},
		{"invalid characters", "metric-with-dashes", true},
		{"valid with colon", "http:request_total", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateMetricNameStrict(tt.metricName)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for metric name '%s' but got none", tt.metricName)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for metric name '%s': %v", tt.metricName, err)
			}
		})
	}
}

func TestValidateMetricLabels(t *testing.T) {
	tests := []struct {
		name           string
		labels         map[string]string
		expectValid    bool
		expectWarnings int
	}{
		{
			name:        "valid labels",
			labels:      map[string]string{"method": "GET", "status": "200"},
			expectValid: true,
		},
		{
			name:        "empty label key",
			labels:      map[string]string{"": "value"},
			expectValid: false,
		},
		{
			name:           "empty label value",
			labels:         map[string]string{"key": ""},
			expectValid:    true,
			expectWarnings: 1,
		},
		{
			name:        "reserved label prefix",
			labels:      map[string]string{"__reserved": "value"},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ValidateMetricLabels(tt.labels)
			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got %v", tt.expectValid, result.Valid)
			}
			if len(result.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d", tt.expectWarnings, len(result.Warnings))
			}
		})
	}
}

func TestValidateTraceID(t *testing.T) {
	tests := []struct {
		name        string
		traceID     string
		expectError bool
	}{
		{"valid trace ID", "abcdef1234567890abcdef1234567890", false},
		{"empty trace ID", "", false}, // Optional field
		{"too short", "abcdef123456789", true},
		{"too long", "abcdef1234567890abcdef12345678901", true},
		{"invalid characters", "ghijkl1234567890abcdef1234567890", true},
		{"all zeros", "00000000000000000000000000000000", true},
		{"valid with numbers", "12345678901234567890123456789012", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateTraceID(tt.traceID)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for trace ID '%s' but got none", tt.traceID)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for trace ID '%s': %v", tt.traceID, err)
			}
		})
	}
}

func TestValidateSpanID(t *testing.T) {
	tests := []struct {
		name        string
		spanID      string
		expectError bool
	}{
		{"valid span ID", "abcdef1234567890", false},
		{"empty span ID", "", false}, // Optional field
		{"too short", "abcdef123456", true},
		{"too long", "abcdef12345678901", true},
		{"invalid characters", "ghijkl1234567890", true},
		{"all zeros", "0000000000000000", true},
		{"valid with numbers", "1234567890123456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateSpanID(tt.spanID)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for span ID '%s' but got none", tt.spanID)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for span ID '%s': %v", tt.spanID, err)
			}
		})
	}
}

func TestValidateSLIValue(t *testing.T) {
	tests := []struct {
		name        string
		sli         pricing.ServiceLevelIndicator
		value       float64
		expectError bool
	}{
		{
			name:  "valid percentage",
			sli:   pricing.ServiceLevelIndicator{Name: "availability", Unit: "percentage"},
			value: 99.9,
		},
		{
			name:        "percentage too high",
			sli:         pricing.ServiceLevelIndicator{Name: "availability", Unit: "percentage"},
			value:       101.0,
			expectError: true,
		},
		{
			name:        "percentage negative",
			sli:         pricing.ServiceLevelIndicator{Name: "availability", Unit: "percentage"},
			value:       -1.0,
			expectError: true,
		},
		{
			name:  "valid seconds",
			sli:   pricing.ServiceLevelIndicator{Name: "latency", Unit: "seconds"},
			value: 1.5,
		},
		{
			name:        "negative seconds",
			sli:         pricing.ServiceLevelIndicator{Name: "latency", Unit: "seconds"},
			value:       -1.0,
			expectError: true,
		},
		{
			name:  "valid ratio",
			sli:   pricing.ServiceLevelIndicator{Name: "error_rate", Unit: "ratio"},
			value: 0.001,
		},
		{
			name:        "ratio too high",
			sli:         pricing.ServiceLevelIndicator{Name: "error_rate", Unit: "ratio"},
			value:       1.5,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateSLIValue(tt.sli, tt.value)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for SLI value %f but got none", tt.value)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for SLI value %f: %v", tt.value, err)
			}
		})
	}
}

func TestCalculateErrorRate(t *testing.T) {
	tests := []struct {
		name       string
		total      int64
		successful int64
		expected   float64
	}{
		{"no requests", 0, 0, 0.0},
		{"all successful", 100, 100, 0.0},
		{"some errors", 100, 95, 5.0},
		{"all errors", 100, 0, 100.0},
		{"single error", 10, 9, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.CalculateErrorRate(tt.total, tt.successful)
			if result != tt.expected {
				t.Errorf("Expected error rate %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestCalculateAvailability(t *testing.T) {
	tests := []struct {
		name     string
		uptime   int64
		total    int64
		expected float64
	}{
		{"no time", 0, 0, 0.0},
		{"full uptime", 3600, 3600, 100.0},
		{"partial uptime", 3540, 3600, 98.333333333333333},
		{"no uptime", 0, 3600, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.CalculateAvailability(tt.uptime, tt.total)
			if result != tt.expected {
				t.Errorf("Expected availability %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"nanoseconds", 500 * time.Nanosecond, "500ns"},
		{"microseconds", 500 * time.Microsecond, "500.0μs"},
		{"milliseconds", 500 * time.Millisecond, "500.0ms"},
		{"seconds", 2 * time.Second, "2.00s"},
		{"minutes", 5 * time.Minute, "5.0m"},
		{"hours", 2 * time.Hour, "2.0h"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("Expected formatted duration '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetObservabilityRequirements(t *testing.T) {
	tests := []struct {
		name          string
		level         pricing.ConformanceLevel
		expectedLevel pricing.ConformanceLevel
		minMetrics    int
		minSLIs       int
		minFeatures   int
	}{
		{
			name:          "basic conformance",
			level:         pricing.ConformanceBasic,
			expectedLevel: pricing.ConformanceBasic,
			minMetrics:    3,
			minSLIs:       2,
			minFeatures:   2,
		},
		{
			name:          "standard conformance",
			level:         pricing.ConformanceStandard,
			expectedLevel: pricing.ConformanceStandard,
			minMetrics:    6,
			minSLIs:       4,
			minFeatures:   4,
		},
		{
			name:          "advanced conformance",
			level:         pricing.ConformanceAdvanced,
			expectedLevel: pricing.ConformanceAdvanced,
			minMetrics:    10,
			minSLIs:       6,
			minFeatures:   7,
		},
		{
			name:          "invalid level defaults to basic",
			level:         "invalid",
			expectedLevel: pricing.ConformanceBasic,
			minMetrics:    3,
			minSLIs:       2,
			minFeatures:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := pricing.GetObservabilityRequirements(tt.level)

			if req.Level != tt.expectedLevel {
				t.Errorf("Expected level %s, got %s", tt.expectedLevel, req.Level)
			}

			if len(req.RequiredMetrics) < tt.minMetrics {
				t.Errorf("Expected at least %d metrics, got %d", tt.minMetrics, len(req.RequiredMetrics))
			}

			if len(req.RequiredSLIs) < tt.minSLIs {
				t.Errorf("Expected at least %d SLIs, got %d", tt.minSLIs, len(req.RequiredSLIs))
			}

			if len(req.RequiredFeatures) < tt.minFeatures {
				t.Errorf("Expected at least %d features, got %d", tt.minFeatures, len(req.RequiredFeatures))
			}
		})
	}
}

func TestValidateTimeRange(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		start       time.Time
		end         time.Time
		expectError bool
	}{
		{
			name:  "valid range",
			start: now.Add(-time.Hour),
			end:   now,
		},
		{
			name:        "zero start time",
			start:       time.Time{},
			end:         now,
			expectError: true,
		},
		{
			name:        "zero end time",
			start:       now.Add(-time.Hour),
			end:         time.Time{},
			expectError: true,
		},
		{
			name:        "end before start",
			start:       now,
			end:         now.Add(-time.Hour),
			expectError: true,
		},
		{
			name:        "too far in past",
			start:       now.Add(-400 * 24 * time.Hour), // 400 days ago
			end:         now,
			expectError: true,
		},
		{
			name:        "too far in future",
			start:       now,
			end:         now.Add(48 * time.Hour), // 2 days in future
			expectError: true,
		},
		{
			name:        "range too long",
			start:       now.Add(-100 * 24 * time.Hour), // 100 days ago
			end:         now,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateTimeRange(tt.start, tt.end)
			if tt.expectError && err == nil {
				t.Errorf("Expected error for time range but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for time range: %v", err)
			}
		})
	}
}

func TestValidateObservabilityMetadata(t *testing.T) {
	tests := []struct {
		name             string
		traceID          string
		spanID           string
		requestID        string
		processingTimeMs int64
		qualityScore     float64
		expectValid      bool
		expectWarnings   int
	}{
		{
			name:             "valid metadata",
			traceID:          "abcdef1234567890abcdef1234567890",
			spanID:           "abcdef1234567890",
			requestID:        "req-123",
			processingTimeMs: 100,
			qualityScore:     0.95,
			expectValid:      true,
		},
		{
			name:             "invalid trace ID",
			traceID:          "invalid",
			spanID:           "",
			requestID:        "req-123",
			processingTimeMs: 100,
			qualityScore:     -1, // negative means not set
			expectValid:      false,
		},
		{
			name:             "span without trace",
			traceID:          "",
			spanID:           "abcdef1234567890",
			requestID:        "req-123",
			processingTimeMs: 100,
			qualityScore:     -1,
			expectValid:      false,
		},
		{
			name:             "empty request ID",
			traceID:          "",
			spanID:           "",
			requestID:        "",
			processingTimeMs: 100,
			qualityScore:     -1,
			expectValid:      true,
			expectWarnings:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := pricing.ValidateObservabilityMetadata(tt.traceID, tt.spanID, tt.requestID,
				tt.processingTimeMs, tt.qualityScore)

			if suite.IsValid() != tt.expectValid {
				t.Errorf("Expected valid=%v, got %v. Errors: %v", tt.expectValid, suite.IsValid(), suite.Errors)
			}

			if len(suite.Warnings) != tt.expectWarnings {
				t.Errorf("Expected %d warnings, got %d: %v", tt.expectWarnings, len(suite.Warnings), suite.Warnings)
			}
		})
	}
}

func TestStandardMetricsConstants(t *testing.T) {
	// Verify all standard metrics are properly defined
	expectedMetrics := []string{
		pricing.StandardMetrics.RequestsTotal,
		pricing.StandardMetrics.RequestDurationSeconds,
		pricing.StandardMetrics.ErrorsTotal,
		pricing.StandardMetrics.LatencyP95Seconds,
		pricing.StandardMetrics.LatencyP99Seconds,
		pricing.StandardMetrics.CacheHitRatePercent,
	}

	for _, metric := range expectedMetrics {
		if metric == "" {
			t.Errorf("Standard metric is empty")
		}

		if err := pricing.ValidateMetricNameStrict(metric); err != nil {
			t.Errorf("Standard metric '%s' fails validation: %v", metric, err)
		}
	}
}

func TestStandardSLIsTargets(t *testing.T) {
	// Verify all standard SLIs have reasonable targets
	slis := []pricing.ServiceLevelIndicator{
		pricing.StandardSLIs.Availability,
		pricing.StandardSLIs.ErrorRate,
		pricing.StandardSLIs.LatencyP99,
		pricing.StandardSLIs.LatencyP95,
		pricing.StandardSLIs.Throughput,
		pricing.StandardSLIs.DataFreshness,
	}

	for _, sli := range slis {
		if sli.Name == "" {
			t.Errorf("SLI name is empty")
		}

		if sli.Target <= 0 {
			t.Errorf("SLI '%s' has non-positive target: %f", sli.Name, sli.Target)
		}

		if err := pricing.ValidateSLIValue(sli, sli.Target); err != nil {
			t.Errorf("SLI '%s' target value fails validation: %v", sli.Name, err)
		}
	}
}
