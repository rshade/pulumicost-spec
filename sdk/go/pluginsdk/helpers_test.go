package pluginsdk_test

import (
	"context"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func TestResourceMatcher(t *testing.T) {
	matcher := pluginsdk.NewResourceMatcher()

	// Add supported providers and resource types
	matcher.AddProvider("aws")
	matcher.AddResourceType("aws:ec2:Instance")

	testCases := []struct {
		name     string
		resource *pbc.ResourceDescriptor
		expected bool
	}{
		{
			name: "supported provider and type",
			resource: &pbc.ResourceDescriptor{
				Provider:     "aws",
				ResourceType: "aws:ec2:Instance",
			},
			expected: true,
		},
		{
			name: "unsupported provider",
			resource: &pbc.ResourceDescriptor{
				Provider:     "azure",
				ResourceType: "azure:compute:VirtualMachine",
			},
			expected: false,
		},
		{
			name: "supported provider, unsupported type",
			resource: &pbc.ResourceDescriptor{
				Provider:     "aws",
				ResourceType: "aws:s3:Bucket",
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := matcher.Supports(tc.resource)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v for resource %v", tc.expected, result, tc.resource)
			}
		})
	}
}

func TestResourceMatcherNoFilters(t *testing.T) {
	// Empty matcher should support everything
	matcher := pluginsdk.NewResourceMatcher()

	resource := &pbc.ResourceDescriptor{
		Provider:     "any",
		ResourceType: "any:resource:Type",
	}

	if !matcher.Supports(resource) {
		t.Errorf("Empty matcher should support all resources")
	}
}

func TestCostCalculator(t *testing.T) {
	calc := pluginsdk.NewCostCalculator()

	// Test hourly to monthly conversion
	hourly := 0.10
	monthly := calc.HourlyToMonthly(hourly)
	expected := 73.0 // 0.10 * 730
	if monthly != expected {
		t.Errorf("Expected monthly cost %f, got %f", expected, monthly)
	}

	// Test monthly to hourly conversion
	monthlyInput := 146.0
	hourlyResult := calc.MonthlyToHourly(monthlyInput)
	expectedHourly := 0.2 // 146.0 / 730
	if hourlyResult != expectedHourly {
		t.Errorf("Expected hourly cost %f, got %f", expectedHourly, hourlyResult)
	}
}

func TestCostCalculatorResponses(t *testing.T) {
	calc := pluginsdk.NewCostCalculator()

	// Test projected cost response
	resp := calc.CreateProjectedCostResponse("USD", 0.05, "Test billing detail")

	if resp.GetCurrency() != "USD" {
		t.Errorf("Expected currency USD, got %s", resp.GetCurrency())
	}

	if resp.GetUnitPrice() != 0.05 {
		t.Errorf("Expected unit price 0.05, got %f", resp.GetUnitPrice())
	}

	if resp.GetCostPerMonth() != 36.5 { // 0.05 * 730
		t.Errorf("Expected cost per month 36.5, got %f", resp.GetCostPerMonth())
	}

	if resp.GetBillingDetail() != "Test billing detail" {
		t.Errorf("Expected billing detail 'Test billing detail', got %s", resp.GetBillingDetail())
	}

	// Test actual cost response
	results := []*pbc.ActualCostResult{
		{Cost: 10.0, Source: "test"},
	}
	actualResp := calc.CreateActualCostResponse(results)

	if len(actualResp.GetResults()) != 1 {
		t.Errorf("Expected 1 result, got %d", len(actualResp.GetResults()))
	}

	if actualResp.GetResults()[0].GetCost() != 10.0 {
		t.Errorf("Expected cost 10.0, got %f", actualResp.GetResults()[0].GetCost())
	}
}

func TestErrorFunctions(t *testing.T) {
	resource := &pbc.ResourceDescriptor{
		Provider:     "test",
		ResourceType: "test:resource:Type",
	}

	// Test NotSupportedError
	err := pluginsdk.NotSupportedError(resource)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedMsg := "resource type test:resource:Type from provider test is not supported"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}

	// Test NoDataError
	err = pluginsdk.NoDataError("test-resource-id")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedMsg = "no cost data available for resource test-resource-id"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestBasePlugin(t *testing.T) {
	plugin := pluginsdk.NewBasePlugin("test-plugin")

	// Test name
	if plugin.Name() != "test-plugin" {
		t.Errorf("Expected name 'test-plugin', got %s", plugin.Name())
	}

	// Test matcher access
	matcher := plugin.Matcher()
	if matcher == nil {
		t.Error("Expected matcher, got nil")
	}

	// Test calculator access
	calc := plugin.Calculator()
	if calc == nil {
		t.Error("Expected calculator, got nil")
	}

	// Test default implementations
	ctx := context.Background()

	// Test default GetProjectedCost (should return not supported error)
	req := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "test",
			ResourceType: "test:resource:Type",
		},
	}

	_, err := plugin.GetProjectedCost(ctx, req)
	if err == nil {
		t.Error("Expected not supported error, got nil")
	}

	// Test default GetActualCost (should return no data error)
	actualReq := &pbc.GetActualCostRequest{
		ResourceId: "test-resource-id",
	}

	_, err = plugin.GetActualCost(ctx, actualReq)
	if err == nil {
		t.Error("Expected no data error, got nil")
	}

	// Test default GetPricingSpec (should return not implemented error)
	pricingReq := &pbc.GetPricingSpecRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "test",
			ResourceType: "test:resource:Type",
		},
	}

	_, err = plugin.GetPricingSpec(ctx, pricingReq)
	if err == nil {
		t.Error("Expected not implemented error, got nil")
	}

	// Test default EstimateCost (should return not implemented error)
	estimateReq := &pbc.EstimateCostRequest{
		ResourceType: "test:resource:Type",
	}

	_, err = plugin.EstimateCost(ctx, estimateReq)
	if err == nil {
		t.Error("Expected not implemented error, got nil")
	}
}

