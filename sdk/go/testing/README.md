# PulumiCost Plugin Testing Framework

This package provides a comprehensive testing framework for PulumiCost plugin
implementations, including integration tests, performance benchmarks, and conformance validation.

## Overview

The testing framework consists of several key components:

- **Test Harness**: In-memory gRPC testing environment using bufconn
- **Mock Plugin**: Configurable mock implementation for testing
- **Conformance Suite**: Four-category test suite with three conformance levels
- **Performance Benchmarks**: Standardized performance testing with latency baselines
- **Concurrency Tests**: Thread safety and parallel request handling validation

### Test Categories

The conformance suite organizes tests into four categories:

| Category            | Description                               | Example Tests                                               |
| ------------------- | ----------------------------------------- | ----------------------------------------------------------- |
| **Spec Validation** | JSON schema and data format compliance    | Billing mode enum, currency format, required fields         |
| **RPC Correctness** | Protocol behavior and response validation | Name response format, error handling, time range validation |
| **Performance**     | Latency thresholds and response time      | Method latency, baseline variance (< 10%)                   |
| **Concurrency**     | Thread safety and parallel execution      | Parallel requests, response consistency                     |

### Conformance Levels

| Level        | Description              | Requirements                                          |
| ------------ | ------------------------ | ----------------------------------------------------- |
| **Basic**    | Required for all plugins | Core functionality, spec compliance                   |
| **Standard** | Production-ready plugins | Reliability, consistency, 10 parallel requests        |
| **Advanced** | High-performance plugins | Strict latency, 50+ parallel requests, variance < 10% |

## Quick Start

### Basic Plugin Testing

```go
package main

import (
    "testing"

    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Assume you have a plugin implementation
func TestMyPlugin(t *testing.T) {
    plugin := &MyPluginImpl{} // Your implementation

    // Run basic integration tests
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    // Test individual methods
    client := harness.Client()
    ctx := context.Background()

    resp, err := client.Name(ctx, &pbc.NameRequest{})
    if err != nil {
        t.Fatalf("Name() failed: %v", err)
    }

    if err := plugintesting.ValidateNameResponse(resp); err != nil {
        t.Errorf("Invalid response: %v", err)
    }
}
```

### Running Conformance Tests

```go
func TestPluginConformance(t *testing.T) {
    plugin := &MyPluginImpl{}

    // Create conformance suite
    suite := plugintesting.NewConformanceSuite()

    // Register all test categories
    plugintesting.RegisterSpecValidationTests(suite)
    plugintesting.RegisterRPCCorrectnessTests(suite)
    plugintesting.RegisterPerformanceTests(suite)
    plugintesting.RegisterConcurrencyTests(suite)

    // Run tests at desired conformance level
    result := suite.Run(plugin, plugintesting.ConformanceLevelStandard)

    // Print detailed report
    plugintesting.PrintReportTo(result, os.Stdout)

    if result.FailedCount > 0 {
        t.Errorf("Plugin failed conformance tests")
    }
}

// Or use convenience functions
func TestPluginConformanceSimple(t *testing.T) {
    plugin := &MyPluginImpl{}

    // Run standard conformance (includes Basic + Standard tests)
    result := plugintesting.RunStandardConformance(plugin)
    plugintesting.PrintReportTo(result, os.Stdout)

    if result.FailedCount > 0 {
        t.Errorf("Plugin failed conformance: %s", result.Summary)
    }
}
```

### Performance Testing

```go
func TestPluginPerformance(t *testing.T) {
    plugin := &MyPluginImpl{}
    suite := plugintesting.NewPerformanceTestSuite(plugin)
    suite.RunPerformanceTests(t)
}

func BenchmarkMyPlugin(b *testing.B) {
    plugin := &MyPluginImpl{}
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(&testing.T{})
    defer harness.Stop()

    client := harness.Client()
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := client.Name(ctx, &pbc.NameRequest{})
        if err != nil {
            b.Fatalf("Name() failed: %v", err)
        }
    }
}
```

## Components

### Conformance Suite

