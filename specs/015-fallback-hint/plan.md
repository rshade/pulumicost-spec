# Implementation Plan: FallbackHint Enum

**Branch**: `001-fallback-hint` | **Date**: 2025-12-06 | **Spec**: [specs/001-fallback-hint/spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-fallback-hint/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.
See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature implements the `FallbackHint` enumeration in the `GetActualCostResponse` protobuf
definition to support plugin orchestration. It enables plugins to explicitly signal when the
core system should attempt to query other plugins (fallback) for a resource. This is critical
for handling "no data" scenarios gracefully and supporting specialized plugins. The
implementation involves updating `costsource.proto`, regenerating the Go SDK, and adding
functional options for hint setting.

## Technical Context

**Language/Version**: Go 1.21+ (SDK), Protobuf 3
**Primary Dependencies**: `google.golang.org/grpc`, `google.golang.org/protobuf`
**Storage**: N/A (API Specification)
**Testing**: Go `testing` package, `buf` (lint/breaking)
**Target Platform**: Cross-platform (Go)
**Project Type**: gRPC Specification & SDK
**Performance Goals**: Zero allocation overhead for default/unspecified hint values.
**Constraints**: strict backward compatibility; default behavior must match existing "no fallback" logic.
**Scale/Scope**: Core protocol change affecting all plugins.

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- **I. gRPC Proto Specification-First**: Passed. Spec defines proto changes first.
- **II. Multi-Provider gRPC Consistency**: Passed. Hint is provider-agnostic.
- **III. Test-First Protocol**: Passed. Test scenarios defined in spec.
- **IV. Protobuf Backward Compatibility**: Passed. New field, default 0 preserves behavior.
- **V. Comprehensive Documentation**: Passed. Spec requires doc updates.
- **VI. Performance**: Passed. Zero overhead for default value.
- **VII. Validation**: Passed. Buf checks included.

## Project Structure

### Documentation (this feature)

```text
specs/001-fallback-hint/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── costsource.proto # Updated proto definition
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
proto/
└── finfocus/
    └── v1/
        └── costsource.proto  # Primary modification target

sdk/
└── go/
    ├── proto/               # Generated code (do not edit manually)
    └── pluginsdk/           # SDK helpers to update
        └── helpers.go       # Functional options implementation
```

**Structure Decision**: Standard repository structure. Modifying existing files.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| None                       | -                  | -                                    |
