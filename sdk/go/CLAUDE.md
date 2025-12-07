# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **Go SDK** for the PulumiCost specification, providing a complete runtime library for implementing and testing
cost source plugins. The SDK consists of six main packages:

- **`currency/`** - ISO 4217 currency validation and metadata with zero-allocation validation
- **`pluginsdk/`** - Plugin development SDK with environment variable handling and gRPC server utilities
- **`pricing/`** - Domain types, validation, and billing mode enumerations
- **`proto/`** - Generated gRPC code from protobuf definitions (do not edit manually)
- **`registry/`** - Plugin registry domain types with optimized zero-allocation validation
- **`testing/`** - Comprehensive testing framework with harness, mocks, and conformance tests

## Build Commands

### Testing

```bash
# Run all tests from repository root
cd ../../ && go test ./...

# Run tests for specific packages (from this directory)
go test ./pricing
go test ./testing
go test -v ./testing -run TestConformance  # Conformance tests only

# Run benchmarks
go test -bench=. -benchmem ./testing/
go test -bench=BenchmarkAllMethods -benchmem ./testing/
```

### Development

```bash
# Build all packages
go build ./...

# Check imports and format
go mod tidy
go fmt ./...

# Run from root directory for full validation
cd ../../ && make test && make lint
```

## Architecture

### Package Structure

**`currency/` Package - ISO 4217 Currency Validation**

- `currency.go` - Currency struct, complete ISO 4217 data (180+ currencies), and metadata functions
- `validate.go` - Zero-allocation IsValid() function for currency code validation
- `doc.go` - Comprehensive package documentation with usage examples
- `currency_test.go` - Table-driven tests covering all validation scenarios
- `benchmark_test.go` - Performance benchmarks targeting <15 ns/op, 0 allocs/op
- Performance: <15 ns/op, 0 B/op, 0 allocs/op for validation
- Pattern: Package-level slice variables for zero-allocation validation (follows registry pattern)

**`registry/` Package - Plugin Registry Domain Types**

- `domain.go` - 8 enum types with optimized zero-allocation validation
- `domain_test.go` - Comprehensive tests and performance benchmarks
- Enum types: Provider, DiscoverySource, PluginStatus, SecurityLevel, InstallationMethod, PluginCapability,
  SystemPermission, AuthMethod
- Performance: 5-12 ns/op, 0 allocs/op across all validation functions
- Pattern: Package-level slice variables for zero-allocation validation

**`pluginsdk/` Package - Plugin Development SDK**

- `sdk.go` - gRPC server setup with `Serve()` function and `ServeConfig` options
- `env.go` - Centralized environment variable handling for all PulumiCost plugins
- `env_test.go` - Comprehensive tests for environment variable functions
- `tracing.go` - Distributed tracing utilities with `TracingUnaryServerInterceptor()`
- `logging.go` - Structured logging helpers with zerolog integration
- `metrics.go` - Prometheus metrics instrumentation for plugins

**Environment Variables (env.go)**:

| Constant | Variable Name | Purpose |
|----------|---------------|---------|
| `EnvPort` | `PULUMICOST_PLUGIN_PORT` | gRPC server port (no fallback) |
| `EnvLogLevel` | `PULUMICOST_LOG_LEVEL` | Log verbosity (debug, info, warn, error) |
| `EnvLogLevelFallback` | `LOG_LEVEL` | Legacy fallback for log level |
| `EnvLogFormat` | `PULUMICOST_LOG_FORMAT` | Log output format (json, text) |
| `EnvLogFile` | `PULUMICOST_LOG_FILE` | Log file path (empty = stdout) |
| `EnvTraceID` | `PULUMICOST_TRACE_ID` | Distributed tracing correlation ID |
| `EnvTestMode` | `PULUMICOST_TEST_MODE` | Enable test mode (only "true" enables) |

**Environment Functions**:

- `GetPort() int` - Returns port from `PULUMICOST_PLUGIN_PORT` or 0 if not set/invalid
- `GetLogLevel() string` - Returns log level (canonical first, then fallback)
- `GetLogFormat() string` - Returns log format or empty string
- `GetLogFile() string` - Returns log file path or empty string
- `GetTraceID() string` - Returns trace ID or empty string
- `GetTestMode() bool` - Returns true only if value is "true" (logs warning for invalid values)
- `IsTestMode() bool` - Same as GetTestMode but without warning logging (for repeated checks)

**`pricing/` Package - Domain Logic**

- `domain.go` - Comprehensive BillingMode and Provider enumerations (44+ billing modes)
- `validate.go` - JSON Schema validation with embedded schema for PricingSpec documents
- Core domain types: `BillingMode`, `Provider` with validation functions

**`proto/` Package - Generated gRPC Code**

- Generated from `../../proto/pulumicost/v1/costsource.proto`
- Contains all message types: `NameRequest/Response`, `SupportsRequest/Response`, etc.
- `CostSourceServiceServer/Client` interfaces for gRPC communication
- **Never edit manually** - regenerated via `make generate` from repo root

**`testing/` Package - Plugin Testing Framework**

- `harness.go` - In-memory gRPC test harness using `bufconn` for fast, isolated tests
- `mock_plugin.go` - Configurable mock plugin with error injection and custom behaviors
- `integration_test.go` - Comprehensive integration tests for all RPC methods
- `conformance_test.go` - Multi-level conformance validation (Basic/Standard/Advanced)
- `benchmark_test.go` - Performance benchmarks with memory profiling

### Key Design Patterns

**Billing Mode System**
The `pricing` package defines 44+ billing modes organized by category:

