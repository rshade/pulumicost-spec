# Implementation Plan: Anomaly Detection via Recommendations

**Branch**: `040-anomaly-detection-recommendations` | **Date**: 2026-01-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/040-anomaly-detection-recommendations/spec.md`

## Summary

Extend the existing `GetRecommendations` RPC to support cost anomaly detection by adding two enum
values: `RECOMMENDATION_CATEGORY_ANOMALY` (=5) and `RECOMMENDATION_ACTION_TYPE_INVESTIGATE` (=12).
This enables plugins to return anomalies as recommendations with an "investigate" action, creating a
unified "action items" view for FinOps practitioners.

**Technical Approach**: Proto-first development with enum additions, SDK code regeneration, and
documentation updates. No new messages, RPCs, or structural changes required.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod), Protocol Buffers v3
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless enum additions)
**Testing**: go test, conformance tests in sdk/go/testing/
**Target Platform**: Linux/macOS/Windows (cross-platform SDK)
**Project Type**: Single (gRPC specification repository with Go SDK)
**Performance Goals**: N/A (enum additions have zero runtime cost)
**Constraints**: Backward compatible; no breaking changes to existing proto wire format
**Scale/Scope**: 2 enum values added to existing proto file

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First**: Does this change start with Proto definitions?
  - YES: Changes begin with `proto/finfocus/v1/costsource.proto` enum additions
- [x] **Spec Consumes**: Does implementation avoid embedding complex pricing logic/calculators?
  - YES: Anomaly detection logic is in plugins/backends, not the spec. Spec only defines enum values.
- [x] **Multi-Provider**: Are examples/patterns provider-agnostic?
  - YES: ANOMALY category applies to AWS, Azure, GCP anomaly APIs uniformly
- [x] **FinFocus Alignment**: Does this respect the identity transition to FinFocus?
  - YES: Uses existing `finfocus.v1` package naming
- [x] **SDK Synchronization**: Are all SDKs (Go, TypeScript) updated when proto changes?
  - YES: Plan includes TypeScript SDK regeneration step
- [x] **Test-First Protocol**: Are conformance tests written first?
  - YES: Plan includes conformance test for anomaly recommendations
- [x] **Backward Compatibility**: Does buf breaking check pass?
  - YES: Adding enum values is backward compatible in proto3

## Project Structure

### Documentation (this feature)

```text
specs/040-anomaly-detection-recommendations/
├── spec.md              # Feature specification (completed)
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal - no unknowns)
├── data-model.md        # Phase 1 output (enum documentation)
├── quickstart.md        # Phase 1 output (plugin developer guide)
├── contracts/           # Phase 1 output (proto diff preview)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
└── costsource.proto     # Add ANOMALY and INVESTIGATE enum values

sdk/go/
├── proto/finfocus/v1/   # Regenerated Go bindings (make generate)
└── testing/             # Add conformance test for anomaly recommendations

sdk/typescript/
└── packages/client/     # Regenerated TypeScript bindings

docs/
└── sdk/                 # Update SDK documentation with anomaly guidance
```

**Structure Decision**: Minimal changes to existing structure. No new directories or files beyond
documentation. Uses established patterns for proto → SDK regeneration.

## Complexity Tracking

> No violations. This is a minimal additive change following established patterns.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| None      | N/A        | N/A                                  |
