// Package testing_test contains conformance tests for the DryRun RPC capability.
//
// TEST-FIRST PROTOCOL: These tests define the expected contract for the DryRun feature.
// Per constitution III, tests MUST be written BEFORE implementation.
//
// Initial state: Tests will FAIL TO COMPILE because proto types don't exist yet.
// After Phase 3 (Proto Definitions): Tests will compile but FAIL (no implementation).
// After Phase 4+ (Implementation): Tests will PASS.
package testing_test

import (
	"context"
	"testing"
	"time"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// =============================================================================
// T003: DryRun RPC Basic Functionality
// =============================================================================

// TestDryRunBasicFunctionality validates that the DryRun RPC returns expected
// field mapping information for a supported resource type.
//
// Expected behavior:
//   - Plugin returns DryRunResponse with field_mappings populated
//   - response.resource_type_supported = true for supported resources
//   - response.configuration_valid = true for valid configuration
//   - Response time < 100ms (no external API calls)
func TestDryRunBasicFunctionality(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a DryRunRequest for a supported resource type
	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	start := time.Now()
	resp, err := harness.Client().DryRun(ctx, req)
	duration := time.Since(start)

	// Verify no error
	if err != nil {
		t.Fatalf("DryRun RPC failed: %v", err)
	}

	// Verify response time < 100ms (dry-run should not make external calls)
	if duration > 100*time.Millisecond {
		t.Errorf("DryRun exceeded 100ms response time: %v", duration)
	}

	// Verify resource type is supported
	if !resp.GetResourceTypeSupported() {
		t.Error("Expected resource_type_supported=true for aws/ec2")
	}

	// Verify configuration is valid
	if !resp.GetConfigurationValid() {
		t.Errorf("Expected configuration_valid=true, got errors: %v", resp.GetConfigurationErrors())
	}

	// Verify field mappings are populated
	if len(resp.GetFieldMappings()) == 0 {
		t.Error("Expected non-empty field_mappings for supported resource")
	}

	// Verify at least some core FOCUS fields are present
	fieldNames := make(map[string]bool)
	for _, fm := range resp.GetFieldMappings() {
		fieldNames[fm.GetFieldName()] = true
	}

	coreFields := []string{"service_category", "billed_cost", "resource_id", "provider_name"}
	for _, field := range coreFields {
		if !fieldNames[field] {
			t.Errorf("Expected core FOCUS field %q in field_mappings", field)
		}
	}
}

// =============================================================================
// T004: Unsupported Resource Type Behavior
// =============================================================================

// TestDryRunUnsupportedResourceType validates behavior when querying an
// unsupported resource type.
//
// Expected behavior:
//   - Plugin returns DryRunResponse (not an error)
//   - response.resource_type_supported = false
//   - response.field_mappings may be empty
//   - response.configuration_valid = true (config itself is valid)
func TestDryRunUnsupportedResourceType(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a DryRunRequest for an unsupported resource type
	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "unsupported_resource_xyz",
			Region:       "us-east-1",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)

	// Should not return an error - just indicate unsupported
	if err != nil {
		t.Fatalf("DryRun should not error for unsupported resource, got: %v", err)
	}

	// Verify resource type is NOT supported
	if resp.GetResourceTypeSupported() {
		t.Error("Expected resource_type_supported=false for unsupported resource")
	}

	// Configuration should still be valid (the plugin itself is configured correctly)
	if !resp.GetConfigurationValid() {
		t.Error("Expected configuration_valid=true even for unsupported resource")
	}
}

// =============================================================================
// T005: FieldSupportStatus Enum Usage
// =============================================================================

