// Copyright 2024-2025 FinFocus Contributors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package pricing provides domain types, validation, billing mode enumerations,
// and audit metadata for FinFocus pricing specifications.
//
// The bypass metadata system allows for structured recording of why certain
// validation rules were skipped, providing an audit trail for security and
// compliance monitoring.
package pricing

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// MaxReasonLength is the maximum allowed length for bypass reason strings (in characters).
// Reasons exceeding this length will be truncated with "..." suffix.
const MaxReasonLength = 500

// BypassSeverity represents the risk level of a bypassed validation.
type BypassSeverity string

const (
	// BypassSeverityWarning indicates a low-risk bypass, informational alert.
	BypassSeverityWarning BypassSeverity = "warning"
	// BypassSeverityError indicates a medium-risk bypass that would have blocked the operation.
	BypassSeverityError BypassSeverity = "error"
	// BypassSeverityCritical indicates a high-risk bypass with security or compliance impact.
	BypassSeverityCritical BypassSeverity = "critical"
)

// allBypassSeverities is a package-level slice containing all valid BypassSeverity values.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allBypassSeverities = []BypassSeverity{
	BypassSeverityWarning,
	BypassSeverityError,
	BypassSeverityCritical,
}

// AllBypassSeverities returns a slice containing every defined BypassSeverity constant.
func AllBypassSeverities() []BypassSeverity {
	return allBypassSeverities
}

// String returns the bypass severity as a lowercase string value.
func (s BypassSeverity) String() string {
	return string(s)
}

// IsValidBypassSeverity checks if the given string represents a valid bypass severity.
func IsValidBypassSeverity(s string) bool {
	severity := BypassSeverity(s)
	for _, valid := range allBypassSeverities {
		if severity == valid {
			return true
		}
	}
	return false
}

// BypassMechanism represents how the bypass was triggered.
type BypassMechanism string

const (
	// BypassMechanismFlag indicates bypass via command-line flag (e.g., --yolo, --force).
	BypassMechanismFlag BypassMechanism = "flag"
	// BypassMechanismEnvVar indicates bypass via environment variable override.
	BypassMechanismEnvVar BypassMechanism = "env_var"
	// BypassMechanismConfig indicates bypass via configuration file setting.
	BypassMechanismConfig BypassMechanism = "config"
	// BypassMechanismProgrammatic indicates bypass via code-level API call.
	BypassMechanismProgrammatic BypassMechanism = "programmatic"
)

// allBypassMechanisms is a package-level slice containing all valid BypassMechanism values.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allBypassMechanisms = []BypassMechanism{
	BypassMechanismFlag,
	BypassMechanismEnvVar,
	BypassMechanismConfig,
	BypassMechanismProgrammatic,
}

// AllBypassMechanisms returns the canonical list of valid bypass mechanisms.
func AllBypassMechanisms() []BypassMechanism {
	return allBypassMechanisms
}

// String returns the bypass mechanism as a lowercase string value.
func (m BypassMechanism) String() string {
	return string(m)
}

// IsValidBypassMechanism reports whether m matches one of the defined BypassMechanism values.
func IsValidBypassMechanism(m string) bool {
	mechanism := BypassMechanism(m)
	for _, valid := range allBypassMechanisms {
		if mechanism == valid {
			return true
		}
	}
	return false
}

// BypassMetadata contains metadata about a validation bypass event for audit trails.
type BypassMetadata struct {
	// Timestamp is when the bypass occurred (UTC recommended).
	Timestamp time.Time `json:"timestamp"`
	// ValidationName is the identifier of the bypassed validation.
	ValidationName string `json:"validation_name"`
	// OriginalError is the error message that would have been shown.
	// WARNING: Ensure this field does not contain sensitive information
	// (API keys, credentials) before recording bypass metadata.
	OriginalError string `json:"original_error"`
	// Reason is the human-readable explanation for why the bypass was performed (max 500 characters).
	Reason string `json:"reason,omitempty"`
	// Operator is who triggered the bypass (user, service account, or "unknown").
	Operator string `json:"operator,omitempty"`
	// Severity is the risk level of the bypassed validation.
	Severity BypassSeverity `json:"severity"`
	// Mechanism is how the bypass was triggered.
	Mechanism BypassMechanism `json:"mechanism"`
	// Truncated indicates whether the reason was truncated due to length limits.
	Truncated bool `json:"truncated,omitempty"`
}

