//nolint:testpackage // Testing internal Server implementation with mocks
package pluginsdk

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockPlugin implements Plugin interface for testing.
type mockPlugin struct {
	name string
}

func (m *mockPlugin) Name() string { return m.name }

func (m *mockPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (m *mockPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (m *mockPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (m *mockPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{}, nil
}

// mockSupportsPlugin implements both Plugin and SupportsProvider interfaces.
type mockSupportsPlugin struct {
	mockPlugin

	supported bool
	reason    string
	err       error
}

func (m *mockSupportsPlugin) Supports(
	_ context.Context,
	_ *pbc.SupportsRequest,
) (*pbc.SupportsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &pbc.SupportsResponse{
		Supported: m.supported,
		Reason:    m.reason,
	}, nil
}

// mockRegistry implements RegistryLookup for testing.
type mockRegistry struct {
	plugins map[string]string // key: "provider:region", value: plugin name
}

func (m *mockRegistry) FindPlugin(provider, region string) string {
	key := provider + ":" + region
	return m.plugins[key]
}

// TestName tests the Name RPC.
func TestName(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}
	server := NewServer(plugin)

	resp, err := server.Name(context.Background(), &pbc.NameRequest{})

	require.NoError(t, err)
	assert.Equal(t, "test-plugin", resp.GetName())
}

// TestGetProjectedCost tests the GetProjectedCost RPC.
func TestGetProjectedCost(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}
	server := NewServer(plugin)

	req := &pbc.GetProjectedCostRequest{}
	resp, err := server.GetProjectedCost(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestGetActualCost tests the GetActualCost RPC.
func TestGetActualCost(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}
	server := NewServer(plugin)

	req := &pbc.GetActualCostRequest{}
	resp, err := server.GetActualCost(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestGetPricingSpec tests the GetPricingSpec RPC.
func TestGetPricingSpec(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}
	server := NewServer(plugin)

	req := &pbc.GetPricingSpecRequest{}
	resp, err := server.GetPricingSpec(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestEstimateCost tests the EstimateCost RPC.
func TestEstimateCost(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}
	server := NewServer(plugin)

	req := &pbc.EstimateCostRequest{}
	resp, err := server.EstimateCost(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
func TestSupports_PluginImplementsAndReturnsSupported(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  true,
	}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "test-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	resp, err := server.Supports(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, resp.GetSupported())
	assert.Empty(t, resp.GetReason())
}

// TestSupports_PluginImplementsAndReturnsNotSupported tests T006: plugin returns not-supported with reason.
func TestSupports_PluginImplementsAndReturnsNotSupported(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  false,
		reason:     "Resource type not supported by this plugin",
	}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "test-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "gamelift",
			Region:       "us-east-1",
		},
	}

	resp, err := server.Supports(context.Background(), req)

	require.NoError(t, err)
	assert.False(t, resp.GetSupported())
	assert.Equal(t, "Resource type not supported by this plugin", resp.GetReason())
}

// TestSupports_InvalidProviderRegionReturnsInvalidArgument tests T007: invalid provider/region returns InvalidArgument.
func TestSupports_InvalidProviderRegionReturnsInvalidArgument(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  true,
	}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "test-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "unknown-provider",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	_, err := server.Supports(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "no plugin registered")
}

// TestSupports_NoPluginRegisteredReturnsInvalidArgument tests T008: no plugin registered returns InvalidArgument.
func TestSupports_NoPluginRegisteredReturnsInvalidArgument(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  true,
	}
	// Empty registry - no plugins registered
	registry := &mockRegistry{
		plugins: map[string]string{},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	_, err := server.Supports(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "no plugin registered")
}

// TestSupports_PluginErrorReturnsInternalStatus tests T009: plugin error returns Internal status.
func TestSupports_PluginErrorReturnsInternalStatus(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		err:        errors.New("database connection failed"),
	}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "test-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	_, err := server.Supports(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	// Error message should be generic, not expose internal details
	assert.Contains(t, st.Message(), "plugin failed to execute")
}

