// Package pluginsdk provides a development SDK for PulumiCost plugins.
package pluginsdk

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// portFlag is the --port command-line flag for specifying the gRPC server port.
// This is registered at package initialization time.
//
//nolint:gochecknoglobals // Package-level flag required for command-line parsing
var portFlag = flag.Int("port", 0, "TCP port for gRPC server (overrides PULUMICOST_PLUGIN_PORT)")

// ParsePortFlag returns the value of the --port command-line flag.
// Returns 0 if the flag was not specified or if flag.Parse() has not been called.
//
// IMPORTANT: The caller must call flag.Parse() before calling this function.
//
// Example usage in plugin main():
//
//	func main() {
//	    flag.Parse()  // Must be called first
//	    port := pluginsdk.ParsePortFlag()
//	    ctx := context.Background()
//	    if err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
//	        Plugin: &MyPlugin{},
//	        Port:   port,
//	    }); err != nil {
//	        log.Fatal(err)
//	    }
//	}
func ParsePortFlag() int {
	if portFlag == nil {
		return 0
	}
	return *portFlag
}

// DefaultSupportsNotImplementedReason is the standardized message returned when
// a plugin does not implement the SupportsProvider interface.
const DefaultSupportsNotImplementedReason = "Supports capability not implemented by this plugin"

// Plugin represents a PulumiCost plugin implementation.
type Plugin interface {
	// Name returns the plugin name identifier.
	Name() string
	// GetProjectedCost calculates projected cost for a resource.
	GetProjectedCost(ctx context.Context, req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error)
	// GetActualCost retrieves actual cost for a resource.
	GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error)
	// GetPricingSpec returns detailed pricing specification for a resource type.
	GetPricingSpec(ctx context.Context, req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error)
	// EstimateCost returns an estimated monthly cost for a resource based on its type and configuration attributes.
	EstimateCost(ctx context.Context, req *pbc.EstimateCostRequest) (*pbc.EstimateCostResponse, error)
}

// SupportsProvider is an optional interface that plugins can implement to indicate
// whether they support pricing for specific resource types. Plugins that do not
// implement this interface will receive a default "not supported" response.
type SupportsProvider interface {
	// Supports checks if the plugin supports pricing for the given resource.
	Supports(ctx context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error)
}

// RecommendationsProvider is an optional interface that plugins can implement
// to provide cost optimization recommendations. Plugins that do not implement
// this interface will return an empty list when GetRecommendations is called.
type RecommendationsProvider interface {
	// GetRecommendations retrieves cost optimization recommendations.
	GetRecommendations(ctx context.Context, req *pbc.GetRecommendationsRequest) (
		*pbc.GetRecommendationsResponse, error)
}

// BudgetsProvider is an optional interface that plugins can implement
// to provide budget information from cloud cost management services.
// Plugins that do not implement this interface will return Unimplemented
// when GetBudgets is called.
type BudgetsProvider interface {
	// GetBudgets retrieves budget information from the cost management service.
	GetBudgets(ctx context.Context, req *pbc.GetBudgetsRequest) (
		*pbc.GetBudgetsResponse, error)
}

// RegistryLookup defines the interface for looking up plugins by provider and region.
// This is used to validate incoming Supports requests against registered plugins.
type RegistryLookup interface {
	// FindPlugin returns the plugin name for the given provider and region.
	// Returns empty string if no plugin is registered for the combination.
	FindPlugin(provider, region string) string
}

// DefaultRegistryLookup provides a no-op registry lookup that always returns empty.
// This causes all Supports() calls to return InvalidArgument since no plugin
// can be found. Use a real RegistryLookup implementation in production.
type DefaultRegistryLookup struct{}

// FindPlugin always returns empty string indicating no plugin is registered.
func (d *DefaultRegistryLookup) FindPlugin(_, _ string) string {
	return ""
}

// Server wraps a Plugin implementation with a gRPC server.
type Server struct {
	pbc.UnimplementedCostSourceServiceServer

	plugin   Plugin
	registry RegistryLookup
	logger   zerolog.Logger
}

