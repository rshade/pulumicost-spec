# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **Go SDK** for the PulumiCost specification, providing a complete runtime library for implementing and testing
cost source plugins. The SDK consists of three main packages:

- **`pricing/`** - Domain types, validation, and billing mode enumerations
- **`proto/`** - Generated gRPC code from protobuf definitions (do not edit manually)
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
