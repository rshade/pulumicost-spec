<!-- markdownlint-disable MD013 -->
# Implementation Plan: Add GetPluginInfo RPC

**Branch**: `029-plugin-info-rpc` | **Date**: 2025-12-30 | **Spec**: [specs/029-plugin-info-rpc/spec.md](spec.md)
**Input**: Feature specification from `/specs/029-plugin-info-rpc/spec.md`

## Summary

Add a `GetPluginInfo` RPC to the `CostSource` service. This RPC allows the core system to request metadata from a loaded plugin, specifically the `spec_version` it was built against, to enable compatibility verification, diagnostics, and graceful degradation. The implementation involves updating the `CostSource` proto definition, regenerating the Go SDK, and adding a default implementation in `pluginsdk` that returns compile-time constants (Name, Version, SpecVersion).

## Technical Context

**Language/Version**: Go 1.25.5+, Protobuf 3
**Primary Dependencies**: `google.golang.org/grpc`, `google.golang.org/protobuf`
**Storage**: N/A
**Testing**: Go standard library `testing`, `pluginsdk.test` framework
**Target Platform**: Linux/macOS/Windows (Cross-platform Go binaries)
**Project Type**: Single project (SDK + Proto definitions)
**Performance Goals**: < 100ms response time for `GetPluginInfo` (metadata retrieval should be instant)
**Constraints**: Must maintain backward compatibility for existing plugins (they will return "Unimplemented").
**Scale/Scope**: Core protocol change affecting all future plugins.

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **I. gRPC Proto Specification-First Development**: Plan starts with updating `costsource.proto`.
- [x] **II. Multi-Provider gRPC Consistency**: `GetPluginInfo` is provider-agnostic metadata.
- [x] **III. Test-First Protocol**: Plan includes creating conformance tests before implementation.
- [x] **IV. Protobuf Backward Compatibility**: Adding a new RPC is a non-breaking change. Existing clients can ignore it; existing servers will return Unimplemented.
- [x] **V. Comprehensive Documentation**: Plan includes updating documentation and examples.
- [x] **VI. Performance as a gRPC Requirement**: `GetPluginInfo` is a lightweight metadata call.
- [x] **VII. Validation at Multiple Levels**: Plan includes SDK validation for metadata fields.

## Project Structure

### Documentation (this feature)

```text
specs/029-plugin-info-rpc/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
# Single project structure
proto/
└── pulumicost/
    └── v1/
        └── costsource.proto # RPC definition

sdk/
└── go/
    ├── proto/               # Generated code
    ├── pluginsdk/           # SDK implementation
    │   ├── base.go          # Default handler
    │   └── version.go       # Version constants
    └── testing/             # Conformance tests
```

**Structure Decision**: Standard Go SDK structure within the existing monorepo.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
| :-------- | :--------- | :----------------------------------- |
| None      | N/A        | N/A                                  |
