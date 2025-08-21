# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **testing package** of the PulumiCost Go SDK, providing a comprehensive testing framework for plugin
implementations. The package enables in-memory gRPC testing, conformance validation, performance benchmarking, and integration
testing for CostSource plugins.

## Core Components

### Test Harness System (`harness.go`)

**In-Memory gRPC Testing**:

- **TestHarness**: Main testing framework using `bufconn.Listener` for isolated tests
- **1MB buffer size** for high-throughput testing scenarios  
- **Automatic server lifecycle management**: Start(), Stop(), Client() methods
- **Validation constants**: Plugin name limits (100 chars), currency codes (3 chars), response times

**Key Functions**:

- `NewTestHarness(impl)` - Create harness with plugin implementation
- `Start(t)` - Initialize gRPC client connection using bufconn
- `Client()` - Get gRPC client for making test requests
- **Validation suite**: `ValidateNameResponse`, `ValidateSupportsResponse`, `ValidateActualCostResponse`, etc.

### Mock Plugin System (`mock_plugin.go`)

**Configurable Mock Implementation**:

- **3 specialized mock types**: `NewMockPlugin()`, `ConfigurableErrorMockPlugin()`, `SlowMockPlugin()`
- **Cross-provider support**: AWS, Azure, GCP, Kubernetes with realistic resource mappings
- **Error injection**: Configurable error behavior for each RPC method
- **Response delays**: Artificial delays for timeout/performance testing (100-500ms ranges)
- **Dynamic data generation**: Cost variation patterns, provider-specific pricing multipliers

**Mock Configuration**:

- **Supported resources**: 13 resource types across 4 providers
- **Cost simulation**: Base rates with realistic variation patterns (±40%)
- **Metric hints**: Provider-appropriate usage metrics (vcpu_hours, storage_gb, invocations)
- **Billing modes**: Automatic selection based on resource type (per_hour, per_gb_month, per_invocation)

### Integration Test Suite (`integration_test.go`)

**Comprehensive RPC Testing**:

- **Basic functionality tests**: All 5 RPC methods with validation
- **Error condition testing**: Network errors, invalid inputs, edge cases  
- **Input validation**: Nil resources, invalid time ranges, missing parameters
- **Multi-provider testing**: Cross-provider consistency validation
- **Concurrency testing**: 10 concurrent requests with race condition detection
- **Data consistency testing**: Multiple calls returning consistent results

**Test Categories**:

- `TestBasicPluginFunctionality` - Core RPC method validation
- `TestErrorHandling` - Error injection and gRPC status code verification
- `TestInputValidation` - Invalid input handling and boundary conditions
- `TestMultipleProviders` - Cross-provider support validation
- `TestConcurrentRequests` - Thread safety and race condition testing
- `TestResponseTimeouts` - Configurable delay and timeout behavior
- `TestDataConsistency` - Idempotent response validation

### Conformance Testing (`conformance_test.go`)

**Three-Tier Validation System**:

- **Basic Conformance**: Core functionality required for all plugins (4 tests)
- **Standard Conformance**: Production-readiness validation (extends Basic + 3 tests)  
- **Advanced Conformance**: High-performance requirements (extends Standard + 3 tests)

**Conformance Levels**:

**Basic** (Required for all plugins):

- Name validation and response format
- Supports handling for valid/invalid resources
- ProjectedCost basic functionality

**Standard** (Production deployment):

- ActualCost time range validation  
- Error handling with proper gRPC codes
- PricingSpec consistency across calls
- 24-hour data handling capability

**Advanced** (High-performance environments):

- Performance baseline (Name < 100ms)
- Concurrent request handling (10+ requests)
- Large dataset support (30-day queries < 10s)

**Usage Patterns**:

```go
// Run conformance tests
result := RunBasicConformanceTests(t, plugin)
PrintConformanceReport(t, result)

// Command-line conformance testing
ConformanceTestMain(plugin, ConformanceStandard)
```

