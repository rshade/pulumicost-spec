package pluginsdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1/pbcconnect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// connectTestPlugin implements the pluginsdk.Plugin interface for connect testing.
type connectTestPlugin struct {
	name string
}

func (m *connectTestPlugin) Name() string { return m.name }

func (m *connectTestPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (m *connectTestPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (m *connectTestPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (m *connectTestPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{
		CostMonthly: 100.0,
		Currency:    "USD",
	}, nil
}

func TestNewConnectHandler_NilServer_Panics(t *testing.T) {
	assert.PanicsWithValue(t, "NewConnectHandler: server cannot be nil", func() {
		pluginsdk.NewConnectHandler(nil)
	})
}

func TestConnectHandler_Name(t *testing.T) {
	plugin := &connectTestPlugin{name: "test-plugin"}
	server := pluginsdk.NewServer(plugin)
	handler := pluginsdk.NewConnectHandler(server)

	ctx := context.Background()
	req := connect.NewRequest(&pbc.NameRequest{})

	resp, err := handler.Name(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, "test-plugin", resp.Msg.GetName())
}

func TestConnectHandler_EstimateCost(t *testing.T) {
	plugin := &connectTestPlugin{name: "test-plugin"}
	server := pluginsdk.NewServer(plugin)
	handler := pluginsdk.NewConnectHandler(server)

	ctx := context.Background()
	req := connect.NewRequest(&pbc.EstimateCostRequest{
		ResourceType: "aws:ec2/instance:Instance",
	})

	resp, err := handler.EstimateCost(ctx, req)
	require.NoError(t, err)
	assert.InDelta(t, 100.0, resp.Msg.GetCostMonthly(), 0.01)
	assert.Equal(t, "USD", resp.Msg.GetCurrency())
}

func TestServeConnect_WithHealthEndpoint(t *testing.T) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &connectTestPlugin{name: "test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in background
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
	resp, err := http.Get("http://" + addr + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "ok", string(body))

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServeConnect_ConnectProtocol(t *testing.T) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &connectTestPlugin{name: "connect-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in background
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
	reqBody := []byte(`{}`)
	resp, err := http.Post(
		"http://"+addr+"/pulumicost.v1.CostSourceService/Name",
		"application/json",
		bytes.NewReader(reqBody),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "connect-test-plugin", result["name"])

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServeConnect_ConnectClient(t *testing.T) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &connectTestPlugin{name: "client-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server in background
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

	// Create Connect client
	client := pbcconnect.NewCostSourceServiceClient(
		http.DefaultClient,
		"http://"+addr,
	)

	// Test Name RPC
	nameResp, err := client.Name(ctx, connect.NewRequest(&pbc.NameRequest{}))
	require.NoError(t, err)
	assert.Equal(t, "client-test-plugin", nameResp.Msg.GetName())

	// Test EstimateCost RPC
	estimateResp, err := client.EstimateCost(ctx, connect.NewRequest(&pbc.EstimateCostRequest{
		ResourceType: "aws:ec2/instance:Instance",
	}))
	require.NoError(t, err)
	assert.InDelta(t, 100.0, estimateResp.Msg.GetCostMonthly(), 0.01)

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestHealthHandler(t *testing.T) {
	handler := pluginsdk.HealthHandler()

	t.Run("GET returns 200 OK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "ok", w.Body.String())
		assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("HEAD returns 200 OK with no body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/healthz", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("POST returns 405 Method Not Allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestWebConfig_Defaults(t *testing.T) {
	config := pluginsdk.DefaultWebConfig()

	assert.False(t, config.Enabled)
	assert.Nil(t, config.AllowedOrigins)
	assert.False(t, config.AllowCredentials)
	assert.False(t, config.EnableHealthEndpoint)
}

func TestWebConfig_WithMethods(t *testing.T) {
	config := pluginsdk.DefaultWebConfig().
		WithWebEnabled(true).
		WithAllowedOrigins([]string{"http://localhost:3000"}).
		WithAllowCredentials(true).
		WithHealthEndpoint(true)

	assert.True(t, config.Enabled)
	assert.Equal(t, []string{"http://localhost:3000"}, config.AllowedOrigins)
	assert.True(t, config.AllowCredentials)
	assert.True(t, config.EnableHealthEndpoint)
}

func TestServeConnect_CORS(t *testing.T) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &connectTestPlugin{name: "cors-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server with CORS enabled
	errCh := make(chan error, 1)
	go func() {
		errCh <- pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin:   plugin,
			Listener: listener,
			Web: pluginsdk.WebConfig{
				Enabled:              true,
				AllowedOrigins:       []string{"http://localhost:3000"},
				AllowCredentials:     true,
				EnableHealthEndpoint: true,
			},
		})
	}()

	// Wait for server to be ready
	addr := listener.Addr().String()
	waitForServer(t, addr)

	t.Run("preflight request returns CORS headers", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodOptions, "http://"+addr+"/pulumicost.v1.CostSourceService/Name", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")

		resp, respErr := http.DefaultClient.Do(req)
		require.NoError(t, respErr)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
	})

	t.Run("actual request includes CORS headers", func(t *testing.T) {
		req, _ := http.NewRequest(
			http.MethodPost,
			"http://"+addr+"/pulumicost.v1.CostSourceService/Name",
			bytes.NewReader([]byte(`{}`)),
		)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Content-Type", "application/json")

		resp, respErr := http.DefaultClient.Do(req)
		require.NoError(t, respErr)
		defer resp.Body.Close()

		assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("disallowed origin gets no CORS headers", func(t *testing.T) {
		req, _ := http.NewRequest(
			http.MethodPost,
			"http://"+addr+"/pulumicost.v1.CostSourceService/Name",
			bytes.NewReader([]byte(`{}`)),
		)
		req.Header.Set("Origin", "http://evil.com")
		req.Header.Set("Content-Type", "application/json")

		resp, respErr := http.DefaultClient.Do(req)
		require.NoError(t, respErr)
		defer resp.Body.Close()

		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	})

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServeGRPC_LegacyMode(t *testing.T) {
	// Create a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	plugin := &connectTestPlugin{name: "legacy-test-plugin"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server with Web.Enabled = false (default)
	errCh := make(chan error, 1)
	go func() {
		errCh <- pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin:   plugin,
			Listener: listener,
			// Web.Enabled defaults to false
		})
	}()

	// Wait for server to be ready using TCP dial (can't use HTTP health check for gRPC)
	addr := listener.Addr().String()
	for range 50 { // 500ms total
		conn, dialErr := net.Dial("tcp", addr)
		if dialErr == nil {
			conn.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// HTTP/1.1 requests to a pure gRPC server will fail because gRPC uses HTTP/2.
	// The server returns an HTTP/2 response that can't be parsed by HTTP/1.1 client.
	_, err = http.Get("http://" + addr + "/healthz")
	// Expect an error because gRPC server speaks HTTP/2, not HTTP/1.1
	require.Error(t, err, "HTTP/1.1 request to gRPC server should fail")

	// Cleanup
	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
		t.Fatal("server did not shut down in time")
	}
}

func TestServe_CORSWildcardCredentialsValidation(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()

	plugin := &connectTestPlugin{name: "cors-validation-test"}
	ctx := context.Background()

	// Attempt to start server with AllowCredentials + wildcard origin
	// This should fail immediately with a validation error
	err = pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
		Plugin:   plugin,
		Listener: listener,
		Web: pluginsdk.WebConfig{
			Enabled:          true,
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "AllowCredentials cannot be used with wildcard origin")
}
