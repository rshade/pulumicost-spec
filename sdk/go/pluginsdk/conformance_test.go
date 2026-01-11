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

//nolint:testpackage // Testing internal conformance adapter functions with mocks
package pluginsdk

import (
	"bytes"
	"context"
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// conformanceMockPlugin implements the Plugin interface for conformance testing.
// It provides realistic responses that will pass conformance tests.
type conformanceMockPlugin struct {
	name string
}

// Name returns the plugin name.
func (m *conformanceMockPlugin) Name() string { return m.name }

// GetProjectedCost returns a valid projected cost response.
func (m *conformanceMockPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{
		UnitPrice:    0.10,
		Currency:     "USD",
		CostPerMonth: 72.0, // 0.10 * 24 * 30
	}, nil
}

// GetActualCost returns a valid actual cost response.
func (m *conformanceMockPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{
		Results: []*pbc.ActualCostResult{},
	}, nil
}

// GetPricingSpec returns a valid pricing spec response.
func (m *conformanceMockPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{
		Spec: &pbc.PricingSpec{
			Provider:     "aws",
			ResourceType: "ec2",
			BillingMode:  "per_hour",
			RatePerUnit:  0.10,
			Currency:     "USD",
			Unit:         "hour",
		},
	}, nil
}

// EstimateCost returns a valid estimate cost response.
func (m *conformanceMockPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{
		Currency:    "USD",
		CostMonthly: 72.0,
	}, nil
}

// =============================================================================
// Phase 2: Foundational (Nil Plugin Validation) Tests - T006, T007
// =============================================================================

// TestValidatePluginNil verifies that validatePlugin returns an error for nil input.
// This test MUST pass before any adapter functions can safely rely on validation.
func TestValidatePluginNil(t *testing.T) {
	err := validatePlugin(nil)

	require.Error(t, err)
	assert.Equal(t, ErrNilPlugin, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

// TestValidatePluginValid verifies that validatePlugin passes for non-nil plugins.
func TestValidatePluginValid(t *testing.T) {
	plugin := &conformanceMockPlugin{name: "test-plugin"}

	err := validatePlugin(plugin)

	require.NoError(t, err)
}

// =============================================================================
// Phase 3: User Story 1 - RunBasicConformance Tests - T009, T010
// =============================================================================

// TestRunBasicConformanceNilPlugin verifies that RunBasicConformance returns
// an error when passed a nil plugin.
func TestRunBasicConformanceNilPlugin(t *testing.T) {
	result, err := RunBasicConformance(nil)

	require.Error(t, err)
	assert.Equal(t, ErrNilPlugin, err)
	assert.Nil(t, result)
}

// TestRunBasicConformanceValidPlugin verifies that RunBasicConformance returns
// a ConformanceResult when passed a valid plugin.
func TestRunBasicConformanceValidPlugin(t *testing.T) {
	plugin := &conformanceMockPlugin{name: "test-basic-plugin"}

	result, err := RunBasicConformance(plugin)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-basic-plugin", result.PluginName)
	assert.NotZero(t, result.Summary.Total)
}

// =============================================================================
// Phase 4: User Story 2 - RunStandardConformance Tests - T015, T016
// =============================================================================

// TestRunStandardConformanceNilPlugin verifies that RunStandardConformance returns
// an error when passed a nil plugin.
func TestRunStandardConformanceNilPlugin(t *testing.T) {
	result, err := RunStandardConformance(nil)

	require.Error(t, err)
	assert.Equal(t, ErrNilPlugin, err)
	assert.Nil(t, result)
}

// TestRunStandardConformanceValidPlugin verifies that RunStandardConformance returns
// a ConformanceResult when passed a valid plugin.
func TestRunStandardConformanceValidPlugin(t *testing.T) {
	plugin := &conformanceMockPlugin{name: "test-standard-plugin"}

	result, err := RunStandardConformance(plugin)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-standard-plugin", result.PluginName)
	assert.NotZero(t, result.Summary.Total)
	// Verify at least one test was executed
	assert.GreaterOrEqual(t, result.Summary.Total, 1)
}

// =============================================================================
// Phase 5: User Story 3 - RunAdvancedConformance Tests - T021, T022
// =============================================================================

// TestRunAdvancedConformanceNilPlugin verifies that RunAdvancedConformance returns
// an error when passed a nil plugin.
func TestRunAdvancedConformanceNilPlugin(t *testing.T) {
	result, err := RunAdvancedConformance(nil)

	require.Error(t, err)
	assert.Equal(t, ErrNilPlugin, err)
	assert.Nil(t, result)
}

// TestRunAdvancedConformanceValidPlugin verifies that RunAdvancedConformance returns
// a ConformanceResult when passed a valid plugin.
func TestRunAdvancedConformanceValidPlugin(t *testing.T) {
	plugin := &conformanceMockPlugin{name: "test-advanced-plugin"}

	result, err := RunAdvancedConformance(plugin)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-advanced-plugin", result.PluginName)
	assert.NotZero(t, result.Summary.Total)
}

// =============================================================================
// Phase 6: User Story 4 - PrintConformanceReport Tests - T027, T028
// =============================================================================

// TestPrintConformanceReportNilResult verifies that PrintConformanceReport does not
// panic when passed a nil result.
func TestPrintConformanceReportNilResult(t *testing.T) {
	// This should not panic - it should log a warning and return
	assert.NotPanics(t, func() {
		PrintConformanceReport(t, nil)
	})
}

// TestPrintConformanceReportValidResult verifies that PrintConformanceReport outputs
// formatted content with expected sections.
func TestPrintConformanceReportValidResult(t *testing.T) {
	plugin := &conformanceMockPlugin{name: "test-report-plugin"}

	result, err := RunBasicConformance(plugin)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Use PrintConformanceReportTo with a buffer to capture output
	var buf bytes.Buffer
	PrintConformanceReportTo(result, &buf)

	output := buf.String()

	// Verify output contains expected sections
	assert.Contains(t, output, "test-report-plugin", "output should contain plugin name")
	assert.Contains(t, output, "Summary", "output should contain summary section")
	assert.Contains(t, output, "Total", "output should contain total count")
	assert.Contains(t, output, "Passed", "output should contain passed count")
}

// TestPrintConformanceReportToNilResult verifies that PrintConformanceReportTo
// returns immediately when passed a nil result without panicking.
func TestPrintConformanceReportToNilResult(t *testing.T) {
	var buf bytes.Buffer

	assert.NotPanics(t, func() {
		PrintConformanceReportTo(nil, &buf)
	})

	// Buffer should be empty since nil result returns early
	assert.Empty(t, buf.String())
}

// =============================================================================
// Additional Edge Case Tests
// =============================================================================

// TestTypeAliasesExist verifies that type aliases are accessible and work correctly.
func TestTypeAliasesExist(t *testing.T) {
	// Verify ConformanceLevel constants are accessible
	assert.Equal(t, ConformanceLevelBasic, ConformanceLevel(0))
	assert.Equal(t, ConformanceLevelStandard, ConformanceLevel(1))
	assert.Equal(t, ConformanceLevelAdvanced, ConformanceLevel(2))
}

// TestConformanceResultPassed verifies the Passed() method works through the alias.
func TestConformanceResultPassed(t *testing.T) {
	plugin := &conformanceMockPlugin{name: "passed-test-plugin"}

	result, err := RunBasicConformance(plugin)
	require.NoError(t, err)
	require.NotNil(t, result)

	// The Passed() method should be accessible through the type alias
	// (may pass or fail depending on mock implementation, but should not panic)
	_ = result.Passed()
}
