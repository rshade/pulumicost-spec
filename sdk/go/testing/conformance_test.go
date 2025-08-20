package testing_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// ConformanceLevel defines the level of conformance testing.
type ConformanceLevel int

const (
	ConformanceBasic ConformanceLevel = iota
	ConformanceStandard
	ConformanceAdvanced
)

// ConformanceResult contains the result of a conformance test run.
type ConformanceResult struct {
	Level        ConformanceLevel
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	Results      []plugintesting.TestResult
	Summary      string
}

func NameValidationTest(harness *plugintesting.TestHarness) plugintesting.TestResult {
	start := time.Now()
	resp, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
	duration := time.Since(start)

	if err != nil {
		return plugintesting.TestResult{
			Method:   "Name",
			Success:  false,
			Error:    err,
			Duration: duration,
			Details:  "RPC call failed",
		}
	}

	if validationErr := plugintesting.ValidateNameResponse(resp); validationErr != nil {
		return plugintesting.TestResult{
			Method:   "Name",
			Success:  false,
			Error:    validationErr,
			Duration: duration,
			Details:  "Response validation failed",
		}
	}

	return plugintesting.TestResult{
		Method:   "Name",
		Success:  true,
		Duration: duration,
		Details:  fmt.Sprintf("Plugin name: %s", resp.GetName()),
	}
}

// RunBasicConformanceTests runs basic conformance tests that all plugins must pass.
func RunBasicConformanceTests(t *testing.T, impl pbc.CostSourceServiceServer) *ConformanceResult {
	suite := plugintesting.NewConformanceSuite()
	addBasicConformanceTests(suite)
	return runConformanceTestSuite(t, impl, suite, ConformanceBasic, "Basic conformance")
}

func addBasicConformanceTests(suite *plugintesting.PluginConformanceSuite) {
	// Basic test cases that all plugins MUST pass
	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "NameReturnsValidResponse",
		Description: "Plugin must return a valid name",
		TestFunc:    NameValidationTest,
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "SupportsHandlesValidInput",
		Description: "Plugin must handle Supports requests correctly",
		TestFunc:    createSupportsValidInputTest(),
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "SupportsHandlesNilResource",
		Description: "Plugin must handle nil resource gracefully",
		TestFunc:    createSupportsNilResourceTest(),
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "GetProjectedCostHandlesValidResource",
		Description: "Plugin must handle GetProjectedCost for valid resources",
		TestFunc:    createProjectedCostValidResourceTest(),
	})
}

func runConformanceTestSuite(
	t *testing.T,
	impl pbc.CostSourceServiceServer,
	suite *plugintesting.PluginConformanceSuite,
	level ConformanceLevel,
	summaryPrefix string,
) *ConformanceResult {
	results := suite.RunTests(t, impl)
	passed, failed := countTestResults(results)

	return &ConformanceResult{
		Level:       level,
		TotalTests:  len(results),
		PassedTests: passed,
		FailedTests: failed,
		Results:     results,
		Summary:     fmt.Sprintf("%s: %d/%d tests passed", summaryPrefix, passed, len(results)),
	}
}

func countTestResults(results []plugintesting.TestResult) (int, int) {
	var passed, failed int
	for _, result := range results {
		if result.Success {
			passed++
		} else {
			failed++
		}
	}
	return passed, failed
}

func createSupportsValidInputTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
		resp, err := harness.Client().Supports(context.Background(), &pbc.SupportsRequest{
			Resource: resource,
		})
		duration := time.Since(start)

		if err != nil {
			return plugintesting.TestResult{
				Method:   "Supports",
				Success:  false,
				Error:    err,
				Duration: duration,
				Details:  "RPC call failed",
			}
		}

		if validationErr := plugintesting.ValidateSupportsResponse(resp); validationErr != nil {
			return plugintesting.TestResult{
				Method:   "Supports",
				Success:  false,
				Error:    validationErr,
				Duration: duration,
				Details:  "Response validation failed",
			}
		}

		return plugintesting.TestResult{
			Method:   "Supports",
			Success:  true,
			Duration: duration,
			Details: fmt.Sprintf(
				"Supported: %v, Reason: %s",
				resp.GetSupported(),
				resp.GetReason(),
			),
		}
	}
}

func createSupportsNilResourceTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()
		resp, err := harness.Client().Supports(context.Background(), &pbc.SupportsRequest{
			Resource: nil,
		})
		duration := time.Since(start)

		if err != nil {
			// Error is acceptable for nil resource
			return plugintesting.TestResult{
				Method:   "Supports",
				Success:  true,
				Duration: duration,
				Details:  "Correctly rejected nil resource with error",
			}
		}

		if resp.GetSupported() {
			return plugintesting.TestResult{
				Method:   "Supports",
				Success:  false,
				Error:    errors.New("plugin incorrectly supports nil resource"),
				Duration: duration,
				Details:  "Should not support nil resource",
			}
		}

		return plugintesting.TestResult{
			Method:   "Supports",
			Success:  true,
			Duration: duration,
			Details:  "Correctly rejected nil resource",
		}
	}
}

func createProjectedCostValidResourceTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

		// Check resource support first
		if !checkResourceSupport(harness, resource) {
			return handleUnsupportedResource(harness, resource, start)
		}

		// Test supported resource
		return testSupportedResource(harness, resource, start)
	}
}

func checkResourceSupport(
	harness *plugintesting.TestHarness,
	resource *pbc.ResourceDescriptor,
) bool {
	supportsResp, err := harness.Client().Supports(context.Background(), &pbc.SupportsRequest{
		Resource: resource,
	})
	return err == nil && supportsResp.GetSupported()
}

func handleUnsupportedResource(
	harness *plugintesting.TestHarness,
	resource *pbc.ResourceDescriptor,
	start time.Time,
) plugintesting.TestResult {
	_, projectedErr := harness.Client().
		GetProjectedCost(context.Background(), &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
	duration := time.Since(start)

	if projectedErr != nil {
		return plugintesting.TestResult{
			Method:   "GetProjectedCost",
			Success:  true,
			Duration: duration,
			Details:  "Correctly rejected unsupported resource",
		}
	}

	return plugintesting.TestResult{
		Method:   "GetProjectedCost",
		Success:  false,
		Error:    errors.New("plugin returned cost for unsupported resource"),
		Duration: duration,
		Details:  "Should reject unsupported resources",
	}
}

func testSupportedResource(
	harness *plugintesting.TestHarness,
	resource *pbc.ResourceDescriptor,
	start time.Time,
) plugintesting.TestResult {
	resp, err := harness.Client().
		GetProjectedCost(context.Background(), &pbc.GetProjectedCostRequest{
			Resource: resource,
		})
	duration := time.Since(start)

	if err != nil {
		return plugintesting.TestResult{
			Method:   "GetProjectedCost",
			Success:  false,
			Error:    err,
			Duration: duration,
			Details:  "RPC call failed for supported resource",
		}
	}

	if validationErr := plugintesting.ValidateProjectedCostResponse(resp); validationErr != nil {
		return plugintesting.TestResult{
			Method:   "GetProjectedCost",
			Success:  false,
			Error:    validationErr,
			Duration: duration,
			Details:  "Response validation failed",
		}
	}

	return plugintesting.TestResult{
		Method:   "GetProjectedCost",
		Success:  true,
		Duration: duration,
		Details: fmt.Sprintf(
			"Unit price: %.6f %s",
			resp.GetUnitPrice(),
			resp.GetCurrency(),
		),
	}
}

// RunStandardConformanceTests runs standard conformance tests for production-ready plugins.
func RunStandardConformanceTests(
	t *testing.T,
	impl pbc.CostSourceServiceServer,
) *ConformanceResult {
	// Run basic tests first
	basicResult := RunBasicConformanceTests(t, impl)
	if basicResult.FailedTests > 0 {
		return createFailedResult(
			ConformanceStandard,
			basicResult,
			"Standard conformance failed: basic tests must pass first",
		)
	}

	suite := plugintesting.NewConformanceSuite()
	addStandardConformanceTests(suite)
	results := suite.RunTests(t, impl)

	return combineConformanceResults(
		basicResult,
		results,
		ConformanceStandard,
		"Standard conformance",
	)
}

func createFailedResult(
	level ConformanceLevel,
	basicResult *ConformanceResult,
	summary string,
) *ConformanceResult {
	return &ConformanceResult{
		Level:       level,
		TotalTests:  basicResult.TotalTests,
		PassedTests: basicResult.PassedTests,
		FailedTests: basicResult.FailedTests,
		Results:     basicResult.Results,
		Summary:     summary,
	}
}

func combineConformanceResults(
	basicResult *ConformanceResult,
	additionalResults []plugintesting.TestResult,
	level ConformanceLevel,
	summaryPrefix string,
) *ConformanceResult {
	allResults := make([]plugintesting.TestResult, len(basicResult.Results)+len(additionalResults))
	copy(allResults, basicResult.Results)
	copy(allResults[len(basicResult.Results):], additionalResults)

	passed := basicResult.PassedTests
	failed := basicResult.FailedTests
	additionalPassed, additionalFailed := countTestResults(additionalResults)
	passed += additionalPassed
	failed += additionalFailed

	return &ConformanceResult{
		Level:       level,
		TotalTests:  len(allResults),
		PassedTests: passed,
		FailedTests: failed,
		Results:     allResults,
		Summary:     fmt.Sprintf("%s: %d/%d tests passed", summaryPrefix, passed, len(allResults)),
	}
}

func addStandardConformanceTests(suite *plugintesting.PluginConformanceSuite) {
	// Standard-level tests
	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "GetActualCostHandlesValidTimeRange",
		Description: "Plugin must handle GetActualCost with valid time ranges",
		TestFunc:    createActualCostValidTimeRangeTest(),
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "GetActualCostRejectsInvalidTimeRange",
		Description: "Plugin must reject invalid time ranges",
		TestFunc:    createActualCostInvalidTimeRangeTest(),
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "GetPricingSpecConsistency",
		Description: "Plugin must return consistent pricing specs for same resource",
		TestFunc:    createPricingSpecConsistencyTest(),
	})
}