The `ConformanceSuite` is the main entry point for plugin validation:

```go
// Create and configure suite
suite := plugintesting.NewConformanceSuite()

// Configure options
suite.WithConfig(plugintesting.SuiteConfig{
    Timeout:          60 * time.Second,
    ParallelRequests: 10,
})

// Register test categories
plugintesting.RegisterSpecValidationTests(suite)
plugintesting.RegisterRPCCorrectnessTests(suite)
plugintesting.RegisterPerformanceTests(suite)
plugintesting.RegisterConcurrencyTests(suite)

// Run at desired level
result := suite.Run(plugin, plugintesting.ConformanceLevelStandard)
```

**Convenience Functions:**

```go
// Run all tests at specific level
result := plugintesting.RunBasicConformance(plugin)
result := plugintesting.RunStandardConformance(plugin)
result := plugintesting.RunAdvancedConformance(plugin)
```

### Test Harness

The `TestHarness` provides an in-memory gRPC testing environment:

```go
type TestHarness struct {
    server   *grpc.Server
    listener *bufconn.Listener
    client   pbc.CostSourceServiceClient
    conn     *grpc.ClientConn
}
```

**Key Methods:**

- `NewTestHarness(impl)`: Create harness for plugin implementation
- `Start(t)`: Initialize client connection
- `Stop()`: Clean up resources
- `Client()`: Get gRPC client for testing

### Mock Plugin

The `MockPlugin` provides a configurable test implementation:

```go
// Basic mock plugin
plugin := plugintesting.NewMockPlugin()

// Configure supported providers/resources
plugin.SupportedProviders = []string{"aws", "azure"}
plugin.SupportedResources["aws"] = []string{"ec2", "s3"}

// Configure error behavior
plugin.ShouldErrorOnName = true

// Configure response delays
plugin.NameDelay = 100 * time.Millisecond
```

**Specialized Mock Plugins:**

- `ConfigurableErrorMockPlugin()`: For error testing
- `SlowMockPlugin()`: For timeout/performance testing

### Validation Functions

The framework provides comprehensive response validation:

```go
// Validate individual responses
err := plugintesting.ValidateNameResponse(nameResp)
err := plugintesting.ValidateSupportsResponse(supportsResp)
err := plugintesting.ValidateActualCostResponse(actualCostResp)
err := plugintesting.ValidateProjectedCostResponse(projectedResp)
err := plugintesting.ValidatePricingSpecResponse(specResp)

// Validate protobuf messages
err := plugintesting.ValidatePricingSpec(spec)
err := plugintesting.ValidateActualCostResult(result)
```

### Conformance Result

The `ConformanceResult` contains detailed test execution information:

```go
type ConformanceResult struct {
    PluginName   string          // Name of the tested plugin
    Level        ConformanceLevel // Tested conformance level
    PassedCount  int             // Number of passed tests
    FailedCount  int             // Number of failed tests
    SkippedCount int             // Number of skipped tests
    TotalTests   int             // Total number of tests
    Duration     time.Duration   // Total execution time
    Results      []TestResult    // Individual test results
    Summary      string          // Human-readable summary
}

type TestResult struct {
    Method   string        // RPC method tested
    Category TestCategory  // Test category
    Success  bool          // Pass/fail status
    Error    error         // Error if failed
    Duration time.Duration // Test execution time
    Details  string        // Additional information
}
```

**Reporting:**

```go
// Print to stdout
plugintesting.PrintReportTo(result, os.Stdout)

// Get JSON report
jsonReport := result.ToJSON()
```

### Helper Functions

```go
// Create test data
resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")
start, end := plugintesting.CreateTimeRange(24) // 24 hours ago to now

// Performance measurement
metrics, err := plugintesting.MeasurePerformance("TestName", 100, func() error {
    _, err := client.Name(ctx, &pbc.NameRequest{})
    return err
})
```

## Testing Levels

### 1. Basic Integration Tests

Tests core functionality of all RPC methods:

