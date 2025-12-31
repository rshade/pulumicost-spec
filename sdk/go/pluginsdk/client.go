// Package pluginsdk provides a development SDK for PulumiCost plugins.
package pluginsdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1/pbcconnect"
)

// Protocol specifies the RPC protocol to use for client connections.
type Protocol int

const (
	// ProtocolConnect uses the Connect protocol (default).
	// This provides the best browser compatibility with simple JSON over HTTP.
	ProtocolConnect Protocol = iota

	// ProtocolGRPC uses the gRPC protocol.
	// Requires HTTP/2 and is ideal for server-to-server communication.
	ProtocolGRPC

	// ProtocolGRPCWeb uses the gRPC-Web protocol.
	// Provides gRPC compatibility for browser clients.
	ProtocolGRPCWeb
)

// allProtocols is the exhaustive list of valid Protocol values.
// Used for zero-allocation validation following the registry package pattern.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allProtocols = []Protocol{ProtocolConnect, ProtocolGRPC, ProtocolGRPCWeb}

// AllProtocols returns all valid Protocol values.
func AllProtocols() []Protocol {
	return allProtocols
}

// IsValidProtocol returns true if the protocol is a known valid value.
func IsValidProtocol(p Protocol) bool {
	for _, valid := range allProtocols {
		if p == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of the Protocol.
func (p Protocol) String() string {
	switch p {
	case ProtocolConnect:
		return "Connect"
	case ProtocolGRPC:
		return "gRPC"
	case ProtocolGRPCWeb:
		return "gRPC-Web"
	default:
		return "Unknown"
	}
}

// DefaultClientTimeout is the default HTTP client timeout for plugin requests.
const DefaultClientTimeout = 30 * time.Second

// Default connection pool settings for high-throughput scenarios.
const (
	// DefaultMaxIdleConns is the max idle connections across all hosts.
	DefaultMaxIdleConns = 100
	// DefaultMaxIdleConnsPerHost is the max idle connections per host.
	DefaultMaxIdleConnsPerHost = 10
	// DefaultIdleConnTimeout is how long idle connections are kept in the pool.
	DefaultIdleConnTimeout = 90 * time.Second
)

// ClientConfig configures a PulumiCost client.
type ClientConfig struct {
	// BaseURL is the server's base URL (e.g., "http://localhost:8080").
	BaseURL string

	// Protocol specifies the RPC protocol to use.
	// Defaults to ProtocolConnect.
	Protocol Protocol

	// HTTPClient is the HTTP client to use for requests.
	// If nil, a default client with 30-second timeout is used.
	HTTPClient *http.Client

	// ConnectOptions allows passing additional connect.ClientOption values.
	ConnectOptions []connect.ClientOption
}

// DefaultClientConfig returns a ClientConfig with sensible defaults.
func DefaultClientConfig(baseURL string) ClientConfig {
	return ClientConfig{
		BaseURL:  baseURL,
		Protocol: ProtocolConnect,
		HTTPClient: &http.Client{
			Timeout: DefaultClientTimeout,
		},
	}
}

// HighThroughputClientConfig returns a ClientConfig optimized for high-throughput scenarios.
// It configures connection pooling for better performance when making many requests.
func HighThroughputClientConfig(baseURL string) ClientConfig {
	transport := &http.Transport{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:     DefaultIdleConnTimeout,
	}

	return ClientConfig{
		BaseURL:  baseURL,
		Protocol: ProtocolConnect,
		HTTPClient: &http.Client{
			Timeout:   DefaultClientTimeout,
			Transport: transport,
		},
	}
}

// Client provides a simplified interface for communicating with PulumiCost plugins.
// It wraps the generated connect client with cleaner method signatures.
//
// Clients should be reused for multiple calls rather than creating a new client
// for each request. When done, call Close() to release connection pool resources.
type Client struct {
	inner      pbcconnect.CostSourceServiceClient
	httpClient *http.Client
	ownsClient bool // true if we created the http.Client and should close it
}

// wrapRPCError wraps an RPC error with context about the operation.
// It distinguishes context cancellation/timeout from other errors for better debugging.
// When the context error is available, it is included using errors.Join to preserve
// both the context error type (DeadlineExceeded vs Canceled) and the original error.
func wrapRPCError(ctx context.Context, operation string, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return errors.Join(
			fmt.Errorf("%s RPC cancelled or timed out", operation),
			ctxErr,
			err,
		)
	}
	return fmt.Errorf("%s RPC failed: %w", operation, err)
}

// NewClient creates a new PulumiCost client with the given configuration.
// If an invalid protocol is specified, it defaults to ProtocolConnect for backward compatibility.
//
// HTTP Client Ownership:
//   - If HTTPClient is nil, an internal HTTP client is created and owned by this Client.
//     Call Close() when done to release connection pool resources.
//   - If HTTPClient is provided, the caller retains ownership. Close() is a no-op;
//     the caller is responsible for closing the HTTP client.
//
// Thread Safety: Client is safe for concurrent use from multiple goroutines.
func NewClient(cfg ClientConfig) *Client {
	// Validate protocol - default to Connect for invalid values
	if !IsValidProtocol(cfg.Protocol) {
		cfg.Protocol = ProtocolConnect
	}

	httpClient := cfg.HTTPClient
	ownsClient := false
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultClientTimeout}
		ownsClient = true
	}

	// Build connect options based on protocol
	opts := make([]connect.ClientOption, 0, len(cfg.ConnectOptions)+1)
	switch cfg.Protocol {
	case ProtocolConnect:
		// Connect is the default protocol, no extra option needed
	case ProtocolGRPC:
		opts = append(opts, connect.WithGRPC())
	case ProtocolGRPCWeb:
		opts = append(opts, connect.WithGRPCWeb())
	}
	opts = append(opts, cfg.ConnectOptions...)

	return &Client{
		inner:      pbcconnect.NewCostSourceServiceClient(httpClient, cfg.BaseURL, opts...),
		httpClient: httpClient,
		ownsClient: ownsClient,
	}
}

