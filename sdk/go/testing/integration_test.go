package testing_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

// TestBasicPluginFunctionality tests all basic RPC methods of a plugin.
func TestBasicPluginFunctionality(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	testNameRPC(ctx, t, client, plugin)
	testSupportsRPC(ctx, t, client)
	testGetActualCostRPC(ctx, t, client)
	testGetProjectedCostRPC(ctx, t, client)
	testGetPricingSpecRPC(ctx, t, client)
}

func testNameRPC(
	ctx context.Context,
	t *testing.T,
	client pbc.CostSourceServiceClient,
	plugin *plugintesting.MockPlugin,
) {
	t.Run("Name", func(t *testing.T) {
		resp, err := client.Name(ctx, &pbc.NameRequest{})
		if err != nil {
			t.Fatalf("Name() failed: %v", err)
		}

		if validationErr := plugintesting.ValidateNameResponse(resp); validationErr != nil {
			t.Errorf("Invalid name response: %v", validationErr)
		}

		if resp.GetName() != plugin.PluginName {
			t.Errorf("Expected name %s, got %s", plugin.PluginName, resp.GetName())
		}
	})
}

func testSupportsRPC(ctx context.Context, t *testing.T, client pbc.CostSourceServiceClient) {
	t.Run("Supports", func(t *testing.T) {
		// Test supported resource
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		resp, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		if err != nil {
			t.Fatalf("Supports() failed: %v", err)
		}

		if validationErr := plugintesting.ValidateSupportsResponse(resp); validationErr != nil {
			t.Errorf("Invalid supports response: %v", validationErr)
		}

		if !resp.GetSupported() {
			t.Errorf("Expected resource to be supported, got: %s", resp.GetReason())
		}
	})

	t.Run("SupportsUnsupportedProvider", func(t *testing.T) {
		// Test unsupported provider
		resource := plugintesting.CreateResourceDescriptor("unsupported", "some_resource", "", "")
		resp, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		if err != nil {
			t.Fatalf("Supports() failed: %v", err)
		}

		if resp.GetSupported() {
			t.Error("Expected unsupported provider to be rejected")
		}

		if resp.GetReason() == "" {
			t.Error("Expected reason for unsupported provider")
		}
	})
}

func testGetActualCostRPC(ctx context.Context, t *testing.T, client pbc.CostSourceServiceClient) {
	t.Run("GetActualCost", func(t *testing.T) {
		start, end := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)
		resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource-123",
			Start:      start,
			End:        end,
			Tags: map[string]string{
				"environment": "test",
			},
		})
		if err != nil {
			t.Fatalf("GetActualCost() failed: %v", err)
		}

		if validationErr := plugintesting.ValidateActualCostResponse(resp); validationErr != nil {
			t.Errorf("Invalid actual cost response: %v", validationErr)
		}

		if len(resp.GetResults()) == 0 {
			t.Error("Expected some cost results")
		}
	})
}

func testGetProjectedCostRPC(
	ctx context.Context,
	t *testing.T,
	client pbc.CostSourceServiceClient,
) {
	t.Run("GetProjectedCost", func(t *testing.T) {
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		resp, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err != nil {
			t.Fatalf("GetProjectedCost() failed: %v", err)
		}

		if validationErr := plugintesting.ValidateProjectedCostResponse(resp); validationErr != nil {
			t.Errorf("Invalid projected cost response: %v", validationErr)
		}

		if resp.GetUnitPrice() <= 0 {
			t.Errorf("Expected positive unit price, got: %f", resp.GetUnitPrice())
		}
	})
}

func testGetPricingSpecRPC(ctx context.Context, t *testing.T, client pbc.CostSourceServiceClient) {
	t.Run("GetPricingSpec", func(t *testing.T) {
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		resp, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: resource,
		})
		if err != nil {
			t.Fatalf("GetPricingSpec() failed: %v", err)
		}

		if validationErr := plugintesting.ValidatePricingSpecResponse(resp); validationErr != nil {
			t.Errorf("Invalid pricing spec response: %v", validationErr)
		}

		spec := resp.GetSpec()
		if spec.GetProvider() != resource.GetProvider() {
			t.Errorf("Expected provider %s, got %s", resource.GetProvider(), spec.GetProvider())
		}
		if spec.GetResourceType() != resource.GetResourceType() {
			t.Errorf(
				"Expected resource type %s, got %s",
				resource.GetResourceType(),
				spec.GetResourceType(),
			)
		}
	})
}

