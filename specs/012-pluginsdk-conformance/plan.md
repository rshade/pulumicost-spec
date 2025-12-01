# Implementation Plan: PluginSDK Conformance Testing Adapters

**Branch**: `012-pluginsdk-conformance` | **Date**: 2025-11-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/012-pluginsdk-conformance/spec.md`

## Summary

Add adapter functions to the `pluginsdk` package that allow plugin developers to run conformance
tests directly on their `Plugin` implementations without manually converting to the raw gRPC
`CostSourceServiceServer` interface. The adapters wrap Plugin→Server conversion and delegate to
existing `sdk/go/testing` conformance functions.

## Technical Context

**Language/Version**: Go 1.25.4
**Primary Dependencies**: `sdk/go/testing` (conformance suite), `sdk/go/pluginsdk` (target package)
**Storage**: N/A (testing utilities only)
**Testing**: Go testing package with `go test`
**Target Platform**: All platforms (cross-platform Go SDK)
**Project Type**: Single - Go SDK library
**Performance Goals**: No overhead beyond existing conformance tests (adapter is thin wrapper)
**Constraints**: Must avoid import cycles between `pluginsdk` and `sdk/go/testing`
**Scale/Scope**: 4 adapter functions + 1 type re-export

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | PASS | No proto changes - SDK helper code only |
| II. Multi-Provider Consistency | N/A | Not provider-specific |
| III. Test-First Protocol | PASS | Will write tests first for adapter functions |
| IV. Protobuf Backward Compatibility | PASS | No proto changes |
| V. Comprehensive Documentation | PASS | Will update README with usage examples |
| VI. Performance as Requirement | PASS | Delegates to existing conformance tests |
| VII. Validation at Multiple Levels | PASS | Unit tests for adapters, integration via existing suite |

**Gate Result**: PASS - All applicable constitutional principles satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/012-pluginsdk-conformance/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - internal API)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── sdk.go               # Existing - Plugin, Server, NewServer
├── testing.go           # Existing - TestServer, TestPlugin
├── conformance.go       # NEW - Adapter functions (FR-001 to FR-006)
├── conformance_test.go  # NEW - Unit tests for adapters
└── README.md            # Update with conformance examples

sdk/go/testing/
├── conformance.go       # Existing - RunBasicConformance, etc.
├── harness.go           # Existing - TestHarness
└── (other files)        # Existing - No changes needed
```

**Structure Decision**: Single new file `conformance.go` in `pluginsdk` package containing all
adapter functions. Tests in corresponding `conformance_test.go`. This follows existing package
patterns (e.g., `testing.go` + `testing_test.go`).

## Complexity Tracking

> No violations - all gates pass.
