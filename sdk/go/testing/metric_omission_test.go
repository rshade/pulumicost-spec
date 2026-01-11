package testing_test

import (
	"context"
	"testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

func TestMetricOmission(t *testing.T) {
	mock := plugintesting.NewMockPlugin()
	mock.SupportedMetrics = []pbc.MetricKind{
		pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT,
		pbc.MetricKind_METRIC_KIND_ENERGY_CONSUMPTION,
	}
	// Omit Energy Consumption for this test
	mock.OmitMetrics = []pbc.MetricKind{
		pbc.MetricKind_METRIC_KIND_ENERGY_CONSUMPTION,
	}

	harness := plugintesting.NewTestHarness(mock)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	t.Run("Plugin omits unavailable metric", func(t *testing.T) {
		resp, err := harness.Client().GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err != nil {
			t.Fatalf("GetProjectedCost failed: %v", err)
		}

		metrics := resp.GetImpactMetrics()
		foundCarbon := false
		foundEnergy := false

		for _, m := range metrics {
			if m.GetKind() == pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT {
				foundCarbon = true
			}
			if m.GetKind() == pbc.MetricKind_METRIC_KIND_ENERGY_CONSUMPTION {
				foundEnergy = true
			}
		}

		if !foundCarbon {
			t.Error("Expected Carbon Footprint to be present")
		}
		if foundEnergy {
			t.Error("Expected Energy Consumption to be omitted")
		}
	})
}
