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

| Constant             | Value            | Description              |
| -------------------- | ---------------- | ------------------------ |
| `FieldTraceID`       | `trace_id`       | Request trace identifier |
| `FieldComponent`     | `component`      | System component         |
| `FieldOperation`     | `operation`      | RPC operation name       |
| `FieldDurationMs`    | `duration_ms`    | Operation duration       |
| `FieldResourceURN`   | `resource_urn`   | Pulumi resource URN      |
| `FieldResourceType`  | `resource_type`  | Resource type            |
| `FieldProvider`      | `provider`       | Cloud provider           |
| `FieldRegion`        | `region`         | Cloud region             |
| `FieldPluginName`    | `plugin_name`    | Plugin identifier        |
| `FieldPluginVersion` | `plugin_version` | Plugin version           |
| `FieldCostMonthly`   | `cost_monthly`   | Monthly cost value       |
| `FieldAdapter`       | `adapter`        | Adapter name             |
| `FieldErrorCode`     | `error_code`     | Error code               |

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

## Prometheus Metrics

The SDK provides optional Prometheus metrics instrumentation for monitoring plugin performance.

### Quick Start

Add the metrics interceptor to your plugin:

```go
config := pluginsdk.ServeConfig{
    Plugin: plugin,
    UnaryInterceptors: []grpc.UnaryServerInterceptor{
        pluginsdk.MetricsUnaryServerInterceptor("my-plugin"),
    },
}
```

### Exposing Metrics

#### Option A: Built-in HTTP Server (Simple)

```go
server, err := pluginsdk.StartMetricsServer(pluginsdk.MetricsServerConfig{
    Port: 9090,  // Default: 9090
})
if err != nil {
    log.Fatal(err)
}
defer server.Shutdown(context.Background())
```

Then scrape `http://localhost:9090/metrics`.

#### Option B: Custom Registry (Production)

```go
metrics := pluginsdk.NewPluginMetrics("my-plugin")
interceptor := pluginsdk.MetricsInterceptorWithRegistry(metrics)

// Add to your existing HTTP server
http.Handle("/metrics", promhttp.HandlerFor(
    metrics.Registry,
    promhttp.HandlerOpts{},
))
```

### Available Metrics

| Metric                                       | Type      | Labels                                    | Description                  |
| -------------------------------------------- | --------- | ----------------------------------------- | ---------------------------- |
| `pulumicost_plugin_requests_total`           | Counter   | `grpc_method`, `grpc_code`, `plugin_name` | Total gRPC requests          |
| `pulumicost_plugin_request_duration_seconds` | Histogram | `grpc_method`, `plugin_name`              | Request latency distribution |

**Histogram Buckets**: 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s

### Example PromQL Queries

```promql
# Request rate by method
sum(rate(pulumicost_plugin_requests_total[5m])) by (grpc_method)

# Error rate
sum(rate(pulumicost_plugin_requests_total{grpc_code!="OK"}[5m]))
/ sum(rate(pulumicost_plugin_requests_total[5m]))

# P99 latency by method
histogram_quantile(0.99, sum(rate(pulumicost_plugin_request_duration_seconds_bucket[5m])) by (le, grpc_method))
```

### Performance

Benchmark results show <1% overhead for realistic handler workloads:

```text
BenchmarkMetricsInterceptor_Overhead    1000000    1391 ns/op    0 B/op    0 allocs/op
BenchmarkMetricsInterceptor_NoMetrics   830516698  1.82 ns/op    0 B/op    0 allocs/op
```

For handlers with typical work (1ms+), the metrics overhead is under 1% of total request time.

### Metrics Functions

| Function                                    | Description                              |
| ------------------------------------------- | ---------------------------------------- |
| `NewPluginMetrics(pluginName)`              | Create metrics with custom registry      |
| `MetricsUnaryServerInterceptor(pluginName)` | Create interceptor with default registry |
| `MetricsInterceptorWithRegistry(metrics)`   | Create interceptor with custom registry  |
| `StartMetricsServer(config)`                | Start optional HTTP metrics server       |

### Metrics Constants

| Constant             | Value        | Description          |
| -------------------- | ------------ | -------------------- |
| `MetricNamespace`    | `pulumicost` | Prometheus namespace |
| `MetricSubsystem`    | `plugin`     | Prometheus subsystem |
| `DefaultMetricsPort` | `9090`       | Default HTTP port    |
| `DefaultMetricsPath` | `/metrics`   | Default URL path     |

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
2. If `PULUMICOST_PLUGIN_PORT` environment variable is set, uses that
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

