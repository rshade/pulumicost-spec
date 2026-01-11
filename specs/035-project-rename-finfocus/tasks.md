---
description: "Task list for Project Rename to FinFocus"
---

# Tasks: Project Rename to FinFocus

**Input**: Design documents from `/specs/035-project-rename-finfocus/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: This feature uses validation commands (make test, make validate, etc.) rather
than unit tests. Test tasks are not included per the mechanical rename nature of the
work.

**Organization**: Tasks are grouped by user story to enable independent validation and incremental delivery.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto files**: `proto/finfocus/v1/`
- **Generated Go code**: `sdk/go/proto/finfocus/v1/`
- **Go SDK**: `sdk/go/`
- **CI/CD**: `.github/workflows/`
- **Documentation**: `README.md`, `docs/`
- **Root files**: `go.mod`, `release-please-config.json`, `.release-please-manifest.json`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project analysis and preparation

- [ ] T001 Verify current state: list all `.proto` files in proto/pulumicost/v1/
- [ ] T002 Verify current state: check go.mod module name is
      github.com/rshade/pulumicost-spec
- [ ] T003 Verify current state: identify all CI/CD workflow files in .github/workflows/
- [ ] T004 Backup current release-please configurations: copy release-please-config.json
      and .release-please-manifest.json to /tmp/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Prepare workspace for rename operations

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 Clean all generated code: run `make clean` to remove existing generated proto code
- [ ] T006 Create proto/finfocus directory structure: mkdir -p proto/finfocus/v1
- [ ] T007 Prepare release-please manifest for v0.5.0 reset: read
      .release-please-manifest.json to understand current version structure
- [ ] T008 Verify no uncommitted changes: run `git status` and ensure workspace is clean before starting rename

**Checkpoint**: Foundation ready - proto rename can begin

---

## Phase 3: User Story 1 - Protocol Branding (Priority: P1) 🎯 MVP

**Goal**: Rename all gRPC definitions to reflect the FinFocus brand (package
pulumicost.v1 → finfocus.v1)

**Independent Test**: Verify that `.proto` files use `package finfocus.v1` and gRPC
services are correctly named (services remain generic)

### Implementation for User Story 1

- [ ] T009 [P] [US1] Update package declaration in proto/pulumicost/v1/costsource.proto from `pulumicost.v1` to `finfocus.v1`
- [ ] T010 [P] [US1] Update go_package option in proto/pulumicost/v1/costsource.proto to `github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1;pbc`
- [ ] T011 [P] [US1] Update package declaration in proto/pulumicost/v1/budget.proto from `pulumicost.v1` to `finfocus.v1`
- [ ] T012 [P] [US1] Update go_package option in proto/pulumicost/v1/budget.proto to `github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1;pbc`
- [ ] T013 [US1] Move proto/pulumicost/v1/costsource.proto to proto/finfocus/v1/costsource.proto
- [ ] T014 [US1] Move proto/pulumicost/v1/budget.proto to proto/finfocus/v1/budget.proto
- [ ] T015 [US1] Remove empty proto/pulumicost/v1/ directory: rmdir proto/pulumicost/v1 proto/pulumicost
- [ ] T016 [P] [US1] Update buf.yaml to reference proto/finfocus/v1 instead of proto/pulumicost/v1 (if applicable)
- [ ] T017 [P] [US1] Update .github/workflows/proto.yml to use proto/finfocus/v1
      paths for buf lint/breaking checks
- [ ] T018 [US1] Regenerate proto code: run `make generate` to populate sdk/go/proto/finfocus/v1/
- [ ] T019 [US1] Verify proto regeneration: check that sdk/go/proto/finfocus/v1/
      directory is populated with generated Go code
- [ ] T020 [US1] Verify service names remain generic: check sdk/go/proto/finfocus/v1/
      generated code confirms CostSourceService, BudgetService are not renamed to
      FinFocus\*Service

**Checkpoint**: At this point, Protocol Branding (User Story 1) is complete and verifiable

---

## Phase 4: User Story 2 - SDK Module Transition (Priority: P2)

**Goal**: Update Go module to github.com/rshade/finfocus-spec and update all internal imports

**Independent Test**: Run `go test ./...` and verify all imports use the new
finfocus-spec path (zero pulumicost-spec imports remain)

### Implementation for User Story 2

- [ ] T021 [US2] Update go.mod module declaration from
      `github.com/rshade/pulumicost-spec` to `github.com/rshade/finfocus-spec`
- [ ] T022 [P] [US2] Update internal imports in sdk/go/pulumiapi/\*.go from
      `github.com/rshade/pulumicost-spec/sdk/go/` to
      `github.com/rshade/finfocus-spec/sdk/go/`
- [ ] T023 [P] [US2] Update internal imports in sdk/go/pluginsdk/\*.go from
      `github.com/rshade/pulumicost-spec/sdk/go/` to
      `github.com/rshade/finfocus-spec/sdk/go/`
- [ ] T024 [P] [US2] Update internal imports in sdk/go/internal/\*.go from
      `github.com/rshade/pulumicost-spec/sdk/go/` to
      `github.com/rshade/finfocus-spec/sdk/go/`
- [ ] T025 [P] [US2] Update proto import in sdk/go/\*_/_\_test.go from
      `pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"` to
      `pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"`
