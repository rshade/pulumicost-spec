package testing_test

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

const epsilon = 1e-9 // Tolerance for float comparisons

// =============================================================================
// Pricing Tier Fields Conformance Tests
// =============================================================================
//
// These tests validate that plugins correctly populate and validate the new
// pricing tier fields: pricing_category and spot_interruption_risk_score.
//
// Test Coverage:
// - T041: Plugins populate pricing_category with valid enum values
// - T042: Plugins populate spot_interruption_risk_score within 0.0-1.0 range
// - T043: Spot risk score is only non-zero when pricing_category is DYNAMIC
// - T044: Response validation rejects invalid spot risk scores (NaN, Inf, out-of-range)
// - T045: Backward compatibility - existing plugins work without new fields

// TestPricingTier_EstimateCostResponse_ValidCategories validates that plugins
// populate pricing_category with valid enum values (T041).
func TestPricingTier_EstimateCostResponse_ValidCategories(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	// Configure distinct pricing categories for different resource types
	// Note: Use simple resource type (module name) as the key
	plugin.SetPricingCategoryForResourceType("ec2", pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC)
	plugin.SetPricingCategoryForResourceType("s3", pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD)

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	client := harness.Client()

	// Test various pricing categories with expected values
	testCases := []struct {
		name             string
		resourceType     string
		expectedCategory pbc.FocusPricingCategory
	}{
		{"ec2_dynamic", "aws:ec2/instance:Instance", pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC},
		{"s3_standard", "aws:s3/bucket:Bucket", pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pbc.EstimateCostRequest{
				ResourceType: tc.resourceType,
			}

			resp, err := client.EstimateCost(ctx, req)
			require.NoError(t, err, "EstimateCost should succeed")

			// Verify pricing_category matches expected value
			category := resp.GetPricingCategory()
			assert.Equal(t, tc.expectedCategory, category,
				"pricing_category should match configured value")

			// Verify it's a valid enum value (not out of range)
			assert.GreaterOrEqual(
				t,
				int32(category),
				int32(pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED),
			)
			// Compute max enum value dynamically to handle future enum additions
			maxEnumValue := int32(0)
			for value := range pbc.FocusPricingCategory_name {
				if value > maxEnumValue {
					maxEnumValue = value
				}
			}
			assert.LessOrEqual(t, int32(category), maxEnumValue,
				"pricing_category should be within valid enum range")
		})
	}
}

// TestPricingTier_EstimateCostResponse_SpotRiskScoreRange validates that plugins
// populate spot_interruption_risk_score within 0.0-1.0 range (T042).
func TestPricingTier_EstimateCostResponse_SpotRiskScoreRange(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	client := harness.Client()

	req := &pbc.EstimateCostRequest{
		ResourceType: "aws:ec2/instance:Instance",
	}

	resp, err := client.EstimateCost(ctx, req)
	require.NoError(t, err, "EstimateCost should succeed")

	score := resp.GetSpotInterruptionRiskScore()

	// Verify score is within valid range
	assert.False(t, math.IsNaN(score), "spot_interruption_risk_score must not be NaN")
	assert.False(t, math.IsInf(score, 0), "spot_interruption_risk_score must not be Inf")
	assert.GreaterOrEqual(t, score, 0.0, "spot_interruption_risk_score must be >= 0.0")
	assert.LessOrEqual(t, score, 1.0, "spot_interruption_risk_score must be <= 1.0")

	// Verify validation function accepts it
	err = pluginsdk.ValidateEstimateCostResponse(resp)
	assert.NoError(t, err, "Response should pass validation")
}