```go
func TestBasicPluginFunctionality(t *testing.T)
func TestErrorHandling(t *testing.T)
func TestInputValidation(t *testing.T)
func TestMultipleProviders(t *testing.T)
func TestConcurrentRequests(t *testing.T)
func TestDataConsistency(t *testing.T)
```

**Run with:**

```bash
go test -v -run TestBasicPluginFunctionality
```

### 2. Performance Benchmarks

Standardized performance tests for all RPC methods:

```go
func BenchmarkName(b *testing.B)
func BenchmarkSupports(b *testing.B)
func BenchmarkGetActualCost(b *testing.B)
func BenchmarkGetProjectedCost(b *testing.B)
func BenchmarkGetPricingSpec(b *testing.B)
func BenchmarkAllMethods(b *testing.B)
func BenchmarkConcurrentRequests(b *testing.B)
```

**Run with:**

```bash
go test -bench=. -benchmem
```

### 3. Conformance Tests

The conformance suite provides multi-level validation across four test categories.

#### Test Categories

**Spec Validation Tests:**

- `SpecValidation_ValidPricingSpec` - Schema-compliant response validation
- `SpecValidation_BillingModeEnum` - Valid billing mode enumeration
- `SpecValidation_RequiredFields` - Required field presence

**RPC Correctness Tests:**

- `RPCCorrectness_NameResponse` - Valid Name response format
- `RPCCorrectness_SupportsValidation` - Supports validation behavior
- `RPCCorrectness_ErrorHandling` - Proper gRPC error codes
- `RPCCorrectness_TimeRangeValidation` - Time range validation
- `RPCCorrectness_ConsistentResponses` - Response consistency

**Performance Tests:**

- `Performance_NameLatency` - Name RPC latency threshold
- `Performance_SupportsLatency` - Supports RPC latency threshold
- `Performance_GetProjectedCostLatency` - GetProjectedCost latency
- `Performance_GetPricingSpecLatency` - GetPricingSpec latency
- `Performance_BaselineVariance` - Variance within 10% (SC-003)

**Concurrency Tests:**

- `Concurrency_ParallelRequests` - Thread-safe parallel execution
- `Concurrency_ConsistentUnderLoad` - Consistent responses under load
- `Concurrency_NoRaceConditions` - Race condition detection

#### Conformance Levels

**Basic Conformance** (Required for all plugins):

- All Spec Validation tests
- Core RPC Correctness tests

**Standard Conformance** (Production-ready):

- All Basic tests
- All RPC Correctness tests
- Standard Performance tests (10 parallel requests)

**Advanced Conformance** (High-performance):

- All Standard tests
- All Performance tests (strict latency thresholds)
- All Concurrency tests (50+ parallel requests)

**Run conformance tests:**

```go
func TestConformance(t *testing.T) {
    plugin := &MyPluginImpl{}

    // Choose conformance level
    result := plugintesting.RunBasicConformance(plugin)
    // result := plugintesting.RunStandardConformance(plugin)
    // result := plugintesting.RunAdvancedConformance(plugin)

    plugintesting.PrintReportTo(result, os.Stdout)

    if result.FailedCount > 0 {
        t.Fatalf("Plugin failed conformance: %s", result.Summary)
    }
}
```

**Custom test selection:**

```go
func TestCustomConformance(t *testing.T) {
    plugin := &MyPluginImpl{}

    suite := plugintesting.NewConformanceSuite()

    // Register only specific categories
    plugintesting.RegisterSpecValidationTests(suite)
    plugintesting.RegisterRPCCorrectnessTests(suite)
    // Skip performance and concurrency tests

    result := suite.Run(plugin, plugintesting.ConformanceLevelBasic)
    plugintesting.PrintReportTo(result, os.Stdout)
}
```

## Test Requirements by RPC Method

### Name()

**Basic Requirements:**

- MUST return non-empty plugin name
- Name MUST be ≤ 100 characters
- MUST respond within 100ms (advanced conformance)

**Validation:**

```go
func ValidateNameResponse(response *pbc.NameResponse) error
```

### Supports()

**Basic Requirements:**

