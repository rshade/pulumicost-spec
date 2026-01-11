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

package pricing_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pricing"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

const tolerance = 0.0001

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance*math.Max(math.Abs(a), math.Abs(b))
}

func floatPtr(f float64) *float64 {
	return &f
}

func TestApplyLinearGrowth(t *testing.T) {
	tests := []struct {
		name     string
		baseCost float64
		rate     float64
		periods  int
		expected float64
	}{
		{"No growth", 100, 0, 3, 100},
		{"10% linear 3 periods", 100, 0.10, 3, 130},
		{"5% linear 12 periods", 100, 0.05, 12, 160},
		{"Zero base cost", 0, 0.10, 5, 0},
		{"One period", 100, 0.10, 1, 110},
		{"Zero periods", 100, 0.10, 0, 100},
		{"Negative rate", 100, -0.10, 3, 70},
		{"High rate", 100, 2.0, 3, 700},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyLinearGrowth(tt.baseCost, tt.rate, tt.periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyLinearGrowth(%v, %v, %d) = %v, want %v",
					tt.baseCost, tt.rate, tt.periods, result, tt.expected)
			}
		})
	}
}

func TestApplyExponentialGrowth(t *testing.T) {
	tests := []struct {
		name     string
		baseCost float64
		rate     float64
		periods  int
		expected float64
	}{
		{"No growth", 100, 0, 3, 100},
		{"10% exponential 3 periods", 100, 0.10, 3, 133.1},
		{"5% exponential 12 periods", 100, 0.05, 12, 179.5856},
		{"Zero base cost", 0, 0.10, 5, 0},
		{"One period", 100, 0.10, 1, 110},
		{"Zero periods", 100, 0.10, 0, 100},
		{"Negative rate", 100, -0.10, 3, 72.9},
		{"High rate", 100, 2.0, 3, 2700},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyExponentialGrowth(tt.baseCost, tt.rate, tt.periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyExponentialGrowth(%v, %v, %d) = %v, want %v",
					tt.baseCost, tt.rate, tt.periods, result, tt.expected)
			}
		})
	}
}

func TestApplyGrowth(t *testing.T) {
	tests := []struct {
		name       string
		baseCost   float64
		growthType pbc.GrowthType
		rate       *float64
		periods    int
		expected   float64
	}{
		{"LINEAR with rate", 100, pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), 3, 130},
		{"EXPONENTIAL with rate", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.10), 3, 133.1},
		{"NONE with rate", 100, pbc.GrowthType_GROWTH_TYPE_NONE, floatPtr(0.10), 3, 100},
		{"UNSPECIFIED with rate", 100, pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, floatPtr(0.10), 3, 100},
		{"LINEAR with nil rate", 100, pbc.GrowthType_GROWTH_TYPE_LINEAR, nil, 3, 100},
		{"EXPONENTIAL with nil rate", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, nil, 3, 100},
		{"NONE with nil rate", 100, pbc.GrowthType_GROWTH_TYPE_NONE, nil, 3, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyGrowth(tt.baseCost, tt.growthType, tt.rate, tt.periods)
			if !almostEqual(result, tt.expected) {
				t.Errorf("ApplyGrowth(%v, %v, %v, %d) = %v, want %v",
					tt.baseCost, tt.growthType, tt.rate, tt.periods, result, tt.expected)
			}
		})
	}
}

