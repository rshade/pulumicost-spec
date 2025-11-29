# Implementation Plan: Trace ID Validation for TracingUnaryServerInterceptor

**Branch**: `008-trace-id-validation` | **Date**: 2025-11-26 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-trace-id-validation/spec.md`

## Summary

Add validation to `TracingUnaryServerInterceptor()` in `sdk/go/pluginsdk/logging.go`
to validate incoming trace_id values against the established format (32 lowercase
hexadecimal characters, not all zeros) before storing them in the request context.
Invalid trace IDs are replaced with newly generated valid ones using
cryptographically secure random generation.

## Technical Context

**Language/Version**: Go 1.24+ (toolchain go1.25.4)
**Primary Dependencies**: google.golang.org/grpc v1.77.0, crypto/rand (stdlib)
**Storage**: N/A (in-memory context propagation only)
**Testing**: go test with bufconn for in-memory gRPC testing
**Target Platform**: Linux server (plugin host environment)
**Project Type**: SDK library (single project)
**Performance Goals**: <1ms validation overhead per request (SC-004)
**Constraints**: Zero external UUID dependencies, backward compatible API
**Scale/Scope**: Interceptor processes all gRPC requests to CostSource plugins

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                           | Status | Notes                                                    |
| ----------------------------------- | ------ | -------------------------------------------------------- |
| I. gRPC Proto Specification-First   | PASS   | No proto changes needed - interceptor is SDK helper code |
| II. Multi-Provider Consistency      | PASS   | Validation applies uniformly to all providers            |
| III. Test-First Protocol            | GATE   | Tests MUST be written before implementation              |
| IV. Protobuf Backward Compatibility | PASS   | No proto changes - SDK-only modification                 |
| V. Comprehensive Documentation      | GATE   | Inline comments and README updates required              |
| VI. Performance as Requirement      | GATE   | Benchmark validation overhead against 1ms target         |
| VII. Validation at Multiple Levels  | PASS   | Enhances service-layer validation                        |

**Gate Requirements**:

- Write interceptor tests with invalid trace IDs that currently pass through
- Verify tests fail before implementation
- Run benchmarks to establish baseline and measure overhead

## Project Structure

### Documentation (this feature)

```text
specs/008-trace-id-validation/
├── plan.md              # This file
├── research.md          # Phase 0 output - validation patterns research
├── data-model.md        # Phase 1 output - trace ID entity model
├── quickstart.md        # Phase 1 output - usage guide
├── contracts/           # Phase 1 output - N/A (no new API contracts)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
sdk/go/
├── pluginsdk/
│   ├── logging.go           # MODIFY: TracingUnaryServerInterceptor
│   ├── logging_test.go      # MODIFY: Add validation test cases
│   └── traceid.go           # ADD: Trace ID generation helper
├── pricing/
│   ├── observability_validate.go  # EXISTING: ValidateTraceID function
│   └── observability_test.go      # EXISTING: Validation tests
└── testing/
    └── harness.go           # USE: For integration testing
```

**Structure Decision**: Single project SDK structure. Modifications contained within
`sdk/go/pluginsdk/` package with reuse of existing validation from `sdk/go/pricing/`.

## Complexity Tracking

> No constitution violations requiring justification.

| Aspect              | Decision                              | Rationale                                       |
| ------------------- | ------------------------------------- | ----------------------------------------------- |
| Validation Location | Interceptor (not separate middleware) | Keeps single responsibility, minimal API change |
| ID Generation       | crypto/rand + hex encoding            | Standard library only, no external deps         |
| API Compatibility   | Same function signature               | Backward compatible - validation is internal    |

## Design Decisions

### Decision 1: Validation Strategy

**Choice**: Option A (Strict validation, always-on)

**Rationale**: The spec prioritizes security by default (User Story 3) and specifies
that validation should be "enabled by default without requiring code changes."
Configurable options (Option B from issue) add complexity without clear benefit -
there's no valid use case for allowing malicious trace IDs.

### Decision 2: Invalid Trace ID Handling

**Choice**: Generate replacement (graceful degradation)

**Rationale**: FR-009 mandates "System MUST NOT return errors to callers due to
invalid trace_id values." Generating a valid replacement maintains request flow
while ensuring observability.

### Decision 3: Generation Method

**Choice**: crypto/rand with hex encoding (32 characters)

**Rationale**:

- Standard library only (no uuid dependency)
- Cryptographically secure random source
- Matches existing trace ID format exactly
- Zero-allocation pattern aligns with registry package conventions

### Decision 4: Validation Reuse

**Choice**: Import existing `ValidateTraceID` from pricing package

**Rationale**:

- Avoids duplication of validation logic
- Leverages tested, production-ready code
- Single source of truth for trace ID format rules

### Decision 5: Observability of Validation Failures

**Choice**: Silent replacement (no logging in interceptor)

**Rationale**:

- Interceptor should be lightweight and not assume logging configuration
- Plugin developers can add logging in their handlers if needed
- Logging invalid input could itself become a log injection vector
- Future: Consider adding optional callback for observability

## Implementation Approach

### Phase 1: Test-First Development

1. Add test cases for invalid trace IDs in `logging_test.go`
2. Tests MUST fail against current implementation
3. Test cases cover all edge cases from spec:
   - Wrong length (too short, too long)
   - Non-hex characters
   - All zeros
   - Unicode/control characters
   - Excessive length (10KB)

### Phase 2: Core Implementation

1. Add `traceid.go` with `GenerateTraceID()` function
2. Modify `TracingUnaryServerInterceptor()` to:
   - Call `ValidateTraceID()` on extracted value
   - Generate replacement if validation fails
   - Store validated/generated trace ID in context

### Phase 3: Validation & Documentation

1. All tests pass (including new validation tests)
2. Benchmark shows <1ms overhead
3. Update package documentation
4. Add usage example in README

## Risk Assessment

| Risk                      | Mitigation                                                         |
| ------------------------- | ------------------------------------------------------------------ |
| Breaking existing plugins | API signature unchanged; only behavior improved                    |
| Performance regression    | Benchmark before/after; optimize regex if needed                   |
| Import cycle              | pricing → pluginsdk import direction verified safe                 |
| Test flakiness            | Use deterministic test inputs; mock random for ID generation tests |
