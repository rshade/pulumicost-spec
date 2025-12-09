package testing_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

// BenchmarkName benchmarks the Name RPC method.
func BenchmarkName(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{}) // Dummy testing.T for harness
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		_, err := client.Name(ctx, &pbc.NameRequest{})
		if err != nil {
			b.Fatalf("Name() failed: %v", err)
		}
	}
}

// BenchmarkSupports benchmarks the Supports RPC method.
func BenchmarkSupports(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	b.ResetTimer()
	for range b.N {
		_, err := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
		if err != nil {
			b.Fatalf("Supports() failed: %v", err)
		}
	}
}

// BenchmarkGetActualCost benchmarks the GetActualCost RPC method.
func BenchmarkGetActualCost(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	start, end := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)

	b.ResetTimer()
	for range b.N {
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

// BenchmarkGetProjectedCost benchmarks the GetProjectedCost RPC method.
func BenchmarkGetProjectedCost(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	b.ResetTimer()
	for range b.N {
		_, err := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
		if err != nil {
			b.Fatalf("GetProjectedCost() failed: %v", err)
		}
	}
}

// BenchmarkGetPricingSpec benchmarks the GetPricingSpec RPC method.
func BenchmarkGetPricingSpec(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	b.ResetTimer()
	for range b.N {
		_, err := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
			Resource: resource,
		})
		if err != nil {
			b.Fatalf("GetPricingSpec() failed: %v", err)
		}
	}
}

// BenchmarkEstimateCost benchmarks the EstimateCost RPC method.
func BenchmarkEstimateCost(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
			ResourceType: "aws:ec2/instance:Instance",
			Attributes:   nil,
		})
		if err != nil {
			b.Fatalf("EstimateCost() failed: %v", err)
		}
	}
}

// BenchmarkAllMethods benchmarks all RPC methods in sequence.
func BenchmarkAllMethods(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
	start, end := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)

	b.ResetTimer()
	for range b.N {
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

		// EstimateCost
		_, err = client.EstimateCost(ctx, &pbc.EstimateCostRequest{
			ResourceType: "aws:ec2/instance:Instance",
			Attributes:   nil,
		})
		if err != nil {
			b.Fatalf("EstimateCost() failed: %v", err)
		}
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent RPC requests.
func BenchmarkConcurrentRequests(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
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

// BenchmarkConcurrentEstimateCost benchmarks concurrent EstimateCost requests.
// This is a key benchmark for T044 - validates 50+ concurrent requests under load.
func BenchmarkConcurrentEstimateCost(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
				ResourceType: "aws:ec2/instance:Instance",
				Attributes:   nil,
			})
			if err != nil {
				b.Fatalf("EstimateCost() failed: %v", err)
			}
		}
	})
}

// BenchmarkConcurrentEstimateCost50 benchmarks exactly 50 concurrent EstimateCost requests.
// This benchmark validates Advanced conformance requirements per T044.
func BenchmarkConcurrentEstimateCost50(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		var wg sync.WaitGroup
		errChan := make(chan error, plugintesting.AdvancedParallelRequests)

		for range plugintesting.AdvancedParallelRequests {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
					ResourceType: "aws:ec2/instance:Instance",
					Attributes:   nil,
				})
				if err != nil {
					errChan <- err
				}
			}()
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			b.Fatalf("EstimateCost() failed: %v", err)
		}
	}
}

// BenchmarkConcurrentEstimateCostLatency measures per-request latency under concurrent load.
// Validates <500ms response time requirement under 50+ concurrent requests.
func BenchmarkConcurrentEstimateCostLatency(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		var wg sync.WaitGroup
		latencies := make(chan time.Duration, plugintesting.AdvancedParallelRequests)
		errChan := make(chan error, plugintesting.AdvancedParallelRequests)

		for range plugintesting.AdvancedParallelRequests {
			wg.Add(1)
			go func() {
				defer wg.Done()
				start := time.Now()
				_, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
					ResourceType: "aws:ec2/instance:Instance",
					Attributes:   nil,
				})
				latencies <- time.Since(start)
				if err != nil {
					errChan <- err
				}
			}()
		}

		wg.Wait()
		close(latencies)
		close(errChan)

		// Check for errors
		for err := range errChan {
			b.Fatalf("EstimateCost() failed: %v", err)
		}

		// Verify all latencies are under 500ms
		for latency := range latencies {
			if latency > 500*time.Millisecond {
				b.Fatalf("Latency %v exceeds 500ms threshold", latency)
			}
		}
	}
}

