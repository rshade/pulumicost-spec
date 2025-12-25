package testing_test

import (
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// FOCUS 1.3 Conformance Tests
// =============================================================================
//
// These tests validate backward compatibility and correct behavior of
// FOCUS 1.3 extensions in the pluginsdk package.

// buildValidFocusRecord creates a minimal valid FOCUS record with all required fields.
// This helper ensures tests focus on FOCUS 1.3 specific behavior.
func buildValidFocusRecord() *pluginsdk.FocusRecordBuilder {
	now := time.Now()
	billingStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	billingEnd := billingStart.AddDate(0, 1, 0)
	chargeStart := now.Add(-24 * time.Hour)
	chargeEnd := now

	return pluginsdk.NewFocusRecordBuilder().
		WithIdentity("AWS", "123456789012", "Production").
		WithBillingPeriod(billingStart, billingEnd, "USD").
		WithChargePeriod(chargeStart, chargeEnd).
		WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2").
		WithChargeDetails(
			pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
			pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
		).
		WithChargeClassification(
			pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
			"EC2 usage",
			pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
		).
		WithUsage(100, "Hours").
		WithFinancials(100.00, 100.00, 95.00, "USD", "INV-2024-001").
		WithContractedCost(95.00)
}

// TestFocus13_BackwardCompatibility_FOCUS12OnlyRecords validates that plugins
// can create valid FOCUS 1.2 records without any FOCUS 1.3 fields (T060, T061).
func TestFocus13_BackwardCompatibility_FOCUS12OnlyRecords(t *testing.T) {
	// Create a pure FOCUS 1.2 record with no FOCUS 1.3 fields
	record, err := buildValidFocusRecord().Build()

	if err != nil {
		t.Fatalf("FOCUS 1.2 record creation failed: %v", err)
	}

	// Verify FOCUS 1.3 fields are empty/zero (default values)
	if record.GetServiceProviderName() != "" {
		t.Errorf("Expected empty service_provider_name, got %q", record.GetServiceProviderName())
	}
	if record.GetHostProviderName() != "" {
		t.Errorf("Expected empty host_provider_name, got %q", record.GetHostProviderName())
	}
	if record.GetAllocatedMethodId() != "" {
		t.Errorf("Expected empty allocated_method_id, got %q", record.GetAllocatedMethodId())
	}
	if record.GetAllocatedResourceId() != "" {
		t.Errorf("Expected empty allocated_resource_id, got %q", record.GetAllocatedResourceId())
	}
	if record.GetContractApplied() != "" {
		t.Errorf("Expected empty contract_applied, got %q", record.GetContractApplied())
	}
	if len(record.GetAllocatedTags()) != 0 {
		t.Errorf("Expected empty allocated_tags, got %v", record.GetAllocatedTags())
	}

	// Verify FOCUS 1.2 fields are correctly populated
	//nolint:staticcheck // SA1019: Testing deprecated provider_name backward compatibility
	if record.GetProviderName() != "AWS" {
		//nolint:staticcheck // SA1019: Testing deprecated provider_name backward compatibility
		t.Errorf("Expected provider_name='AWS', got %q", record.GetProviderName())
	}
	if record.GetBillingAccountId() != "123456789012" {
		t.Errorf("Expected billing_account_id='123456789012', got %q", record.GetBillingAccountId())
	}
}

// TestFocus13_BackwardCompatibility_DeprecatedFieldsWork validates that
// deprecated fields (provider_name, publisher) continue to function (T062).
func TestFocus13_BackwardCompatibility_DeprecatedFieldsWork(t *testing.T) {
	t.Run("provider_name continues to work", func(t *testing.T) {
		record, err := buildValidFocusRecord().Build()

		if err != nil {
			t.Fatalf("Record creation failed: %v", err)
		}

		//nolint:staticcheck // SA1019: Testing deprecated provider_name backward compatibility
		if record.GetProviderName() != "AWS" {
			//nolint:staticcheck // SA1019: Testing deprecated provider_name backward compatibility
			t.Errorf("Expected provider_name='AWS', got %q", record.GetProviderName())
		}
	})

	t.Run("publisher continues to work", func(t *testing.T) {
		//nolint:staticcheck // SA1019: Testing deprecated WithPublisher backward compatibility
		record, err := buildValidFocusRecord().
			WithPublisher("AWS Inc.").
			Build()

		if err != nil {
			t.Fatalf("Record creation failed: %v", err)
		}

		//nolint:staticcheck // SA1019: Testing deprecated publisher backward compatibility
		if record.GetPublisher() != "AWS Inc." {
			//nolint:staticcheck // SA1019: Testing deprecated publisher backward compatibility
			t.Errorf("Expected publisher='AWS Inc.', got %q", record.GetPublisher())
		}
	})
}

// TestFocus13_NewFieldsDefaultValues validates that all new FOCUS 1.3 fields
// have sensible default values that don't affect existing behavior (T063).
func TestFocus13_NewFieldsDefaultValues(t *testing.T) {
	record, err := buildValidFocusRecord().Build()

	if err != nil {
		t.Fatalf("Record creation with defaults failed: %v", err)
	}

	// All FOCUS 1.3 string fields should default to empty string
	focus13StringFields := map[string]string{
		"service_provider_name":    record.GetServiceProviderName(),
		"host_provider_name":       record.GetHostProviderName(),
		"allocated_method_id":      record.GetAllocatedMethodId(),
		"allocated_method_details": record.GetAllocatedMethodDetails(),
		"allocated_resource_id":    record.GetAllocatedResourceId(),
		"allocated_resource_name":  record.GetAllocatedResourceName(),
		"contract_applied":         record.GetContractApplied(),
	}

	for fieldName, value := range focus13StringFields {
		if value != "" {
			t.Errorf("FOCUS 1.3 field %s should default to empty, got %q", fieldName, value)
		}
	}

	// AllocatedTags should default to nil or empty map
	if tags := record.GetAllocatedTags(); len(tags) > 0 {
		t.Errorf("allocated_tags should be empty by default, got %v", tags)
	}
}

// TestFocus13_AllocationFields validates the new cost allocation builder methods.
func TestFocus13_AllocationFields(t *testing.T) {
	t.Run("allocation fields set correctly", func(t *testing.T) {
		record, err := buildValidFocusRecord().
			WithAllocation("proportional-cpu", "Allocated by CPU usage percentage").
			WithAllocatedResource("workload-frontend", "Frontend Application").
			WithAllocatedTags(map[string]string{
				"team":        "frontend",
				"environment": "production",
			}).
			Build()

		if err != nil {
			t.Fatalf("Record creation failed: %v", err)
		}

		if record.GetAllocatedMethodId() != "proportional-cpu" {
			t.Errorf("Expected allocated_method_id='proportional-cpu', got %q",
				record.GetAllocatedMethodId())
		}
		if record.GetAllocatedMethodDetails() != "Allocated by CPU usage percentage" {
			t.Errorf("Expected allocated_method_details='Allocated by CPU usage percentage', got %q",
				record.GetAllocatedMethodDetails())
		}
		if record.GetAllocatedResourceId() != "workload-frontend" {
			t.Errorf("Expected allocated_resource_id='workload-frontend', got %q",
				record.GetAllocatedResourceId())
		}
		if record.GetAllocatedResourceName() != "Frontend Application" {
			t.Errorf("Expected allocated_resource_name='Frontend Application', got %q",
				record.GetAllocatedResourceName())
		}

		tags := record.GetAllocatedTags()
		if tags["team"] != "frontend" {
			t.Errorf("Expected allocated_tags['team']='frontend', got %q", tags["team"])
		}
		if tags["environment"] != "production" {
			t.Errorf("Expected allocated_tags['environment']='production', got %q", tags["environment"])
		}
	})
}

// TestFocus13_ProviderFields validates service_provider_name and host_provider_name.
func TestFocus13_ProviderFields(t *testing.T) {
	t.Run("marketplace scenario", func(t *testing.T) {
		record, err := buildValidFocusRecord().
			WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_MANAGEMENT, "Datadog APM").
			WithServiceProvider("Datadog"). // ISV selling via marketplace
			WithHostProvider("AWS").        // Cloud platform hosting
			Build()

		if err != nil {
			t.Fatalf("Marketplace record creation failed: %v", err)
		}

		if record.GetServiceProviderName() != "Datadog" {
			t.Errorf("Expected service_provider_name='Datadog', got %q",
				record.GetServiceProviderName())
		}
		if record.GetHostProviderName() != "AWS" {
			t.Errorf("Expected host_provider_name='AWS', got %q",
				record.GetHostProviderName())
		}
	})

	t.Run("same provider scenario", func(t *testing.T) {
		record, err := buildValidFocusRecord().
			WithServiceProvider("AWS"). // AWS is both service and host provider
			WithHostProvider("AWS").
			Build()

		if err != nil {
			t.Fatalf("Same provider record creation failed: %v", err)
		}

		if record.GetServiceProviderName() != "AWS" {
			t.Errorf("Expected service_provider_name='AWS', got %q",
				record.GetServiceProviderName())
		}
		if record.GetHostProviderName() != "AWS" {
			t.Errorf("Expected host_provider_name='AWS', got %q",
				record.GetHostProviderName())
		}
	})
}

