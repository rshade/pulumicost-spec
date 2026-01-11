package testing

import (
	"context"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// MockLegacyPlugin implements the CostSourceServiceServer interface directly
// without PluginInfoProvider for testing legacy plugin behavior.
type MockLegacyPlugin struct {
	pbc.UnimplementedCostSourceServiceServer

	name string
}

// NewMockLegacyPlugin creates a new MockLegacyPlugin.
func NewMockLegacyPlugin(name string) *MockLegacyPlugin {
	return &MockLegacyPlugin{name: name}
}

func (p *MockLegacyPlugin) Name(_ context.Context, _ *pbc.NameRequest) (*pbc.NameResponse, error) {
	return &pbc.NameResponse{Name: p.name}, nil
}

func (p *MockLegacyPlugin) Supports(_ context.Context, _ *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
	return &pbc.SupportsResponse{Supported: true}, nil
}

func (p *MockLegacyPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (p *MockLegacyPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (p *MockLegacyPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (p *MockLegacyPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{}, nil
}

func (p *MockLegacyPlugin) GetRecommendations(
	_ context.Context,
	_ *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	return &pbc.GetRecommendationsResponse{}, nil
}

func (p *MockLegacyPlugin) GetBudgets(
	_ context.Context,
	_ *pbc.GetBudgetsRequest,
) (*pbc.GetBudgetsResponse, error) {
	return &pbc.GetBudgetsResponse{}, nil
}

func (p *MockLegacyPlugin) DismissRecommendation(
	_ context.Context,
	_ *pbc.DismissRecommendationRequest,
) (*pbc.DismissRecommendationResponse, error) {
	return &pbc.DismissRecommendationResponse{}, nil
}