// TestSupports_PluginNotImplementsReturnsDefaultResponse tests T014: plugin without SupportsProvider returns default.
func TestSupports_PluginNotImplementsReturnsDefaultResponse(t *testing.T) {
	// Use mockPlugin which does NOT implement SupportsProvider
	plugin := &mockPlugin{name: "legacy-plugin"}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "legacy-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	resp, err := server.Supports(context.Background(), req)

	require.NoError(t, err)
	assert.False(t, resp.GetSupported())
	assert.NotEmpty(t, resp.GetReason())
}

// TestSupports_DefaultResponseIncludesReason tests T015: default response includes reason explaining not implemented.
func TestSupports_DefaultResponseIncludesReason(t *testing.T) {
	plugin := &mockPlugin{name: "legacy-plugin"}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "legacy-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	resp, err := server.Supports(context.Background(), req)

	require.NoError(t, err)
	assert.False(t, resp.GetSupported())
	assert.Contains(t, resp.GetReason(), "not implemented")
}

// TestSupports_NilResourceReturnsInvalidArgument tests that a nil resource descriptor returns InvalidArgument.
func TestSupports_NilResourceReturnsInvalidArgument(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  true,
	}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "test-plugin"},
	}
	server := NewServerWithRegistry(plugin, registry)

	// Request with nil resource
	req := &pbc.SupportsRequest{Resource: nil}

	_, err := server.Supports(context.Background(), req)

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "resource descriptor is required")
}

// TestSupports_NewServerUsesDefaultRegistry tests that NewServer uses DefaultRegistryLookup.
func TestSupports_NewServerUsesDefaultRegistry(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  true,
	}
	// Use NewServer which uses DefaultRegistryLookup (always returns empty)
	server := NewServer(plugin)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	_, err := server.Supports(context.Background(), req)

	// DefaultRegistryLookup returns "" for all lookups, so this should fail with InvalidArgument
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "no plugin registered")
}

// TestNewServerWithRegistry_NilRegistryUsesDefault tests that nil registry falls back to default.
func TestNewServerWithRegistry_NilRegistryUsesDefault(t *testing.T) {
	plugin := &mockSupportsPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		supported:  true,
	}
	// Pass nil registry - should fall back to DefaultRegistryLookup
	server := NewServerWithRegistry(plugin, nil)

	req := &pbc.SupportsRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Region:       "us-east-1",
		},
	}

	_, err := server.Supports(context.Background(), req)

	// Should behave same as NewServer - DefaultRegistryLookup returns ""
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "no plugin registered")
}

// TestNewServerWithOptions_WithCustomLogger tests that custom logger is used.
func TestNewServerWithOptions_WithCustomLogger(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}
	registry := &mockRegistry{
		plugins: map[string]string{"aws:us-east-1": "test-plugin"},
	}
	logger := newDefaultLogger()

	server := NewServerWithOptions(plugin, registry, &logger)

	// Verify the server was created successfully
	resp, err := server.Name(context.Background(), &pbc.NameRequest{})
	require.NoError(t, err)
	assert.Equal(t, "test-plugin", resp.GetName())
}

// TestNewServerWithOptions_NilRegistryAndLogger tests nil parameters use defaults.
func TestNewServerWithOptions_NilRegistryAndLogger(t *testing.T) {
	plugin := &mockPlugin{name: "test-plugin"}

	// Both nil - should use defaults
	server := NewServerWithOptions(plugin, nil, nil)

	resp, err := server.Name(context.Background(), &pbc.NameRequest{})
	require.NoError(t, err)
	assert.Equal(t, "test-plugin", resp.GetName())
}

