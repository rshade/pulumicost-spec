# Quickstart: Plugin Conformance Test Suite

**Feature**: 011-plugin-conformance-suite

## Overview

The Plugin Conformance Test Suite validates that your CostSource plugin implementation meets the
pulumicost-spec requirements. Run these tests to ensure your plugin is ready for production
deployment.

## Installation

Add the testing package to your plugin's test dependencies:

```bash
go get github.com/rshade/pulumicost-spec/sdk/go/testing
```

## Basic Usage (< 10 lines)

```go
package myplugin_test

import (
    "testing"

    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
    "github.com/yourorg/myplugin"
)

func TestConformance(t *testing.T) {
    plugin := myplugin.New()
    result, err := plugintesting.RunStandardConformance(plugin)
    if err != nil {
        t.Fatalf("Suite execution failed: %v", err)
    }
    if !result.Passed() {
        plugintesting.PrintReport(result)
        t.Fatalf("Plugin failed conformance: %s", result.LevelAchieved)
    }
}
```

## Running Tests

```bash
# Run conformance tests
go test -v -run TestConformance

# Run with race detector (for concurrency validation)
go test -race -run TestConformance

# Run benchmarks
go test -bench=. -benchmem
```

## Conformance Levels

### Basic (Required)

All plugins MUST pass Basic conformance:

- Name validation and response format
- Supports handling for valid/invalid resources
- GetProjectedCost basic functionality
- GetPricingSpec schema compliance

```go
result, _ := plugintesting.RunBasicConformance(plugin)
```

### Standard (Recommended)

Production-ready plugins SHOULD pass Standard conformance:

- All Basic tests
- Error handling with proper gRPC codes
- Data consistency across calls
- 24-hour data handling
- 10 concurrent request handling

```go
result, _ := plugintesting.RunStandardConformance(plugin)
```

### Advanced (Optional)

High-performance plugins MAY pass Advanced conformance:

- All Standard tests
- Latency thresholds (Name < 50ms, etc.)
- 50 concurrent request handling
- 30-day data queries < 10s

```go
result, _ := plugintesting.RunAdvancedConformance(plugin)
```

## Custom Configuration

```go
suite := plugintesting.NewConformanceSuiteWithConfig(plugintesting.SuiteConfig{
    TargetLevel:      plugintesting.ConformanceLevelAdvanced,
    Timeout:          120 * time.Second,
    ParallelRequests: 50,
    EnableBenchmarks: true,
})

result, err := suite.Run(plugin)
```

## Running Individual Categories

```go
// Run only spec validation tests
result, err := suite.RunCategory(plugin, plugintesting.CategorySpecValidation)

// Run only performance benchmarks
result, err := suite.RunCategory(plugin, plugintesting.CategoryPerformance)
```

## JSON Report Output

For CI/CD integration, generate a JSON report:

```go
result, _ := plugintesting.RunStandardConformance(plugin)
jsonBytes, _ := result.ToJSON()
os.WriteFile("conformance-report.json", jsonBytes, 0644)
```

Example output:

```json
{
  "version": "1.0.0",
  "timestamp": "2025-11-28T12:00:00Z",
  "plugin_name": "my-aws-plugin",
  "level_achieved": "Standard",
  "summary": {
    "total": 15,
    "passed": 15,
    "failed": 0,
    "skipped": 0
  },
  "categories": {
    "spec_validation": { "passed": 4, "failed": 0 },
    "rpc_correctness": { "passed": 5, "failed": 0 },
    "performance": { "passed": 3, "failed": 0 },
    "concurrency": { "passed": 3, "failed": 0 }
  },
  "duration": "45.2s"
}
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Plugin Conformance

on: [push, pull_request]

jobs:
  conformance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"

      - name: Run Conformance Tests
        run: go test -v -run TestConformance

      - name: Run with Race Detector
        run: go test -race -run TestConformance

      - name: Generate Report
        run: |
          go test -v -run TestConformance -json > conformance-report.json
```

### Makefile

```makefile
.PHONY: conformance conformance-basic conformance-advanced

conformance:
    go test -v -run TestConformance

conformance-basic:
    go test -v -run TestConformance/Basic

conformance-advanced:
    go test -race -v -run TestConformance
    go test -bench=. -benchmem
```

## Troubleshooting

### Common Failures

#### Plugin name is empty

- Ensure your `Name()` RPC returns a non-empty string

#### Unsupported resource should have a reason

- When `Supports()` returns false, provide a reason in the response

#### Currency should be 3-character ISO code

- Use standard ISO 4217 codes: USD, EUR, GBP, etc.

#### Cost cannot be negative

- All cost values must be >= 0

### Race Condition Detection

If tests pass normally but fail with `-race`:

```bash
# Run with verbose race output
GORACE="history_size=7" go test -race -run TestConformance
```

### Performance Baseline Failures

If your plugin is slower than expected, check:

1. External API calls are properly cached
2. Large datasets use streaming or pagination
3. No unnecessary allocations in hot paths

## Next Steps

1. Run Basic conformance to validate core functionality
2. Fix any failures before proceeding
3. Run Standard conformance for production readiness
4. Integrate into CI/CD pipeline
5. Optionally target Advanced for high-performance environments
