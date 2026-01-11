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

// Package pricing provides domain types and validation for PulumiCost pricing specifications.
package pricing

import (
	"errors"
	"fmt"
	"math"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Growth thresholds and limits.
const (
	// HighGrowthRateThreshold is the rate above which warnings are generated (100% per period).
	// Rates exceeding this may indicate unrealistic hyper-growth assumptions.
	HighGrowthRateThreshold = 1.0

	// LongProjectionThreshold is the number of periods (months) above which warnings are
	// generated for exponential growth. Projections beyond this become increasingly unreliable.
	LongProjectionThreshold = 36

	// MinValidGrowthRate is the minimum allowed growth rate (-100%).
	// Rates below this would result in negative costs, which is invalid.
	MinValidGrowthRate = -1.0
)

// Growth validation errors.
var (
	// ErrMissingGrowthRate is returned when LINEAR or EXPONENTIAL growth type is specified
	// without a growth_rate value.
	ErrMissingGrowthRate = errors.New("growth_rate required for LINEAR/EXPONENTIAL growth type")

	// ErrInvalidGrowthRate is returned when growth_rate is less than -1.0,
	// which would result in negative costs.
	ErrInvalidGrowthRate = errors.New("growth_rate must be >= -1.0")
)

// ApplyLinearGrowth calculates the cost at period n using linear growth.
// Formula: cost_at_n = baseCost * (1 + rate * n)
//
// Parameters:
//   - baseCost: The starting cost (e.g., current monthly cost)
//   - rate: Growth rate as a decimal fraction (e.g., 0.10 for 10%)
//   - periods: Number of periods (e.g., months) to project
//
// Returns the projected cost at period n.
//
// Example:
//
//	base=$100, rate=0.10, periods=3
//	Period 0: $100 * (1 + 0.10 * 0) = $100.00
//	Period 1: $100 * (1 + 0.10 * 1) = $110.00
//	Period 2: $100 * (1 + 0.10 * 2) = $120.00
//	Period 3: $100 * (1 + 0.10 * 3) = $130.00
func ApplyLinearGrowth(baseCost, rate float64, periods int) float64 {
	return baseCost * (1 + rate*float64(periods))
}

// ApplyExponentialGrowth calculates the cost at period n using exponential (compounding) growth.
// Formula: cost_at_n = baseCost * (1 + rate)^n
//
// Parameters:
//   - baseCost: The starting cost (e.g., current monthly cost)
//   - rate: Growth rate as a decimal fraction (e.g., 0.05 for 5%)
//   - periods: Number of periods (e.g., months) to project
//
// Returns the projected cost at period n.
//
// Example:
//
//	base=$100, rate=0.10, periods=3
//	Period 0: $100 * (1.10)^0 = $100.00
//	Period 1: $100 * (1.10)^1 = $110.00
//	Period 2: $100 * (1.10)^2 = $121.00
//	Period 3: $100 * (1.10)^3 = $133.10
func ApplyExponentialGrowth(baseCost, rate float64, periods int) float64 {
	return baseCost * math.Pow(1+rate, float64(periods))
}

// ApplyGrowth calculates the projected cost based on the growth type.
// This is a convenience function that dispatches to the appropriate growth function.
//
// Parameters:
//   - baseCost: The starting cost
//   - growthType: The type of growth model to apply
//   - rate: Growth rate (pointer, nil treated as 0)
//   - periods: Number of periods to project
//
// For GROWTH_TYPE_NONE or GROWTH_TYPE_UNSPECIFIED, returns the base cost unchanged.
// For GROWTH_TYPE_LINEAR, applies linear growth formula.
// For GROWTH_TYPE_EXPONENTIAL, applies exponential growth formula.
func ApplyGrowth(baseCost float64, growthType pbc.GrowthType, rate *float64, periods int) float64 {
	// Handle nil rate as 0 (no growth)
	r := 0.0
	if rate != nil {
		r = *rate
	}

	switch growthType {
	case pbc.GrowthType_GROWTH_TYPE_LINEAR:
		return ApplyLinearGrowth(baseCost, r, periods)
	case pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL:
		return ApplyExponentialGrowth(baseCost, r, periods)
	case pbc.GrowthType_GROWTH_TYPE_NONE, pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED:
		return baseCost
	default:
		// Unknown growth type - return base cost unchanged
		return baseCost
	}
}

// ResolveGrowthType returns the effective growth type, treating UNSPECIFIED as NONE.
// This ensures consistent default behavior across the codebase.
func ResolveGrowthType(gt pbc.GrowthType) pbc.GrowthType {
	if gt == pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED {
		return pbc.GrowthType_GROWTH_TYPE_NONE
	}
	return gt
}

// ResolveGrowthParams merges request-level growth parameters with resource-level defaults.
// Request-level parameters take precedence when set.
//
// Parameters:
//   - reqType: Growth type from request (overrides if not UNSPECIFIED)
//   - reqRate: Growth rate from request (overrides if not nil)
//   - resType: Default growth type from resource
//   - resRate: Default growth rate from resource
//
// Returns the effective growth type and rate to use.
func ResolveGrowthParams(
	reqType pbc.GrowthType, reqRate *float64,
	resType pbc.GrowthType, resRate *float64,
) (pbc.GrowthType, *float64) {
	// Start with resource defaults
	effectiveType := resType
	effectiveRate := resRate

	// Override with request-level values if set
	if reqType != pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED {
		effectiveType = reqType
	}
	if reqRate != nil {
		effectiveRate = reqRate
	}

	// Normalize UNSPECIFIED to NONE
	effectiveType = ResolveGrowthType(effectiveType)

	return effectiveType, effectiveRate
}

// ValidateGrowthParams validates the growth_type and growth_rate combination.
// Returns an error if the parameters are invalid:
//   - LINEAR or EXPONENTIAL without growth_rate: returns ErrMissingGrowthRate
//   - growth_rate < -1.0: returns ErrInvalidGrowthRate
//   - NONE or UNSPECIFIED: always valid (rate is ignored if present)
//
// This function is designed for use in plugin implementations to validate
// incoming requests before processing.
func ValidateGrowthParams(growthType pbc.GrowthType, growthRate *float64) error {
	// Check if rate is required but missing
	if growthType == pbc.GrowthType_GROWTH_TYPE_LINEAR ||
		growthType == pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL {
		if growthRate == nil {
			return fmt.Errorf("%w: %s", ErrMissingGrowthRate, growthType.String())
		}
	}

	// Check if rate is in valid range (if provided)
	if growthRate != nil && *growthRate < MinValidGrowthRate {
		return fmt.Errorf("%w: got %.4f (minimum allowed: %.1f)", ErrInvalidGrowthRate, *growthRate, MinValidGrowthRate)
	}

	return nil
}

// GrowthWarning represents a warning condition for growth parameter usage.
type GrowthWarning struct {
	Code    string
	Message string
	Rate    float64
	Periods int
}

// CheckGrowthWarnings examines growth parameters and returns any warnings.
// These are non-fatal conditions that may indicate unrealistic projections.
//
// Warnings are returned for:
//   - High growth rate (>100% per period): May indicate hyper-growth assumptions
//   - Long exponential projections (>36 months): Results become increasingly unreliable
//
// Note: For overflow detection, use CheckGrowthWarningsWithCost which includes
// OVERFLOW_RISK warnings when baseCost is provided.
func CheckGrowthWarnings(growthType pbc.GrowthType, growthRate *float64, periods int) []GrowthWarning {
	return checkGrowthWarningsInternal(0, growthType, growthRate, periods, false)
}

// CheckGrowthWarningsWithCost examines growth parameters including baseCost and returns warnings.
// This extended version includes OVERFLOW_RISK detection in addition to standard warnings.
//
// Warnings are returned for:
//   - OVERFLOW_RISK: Calculation would overflow to +Inf
//   - HIGH_GROWTH_RATE: Rate exceeds 100% per period
//   - LONG_PROJECTION: Exponential projection over 36 months
func CheckGrowthWarningsWithCost(
	baseCost float64,
	growthType pbc.GrowthType,
	growthRate *float64,
	periods int,
) []GrowthWarning {
	return checkGrowthWarningsInternal(baseCost, growthType, growthRate, periods, true)
}

// checkGrowthWarningsInternal is the shared implementation for warning checks.
// Optimized to count warnings first, then allocate slice with exact capacity.
func checkGrowthWarningsInternal(
	baseCost float64,
	growthType pbc.GrowthType,
	growthRate *float64,
	periods int,
	checkOverflow bool,
) []GrowthWarning {
	// Check conditions and count warnings first to avoid reallocation
	hasOverflowRisk := checkOverflow && CheckOverflowRisk(baseCost, growthType, growthRate, periods)

	// Early return if no rate - only overflow warning possible
	if growthRate == nil {
		if hasOverflowRisk {
			return []GrowthWarning{{
				Code: "OVERFLOW_RISK",
				Message: fmt.Sprintf(
					"projection with baseCost=%.2f, rate=0.0000, periods=%d would overflow to +Inf",
					baseCost,
					periods,
				),
				Rate:    0,
				Periods: periods,
			}}
		}
		return nil
	}

	rate := *growthRate
	hasHighGrowth := rate > HighGrowthRateThreshold
	hasLongProjection := growthType == pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL && periods > LongProjectionThreshold

	// Count warnings to allocate exact capacity
	count := 0
	if hasOverflowRisk {
		count++
	}
	if hasHighGrowth {
		count++
	}
	if hasLongProjection {
		count++
	}

	// No warnings - return nil (zero allocation)
	if count == 0 {
		return nil
	}

	// Allocate slice with exact capacity
	warnings := make([]GrowthWarning, 0, count)

	// Add warnings in order
	if hasOverflowRisk {
		warnings = append(warnings, GrowthWarning{
			Code: "OVERFLOW_RISK",
			Message: fmt.Sprintf(
				"projection with baseCost=%.2f, rate=%.4f, periods=%d would overflow to +Inf",
				baseCost,
				rate,
				periods,
			),
			Rate:    rate,
			Periods: periods,
		})
	}

	const percentMultiplier = 100.0
	if hasHighGrowth {
		warnings = append(warnings, GrowthWarning{
			Code: "HIGH_GROWTH_RATE",
			Message: fmt.Sprintf(
				"growth_rate %.2f (%.0f%%) exceeds 100%% per period - verify assumptions",
				rate,
				rate*percentMultiplier,
			),
			Rate:    rate,
			Periods: periods,
		})
	}

	if hasLongProjection {
		multiplier := math.Pow(1+rate, float64(periods))
		warnings = append(warnings, GrowthWarning{
			Code: "LONG_PROJECTION",
			Message: fmt.Sprintf(
				"exponential projection over %d periods results in %.1fx multiplier - results may be unreliable",
				periods,
				multiplier,
			),
			Rate:    rate,
			Periods: periods,
		})
	}

	return warnings
}

// IsHighGrowthRate returns true if the rate exceeds HighGrowthRateThreshold (100% per period).
// This is a convenience function for checking if a warning log should be emitted.
func IsHighGrowthRate(rate float64) bool {
	return rate > HighGrowthRateThreshold
}

// IsLongProjection returns true if exponential growth is projected over LongProjectionThreshold months.
// This is a convenience function for checking if a warning log should be emitted.
func IsLongProjection(growthType pbc.GrowthType, periods int) bool {
	return growthType == pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL && periods > LongProjectionThreshold
}

// CheckOverflowRisk returns true if the growth calculation would likely overflow to +Inf.
// This allows callers to detect potential overflow before performing the calculation.
//
// For exponential growth, overflow occurs when (1+rate)^periods exceeds MaxFloat64/baseCost.
// For linear growth, overflow occurs when baseCost * rate * periods exceeds MaxFloat64.
//
// Parameters:
//   - baseCost: The starting cost (must be > 0 for overflow to be possible)
//   - growthType: The type of growth model to check
//   - rate: Growth rate (nil is treated as 0, no overflow possible)
//   - periods: Number of periods to project
//
// Returns true if overflow is likely, false otherwise.
func CheckOverflowRisk(
	baseCost float64,
	growthType pbc.GrowthType,
	rate *float64,
	periods int,
) bool {
	// No risk for no-growth types or non-positive base/periods
	if growthType == pbc.GrowthType_GROWTH_TYPE_NONE ||
		growthType == pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED {
		return false
	}

	if baseCost <= 0 || periods <= 0 || rate == nil {
		return false
	}

	r := *rate
	if r <= 0 {
		// Negative or zero growth cannot overflow
		return false
	}

	switch growthType {
	case pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL:
		// Overflow when: baseCost * (1+r)^periods > MaxFloat64
		// Rearranged: periods > log(MaxFloat64/baseCost) / log(1+r)
		if baseCost >= math.MaxFloat64 {
			return true // Already at max
		}
		logBase := math.Log(1 + r)
		if logBase <= 0 {
			return false // Rate near 0, won't overflow
		}
		maxSafePeriods := math.Log(math.MaxFloat64/baseCost) / logBase
		return float64(periods) > maxSafePeriods

	case pbc.GrowthType_GROWTH_TYPE_LINEAR:
		// Overflow when: baseCost * (1 + r * periods) > MaxFloat64
		// Rearranged: baseCost * r * periods > MaxFloat64 - baseCost
		if baseCost >= math.MaxFloat64 {
			return true // Already at max
		}
		// Check if the multiplication would overflow
		// baseCost * r * periods
		product := baseCost * r * float64(periods)
		return math.IsInf(product, 1) || product > math.MaxFloat64-baseCost

	case pbc.GrowthType_GROWTH_TYPE_NONE, pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED:
		// No growth types never overflow
		return false
	}

	return false
}

// ErrOverflow is returned when a growth projection would overflow to +Inf or NaN.
var ErrOverflow = errors.New("growth projection resulted in overflow")

// ProjectCostSafely calculates projected cost with built-in validation and warning checks.
// This is a convenience function that combines validation, overflow detection, and calculation.
//
// Unlike ApplyGrowth, this function:
//   - Validates growth parameters before calculation
//   - Checks for overflow risk and returns ErrOverflow instead of +Inf
//   - Returns any growth warnings detected
//
// Returns:
//   - cost: The projected cost (0 if validation/overflow fails)
//   - warnings: Any growth warnings detected (may be nil)
//   - err: Validation or overflow errors (nil on success)
func ProjectCostSafely(
	baseCost float64,
	growthType pbc.GrowthType,
	rate *float64,
	periods int,
) (float64, []GrowthWarning, error) {
	// Validate parameters first
	if validErr := ValidateGrowthParams(growthType, rate); validErr != nil {
		return 0, nil, validErr
	}

	// Collect warnings
	warnings := CheckGrowthWarnings(growthType, rate, periods)

	// Check for overflow risk before calculation
	if CheckOverflowRisk(baseCost, growthType, rate, periods) {
		return 0, warnings, ErrOverflow
	}

	// Perform the calculation
	cost := ApplyGrowth(baseCost, growthType, rate, periods)

	// Final check for unexpected overflow/NaN (shouldn't happen if CheckOverflowRisk is accurate)
	if math.IsInf(cost, 0) || math.IsNaN(cost) {
		return 0, warnings, ErrOverflow
	}

	return cost, warnings, nil
}
