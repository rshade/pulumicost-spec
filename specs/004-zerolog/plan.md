# Implementation Plan: Zerolog SDK Logging Utilities

**Branch**: `001-zerolog` | **Date**: 2025-11-24 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-zerolog/spec.md`

## Summary

Add standardized logging utilities to the Plugin SDK using zerolog v1.34.0+ to
enable consistent structured logging and distributed tracing across all
PulumiCost plugins. The SDK will provide logger construction, gRPC interceptors
for trace ID propagation, standard field name constants, and operation timing
helpers.

## Technical Context

**Language/Version**: Go 1.24+ (matches existing SDK)
**Primary Dependencies**: zerolog v1.34.0+, google.golang.org/grpc,
google.golang.org/grpc/metadata
**Storage**: N/A (logging to stderr/files)
**Testing**: go test with table-driven tests, bufconn for gRPC integration tests
**Target Platform**: Linux server (cross-platform Go)
**Project Type**: Single project (SDK library extension)
**Performance Goals**: Zero-allocation logging (zerolog default), <1μs per log
call overhead
**Constraints**: Must not break existing SDK API, maintain backward compatibility
**Scale/Scope**: SDK utility package (~300 LOC), 90%+ test coverage

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Proto Spec-First | ✅ PASS | No proto changes - SDK helper package only |
| II. Multi-Provider Consistency | ✅ PASS | Logging utilities are provider-agnostic |
| III. Test-First Protocol | ✅ PASS | Tests will be written first per TDD |
| IV. Backward Compatibility | ✅ PASS | New package, no breaking changes |
| V. Comprehensive Documentation | ✅ PASS | Documentation required in spec |
| VI. Performance Requirements | ✅ PASS | Benchmarks will track logging overhead |
| VII. Multi-Level Validation | ✅ PASS | Unit + integration tests planned |

**Gate Result**: PASS - All principles satisfied. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/004-zerolog/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0: Technology research
├── data-model.md        # Phase 1: Data structures
├── quickstart.md        # Phase 1: Usage guide
├── contracts/           # Phase 1: API contracts
└── tasks.md             # Phase 2: Implementation tasks
```

### Source Code (repository root)

```text
sdk/go/
├── pluginsdk/           # NEW: Plugin SDK utilities package
│   ├── logging.go       # Logger construction, interceptors, helpers
│   └── logging_test.go  # Unit and integration tests
├── pricing/             # Existing: Billing modes and validation
├── proto/               # Existing: Generated protobuf code
├── registry/            # Existing: Plugin registry types
└── testing/             # Existing: Test harness and mocks
```

**Structure Decision**: Extend existing SDK with new `pluginsdk` package for
plugin development utilities. This follows the established pattern of separate
packages for different concerns (pricing, registry, testing).

## Complexity Tracking

No violations requiring justification. Feature is a straightforward SDK utility
addition that follows existing patterns and all constitutional principles.
