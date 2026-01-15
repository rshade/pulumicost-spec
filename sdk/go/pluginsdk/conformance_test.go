//nolint:testpackage // Testing internal conformance adapter functions with mocks
package pluginsdk

import (
	"bytes"
	"context"
	"testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
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

// TestGetPluginInfoCapabilitiesDiscovery verifies that GetPluginInfo auto-discovers
// capabilities based on implemented interfaces.
func TestGetPluginInfoCapabilitiesDiscovery(t *testing.T) {
	// Create a test plugin that implements specific interfaces
	testPlugin := &capabilityTestPlugin{
		name: "capability-test-plugin",
	}

	// Create server with basic plugin info to enable GetPluginInfo
	pluginInfo := NewPluginInfo("capability-test-plugin", "v1.0.0")
	server := NewServerWithOptions(testPlugin, nil, nil, pluginInfo)

	ctx := context.Background()
	req := &pbc.GetPluginInfoRequest{}

	resp, err := server.GetPluginInfo(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify auto-discovered capabilities - includes all optional interfaces + base interface capabilities
	expectedCapabilities := []pbc.PluginCapability{
		pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_PRICING_SPEC,
		pbc.PluginCapability_PLUGIN_CAPABILITY_ESTIMATE_COST,
		pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS,
		pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
	}

	assert.ElementsMatch(t, expectedCapabilities, resp.GetCapabilities(),
		"GetPluginInfo should auto-discover capabilities based on implemented interfaces")
}

func TestGetPluginInfoCapabilitiesOverride(t *testing.T) {
	testPlugin := &capabilityTestPlugin{
		name: "override-test-plugin",
	}

	pluginInfo := NewPluginInfo("override-test-plugin", "v1.0.0",
		WithCapabilities(pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS))
	server := NewServerWithOptions(testPlugin, nil, nil, pluginInfo)

	ctx := context.Background()
	req := &pbc.GetPluginInfoRequest{}

	resp, err := server.GetPluginInfo(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.ElementsMatch(t,
		[]pbc.PluginCapability{pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS},
		resp.GetCapabilities(),
		"GetPluginInfo should honor explicit capabilities override")
}

// capabilityTestPlugin implements multiple interfaces for testing capability discovery.
type capabilityTestPlugin struct {
	name string
}

func (p *capabilityTestPlugin) Name() string { return p.name }

// Implements Plugin interface.
func (p *capabilityTestPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (p *capabilityTestPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (p *capabilityTestPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (p *capabilityTestPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{}, nil
}

// Implements RecommendationsProvider.
func (p *capabilityTestPlugin) GetRecommendations(
	_ context.Context,
	_ *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	return &pbc.GetRecommendationsResponse{}, nil
}

// Implements BudgetsProvider.
func (p *capabilityTestPlugin) GetBudgets(
	_ context.Context,
	_ *pbc.GetBudgetsRequest,
) (*pbc.GetBudgetsResponse, error) {
	return &pbc.GetBudgetsResponse{}, nil
}

// Implements DismissProvider.
func (p *capabilityTestPlugin) DismissRecommendation(
	_ context.Context,
	_ *pbc.DismissRecommendationRequest,
) (*pbc.DismissRecommendationResponse, error) {
	return &pbc.DismissRecommendationResponse{}, nil
}

// Implements DryRunHandler.
func (p *capabilityTestPlugin) HandleDryRun(
	_ *pbc.DryRunRequest,
) (*pbc.DryRunResponse, error) {
	return &pbc.DryRunResponse{}, nil
}