// BenchmarkActualCostDataSizes benchmarks GetActualCost with different data sizes.
func BenchmarkActualCostDataSizes(b *testing.B) {
	testCases := []struct {
		name       string
		hours      int
		dataPoints int
	}{
		{"1Hour", 1, 1},
		{"24Hours", plugintesting.HoursPerDay, plugintesting.HoursPerDay},
		{"7Days", 168, 168},
		{"30Days", plugintesting.HoursIn30Days, plugintesting.HoursIn30Days},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			plugin := plugintesting.NewMockPlugin()
			plugin.ActualCostDataPoints = tc.dataPoints

			harness := plugintesting.NewTestHarness(plugin)
			harness.Start(&testing.T{})
			defer harness.Stop()

			client := harness.Client()
			ctx := context.Background()
			start, end := plugintesting.CreateTimeRange(tc.hours)

			b.ResetTimer()
			for range b.N {
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

// BenchmarkDifferentProviders benchmarks requests across different providers.
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
			plugin := plugintesting.NewMockPlugin()
			harness := plugintesting.NewTestHarness(plugin)
			harness.Start(&testing.T{})
			defer harness.Stop()

			client := harness.Client()
			ctx := context.Background()
			resource := plugintesting.CreateResourceDescriptor(
				p.provider,
				p.resourceType,
				p.sku,
				"us-east-1",
			)

			b.ResetTimer()
			for range b.N {
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

// PerformanceTestSuite provides standardized performance testing.
type PerformanceTestSuite struct {
	impl pbc.CostSourceServiceServer
}

// NewPerformanceTestSuite creates a new performance test suite.
func NewPerformanceTestSuite(impl pbc.CostSourceServiceServer) *PerformanceTestSuite {
	return &PerformanceTestSuite{impl: impl}
}

// RunPerformanceTests runs a standardized set of performance tests.
func (pts *PerformanceTestSuite) RunPerformanceTests(t *testing.T) {
	harness := plugintesting.NewTestHarness(pts.impl)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	// Test Name performance
	nameMetrics, err := plugintesting.MeasurePerformance(
		"Name",
		plugintesting.NumPerformanceIterations,
		func() error {
			_, err := client.Name(ctx, &pbc.NameRequest{})
			return err
		},
	)
	if err != nil {
		t.Fatalf("Name performance test failed: %v", err)
	}
	t.Logf("Name Performance: avg=%v, min=%v, max=%v",
		nameMetrics.AvgDuration, nameMetrics.MinDuration, nameMetrics.MaxDuration)

	// Test Supports performance
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
	supportsMetrics, err := plugintesting.MeasurePerformance(
		"Supports",
		plugintesting.NumPerformanceIterations,
		func() error {
			_, callErr := client.Supports(ctx, &pbc.SupportsRequest{Resource: resource})
			return callErr
		},
	)
	if err != nil {
		t.Fatalf("Supports performance test failed: %v", err)
	}
	t.Logf("Supports Performance: avg=%v, min=%v, max=%v",
		supportsMetrics.AvgDuration, supportsMetrics.MinDuration, supportsMetrics.MaxDuration)

	// Test GetProjectedCost performance
	projectedMetrics, err := plugintesting.MeasurePerformance(
		"GetProjectedCost",
		plugintesting.NumPerformanceIterations,
		func() error {
			_, callErr := client.GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
				Resource: resource,
			})
			return callErr
		},
	)
	if err != nil {
		t.Fatalf("GetProjectedCost performance test failed: %v", err)
	}
	t.Logf("GetProjectedCost Performance: avg=%v, min=%v, max=%v",
		projectedMetrics.AvgDuration, projectedMetrics.MinDuration, projectedMetrics.MaxDuration)

	// Test GetActualCost performance
	start, end := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)
	actualMetrics, err := plugintesting.MeasurePerformance(
		"GetActualCost",
		plugintesting.ReducedIterations,
		func() error {
			_, callErr := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
				ResourceId: "test-resource",
				Start:      start,
				End:        end,
			})
			return callErr
		},
	)
	if err != nil {
		t.Fatalf("GetActualCost performance test failed: %v", err)
	}
	t.Logf("GetActualCost Performance: avg=%v, min=%v, max=%v",
		actualMetrics.AvgDuration, actualMetrics.MinDuration, actualMetrics.MaxDuration)

	// Test GetPricingSpec performance
	specMetrics, err := plugintesting.MeasurePerformance(
		"GetPricingSpec",
		plugintesting.NumPerformanceIterations,
		func() error {
			_, callErr := client.GetPricingSpec(ctx, &pbc.GetPricingSpecRequest{
				Resource: resource,
			})
			return callErr
		},
	)
	if err != nil {
		t.Fatalf("GetPricingSpec performance test failed: %v", err)
	}
	t.Logf("GetPricingSpec Performance: avg=%v, min=%v, max=%v",
		specMetrics.AvgDuration, specMetrics.MinDuration, specMetrics.MaxDuration)
}

// largeResultSetRecommendationCount is the number of recommendations for large result set testing.
const largeResultSetRecommendationCount = 10000

// BenchmarkGetRecommendations_LargeResultSet benchmarks GetRecommendations with 10,000 recommendations
// returned in a single response. Per SC-005, this should complete in <500ms.
func BenchmarkGetRecommendations_LargeResultSet(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	// Generate 10,000 recommendations
	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: plugintesting.GenerateSampleRecommendations(largeResultSetRecommendationCount),
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, err := client.GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
			// Request the full configured set in a single response for SC-005.
			PageSize: largeResultSetRecommendationCount,
		})
		if err != nil {
			b.Fatalf("GetRecommendations() failed: %v", err)
		}
	}
}

