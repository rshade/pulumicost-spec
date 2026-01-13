# Implementation Plan: Extend RecommendationActionType Enum

**Branch**: `019-recommendation-action-types` | **Date**: 2025-12-17 |
**Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/024-recommendation-action-types/spec.md`

## Summary

Extend the `RecommendationActionType` enum in the protobuf definition with 5 new action types
(MIGRATE, CONSOLIDATE, SCHEDULE, REFACTOR, OTHER) to provide comprehensive coverage of FinOps
platform recommendation categories. This is a backward-compatible, additive change that
requires proto update, SDK regeneration, and documentation.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless enum extension)
**Testing**: go test, buf lint, buf breaking, conformance tests
**Target Platform**: Go SDK generation, gRPC wire protocol
**Project Type**: Single (proto specification repository)
**Performance Goals**: N/A (enum lookup is O(1))
**Constraints**: Must maintain proto3 backward compatibility
**Scale/Scope**: 5 new enum values (positions 7-11)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | ✅ PASS | Proto changes come first, SDK regenerated |
| II. Multi-Provider Consistency | ✅ PASS | New types sourced from AWS/Azure/GCP/K8s platforms |
| III. Test-First Protocol | ✅ PASS | Conformance tests to be written before proto change |
| IV. Backward Compatibility | ✅ PASS | Additive enum extension, no breaking changes |
| V. Comprehensive Documentation | ✅ PASS | Proto comments required for each new value |
| VI. Performance Requirement | ✅ PASS | Enum extension has no performance impact |
| VII. Validation at Multiple Levels | ✅ PASS | buf lint + conformance tests validate |

**Gate Status**: PASSED - Proceed to Phase 0

## Project Structure

### Documentation (this feature)

```text
specs/019-recommendation-action-types/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── checklists/          # Quality checklists
    └── requirements.md  # Spec quality checklist (complete)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
└── costsource.proto     # RecommendationActionType enum (lines 613-621)

sdk/go/proto/finfocus/v1/
└── costsource.pb.go     # Generated Go code (regenerate with buf)

sdk/go/testing/
├── conformance_test.go  # Add tests for new action types
├── mock_plugin.go       # Update mock to support new types
└── integration_test.go  # Add integration tests
```

**Structure Decision**: This is a proto-first specification repository. Changes focus on
`proto/` definition with automatic SDK regeneration via `make generate`.