// RunAdvancedConformanceTests runs advanced conformance tests for high-performance plugins.
func RunAdvancedConformanceTests(
	t *testing.T,
	impl pbc.CostSourceServiceServer,
) *ConformanceResult {
	// First run standard tests
	standardResult := RunStandardConformanceTests(t, impl)
	if standardResult.FailedTests > 0 {
		return createFailedResult(
			ConformanceAdvanced,
			standardResult,
			"Advanced conformance failed: standard tests must pass first",
		)
	}

	suite := plugintesting.NewConformanceSuite()
	addAdvancedConformanceTests(suite)
	results := suite.RunTests(t, impl)

	return combineConformanceResults(
		standardResult,
		results,
		ConformanceAdvanced,
		"Advanced conformance",
	)
}

func addAdvancedConformanceTests(suite *plugintesting.PluginConformanceSuite) {
	// Advanced performance and reliability tests
	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "PerformanceBaseline",
		Description: "Plugin must meet minimum performance requirements",
		TestFunc:    createPerformanceBaselineTest(),
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "ConcurrentRequestHandling",
		Description: "Plugin must handle concurrent requests safely",
		TestFunc:    createConcurrentRequestTest(),
	})

	suite.AddTest(plugintesting.ConformanceTest{
		Name:        "LargeDataHandling",
		Description: "Plugin must handle large datasets efficiently",
		TestFunc:    createLargeDataHandlingTest(),
	})
}

// PrintConformanceReport prints a detailed conformance test report.
func PrintConformanceReport(t *testing.T, result *ConformanceResult) {
	t.Log("\n=== CONFORMANCE TEST REPORT ===")
	t.Logf("Level: %s", conformanceLevelString(result.Level))
	t.Logf("Total Tests: %d", result.TotalTests)
	t.Logf("Passed: %d", result.PassedTests)
	t.Logf("Failed: %d", result.FailedTests)
	if result.SkippedTests > 0 {
		t.Logf("Skipped: %d", result.SkippedTests)
	}
	t.Logf(
		"Success Rate: %.1f%%",
		float64(result.PassedTests)/float64(result.TotalTests)*plugintesting.SuccessRateMultiplier,
	)
	t.Logf("Summary: %s", result.Summary)

	if result.FailedTests > 0 {
		t.Log("\n--- FAILED TESTS ---")
		for _, testResult := range result.Results {
			if !testResult.Success {
				t.Logf("❌ %s: %v (%s)", testResult.Method, testResult.Error, testResult.Details)
			}
		}
	}

	t.Log("\n--- ALL TEST RESULTS ---")
	for _, testResult := range result.Results {
		status := "✅"
		if !testResult.Success {
			status = "❌"
		}
		t.Logf(
			"%s %s (%v) - %s",
			status,
			testResult.Method,
			testResult.Duration,
			testResult.Details,
		)
	}
	t.Log("===============================")
}

func conformanceLevelString(level ConformanceLevel) string {
	switch level {
	case ConformanceBasic:
		return "Basic"
	case ConformanceStandard:
		return "Standard"
	case ConformanceAdvanced:
		return "Advanced"
	default:
		return "Unknown"
	}
}