// NewServer creates a Server that exposes the provided Plugin over gRPC.
// Uses DefaultRegistryLookup which returns empty for all lookups, causing
// Supports() calls to return InvalidArgument. Use NewServerWithRegistry
// to provide a real registry for production use.
func NewServer(plugin Plugin) *Server {
	return &Server{
		plugin:   plugin,
		registry: &DefaultRegistryLookup{},
		logger:   newDefaultLogger(),
	}
}

// NewServerWithRegistry creates a Server with a custom registry lookup.
// If registry is nil, DefaultRegistryLookup is used.
func NewServerWithRegistry(plugin Plugin, registry RegistryLookup) *Server {
	return NewServerWithOptions(plugin, registry, nil)
}

// NewServerWithOptions creates a Server with custom registry and logger.
// If registry is nil, DefaultRegistryLookup is used.
// If logger is nil, a default logger is used.
func NewServerWithOptions(plugin Plugin, registry RegistryLookup, logger *zerolog.Logger) *Server {
	if registry == nil {
		registry = &DefaultRegistryLookup{}
	}
	var log zerolog.Logger
	if logger != nil {
		log = *logger
	} else {
		log = newDefaultLogger()
	}
	return &Server{
		plugin:   plugin,
		registry: registry,
		logger:   log,
	}
}

// Name implements the gRPC Name method.
func (s *Server) Name(_ context.Context, _ *pbc.NameRequest) (*pbc.NameResponse, error) {
	return &pbc.NameResponse{Name: s.plugin.Name()}, nil
}

// GetProjectedCost implements the gRPC GetProjectedCost method.
func (s *Server) GetProjectedCost(
	ctx context.Context,
	req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	return s.plugin.GetProjectedCost(ctx, req)
}

// GetActualCost implements the gRPC GetActualCost method.
func (s *Server) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
	return s.plugin.GetActualCost(ctx, req)
}

// GetPricingSpec implements the gRPC GetPricingSpec method.
func (s *Server) GetPricingSpec(
	ctx context.Context,
	req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	return s.plugin.GetPricingSpec(ctx, req)
}

// EstimateCost implements the gRPC EstimateCost method.
func (s *Server) EstimateCost(
	ctx context.Context,
	req *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	return s.plugin.EstimateCost(ctx, req)
}

// Supports implements the gRPC Supports method.
// It performs two-step validation: first checks registry for plugin by provider/region,
// then delegates to the plugin's Supports method if implemented.
func (s *Server) Supports(ctx context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
	// Validate request has resource descriptor
	if req.GetResource() == nil {
		return nil, status.Error(codes.InvalidArgument, "resource descriptor is required")
	}

	resource := req.GetResource()
	provider := resource.GetProvider()
	region := resource.GetRegion()

	// Step 1: Registry lookup - validate provider/region combination
	pluginName := s.registry.FindPlugin(provider, region)
	if pluginName == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"no plugin registered for provider %q and region %q",
			provider,
			region,
		)
	}

	// Step 2: Check if plugin implements SupportsProvider
	supportsProvider, ok := s.plugin.(SupportsProvider)
	if !ok {
		// Plugin does not implement SupportsProvider - return default response
		return &pbc.SupportsResponse{
				Supported: false,
				Reason:    DefaultSupportsNotImplementedReason,
			},
			nil
	}

	// Delegate to plugin's Supports method
	resp, err := supportsProvider.Supports(ctx, req)
	if err != nil {
		// Log the detailed error server-side for debugging
		s.logger.Error().
			Str(FieldResourceType, resource.GetResourceType()).
			Str(FieldProvider, provider).
			Str(FieldRegion, region).
			Err(err).
			Msg("Supports handler error")
		// Return generic message to client (internal error details not exposed)
		return nil, status.Error(codes.Internal, "plugin failed to execute")
	}

	return resp, nil
}

