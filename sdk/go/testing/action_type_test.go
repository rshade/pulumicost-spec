// Copyright 2026 PulumiCost/FinFocus Authors
//
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
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// Extended RecommendationActionType Enum Tests (T005-T010)
//
// These tests validate the 5 new action types added in feature 019:
// - MIGRATE (7)
// - CONSOLIDATE (8)
// - SCHEDULE (9)
// - REFACTOR (10)
// - OTHER (11)
//
// Per Constitution III (Test-First Protocol), these tests MUST be written
// and fail BEFORE the proto changes are implemented.
// =============================================================================

// T005: Test MIGRATE action type serialization.
// MIGRATE is for moving workloads to different regions/zones/SKUs.
func TestActionTypeMigrateSerialization(t *testing.T) {
	// This test will fail until RECOMMENDATION_ACTION_TYPE_MIGRATE is added to proto
	actionType := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE

	// Verify enum value is 7
	if actionType.Number() != 7 {
		t.Errorf("MIGRATE should have value 7, got %d", actionType.Number())
	}

	// Verify string representation
	expected := "RECOMMENDATION_ACTION_TYPE_MIGRATE"
	if actionType.String() != expected {
		t.Errorf("MIGRATE string should be %q, got %q", expected, actionType.String())
	}

	// Verify round-trip through enum value
	roundTripped := pbc.RecommendationActionType(7)
	if roundTripped != actionType {
		t.Errorf("Round-trip failed: %v != %v", roundTripped, actionType)
	}
}

// T006: Test CONSOLIDATE action type serialization.
// CONSOLIDATE is for combining multiple resources into fewer, larger ones.
func TestActionTypeConsolidateSerialization(t *testing.T) {
	// This test will fail until RECOMMENDATION_ACTION_TYPE_CONSOLIDATE is added to proto
	actionType := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE

	// Verify enum value is 8
	if actionType.Number() != 8 {
		t.Errorf("CONSOLIDATE should have value 8, got %d", actionType.Number())
	}

	// Verify string representation
	expected := "RECOMMENDATION_ACTION_TYPE_CONSOLIDATE"
	if actionType.String() != expected {
		t.Errorf("CONSOLIDATE string should be %q, got %q", expected, actionType.String())
	}

	// Verify round-trip through enum value
	roundTripped := pbc.RecommendationActionType(8)
	if roundTripped != actionType {
		t.Errorf("Round-trip failed: %v != %v", roundTripped, actionType)
	}
}

// T007: Test SCHEDULE action type serialization.
// SCHEDULE is for start/stop resources on schedule (dev/test environments).
func TestActionTypeScheduleSerialization(t *testing.T) {
	// This test will fail until RECOMMENDATION_ACTION_TYPE_SCHEDULE is added to proto
	actionType := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE

	// Verify enum value is 9
	if actionType.Number() != 9 {
		t.Errorf("SCHEDULE should have value 9, got %d", actionType.Number())
	}

	// Verify string representation
	expected := "RECOMMENDATION_ACTION_TYPE_SCHEDULE"
	if actionType.String() != expected {
		t.Errorf("SCHEDULE string should be %q, got %q", expected, actionType.String())
	}

	// Verify round-trip through enum value
	roundTripped := pbc.RecommendationActionType(9)
	if roundTripped != actionType {
		t.Errorf("Round-trip failed: %v != %v", roundTripped, actionType)
	}
}

// T008: Test REFACTOR action type serialization.
// REFACTOR is for architectural changes (e.g., move to serverless).
func TestActionTypeRefactorSerialization(t *testing.T) {
	// This test will fail until RECOMMENDATION_ACTION_TYPE_REFACTOR is added to proto
	actionType := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR

	// Verify enum value is 10
	if actionType.Number() != 10 {
		t.Errorf("REFACTOR should have value 10, got %d", actionType.Number())
	}

	// Verify string representation
	expected := "RECOMMENDATION_ACTION_TYPE_REFACTOR"
	if actionType.String() != expected {
		t.Errorf("REFACTOR string should be %q, got %q", expected, actionType.String())
	}

	// Verify round-trip through enum value
	roundTripped := pbc.RecommendationActionType(10)
	if roundTripped != actionType {
		t.Errorf("Round-trip failed: %v != %v", roundTripped, actionType)
	}
}

// T009: Test OTHER action type serialization.
// OTHER is for provider-specific recommendations not fitting other categories.
func TestActionTypeOtherSerialization(t *testing.T) {
	// This test will fail until RECOMMENDATION_ACTION_TYPE_OTHER is added to proto
	actionType := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER

	// Verify enum value is 11
	if actionType.Number() != 11 {
		t.Errorf("OTHER should have value 11, got %d", actionType.Number())
	}

	// Verify string representation
	expected := "RECOMMENDATION_ACTION_TYPE_OTHER"
	if actionType.String() != expected {
		t.Errorf("OTHER string should be %q, got %q", expected, actionType.String())
	}

	// Verify round-trip through enum value
	roundTripped := pbc.RecommendationActionType(11)
	if roundTripped != actionType {
		t.Errorf("Round-trip failed: %v != %v", roundTripped, actionType)
	}
}