// NewConnectClient creates a client using the Connect protocol (JSON over HTTP).
// This is the recommended protocol for browser clients.
//
// The SDK owns the created HTTP client. Call Close() when done to release
// connection pool resources.
func NewConnectClient(baseURL string) *Client {
	cfg := DefaultClientConfig(baseURL)
	cfg.Protocol = ProtocolConnect
	client := NewClient(cfg)
	// Clients created via convenience constructors are SDK-owned:
	// Close() should release connection pool resources.
	client.ownsClient = true
	return client
}

// NewGRPCClient creates a client using the gRPC protocol.
// Requires HTTP/2 support on both client and server.
//
// The SDK owns the created HTTP client. Call Close() when done to release
// connection pool resources.
func NewGRPCClient(baseURL string) *Client {
	cfg := DefaultClientConfig(baseURL)
	cfg.Protocol = ProtocolGRPC
	client := NewClient(cfg)
	// Clients created via convenience constructors are SDK-owned:
	// Close() should release connection pool resources.
	client.ownsClient = true
	return client
}

// NewGRPCWebClient creates a client using the gRPC-Web protocol.
// Useful for browser clients that need gRPC compatibility.
//
// The SDK owns the created HTTP client. Call Close() when done to release
// connection pool resources.
func NewGRPCWebClient(baseURL string) *Client {
	cfg := DefaultClientConfig(baseURL)
	cfg.Protocol = ProtocolGRPCWeb
	client := NewClient(cfg)
	// Clients created via convenience constructors are SDK-owned:
	// Close() should release connection pool resources.
	client.ownsClient = true
	return client
}

// Name returns the display name of the cost source plugin.
func (c *Client) Name(ctx context.Context) (string, error) {
	resp, err := c.inner.Name(ctx, connect.NewRequest(&pbc.NameRequest{}))
	if err != nil {
		return "", wrapRPCError(ctx, "Name", err)
	}
	return resp.Msg.GetName(), nil
}

// Supports checks if the cost source supports pricing for a given resource.
func (c *Client) Supports(ctx context.Context, resource *pbc.ResourceDescriptor) (*pbc.SupportsResponse, error) {
	if resource == nil {
		return nil, errors.New("resource descriptor cannot be nil")
	}
	if resource.GetResourceType() == "" {
		return nil, errors.New("resource type is required")
	}
	resp, err := c.inner.Supports(ctx, connect.NewRequest(&pbc.SupportsRequest{
		Resource: resource,
	}))
	if err != nil {
		return nil, wrapRPCError(ctx, "Supports", err)
	}
	return resp.Msg, nil
}

// SupportsResourceType is a convenience method to check support by resource type string.
func (c *Client) SupportsResourceType(ctx context.Context, resourceType string) (bool, error) {
	resp, err := c.Supports(ctx, &pbc.ResourceDescriptor{
		ResourceType: resourceType,
	})
	if err != nil {
		return false, err
	}
	return resp.GetSupported(), nil
}

// EstimateCost returns an estimated monthly cost for a resource.
func (c *Client) EstimateCost(ctx context.Context, req *pbc.EstimateCostRequest) (*pbc.EstimateCostResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.EstimateCost(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "EstimateCost", err)
	}
	return resp.Msg, nil
}

// GetActualCost retrieves historical cost data for a specific resource.
func (c *Client) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.GetActualCost(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "GetActualCost", err)
	}
	return resp.Msg, nil
}

// GetProjectedCost calculates projected cost information for a resource.
func (c *Client) GetProjectedCost(
	ctx context.Context,
	req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.GetProjectedCost(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "GetProjectedCost", err)
	}
	return resp.Msg, nil
}

// GetPricingSpec returns detailed pricing specification for a resource type.
func (c *Client) GetPricingSpec(
	ctx context.Context,
	req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.GetPricingSpec(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "GetPricingSpec", err)
	}
	return resp.Msg, nil
}

// GetRecommendations retrieves cost optimization recommendations.
func (c *Client) GetRecommendations(
	ctx context.Context,
	req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.GetRecommendations(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "GetRecommendations", err)
	}
	return resp.Msg, nil
}

// DismissRecommendation marks a recommendation as dismissed/ignored.
func (c *Client) DismissRecommendation(
	ctx context.Context,
	req *pbc.DismissRecommendationRequest,
) (*pbc.DismissRecommendationResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.DismissRecommendation(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "DismissRecommendation", err)
	}
	return resp.Msg, nil
}

// GetBudgets returns budget information from the cost management service.
func (c *Client) GetBudgets(ctx context.Context, req *pbc.GetBudgetsRequest) (*pbc.GetBudgetsResponse, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	resp, err := c.inner.GetBudgets(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, wrapRPCError(ctx, "GetBudgets", err)
	}
	return resp.Msg, nil
}

// Inner returns the underlying connect client for advanced use cases.
func (c *Client) Inner() pbcconnect.CostSourceServiceClient {
	return c.inner
}

// Close releases resources associated with the client.
// It closes idle connections in the underlying HTTP transport pool.
// Active connections are not forcibly closed - they will be cleaned up
// after the current requests complete. For immediate cleanup, cancel
// the request contexts before calling Close().
//
// This is a no-op if the client was created with a user-provided HTTPClient
// (in that case, the caller is responsible for closing it).
//
// Close is safe to call multiple times.
func (c *Client) Close() {
	if !c.ownsClient || c.httpClient == nil {
		return
	}
	c.httpClient.CloseIdleConnections()
}
