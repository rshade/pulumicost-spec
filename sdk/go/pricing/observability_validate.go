package pricing

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

var (
	// metricNameRegex validates metric names according to Prometheus conventions
	metricNameRegex = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)

	// labelNameRegex validates label names
	labelNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

	// traceIDRegex validates OpenTelemetry trace IDs (32 hex characters)
	traceIDRegex = regexp.MustCompile(`^[0-9a-f]{32}$`)

	// spanIDRegex validates OpenTelemetry span IDs (16 hex characters)
	spanIDRegex = regexp.MustCompile(`^[0-9a-f]{16}$`)
)

// ValidateMetricNameStrict performs comprehensive metric name validation.
func ValidateMetricNameStrict(name string) error {
	if name == "" {
		return fmt.Errorf("metric name cannot be empty")
	}

	if !metricNameRegex.MatchString(name) {
		return fmt.Errorf("metric name '%s' must match pattern %s", name, metricNameRegex.String())
	}

	if len(name) > 200 {
		return fmt.Errorf("metric name '%s' exceeds maximum length of 200 characters", name)
	}

	// Check for reserved prefixes
	reservedPrefixes := []string{"__", "prometheus_", "go_", "process_"}
	for _, prefix := range reservedPrefixes {
		if strings.HasPrefix(name, prefix) {
			return fmt.Errorf("metric name '%s' uses reserved prefix '%s'", name, prefix)
		}
	}

	return nil
}

// ValidateLabelNameStrict performs comprehensive label name validation.
func ValidateLabelNameStrict(name string) error {
	if name == "" {
		return fmt.Errorf("label name cannot be empty")
	}

	if !labelNameRegex.MatchString(name) {
		return fmt.Errorf("label name '%s' must match pattern %s", name, labelNameRegex.String())
	}

	if len(name) > 100 {
		return fmt.Errorf("label name '%s' exceeds maximum length of 100 characters", name)
	}

	// Check for reserved names
	reservedNames := []string{"__name__", "__value__", "job", "instance"}
	for _, reserved := range reservedNames {
		if name == reserved {
			return fmt.Errorf("label name '%s' is reserved", name)
		}
	}

	return nil
}

// ValidateTraceID validates OpenTelemetry trace identifiers.
func ValidateTraceID(traceID string) error {
	if traceID == "" {
		return nil // trace ID is optional
	}

	if !traceIDRegex.MatchString(traceID) {
		return fmt.Errorf("trace ID '%s' must be 32 hexadecimal characters", traceID)
	}

	// Check for all-zero trace ID (invalid)
	if traceID == "00000000000000000000000000000000" {
		return fmt.Errorf("trace ID cannot be all zeros")
	}

	return nil
}

// ValidateSpanID validates OpenTelemetry span identifiers.
func ValidateSpanID(spanID string) error {
	if spanID == "" {
		return nil // span ID is optional
	}

	if !spanIDRegex.MatchString(spanID) {
		return fmt.Errorf("span ID '%s' must be 16 hexadecimal characters", spanID)
	}

	// Check for all-zero span ID (invalid)
	if spanID == "0000000000000000" {
		return fmt.Errorf("span ID cannot be all zeros")
	}

	return nil
}