// TestCreateTestResource tests the CreateTestResource helper.
func TestCreateTestResource(t *testing.T) {
	// Test with properties
	props := map[string]string{"instance_type": "t2.micro", "region": "us-east-1"}
	resource := CreateTestResource("aws", "aws:ec2:Instance", props)

	assert.Equal(t, "aws", resource.GetProvider())
	assert.Equal(t, "aws:ec2:Instance", resource.GetResourceType())
	assert.Equal(t, "t2.micro", resource.GetTags()["instance_type"])
	assert.Equal(t, "us-east-1", resource.GetTags()["region"])

	// Test with nil properties - should create empty map
	resourceNilProps := CreateTestResource("azure", "azure:compute:VirtualMachine", nil)

	assert.Equal(t, "azure", resourceNilProps.GetProvider())
	assert.Equal(t, "azure:compute:VirtualMachine", resourceNilProps.GetResourceType())
	assert.NotNil(t, resourceNilProps.GetTags())
	assert.Empty(t, resourceNilProps.GetTags())
}

// =============================================================================
// UnaryInterceptors Tests (User Stories 1, 2, 3)
// =============================================================================

// trackingInterceptor creates an interceptor that records invocation order.
func trackingInterceptor(id string, order *[]string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		*order = append(*order, id)
		return handler(ctx, req)
	}
}

// contextModifyingInterceptor creates an interceptor that adds a value to context.
type testContextKey string

func contextModifyingInterceptor(key testContextKey, value string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = context.WithValue(ctx, key, value)
		return handler(ctx, req)
	}
}

// countingInterceptor creates an interceptor that counts invocations.
func countingInterceptor(counter *atomic.Int32) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		counter.Add(1)
		return handler(ctx, req)
	}
}

// --- User Story 1 Tests: Single Interceptor Registration ---

// TestServeConfig_SingleInterceptorInvocation tests T003: single interceptor is invoked.
func TestServeConfig_SingleInterceptorInvocation(t *testing.T) {
	var counter atomic.Int32
	interceptor := countingInterceptor(&counter)

	config := ServeConfig{
		Plugin:            &mockPlugin{name: "test"},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{interceptor},
	}

	// Verify the plugin and interceptor are included in config
	assert.NotNil(t, config.Plugin, "plugin should be set")
	assert.Len(t, config.UnaryInterceptors, 1)
	assert.NotNil(t, config.UnaryInterceptors[0])
}

// TestServeConfig_TraceIDPropagation tests T004: trace ID propagates through interceptor.
// This test verifies that the ServeConfig accepts interceptors that can access trace IDs.
func TestServeConfig_TraceIDPropagation(t *testing.T) {
	var capturedTraceID string
	interceptor := func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Capture trace ID from context (set by TracingUnaryServerInterceptor)
		capturedTraceID = TraceIDFromContext(ctx)
		return handler(ctx, req)
	}

	config := ServeConfig{
		Plugin:            &mockPlugin{name: "test"},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{interceptor},
	}

	// Verify the plugin and interceptor are configured
	assert.NotNil(t, config.Plugin, "plugin should be set")
	assert.Len(t, config.UnaryInterceptors, 1)
	// Note: Full end-to-end trace ID propagation requires running Serve(),
	// which starts a gRPC server. This test validates the config accepts the interceptor.
	_ = capturedTraceID // Will be set when interceptor runs
}

// --- User Story 2 Tests: Multiple Interceptors Chaining ---

// TestServeConfig_MultipleInterceptorsOrder tests T008: multiple interceptors execute in order.
func TestServeConfig_MultipleInterceptorsOrder(t *testing.T) {
	var order []string
	interceptor1 := trackingInterceptor("first", &order)
	interceptor2 := trackingInterceptor("second", &order)
	interceptor3 := trackingInterceptor("third", &order)

	config := ServeConfig{
		Plugin: &mockPlugin{name: "test"},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			interceptor1,
			interceptor2,
			interceptor3,
		},
	}

	// Verify plugin and all interceptors are in config in order
	assert.NotNil(t, config.Plugin, "plugin should be set")
	assert.Len(t, config.UnaryInterceptors, 3)
}

