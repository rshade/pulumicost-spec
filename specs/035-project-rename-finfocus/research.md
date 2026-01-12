# Research: Project Rename to FinFocus

**Date**: 2026-01-11
**Feature**: 035-project-rename-finfocus

## Overview

This document consolidates research findings for the project rename from FinFocus to
FinFocus. All clarifications from the feature specification have been resolved.

## Research Decisions

### Decision 1: Service Naming Convention

**Decision**: Keep generic service names (e.g., `CostSourceService`) in the new `finfocus.v1` package.

**Rationale**:

- Branding should be limited to the package namespace, not individual service implementations
- Generic service names improve discoverability and reduce cognitive load for developers
- Aligns with the "Multi-Provider gRPC Consistency" principle in the constitution
- Consistent with the established pattern in the existing codebase

**Alternatives Considered**:

- Option B: Rename services to include "FinFocus" (e.g., `FinFocusCostSourceService`)
  - Rejected: Creates verbose, redundant naming (`finfocus.v1.FinFocusCostSourceService`)

### Decision 2: Version for Rename Release

**Decision**: Use v0.5.0 for the rename release.

**Rationale**:

- Signals a significant breaking change (module/package rename)
- Provides clear semantic versioning distinction from previous releases
- Aligns with the requirement to reset the current package version to 0.5.0

**Alternatives Considered**:

- Option A: v0.2.0 (next sequential patch)
  - Rejected: Underestimates the impact of a module-level breaking change
- Option B: v1.0.0
  - Rejected: The project is still in pre-release phase (< 1.0.0)

### Decision 3: Go Package Alias for Generated Protos

**Decision**: Keep the "pbc" alias for the generated proto package.

**Rationale**:

- Minimal churn for plugin implementations
- "pbc" stands for "protobuf client" and remains semantically valid
- Reduces risk of merge conflicts during the transition

**Alternatives Considered**:

- Option A: Rename to "ff" (finfocus)
  - Rejected: Requires updating all plugin implementations, increases scope

### Decision 4: Release-Please Manifest Update Strategy

**Decision**: Reset current package version to 0.5.0 in the `.release-please-manifest.json`.

**Rationale**:

- Provides a clean slate for the new project identity
- Ensures `release-please` generates the correct next version (v0.5.0 or v0.6.0)
- Aligns with the version decision above

**Alternatives Considered**:

- Option B: Append entry with old name and version
  - Rejected: Creates confusing history and complicates version tracking

### Decision 5: Documentation Replacement Scope

**Decision**: Contextual replacement - preserve historical context and legacy notes where appropriate.

**Rationale**:

- Maintains project history and acknowledgments
- Prevents rewriting of the past
- Allows for smooth transition for existing users

**Examples of Contextual Replacement**:

- In CHANGELOG.md: Keep historical entries mentioning "FinFocus" unchanged
- In docs/: Preserve legacy notes like "formerly known as FinFocus" in migration guides
- In README.md: Update active branding but keep acknowledgments of original project name
- In specs/: Leave historical references in feature specifications unchanged

**Alternatives Considered**:

- Option A: Replace all occurrences
  - Rejected: Erases project history and creates confusion

## Technical Approach Research

### Global Rename Strategy

We will perform a contextual find/replace of "FinFocus" and "finfocus" across the
codebase, documentation, and CI/CD configurations.

**Replacements Map**:

| Category      | Search Term                         | Replace With                            | Notes                              |
| ------------- | ----------------------------------- | --------------------------------------- | ---------------------------------- |
| Proto Package | `finfocus.v1`                     | `finfocus.v1`                           | Wire breaking change               |
| Go Module     | `github.com/rshade/finfocus-spec` | `github.com/rshade/finfocus-spec`       | Import breaking change             |
| Brand (Title) | `FinFocus`                        | `FinFocus`                              |                                    |
| Path Alias    | `pbc`                               | `pbc`                                   | Preserved for internal consistency |
| Directory     | `proto/finfocus`                  | `proto/finfocus`                        |                                    |
| Tagline       | `Production-ready specification...` | `FinFocus: Focusing your finances left` |                                    |

**Rationale**: To establish the "FinFocus" identity, a complete transition is required for
the protocol package and Go module. Contextual replacement ensures we maintain
historical accuracy in specific files (like the rename plan itself) while providing
a clean brand for all active development.

### Rename Order of Operations

Based on best practices for Go module and protobuf package renames, the optimal order is:

1. **Protobuf package rename** (first)
   - Update `package finfocus.v1` to `package finfocus.v1` in all `.proto` files
   - Move `proto/finfocus/` directory to `proto/finfocus/`
   - Update `go_package` options to reference new paths

2. **Go module rename**
   - Update `go.mod` module declaration to `github.com/rshade/finfocus-spec`
   - Update all internal imports to match the new module path

3. **Global text replacement**
   - Case-sensitive find/replace of `finfocus` → `finfocus` (lowercase)
   - Case-sensitive find/replace of `FinFocus` → `FinFocus` (PascalCase)
   - Case-sensitive find/replace of `PULUMICOST` → `FINFOCUS` (UPPERCASE)
   - Preserve historical context in documentation

4. **CI/CD and release configuration**
   - Update `release-please-config.json` and `.release-please-manifest.json`
   - Update workflow files (`.github/workflows/*.yml`)
   - Update documentation (README.md, docs/)

5. **Regeneration and validation**
   - Run `make clean` and `make generate` to regenerate proto code
   - Run `make test` to verify all tests pass
   - Run `make validate` to ensure linting and schema validation pass

### Release Management Analysis

**Release-Please Configuration**: We will reset the manifest to `v0.5.0`.

- **File**: `.release-please-manifest.json`
- **Action**: Set root path version to `0.5.0`.
- **File**: `release-please-config.json`
- **Action**: Ensure package name is updated to `finfocus-spec`.

### CI/CD Impact

All workflows in `.github/workflows/` must be updated:

- Update Go module name in cache keys.
- Update paths in `buf` lint/breaking actions.
- Update repository references.

### Breaking Change Justification

The rename constitutes a **breaking change** at the following levels:

1. **Go module path**: External consumers must update their `go.mod` files to use the new module path
2. **Protobuf package namespace**: Generated code changes package names, requiring
   plugin implementations to update imports
3. **Wire protocol compatibility**: While the gRPC wire format remains unchanged, the
   package namespace change affects client code generation

This justifies the version bump to v0.5.0 (MINOR version in pre-1.0 semantic versioning).

### Risk Assessment and Mitigation

| Risk                                       | Likelihood | Impact | Mitigation                                                         |
| ------------------------------------------ | ---------- | ------ | ------------------------------------------------------------------ |
| Missed imports causing build failures      | Medium     | High   | Run `go test ./...` and `make validate` before committing          |
| Documentation drift (old brand references) | High       | Low    | Use `rg` (ripgrep) to search for remaining `finfocus` references |
| Release-please version confusion           | Low        | Medium | Test with `release-please` dry-run mode before merging             |
| Generated code conflicts                   | Low        | High   | Ensure `make clean` is run before `make generate`                  |

## Summary

All clarifications from the feature specification have been resolved. The technical
approach is straightforward mechanical refactoring with clear justification for the
breaking change (module/package rename). The rename order of operations ensures minimal
disruption and validates changes incrementally.

No additional research is required; Phase 1 design can proceed.