// TestErrorHandling tests various error conditions.
func TestErrorHandling(t *testing.T) {
	plugin := plugintesting.ConfigurableErrorMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	t.Run("NameError", func(t *testing.T) {
		plugin.ShouldErrorOnName = true
		_, err := client.Name(ctx, &pbc.NameRequest{})
		if err == nil {
			t.Error("Expected error from Name(), got nil")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Error("Expected gRPC status error")
		}
		if st.Code() != codes.Internal {
			t.Errorf("Expected Internal error, got %v", st.Code())
		}
	})

	t.Run("SupportsError", func(t *testing.T) {
		plugin.ShouldErrorOnSupports = true
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		_, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		if err == nil {
			t.Error("Expected error from Supports(), got nil")
		}
	})

	t.Run("ActualCostError", func(t *testing.T) {
		plugin.ShouldErrorOnActualCost = true
		start, end := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)
		_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      start,
			End:        end,
		})
		if err == nil {
			t.Error("Expected error from GetActualCost(), got nil")
		}
	})

	t.Run("ProjectedCostError", func(t *testing.T) {
		plugin.ShouldErrorOnProjectedCost = true
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		_, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err == nil {
			t.Error("Expected error from GetProjectedCost(), got nil")
		}
	})

	t.Run("PricingSpecError", func(t *testing.T) {
		plugin.ShouldErrorOnPricingSpec = true
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		_, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: resource,
		})
		if err == nil {
			t.Error("Expected error from GetPricingSpec(), got nil")
		}
	})
}

// TestInputValidation tests input validation for all methods.
func TestInputValidation(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	t.Run("SupportsNilResource", func(t *testing.T) {
		resp, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: nil})
		if err != nil {
			t.Fatalf("Supports() failed: %v", err)
		}

		if resp.GetSupported() {
			t.Error("Expected nil resource to be unsupported")
		}
	})

	t.Run("ActualCostMissingTimestamps", func(t *testing.T) {
		_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			// Missing Start and End timestamps
		})
		if err == nil {
			t.Error("Expected error for missing timestamps")
		}
	})

	t.Run("ActualCostInvalidTimeRange", func(t *testing.T) {
		end, start := plugintesting.CreateTimeRange(plugintesting.HoursPerDay) // Swapped start and end
		_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      start,
			End:        end,
		})
		if err == nil {
			t.Error("Expected error for invalid time range")
		}
	})

	t.Run("ProjectedCostNilResource", func(t *testing.T) {
		_, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: nil,
		})
		if err == nil {
			t.Error("Expected error for nil resource")
		}
	})

	t.Run("PricingSpecNilResource", func(t *testing.T) {
		_, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: nil,
		})
		if err == nil {
			t.Error("Expected error for nil resource")
		}
	})
}

// TestMultipleProviders tests plugin behavior with different providers.
func TestMultipleProviders(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	providers := []struct {
		name         string
		provider     string
		resourceType string
		sku          string
	}{
		{"AWS EC2", "aws", "ec2", "t3.micro"},
		{"Azure VM", "azure", "vm", "Standard_B1s"},
		{"GCP Compute", "gcp", "compute_engine", "n1-standard-1"},
		{"Kubernetes", "kubernetes", "namespace", ""},
	}

	for _, p := range providers {
		t.Run(p.name, func(t *testing.T) {
			testProviderSupport(ctx, t, client, p.provider, p.resourceType, p.sku)
		})
	}
}

func testProviderSupport(
	ctx context.Context,
	t *testing.T,
	client pbc.CostSourceServiceClient,
	provider, resourceType, sku string,
) {
	resource := plugintesting.CreateResourceDescriptor(provider, resourceType, sku, "us-east-1")

	// Test Supports
	supportsResp, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
	if err != nil {
		t.Fatalf("Supports() failed: %v", err)
	}

	if !supportsResp.GetSupported() {
		t.Errorf("Provider %s should be supported: %s", provider, supportsResp.GetReason())
	}

	// Test GetProjectedCost
	projectedResp, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
		Resource: resource,
	})
	if err != nil {
		t.Fatalf("GetProjectedCost() failed: %v", err)
	}

	if validationErr := plugintesting.ValidateProjectedCostResponse(projectedResp); validationErr != nil {
		t.Errorf("Invalid projected cost response for %s: %v", provider, validationErr)
	}

	// Test GetPricingSpec
	specResp, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
		Resource: resource,
	})
	if err != nil {
		t.Fatalf("GetPricingSpec() failed: %v", err)
	}

	spec := specResp.GetSpec()
	if spec.GetProvider() != provider {
		t.Errorf("Expected provider %s, got %s", provider, spec.GetProvider())
	}
	if spec.GetResourceType() != resourceType {
		t.Errorf("Expected resource type %s, got %s", resourceType, spec.GetResourceType())
	}
}