// T010: Test round-trip serialization of all 12 action type enum values.
// This validates the complete enum from UNSPECIFIED (0) through OTHER (11).
func TestActionTypeAllValuesRoundTrip(t *testing.T) {
	testCases := []struct {
		name     string
		value    int32
		enumName string
	}{
		{"UNSPECIFIED", 0, "RECOMMENDATION_ACTION_TYPE_UNSPECIFIED"},
		{"RIGHTSIZE", 1, "RECOMMENDATION_ACTION_TYPE_RIGHTSIZE"},
		{"TERMINATE", 2, "RECOMMENDATION_ACTION_TYPE_TERMINATE"},
		{"PURCHASE_COMMITMENT", 3, "RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT"},
		{"ADJUST_REQUESTS", 4, "RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS"},
		{"MODIFY", 5, "RECOMMENDATION_ACTION_TYPE_MODIFY"},
		{"DELETE_UNUSED", 6, "RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED"},
		// New action types (T005-T009) - will fail until proto changes
		{"MIGRATE", 7, "RECOMMENDATION_ACTION_TYPE_MIGRATE"},
		{"CONSOLIDATE", 8, "RECOMMENDATION_ACTION_TYPE_CONSOLIDATE"},
		{"SCHEDULE", 9, "RECOMMENDATION_ACTION_TYPE_SCHEDULE"},
		{"REFACTOR", 10, "RECOMMENDATION_ACTION_TYPE_REFACTOR"},
		{"OTHER", 11, "RECOMMENDATION_ACTION_TYPE_OTHER"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create enum from numeric value
			actionType := pbc.RecommendationActionType(tc.value)

			// Verify string representation matches expected name
			if actionType.String() != tc.enumName {
				t.Errorf("RecommendationActionType(%d).String() = %q, want %q",
					tc.value, actionType.String(), tc.enumName)
			}

			// Verify numeric value round-trips correctly
			if int32(actionType.Number()) != tc.value {
				t.Errorf("RecommendationActionType(%d).Number() = %d, want %d",
					tc.value, actionType.Number(), tc.value)
			}
		})
	}
}

// TestActionTypeEnumCount validates the expected total count of action types.
// After feature 019, there should be 12 action types (0-11).
func TestActionTypeEnumCount(t *testing.T) {
	// Count known enum values by checking if they have proper string names.
	// Known enum values have names like "RECOMMENDATION_ACTION_TYPE_MIGRATE".
	// Unknown values return just their numeric value as a string (e.g., "12").
	knownCount := 0

	for i := range int32(20) { // Check up to 19 to catch any extras
		actionType := pbc.RecommendationActionType(i)
		str := actionType.String()

		// Known enum values start with "RECOMMENDATION_ACTION_TYPE_".
		// Unknown values are just numeric strings (e.g., "12", "15").
		if isKnownActionType(str) {
			knownCount++
		}
	}

	// After proto changes, we expect 12 action types (0-11)
	expectedCount := 12
	if knownCount != expectedCount {
		t.Errorf("Expected %d action types, found %d", expectedCount, knownCount)
	}
}

// isKnownActionType checks if an enum string represents a defined action type.
// Known action types have the prefix "RECOMMENDATION_ACTION_TYPE_".
// Unknown values are represented as just their numeric value (e.g., "12").
func isKnownActionType(s string) bool {
	const prefix = "RECOMMENDATION_ACTION_TYPE_"
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// TestActionTypeInRecommendation validates that new action types can be used
// in Recommendation messages correctly.
func TestActionTypeInRecommendation(t *testing.T) {
	// Test that new action types can be assigned to Recommendation.ActionType
	// This will compile-fail if the enum constants don't exist

	testCases := []struct {
		name       string
		actionType pbc.RecommendationActionType
		wantValue  int32
	}{
		{
			name:       "MIGRATE in recommendation",
			actionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			wantValue:  7,
		},
		{
			name:       "CONSOLIDATE in recommendation",
			actionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
			wantValue:  8,
		},
		{
			name:       "SCHEDULE in recommendation",
			actionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
			wantValue:  9,
		},
		{
			name:       "REFACTOR in recommendation",
			actionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
			wantValue:  10,
		},
		{
			name:       "OTHER in recommendation",
			actionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
			wantValue:  11,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a recommendation with the new action type
			rec := &pbc.Recommendation{
				Id:         "test-rec-" + tc.name,
				ActionType: tc.actionType,
			}

			// Verify the action type was set correctly
			if rec.GetActionType() != tc.actionType {
				t.Errorf("Recommendation.ActionType = %v, want %v",
					rec.GetActionType(), tc.actionType)
			}

			// Verify the numeric value
			if int32(rec.GetActionType().Number()) != tc.wantValue {
				t.Errorf("Recommendation.ActionType.Number() = %d, want %d",
					rec.GetActionType().Number(), tc.wantValue)
			}
		})
	}
}

// TestActionTypeInRecommendationFilter validates that new action types work
// correctly in RecommendationFilter for filtering operations.
func TestActionTypeInRecommendationFilter(t *testing.T) {
	// Test that new action types can be used in RecommendationFilter.ActionType

	newActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	for _, actionType := range newActionTypes {
		t.Run(actionType.String(), func(t *testing.T) {
			// Create a filter with the new action type
			filter := &pbc.RecommendationFilter{
				ActionType: actionType,
			}

			// Verify the filter was set correctly
			if filter.GetActionType() != actionType {
				t.Errorf("RecommendationFilter.ActionType = %v, want %v",
					filter.GetActionType(), actionType)
			}
		})
	}
}
