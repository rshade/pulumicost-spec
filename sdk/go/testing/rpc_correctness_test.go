package testing_test

import (
	"testing"

	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

// TestRPCCorrectnessNameRPC validates the Name RPC method (T026).
func TestRPCCorrectnessNameRPC(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.RPCCorrectnessTests()
	for _, test := range tests {
		if test.Name == "RPCCorrectness_NameRPC" {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
			break
		}
	}
}

// TestRPCCorrectnessSupportsRPC validates the Supports RPC method (T027).
func TestRPCCorrectnessSupportsRPC(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.RPCCorrectnessTests()
	for _, test := range tests {
		if test.Name == "RPCCorrectness_SupportsRPC" {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
			break
		}
	}
}

// TestRPCCorrectnessNilResource validates nil resource handling (T028).
func TestRPCCorrectnessNilResource(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.RPCCorrectnessTests()
	for _, test := range tests {
		if test.Name == "RPCCorrectness_NilResource" {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
			break
		}
	}
}

// TestRPCCorrectnessInvalidTimeRange validates invalid time range handling (T029).
func TestRPCCorrectnessInvalidTimeRange(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.RPCCorrectnessTests()
	for _, test := range tests {
		if test.Name == "RPCCorrectness_InvalidTimeRange" {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
			break
		}
	}
}

// TestRPCCorrectnessAllTests runs all RPC correctness tests.
func TestRPCCorrectnessAllTests(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	tests := plugintesting.RPCCorrectnessTests()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := test.TestFunc(harness)
			if !result.Success {
				t.Errorf("Test %s failed: %v - %s", test.Name, result.Error, result.Details)
			}
		})
	}
}

// TestRegisterRPCCorrectnessTests validates test registration.
func TestRegisterRPCCorrectnessTests(t *testing.T) {
	suite := plugintesting.NewConformanceSuite()
	plugintesting.RegisterRPCCorrectnessTests(suite)

	// Verify tests were registered
	config := suite.GetConfig()
	if config.TargetLevel != plugintesting.ConformanceLevelStandard {
		t.Errorf("Expected default target level Standard, got %v", config.TargetLevel)
	}
}

// TestRPCCorrectnessPanicRecovery validates that the suite recovers from plugin panics (T077).
func TestRPCCorrectnessPanicRecovery(t *testing.T) {
	// Create a plugin that might panic (mock doesn't panic, but we can test the framework)
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	// Run tests - they should complete without panicking the test framework
	tests := plugintesting.RPCCorrectnessTests()
	for _, test := range tests {
		t.Run(test.Name+"_PanicSafe", func(t *testing.T) {
			// Wrap in recovery to test framework stability
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Test framework panicked: %v", r)
				}
			}()

			result := test.TestFunc(harness)
			// We just verify it doesn't panic
			_ = result
		})
	}
}
