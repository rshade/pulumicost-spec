# PulumiCost Plugin Testing Framework

This package provides a comprehensive testing framework for PulumiCost plugin
implementations, including integration tests, performance benchmarks, and conformance validation.

## Overview

The testing framework consists of several key components:

- **Test Harness**: In-memory gRPC testing environment
- **Mock Plugin**: Configurable mock implementation for testing
- **Integration Tests**: Comprehensive RPC method validation
- **Performance Benchmarks**: Standardized performance testing
- **Conformance Tests**: Multi-level plugin validation

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
    
    // Run standard conformance tests
    result := plugintesting.RunStandardConformanceTests(t, plugin)
    plugintesting.PrintConformanceReport(result)
    
    if result.FailedTests > 0 {
        t.Errorf("Plugin failed conformance tests: %s", result.Summary)
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

Multi-level validation for plugin certification:

#### Basic Conformance

- All plugins MUST pass these tests
- Validates core functionality and error handling
- Required for plugin submission

#### Standard Conformance  

- Production-ready plugins should pass these tests
- Includes data consistency and reliability tests
- Recommended for enterprise deployments

#### Advanced Conformance

- High-performance plugins should pass these tests  
- Includes performance, concurrency, and scalability tests
- Required for high-throughput environments

**Run conformance tests:**

```go
// In your test file
func TestConformance(t *testing.T) {
    plugin := &MyPluginImpl{}
    
    // Choose conformance level
    result := plugintesting.RunBasicConformanceTests(t, plugin)
    // result := plugintesting.RunStandardConformanceTests(t, plugin)  
    // result := plugintesting.RunAdvancedConformanceTests(t, plugin)
    
    plugintesting.PrintConformanceReport(result)
    
    if result.FailedTests > 0 {
        t.Fatalf("Plugin failed conformance: %s", result.Summary)
    }
}
```

**Command-line conformance testing:**

```go
// main.go
package main

import (
    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

func main() {
    plugin := &MyPluginImpl{}
    plugintesting.ConformanceTestMain(plugin, plugintesting.ConformanceStandard)
}
```

```bash
go run main.go
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

### Response Time Requirements

- **Name()**: < 100ms (advanced)
- **Supports()**: < 50ms (standard), < 25ms (advanced)
- **GetProjectedCost()**: < 200ms (standard), < 100ms (advanced)
- **GetPricingSpec()**: < 200ms (standard), < 100ms (advanced)
- **GetActualCost()**: < 2s for 24h data (standard), < 10s for 30d data (advanced)

### Concurrency Requirements

- **Standard**: Must handle 10 concurrent requests
- **Advanced**: Must handle 50+ concurrent requests safely

### Memory Requirements

- Should not consume excessive memory for normal operations
- Must handle large datasets (30+ days) without memory issues

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
        go-version: '1.21'
    
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
    "testing"
    
    "github.com/myplugin/internal"
    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

func TestMyPlugin(t *testing.T) {
    plugin := internal.NewMyPlugin()
    
    // Basic integration tests
    plugintesting.TestBasicPluginFunctionality(t, plugin)
    plugintesting.TestErrorHandling(t, plugin) 
    plugintesting.TestInputValidation(t, plugin)
}

func TestMyPluginConformance(t *testing.T) {
    plugin := internal.NewMyPlugin()
    result := plugintesting.RunStandardConformanceTests(t, plugin)
    
    if result.FailedTests > 0 {
        t.Errorf("Plugin failed conformance: %s", result.Summary)
        plugintesting.PrintConformanceReport(result)
    }
}

func TestMyPluginPerformance(t *testing.T) {
    plugin := internal.NewMyPlugin()
    suite := plugintesting.NewPerformanceTestSuite(plugin)
    suite.RunPerformanceTests(t)
}

func BenchmarkMyPlugin(b *testing.B) {
    plugin := internal.NewMyPlugin()
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(&testing.T{})
    defer harness.Stop()
    
    // Add your benchmarks here
}
```

This testing framework ensures your plugin meets the PulumiCost specification requirements and
provides a reliable, performant cost source for users.
