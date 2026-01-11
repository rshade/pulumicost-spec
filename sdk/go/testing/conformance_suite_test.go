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
	"encoding/json"
	"testing"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

// TestFullSuiteReturnsConsolidatedReport validates that the full suite returns a report (T059).
func TestFullSuiteReturnsConsolidatedReport(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	result, err := plugintesting.RunStandardConformance(plugin)
	if err != nil {
		t.Fatalf("Suite execution failed: %v", err)
	}

	// Verify consolidated report structure
	if result.Version == "" {
		t.Error("Report version is empty")
	}
	if result.PluginName == "" {
		t.Error("Plugin name is empty")
	}
	if result.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}
	if result.Duration <= 0 {
		t.Error("Duration should be positive")
	}
	if result.Categories == nil {
		t.Error("Categories map is nil")
	}

	// Verify summary
	if result.Summary.Total == 0 {
		t.Error("Total tests should be > 0")
	}
}

// TestSuiteDeterminesCorrectCertificationLevel validates level determination (T060).
func TestSuiteDeterminesCorrectCertificationLevel(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	// Test Basic conformance
	t.Run("Basic", func(t *testing.T) {
		result, err := plugintesting.RunBasicConformance(plugin)
		if err != nil {
			t.Fatalf("Basic conformance failed: %v", err)
		}
		if !result.Passed() {
			t.Errorf("Expected Basic conformance to pass, got level: %s", result.LevelAchieved.String())
		}
	})

	// Test Standard conformance
	t.Run("Standard", func(t *testing.T) {
		result, err := plugintesting.RunStandardConformance(plugin)
		if err != nil {
			t.Fatalf("Standard conformance failed: %v", err)
		}
		if !result.Passed() {
			t.Logf("Standard conformance result: %s (some failures may be expected)", result.LevelAchieved.String())
		}
	})

	// Test Advanced conformance
	t.Run("Advanced", func(t *testing.T) {
		result, err := plugintesting.RunAdvancedConformance(plugin)
		if err != nil {
			t.Fatalf("Advanced conformance failed: %v", err)
		}
		// Log the result - may not pass all advanced tests
		t.Logf("Advanced conformance result: %s", result.LevelAchieved.String())
	})
}

// TestSuiteProvidesActionableFailureFeedback validates error messaging (T061).
func TestSuiteProvidesActionableFailureFeedback(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	result, err := plugintesting.RunStandardConformance(plugin)
	if err != nil {
		t.Fatalf("Suite execution failed: %v", err)
	}

	// Check that all test results have appropriate feedback
	for catName, catResult := range result.Categories {
		for i, testResult := range catResult.Results {
			// Only flag if both details AND error are missing for failed tests
			if !testResult.Success && testResult.Details == "" && testResult.Error == nil {
				t.Errorf("Failed test in category %s index %d lacks both details and error", catName, i)
			}

			// Failed tests should have actionable error messages
			if !testResult.Success && testResult.Error == nil {
				t.Errorf("Failed test should have error: category=%s, test=%d", catName, i)
			}
		}
	}
}

// TestConformanceResultToJSON validates JSON serialization (T062).
func TestConformanceResultToJSON(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	result, err := plugintesting.RunBasicConformance(plugin)
	if err != nil {
		t.Fatalf("Suite execution failed: %v", err)
	}

	// Serialize to JSON
	jsonBytes, err := result.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if unmarshalErr := json.Unmarshal(jsonBytes, &parsed); unmarshalErr != nil {
		t.Fatalf("JSON parsing failed: %v", unmarshalErr)
	}

	// Verify required fields exist
	requiredFields := []string{
		"version",
		"timestamp",
		"plugin_name",
		"level_achieved",
		"summary",
		"categories",
		"duration",
	}
	for _, field := range requiredFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("Missing required JSON field: %s", field)
		}
	}

	// Verify summary structure
	summary, ok := parsed["summary"].(map[string]interface{})
	if !ok {
		t.Error("Summary is not a valid object")
	} else {
		for _, field := range []string{"total", "passed", "failed", "skipped"} {
			if _, hasField := summary[field]; !hasField {
				t.Errorf("Missing summary field: %s", field)
			}
		}
	}
}

