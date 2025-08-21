package testing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pricing"
)

// ObservabilityTestSuite provides comprehensive testing for plugin observability features.
type ObservabilityTestSuite struct {
	Plugin CostSourcePlugin
	T      *testing.T
}

// NewObservabilityTestSuite creates a new observability test suite.
func NewObservabilityTestSuite(plugin CostSourcePlugin, t *testing.T) *ObservabilityTestSuite {
	return &ObservabilityTestSuite{
		Plugin: plugin,
		T:      t,
	}
}

// RunBasicObservabilityTests runs basic observability conformance tests.
func (suite *ObservabilityTestSuite) RunBasicObservabilityTests() bool {
	suite.T.Helper()
	
	passed := true
	
	// Test health check availability
	if !suite.testHealthCheckBasic() {
		passed = false
	}
	
	// Test basic metrics collection
	if !suite.testBasicMetrics() {
		passed = false
	}
	
	// Test error handling with observability
	if !suite.testErrorObservability() {
		passed = false
	}
	
	return passed
}

// RunStandardObservabilityTests runs standard observability conformance tests.
func (suite *ObservabilityTestSuite) RunStandardObservabilityTests() bool {
	suite.T.Helper()
	
	if !suite.RunBasicObservabilityTests() {
		return false
	}
	
	passed := true
	
	// Test metrics endpoint functionality
	if !suite.testMetricsEndpoint() {
		passed = false
	}
	
	// Test SLI reporting
	if !suite.testSLIReporting() {
		passed = false
	}
	
	// Test tracing context
	if !suite.testTracingContext() {
		passed = false
	}
	
	// Test structured logging
	if !suite.testStructuredLogging() {
		passed = false
	}
	
	return passed
}

// RunAdvancedObservabilityTests runs advanced observability conformance tests.
func (suite *ObservabilityTestSuite) RunAdvancedObservabilityTests() bool {
	suite.T.Helper()
	
	if !suite.RunStandardObservabilityTests() {
		return false
	}
	
	passed := true
	
	// Test performance monitoring
	if !suite.testPerformanceMonitoring() {
		passed = false
	}
	
	// Test custom metrics
	if !suite.testCustomMetrics() {
		passed = false
	}
	
	// Test distributed tracing
	if !suite.testDistributedTracing() {
		passed = false
	}
	
	// Test observability under load
	if !suite.testObservabilityUnderLoad() {
		passed = false
	}
	
	return passed
}

func (suite *ObservabilityTestSuite) testHealthCheckBasic() bool {
	suite.T.Helper()
	
	// Check if plugin implements observability interface
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		suite.T.Error("Plugin does not implement ObservabilityPlugin interface")
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test basic health check
	health, err := obsPlugin.HealthCheck(ctx, "")
	if err != nil {
		suite.T.Errorf("Health check failed: %v", err)
		return false
	}
	
	if health.Status == "" {
		suite.T.Error("Health check returned empty status")
		return false
	}
	
	// Validate health status
	if err := pricing.ValidateHealthStatus(health.Status); err != nil {
		suite.T.Errorf("Invalid health status: %v", err)
		return false
	}
	
	suite.T.Log("✓ Basic health check passed")
	return true
}

func (suite *ObservabilityTestSuite) testBasicMetrics() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test basic metrics collection
	metrics, err := obsPlugin.GetMetrics(ctx, []string{})
	if err != nil {
		suite.T.Errorf("Failed to get metrics: %v", err)
		return false
	}
	
	if len(metrics.Metrics) == 0 {
		suite.T.Error("No metrics returned")
		return false
	}
	
	// Validate each metric
	requiredMetrics := pricing.GetObservabilityRequirements(pricing.ConformanceBasic).RequiredMetrics
	foundMetrics := make(map[string]bool)
	
	for _, metric := range metrics.Metrics {
		if err := pricing.ValidateMetricNameStrict(metric.Name); err != nil {
			suite.T.Errorf("Invalid metric name '%s': %v", metric.Name, err)
			return false
		}
		
		foundMetrics[metric.Name] = true
		
		// Validate metric samples
		for _, sample := range metric.Samples {
			if err := pricing.ValidateMetricValue(sample.Value, pricing.MetricType(metric.Type.String())); err != nil {
				suite.T.Errorf("Invalid metric value for '%s': %v", metric.Name, err)
				return false
			}
			
			// Validate labels
			result := pricing.ValidateMetricLabels(sample.Labels)
			if !result.Valid {
				suite.T.Errorf("Invalid metric labels for '%s': %v", metric.Name, result.Errors)
				return false
			}
		}
	}
	
	// Check required metrics are present
	for _, required := range requiredMetrics {
		if !foundMetrics[required] {
			suite.T.Errorf("Required metric '%s' not found", required)
			return false
		}
	}
	
	suite.T.Log("✓ Basic metrics collection passed")
	return true
}

