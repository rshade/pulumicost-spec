# Implementation Plan: ISO 4217 Currency Validation Package

**Branch**: `013-iso4217-currency` | **Date**: 2025-11-30 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/013-iso4217-currency/spec.md`

## Summary

Extract the ISO 4217 currency validation logic currently embedded in `sdk/go/pluginsdk/focus_conformance.go`
into a standalone, reusable `sdk/go/currency/` package. The package will provide zero-allocation currency
code validation (following the established registry package pattern), complete ISO 4217 currency metadata
(code, name, numeric code, minor units), and maintain backward compatibility with existing FOCUS conformance
tests.

## Technical Context

**Language/Version**: Go 1.25.4 (toolchain as specified in go.mod)
**Primary Dependencies**: None (stdlib only - no external dependencies required)
**Storage**: N/A (static in-memory data structures)
**Testing**: `go test` with table-driven tests and benchmarks
**Target Platform**: Cross-platform Go library (Linux, macOS, Windows)
**Project Type**: Single Go package within existing SDK structure
**Performance Goals**: <15 ns/op validation, 0 B/op, 0 allocs/op (matching registry package)
**Constraints**: Must not break existing FOCUS conformance tests
**Scale/Scope**: 180+ ISO 4217 currency codes, single package addition

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | N/A | No proto changes - SDK helper package |
| II. Multi-Provider gRPC Consistency | N/A | Currency validation is provider-agnostic |
| III. Test-First Protocol | PASS | Will write tests before implementation |
| IV. Protobuf Backward Compatibility | N/A | No proto changes |
| V. Comprehensive Documentation | PASS | Package docs, README, CLAUDE.md planned |
| VI. Performance as a gRPC Requirement | PASS | Benchmarks required per spec |
| VII. Validation at Multiple Levels | PASS | Unit tests + integration via existing conformance |

**Gate Status**: PASS - All applicable principles satisfied. This is a pure SDK enhancement
with no proto changes required.

## Project Structure

### Documentation (this feature)

```text
specs/013-iso4217-currency/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (Go package API)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
sdk/go/
├── currency/            # NEW: ISO 4217 currency package
│   ├── currency.go      # Currency type, metadata, and data
│   ├── validate.go      # IsValid() validation function
│   ├── doc.go           # Package documentation comments
│   ├── currency_test.go # Unit tests with table-driven cases
│   ├── benchmark_test.go # Performance benchmarks
│   ├── CLAUDE.md        # Package-specific guidance
│   └── README.md        # Package documentation
├── pluginsdk/
│   └── focus_conformance.go  # MODIFIED: Import currency package
├── registry/            # REFERENCE: Existing pattern to follow
│   └── domain.go
└── ...
```

**Structure Decision**: Single new package `sdk/go/currency/` following the existing SDK package
organization. The package follows the established zero-allocation validation pattern from
`sdk/go/registry/domain.go`.

## Complexity Tracking

> No constitutional violations requiring justification.
