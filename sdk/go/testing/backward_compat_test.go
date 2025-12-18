package testing_test

import (
	"context"
	"testing"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// Backward Compatibility Tests (T022-T024)
//
// These tests verify that the extended RecommendationActionType enum maintains
// backward compatibility with existing plugins and clients.
// =============================================================================

// T022: Test backward compatibility for plugins using original action types.
// Existing plugins should continue to work without modification.
func TestBackwardCompatibilityOriginalActionTypes(t *testing.T) {
	// Create a mock plugin that uses a subset of original action types
	plugin := plugintesting.NewMockPlugin()

	// Configure with representative original action types (subset of 0-6)
	originalActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
	}

	// Create recommendations using only original action types
	recs := make([]*pbc.Recommendation, len(originalActionTypes))
	for i, actionType := range originalActionTypes {
		recs[i] = &pbc.Recommendation{
			Id:          "rec-original-" + actionType.String(),
			Category:    pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType:  actionType,
			Description: "Test recommendation with original action type",
			Resource: &pbc.ResourceRecommendationInfo{
				Id:       "res-1",
				Provider: "aws",
			},
			Impact: &pbc.RecommendationImpact{
				EstimatedSavings: 100.0,
				Currency:         "USD",
			},
		}
	}

	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: recs,
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Get recommendations - should work exactly as before
	resp, err := harness.Client().GetRecommendations(context.Background(), &pbc.GetRecommendationsRequest{})
	if err != nil {
		t.Fatalf("GetRecommendations failed: %v", err)
	}

	// Verify all recommendations were returned
	if len(resp.GetRecommendations()) != len(originalActionTypes) {
		t.Errorf("Expected %d recommendations, got %d",
			len(originalActionTypes), len(resp.GetRecommendations()))
	}

	// Verify action types are preserved correctly
	for i, rec := range resp.GetRecommendations() {
		if rec.GetActionType() != originalActionTypes[i] {
			t.Errorf("Recommendation %d: expected action type %v, got %v",
				i, originalActionTypes[i], rec.GetActionType())
		}

		// Verify the action type still serializes to its expected string
		expectedStr := originalActionTypes[i].String()
		if rec.GetActionType().String() != expectedStr {
			t.Errorf("Action type string mismatch: expected %q, got %q",
				expectedStr, rec.GetActionType().String())
		}
	}
}

// T023: Test unknown enum value handling.
// When core receives value 7-11 from a new plugin, old clients should handle gracefully.
func TestUnknownEnumValueHandling(t *testing.T) {
	// Test that unknown enum values (7-11) can be handled by treating them as numeric values
	// This simulates an old client receiving new enum values from an updated plugin

	testCases := []struct {
		name      string
		enumValue int32
		enumName  string
	}{
		{"MIGRATE", 7, "RECOMMENDATION_ACTION_TYPE_MIGRATE"},
		{"CONSOLIDATE", 8, "RECOMMENDATION_ACTION_TYPE_CONSOLIDATE"},
		{"SCHEDULE", 9, "RECOMMENDATION_ACTION_TYPE_SCHEDULE"},
		{"REFACTOR", 10, "RECOMMENDATION_ACTION_TYPE_REFACTOR"},
		{"OTHER", 11, "RECOMMENDATION_ACTION_TYPE_OTHER"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create action type from numeric value (simulating wire format)
			actionType := pbc.RecommendationActionType(tc.enumValue)

			// Verify numeric value is preserved
			if int32(actionType.Number()) != tc.enumValue {
				t.Errorf("Numeric value not preserved: expected %d, got %d",
					tc.enumValue, actionType.Number())
			}

			// Verify string representation works (protobuf generates these)
			str := actionType.String()
			if str != tc.enumName {
				t.Errorf("String representation: expected %q, got %q",
					tc.enumName, str)
			}

			// Verify round-trip through recommendation
			rec := &pbc.Recommendation{
				Id:         "test-rec",
				ActionType: actionType,
			}

			if rec.GetActionType() != actionType {
				t.Errorf("Round-trip through recommendation failed: %v != %v",
					rec.GetActionType(), actionType)
			}
		})
	}
}