- MUST handle nil resource gracefully (return false or error)
- MUST provide reason when resource not supported
- MUST be consistent (same input → same output)

**Standard Requirements:**

- SHOULD respond within 50ms for supported providers
- MUST validate all resource descriptor fields

### GetActualCost()

**Basic Requirements:**

- MUST validate time range (end > start)
- MUST return valid cost data structure
- Cost values MUST be non-negative
- MUST include source in results

**Standard Requirements:**

- SHOULD handle missing data gracefully (return empty or appropriate error)
- SHOULD support reasonable time ranges (up to 30 days)
- MUST provide meaningful error messages

**Advanced Requirements:**

- SHOULD handle large datasets efficiently (30+ days)
- MUST respond within 10 seconds for large queries

### GetProjectedCost()

**Basic Requirements:**

- MUST return valid pricing data
- Unit price MUST be non-negative
- Currency MUST be valid 3-character ISO code
- MUST reject unsupported resources

**Standard Requirements:**

- SHOULD be consistent across calls
- SHOULD provide cost per month calculation
- MUST include billing detail information

### GetPricingSpec()

**Basic Requirements:**

- MUST return complete pricing specification
- MUST include all required fields (provider, resource_type, billing_mode, rate_per_unit, currency)
- Rate per unit MUST be non-negative
- Currency MUST be valid 3-character ISO code

**Standard Requirements:**

- SHOULD include metric hints for cost calculation
- SHOULD be consistent across calls
- SHOULD provide meaningful descriptions

## Error Handling Requirements

### Expected Error Codes

- `InvalidArgument`: Invalid input parameters
- `NotFound`: Resource or data not found
- `Unavailable`: Service temporarily unavailable
- `PermissionDenied`: Access denied to resource/data
- `ResourceExhausted`: Rate limit or quota exceeded
- `Internal`: Internal service error

### Error Testing

```go
// Test error conditions
plugin.ShouldErrorOnName = true
_, err := client.Name(ctx, &pbc.NameRequest{})
if err == nil {
    t.Error("Expected error")
}

// Check error code
st, ok := status.FromError(err)
if !ok {
    t.Error("Expected gRPC status error")
}
if st.Code() != codes.Internal {
    t.Errorf("Expected Internal error, got %v", st.Code())
}
```

## Performance Requirements

### Latency Baselines

The conformance suite validates against these latency thresholds:

| Method              | Standard | Advanced |
| ------------------- | -------- | -------- |
| Name()              | 100ms    | 50ms     |
| Supports()          | 50ms     | 25ms     |
| GetProjectedCost()  | 200ms    | 100ms    |
| GetPricingSpec()    | 200ms    | 100ms    |
| GetActualCost (24h) | 2000ms   | 1000ms   |
| GetActualCost (30d) | N/A      | 10000ms  |

### Variance Requirements (SC-003)

- Advanced conformance requires benchmark variance ≤ 10% from baseline
- Measured over multiple iterations (default: 50 iterations)
- Helps ensure consistent performance under load

### Concurrency Requirements

| Level    | Parallel Requests | Requirements                    |
| -------- | ----------------- | ------------------------------- |
| Standard | 10                | Thread-safe, no race conditions |
| Advanced | 50                | Consistent responses under load |

### Memory Requirements

- Should not consume excessive memory for normal operations
- Must handle large datasets (30+ days) without memory issues
- Allocations tracked via benchmark tests

## CI/CD Integration

### GitHub Actions

Add to your `.github/workflows/test.yml`:

```yaml
name: Plugin Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.21"

      - name: Run Integration Tests
        run: go test -v ./...

      - name: Run Benchmarks
        run: go test -bench=. -benchmem

      - name: Run Conformance Tests
        run: go test -v -run TestConformance
```

### Makefile Integration

```makefile
.PHONY: test-integration test-conformance test-performance

test-integration:
 go test -v -run TestBasic

test-conformance:
 go test -v -run TestConformance

test-performance:
 go test -bench=. -benchmem

test-all: test-integration test-conformance test-performance
```

## Best Practices

### Plugin Implementation

