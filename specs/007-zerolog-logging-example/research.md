# Research: Structured Logging Example for EstimateCost

**Feature**: 007-zerolog-logging-example
**Date**: 2025-11-26

## Overview

This research consolidates findings for implementing a structured logging example demonstrating
zerolog integration for the EstimateCost RPC in the FinFocus SDK testing package.

## Dependencies Analysis

### 005-zerolog SDK Logging Utilities

**Status**: Specification complete (005-zerolog/spec.md), implementation pending

**Key Functions Required**:

| Function             | Purpose                          | Signature Pattern                                                         |
| -------------------- | -------------------------------- | ------------------------------------------------------------------------- |
| `NewPluginLogger`    | Create configured zerolog.Logger | `(name, version string, level zerolog.Level, w io.Writer) zerolog.Logger` |
| `TraceIDFromContext` | Extract trace_id from context    | `(ctx context.Context) string`                                            |
| `ContextWithTraceID` | Add trace_id to context          | `(ctx context.Context, traceID string) context.Context`                   |
| `LogOperation`       | Timing helper for operations     | `(logger zerolog.Logger, operation string) func()`                        |

**Field Constants** (from FR-006 of 005-zerolog):

- `FieldTraceID` = "trace_id"
- `FieldComponent` = "component"
- `FieldOperation` = "operation"
- `FieldDurationMs` = "duration_ms"
- `FieldResourceURN` = "resource_urn"
- `FieldResourceType` = "resource_type"
- `FieldPluginName` = "plugin_name"
- `FieldPluginVersion` = "plugin_version"
- `FieldCostMonthly` = "cost_monthly"
- `FieldAdapter` = "adapter"
- `FieldErrorCode` = "error_code"

**Decision**: Example will use these utilities from `sdk/go/logging/` package when implemented.
**Rationale**: Spec 005-zerolog defines the standard utilities; using them ensures consistency.
**Alternatives**: Could inline zerolog usage, but that defeats the purpose of showing SDK patterns.

### 006-estimate-cost EstimateCost RPC

**Status**: Implemented in proto/finfocus/v1/costsource.proto (lines 29-42, 440-484)

**Key Types**:

- `EstimateCostRequest`: resource_type (string), attributes (google.protobuf.Struct)
- `EstimateCostResponse`: currency (string), cost_monthly (double)

**Error Cases** (from proto comments):

- `InvalidArgument`: Invalid resource_type format or missing required attributes
- `NotFound`: Unsupported resource_type for this plugin
- `Unavailable`: Pricing source temporarily unavailable

**Decision**: Example will use existing MockPlugin which already has EstimateCost support.
**Rationale**: MockPlugin in mock_plugin.go already implements EstimateCost with configurable behavior.
**Alternatives**: Create new mock - rejected as existing infrastructure is sufficient.

## Existing Testing Patterns

### Integration Test Patterns (sdk/go/testing/integration_test.go)

#### Pattern 1: Basic Test Structure

```go
func TestFeatureName(t *testing.T) {
    plugin := plugintesting.NewMockPlugin()
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    client := harness.Client()
    ctx := context.Background()

    // Test logic
}
```

#### Pattern 2: Subtests with t.Run

```go
t.Run("SubtestName", func(t *testing.T) {
    // Subtest logic
})
```

#### Pattern 3: Error Condition Testing

```go
plugin := plugintesting.ConfigurableErrorMockPlugin()
plugin.ShouldErrorOnName = true
// Test error handling
```

**Decision**: Example will follow existing patterns with subtests for each logging scenario.
**Rationale**: Consistency with existing tests improves maintainability.
**Alternatives**: Separate test functions - rejected for cohesion of logging examples.

### Mock Plugin Capabilities (mock_plugin.go)

**EstimateCost Support**:

- `ShouldErrorOnEstimateCost` flag for error injection
- `EstimateCostDelay` for timeout testing
- Returns simulated cost based on resource attributes

