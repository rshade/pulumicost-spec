package pluginsdk_test

import (
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func TestFocusRecordBuilder_Build_HappyPath(t *testing.T) {
	now := time.Now()
	billingStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingEnd := billingStart.AddDate(0, 1, 0)
	chargeStart := now.Add(-24 * time.Hour)
	chargeEnd := now

	builder := pluginsdk.NewFocusRecordBuilder()
	builder.WithIdentity("AWS", "acc-123", "My Account")
	builder.WithBillingPeriod(billingStart, billingEnd, "USD")
	builder.WithChargePeriod(chargeStart, chargeEnd)
	builder.WithChargeDetails(
		pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
		pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
	)
	builder.WithChargeClassification(
		pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
		"EC2 Instance Usage",
		pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
	)
	builder.WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2")
	builder.WithFinancials(10.5, 12.0, 10.0, "USD", "inv-001")
	builder.WithUsage(1.0, "Hour")

	record, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if record.GetBillingAccountId() != "acc-123" {
		t.Errorf("Expected BillingAccountId 'acc-123', got %s", record.GetBillingAccountId())
	}
	if record.GetBilledCost() != 10.5 {
		t.Errorf("Expected BilledCost 10.5, got %f", record.GetBilledCost())
	}
	if record.GetServiceCategory() != pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE {
		t.Errorf("Expected Compute category, got %v", record.GetServiceCategory())
	}
}

func TestFocusRecordBuilder_Build_MissingFields(t *testing.T) {
	builder := pluginsdk.NewFocusRecordBuilder()
	// Only set some fields
	builder.WithIdentity("provider-aws", "acc-123", "My Account")

	_, err := builder.Build()
	if err == nil {
		t.Error("Expected error due to missing fields, got nil")
	}
}

func TestFocusRecordBuilder_WithExtension(t *testing.T) {
	builder := createValidBuilder()

	builder.WithExtension("MyKey", "MyValue")
	builder.WithExtension("AnotherKey", "AnotherValue")

	record, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if record.GetExtendedColumns()["MyKey"] != "MyValue" {
		t.Errorf("Expected MyKey=MyValue, got %s", record.GetExtendedColumns()["MyKey"])
	}
	if record.GetExtendedColumns()["AnotherKey"] != "AnotherValue" {
		t.Errorf("Expected AnotherKey=AnotherValue, got %s", record.GetExtendedColumns()["AnotherKey"])
	}
}

// =============================================================================
// FOCUS 1.2 New Column Builder Tests (Phase 5: US3)
// =============================================================================

// createValidBuilder creates a builder with all mandatory fields set for testing new methods.
// Updated to include all 14 FOCUS 1.2 mandatory fields.
func createValidBuilder() *pluginsdk.FocusRecordBuilder {
	now := time.Now()
	billingStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingEnd := billingStart.AddDate(0, 1, 0)

	builder := pluginsdk.NewFocusRecordBuilder()
	builder.WithIdentity("AWS", "acc-123", "Account Name")
	builder.WithBillingPeriod(billingStart, billingEnd, "USD")
	builder.WithChargePeriod(now.Add(-time.Hour), now)
	builder.WithChargeDetails(
		pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
		pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
	)
	builder.WithChargeClassification(
		pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
		"EC2 Instance Usage",
		pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
	)
	builder.WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2")
	builder.WithFinancials(10.0, 12.0, 10.0, "USD", "inv-001")
	builder.WithUsage(1.0, "Hour")
	return builder
}

// TestFocusRecordBuilder_WithContractedCost tests the WithContractedCost builder method.
func TestFocusRecordBuilder_WithContractedCost(t *testing.T) {
	tests := []struct {
		name     string
		cost     float64
		expected float64
	}{
		{"zero cost", 0.0, 0.0},
		{"positive cost", 100.50, 100.50},
		{"large cost", 999999.99, 999999.99},
		{"small cost", 0.001, 0.001},
		{"negative cost (credit)", -50.25, -50.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithContractedCost(tt.cost)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetContractedCost() != tt.expected {
				t.Errorf("Expected ContractedCost %f, got %f", tt.expected, record.GetContractedCost())
			}
		})
	}
}

