# Implementation Plan: FOCUS 1.2 Column Audit

**Branch**: `010-focus-column-audit` | **Date**: 2025-11-28 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/010-focus-column-audit/spec.md`

## Summary

Complete the FOCUS 1.2 column implementation by adding 19 missing columns (1 mandatory,
18 conditional) to `focus.proto`, updating the `FocusRecordBuilder` with corresponding
methods, creating an automated audit script, and providing comprehensive documentation
with 80%+ godoc coverage.

## Technical Context

**Language/Version**: Go 1.24+ (toolchain go1.25.4), Protocol Buffers v3
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (specification repository - no runtime storage)
**Testing**: go test, conformance tests via bufconn harness, buf lint/breaking
**Target Platform**: Cross-platform gRPC specification and Go SDK
**Project Type**: gRPC specification repository with generated SDK
**Performance Goals**: Zero-allocation enum validation (5-12 ns/op), conformance test pass
**Constraints**: Backward compatible proto changes, buf breaking check must pass
**Scale/Scope**: 57 total FOCUS columns, 19 new proto fields, ~20 new builder methods

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

### Principle I: gRPC Proto Specification-First Development ✅

- [x] Proto definitions updated before SDK code
- [x] Proto changes drive SDK generation via `make generate`
- [x] Comprehensive validation via buf lint

### Principle II: Multi-Provider gRPC Consistency ✅

- [x] New columns support all providers (AWS, Azure, GCP, Kubernetes)
- [x] FOCUS 1.2 is inherently provider-agnostic specification
- [x] Examples will demonstrate cross-provider mappings

### Principle III: Test-First Protocol ✅

- [x] Conformance tests will validate new column behavior
- [x] Tests defined before proto implementation
- [x] Red-Green-Refactor enforced

### Principle IV: Protobuf Backward Compatibility ✅

- [x] Adding new fields is backward compatible (non-breaking)
- [x] No field removals or type changes
- [x] buf breaking check will pass
- [x] Field numbers carefully assigned to avoid conflicts

### Principle V: Comprehensive Documentation ✅

- [x] 80%+ godoc coverage required (FR-010)
- [x] Proto comments with FOCUS section references (FR-007)
- [x] Developer guide and user reference planned (FR-012, FR-013)

### Principle VI: Performance as a gRPC Requirement ✅

- [x] Zero-allocation enum validation pattern maintained
- [x] Benchmarks for new enum types
- [x] No performance regression expected

### Principle VII: Validation at Multiple Levels ✅

- [x] Protobuf layer: buf lint/breaking
- [x] Service layer: conformance tests
- [x] SDK layer: integration tests
- [x] CI layer: all validation in GitHub Actions

**Constitution Gate Status**: ✅ PASSED - All principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/010-focus-column-audit/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto field definitions)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
├── focus.proto          # Add 19 missing FOCUS columns
└── enums.proto          # Add new enum types (CommitmentDiscountStatus, etc.)

sdk/go/
├── proto/               # Generated protobuf code (via buf generate)
│   └── finfocus/v1/
│       ├── focus.pb.go
│       └── enums.pb.go
└── pluginsdk/
    ├── focus_builder.go       # Add ~20 new With* methods
    ├── focus_builder_test.go  # Tests for new methods
    └── README.md              # Developer guide

docs/
└── focus-columns.md     # User-facing FOCUS column reference

examples/plugins/
└── focus_example.go     # Complete FOCUS record example

scripts/
└── audit_focus_columns.go  # Automated column audit script
```

**Structure Decision**: Single project structure - this is a specification repository
with proto definitions and generated SDK code. No web/mobile components.

## Complexity Tracking

> No constitution violations requiring justification.

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| N/A       | N/A        | N/A                                  |
