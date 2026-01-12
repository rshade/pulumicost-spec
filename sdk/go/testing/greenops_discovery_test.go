package testing_test

import (
	"context"
	"testing"
	"time"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

// countMatchingMetrics returns the number of expected metrics found in the supported list.
func countMatchingMetrics(expected, supported []pbc.MetricKind) int {
	count := 0
	for _, exp := range expected {
		for _, got := range supported {
			if exp == got {
				count++
				break
			}
		}
	}
	return count
}

// TestGreenOpsDiscovery validates that plugins can advertise supported GreenOps metrics.
func TestGreenOpsDiscovery(t *testing.T) {
	// Setup mock plugin with GreenOps metrics
	mock := plugintesting.NewMockPlugin()
	expectedMetrics := []pbc.MetricKind{
		pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT,
		pbc.MetricKind_METRIC_KIND_ENERGY_CONSUMPTION,
		pbc.MetricKind_METRIC_KIND_WATER_USAGE,
	}
	mock.SupportedMetrics = expectedMetrics

	harness := plugintesting.NewTestHarness(mock)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	t.Run("Plugin advertises GreenOps metrics", func(t *testing.T) {
		start := time.Now()
		resp, err := harness.Client().Supports(ctx, &pbc.SupportsRequest{
			Resource: resource,
		})
		duration := time.Since(start)

		if err != nil {
			t.Fatalf("Supports RPC failed: %v", err)
		}

		if !resp.GetSupported() {
			t.Fatal("Expected resource to be supported")
		}

		supported := resp.GetSupportedMetrics()
		if len(supported) != len(expectedMetrics) {
			t.Errorf("Expected %d metrics, got %d", len(expectedMetrics), len(supported))
		}

		foundCount := countMatchingMetrics(expectedMetrics, supported)
		if foundCount != len(expectedMetrics) {
			t.Errorf("Metric mismatch: found %d/%d expected metrics", foundCount, len(expectedMetrics))
		}

		t.Logf("GreenOps discovery completed in %v", duration)
	})

	t.Run("Plugin returns empty metrics when none supported", func(t *testing.T) {
		mock.SupportedMetrics = nil
		resp, err := harness.Client().Supports(ctx, &pbc.SupportsRequest{
			Resource: resource,
		})
		if err != nil {
			t.Fatalf("Supports RPC failed: %v", err)
		}

		if len(resp.GetSupportedMetrics()) != 0 {
			t.Errorf("Expected 0 metrics, got %d", len(resp.GetSupportedMetrics()))
		}
	})
}