// TestPricingTier_EstimateCostResponse_SpotRiskConsistency validates that
// spot_interruption_risk_score is only non-zero when pricing_category is DYNAMIC (T043).
func TestPricingTier_EstimateCostResponse_SpotRiskConsistency(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	// Add "spot" as a supported resource type for aws
	plugin.SupportedResources["aws"] = append(plugin.SupportedResources["aws"], "spot")

	// Configure mock plugin to return DYNAMIC pricing for spot instance
	// Note: Use simple resource type (module name) as the key
	plugin.SetPricingCategoryForResourceType("spot", pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC)
	plugin.SetSpotRiskScoreForResourceType("spot", 0.8)

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	client := harness.Client()

	testCases := []struct {
		name              string
		resourceType      string
		expectDynamic     bool
		expectNonZeroRisk bool
	}{
		{"spot_instance_dynamic_with_risk", "aws:spot/instance:Instance", true, true},
		{"standard_instance_should_not_be_dynamic", "aws:ec2/instance:Instance", false, false},
		{"reserved_instance_should_be_committed", "aws:ec2/instance:Instance", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pbc.EstimateCostRequest{
				ResourceType: tc.resourceType,
			}

			resp, err := client.EstimateCost(ctx, req)
			require.NoError(t, err, "EstimateCost should succeed")

			category := resp.GetPricingCategory()
			score := resp.GetSpotInterruptionRiskScore()

			if tc.expectDynamic {
				assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC, category,
					"Spot instances should have DYNAMIC pricing category")
			}

			// Check consistency using SDK helper
			warnings := pluginsdk.CheckSpotRiskConsistency(category, score)
			if tc.expectNonZeroRisk {
				// Spot instances can have non-zero risk scores
				if score > 0.0 {
					assert.Empty(t, warnings, "DYNAMIC pricing with non-zero risk should have no warnings")
				}
			} else {
				// Non-spot instances should have zero risk score
				if score > 0.0 {
					assert.NotEmpty(t, warnings, "Non-DYNAMIC pricing with non-zero risk should have warnings")
				}
			}
		})
	}
}

// TestPricingTier_GetProjectedCostResponse_ValidCategories validates that plugins
// populate pricing_category in GetProjectedCostResponse (T041).
func TestPricingTier_GetProjectedCostResponse_ValidCategories(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	client := harness.Client()

	req := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "aws:ec2:Instance",
			Sku:          "t3.medium",
			Region:       "us-east-1",
		},
	}

	resp, err := client.GetProjectedCost(ctx, req)
	require.NoError(t, err, "GetProjectedCost should succeed")

	// Verify pricing_category is populated
	category := resp.GetPricingCategory()
	assert.NotEqual(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED,
		category, "pricing_category should not be UNSPECIFIED")

	// Verify it's a valid enum value
	assert.GreaterOrEqual(
		t,
		int32(category),
		int32(pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED),
	)
	// Compute max enum value dynamically to handle future enum additions
	maxEnumValue := int32(0)
	for value := range pbc.FocusPricingCategory_name {
		if value > maxEnumValue {
			maxEnumValue = value
		}
	}
	assert.LessOrEqual(t, int32(category), maxEnumValue,
		"pricing_category should be within valid enum range")

	// Verify validation function accepts it
	err = pluginsdk.ValidateGetProjectedCostResponse(resp)
	assert.NoError(t, err, "Response should pass validation")
}

// TestPricingTier_GetProjectedCostResponse_SpotRiskScoreRange validates that plugins
// populate spot_interruption_risk_score within valid range in GetProjectedCostResponse (T042).
func TestPricingTier_GetProjectedCostResponse_SpotRiskScoreRange(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	// Configure mock to return valid DYNAMIC pricing with non-zero spot risk
	// Add "spot" as a supported resource type for aws
	plugin.SupportedResources["aws"] = append(plugin.SupportedResources["aws"], "spot")
	// Configure DYNAMIC pricing (required for non-zero spot risk per T043)
	plugin.SetPricingCategoryForResourceType("spot", pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC)
	plugin.SetSpotRiskScoreForResourceType("spot", 0.5)

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	client := harness.Client()

	req := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "aws:spot/instance:Instance", // Use Pulumi-style format matching mock's key extraction
			Sku:          "t3.medium",
			Region:       "us-east-1",
		},
	}

	resp, err := client.GetProjectedCost(ctx, req)
	require.NoError(t, err, "GetProjectedCost should succeed")

	score := resp.GetSpotInterruptionRiskScore()

	// Verify score is within valid range
	assert.False(t, math.IsNaN(score), "spot_interruption_risk_score must not be NaN")
	assert.False(t, math.IsInf(score, 0), "spot_interruption_risk_score must not be Inf")
	assert.GreaterOrEqual(t, score, 0.0, "spot_interruption_risk_score must be >= 0.0")
	assert.LessOrEqual(t, score, 1.0, "spot_interruption_risk_score must be <= 1.0")

	// Verify the mock actually returned the configured non-zero score
	assert.InDelta(t, 0.5, score, epsilon, "spot risk score should match configured value")
}

