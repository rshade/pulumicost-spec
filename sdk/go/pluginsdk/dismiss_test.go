//nolint:testpackage // Testing internal Server implementation with mocks
package pluginsdk

import (
	"context"
	"errors"
	"testing"
	"time"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// requireGRPCError asserts that err is a gRPC status error with the expected code
// and that the message contains the expected substring.
func requireGRPCError(t *testing.T, err error, expectedCode codes.Code, msgContains string) {
	t.Helper()
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "error should be a gRPC status error")
	assert.Equal(t, expectedCode, st.Code())
	assert.Contains(t, st.Message(), msgContains)
}

// mockDismissPlugin implements both Plugin and DismissProvider interfaces.
type mockDismissPlugin struct {
	mockPlugin

	success       bool          // Controls Success field in response
	err           error         // Error to return from DismissRecommendation
	returnNil     bool          // Forces nil response to test server error handling
	checkContext  bool          // Whether to check context before processing
	simulateDelay time.Duration // Simulated processing delay (respects context)
}

func (m *mockDismissPlugin) DismissRecommendation(
	ctx context.Context,
	_ *pbc.DismissRecommendationRequest,
) (*pbc.DismissRecommendationResponse, error) {
	// Check context first if configured (simulates real plugin behavior)
	if m.checkContext {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
	}

	// Simulate processing delay that respects context
	if m.simulateDelay > 0 {
		select {
		case <-time.After(m.simulateDelay):
			// Delay completed
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if m.err != nil {
		return nil, m.err
	}
	if m.returnNil {
		//nolint:nilnil // Intentional nil return to test server error handling
		return nil, nil
	}
	return &pbc.DismissRecommendationResponse{
		Success: m.success,
	}, nil
}

func TestDismissRecommendation_PluginImplements(t *testing.T) {
	plugin := &mockDismissPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		success:    true,
	}
	server := NewServer(plugin)

	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}
	resp, err := server.DismissRecommendation(context.Background(), req)

	require.NoError(t, err)
	assert.True(t, resp.GetSuccess())
}

func TestDismissRecommendation_PluginNotImplements(t *testing.T) {
	// mockPlugin does not implement DismissProvider
	plugin := &mockPlugin{name: "test-plugin"}
	server := NewServer(plugin)

	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}
	_, err := server.DismissRecommendation(context.Background(), req)

	requireGRPCError(t, err, codes.Unimplemented, "plugin does not support DismissRecommendation")
}

func TestDismissRecommendation_PluginError(t *testing.T) {
	plugin := &mockDismissPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		err:        errors.New("db error"),
	}
	server := NewServer(plugin)

	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}
	_, err := server.DismissRecommendation(context.Background(), req)

	requireGRPCError(t, err, codes.Internal, "plugin failed to execute DismissRecommendation")
}

func TestDismissRecommendation_NilResponse(t *testing.T) {
	plugin := &mockDismissPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		returnNil:  true,
	}
	server := NewServer(plugin)

	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}
	_, err := server.DismissRecommendation(context.Background(), req)

	requireGRPCError(t, err, codes.Internal, "plugin returned a nil response")
}

func TestDismissRecommendation_ContextCanceled(t *testing.T) {
	plugin := &mockDismissPlugin{
		mockPlugin:   mockPlugin{name: "test-plugin"},
		checkContext: true, // Plugin checks context before processing
		success:      true,
	}
	server := NewServer(plugin)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}
	_, err := server.DismissRecommendation(ctx, req)

	// Plugin returns context.Canceled, server wraps as Internal error
	requireGRPCError(t, err, codes.Internal, "plugin failed to execute DismissRecommendation")
}

func TestDismissRecommendation_ContextTimeout(t *testing.T) {
	plugin := &mockDismissPlugin{
		mockPlugin:    mockPlugin{name: "test-plugin"},
		simulateDelay: 100 * time.Millisecond, // Slow operation
		success:       true,
	}
	server := NewServer(plugin)

	// Short timeout that expires before delay completes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}
	_, err := server.DismissRecommendation(ctx, req)

	// Plugin returns context.DeadlineExceeded, server wraps as Internal error
	requireGRPCError(t, err, codes.Internal, "plugin failed to execute DismissRecommendation")
}

func BenchmarkDismissRecommendation(b *testing.B) {
	plugin := &mockDismissPlugin{
		mockPlugin: mockPlugin{name: "test-plugin"},
		success:    true,
	}
	server := NewServer(plugin)

	ctx := context.Background()
	req := &pbc.DismissRecommendationRequest{
		RecommendationId: "rec-123",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = server.DismissRecommendation(ctx, req)
	}
}