// TestFocusRecordBuilder_WithBillingAccountType tests the WithBillingAccountType builder method.
func TestFocusRecordBuilder_WithBillingAccountType(t *testing.T) {
	tests := []struct {
		name        string
		accountType string
		expected    string
	}{
		{"empty type", "", ""},
		{"enterprise", "Enterprise", "Enterprise"},
		{"pay as you go", "PayAsYouGo", "PayAsYouGo"},
		{"linked account", "LinkedAccount", "LinkedAccount"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithBillingAccountType(tt.accountType)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetBillingAccountType() != tt.expected {
				t.Errorf("Expected BillingAccountType %q, got %q", tt.expected, record.GetBillingAccountType())
			}
		})
	}
}

// TestFocusRecordBuilder_WithSubAccountType tests the WithSubAccountType builder method.
func TestFocusRecordBuilder_WithSubAccountType(t *testing.T) {
	tests := []struct {
		name        string
		accountType string
		expected    string
	}{
		{"empty type", "", ""},
		{"subscription", "Subscription", "Subscription"},
		{"project", "Project", "Project"},
		{"account", "Account", "Account"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithSubAccountType(tt.accountType)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetSubAccountType() != tt.expected {
				t.Errorf("Expected SubAccountType %q, got %q", tt.expected, record.GetSubAccountType())
			}
		})
	}
}

// TestFocusRecordBuilder_WithCapacityReservation tests the WithCapacityReservation builder method.
func TestFocusRecordBuilder_WithCapacityReservation(t *testing.T) {
	tests := []struct {
		name           string
		reservationID  string
		status         pbc.FocusCapacityReservationStatus
		expectedID     string
		expectedStatus pbc.FocusCapacityReservationStatus
	}{
		{
			"empty reservation",
			"",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_UNSPECIFIED,
			"",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_UNSPECIFIED,
		},
		{
			"used reservation",
			"cr-12345",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_USED,
			"cr-12345",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_USED,
		},
		{
			"unused reservation",
			"cr-67890",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_UNUSED,
			"cr-67890",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_UNUSED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithCapacityReservation(tt.reservationID, tt.status)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetCapacityReservationId() != tt.expectedID {
				t.Errorf("Expected CapacityReservationId %q, got %q", tt.expectedID, record.GetCapacityReservationId())
			}
			if record.GetCapacityReservationStatus() != tt.expectedStatus {
				t.Errorf(
					"Expected CapacityReservationStatus %v, got %v",
					tt.expectedStatus,
					record.GetCapacityReservationStatus(),
				)
			}
		})
	}
}

// TestFocusRecordBuilder_WithCommitmentDiscountDetails tests WithCommitmentDiscountDetails.
func TestFocusRecordBuilder_WithCommitmentDiscountDetails(t *testing.T) {
	tests := []struct {
		name             string
		quantity         float64
		status           pbc.FocusCommitmentDiscountStatus
		discountType     string
		unit             string
		expectedQuantity float64
		expectedStatus   pbc.FocusCommitmentDiscountStatus
		expectedType     string
		expectedUnit     string
	}{
		{
			"empty details",
			0.0,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_UNSPECIFIED,
			"",
			"",
			0.0,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_UNSPECIFIED,
			"",
			"",
		},
		{
			"used savings plan",
			100.0,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_USED,
			"SavingsPlan",
			"USD/Hour",
			100.0,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_USED,
			"SavingsPlan",
			"USD/Hour",
		},
		{
			"unused reserved instance",
			50.5,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_UNUSED,
			"ReservedInstance",
			"Hours",
			50.5,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_UNUSED,
			"ReservedInstance",
			"Hours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithCommitmentDiscountDetails(tt.quantity, tt.status, tt.discountType, tt.unit)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetCommitmentDiscountQuantity() != tt.expectedQuantity {
				t.Errorf(
					"Expected CommitmentDiscountQuantity %f, got %f",
					tt.expectedQuantity,
					record.GetCommitmentDiscountQuantity(),
				)
			}
			if record.GetCommitmentDiscountStatus() != tt.expectedStatus {
				t.Errorf(
					"Expected CommitmentDiscountStatus %v, got %v",
					tt.expectedStatus,
					record.GetCommitmentDiscountStatus(),
				)
			}
			if record.GetCommitmentDiscountType() != tt.expectedType {
				t.Errorf(
					"Expected CommitmentDiscountType %q, got %q",
					tt.expectedType,
					record.GetCommitmentDiscountType(),
				)
			}
			if record.GetCommitmentDiscountUnit() != tt.expectedUnit {
				t.Errorf(
					"Expected CommitmentDiscountUnit %q, got %q",
					tt.expectedUnit,
					record.GetCommitmentDiscountUnit(),
				)
			}
		})
	}
}

