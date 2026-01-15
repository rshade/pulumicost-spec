package pluginsdk_test

import (
	"context"
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// BenchmarkInferCapabilities benchmarks the capability inference logic
// by measuring GetPluginInfo performance with auto-discovery enabled.
func BenchmarkInferCapabilities(b *testing.B) {
	// Create a plugin that implements multiple interfaces for realistic benchmarking
	plugin := &benchmarkCapabilityPlugin{
		name: "benchmark-plugin",
	}

	// Create server with minimal plugin info to enable GetPluginInfo with auto-discovery
	pluginInfo := pluginsdk.NewPluginInfo("benchmark-plugin", "v1.0.0")
	server := pluginsdk.NewServerWithOptions(plugin, nil, nil, pluginInfo)

	ctx := context.Background()
	req := &pbc.GetPluginInfoRequest{}

	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		resp, err := server.GetPluginInfo(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
		if len(resp.GetCapabilities()) != 4 {
			b.Fatalf("Expected 4 capabilities, got %d", len(resp.GetCapabilities()))
		}
	}
}

// BenchmarkInferCapabilitiesExplicit benchmarks with explicitly set capabilities
// to compare performance difference.
func BenchmarkInferCapabilitiesExplicit(b *testing.B) {
	plugin := &benchmarkCapabilityPlugin{
		name: "benchmark-plugin",
	}

	// Create server with explicit capabilities (no auto-discovery)
	pluginInfo := pluginsdk.NewPluginInfo("benchmark-plugin", "v1.0.0",
		pluginsdk.WithCapabilities(
			pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS,
			pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS,
			pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS,
			pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS,
		))
	server := pluginsdk.NewServerWithOptions(plugin, nil, nil, pluginInfo)

	ctx := context.Background()
	req := &pbc.GetPluginInfoRequest{}

	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		resp, err := server.GetPluginInfo(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
		if len(resp.GetCapabilities()) != 4 {
			b.Fatalf("Expected 4 capabilities, got %d", len(resp.GetCapabilities()))
		}
	}
}

// benchmarkCapabilityPlugin implements multiple interfaces for benchmarking.
type benchmarkCapabilityPlugin struct {
	name string
}

func (p *benchmarkCapabilityPlugin) Name() string { return p.name }

// Implements Plugin interface.
func (p *benchmarkCapabilityPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (p *benchmarkCapabilityPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (p *benchmarkCapabilityPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (p *benchmarkCapabilityPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{}, nil
}

// Implements RecommendationsProvider.
func (p *benchmarkCapabilityPlugin) GetRecommendations(
	_ context.Context,
	_ *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	return &pbc.GetRecommendationsResponse{}, nil
}

// Implements BudgetsProvider.
func (p *benchmarkCapabilityPlugin) GetBudgets(
	_ context.Context,
	_ *pbc.GetBudgetsRequest,
) (*pbc.GetBudgetsResponse, error) {
	return &pbc.GetBudgetsResponse{}, nil
}