func (suite *ObservabilityTestSuite) testErrorObservability() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	// Force an error condition and check observability
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// This should timeout and generate error metrics
	_, err := obsPlugin.GetMetrics(ctx, []string{})
	if err == nil {
		suite.T.Log("Expected timeout error, but operation succeeded")
	}
	
	// Check error metrics after the failed operation
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	
	metrics, err := obsPlugin.GetMetrics(ctx2, []string{pricing.StandardMetrics.ErrorsTotal})
	if err != nil {
		suite.T.Errorf("Failed to get error metrics: %v", err)
		return false
	}
	
	// Should have error metrics now
	errorMetricFound := false
	for _, metric := range metrics.Metrics {
		if metric.Name == pricing.StandardMetrics.ErrorsTotal {
			errorMetricFound = true
			break
		}
	}
	
	if !errorMetricFound {
		suite.T.Error("Error metrics not found after error condition")
		return false
	}
	
	suite.T.Log("✓ Error observability passed")
	return true
}

func (suite *ObservabilityTestSuite) testMetricsEndpoint() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test different metric formats
	formats := []string{"", "prometheus", "json"}
	
	for _, format := range formats {
		metrics, err := obsPlugin.GetMetrics(ctx, []string{})
		if err != nil {
			suite.T.Errorf("Failed to get metrics in format '%s': %v", format, err)
			return false
		}
		
		if metrics.Timestamp == nil {
			suite.T.Error("Metrics response missing timestamp")
			return false
		}
		
		// Check timestamp is recent
		if time.Since(metrics.Timestamp.AsTime()) > 10*time.Second {
			suite.T.Error("Metrics timestamp is too old")
			return false
		}
	}
	
	suite.T.Log("✓ Metrics endpoint testing passed")
	return true
}

func (suite *ObservabilityTestSuite) testSLIReporting() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Test SLI collection
	slis, err := obsPlugin.GetServiceLevelIndicators(ctx, nil, []string{})
	if err != nil {
		suite.T.Errorf("Failed to get SLIs: %v", err)
		return false
	}
	
	if len(slis.Slis) == 0 {
		suite.T.Error("No SLIs returned")
		return false
	}
	
	requiredSLIs := pricing.GetObservabilityRequirements(pricing.ConformanceStandard).RequiredSLIs
	foundSLIs := make(map[string]bool)
	
	for _, sli := range slis.Slis {
		if sli.Name == "" {
			suite.T.Error("SLI missing name")
			return false
		}
		
		foundSLIs[sli.Name] = true
		
		// Validate SLI value
		standardSLI := pricing.ServiceLevelIndicator{Name: sli.Name, Unit: sli.Unit}
		if err := pricing.ValidateSLIValue(standardSLI, sli.Value); err != nil {
			suite.T.Errorf("Invalid SLI value for '%s': %v", sli.Name, err)
			return false
		}
	}
	
	// Check required SLIs
	for _, required := range requiredSLIs {
		if !foundSLIs[required] {
			suite.T.Errorf("Required SLI '%s' not found", required)
			return false
		}
	}
	
	suite.T.Log("✓ SLI reporting passed")
	return true
}

