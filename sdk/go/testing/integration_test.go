package testing_test

import (
	"context"
	"testing"
	"time"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
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
