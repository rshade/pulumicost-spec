package pbc_test

import (
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func TestRecommendationActionTypeValidation(t *testing.T) {
	tests := []struct {
		at       pbc.RecommendationActionType
		expected bool
	}{
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR, true},
		{pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER, true},
		{pbc.RecommendationActionType(12), false},
		{pbc.RecommendationActionType(-1), false},
	}

	for _, test := range tests {
		result := pbc.IsValidRecommendationActionType(test.at)
		if result != test.expected {
			t.Errorf("pbc.IsValidRecommendationActionType(%v) = %v, expected %v", test.at, result, test.expected)
		}
	}
}

func TestAllRecommendationActionTypes(t *testing.T) {
	expectedCount := 12
	all := pbc.AllRecommendationActionTypes()
	if len(all) != expectedCount {
		t.Errorf("pbc.AllRecommendationActionTypes() returned %d items, expected %d", len(all), expectedCount)
	}

	// Verify all items are unique and valid
	seen := make(map[pbc.RecommendationActionType]bool)
	for _, at := range all {
		if seen[at] {
			t.Errorf("Duplicate RecommendationActionType found in AllRecommendationActionTypes(): %v", at)
		}
		seen[at] = true
		if !pbc.IsValidRecommendationActionType(at) {
			t.Errorf("Invalid RecommendationActionType found in AllRecommendationActionTypes(): %v", at)
		}
	}
}

func BenchmarkIsValidRecommendationActionType(b *testing.B) {
	testCases := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
		pbc.RecommendationActionType(12),
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED,
	}
	b.ResetTimer()
	for i := range b.N {
		_ = pbc.IsValidRecommendationActionType(testCases[i%len(testCases)])
	}
}