func (suite *ObservabilityTestSuite) testTracingContext() bool {
	suite.T.Helper()
	
	// Test if responses include tracing metadata
	harness := NewTestHarness(suite.Plugin)
	harness.Start(suite.T)
	defer harness.Stop()
	
	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Make a request and check for tracing metadata
	resp, err := client.Name(ctx, &NameRequest{})
	if err != nil {
		suite.T.Errorf("Name request failed: %v", err)
		return false
	}
	
	// Check if response has telemetry metadata (if supported)
	if telemetryResp, ok := resp.(interface{ GetTelemetry() *TelemetryMetadata }); ok {
		metadata := telemetryResp.GetTelemetry()
		if metadata != nil {
			// Validate tracing IDs if present
			suite := pricing.ValidateObservabilityMetadata(
				metadata.TraceId,
				metadata.SpanId,
				metadata.RequestId,
				metadata.ProcessingTimeMs,
				metadata.QualityScore,
			)
			
			if !suite.IsValid() {
				suite.T.Errorf("Invalid telemetry metadata: %v", suite.Errors)
				return false
			}
		}
	}
	
	suite.T.Log("✓ Tracing context passed")
	return true
}

func (suite *ObservabilityTestSuite) testStructuredLogging() bool {
	suite.T.Helper()
	
	// This is a behavioral test - we can't directly test logging output
	// but we can validate that the plugin follows logging best practices
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	// Check if plugin provides logging configuration
	if logPlugin, ok := obsPlugin.(interface {
		GetLogLevel() string
		SetLogLevel(string) error
	}); ok {
		level := logPlugin.GetLogLevel()
		if err := pricing.ValidateLogLevel(level); err != nil {
			suite.T.Errorf("Invalid log level: %v", err)
			return false
		}
		
		// Test setting log level
		if err := logPlugin.SetLogLevel("INFO"); err != nil {
			suite.T.Errorf("Failed to set log level: %v", err)
			return false
		}
	}
	
	suite.T.Log("✓ Structured logging passed")
	return true
}

func (suite *ObservabilityTestSuite) testPerformanceMonitoring() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Test performance metrics collection
	metrics, err := obsPlugin.GetMetrics(ctx, []string{
		pricing.StandardMetrics.LatencyP99Seconds,
		pricing.StandardMetrics.MemoryUsageBytes,
		pricing.StandardMetrics.CpuUsagePercent,
	})
	if err != nil {
		suite.T.Errorf("Failed to get performance metrics: %v", err)
		return false
	}
	
	performanceMetricsFound := 0
	for _, metric := range metrics.Metrics {
		switch metric.Name {
		case pricing.StandardMetrics.LatencyP99Seconds,
			 pricing.StandardMetrics.MemoryUsageBytes,
			 pricing.StandardMetrics.CpuUsagePercent:
			performanceMetricsFound++
		}
	}
	
	if performanceMetricsFound < 2 {
		suite.T.Errorf("Expected at least 2 performance metrics, found %d", performanceMetricsFound)
		return false
	}
	
	suite.T.Log("✓ Performance monitoring passed")
	return true
}

func (suite *ObservabilityTestSuite) testCustomMetrics() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	// Test custom metric collection
	if customPlugin, ok := obsPlugin.(interface {
		RegisterCustomMetric(name, help string, metricType pricing.MetricType) error
		IncrementCustomMetric(name string, labels map[string]string) error
	}); ok {
		// Register a custom metric
		err := customPlugin.RegisterCustomMetric(
			"pulumicost_test_custom_counter",
			"Test custom counter metric",
			pricing.MetricTypeCounter,
		)
		if err != nil {
			suite.T.Errorf("Failed to register custom metric: %v", err)
			return false
		}
		
		// Increment the custom metric
		err = customPlugin.IncrementCustomMetric(
			"pulumicost_test_custom_counter",
			map[string]string{"test": "true"},
		)
		if err != nil {
			suite.T.Errorf("Failed to increment custom metric: %v", err)
			return false
		}
		
		// Verify the custom metric appears in metrics output
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		metrics, err := obsPlugin.GetMetrics(ctx, []string{"pulumicost_test_custom_counter"})
		if err != nil {
			suite.T.Errorf("Failed to get custom metrics: %v", err)
			return false
		}
		
		customMetricFound := false
		for _, metric := range metrics.Metrics {
			if metric.Name == "pulumicost_test_custom_counter" {
				customMetricFound = true
				break
			}
		}
		
		if !customMetricFound {
			suite.T.Error("Custom metric not found in metrics output")
			return false
		}
	}
	
	suite.T.Log("✓ Custom metrics passed")
	return true
}