func TestBasePluginNilRequests(t *testing.T) {
	plugin := pluginsdk.NewBasePlugin("test-plugin")
	ctx := context.Background()

	// Test nil GetProjectedCostRequest
	_, err := plugin.GetProjectedCost(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}

	// Test GetProjectedCost with nil resource
	_, err = plugin.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{})
	if err == nil {
		t.Error("Expected error for nil resource, got nil")
	}

	// Test nil GetActualCostRequest
	_, err = plugin.GetActualCost(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}

	// Test nil GetPricingSpecRequest
	_, err = plugin.GetPricingSpec(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}

	// Test nil EstimateCostRequest
	_, err = plugin.EstimateCost(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}
}

func TestResourceMatcherEmptyStrings(t *testing.T) {
	matcher := pluginsdk.NewResourceMatcher()

	// Empty strings should be ignored
	matcher.AddProvider("")
	matcher.AddResourceType("")

	// Should not have added empty strings - matcher should still support all
	resource := &pbc.ResourceDescriptor{
		Provider:     "any",
		ResourceType: "any:resource:Type",
	}

	if !matcher.Supports(resource) {
		t.Error("Empty string filters should be ignored, matcher should support all")
	}
}

func TestResourceMatcherNilResource(t *testing.T) {
	matcher := pluginsdk.NewResourceMatcher()
	matcher.AddProvider("aws")

	// Nil resource should return false
	if matcher.Supports(nil) {
		t.Error("Expected false for nil resource")
	}
}

func TestResourceMatcherNilReceiver(t *testing.T) {
	var matcher *pluginsdk.ResourceMatcher

	// Nil receiver should return false
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "aws:ec2:Instance",
	}

	if matcher.Supports(resource) {
		t.Error("Expected false for nil matcher")
	}
}

func TestHoursPerMonthExported(t *testing.T) {
	// Verify the constant is exported and has the correct value
	if pluginsdk.HoursPerMonth != 730.0 {
		t.Errorf("Expected HoursPerMonth to be 730.0, got %f", pluginsdk.HoursPerMonth)
	}
}

// TestFallbackHintDefaultValue tests default hint value (unspecified = 0) behavior.
func TestFallbackHintDefaultValue(t *testing.T) {
	// Verify that FallbackHint_FALLBACK_HINT_UNSPECIFIED is 0 (proto3 default)
	if pbc.FallbackHint_FALLBACK_HINT_UNSPECIFIED != 0 {
		t.Errorf("Expected FALLBACK_HINT_UNSPECIFIED to be 0, got %d",
			pbc.FallbackHint_FALLBACK_HINT_UNSPECIFIED)
	}

	// A new GetActualCostResponse should have unspecified (0) as default
	resp := &pbc.GetActualCostResponse{}
	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_UNSPECIFIED {
		t.Errorf("Expected default hint to be UNSPECIFIED, got %v", resp.GetFallbackHint())
	}
}

// TestCreateActualCostResponseDefaultHint tests that CreateActualCostResponse returns unspecified hint by default.
func TestCreateActualCostResponseDefaultHint(t *testing.T) {
	calc := pluginsdk.NewCostCalculator()

	results := []*pbc.ActualCostResult{
		{Cost: 10.0, Source: "test"},
	}
	resp := calc.CreateActualCostResponse(results)

	// Existing CreateActualCostResponse should work with default (unspecified) hint
	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_UNSPECIFIED {
		t.Errorf("Expected default hint to be UNSPECIFIED, got %v", resp.GetFallbackHint())
	}

	// Results should still be set correctly
	if len(resp.GetResults()) != 1 {
		t.Errorf("Expected 1 result, got %d", len(resp.GetResults()))
	}
}

