# Implementation Plan: Plugin Conformance Test Suite

**Branch**: `011-plugin-conformance-suite` | **Date**: 2025-11-28 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/010-plugin-conformance-suite/spec.md`

## Summary

Formalize and extend the existing `sdk/go/testing/` framework into a documented Plugin Conformance
Test Suite. The suite will be a Go library that plugin authors import to validate their
implementations against the pulumicost-spec at three conformance levels (Basic, Standard, Advanced)
with structured JSON output for CI/CD integration.

## Technical Context

**Language/Version**: Go 1.24+ (matches existing SDK)
**Primary Dependencies**: google.golang.org/grpc, google.golang.org/protobuf, bufconn (existing)
**Storage**: N/A (in-memory testing only)
**Testing**: go test, go test -race, go test -bench (standard Go tooling)
**Target Platform**: Any platform supporting Go (Linux, macOS, Windows)
**Project Type**: Single project - extending existing sdk/go/testing/ package
**Performance Goals**: Basic conformance < 60s, benchmark variance < 10%
**Constraints**: Must be backward-compatible with existing TestHarness and MockPlugin APIs
**Scale/Scope**: 4 test categories, 3 conformance levels, ~15-20 conformance tests total

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                           | Status | Evidence                                                        |
| ----------------------------------- | ------ | --------------------------------------------------------------- |
| I. Proto-First Development          | PASS   | Extends SDK generated from proto; no proto changes needed       |
| II. Multi-Provider Consistency      | PASS   | Tests validate cross-provider behavior via existing mocks       |
| III. Test-First Protocol            | PASS   | Feature IS a testing framework; tests validate the suite itself |
| IV. Protobuf Backward Compatibility | PASS   | No proto changes; SDK extension only                            |
| V. Comprehensive Documentation      | PASS   | Spec requires documentation and quickstart guide                |
| VI. Performance as Requirement      | PASS   | Performance benchmarks are core feature (FR-006, FR-007)        |
| VII. Validation at Multiple Levels  | PASS   | Suite implements conformance testing layer                      |

**Gate Result**: PASS - All principles satisfied. No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/011-plugin-conformance-suite/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (Go interfaces)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/testing/
├── harness.go           # Existing - TestHarness, validation functions
├── mock_plugin.go       # Existing - MockPlugin implementations
├── conformance.go       # NEW - ConformanceSuite, ConformanceResult, ConformanceLevel
├── spec_validation.go   # NEW - Spec validation tests (FR-001, FR-002, FR-003)
├── rpc_correctness.go   # NEW - RPC correctness tests (FR-004, FR-005)
├── performance.go       # NEW - Performance benchmark infrastructure (FR-006, FR-007)
├── concurrency.go       # NEW - Concurrency tests (FR-009, FR-010)
├── report.go            # NEW - JSON report generation (FR-013)
├── conformance_test.go  # Existing - Extended with new tests
├── benchmark_test.go    # Existing - Extended with baseline comparisons
└── integration_test.go  # Existing - Extended with suite self-tests
```

**Structure Decision**: Extend existing `sdk/go/testing/` package with new files for each test
category. Maintains backward compatibility while adding formal conformance structure.

## Complexity Tracking

> No violations - table not required.
