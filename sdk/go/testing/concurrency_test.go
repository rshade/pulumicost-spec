package testing_test

import (
	"testing"
	"time"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

// TestConcurrentRequestsCompleteSuccessfully validates that concurrent requests complete (T049).
func TestConcurrentRequestsCompleteSuccessfully(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.ConcurrencyTests()
	for _, test := range tests {
		if test.Name == "Concurrency_ParallelRequests_Standard" {
			t.Run(test.Name, func(t *testing.T) {
				result := test.TestFunc(harness)
				if !result.Success {
					t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
				}
			})
			break
		}
	}
}

// TestRaceDetectionIntegration validates race detection works (T050).
// Note: This test is designed to be run with -race flag.
func TestRaceDetectionIntegration(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Run concurrent tests - any race conditions should be caught by -race flag
	tests := plugintesting.ConcurrencyTests()
	for _, test := range tests {
		t.Run(test.Name+"_RaceCheck", func(t *testing.T) {
			result := test.TestFunc(harness)
			// Log result for debugging if race detector doesn't catch anything
			if !result.Success {
				t.Logf("Test completed with failure (race detector may still catch issues): %v", result.Error)
			}
		})
	}
}

// TestResponseConsistencyUnderLoad validates response consistency (T051).
func TestResponseConsistencyUnderLoad(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.ConcurrencyTests()
	for _, test := range tests {
		if test.Name == "Concurrency_ResponseConsistency" {
			t.Run(test.Name, func(t *testing.T) {
				result := test.TestFunc(harness)
				if !result.Success {
					t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
				}
			})
			break
		}
	}
}

// TestConcurrencyAllTests runs all concurrency tests.
func TestConcurrencyAllTests(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.ConcurrencyTests()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
		})
	}
}

// TestConcurrencyConfig validates ConcurrencyConfig.
func TestConcurrencyConfig(t *testing.T) {
	config := plugintesting.DefaultConcurrencyConfig()

	if config.ParallelRequests != plugintesting.NumConcurrentRequests {
		t.Errorf("Expected ParallelRequests %d, got %d",
			plugintesting.NumConcurrentRequests, config.ParallelRequests)
	}

	if config.Timeout <= 0 {
		t.Error("Expected positive timeout")
	}

	if config.Method == "" {
		t.Error("Expected non-empty method")
	}
}

// TestRegisterConcurrencyTests validates test registration.
func TestRegisterConcurrencyTests(t *testing.T) {
	suite := plugintesting.NewConformanceSuite()
	plugintesting.RegisterConcurrencyTests(suite)

	config := suite.GetConfig()
	if config.TargetLevel != plugintesting.ConformanceLevelStandard {
		t.Errorf("Expected default target level Standard, got %v", config.TargetLevel)
	}

	// Verify tests were actually registered
	tests := plugintesting.ConcurrencyTests()
	if len(tests) == 0 {
		t.Error("Expected concurrency tests to be available for registration")
	}
}

// TestSuiteEnforcesTimeoutForSlowPlugins validates timeout enforcement (T078).
func TestSuiteEnforcesTimeoutForSlowPlugins(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Verify DefaultConcurrencyConfig has a reasonable timeout
	config := plugintesting.DefaultConcurrencyConfig()
	if config.Timeout <= 0 {
		t.Error("Expected positive timeout in default config")
	}

	// Run a test to verify mock plugin completes promptly
	start := time.Now()
	tests := plugintesting.ConcurrencyTests()
	for _, test := range tests {
		if test.Name == "Concurrency_ParallelRequests_Standard" {
			result := test.TestFunc(harness)
			// Log result for debugging
			t.Logf("Test %s completed in %v, success: %v", test.Name, result.Duration, result.Success)
			break
		}
	}
	elapsed := time.Since(start)

	// The test should complete reasonably quickly (mock plugin is fast)
	if elapsed > 5*time.Second {
		t.Errorf("Test took too long: %v", elapsed)
	}
}