// TestWithFallbackHintNone tests WithFallbackHint(NONE) option.
func TestWithFallbackHintNone(t *testing.T) {
	resp := pluginsdk.NewActualCostResponse(
		pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
	)

	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_NONE {
		t.Errorf("Expected FALLBACK_HINT_NONE, got %v", resp.GetFallbackHint())
	}
}

// TestNewActualCostResponseWithResultsAndNoneHint tests NewActualCostResponse with results and explicit NONE hint.
func TestNewActualCostResponseWithResultsAndNoneHint(t *testing.T) {
	results := []*pbc.ActualCostResult{
		{Cost: 25.50, Source: "aws-ce"},
		{Cost: 10.00, Source: "aws-ce"},
	}

	resp := pluginsdk.NewActualCostResponse(
		pluginsdk.WithResults(results),
		pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
	)

	// Check results
	if len(resp.GetResults()) != 2 {
		t.Errorf("Expected 2 results, got %d", len(resp.GetResults()))
	}
	if resp.GetResults()[0].GetCost() != 25.50 {
		t.Errorf("Expected first cost 25.50, got %f", resp.GetResults()[0].GetCost())
	}

	// Check hint
	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_NONE {
		t.Errorf("Expected FALLBACK_HINT_NONE, got %v", resp.GetFallbackHint())
	}
}

// TestWithFallbackHintRecommendedEmptyResults tests WithFallbackHint(RECOMMENDED) with empty results.
func TestWithFallbackHintRecommendedEmptyResults(t *testing.T) {
	resp := pluginsdk.NewActualCostResponse(
		pluginsdk.WithResults(nil),
		pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
	)

	// Should have no results
	if len(resp.GetResults()) != 0 {
		t.Errorf("Expected 0 results, got %d", len(resp.GetResults()))
	}

	// Should have RECOMMENDED hint
	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED {
		t.Errorf("Expected FALLBACK_HINT_RECOMMENDED, got %v", resp.GetFallbackHint())
	}
}

// TestNewActualCostResponseNilResultsRecommended tests NewActualCostResponse with nil results and RECOMMENDED hint.
func TestNewActualCostResponseNilResultsRecommended(t *testing.T) {
	resp := pluginsdk.NewActualCostResponse(
		pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
	)

	// Results should be nil/empty
	if resp.GetResults() != nil && len(resp.GetResults()) != 0 {
		t.Errorf("Expected nil or empty results, got %v", resp.GetResults())
	}

	// Hint should be RECOMMENDED
	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED {
		t.Errorf("Expected FALLBACK_HINT_RECOMMENDED, got %v", resp.GetFallbackHint())
	}
}

// TestWithFallbackHintRequired tests WithFallbackHint(REQUIRED).
func TestWithFallbackHintRequired(t *testing.T) {
	resp := pluginsdk.NewActualCostResponse(
		pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_REQUIRED),
	)

	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_REQUIRED {
		t.Errorf("Expected FALLBACK_HINT_REQUIRED, got %v", resp.GetFallbackHint())
	}
}

// TestNewActualCostResponseRequiredForUnsupportedType tests NewActualCostResponse with REQUIRED hint for unsupported type.
func TestNewActualCostResponseRequiredForUnsupportedType(t *testing.T) {
	// Simulate a plugin that cannot handle a specific resource type
	// It returns empty results with REQUIRED hint
	resp := pluginsdk.NewActualCostResponse(
		pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_REQUIRED),
	)

	// No results (plugin doesn't handle this type)
	if len(resp.GetResults()) != 0 {
		t.Errorf("Expected 0 results for unsupported type, got %d", len(resp.GetResults()))
	}

	// REQUIRED hint signals core must try fallback
	if resp.GetFallbackHint() != pbc.FallbackHint_FALLBACK_HINT_REQUIRED {
		t.Errorf("Expected FALLBACK_HINT_REQUIRED, got %v", resp.GetFallbackHint())
	}
}

// TestGetActualCostReturnsErrorForAPIFailures tests that GetActualCost should return error for API failures, not hint.
func TestGetActualCostReturnsErrorForAPIFailures(t *testing.T) {
	// This test documents the expected behavior: errors should be returned
	// as gRPC errors, not as success responses with hints.
	//
	// The BasePlugin.GetActualCost returns an error (NoDataError), not a hint.
	// This is the correct pattern - system failures use error path, not fallback path.

	plugin := pluginsdk.NewBasePlugin("test-plugin")
	ctx := context.Background()

	req := &pbc.GetActualCostRequest{
		ResourceId: "test-resource-id",
	}

	// Default implementation returns an error, not a response with hint
	resp, err := plugin.GetActualCost(ctx, req)

	// Should return error
	if err == nil {
		t.Error("Expected error for default GetActualCost, got nil")
	}

	// Response should be nil when error occurs
	if resp != nil {
		t.Errorf("Expected nil response when error occurs, got %v", resp)
	}

	// This demonstrates the pattern: errors (API failures, network issues) should
	// return gRPC errors, not success responses with hints.
}
