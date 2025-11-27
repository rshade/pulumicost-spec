# pluginsdk - PulumiCost Plugin Development SDK

The `pluginsdk` package provides a comprehensive development SDK for building PulumiCost plugins. It includes
the core plugin interface, helper utilities for cost calculations, structured logging with zerolog, and testing
utilities for plugin development.

## Installation

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
```

## Quick Start

### 1. Implement the Plugin Interface

```go
package main

import (
    "context"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

type MyPlugin struct {
    *pluginsdk.BasePlugin
}

func NewMyPlugin() *MyPlugin {
    plugin := &MyPlugin{
        BasePlugin: pluginsdk.NewBasePlugin("my-plugin"),
    }

    // Configure supported providers and resource types
    plugin.Matcher().AddProvider("aws")
    plugin.Matcher().AddResourceType("aws:ec2:Instance")

    return plugin
}

// Override GetProjectedCost to implement your pricing logic
func (p *MyPlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    resource := req.GetResource()

    // Check if resource is supported
    if !p.Matcher().Supports(resource) {
        return nil, pluginsdk.NotSupportedError(resource)
    }

    // Calculate cost (example: $0.10/hour for t2.micro)
    hourlyRate := 0.10

    return p.Calculator().CreateProjectedCostResponse(
        "USD",
        hourlyRate,
        "On-demand EC2 instance pricing",
    ), nil
}
```

### 2. Serve Your Plugin

```go
func main() {
    plugin := NewMyPlugin()

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    if err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: plugin,
        Port:   50051, // Optional: use 0 for automatic port
    }); err != nil {
        log.Fatal(err)
    }
}
```

## Core Components

### Plugin Interface

The `Plugin` interface defines the core methods that every PulumiCost plugin must implement:

```go
type Plugin interface {
    Name() string
    GetProjectedCost(ctx context.Context, req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error)
    GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error)
    GetPricingSpec(ctx context.Context, req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error)
    EstimateCost(ctx context.Context, req *pbc.EstimateCostRequest) (*pbc.EstimateCostResponse, error)
}
```

### BasePlugin

`BasePlugin` provides a scaffold with default implementations for all methods. Extend it and override
the methods you need:

```go
plugin := pluginsdk.NewBasePlugin("my-plugin")

// Access helpers
matcher := plugin.Matcher()      // ResourceMatcher for filtering resources
calc := plugin.Calculator()      // CostCalculator for cost computations
name := plugin.Name()            // Returns "my-plugin"
```

### ResourceMatcher

Helps filter which resources your plugin supports:

```go
matcher := pluginsdk.NewResourceMatcher()
matcher.AddProvider("aws")
matcher.AddProvider("azure")
matcher.AddResourceType("aws:ec2:Instance")
matcher.AddResourceType("azure:compute:VirtualMachine")

// Check if a resource is supported
if matcher.Supports(resource) {
    // Process the resource
}
```

**Thread Safety**: ResourceMatcher is NOT safe for concurrent use. Configure it during plugin
initialization before calling `Serve()`.

### CostCalculator

Provides utilities for cost calculations:

```go
calc := pluginsdk.NewCostCalculator()

// Convert hourly to monthly (using 730 hours/month)
monthly := calc.HourlyToMonthly(0.10)  // Returns 73.0

// Convert monthly to hourly
hourly := calc.MonthlyToHourly(73.0)   // Returns 0.10

// Create standard response
resp := calc.CreateProjectedCostResponse("USD", 0.10, "Hourly pricing")
```

### Constants

```go
pluginsdk.HoursPerMonth  // 730.0 - standard hours for monthly calculations
```

## Structured Logging

The SDK uses zerolog for structured logging with standardized field names:

```go
import "github.com/rs/zerolog"

// Create a configured plugin logger
logger := pluginsdk.NewPluginLogger(
    "my-plugin",           // Plugin name
    "v1.0.0",              // Version
    zerolog.InfoLevel,     // Log level
    os.Stderr,             // Output writer (nil for os.Stderr)
)

// Standard field constants
logger.Info().
    Str(pluginsdk.FieldTraceID, traceID).
    Str(pluginsdk.FieldOperation, "GetProjectedCost").
    Str(pluginsdk.FieldResourceType, "aws:ec2:Instance").
    Float64(pluginsdk.FieldCostMonthly, 73.0).
    Msg("calculated cost")
```

### Available Field Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `FieldTraceID` | `trace_id` | Request trace identifier |
| `FieldComponent` | `component` | System component |
| `FieldOperation` | `operation` | RPC operation name |
| `FieldDurationMs` | `duration_ms` | Operation duration |
| `FieldResourceURN` | `resource_urn` | Pulumi resource URN |
| `FieldResourceType` | `resource_type` | Resource type |
| `FieldProvider` | `provider` | Cloud provider |
| `FieldRegion` | `region` | Cloud region |
| `FieldPluginName` | `plugin_name` | Plugin identifier |
| `FieldPluginVersion` | `plugin_version` | Plugin version |
| `FieldCostMonthly` | `cost_monthly` | Monthly cost value |
| `FieldAdapter` | `adapter` | Adapter name |
| `FieldErrorCode` | `error_code` | Error code |

### Trace ID Propagation

The SDK includes a gRPC interceptor for trace ID propagation:

```go
// Interceptor automatically extracts/generates trace IDs
// Already integrated into Serve() function

// Access trace ID in your handlers
traceID := pluginsdk.TraceIDFromContext(ctx)