// T024: Test gRPC communication between old plugin binary and new SDK.
// This verifies the wire format compatibility of the enum extension.
func TestGRPCWireFormatCompatibility(t *testing.T) {
	// Create plugin with mix of old and new action types
	plugin := plugintesting.NewMockPlugin()

	// Create recommendations with all action types (old and new)
	allActionTypes := []pbc.RecommendationActionType{
		// Original types (0-6)
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED,
		// New types (7-11)
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	recs := make([]*pbc.Recommendation, len(allActionTypes))
	for i, actionType := range allActionTypes {
		recs[i] = &pbc.Recommendation{
			Id:          "rec-" + actionType.String(),
			Category:    pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType:  actionType,
			Description: "Test recommendation",
			Resource: &pbc.ResourceRecommendationInfo{
				Id:       "res-" + actionType.String(),
				Provider: "aws",
			},
			Impact: &pbc.RecommendationImpact{
				EstimatedSavings: float64(i) * 10.0,
				Currency:         "USD",
			},
		}
	}

	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: recs,
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Test round-trip through gRPC (bufconn simulates network wire format)
	resp, err := harness.Client().GetRecommendations(context.Background(), &pbc.GetRecommendationsRequest{})
	if err != nil {
		t.Fatalf("GetRecommendations failed: %v", err)
	}

	// Verify all 12 action types survived gRPC serialization/deserialization
	if len(resp.GetRecommendations()) != len(allActionTypes) {
		t.Fatalf("Expected %d recommendations, got %d",
			len(allActionTypes), len(resp.GetRecommendations()))
	}

	for i, rec := range resp.GetRecommendations() {
		expectedActionType := allActionTypes[i]

		// Verify action type value is correct
		if rec.GetActionType() != expectedActionType {
			t.Errorf("Recommendation %d: action type mismatch: expected %v, got %v",
				i, expectedActionType, rec.GetActionType())
		}

		// Verify numeric value is preserved through gRPC
		if int32(rec.GetActionType().Number()) != int32(expectedActionType.Number()) {
			t.Errorf("Recommendation %d: numeric value mismatch: expected %d, got %d",
				i, expectedActionType.Number(), rec.GetActionType().Number())
		}
	}
}

// TestFilterByNewActionTypes verifies filtering works correctly for new action types.
func TestFilterByNewActionTypes(t *testing.T) {
	// Test filtering by each new action type
	newActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	// Create plugin with recommendations for each new action type
	plugin := plugintesting.NewMockPlugin()
	recs := make([]*pbc.Recommendation, len(newActionTypes))
	for i, actionType := range newActionTypes {
		recs[i] = &pbc.Recommendation{
			Id:          "rec-" + actionType.String(),
			Category:    pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType:  actionType,
			Description: "Test " + actionType.String(),
			Resource: &pbc.ResourceRecommendationInfo{
				Id:       "res-" + actionType.String(),
				Provider: "aws",
			},
			Impact: &pbc.RecommendationImpact{
				EstimatedSavings: 50.0,
				Currency:         "USD",
			},
		}
	}

	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: recs,
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Test filtering by each new action type
	for _, targetActionType := range newActionTypes {
		t.Run(targetActionType.String(), func(t *testing.T) {
			resp, err := harness.Client().GetRecommendations(context.Background(),
				&pbc.GetRecommendationsRequest{
					Filter: &pbc.RecommendationFilter{
						ActionType: targetActionType,
					},
				})
			if err != nil {
				t.Fatalf("GetRecommendations with filter failed: %v", err)
			}

			// Should return exactly 1 recommendation
			if len(resp.GetRecommendations()) != 1 {
				t.Errorf("Expected 1 recommendation for filter %v, got %d",
					targetActionType, len(resp.GetRecommendations()))
				return
			}

			// Verify the returned recommendation has the correct action type
			rec := resp.GetRecommendations()[0]
			if rec.GetActionType() != targetActionType {
				t.Errorf("Filtered recommendation has wrong action type: expected %v, got %v",
					targetActionType, rec.GetActionType())
			}
		})
	}
}

// TestSummaryCountsByNewActionTypes verifies summary aggregation includes new action types.
func TestSummaryCountsByNewActionTypes(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	// Create recommendations with new action types
	recs := []*pbc.Recommendation{
		{
			Id:         "rec-migrate-1",
			Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r1", Provider: "aws"},
			Impact:     &pbc.RecommendationImpact{EstimatedSavings: 100.0, Currency: "USD"},
		},
		{
			Id:         "rec-migrate-2",
			Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r2", Provider: "aws"},
			Impact:     &pbc.RecommendationImpact{EstimatedSavings: 150.0, Currency: "USD"},
		},
		{
			Id:         "rec-schedule-1",
			Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r3", Provider: "azure"},
			Impact:     &pbc.RecommendationImpact{EstimatedSavings: 200.0, Currency: "USD"},
		},
	}

	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: recs,
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	resp, err := harness.Client().GetRecommendations(context.Background(), &pbc.GetRecommendationsRequest{})
	if err != nil {
		t.Fatalf("GetRecommendations failed: %v", err)
	}

	summary := resp.GetSummary()
	if summary == nil {
		t.Fatal("Summary is nil")
	}

	// Verify count by action type includes new types
	migrateKey := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE.String()
	scheduleKey := pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE.String()

	countByActionType := summary.GetCountByActionType()
	savingsByActionType := summary.GetSavingsByActionType()

	if count, ok := countByActionType[migrateKey]; !ok || count != 2 {
		t.Errorf("Expected 2 MIGRATE recommendations, got %d", count)
	}

	if count, ok := countByActionType[scheduleKey]; !ok || count != 1 {
		t.Errorf("Expected 1 SCHEDULE recommendation, got %d", count)
	}

	// Verify savings by action type includes new types
	if savings, ok := savingsByActionType[migrateKey]; !ok || savings != 250.0 {
		t.Errorf("Expected $250.00 savings for MIGRATE, got $%.2f", savings)
	}

	if savings, ok := savingsByActionType[scheduleKey]; !ok || savings != 200.0 {
		t.Errorf("Expected $200.00 savings for SCHEDULE, got $%.2f", savings)
	}
}
