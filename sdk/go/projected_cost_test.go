package finfocus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

const epsilon = 1e-9 // Tolerance for float comparisons

// Benchmark sink variable to prevent compiler optimizations.
//
//nolint:gochecknoglobals // Intentional global for benchmark sink pattern
var benchmarkSink *pbc.GetProjectedCostResponse

func TestGetProjectedCostResponse_SpotRisk(t *testing.T) {
	// T009: Implement unit test that constructs GetProjectedCostResponse with FOCUS_PRICING_CATEGORY_DYNAMIC and risk score 0.8
	resp := &pbc.GetProjectedCostResponse{
		UnitPrice:                 0.05,
		Currency:                  "USD",
		CostPerMonth:              36.50,
		BillingDetail:             "spot-instance",
		PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
		SpotInterruptionRiskScore: 0.8,
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC, resp.GetPricingCategory())
	assert.InDelta(t, 0.8, resp.GetSpotInterruptionRiskScore(), epsilon)
}

func TestGetProjectedCostResponse_Committed(t *testing.T) {
	// T012: Add test case for FOCUS_PRICING_CATEGORY_COMMITTED scenario (Savings Plan)
	resp := &pbc.GetProjectedCostResponse{
		UnitPrice:       0.04,
		Currency:        "USD",
		CostPerMonth:    29.20,
		BillingDetail:   "savings-plan",
		PricingCategory: pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_COMMITTED,
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_COMMITTED, resp.GetPricingCategory())
	assert.InDelta(t, 0.0, resp.GetSpotInterruptionRiskScore(), epsilon) // Should default to 0.0
}

func TestGetProjectedCostResponse_SpotRisk_Boundaries(t *testing.T) {
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
			resp := &pbc.GetProjectedCostResponse{
				UnitPrice:                 0.05,
				Currency:                  "USD",
				CostPerMonth:              36.50,
				BillingDetail:             "spot-instance",
				PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
				SpotInterruptionRiskScore: tc.riskScore,
			}

			assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC, resp.GetPricingCategory())
			assert.InDelta(t, tc.riskScore, resp.GetSpotInterruptionRiskScore(), epsilon)
		})
	}
}

func TestGetProjectedCostResponse_SpotRisk_IgnoredForStandard(t *testing.T) {
	// Spot risk score should be semantically ignored when pricing category is not DYNAMIC
	resp := &pbc.GetProjectedCostResponse{
		UnitPrice:                 0.06,
		Currency:                  "USD",
		CostPerMonth:              43.80,
		BillingDetail:             "standard",
		PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
		SpotInterruptionRiskScore: 0.5, // Non-zero value should be ignored
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD, resp.GetPricingCategory())
	// Note: The field is still accessible but semantically should be treated as 0.0 per service contract
	assert.InDelta(t, 0.5, resp.GetSpotInterruptionRiskScore(), epsilon) // Field value is present
}

func TestGetProjectedCostResponse_UnspecifiedCategory(t *testing.T) {
	resp := &pbc.GetProjectedCostResponse{
		UnitPrice:       0.05,
		Currency:        "USD",
		CostPerMonth:    36.50,
		BillingDetail:   "unspecified",
		PricingCategory: pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED,
	}

	assert.Equal(t, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED, resp.GetPricingCategory())
	assert.InDelta(t, 0.0, resp.GetSpotInterruptionRiskScore(), epsilon) // Should default to 0.0
}

func BenchmarkGetProjectedCostResponse_Construction(b *testing.B) {
	b.ReportAllocs()
	var resp *pbc.GetProjectedCostResponse
	for range b.N {
		resp = &pbc.GetProjectedCostResponse{
			UnitPrice:                 0.05,
			Currency:                  "USD",
			CostPerMonth:              36.50,
			BillingDetail:             "spot-instance",
			PricingCategory:           pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
			SpotInterruptionRiskScore: 0.8,
		}
	}
	benchmarkSink = resp
}
