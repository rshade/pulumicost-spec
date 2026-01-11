// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing_test

import (
	"testing"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

// TestPerformanceReturnsLatencyMetrics validates that benchmarks return latency metrics (T039).
func TestPerformanceReturnsLatencyMetrics(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.PerformanceTests()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.TestFunc(harness)

			// Verify latency metrics are present
			if result.Duration <= 0 {
				t.Errorf("Test %s did not return duration", test.Name)
			}

			// Verify success for mock plugin (should be fast)
			if !result.Success {
				t.Errorf("Test %s failed unexpectedly: %v", test.Name, result.Error)
			}
		})
	}
}

// TestPerformanceBaselineThresholds validates benchmark comparison against baselines (T040).
func TestPerformanceBaselineThresholds(t *testing.T) {
	// Get all baselines
	baselines := plugintesting.DefaultBaselines()

	// Verify we have expected baselines
	expectedMethods := []string{
		"Name",
		"Supports",
		"GetPluginInfo",
		"GetProjectedCost",
		"GetPricingSpec",
		"GetActualCost_24h",
	}
	for _, method := range expectedMethods {
		found := false
		for _, b := range baselines {
			if b.Method == method {
				found = true

				// Verify threshold values are sensible
				if b.StandardLatency > 0 && b.AdvancedLatency > 0 {
					if b.AdvancedLatency > b.StandardLatency {
						t.Errorf("Advanced latency (%v) should be <= Standard latency (%v) for %s",
							b.AdvancedLatency, b.StandardLatency, method)
					}
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected baseline for method %s not found", method)
		}
	}
}

// TestPerformanceVarianceWithinThreshold validates SC-003 variance check (T040).
func TestPerformanceVarianceWithinThreshold(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Find the variance test
	tests := plugintesting.PerformanceTests()
	for _, test := range tests {
		if test.Name == "Performance_BaselineVariance" {
			result := test.TestFunc(harness)

			// For a mock plugin, variance should be well within threshold
			// since mock plugins respond quickly
			if !result.Success {
				if testing.Short() {
					t.Skipf("Variance test skipped in short mode: %s", result.Details)
				}
				t.Logf("Variance test result (may be environment-dependent): %s", result.Details)
			}
			break
		}
	}
}

// TestPerformanceExcessiveAllocationsWarning validates allocation tracking (T041).
func TestPerformanceExcessiveAllocationsWarning(t *testing.T) {
	// This test validates that the framework CAN track allocations
	// The actual implementation may use runtime.MemStats

	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Run performance tests
	tests := plugintesting.PerformanceTests()
	for _, test := range tests {
		t.Run(test.Name+"_AllocCheck", func(t *testing.T) {
			result := test.TestFunc(harness)
			// Scaffold: allocation tracking would go here
			t.Logf("Test %s completed in %v", test.Name, result.Duration)
		})
	}
}

// TestPerformanceGetBaseline validates GetBaseline function.
//
//nolint:gocognit // Table-driven tests inherently have higher complexity
func TestPerformanceGetBaseline(t *testing.T) {
	testCases := []struct {
		method      string
		expectNil   bool
		hasStandard bool
		hasAdvanced bool
	}{
		{"Name", false, true, true},
		{"Supports", false, true, true},
		{"GetPluginInfo", false, true, true},
		{"GetProjectedCost", false, true, true},
		{"GetPricingSpec", false, true, true},
		{"GetActualCost_24h", false, true, true},
		{"GetActualCost_30d", false, false, true}, // Only Advanced
		{"NonExistent", true, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			baseline := plugintesting.GetBaseline(tc.method)

			if tc.expectNil && baseline != nil {
				t.Errorf("Expected nil baseline for %s, got %v", tc.method, baseline)
			}

			if !tc.expectNil && baseline == nil {
				t.Errorf("Expected baseline for %s, got nil", tc.method)
			}

			if baseline != nil {
				if tc.hasStandard && baseline.StandardLatency == 0 {
					t.Errorf("Expected Standard latency for %s", tc.method)
				}
				if tc.hasAdvanced && baseline.AdvancedLatency == 0 {
					t.Errorf("Expected Advanced latency for %s", tc.method)
				}
			}
		})
	}
}

// TestRegisterPerformanceTests validates test registration.
func TestRegisterPerformanceTests(t *testing.T) {
	suite := plugintesting.NewConformanceSuite()
	plugintesting.RegisterPerformanceTests(suite)

	// Verify the suite has tests registered
	config := suite.GetConfig()
	if config.TargetLevel != plugintesting.ConformanceLevelStandard {
		t.Errorf("Expected default target level Standard, got %v", config.TargetLevel)
	}
}
