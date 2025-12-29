// tiered_pricing_test.go provides tests for the TieredPricingCalculator
// demonstrating proper testing patterns for pricing calculation logic.
package main

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestTieredPricingCalculator_BasicCalculations tests fundamental cost calculations.
func TestTieredPricingCalculator_BasicCalculations(t *testing.T) {
	// Standard AWS S3-style pricing tiers
	tiers := []PricingTier{
		{MinQuantity: 0, MaxQuantity: 50000, RatePerUnit: 0.023, Description: "First 50 TB"},
		{MinQuantity: 50000, MaxQuantity: 450000, RatePerUnit: 0.022, Description: "Next 400 TB"},
		{MinQuantity: 450000, MaxQuantity: 0, RatePerUnit: 0.021, Description: "Over 450 TB"},
	}
	calc := MustNewTieredPricingCalculator(tiers, "USD")

	tests := []struct {
		name          string
		usageGB       float64
		expectedCost  float64
		expectedTiers int
	}{
		{
			name:          "zero usage",
			usageGB:       0,
			expectedCost:  0,
			expectedTiers: 0,
		},
		{
			name:          "small usage within first tier",
			usageGB:       100,
			expectedCost:  100 * 0.023,
			expectedTiers: 1,
		},
		{
			name:          "exactly at first tier boundary",
			usageGB:       50000,
			expectedCost:  50000 * 0.023,
			expectedTiers: 1,
		},
		{
			name:          "cross first tier boundary",
			usageGB:       60000,
			expectedCost:  (50000 * 0.023) + (10000 * 0.022),
			expectedTiers: 2,
		},
		{
			name:          "exactly at second tier boundary",
			usageGB:       450000,
			expectedCost:  (50000 * 0.023) + (400000 * 0.022),
			expectedTiers: 2,
		},
		{
			name:          "cross into third tier",
			usageGB:       500000,
			expectedCost:  (50000 * 0.023) + (400000 * 0.022) + (50000 * 0.021),
			expectedTiers: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cost, breakdown := calc.CalculateCost(tc.usageGB)

			// Compare with tolerance for floating point
			if diff := cost - tc.expectedCost; diff > 0.0001 || diff < -0.0001 {
				t.Errorf("CalculateCost(%v) = %v, want %v", tc.usageGB, cost, tc.expectedCost)
			}

			if len(breakdown) != tc.expectedTiers {
				t.Errorf("CalculateCost(%v) returned %d tiers, want %d", tc.usageGB, len(breakdown), tc.expectedTiers)
			}
		})
	}
}

// TestTieredPricingCalculator_BreakdownDetails verifies the tier breakdown contains correct details.
func TestTieredPricingCalculator_BreakdownDetails(t *testing.T) {
	tiers := []PricingTier{
		{MinQuantity: 0, MaxQuantity: 1000, RatePerUnit: 0.10, Description: "Tier 1"},
		{MinQuantity: 1000, MaxQuantity: 0, RatePerUnit: 0.05, Description: "Tier 2"},
	}
	calc := MustNewTieredPricingCalculator(tiers, "USD")

	_, breakdown := calc.CalculateCost(1500)

	if len(breakdown) != 2 {
		t.Fatalf("expected 2 tiers in breakdown, got %d", len(breakdown))
	}

	// Check first tier
	if breakdown[0].UsageInTier != 1000 {
		t.Errorf("tier 1 usage = %v, want 1000", breakdown[0].UsageInTier)
	}
	if breakdown[0].TierCost != 100.0 { // 1000 * 0.10
		t.Errorf("tier 1 cost = %v, want 100.0", breakdown[0].TierCost)
	}

	// Check second tier
	if breakdown[1].UsageInTier != 500 {
		t.Errorf("tier 2 usage = %v, want 500", breakdown[1].UsageInTier)
	}
	if breakdown[1].TierCost != 25.0 { // 500 * 0.05
		t.Errorf("tier 2 cost = %v, want 25.0", breakdown[1].TierCost)
	}
}

