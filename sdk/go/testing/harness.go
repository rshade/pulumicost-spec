package testing

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

const (
	bufSize = 1024 * 1024

	// Validation limits.
	maxPluginNameLength = 100 // Maximum allowed plugin name length
	currencyCodeLength  = 3   // ISO currency code length (e.g., USD, EUR)

	// HoursPerDay represents hours in a day for time range calculations.
	HoursPerDay              = 24
	MaxResponseTimeMs        = 100 // Maximum response time in milliseconds for performance tests
	NumConcurrentRequests    = 10  // Number of concurrent requests for concurrency tests
	NumPerformanceIterations = 100 // Number of iterations for performance measurements
	ReducedIterations        = 50  // Reduced iterations for expensive operations
	HoursIn30Days            = 720 // Hours in 30 days (24 * 30)
	MaxLargeQueryTimeSeconds = 10  // Maximum time for large dataset queries in seconds
	NumConsistencyChecks     = 3   // Number of consistency check iterations
	SuccessRateMultiplier    = 100 // Multiplier for success rate percentage calculation
)

// TestHarness provides a testing framework for CostSource plugin implementations.
type TestHarness struct {
	server   *grpc.Server
	listener *bufconn.Listener
	client   pbc.CostSourceServiceClient
	conn     *grpc.ClientConn
}

// NewTestHarness creates a new test harness for the given CostSource implementation.
func NewTestHarness(impl pbc.CostSourceServiceServer) *TestHarness {
	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	pbc.RegisterCostSourceServiceServer(server, impl)

	go func() {
		_ = server.Serve(listener)
	}()

	return &TestHarness{
		server:   server,
		listener: listener,
	}
}

// Start initializes the client connection to the test server.
func (h *TestHarness) Start(t *testing.T) {
	//nolint:staticcheck // grpc.NewClient doesn't work with bufconn
	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return h.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	h.conn = conn
	h.client = pbc.NewCostSourceServiceClient(conn)
}

// Stop cleans up the test harness.
func (h *TestHarness) Stop() {
	if h.conn != nil {
		_ = h.conn.Close()
	}
	if h.server != nil {
		h.server.Stop()
	}
}

// Client returns the gRPC client for making requests.
func (h *TestHarness) Client() pbc.CostSourceServiceClient {
	return h.client
}

// TestResult represents the result of a test operation.
type TestResult struct {
	Method   string
	Success  bool
	Error    error
	Duration time.Duration
	Details  string
}

// ConformanceTest represents a single conformance test case.
type ConformanceTest struct {
	Name        string
	Description string
	TestFunc    func(*TestHarness) TestResult
}

// PluginConformanceSuite contains all conformance tests for a plugin.
type PluginConformanceSuite struct {
	tests []ConformanceTest
}

// NewConformanceSuite creates a new conformance test suite.
func NewConformanceSuite() *PluginConformanceSuite {
	return &PluginConformanceSuite{
		tests: make([]ConformanceTest, 0),
	}
}

// AddTest adds a test to the conformance suite.
func (s *PluginConformanceSuite) AddTest(test ConformanceTest) {
	s.tests = append(s.tests, test)
}

// RunTests executes all conformance tests against the given plugin implementation.
func (s *PluginConformanceSuite) RunTests(t *testing.T, impl pbc.CostSourceServiceServer) []TestResult {
	harness := NewTestHarness(impl)
	harness.Start(t)
	defer harness.Stop()

	results := make([]TestResult, 0, len(s.tests))

	for _, test := range s.tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.TestFunc(harness)
			results = append(results, result)

			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
		})
	}

	return results
}

// Standard test helpers

// CreateResourceDescriptor creates a standard resource descriptor for testing.
func CreateResourceDescriptor(provider, resourceType, sku, region string) *pbc.ResourceDescriptor {
	return &pbc.ResourceDescriptor{
		Provider:     provider,
		ResourceType: resourceType,
		Sku:          sku,
		Region:       region,
		Tags: map[string]string{
			"environment": "test",
			"app":         "integration-test",
		},
	}
}

// CreateTimeRange creates a standard time range for testing.
func CreateTimeRange(hoursBack int) (*timestamppb.Timestamp, *timestamppb.Timestamp) {
	end := time.Now()
	start := end.Add(-time.Duration(hoursBack) * time.Hour)
	return timestamppb.New(start), timestamppb.New(end)
}

// ValidateNameResponse validates a Name RPC response.
func ValidateNameResponse(response *pbc.NameResponse) error {
	if response == nil {
		return errors.New("response is nil")
	}
	if response.GetName() == "" {
		return errors.New("plugin name is empty")
	}
	if len(response.GetName()) > maxPluginNameLength {
		return fmt.Errorf("plugin name too long: %d characters", len(response.GetName()))
	}
	return nil
}

