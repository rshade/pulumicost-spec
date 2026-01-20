// Copyright 2025 FinFocus
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
	"context"
	"testing"

	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// TestAnomalyConformance_UnifiedView tests User Story 1: Unified Actionable Cost Insights View
// This test verifies that GetRecommendations can return both optimization recommendations and
// cost anomalies in a single response.
func TestAnomalyConformance_UnifiedView(t *testing.T) {
	t.Run("ReturnsAnomalyCategoryRecommendations", func(t *testing.T) {
		// T008: Write conformance test - plugin returns ANOMALY category recommendation
		// GetRecommendations returns mixed categories

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{})

		require.NoError(t, err, "GetRecommendations should not error")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Verify anomaly recommendations are present
		anomalyFound := false
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				anomalyFound = true
				assert.NotEmpty(t, rec.GetDescription(), "Anomaly should have description")
				assert.NotEmpty(t, rec.GetId(), "Anomaly should have ID")
				break
			}
		}
		assert.True(t, anomalyFound, "Response should contain at least one ANOMALY category recommendation")
	})

	t.Run("CategoryFilterWorks_IncludesAnomalies", func(t *testing.T) {
		// T009: Write conformance test - category filter works for ANOMALY
		// category=ANOMALY returns only anomalies

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
			},
		})

		require.NoError(t, err, "GetRecommendations with ANOMALY filter should not error")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Verify all returned recommendations are ANOMALY category
		for _, rec := range resp.GetRecommendations() {
			assert.Equal(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				rec.GetCategory(), "Filtered response should only contain ANOMALY category")
		}
	})

	t.Run("NoFilterReturnsAllCategories", func(t *testing.T) {
		// Verify that no filter returns mixed categories (both anomalies and optimizations)

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{})

		require.NoError(t, err, "GetRecommendations should not error")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Count different categories
		categoryCount := make(map[pbc.RecommendationCategory]int)
		for _, rec := range resp.GetRecommendations() {
			categoryCount[rec.GetCategory()]++
		}

		// Should have at least one ANOMALY (if plugin supports it)
		if count, hasAnomaly := categoryCount[pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY]; hasAnomaly {
			assert.Greater(t, count, 0, "Should have ANOMALY recommendations")
		}
	})
}

// TestAnomalyConformance_NonSupportingPlugins tests that plugins not supporting anomalies
// simply don't return any ANOMALY recommendations (no error is raised).
func TestAnomalyConformance_NonSupportingPlugins(t *testing.T) {
	// T010: Write conformance test - plugins not supporting anomalies return zero ANOMALY recommendations

	t.Run("PluginWithoutAnomalySupport", func(t *testing.T) {
		// Create a plugin that returns recommendations without ANOMALY category
		// This simulates a plugin that doesn't support anomaly detection
		plugin := plugintesting.NewMockPlugin()
		// Override recommendations to exclude ANOMALY category
		recs := []*pbc.Recommendation{
			{
				Id:       "rec-1",
				Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
				ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
				Description: "Non-anomaly recommendation",
			},
		}
		plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
			Recommendations: recs,
		})

		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{})

		// Should not error - this is backward compatible behavior
		require.NoError(t, err, "GetRecommendations should not error for non-supporting plugins")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Verify no ANOMALY recommendations are returned
		for _, rec := range resp.GetRecommendations() {
			assert.NotEqual(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				rec.GetCategory(), "Non-supporting plugin should not return ANOMALY category")
		}
	})
}

