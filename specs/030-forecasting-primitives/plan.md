# Implementation Plan: Forecasting Primitives

**Branch**: `030-forecasting-primitives` | **Date**: 2025-12-30 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/030-forecasting-primitives/spec.md`

## Summary

Add forecasting primitives (`GrowthType` enum and `growth_rate` field) to enable forward-looking
cost projections with linear or exponential growth models. Fields are added to both
`ResourceDescriptor` (defaults) and `GetProjectedCostRequest` (overrides) to support flexible
projection scenarios while maintaining backward compatibility.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod) + Protocol Buffers v3
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless proto definitions)
**Testing**: go test, conformance tests via sdk/go/testing harness
**Target Platform**: gRPC service specification (cross-platform)
**Project Type**: gRPC proto specification repository
**Performance Goals**: Validation <100ms, projection calculations with 0.01% accuracy
**Constraints**: Backward compatible (optional fields only), buf breaking check must pass
**Scale/Scope**: Extends existing ResourceDescriptor and GetProjectedCostRequest messages

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. gRPC Proto Specification-First | PASS | Proto changes defined before SDK implementation |
| II. Multi-Provider Consistency | PASS | Growth fields are provider-agnostic |
| III. Test-First Protocol | PASS | Conformance tests will define expected behavior |
| IV. Protobuf Backward Compatibility | PASS | All new fields are optional, no removals |
| V. Comprehensive Documentation | PENDING | Proto comments required during implementation |
| VI. Performance Requirements | PASS | Validation <100ms specified in SC-005 |
| VII. Multi-Layer Validation | PASS | buf lint + schema + conformance planned |

**Gate Status**: PASS - No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/030-forecasting-primitives/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto definitions)
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
├── costsource.proto     # Extended: ResourceDescriptor, GetProjectedCostRequest
└── enums.proto          # New: GrowthType enum

sdk/go/
├── proto/               # Generated code (via buf generate)
├── pricing/             # SDK helpers (validation, growth calculation)
└── testing/             # Conformance tests for forecasting

examples/
└── requests/            # Sample gRPC request payloads with growth parameters
```

**Structure Decision**: Single project (gRPC spec repo). Proto changes in
`proto/finfocus/v1/`, generated SDK in `sdk/go/proto/`, helper code in `sdk/go/pricing/`.

## Complexity Tracking

> No Constitution Check violations requiring justification.