// ValidateSupportsResponse validates a Supports RPC response.
func ValidateSupportsResponse(response *pbc.SupportsResponse) error {
	if response == nil {
		return errors.New("response is nil")
	}
	// If not supported, should have a reason
	if !response.GetSupported() && response.GetReason() == "" {
		return errors.New("unsupported resource should have a reason")
	}
	return nil
}

// ValidateActualCostResponse validates a GetActualCost RPC response.
func ValidateActualCostResponse(response *pbc.GetActualCostResponse) error {
	if response == nil {
		return errors.New("response is nil")
	}

	results := response.GetResults()
	if len(results) == 0 {
		// Empty results are valid (no cost data available)
		return nil
	}

	for i, result := range results {
		if err := ValidateActualCostResult(result); err != nil {
			return fmt.Errorf("result[%d]: %w", i, err)
		}
	}

	return nil
}

// ValidateActualCostResult validates a single ActualCostResult.
func ValidateActualCostResult(result *pbc.ActualCostResult) error {
	if result == nil {
		return errors.New("result is nil")
	}

	if result.GetTimestamp() == nil {
		return errors.New("timestamp is required")
	}

	if result.GetCost() < 0 {
		return fmt.Errorf("cost cannot be negative: %f", result.GetCost())
	}

	if result.GetUsageAmount() < 0 {
		return fmt.Errorf("usage amount cannot be negative: %f", result.GetUsageAmount())
	}

	if result.GetSource() == "" {
		return errors.New("source is required")
	}

	return nil
}

// ValidateProjectedCostResponse validates a GetProjectedCost RPC response.
func ValidateProjectedCostResponse(response *pbc.GetProjectedCostResponse) error {
	if response == nil {
		return errors.New("response is nil")
	}

	if response.GetUnitPrice() < 0 {
		return fmt.Errorf("unit price cannot be negative: %f", response.GetUnitPrice())
	}

	if response.GetCurrency() == "" {
		return errors.New("currency is required")
	}

	// Currency should be 3-character ISO code
	if len(response.GetCurrency()) != currencyCodeLength {
		return fmt.Errorf("currency should be 3-character ISO code, got: %s", response.GetCurrency())
	}

	if response.GetCostPerMonth() < 0 {
		return fmt.Errorf("cost per month cannot be negative: %f", response.GetCostPerMonth())
	}

	return nil
}

// ValidatePricingSpecResponse validates a GetPricingSpec RPC response.
func ValidatePricingSpecResponse(response *pbc.GetPricingSpecResponse) error {
	if response == nil {
		return errors.New("response is nil")
	}

	spec := response.GetSpec()
	if spec == nil {
		return errors.New("pricing spec is nil")
	}

	return ValidatePricingSpec(spec)
}

// ValidatePricingSpec validates a PricingSpec message.
func ValidatePricingSpec(spec *pbc.PricingSpec) error {
	if spec == nil {
		return errors.New("spec is nil")
	}

	if spec.GetProvider() == "" {
		return errors.New("provider is required")
	}

	if spec.GetResourceType() == "" {
		return errors.New("resource type is required")
	}

	if spec.GetBillingMode() == "" {
		return errors.New("billing mode is required")
	}

	if spec.GetRatePerUnit() < 0 {
		return fmt.Errorf("rate per unit cannot be negative: %f", spec.GetRatePerUnit())
	}

	if spec.GetCurrency() == "" {
		return errors.New("currency is required")
	}

	// Currency should be 3-character ISO code
	if len(spec.GetCurrency()) != currencyCodeLength {
		return fmt.Errorf("currency should be 3-character ISO code, got: %s", spec.GetCurrency())
	}

	return nil
}

// PerformanceMetrics contains performance measurement data.
type PerformanceMetrics struct {
	Method        string
	RequestCount  int
	TotalDuration time.Duration
	AvgDuration   time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
}

// MeasurePerformance measures the performance of a function.
func MeasurePerformance(name string, iterations int, testFunc func() error) (*PerformanceMetrics, error) {
	metrics := &PerformanceMetrics{
		Method:       name,
		RequestCount: iterations,
		MinDuration:  time.Hour, // Start with a large value
	}

	var totalErr error
	for range iterations {
		start := time.Now()
		err := testFunc()
		duration := time.Since(start)

		if err != nil && totalErr == nil {
			totalErr = err
		}

		metrics.TotalDuration += duration
		if duration < metrics.MinDuration {
			metrics.MinDuration = duration
		}
		if duration > metrics.MaxDuration {
			metrics.MaxDuration = duration
		}
	}

	if iterations > 0 {
		metrics.AvgDuration = metrics.TotalDuration / time.Duration(iterations)
	}

	return metrics, totalErr
}