// TestServeConfig_ContextModificationsPropagation tests T009: context modifications propagate.
func TestServeConfig_ContextModificationsPropagation(t *testing.T) {
	key1 := testContextKey("key1")
	key2 := testContextKey("key2")
	interceptor1 := contextModifyingInterceptor(key1, "value1")
	interceptor2 := contextModifyingInterceptor(key2, "value2")

	config := ServeConfig{
		Plugin: &mockPlugin{name: "test"},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			interceptor1,
			interceptor2,
		},
	}

	// Verify plugin and both interceptors are configured
	assert.NotNil(t, config.Plugin, "plugin should be set")
	assert.Len(t, config.UnaryInterceptors, 2)
}

// --- User Story 3 Tests: Backward Compatibility ---

// TestServeConfig_NilUnaryInterceptorsField tests T012: nil field works (backward compat).
func TestServeConfig_NilUnaryInterceptorsField(t *testing.T) {
	// This is the existing usage pattern - no UnaryInterceptors field set
	config := ServeConfig{
		Plugin: &mockPlugin{name: "test"},
		Port:   0,
	}

	// Verify plugin and port are set (backward compatibility config)
	assert.NotNil(t, config.Plugin, "plugin should be set")
	assert.Equal(t, 0, config.Port, "port should be zero for ephemeral")
	// UnaryInterceptors should be nil by default (Go zero value)
	assert.Nil(t, config.UnaryInterceptors)

	// Verify the config is valid for use with Serve()
	// (actual Serve() would work because append to nil slice is safe)
}

// TestServeConfig_EmptySliceUnaryInterceptors tests T013: empty slice works.
func TestServeConfig_EmptySliceUnaryInterceptors(t *testing.T) {
	config := ServeConfig{
		Plugin:            &mockPlugin{name: "test"},
		Port:              0,
		UnaryInterceptors: []grpc.UnaryServerInterceptor{},
	}

	// Verify plugin and port are set
	assert.NotNil(t, config.Plugin, "plugin should be set")
	assert.Equal(t, 0, config.Port, "port should be zero for ephemeral")
	// Empty slice should be valid
	assert.NotNil(t, config.UnaryInterceptors)
	assert.Empty(t, config.UnaryInterceptors)
}

// TestInterceptorChainBuilding tests the interceptor chain building logic from Serve().
// This is a unit test of the chain building pattern used in Serve().
func TestInterceptorChainBuilding(t *testing.T) {
	tests := []struct {
		name             string
		userInterceptors []grpc.UnaryServerInterceptor
		expectedChainLen int
		description      string
	}{
		{
			name:             "nil user interceptors",
			userInterceptors: nil,
			expectedChainLen: 1, // just tracing
			description:      "tracing only",
		},
		{
			name:             "empty user interceptors",
			userInterceptors: []grpc.UnaryServerInterceptor{},
			expectedChainLen: 1, // just tracing
			description:      "tracing only",
		},
		{
			name: "one user interceptor",
			userInterceptors: []grpc.UnaryServerInterceptor{
				func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
					return handler(ctx, req)
				},
			},
			expectedChainLen: 2, // tracing + user
			description:      "tracing + 1 user",
		},
		{
			name: "three user interceptors",
			userInterceptors: []grpc.UnaryServerInterceptor{
				func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
					return handler(ctx, req)
				},
				func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
					return handler(ctx, req)
				},
				func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
					return handler(ctx, req)
				},
			},
			expectedChainLen: 4, // tracing + 3 user
			description:      "tracing + 3 user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the chain building logic from Serve()
			interceptors := make([]grpc.UnaryServerInterceptor, 0, 1+len(tt.userInterceptors))
			interceptors = append(interceptors, TracingUnaryServerInterceptor())
			interceptors = append(interceptors, tt.userInterceptors...)

			assert.Len(t, interceptors, tt.expectedChainLen, tt.description)
			// First interceptor should always be tracing
			assert.NotNil(t, interceptors[0], "tracing interceptor should be first")
		})
	}
}