func TestResolveGrowthType(t *testing.T) {
	tests := []struct {
		name     string
		input    pbc.GrowthType
		expected pbc.GrowthType
	}{
		{"UNSPECIFIED becomes NONE", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, pbc.GrowthType_GROWTH_TYPE_NONE},
		{"NONE stays NONE", pbc.GrowthType_GROWTH_TYPE_NONE, pbc.GrowthType_GROWTH_TYPE_NONE},
		{"LINEAR stays LINEAR", pbc.GrowthType_GROWTH_TYPE_LINEAR, pbc.GrowthType_GROWTH_TYPE_LINEAR},
		{
			"EXPONENTIAL stays EXPONENTIAL",
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ResolveGrowthType(tt.input)
			if result != tt.expected {
				t.Errorf("ResolveGrowthType(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResolveGrowthParams(t *testing.T) {
	tests := []struct {
		name         string
		reqType      pbc.GrowthType
		reqRate      *float64
		resType      pbc.GrowthType
		resRate      *float64
		expectedType pbc.GrowthType
		expectedRate *float64
	}{
		{
			name:         "Request overrides resource type",
			reqType:      pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			reqRate:      nil,
			resType:      pbc.GrowthType_GROWTH_TYPE_LINEAR,
			resRate:      floatPtr(0.10),
			expectedType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			expectedRate: floatPtr(0.10),
		},
		{
			name:         "Request overrides resource rate",
			reqType:      pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED,
			reqRate:      floatPtr(0.20),
			resType:      pbc.GrowthType_GROWTH_TYPE_LINEAR,
			resRate:      floatPtr(0.10),
			expectedType: pbc.GrowthType_GROWTH_TYPE_LINEAR,
			expectedRate: floatPtr(0.20),
		},
		{
			name:         "Request overrides both",
			reqType:      pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			reqRate:      floatPtr(0.25),
			resType:      pbc.GrowthType_GROWTH_TYPE_LINEAR,
			resRate:      floatPtr(0.10),
			expectedType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			expectedRate: floatPtr(0.25),
		},
		{
			name:         "Use resource defaults",
			reqType:      pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED,
			reqRate:      nil,
			resType:      pbc.GrowthType_GROWTH_TYPE_LINEAR,
			resRate:      floatPtr(0.10),
			expectedType: pbc.GrowthType_GROWTH_TYPE_LINEAR,
			expectedRate: floatPtr(0.10),
		},
		{
			name:         "Both unspecified becomes NONE",
			reqType:      pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED,
			reqRate:      nil,
			resType:      pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED,
			resRate:      nil,
			expectedType: pbc.GrowthType_GROWTH_TYPE_NONE,
			expectedRate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultType, resultRate := pricing.ResolveGrowthParams(tt.reqType, tt.reqRate, tt.resType, tt.resRate)

			if resultType != tt.expectedType {
				t.Errorf("ResolveGrowthParams type = %v, want %v", resultType, tt.expectedType)
			}

			if tt.expectedRate == nil && resultRate != nil {
				t.Errorf("ResolveGrowthParams rate = %v, want nil", *resultRate)
			}
			if tt.expectedRate != nil && resultRate == nil {
				t.Errorf("ResolveGrowthParams rate = nil, want %v", *tt.expectedRate)
			}
			if tt.expectedRate != nil && resultRate != nil && *resultRate != *tt.expectedRate {
				t.Errorf("ResolveGrowthParams rate = %v, want %v", *resultRate, *tt.expectedRate)
			}
		})
	}
}

func BenchmarkApplyLinearGrowth(b *testing.B) {
	for range b.N {
		_ = pricing.ApplyLinearGrowth(100.0, 0.10, 12)
	}
}

func BenchmarkApplyExponentialGrowth(b *testing.B) {
	for range b.N {
		_ = pricing.ApplyExponentialGrowth(100.0, 0.10, 12)
	}
}

func BenchmarkApplyGrowth(b *testing.B) {
	rate := floatPtr(0.10)
	for range b.N {
		_ = pricing.ApplyGrowth(100.0, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, rate, 12)
	}
}

func BenchmarkResolveGrowthParams(b *testing.B) {
	reqRate := floatPtr(0.20)
	resRate := floatPtr(0.10)
	for range b.N {
		_, _ = pricing.ResolveGrowthParams(
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, reqRate,
			pbc.GrowthType_GROWTH_TYPE_LINEAR, resRate,
		)
	}
}

func TestValidateGrowthParams(t *testing.T) {
	tests := []struct {
		name       string
		growthType pbc.GrowthType
		rate       *float64
		wantErr    bool
		errType    error
	}{
		// Valid cases
		{"UNSPECIFIED nil", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, nil, false, nil},
		{"NONE nil", pbc.GrowthType_GROWTH_TYPE_NONE, nil, false, nil},
		{"LINEAR with rate", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), false, nil},
		{"EXPONENTIAL with rate", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.05), false, nil},
		{"LINEAR with zero rate", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.0), false, nil},
		{"LINEAR with -1.0 rate", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(-1.0), false, nil},

		// Invalid: missing rate
		{"LINEAR nil rate", pbc.GrowthType_GROWTH_TYPE_LINEAR, nil, true, pricing.ErrMissingGrowthRate},
		{"EXPONENTIAL nil rate", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, nil, true, pricing.ErrMissingGrowthRate},

		// Invalid: rate too low
		{"LINEAR rate < -1.0", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(-1.5), true, pricing.ErrInvalidGrowthRate},
		{
			"EXPONENTIAL rate < -1.0",
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(-2.0),
			true,
			pricing.ErrInvalidGrowthRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pricing.ValidateGrowthParams(tt.growthType, tt.rate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGrowthParams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errType != nil && err != nil {
				if !errors.Is(err, tt.errType) {
					t.Errorf("ValidateGrowthParams() error type = %v, want %v", err, tt.errType)
				}
			}
		})
	}
}