// TestDryRunFieldSupportStatusEnumValues validates that FieldSupportStatus enum
// values are used correctly in field mappings.
//
// Expected enum values:
//   - FIELD_SUPPORT_STATUS_UNSPECIFIED (0) - default/unknown
//   - FIELD_SUPPORT_STATUS_SUPPORTED (1) - always populated
//   - FIELD_SUPPORT_STATUS_UNSUPPORTED (2) - never populated
//   - FIELD_SUPPORT_STATUS_CONDITIONAL (3) - depends on resource config
//   - FIELD_SUPPORT_STATUS_DYNAMIC (4) - requires runtime data
func TestDryRunFieldSupportStatusEnumValues(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)
	if err != nil {
		t.Fatalf("DryRun RPC failed: %v", err)
	}

	// Verify enum values are valid (not UNSPECIFIED for actual fields)
	for _, fm := range resp.GetFieldMappings() {
		status := fm.GetSupportStatus()

		// Verify status is a valid enum value
		switch status {
		case pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED,
			pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_UNSUPPORTED,
			pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL,
			pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_DYNAMIC:
			// Valid status
		case pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_UNSPECIFIED:
			t.Errorf("Field %q has UNSPECIFIED status (should be explicit)", fm.GetFieldName())
		default:
			t.Errorf("Field %q has unknown status value: %d", fm.GetFieldName(), status)
		}

		// Verify CONDITIONAL and DYNAMIC have condition_description
		if status == pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL ||
			status == pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_DYNAMIC {
			if fm.GetConditionDescription() == "" {
				t.Errorf("Field %q has %v status but empty condition_description",
					fm.GetFieldName(), status)
			}
		}
	}
}

// TestDryRunUnimplementedPlugin validates backward compatibility behavior
// for legacy plugins that don't implement DryRun.
//
// Expected behavior:
//   - Plugin returns codes.Unimplemented error
//   - Host can detect this and fall back to other discovery methods
func TestDryRunUnimplementedPlugin(t *testing.T) {
	// Create a minimal plugin that doesn't implement DryRun
	plugin := &minimalPlugin{}
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	_, err := harness.Client().DryRun(ctx, req)

	// Should return Unimplemented error
	if err == nil {
		t.Fatal("Expected Unimplemented error from minimal plugin")
	}

	if status.Code(err) != codes.Unimplemented {
		t.Errorf("Expected codes.Unimplemented, got: %v", status.Code(err))
	}
}

// minimalPlugin implements only the minimum required interface without DryRun.
type minimalPlugin struct {
	pbc.UnimplementedCostSourceServiceServer
}

func (p *minimalPlugin) Name(context.Context, *pbc.NameRequest) (*pbc.NameResponse, error) {
	return &pbc.NameResponse{Name: "minimal-plugin"}, nil
}

// =============================================================================
// Conformance Test Registration
// =============================================================================

// DryRunBasicConformanceTest validates basic DryRun RPC functionality.
func DryRunBasicConformanceTest(harness *plugintesting.TestHarness) plugintesting.TestResult {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)
	duration := time.Since(start)

	if err != nil {
		return plugintesting.TestResult{
			Method:   "DryRun",
			Success:  false,
			Error:    err,
			Duration: duration,
			Details:  "DryRun RPC call failed",
		}
	}

	// Validate response
	if !resp.GetResourceTypeSupported() {
		return plugintesting.TestResult{
			Method:   "DryRun",
			Success:  false,
			Duration: duration,
			Details:  "Expected resource_type_supported=true for aws/ec2",
		}
	}

	if len(resp.GetFieldMappings()) == 0 {
		return plugintesting.TestResult{
			Method:   "DryRun",
			Success:  false,
			Duration: duration,
			Details:  "Expected non-empty field_mappings",
		}
	}

	// Verify response time < 100ms
	if duration > 100*time.Millisecond {
		return plugintesting.TestResult{
			Method:   "DryRun",
			Success:  false,
			Duration: duration,
			Details:  "Response time exceeded 100ms requirement",
		}
	}

	return plugintesting.TestResult{
		Method:   "DryRun",
		Success:  true,
		Duration: duration,
		Details:  "DryRun basic functionality validated",
	}
}

// =============================================================================
// T029: Configuration Validation - Valid Config
// =============================================================================

