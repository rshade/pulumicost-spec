# Data Model: PluginSDK Conformance Testing Adapters

**Date**: 2025-11-30
**Feature**: 012-pluginsdk-conformance

## Overview

This feature adds adapter functions to the `pluginsdk` package. The data model primarily consists
of type aliases and function signatures that bridge the `Plugin` interface with the existing
conformance testing framework.

## Entities

### 1. Plugin (existing)

**Source**: `sdk/go/pluginsdk/sdk.go`
**Role**: Input to adapter functions

```go
type Plugin interface {
    Name() string
    GetProjectedCost(ctx, req) (*GetProjectedCostResponse, error)
    GetActualCost(ctx, req) (*GetActualCostResponse, error)
    GetPricingSpec(ctx, req) (*GetPricingSpecResponse, error)
    EstimateCost(ctx, req) (*EstimateCostResponse, error)
}
```

**Validation Rules**:

- MUST NOT be nil (adapter functions return error if nil)
- SHOULD implement all interface methods without panicking

### 2. Server (existing)

**Source**: `sdk/go/pluginsdk/sdk.go`
**Role**: Intermediate conversion type

```go
type Server struct {
    pbc.UnimplementedCostSourceServiceServer
    plugin   Plugin
    registry RegistryLookup
    logger   zerolog.Logger
}
```

**Relationships**:

- Created from Plugin via `NewServer(plugin)`
- Implements `pbc.CostSourceServiceServer` interface
- Passed to conformance functions

### 3. ConformanceResult (re-exported)

**Source**: `sdk/go/testing/conformance.go`
**Role**: Output of adapter functions (type alias)

```go
// Type alias in pluginsdk/conformance.go
type ConformanceResult = plugintesting.ConformanceResult

// Original definition
type ConformanceResult struct {
    Version          string                           // Report schema version
    Timestamp        time.Time                        // Execution timestamp
    PluginName       string                           // From plugin's Name() method
    LevelAchieved    ConformanceLevel                 // Highest level passed
    LevelAchievedStr string                           // String representation
    Summary          ResultSummary                    // Aggregate counts
    Categories       map[TestCategory]*CategoryResult // Results by category
    Duration         time.Duration                    // Total execution time
    DurationStr      string                           // String representation
}
```

**Validation Rules**:

- Version follows semantic versioning
- LevelAchieved is one of: Basic, Standard, Advanced
- Summary.Total = Summary.Passed + Summary.Failed + Summary.Skipped

### 4. ConformanceLevel (re-exported)

**Source**: `sdk/go/testing/conformance.go`
**Role**: Enumeration for certification levels (type alias)

```go
// Type alias in pluginsdk/conformance.go
type ConformanceLevel = plugintesting.ConformanceLevel

// Constants (from source package)
const (
    ConformanceLevelBasic    ConformanceLevel = iota  // Core functionality
    ConformanceLevelStandard                          // Production readiness
    ConformanceLevelAdvanced                          // High performance
)
```

## Function Signatures

### Adapter Functions (new)

```go
// RunBasicConformance runs basic conformance tests against a Plugin implementation.
// Returns error if plugin is nil.
func RunBasicConformance(plugin Plugin) (*ConformanceResult, error)

// RunStandardConformance runs standard conformance tests against a Plugin implementation.
// Returns error if plugin is nil.
func RunStandardConformance(plugin Plugin) (*ConformanceResult, error)

// RunAdvancedConformance runs advanced conformance tests against a Plugin implementation.
// Returns error if plugin is nil.
func RunAdvancedConformance(plugin Plugin) (*ConformanceResult, error)

// PrintConformanceReport prints a formatted conformance report to the test log.
func PrintConformanceReport(t *testing.T, result *ConformanceResult)
```

## Entity Relationships

```text
┌──────────────────────────────────────────────────────────────┐
│                      pluginsdk package                       │
│                                                              │
│  ┌─────────┐                                                 │
│  │ Plugin  │ ─────────────────────┐                          │
│  └────┬────┘                      │                          │
│       │                           ▼                          │
│       │ NewServer()        ┌─────────────┐                   │
│       └────────────────────│   Server    │                   │
│                            └──────┬──────┘                   │
│                                   │                          │
│                                   │ implements               │
│                                   ▼                          │
│                     ┌─────────────────────────┐              │
│                     │ CostSourceServiceServer │              │
│                     └────────────┬────────────┘              │
│                                  │                           │
└──────────────────────────────────┼───────────────────────────┘
                                   │
                                   │ passed to
                                   ▼
┌──────────────────────────────────────────────────────────────┐
│                      testing package                         │
│                                                              │
│  ┌─────────────────────┐         ┌───────────────────────┐   │
│  │ RunBasicConformance │─────────│                       │   │
│  │ RunStandardConf...  │─────────│  ConformanceResult    │   │
│  │ RunAdvancedConf...  │─────────│                       │   │
│  └─────────────────────┘         └───────────────────────┘   │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

## State Transitions

No stateful entities. All functions are stateless operations that:

1. Accept Plugin input
2. Convert to Server (if not nil)
3. Execute conformance tests
4. Return immutable ConformanceResult

## Data Volume / Scale Assumptions

- **Plugin implementations**: Typically 1 per test file
- **ConformanceResult size**: ~1-5 KB depending on number of test categories
- **Execution time**: 1-60 seconds depending on conformance level