// ValidateLogLevel checks if a log level is valid.
func ValidateLogLevel(level string) error {
	validLevels := []string{
		string(LogLevelDebug),
		string(LogLevelInfo),
		string(LogLevelWarn),
		string(LogLevelError),
		string(LogLevelFatal),
	}

	for _, valid := range validLevels {
		if level == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid log level '%s', must be one of: %v", level, validLevels)
}

// ValidateHealthStatus checks if a health status is valid.
func ValidateHealthStatus(status string) error {
	validStatuses := []string{
		string(HealthStatusUnknown),
		string(HealthStatusServing),
		string(HealthStatusNotServing),
		string(HealthStatusServiceUnknown),
	}

	for _, valid := range validStatuses {
		if status == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid health status '%s', must be one of: %v", status, validStatuses)
}

// ValidateMetricValue checks if a metric value is valid.
func ValidateMetricValue(value float64, metricType MetricType) error {
	// Check for invalid float values
	if value != value { // NaN check
		return fmt.Errorf("metric value cannot be NaN")
	}

	// Counter values must be non-negative and finite
	if metricType == MetricTypeCounter {
		if value < 0 {
			return fmt.Errorf("counter metric value cannot be negative: %f", value)
		}
		if math.IsInf(value, 0) { // Infinity check
			return fmt.Errorf("counter metric value cannot be infinite: %f", value)
		}
	}

	return nil
}

// ValidateTimeRange ensures a time range is valid for metrics collection.
func ValidateTimeRange(start, end time.Time) error {
	if start.IsZero() || end.IsZero() {
		return fmt.Errorf("start and end times cannot be zero")
	}

	if end.Before(start) {
		return fmt.Errorf("end time cannot be before start time")
	}

	// Check for reasonable time ranges (not too far in the past or future)
	now := time.Now()
	maxPastDuration := 365 * 24 * time.Hour // 1 year
	maxFutureDuration := 24 * time.Hour     // 1 day

	if start.Before(now.Add(-maxPastDuration)) {
		return fmt.Errorf("start time cannot be more than %v in the past", maxPastDuration)
	}

	if end.After(now.Add(maxFutureDuration)) {
		return fmt.Errorf("end time cannot be more than %v in the future", maxFutureDuration)
	}

	// Check for excessively long time ranges
	duration := end.Sub(start)
	maxDuration := 90 * 24 * time.Hour // 90 days
	if duration > maxDuration {
		return fmt.Errorf("time range duration (%v) exceeds maximum allowed (%v)", duration, maxDuration)
	}

	return nil
}

// ValidateErrorCode checks if an error code follows conventions.
func ValidateErrorCode(code string) error {
	if code == "" {
		return fmt.Errorf("error code cannot be empty")
	}

	// Error codes should be upper case with underscores
	if strings.ToUpper(code) != code {
		return fmt.Errorf("error code '%s' should be uppercase", code)
	}

	// Basic pattern validation
	validPattern := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	if !validPattern.MatchString(code) {
		return fmt.Errorf("error code '%s' must start with letter and contain only uppercase letters, numbers, and underscores", code)
	}

	if len(code) > 50 {
		return fmt.Errorf("error code '%s' exceeds maximum length of 50 characters", code)
	}

	return nil
}

// ValidateQualityScore ensures quality scores are in valid range.
func ValidateQualityScore(score float64) error {
	if score < 0.0 || score > 1.0 {
		return fmt.Errorf("quality score must be between 0.0 and 1.0, got %f", score)
	}

	if score != score { // NaN check
		return fmt.Errorf("quality score cannot be NaN")
	}

	return nil
}

// ValidateProcessingTime checks if processing time is reasonable.
func ValidateProcessingTime(processingTimeMs int64) error {
	if processingTimeMs < 0 {
		return fmt.Errorf("processing time cannot be negative: %d", processingTimeMs)
	}

	// Flag unreasonably long processing times (>10 minutes)
	maxProcessingTimeMs := int64(10 * 60 * 1000) // 10 minutes
	if processingTimeMs > maxProcessingTimeMs {
		return fmt.Errorf("processing time %dms exceeds reasonable maximum of %dms", processingTimeMs, maxProcessingTimeMs)
	}

	return nil
}

// ValidateComponent checks if a component name is valid for logging.
func ValidateComponent(component string) error {
	if component == "" {
		return fmt.Errorf("component name cannot be empty")
	}

	// Component names should be lowercase with dots or dashes for hierarchy
	validPattern := regexp.MustCompile(`^[a-z][a-z0-9\-\.]*$`)
	if !validPattern.MatchString(component) {
		return fmt.Errorf("component name '%s' must be lowercase and contain only letters, numbers, dots, and dashes", component)
	}

	if len(component) > 100 {
		return fmt.Errorf("component name '%s' exceeds maximum length of 100 characters", component)
	}

	return nil
}

// ObservabilityValidationSuite runs comprehensive validation on observability data.
type ObservabilityValidationSuite struct {
	Errors   []string
	Warnings []string
}

// AddError adds an error to the validation suite.
func (suite *ObservabilityValidationSuite) AddError(err error) {
	if err != nil {
		suite.Errors = append(suite.Errors, err.Error())
	}
}

// AddWarning adds a warning to the validation suite.
func (suite *ObservabilityValidationSuite) AddWarning(warning string) {
	suite.Warnings = append(suite.Warnings, warning)
}

// IsValid returns true if there are no validation errors.
func (suite *ObservabilityValidationSuite) IsValid() bool {
	return len(suite.Errors) == 0
}

// Summary returns a formatted summary of validation results.
func (suite *ObservabilityValidationSuite) Summary() string {
	var parts []string

	if len(suite.Errors) > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", len(suite.Errors)))
	}

	if len(suite.Warnings) > 0 {
		parts = append(parts, fmt.Sprintf("%d warnings", len(suite.Warnings)))
	}

	if len(parts) == 0 {
		return "validation passed"
	}

	return strings.Join(parts, ", ")
}

// ValidateObservabilityMetadata performs comprehensive validation of telemetry metadata.
func ValidateObservabilityMetadata(traceID, spanID, requestID string, processingTimeMs int64, qualityScore float64) *ObservabilityValidationSuite {
	suite := &ObservabilityValidationSuite{}

	suite.AddError(ValidateTraceID(traceID))
	suite.AddError(ValidateSpanID(spanID))
	suite.AddError(ValidateProcessingTime(processingTimeMs))

	if qualityScore >= 0 { // Only validate if provided (negative means not set)
		suite.AddError(ValidateQualityScore(qualityScore))
	}

	if requestID == "" {
		suite.AddWarning("request ID is empty - recommended for request correlation")
	}

	if traceID == "" && spanID != "" {
		suite.AddError(fmt.Errorf("span ID provided without trace ID"))
	}

	return suite
}

// ValidateLogEntry performs comprehensive validation of log entries.
func ValidateLogEntry(level, message, component, traceID, spanID string, fields map[string]string) *ObservabilityValidationSuite {
	suite := &ObservabilityValidationSuite{}

	suite.AddError(ValidateLogLevel(level))
	suite.AddError(ValidateComponent(component))
	suite.AddError(ValidateTraceID(traceID))
	suite.AddError(ValidateSpanID(spanID))

	if message == "" {
		suite.AddError(fmt.Errorf("log message cannot be empty"))
	}

	if len(message) > 10000 {
		suite.AddError(fmt.Errorf("log message exceeds maximum length of 10000 characters"))
	}

	// Validate log fields
	for key, value := range fields {
		suite.AddError(ValidateLabelNameStrict(key))

		if len(value) > 1000 {
			suite.AddWarning(fmt.Sprintf("log field '%s' value exceeds recommended length of 1000 characters", key))
		}
	}

	if len(fields) > 50 {
		suite.AddWarning("log entry has more than 50 fields - consider reducing for performance")
	}

	return suite
}
