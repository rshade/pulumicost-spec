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

package pluginsdk //nolint:testpackage // White-box testing required

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// SlowTestPlugin implements a plugin that delays responses for timeout testing.
type SlowTestPlugin struct {
	delay time.Duration
}

func (p *SlowTestPlugin) Name() string {
	return "slow-test-plugin"
}

func (p *SlowTestPlugin) EstimateCost(
	ctx context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	select {
	case <-time.After(p.delay):
		return &pbc.EstimateCostResponse{}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *SlowTestPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (p *SlowTestPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (p *SlowTestPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

// StartSlowTestServer starts a test server with a slow plugin for timeout testing.
func StartSlowTestServer(delay time.Duration) (string, func(), error) {
	plugin := &SlowTestPlugin{delay: delay}

	// Create listener on random port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}

	// Serve in background using pluginsdk.Serve
	// We use Web.Enabled=true so it supports Connect protocol
	config := ServeConfig{
		Plugin:   plugin,
		Listener: lis,
		Web: WebConfig{
			Enabled: true,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		_ = Serve(ctx, config)
		close(done)
	}()

	// Return address and cleanup function
	addr := lis.Addr().String()
	cleanup := func() {
		cancel()
		<-done
		lis.Close()
	}

	return addr, cleanup, nil
}

// CreateTimeoutTestClient creates a client with specified timeout for testing.
func CreateTimeoutTestClient(serverAddr string, timeout time.Duration) (*Client, func()) {
	// Create client config with timeout
	// Default protocol is Connect, which works with pluginsdk.Serve(Web=true)
	cfg := DefaultClientConfig("http://" + serverAddr).WithTimeout(timeout)

	client := NewClient(cfg)

	cleanup := func() {
		client.Close()
	}

	return client, cleanup
}

func TestClientTimeout_ExceedsConfiguredTimeout(t *testing.T) {
	// Start server with 10-second delay
	serverAddr, serverCleanup, err := StartSlowTestServer(10 * time.Second)
	require.NoError(t, err)
	defer serverCleanup()

	// Create client with 2-second timeout
	client, clientCleanup := CreateTimeoutTestClient(serverAddr, 2*time.Second)
	defer clientCleanup()

	// Call EstimateCost - should timeout after 2 seconds (not 10)
	start := time.Now()
	_, err = client.EstimateCost(context.Background(), &pbc.EstimateCostRequest{})
	duration := time.Since(start)

	// Should get timeout error and complete within reasonable time
	require.Error(t, err)
	assert.Less(t, duration, 5*time.Second, "Should timeout within 5 seconds, took %v", duration)
	assert.Greater(t, duration, 1*time.Second, "Should take at least 1 second (client timeout), took %v", duration)
}

func TestClientTimeout_ContextDeadlinePrecedence(t *testing.T) {
	// Start server with 10-second delay
	serverAddr, serverCleanup, err := StartSlowTestServer(10 * time.Second)
	require.NoError(t, err)
	defer serverCleanup()

	// Create client with 30-second timeout (should not matter)
	client, clientCleanup := CreateTimeoutTestClient(serverAddr, 30*time.Second)
	defer clientCleanup()

	// Set context deadline of 1 second (should take precedence)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Call EstimateCost - should timeout after 1 second (context deadline)
	start := time.Now()
	_, err = client.EstimateCost(ctx, &pbc.EstimateCostRequest{})
	duration := time.Since(start)

	// Should get context deadline exceeded error
	require.Error(t, err)
	assert.Less(t, duration, 2*time.Second, "Should timeout within 2 seconds, took %v", duration)
	assert.Greater(t, duration, 500*time.Millisecond, "Should take at least 500ms, took %v", duration)

	// Verify it's a context deadline error (or wrapped)
	assert.Contains(t, err.Error(), "deadline exceeded")
}

func TestClientTimeout_DefaultValue(t *testing.T) {
	// Start server with 40-second delay (longer than default timeout)
	serverAddr, serverCleanup, err := StartSlowTestServer(40 * time.Second)
	require.NoError(t, err)
	defer serverCleanup()

	// Create client with zero timeout (should use default 30 seconds)
	cfg := DefaultClientConfig("http://" + serverAddr)
	// Don't call WithTimeout - timeout should be 0, triggering default

	client := NewClient(cfg)
	defer client.Close()

	// Call EstimateCost - should timeout after 30 seconds (default)
	start := time.Now()
	_, err = client.EstimateCost(context.Background(), &pbc.EstimateCostRequest{})
	duration := time.Since(start)

	// Should get timeout error and complete within reasonable time of 30 seconds
	require.Error(t, err)
	assert.Less(t, duration, 35*time.Second, "Should timeout within 35 seconds, took %v", duration)
	assert.Greater(t, duration, 25*time.Second, "Should take at least 25 seconds, took %v", duration)
}

func TestClientTimeout_CustomHTTPClientPrecedence(t *testing.T) {
	// Start server with 10-second delay
	serverAddr, serverCleanup, err := StartSlowTestServer(10 * time.Second)
	require.NoError(t, err)
	defer serverCleanup()

	// Create custom HTTP client with 3-second timeout
	customClient := &http.Client{Timeout: 3 * time.Second}

	// Create client config with custom HTTP client AND timeout field set
	// HTTPClient should take precedence
	cfg := ClientConfig{
		BaseURL:    "http://" + serverAddr,
		HTTPClient: customClient,
		Timeout:    30 * time.Second, // Should be ignored
	}

	client := NewClient(cfg)
	defer client.Close()

	// Call EstimateCost - should timeout after 3 seconds (custom HTTP client)
	start := time.Now()
	_, err = client.EstimateCost(context.Background(), &pbc.EstimateCostRequest{})
	duration := time.Since(start)

	// Should get timeout error and complete within reasonable time
	require.Error(t, err)
	assert.Less(t, duration, 6*time.Second, "Should timeout within 6 seconds, took %v", duration)
	assert.Greater(t, duration, 2*time.Second, "Should take at least 2 seconds, took %v", duration)
}
