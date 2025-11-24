# Quickstart: Zerolog SDK Logging Utilities

**Feature**: 001-zerolog | **SDK Package**: `github.com/rshade/pulumicost-spec/sdk/go/pluginsdk`

## Installation

Add the SDK to your plugin's dependencies:

```bash
go get github.com/rshade/pulumicost-spec/sdk/go@latest
```

## Basic Usage

### 1. Create a Logger

```go
package main

import (
    "os"
    "github.com/rs/zerolog"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    // Create logger with plugin metadata
    logger := pluginsdk.NewPluginLogger(
        "aws-public",      // plugin name
        "v1.0.0",          // version
        zerolog.InfoLevel, // minimum log level
        os.Stderr,         // output (nil defaults to stderr)
    )

    logger.Info().Msg("Plugin started")
}
```

### 2. Configure gRPC Server with Tracing

```go
import (
    "google.golang.org/grpc"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    // ... logger setup ...

    server := grpc.NewServer(
        grpc.UnaryInterceptor(pluginsdk.TracingUnaryServerInterceptor()),
    )

    // Register your service
    pb.RegisterCostSourceServiceServer(server, &myPlugin{logger: logger})
}
```

### 3. Use Trace ID in Handlers

```go
func (p *myPlugin) GetProjectedCost(
    ctx context.Context,
    req *pb.GetProjectedCostRequest,
) (*pb.GetProjectedCostResponse, error) {

    // Extract trace ID from context
    traceID := pluginsdk.TraceIDFromContext(ctx)

    // Log with trace ID and standard fields
    p.logger.Info().
        Str(pluginsdk.FieldTraceID, traceID).
        Str(pluginsdk.FieldOperation, "GetProjectedCost").
        Str(pluginsdk.FieldResourceType, req.ResourceType).
        Msg("Processing cost request")

    // ... implementation ...
}
```

### 4. Time Operations

```go
func (p *myPlugin) GetProjectedCost(
    ctx context.Context,
    req *pb.GetProjectedCostRequest,
) (*pb.GetProjectedCostResponse, error) {

    // Start timing - logs duration when done() is called
    done := pluginsdk.LogOperation(p.logger, "GetProjectedCost")
    defer done()

    // ... implementation ...

    return &pb.GetProjectedCostResponse{...}, nil
}
```

## Complete Plugin Example

```go
package main

import (
    "context"
    "net"
    "os"

    "github.com/rs/zerolog"
    "google.golang.org/grpc"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pb "github.com/rshade/pulumicost-spec/sdk/go/proto"
)

type awsPlugin struct {
    pb.UnimplementedCostSourceServiceServer
    logger zerolog.Logger
}

func (p *awsPlugin) Name(
    ctx context.Context,
    req *emptypb.Empty,
) (*pb.NameResponse, error) {
    return &pb.NameResponse{Name: "aws-public"}, nil
}

func (p *awsPlugin) GetProjectedCost(
    ctx context.Context,
    req *pb.GetProjectedCostRequest,
) (*pb.GetProjectedCostResponse, error) {

    traceID := pluginsdk.TraceIDFromContext(ctx)
    done := pluginsdk.LogOperation(p.logger, "GetProjectedCost")
    defer done()

    p.logger.Info().
        Str(pluginsdk.FieldTraceID, traceID).
        Str(pluginsdk.FieldResourceType, req.ResourceType).
        Msg("Processing projected cost request")

    // Fetch pricing data
    cost, err := p.fetchPricing(ctx, req)
    if err != nil {
        p.logger.Error().
            Str(pluginsdk.FieldTraceID, traceID).
            Err(err).
            Str(pluginsdk.FieldErrorCode, "PRICING_FETCH_FAILED").
            Msg("Failed to fetch pricing")
        return nil, err
    }

    p.logger.Info().
        Str(pluginsdk.FieldTraceID, traceID).
        Float64(pluginsdk.FieldCostMonthly, cost.Monthly).
        Msg("Cost calculated successfully")

    return &pb.GetProjectedCostResponse{...}, nil
}

func main() {
    logger := pluginsdk.NewPluginLogger(
        "aws-public",
        "v1.0.0",
        zerolog.InfoLevel,
        os.Stderr,
    )

    logger.Info().Msg("Starting AWS public pricing plugin")

    server := grpc.NewServer(
        grpc.UnaryInterceptor(pluginsdk.TracingUnaryServerInterceptor()),
    )

    pb.RegisterCostSourceServiceServer(server, &awsPlugin{logger: logger})

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to listen")
    }

    logger.Info().Str("address", ":50051").Msg("Server listening")

    if err := server.Serve(lis); err != nil {
        logger.Fatal().Err(err).Msg("Server failed")
    }
}
```

## Standard Field Names

Use these constants for consistent field naming across all plugins:

| Constant | Field Name | Usage |
|----------|------------|-------|
| `FieldTraceID` | trace_id | Request correlation |
| `FieldOperation` | operation | RPC method name |
| `FieldDurationMs` | duration_ms | Operation timing |
| `FieldResourceType` | resource_type | Cloud resource type |
| `FieldResourceURN` | resource_urn | Pulumi resource ID |
| `FieldPluginName` | plugin_name | Auto-added by logger |
| `FieldPluginVersion` | plugin_version | Auto-added by logger |
| `FieldCostMonthly` | cost_monthly | Monthly cost value |
| `FieldAdapter` | adapter | Data source adapter |
| `FieldErrorCode` | error_code | Error classification |
| `FieldComponent` | component | System component |

## File-Based Logging

For production deployments with file output:

```go
// Open log file
logFile, err := os.OpenFile(
    "/var/log/pulumicost/aws-public.log",
    os.O_APPEND|os.O_CREATE|os.O_WRONLY,
    0644,
)
if err != nil {
    log.Fatal(err)
}
defer logFile.Close()

// Create logger with file output
logger := pluginsdk.NewPluginLogger(
    "aws-public",
    "v1.0.0",
    zerolog.InfoLevel,
    logFile,
)
```

Configure log rotation externally using logrotate or systemd.

## Log Levels

```go
zerolog.TraceLevel // -1: Most verbose
zerolog.DebugLevel // 0: Development details
zerolog.InfoLevel  // 1: Normal operations (recommended default)
zerolog.WarnLevel  // 2: Warnings
zerolog.ErrorLevel // 3: Errors
zerolog.FatalLevel // 4: Fatal errors (calls os.Exit)
zerolog.PanicLevel // 5: Panics
```

## Testing with Logging

```go
func TestMyHandler(t *testing.T) {
    // Capture logs in buffer for testing
    var buf bytes.Buffer
    logger := pluginsdk.NewPluginLogger(
        "test-plugin",
        "v0.0.1",
        zerolog.DebugLevel,
        &buf,
    )

    // Run handler with logger
    plugin := &myPlugin{logger: logger}
    resp, err := plugin.GetProjectedCost(ctx, req)

    // Assert on captured logs
    logOutput := buf.String()
    assert.Contains(t, logOutput, "Processing cost request")
}
```

## Best Practices

1. **Create logger once** at plugin startup, pass by value (loggers are
   immutable)
2. **Always include trace_id** when available for distributed tracing
3. **Use standard field constants** for ecosystem-wide log analysis
4. **Log at appropriate levels**: Debug for development, Info for operations
5. **Include operation context**: resource type, URN, adapter name
6. **Time expensive operations** using LogOperation for performance monitoring
7. **Use Err() for errors**: `logger.Error().Err(err).Msg("...")`