// GetRecommendations implements the gRPC GetRecommendations method.
// GetRecommendations handles GetRecommendations RPC requests.
// If the plugin implements RecommendationsProvider, delegates to it.
// Otherwise returns an empty list (not an error) per FR-012.
func (s *Server) GetRecommendations(
	ctx context.Context,
	req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	// Log incoming request with filter details
	filter := req.GetFilter()
	s.logger.Debug().
		Str(FieldFilterCategory, filter.GetCategory().String()).
		Str(FieldFilterActionType, filter.GetActionType().String()).
		Int32(FieldPageSize, req.GetPageSize()).
		Msg("GetRecommendations request received")

	// Check if plugin implements RecommendationsProvider
	recProvider, ok := s.plugin.(RecommendationsProvider)
	if !ok {
		// Plugin does not implement recommendations - return empty list per FR-012
		// Include projection_period from request for client consistency
		s.logger.Debug().
			Int(FieldRecommendationCount, 0).
			Msg("GetRecommendations returning empty response (not implemented)")
		return &pbc.GetRecommendationsResponse{
			Recommendations: []*pbc.Recommendation{},
			Summary: &pbc.RecommendationSummary{
				TotalRecommendations: 0,
				ProjectionPeriod:     req.GetProjectionPeriod(),
			},
			NextPageToken: "",
		}, nil
	}

	// Delegate to plugin's GetRecommendations method
	resp, err := recProvider.GetRecommendations(ctx, req)
	if err != nil {
		s.logger.Error().
			Str(FieldFilterCategory, filter.GetCategory().String()).
			Str(FieldFilterActionType, filter.GetActionType().String()).
			Err(err).
			Msg("GetRecommendations handler error")
		return nil, status.Error(codes.Internal, "plugin failed to execute GetRecommendations")
	}

	// Guard against nil response from plugin.
	if resp == nil {
		s.logger.Error().
			Str(FieldFilterCategory, filter.GetCategory().String()).
			Str(FieldFilterActionType, filter.GetActionType().String()).
			Msg("GetRecommendations handler returned a nil response")
		return nil, status.Error(codes.Internal, "plugin returned a nil response")
	}

	// Log successful response with summary
	summary := resp.GetSummary()
	if summary == nil {
		s.logger.Error().
			Str(FieldFilterCategory, filter.GetCategory().String()).
			Str(FieldFilterActionType, filter.GetActionType().String()).
			Msg("GetRecommendations response has nil summary")
		return nil, status.Error(codes.Internal, "plugin returned response with nil summary")
	}
	s.logger.Info().
		Int32(FieldRecommendationCount, summary.GetTotalRecommendations()).
		Float64(FieldTotalSavings, summary.GetTotalEstimatedSavings()).
		Msg("GetRecommendations completed")

	return resp, nil
}

// GetBudgets implements the gRPC GetBudgets method.
// GetBudgets handles GetBudgets RPC requests.
// If the plugin implements BudgetsProvider, delegates to it.
// Otherwise returns Unimplemented error per specification.
func (s *Server) GetBudgets(
	ctx context.Context,
	req *pbc.GetBudgetsRequest,
) (*pbc.GetBudgetsResponse, error) {
	// Log incoming request
	s.logger.Debug().
		Bool(FieldIncludeStatus, req.GetIncludeStatus()).
		Msg("GetBudgets request received")

	// Check if plugin implements BudgetsProvider
	budgetsProvider, ok := s.plugin.(BudgetsProvider)
	if !ok {
		// Plugin does not implement budgets - return Unimplemented per spec
		s.logger.Debug().Msg("GetBudgets returning Unimplemented (not supported by plugin)")
		return nil, status.Error(codes.Unimplemented, "plugin does not support GetBudgets")
	}

	// Delegate to plugin's GetBudgets method
	resp, err := budgetsProvider.GetBudgets(ctx, req)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("GetBudgets handler error")
		return nil, status.Error(codes.Internal, "plugin failed to execute GetBudgets")
	}

	// Guard against nil response from plugin
	if resp == nil {
		s.logger.Error().Msg("GetBudgets handler returned a nil response")
		return nil, status.Error(codes.Internal, "plugin returned a nil response")
	}

	// Log successful response
	summary := resp.GetSummary()
	if summary != nil {
		s.logger.Info().
			Int32(FieldTotalBudgets, summary.GetTotalBudgets()).
			Int32(FieldBudgetsOk, summary.GetBudgetsOk()).
			Int32(FieldBudgetsWarning, summary.GetBudgetsWarning()).
			Int32(FieldBudgetsCritical, summary.GetBudgetsCritical()).
			Int32(FieldBudgetsExceeded, summary.GetBudgetsExceeded()).
			Msg("GetBudgets completed")
	}

	return resp, nil
}