// TestFocusRecordBuilder_WithContractedUnitPrice tests the WithContractedUnitPrice builder method.
func TestFocusRecordBuilder_WithContractedUnitPrice(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		expected float64
	}{
		{"zero price", 0.0, 0.0},
		{"positive price", 0.0023, 0.0023},
		{"large price", 1500.99, 1500.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithContractedUnitPrice(tt.price)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetContractedUnitPrice() != tt.expected {
				t.Errorf("Expected ContractedUnitPrice %f, got %f", tt.expected, record.GetContractedUnitPrice())
			}
		})
	}
}

// TestFocusRecordBuilder_WithPricingCurrency tests the WithPricingCurrency builder method.
func TestFocusRecordBuilder_WithPricingCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		expected string
	}{
		{"empty currency", "", ""},
		{"USD", "USD", "USD"},
		{"EUR", "EUR", "EUR"},
		{"GBP", "GBP", "GBP"},
		{"JPY", "JPY", "JPY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithPricingCurrency(tt.currency)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetPricingCurrency() != tt.expected {
				t.Errorf("Expected PricingCurrency %q, got %q", tt.expected, record.GetPricingCurrency())
			}
		})
	}
}

// TestFocusRecordBuilder_WithPricingCurrencyPrices tests the WithPricingCurrencyPrices method.
func TestFocusRecordBuilder_WithPricingCurrencyPrices(t *testing.T) {
	tests := []struct {
		name                string
		contractedUnitPrice float64
		effectiveCost       float64
		listUnitPrice       float64
		expectedContracted  float64
		expectedEffective   float64
		expectedList        float64
	}{
		{"all zeros", 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{"all positive", 10.50, 100.00, 12.00, 10.50, 100.00, 12.00},
		{"mixed values", 0.023, 50.5, 0.025, 0.023, 50.5, 0.025},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithPricingCurrencyPrices(tt.contractedUnitPrice, tt.effectiveCost, tt.listUnitPrice)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetPricingCurrencyContractedUnitPrice() != tt.expectedContracted {
				t.Errorf(
					"Expected PricingCurrencyContractedUnitPrice %f, got %f",
					tt.expectedContracted,
					record.GetPricingCurrencyContractedUnitPrice(),
				)
			}
			if record.GetPricingCurrencyEffectiveCost() != tt.expectedEffective {
				t.Errorf(
					"Expected PricingCurrencyEffectiveCost %f, got %f",
					tt.expectedEffective,
					record.GetPricingCurrencyEffectiveCost(),
				)
			}
			if record.GetPricingCurrencyListUnitPrice() != tt.expectedList {
				t.Errorf(
					"Expected PricingCurrencyListUnitPrice %f, got %f",
					tt.expectedList,
					record.GetPricingCurrencyListUnitPrice(),
				)
			}
		})
	}
}

// TestFocusRecordBuilder_WithPublisher tests the WithPublisher builder method.
func TestFocusRecordBuilder_WithPublisher(t *testing.T) {
	tests := []struct {
		name      string
		publisher string
		expected  string
	}{
		{"empty publisher", "", ""},
		{"AWS", "Amazon Web Services", "Amazon Web Services"},
		{"Azure", "Microsoft", "Microsoft"},
		{"GCP", "Google", "Google"},
		{"third party", "Datadog", "Datadog"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithPublisher(tt.publisher)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetPublisher() != tt.expected {
				t.Errorf("Expected Publisher %q, got %q", tt.expected, record.GetPublisher())
			}
		})
	}
}

