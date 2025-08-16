package testing

import (
	"context"
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// BenchmarkName benchmarks the Name RPC method
func BenchmarkName(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{}) // Dummy testing.T for harness
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Name(ctx, &pbc.NameRequest{})
		if err != nil {
			b.Fatalf("Name() failed: %v", err)
		}
	}
}

// BenchmarkSupports benchmarks the Supports RPC method
func BenchmarkSupports(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		if err != nil {
			b.Fatalf("Supports() failed: %v", err)
		}
	}
}

// BenchmarkGetActualCost benchmarks the GetActualCost RPC method
func BenchmarkGetActualCost(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	start, end := CreateTimeRange(24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      start,
			End:        end,
		})
		if err != nil {
			b.Fatalf("GetActualCost() failed: %v", err)
		}
	}
}

// BenchmarkGetProjectedCost benchmarks the GetProjectedCost RPC method
func BenchmarkGetProjectedCost(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err != nil {
			b.Fatalf("GetProjectedCost() failed: %v", err)
		}
	}
}

// BenchmarkGetPricingSpec benchmarks the GetPricingSpec RPC method
func BenchmarkGetPricingSpec(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: resource,
		})
		if err != nil {
			b.Fatalf("GetPricingSpec() failed: %v", err)
		}
	}
}

// BenchmarkAllMethods benchmarks all RPC methods in sequence
func BenchmarkAllMethods(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
	start, end := CreateTimeRange(24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Name
		_, err := client.Name(ctx, &pbc.NameRequest{})
		if err != nil {
			b.Fatalf("Name() failed: %v", err)
		}

		// Supports
		_, err = client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		if err != nil {
			b.Fatalf("Supports() failed: %v", err)
		}

		// GetActualCost
		_, err = client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      start,
			End:        end,
		})
		if err != nil {
			b.Fatalf("GetActualCost() failed: %v", err)
		}

		// GetProjectedCost
		_, err = client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err != nil {
			b.Fatalf("GetProjectedCost() failed: %v", err)
		}

		// GetPricingSpec
		_, err = client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: resource,
		})
		if err != nil {
			b.Fatalf("GetPricingSpec() failed: %v", err)
		}
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent RPC requests
func BenchmarkConcurrentRequests(b *testing.B) {
	plugin := NewMockPlugin()
	harness := NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Name(ctx, &pbc.NameRequest{})
			if err != nil {
				b.Fatalf("Name() failed: %v", err)
			}
		}
	})
}

// BenchmarkActualCostDataSizes benchmarks GetActualCost with different data sizes
func BenchmarkActualCostDataSizes(b *testing.B) {
	testCases := []struct {
		name      string
		hours     int
		dataPoints int
	}{
		{"1Hour", 1, 1},
		{"24Hours", 24, 24},
		{"7Days", 168, 168},
		{"30Days", 720, 720},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			plugin := NewMockPlugin()
			plugin.ActualCostDataPoints = tc.dataPoints
			
			harness := NewTestHarness(plugin)
			harness.Start(&testing.T{})
			defer harness.Stop()

			client := harness.Client()
			ctx := context.Background()
			start, end := CreateTimeRange(tc.hours)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
					ResourceId: "test-resource",
					Start:      start,
					End:        end,
				})
				if err != nil {
					b.Fatalf("GetActualCost() failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkDifferentProviders benchmarks requests across different providers
func BenchmarkDifferentProviders(b *testing.B) {
	providers := []struct {
		name         string
		provider     string
		resourceType string
		sku          string
	}{
		{"AWS", "aws", "ec2", "t3.micro"},
		{"Azure", "azure", "vm", "Standard_B1s"},
		{"GCP", "gcp", "compute_engine", "n1-standard-1"},
		{"Kubernetes", "kubernetes", "namespace", ""},
	}

	for _, p := range providers {
		b.Run(p.name, func(b *testing.B) {
			plugin := NewMockPlugin()
			harness := NewTestHarness(plugin)
			harness.Start(&testing.T{})
			defer harness.Stop()

			client := harness.Client()
			ctx := context.Background()
			resource := CreateResourceDescriptor(p.provider, p.resourceType, p.sku, "us-east-1")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
					Resource: resource,
				})
				if err != nil {
					b.Fatalf("GetProjectedCost() failed: %v", err)
				}
			}
		})
	}
}

// PerformanceTestSuite provides standardized performance testing
type PerformanceTestSuite struct {
	impl pbc.CostSourceServiceServer
}

// NewPerformanceTestSuite creates a new performance test suite
func NewPerformanceTestSuite(impl pbc.CostSourceServiceServer) *PerformanceTestSuite {
	return &PerformanceTestSuite{impl: impl}
}

// RunPerformanceTests runs a standardized set of performance tests
func (pts *PerformanceTestSuite) RunPerformanceTests(t *testing.T) {
	harness := NewTestHarness(pts.impl)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	// Test Name performance
	nameMetrics, err := MeasurePerformance("Name", 100, func() error {
		_, err := client.Name(ctx, &pbc.NameRequest{})
		return err
	})
	if err != nil {
		t.Fatalf("Name performance test failed: %v", err)
	}
	t.Logf("Name Performance: avg=%v, min=%v, max=%v",
		nameMetrics.AvgDuration, nameMetrics.MinDuration, nameMetrics.MaxDuration)

	// Test Supports performance
	resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
	supportsMetrics, err := MeasurePerformance("Supports", 100, func() error {
		_, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		return err
	})
	if err != nil {
		t.Fatalf("Supports performance test failed: %v", err)
	}
	t.Logf("Supports Performance: avg=%v, min=%v, max=%v",
		supportsMetrics.AvgDuration, supportsMetrics.MinDuration, supportsMetrics.MaxDuration)

	// Test GetProjectedCost performance
	projectedMetrics, err := MeasurePerformance("GetProjectedCost", 100, func() error {
		_, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		return err
	})
	if err != nil {
		t.Fatalf("GetProjectedCost performance test failed: %v", err)
	}
	t.Logf("GetProjectedCost Performance: avg=%v, min=%v, max=%v",
		projectedMetrics.AvgDuration, projectedMetrics.MinDuration, projectedMetrics.MaxDuration)

	// Test GetActualCost performance
	start, end := CreateTimeRange(24)
	actualMetrics, err := MeasurePerformance("GetActualCost", 50, func() error {
		_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      start,
			End:        end,
		})
		return err
	})
	if err != nil {
		t.Fatalf("GetActualCost performance test failed: %v", err)
	}
	t.Logf("GetActualCost Performance: avg=%v, min=%v, max=%v",
		actualMetrics.AvgDuration, actualMetrics.MinDuration, actualMetrics.MaxDuration)

	// Test GetPricingSpec performance
	specMetrics, err := MeasurePerformance("GetPricingSpec", 100, func() error {
		_, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: resource,
		})
		return err
	})
	if err != nil {
		t.Fatalf("GetPricingSpec performance test failed: %v", err)
	}
	t.Logf("GetPricingSpec Performance: avg=%v, min=%v, max=%v",
		specMetrics.AvgDuration, specMetrics.MinDuration, specMetrics.MaxDuration)
}