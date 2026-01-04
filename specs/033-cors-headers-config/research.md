# Research: Configurable CORS Headers

**Feature**: 033-cors-headers-config
**Date**: 2026-01-04

## Research Summary

No NEEDS CLARIFICATION items in Technical Context. Research focuses on validating the existing
patterns and CORS best practices to ensure the implementation aligns with standards.

## Decision 1: Field Type for Header Configuration

**Decision**: Use `[]string` for both `AllowedHeaders` and `ExposedHeaders` fields.

**Rationale**:

- Matches existing `AllowedOrigins` field type in `WebConfig`
- Simple, idiomatic Go type that's easy to work with
- Allows natural comma-joining for HTTP header value construction
- No need for custom types or validation wrappers

**Alternatives Considered**:

- `map[string]bool`: Rejected - adds complexity for no benefit (no lookup needed)
- Custom `HeaderSet` type: Rejected - over-engineering for simple string list

## Decision 2: Nil vs Empty Slice Semantics

**Decision**: `nil` means "use defaults", empty `[]string{}` means "no custom headers".

**Rationale**:

- Follows Go idiom where nil slice means "not set" vs empty means "explicitly empty"
- Matches existing `AllowedOrigins` behavior in the codebase
- Enables backward compatibility - existing code with nil fields gets default headers
- Clear semantic distinction for API users

**Alternatives Considered**:

- Pointer to slice `*[]string`: Rejected - adds nil-check complexity for no benefit
- Separate "UseDefaults" bool: Rejected - redundant with nil check

## Decision 3: Builder Method Pattern

**Decision**: Follow exact pattern from `WithAllowedOrigins()` for new methods.

**Rationale**:

- Consistency with existing codebase
- Defensive copy prevents external mutation
- Value receiver + return copy enables chaining

**Implementation Pattern** (from existing code):

```go
func (c WebConfig) WithAllowedHeaders(headers []string) WebConfig {
    if headers != nil {
        c.AllowedHeaders = make([]string, len(headers))
        copy(c.AllowedHeaders, headers)
    } else {
        c.AllowedHeaders = nil
    }
    return c
}
```

**Alternatives Considered**:

- Pointer receiver: Rejected - breaks value semantics and chaining
- Append mode: Rejected - replace is simpler and more predictable

## Decision 4: Default Header Constants

**Decision**: Define package-level constants for default headers.

**Rationale**:

- Centralizes the default values (DRY)
- Makes defaults documentable and testable
- Allows users to reference defaults in their configuration

**Implementation**:

```go
// DefaultAllowedHeaders contains the CORS allowed headers for Connect/gRPC-Web.
// These headers are used when WebConfig.AllowedHeaders is nil.
const DefaultAllowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, " +
    "Authorization, X-CSRF-Token, X-Requested-With, Connect-Protocol-Version, " +
    "Connect-Timeout-Ms, Grpc-Timeout, X-Grpc-Web, X-User-Agent"

// DefaultExposedHeaders contains the CORS exposed headers for Connect/gRPC-Web.
// These headers are used when WebConfig.ExposedHeaders is nil.
const DefaultExposedHeaders = "Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin, " +
    "Connect-Content-Encoding, Connect-Content-Type"
```

**Alternatives Considered**:

- Slice constants: Not possible in Go (slices can't be const)
- Package-level variables: Could work but constants are safer and documented

## Decision 5: Header Joining Strategy

**Decision**: Use `strings.Join(headers, ", ")` for HTTP header value construction.

**Rationale**:

- Standard HTTP header format uses comma-space separation
- `strings.Join` is efficient (single allocation)
- Handles empty slice correctly (returns "")

**Alternatives Considered**:

- Manual loop with StringBuilder: Rejected - strings.Join is cleaner and faster
- Pre-joining at config time: Rejected - adds cached state complexity

## CORS Standard Reference

Per MDN and W3C CORS specification:

- `Access-Control-Allow-Headers`: Comma-separated list of header names browser may send
- `Access-Control-Expose-Headers`: Comma-separated list of headers browser JavaScript can read
- Case-insensitive matching (no normalization needed)
- Duplicates allowed (browser deduplicates internally)

## Performance Considerations

The performance impact should be minimal:

1. **nil check**: ~1 nanosecond (branch prediction)
2. **strings.Join**: ~50-100 nanoseconds for 12 headers
3. **Total expected**: <500 nanoseconds, well under 1 microsecond target

Benchmark verification will be added to confirm SC-005 compliance.

## Test Strategy

1. **Unit tests for builder methods**: Verify defensive copy, nil handling, chaining
2. **Unit tests for corsMiddleware**: Verify header selection logic
3. **Integration tests**: Verify end-to-end CORS response headers
4. **Benchmark tests**: Verify <1μs overhead per request