func (suite *ObservabilityTestSuite) testDistributedTracing() bool {
	suite.T.Helper()
	
	// Test distributed tracing context propagation
	harness := NewTestHarness(suite.Plugin)
	harness.Start(suite.T)
	defer harness.Stop()
	
	client := harness.Client()
	
	// Create context with tracing metadata
	ctx := context.Background()
	traceID := "abcdef1234567890abcdef1234567890"
	spanID := "abcdef1234567890"
	
	// Add trace context to request (this would be done by tracing framework)
	if tracingClient, ok := client.(interface {
		NameWithTracing(ctx context.Context, req *NameRequest, traceID, spanID string) (*NameResponse, error)
	}); ok {
		resp, err := tracingClient.NameWithTracing(ctx, &NameRequest{}, traceID, spanID)
		if err != nil {
			suite.T.Errorf("Name request with tracing failed: %v", err)
			return false
		}
		
		// Check if response includes trace context
		if telemetryResp, ok := resp.(interface{ GetTelemetry() *TelemetryMetadata }); ok {
			metadata := telemetryResp.GetTelemetry()
			if metadata != nil && metadata.TraceId != traceID {
				suite.T.Errorf("Trace ID not propagated correctly: expected %s, got %s", traceID, metadata.TraceId)
				return false
			}
		}
	}
	
	suite.T.Log("✓ Distributed tracing passed")
	return true
}

func (suite *ObservabilityTestSuite) testObservabilityUnderLoad() bool {
	suite.T.Helper()
	
	obsPlugin, ok := suite.Plugin.(ObservabilityPlugin)
	if !ok {
		return false
	}
	
	// Test observability under concurrent load
	const numGoroutines = 10
	const requestsPerGoroutine = 10
	
	errChan := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			
			for j := 0; j < requestsPerGoroutine; j++ {
				// Collect metrics under load
				_, err := obsPlugin.GetMetrics(ctx, []string{})
				if err != nil {
					errChan <- fmt.Errorf("metrics collection failed under load: %v", err)
					return
				}
				
				// Check health under load
				_, err = obsPlugin.HealthCheck(ctx, "")
				if err != nil {
					errChan <- fmt.Errorf("health check failed under load: %v", err)
					return
				}
			}
			
			errChan <- nil
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		if err := <-errChan; err != nil {
			suite.T.Errorf("Observability under load test failed: %v", err)
			return false
		}
	}
	
	suite.T.Log("✓ Observability under load passed")
	return true
}

// ObservabilityPlugin interface that plugins can implement for observability testing.
type ObservabilityPlugin interface {
	HealthCheck(ctx context.Context, serviceName string) (*HealthCheckResponse, error)
	GetMetrics(ctx context.Context, metricNames []string) (*GetMetricsResponse, error)
	GetServiceLevelIndicators(ctx context.Context, timeRange *TimeRange, sliNames []string) (*GetSLIResponse, error)
}

// Helper types for testing (these would be generated from the protobuf)
type HealthCheckResponse struct {
	Status      string
	Message     string
	LastCheckTime time.Time
}

type GetMetricsResponse struct {
	Metrics   []Metric
	Timestamp *time.Time
	Format    string
}

type Metric struct {
	Name    string
	Help    string
	Type    MetricType
	Samples []MetricSample
}

type MetricType struct{}

func (mt MetricType) String() string {
	return "counter" // simplified for testing
}

type MetricSample struct {
	Labels    map[string]string
	Value     float64
	Timestamp time.Time
}

type GetSLIResponse struct {
	Slis            []ServiceLevelIndicator
	MeasurementTime time.Time
}

type ServiceLevelIndicator struct {
	Name         string
	Description  string
	Value        float64
	Unit         string
	TargetValue  float64
	Status       string
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type TelemetryMetadata struct {
	TraceId          string
	SpanId           string
	RequestId        string
	ProcessingTimeMs int64
	QualityScore     float64
}

type NameRequest struct{}

type NameResponse struct {
	Name string
}