// TestDryRunConfigurationValidation validates that DryRun properly reports
// valid configuration.
//
// Expected behavior:
//   - configuration_valid = true when plugin is properly configured
//   - configuration_errors is empty or nil for valid config
func TestDryRunConfigurationValidation(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.DryRunConfigValid = true
	plugin.DryRunConfigErrors = nil

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)
	if err != nil {
		t.Fatalf("DryRun RPC failed: %v", err)
	}

	if !resp.GetConfigurationValid() {
		t.Error("Expected configuration_valid=true for valid config")
	}

	if len(resp.GetConfigurationErrors()) > 0 {
		t.Errorf("Expected no configuration_errors, got: %v", resp.GetConfigurationErrors())
	}
}

// =============================================================================
// T030: Configuration Validation - Invalid Config with Error Reporting
// =============================================================================

// TestDryRunConfigurationErrors validates that DryRun properly reports
// configuration errors with descriptive messages.
//
// Expected behavior:
//   - configuration_valid = false when plugin is misconfigured
//   - configuration_errors contains descriptive error messages
//   - Multiple errors can be reported at once
func TestDryRunConfigurationErrors(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.DryRunConfigValid = false
	plugin.DryRunConfigErrors = []string{
		"Missing required API key",
		"Invalid endpoint URL: must start with https://",
	}

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)
	if err != nil {
		t.Fatalf("DryRun RPC should not error for config issues: %v", err)
	}

	// Verify configuration is reported as invalid
	if resp.GetConfigurationValid() {
		t.Error("Expected configuration_valid=false for invalid config")
	}

	// Verify error messages are present
	configErrors := resp.GetConfigurationErrors()
	if len(configErrors) == 0 {
		t.Error("Expected configuration_errors to contain error messages")
	}

	// Verify expected error count
	if len(configErrors) != 2 {
		t.Errorf("Expected 2 configuration errors, got %d", len(configErrors))
	}

	// Verify error messages are descriptive (not empty)
	for i, errMsg := range configErrors {
		if errMsg == "" {
			t.Errorf("Configuration error %d is empty", i)
		}
	}
}

// DryRunCapabilityTest validates that plugins advertise DryRun capability.
func DryRunCapabilityTest(harness *plugintesting.TestHarness) plugintesting.TestResult {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp, err := harness.Client().Supports(ctx, req)
	duration := time.Since(start)

	if err != nil {
		return plugintesting.TestResult{
			Method:   "Supports",
			Success:  false,
			Error:    err,
			Duration: duration,
			Details:  "Supports RPC call failed",
		}
	}

	// Check for dry_run capability
	if caps := resp.GetCapabilities(); caps != nil {
		if supported, ok := caps["dry_run"]; ok && supported {
			return plugintesting.TestResult{
				Method:   "Supports",
				Success:  true,
				Duration: duration,
				Details:  "Plugin advertises dry_run capability",
			}
		}
	}

	return plugintesting.TestResult{
		Method:   "Supports",
		Success:  false,
		Duration: duration,
		Details:  "Plugin does not advertise dry_run capability in Supports response",
	}
}

// =============================================================================
// T035: CONDITIONAL Field Status with Description
// =============================================================================

// TestDryRunConditionalFieldStatus validates that CONDITIONAL fields have
// proper condition_description explaining when they are populated.
//
// Expected behavior:
//   - Fields with CONDITIONAL status must have non-empty condition_description
//   - Description explains the conditions under which field is populated
func TestDryRunConditionalFieldStatus(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	// Configure plugin with a field that has CONDITIONAL status
	plugin.DryRunFieldMappings = []*pbc.FieldMapping{
		{
			FieldName:            "availability_zone",
			SupportStatus:        pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL,
			ConditionDescription: "Only populated for resources deployed in multi-AZ configurations",
		},
		{
			FieldName:            "service_category",
			SupportStatus:        pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_SUPPORTED,
			ConditionDescription: "", // Not required for SUPPORTED
		},
	}

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)
	if err != nil {
		t.Fatalf("DryRun RPC failed: %v", err)
	}

	// Find the CONDITIONAL field and verify description
	foundConditional := false
	for _, fm := range resp.GetFieldMappings() {
		if fm.GetSupportStatus() == pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_CONDITIONAL {
			foundConditional = true
			if fm.GetConditionDescription() == "" {
				t.Errorf("CONDITIONAL field %q missing condition_description", fm.GetFieldName())
			}
		}
	}

	if !foundConditional {
		t.Error("Expected at least one CONDITIONAL field in response")
	}
}

