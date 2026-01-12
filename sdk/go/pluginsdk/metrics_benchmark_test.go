package pluginsdk_test

import (
	"context"
	"testing"
	"time"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// ============================================================================
// Performance Benchmarks for Metrics Interceptor
// ============================================================================

// BenchmarkMetricsInterceptor_Overhead measures the overhead of the metrics interceptor.
func BenchmarkMetricsInterceptor_Overhead(b *testing.B) {
	metrics := pluginsdk.NewPluginMetrics("benchmark-plugin")
	interceptor := pluginsdk.MetricsInterceptorWithRegistry(metrics)

	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/finfocus.v1.CostSource/GetProjectedCost",
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = interceptor(ctx, "request", info, handler)
	}
}

// BenchmarkMetricsInterceptor_NoMetrics provides baseline without metrics interceptor.
func BenchmarkMetricsInterceptor_NoMetrics(b *testing.B) {
	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return "response", nil
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = handler(ctx, "request")
	}
}

// BenchmarkMetricsInterceptor_WithError measures overhead when handler returns error.
func BenchmarkMetricsInterceptor_WithError(b *testing.B) {
	metrics := pluginsdk.NewPluginMetrics("benchmark-error-plugin")
	interceptor := pluginsdk.MetricsInterceptorWithRegistry(metrics)

	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, context.DeadlineExceeded
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/finfocus.v1.CostSource/GetActualCost",
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = interceptor(ctx, "request", info, handler)
	}
}

// BenchmarkNewPluginMetrics measures the cost of creating new metrics.
func BenchmarkNewPluginMetrics(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		_ = pluginsdk.NewPluginMetrics("benchmark-plugin")
	}
}

// ============================================================================
// Latency Accuracy Test
// ============================================================================

// TestMetrics_LatencyAccuracy verifies recorded duration is accurate within tolerance.
func TestMetrics_LatencyAccuracy(t *testing.T) {
	metrics := pluginsdk.NewPluginMetrics("latency-test")
	interceptor := pluginsdk.MetricsInterceptorWithRegistry(metrics)

	const expectedDuration = 100 * time.Millisecond

	handler := func(_ context.Context, _ interface{}) (interface{}, error) {
		time.Sleep(expectedDuration)
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/finfocus.v1.CostSource/GetProjectedCost",
	}

	ctx := context.Background()
	_, err := interceptor(ctx, "request", info, handler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get histogram values
	count, sum, err := getHistogramValues(
		metrics.RequestDuration,
		"finfocus.v1.CostSource/GetProjectedCost",
		"latency-test",
	)
	if err != nil {
		t.Fatalf("failed to get histogram values: %v", err)
	}

	assert.Equal(t, uint64(1), count)

	// Verify latency is 100ms ± 10ms (10% tolerance for timing jitter)
	expectedSeconds := expectedDuration.Seconds()
	tolerance := 0.010 // 10ms tolerance
	assert.InDelta(t, expectedSeconds, sum, tolerance,
		"recorded duration should be 100ms ± 10ms")
}