**Decision**: Use ConfigurableErrorMockPlugin for error logging examples.
**Rationale**: Built-in error injection makes testing error logging patterns straightforward.
**Alternatives**: Custom mock - rejected as existing mock is sufficient.

## Logging Best Practices Research

### zerolog Builder Pattern

**Standard Usage**:

```go
logger.Info().
    Str(FieldOperation, "EstimateCost").
    Str(FieldResourceType, resourceType).
    Msg("Processing request")
```

**Decision**: Example will use builder pattern with standard field constants.
**Rationale**: This is idiomatic zerolog usage and matches 005-zerolog spec.

### Correlation ID Propagation

**Pattern**:

1. Extract trace_id from context at request entry
2. Add trace_id to all log entries for that request
3. Use logger.With() for persistent fields

**Decision**: Use ContextWithTraceID and TraceIDFromContext from 005-zerolog.
**Rationale**: Standard SDK utilities ensure consistent tracing across plugins.

### Operation Timing

**Pattern**:

```go
done := LogOperation(logger, "EstimateCost")
defer done()
// Operation logic
```

**Decision**: Demonstrate LogOperation helper from 005-zerolog.
**Rationale**: Timing is critical for performance monitoring (NFR-001).

### Sensitive Data Handling

**Best Practice**: Never log attribute values directly - log count and keys only.

**Decision**: Example will demonstrate logging `len(attributes.Fields)` not values.
**Rationale**: Prevents accidental credential/secret exposure in logs.

## Example Code Structure

### Proposed Test Function Structure

```go
func TestStructuredLoggingExample(t *testing.T) {
    // Setup with logging

    t.Run("RequestLogging", func(t *testing.T) {
        // Demonstrate logging incoming requests
    })

    t.Run("SuccessResponseLogging", func(t *testing.T) {
        // Demonstrate logging successful responses
    })

    t.Run("ErrorLogging", func(t *testing.T) {
        // Demonstrate error logging with context
    })

    t.Run("CorrelationIDPropagation", func(t *testing.T) {
        // Demonstrate trace_id across logs
    })
}
```

**Decision**: Single test function with subtests covering all scenarios.
**Rationale**: Keeps all logging examples together for easy reference.

## Technical Considerations

### Log Output Capture

**Challenge**: Tests need to verify log output structure.

**Solution**: Use `bytes.Buffer` as output writer for zerolog:

```go
var buf bytes.Buffer
logger := zerolog.New(&buf).With().Timestamp().Logger()
```

**Decision**: Use buffer-based logging for testable output.
**Rationale**: Allows assertions on JSON log structure without external dependencies.

### JSON Log Parsing

**Approach**: Parse log buffer as JSON and assert on fields:

```go
var logEntry map[string]interface{}
json.Unmarshal(buf.Bytes(), &logEntry)
// Assert on logEntry["operation"], logEntry["trace_id"], etc.
```

**Decision**: Use standard encoding/json for log verification.
**Rationale**: Simple, no external dependencies, matches production parsing patterns.

## Resolved Questions

| Topic                 | Resolution                                          |
| --------------------- | --------------------------------------------------- |
| Which test file?      | `sdk/go/testing/integration_test.go` per FR-008     |
| Which mock to use?    | `ConfigurableErrorMockPlugin` for error scenarios   |
| How to capture logs?  | `bytes.Buffer` as zerolog output writer             |
| How to verify output? | JSON parsing with standard library                  |
| Field naming?         | Use constants from 005-zerolog (FieldTraceID, etc.) |
| Sensitive data?       | Log attribute count, never attribute values         |

## Risks and Mitigations

| Risk                            | Impact   | Mitigation                                                      |
| ------------------------------- | -------- | --------------------------------------------------------------- |
| 005-zerolog not yet implemented | Blocking | Document expected API; example can be written against interface |
| Mock plugin behavior changes    | Low      | Test against stable public API only                             |
| Log format changes              | Low      | Use field constants for forward compatibility                   |
