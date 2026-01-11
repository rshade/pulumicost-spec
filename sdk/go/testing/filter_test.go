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

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// RecommendationFilter Tests for New Action Types (T032)
// =============================================================================

// T032: Verify RecommendationFilter.action_type accepts new values.
func TestFilterActionTypeNewValues(t *testing.T) {
	// Create a diverse set of recommendations
	recs := []*pbc.Recommendation{
		{
			Id:         "rec-rightsize",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r1", Provider: "aws"},
		},
		{
			Id:         "rec-migrate",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r2", Provider: "azure"},
		},
		{
			Id:         "rec-consolidate",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r3", Provider: "kubernetes"},
		},
		{
			Id:         "rec-schedule",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r4", Provider: "aws"},
		},
		{
			Id:         "rec-refactor",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r5", Provider: "gcp"},
		},
		{
			Id:         "rec-other",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r6", Provider: "aws"},
		},
	}

	// Test filtering by each new action type
	newActionTypeTests := []struct {
		name          string
		actionType    pbc.RecommendationActionType
		expectedID    string
		expectedCount int
	}{
		{
			name:          "Filter by MIGRATE",
			actionType:    pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			expectedID:    "rec-migrate",
			expectedCount: 1,
		},
		{
			name:          "Filter by CONSOLIDATE",
			actionType:    pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
			expectedID:    "rec-consolidate",
			expectedCount: 1,
		},
		{
			name:          "Filter by SCHEDULE",
			actionType:    pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
			expectedID:    "rec-schedule",
			expectedCount: 1,
		},
		{
			name:          "Filter by REFACTOR",
			actionType:    pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
			expectedID:    "rec-refactor",
			expectedCount: 1,
		},
		{
			name:          "Filter by OTHER",
			actionType:    pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
			expectedID:    "rec-other",
			expectedCount: 1,
		},
	}

	for _, tc := range newActionTypeTests {
		t.Run(tc.name, func(t *testing.T) {
			filter := &pbc.RecommendationFilter{
				ActionType: tc.actionType,
			}

			result := pluginsdk.ApplyRecommendationFilter(recs, filter)

			if len(result) != tc.expectedCount {
				t.Errorf("Filter returned %d recommendations, want %d",
					len(result), tc.expectedCount)
				return
			}

			if result[0].GetId() != tc.expectedID {
				t.Errorf("Filtered recommendation ID = %q, want %q",
					result[0].GetId(), tc.expectedID)
			}

			if result[0].GetActionType() != tc.actionType {
				t.Errorf("Filtered recommendation ActionType = %v, want %v",
					result[0].GetActionType(), tc.actionType)
			}
		})
	}
}

// TestFilterActionTypeUnspecifiedReturnsAll verifies that UNSPECIFIED doesn't filter.
func TestFilterActionTypeUnspecifiedReturnsAll(t *testing.T) {
	recs := []*pbc.Recommendation{
		{
			Id:         "rec-1",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r1", Provider: "aws"},
		},
		{
			Id:         "rec-2",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r2", Provider: "aws"},
		},
	}

	// UNSPECIFIED should not filter - return all
	filter := &pbc.RecommendationFilter{
		ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED,
	}

	result := pluginsdk.ApplyRecommendationFilter(recs, filter)

	if len(result) != len(recs) {
		t.Errorf("UNSPECIFIED filter returned %d recommendations, want %d (all)",
			len(result), len(recs))
	}
}

// TestFilterCombinedWithNewActionType tests combining action type with other filters.
func TestFilterCombinedWithNewActionType(t *testing.T) {
	recs := []*pbc.Recommendation{
		{
			Id:         "rec-migrate-aws",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r1", Provider: "aws", Region: "us-east-1"},
		},
		{
			Id:         "rec-migrate-azure",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r2", Provider: "azure", Region: "eastus"},
		},
		{
			Id:         "rec-schedule-aws",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r3", Provider: "aws", Region: "us-east-1"},
		},
	}

	// Filter by MIGRATE + aws provider
	filter := &pbc.RecommendationFilter{
		ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		Provider:   "aws",
	}

	result := pluginsdk.ApplyRecommendationFilter(recs, filter)

	if len(result) != 1 {
		t.Fatalf("Combined filter returned %d recommendations, want 1", len(result))
	}

	if result[0].GetId() != "rec-migrate-aws" {
		t.Errorf("Combined filter returned ID = %q, want %q",
			result[0].GetId(), "rec-migrate-aws")
	}
}

// TestFilterNoMatchForNewActionType verifies empty result when no matches.
func TestFilterNoMatchForNewActionType(t *testing.T) {
	// All recommendations use old action types
	recs := []*pbc.Recommendation{
		{
			Id:         "rec-rightsize",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r1", Provider: "aws"},
		},
		{
			Id:         "rec-terminate",
			ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE,
			Resource:   &pbc.ResourceRecommendationInfo{Id: "r2", Provider: "aws"},
		},
	}

	// Filter by new action type that doesn't exist in the data
	filter := &pbc.RecommendationFilter{
		ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
	}

	result := pluginsdk.ApplyRecommendationFilter(recs, filter)

	if len(result) != 0 {
		t.Errorf("Filter for non-existent action type returned %d recommendations, want 0",
			len(result))
	}
}

// TestFilterValidationAcceptsNewActionTypes verifies filter validation works with new types.
func TestFilterValidationAcceptsNewActionTypes(t *testing.T) {
	newActionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
	}

	for _, actionType := range newActionTypes {
		t.Run(actionType.String(), func(t *testing.T) {
			filter := &pbc.RecommendationFilter{
				ActionType: actionType,
			}

			// Validate should not error for new action types
			err := pluginsdk.ValidateRecommendationFilter(filter)
			if err != nil {
				t.Errorf("ValidateRecommendationFilter failed for %v: %v", actionType, err)
			}
		})
	}
}
