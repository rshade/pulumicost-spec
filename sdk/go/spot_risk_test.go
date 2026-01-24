package finfocus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

const epsilonSpot = 1e-9 // Tolerance for float comparisons

func TestEstimateCostResponse_SpotRisk(t *testing.T) {
	// T008: Implement unit test that constructs EstimateCostResponse with FOCUS_PRICING_CATEGORY_DYNAMIC and risk score 0.8
	resp := &pbc.EstimateCostResponse{
		Currency:                  "USD",
		CostMonthly:               50.0,
		PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
		SpotInterruptionRiskScore: 0.8,
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC, resp.GetPricingCategory())
	assert.InDelta(t, 0.8, resp.GetSpotInterruptionRiskScore(), epsilonSpot)
}

func TestEstimateCostResponse_Committed(t *testing.T) {
	// T011: Add test case for FOCUS_PRICING_CATEGORY_COMMITTED scenario (Reserved Instance)
	resp := &pbc.EstimateCostResponse{
		Currency:        "USD",
		CostMonthly:     30.0,
		PricingCategory: pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_COMMITTED,
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_COMMITTED, resp.GetPricingCategory())
	assert.InDelta(t, 0.0, resp.GetSpotInterruptionRiskScore(), epsilonSpot) // Should default to 0.0
}

func TestEstimateCostResponse_SpotRisk_Boundaries(t *testing.T) {
	testCases := []struct {
		name      string
		riskScore float64
	}{
		{"zero_risk", 0.0},
		{"medium_risk", 0.5},
		{"max_risk", 1.0},
		{"low_risk", 0.0001},
		{"near_max_risk", 0.9999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &pbc.EstimateCostResponse{
				Currency:                  "USD",
				CostMonthly:               50.0,
				PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
				SpotInterruptionRiskScore: tc.riskScore,
			}

			assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC, resp.GetPricingCategory())
			assert.InDelta(t, tc.riskScore, resp.GetSpotInterruptionRiskScore(), epsilonSpot)
		})
	}
}

func TestEstimateCostResponse_SpotRisk_IgnoredForStandard(t *testing.T) {
	// Spot risk score should be semantically ignored when pricing category is not DYNAMIC
	resp := &pbc.EstimateCostResponse{
		Currency:                  "USD",
		CostMonthly:               40.0,
		PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
		SpotInterruptionRiskScore: 0.5, // Non-zero value should be ignored
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD, resp.GetPricingCategory())
	// Note: The field is still accessible but semantically should be treated as 0.0 per service contract
	assert.InDelta(t, 0.5, resp.GetSpotInterruptionRiskScore(), epsilonSpot) // Field value is present
}

func TestEstimateCostResponse_UnspecifiedCategory(t *testing.T) {
	resp := &pbc.EstimateCostResponse{
		Currency:        "USD",
		CostMonthly:     35.0,
		PricingCategory: pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED,
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED, resp.GetPricingCategory())
	assert.InDelta(t, 0.0, resp.GetSpotInterruptionRiskScore(), epsilonSpot) // Should default to 0.0
}
