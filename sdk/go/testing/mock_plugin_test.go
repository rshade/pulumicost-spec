package testing_test

import (
	"testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	pktesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

// assertSummaryFields asserts scalar fields of RecommendationSummary using Get*() methods.
func assertSummaryFields(t *testing.T, result, expected *pbc.RecommendationSummary) {
	t.Helper()

	if result.GetTotalRecommendations() != expected.GetTotalRecommendations() {
		t.Errorf("TotalRecommendations = %d, want %d",
			result.GetTotalRecommendations(), expected.GetTotalRecommendations())
	}

	if result.GetProjectionPeriod() != expected.GetProjectionPeriod() {
		t.Errorf(
			"ProjectionPeriod = %s, want %s",
			result.GetProjectionPeriod(),
			expected.GetProjectionPeriod(),
		)
	}

	if result.GetTotalEstimatedSavings() != expected.GetTotalEstimatedSavings() {
		t.Errorf(
			"TotalEstimatedSavings = %f, want %f",
			result.GetTotalEstimatedSavings(),
			expected.GetTotalEstimatedSavings(),
		)
	}

	if result.GetCurrency() != expected.GetCurrency() {
		t.Errorf("Currency = %s, want %s", result.GetCurrency(), expected.GetCurrency())
	}
}

// assertMapFieldsInt32 asserts map[string]int32 fields.
func assertMapFieldsInt32(t *testing.T, name string, result, expected map[string]int32) {
	t.Helper()

	if len(result) != len(expected) {
		t.Errorf("%s length = %d, want %d", name, len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("%s[%s] = %d, want %d", name, k, result[k], v)
		}
	}
}

// assertMapFieldsFloat64 asserts map[string]float64 fields.
func assertMapFieldsFloat64(t *testing.T, name string, result, expected map[string]float64) {
	t.Helper()

	if len(result) != len(expected) {
		t.Errorf("%s length = %d, want %d", name, len(result), len(expected))
	}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("%s[%s] = %f, want %f", name, k, result[k], v)
		}
	}
}