- [ ] T026 [P] [US2] Update examples/\*_/_.go imports from
      `github.com/rshade/pulumicost-spec/sdk/go/` to
      `github.com/rshade/finfocus-spec/sdk/go/`
- [ ] T027 [US2] Run go mod tidy to resolve dependencies: `go mod tidy` at repository root
- [ ] T028 [US2] Verify module rename: run `go list -m` and confirm output is
      `github.com/rshade/finfocus-spec`
- [ ] T029 [US2] Search for remaining pulumicost-spec imports: run
      `rg "github.com/rshade/pulumicost-spec"` and verify only historical
      references remain
- [ ] T030 [US2] Verify tests pass: run `make test` and confirm all tests pass
      with new module path

**Checkpoint**: At this point, User Story 2 should be complete and all tests pass

---

## Phase 5: User Story 3 - Automated Release Management (Priority: P3)

**Goal**: Update release-please configuration and CI/CD workflows for the new project name

**Independent Test**: Run `release-please` in dry-run mode and verify the proposed version (v0.5.0) and manifest updates

### Implementation for User Story 3

- [ ] T031 [P] [US3] Update package-name in release-please-config.json from `pulumicost-spec` to `finfocus-spec`
- [ ] T032 [US3] Reset root version in .release-please-manifest.json to `0.5.0` (complete reset as per research decision)
- [ ] T033 [P] [US3] Update .github/workflows/ci.yml to use new module name in
      cache keys (e.g., go-mod-cache: github.com/rshade/finfocus-spec)
- [ ] T034 [P] [US3] Update .github/workflows/release-please.yml to reference
      finfocus-spec package name
- [ ] T035 [P] [US3] Update any other workflow files in .github/workflows/ that
      reference the old module name
- [ ] T036 [US3] Test release-please dry-run: verify proposed version is v0.5.0
      and configuration is valid
- [ ] T037 [US3] Verify CI/CD workflow paths: check that buf lint/breaking
      actions reference proto/finfocus/v1

**Checkpoint**: At this point, User Story 3 is complete and release configuration is updated

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Global text replacement, documentation updates, and final validation

- [ ] T038 [P] Perform global case-sensitive find/replace: `pulumicost` →
      `finfocus` (lowercase) across entire repository (excluding historical
      files like specs/)
- [ ] T039 [P] Perform global case-sensitive find/replace: `PulumiCost` →
      `FinFocus` (PascalCase) across entire repository (excluding historical
      files like specs/)
- [ ] T040 [P] Perform global case-sensitive find/replace: `PULUMICOST` →
      `FINFOCUS` (UPPERCASE) across entire repository (excluding historical
      files like specs/)
- [ ] T041 Update project tagline in README.md to: "FinFocus: Focusing your
      finances left"
- [ ] T042 [P] Update docs/ directory: search and replace brand references while
      preserving historical context (contextual replacement per research
      decision)
- [ ] T043 [P] Update CLAUDE.md references to pulumicost-spec where appropriate (preserve historical notes)
- [ ] T044 [P] Update GEMINI.md references to pulumicost-spec where appropriate (preserve historical notes)
- [ ] T045 [P] Update AGENTS.md references to pulumicost-spec where appropriate (preserve historical notes)
- [ ] T046 [P] Update any other documentation files (docs/\*.md) with FinFocus branding (contextual replacement)
- [ ] T047 Final validation: run `make generate` and confirm clean proto code generation
- [ ] T048 Final validation: run `make test` and confirm all tests pass
- [ ] T049 Final validation: run `make validate` (linting and schema
      validation) and confirm all checks pass
- [ ] T050 [P] Search for remaining pulumicost brand references: run
      `rg -i "pulumicost"` and verify only historical context remains (in
      specs/, docs with legacy notes)
- [ ] T051 [P] Search for remaining PulumiCost brand references: run
      `rg "PulumiCost"` and verify only historical context remains
- [ ] T052 Verify wire protocol compatibility: run `make generate` and confirm
      generated proto files maintain service and message structure
- [ ] T053 Validate success criteria: verify SC-001 through SC-005 from spec.md
      are met
- [ ] T054 Run markdownlint on all updated \*.md files to ensure markdown compliance

**Checkpoint**: All rename work is complete, validated, and ready for commit

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational - Must complete before US2
- **User Story 2 (Phase 4)**: Depends on US1 completion - Must complete before US3
- **User Story 3 (Phase 5)**: Depends on US2 completion
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

**Critical Dependency Chain**: US1 → US2 → US3 (Sequential)

- **User Story 1 (P1)**: Protocol Branding
  - Can start after Foundational (Phase 2)
  - Must complete before US2 (module rename depends on proto package being renamed first)
  - Independently testable: Verify proto files use `finfocus.v1` package