// TestConcurrentRequests tests plugin behavior under concurrent load.
func TestConcurrentRequests(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	const numConcurrentRequestsLocal = plugintesting.NumConcurrentRequests
	errors := make(chan error, numConcurrentRequestsLocal)

	// Run concurrent Name requests
	for range numConcurrentRequestsLocal {
		go func() {
			_, err := client.Name(ctx, &pbc.NameRequest{})
			errors <- err
		}()
	}

	// Check all requests completed successfully
	for i := range numConcurrentRequestsLocal {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent request %d failed: %v", i, err)
		}
	}
}

// TestResponseTimeouts tests plugin behavior with configured delays.
func TestResponseTimeouts(t *testing.T) {
	plugin := plugintesting.SlowMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	t.Run("NameWithDelay", func(t *testing.T) {
		start := time.Now()
		resp, err := client.Name(ctx, &pbc.NameRequest{})
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Name() failed: %v", err)
		}

		if duration < plugin.NameDelay {
			t.Errorf("Expected delay of at least %v, got %v", plugin.NameDelay, duration)
		}

		if resp.GetName() != plugin.PluginName {
			t.Errorf("Expected name %s, got %s", plugin.PluginName, resp.GetName())
		}
	})

	t.Run("SupportsWithDelay", func(t *testing.T) {
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		start := time.Now()
		_, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Supports() failed: %v", err)
		}

		if duration < plugin.SupportsDelay {
			t.Errorf("Expected delay of at least %v, got %v", plugin.SupportsDelay, duration)
		}
	})
}

// TestDataConsistency tests that plugin returns consistent data.
func TestDataConsistency(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	// Get projected cost multiple times - should be consistent
	var firstResponse *pbc.GetProjectedCostResponse
	for i := range 5 {
		resp, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err != nil {
			t.Fatalf("GetProjectedCost() failed on iteration %d: %v", i, err)
		}

		if i == 0 {
			firstResponse = resp
		} else {
			// Responses should be identical
			if resp.GetUnitPrice() != firstResponse.GetUnitPrice() {
				t.Errorf("Unit price inconsistent: iteration 0 = %f, iteration %d = %f",
					firstResponse.GetUnitPrice(), i, resp.GetUnitPrice())
			}
			if resp.GetCurrency() != firstResponse.GetCurrency() {
				t.Errorf("Currency inconsistent: iteration 0 = %s, iteration %d = %s",
					firstResponse.GetCurrency(), i, resp.GetCurrency())
			}
			if resp.GetCostPerMonth() != firstResponse.GetCostPerMonth() {
				t.Errorf("Cost per month inconsistent: iteration 0 = %f, iteration %d = %f",
					firstResponse.GetCostPerMonth(), i, resp.GetCostPerMonth())
			}
		}
	}
}

