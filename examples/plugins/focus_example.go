package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Example cost values for demonstration.
const (
	exampleBilledCost    = 0.12
	exampleListCost      = 0.12
	exampleEffectiveCost = 0.12
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Demonstrating full Builder usage
	builder := pluginsdk.NewFocusRecordBuilder()

	start := time.Now()
	end := start.Add(1 * time.Hour)

	builder.WithIdentity("aws", "123456789012", "Production")
	builder.WithChargePeriod(start, end)
	builder.WithServiceCategory(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE)
	builder.WithChargeDetails(
		pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
		pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
	)
	builder.WithFinancials(exampleBilledCost, exampleListCost, exampleEffectiveCost, "USD", "inv-2025-11")
	builder.WithUsage(1.0, "Hour")

	// Backpack usage (extensions)
	builder.WithExtension("Environment", "Prod")
	builder.WithExtension("CostCenter", "CC-999")

	record, err := builder.Build()
	if err != nil {
		logger.Error("Failed to build record", "error", err)
		os.Exit(1)
	}

	logger.Info("Generated Record",
		"provider", record.GetProviderName(),
		"cost", record.GetBilledCost(),
		"currency", record.GetCurrency(),
		"extensions", record.GetExtendedColumns(),
	)
}