- Time-based: `per_hour`, `per_minute`, `per_second`, etc.
- Storage-based: `per_gb_month`, `per_gb_hour`, etc.
- Usage-based: `per_request`, `per_operation`, etc.
- Compute-based: `per_cpu_hour`, `per_memory_gb_hour`, etc.
- Database-specific: `per_rcu`, `per_wcu`, `per_dtu`, etc.
- Pricing models: `on_demand`, `reserved`, `spot`, etc.

**Schema Validation**:

- Embedded JSON schema in `validate.go` for runtime PricingSpec validation
- `ValidatePricingSpec([]byte) error` function for document validation
- Schema synchronized with `../../schemas/pricing_spec.schema.json`

**Testing Framework Architecture**:

- **TestHarness**: In-memory gRPC server using `bufconn.Listener` for fast, isolated testing
- **MockPlugin**: Configurable mock with error injection, delays, and custom responses
- **Conformance Testing**: Three-tier validation (Basic/Standard/Advanced) with performance requirements
- **Performance Benchmarks**: Memory-profiled benchmarks for all RPC methods

### Enum Validation Pattern (Registry Package)

The registry package implements **optimized zero-allocation validation** for all enum types:

**Pattern**: Package-level slice variables

```go
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allProviders = []Provider{ProviderAWS, ProviderAzure, ProviderGCP, ProviderKubernetes, ProviderCustom}

func AllProviders() []Provider {
    return allProviders  // Zero allocation
}

func IsValidProvider(p string) bool {
    provider := Provider(p)
    for _, valid := range allProviders {  // Direct slice access
        if provider == valid {
            return true
        }
    }
    return false
}
```

**Performance**:

- 5-12 ns/op across all 8 enum types
- 0 B/op, 0 allocs/op (zero allocation)
- 2x faster than map-based alternatives
- Memory footprint: ~608 bytes total for all enums

**Documentation**: See `../specs/001-domain-enum-optimization/validation-pattern.md` for complete pattern guide.

**Status**: Registry package fully optimized ✅, currency package fully optimized ✅, pricing package pending future optimization

### Plugin Implementation Flow

1. **Implement** `proto.CostSourceServiceServer` interface
2. **Test** using `TestHarness` with your implementation
3. **Validate** using conformance tests (`RunBasicConformanceTests`, etc.)
4. **Benchmark** performance using provided benchmark suite
5. **Integrate** with validation using `pricing.ValidatePricingSpec`

## Common Development Tasks

### Creating a New Plugin

```go
// Implement the gRPC interface
type MyPlugin struct {
    proto.UnimplementedCostSourceServiceServer
}

// Test with harness
func TestMyPlugin(t *testing.T) {
    plugin := &MyPlugin{}
    harness := testing.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    // Run conformance tests
    result := testing.RunBasicConformanceTests(t, plugin)
    if result.FailedTests > 0 {
        t.Errorf("Conformance failed: %s", result.Summary)
    }
}
```

### Validating PricingSpec Documents

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pricing"

// Validate JSON document
if err := pricing.ValidatePricingSpec(jsonData); err != nil {
    return fmt.Errorf("invalid pricing spec: %w", err)
}
```

### Currency Validation

```go
import "github.com/rshade/pulumicost-spec/sdk/go/currency"

// Validate currency codes
if !currency.IsValid("USD") {
    return errors.New("invalid currency")
}

// Get currency metadata
usd, err := currency.GetCurrency("USD")
if err != nil {
    return err
}
fmt.Printf("%s uses %d decimal places\n", usd.Name, usd.MinorUnits)

// List all currencies
for _, c := range currency.AllCurrencies() {
    fmt.Printf("%s: %s\n", c.Code, c.Name)
}
```

### Using Environment Variables

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"

// Get plugin port (returns 0 if not set)
port := pluginsdk.GetPort()
if port == 0 {
    port = 8080 // use default
}

// Get log configuration with fallback support
logLevel := pluginsdk.GetLogLevel()  // checks PULUMICOST_LOG_LEVEL, then LOG_LEVEL
logFormat := pluginsdk.GetLogFormat() // empty string if not set
logFile := pluginsdk.GetLogFile()     // empty string means stdout

// Get trace ID for distributed tracing
traceID := pluginsdk.GetTraceID()
if traceID != "" {
    // Include in logs and responses
}

// Check test mode (use IsTestMode for repeated checks to avoid log spam)
if pluginsdk.IsTestMode() {
    // Enable test-specific behavior
}
```

### Running Specific Test Suites

```bash
# Integration tests only
go test -v ./testing -run TestIntegration

# Performance benchmarks
go test -bench=. ./testing

# Conformance tests
go test -v ./testing -run TestConformance

# Single RPC method tests
go test -v ./testing -run TestName
```

### Mock Plugin Configuration

```go
// Create configurable mock for testing
mock := testing.NewMockPlugin()
mock.SetNameResponse("test-plugin")
mock.SetSupportsResponse("aws", true)
mock.InjectError("GetActualCost", errors.New("simulated error"))

harness := testing.NewTestHarness(mock)
// Test error handling scenarios
```

## Performance Requirements

- **Response Times**: Name() < 100ms, others vary by complexity
- **Concurrency**: Must handle 10+ concurrent requests (Standard conformance)
- **Memory**: Efficient memory usage tracked via benchmarks
- **Consistency**: Consistent responses across multiple calls

## Code Generation

The `proto/` directory contains generated code. To regenerate:

```bash
cd ../../ && make generate
```

Never edit generated files directly - modify the source `.proto` files instead.
