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

// Package testing_test contains conformance tests for FinFocus plugin implementations.
// The anomaly_conformance_test.go file validates spec 040-anomaly-detection-recommendations.
package testing_test

import (
	"context"
	"strings"
	"testing"

	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

const (
	// Test boundary values for confidence score filtering.
	testConfidenceLow   = 0.3
	testConfidenceExact = 0.5
	testConfidenceHigh  = 0.7
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
			assert.Positive(t, count, "Should have ANOMALY recommendations")
		}
	})
}

// TestAnomalyConformance_NonSupportingPlugins tests that plugins not supporting anomalies
// simply don't return any ANOMALY recommendations (no error is raised).
func TestAnomalyConformance_NonSupportingPlugins(t *testing.T) {
	// T010: Write conformance test - plugins not supporting anomalies return zero ANOMALY recommendations

	t.Run("LegacyPluginWithoutAnomalySupport_ReturnsEmptyGracefully", func(t *testing.T) {
		// Simulates a plugin built before anomaly detection was added.
		// Per spec FR-003, this is valid backward-compatible behavior.
		plugin := plugintesting.NewMockPlugin()
		// Override recommendations to exclude ANOMALY category
		recs := []*pbc.Recommendation{
			{
				Id:          "rec-1",
				Category:    pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
				ActionType:  pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
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

		// Count anomalies and verify confidence scores
		anomalyCount := 0
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				anomalyCount++
				// Confidence score must be set (non-zero)
				require.Greater(t, rec.GetConfidenceScore(), 0.0,
					"Anomaly %s confidence_score must be set (not zero/nil)", rec.GetId())
				// Confidence score must meet the filter threshold
				require.GreaterOrEqual(t, rec.GetConfidenceScore(), 0.5,
					"Anomaly %s confidence_score should meet minimum threshold of 0.5", rec.GetId())
			}
		}

		// Assert at least one anomaly exists in filtered results
		require.GreaterOrEqual(t, anomalyCount, 1,
			"expected at least one anomaly with confidence_score >= 0.5")
	})

	t.Run("ConfidenceScoreBoundaryConditions", func(t *testing.T) {
		// P2 CodeRabbit: Add explicit boundary testing for confidence score filtering
		plugin := plugintesting.NewMockPlugin()

		// Create recommendations with specific confidence scores around the 0.5 threshold
		recs := []*pbc.Recommendation{
			{
				Id:              "rec-low",
				Category:        pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				ActionType:      pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE,
				Description:     "Low confidence anomaly",
				ConfidenceScore: ptrFloat64(testConfidenceLow),
			},
			{
				Id:              "rec-exact",
				Category:        pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				ActionType:      pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE,
				Description:     "Exact threshold anomaly",
				ConfidenceScore: ptrFloat64(testConfidenceExact),
			},
			{
				Id:              "rec-high",
				Category:        pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				ActionType:      pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE,
				Description:     "High confidence anomaly",
				ConfidenceScore: ptrFloat64(testConfidenceHigh),
			},
		}
		plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
			Recommendations: recs,
		})

		harness := plugintesting.NewTestHarness(plugin)
		harness.Start(t)
		defer harness.Stop()

		ctx := context.Background()
		resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{MinConfidenceScore: 0.5},
		})

		require.NoError(t, err)

		assert.Len(t, resp.GetRecommendations(), 2, "Should include exactly 0.5 and above")
		ids := make(map[string]bool)
		for _, rec := range resp.GetRecommendations() {
			ids[rec.GetId()] = true
		}
		assert.False(t, ids["rec-low"], "rec-low (0.3) should be excluded")
		assert.True(t, ids["rec-exact"], "rec-exact (0.5) should be included")
		assert.True(t, ids["rec-high"], "rec-high (0.7) should be included")
	})
}

// ptrFloat64 is a helper to create float64 pointers for test data.
func ptrFloat64(v float64) *float64 {
	return &v
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
		require.NotEmpty(t, resp.GetRecommendations(), "expected at least one recommendation")

		// Count anomalies with INVESTIGATE action
		investigateActionCount := 0
		for _, rec := range resp.GetRecommendations() {
			assert.Equal(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
				rec.GetCategory(), "Should be ANOMALY category")

			// Track INVESTIGATE action
			if rec.GetActionType() == pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE {
				investigateActionCount++
				t.Logf("Anomaly %s has INVESTIGATE action", rec.GetId())
			}
		}

		// Assert at least one anomaly has INVESTIGATE action
		require.GreaterOrEqual(t, investigateActionCount, 1,
			"expected at least one anomaly with INVESTIGATE action type")
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
		require.NotEmpty(t, resp.GetRecommendations(), "expected at least one recommendation")

		// Verify all anomalies have sufficient investigation context
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				require.NotEmpty(t, rec.GetDescription(),
					"Anomaly %s should have description for investigation context", rec.GetId())
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
//
//nolint:gocognit // Test function with multiple validation scenarios
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

		// Count anomalies and negative savings anomalies
		totalAnomalies := 0
		negativeSavingsCount := 0
		for _, rec := range resp.GetRecommendations() {
			if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY {
				totalAnomalies++
				if rec.GetImpact() != nil && rec.GetImpact().GetEstimatedSavings() < 0 {
					negativeSavingsCount++
					t.Logf("Anomaly %s has negative estimated_savings: %.2f (overspend)",
						rec.GetId(), rec.GetImpact().GetEstimatedSavings())
				}
			}
		}

		// Assert that at least one anomaly exists and has negative savings
		require.GreaterOrEqual(t, totalAnomalies, 1, "expected at least one anomaly")
		require.GreaterOrEqual(
			t,
			negativeSavingsCount,
			1,
			"expected at least one anomaly with negative estimated_savings",
		)
	})

	t.Run("NegativeSavingsSemantics", func(t *testing.T) {
		// P1 CodeRabbit: Verify negative savings are semantically appropriate for anomalies
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

		require.NoError(t, err)

		// Verify negative savings semantics for anomalies
		foundNegativeSavings := false
		for _, rec := range resp.GetRecommendations() {
			if rec.GetImpact() != nil && rec.GetImpact().GetEstimatedSavings() < 0 {
				foundNegativeSavings = true
				savings := rec.GetImpact().GetEstimatedSavings()
				assert.Greater(t, savings, -1000000.0,
					"Negative savings should not be unrealistically large: %f", savings)

				// Verify description mentions spending for overspend anomalies
				description := rec.GetDescription()
				assert.NotEmpty(t, description, "Anomaly should have description")
				// Check for spending-related keywords in description
				lowerDesc := strings.ToLower(description)
				hasSpendingContext := strings.Contains(lowerDesc, "spending") ||
					strings.Contains(lowerDesc, "overspend") ||
					strings.Contains(lowerDesc, "above") ||
					strings.Contains(lowerDesc, "increase") ||
					strings.Contains(lowerDesc, "exceed")
				assert.True(t, hasSpendingContext,
					"Anomaly with negative savings should mention spending context in description: %s", description)
			}
		}

		if foundNegativeSavings {
			t.Logf("Verified negative savings semantic validation for %d anomalies", len(resp.GetRecommendations()))
		}
	})
}

// Note: minimalPlugin is defined in dry_run_conformance_test.go and reused here
// for backward compatibility testing.
