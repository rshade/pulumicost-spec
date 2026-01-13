# Implementation Plan: Configurable CORS Headers

**Branch**: `033-cors-headers-config` | **Date**: 2026-01-04 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from GitHub Issue #228

## Summary

Add two new optional fields (`AllowedHeaders` and `ExposedHeaders`) to the `WebConfig` struct in the
pluginsdk package, enabling plugin developers to customize CORS header configuration. The implementation
follows the existing builder pattern with `With*` methods and uses nil-defaults for backward compatibility.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: net/http (stdlib), strings (stdlib) - no new dependencies
**Storage**: N/A (stateless configuration extension)
**Testing**: Go testing package with table-driven tests, existing TestHarness for integration
**Target Platform**: Linux/macOS/Windows servers running FinFocus plugins
**Project Type**: Single SDK package extension
**Performance Goals**: <1 microsecond overhead per CORS request (SC-005)
**Constraints**: 100% backward compatibility, defensive slice copying
**Scale/Scope**: 2 new struct fields, 2 new builder methods, 1 middleware update

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
| --------- | ------ | ----- |
| I. gRPC Proto-First | N/A | No proto changes - WebConfig is Go-only |
| II. Multi-Provider Consistency | N/A | Feature is provider-agnostic |
| III. Test-First (NON-NEGOTIABLE) | PASS | Tests written before implementation |
| IV. Backward Compatibility | PASS | nil-defaults preserve existing behavior |
| V. Documentation (NON-NEGOTIABLE) | PASS | README and godoc updates included |
| VI. Performance Requirements | PASS | Benchmark tests verify <1μs overhead |
| VII. Multi-Layer Validation | PASS | Unit + integration tests provided |

**Result**: All applicable gates PASS. No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/033-cors-headers-config/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── options.go           # WebConfig struct + builder methods (MODIFY)
├── options_test.go      # Builder method tests (CREATE or EXTEND)
├── sdk.go               # corsMiddleware function (MODIFY)
├── sdk_test.go          # CORS middleware tests (EXTEND)
└── README.md            # Package documentation (MODIFY)
```

**Structure Decision**: This feature extends existing files in `sdk/go/pluginsdk/`. No new files
created except potentially `cors_test.go` for focused CORS testing if the existing test file is large.

## Complexity Tracking

No violations to justify. Implementation follows existing patterns:

- Builder methods follow `WithAllowedOrigins()` pattern exactly
- Middleware extension is minimal (conditional header selection)
- No new abstractions or patterns introduced