// ServeConfig holds configuration for serving a plugin.
type ServeConfig struct {
	Plugin   Plugin
	Port     int             // If 0, will use PULUMICOST_PLUGIN_PORT env var or random port
	Registry RegistryLookup  // Optional; if nil, DefaultRegistryLookup is used
	Logger   *zerolog.Logger // Optional; if nil, a default logger is used
	// Listener is an optional pre-configured listener.
	// If provided, Port is ignored and this listener is used for serving.
	Listener net.Listener
	// UnaryInterceptors is an optional list of gRPC unary server interceptors
	// to chain after the built-in TracingUnaryServerInterceptor.
	// Interceptors execute in order: tracing first, then each interceptor
	// in the order provided. If nil or empty, only the tracing interceptor runs.
	// Note: Passing nil elements in the slice will cause a panic (standard gRPC behavior).
	UnaryInterceptors []grpc.UnaryServerInterceptor
}

// resolvePort determines the port to use with the following priority:
//  1. requested (set from --port flag via ParsePortFlag(), or explicitly configured in ServeConfig.Port)
//  2. PULUMICOST_PLUGIN_PORT env var (via GetPort())
//  3. 0 (ephemeral port - OS assigns available port)
//
// Note: The generic PORT env var is NOT supported to avoid multi-plugin conflicts.
// When pulumicost-core spawns multiple plugins (e.g., aws-public + aws-ce), each
// needs a unique port. Using --port flag allows the core to allocate distinct ports.
func resolvePort(requested int) int {
	if requested > 0 {
		return requested
	}
	// Use centralized GetPort() which reads PULUMICOST_PLUGIN_PORT only (no fallback)
	return GetPort()
}

func listenOnLoopback(ctx context.Context, port int) (net.Listener, *net.TCPAddr, error) {
	address := "127.0.0.1:0"
	if port > 0 {
		address = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	}
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen: %w", err)
	}

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		closeErr := listener.Close()
		if closeErr != nil {
			return nil, nil, errors.Join(
				errors.New("listener address is not TCP"),
				fmt.Errorf("closing listener: %w", closeErr),
			)
		}
		return nil, nil, errors.New("listener address is not TCP")
	}

	return listener, tcpAddr, nil
}

func announcePort(listener net.Listener, addr *net.TCPAddr) error {
	if _, err := fmt.Fprintf(os.Stdout, "PORT=%d\n", addr.Port); err != nil {
		closeErr := listener.Close()
		if closeErr != nil {
			return errors.Join(
				fmt.Errorf("writing port: %w", err),
				fmt.Errorf("closing listener: %w", closeErr),
			)
		}
		return fmt.Errorf("writing port: %w", err)
	}
	return nil
}

// Serve starts the gRPC server for the provided plugin and prints the chosen port as PORT=<port> to stdout.
//
// It uses config.Port when > 0; if config.Port is 0 it reads the PULUMICOST_PLUGIN_PORT environment variable
// and falls back to an ephemeral port when none is provided. The function registers the plugin's service, begins
// serving on the selected port, and performs a graceful stop when the context is cancelled.
//
// Returns an error if the listener cannot be created or if the gRPC server fails to serve.
func Serve(ctx context.Context, config ServeConfig) error {
	var listener net.Listener
	var tcpAddr *net.TCPAddr
	var err error

	if config.Listener != nil {
		listener = config.Listener
		var ok bool
		tcpAddr, ok = listener.Addr().(*net.TCPAddr)
		if !ok {
			return errors.New("provided listener address is not TCP")
		}
	} else {
		port := resolvePort(config.Port)
		listener, tcpAddr, err = listenOnLoopback(ctx, port)
		if err != nil {
			return err
		}
	}

	if announceErr := announcePort(listener, tcpAddr); announceErr != nil {
		return announceErr
	}

	// Build interceptor chain: tracing first, then user interceptors
	interceptors := make([]grpc.UnaryServerInterceptor, 0, 1+len(config.UnaryInterceptors))
	interceptors = append(interceptors, TracingUnaryServerInterceptor())
	interceptors = append(interceptors, config.UnaryInterceptors...)

	// Create and register server with interceptor chain
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	server := NewServerWithOptions(config.Plugin, config.Registry, config.Logger)
	pbc.RegisterCostSourceServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	// Start serving
	err = grpcServer.Serve(listener)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	return nil
}