func BenchmarkValidateGrowthParams(b *testing.B) {
	rate := floatPtr(0.10)
	for range b.N {
		_ = pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, rate)
	}
}

func TestCheckGrowthWarnings(t *testing.T) {
	tests := []struct {
		name          string
		growthType    pbc.GrowthType
		growthRate    *float64
		periods       int
		expectedCodes []string
	}{
		{"Nil rate returns no warnings", pbc.GrowthType_GROWTH_TYPE_LINEAR, nil, 12, []string{}},
		{"Normal rate no warnings", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), 12, []string{}},
		{"High rate warning", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(1.5), 12, []string{"HIGH_GROWTH_RATE"}},
		{"200% rate warning", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(2.0), 12, []string{"HIGH_GROWTH_RATE"}},
		{
			"Long projection warning",
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(0.05),
			48,
			[]string{"LONG_PROJECTION"},
		},
		{
			"Both warnings",
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(2.0),
			48,
			[]string{"HIGH_GROWTH_RATE", "LONG_PROJECTION"},
		},
		{"Linear long projection no warning", pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.05), 48, []string{}},
		{"Exactly 36 periods no warning", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.05), 36, []string{}},
		{"37 periods warning", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.05), 37, []string{"LONG_PROJECTION"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := pricing.CheckGrowthWarnings(tt.growthType, tt.growthRate, tt.periods)

			if len(warnings) != len(tt.expectedCodes) {
				t.Errorf(
					"CheckGrowthWarnings() returned %d warnings, want %d",
					len(warnings),
					len(tt.expectedCodes),
				)
				return
			}

			for i, expectedCode := range tt.expectedCodes {
				if warnings[i].Code != expectedCode {
					t.Errorf("Warning[%d].Code = %s, want %s", i, warnings[i].Code, expectedCode)
				}
			}
		})
	}
}

func TestIsHighGrowthRate(t *testing.T) {
	tests := []struct {
		name     string
		rate     float64
		expected bool
	}{
		{"Zero rate", 0.0, false},
		{"10% rate", 0.10, false},
		{"100% rate boundary", 1.0, false},
		{"Just over 100%", 1.01, true},
		{"200% rate", 2.0, true},
		{"Negative rate", -0.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.IsHighGrowthRate(tt.rate)
			if result != tt.expected {
				t.Errorf("IsHighGrowthRate(%v) = %v, want %v", tt.rate, result, tt.expected)
			}
		})
	}
}

func TestIsLongProjection(t *testing.T) {
	tests := []struct {
		name       string
		growthType pbc.GrowthType
		periods    int
		expected   bool
	}{
		{"Linear 12 periods", pbc.GrowthType_GROWTH_TYPE_LINEAR, 12, false},
		{"Linear 48 periods", pbc.GrowthType_GROWTH_TYPE_LINEAR, 48, false},
		{"Exponential 12 periods", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, 12, false},
		{"Exponential 36 periods boundary", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, 36, false},
		{"Exponential 37 periods", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, 37, true},
		{"Exponential 48 periods", pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, 48, true},
		{"None 48 periods", pbc.GrowthType_GROWTH_TYPE_NONE, 48, false},
		{"Unspecified 48 periods", pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED, 48, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.IsLongProjection(tt.growthType, tt.periods)
			if result != tt.expected {
				t.Errorf("IsLongProjection(%v, %d) = %v, want %v", tt.growthType, tt.periods, result, tt.expected)
			}
		})
	}
}

func BenchmarkCheckGrowthWarnings(b *testing.B) {
	rate := floatPtr(2.0)
	for range b.N {
		_ = pricing.CheckGrowthWarnings(pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, rate, 48)
	}
}

