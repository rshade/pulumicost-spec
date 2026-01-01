// Package pluginsdk_test contains unit tests for dry-run helpers.
//
// TEST-FIRST PROTOCOL: These tests define the expected contract for dry-run helpers.
// Per constitution III, tests MUST be written BEFORE implementation.
//
// Initial state: Tests will FAIL TO COMPILE because dry_run.go doesn't exist yet.
// After Phase 4 (Implementation): Tests will PASS.
package pluginsdk_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// T006: FocusFieldNames Helper Tests
// =============================================================================

// TestFocusFieldNames validates that FocusFieldNames() returns all expected
// FOCUS 1.2/1.3 field names.
//
// Expected behavior:
//   - Returns slice of ~50-66 FOCUS field names
//   - Contains all core fields (service_category, billed_cost, resource_id, etc.)
//   - Contains FOCUS 1.3 additions (service_provider_name, allocated_method_id, etc.)
//   - Field names match FocusCostRecord message field names
func TestFocusFieldNames(t *testing.T) {
	fields := pluginsdk.FocusFieldNames()

	// Verify we have a reasonable number of fields (FOCUS 1.2/1.3 has ~66 fields)
	if len(fields) < 40 {
		t.Errorf("Expected at least 40 FOCUS fields, got %d", len(fields))
	}

	// Create map for easy lookup
	fieldSet := make(map[string]bool)
	for _, f := range fields {
		fieldSet[f] = true
	}

	// Verify core FOCUS 1.2 fields are present
	coreFields := []string{
		// Identity & Hierarchy
		"provider_name",
		"billing_account_id",
		"sub_account_id",
		// Billing Period
		"billing_period_start",
		"billing_period_end",
		"billing_currency",
		// Charge Period
		"charge_period_start",
		"charge_period_end",
		// Charge Details
		"charge_category",
		"charge_class",
		"charge_description",
		// Service & Product
		"service_category",
		"service_name",
		// Resource Details
		"resource_id",
		"resource_name",
		"resource_type",
		// Location
		"region_id",
		"region_name",
		// Financial Amounts
		"billed_cost",
		"list_cost",
		"effective_cost",
		// Consumption/Usage
		"consumed_quantity",
		"consumed_unit",
		// Tags
		"tags",
	}

	for _, field := range coreFields {
		if !fieldSet[field] {
			t.Errorf("Missing core FOCUS field: %s", field)
		}
	}

	// Verify FOCUS 1.3 additions are present
	focus13Fields := []string{
		"service_provider_name",
		"host_provider_name",
		"allocated_method_id",
		"allocated_method_details",
		"allocated_resource_id",
		"allocated_resource_name",
		"allocated_tags",
		"contract_applied",
	}

	for _, field := range focus13Fields {
		if !fieldSet[field] {
			t.Errorf("Missing FOCUS 1.3 field: %s", field)
		}
	}
}

// TestFocusFieldNamesNoDuplicates verifies no duplicate field names.
func TestFocusFieldNamesNoDuplicates(t *testing.T) {
	fields := pluginsdk.FocusFieldNames()

	seen := make(map[string]bool)
	for _, f := range fields {
		if seen[f] {
			t.Errorf("Duplicate field name: %s", f)
		}
		seen[f] = true
	}
}

// =============================================================================
// NewFieldMapping Helper Tests
// =============================================================================

// TestNewFieldMapping validates the NewFieldMapping helper constructor.
func TestNewFieldMapping(t *testing.T) {
	fm := pluginsdk.NewFieldMapping("service_category", pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED)

	if fm.GetFieldName() != "service_category" {
		t.Errorf("Expected field_name='service_category', got %q", fm.GetFieldName())
	}

	if fm.GetSupportStatus() != pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED {
		t.Errorf("Expected status=SUPPORTED, got %v", fm.GetSupportStatus())
	}
}

// TestNewFieldMappingWithOptions validates optional fields.
func TestNewFieldMappingWithOptions(t *testing.T) {
	fm := pluginsdk.NewFieldMapping(
		"availability_zone",
		pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL,
		pluginsdk.WithConditionDescription("Only populated for regional resources"),
		pluginsdk.WithExpectedType("string"),
	)

	if fm.GetConditionDescription() != "Only populated for regional resources" {
		t.Errorf("Expected condition_description set, got %q", fm.GetConditionDescription())
	}

	if fm.GetExpectedType() != "string" {
		t.Errorf("Expected expected_type='string', got %q", fm.GetExpectedType())
	}
}

// =============================================================================
// DryRunResponse Builder Tests
// =============================================================================

// TestNewDryRunResponse validates the DryRunResponse builder.
func TestNewDryRunResponse(t *testing.T) {
	mappings := []*pbc.FieldMapping{
		pluginsdk.NewFieldMapping("service_category", pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED),
		pluginsdk.NewFieldMapping("billed_cost", pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED),
	}

	resp := pluginsdk.NewDryRunResponse(
		pluginsdk.WithFieldMappings(mappings),
		pluginsdk.WithResourceTypeSupported(true),
		pluginsdk.WithConfigurationValid(true),
	)

	if !resp.GetResourceTypeSupported() {
		t.Error("Expected resource_type_supported=true")
	}

	if !resp.GetConfigurationValid() {
		t.Error("Expected configuration_valid=true")
	}

	if len(resp.GetFieldMappings()) != 2 {
		t.Errorf("Expected 2 field mappings, got %d", len(resp.GetFieldMappings()))
	}
}

// TestNewDryRunResponseWithErrors validates error reporting.
func TestNewDryRunResponseWithErrors(t *testing.T) {
	resp := pluginsdk.NewDryRunResponse(
		pluginsdk.WithConfigurationValid(false),
		pluginsdk.WithConfigurationErrors([]string{
			"Missing API key",
			"Invalid endpoint URL",
		}),
	)

	if resp.GetConfigurationValid() {
		t.Error("Expected configuration_valid=false")
	}

	errors := resp.GetConfigurationErrors()
	if len(errors) != 2 {
		t.Errorf("Expected 2 configuration errors, got %d", len(errors))
	}
}

// =============================================================================
// All Fields Helper Tests
// =============================================================================

// TestAllFieldsMapping validates creating mappings for all FOCUS fields.
func TestAllFieldsMapping(t *testing.T) {
	// Create mappings for all fields with default SUPPORTED status
	mappings := pluginsdk.AllFieldsWithStatus(pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED)

	// Should have mapping for each field
	fieldCount := len(pluginsdk.FocusFieldNames())
	if len(mappings) != fieldCount {
		t.Errorf("Expected %d mappings, got %d", fieldCount, len(mappings))
	}

	// Verify all have SUPPORTED status
	for _, fm := range mappings {
		if fm.GetSupportStatus() != pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED {
			t.Errorf("Field %q has unexpected status: %v",
				fm.GetFieldName(), fm.GetSupportStatus())
		}
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

// BenchmarkFocusFieldNames measures FocusFieldNames() performance.
func BenchmarkFocusFieldNames(b *testing.B) {
	for range b.N {
		_ = pluginsdk.FocusFieldNames()
	}
}

// BenchmarkNewFieldMapping measures NewFieldMapping constructor performance.
func BenchmarkNewFieldMapping(b *testing.B) {
	for range b.N {
		_ = pluginsdk.NewFieldMapping("service_category", pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED)
	}
}

// BenchmarkAllFieldsWithStatus measures bulk field mapping creation.
func BenchmarkAllFieldsWithStatus(b *testing.B) {
	for range b.N {
		_ = pluginsdk.AllFieldsWithStatus(pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED)
	}
}
