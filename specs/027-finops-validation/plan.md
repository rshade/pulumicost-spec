# Implementation Plan: Contextual FinOps Validation

**Branch**: `027-finops-validation` | **Date**: 2025-12-25 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/027-finops-validation/spec.md`

## Summary

Extend the existing `ValidateFocusRecord` function in `sdk/go/pluginsdk/focus_conformance.go` to
include contextual business logic validation for FinOps cost records. This includes cost
relationship validation (EffectiveCost <= BilledCost <= ListCost), commitment discount field
consistency, capacity reservation consistency, and pricing model validation. The implementation
adds a new `ValidationError` structured type and supports both fail-fast (default) and
aggregate-all-errors modes while maintaining zero-allocation performance on the happy path.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc (existing)
**Storage**: N/A (stateless validation functions)
**Testing**: go test with stretchr/testify (existing pattern)
**Target Platform**: Plugin SDK library (cross-platform Go)
**Project Type**: SDK extension (single package modification)
**Performance Goals**: <100ns validation, 0 allocs/op on happy path (per SC-001)
**Constraints**: Backward compatibility with existing `ValidateFocusRecord` (per FR-008)
**Scale/Scope**: Single package extension (~200-300 lines of new code)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                          | Status | Notes                                               |
| ---------------------------------- | ------ | --------------------------------------------------- |
| I. Proto Specification-First       | PASS   | Uses existing FocusCostRecord proto, no proto changes needed |
| II. Multi-Provider Consistency     | PASS   | Validation rules apply uniformly across all providers |
| III. Test-First Protocol           | PASS   | Conformance tests will define expected behavior first |
| IV. Protobuf Backward Compatibility | PASS   | No proto changes, SDK-only enhancement |
| V. Comprehensive Documentation     | PASS   | Will update pluginsdk README and add examples |
| VI. Performance as Requirement     | PASS   | Zero-allocation target specified (SC-001) |
| VII. Validation at Multiple Levels | PASS   | Extends existing conformance test framework |

**Gate Status**: ✅ All gates pass. Proceeding to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/027-finops-validation/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - no new APIs)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── focus_conformance.go       # EXTEND: Add contextual validation rules
├── focus_conformance_test.go  # EXTEND: Add test cases for new rules
├── validation_error.go        # NEW: Structured ValidationError type
├── validation_error_test.go   # NEW: ValidationError tests
├── validation_options.go      # NEW: ValidationMode and options
└── validation_options_test.go # NEW: Options tests

sdk/go/testing/
├── focus13_conformance_test.go # EXTEND: Add contextual validation conformance tests
└── README.md                   # UPDATE: Document new validation capabilities
```

**Structure Decision**: Extend existing `sdk/go/pluginsdk/` package with new files for the
ValidationError type and validation options. Keep validation logic in `focus_conformance.go`
for cohesion with existing business rules.

## Complexity Tracking

> No constitution violations requiring justification.