func BenchmarkCheckGrowthWarningsNoWarnings(b *testing.B) {
	rate := floatPtr(0.10) // Normal rate, no warnings
	for range b.N {
		_ = pricing.CheckGrowthWarnings(pbc.GrowthType_GROWTH_TYPE_LINEAR, rate, 12)
	}
}

// TestExponentialGrowthOverflow documents behavior with extreme values that cause float64 overflow.
// This is important for understanding edge cases in production systems.
func TestExponentialGrowthOverflow(t *testing.T) {
	tests := []struct {
		name     string
		baseCost float64
		rate     float64
		periods  int
		checkInf bool // whether to verify result is +Inf
	}{
		{
			name:     "extreme rate causes overflow",
			baseCost: 100.0,
			rate:     100.0, // 10000% growth per period
			periods:  1000,
			checkInf: true,
		},
		{
			name:     "large periods cause overflow",
			baseCost: 100.0,
			rate:     0.10,   // 10% growth
			periods:  100000, // very long projection
			checkInf: true,
		},
		{
			name:     "max float64 base with any growth overflows",
			baseCost: math.MaxFloat64,
			rate:     0.01,
			periods:  1,
			checkInf: true,
		},
		{
			name:     "realistic high growth does not overflow",
			baseCost: 1000000.0, // $1M
			rate:     0.50,      // 50% annual growth
			periods:  120,       // 10 years monthly
			checkInf: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyExponentialGrowth(tt.baseCost, tt.rate, tt.periods)

			if tt.checkInf {
				if !math.IsInf(result, 1) {
					t.Errorf("Expected +Inf for extreme values, got %v", result)
				}
				t.Logf("Confirmed overflow to +Inf: baseCost=%v, rate=%v, periods=%d",
					tt.baseCost, tt.rate, tt.periods)
			} else if math.IsInf(result, 0) || math.IsNaN(result) {
				t.Errorf("Unexpected overflow for realistic values: got %v", result)
			}
		})
	}
}

// TestLinearGrowthOverflow documents behavior with extreme values for linear growth.
func TestLinearGrowthOverflow(t *testing.T) {
	tests := []struct {
		name     string
		baseCost float64
		rate     float64
		periods  int
		checkInf bool
	}{
		{
			name:     "extreme rate and periods causes overflow",
			baseCost: 1e300,   // very large base
			rate:     1e10,    // very large rate
			periods:  1000000, // many periods
			checkInf: true,
		},
		{
			name:     "max float64 base overflows on any positive growth",
			baseCost: math.MaxFloat64,
			rate:     0.01,
			periods:  1,
			checkInf: true,
		},
		{
			name:     "realistic linear growth does not overflow",
			baseCost: 1000000.0, // $1M
			rate:     0.10,      // 10% linear growth
			periods:  1200,      // 100 years monthly
			checkInf: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyLinearGrowth(tt.baseCost, tt.rate, tt.periods)

			if tt.checkInf {
				if !math.IsInf(result, 1) {
					t.Errorf("Expected +Inf for extreme values, got %v", result)
				}
			} else {
				if math.IsInf(result, 0) || math.IsNaN(result) {
					t.Errorf("Unexpected overflow for realistic values: got %v", result)
				}
			}
		})
	}
}

// TestApplyGrowthWithExtremeValues tests ApplyGrowth wrapper with extreme inputs.
func TestApplyGrowthWithExtremeValues(t *testing.T) {
	tests := []struct {
		name       string
		baseCost   float64
		growthType pbc.GrowthType
		rate       *float64
		periods    int
		checkInf   bool
	}{
		{
			name:       "exponential overflow via ApplyGrowth",
			baseCost:   100.0,
			growthType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			rate:       floatPtr(100.0),
			periods:    500,
			checkInf:   true,
		},
		{
			name:       "linear overflow via ApplyGrowth",
			baseCost:   math.MaxFloat64,
			growthType: pbc.GrowthType_GROWTH_TYPE_LINEAR,
			rate:       floatPtr(0.01),
			periods:    1,
			checkInf:   true,
		},
		{
			name:       "NONE never overflows",
			baseCost:   math.MaxFloat64,
			growthType: pbc.GrowthType_GROWTH_TYPE_NONE,
			rate:       floatPtr(100.0),
			periods:    1000000,
			checkInf:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyGrowth(tt.baseCost, tt.growthType, tt.rate, tt.periods)

			if tt.checkInf {
				if !math.IsInf(result, 1) {
					t.Errorf("Expected +Inf for extreme values, got %v", result)
				}
			} else {
				if math.IsInf(result, 0) || math.IsNaN(result) {
					t.Errorf("Unexpected overflow: got %v", result)
				}
			}
		})
	}
}

