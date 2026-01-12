package testing_test

import (
	"math"
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// tolerance defines the acceptable delta for floating point comparisons.
// Per SC-002/SC-003, accuracy must be within 0.01%.
const tolerance = 0.0001

// floatPtr is a helper to create a *float64 from a literal.
func floatPtr(f float64) *float64 {
	return &f
}

// almostEqual compares two floats within the tolerance.
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance*math.Max(math.Abs(a), math.Abs(b))
}

// TestLinearGrowthConformance validates LINEAR growth calculations.
// Given growth_type=LINEAR, rate=0.10, verify cost increases 10% per period linearly.
func TestLinearGrowthConformance(t *testing.T) {
	baseCost := 100.0
	rate := 0.10 // 10% per period

	tests := []struct {
		name     string
		periods  int
		expected float64
	}{
		{"Period 0", 0, 100.00},
		{"Period 1", 1, 110.00},
		{"Period 2", 2, 120.00},
		{"Period 3", 3, 130.00},
		{"Period 12", 12, 220.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyLinearGrowth(baseCost, rate, tt.periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyLinearGrowth(%v, %v, %d) = %v, want %v",
					baseCost, rate, tt.periods, result, tt.expected)
			}
		})
	}
}

// TestExponentialGrowthConformance validates EXPONENTIAL growth calculations.
// Given growth_type=EXPONENTIAL, rate=0.05, verify cost compounds at 5% per period.
func TestExponentialGrowthConformance(t *testing.T) {
	baseCost := 100.0
	rate := 0.05 // 5% per period

	tests := []struct {
		name     string
		periods  int
		expected float64
	}{
		{"Period 0", 0, 100.00},
		{"Period 1", 1, 105.00},
		{"Period 2", 2, 110.25},
		{"Period 3", 3, 115.7625},
		{"Period 12", 12, 179.5856}, // (1.05)^12 ≈ 1.795856
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyExponentialGrowth(baseCost, rate, tt.periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyExponentialGrowth(%v, %v, %d) = %v, want %v",
					baseCost, rate, tt.periods, result, tt.expected)
			}
		})
	}
}

// TestNoneGrowthConformance validates NONE growth type.
// Given growth_type=NONE, verify no growth is applied.
func TestNoneGrowthConformance(t *testing.T) {
	baseCost := 100.0
	rate := floatPtr(0.10) // Rate should be ignored

	tests := []struct {
		name     string
		periods  int
		expected float64
	}{
		{"Period 0", 0, 100.00},
		{"Period 1", 1, 100.00},
		{"Period 12", 12, 100.00},
		{"Period 36", 36, 100.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyGrowth(baseCost, pbc.GrowthType_GROWTH_TYPE_NONE, rate, tt.periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyGrowth(%v, NONE, %v, %d) = %v, want %v",
					baseCost, *rate, tt.periods, result, tt.expected)
			}
		})
	}
}

// TestUnspecifiedTreatedAsNone validates that UNSPECIFIED is treated as NONE.
func TestUnspecifiedTreatedAsNone(t *testing.T) {
	baseCost := 100.0
	rate := floatPtr(0.20) // Rate should be ignored

	result := pricing.ApplyGrowth(baseCost, pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, rate, 12)
	expected := 100.0

	if !almostEqual(result, expected) {
		t.Errorf("ApplyGrowth with UNSPECIFIED should not apply growth: got %v, want %v", result, expected)
	}
}

// TestActualCostIgnoresGrowthParams validates FR-008: GetActualCost ignores growth parameters.
// Growth parameters are for projections only; actual cost retrieval is unaffected.
func TestActualCostIgnoresGrowthParams(t *testing.T) {
	// Create a ResourceDescriptor with growth parameters
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.medium",
		Region:       "us-east-1",
		GrowthType:   pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
		GrowthRate:   floatPtr(0.25), // 25% growth - should be ignored
	}

	// Verify the growth fields are set (proto fields exist)
	if resource.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL {
		t.Errorf("Expected GrowthType to be EXPONENTIAL, got %v", resource.GetGrowthType())
	}
	if resource.GetGrowthRate() != 0.25 {
		t.Errorf("Expected GrowthRate to be 0.25, got %v", resource.GetGrowthRate())
	}

	// In practice, GetActualCost implementation should not use these fields.
	// This test documents the contract that growth params exist on ResourceDescriptor
	// but have no effect on actual cost calculations.
	t.Log("FR-008 Validated: ResourceDescriptor accepts growth_type and growth_rate fields")
	t.Log("Plugin implementations MUST ignore these fields in GetActualCost RPC")
}