// TestFocusRecordBuilder_WithServiceSubcategory tests the WithServiceSubcategory builder method.
func TestFocusRecordBuilder_WithServiceSubcategory(t *testing.T) {
	tests := []struct {
		name        string
		subcategory string
		expected    string
	}{
		{"empty subcategory", "", ""},
		{"virtual machine", "Virtual Machine", "Virtual Machine"},
		{"container", "Container", "Container"},
		{"serverless", "Serverless", "Serverless"},
		{"kubernetes", "Kubernetes", "Kubernetes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithServiceSubcategory(tt.subcategory)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetServiceSubcategory() != tt.expected {
				t.Errorf("Expected ServiceSubcategory %q, got %q", tt.expected, record.GetServiceSubcategory())
			}
		})
	}
}

// TestFocusRecordBuilder_WithSkuDetails tests the WithSkuDetails builder method.
func TestFocusRecordBuilder_WithSkuDetails(t *testing.T) {
	tests := []struct {
		name            string
		meter           string
		priceDetails    string
		expectedMeter   string
		expectedDetails string
	}{
		{"empty details", "", "", "", ""},
		{"compute meter", "compute-hours", "On-Demand", "compute-hours", "On-Demand"},
		{"storage meter", "storage-gb-month", "Standard tier", "storage-gb-month", "Standard tier"},
		{"network meter", "data-transfer-gb", "Outbound", "data-transfer-gb", "Outbound"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := createValidBuilder()
			builder.WithSkuDetails(tt.meter, tt.priceDetails)
			record, err := builder.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}
			if record.GetSkuMeter() != tt.expectedMeter {
				t.Errorf("Expected SkuMeter %q, got %q", tt.expectedMeter, record.GetSkuMeter())
			}
			if record.GetSkuPriceDetails() != tt.expectedDetails {
				t.Errorf("Expected SkuPriceDetails %q, got %q", tt.expectedDetails, record.GetSkuPriceDetails())
			}
		})
	}
}

// TestFocusRecordBuilder_ChainNewMethods tests chaining multiple new builder methods.
func TestFocusRecordBuilder_ChainNewMethods(t *testing.T) {
	builder := createValidBuilder()

	// Chain all new methods
	record, err := builder.
		WithContractedCost(100.50).
		WithBillingAccountType("Enterprise").
		WithSubAccountType("Subscription").
		WithCapacityReservation("cr-123", pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_USED).
		WithCommitmentDiscountDetails(50.0, pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_USED, "SavingsPlan", "USD/Hour").
		WithContractedUnitPrice(0.05).
		WithPricingCurrency("EUR").
		WithPricingCurrencyPrices(0.045, 90.0, 0.055).
		WithPublisher("Amazon Web Services").
		WithServiceSubcategory("Virtual Machine").
		WithSkuDetails("compute-hours", "On-Demand").
		Build()

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify all fields were set
	if record.GetContractedCost() != 100.50 {
		t.Errorf("Expected ContractedCost 100.50, got %f", record.GetContractedCost())
	}
	if record.GetBillingAccountType() != "Enterprise" {
		t.Errorf("Expected BillingAccountType 'Enterprise', got %q", record.GetBillingAccountType())
	}
	if record.GetSubAccountType() != "Subscription" {
		t.Errorf("Expected SubAccountType 'Subscription', got %q", record.GetSubAccountType())
	}
	if record.GetCapacityReservationId() != "cr-123" {
		t.Errorf("Expected CapacityReservationId 'cr-123', got %q", record.GetCapacityReservationId())
	}
	if record.GetCommitmentDiscountQuantity() != 50.0 {
		t.Errorf("Expected CommitmentDiscountQuantity 50.0, got %f", record.GetCommitmentDiscountQuantity())
	}
	if record.GetContractedUnitPrice() != 0.05 {
		t.Errorf("Expected ContractedUnitPrice 0.05, got %f", record.GetContractedUnitPrice())
	}
	if record.GetPricingCurrency() != "EUR" {
		t.Errorf("Expected PricingCurrency 'EUR', got %q", record.GetPricingCurrency())
	}
	if record.GetPublisher() != "Amazon Web Services" {
		t.Errorf("Expected Publisher 'Amazon Web Services', got %q", record.GetPublisher())
	}
	if record.GetServiceSubcategory() != "Virtual Machine" {
		t.Errorf("Expected ServiceSubcategory 'Virtual Machine', got %q", record.GetServiceSubcategory())
	}
	if record.GetSkuMeter() != "compute-hours" {
		t.Errorf("Expected SkuMeter 'compute-hours', got %q", record.GetSkuMeter())
	}
}