func TestCheckOverflowRisk(t *testing.T) {
	tests := []struct {
		name       string
		baseCost   float64
		growthType pbc.GrowthType
		rate       *float64
		periods    int
		expected   bool
	}{
		// No risk cases
		{"NONE never overflows", math.MaxFloat64, pbc.GrowthType_GROWTH_TYPE_NONE, floatPtr(100.0), 1000000, false},
		{
			"UNSPECIFIED never overflows",
			math.MaxFloat64,
			pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED,
			floatPtr(100.0),
			1000000,
			false,
		},
		{"zero base cost", 0, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(100.0), 1000, false},
		{"negative base cost", -100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(100.0), 1000, false},
		{"zero periods", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(100.0), 0, false},
		{"negative periods", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(100.0), -5, false},
		{"nil rate", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, nil, 1000, false},
		{"zero rate", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0), 1000000, false},
		{"negative rate", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(-0.5), 1000000, false},

		// Realistic safe cases
		{"realistic exponential safe", 1000000, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.50), 120, false},
		{"realistic linear safe", 1000000, pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), 1200, false},

		// Overflow cases - exponential
		{"exponential extreme rate overflow", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(100.0), 1000, true},
		{
			"exponential long periods overflow",
			100,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(0.10),
			100000,
			true,
		},
		{
			"exponential max base overflow",
			math.MaxFloat64,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(0.01),
			1,
			true,
		},

		// Overflow cases - linear
		{"linear extreme values overflow", 1e300, pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(1e10), 1000000, true},
		{"linear max base overflow", math.MaxFloat64, pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.01), 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.CheckOverflowRisk(tt.baseCost, tt.growthType, tt.rate, tt.periods)
			if result != tt.expected {
				t.Errorf("CheckOverflowRisk() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckGrowthWarningsWithCost(t *testing.T) {
	tests := []struct {
		name          string
		baseCost      float64
		growthType    pbc.GrowthType
		growthRate    *float64
		periods       int
		expectedCodes []string
	}{
		{"No overflow normal case", 100, pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), 12, []string{}},
		// 1000 periods with 100x rate overflows AND is long projection
		{
			"Overflow with long projection",
			100,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(100.0),
			1000,
			[]string{"OVERFLOW_RISK", "HIGH_GROWTH_RATE", "LONG_PROJECTION"},
		},
		// 50 periods with 100x rate: high rate + long projection but may not overflow at 50 periods
		{
			"High rate and long projection",
			100,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(2.0),
			48,
			[]string{"HIGH_GROWTH_RATE", "LONG_PROJECTION"},
		},
		// Overflow with extreme rate at boundary (>36 for long projection)
		{
			"Overflow extreme rate short projection",
			100,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(100.0),
			200,
			[]string{"OVERFLOW_RISK", "HIGH_GROWTH_RATE", "LONG_PROJECTION"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := pricing.CheckGrowthWarningsWithCost(tt.baseCost, tt.growthType, tt.growthRate, tt.periods)

			if len(warnings) != len(tt.expectedCodes) {
				t.Errorf(
					"CheckGrowthWarningsWithCost() returned %d warnings, want %d",
					len(warnings),
					len(tt.expectedCodes),
				)
				for i, w := range warnings {
					t.Logf("  warning[%d]: %s", i, w.Code)
				}
				return
			}

			for i, expectedCode := range tt.expectedCodes {
				if warnings[i].Code != expectedCode {
					t.Errorf("Warning[%d].Code = %s, want %s", i, warnings[i].Code, expectedCode)
				}
			}
		})
	}
}