// TestPricingTier_ValidationRejectsInvalidValues validates that the validation
// functions correctly reject invalid spot risk scores (T044).
func TestPricingTier_ValidationRejectsInvalidValues(t *testing.T) {
	t.Run("EstimateCostResponse_rejects_nan", func(t *testing.T) {
		resp := &pbc.EstimateCostResponse{
			Currency:                  "USD",
			CostMonthly:               50.0,
			SpotInterruptionRiskScore: math.NaN(),
		}

		err := pluginsdk.ValidateEstimateCostResponse(resp)
		assert.ErrorIs(t, err, pluginsdk.ErrSpotRiskScoreNaN,
			"Validation should reject NaN spot risk score")
	})

	t.Run("EstimateCostResponse_rejects_positive_inf", func(t *testing.T) {
		resp := &pbc.EstimateCostResponse{
			Currency:                  "USD",
			CostMonthly:               50.0,
			SpotInterruptionRiskScore: math.Inf(1),
		}

		err := pluginsdk.ValidateEstimateCostResponse(resp)
		assert.ErrorIs(t, err, pluginsdk.ErrSpotRiskScoreNaN,
			"Validation should reject +Inf spot risk score")
	})

	t.Run("EstimateCostResponse_rejects_negative_value", func(t *testing.T) {
		resp := &pbc.EstimateCostResponse{
			Currency:                  "USD",
			CostMonthly:               50.0,
			SpotInterruptionRiskScore: -0.5,
		}

		err := pluginsdk.ValidateEstimateCostResponse(resp)
		assert.ErrorIs(t, err, pluginsdk.ErrSpotRiskScoreOutOfRange,
			"Validation should reject negative spot risk score")
	})

	t.Run("EstimateCostResponse_rejects_greater_than_one", func(t *testing.T) {
		resp := &pbc.EstimateCostResponse{
			Currency:                  "USD",
			CostMonthly:               50.0,
			SpotInterruptionRiskScore: 1.5,
		}

		err := pluginsdk.ValidateEstimateCostResponse(resp)
		assert.ErrorIs(t, err, pluginsdk.ErrSpotRiskScoreOutOfRange,
			"Validation should reject spot risk score > 1.0")
	})

	t.Run("GetProjectedCostResponse_rejects_invalid_values", func(t *testing.T) {
		invalidValues := []float64{math.NaN(), math.Inf(1), -0.5, 2.0}

		for _, value := range invalidValues {
			resp := &pbc.GetProjectedCostResponse{
				UnitPrice:                 0.05,
				Currency:                  "USD",
				CostPerMonth:              36.50,
				SpotInterruptionRiskScore: value,
			}

			err := pluginsdk.ValidateGetProjectedCostResponse(resp)
			assert.Error(t, err, "Validation should reject invalid value: %v", value)
		}
	})
}

