# Quickstart: pluginsdk.Serve()

**Feature**: 001-pluginsdk-serve-docs
**Date**: 2025-12-08

This quickstart provides copy-paste-ready examples for using `pluginsdk.Serve()` to start a
FinFocus plugin gRPC server.

## Minimal Plugin Example

The simplest possible plugin implementation:

```go
package main

import (
    "context"
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// MyPlugin implements the pluginsdk.Plugin interface.
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-cost-plugin"
}

func (p *MyPlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    // Implement projected cost calculation
    return &pbc.GetProjectedCostResponse{}, nil
}

func (p *MyPlugin) GetActualCost(
    ctx context.Context,
    req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
    // Implement actual cost retrieval
    return &pbc.GetActualCostResponse{}, nil
}

func (p *MyPlugin) GetPricingSpec(
    ctx context.Context,
    req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
    // Implement pricing spec lookup
    return &pbc.GetPricingSpecResponse{}, nil
}

func (p *MyPlugin) EstimateCost(
    ctx context.Context,
    req *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
    // Implement cost estimation
    return &pbc.EstimateCostResponse{}, nil
}

func main() {
    // Parse command-line flags (required before ParsePortFlag)
    flag.Parse()

    // Create cancellable context for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle SIGINT and SIGTERM for graceful shutdown
    ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
    defer stop() // Release resources associated with the context

    // Start the gRPC server
    if err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: &MyPlugin{},
        Port:   pluginsdk.ParsePortFlag(), // Uses --port flag if provided
    }); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

## Port Configuration Examples

### Using --port Flag

```bash
# Start plugin on specific port
./my-plugin --port 50051
```

### Using Environment Variable

```bash
# Start plugin using environment variable
export PULUMICOST_PLUGIN_PORT=50052
./my-plugin
```

### Using Ephemeral Port

```bash
# Start plugin with OS-assigned port (outputs PORT=<number> to stdout)
./my-plugin
# Output: PORT=54321
```

## Port Resolution Priority

The port is determined in this order:

1. **--port flag** (if provided and > 0)
2. **PULUMICOST_PLUGIN_PORT** environment variable (if set and valid)
3. **Ephemeral port**: The OS assigns an available port (when Port=0 or both above are 0/unset)

## Full Configuration Example

```go
func main() {
    flag.Parse()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Create custom logger
    logger := zerolog.New(os.Stderr).With().
        Str("plugin", "my-cost-plugin").
        Timestamp().
        Logger()

    // Create custom registry (for Supports() validation)
    registry := &MyRegistryLookup{}

    // Configure and start server
    if err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin:   &MyPlugin{},
        Port:     pluginsdk.ParsePortFlag(),
        Registry: registry,      // Optional: custom registry lookup
        Logger:   &logger,       // Optional: custom logger
        UnaryInterceptors: []grpc.UnaryServerInterceptor{
            // Optional: custom interceptors (after built-in tracing)
            myAuthInterceptor,
            myMetricsInterceptor,
        },
    }); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

## Environment Variables Reference

| Variable | Purpose | Default |
|----------|---------|---------|
| PULUMICOST_PLUGIN_PORT | gRPC server port | Ephemeral |
| PULUMICOST_LOG_LEVEL | Log verbosity (debug/info/warn/error) | info |
| PULUMICOST_LOG_FORMAT | Log format (json/text) | json |
| PULUMICOST_LOG_FILE | Log file path | stderr |
| PULUMICOST_TRACE_ID | Distributed trace correlation ID | (none) |
| PULUMICOST_TEST_MODE | Enable test mode (true/false) | false |

## Why PORT is Not Supported

The generic `PORT` environment variable is intentionally not supported to prevent multi-plugin
port conflicts. When finfocus-core spawns multiple plugins (e.g., aws-public + aws-ce), each
needs a unique port. Using `--port` flag allows the core to allocate distinct ports for each
plugin subprocess.

## Error Handling

`Serve()` returns errors for these conditions:

- **Port in use**: `"failed to listen: address already in use"`
- **Permission denied**: `"failed to listen: permission denied"`
- **Context canceled**: `context.Canceled` or `context.DeadlineExceeded`
- **gRPC failure**: Underlying gRPC server error

Always check the returned error:

```go
if err := pluginsdk.Serve(ctx, config); err != nil {
    if errors.Is(err, context.Canceled) {
        log.Println("Server shutdown gracefully")
    } else {
        log.Fatalf("Server error: %v", err)
    }
}
```

## Graceful Shutdown

When the context is canceled:

1. `Serve()` calls `grpcServer.GracefulStop()`
2. Existing in-flight requests complete
3. New connections are rejected
4. Server shuts down cleanly
5. `Serve()` returns `context.Canceled`

This ensures no requests are dropped during shutdown.
