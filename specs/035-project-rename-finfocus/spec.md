# Feature Specification: Project Rename to FinFocus

**Feature Branch**: `035-project-rename-finfocus`  
**Created**: 2026-01-11  
**Status**: Draft  
**Input**: User description: "Rename the project to FinFocus and update release-please for this transition."

## User Scenarios & Testing _(mandatory)_

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE.
-->

## Clarifications

### Session 2026-01-11

- Q: Should services be renamed to include the "FinFocus" brand? → A: Keep generic
  names (e.g., CostSourceService) in the new package (Option A).
- Q: What version should be used for the rename release? → A: v0.5.0 (C).
- Q: What Go package alias should be used for generated protos? → A: Keep "pbc" alias (Option B).
- Q: How should the release-please manifest be updated? → A: Reset current package version to 0.5.0 (Option A).
- Q: What is the scope for documentation replacements? → A: Contextual replacement (Keep history/legacy notes) (Option B).

### User Story 1 - Protocol Branding (Priority: P1)

As a protocol consumer, I want all gRPC definitions to reflect the FinFocus brand so that
the wire format and client code are consistent with the new identity.

**Why this priority**: High. This is the foundation of the rename and impacts every downstream component.

**Independent Test**: Can be tested by verifying that `.proto` files use `package
finfocus.v1` and gRPC services are correctly named.

**Acceptance Scenarios**:

1. **Given** the proto directory, **When** examining `costsource.proto`, **Then** the package MUST be `finfocus.v1`.
2. **Given** the proto directory, **When** examining service definitions, **Then**
   service names MUST remain generic (e.g., `CostSourceService`).

---

### User Story 2 - SDK Module Transition (Priority: P2)

As a Go developer, I want to import the SDK using the `finfocus-spec` module path so that
my project dependencies reflect the current project name.

**Why this priority**: High. Prevents confusion and aligns the codebase with the new repository name.

**Independent Test**: Can be tested by running `go test ./...` and verifying all imports use the new `finfocus-spec` path.

**Acceptance Scenarios**:

1. **Given** `go.mod`, **When** the module is updated, **Then** all internal imports MUST
   be updated to match.
2. **Given** a plugin implementation, **When** importing the generated proto, **Then**
   it MUST use the `finfocus` package path.

---

### User Story 3 - Automated Release Management (Priority: P3)

As a project maintainer, I want `release-please` to continue managing our releases
seamlessly after the rename so that we don't lose automated changelogs and versioning.

**Why this priority**: Medium. Essential for long-term maintenance and CI/CD stability.

**Independent Test**: Can be tested by running `release-please` in dry-run mode and
verifying the proposed version and manifest updates.

**Acceptance Scenarios**:

1. **Given** the repository root, **When** `release-please-config.json` is updated,
   **Then** it MUST reference the new package name and repository.
2. **Given** a new commit, **When** CI runs, **Then** `release-please` MUST correctly
   identify the next version (v0.2.0 or v0.1.0).

---

### Edge Cases

- **Mixed Imports**: Ensuring no `pulumicost` imports remain in the `sdk/go` directory.
- **Documentation Drift**: Ensuring `README.md` and `docs/` are fully updated to prevent branding inconsistency.
- **Generated Code Persistence**: Ensuring `make clean` and `make generate` work correctly with the new structure.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST rename all Protobuf packages from `pulumicost.v1` to `finfocus.v1`.
- **FR-002**: System MUST update `go.mod` to `github.com/rshade/finfocus-spec`.
- **FR-003**: System MUST perform a global case-insensitive find/replace of
  `pulumicost` to `finfocus` (respecting casing: `PulumiCost` -> `FinFocus`).
- **FR-004**: System MUST move `proto/pulumicost/` directory to `proto/finfocus/`.
- **FR-005**: System MUST update `release-please-config.json` and
  `.release-please-manifest.json` to reflect the new project name and version
  (v0.5.0).
- **FR-006**: System MUST update all CI/CD workflows (`.github/workflows/*.yml`) to use
  the new module name and paths.
- **FR-007**: System MUST update the Go SDK `go_package` options in `.proto` files to
  `github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1;pbc`.
- **FR-008**: System MUST update the project tagline in the root `README.md` to:
  _"FinFocus: Focusing your finances left"_.

### Key Entities

- **Protobuf Contract**: The wire protocol definition transitioning from `pulumicost.v1` to `finfocus.v1`.
- **Go SDK Module**: The Go module identity in `go.mod`.
- **Release Manifest**: The `release-please` state files tracking versions.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: 100% of `.proto` files use the `finfocus.v1` package name.
- **SC-002**: `make generate` completes without errors and populates `sdk/go/proto/finfocus/v1`.
- **SC-003**: `make test` passes with zero remaining imports of `github.com/rshade/pulumicost-spec`.
- **SC-004**: All occurrences of the string "PulumiCost" (case-insensitive) in the
  `docs/` directory are replaced with "FinFocus", unless referring to historical
  context.
- **SC-005**: `release-please` configuration is valid and points to the new repository and package names.
