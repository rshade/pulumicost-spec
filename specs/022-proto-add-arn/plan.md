# Implementation Plan: Add ARN to GetActualCostRequest

**Branch**: `018-proto-add-arn` | **Date**: 2025-12-14 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/018-proto-add-arn/spec.md`

## Summary

Update `GetActualCostRequest` in `costsource.proto` to add an optional `arn` field (Canonical Cloud
Identifier). This addresses ambiguity with `resource_id` for cloud plugins like AWS Cost Explorer,
enabling them to robustly identify resources across services and contexts. The field will be
`optional string arn = 5;`.

## Technical Context

**Language/Version**: Go 1.25.5+, Protobuf 3
**Primary Dependencies**: `google.golang.org/protobuf`, `google.golang.org/grpc`
**Storage**: N/A (Proto definition)
**Testing**: `buf generate` (validation), Go unit tests, `buf breaking`
**Target Platform**: Cross-platform (Proto definition)
**Project Type**: Library/Spec
**Performance Goals**: Minimal overhead (standard protobuf field)
**Constraints**: Backward compatibility (optional field)
**Scale/Scope**: Core protocol change

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **I. gRPC Proto Specification-First Development**: Change starts with `costsource.proto`.
- [x] **II. Multi-Provider gRPC Consistency**: `arn` (or canonical ID) is applicable to AWS, Azure, GCP.
- [x] **III. Test-First Protocol**: Plan includes writing a test case in `sdk/go/testing` to verify
  the new field before full integration.
- [x] **IV. Protobuf Backward Compatibility**: Adding an optional field is non-breaking. `buf breaking` will verify.
- [x] **V. Comprehensive Documentation**: Will document the field in proto comments.
- [x] **VI. Performance**: Field number 5 is within the optimized 1-15 range.
- [x] **VII. Validation**: `buf lint` and `buf generate` will run in CI.

## Project Structure

### Documentation (this feature)

```text
specs/018-proto-add-arn/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (Proto snippet)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
proto/
└── finfocus/
    └── v1/
        └── costsource.proto  # Target file for modification

sdk/
└── go/
    ├── proto/                # Generated code location
    └── testing/              # Test case location
```

**Structure Decision**: Standard repository structure for `finfocus-spec`.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
| :--- | :--- | :--- |
| N/A | N/A | N/A |
