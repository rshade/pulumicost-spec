package pluginsdk_test

import (
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func TestFocusRecordBuilder_Build_HappyPath(t *testing.T) {
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	builder := pluginsdk.NewFocusRecordBuilder()
	builder.WithIdentity("provider-aws", "acc-123", "My Account")
	builder.WithChargePeriod(start, end)
	builder.WithServiceCategory(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE)
	builder.WithChargeDetails(
		pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
		pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
	)
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
	builder := pluginsdk.NewFocusRecordBuilder()
	// Set mandatory fields
	builder.WithIdentity("p", "acc", "n")
	builder.WithChargePeriod(time.Now(), time.Now())
	builder.WithServiceCategory(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE)
	builder.WithChargeDetails(
		pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
		pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
	)
	builder.WithFinancials(1.0, 1.0, 1.0, "USD", "i")
	builder.WithUsage(1.0, "unit")

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