// =============================================================================
// STRUCTURED LOGGING EXAMPLE
// =============================================================================
//
// TestStructuredLoggingExample demonstrates structured logging patterns for the
// EstimateCost RPC, per NFR-001 of spec 006-estimate-cost.
//
// This example serves as the canonical reference for plugin developers to
// understand how to properly integrate zerolog structured logging with
// PulumiCost plugin operations.
//
// Key patterns demonstrated:
//   - Creating a configured logger with plugin metadata (FR-001)
//   - Logging requests with resource context (FR-002)
//   - Logging successful responses with cost details (FR-003)
//   - Logging errors with error codes and context (FR-004)
//   - Correlation ID (trace_id) propagation (FR-005)
//   - Using standard field name constants (FR-006)
//   - Operation timing with LogOperation helper (FR-009)
//
// Best Practices:
//   - ALWAYS include trace_id when available for distributed tracing
//   - NEVER log attribute values directly - log count only to prevent credential exposure
//   - Use standard field constants from pluginsdk for consistent naming
//   - Include operation name in every log entry for filterability
//   - Log at appropriate levels: Info for normal flow, Error for failures
//
//nolint:gocognit // Intentional: educational test with comprehensive inline documentation
func TestStructuredLoggingExample(t *testing.T) {
	// Setup plugin and harness for all subtests
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()

	// =========================================================================
	// User Story 1: Plugin Developer Learns Logging Patterns
	// =========================================================================

	// T008, T009, T010: RequestLogging subtest
	t.Run("RequestLogging", func(t *testing.T) {
		// Create a buffer to capture log output for verification
		var buf bytes.Buffer

		// FR-001: Create a configured logger with plugin name and version
		// Best Practice: Use NewPluginLogger to ensure consistent metadata fields
		logger := pluginsdk.NewPluginLogger(
			"example-plugin",  // plugin name - identifies the plugin in logs
			"v1.0.0",          // version - helps correlate logs with deployments
			zerolog.InfoLevel, // minimum log level
			&buf,              // output writer (use os.Stderr in production)
		)

		// Simulate a trace ID from incoming request context
		// Best Practice: Extract trace_id from gRPC metadata using TracingUnaryServerInterceptor
		traceID := "trace-abc123"
		ctx := pluginsdk.ContextWithTraceID(context.Background(), traceID)

		// Create sample request data
		resourceType := "aws:ec2/instance:Instance"
		attrs, _ := structpb.NewStruct(map[string]interface{}{
			"instanceType": "t3.micro",
			"region":       "us-east-1",
		})

		// FR-002: Log incoming request with resource context
		// Best Practice: Log attribute COUNT not values to prevent credential exposure
		logger.Info().
			Str(pluginsdk.FieldTraceID, pluginsdk.TraceIDFromContext(ctx)).
			Str(pluginsdk.FieldOperation, "EstimateCost").
			Str(pluginsdk.FieldResourceType, resourceType).
			Int("attribute_count", len(attrs.GetFields())).
			Msg("Processing cost estimation request")

		// Verify log output contains expected fields
		logOutput := buf.String()
		assertLogContains(t, logOutput, pluginsdk.FieldTraceID, "trace_id field missing")
		assertLogContains(t, logOutput, pluginsdk.FieldOperation, "operation field missing")
		assertLogContains(t, logOutput, pluginsdk.FieldResourceType, "resource_type field missing")
		assertLogContains(t, logOutput, "attribute_count", "attribute_count field missing")
		assertLogContains(t, logOutput, traceID, "trace_id value missing")
		assertLogContains(t, logOutput, resourceType, "resource_type value missing")

		// Verify we did NOT log sensitive attribute values
		assertLogNotContains(t, logOutput, "t3.micro", "sensitive attribute value should not be logged")
	})

	// T011, T012, T013: SuccessResponseLogging subtest
	t.Run("SuccessResponseLogging", func(t *testing.T) {
		var buf bytes.Buffer
		logger := pluginsdk.NewPluginLogger("example-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

		traceID := "trace-def456"
		ctx := pluginsdk.ContextWithTraceID(context.Background(), traceID)
		resourceType := "aws:ec2/instance:Instance"

		// FR-009: Use LogOperation helper for automatic timing measurement
		// Note: In production, use `defer done()` to ensure timing is logged even on panic.
		// Here we call done() explicitly to verify log output within the test.
		done := pluginsdk.LogOperation(logger, "EstimateCost")

		// Perform the actual EstimateCost RPC call
		attrs, _ := structpb.NewStruct(map[string]interface{}{
			"instanceType": "t3.micro",
			"region":       "us-east-1",
		})
		resp, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
			ResourceType: resourceType,
			Attributes:   attrs,
		})
		if err != nil {
			t.Fatalf("EstimateCost() failed: %v", err)
		}

		// FR-003: Log successful response with cost details
		// Best Practice: Include all relevant business data for operational visibility
		logger.Info().
			Str(pluginsdk.FieldTraceID, pluginsdk.TraceIDFromContext(ctx)).
			Str(pluginsdk.FieldOperation, "EstimateCost").
			Str(pluginsdk.FieldResourceType, resourceType).
			Float64(pluginsdk.FieldCostMonthly, resp.GetCostMonthly()).
			Str("currency", resp.GetCurrency()).
			Msg("Cost estimation completed")

		// Log operation timing (this will add duration_ms)
		done()

		// Verify log output
		logOutput := buf.String()
		assertLogContains(t, logOutput, pluginsdk.FieldCostMonthly, "cost_monthly field missing")
		assertLogContains(t, logOutput, "currency", "currency field missing")
		assertLogContains(t, logOutput, pluginsdk.FieldDurationMs, "duration_ms field missing")
		assertLogContains(t, logOutput, "Cost estimation completed", "completion message missing")

		// Verify cost value is present (even if zero - valid business case)
		entries := parseMultipleLogEntries(t, logOutput)
		foundCost := false
		for _, entry := range entries {
			if _, ok := entry[pluginsdk.FieldCostMonthly]; ok {
				foundCost = true
				break
			}
		}
		if !foundCost {
			t.Error("cost_monthly field not found in any log entry")
		}
	})

	// =========================================================================
	// User Story 2: Plugin Developer Implements Error Logging
	// =========================================================================

	// T015, T016, T017: ErrorLogging subtest
	t.Run("ErrorLogging", func(t *testing.T) {
		// Use configurable error mock for error injection
		errorPlugin := plugintesting.ConfigurableErrorMockPlugin()
		errorPlugin.ShouldErrorOnEstimateCost = true

		errorHarness := plugintesting.NewTestHarness(errorPlugin)
		errorHarness.Start(t)
		defer errorHarness.Stop()

		errorClient := errorHarness.Client()

		var buf bytes.Buffer
		logger := pluginsdk.NewPluginLogger("example-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

		traceID := "trace-error789"
		ctx := pluginsdk.ContextWithTraceID(context.Background(), traceID)
		resourceType := "invalid:resource/type:Type"

		// Attempt the call (will fail)
		attrs, _ := structpb.NewStruct(map[string]interface{}{})
		_, err := errorClient.EstimateCost(ctx, &pbc.EstimateCostRequest{
			ResourceType: resourceType,
			Attributes:   attrs,
		})

		// FR-004: Log errors with error code, message, and original request context
		// Best Practice: Include enough context to diagnose without re-running the request
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error().
				Err(err).
				Str(pluginsdk.FieldTraceID, pluginsdk.TraceIDFromContext(ctx)).
				Str(pluginsdk.FieldOperation, "EstimateCost").
				Str(pluginsdk.FieldResourceType, resourceType).
				Str(pluginsdk.FieldErrorCode, st.Code().String()).
				Msg("Cost estimation failed")
		}

		// Verify error log output
		logOutput := buf.String()
		assertLogContains(t, logOutput, "error", "error level/field missing")
		assertLogContains(t, logOutput, pluginsdk.FieldErrorCode, "error_code field missing")
		assertLogContains(t, logOutput, pluginsdk.FieldTraceID, "trace_id missing in error log")
		assertLogContains(t, logOutput, pluginsdk.FieldResourceType, "resource_type missing in error log")
		assertLogContains(t, logOutput, "Cost estimation failed", "error message missing")
	})

	// T018, T019, T020, T021: CorrelationIDPropagation subtest
	t.Run("CorrelationIDPropagation", func(t *testing.T) {
		var buf bytes.Buffer
		logger := pluginsdk.NewPluginLogger("example-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

		// FR-005: Demonstrate correlation ID propagation
		// Best Practice: Use ContextWithTraceID to store, TraceIDFromContext to retrieve
		traceID := "trace-correlation-xyz"
		ctx := pluginsdk.ContextWithTraceID(context.Background(), traceID)

		// Verify trace ID can be retrieved from context
		retrievedTraceID := pluginsdk.TraceIDFromContext(ctx)
		if retrievedTraceID != traceID {
			t.Errorf("Expected trace ID %s, got %s", traceID, retrievedTraceID)
		}

		// Log multiple operations with the same trace_id
		for i, op := range []string{"validate", "estimate", "respond"} {
			logger.Info().
				Str(pluginsdk.FieldTraceID, pluginsdk.TraceIDFromContext(ctx)).
				Str(pluginsdk.FieldOperation, op).
				Int("step", i+1).
				Msg("Processing step")
		}

		// Verify all log entries contain the same trace_id
		logOutput := buf.String()
		entries := parseMultipleLogEntries(t, logOutput)

		// Best Practice: Verify trace_id appears in ALL related log entries
		for i, entry := range entries {
			entryTraceID, ok := entry[pluginsdk.FieldTraceID].(string)
			if !ok {
				t.Errorf("Log entry %d missing trace_id", i)
				continue
			}
			if entryTraceID != traceID {
				t.Errorf("Log entry %d has wrong trace_id: expected %s, got %s", i, traceID, entryTraceID)
			}
		}

		// Test graceful degradation: empty context returns empty string
		emptyCtx := context.Background()
		emptyTraceID := pluginsdk.TraceIDFromContext(emptyCtx)
		if emptyTraceID != "" {
			t.Errorf("Expected empty trace_id from empty context, got %s", emptyTraceID)
		}
	})

	// =========================================================================
	// User Story 3: Operator Monitors EstimateCost Health
	// =========================================================================

	// T022, T023, T024, T025: LogStructureValidation subtest
	t.Run("LogStructureValidation", func(t *testing.T) {
		var buf bytes.Buffer
		logger := pluginsdk.NewPluginLogger("example-plugin", "v1.0.0", zerolog.InfoLevel, &buf)

		traceID := "trace-structure-test"
		ctx := pluginsdk.ContextWithTraceID(context.Background(), traceID)

		// Generate multiple log entries for different operations
		operations := []string{"EstimateCost", "GetProjectedCost", "GetActualCost"}
		resourceTypes := []string{
			"aws:ec2/instance:Instance",
			"azure:compute/virtualMachine:VirtualMachine",
			"gcp:compute/instance:Instance",
		}

		for i, op := range operations {
			logger.Info().
				Str(pluginsdk.FieldTraceID, pluginsdk.TraceIDFromContext(ctx)).
				Str(pluginsdk.FieldOperation, op).
				Str(pluginsdk.FieldResourceType, resourceTypes[i]).
				Float64(pluginsdk.FieldCostMonthly, float64(i+1)*10.50).
				Msg("Operation completed")
		}

		// Parse all log entries
		logOutput := buf.String()
		entries := parseMultipleLogEntries(t, logOutput)

		if len(entries) != len(operations) {
			t.Errorf("Expected %d log entries, got %d", len(operations), len(entries))
		}

		// Verify all entries are valid JSON and use standard field names
		// Best Practice: Consistent field names enable cross-plugin log aggregation
		for i, entry := range entries {
			// Check required fields exist
			if _, ok := entry[pluginsdk.FieldTraceID]; !ok {
				t.Errorf("Entry %d missing %s", i, pluginsdk.FieldTraceID)
			}
			if _, ok := entry[pluginsdk.FieldOperation]; !ok {
				t.Errorf("Entry %d missing %s", i, pluginsdk.FieldOperation)
			}

			// Verify filterability: operation field should allow filtering
			opValue, ok := entry[pluginsdk.FieldOperation].(string)
			if !ok {
				t.Errorf("Entry %d: operation field not a string", i)
				continue
			}
			if opValue != operations[i] {
				t.Errorf("Entry %d: expected operation %s, got %s", i, operations[i], opValue)
			}
		}

		// Demonstrate filterability by operation
		// Operators can query: jq 'select(.operation == "EstimateCost")' logs.json
		estimateCostCount := 0
		for _, entry := range entries {
			if op, ok := entry[pluginsdk.FieldOperation].(string); ok && op == "EstimateCost" {
				estimateCostCount++
			}
		}
		if estimateCostCount != 1 {
			t.Errorf("Expected 1 EstimateCost operation, found %d", estimateCostCount)
		}
	})
}

// =============================================================================
// Helper Functions for Log Verification
// =============================================================================

// parseMultipleLogEntries parses newline-delimited JSON log entries from a buffer.
// Each line is expected to be a valid JSON object.
func parseMultipleLogEntries(t *testing.T, logOutput string) []map[string]interface{} {
	t.Helper()
	var entries []map[string]interface{}

	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Errorf("Failed to parse log entry %d as JSON: %v\nLine: %s", i, err, line)
			continue
		}
		entries = append(entries, entry)
	}

	return entries
}

// assertLogContains verifies that the log output contains the expected substring.
func assertLogContains(t *testing.T, logOutput, expected, errMsg string) {
	t.Helper()
	if !strings.Contains(logOutput, expected) {
		t.Errorf("%s: expected '%s' in log output:\n%s", errMsg, expected, logOutput)
	}
}

// assertLogNotContains verifies that the log output does NOT contain the unexpected substring.
// Use this to verify sensitive data is not being logged.
func assertLogNotContains(t *testing.T, logOutput, unexpected, errMsg string) {
	t.Helper()
	if strings.Contains(logOutput, unexpected) {
		t.Errorf("%s: unexpected '%s' found in log output:\n%s", errMsg, unexpected, logOutput)
	}
}