### Conformance Testing

The pluginsdk package provides adapter functions for running conformance tests directly on your
`Plugin` implementation without manual conversion to the gRPC server interface:

```go
func TestPluginConformance(t *testing.T) {
    plugin := NewMyPlugin()

    // Run basic conformance (core functionality)
    result, err := pluginsdk.RunBasicConformance(plugin)
    if err != nil {
        t.Fatalf("Conformance test error: %v", err)
    }

    // Print formatted report
    pluginsdk.PrintConformanceReport(t, result)

    if !result.Passed() {
        t.Errorf("Basic conformance failed: %d/%d tests passed",
            result.Summary.Passed, result.Summary.Total)
    }
}
```

**Conformance Levels**:

| Level                            | Description                                | Use Case                       |
| -------------------------------- | ------------------------------------------ | ------------------------------ |
| `RunBasicConformance(plugin)`    | Core functionality, required for all       | Minimum validation             |
| `RunStandardConformance(plugin)` | Production-ready (includes error handling) | Production deployments         |
| `RunAdvancedConformance(plugin)` | High performance (strict latency limits)   | Performance-critical scenarios |

**Type Aliases**:

The package re-exports key types from `sdk/go/testing` for convenience:

```go
// Use directly from pluginsdk
var level pluginsdk.ConformanceLevel = pluginsdk.ConformanceLevelStandard
var result *pluginsdk.ConformanceResult
var summary pluginsdk.ResultSummary
```

**Complete Example**:

```go
func TestConformance(t *testing.T) {
    plugin := NewMyPlugin()

    t.Run("Basic", func(t *testing.T) {
        result, err := pluginsdk.RunBasicConformance(plugin)
        if err != nil {
            t.Fatalf("Error: %v", err)
        }
        if result.Summary.Failed > 0 {
            pluginsdk.PrintConformanceReport(t, result)
            t.Fail()
        }
    })

    t.Run("Standard", func(t *testing.T) {
        result, err := pluginsdk.RunStandardConformance(plugin)
        if err != nil {
            t.Fatalf("Error: %v", err)
        }
        pluginsdk.PrintConformanceReport(t, result)
        if result.LevelAchieved < pluginsdk.ConformanceLevelStandard {
            t.Errorf("Expected Standard conformance, achieved: %s",
                result.LevelAchievedStr)
        }
    })
}
```

### Advanced Testing with sdk/go/testing

For more advanced testing scenarios including mock plugins with error injection,
custom configurations, and performance benchmarks, use the `sdk/go/testing` package directly:

