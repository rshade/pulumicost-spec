//go:build ignore

// audit_focus_columns.go verifies that all 57 FOCUS 1.2 columns are implemented
// in the FocusCostRecord proto message.
//
// Usage: go run scripts/audit_focus_columns.go
//
// Reference: https://focus.finops.org/focus-specification/v1-2/
package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Column represents a FOCUS 1.2 column specification.
type Column struct {
	Name       string // FOCUS column name (PascalCase)
	ProtoField string // Expected proto field name (snake_case)
	Level      string // Mandatory, Recommended, or Conditional
	Category   string // Column category for grouping
}

// FOCUS 1.2 columns: 14 mandatory, 1 recommended, 42 conditional = 57 total
var focus12Columns = []Column{
	// MANDATORY (14)
	{"BilledCost", "billed_cost", "Mandatory", "Financial"},
	{"BillingAccountId", "billing_account_id", "Mandatory", "Identity"},
	{"BillingAccountName", "billing_account_name", "Mandatory", "Identity"},
	{"BillingCurrency", "billing_currency", "Mandatory", "Financial"},
	{"BillingPeriodEnd", "billing_period_end", "Mandatory", "Billing Period"},
	{"BillingPeriodStart", "billing_period_start", "Mandatory", "Billing Period"},
	{"ChargeCategory", "charge_category", "Mandatory", "Charge"},
	{"ChargeClass", "charge_class", "Mandatory", "Charge"},
	{"ChargeDescription", "charge_description", "Mandatory", "Charge"},
	{"ChargePeriodEnd", "charge_period_end", "Mandatory", "Charge Period"},
	{"ChargePeriodStart", "charge_period_start", "Mandatory", "Charge Period"},
	{"ContractedCost", "contracted_cost", "Mandatory", "Financial"},
	{"ProviderName", "provider_name", "Mandatory", "Identity"},
	{"ServiceName", "service_name", "Mandatory", "Service"},

	// RECOMMENDED (1)
	{"ChargeFrequency", "charge_frequency", "Recommended", "Charge"},

	// CONDITIONAL (42)
	{"AvailabilityZone", "availability_zone", "Conditional", "Location"},
	{"BillingAccountType", "billing_account_type", "Conditional", "Identity"},
	{"CapacityReservationId", "capacity_reservation_id", "Conditional", "Capacity Reservation"},
	{"CapacityReservationStatus", "capacity_reservation_status", "Conditional", "Capacity Reservation"},
	{"CommitmentDiscountCategory", "commitment_discount_category", "Conditional", "Commitment Discount"},
	{"CommitmentDiscountId", "commitment_discount_id", "Conditional", "Commitment Discount"},
	{"CommitmentDiscountName", "commitment_discount_name", "Conditional", "Commitment Discount"},
	{"CommitmentDiscountQuantity", "commitment_discount_quantity", "Conditional", "Commitment Discount"},
	{"CommitmentDiscountStatus", "commitment_discount_status", "Conditional", "Commitment Discount"},
	{"CommitmentDiscountType", "commitment_discount_type", "Conditional", "Commitment Discount"},
	{"CommitmentDiscountUnit", "commitment_discount_unit", "Conditional", "Commitment Discount"},
	{"ConsumedQuantity", "consumed_quantity", "Conditional", "Usage"},
	{"ConsumedUnit", "consumed_unit", "Conditional", "Usage"},
	{"ContractedUnitPrice", "contracted_unit_price", "Conditional", "Financial"},
	{"EffectiveCost", "effective_cost", "Conditional", "Financial"},
	{"InvoiceId", "invoice_id", "Conditional", "Invoice"},
	{"InvoiceIssuer", "invoice_issuer", "Conditional", "Invoice"},
	{"ListCost", "list_cost", "Conditional", "Financial"},
	{"ListUnitPrice", "list_unit_price", "Conditional", "Pricing"},
	{"PricingCategory", "pricing_category", "Conditional", "Pricing"},
	{"PricingCurrency", "pricing_currency", "Conditional", "Pricing"},
	{"PricingCurrencyContractedUnitPrice", "pricing_currency_contracted_unit_price", "Conditional", "Pricing"},
	{"PricingCurrencyEffectiveCost", "pricing_currency_effective_cost", "Conditional", "Pricing"},
	{"PricingCurrencyListUnitPrice", "pricing_currency_list_unit_price", "Conditional", "Pricing"},
	{"PricingQuantity", "pricing_quantity", "Conditional", "Pricing"},
	{"PricingUnit", "pricing_unit", "Conditional", "Pricing"},
	{"Publisher", "publisher", "Conditional", "Service"},
	{"RegionId", "region_id", "Conditional", "Location"},
	{"RegionName", "region_name", "Conditional", "Location"},
	{"ResourceId", "resource_id", "Conditional", "Resource"},
	{"ResourceName", "resource_name", "Conditional", "Resource"},
	{"ResourceType", "resource_type", "Conditional", "Resource"},
	{"ServiceCategory", "service_category", "Conditional", "Service"},
	{"ServiceSubcategory", "service_subcategory", "Conditional", "Service"},
	{"SkuId", "sku_id", "Conditional", "SKU"},
	{"SkuMeter", "sku_meter", "Conditional", "SKU"},
	{"SkuPriceDetails", "sku_price_details", "Conditional", "SKU"},
	{"SkuPriceId", "sku_price_id", "Conditional", "SKU"},
	{"SubAccountId", "sub_account_id", "Conditional", "Identity"},
	{"SubAccountName", "sub_account_name", "Conditional", "Identity"},
	{"SubAccountType", "sub_account_type", "Conditional", "Identity"},
	{"Tags", "tags", "Conditional", "Metadata"},
}

