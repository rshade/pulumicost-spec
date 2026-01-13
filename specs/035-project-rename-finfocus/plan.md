# Implementation Plan: Project Rename to FinFocus

**Branch**: `035-project-rename-finfocus` | **Date**: 2026-01-11 | **Spec**: /specs/035-project-rename-finfocus/spec.md
**Input**: Feature specification from `/specs/035-project-rename-finfocus/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See
`.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature implements a comprehensive rename of the project from FinFocus to
FinFocus, aligning the repository with the FinOps FOCUS specification. The
technical approach involves:

1. **Protocol rename**: Updating all protobuf package declarations from `finfocus.v1` to `finfocus.v1`
2. **Module rename**: Changing the Go module from `github.com/rshade/finfocus-spec` to `github.com/rshade/finfocus-spec`
3. **Global text replacement**: Performing case-sensitive find/replace of branding references
4. **Release management**: Updating `release-please` configuration and resetting version to v0.5.0
5. **CI/CD updates**: Ensuring all workflows reference the new module paths

The rename preserves the gRPC wire protocol's semantic structure while changing
only the package namespace, maintaining backward compatibility at the service
level while introducing a breaking change at the module/package level (justifying
the version bump to v0.5.0).

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: gRPC, protobuf, buf v1.32.1
**Storage**: N/A (SDK does not manage persistent storage)
**Testing**: go test, make test
**Target Platform**: Linux server (gRPC service specification)
**Project Type**: single (SDK/specification project)
**Performance Goals**: Zero-allocation goal for common operations (validation, enum lookups)
**Constraints**: Protobuf backward compatibility, gRPC wire format compatibility, Apache 2.0 license compliance
**Scale/Scope**: Protocol specification and Go SDK (~1M LOC across ecosystem)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First**: Does this change start with Proto definitions?
  - **Status**: PASS - The rename begins with updating `.proto` package declarations from
    `finfocus.v1` to `finfocus.v1`, consistent with the "Proto definitions are the
    source of truth" principle.
- [x] **Spec Consumes**: Does implementation avoid embedding complex pricing logic/calculators?
  - **Status**: PASS - This is a purely mechanical rename; no pricing logic or calculations are added or modified.
- [x] **Multi-Provider**: Are examples/patterns provider-agnostic?
  - **Status**: PASS - The rename does not affect the provider-agnostic nature of
    existing examples; all provider-specific patterns are preserved.
- [x] **FinFocus Alignment**: Does this respect the identity transition to FinFocus?
  - **Status**: PASS - This feature is the primary implementation of the FinFocus
    transition, as documented in the constitution's "Comprehensive Documentation &
    Identity Transition" section.
- [x] **Backward Compatibility**: Are breaking changes justified and properly versioned?
  - **Status**: PASS - The rename is a breaking change at the module/package level
    (justifying the version bump to v0.5.0). The gRPC wire protocol semantics
    remain unchanged, only the package namespace is updated.
- [x] **Documentation Complete**: Will documentation be updated in the same PR?
  - **Status**: PASS - All references in README.md, docs/, and inline comments will be updated to reflect the FinFocus brand.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
proto/
├── finfocus/v1/         # Renamed from finfocus/v1/
│   ├── costsource.proto
│   └── budget.proto
├── buf.yaml
└── buf.lock

sdk/go/
├── proto/finfocus/v1/   # Generated Go code from proto files
├── internal/
├── pulumiapi/
├── pluginsdk/
├── go.mod               # Will be updated to github.com/rshade/finfocus-spec
└── README.md

examples/
├── pricing/
│   └── aws/
├── budget/
│   └── aws/
└── validation/

.github/
├── workflows/
│   ├── ci.yml
│   ├── release-please.yml
│   └── proto.yml
└── workflows/

docs/
├── DESIGN.md
├── DEVELOPMENT.md
└── [additional documentation]

root files:
├── go.mod              # Module rename to finfocus-spec
├── release-please-config.json       # Updated project name
├── .release-please-manifest.json    # Reset to v0.5.0
├── README.md           # Updated branding and tagline
├── CLAUDE.md
├── GEMINI.md
└── AGENTS.md
```

**Structure Decision**: Single project structure (SDK/specification repository). The rename
maintains the existing directory layout while updating package names and module
paths. All directories remain in place; only package identifiers and references
are updated.

## Constitution Check (Post-Design)

_GATE: Re-evaluated after Phase 1 design completion._

- [x] **Contract First**: Does this change start with Proto definitions?
  - **Status**: PASS - Design confirms that protobuf package updates occur first in
    the rename order of operations.
- [x] **Spec Consumes**: Does implementation avoid embedding complex pricing logic/calculators?
  - **Status**: PASS - This is a mechanical rename; no business logic changes are introduced.
- [x] **Multi-Provider**: Are examples/patterns provider-agnostic?
  - **Status**: PASS - The rename preserves all existing provider-agnostic examples and patterns.
- [x] **FinFocus Alignment**: Does this respect the identity transition to FinFocus?
  - **Status**: PASS - The design implements the complete identity transition as
    specified in the constitution.
- [x] **Backward Compatibility**: Are breaking changes justified and properly versioned?
  - **Status**: PASS - The package/module rename is a breaking change justifying v0.5.0. Wire protocol remains compatible.
- [x] **Documentation Complete**: Will documentation be updated in the same PR?
  - **Status**: PASS - Design includes comprehensive documentation updates including
    quickstart guide, data model, and API contracts.

**Overall Constitution Compliance**: ✅ PASS - All constitutional gates satisfied. No violations requiring justification.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| N/A       | N/A        | N/A                                  |
