package pluginsdk_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// mockPlugin implements pluginsdk.Plugin for testing.
type mockPlugin struct{}

func (p *mockPlugin) Name() string { return "mock-plugin" }
func (p *mockPlugin) GetProjectedCost(ctx context.Context, req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error) {
	return nil, nil
}
func (p *mockPlugin) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
	return nil, nil
}
func (p *mockPlugin) GetPricingSpec(ctx context.Context, req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error) {
	return nil, nil
}
func (p *mockPlugin) EstimateCost(ctx context.Context, req *pbc.EstimateCostRequest) (*pbc.EstimateCostResponse, error) {
	return nil, nil
}

// TestServeReflection verifies that the gRPC server started by Serve supports reflection.
func TestServeReflection(t *testing.T) {
	// 1. Get a free port
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()

	// 2. Start server in goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin: &mockPlugin{},
			Port:   port,
		})
		if err != nil && ctx.Err() == nil {
			errCh <- err
		}
	}()

	// 3. Wait for server to be ready (poll)
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	// Poll until connected
	deadline := time.Now().Add(5 * time.Second)
	for {
		if time.Now().After(deadline) {
			t.Fatal("Timeout waiting for server to start")
		}
		state := conn.GetState()
		if state == connectivity.Ready {
			break
		}
		// Try to dial to check if it's accepting connections, 
		// because GetState might not transition if we don't make a call?
		// NewClient uses lazy connection. conn.Connect() triggers it.
		conn.Connect() 
		
		// Actually, let's just try to create the reflection client and call it.
		// If it fails with "Unavailable", we retry.
		// If it fails with "Unimplemented", we know reflection is missing.
		time.Sleep(100 * time.Millisecond)
	}

	// 4. Check for reflection service
	refClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := refClient.ServerReflectionInfo(ctx)
	if err != nil {
		// This might happen if connection is not ready yet
		t.Logf("Failed to create stream: %v", err)
	}
	
	// Try to send a request to list services
	err = stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host: "",
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
			ListServices: "*",
		},
	})
	
	// If the server doesn't support reflection, the stream might work but return error on Recv
	// or Send might fail if the service is not found.
	
	if err == nil {
		_, err = stream.Recv()
	}

	if err != nil {
		t.Fatalf("Reflection failed: %v", err)
	}
	
	t.Log("Reflection works as expected")
}