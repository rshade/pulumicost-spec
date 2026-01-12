package pluginsdk_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// mockPlugin implements pluginsdk.Plugin for testing.
type mockPlugin struct{}

func (p *mockPlugin) Name() string { return "mock-plugin" }

func (p *mockPlugin) GetProjectedCost(
	_ context.Context,
	_ *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return &pbc.GetProjectedCostResponse{}, nil
}

func (p *mockPlugin) GetActualCost(
	_ context.Context,
	_ *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	return &pbc.GetActualCostResponse{}, nil
}

func (p *mockPlugin) GetPricingSpec(
	_ context.Context,
	_ *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return &pbc.GetPricingSpecResponse{}, nil
}

func (p *mockPlugin) EstimateCost(
	_ context.Context,
	_ *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return &pbc.EstimateCostResponse{}, nil
}

// TestServeReflection verifies that the gRPC server started by Serve supports reflection.
func TestServeReflection(t *testing.T) {
	// 1. Get a free port by listening (and keeping it open to pass to Serve)
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port

	// 2. Start server in goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		// Pass the listener directly to Serve to avoid race condition
		serveErr := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
			Plugin:   &mockPlugin{},
			Listener: l,
		})
		if serveErr != nil && !errors.Is(serveErr, context.Canceled) {
			errCh <- serveErr
		}
		close(errCh)
	}()

	// 3. Connect and verify reflection service
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	// Check for immediate server startup failure before entering the loop
	select {
	case serveErr := <-errCh:
		t.Fatalf("Server exited immediately: %v", serveErr)
	default:
	}

	// Poll until success or timeout
	deadline := time.Now().Add(5 * time.Second)
	var lastErr error

	for time.Now().Before(deadline) {
		select {
		case serveErr := <-errCh:
			t.Fatalf("Server exited unexpectedly: %v", serveErr)
		default:
		}

		if checkErr := checkReflection(ctx, conn); checkErr != nil {
			lastErr = checkErr
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return // Success
	}

	t.Fatalf("Reflection test failed after timeout. Last error: %v", lastErr)
}

//nolint:staticcheck // Validating legacy reflection API
func checkReflection(ctx context.Context, conn grpc.ClientConnInterface) error {
	refClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	stream, err := refClient.ServerReflectionInfo(ctx)
	if err != nil {
		return err
	}

	if sendErr := stream.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host: "",
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
			ListServices: "*",
		},
	}); sendErr != nil {
		return sendErr
	}

	resp, err := stream.Recv()
	if err != nil {
		return err
	}

	// Verify response contains the expected service
	listResp := resp.GetListServicesResponse()
	if listResp == nil {
		return errors.New("received nil ListServicesResponse")
	}

	for _, svc := range listResp.GetService() {
		if svc.GetName() == "finfocus.v1.CostSourceService" {
			return nil
		}
	}

	return errors.New("CostSourceService not found in reflection response")
}