// TestFocus13_ContractApplied validates the ContractApplied field linking.
func TestFocus13_ContractApplied(t *testing.T) {
	t.Run("links to contract commitment", func(t *testing.T) {
		// First create a contract commitment
		commitment, err := pluginsdk.NewContractCommitmentBuilder().
			WithIdentity("commit-ri-001", "contract-ea-2025").
			WithCategory(pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND).
			WithType("Reserved Instance").
			WithFinancials(10000.00, 0, "", "USD").
			Build()

		if err != nil {
			t.Fatalf("Contract commitment creation failed: %v", err)
		}

		// Then link a cost record to it
		record, err := buildValidFocusRecord().
			WithContractApplied(commitment.GetContractCommitmentId()). // Link to commitment
			Build()

		if err != nil {
			t.Fatalf("Record with contract applied failed: %v", err)
		}

		if record.GetContractApplied() != "commit-ri-001" {
			t.Errorf("Expected contract_applied='commit-ri-001', got %q",
				record.GetContractApplied())
		}
	})

	t.Run("accepts any string as opaque reference", func(t *testing.T) {
		// ContractApplied accepts any string - no validation against commitment dataset
		record, err := buildValidFocusRecord().
			WithContractApplied("any-arbitrary-commitment-id"). // No validation
			Build()

		if err != nil {
			t.Fatalf("Record with arbitrary contract_applied failed: %v", err)
		}

		if record.GetContractApplied() != "any-arbitrary-commitment-id" {
			t.Errorf("Expected contract_applied='any-arbitrary-commitment-id', got %q",
				record.GetContractApplied())
		}
	})
}

