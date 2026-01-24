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

package testing_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
)

// TestBypassMetadata_JSONConformance verifies that BypassMetadata serializes
// and deserializes correctly according to the specification.
func TestBypassMetadata_JSONConformance(t *testing.T) {
	// Create bypass metadata with all fields populated
	original := pricing.NewBypassMetadata(
		"budget_limit",
		"Cost exceeds budget by $500",
		pricing.WithReason("Emergency deployment approved by manager"),
		pricing.WithOperator("user@example.com"),
		pricing.WithSeverity(pricing.BypassSeverityError),
		pricing.WithMechanism(pricing.BypassMechanismFlag),
	)

	// Serialize to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal BypassMetadata: %v", err)
	}

	// Verify JSON field naming (snake_case)
	jsonStr := string(data)
	expectedFields := []string{
		`"timestamp"`,
		`"validation_name"`,
		`"original_error"`,
		`"reason"`,
		`"operator"`,
		`"severity"`,
		`"mechanism"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON should contain field %s, got: %s", field, jsonStr)
		}
	}

	// Deserialize and verify round-trip
	var restored pricing.BypassMetadata
	if unmarshalErr := json.Unmarshal(data, &restored); unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal BypassMetadata: %v", unmarshalErr)
	}

	// Verify all fields are preserved
	if !restored.Timestamp.Equal(original.Timestamp) {
		t.Errorf("Timestamp not preserved: got %v, want %v", restored.Timestamp, original.Timestamp)
	}
	if restored.ValidationName != original.ValidationName {
		t.Errorf("ValidationName not preserved: got %q, want %q", restored.ValidationName, original.ValidationName)
	}
	if restored.OriginalError != original.OriginalError {
		t.Errorf("OriginalError not preserved: got %q, want %q", restored.OriginalError, original.OriginalError)
	}
	if restored.Reason != original.Reason {
		t.Errorf("Reason not preserved: got %q, want %q", restored.Reason, original.Reason)
	}
	if restored.Operator != original.Operator {
		t.Errorf("Operator not preserved: got %q, want %q", restored.Operator, original.Operator)
	}
	if restored.Severity != original.Severity {
		t.Errorf("Severity not preserved: got %q, want %q", restored.Severity, original.Severity)
	}
	if restored.Mechanism != original.Mechanism {
		t.Errorf("Mechanism not preserved: got %q, want %q", restored.Mechanism, original.Mechanism)
	}
}

// TestValidationResult_BackwardCompatibility verifies that ValidationResult
// with no bypasses serializes the same as before (omits bypasses field).
func TestValidationResult_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name           string
		result         pricing.ValidationResult
		wantBypasses   bool
		wantFieldCount int
	}{
		{
			name: "no bypasses - nil",
			result: pricing.ValidationResult{
				Valid:    true,
				Bypasses: nil,
			},
			wantBypasses: false,
		},
		{
			name: "no bypasses - empty slice",
			result: pricing.ValidationResult{
				Valid:    true,
				Bypasses: []pricing.BypassMetadata{},
			},
			wantBypasses: false,
		},
		{
			name: "with errors, no bypasses",
			result: pricing.ValidationResult{
				Valid:  false,
				Errors: []string{"validation failed"},
			},
			wantBypasses: false,
		},
		{
			name: "with bypasses",
			result: pricing.ValidationResult{
				Valid: true,
				Bypasses: []pricing.BypassMetadata{
					pricing.NewBypassMetadata("test", "error"),
				},
			},
			wantBypasses: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("Failed to marshal ValidationResult: %v", err)
			}

			jsonStr := string(data)
			hasBypasses := strings.Contains(jsonStr, `"bypasses"`)

			if hasBypasses != tt.wantBypasses {
				t.Errorf("JSON bypasses field presence = %v, want %v\nJSON: %s", hasBypasses, tt.wantBypasses, jsonStr)
			}
		})
	}
}

// TestBypassMetadata_CrossServiceBoundary simulates bypass metadata
// crossing a service boundary (serialization → transmission → deserialization).
func TestBypassMetadata_CrossServiceBoundary(t *testing.T) {
	// Simulate service A creating a validation result with bypasses
	serviceAResult := pricing.ValidationResult{
		Valid: true,
		Bypasses: []pricing.BypassMetadata{
			pricing.NewBypassMetadata(
				"budget_limit",
				"Cost exceeds budget by $500",
				pricing.WithReason("Emergency deployment approved by manager"),
				pricing.WithOperator("user@example.com"),
				pricing.WithSeverity(pricing.BypassSeverityError),
				pricing.WithMechanism(pricing.BypassMechanismFlag),
			),
			pricing.NewBypassMetadata(
				"region_policy",
				"Region not in allowlist",
				pricing.WithReason("Temporary exception for DR test"),
				pricing.WithOperator("admin@example.com"),
				pricing.WithSeverity(pricing.BypassSeverityWarning),
				pricing.WithMechanism(pricing.BypassMechanismConfig),
			),
		},
	}

	// Serialize (simulating transmission)
	wireData, err := json.Marshal(serviceAResult)
	if err != nil {
		t.Fatalf("Service A: Failed to serialize: %v", err)
	}

	// Deserialize in Service B
	var serviceBResult pricing.ValidationResult
	if unmarshalErr := json.Unmarshal(wireData, &serviceBResult); unmarshalErr != nil {
		t.Fatalf("Service B: Failed to deserialize: %v", unmarshalErr)
	}

	// Verify all bypass metadata is preserved
	if len(serviceBResult.Bypasses) != len(serviceAResult.Bypasses) {
		t.Fatalf("Bypass count mismatch: got %d, want %d",
			len(serviceBResult.Bypasses), len(serviceAResult.Bypasses))
	}

	for i, original := range serviceAResult.Bypasses {
		restored := serviceBResult.Bypasses[i]

		if !restored.Timestamp.Equal(original.Timestamp) {
			t.Errorf("Bypass[%d].Timestamp not preserved", i)
		}
		if restored.ValidationName != original.ValidationName {
			t.Errorf("Bypass[%d].ValidationName: got %q, want %q",
				i, restored.ValidationName, original.ValidationName)
		}
		if restored.OriginalError != original.OriginalError {
			t.Errorf("Bypass[%d].OriginalError: got %q, want %q",
				i, restored.OriginalError, original.OriginalError)
		}
		if restored.Reason != original.Reason {
			t.Errorf("Bypass[%d].Reason: got %q, want %q",
				i, restored.Reason, original.Reason)
		}
		if restored.Operator != original.Operator {
			t.Errorf("Bypass[%d].Operator: got %q, want %q",
				i, restored.Operator, original.Operator)
		}
		if restored.Severity != original.Severity {
			t.Errorf("Bypass[%d].Severity: got %q, want %q",
				i, restored.Severity, original.Severity)
		}
		if restored.Mechanism != original.Mechanism {
			t.Errorf("Bypass[%d].Mechanism: got %q, want %q",
				i, restored.Mechanism, original.Mechanism)
		}
	}
}

// TestBypassMetadata_TimestampPrecision verifies that timestamps survive
// JSON round-trip with sufficient precision for audit purposes.
func TestBypassMetadata_TimestampPrecision(t *testing.T) {
	// Create with specific timestamp
	ts := time.Date(2026, 1, 24, 10, 30, 45, 123456789, time.UTC)
	original := pricing.BypassMetadata{
		Timestamp:      ts,
		ValidationName: "test",
		OriginalError:  "error",
		Severity:       pricing.BypassSeverityError,
		Mechanism:      pricing.BypassMechanismFlag,
	}

	data, marshalErr := json.Marshal(original)
	if marshalErr != nil {
		t.Fatalf("Marshal failed: %v", marshalErr)
	}

	var restored pricing.BypassMetadata
	if unmarshalErr := json.Unmarshal(data, &restored); unmarshalErr != nil {
		t.Fatalf("Unmarshal failed: %v", unmarshalErr)
	}

	// JSON timestamp has second precision (RFC3339)
	// Verify year, month, day, hour, minute, second are preserved
	if restored.Timestamp.Year() != original.Timestamp.Year() ||
		restored.Timestamp.Month() != original.Timestamp.Month() ||
		restored.Timestamp.Day() != original.Timestamp.Day() ||
		restored.Timestamp.Hour() != original.Timestamp.Hour() ||
		restored.Timestamp.Minute() != original.Timestamp.Minute() ||
		restored.Timestamp.Second() != original.Timestamp.Second() {
		t.Errorf("Timestamp precision lost: got %v, want %v (to second)",
			restored.Timestamp.Format(time.RFC3339),
			original.Timestamp.Format(time.RFC3339))
	}
}

// TestBypassSeverity_AllValuesValid verifies that all severity enum values
// pass validation.
func TestBypassSeverity_AllValuesValid(t *testing.T) {
	for _, severity := range pricing.AllBypassSeverities() {
		if !pricing.IsValidBypassSeverity(string(severity)) {
			t.Errorf("Severity %q should be valid", severity)
		}
	}
}

// TestBypassMechanism_AllValuesValid verifies that all mechanism enum values
// pass validation.
func TestBypassMechanism_AllValuesValid(t *testing.T) {
	for _, mechanism := range pricing.AllBypassMechanisms() {
		if !pricing.IsValidBypassMechanism(string(mechanism)) {
			t.Errorf("Mechanism %q should be valid", mechanism)
		}
	}
}
