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
	"fmt"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// RecommendationSummary Tests for New Action Types (T031)
// =============================================================================

// T031: Verify RecommendationSummary.count_by_action_type handles new values correctly.
func TestSummaryCountByActionTypeNewValues(t *testing.T) {
	// Create recommendations using all 12 action types
	allActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	recs := make([]*pbc.Recommendation, len(allActionTypes))
	for i, actionType := range allActionTypes {
		recs[i] = &pbc.Recommendation{
			Id:         actionType.String(),
			Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType: actionType,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "res", Provider: "aws"},
			Impact: &pbc.RecommendationImpact{
				EstimatedSavings: float64((i + 1) * 100),
				Currency:         "USD",
			},
		}
	}

	// Calculate summary
	summary := pluginsdk.CalculateRecommendationSummary(recs, "monthly")

	// Verify total count
	if summary.GetTotalRecommendations() != int32(len(allActionTypes)) {
		t.Errorf("TotalRecommendations = %d, want %d",
			summary.GetTotalRecommendations(), len(allActionTypes))
	}

	// Verify each action type is counted
	countByActionType := summary.GetCountByActionType()
	for _, actionType := range allActionTypes {
		key := actionType.String()
		count, ok := countByActionType[key]
		if !ok {
			t.Errorf("CountByActionType missing key %q", key)
			continue
		}
		if count != 1 {
			t.Errorf("CountByActionType[%q] = %d, want 1", key, count)
		}
	}

	// Verify new action types specifically
	newActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	for _, actionType := range newActionTypes {
		key := actionType.String()
		count, ok := countByActionType[key]
		if !ok {
			t.Errorf("New action type %q not in CountByActionType", key)
		}
		if count != 1 {
			t.Errorf("CountByActionType[%q] = %d, want 1", key, count)
		}
	}
}

// TestSummarySavingsByActionTypeNewValues verifies savings aggregation for new action types.
func TestSummarySavingsByActionTypeNewValues(t *testing.T) {
	// Create multiple recommendations per new action type with known savings
	testData := []struct {
		actionType pbc.RecommendationActionType
		savings    []float64
	}{
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE, []float64{100.0, 200.0}},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE, []float64{150.0}},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE, []float64{50.0, 75.0, 125.0}},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR, []float64{500.0}},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER, []float64{25.0, 25.0}},
	}

	var recs []*pbc.Recommendation
	expectedSavings := make(map[string]float64)

	for _, td := range testData {
		key := td.actionType.String()
		var totalForType float64
		for i, savings := range td.savings {
			recs = append(recs, &pbc.Recommendation{
				Id:         fmt.Sprintf("%s-%d", td.actionType.String(), i),
				Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
				ActionType: td.actionType,
				Resource:   &pbc.ResourceRecommendationInfo{Id: "res", Provider: "aws"},
				Impact: &pbc.RecommendationImpact{
					EstimatedSavings: savings,
					Currency:         "USD",
				},
			})
			totalForType += savings
		}
		expectedSavings[key] = totalForType
	}

	// Calculate summary
	summary := pluginsdk.CalculateRecommendationSummary(recs, "monthly")

	// Verify savings by action type
	savingsByActionType := summary.GetSavingsByActionType()
	for actionTypeKey, expected := range expectedSavings {
		actual, ok := savingsByActionType[actionTypeKey]
		if !ok {
			t.Errorf("SavingsByActionType missing key %q", actionTypeKey)
			continue
		}
		if actual != expected {
			t.Errorf("SavingsByActionType[%q] = %.2f, want %.2f", actionTypeKey, actual, expected)
		}
	}

	// Verify total savings
	var expectedTotal float64
	for _, savings := range expectedSavings {
		expectedTotal += savings
	}
	if summary.GetTotalEstimatedSavings() != expectedTotal {
		t.Errorf("TotalEstimatedSavings = %.2f, want %.2f",
			summary.GetTotalEstimatedSavings(), expectedTotal)
	}
}

// TestSummaryEmptyForNewActionTypes verifies empty results for action types with no recommendations.
func TestSummaryEmptyForNewActionTypes(t *testing.T) {
	// Create recommendations with only OLD action types
	oldRecs := []*pbc.Recommendation{
		{
			Id:         "old-1",
			Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "res", Provider: "aws"},
			Impact:     &pbc.RecommendationImpact{EstimatedSavings: 100.0, Currency: "USD"},
		},
	}

	summary := pluginsdk.CalculateRecommendationSummary(oldRecs, "monthly")

	// New action types should not appear in the summary (no zero-fill)
	newActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	countByActionType := summary.GetCountByActionType()
	savingsByActionType := summary.GetSavingsByActionType()

	for _, actionType := range newActionTypes {
		key := actionType.String()
		if count, ok := countByActionType[key]; ok && count != 0 {
			t.Errorf("CountByActionType[%q] = %d, want 0 or not present", key, count)
		}
		if savings, ok := savingsByActionType[key]; ok && savings != 0 {
			t.Errorf("SavingsByActionType[%q] = %.2f, want 0 or not present", key, savings)
		}
	}
}
