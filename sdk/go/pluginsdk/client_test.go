package pluginsdk_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// waitForServer polls the server until it responds or times out.
func waitForServer(t *testing.T, addr string) {
	t.Helper()
	for range 50 { // 500ms total
		_, err := http.Get("http://" + addr + "/healthz")
		if err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("server did not start in time")
}

// clientTestPlugin implements the pluginsdk.Plugin interface for client testing.
type clientTestPlugin struct {
	name string
}

func (m *clientTestPlugin) Name() string { return m.name }

func (m *clientTestPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (m *clientTestPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (m *clientTestPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (m *clientTestPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{
		CostMonthly: 50.0,
		Currency:    "USD",
	}, nil
}

func TestClient_Name(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &clientTestPlugin{name: "client-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin:   plugin,
			Listener: listener,
			Web: pluginsdk.WebConfig{
				Enabled:              true,
				EnableHealthEndpoint: true,
			},
		})
	}()

	// Wait for server to be ready
	addr := listener.Addr().String()
	waitForServer(t, addr)

	// Create client and test
	client := pluginsdk.NewConnectClient("http://" + addr)

	name, err := client.Name(ctx)
	require.NoError(t, err)
	assert.Equal(t, "client-test-plugin", name)

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestClient_EstimateCost(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &clientTestPlugin{name: "estimate-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin:   plugin,
			Listener: listener,
			Web: pluginsdk.WebConfig{
				Enabled:              true,
				EnableHealthEndpoint: true,
			},
		})
	}()

	// Wait for server to be ready
	addr := listener.Addr().String()
	waitForServer(t, addr)

	// Create client and test
	client := pluginsdk.NewConnectClient("http://" + addr)

	resp, err := client.EstimateCost(ctx, &pbc.EstimateCostRequest{
		ResourceType: "aws:ec2/instance:Instance",
	})
	require.NoError(t, err)
	assert.InDelta(t, 50.0, resp.GetCostMonthly(), 0.01)
	assert.Equal(t, "USD", resp.GetCurrency())

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestClient_Supports(t *testing.T) {
	// Start server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &clientTestPlugin{name: "supports-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin:   plugin,
			Listener: listener,
			Web: pluginsdk.WebConfig{
				Enabled:              true,
				EnableHealthEndpoint: true,
			},
		})
	}()

	// Wait for server to be ready
	addr := listener.Addr().String()
	waitForServer(t, addr)

	// Create client and test
	client := pluginsdk.NewConnectClient("http://" + addr)

	// The Supports method requires a registry lookup, which will fail for this test.
	// This test verifies the client can make the call and properly receive error responses.
	_, err = client.Supports(ctx, &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "aws:ec2/instance:Instance",
		Region:       "us-east-1",
	})
	// Error expected because no plugin is registered in the registry
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no plugin registered")

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestClient_Protocol(t *testing.T) {
	t.Run("ProtocolConnect creates connect client", func(t *testing.T) {
		cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
		cfg.Protocol = pluginsdk.ProtocolConnect
		client := pluginsdk.NewClient(cfg)
		assert.NotNil(t, client)
		assert.NotNil(t, client.Inner())
	})

	t.Run("ProtocolGRPC creates gRPC client", func(t *testing.T) {
		cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
		cfg.Protocol = pluginsdk.ProtocolGRPC
		client := pluginsdk.NewClient(cfg)
		assert.NotNil(t, client)
	})

	t.Run("ProtocolGRPCWeb creates gRPC-Web client", func(t *testing.T) {
		cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
		cfg.Protocol = pluginsdk.ProtocolGRPCWeb
		client := pluginsdk.NewClient(cfg)
		assert.NotNil(t, client)
	})

	t.Run("invalid protocol defaults to Connect", func(t *testing.T) {
		cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
		cfg.Protocol = pluginsdk.Protocol(999) // Invalid protocol value
		client := pluginsdk.NewClient(cfg)
		assert.NotNil(t, client)
		// Client should be created successfully with Connect protocol as fallback
		assert.NotNil(t, client.Inner())
	})
}

func TestDefaultClientConfig(t *testing.T) {
	cfg := pluginsdk.DefaultClientConfig("http://example.com")

	assert.Equal(t, "http://example.com", cfg.BaseURL)
	assert.Equal(t, pluginsdk.ProtocolConnect, cfg.Protocol)
	assert.NotNil(t, cfg.HTTPClient)
}

func TestHighThroughputClientConfig(t *testing.T) {
	cfg := pluginsdk.HighThroughputClientConfig("http://example.com")

	assert.Equal(t, "http://example.com", cfg.BaseURL)
	assert.Equal(t, pluginsdk.ProtocolConnect, cfg.Protocol)
	assert.NotNil(t, cfg.HTTPClient)
	// Verify transport is configured (should have custom transport, not nil)
	assert.NotNil(t, cfg.HTTPClient.Transport)
}

func TestNewConnectClient(t *testing.T) {
	client := pluginsdk.NewConnectClient("http://localhost:8080")
	assert.NotNil(t, client)
}

func TestNewGRPCClient(t *testing.T) {
	client := pluginsdk.NewGRPCClient("http://localhost:8080")
	assert.NotNil(t, client)
}

func TestNewGRPCWebClient(t *testing.T) {
	client := pluginsdk.NewGRPCWebClient("http://localhost:8080")
	assert.NotNil(t, client)
}

func TestClient_Supports_NilResource(t *testing.T) {
	client := pluginsdk.NewConnectClient("http://localhost:8080")

	// Calling Supports with nil resource should return an error
	_, err := client.Supports(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resource descriptor cannot be nil")
}

func TestClient_Supports_EmptyResourceType(t *testing.T) {
	client := pluginsdk.NewConnectClient("http://localhost:8080")

	// Calling Supports with empty resource type should return an error
	_, err := client.Supports(context.Background(), &pbc.ResourceDescriptor{
		ResourceType: "",
		Provider:     "aws",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resource type is required")
}