### Performance Testing (`benchmark_test.go`)

**Comprehensive Benchmark Suite**:

- **Individual method benchmarks**: All 5 RPC methods with memory profiling
- **Combined workload benchmarks**: `BenchmarkAllMethods` for realistic usage patterns  
- **Concurrency benchmarks**: `BenchmarkConcurrentRequests` with parallel execution
- **Data size benchmarks**: 1 hour to 30 days with varying dataset sizes
- **Provider-specific benchmarks**: Cross-provider performance comparison

**Performance Test Suite**:

- `PerformanceTestSuite` - Standardized performance measurement framework
- **Iteration counts**: 100 iterations (50 for expensive operations)
- **Measurement precision**: Min/avg/max duration tracking with memory profiling
- **Cross-provider testing**: AWS, Azure, GCP, Kubernetes performance comparison

## Build Commands

### Testing from Repository Root

```bash
# Run all tests (recommended)
cd ../../../ && make test

# Run comprehensive validation
cd ../../../ && make validate

# Run linting with Go and buf
cd ../../../ && make lint
```

### Testing from This Directory

```bash
# Run all testing package tests
go test

# Run specific test categories
go test -run TestBasicPluginFunctionality
go test -run TestConformance
go test -run TestErrorHandling

# Run with coverage analysis
go test -cover
go test -coverprofile=coverage.out && go tool cover -html=coverage.out

# Run integration tests with verbose output
go test -v -run TestMultipleProviders
```

### Benchmark Execution

```bash
# Run all benchmarks
go test -bench=.

# Run benchmarks with memory profiling
go test -bench=. -benchmem

# Run specific benchmark categories
go test -bench=BenchmarkName
go test -bench=BenchmarkAllMethods
go test -bench=BenchmarkConcurrentRequests

# Run data size benchmarks
go test -bench=BenchmarkActualCostDataSizes

# Run provider-specific benchmarks
go test -bench=BenchmarkDifferentProviders
```

## Architecture Patterns

### In-Memory Testing Strategy

The framework uses **bufconn-based in-memory gRPC** for isolated, fast testing:

1. **Isolated Testing**: Each test gets dedicated server instance
2. **No Network Dependencies**: Pure in-memory communication via bufconn
3. **Deterministic Behavior**: Consistent test execution without network variability
4. **High Performance**: Sub-millisecond setup/teardown for rapid test cycles

### Mock Plugin Architecture

**Configurable Behavior System**:

- **Error Injection**: Per-method error configuration with realistic gRPC status codes
- **Response Delays**: Configurable delays mimicking real network/processing latency
- **Data Generation**: Deterministic cost data with realistic variation patterns
- **Provider Simulation**: Authentic cross-provider behavior differences

### Conformance Test Design

**Progressive Validation Approach**:

- **Hierarchical Requirements**: Each level builds on previous requirements
- **Failure Fast**: Higher levels fail immediately if lower levels don't pass  
- **Detailed Reporting**: Comprehensive test result analysis and failure diagnostics
- **Extensible Framework**: Easy addition of new conformance tests

### Performance Testing Strategy

**Multi-Dimensional Benchmarking**:

- **Method-Level**: Individual RPC performance characteristics
- **Workload-Level**: Realistic usage pattern simulation
- **Concurrency-Level**: Thread safety and scalability validation  
- **Data-Level**: Performance scaling with dataset size

## Common Development Patterns

### Plugin Testing Workflow

**Standard Testing Pattern**:

```go
func TestMyPlugin(t *testing.T) {
    plugin := &MyPluginImpl{}
    
    // Basic integration testing
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()
    
    // Individual method testing
    client := harness.Client()
    resp, err := client.Name(ctx, &pbc.NameRequest{})
    // ... validation
}
```

**Conformance Validation**:

