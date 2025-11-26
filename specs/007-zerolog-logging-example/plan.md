# Implementation Plan: Structured Logging Example for EstimateCost

**Branch**: `007-zerolog-logging-example` | **Date**: 2025-11-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/007-zerolog-logging-example/spec.md`

## Summary

Add a structured logging example in `sdk/go/testing/integration_test.go` demonstrating zerolog
integration for the EstimateCost RPC. The example will serve as the canonical reference for
NFR-001 compliance, showing request/response/error logging patterns with correlation ID
propagation and timing measurement.

## Technical Context

**Language/Version**: Go 1.24+ (toolchain go1.25.4)
**Primary Dependencies**: zerolog v1.34.0+, google.golang.org/grpc, sdk/go/testing harness
**Storage**: N/A (example code, no data persistence)
**Testing**: go test (integration test function added to existing file)
**Target Platform**: Linux/macOS/Windows (cross-platform Go SDK)
**Project Type**: Single (SDK addition to existing test file)
**Performance Goals**: N/A (documentation example, not performance-critical)
**Constraints**: Must follow existing integration_test.go patterns; must use 005-zerolog utilities
**Scale/Scope**: Single test function with comprehensive logging examples

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | PASS | No proto changes required; uses existing EstimateCost RPC |
| II. Multi-Provider gRPC Consistency | PASS | Example uses mock plugin supporting all providers |
| III. Test-First Protocol | PASS | Example IS a test function demonstrating logging |
| IV. Protobuf Backward Compatibility | PASS | No proto changes; documentation only |
| V. Comprehensive Documentation | PASS | Example serves as documentation with code comments |
| VI. Performance as gRPC Requirement | PASS | Not applicable; documentation example |
| VII. Validation at Multiple Levels | PASS | Example validates logging output structure |

**Gate Status**: ALL PASS - No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/007-zerolog-logging-example/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal for this feature)
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
sdk/go/
├── testing/
│   ├── integration_test.go    # Target file - add logging example test
│   ├── harness.go             # Existing test harness (used by example)
│   ├── mock_plugin.go         # Existing mock plugin (used by example)
│   └── CLAUDE.md              # Testing package documentation
├── logging/                   # From 005-zerolog (dependency)
│   ├── logger.go              # NewPluginLogger, field constants
│   ├── tracing.go             # TraceIDFromContext, ContextWithTraceID
│   └── timing.go              # LogOperation helper
└── proto/                     # Generated gRPC code (existing)
```

**Structure Decision**: Single file addition to existing `sdk/go/testing/integration_test.go`.
No new directories or files required beyond the test function itself. Uses existing testing
framework infrastructure and depends on 005-zerolog logging utilities.

## Complexity Tracking

No violations to justify - all constitutional gates pass.
