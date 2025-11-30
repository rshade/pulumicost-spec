# Research: Plugin Conformance Test Suite

**Feature**: 011-plugin-conformance-suite
**Date**: 2025-11-28

## Research Topics

### 1. Existing Testing Framework Analysis

**Decision**: Extend existing `sdk/go/testing/` package rather than creating new package

**Rationale**:

- Existing TestHarness, MockPlugin, and validation functions are mature and tested
- ConformanceSuite type already exists in harness.go with AddTest/RunTests methods
- conformance_test.go already implements Basic/Standard/Advanced level pattern
- Extending maintains backward compatibility with existing plugin test code

**Alternatives Considered**:

- New `sdk/go/conformance/` package - rejected (would duplicate harness infrastructure)
- Standalone CLI tool - rejected (spec requires library import pattern)

### 2. JSON Report Structure Best Practices

**Decision**: Use structured ConformanceReport type with nested results

**Rationale**:

- Go's encoding/json provides native marshaling with struct tags
- Hierarchical structure matches test category organization
- Enables easy CI/CD integration (jq queries, GitHub Actions parsing)
- Matches existing TestResult struct pattern in harness.go

**Report Schema**:

```json
{
  "version": "1.0.0",
  "timestamp": "2025-11-28T12:00:00Z",
  "plugin_name": "example-plugin",
  "level_achieved": "Standard",
  "summary": {
    "total_tests": 15,
    "passed": 14,
    "failed": 1,
    "skipped": 0
  },
  "categories": {
    "spec_validation": { ... },
    "rpc_correctness": { ... },
    "performance": { ... },
    "concurrency": { ... }
  },
  "details": [ ... ]
}
```

**Alternatives Considered**:

- YAML output - rejected (JSON more universal for CI/CD)
- Plain text only - rejected (not programmatically parseable)

### 3. Performance Baseline Thresholds

**Decision**: Reference existing thresholds from sdk/go/testing/README.md and harness.go constants

**Rationale**:

- Thresholds already defined and documented:
  - `MaxResponseTimeMs = 100` (Name RPC)
  - `MaxLargeQueryTimeSeconds = 10` (30-day queries)
  - `NumConcurrentRequests = 10` (Standard concurrency)
- Avoids duplication and potential drift
- Plugin developers already familiar with these values

**Threshold Matrix** (from existing documentation):

| RPC Method         | Basic | Standard | Advanced |
| ------------------ | ----- | -------- | -------- |
| Name()             | N/A   | < 100ms  | < 50ms   |
| Supports()         | N/A   | < 50ms   | < 25ms   |
| GetProjectedCost() | N/A   | < 200ms  | < 100ms  |
| GetPricingSpec()   | N/A   | < 200ms  | < 100ms  |
| GetActualCost(24h) | N/A   | < 2s     | < 1s     |
| GetActualCost(30d) | N/A   | N/A      | < 10s    |
| Concurrency        | N/A   | 10 req   | 50 req   |

### 4. Race Condition Detection Integration

**Decision**: Use `go test -race` flag with explicit concurrency test functions

**Rationale**:

- Go's race detector is the standard tool (FR-010)
- Suite provides concurrent request functions that trigger detection
- No custom race detection needed - leverage existing tooling
- Document usage pattern in quickstart guide

**Integration Pattern**:

```go
// Test concurrency with race detector enabled
// go test -race -run TestConcurrency

func TestConcurrency(t *testing.T) {
    result := RunConcurrencyTests(t, plugin, ConcurrencyConfig{
        ParallelRequests: 10,
        Duration: 5 * time.Second,
    })
}
```

### 5. Conformance Level Hierarchy

**Decision**: Progressive levels where each includes all tests from lower levels

**Rationale**:

- Matches existing pattern in conformance_test.go
- Clear progression: Basic (required) → Standard (recommended) → Advanced (optional)
- Plugins can target specific level based on deployment needs

**Test Distribution**:

- **Basic** (6 tests): Core functionality - Name, Supports, GetProjectedCost, GetPricingSpec basics
- **Standard** (+5 tests): Production readiness - Error handling, consistency, 24h data, 10 concurrent
- **Advanced** (+4 tests): High performance - Latency thresholds, 50 concurrent, 30-day data

### 6. Error Message Formatting

**Decision**: Use structured error types with field-level detail

**Rationale**:

- FR-012 requires "clear, actionable error messages"
- SC-002 requires 95% of failures identify exact field/value
- Go error wrapping provides context chain

**Pattern**:

```go
type ValidationError struct {
    Field    string
    Value    interface{}
    Expected string
    Message  string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed: field=%s value=%v expected=%s: %s",
        e.Field, e.Value, e.Expected, e.Message)
}
```

## Dependencies Analysis

### Existing Dependencies (no changes needed)

- `google.golang.org/grpc` - gRPC client/server
- `google.golang.org/grpc/test/bufconn` - In-memory testing
- `google.golang.org/protobuf` - Proto message handling
- `testing` (stdlib) - Test framework integration

### New Internal Dependencies

- `encoding/json` (stdlib) - JSON report generation
- `sync` (stdlib) - Concurrency primitives
- `time` (stdlib) - Already used for timing

## Conclusion

All research topics resolved with clear decisions. No external dependencies needed.
Implementation should proceed with Phase 1 design artifacts.