// BypassOption is a functional option for configuring BypassMetadata.
type BypassOption func(*BypassMetadata)

// NewBypassMetadata creates a new BypassMetadata with required fields and optional configuration.
// been produced (avoid embedding sensitive data in originalError).
func NewBypassMetadata(validationName, originalError string, opts ...BypassOption) BypassMetadata {
	m := BypassMetadata{
		Timestamp:      time.Now().UTC(),
		ValidationName: validationName,
		OriginalError:  originalError,
		Operator:       "unknown",
		Severity:       BypassSeverityError,
		Mechanism:      BypassMechanismFlag,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

// WithReason sets the bypass reason, truncating to 500 characters (runes) if necessary.
// WithReason returns a BypassOption that sets the metadata Reason and marks it as truncated when necessary.
// 
// The provided reason is stored on the BypassMetadata. If the reason exceeds MaxReasonLength characters it
// is truncated to fit and "..." is appended; in that case Truncated is set to true. Truncation is performed
// safely on runes to preserve UTF-8 correctness.
func WithReason(reason string) BypassOption {
	return func(m *BypassMetadata) {
		runes := []rune(reason)
		if len(runes) > MaxReasonLength {
			m.Reason = string(runes[:MaxReasonLength-3]) + "..."
			m.Truncated = true
		} else {
			m.Reason = reason
		}
	}
}

// WithOperator returns a BypassOption that sets the Operator field on a BypassMetadata
// to the provided operator if the string is non-empty.
func WithOperator(operator string) BypassOption {
	return func(m *BypassMetadata) {
		if operator != "" {
			m.Operator = operator
		}
	}
}

// WithSeverity sets the severity level for the bypass.
// WithSeverity returns a BypassOption that sets the BypassMetadata Severity to the provided value.
// The option does not validate the severity value.
func WithSeverity(severity BypassSeverity) BypassOption {
	return func(m *BypassMetadata) {
		m.Severity = severity
	}
}

// WithMechanism sets the mechanism type for the bypass.
// WithMechanism returns a BypassOption that sets the Mechanism field on a BypassMetadata.
// No validation of the provided mechanism is performed by this option.
func WithMechanism(mechanism BypassMechanism) BypassOption {
	return func(m *BypassMetadata) {
		m.Mechanism = mechanism
	}
}

// ValidateBypassMetadata validates a BypassMetadata and returns an error describing any violations.
// It ensures Timestamp is not zero, ValidationName and OriginalError are non-empty, and that Severity
// and Mechanism are one of the allowed values. If one or more checks fail, the returned error
// aggregates all violation messages joined by "; ".
func ValidateBypassMetadata(m BypassMetadata) error {
	var errs []string

	if m.Timestamp.IsZero() {
		errs = append(errs, "timestamp must not be zero")
	}

	if m.ValidationName == "" {
		errs = append(errs, "validation_name is required")
	}

	if m.OriginalError == "" {
		errs = append(errs, "original_error is required")
	}

	if !IsValidBypassSeverity(string(m.Severity)) {
		errs = append(errs, fmt.Sprintf("invalid severity: %s", m.Severity))
	}

	if !IsValidBypassMechanism(string(m.Mechanism)) {
		errs = append(errs, fmt.Sprintf("invalid mechanism: %s", m.Mechanism))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

// HasBypasses reports whether the ValidationResult contains any bypass metadata.
func HasBypasses(r ValidationResult) bool {
	return len(r.Bypasses) > 0
}

// CountBypassesBySeverity counts bypass metadata entries by their severity.
// It returns a map from BypassSeverity to the number of occurrences for each severity.
func CountBypassesBySeverity(bypasses []BypassMetadata) map[BypassSeverity]int {
	counts := make(map[BypassSeverity]int)
	for _, b := range bypasses {
		counts[b.Severity]++
	}
	return counts
}

// FormatBypassSummary builds a short human-readable summary of bypass counts by severity for CLI output.
// It returns an empty string if the provided slice has no bypasses.
// The summary has the form "Bypassed validations: X critical, Y error, Z warning" and includes only severities with non-zero counts.
func FormatBypassSummary(bypasses []BypassMetadata) string {
	if len(bypasses) == 0 {
		return ""
	}

	counts := CountBypassesBySeverity(bypasses)
	var parts []string

	if c := counts[BypassSeverityCritical]; c > 0 {
		parts = append(parts, fmt.Sprintf("%d critical", c))
	}
	if c := counts[BypassSeverityError]; c > 0 {
		parts = append(parts, fmt.Sprintf("%d error", c))
	}
	if c := counts[BypassSeverityWarning]; c > 0 {
		parts = append(parts, fmt.Sprintf("%d warning", c))
	}

	return fmt.Sprintf("Bypassed validations: %s", strings.Join(parts, ", "))
}

// FormatBypassDetail formats a single BypassMetadata into a multi-line,
// human-readable detail suitable for CLI output.
// The output includes a header with uppercased severity and validation name,
// the original error, an optional reason if present, the operator, the mechanism,
// and a timestamp formatted as RFC3339.
func FormatBypassDetail(b BypassMetadata) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s] %s\n", strings.ToUpper(string(b.Severity)), b.ValidationName))
	sb.WriteString(fmt.Sprintf("  Original error: %s\n", b.OriginalError))

	if b.Reason != "" {
		sb.WriteString(fmt.Sprintf("  Reason: %s\n", b.Reason))
	}

	sb.WriteString(fmt.Sprintf("  Operator: %s\n", b.Operator))
	sb.WriteString(fmt.Sprintf("  Mechanism: %s\n", b.Mechanism))
	sb.WriteString(fmt.Sprintf("  Time: %s\n", b.Timestamp.Format(time.RFC3339)))

	return sb.String()
}

// MarshalZerologObject implements the zerolog.LogObjectMarshaler interface.
// This allows for zero-allocation logging of BypassMetadata with zerolog.
func (m BypassMetadata) MarshalZerologObject(e *zerolog.Event) {
	e.Time("timestamp", m.Timestamp)
	e.Str("validation_name", m.ValidationName)
	e.Str("original_error", m.OriginalError)
	e.Str("severity", string(m.Severity))
	e.Str("mechanism", string(m.Mechanism))

	if m.Reason != "" {
		e.Str("reason", m.Reason)
	}
	if m.Operator != "" {
		e.Str("operator", m.Operator)
	}
	if m.Truncated {
		e.Bool("truncated", true)
	}
}

// FilterByTimeRange returns the subset of bypass metadata whose Timestamp falls within the inclusive time range defined by start and end.
// The original slice order is preserved in the returned slice.
func FilterByTimeRange(bypasses []BypassMetadata, start, end time.Time) []BypassMetadata {
	var result []BypassMetadata
	for _, b := range bypasses {
		if (b.Timestamp.Equal(start) || b.Timestamp.After(start)) &&
			(b.Timestamp.Equal(end) || b.Timestamp.Before(end)) {
			result = append(result, b)
		}
	}
	return result
}

// FilterByOperator filters bypass metadata and returns entries whose Operator equals the provided operator.
func FilterByOperator(bypasses []BypassMetadata, operator string) []BypassMetadata {
	var result []BypassMetadata
	for _, b := range bypasses {
		if b.Operator == operator {
			result = append(result, b)
		}
	}
	return result
}

// FilterBySeverity returns all BypassMetadata entries whose Severity matches the provided severity.
func FilterBySeverity(bypasses []BypassMetadata, severity BypassSeverity) []BypassMetadata {
	var result []BypassMetadata
	for _, b := range bypasses {
		if b.Severity == severity {
			result = append(result, b)
		}
	}
	return result
}

// FilterByMechanism filters the provided bypass metadata to those whose Mechanism equals the specified mechanism.
func FilterByMechanism(bypasses []BypassMetadata, mechanism BypassMechanism) []BypassMetadata {
	var result []BypassMetadata
	for _, b := range bypasses {
		if b.Mechanism == mechanism {
			result = append(result, b)
		}
	}
	return result
}