```go
import plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

// Create configurable mock
mock := plugintesting.NewMockPlugin()
mock.ShouldErrorOnName = true  // Inject errors

// Convert plugin to server and use test harness for in-memory gRPC
plugin := NewMyPlugin()
server := pluginsdk.NewServer(plugin)
harness := plugintesting.NewTestHarness(server)
harness.Start(t)
defer harness.Stop()
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

## FOCUS 1.2 Cost Records

The SDK includes a comprehensive `FocusRecordBuilder` for constructing FinOps FOCUS 1.2 compliant cost records.
FOCUS (FinOps Open Cost and Usage Specification) is a standard for cloud billing data.

### Quick Start

```go
import (
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

builder := pluginsdk.NewFocusRecordBuilder()

// Set mandatory fields
builder.WithIdentity("AWS", "123456789012", "Production Account")
builder.WithBillingPeriod(billingStart, billingEnd, "USD")
builder.WithChargePeriod(chargeStart, chargeEnd)
builder.WithChargeDetails(
    pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
    pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
)
builder.WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2")
builder.WithFinancials(73.0, 80.0, 70.0, "USD", "INV-2025-11")
builder.WithContractedCost(65.0)  // New in FOCUS 1.2

record, err := builder.Build()
```

### FOCUS 1.2 Column Coverage

The `FocusCostRecord` proto message implements all 57 columns defined in FOCUS 1.2:

| Category             | Mandatory | Recommended | Conditional |
| -------------------- | --------- | ----------- | ----------- |
| Identity & Hierarchy | 3         | 0           | 4           |
| Billing Period       | 3         | 0           | 0           |
| Charge Period        | 2         | 0           | 0           |
| Charge Details       | 3         | 1           | 0           |
| Pricing              | 0         | 0           | 10          |
| Service & Product    | 1         | 0           | 3           |
| Resource Details     | 0         | 0           | 3           |
| SKU Details          | 0         | 0           | 4           |
| Location             | 0         | 0           | 3           |
| Financial Amounts    | 2         | 0           | 3           |
| Consumption/Usage    | 0         | 0           | 2           |
| Commitment Discounts | 0         | 0           | 7           |
| Capacity Reservation | 0         | 0           | 2           |
| Invoice              | 0         | 0           | 2           |
| Metadata             | 0         | 0           | 1           |
| **Total**            | **14**    | **1**       | **42**      |

### Builder Methods by Category

**Mandatory Fields (14 columns)**:

| Method                                                         | Fields Set                                            | FOCUS Section |
| -------------------------------------------------------------- | ----------------------------------------------------- | ------------- |
| `WithIdentity(provider, accountID, accountName)`               | ProviderName, BillingAccountId, BillingAccountName    | 2.1           |
| `WithBillingPeriod(start, end, currency)`                      | BillingPeriodStart, BillingPeriodEnd, BillingCurrency | 2.2           |
| `WithChargePeriod(start, end)`                                 | ChargePeriodStart, ChargePeriodEnd                    | 2.3           |
| `WithChargeDetails(chargeCat, pricingCat)`                     | ChargeCategory, PricingCategory                       | 2.4           |
| `WithChargeClassification(class, desc, freq)`                  | ChargeClass, ChargeDescription                        | 2.4           |
| `WithService(category, name)`                                  | ServiceCategory, ServiceName                          | 2.6           |
| `WithFinancials(billed, list, effective, currency, invoiceID)` | BilledCost, ListCost, EffectiveCost                   | 2.10          |
| `WithContractedCost(cost)`                                     | ContractedCost                                        | 3.20          |

**Conditional Fields - New in FOCUS 1.2**:

| Method                                                   | Fields Set                                       | FOCUS Section |
| -------------------------------------------------------- | ------------------------------------------------ | ------------- |
| `WithBillingAccountType(type)`                           | BillingAccountType                               | 3.3           |
| `WithSubAccountType(type)`                               | SubAccountType                                   | 3.45          |
| `WithCapacityReservation(id, status)`                    | CapacityReservationId, CapacityReservationStatus | 3.6, 3.7      |
| `WithCommitmentDiscountDetails(qty, status, type, unit)` | CommitmentDiscountQuantity, Status, Type, Unit   | 3.14-3.19     |
| `WithContractedUnitPrice(price)`                         | ContractedUnitPrice                              | 3.21          |
| `WithPricingCurrency(currency)`                          | PricingCurrency                                  | 3.34          |
| `WithPricingCurrencyPrices(contracted, effective, list)` | PricingCurrency\*Prices                          | 3.35-3.37     |
| `WithPublisher(publisher)`                               | Publisher                                        | 3.39          |
| `WithServiceSubcategory(subcategory)`                    | ServiceSubcategory                               | 3.43          |
| `WithSkuDetails(meter, priceDetails)`                    | SkuMeter, SkuPriceDetails                        | 3.46, 3.48    |

**Other Conditional Fields**:

| Method                                       | Fields Set                                  | FOCUS Section |
| -------------------------------------------- | ------------------------------------------- | ------------- |
| `WithSubAccount(id, name)`                   | SubAccountId, SubAccountName                | 2.1           |
| `WithLocation(regionID, regionName, az)`     | RegionId, RegionName, AvailabilityZone      | 2.9           |
| `WithResource(id, name, type)`               | ResourceId, ResourceName, ResourceType      | 2.7           |
| `WithSKU(skuID, skuPriceID)`                 | SkuId, SkuPriceId                           | 2.8           |
| `WithPricing(quantity, unit, listUnitPrice)` | PricingQuantity, PricingUnit, ListUnitPrice | 2.5           |
| `WithUsage(quantity, unit)`                  | ConsumedQuantity, ConsumedUnit              | 2.11          |
| `WithCommitmentDiscount(category, id, name)` | CommitmentDiscount\*                        | 2.12          |
| `WithInvoice(id, issuer)`                    | InvoiceId, InvoiceIssuer                    | 2.13          |
| `WithTag(key, value)` / `WithTags(map)`      | Tags                                        | 2.14          |
| `WithExtension(key, value)`                  | ExtendedColumns                             | 2.14          |

### Validation

The builder validates mandatory fields and business rules on `Build()`:

```go
record, err := builder.Build()
if err != nil {
    // Handle validation error
    // e.g., "billing_account_id is required"
}
```

**Validated Rules**:

- BillingAccountId is required
- ChargePeriod (start/end) is required
- ServiceCategory must be specified (not UNSPECIFIED)
- ChargeCategory must be specified (not UNSPECIFIED)
- BillingCurrency is required
- Usage records must have non-zero ConsumedQuantity

### Migration Guide

If you're updating from an earlier version without FOCUS 1.2 support:

**New Mandatory Field**: `ContractedCost` is now mandatory in FOCUS 1.2. Add it to all records:

```go
builder.WithContractedCost(65.0)
```

**New Builder Methods**: 11 new builder methods for FOCUS 1.2 columns:

- `WithContractedCost`
- `WithBillingAccountType`
- `WithSubAccountType`
- `WithCapacityReservation`
- `WithCommitmentDiscountDetails`
- `WithContractedUnitPrice`
- `WithPricingCurrency`
- `WithPricingCurrencyPrices`
- `WithPublisher`
- `WithServiceSubcategory`
- `WithSkuDetails`

### Troubleshooting

**Common Validation Errors**:

| Error                                          | Cause                                     | Solution                                         |
| ---------------------------------------------- | ----------------------------------------- | ------------------------------------------------ |
| `billing_account_id is required`               | Missing identity                          | Call `WithIdentity()`                            |
| `charge_period (start/end) is required`        | Missing charge period                     | Call `WithChargePeriod()`                        |
| `service_category is required`                 | Unspecified category                      | Call `WithService()` or `WithServiceCategory()`  |
| `charge_category is required`                  | Unspecified category                      | Call `WithChargeDetails()`                       |
| `billing_currency is required`                 | Missing currency                          | Call `WithBillingPeriod()` or `WithFinancials()` |
| `consumed_quantity must be positive for usage` | Zero or negative usage for Usage category | Set positive usage via `WithUsage()`             |

### Complete Example

See `examples/plugins/focus_example.go` for a comprehensive example demonstrating all 57 FOCUS 1.2 columns.

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

| Type                | Description                                  |
| ------------------- | -------------------------------------------- |
| `Plugin`            | Core interface for plugin implementations    |
| `SupportsProvider`  | Optional interface for Supports() capability |
| `BasePlugin`        | Scaffold with default implementations        |
| `ResourceMatcher`   | Resource filtering by provider/type          |
| `CostCalculator`    | Cost calculation utilities                   |
| `Server`            | gRPC server wrapper                          |
| `ServeConfig`       | Configuration for Serve()                    |
| `TestServer`        | Testing server with cleanup                  |
| `TestPlugin`        | High-level testing utilities                 |
| `ValidationErrors`  | Multiple validation errors                   |
| `ConformanceResult` | Result of conformance suite execution        |
| `ConformanceLevel`  | Conformance certification level              |
| `ResultSummary`     | Aggregate test counts                        |

### Functions

| Function                                         | Description                         |
| ------------------------------------------------ | ----------------------------------- |
| `NewServer(plugin)`                              | Create server with default registry |
| `NewServerWithRegistry(plugin, registry)`        | Create server with custom registry  |
| `NewServerWithOptions(plugin, registry, logger)` | Create server with all options      |
| `Serve(ctx, config)`                             | Start gRPC server                   |
| `NewPluginLogger(name, version, level, writer)`  | Create configured logger            |
| `TracingUnaryServerInterceptor()`                | gRPC interceptor for trace IDs      |
| `TraceIDFromContext(ctx)`                        | Extract trace ID from context       |
| `ContextWithTraceID(ctx, traceID)`               | Inject trace ID into context        |
| `GenerateTraceID()`                              | Generate new trace ID               |
| `LogOperation(logger, operation)`                | Log operation with timing           |
| `NotSupportedError(resource)`                    | Create not-supported error          |
| `NoDataError(resourceID)`                        | Create no-data error                |
| `LoadManifest(path)`                             | Load manifest from file             |
| `SaveManifest(path, manifest)`                   | Save manifest to file               |
| `NewTestServer(t, plugin)`                       | Create test server                  |
| `NewTestPlugin(t, plugin)`                       | Create test plugin helper           |
| `CreateTestResource(provider, type, props)`      | Create test resource                |
| `RunBasicConformance(plugin)`                    | Run basic conformance tests         |
| `RunStandardConformance(plugin)`                 | Run standard conformance tests      |
| `RunAdvancedConformance(plugin)`                 | Run advanced conformance tests      |
| `PrintConformanceReport(t, result)`              | Print formatted report to test log  |
| `PrintConformanceReportTo(result, writer)`       | Print formatted report to io.Writer |
