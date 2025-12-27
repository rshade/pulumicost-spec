// tiered_pricing.go demonstrates tiered pricing calculation patterns for
// cloud resources like storage that have volume-based pricing tiers.
//
// This example implements the TieredPricingCalculator pattern described
// in docs/ADVANCED_PATTERNS.md and shows how to:
// - Parse pricing tiers from a PricingSpec
// - Calculate costs across multiple tiers
// - Generate detailed cost breakdowns
//
// Reference: AWS S3 Standard storage pricing as of 2024
// https://aws.amazon.com/s3/pricing/
package main

import (
	"encoding/json"
	"log/slog"
	"os"
)

// PricingTier represents a single pricing tier with quantity bounds and rate.
type PricingTier struct {
	MinQuantity float64 `json:"min_quantity"`
	MaxQuantity float64 `json:"max_quantity"` // 0 means unlimited
	RatePerUnit float64 `json:"rate_per_unit"`
	Description string  `json:"description"`
}

// TierBreakdown provides details about cost in each tier.
type TierBreakdown struct {
	Tier        PricingTier
	UsageInTier float64
	TierCost    float64
}

// TieredPricingCalculator calculates costs across volume-based pricing tiers.
type TieredPricingCalculator struct {
	tiers    []PricingTier
	currency string
}

// NewTieredPricingCalculator creates a calculator with the given tiers.
// Tiers should be provided in ascending order by MinQuantity.
func NewTieredPricingCalculator(tiers []PricingTier, currency string) *TieredPricingCalculator {
	return &TieredPricingCalculator{
		tiers:    tiers,
		currency: currency,
	}
}

// CalculateCost computes the total cost for a given usage quantity.
// Returns the total cost and a breakdown by tier.
func (c *TieredPricingCalculator) CalculateCost(usage float64) (float64, []TierBreakdown) {
	var totalCost float64
	var breakdown []TierBreakdown

	for _, tier := range c.tiers {
		// Calculate usage within this tier
		tierMax := tier.MaxQuantity
		if tierMax == 0 {
			tierMax = usage + 1 // Unlimited tier
		}

		// Skip if usage hasn't reached this tier yet
		if usage <= tier.MinQuantity {
			continue
		}

		// Calculate quantity in this tier
		usageInTier := min(usage, tierMax) - tier.MinQuantity
		if usageInTier <= 0 {
			continue
		}

		tierCost := usageInTier * tier.RatePerUnit
		totalCost += tierCost
		breakdown = append(breakdown, TierBreakdown{
			Tier:        tier,
			UsageInTier: usageInTier,
			TierCost:    tierCost,
		})
	}

	return totalCost, breakdown
}

// PricingSpec represents a pricing specification with tiered pricing.
type PricingSpec struct {
	Provider     string        `json:"provider"`
	ResourceType string        `json:"resource_type"`
	BillingMode  string        `json:"billing_mode"`
	RatePerUnit  float64       `json:"rate_per_unit"`
	Currency     string        `json:"currency"`
	PricingTiers []PricingTier `json:"pricing_tiers"`
}

// demonstrateTieredPricing shows tiered pricing calculation with different usage levels.
//
//nolint:mnd // pricing data is intentionally hard-coded for demonstration purposes
func demonstrateTieredPricing(logger *slog.Logger) {
	// AWS S3 Standard storage pricing tiers (us-east-1, as of 2024)
	tiers := []PricingTier{
		{
			MinQuantity: 0,
			MaxQuantity: 50000, // First 50 TB
			RatePerUnit: 0.023,
			Description: "First 50 TB / Month",
		},
		{
			MinQuantity: 50000,
			MaxQuantity: 450000, // Next 400 TB
			RatePerUnit: 0.022,
			Description: "Next 400 TB / Month",
		},
		{
			MinQuantity: 450000,
			MaxQuantity: 0, // Unlimited (over 450 TB)
			RatePerUnit: 0.021,
			Description: "Over 450 TB / Month",
		},
	}

	calc := NewTieredPricingCalculator(tiers, "USD")

	// Test cases with different usage levels
	testCases := []struct {
		name    string
		usageGB float64
	}{
		{"Small (10 GB)", 10},
		{"Medium (100 GB)", 100},
		{"Large (1 TB)", 1000},
		{"First tier boundary (50 TB)", 50000},
		{"Cross tier (100 TB)", 100000},
		{"High volume (500 TB)", 500000},
	}

	logger.Info("AWS S3 Standard Storage Tiered Pricing Demo")

	for _, tc := range testCases {
		totalCost, breakdown := calc.CalculateCost(tc.usageGB)

		// Build breakdown attributes
		breakdownAttrs := make([]any, 0, len(breakdown)*3)
		for i, b := range breakdown {
			breakdownAttrs = append(breakdownAttrs,
				slog.Group("tier_"+string(rune('1'+i)),
					slog.String("description", b.Tier.Description),
					slog.Float64("usage_gb", b.UsageInTier),
					slog.Float64("cost", b.TierCost),
				),
			)
		}

		logger.Info("Cost calculation",
			slog.String("usage", tc.name),
			slog.Float64("usage_gb", tc.usageGB),
			slog.Float64("total_cost", totalCost),
			slog.Group("breakdown", breakdownAttrs...),
		)
	}
}