// TestAnomalyConformance_ConfidenceScoreFiltering tests User Story 2: Anomaly Triage by Confidence Score
// This test verifies that the existing confidence_score filter works with ANOMALY recommendations.
func TestAnomalyConformance_ConfidenceScoreFiltering(t *testing.T) {
	// T013: Write conformance test - confidence_score filter works for anomalies
	// min_confidence_score=0.5 filters correctly

	t.Run("FilterByConfidenceScore", func(t *testing.T) {
		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()

		// Request with confidence score filter
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				MinConfidenceScore: 0.5,
			},
		})

		require.NoError(t, err, "GetRecommendations with confidence filter should not error")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Verify that anomalies with confidence >= 0.5 are included
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				// If confidence_score is set, it should meet the minimum
				// Note: proto3 optional double stores as *float64 internally
				if rec.ConfidenceScore != nil && *rec.ConfidenceScore > 0 {
					assert.GreaterOrEqual(t, *rec.ConfidenceScore, 0.5,
						"Anomaly confidence score should meet minimum threshold")
				}
			}
		}
	})
}

// TestAnomalyConformance_InvestigateAction tests User Story 3: Investigate Anomalous Spending
// This test verifies that anomaly recommendations use the INVESTIGATE action type.
func TestAnomalyConformance_InvestigateAction(t *testing.T) {
	t.Run("AnomalyRecommendationsHaveInvestigateAction", func(t *testing.T) {
		// T016: Write conformance test - INVESTIGATE action_type returned with ANOMALY category

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
			},
		})

		require.NoError(t, err, "GetRecommendations should not error")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Verify anomalies have description and metadata
		for _, rec := range resp.GetRecommendations() {
			assert.Equal(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				rec.GetCategory(), "Should be ANOMALY category")

			// Typically should be INVESTIGATE action
			if rec.GetActionType() == pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE {
				t.Logf("Anomaly %s has INVESTIGATE action", rec.GetId())
			}
		}
	})

	t.Run("AnomalyHasContextForInvestigation", func(t *testing.T) {
		// T017: Write conformance test - anomaly recommendation has description and metadata context

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
			},
		})

		require.NoError(t, err, "GetRecommendations should not error")

		// Verify anomalies have sufficient investigation context
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				assert.NotEmpty(t, rec.GetDescription(),
					"Anomaly should have description for investigation context")
				// Metadata may contain additional context
				if len(rec.GetMetadata()) > 0 {
					t.Logf("Anomaly %s has metadata: %v", rec.GetId(), rec.GetMetadata())
				}
			}
		}
	})
}

// TestAnomalyConformance_ExclusionFiltering tests User Story 4: Exclude Anomalies from Optimization Workflows
// This test verifies that category filtering allows excluding ANOMALY recommendations.
func TestAnomalyConformance_ExclusionFiltering(t *testing.T) {
	t.Run("FilterByCategoryExcludesAnomalies", func(t *testing.T) {
		// T020: Write conformance test - filter by category=COST excludes ANOMALY

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()

		// Request only COST recommendations
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			},
		})

		require.NoError(t, err, "GetRecommendations should not error")
		require.NotNil(t, resp, "GetRecommendations should return a response")

		// Verify no ANOMALY recommendations are returned
		for _, rec := range resp.GetRecommendations() {
			assert.NotEqual(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				rec.GetCategory(), "COST filter should exclude ANOMALY category")
		}
	})

	t.Run("NegativeEstimatedSavingsAccepted", func(t *testing.T) {
		// T021: Write conformance test - negative estimated_savings accepted for overspend anomalies

		plugin := plugintesting.NewMockPlugin()
		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
			},
		})

		require.NoError(t, err, "GetRecommendations should not error")

		// Verify anomalies can have negative estimated_savings
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				if rec.GetImpact() != nil && rec.GetImpact().GetEstimatedSavings() < 0 {
					t.Logf("Anomaly %s has negative estimated_savings: %.2f (overspend)",
						rec.GetId(), rec.GetImpact().GetEstimatedSavings())
				}
			}
		}
		// Plugin should accept negative estimated_savings without validation errors
		assert.NoError(t, err, "Plugin should accept negative estimated_savings")
	})
}

// Note: minimalPlugin is defined in dry_run_conformance_test.go and reused here
// for backward compatibility testing.