// TestFocus13_ContractCommitmentBuilder_SPEND validates SPEND category commitments.
func TestFocus13_ContractCommitmentBuilder_SPEND(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	commitment, err := pluginsdk.NewContractCommitmentBuilder().
		WithIdentity("commit-sp-001", "aws-ea-2025").
		WithCategory(pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND).
		WithType("Savings Plan").
		WithCommitmentPeriod(start, end).
		WithContractPeriod(start, end.AddDate(2, 0, 0)).
		WithFinancials(100000.00, 0, "", "USD").
		Build()

	if err != nil {
		t.Fatalf("SPEND commitment creation failed: %v", err)
	}

	if commitment.GetContractCommitmentId() != "commit-sp-001" {
		t.Errorf("Expected contract_commitment_id='commit-sp-001', got %q",
			commitment.GetContractCommitmentId())
	}
	if commitment.GetContractId() != "aws-ea-2025" {
		t.Errorf("Expected contract_id='aws-ea-2025', got %q",
			commitment.GetContractId())
	}
	if commitment.GetContractCommitmentCategory() !=
		pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND {
		t.Errorf("Expected SPEND category, got %v",
			commitment.GetContractCommitmentCategory())
	}
	if commitment.GetContractCommitmentCost() != 100000.00 {
		t.Errorf("Expected cost=100000.00, got %f",
			commitment.GetContractCommitmentCost())
	}
}

