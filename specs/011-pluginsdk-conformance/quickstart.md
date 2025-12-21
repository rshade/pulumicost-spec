# Quickstart: PluginSDK Conformance Testing Adapters

**Date**: 2025-11-30
**Feature**: 012-pluginsdk-conformance

## Overview

This guide shows how to use the conformance testing adapter functions to validate your
`pluginsdk.Plugin` implementation against the PulumiCost specification.

## Prerequisites

- Go 1.20 or later
- A plugin implementing `pluginsdk.Plugin` interface
- Import: `github.com/rshade/pulumicost-spec/sdk/go/pluginsdk`

## Basic Usage

### 1. Run Basic Conformance

Basic conformance validates core plugin functionality:

```go
package myplugin_test

import (
    "testing"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func TestPluginBasicConformance(t *testing.T) {
    plugin := NewMyPlugin() // Your Plugin implementation

    result, err := pluginsdk.RunBasicConformance(plugin)
    if err != nil {
        t.Fatalf("Conformance test error: %v", err)
    }

    if !result.Passed() {
        pluginsdk.PrintConformanceReport(t, result)
        t.Errorf("Basic conformance failed: %d/%d tests passed",
            result.Summary.Passed, result.Summary.Total)
    }
}
```

### 2. Run Standard Conformance

Standard conformance is recommended for production-ready plugins:

```go
func TestPluginStandardConformance(t *testing.T) {
    plugin := NewMyPlugin()

    result, err := pluginsdk.RunStandardConformance(plugin)
    if err != nil {
        t.Fatalf("Conformance test error: %v", err)
    }

    // Print report regardless of outcome for visibility
    pluginsdk.PrintConformanceReport(t, result)

    if result.LevelAchieved < pluginsdk.ConformanceLevelStandard {
        t.Errorf("Expected Standard conformance, achieved: %s",
            result.LevelAchievedStr)
    }
}
```

### 3. Run Advanced Conformance

Advanced conformance validates high-performance requirements:

```go
func TestPluginAdvancedConformance(t *testing.T) {
    plugin := NewMyPlugin()

    result, err := pluginsdk.RunAdvancedConformance(plugin)
    if err != nil {
        t.Fatalf("Conformance test error: %v", err)
    }

    pluginsdk.PrintConformanceReport(t, result)

    // Advanced is optional - log result without failing
    t.Logf("Conformance level achieved: %s", result.LevelAchievedStr)
}
```

## Complete Test File Example

```go
package myplugin_test

import (
    "testing"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

// TestConformance runs the full conformance suite
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
        if result.Summary.Failed > 0 {
            t.Fail()
        }
    })

    t.Run("Advanced", func(t *testing.T) {
        result, err := pluginsdk.RunAdvancedConformance(plugin)
        if err != nil {
            t.Fatalf("Error: %v", err)
        }
        pluginsdk.PrintConformanceReport(t, result)
        // Advanced is informational - don't fail
        t.Logf("Level: %s", result.LevelAchievedStr)
    })
}
```

## Understanding Results

### ConformanceResult Fields

| Field | Description |
|-------|-------------|
| `PluginName` | Name returned by your plugin |
| `LevelAchieved` | Highest conformance level passed |
| `Summary.Total` | Total number of tests run |
| `Summary.Passed` | Tests that passed |
| `Summary.Failed` | Tests that failed |
| `Categories` | Results grouped by test category |

### Conformance Levels

| Level | Description | When to Target |
|-------|-------------|----------------|
| Basic | Core functionality | Minimum for any plugin |
| Standard | Production readiness | Production deployments |
| Advanced | High performance | Performance-critical use cases |

### Test Categories

| Category | What It Tests |
|----------|---------------|
| `spec_validation` | PricingSpec schema compliance |
| `rpc_correctness` | RPC method behavior |
| `performance` | Latency and memory benchmarks |
| `concurrency` | Parallel request handling |

## Error Handling

```go
result, err := pluginsdk.RunBasicConformance(nil)
if err != nil {
    // Error: "plugin cannot be nil"
    t.Fatalf("Setup error: %v", err)
}
```

## Migration from Direct sdk/go/testing Usage

If you were previously using `sdk/go/testing` directly:

**Before** (manual conversion):

```go
import plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"

plugin := NewMyPlugin()
server := pluginsdk.NewServer(plugin)  // Manual conversion
result, err := plugintesting.RunStandardConformance(server)
```

**After** (using adapters):

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"

plugin := NewMyPlugin()
result, err := pluginsdk.RunStandardConformance(plugin)  // Direct!
```

## Running Tests

```bash
# Run all conformance tests
go test -v -run TestConformance ./...

# Run specific level
go test -v -run TestConformance/Basic ./...
go test -v -run TestConformance/Standard ./...
go test -v -run TestConformance/Advanced ./...

# With race detection (recommended for concurrency tests)
go test -v -race -run TestConformance ./...
```

## Next Steps

1. Start with Basic conformance to validate core functionality
2. Progress to Standard for production readiness
3. Optionally target Advanced for performance-critical scenarios
4. Address any failed tests using the detailed error messages in the report
