package testing

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// ConformanceLevel defines the level of conformance testing
type ConformanceLevel int

const (
	ConformanceBasic ConformanceLevel = iota
	ConformanceStandard
	ConformanceAdvanced
)

// ConformanceResult contains the result of a conformance test run
type ConformanceResult struct {
	Level        ConformanceLevel
	TotalTests   int
	PassedTests  int
	FailedTests  int
	SkippedTests int
	Results      []TestResult
	Summary      string
}

// RunBasicConformanceTests runs basic conformance tests that all plugins must pass
func RunBasicConformanceTests(t *testing.T, impl pbc.CostSourceServiceServer) *ConformanceResult {
	suite := NewConformanceSuite()

	// Basic test cases that all plugins MUST pass
	suite.AddTest(ConformanceTest{
		Name:        "NameReturnsValidResponse",
		Description: "Plugin must return a valid name",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			resp, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
			duration := time.Since(start)

			if err != nil {
				return TestResult{
					Method:   "Name",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "RPC call failed",
				}
			}

			if err := ValidateNameResponse(resp); err != nil {
				return TestResult{
					Method:   "Name",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "Response validation failed",
				}
			}

			return TestResult{
				Method:   "Name",
				Success:  true,
				Duration: duration,
				Details:  fmt.Sprintf("Plugin name: %s", resp.GetName()),
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "SupportsHandlesValidInput",
		Description: "Plugin must handle Supports requests correctly",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
			resp, err := harness.Client().Supports(context.Background(), &pbc.SupportsRequest{
				Resource: resource,
			})
			duration := time.Since(start)

			if err != nil {
				return TestResult{
					Method:   "Supports",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "RPC call failed",
				}
			}

			if err := ValidateSupportsResponse(resp); err != nil {
				return TestResult{
					Method:   "Supports",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "Response validation failed",
				}
			}

			return TestResult{
				Method:   "Supports",
				Success:  true,
				Duration: duration,
				Details:  fmt.Sprintf("Supported: %v, Reason: %s", resp.GetSupported(), resp.GetReason()),
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "SupportsHandlesNilResource",
		Description: "Plugin must handle nil resource gracefully",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			resp, err := harness.Client().Supports(context.Background(), &pbc.SupportsRequest{
				Resource: nil,
			})
			duration := time.Since(start)

			if err != nil {
				// Error is acceptable for nil resource
				return TestResult{
					Method:   "Supports",
					Success:  true,
					Duration: duration,
					Details:  "Correctly rejected nil resource with error",
				}
			}

			if resp.GetSupported() {
				return TestResult{
					Method:   "Supports",
					Success:  false,
					Error:    fmt.Errorf("plugin incorrectly supports nil resource"),
					Duration: duration,
					Details:  "Should not support nil resource",
				}
			}

			return TestResult{
				Method:   "Supports",
				Success:  true,
				Duration: duration,
				Details:  "Correctly rejected nil resource",
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "GetProjectedCostHandlesValidResource",
		Description: "Plugin must handle GetProjectedCost for valid resources",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
			
			// First check if resource is supported
			supportsResp, err := harness.Client().Supports(context.Background(), &pbc.SupportsRequest{
				Resource: resource,
			})
			if err != nil {
				return TestResult{
					Method:   "GetProjectedCost",
					Success:  false,
					Error:    err,
					Duration: time.Since(start),
					Details:  "Failed to check resource support",
				}
			}

			if !supportsResp.GetSupported() {
				// If resource is not supported, GetProjectedCost should error
				_, err := harness.Client().GetProjectedCost(context.Background(), &pbc.GetProjectedCostRequest{
					Resource: resource,
				})
				duration := time.Since(start)

				if err != nil {
					return TestResult{
						Method:   "GetProjectedCost",
						Success:  true,
						Duration: duration,
						Details:  "Correctly rejected unsupported resource",
					}
				}

				return TestResult{
					Method:   "GetProjectedCost",
					Success:  false,
					Error:    fmt.Errorf("plugin returned cost for unsupported resource"),
					Duration: duration,
					Details:  "Should reject unsupported resources",
				}
			}

			// Resource is supported, should return valid cost
			resp, err := harness.Client().GetProjectedCost(context.Background(), &pbc.GetProjectedCostRequest{
				Resource: resource,
			})
			duration := time.Since(start)

			if err != nil {
				return TestResult{
					Method:   "GetProjectedCost",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "RPC call failed for supported resource",
				}
			}

			if err := ValidateProjectedCostResponse(resp); err != nil {
				return TestResult{
					Method:   "GetProjectedCost",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "Response validation failed",
				}
			}

			return TestResult{
				Method:   "GetProjectedCost",
				Success:  true,
				Duration: duration,
				Details:  fmt.Sprintf("Unit price: %.6f %s", resp.GetUnitPrice(), resp.GetCurrency()),
			}
		},
	})

	results := suite.RunTests(t, impl)
	
	passed := 0
	failed := 0
	for _, result := range results {
		if result.Success {
			passed++
		} else {
			failed++
		}
	}

	return &ConformanceResult{
		Level:       ConformanceBasic,
		TotalTests:  len(results),
		PassedTests: passed,
		FailedTests: failed,
		Results:     results,
		Summary:     fmt.Sprintf("Basic conformance: %d/%d tests passed", passed, len(results)),
	}
}

// RunStandardConformanceTests runs standard conformance tests for production-ready plugins
func RunStandardConformanceTests(t *testing.T, impl pbc.CostSourceServiceServer) *ConformanceResult {
	suite := NewConformanceSuite()

	// Run basic tests first
	basicResult := RunBasicConformanceTests(t, impl)
	if basicResult.FailedTests > 0 {
		return &ConformanceResult{
			Level:       ConformanceStandard,
			TotalTests:  basicResult.TotalTests,
			PassedTests: basicResult.PassedTests,
			FailedTests: basicResult.FailedTests,
			Results:     basicResult.Results,
			Summary:     "Standard conformance failed: basic tests must pass first",
		}
	}

	// Standard-level tests
	suite.AddTest(ConformanceTest{
		Name:        "GetActualCostHandlesValidTimeRange",
		Description: "Plugin must handle GetActualCost with valid time ranges",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			timeStart, timeEnd := CreateTimeRange(24)
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
					return TestResult{
						Method:   "GetActualCost",
						Success:  true,
						Duration: duration,
						Details:  "Correctly indicated no data available",
					}
				}

				return TestResult{
					Method:   "GetActualCost",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "Unexpected error type",
				}
			}

			if err := ValidateActualCostResponse(resp); err != nil {
				return TestResult{
					Method:   "GetActualCost",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "Response validation failed",
				}
			}

			return TestResult{
				Method:   "GetActualCost",
				Success:  true,
				Duration: duration,
				Details:  fmt.Sprintf("Returned %d cost data points", len(resp.GetResults())),
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "GetActualCostRejectsInvalidTimeRange",
		Description: "Plugin must reject invalid time ranges",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			// Swap start and end to create invalid range
			timeEnd, timeStart := CreateTimeRange(24)
			_, err := harness.Client().GetActualCost(context.Background(), &pbc.GetActualCostRequest{
				ResourceId: "test-resource",
				Start:      timeStart,
				End:        timeEnd,
			})
			duration := time.Since(start)

			if err == nil {
				return TestResult{
					Method:   "GetActualCost",
					Success:  false,
					Error:    fmt.Errorf("plugin accepted invalid time range"),
					Duration: duration,
					Details:  "Should reject end time before start time",
				}
			}

			return TestResult{
				Method:   "GetActualCost",
				Success:  true,
				Duration: duration,
				Details:  "Correctly rejected invalid time range",
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "GetPricingSpecConsistency",
		Description: "Plugin must return consistent pricing specs for same resource",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			resource := CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

			// Get pricing spec multiple times
			var firstSpec *pbc.PricingSpec
			for i := 0; i < 3; i++ {
				resp, err := harness.Client().GetPricingSpec(context.Background(), &pbc.GetPricingSpecRequest{
					Resource: resource,
				})
				if err != nil {
					return TestResult{
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
					spec := resp.GetSpec()
					if spec.GetRatePerUnit() != firstSpec.GetRatePerUnit() {
						return TestResult{
							Method:   "GetPricingSpec",
							Success:  false,
							Error:    fmt.Errorf("inconsistent rate: %.6f vs %.6f", firstSpec.GetRatePerUnit(), spec.GetRatePerUnit()),
							Duration: time.Since(start),
							Details:  "Rate per unit should be consistent",
						}
					}
					if spec.GetCurrency() != firstSpec.GetCurrency() {
						return TestResult{
							Method:   "GetPricingSpec",
							Success:  false,
							Error:    fmt.Errorf("inconsistent currency: %s vs %s", firstSpec.GetCurrency(), spec.GetCurrency()),
							Duration: time.Since(start),
							Details:  "Currency should be consistent",
						}
					}
				}
			}

			return TestResult{
				Method:   "GetPricingSpec",
				Success:  true,
				Duration: time.Since(start),
				Details:  "Pricing spec is consistent across multiple calls",
			}
		},
	})

	results := suite.RunTests(t, impl)
	
	// Combine with basic results
	allResults := append(basicResult.Results, results...)
	passed := basicResult.PassedTests
	failed := basicResult.FailedTests
	
	for _, result := range results {
		if result.Success {
			passed++
		} else {
			failed++
		}
	}

	return &ConformanceResult{
		Level:       ConformanceStandard,
		TotalTests:  len(allResults),
		PassedTests: passed,
		FailedTests: failed,
		Results:     allResults,
		Summary:     fmt.Sprintf("Standard conformance: %d/%d tests passed", passed, len(allResults)),
	}
}

// RunAdvancedConformanceTests runs advanced conformance tests for high-performance plugins
func RunAdvancedConformanceTests(t *testing.T, impl pbc.CostSourceServiceServer) *ConformanceResult {
	// First run standard tests
	standardResult := RunStandardConformanceTests(t, impl)
	if standardResult.FailedTests > 0 {
		return &ConformanceResult{
			Level:       ConformanceAdvanced,
			TotalTests:  standardResult.TotalTests,
			PassedTests: standardResult.PassedTests,
			FailedTests: standardResult.FailedTests,
			Results:     standardResult.Results,
			Summary:     "Advanced conformance failed: standard tests must pass first",
		}
	}

	suite := NewConformanceSuite()

	// Advanced performance and reliability tests
	suite.AddTest(ConformanceTest{
		Name:        "PerformanceBaseline",
		Description: "Plugin must meet minimum performance requirements",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			
			// Test Name performance (should be fast)
			nameStart := time.Now()
			_, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
			nameDuration := time.Since(nameStart)
			
			if err != nil {
				return TestResult{
					Method:   "Performance",
					Success:  false,
					Error:    err,
					Duration: time.Since(start),
					Details:  "Name RPC failed",
				}
			}

			// Name should respond within 100ms
			if nameDuration > 100*time.Millisecond {
				return TestResult{
					Method:   "Performance",
					Success:  false,
					Error:    fmt.Errorf("name RPC too slow: %v", nameDuration),
					Duration: time.Since(start),
					Details:  "Name RPC should respond within 100ms",
				}
			}

			return TestResult{
				Method:   "Performance",
				Success:  true,
				Duration: time.Since(start),
				Details:  fmt.Sprintf("Name RPC responded in %v", nameDuration),
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "ConcurrentRequestHandling",
		Description: "Plugin must handle concurrent requests safely",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			
			const numConcurrent = 10
			errors := make(chan error, numConcurrent)
			
			// Launch concurrent Name requests
			for i := 0; i < numConcurrent; i++ {
				go func() {
					_, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
					errors <- err
				}()
			}

			// Collect results
			for i := 0; i < numConcurrent; i++ {
				if err := <-errors; err != nil {
					return TestResult{
						Method:   "Concurrency",
						Success:  false,
						Error:    err,
						Duration: time.Since(start),
						Details:  fmt.Sprintf("Concurrent request %d failed", i),
					}
				}
			}

			return TestResult{
				Method:   "Concurrency",
				Success:  true,
				Duration: time.Since(start),
				Details:  fmt.Sprintf("Successfully handled %d concurrent requests", numConcurrent),
			}
		},
	})

	suite.AddTest(ConformanceTest{
		Name:        "LargeDataHandling",
		Description: "Plugin must handle large datasets efficiently",
		TestFunc: func(harness *TestHarness) TestResult {
			start := time.Now()
			
			// Request 30 days of data (should be a reasonable large dataset)
			timeStart, timeEnd := CreateTimeRange(720) // 30 days
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
					return TestResult{
						Method:   "LargeData",
						Success:  true,
						Duration: duration,
						Details:  "Correctly indicated large dataset not supported",
					}
				}

				return TestResult{
					Method:   "LargeData",
					Success:  false,
					Error:    err,
					Duration: duration,
					Details:  "Unexpected error for large dataset",
				}
			}

			// Should respond within reasonable time (10 seconds)
			if duration > 10*time.Second {
				return TestResult{
					Method:   "LargeData",
					Success:  false,
					Error:    fmt.Errorf("large dataset query too slow: %v", duration),
					Duration: duration,
					Details:  "Large dataset queries should complete within 10 seconds",
				}
			}

			return TestResult{
				Method:   "LargeData",
				Success:  true,
				Duration: duration,
				Details:  fmt.Sprintf("Handled large dataset (%d points) in %v", len(resp.GetResults()), duration),
			}
		},
	})

	results := suite.RunTests(t, impl)
	
	// Combine with standard results
	allResults := append(standardResult.Results, results...)
	passed := standardResult.PassedTests
	failed := standardResult.FailedTests
	
	for _, result := range results {
		if result.Success {
			passed++
		} else {
			failed++
		}
	}

	return &ConformanceResult{
		Level:       ConformanceAdvanced,
		TotalTests:  len(allResults),
		PassedTests: passed,
		FailedTests: failed,
		Results:     allResults,
		Summary:     fmt.Sprintf("Advanced conformance: %d/%d tests passed", passed, len(allResults)),
	}
}

