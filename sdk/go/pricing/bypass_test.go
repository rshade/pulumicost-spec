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

package pricing_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
)

// =============================================================================
// Phase 2: Enum Validation Tests (T011, T012)
// =============================================================================

func TestIsValidBypassSeverity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid warning", "warning", true},
		{"valid error", "error", true},
		{"valid critical", "critical", true},
		{"invalid empty", "", false},
		{"invalid unknown", "unknown", false},
		{"invalid uppercase", "WARNING", false},
		{"invalid mixed case", "Warning", false},
		{"invalid info", "info", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.IsValidBypassSeverity(tt.input)
			if result != tt.expected {
				t.Errorf("pricing.IsValidBypassSeverity(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidBypassMechanism(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid flag", "flag", true},
		{"valid env_var", "env_var", true},
		{"valid config", "config", true},
		{"valid programmatic", "programmatic", true},
		{"invalid empty", "", false},
		{"invalid unknown", "unknown", false},
		{"invalid uppercase", "FLAG", false},
		{"invalid mixed case", "Flag", false},
		{"invalid envvar no underscore", "envvar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.IsValidBypassMechanism(tt.input)
			if result != tt.expected {
				t.Errorf("pricing.IsValidBypassMechanism(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAllBypassSeverities(t *testing.T) {
	severities := pricing.AllBypassSeverities()

	// Verify count
	expectedCount := 3
	if len(severities) != expectedCount {
		t.Errorf(
			"pricing.AllBypassSeverities() returned %d severities, want %d",
			len(severities),
			expectedCount,
		)
	}

	// Verify all are valid
	for _, s := range severities {
		if !pricing.IsValidBypassSeverity(string(s)) {
			t.Errorf("pricing.AllBypassSeverities() contains invalid severity: %s", s)
		}
	}

	// Verify expected values present
	expectedValues := map[pricing.BypassSeverity]bool{
		pricing.BypassSeverityWarning:  false,
		pricing.BypassSeverityError:    false,
		pricing.BypassSeverityCritical: false,
	}

	for _, s := range severities {
		expectedValues[s] = true
	}

	for s, found := range expectedValues {
		if !found {
			t.Errorf("pricing.AllBypassSeverities() missing expected severity: %s", s)
		}
	}
}

func TestAllBypassMechanisms(t *testing.T) {
	mechanisms := pricing.AllBypassMechanisms()

	// Verify count
	expectedCount := 4
	if len(mechanisms) != expectedCount {
		t.Errorf(
			"pricing.AllBypassMechanisms() returned %d mechanisms, want %d",
			len(mechanisms),
			expectedCount,
		)
	}

	// Verify all are valid
	for _, m := range mechanisms {
		if !pricing.IsValidBypassMechanism(string(m)) {
			t.Errorf("pricing.AllBypassMechanisms() contains invalid mechanism: %s", m)
		}
	}

	// Verify expected values present
	expectedValues := map[pricing.BypassMechanism]bool{
		pricing.BypassMechanismFlag:         false,
		pricing.BypassMechanismEnvVar:       false,
		pricing.BypassMechanismConfig:       false,
		pricing.BypassMechanismProgrammatic: false,
	}

	for _, m := range mechanisms {
		expectedValues[m] = true
	}

	for m, found := range expectedValues {
		if !found {
			t.Errorf("pricing.AllBypassMechanisms() missing expected mechanism: %s", m)
		}
	}
}

func TestBypassSeverityString(t *testing.T) {
	tests := []struct {
		severity pricing.BypassSeverity
		expected string
	}{
		{pricing.BypassSeverityWarning, "warning"},
		{pricing.BypassSeverityError, "error"},
		{pricing.BypassSeverityCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.severity.String(); got != tt.expected {
				t.Errorf("BypassSeverity.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestBypassMechanismString(t *testing.T) {
	tests := []struct {
		mechanism pricing.BypassMechanism
		expected  string
	}{
		{pricing.BypassMechanismFlag, "flag"},
		{pricing.BypassMechanismEnvVar, "env_var"},
		{pricing.BypassMechanismConfig, "config"},
		{pricing.BypassMechanismProgrammatic, "programmatic"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.mechanism.String(); got != tt.expected {
				t.Errorf("BypassMechanism.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// =============================================================================
// Phase 2: Benchmark Tests (T013)
// =============================================================================

func BenchmarkIsValidBypassSeverity(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		pricing.IsValidBypassSeverity("error")
	}
}

func BenchmarkIsValidBypassSeverity_Invalid(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		pricing.IsValidBypassSeverity("invalid")
	}
}

func BenchmarkIsValidBypassMechanism(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		pricing.IsValidBypassMechanism("flag")
	}
}

func BenchmarkIsValidBypassMechanism_Invalid(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		pricing.IsValidBypassMechanism("invalid")
	}
}

func BenchmarkAllBypassSeverities(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pricing.AllBypassSeverities()
	}
}

func BenchmarkAllBypassMechanisms(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pricing.AllBypassMechanisms()
	}
}

// =============================================================================
// Phase 3: BypassMetadata Constructor and Validation Tests (T023)
// =============================================================================

func TestNewBypassMetadata(t *testing.T) {
	validationName := "budget_limit"
	originalError := "Cost exceeds budget by $500"

	m := pricing.NewBypassMetadata(validationName, originalError)

	if m.ValidationName != validationName {
		t.Errorf("ValidationName = %q, want %q", m.ValidationName, validationName)
	}
	if m.OriginalError != originalError {
		t.Errorf("OriginalError = %q, want %q", m.OriginalError, originalError)
	}
	if m.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	if m.Operator != "unknown" {
		t.Errorf("Operator = %q, want %q", m.Operator, "unknown")
	}
	if m.Severity != pricing.BypassSeverityError {
		t.Errorf("Severity = %q, want %q", m.Severity, pricing.BypassSeverityError)
	}
	if m.Mechanism != pricing.BypassMechanismFlag {
		t.Errorf("Mechanism = %q, want %q", m.Mechanism, pricing.BypassMechanismFlag)
	}
}

func TestNewBypassMetadata_WithOptions(t *testing.T) {
	m := pricing.NewBypassMetadata(
		"budget_limit",
		"Cost exceeds budget",
		pricing.WithReason("Emergency deployment"),
		pricing.WithOperator("user@example.com"),
		pricing.WithSeverity(pricing.BypassSeverityCritical),
		pricing.WithMechanism(pricing.BypassMechanismEnvVar),
	)

	if m.Reason != "Emergency deployment" {
		t.Errorf("Reason = %q, want %q", m.Reason, "Emergency deployment")
	}
	if m.Operator != "user@example.com" {
		t.Errorf("Operator = %q, want %q", m.Operator, "user@example.com")
	}
	if m.Severity != pricing.BypassSeverityCritical {
		t.Errorf("Severity = %q, want %q", m.Severity, pricing.BypassSeverityCritical)
	}
	if m.Mechanism != pricing.BypassMechanismEnvVar {
		t.Errorf("Mechanism = %q, want %q", m.Mechanism, pricing.BypassMechanismEnvVar)
	}
	if m.Truncated {
		t.Error("Truncated should be false for short reason")
	}
}

func TestWithOperator_Empty(t *testing.T) {
	m := pricing.NewBypassMetadata("test", "error", pricing.WithOperator(""))

	if m.Operator != "unknown" {
		t.Errorf("Operator = %q, want %q for empty input", m.Operator, "unknown")
	}
}

func TestValidateBypassMetadata_Valid(t *testing.T) {
	m := pricing.NewBypassMetadata(
		"budget_limit",
		"Cost exceeds budget",
		pricing.WithReason("Emergency deployment"),
	)

	if err := pricing.ValidateBypassMetadata(m); err != nil {
		t.Errorf("pricing.ValidateBypassMetadata() returned unexpected error: %v", err)
	}
}

func TestValidateBypassMetadata_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		metadata pricing.BypassMetadata
		wantErr  string
	}{
		{
			name: "zero timestamp",
			metadata: pricing.BypassMetadata{
				ValidationName: "test",
				OriginalError:  "err",
				Severity:       pricing.BypassSeverityError,
				Mechanism:      pricing.BypassMechanismFlag,
			},
			wantErr: "timestamp must not be zero",
		},
		{
			name: "empty validation name",
			metadata: pricing.BypassMetadata{
				Timestamp:     time.Now(),
				OriginalError: "err",
				Severity:      pricing.BypassSeverityError,
				Mechanism:     pricing.BypassMechanismFlag,
			},
			wantErr: "validation_name is required",
		},
		{
			name: "empty original error",
			metadata: pricing.BypassMetadata{
				Timestamp:      time.Now(),
				ValidationName: "test",
				Severity:       pricing.BypassSeverityError,
				Mechanism:      pricing.BypassMechanismFlag,
			},
			wantErr: "original_error is required",
		},
		{
			name: "invalid severity",
			metadata: pricing.BypassMetadata{
				Timestamp:      time.Now(),
				ValidationName: "test",
				OriginalError:  "err",
				Severity:       "invalid",
				Mechanism:      pricing.BypassMechanismFlag,
			},
			wantErr: "invalid severity",
		},
		{
			name: "invalid mechanism",
			metadata: pricing.BypassMetadata{
				Timestamp:      time.Now(),
				ValidationName: "test",
				OriginalError:  "err",
				Severity:       pricing.BypassSeverityError,
				Mechanism:      "invalid",
			},
			wantErr: "invalid mechanism",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateBypassMetadata(tt.metadata)
			if err == nil {
				t.Error("pricing.ValidateBypassMetadata() should return error")
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf(
					"pricing.ValidateBypassMetadata() error = %q, want containing %q",
					err.Error(),
					tt.wantErr,
				)
			}
		})
	}
}

// =============================================================================
// Phase 3: Reason Truncation Tests (T024)
// =============================================================================

func TestWithReason_ExactlyMaxLength(t *testing.T) {
	reason := strings.Repeat("a", pricing.MaxReasonLength)
	m := pricing.NewBypassMetadata("test", "error", pricing.WithReason(reason))

	if len(m.Reason) != pricing.MaxReasonLength {
		t.Errorf("Reason length = %d, want %d", len(m.Reason), pricing.MaxReasonLength)
	}
	if m.Truncated {
		t.Error("Truncated should be false for exactly max length reason")
	}
}

func TestWithReason_OverMaxLength(t *testing.T) {
	reason := strings.Repeat("a", pricing.MaxReasonLength+100)
	m := pricing.NewBypassMetadata("test", "error", pricing.WithReason(reason))

	if len(m.Reason) != pricing.MaxReasonLength {
		t.Errorf("Reason length = %d, want %d", len(m.Reason), pricing.MaxReasonLength)
	}
	if !strings.HasSuffix(m.Reason, "...") {
		t.Error("Truncated reason should end with '...'")
	}
	if !m.Truncated {
		t.Error("Truncated should be true for over max length reason")
	}
}

func TestWithReason_UnderMaxLength(t *testing.T) {
	reason := "Short reason"
	m := pricing.NewBypassMetadata("test", "error", pricing.WithReason(reason))

	if m.Reason != reason {
		t.Errorf("Reason = %q, want %q", m.Reason, reason)
	}
	if m.Truncated {
		t.Error("Truncated should be false for under max length reason")
	}
}

func TestWithReason_Empty(t *testing.T) {
	m := pricing.NewBypassMetadata("test", "error", pricing.WithReason(""))

	if m.Reason != "" {
		t.Errorf("Reason = %q, want empty string", m.Reason)
	}
	if m.Truncated {
		t.Error("Truncated should be false for empty reason")
	}
}

func TestWithReason_MultiByte(t *testing.T) {
	// Create a reason string with multi-byte characters (Japanese "日", 3 bytes each)
	// repeated enough times to exceed MaxReasonLength (500 runes).
	// 600 chars > 500 chars limit.
	reason := strings.Repeat("日", 600)

	m := pricing.NewBypassMetadata("test", "error", pricing.WithReason(reason))

	// Assert truncation occurred
	if !m.Truncated {
		t.Error("Truncated should be true for multi-byte string exceeding max length")
	}

	// Assert valid UTF-8
	if !utf8.ValidString(m.Reason) {
		t.Error("Truncated reason contains invalid UTF-8")
	}

	// Assert length in runes is at most MaxReasonLength
	runeCount := utf8.RuneCountInString(m.Reason)
	if runeCount > pricing.MaxReasonLength {
		t.Errorf("Reason rune count = %d, want <= %d", runeCount, pricing.MaxReasonLength)
	}

	// Assert specific structure (should end in "...")
	if !strings.HasSuffix(m.Reason, "...") {
		t.Errorf("Truncated reason should end with '...', got suffix: %q", m.Reason[len(m.Reason)-10:])
	}
}

// =============================================================================
// Phase 3: JSON Round-Trip Tests (T025)
// =============================================================================

func TestBypassMetadata_JSONRoundTrip(t *testing.T) {
	original := pricing.NewBypassMetadata(
		"budget_limit",
		"Cost exceeds budget by $500",
		pricing.WithReason("Emergency deployment approved by manager"),
		pricing.WithOperator("user@example.com"),
		pricing.WithSeverity(pricing.BypassSeverityError),
		pricing.WithMechanism(pricing.BypassMechanismFlag),
	)

	// Serialize
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error: %v", err)
	}

	// Deserialize
	var restored pricing.BypassMetadata
	if unmarshalErr := json.Unmarshal(data, &restored); unmarshalErr != nil {
		t.Fatalf("json.Unmarshal() error: %v", unmarshalErr)
	}

	// Compare fields
	if !restored.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Timestamp mismatch: got %v, want %v", restored.Timestamp, original.Timestamp)
	}
	if restored.ValidationName != original.ValidationName {
		t.Errorf("ValidationName = %q, want %q", restored.ValidationName, original.ValidationName)
	}
	if restored.OriginalError != original.OriginalError {
		t.Errorf("OriginalError = %q, want %q", restored.OriginalError, original.OriginalError)
	}
	if restored.Reason != original.Reason {
		t.Errorf("Reason = %q, want %q", restored.Reason, original.Reason)
	}
	if restored.Operator != original.Operator {
		t.Errorf("Operator = %q, want %q", restored.Operator, original.Operator)
	}
	if restored.Severity != original.Severity {
		t.Errorf("Severity = %q, want %q", restored.Severity, original.Severity)
	}
	if restored.Mechanism != original.Mechanism {
		t.Errorf("Mechanism = %q, want %q", restored.Mechanism, original.Mechanism)
	}
}

func TestValidationResult_WithBypasses_JSONRoundTrip(t *testing.T) {
	original := pricing.ValidationResult{
		Valid: true,
		Bypasses: []pricing.BypassMetadata{
			pricing.NewBypassMetadata(
				"budget_limit",
				"Cost exceeds budget",
				pricing.WithReason("Emergency deployment"),
				pricing.WithOperator("user@example.com"),
				pricing.WithSeverity(pricing.BypassSeverityError),
				pricing.WithMechanism(pricing.BypassMechanismFlag),
			),
			pricing.NewBypassMetadata(
				"region_policy",
				"Region not in allowlist",
				pricing.WithReason("Temporary exception"),
				pricing.WithOperator("admin@example.com"),
				pricing.WithSeverity(pricing.BypassSeverityWarning),
				pricing.WithMechanism(pricing.BypassMechanismConfig),
			),
		},
	}

	// Serialize
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error: %v", err)
	}

	// Verify JSON structure
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"bypasses"`) {
		t.Error("JSON should contain 'bypasses' field")
	}
	if !strings.Contains(jsonStr, `"validation_name"`) {
		t.Error("JSON should contain 'validation_name' field (snake_case)")
	}
	if !strings.Contains(jsonStr, `"original_error"`) {
		t.Error("JSON should contain 'original_error' field (snake_case)")
	}

	// Deserialize
	var restored pricing.ValidationResult
	if unmarshalErr := json.Unmarshal(data, &restored); unmarshalErr != nil {
		t.Fatalf("json.Unmarshal() error: %v", unmarshalErr)
	}

	// Compare
	if restored.Valid != original.Valid {
		t.Errorf("Valid = %v, want %v", restored.Valid, original.Valid)
	}
	if len(restored.Bypasses) != len(original.Bypasses) {
		t.Fatalf("Bypasses length = %d, want %d", len(restored.Bypasses), len(original.Bypasses))
	}

	for i, b := range restored.Bypasses {
		if b.ValidationName != original.Bypasses[i].ValidationName {
			t.Errorf(
				"Bypass[%d].ValidationName = %q, want %q",
				i,
				b.ValidationName,
				original.Bypasses[i].ValidationName,
			)
		}
		if b.Severity != original.Bypasses[i].Severity {
			t.Errorf(
				"Bypass[%d].Severity = %q, want %q",
				i,
				b.Severity,
				original.Bypasses[i].Severity,
			)
		}
	}
}

func TestValidationResult_EmptyBypasses_JSONOmitsField(t *testing.T) {
	result := pricing.ValidationResult{
		Valid:    true,
		Bypasses: nil,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error: %v", err)
	}

	jsonStr := string(data)
	if strings.Contains(jsonStr, `"bypasses"`) {
		t.Error("JSON should omit 'bypasses' field when nil (omitempty)")
	}
}

// =============================================================================
// Phase 4: Format Functions Tests (T026, T027)
// =============================================================================

func TestFormatBypassSummary(t *testing.T) {
	tests := []struct {
		name     string
		bypasses []pricing.BypassMetadata
		wantPart string
	}{
		{
			name:     "empty bypasses",
			bypasses: nil,
			wantPart: "",
		},
		{
			name: "single critical",
			bypasses: []pricing.BypassMetadata{
				{Severity: pricing.BypassSeverityCritical},
			},
			wantPart: "1 critical",
		},
		{
			name: "mixed severities",
			bypasses: []pricing.BypassMetadata{
				{Severity: pricing.BypassSeverityCritical},
				{Severity: pricing.BypassSeverityError},
				{Severity: pricing.BypassSeverityError},
				{Severity: pricing.BypassSeverityWarning},
			},
			wantPart: "2 error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.FormatBypassSummary(tt.bypasses)
			if tt.wantPart == "" {
				if result != "" {
					t.Errorf("pricing.FormatBypassSummary() = %q, want empty string", result)
				}
			} else if !strings.Contains(result, tt.wantPart) {
				t.Errorf("pricing.FormatBypassSummary() = %q, want containing %q", result, tt.wantPart)
			}
		})
	}
}

func TestFormatBypassDetail(t *testing.T) {
	ts := time.Date(2026, 1, 24, 10, 30, 0, 0, time.UTC)
	b := pricing.BypassMetadata{
		Timestamp:      ts,
		ValidationName: "budget_limit",
		OriginalError:  "Cost exceeds budget",
		Reason:         "Emergency deployment",
		Operator:       "user@example.com",
		Severity:       pricing.BypassSeverityError,
		Mechanism:      pricing.BypassMechanismFlag,
	}

	result := pricing.FormatBypassDetail(b)

	expectedParts := []string{
		"[ERROR]",
		"budget_limit",
		"Original error: Cost exceeds budget",
		"Reason: Emergency deployment",
		"Operator: user@example.com",
		"Mechanism: flag",
		"2026-01-24T10:30:00Z",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("pricing.FormatBypassDetail() missing %q in output:\n%s", part, result)
		}
	}
}

func TestBypassMetadata_MarshalZerologObject(t *testing.T) {
	ts := time.Date(2026, 1, 24, 10, 30, 0, 0, time.UTC)
	b := pricing.BypassMetadata{
		Timestamp:      ts,
		ValidationName: "budget_limit",
		OriginalError:  "Cost exceeds budget",
		Reason:         "Emergency deployment",
		Operator:       "user@example.com",
		Severity:       pricing.BypassSeverityError,
		Mechanism:      pricing.BypassMechanismFlag,
		Truncated:      true,
	}

	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	logger.Info().Object("bypass", b).Msg("Bypass event")

	// Verify JSON structure and fields
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	bypassObj, ok := logEntry["bypass"].(map[string]interface{})
	if !ok {
		t.Fatal("Log output missing 'bypass' object")
	}

	// Helper to check string fields
	checkField := func(key, expected string) {
		val, okVal := bypassObj[key].(string)
		if !okVal || val != expected {
			t.Errorf("Field %q = %v, want %q", key, val, expected)
		}
	}

	checkField("validation_name", "budget_limit")
	checkField("original_error", "Cost exceeds budget")
	checkField("reason", "Emergency deployment")
	checkField("operator", "user@example.com")
	checkField("severity", "error")
	checkField("mechanism", "flag")
	checkField("timestamp", ts.Format(time.RFC3339))

	if truncated, okVal := bypassObj["truncated"].(bool); !okVal || !truncated {
		t.Error("Field 'truncated' should be true")
	}
}

// =============================================================================
// Phase 5: Filter Functions Tests (T032, T033, T034)
// =============================================================================

func TestFilterByTimeRange(t *testing.T) {
	now := time.Now().UTC()
	bypasses := []pricing.BypassMetadata{
		{Timestamp: now.Add(-48 * time.Hour), ValidationName: "old"},
		{Timestamp: now.Add(-24 * time.Hour), ValidationName: "yesterday"},
		{Timestamp: now, ValidationName: "today"},
		{Timestamp: now.Add(24 * time.Hour), ValidationName: "tomorrow"},
	}

	// Filter for last 25 hours
	start := now.Add(-25 * time.Hour)
	end := now.Add(time.Hour)

	filtered := pricing.FilterByTimeRange(bypasses, start, end)

	if len(filtered) != 2 {
		t.Errorf("pricing.FilterByTimeRange() returned %d items, want 2", len(filtered))
	}

	for _, b := range filtered {
		if b.ValidationName == "old" || b.ValidationName == "tomorrow" {
			t.Errorf("pricing.FilterByTimeRange() should not include %q", b.ValidationName)
		}
	}
}

func TestFilterByOperator(t *testing.T) {
	bypasses := []pricing.BypassMetadata{
		{Operator: "user1@example.com", ValidationName: "test1"},
		{Operator: "user2@example.com", ValidationName: "test2"},
		{Operator: "user1@example.com", ValidationName: "test3"},
	}

	filtered := pricing.FilterByOperator(bypasses, "user1@example.com")

	if len(filtered) != 2 {
		t.Errorf("pricing.FilterByOperator() returned %d items, want 2", len(filtered))
	}

	for _, b := range filtered {
		if b.Operator != "user1@example.com" {
			t.Errorf("pricing.FilterByOperator() returned wrong operator: %q", b.Operator)
		}
	}
}

func TestFilterBySeverity(t *testing.T) {
	bypasses := []pricing.BypassMetadata{
		{Severity: pricing.BypassSeverityWarning, ValidationName: "test1"},
		{Severity: pricing.BypassSeverityError, ValidationName: "test2"},
		{Severity: pricing.BypassSeverityCritical, ValidationName: "test3"},
		{Severity: pricing.BypassSeverityError, ValidationName: "test4"},
	}

	filtered := pricing.FilterBySeverity(bypasses, pricing.BypassSeverityError)

	if len(filtered) != 2 {
		t.Errorf("pricing.FilterBySeverity() returned %d items, want 2", len(filtered))
	}

	for _, b := range filtered {
		if b.Severity != pricing.BypassSeverityError {
			t.Errorf("pricing.FilterBySeverity() returned wrong severity: %q", b.Severity)
		}
	}
}

func TestFilterByMechanism(t *testing.T) {
	bypasses := []pricing.BypassMetadata{
		{Mechanism: pricing.BypassMechanismFlag, ValidationName: "test1"},
		{Mechanism: pricing.BypassMechanismEnvVar, ValidationName: "test2"},
		{Mechanism: pricing.BypassMechanismFlag, ValidationName: "test3"},
		{Mechanism: pricing.BypassMechanismConfig, ValidationName: "test4"},
	}

	filtered := pricing.FilterByMechanism(bypasses, pricing.BypassMechanismFlag)

	if len(filtered) != 2 {
		t.Errorf("pricing.FilterByMechanism() returned %d items, want 2", len(filtered))
	}

	for _, b := range filtered {
		if b.Mechanism != pricing.BypassMechanismFlag {
			t.Errorf("pricing.FilterByMechanism() returned wrong mechanism: %q", b.Mechanism)
		}
	}
}

func TestFilterFunctions_NilAndEmptyInput(t *testing.T) {
	tests := []struct {
		name     string
		filterFn func([]pricing.BypassMetadata) []pricing.BypassMetadata
	}{
		{"FilterByTimeRange", func(b []pricing.BypassMetadata) []pricing.BypassMetadata {
			return pricing.FilterByTimeRange(b, time.Now(), time.Now().Add(time.Hour))
		}},
		{"FilterByOperator", func(b []pricing.BypassMetadata) []pricing.BypassMetadata {
			return pricing.FilterByOperator(b, "test@example.com")
		}},
		{"FilterBySeverity", func(b []pricing.BypassMetadata) []pricing.BypassMetadata {
			return pricing.FilterBySeverity(b, pricing.BypassSeverityError)
		}},
		{"FilterByMechanism", func(b []pricing.BypassMetadata) []pricing.BypassMetadata {
			return pricing.FilterByMechanism(b, pricing.BypassMechanismFlag)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name+" nil input", func(t *testing.T) {
			result := tt.filterFn(nil)
			if len(result) != 0 {
				t.Error("Filter should return empty/nil slice for nil input")
			}
		})
		t.Run(tt.name+" empty input", func(t *testing.T) {
			result := tt.filterFn([]pricing.BypassMetadata{})
			if len(result) != 0 {
				t.Error("Filter should return empty slice for empty input")
			}
		})
	}
}

// =============================================================================
// Helper Function Tests
// =============================================================================

func TestHasBypasses(t *testing.T) {
	tests := []struct {
		name   string
		result pricing.ValidationResult
		want   bool
	}{
		{
			name:   "nil bypasses",
			result: pricing.ValidationResult{Bypasses: nil},
			want:   false,
		},
		{
			name:   "empty bypasses",
			result: pricing.ValidationResult{Bypasses: []pricing.BypassMetadata{}},
			want:   false,
		},
		{
			name:   "has bypasses",
			result: pricing.ValidationResult{Bypasses: []pricing.BypassMetadata{{}}},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pricing.HasBypasses(tt.result); got != tt.want {
				t.Errorf("pricing.HasBypasses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountBypassesBySeverity(t *testing.T) {
	bypasses := []pricing.BypassMetadata{
		{Severity: pricing.BypassSeverityWarning},
		{Severity: pricing.BypassSeverityError},
		{Severity: pricing.BypassSeverityError},
		{Severity: pricing.BypassSeverityCritical},
	}

	counts := pricing.CountBypassesBySeverity(bypasses)

	if counts[pricing.BypassSeverityWarning] != 1 {
		t.Errorf("Warning count = %d, want 1", counts[pricing.BypassSeverityWarning])
	}
	if counts[pricing.BypassSeverityError] != 2 {
		t.Errorf("Error count = %d, want 2", counts[pricing.BypassSeverityError])
	}
	if counts[pricing.BypassSeverityCritical] != 1 {
		t.Errorf("Critical count = %d, want 1", counts[pricing.BypassSeverityCritical])
	}
}

// Example_operatorNormalization demonstrates the recommended pattern for
// normalizing operator identifiers to ensure consistent exact-match filtering.
func Example_operatorNormalization() {
	// Raw operator input from user or system (may have mixed case, whitespace)
	rawOperator := "  Admin@Example.COM  "

	// Normalize at write time for consistent filtering later
	normalizedOperator := strings.ToLower(strings.TrimSpace(rawOperator))

	// Create bypass metadata with normalized operator
	bypass := pricing.NewBypassMetadata(
		"budget_limit",
		"Cost exceeds budget",
		pricing.WithReason("Emergency deployment"),
		pricing.WithOperator(normalizedOperator),
	)

	// Later, when filtering for audit queries, use the same normalized form
	bypasses := []pricing.BypassMetadata{bypass}
	filtered := pricing.FilterByOperator(bypasses, "admin@example.com")

	fmt.Printf("Normalized operator: %s\n", bypass.Operator)
	fmt.Printf("Found %d matching bypass(es)\n", len(filtered))

	// Output:
	// Normalized operator: admin@example.com
	// Found 1 matching bypass(es)
}