```go
func TestPluginConformance(t *testing.T) {
    plugin := &MyPluginImpl{}
    
    // Run standard conformance tests
    result := plugintesting.RunStandardConformanceTests(t, plugin)
    plugintesting.PrintConformanceReport(t, result)
    
    if result.FailedTests > 0 {
        t.Errorf("Plugin failed conformance: %s", result.Summary)
    }
}
```

### Mock-Driven Development

**Error Condition Testing**:

```go
func TestErrorHandling(t *testing.T) {
    plugin := plugintesting.ConfigurableErrorMockPlugin()
    
    // Configure specific error behavior
    plugin.ShouldErrorOnName = true
    
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()
    
    // Test error conditions
    _, err := client.Name(ctx, &pbc.NameRequest{})
    // ... verify error handling
}
```

**Performance Testing**:

```go
func TestPluginPerformance(t *testing.T) {
    plugin := &MyPluginImpl{}
    suite := plugintesting.NewPerformanceTestSuite(plugin)
    suite.RunPerformanceTests(t)
}
```

### Custom Test Extensions

**Adding New Conformance Tests**:

```go
func addCustomConformanceTests(suite *plugintesting.PluginConformanceSuite) {
    suite.AddTest(plugintesting.ConformanceTest{
        Name:        "CustomValidation",
        Description: "Custom validation logic",
        TestFunc:    createCustomValidationTest(),
    })
}
```

## Performance Requirements

### Response Time Baselines

- **Name()**: < 100ms (Advanced conformance)
- **Supports()**: < 50ms (Standard), < 25ms (Advanced)  
- **GetProjectedCost()**: < 200ms (Standard), < 100ms (Advanced)
- **GetPricingSpec()**: < 200ms (Standard), < 100ms (Advanced)
- **GetActualCost()**: < 2s for 24h data, < 10s for 30d data

### Concurrency Requirements

- **Standard**: Handle 10 concurrent requests safely
- **Advanced**: Handle 50+ concurrent requests with consistent performance

### Memory Requirements

- **Efficient memory usage** for normal operations
- **Large dataset handling** without memory leaks or excessive allocation
- **Bounded memory growth** for long-running operations

## Key Design Decisions

### Bufconn vs Network Testing

- **In-memory advantages**: Speed, isolation, determinism  
- **Trade-offs**: No real network conditions, simplified error scenarios
- **Best practice**: Use bufconn for unit/integration, real network for system tests

### Mock Plugin Flexibility

- **Configuration over inheritance**: Behavior modification through properties
- **Realistic simulation**: Provider-specific differences and authentic error conditions
- **Test data generation**: Deterministic yet varied cost data for comprehensive testing

### Conformance Test Hierarchy

- **Progressive complexity**: Basic → Standard → Advanced validation levels
- **Clear requirements**: Explicit performance and functionality thresholds
- **Production readiness**: Standard level indicates deployment readiness

### Performance Test Comprehensiveness

- **Multi-dimensional coverage**: Methods, workloads, concurrency, data sizes
- **Realistic scenarios**: Cross-provider testing with authentic resource configurations  
- **Actionable metrics**: Min/avg/max measurements with memory profiling

## Test Execution Examples

### Full Plugin Validation

```bash
# Complete validation pipeline
cd ../../../ && make test      # Unit tests
go test -v -run TestConformance # Conformance validation
go test -bench=.               # Performance benchmarks
cd ../../../ && make lint      # Code quality
```

### Targeted Testing

```bash
# Focus on specific functionality
go test -v -run TestBasicPluginFunctionality
go test -v -run TestMultipleProviders
go test -bench=BenchmarkName

# Performance analysis
go test -bench=BenchmarkAllMethods -benchmem
go test -bench=BenchmarkActualCostDataSizes
```

### CI/CD Integration

```bash
# Conformance validation for CI
go test -v -run TestConformance
# Performance regression testing  
go test -bench=. -benchtime=10s
# Coverage analysis
go test -cover -coverprofile=coverage.out
```
