# Implementation Plan: PluginSDK Request Validation Helpers

**Branch**: `017-pluginsdk-validation` | **Date**: 2025-12-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md`
for the execution workflow.

## Summary

Implement `ValidateProjectedCostRequest` and `ValidateActualCostRequest` functions in
`sdk/go/pluginsdk/validation.go` to provide a standardized, high-performance validation layer for both Core and
Plugins. The system will ensure consistent error messages that guide developers to use specific `mapping` package
helpers (e.g., `mapping.ExtractAWSSKU`).

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.21+
**Primary Dependencies**: `google.golang.org/grpc`, `google.golang.org/protobuf` (existing in project)
**Storage**: N/A
**Testing**: `go test`, table-driven tests
**Target Platform**: Linux / Cross-platform (Go Library)
**Performance Goals**: <100ns execution time, zero heap allocations (SC-001)
**Constraints**: No external dependencies; must use existing `mapping` package for error guidance
**Scale/Scope**: Core shared library function, used by all plugins and core

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- **I. gRPC Proto Specification-First**: Passed. No proto changes required; implementing validation for existing
  protos.
- **II. Multi-Provider**: Passed. Validation logic is provider-agnostic (checks basic fields like Sku, Region
  presence) and delegates provider specifics to mapping helpers via error messages.
- **III. Test-First Protocol**: Passed. Plan includes SC-003 (100% test coverage) and follows TDD.
- **IV. Protobuf Backward Compatibility**: N/A.
- **V. Comprehensive Documentation**: Passed. Will include GoDoc and `quickstart.md`.
- **VI. Performance**: Passed. Explicit goal SC-001 (Zero allocation).
- **VII. Validation**: Passed. This feature _is_ the validation.

## Project Structure

### Documentation (this feature)

```text
specs/017-pluginsdk-validation/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── api.go.txt       # Proposed function signatures
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
sdk/go/pluginsdk/
├── validation.go        # Implementation of validation logic
└── validation_test.go   # Table-driven unit tests
```

**Structure Decision**: Extending the existing `pluginsdk` Go package to keep validation helpers close to the
mapping helpers they reference.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| None      |            |                                      |