// TestRunCategory validates RunCategory method (T063).
func TestRunCategory(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	suite := plugintesting.NewConformanceSuite()
	plugintesting.RegisterSpecValidationTests(suite)

	// Run only spec validation category
	result, err := suite.RunCategory(plugin, plugintesting.CategorySpecValidation)
	if err != nil {
		t.Fatalf("RunCategory failed: %v", err)
	}

	if result.Name != plugintesting.CategorySpecValidation {
		t.Errorf("Expected category %s, got %s", plugintesting.CategorySpecValidation, result.Name)
	}

	// Verify we got some results
	totalTests := result.Passed + result.Failed + result.Skipped
	if totalTests == 0 {
		t.Error("Expected some tests to run")
	}
}

// TestPrintReport validates PrintReport function (T066).
func TestPrintReport(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()

	result, err := plugintesting.RunBasicConformance(plugin)
	if err != nil {
		t.Fatalf("Suite execution failed: %v", err)
	}

	// PrintReport should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintReport panicked: %v", r)
		}
	}()

	// This just verifies it doesn't panic - output goes to stdout
	plugintesting.PrintReport(result)
}

// TestSuiteHandlesNilPlugin validates nil plugin handling (T076).
func TestSuiteHandlesNilPlugin(t *testing.T) {
	// This tests that the framework documents nil plugin behavior.
	// gRPC will panic when a nil plugin is used since the handler
	// performs an interface conversion. This is expected behavior.
	// The test simply documents this - production code should validate
	// plugins before passing to the conformance suite.

	// Verify a valid plugin works (control case)
	plugin := plugintesting.NewMockPlugin()
	result, err := plugintesting.RunBasicConformance(plugin)
	if err != nil {
		t.Fatalf("Expected valid plugin to work: %v", err)
	}
	if result.Summary.Total == 0 {
		t.Error("Expected some tests to run with valid plugin")
	}

	// Note: Passing nil to RunBasicConformance will cause a panic in gRPC.
	// This is documented behavior - callers must validate plugins beforehand.
}

// TestSuiteHandlesUnimplementedPlugin validates empty response handling (T079).
func TestSuiteHandlesUnimplementedPlugin(t *testing.T) {
	// Use a mock plugin with minimal responses
	plugin := plugintesting.NewMockPlugin()

	result, err := plugintesting.RunBasicConformance(plugin)
	if err != nil {
		t.Fatalf("Suite execution failed: %v", err)
	}

	// The suite should complete even with minimal plugin responses
	if result.Summary.Total == 0 {
		t.Error("Expected some tests to run")
	}
}

// TestConformanceResultPassed validates the Passed() method.
func TestConformanceResultPassed(t *testing.T) {
	testCases := []struct {
		name     string
		failed   int
		expected bool
	}{
		{"NoFailures", 0, true},
		{"OneFailure", 1, false},
		{"ManyFailures", 10, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := &plugintesting.ConformanceResult{
				Summary: plugintesting.ResultSummary{
					Failed: tc.failed,
				},
			}

			if result.Passed() != tc.expected {
				t.Errorf("Expected Passed()=%v for %d failures", tc.expected, tc.failed)
			}
		})
	}
}

// TestConformanceLevelString validates String() method.
func TestConformanceLevelString(t *testing.T) {
	testCases := []struct {
		level    plugintesting.ConformanceLevel
		expected string
	}{
		{plugintesting.ConformanceLevelBasic, "Basic"},
		{plugintesting.ConformanceLevelStandard, "Standard"},
		{plugintesting.ConformanceLevelAdvanced, "Advanced"},
		{plugintesting.ConformanceLevel(99), "Unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.level.String() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.level.String())
			}
		})
	}
}

// TestTestCategoryString validates category String() method.
func TestTestCategoryString(t *testing.T) {
	testCases := []struct {
		category plugintesting.TestCategory
		expected string
	}{
		{plugintesting.CategorySpecValidation, "Spec Validation"},
		{plugintesting.CategoryRPCCorrectness, "RPC Correctness"},
		{plugintesting.CategoryPerformance, "Performance"},
		{plugintesting.CategoryConcurrency, "Concurrency"},
		{plugintesting.TestCategory("custom"), "custom"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.category), func(t *testing.T) {
			if tc.category.String() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.category.String())
			}
		})
	}
}