func TestProjectCostSafely(t *testing.T) {
	tests := []struct {
		name         string
		baseCost     float64
		growthType   pbc.GrowthType
		rate         *float64
		periods      int
		expectedCost float64
		wantErr      bool
		errType      error
	}{
		// Success cases
		{"linear growth", 100, pbc.GrowthType_GROWTH_TYPE_LINEAR, floatPtr(0.10), 3, 130, false, nil},
		{"exponential growth", 100, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, floatPtr(0.10), 3, 133.1, false, nil},
		{"no growth", 100, pbc.GrowthType_GROWTH_TYPE_NONE, floatPtr(0.10), 3, 100, false, nil},

		// Validation errors
		{
			"missing rate for LINEAR",
			100,
			pbc.GrowthType_GROWTH_TYPE_LINEAR,
			nil,
			3,
			0,
			true,
			pricing.ErrMissingGrowthRate,
		},
		{
			"invalid rate below -1.0",
			100,
			pbc.GrowthType_GROWTH_TYPE_LINEAR,
			floatPtr(-1.5),
			3,
			0,
			true,
			pricing.ErrInvalidGrowthRate,
		},

		// Overflow error
		{
			"overflow detection",
			100,
			pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			floatPtr(100.0),
			1000,
			0,
			true,
			pricing.ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost, _, err := pricing.ProjectCostSafely(tt.baseCost, tt.growthType, tt.rate, tt.periods)

			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectCostSafely() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errType != nil && err != nil && !errors.Is(err, tt.errType) {
				t.Errorf("ProjectCostSafely() error type = %v, want %v", err, tt.errType)
			}

			if !tt.wantErr && !almostEqual(cost, tt.expectedCost) {
				t.Errorf("ProjectCostSafely() cost = %v, want %v", cost, tt.expectedCost)
			}
		})
	}
}

func TestProjectCostSafelyReturnsWarnings(t *testing.T) {
	// Test that warnings are returned even on success
	rate := 2.0 // High growth rate
	cost, warnings, err := pricing.ProjectCostSafely(100, pbc.GrowthType_GROWTH_TYPE_LINEAR, &rate, 48)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if cost == 0 {
		t.Error("Expected non-zero cost")
	}

	// Should have HIGH_GROWTH_RATE warning
	found := false
	for _, w := range warnings {
		if w.Code == "HIGH_GROWTH_RATE" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected HIGH_GROWTH_RATE warning")
	}
}

func BenchmarkCheckOverflowRisk(b *testing.B) {
	rate := floatPtr(0.10)
	for range b.N {
		_ = pricing.CheckOverflowRisk(100.0, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, rate, 120)
	}
}

func BenchmarkProjectCostSafely(b *testing.B) {
	rate := floatPtr(0.10)
	for range b.N {
		_, _, _ = pricing.ProjectCostSafely(100.0, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, rate, 12)
	}
}

// TestRateBoundaryAtMinusOne verifies behavior at the exact -1.0 boundary.
func TestRateBoundaryAtMinusOne(t *testing.T) {
	// Rate of exactly -1.0 should be valid and produce zero cost
	rate := -1.0

	// Validate that -1.0 is accepted
	err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, &rate)
	if err != nil {
		t.Errorf("Rate -1.0 should be valid, got error: %v", err)
	}

	// Linear: 100 * (1 + -1.0 * n) approaches 0 as n increases
	// At n=1: 100 * (1 - 1) = 0
	linearResult := pricing.ApplyLinearGrowth(100.0, rate, 1)
	if linearResult != 0 {
		t.Errorf("Linear at -1.0 rate, 1 period should be 0, got %v", linearResult)
	}

	// Exponential: 100 * (1 + -1.0)^n = 100 * 0^n = 0 for n > 0
	expResult := pricing.ApplyExponentialGrowth(100.0, rate, 1)
	if expResult != 0 {
		t.Errorf("Exponential at -1.0 rate, 1 period should be 0, got %v", expResult)
	}

	// At 0 periods, both should return base cost
	linearZero := pricing.ApplyLinearGrowth(100.0, rate, 0)
	if linearZero != 100.0 {
		t.Errorf("Linear at -1.0 rate, 0 periods should be 100, got %v", linearZero)
	}

	expZero := pricing.ApplyExponentialGrowth(100.0, rate, 0)
	if expZero != 100.0 {
		t.Errorf("Exponential at -1.0 rate, 0 periods should be 100, got %v", expZero)
	}
}