// TestFocus13_ContractCommitmentBuilder_USAGE validates USAGE category commitments.
func TestFocus13_ContractCommitmentBuilder_USAGE(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	commitment, err := pluginsdk.NewContractCommitmentBuilder().
		WithIdentity("commit-cud-001", "gcp-cud-2025").
		WithCategory(pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_USAGE).
		WithType("Committed Use Discount").
		WithCommitmentPeriod(start, end).
		WithFinancials(0, 10000, "vCPU-Hours", "USD").
		Build()

	if err != nil {
		t.Fatalf("USAGE commitment creation failed: %v", err)
	}

	if commitment.GetContractCommitmentCategory() !=
		pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_USAGE {
		t.Errorf("Expected USAGE category, got %v",
			commitment.GetContractCommitmentCategory())
	}
	if commitment.GetContractCommitmentQuantity() != 10000 {
		t.Errorf("Expected quantity=10000, got %f",
			commitment.GetContractCommitmentQuantity())
	}
	if commitment.GetContractCommitmentUnit() != "vCPU-Hours" {
		t.Errorf("Expected unit='vCPU-Hours', got %q",
			commitment.GetContractCommitmentUnit())
	}
}

// TestFocus13_ContractCommitmentBuilder_Validation validates builder error handling.
func TestFocus13_ContractCommitmentBuilder_Validation(t *testing.T) {
	t.Run("missing required fields", func(t *testing.T) {
		_, err := pluginsdk.NewContractCommitmentBuilder().
			WithFinancials(100.0, 0, "", "USD").
			Build()

		if err == nil {
			t.Error("Expected error for missing contract_commitment_id")
		}
	})

	t.Run("invalid currency", func(t *testing.T) {
		_, err := pluginsdk.NewContractCommitmentBuilder().
			WithIdentity("commit-001", "contract-001").
			WithFinancials(100.0, 0, "", "ZZZ"). // Invalid currency
			Build()

		if err == nil {
			t.Error("Expected error for invalid currency")
		}
	})

	t.Run("invalid period", func(t *testing.T) {
		start := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

		_, err := pluginsdk.NewContractCommitmentBuilder().
			WithIdentity("commit-001", "contract-001").
			WithCommitmentPeriod(start, end). // End before start
			WithFinancials(100.0, 0, "", "USD").
			Build()

		if err == nil {
			t.Error("Expected error for invalid commitment period")
		}
	})

	t.Run("negative cost", func(t *testing.T) {
		_, err := pluginsdk.NewContractCommitmentBuilder().
			WithIdentity("commit-001", "contract-001").
			WithCategory(pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND).
			WithFinancials(-100.0, 0, "", "USD"). // Negative cost
			Build()

		if err == nil {
			t.Error("Expected error for negative cost")
		}
	})

	t.Run("negative quantity", func(t *testing.T) {
		_, err := pluginsdk.NewContractCommitmentBuilder().
			WithIdentity("commit-001", "contract-001").
			WithCategory(pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_USAGE).
			WithFinancials(0, -100.0, "vCPU-Hours", "USD"). // Negative quantity
			Build()

		if err == nil {
			t.Error("Expected error for negative quantity")
		}
	})
}

// TestFocus13_MixedVersionRecords validates records with both FOCUS 1.2 and 1.3 fields.
func TestFocus13_MixedVersionRecords(t *testing.T) {
	t.Run("all FOCUS 1.3 fields populated", func(t *testing.T) {
		record, err := buildValidFocusRecord().
			// FOCUS 1.3 fields
			WithServiceProvider("AWS").
			WithHostProvider("AWS").
			WithAllocation("proportional-cost", "Cost split by resource usage").
			WithAllocatedResource("app-backend", "Backend Service").
			WithAllocatedTags(map[string]string{
				"cost-center": "engineering",
				"project":     "infrastructure",
			}).
			WithContractApplied("commit-ri-backend-001").
			Build()

		if err != nil {
			t.Fatalf("Mixed version record creation failed: %v", err)
		}

		// Verify both 1.2 and 1.3 fields are populated
		//nolint:staticcheck // SA1019: Testing deprecated provider_name backward compatibility
		if record.GetProviderName() != "AWS" {
			t.Error("FOCUS 1.2 provider_name not set")
		}
		if record.GetServiceProviderName() != "AWS" {
			t.Error("FOCUS 1.3 service_provider_name not set")
		}
		if record.GetAllocatedMethodId() != "proportional-cost" {
			t.Error("FOCUS 1.3 allocated_method_id not set")
		}
		if record.GetContractApplied() != "commit-ri-backend-001" {
			t.Error("FOCUS 1.3 contract_applied not set")
		}
	})
}