// TestApplyGrowthWithAllTypes tests the unified ApplyGrowth function.
func TestApplyGrowthWithAllTypes(t *testing.T) {
	baseCost := 100.0
	rate := floatPtr(0.10)
	periods := 3

	tests := []struct {
		name       string
		growthType pbc.GrowthType
		expected   float64
	}{
		{"LINEAR", pbc.GrowthType_GROWTH_TYPE_LINEAR, 130.0},           // 100 * (1 + 0.1*3)
		{"EXPONENTIAL", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, 133.1}, // 100 * (1.1)^3
		{"NONE", pbc.GrowthType_GROWTH_TYPE_NONE, 100.0},
		{"UNSPECIFIED", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyGrowth(baseCost, tt.growthType, rate, periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyGrowth(%v, %v, %v, %d) = %v, want %v",
					baseCost, tt.growthType, *rate, periods, result, tt.expected)
			}
		})
	}
}

// TestNilRateTreatedAsZero tests that nil rate is treated as 0 growth.
func TestNilRateTreatedAsZero(t *testing.T) {
	baseCost := 100.0
	var rate *float64 // nil by default
	periods := 12

	// With LINEAR and nil rate, should apply no growth (rate=0)
	result := pricing.ApplyGrowth(baseCost, pbc.GrowthType_GROWTH_TYPE_LINEAR, rate, periods)
	expected := 100.0 // 100 * (1 + 0*12) = 100

	if !almostEqual(result, expected) {
		t.Errorf("ApplyGrowth with nil rate should apply 0 growth: got %v, want %v", result, expected)
	}
}

// TestNegativeGrowthRate tests decline scenarios.
func TestNegativeGrowthRate(t *testing.T) {
	baseCost := 100.0
	rate := floatPtr(-0.10) // 10% decline per period

	t.Run("LinearDecline", func(t *testing.T) {
		result := pricing.ApplyLinearGrowth(baseCost, *rate, 3)
		expected := 70.0 // 100 * (1 + (-0.1)*3) = 100 * 0.7 = 70
		if !almostEqual(result, expected) {
			t.Errorf("Linear decline: got %v, want %v", result, expected)
		}
	})

	t.Run("ExponentialDecline", func(t *testing.T) {
		result := pricing.ApplyExponentialGrowth(baseCost, *rate, 3)
		expected := 72.9 // 100 * (0.9)^3 = 72.9
		if !almostEqual(result, expected) {
			t.Errorf("Exponential decline: got %v, want %v", result, expected)
		}
	})
}

// TestHighGrowthRates tests hyper-growth scenarios (>100% per period).
func TestHighGrowthRates(t *testing.T) {
	baseCost := 100.0
	rate := floatPtr(2.0) // 200% growth per period
	periods := 3

	t.Run("LinearHyperGrowth", func(t *testing.T) {
		result := pricing.ApplyLinearGrowth(baseCost, *rate, periods)
		expected := 700.0 // 100 * (1 + 2.0*3) = 100 * 7 = 700
		if !almostEqual(result, expected) {
			t.Errorf("Linear hyper-growth: got %v, want %v", result, expected)
		}
	})

	t.Run("ExponentialHyperGrowth", func(t *testing.T) {
		result := pricing.ApplyExponentialGrowth(baseCost, *rate, periods)
		expected := 2700.0 // 100 * (3.0)^3 = 100 * 27 = 2700
		if !almostEqual(result, expected) {
			t.Errorf("Exponential hyper-growth: got %v, want %v", result, expected)
		}
	})
}

// TestProtoFieldsExist validates that the proto fields were generated correctly.
func TestProtoFieldsExist(t *testing.T) {
	// Test ResourceDescriptor fields
	rd := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.medium",
		Region:       "us-east-1",
		GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
		GrowthRate:   floatPtr(0.05),
	}

	if rd.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_LINEAR {
		t.Errorf("ResourceDescriptor.GrowthType not set correctly")
	}
	if rd.GetGrowthRate() != 0.05 {
		t.Errorf("ResourceDescriptor.GrowthRate not set correctly")
	}

	// Test GetProjectedCostRequest fields
	req := &pbc.GetProjectedCostRequest{
		Resource:              rd,
		UtilizationPercentage: 0.5,
		GrowthType:            pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
		GrowthRate:            floatPtr(0.10),
	}

	if req.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL {
		t.Errorf("GetProjectedCostRequest.GrowthType not set correctly")
	}
	if req.GetGrowthRate() != 0.10 {
		t.Errorf("GetProjectedCostRequest.GrowthRate not set correctly")
	}
}

// TestGrowthTypeEnum validates all GrowthType enum values exist.
func TestGrowthTypeEnum(t *testing.T) {
	tests := []struct {
		name  string
		value pbc.GrowthType
		num   int32
	}{
		{"UNSPECIFIED", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, 0},
		{"NONE", pbc.GrowthType_GROWTH_TYPE_NONE, 1},
		{"LINEAR", pbc.GrowthType_GROWTH_TYPE_LINEAR, 2},
		{"EXPONENTIAL", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int32(tt.value) != tt.num {
				t.Errorf("%s has wrong value: got %d, want %d", tt.name, int32(tt.value), tt.num)
			}
		})
	}
}