// Manually inject trace ID (for testing)
ctx = pluginsdk.ContextWithTraceID(ctx, "abc123def456...")

// Generate a new trace ID
traceID, err := pluginsdk.GenerateTraceID()
```

### Operation Timing

```go
done := pluginsdk.LogOperation(logger, "GetProjectedCost")
defer done()  // Logs operation completion with duration
```

## Server Configuration

### ServeConfig Options

```go
pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
    Plugin:   plugin,           // Required: your Plugin implementation
    Port:     50051,            // Optional: 0 for PORT env or ephemeral
    Registry: customRegistry,   // Optional: for Supports() validation
    Logger:   &customLogger,    // Optional: custom zerolog.Logger
})
```

### Port Resolution

1. If `Port > 0`, uses that port
2. If `PORT` environment variable is set, uses that
3. Otherwise, uses an ephemeral port

The selected port is printed to stdout as `PORT=<port>`.

## Testing Utilities

The SDK provides testing utilities for plugin development:

### TestPlugin (Quick Testing)

```go
func TestMyPlugin(t *testing.T) {
    plugin := NewMyPlugin()

    // Creates an in-process gRPC server with cleanup
    tp := pluginsdk.NewTestPlugin(t, plugin)

    // Test name
    tp.TestName("my-plugin")

    // Test projected cost (expectError=false for success)
    resource := pluginsdk.CreateTestResource("aws", "aws:ec2:Instance", nil)
    resp := tp.TestProjectedCost(resource, false)

    // Validate response
    if resp.GetCostPerMonth() != 73.0 {
        t.Errorf("unexpected cost: %f", resp.GetCostPerMonth())
    }
}
```

### TestServer (Lower-Level Control)

```go
func TestWithServer(t *testing.T) {
    plugin := NewMyPlugin()
    ts := pluginsdk.NewTestServer(t, plugin)
    defer ts.Close()

    // Use the gRPC client directly
    client := ts.Client()
    resp, err := client.Name(ctx, &pbc.NameRequest{})
    // ...
}
```

### CreateTestResource

```go
// Create a test resource with properties
resource := pluginsdk.CreateTestResource(
    "aws",
    "aws:ec2:Instance",
    map[string]string{
        "instance_type": "t2.micro",
        "region":        "us-east-1",
    },
)
```

### Comprehensive Testing

For more comprehensive testing including conformance tests, mock plugins with error injection,
and performance benchmarks, use the `sdk/go/testing` package:

```go
import plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

// Run conformance tests
result := plugintesting.RunStandardConformanceTests(t, plugin)
plugintesting.PrintConformanceReport(t, result)
```

## Error Helpers

```go
// Resource not supported
err := pluginsdk.NotSupportedError(resource)
// Returns: "resource type X from provider Y is not supported"

// No cost data available
err := pluginsdk.NoDataError("resource-id")
// Returns: "no cost data available for resource resource-id"
```

## Manifest Management

Load and save plugin manifests:

```go
// Load from YAML or JSON
manifest, err := pluginsdk.LoadManifest("plugin-manifest.yaml")

// Save to file (format determined by extension)
err := pluginsdk.SaveManifest("plugin-manifest.json", manifest)
```

## Migration from pulumicost-core

If you're migrating from `pulumicost-core`, update your imports:

```go
// Before
import "github.com/rshade/pulumicost-core/pluginsdk"

// After
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
```

Key changes:

1. **Logging**: Now uses zerolog instead of slog
2. **Server.logger**: Changed from `*slog.Logger` to `zerolog.Logger`
3. **HoursPerMonth**: Now exported (was `hoursPerMonth`)
4. **ServeConfig.Logger**: New field for custom logger injection
5. **TracingUnaryServerInterceptor**: Now automatically integrated in `Serve()`

## API Reference

### Types

| Type | Description |
|------|-------------|
| `Plugin` | Core interface for plugin implementations |
| `SupportsProvider` | Optional interface for Supports() capability |
| `BasePlugin` | Scaffold with default implementations |
| `ResourceMatcher` | Resource filtering by provider/type |
| `CostCalculator` | Cost calculation utilities |
| `Server` | gRPC server wrapper |
| `ServeConfig` | Configuration for Serve() |
| `TestServer` | Testing server with cleanup |
| `TestPlugin` | High-level testing utilities |
| `ValidationErrors` | Multiple validation errors |

### Functions

| Function | Description |
|----------|-------------|
| `NewServer(plugin)` | Create server with default registry |
| `NewServerWithRegistry(plugin, registry)` | Create server with custom registry |
| `NewServerWithOptions(plugin, registry, logger)` | Create server with all options |
| `Serve(ctx, config)` | Start gRPC server |
| `NewPluginLogger(name, version, level, writer)` | Create configured logger |
| `TracingUnaryServerInterceptor()` | gRPC interceptor for trace IDs |
| `TraceIDFromContext(ctx)` | Extract trace ID from context |
| `ContextWithTraceID(ctx, traceID)` | Inject trace ID into context |
| `GenerateTraceID()` | Generate new trace ID |
| `LogOperation(logger, operation)` | Log operation with timing |
| `NotSupportedError(resource)` | Create not-supported error |
| `NoDataError(resourceID)` | Create no-data error |
| `LoadManifest(path)` | Load manifest from file |
| `SaveManifest(path, manifest)` | Save manifest to file |
| `NewTestServer(t, plugin)` | Create test server |
| `NewTestPlugin(t, plugin)` | Create test plugin helper |
| `CreateTestResource(provider, type, props)` | Create test resource |