1. **Implement all RPC methods** even if you return "not supported" errors
2. **Validate inputs** and return appropriate error codes
3. **Handle concurrency** safely - avoid shared mutable state
4. **Use timeouts** for external service calls
5. **Implement caching** for frequently requested data
6. **Provide meaningful error messages** to help users debug issues

### Testing Implementation

1. **Test all supported providers and resource types**
2. **Test error conditions** and edge cases
3. **Test with realistic data sizes** (don't just test with 1-2 data points)
4. **Test concurrency** if your plugin will handle concurrent requests
5. **Test performance** under expected load

### Debugging Tips

1. **Use the mock plugin** to understand expected behavior
2. **Check validation errors** - they provide specific failure details
3. **Run conformance tests** to identify compliance issues
4. **Use performance tests** to identify bottlenecks
5. **Enable debug logging** in your plugin during testing

## Example Test Suite

```go
package myplugin_test

import (
    "os"
    "testing"

    "github.com/myplugin/internal"
    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

func TestMyPlugin(t *testing.T) {
    plugin := internal.NewMyPlugin()

    // Create test harness
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    // Test individual methods
    t.Run("Name", func(t *testing.T) {
        resp, err := harness.Client().Name(context.Background(), &pbc.NameRequest{})
        if err != nil {
            t.Fatalf("Name() failed: %v", err)
        }
        if err := plugintesting.ValidateNameResponse(resp); err != nil {
            t.Errorf("Invalid response: %v", err)
        }
    })
}

func TestMyPluginConformance(t *testing.T) {
    plugin := internal.NewMyPlugin()

    // Run full conformance suite at Standard level
    result := plugintesting.RunStandardConformance(plugin)

    // Print detailed report
    plugintesting.PrintReportTo(result, os.Stdout)

    if result.FailedCount > 0 {
        t.Errorf("Plugin failed conformance: %s", result.Summary)
    }
}

func TestMyPluginCustomConformance(t *testing.T) {
    plugin := internal.NewMyPlugin()

    // Create custom suite with specific tests
    suite := plugintesting.NewConformanceSuite()
    plugintesting.RegisterSpecValidationTests(suite)
    plugintesting.RegisterRPCCorrectnessTests(suite)

    result := suite.Run(plugin, plugintesting.ConformanceLevelBasic)

    if result.FailedCount > 0 {
        t.Errorf("Plugin failed: %d/%d tests failed", result.FailedCount, result.TotalTests)
    }
}

func BenchmarkMyPlugin(b *testing.B) {
    plugin := internal.NewMyPlugin()
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(&testing.T{})
    defer harness.Stop()

    client := harness.Client()
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = client.Name(ctx, &pbc.NameRequest{})
    }
}
```

## Conformance Report Example

Running conformance tests produces a detailed report:

```text
================================================================================
                     PulumiCost Plugin Conformance Report
================================================================================

Plugin:     my-cost-plugin
Level:      Standard
Duration:   1.234s

--------------------------------------------------------------------------------
                                Test Results
--------------------------------------------------------------------------------

Category: spec_validation
  ✓ SpecValidation_ValidPricingSpec                     [PASS]     12ms
  ✓ SpecValidation_BillingModeEnum                      [PASS]      8ms
  ✓ SpecValidation_RequiredFields                       [PASS]      5ms

Category: rpc_correctness
  ✓ RPCCorrectness_NameResponse                         [PASS]     15ms
  ✓ RPCCorrectness_SupportsValidation                   [PASS]     22ms
  ✓ RPCCorrectness_ErrorHandling                        [PASS]     18ms

Category: performance
  ✓ Performance_NameLatency                             [PASS]     45ms
  ✓ Performance_SupportsLatency                         [PASS]     32ms

--------------------------------------------------------------------------------
                                   Summary
--------------------------------------------------------------------------------

Total:    8 tests
Passed:   8
Failed:   0
Skipped:  0

Result:   PASS - Plugin conforms to Standard level
================================================================================
```

This testing framework ensures your plugin meets the PulumiCost specification requirements and
provides a reliable, performant cost source for users.