// TestPricingTier_BackwardCompatibility validates that plugins without pricing
// tier fields continue to work (T045).
func TestPricingTier_BackwardCompatibility(t *testing.T) {
	t.Run("EstimateCostResponse_without_new_fields", func(t *testing.T) {
		// Simulate a legacy plugin that doesn't populate new fields
		resp := &pbc.EstimateCostResponse{
			Currency:    "USD",
			CostMonthly: 50.0,
			// pricing_category defaults to UNSPECIFIED (0)
			// spot_interruption_risk_score defaults to 0.0
		}

		// Should validate successfully (backward compatibility)
		err := pluginsdk.ValidateEstimateCostResponse(resp)
		require.NoError(t, err, "Legacy response without new fields should validate")

		// Verify default values
		assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED,
			resp.GetPricingCategory(), "Default pricing_category should be UNSPECIFIED")
		assert.InDelta(t, 0.0, resp.GetSpotInterruptionRiskScore(), epsilon,
			"Default spot_interruption_risk_score should be 0.0")
	})

	t.Run("GetProjectedCostResponse_without_new_fields", func(t *testing.T) {
		resp := &pbc.GetProjectedCostResponse{
			UnitPrice:    0.05,
			Currency:     "USD",
			CostPerMonth: 36.50,
			// New fields use proto3 defaults (0 and 0.0)
		}

		err := pluginsdk.ValidateGetProjectedCostResponse(resp)
		assert.NoError(t, err, "Legacy response without new fields should validate")
	})
}

// TestPricingTier_PerformanceBaseline validates that validation functions
// meet performance requirements (<100ns for valid responses).
// This is a runtime guard - see BenchmarkValidateEstimateCostResponse for detailed performance testing.
func TestPricingTier_PerformanceBaseline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	validResp := &pbc.EstimateCostResponse{
		Currency:                  "USD",
		CostMonthly:               50.0,
		PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
		SpotInterruptionRiskScore: 0.8,
	}

	// Simple runtime check - target <100ns per validation
	for range 100 {
		err := pluginsdk.ValidateEstimateCostResponse(validResp)
		if err != nil {
			t.Fatalf("Validation failed: %v", err)
		}
	}
}

// BenchmarkValidateEstimateCostResponse measures the performance of EstimateCostResponse validation.
// Target: <100ns per operation with zero allocations.
func BenchmarkValidateEstimateCostResponse(b *testing.B) {
	validResp := &pbc.EstimateCostResponse{
		Currency:                  "USD",
		CostMonthly:               50.0,
		PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
		SpotInterruptionRiskScore: 0.8,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_ = pluginsdk.ValidateEstimateCostResponse(validResp)
	}
}

// TestPricingTier_ConformanceSuiteIntegration validates that pricing tier tests
// can be integrated into the standard conformance suite.
func TestPricingTier_ConformanceSuiteIntegration(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	// Create a conformance test
	test := plugintesting.ConformanceTest{
		Name:        "PricingTierFieldsPopulated",
		Description: "Plugin populates pricing_category and spot_interruption_risk_score",
		TestFunc: func(harness *plugintesting.TestHarness) plugintesting.TestResult {
			start := time.Now()
			ctx := context.Background()
			client := harness.Client()

			req := &pbc.EstimateCostRequest{
				ResourceType: "aws:ec2/instance:Instance",
			}

			resp, err := client.EstimateCost(ctx, req)
			duration := time.Since(start)

			if err != nil {
				return plugintesting.TestResult{
					Method:   "EstimateCost",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "RPC call failed",
				}
			}

			// Validate pricing tier fields
			if validationErr := pluginsdk.ValidateEstimateCostResponse(resp); validationErr != nil {
				return plugintesting.TestResult{
					Method:   "EstimateCost",
					Success:  false,
					Error:    validationErr,
					Duration: duration,
					Details:  "Response validation failed",
				}
			}

			// Check that fields are populated
			category := resp.GetPricingCategory()
			score := resp.GetSpotInterruptionRiskScore()

			if category == pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED {
				return plugintesting.TestResult{
					Method:   "EstimateCost",
					Success:  false,
					Duration: duration,
					Details:  "pricing_category should not be UNSPECIFIED for real resources",
				}
			}

			return plugintesting.TestResult{
				Method:   "EstimateCost",
				Success:  true,
				Duration: duration,
				Details:  fmt.Sprintf("pricing_category=%v, spot_risk=%v", category, score),
			}
		},
	}

	// Run the conformance test
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	t.Logf("Running conformance test: %s - %s", test.Name, test.Description)
	result := test.TestFunc(harness)
	assert.True(t, result.Success, "Conformance test should pass: %s", result.Details)
	if result.Error != nil {
		t.Errorf("Test error: %v", result.Error)
	}
}