// PrintConformanceReport prints a detailed conformance test report
func PrintConformanceReport(result *ConformanceResult) {
	fmt.Printf("\n=== CONFORMANCE TEST REPORT ===\n")
	fmt.Printf("Level: %s\n", conformanceLevelString(result.Level))
	fmt.Printf("Total Tests: %d\n", result.TotalTests)
	fmt.Printf("Passed: %d\n", result.PassedTests)
	fmt.Printf("Failed: %d\n", result.FailedTests)
	if result.SkippedTests > 0 {
		fmt.Printf("Skipped: %d\n", result.SkippedTests)
	}
	fmt.Printf("Success Rate: %.1f%%\n", float64(result.PassedTests)/float64(result.TotalTests)*100)
	fmt.Printf("Summary: %s\n", result.Summary)

	if result.FailedTests > 0 {
		fmt.Printf("\n--- FAILED TESTS ---\n")
		for _, testResult := range result.Results {
			if !testResult.Success {
				fmt.Printf("❌ %s: %v (%s)\n", testResult.Method, testResult.Error, testResult.Details)
			}
		}
	}

	fmt.Printf("\n--- ALL TEST RESULTS ---\n")
	for _, testResult := range result.Results {
		status := "✅"
		if !testResult.Success {
			status = "❌"
		}
		fmt.Printf("%s %s (%v) - %s\n", status, testResult.Method, testResult.Duration, testResult.Details)
	}
	fmt.Printf("===============================\n\n")
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

// ConformanceTestMain provides a standard main function for plugin conformance testing
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
	
	PrintConformanceReport(result)
	
	if result.FailedTests > 0 {
		fmt.Printf("❌ Conformance tests failed. Plugin does not meet %s conformance requirements.\n", 
			strings.ToLower(conformanceLevelString(level)))
		return
	}
	
	fmt.Printf("✅ Plugin successfully meets %s conformance requirements!\n", 
		strings.ToLower(conformanceLevelString(level)))
}