- **User Story 2 (P2)**: SDK Module Transition
  - **BLOCKS ON**: US1 completion (proto must be renamed before Go module)
  - Must complete before US3 (release config should be done after module rename)
  - Independently testable: Run `go test ./...` and verify zero pulumicost-spec imports

- **User Story 3 (P3)**: Automated Release Management
  - **BLOCKS ON**: US2 completion (module path must be updated before release config)
  - Independently testable: Run `release-please` dry-run and verify v0.5.0

### Within Each User Story

- Proto file updates can run in parallel (T009-T012)
- Proto file moves must be sequential (T013-T015)
- Internal import updates can run in parallel (T022-T026)
- Release config updates can run in parallel (T031, T033-T035)
- Global text replacement can run in parallel (T038-T046)
- Final validation must be sequential (T047-T054)

### Parallel Opportunities

- **Setup phase**: T001-T004 can run in parallel (verification tasks)
- **Foundational phase**: T006-T007 can run in parallel
- **User Story 1**: T009-T012 (proto file updates) can run in parallel; T016-T017 (config updates) can run in parallel
- **User Story 2**: T022-T026 (import updates) can run in parallel
- **User Story 3**: T031, T033-T035 (config updates) can run in parallel
- **Polish phase**: T038-T046 (global replacements) can run in parallel; T050-T051 (search verification) can run in parallel

---

## Parallel Example: User Story 2 (SDK Module Transition)

```bash
# Launch all import updates together (parallel):
Task: "Update internal imports in sdk/go/pulumiapi/*.go"
Task: "Update internal imports in sdk/go/pluginsdk/*.go"
Task: "Update internal imports in sdk/go/internal/*.go"
Task: "Update proto import in sdk/go/**/*_test.go"
Task: "Update examples/**/*.go imports"
```

---

## Parallel Example: Polish Phase (Global Replacement)

```bash
# Launch all global replacements together (parallel):
Task: "Perform global case-sensitive find/replace: pulumicost →
      finfocus"
Task: "Perform global case-sensitive find/replace: PulumiCost →
      FinFocus"
Task: "Perform global case-sensitive find/replace: PULUMICOST →
      FINFOCUS"
Task: "Update docs/ directory with FinFocus branding"
Task: "Update CLAUDE.md references to finfocus-spec"
Task: "Update GEMINI.md references to finfocus-spec"
Task: "Update AGENTS.md references to finfocus-spec"
Task: "Update any other documentation files (docs/*.md)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

**Note**: Due to the sequential dependency chain (US1 → US2 → US3), a true MVP
delivering the full rename requires completing all three user stories. However,
User Story 1 can be validated independently.

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Protocol Branding)
4. **STOP and VALIDATE**: Verify proto files use `finfocus.v1` package
5. Continue with US2 and US3 to complete full rename

### Incremental Delivery

Due to the mechanical rename nature of this feature, the entire feature (US1 + US2 + US3) should be delivered together:

1. Complete Setup + Foundational → Foundation ready
2. Complete US1 (Protocol Branding) → Verify proto rename
3. Complete US2 (SDK Module Transition) → Verify module rename
4. Complete US3 (Release Management) → Verify release config
5. Complete Polish → Final validation and documentation updates
6. **Single PR** containing the entire rename work

**Rationale**: The rename is atomic - partially renaming (e.g., proto but not module)
would leave the repository in an inconsistent, broken state.

### Sequential Team Strategy

With multiple developers, the dependency chain limits parallelism:

1. Team completes Setup + Foundational together (Phase 1-2)
2. Once Foundational is done:
   - All developers focus on US1 (Protocol Branding) - parallelize proto
     file updates
   - Move together to US2 (SDK Module Transition) - parallelize import
     updates
   - Move together to US3 (Release Management) - parallelize config
     updates
   - Move together to Polish - parallelize global replacements
3. Stories complete sequentially due to dependency constraints

---

## Notes

- **Sequential Dependency**: US1 → US2 → US3 is required by the nature of Go
  module/protobuf renaming (proto must be renamed before Go module, which
  must be renamed before release config)
- **[P] tasks**: Different files, no dependencies, can run in parallel within
  a user story
- **[Story] label**: Maps task to specific user story for traceability
- **No Partial Delivery**: Do not deliver partially renamed state; entire rename
  should be atomic
- **Historical Context Preservation**: When doing global replacements, preserve
  historical references in docs/ and specs/ (contextual replacement per
  research decision)
- **Validation at Each Checkpoint**: Run validation commands after each phase to
  catch issues early
- **Commit Strategy**: Consider committing after each user story phase for better
  rollback capability, but do not merge until all phases are complete
- **Breaking Change**: This is a breaking change at the module/package level
  (justified v0.5.0), but wire protocol semantics remain compatible
- **Copyright Headers**: Add/update Apache 2.0 copyright headers to all source
  files per constitution (Go, Proto, Script, Schema)