// ConformanceTestMain provides a standard main function for plugin conformance testing.
func ConformanceTestMain(impl pbc.CostSourceServiceServer, level ConformanceLevel) {
	t := &testing.T{}

	var result *ConformanceResult

	switch level {
	case ConformanceBasic:
		result = RunBasicConformanceTests(t, impl)
	case ConformanceStandard:
		result = RunStandardConformanceTests(t, impl)
	case ConformanceAdvanced:
		result = RunAdvancedConformanceTests(t, impl)
	}

	PrintConformanceReport(t, result)

	if result.FailedTests > 0 {
		t.Logf("❌ Conformance tests failed. Plugin does not meet %s conformance requirements.",
			strings.ToLower(conformanceLevelString(level)))
		return
	}

	t.Logf("✅ Plugin successfully meets %s conformance requirements!",
		strings.ToLower(conformanceLevelString(level)))
}

func createActualCostValidTimeRangeTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()
		timeStart, timeEnd := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)
		resp, err := harness.Client().GetActualCost(context.Background(), &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      timeStart,
			End:        timeEnd,
		})
		duration := time.Since(start)

		if err != nil {
			// Error is acceptable if no data is available
			st, ok := status.FromError(err)
			if ok && (st.Code() == codes.NotFound || st.Code() == codes.Unavailable) {
				return plugintesting.TestResult{
					Method:   "GetActualCost",
					Success:  true,
					Duration: duration,
					Details:  "Correctly indicated no data available",
				}
			}

			return plugintesting.TestResult{
				Method:   "GetActualCost",
				Success:  false,
				Error:    err,
				Duration: duration,
				Details:  "Unexpected error type",
			}
		}

		if validationErr := plugintesting.ValidateActualCostResponse(resp); validationErr != nil {
			return plugintesting.TestResult{
				Method:   "GetActualCost",
				Success:  false,
				Error:    validationErr,
				Duration: duration,
				Details:  "Response validation failed",
			}
		}

		return plugintesting.TestResult{
			Method:   "GetActualCost",
			Success:  true,
			Duration: duration,
			Details:  fmt.Sprintf("Returned %d cost data points", len(resp.GetResults())),
		}
	}
}

func createActualCostInvalidTimeRangeTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()
		// Swap start and end to create invalid range
		timeEnd, timeStart := plugintesting.CreateTimeRange(plugintesting.HoursPerDay)
		_, err := harness.Client().GetActualCost(context.Background(), &pbc.GetActualCostRequest{
			ResourceId: "test-resource",
			Start:      timeStart,
			End:        timeEnd,
		})
		duration := time.Since(start)

		if err == nil {
			return plugintesting.TestResult{
				Method:   "GetActualCost",
				Success:  false,
				Error:    errors.New("plugin accepted invalid time range"),
				Duration: duration,
				Details:  "Should reject end time before start time",
			}
		}

		return plugintesting.TestResult{
			Method:   "GetActualCost",
			Success:  true,
			Duration: duration,
			Details:  "Correctly rejected invalid time range",
		}
	}
}

func createPricingSpecConsistencyTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()
		resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

		// Get pricing spec multiple times
		var firstSpec *pbc.PricingSpec
		for i := range plugintesting.NumConsistencyChecks {
			resp, err := harness.Client().
				GetPricingSpec(context.Background(), &pbc.GetPricingSpecRequest{
					Resource: resource,
				})
			if err != nil {
				return plugintesting.TestResult{
					Method:   "GetPricingSpec",
					Success:  false,
					Error:    err,
					Duration: time.Since(start),
					Details:  fmt.Sprintf("Failed on iteration %d", i),
				}
			}

			if i == 0 {
				firstSpec = resp.GetSpec()
			} else {
				if inconsistencyErr := checkSpecConsistency(firstSpec, resp.GetSpec(), start); inconsistencyErr != nil {
					return *inconsistencyErr
				}
			}
		}

		return plugintesting.TestResult{
			Method:   "GetPricingSpec",
			Success:  true,
			Duration: time.Since(start),
			Details:  "Pricing spec is consistent across multiple calls",
		}
	}
}