// =============================================================================
// T036: DYNAMIC Field Status with Description
// =============================================================================

// TestDryRunDynamicFieldStatus validates that DYNAMIC fields have proper
// condition_description explaining what runtime data is needed.
//
// Expected behavior:
//   - Fields with DYNAMIC status should have condition_description
//   - Description explains what runtime data determines the value
func TestDryRunDynamicFieldStatus(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	// Configure plugin with a field that has DYNAMIC status
	plugin.DryRunFieldMappings = []*pbc.FieldMapping{
		{
			FieldName:            "billed_cost",
			SupportStatus:        pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_DYNAMIC,
			ConditionDescription: "Requires billing API call at runtime; value varies by usage",
		},
		{
			FieldName:            "effective_cost",
			SupportStatus:        pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_DYNAMIC,
			ConditionDescription: "Computed from billed_cost minus applicable discounts",
		},
	}

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp, err := harness.Client().DryRun(ctx, req)
	if err != nil {
		t.Fatalf("DryRun RPC failed: %v", err)
	}

	// Verify all DYNAMIC fields have descriptions
	dynamicCount := 0
	for _, fm := range resp.GetFieldMappings() {
		if fm.GetSupportStatus() == pbc.FieldSupportStatus_FIELD_SUPPORT_STATUS_DYNAMIC {
			dynamicCount++
			if fm.GetConditionDescription() == "" {
				t.Errorf("DYNAMIC field %q should have condition_description", fm.GetFieldName())
			}
		}
	}

	if dynamicCount == 0 {
		t.Error("Expected at least one DYNAMIC field in response")
	}
}

// =============================================================================
// T037: Simulation Parameters Affect Field Status
// =============================================================================

// =============================================================================
// T042: GetActualCost with dry_run=true
// =============================================================================

// TestGetActualCostDryRun validates that GetActualCost with dry_run=true
// returns a DryRunResponse without making actual cost retrieval.
//
// Expected behavior:
//   - dry_run_result is populated with field mappings
//   - No actual cost data is returned (empty results)
//   - Response time < 100ms (no external API calls)
func TestGetActualCostDryRun(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	req := &pbc.GetActualCostRequest{
		ResourceId: "i-1234567890abcdef0", // Mock resource ID
		Start:      timestamppb.New(now.Add(-24 * time.Hour)),
		End:        timestamppb.New(now),
		DryRun:     true, // Enable dry-run mode
	}

	start := time.Now()
	resp, err := harness.Client().GetActualCost(ctx, req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetActualCost with dry_run=true failed: %v", err)
	}

	// Verify dry_run_result is populated
	dryRunResult := resp.GetDryRunResult()
	if dryRunResult == nil {
		t.Fatal("Expected dry_run_result to be populated when dry_run=true")
	}

	// Verify field mappings are present
	if len(dryRunResult.GetFieldMappings()) == 0 {
		t.Error("Expected non-empty field_mappings in dry_run_result")
	}

	// Verify resource type is supported
	if !dryRunResult.GetResourceTypeSupported() {
		t.Error("Expected resource_type_supported=true for aws/ec2")
	}

	// Verify response time < 100ms (no external calls)
	if duration > 100*time.Millisecond {
		t.Errorf("GetActualCost dry_run exceeded 100ms: %v", duration)
	}
}

// =============================================================================
// T043: GetProjectedCost with dry_run=true
// =============================================================================

