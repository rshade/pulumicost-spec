# Research: Trace ID Validation Patterns

**Feature**: 008-trace-id-validation
**Date**: 2025-11-26
**Purpose**: Resolve technical unknowns and document best practices for implementation

## Research Summary

All technical unknowns have been resolved through codebase exploration and
industry best practices analysis.

## Topic 1: Existing Trace ID Validation Logic

### Decision

Reuse existing `ValidateTraceID()` function from `sdk/go/pricing/observability_validate.go`.

### Rationale

The existing function provides comprehensive validation:

- Format check: regex `^[0-9a-f]{32}$` (32 lowercase hexadecimal characters)
- All-zeros rejection: explicit check for `00000000000000000000000000000000`
- Empty string handling: returns `nil` (valid - trace ID is optional)

### Alternatives Considered

| Alternative                       | Rejected Because                       |
| --------------------------------- | -------------------------------------- |
| Duplicate validation in pluginsdk | Violates DRY, maintenance burden       |
| Move validation to pluginsdk      | Breaking change to pricing package API |
| Create shared validation package  | Over-engineering for single function   |

### Code Reference

```go
// sdk/go/pricing/observability_validate.go:96-112
func ValidateTraceID(traceID string) error {
    if traceID == "" {
        return nil // trace ID is optional
    }
    if !traceIDRegex.MatchString(traceID) {
        return fmt.Errorf("trace ID '%s' must be 32 hexadecimal characters", traceID)
    }
    if traceID == "00000000000000000000000000000000" {
        return errors.New("trace ID cannot be all zeros")
    }
    return nil
}
```

## Topic 2: Trace ID Generation Pattern

### Decision

Use `crypto/rand` with hex encoding to generate 16 random bytes (32 hex characters).

### Rationale

- **Standard library only**: No external UUID dependencies
- **Cryptographically secure**: Uses OS-level entropy source
- **Collision resistance**: 128 bits of randomness (2^128 possibilities)
- **Format compliance**: Produces exactly 32 lowercase hex characters
- **Performance**: Single allocation, ~150ns per generation

### Alternatives Considered

| Alternative               | Rejected Because                                |
| ------------------------- | ----------------------------------------------- |
| google/uuid package       | External dependency, overkill for simple format |
| time-based IDs            | Not cryptographically random, predictable       |
| incrementing counter      | Requires state, not thread-safe without locks   |
| Span ID format (16 chars) | Wrong format for trace ID                       |

### Implementation Pattern

```go
import (
    "crypto/rand"
    "encoding/hex"
)

// GenerateTraceID creates a new valid trace ID (32 hex characters).
func GenerateTraceID() string {
    b := make([]byte, 16)
    _, _ = rand.Read(b) // crypto/rand.Read never returns error on supported systems
    return hex.EncodeToString(b)
}
```

### Performance Characteristics

- Generation: ~150ns per call (based on crypto/rand benchmarks)
- Memory: Single 16-byte allocation + 32-byte string
- Thread-safe: crypto/rand is concurrent-safe

## Topic 3: gRPC Interceptor Integration Pattern

### Decision

Modify `TracingUnaryServerInterceptor()` inline without changing function signature.

### Rationale

- **Backward compatible**: No API changes for existing users
- **Single location**: All trace ID processing in one interceptor
- **Fail-safe**: Invalid inputs become valid outputs automatically

### Alternatives Considered

| Alternative                     | Rejected Because                                     |
| ------------------------------- | ---------------------------------------------------- |
| New interceptor function        | Would require code changes in existing plugins       |
| Middleware chain                | Over-engineering, gRPC already has interceptor chain |
| Separate validation interceptor | Additional stack frame, harder to reason about       |

### Current Implementation (to be modified)

```go
// sdk/go/pluginsdk/logging.go:61-80
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{},
        _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if values := md.Get(TraceIDMetadataKey); len(values) > 0 {
                ctx = ContextWithTraceID(ctx, values[0])  // CHANGE: Add validation here
            }
        }
        return handler(ctx, req)
    }
}
```

### Proposed Implementation

```go
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{},
        _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        var traceID string
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if values := md.Get(TraceIDMetadataKey); len(values) > 0 {
                traceID = values[0]
            }
        }
        // Validate and generate replacement if invalid or missing
        if err := pricing.ValidateTraceID(traceID); err != nil || traceID == "" {
            traceID = GenerateTraceID()
        }
        ctx = ContextWithTraceID(ctx, traceID)
        return handler(ctx, req)
    }
}
```

## Topic 4: Import Dependency Direction

### Decision

Import `pricing` package from `pluginsdk` package.

### Rationale

- **No cycle**: `pluginsdk` → `pricing` direction is safe
- **Existing pattern**: `pluginsdk` already imports gRPC packages
- **Logical hierarchy**: SDK helpers can depend on domain validation

### Verification

```text
Current dependencies:
  sdk/go/pricing/       → (no internal SDK imports)
  sdk/go/pluginsdk/     → (can safely import pricing)
  sdk/go/testing/       → imports pricing, pluginsdk
```

The import graph remains acyclic with this change.

## Topic 5: Test Strategy

### Decision

Add table-driven tests with edge cases to `logging_test.go`.

### Test Cases Required

| Test Case        | Input                                | Expected Behavior   |
| ---------------- | ------------------------------------ | ------------------- |
| Valid trace ID   | `"abcdef1234567890abcdef1234567890"` | Preserved unchanged |
| Empty trace ID   | `""`                                 | New ID generated    |
| Missing metadata | No header                            | New ID generated    |
| Too short        | `"abcdef"`                           | New ID generated    |
| Too long         | `"abcdef...01"` (33+ chars)          | New ID generated    |
| Non-hex chars    | `"gggggggggggggggggggggggggggggggg"` | New ID generated    |
| All zeros        | `"00000000000000000000000000000000"` | New ID generated    |
| Control chars    | `"abc\ndef..."`                      | New ID generated    |
| Unicode          | `"ąbcdef..."`                        | New ID generated    |
| Excessive length | `strings.Repeat("a", 10240)`         | New ID generated    |
| Multiple headers | `["valid1", "invalid2"]`             | First validated     |

### Benchmark Requirements

```go
func BenchmarkTracingUnaryServerInterceptor(b *testing.B) {
    // Baseline: valid trace ID (no generation)
    // Worst case: invalid trace ID (validation + generation)
}
```

Target: Both scenarios under 1ms (SC-004).

## Research Conclusions

All technical decisions are resolved:

1. **Validation**: Reuse existing `ValidateTraceID()` from pricing package
2. **Generation**: Use `crypto/rand` + hex encoding (stdlib only)
3. **Integration**: Modify existing interceptor inline (backward compatible)
4. **Imports**: `pluginsdk` → `pricing` direction is safe
5. **Testing**: Table-driven tests covering all 11 edge cases from spec

No NEEDS CLARIFICATION items remain. Ready for Phase 1 design artifacts.