func TestCalculateMockSummary(t *testing.T) {
	tests := []struct {
		name             string
		recommendations  []*pbc.Recommendation
		projectionPeriod string
		expected         *pbc.RecommendationSummary
	}{
		{
			name:             "empty recommendations",
			recommendations:  []*pbc.Recommendation{},
			projectionPeriod: "monthly",
			expected: &pbc.RecommendationSummary{
				TotalRecommendations:  0,
				CountByCategory:       make(map[string]int32),
				SavingsByCategory:     make(map[string]float64),
				CountByActionType:     make(map[string]int32),
				SavingsByActionType:   make(map[string]float64),
				ProjectionPeriod:      "monthly",
				TotalEstimatedSavings: 0,
				Currency:              "",
			},
		},
		{
			name: "single recommendation with savings",
			recommendations: []*pbc.Recommendation{
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 100.0,
						Currency:         "USD",
					},
				},
			},
			projectionPeriod: "monthly",
			expected: &pbc.RecommendationSummary{
				TotalRecommendations: 1,
				CountByCategory: map[string]int32{
					"RECOMMENDATION_CATEGORY_COST": 1,
				},
				SavingsByCategory: map[string]float64{
					"RECOMMENDATION_CATEGORY_COST": 100.0,
				},
				CountByActionType: map[string]int32{
					"RECOMMENDATION_ACTION_TYPE_MODIFY": 1,
				},
				SavingsByActionType: map[string]float64{
					"RECOMMENDATION_ACTION_TYPE_MODIFY": 100.0,
				},
				ProjectionPeriod:      "monthly",
				TotalEstimatedSavings: 100.0,
				Currency:              "USD",
			},
		},
		{
			name: "multiple recommendations with consistent currency",
			recommendations: []*pbc.Recommendation{
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 50.0,
						Currency:         "USD",
					},
				},
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_PERFORMANCE,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 75.0,
						Currency:         "USD",
					},
				},
			},
			projectionPeriod: "yearly",
			expected: &pbc.RecommendationSummary{
				TotalRecommendations: 2,
				CountByCategory: map[string]int32{
					"RECOMMENDATION_CATEGORY_COST":        1,
					"RECOMMENDATION_CATEGORY_PERFORMANCE": 1,
				},
				SavingsByCategory: map[string]float64{
					"RECOMMENDATION_CATEGORY_COST":        50.0,
					"RECOMMENDATION_CATEGORY_PERFORMANCE": 75.0,
				},
				CountByActionType: map[string]int32{
					"RECOMMENDATION_ACTION_TYPE_MODIFY":    1,
					"RECOMMENDATION_ACTION_TYPE_RIGHTSIZE": 1,
				},
				SavingsByActionType: map[string]float64{
					"RECOMMENDATION_ACTION_TYPE_MODIFY":    50.0,
					"RECOMMENDATION_ACTION_TYPE_RIGHTSIZE": 75.0,
				},
				ProjectionPeriod:      "yearly",
				TotalEstimatedSavings: 125.0,
				Currency:              "USD",
			},
		},
		{
			name: "mixed currencies clears currency field",
			recommendations: []*pbc.Recommendation{
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 50.0,
						Currency:         "USD",
					},
				},
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_PERFORMANCE,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 75.0,
						Currency:         "EUR",
					},
				},
			},
			projectionPeriod: "monthly",
			expected: &pbc.RecommendationSummary{
				TotalRecommendations: 2,
				CountByCategory: map[string]int32{
					"RECOMMENDATION_CATEGORY_COST":        1,
					"RECOMMENDATION_CATEGORY_PERFORMANCE": 1,
				},
				SavingsByCategory: map[string]float64{
					"RECOMMENDATION_CATEGORY_COST":        50.0,
					"RECOMMENDATION_CATEGORY_PERFORMANCE": 75.0,
				},
				CountByActionType: map[string]int32{
					"RECOMMENDATION_ACTION_TYPE_MODIFY":    1,
					"RECOMMENDATION_ACTION_TYPE_RIGHTSIZE": 1,
				},
				SavingsByActionType: map[string]float64{
					"RECOMMENDATION_ACTION_TYPE_MODIFY":    50.0,
					"RECOMMENDATION_ACTION_TYPE_RIGHTSIZE": 75.0,
				},
				ProjectionPeriod:      "monthly",
				TotalEstimatedSavings: 125.0,
				Currency:              "", // Should be empty due to currency mismatch
			},
		},
		{
			name: "recommendation without impact",
			recommendations: []*pbc.Recommendation{
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
					Impact:     nil,
				},
			},
			projectionPeriod: "monthly",
			expected: &pbc.RecommendationSummary{
				TotalRecommendations: 1,
				CountByCategory: map[string]int32{
					"RECOMMENDATION_CATEGORY_COST": 1,
				},
				SavingsByCategory: make(map[string]float64), // Empty when no impact
				CountByActionType: map[string]int32{
					"RECOMMENDATION_ACTION_TYPE_MODIFY": 1,
				},
				SavingsByActionType:   make(map[string]float64), // Empty when no impact
				ProjectionPeriod:      "monthly",
				TotalEstimatedSavings: 0,
				Currency:              "",
			},
		},
		{
			name: "multiple recommendations same category and action type",
			recommendations: []*pbc.Recommendation{
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 25.0,
						Currency:         "USD",
					},
				},
				{
					Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
					ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
					Impact: &pbc.RecommendationImpact{
						EstimatedSavings: 35.0,
						Currency:         "USD",
					},
				},
			},
			projectionPeriod: "monthly",
			expected: &pbc.RecommendationSummary{
				TotalRecommendations: 2,
				CountByCategory: map[string]int32{
					"RECOMMENDATION_CATEGORY_COST": 2,
				},
				SavingsByCategory: map[string]float64{
					"RECOMMENDATION_CATEGORY_COST": 60.0,
				},
				CountByActionType: map[string]int32{
					"RECOMMENDATION_ACTION_TYPE_MODIFY": 2,
				},
				SavingsByActionType: map[string]float64{
					"RECOMMENDATION_ACTION_TYPE_MODIFY": 60.0,
				},
				ProjectionPeriod:      "monthly",
				TotalEstimatedSavings: 60.0,
				Currency:              "USD",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pktesting.CalculateMockSummary(tt.recommendations, tt.projectionPeriod)

			assertSummaryFields(t, result, tt.expected)

			assertMapFieldsInt32(
				t,
				"CountByCategory",
				result.GetCountByCategory(),
				tt.expected.GetCountByCategory(),
			)
			assertMapFieldsFloat64(
				t,
				"SavingsByCategory",
				result.GetSavingsByCategory(),
				tt.expected.GetSavingsByCategory(),
			)
			assertMapFieldsInt32(
				t,
				"CountByActionType",
				result.GetCountByActionType(),
				tt.expected.GetCountByActionType(),
			)
			assertMapFieldsFloat64(
				t,
				"SavingsByActionType",
				result.GetSavingsByActionType(),
				tt.expected.GetSavingsByActionType(),
			)
		})
	}
}