// BenchmarkGetRecommendations_LargeResultSetPagination benchmarks paginating through 10,000 recommendations.
func BenchmarkGetRecommendations_LargeResultSetPagination(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: plugintesting.GenerateSampleRecommendations(largeResultSetRecommendationCount),
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(&testing.T{})
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		var token string
		for {
			resp, err := client.GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
				PageSize:  100,
				PageToken: token,
			})
			if err != nil {
				b.Fatalf("GetRecommendations() failed: %v", err)
			}
			token = resp.GetNextPageToken()
			if token == "" {
				break
			}
		}
	}
}

// TestGetRecommendations_LargeResultSetLatency validates <500ms requirement for large result sets.
// This tests fetching all 10k recommendations in a single response per SC-005.
func TestGetRecommendations_LargeResultSetLatency(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetRecommendationsConfig(plugintesting.RecommendationsConfig{
		Recommendations: plugintesting.GenerateSampleRecommendations(largeResultSetRecommendationCount),
	})

	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	// Measure single request latency for the full result set
	start := time.Now()
	_, err := client.GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
		PageSize: largeResultSetRecommendationCount,
	})
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("GetRecommendations() failed: %v", err)
	}

	maxLatency := 500 * time.Millisecond
	if duration > maxLatency {
		t.Errorf("GetRecommendations latency %v exceeds %v requirement", duration, maxLatency)
	}
	t.Logf("GetRecommendations latency for 10k recommendations: %v", duration)
}

// =============================================================================
// FallbackHint Performance Benchmarks
// =============================================================================

// BenchmarkFallbackHintResponseCreation benchmarks creating responses with different FallbackHint values.
// Measures overhead of the functional options pattern for response construction.
func BenchmarkFallbackHintResponseCreation(b *testing.B) {
	b.Run("WithoutHint", func(b *testing.B) {
		results := []*pbc.ActualCostResult{
			{Cost: 10.0, Source: "aws-ce"},
		}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			resp := &pbc.GetActualCostResponse{Results: results}
			_ = resp.GetResults() // Use the result to avoid unused write warning
		}
	})

	b.Run("WithFallbackHintNone", func(b *testing.B) {
		results := []*pbc.ActualCostResult{
			{Cost: 10.0, Source: "aws-ce"},
		}
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_ = plugintesting.NewActualCostResponseWithHint(
				results,
				pbc.FallbackHint_FALLBACK_HINT_NONE,
			)
		}
	})

	b.Run("WithFallbackHintRecommended", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_ = plugintesting.NewActualCostResponseWithHint(
				nil,
				pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED,
			)
		}
	})

	b.Run("WithFallbackHintRequired", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			_ = plugintesting.NewActualCostResponseWithHint(
				nil,
				pbc.FallbackHint_FALLBACK_HINT_REQUIRED,
			)
		}
	})
}

// BenchmarkGetActualCostWithFallbackHint benchmarks GetActualCost RPC with different FallbackHint values.
// Validates that FallbackHint adds minimal overhead to RPC performance.
func BenchmarkGetActualCostWithFallbackHint(b *testing.B) {
	testCases := []struct {
		name string
		hint pbc.FallbackHint
	}{
		{"Unspecified", pbc.FallbackHint_FALLBACK_HINT_UNSPECIFIED},
		{"None", pbc.FallbackHint_FALLBACK_HINT_NONE},
		{"Recommended", pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED},
		{"Required", pbc.FallbackHint_FALLBACK_HINT_REQUIRED},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			plugin := plugintesting.NewMockPlugin()
			plugin.SetFallbackHint(tc.hint)

			harness := plugintesting.NewTestHarness(plugin)
			harness.Start(&testing.T{})
			defer harness.Stop()

			client := harness.Client()
			ctx := context.Background()
			start, end := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
					ResourceId: "test-resource",
					Start:      start,
					End:        end,
				})
				if err != nil {
					b.Fatalf("GetActualCost() failed: %v", err)
				}
				if resp.GetFallbackHint() != tc.hint {
					b.Fatalf("Expected hint %v, got %v", tc.hint, resp.GetFallbackHint())
				}
			}
		})
	}
}

// BenchmarkValidateActualCostResponse benchmarks response validation with various data sizes.
func BenchmarkValidateActualCostResponse(b *testing.B) {
	testCases := []struct {
		name        string
		resultCount int
	}{
		{"Empty", 0},
		{"1Result", 1},
		{"10Results", 10},
		{"100Results", 100},
		{"1000Results", 1000},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			results := make([]*pbc.ActualCostResult, tc.resultCount)
			for i := range tc.resultCount {
				results[i] = &pbc.ActualCostResult{
					Cost:   float64(i) * 1.5,
					Source: "aws-ce",
				}
			}
			resp := &pbc.GetActualCostResponse{
				Results:      results,
				FallbackHint: pbc.FallbackHint_FALLBACK_HINT_NONE,
			}

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				_ = pluginsdk.ValidateActualCostResponse(resp)
			}
		})
	}
}