// demonstratePricingSpecParsing shows how to parse tiered pricing from a JSON spec.
func demonstratePricingSpecParsing(logger *slog.Logger) {
	// Example PricingSpec JSON (like examples/specs/aws-s3-tiered-pricing.json)
	specJSON := `{
		"provider": "aws",
		"resource_type": "s3",
		"billing_mode": "tiered",
		"rate_per_unit": 0.023,
		"currency": "USD",
		"pricing_tiers": [
			{"min_quantity": 0, "max_quantity": 50000, "rate_per_unit": 0.023, "description": "First 50 TB / Month"},
			{"min_quantity": 50000, "max_quantity": 450000, "rate_per_unit": 0.022, "description": "Next 400 TB / Month"},
			{"min_quantity": 450000, "max_quantity": 0, "rate_per_unit": 0.021, "description": "Over 450 TB / Month"}
		]
	}`

	var spec PricingSpec
	if err := json.Unmarshal([]byte(specJSON), &spec); err != nil {
		logger.Error("Failed to parse pricing spec", "error", err)
		os.Exit(1)
	}

	logger.Info("Parsed PricingSpec",
		slog.String("provider", spec.Provider),
		slog.String("resource_type", spec.ResourceType),
		slog.String("billing_mode", spec.BillingMode),
		slog.Int("tier_count", len(spec.PricingTiers)),
	)

	// Create calculator from parsed spec
	calc := NewTieredPricingCalculator(spec.PricingTiers, spec.Currency)

	// Calculate cost for 75 TB
	usage := 75000.0 // 75 TB in GB
	totalCost, breakdown := calc.CalculateCost(usage)

	logger.Info("Cost for 75 TB storage",
		slog.Float64("total_cost", totalCost),
		slog.String("currency", spec.Currency),
		slog.Int("tiers_used", len(breakdown)),
	)

	for _, b := range breakdown {
		logger.Info("Tier breakdown",
			slog.String("tier", b.Tier.Description),
			slog.Float64("cost", b.TierCost),
		)
	}
}

// demonstrateGCPStorageTiers shows GCP Cloud Storage tiered pricing.
//
//nolint:mnd // pricing data is intentionally hard-coded for demonstration purposes
func demonstrateGCPStorageTiers(logger *slog.Logger) {
	// GCP Cloud Storage Standard pricing tiers (us-multi-region, as of 2024)
	tiers := []PricingTier{
		{
			MinQuantity: 0,
			MaxQuantity: 50000, // First 50 TB
			RatePerUnit: 0.020,
			Description: "First 50 TB / Month",
		},
		{
			MinQuantity: 50000,
			MaxQuantity: 450000, // Next 400 TB
			RatePerUnit: 0.019,
			Description: "Next 400 TB / Month",
		},
		{
			MinQuantity: 450000,
			MaxQuantity: 0, // Over 450 TB
			RatePerUnit: 0.018,
			Description: "Over 450 TB / Month",
		},
	}

	calc := NewTieredPricingCalculator(tiers, "USD")

	logger.Info("GCP Cloud Storage Standard Tiered Pricing Demo")

	// Compare cost for same usage across providers
	usage := 100000.0 // 100 TB
	gcpCost, gcpBreakdown := calc.CalculateCost(usage)

	logger.Info("GCP cost for 100 TB",
		slog.Float64("total_cost", gcpCost),
		slog.Int("tiers_used", len(gcpBreakdown)),
	)

	for _, b := range gcpBreakdown {
		logger.Info("Tier breakdown",
			slog.String("tier", b.Tier.Description),
			slog.Float64("usage_gb", b.UsageInTier),
			slog.Float64("rate", b.Tier.RatePerUnit),
			slog.Float64("cost", b.TierCost),
		)
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Check for quiet mode (for testing)
	quiet := len(os.Args) > 1 && os.Args[1] == "--quiet"
	if quiet {
		logger.Info("Tiered pricing examples executed successfully")
		return
	}

	demonstrateTieredPricing(logger)
	demonstratePricingSpecParsing(logger)
	demonstrateGCPStorageTiers(logger)

	logger.Info("Key Takeaways",
		slog.String("point_1", "Tiered pricing rewards higher usage with lower per-unit rates"),
		slog.String("point_2", "Always calculate cost per tier, not using a single rate"),
		slog.String("point_3", "Parse pricing_tiers from PricingSpec for accurate calculations"),
		slog.String("point_4", "Consider cross-provider comparison for cost optimization"),
	)
}