// TestPricingSpecParsing verifies JSON parsing of PricingSpec with tiers.
func TestPricingSpecParsing(t *testing.T) {
	specJSON := `{
		"provider": "aws",
		"resource_type": "s3",
		"billing_mode": "tiered",
		"rate_per_unit": 0.023,
		"currency": "USD",
		"pricing_tiers": [
			{"min_quantity": 0, "max_quantity": 50000, "rate_per_unit": 0.023, "description": "First 50 TB"},
			{"min_quantity": 50000, "max_quantity": 0, "rate_per_unit": 0.022, "description": "Over 50 TB"}
		]
	}`

	var spec PricingSpec
	if err := json.Unmarshal([]byte(specJSON), &spec); err != nil {
		t.Fatalf("failed to parse spec: %v", err)
	}

	if spec.Provider != "aws" {
		t.Errorf("provider = %q, want %q", spec.Provider, "aws")
	}
	if spec.BillingMode != "tiered" {
		t.Errorf("billing_mode = %q, want %q", spec.BillingMode, "tiered")
	}
	if len(spec.PricingTiers) != 2 {
		t.Errorf("pricing_tiers count = %d, want 2", len(spec.PricingTiers))
	}

	// Verify calculator works with parsed tiers
	calc := MustNewTieredPricingCalculator(spec.PricingTiers, spec.Currency)
	cost, _ := calc.CalculateCost(75000)

	// 50000 * 0.023 + 25000 * 0.022 = 1150 + 550 = 1700
	expectedCost := 1700.0
	if diff := cost - expectedCost; diff > 0.0001 || diff < -0.0001 {
		t.Errorf("cost for 75TB = %v, want %v", cost, expectedCost)
	}
}

// TestTieredPricingCalculator_EdgeCases tests edge cases and boundary conditions.
func TestTieredPricingCalculator_EdgeCases(t *testing.T) {
	t.Run("single unlimited tier", func(t *testing.T) {
		tiers := []PricingTier{
			{MinQuantity: 0, MaxQuantity: 0, RatePerUnit: 0.05, Description: "Flat rate"},
		}
		calc := MustNewTieredPricingCalculator(tiers, "USD")
		cost, breakdown := calc.CalculateCost(1000)

		if cost != 50.0 {
			t.Errorf("cost = %v, want 50.0", cost)
		}
		if len(breakdown) != 1 {
			t.Errorf("breakdown length = %d, want 1", len(breakdown))
		}
	})

	t.Run("empty tiers", func(t *testing.T) {
		calc := MustNewTieredPricingCalculator(nil, "USD")
		cost, breakdown := calc.CalculateCost(1000)

		if cost != 0 {
			t.Errorf("cost = %v, want 0", cost)
		}
		if len(breakdown) != 0 {
			t.Errorf("breakdown length = %d, want 0", len(breakdown))
		}
	})

	t.Run("negative usage returns zero", func(t *testing.T) {
		tiers := []PricingTier{
			{MinQuantity: 0, MaxQuantity: 1000, RatePerUnit: 0.10},
		}
		calc := MustNewTieredPricingCalculator(tiers, "USD")
		cost, breakdown := calc.CalculateCost(-100)

		if cost != 0 {
			t.Errorf("cost = %v, want 0", cost)
		}
		if len(breakdown) != 0 {
			t.Errorf("breakdown length = %d, want 0", len(breakdown))
		}
	})

	t.Run("empty currency returns error", func(t *testing.T) {
		tiers := []PricingTier{
			{MinQuantity: 0, MaxQuantity: 1000, RatePerUnit: 0.10},
		}
		_, err := NewTieredPricingCalculator(tiers, "")
		if !errors.Is(err, ErrEmptyCurrency) {
			t.Errorf("expected ErrEmptyCurrency, got %v", err)
		}
	})

	t.Run("unsorted tiers returns error", func(t *testing.T) {
		tiers := []PricingTier{
			{MinQuantity: 1000, MaxQuantity: 0, RatePerUnit: 0.05},
			{MinQuantity: 0, MaxQuantity: 1000, RatePerUnit: 0.10},
		}
		_, err := NewTieredPricingCalculator(tiers, "USD")
		if !errors.Is(err, ErrTiersNotSorted) {
			t.Errorf("expected ErrTiersNotSorted, got %v", err)
		}
	})
}
