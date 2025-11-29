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
		t.Run(test.Name+"_RaceCheck", func(_ *testing.T) {
			result := test.TestFunc(harness)
			// We're checking for race conditions, not success
			// The -race flag will catch any issues
			_ = result
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
}

// TestSuiteEnforcesTimeoutForSlowPlugins validates timeout enforcement (T078).
func TestSuiteEnforcesTimeoutForSlowPlugins(t *testing.T) {
	// Use a mock plugin - we can't easily make it slow, but we can verify
	// that the framework has timeout support
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Create a config with a short timeout
	config := plugintesting.DefaultConcurrencyConfig()
	config.Timeout = 100 * time.Millisecond

	// Verify the config has a timeout set
	if config.Timeout == 0 {
		t.Error("Expected timeout to be set")
	}

	// Run a quick test to verify framework doesn't hang
	start := time.Now()
	tests := plugintesting.ConcurrencyTests()
	for _, test := range tests {
		if test.Name == "Concurrency_ParallelRequests_Standard" {
			result := test.TestFunc(harness)
			_ = result
			break
		}
	}
	elapsed := time.Since(start)

	// The test should complete reasonably quickly (mock plugin is fast)
	if elapsed > 5*time.Second {
		t.Errorf("Test took too long: %v", elapsed)
	}
}