// =============================================================================
// Benchmarks for Builder Operations
// =============================================================================

// BenchmarkNewFocusRecordBuilder measures builder creation performance.
func BenchmarkNewFocusRecordBuilder(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.NewFocusRecordBuilder()
	}
}

// BenchmarkFocusRecordBuilder_Build measures the full build cycle performance.
func BenchmarkFocusRecordBuilder_Build(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		builder := createValidBuilder()
		_, _ = builder.Build()
	}
}

// BenchmarkFocusRecordBuilder_ChainedBuild measures chained method performance.
func BenchmarkFocusRecordBuilder_ChainedBuild(b *testing.B) {
	now := time.Now()
	billingStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingEnd := billingStart.AddDate(0, 1, 0)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, _ = pluginsdk.NewFocusRecordBuilder().
			WithIdentity("AWS", "acc-123", "Account").
			WithBillingPeriod(billingStart, billingEnd, "USD").
			WithChargePeriod(now.Add(-time.Hour), now).
			WithChargeDetails(
				pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
				pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
			).
			WithChargeClassification(
				pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
				"Usage",
				pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
			).
			WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "EC2").
			WithFinancials(10.0, 12.0, 10.0, "USD", "inv-001").
			WithUsage(1.0, "Hour").
			Build()
	}
}

// BenchmarkFocusRecordBuilder_FullRecord measures building a complete 57-column record.
func BenchmarkFocusRecordBuilder_FullRecord(b *testing.B) {
	now := time.Now()
	billingStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingEnd := billingStart.AddDate(0, 1, 0)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		builder := pluginsdk.NewFocusRecordBuilder()

		// Mandatory columns
		builder.WithIdentity("AWS", "acc-123", "Account")
		builder.WithBillingPeriod(billingStart, billingEnd, "USD")
		builder.WithChargePeriod(now.Add(-time.Hour), now)
		builder.WithChargeDetails(
			pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
			pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
		)
		builder.WithChargeClassification(
			pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
			"Usage",
			pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
		)
		builder.WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "EC2")
		builder.WithFinancials(10.0, 12.0, 10.0, "USD", "inv-001")
		builder.WithUsage(1.0, "Hour")

		// Conditional columns
		builder.WithSubAccount("sub-123", "SubAccount")
		builder.WithBillingAccountType("Enterprise")
		builder.WithSubAccountType("Subscription")
		builder.WithLocation("us-east-1", "US East", "us-east-1a")
		builder.WithResource("i-123", "instance", "m5.large")
		builder.WithSKU("sku-123", "price-123")
		builder.WithSkuDetails("meter", "details")
		builder.WithPricing(100.0, "Hours", 0.10)
		builder.WithContractedCost(10.0)
		builder.WithContractedUnitPrice(0.10)
		builder.WithPricingCurrency("USD")
		builder.WithPricingCurrencyPrices(0.10, 10.0, 0.12)
		builder.WithPublisher("AWS")
		builder.WithServiceSubcategory("VM")
		builder.WithCommitmentDiscount(
			pbc.FocusCommitmentDiscountCategory_FOCUS_COMMITMENT_DISCOUNT_CATEGORY_USAGE,
			"ri-123", "Reserved Instance",
		)
		builder.WithCommitmentDiscountDetails(
			100.0,
			pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_USED,
			"Standard RI", "Hours",
		)
		builder.WithCapacityReservation(
			"cr-123",
			pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_USED,
		)
		builder.WithInvoice("inv-001", "AWS Inc")
		builder.WithTags(map[string]string{"env": "prod", "team": "platform"})
		builder.WithExtension("custom", "value")

		_, _ = builder.Build()
	}
}