func main() {
	fmt.Println("FOCUS 1.2 Column Audit")
	fmt.Println("======================")
	fmt.Println()

	recordType := reflect.TypeOf(pbc.FocusCostRecord{})

	var missing []Column
	var present []Column
	mandatoryMissing := 0

	for _, col := range focus12Columns {
		goFieldName := snakeToPascal(col.ProtoField)
		if _, found := recordType.FieldByName(goFieldName); found {
			present = append(present, col)
		} else {
			missing = append(missing, col)
			if col.Level == "Mandatory" {
				mandatoryMissing++
			}
		}
	}

	// Summary
	fmt.Printf("Total columns: %d\n", len(focus12Columns))
	fmt.Printf("Implemented: %d\n", len(present))
	fmt.Printf("Missing: %d\n", len(missing))
	fmt.Println()

	// Breakdown by level
	mandatoryCount := 0
	recommendedCount := 0
	conditionalCount := 0
	for _, col := range focus12Columns {
		switch col.Level {
		case "Mandatory":
			mandatoryCount++
		case "Recommended":
			recommendedCount++
		case "Conditional":
			conditionalCount++
		}
	}

	mandatoryPresent := 0
	recommendedPresent := 0
	for _, col := range present {
		switch col.Level {
		case "Mandatory":
			mandatoryPresent++
		case "Recommended":
			recommendedPresent++
		}
	}

	conditionalPresent := len(present) - mandatoryPresent - recommendedPresent
	fmt.Printf("Mandatory: %d/%d\n", mandatoryPresent, mandatoryCount)
	fmt.Printf("Recommended: %d/%d\n", recommendedPresent, recommendedCount)
	fmt.Printf("Conditional: %d/%d\n", conditionalPresent, conditionalCount)
	fmt.Println()

	if len(missing) > 0 {
		fmt.Println("❌ Missing columns:")
		for _, col := range missing {
			fmt.Printf("  - %s (%s) [%s]\n", col.Name, col.ProtoField, col.Level)
		}
		fmt.Println()

		if mandatoryMissing > 0 {
			fmt.Printf("⚠️  %d MANDATORY column(s) missing!\n", mandatoryMissing)
		}

		os.Exit(1)
	}

	fmt.Println("✅ All FOCUS 1.2 columns are implemented")
	os.Exit(0)
}

// snakeToPascal converts snake_case to PascalCase for Go struct field lookup.
func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
