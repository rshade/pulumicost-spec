package pluginsdk //nolint:testpackage // White-box testing required

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// ErrorTestPlugin implements Plugin and PluginInfoProvider interfaces for testing error conditions.
type ErrorTestPlugin struct {
	returnNil    bool
	incomplete   bool
	invalidSpec  bool
	handlerError bool
}

func (p *ErrorTestPlugin) Name() string {
	return "error-test-plugin"
}

func (p *ErrorTestPlugin) GetPluginInfo(
	_ context.Context,
	_ *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error) {
	if p.handlerError {
		return nil, &testError{message: "handler error"}
	}
	if p.returnNil {
		//nolint:nilnil // Intentional nil return for testing
		return nil, nil
	}
	if p.incomplete {
		return &pbc.GetPluginInfoResponse{
			Name:        "", // Empty name
			Version:     "1.0.0",
			SpecVersion: "v1.2.0",
		}, nil
	}
	if p.invalidSpec {
		return &pbc.GetPluginInfoResponse{
			Name:        "test-plugin",
			Version:     "1.0.0",
			SpecVersion: "invalid-format", // Invalid spec version
		}, nil
	}
	return &pbc.GetPluginInfoResponse{
		Name:        "error-test-plugin",
		Version:     "1.0.0",
		SpecVersion: "v1.2.0",
	}, nil
}

func (p *ErrorTestPlugin) Supports(
	_ context.Context,
	_ *pbc.SupportsRequest,
) (*pbc.SupportsResponse, error) {
	return &pbc.SupportsResponse{Supported: true}, nil
}

func (p *ErrorTestPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (p *ErrorTestPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (p *ErrorTestPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (p *ErrorTestPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{}, nil
}

func (p *ErrorTestPlugin) GetRecommendations(
	_ context.Context,
	_ *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	return &pbc.GetRecommendationsResponse{}, nil
}

func (p *ErrorTestPlugin) GetBudgets(
	_ context.Context,
	_ *pbc.GetBudgetsRequest,
) (*pbc.GetBudgetsResponse, error) {
	return &pbc.GetBudgetsResponse{}, nil
}

func (p *ErrorTestPlugin) DismissRecommendation(
	_ context.Context,
	_ *pbc.DismissRecommendationRequest,
) (*pbc.DismissRecommendationResponse, error) {
	return &pbc.DismissRecommendationResponse{}, nil
}

// testError implements error for testing.
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

// CreateErrorTestServer creates a test server with error test plugin.
func CreateErrorTestServer(plugin *ErrorTestPlugin) *Server {
	return NewServer(plugin)
}

func TestGetPluginInfoError_NilResponse(t *testing.T) {
	plugin := &ErrorTestPlugin{returnNil: true}
	server := CreateErrorTestServer(plugin)

	_, err := server.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "unable to retrieve plugin metadata")
}