func checkSpecConsistency(
	firstSpec, currentSpec *pbc.PricingSpec,
	start time.Time,
) *plugintesting.TestResult {
	if currentSpec.GetRatePerUnit() != firstSpec.GetRatePerUnit() {
		return &plugintesting.TestResult{
			Method:  "GetPricingSpec",
			Success: false,
			Error: fmt.Errorf(
				"inconsistent rate: %.6f vs %.6f",
				firstSpec.GetRatePerUnit(),
				currentSpec.GetRatePerUnit(),
			),
			Duration: time.Since(start),
			Details:  "Rate per unit should be consistent",
		}
	}
	if currentSpec.GetCurrency() != firstSpec.GetCurrency() {
		return &plugintesting.TestResult{
			Method:  "GetPricingSpec",
			Success: false,
			Error: fmt.Errorf(
				"inconsistent currency: %s vs %s",
				firstSpec.GetCurrency(),
				currentSpec.GetCurrency(),
			),
			Duration: time.Since(start),
			Details:  "Currency should be consistent",
		}
	}
	return nil
}

func createPerformanceBaselineTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()

		// Test Name performance (should be fast)
		nameStart := time.Now()
		_, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
		nameDuration := time.Since(nameStart)

		if err != nil {
			return plugintesting.TestResult{
				Method:   "Performance",
				Success:  false,
				Error:    err,
				Duration: time.Since(start),
				Details:  "Name RPC failed",
			}
		}

		// Name should respond within 100ms
		if nameDuration > plugintesting.MaxResponseTimeMs*time.Millisecond {
			return plugintesting.TestResult{
				Method:   "Performance",
				Success:  false,
				Error:    fmt.Errorf("name RPC too slow: %v", nameDuration),
				Duration: time.Since(start),
				Details:  "Name RPC should respond within 100ms",
			}
		}

		return plugintesting.TestResult{
			Method:   "Performance",
			Success:  true,
			Duration: time.Since(start),
			Details:  fmt.Sprintf("Name RPC responded in %v", nameDuration),
		}
	}
}

func createConcurrentRequestTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()

		const numConcurrent = plugintesting.NumConcurrentRequests
		errors := make(chan error, numConcurrent)

		// Launch concurrent Name requests
		for range numConcurrent {
			go func() {
				_, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
				errors <- err
			}()
		}

		// Collect results
		for i := range numConcurrent {
			if err := <-errors; err != nil {
				return plugintesting.TestResult{
					Method:   "Concurrency",
					Success:  false,
					Error:    err,
					Duration: time.Since(start),
					Details:  fmt.Sprintf("Concurrent request %d failed", i),
				}
			}
		}

		return plugintesting.TestResult{
			Method:   "Concurrency",
			Success:  true,
			Duration: time.Since(start),
			Details:  fmt.Sprintf("Successfully handled %d concurrent requests", numConcurrent),
		}
	}
}

func createLargeDataHandlingTest() func(*plugintesting.TestHarness) plugintesting.TestResult {
	return func(harness *plugintesting.TestHarness) plugintesting.TestResult {
		start := time.Now()

		// Request 30 days of data (should be a reasonable large dataset)
		timeStart, timeEnd := plugintesting.CreateTimeRange(plugintesting.HoursIn30Days) // 30 days
		resp, err := harness.Client().GetActualCost(context.Background(), &pbc.GetActualCostRequest{
			ResourceId: "large-dataset-test",
			Start:      timeStart,
			End:        timeEnd,
		})
		duration := time.Since(start)

		if err != nil {
			// Error is acceptable if plugin doesn't support large datasets
			st, ok := status.FromError(err)
			if ok && (st.Code() == codes.InvalidArgument || st.Code() == codes.ResourceExhausted) {
				return plugintesting.TestResult{
					Method:   "LargeData",
					Success:  true,
					Duration: duration,
					Details:  "Correctly indicated large dataset not supported",
				}
			}

			return plugintesting.TestResult{
				Method:   "LargeData",
				Success:  false,
				Error:    err,
				Duration: duration,
				Details:  "Unexpected error for large dataset",
			}
		}

		// Should respond within reasonable time (10 seconds)
		if duration > plugintesting.MaxLargeQueryTimeSeconds*time.Second {
			return plugintesting.TestResult{
				Method:   "LargeData",
				Success:  false,
				Error:    fmt.Errorf("large dataset query too slow: %v", duration),
				Duration: duration,
				Details:  "Large dataset queries should complete within 10 seconds",
			}
		}

		return plugintesting.TestResult{
			Method:   "LargeData",
			Success:  true,
			Duration: duration,
			Details: fmt.Sprintf(
				"Handled large dataset (%d points) in %v",
				len(resp.GetResults()),
				duration,
			),
		}
	}
}