// TestGetProjectedCostDryRun validates that GetProjectedCost with dry_run=true
// returns a DryRunResponse without making actual cost calculations.
//
// Expected behavior:
//   - dry_run_result is populated with field mappings
//   - Response time < 100ms (no external API calls)
func TestGetProjectedCostDryRun(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
		DryRun: true, // Enable dry-run mode
	}

	start := time.Now()
	resp, err := harness.Client().GetProjectedCost(ctx, req)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetProjectedCost with dry_run=true failed: %v", err)
	}

	// Verify dry_run_result is populated
	dryRunResult := resp.GetDryRunResult()
	if dryRunResult == nil {
		t.Fatal("Expected dry_run_result to be populated when dry_run=true")
	}

	// Verify field mappings are present
	if len(dryRunResult.GetFieldMappings()) == 0 {
		t.Error("Expected non-empty field_mappings in dry_run_result")
	}

	// Verify resource type is supported
	if !dryRunResult.GetResourceTypeSupported() {
		t.Error("Expected resource_type_supported=true for aws/ec2")
	}

	// Verify response time < 100ms (no external calls)
	if duration > 100*time.Millisecond {
		t.Errorf("GetProjectedCost dry_run exceeded 100ms: %v", duration)
	}
}

// =============================================================================
// T044: dry_run=false Default Behavior
// =============================================================================

// TestCostRPCsNormalBehavior validates that cost RPCs work normally when
// dry_run=false (or not specified).
//
// Expected behavior:
//   - Normal cost data is returned
//   - dry_run_result is nil
func TestCostRPCsNormalBehavior(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()

	// Test GetActualCost with dry_run=false (default)
	actualReq := &pbc.GetActualCostRequest{
		ResourceId: "i-1234567890abcdef0", // Mock resource ID
		Start:      timestamppb.New(now.Add(-24 * time.Hour)),
		End:        timestamppb.New(now),
		DryRun:     false, // Explicitly false
	}

	actualResp, err := harness.Client().GetActualCost(ctx, actualReq)
	if err != nil {
		t.Fatalf("GetActualCost failed: %v", err)
	}

	// Normal response should have cost results
	if len(actualResp.GetResults()) == 0 {
		t.Error("Expected cost results when dry_run=false")
	}

	// dry_run_result should be nil
	if actualResp.GetDryRunResult() != nil {
		t.Error("Expected dry_run_result=nil when dry_run=false")
	}

	// Test GetProjectedCost with dry_run unset (default=false)
	projectedReq := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
		// DryRun not set, should default to false
	}

	projectedResp, err := harness.Client().GetProjectedCost(ctx, projectedReq)
	if err != nil {
		t.Fatalf("GetProjectedCost failed: %v", err)
	}

	// Normal response should have cost data
	if projectedResp.GetCostPerMonth() == 0 {
		t.Error("Expected non-zero cost_per_month when dry_run=false")
	}

	// dry_run_result should be nil
	if projectedResp.GetDryRunResult() != nil {
		t.Error("Expected dry_run_result=nil when dry_run=false/unset")
	}
}

// TestDryRunSimulationParameters validates that simulation_parameters can
// affect the returned field mappings.
//
// Expected behavior:
//   - Plugins can use simulation_parameters to adjust field status
//   - Different parameters may result in different field statuses
func TestDryRunSimulationParameters(t *testing.T) {
	// This test validates the simulation_parameters field is passed through.
	// In a real plugin, parameters like "region=multi-az" could change
	// availability_zone from CONDITIONAL to SUPPORTED.

	plugin := plugintesting.NewMockPlugin()

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First request without simulation parameters
	req1 := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
	}

	resp1, err := harness.Client().DryRun(ctx, req1)
	if err != nil {
		t.Fatalf("DryRun RPC failed: %v", err)
	}

	// Second request with simulation parameters
	req2 := &pbc.DryRunRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
		},
		SimulationParameters: map[string]string{
			"deployment_mode": "multi-az",
			"pricing_tier":    "reserved",
		},
	}

	resp2, err := harness.Client().DryRun(ctx, req2)
	if err != nil {
		t.Fatalf("DryRun RPC with simulation_parameters failed: %v", err)
	}

	// Both responses should be valid (simulation parameters don't break basic functionality)
	if !resp1.GetResourceTypeSupported() || !resp2.GetResourceTypeSupported() {
		t.Error("Both requests should return resource_type_supported=true")
	}

	// Verify simulation parameters are accepted (don't cause error)
	if !resp2.GetConfigurationValid() {
		t.Errorf("Simulation parameters should not cause config errors: %v",
			resp2.GetConfigurationErrors())
	}
}
