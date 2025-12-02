# Research: PluginSDK Conformance Testing Adapters

**Date**: 2025-11-30
**Feature**: 012-pluginsdk-conformance

## Research Tasks

### 1. Import Cycle Analysis

**Question**: Can `pluginsdk` import `sdk/go/testing` without creating cycles?

**Findings**:

- `sdk/go/testing` package does NOT import `pluginsdk`
- Current `sdk/go/testing` imports: `pbc` (proto), standard library only
- `pluginsdk` can safely import `sdk/go/testing` as an alias (e.g., `plugintesting`)

**Decision**: Safe to import - no cycle risk.

**Rationale**: The testing package is designed as a standalone conformance framework that depends
only on the proto definitions, not on the SDK implementation packages.

**Alternatives Considered**:

- Interface-based decoupling (not needed - no cycle)
- Separate conformance package (adds complexity without benefit)

### 2. Existing Conformance API Analysis

**Question**: What is the signature of existing conformance functions to wrap?

**Findings**:

From `sdk/go/testing/conformance.go`:

```go
// Three conformance levels
func RunBasicConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)
func RunStandardConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)
func RunAdvancedConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)

// Result type
type ConformanceResult struct {
    Version          string
    Timestamp        time.Time
    PluginName       string
    LevelAchieved    ConformanceLevel
    LevelAchievedStr string
    Summary          ResultSummary
    Categories       map[TestCategory]*CategoryResult
    Duration         time.Duration
    DurationStr      string
}

// Additional types needed
type ConformanceLevel int  // Basic, Standard, Advanced
type ResultSummary struct { Total, Passed, Failed, Skipped int }
```

**Decision**: Wrap all three `Run*Conformance` functions with Plugin→Server conversion.

**Rationale**: Direct 1:1 mapping preserves existing API semantics while adding Plugin support.

### 3. Error Handling Strategy

**Question**: How should adapters handle nil plugins and panics?

**Findings**:

- Current `NewServer(nil)` behavior: Creates server with nil plugin (will panic on method calls)
- Panics in plugin methods: Not recovered by existing conformance tests
- Go testing convention: `t.Fatal()` for unrecoverable errors, `t.Error()` for recoverable

**Decision**:

1. **Nil plugin**: Return early with descriptive error before creating server
2. **Panics**: Let underlying conformance tests handle (they use recover internally)
3. **Server creation failures**: Return error with context

**Rationale**: Fail-fast on nil prevents confusing downstream panics. Panic recovery is already
handled by the conformance harness which uses `defer recover()` in test execution.

**Alternatives Considered**:

- Wrapping all calls in recover() - Adds complexity, duplicates harness behavior
- Returning partial results on panic - Inconsistent with test-or-fail paradigm

### 4. Type Re-export Strategy

**Question**: How should `ConformanceResult` be exposed from `pluginsdk`?

**Findings**:

Go options for type exposure:

1. **Type alias**: `type ConformanceResult = plugintesting.ConformanceResult`
2. **Re-export via embed**: Not applicable (not embedding)
3. **Import for user**: Document that users import both packages

**Decision**: Use type aliases for key types (ConformanceResult, ConformanceLevel).

**Rationale**: Type aliases provide seamless interoperability - users can pass results to either
package's functions without conversion. This is cleaner than requiring dual imports.

```go
// In pluginsdk/conformance.go
type ConformanceResult = plugintesting.ConformanceResult
type ConformanceLevel = plugintesting.ConformanceLevel
```

### 5. PrintConformanceReport Implementation

**Question**: Should PrintConformanceReport wrap existing functionality or be new?

**Findings**:

From `sdk/go/testing/conformance_test.go`:

```go
func PrintConformanceReport(t *testing.T, result *ConformanceResult) {
    // Existing implementation using t.Log()
}
```

This function exists but is in a `_test.go` file (not exported from package).

From `sdk/go/testing/report.go`:

```go
func PrintReportTo(result *ConformanceResult, w io.Writer)
```

This is exported and writes to any io.Writer.

**Decision**: Create `PrintConformanceReport(t, result)` in pluginsdk that delegates to
`PrintReportTo` using a test log writer adapter.

**Rationale**: Reuses existing formatting logic while providing test-friendly API.

## Summary of Decisions

| Topic | Decision | Key Reason |
|-------|----------|------------|
| Import strategy | Direct import with alias | No cycle risk |
| API design | 1:1 function wrapping | Preserves existing semantics |
| Nil handling | Return error before server creation | Fail-fast, clear errors |
| Panic handling | Delegate to conformance harness | Avoid duplication |
| Type exposure | Type aliases for ConformanceResult, ConformanceLevel | Seamless interop |
| Report printing | Delegate to PrintReportTo | Reuse existing formatting |

## Dependencies Identified

- `sdk/go/testing` - Conformance suite (import as `plugintesting`)
- No new external dependencies required

## Unknowns Resolved

All technical unknowns from Technical Context have been resolved:

- ✅ Import cycle risk: None - safe to proceed
- ✅ Error handling strategy: Nil check + delegate panic handling
- ✅ Type re-export approach: Type aliases