// TestChainedProjections tests applying growth multiple times sequentially.
func TestChainedProjections(t *testing.T) {
	rate := 0.10 // 10% growth

	// Scenario: Project cost for 3 periods, then use that as base for 3 more
	baseCost := 100.0

	// First projection: 3 periods
	cost1 := pricing.ApplyExponentialGrowth(baseCost, rate, 3)
	// Expected: 100 * 1.10^3 = 133.1

	// Second projection: 3 more periods starting from cost1
	cost2 := pricing.ApplyExponentialGrowth(cost1, rate, 3)
	// Expected: 133.1 * 1.10^3 = 177.16

	// This should equal single projection of 6 periods
	costDirect := pricing.ApplyExponentialGrowth(baseCost, rate, 6)
	// Expected: 100 * 1.10^6 = 177.16

	if !almostEqual(cost2, costDirect) {
		t.Errorf("Chained projection mismatch: chained=%v, direct=%v", cost2, costDirect)
	}

	// Verify intermediate value
	expected1 := 133.1
	if !almostEqual(cost1, expected1) {
		t.Errorf("First projection: got %v, want %v", cost1, expected1)
	}

	t.Logf("Chained projections verified: $%.2f → $%.2f → $%.2f (direct: $%.2f)",
		baseCost, cost1, cost2, costDirect)
}

// TestConcurrentGrowthCalculations verifies thread safety.
func TestConcurrentGrowthCalculations(t *testing.T) {
	const numGoroutines = 100
	const iterations = 1000

	rate := floatPtr(0.10)
	expectedLinear := 130.0      // 100 * (1 + 0.10 * 3)
	expectedExponential := 133.1 // 100 * 1.10^3

	errChan := make(chan error, numGoroutines*2)

	// Test ApplyGrowth concurrently
	for i := range numGoroutines {
		go func(id int) {
			for range iterations {
				result := pricing.ApplyGrowth(100.0, pbc.GrowthType_GROWTH_TYPE_LINEAR, rate, 3)
				if !almostEqual(result, expectedLinear) {
					errChan <- fmt.Errorf("goroutine %d: linear got %v, want %v", id, result, expectedLinear)
					return
				}
			}
			errChan <- nil
		}(i)

		go func(id int) {
			for range iterations {
				result := pricing.ApplyGrowth(100.0, pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, rate, 3)
				if !almostEqual(result, expectedExponential) {
					errChan <- fmt.Errorf("goroutine %d: exponential got %v, want %v", id, result, expectedExponential)
					return
				}
			}
			errChan <- nil
		}(i)
	}

	// Collect results
	for range numGoroutines * 2 {
		if err := <-errChan; err != nil {
			t.Error(err)
		}
	}
}

// TestNegativePeriods documents behavior with negative periods.
// Current implementation does not validate periods, treating negative as valid input.
func TestNegativePeriods(t *testing.T) {
	tests := []struct {
		name       string
		growthType pbc.GrowthType
		periods    int
		validate   func(t *testing.T, result float64)
	}{
		{
			name:       "linear negative periods produces negative multiplier",
			growthType: pbc.GrowthType_GROWTH_TYPE_LINEAR,
			periods:    -5,
			validate: func(t *testing.T, result float64) {
				// LINEAR: 100 * (1 + 0.10 * -5) = 100 * 0.5 = 50
				expected := 50.0
				if !almostEqual(result, expected) {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:       "exponential negative periods produces fractional result",
			growthType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
			periods:    -5,
			validate: func(t *testing.T, result float64) {
				// EXPONENTIAL: 100 * (1.10)^-5 = 100 * 0.6209... ≈ 62.09
				expected := 100.0 * math.Pow(1.10, -5)
				if !almostEqual(result, expected) {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			},
		},
		{
			name:       "NONE ignores negative periods",
			growthType: pbc.GrowthType_GROWTH_TYPE_NONE,
			periods:    -100,
			validate: func(t *testing.T, result float64) {
				if result != 100.0 {
					t.Errorf("Expected 100.0, got %v", result)
				}
			},
		},
	}

	baseCost := 100.0
	rate := floatPtr(0.10)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pricing.ApplyGrowth(baseCost, tt.growthType, rate, tt.periods)
			tt.validate(t, result)
			t.Logf("Negative periods behavior: type=%v, periods=%d, result=%v",
				tt.growthType, tt.periods, result)
		})
	}
}
