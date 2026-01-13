# Tasks: Add Migration Documentation

**Input**: Design documents from `/specs/001-migration-docs/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: No tests requested in specification - focusing on documentation creation and validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Documentation**: Root-level files (MIGRATION.md, CHANGELOG.md, README.md, llm-migration.json)
- Paths shown are absolute for clarity

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for documentation work

- [x] T001 Install markdownlint-cli2 for validation in package.json
- [x] T002 Verify JSON validation tools are available

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

N/A for documentation feature - no blocking prerequisites beyond basic tool setup.

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Environment Variable Migration Guide (Priority: P1) 🎯 MVP

**Goal**: Create comprehensive documentation for environment variable renames from PULUMICOST*\* to FINFOCUS*\*.

**Independent Test**: Can be fully tested by verifying that all 6 environment variable
renames are documented with before/after examples in MIGRATION.md and users can successfully
update their configurations.

### Implementation for User Story 1

- [x] T003 [US1] Create MIGRATION.md with overview section (use markdown table format for
  variable mappings)
- [x] T004 [US1] Add environment variables mapping table to MIGRATION.md
- [x] T005 [US1] Add migration steps section to MIGRATION.md including variable updates
- [x] T006 [US1] Add backwards compatibility note for SDK support to MIGRATION.md

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Plugin Directory Migration (Priority: P1)

**Goal**: Document plugin directory migration from ~/.pulumicost/plugins/ to ~/.finfocus/plugins/.

**Independent Test**: Can be fully tested by verifying that the migration command correctly
moves plugins from the old directory to the new one and that plugins are subsequently
discovered.

### Implementation for User Story 2

- [x] T007 [US2] Add plugin discovery paths section to MIGRATION.md
- [x] T008 [US2] Add plugin migration command example to MIGRATION.md
- [x] T009 [US2] Add plugin configuration file update guidance to MIGRATION.md

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - LLM-Friendly Migration Manifest (Priority: P2)

**Goal**: Create machine-readable JSON manifest for AI-assisted migration across repositories.

**Independent Test**: Can be fully tested by verifying that an AI assistant can parse the
manifest and correctly identify all required changes without human intervention.

### Implementation for User Story 3

- [x] T010 [US3] Create llm-migration.json with project rename section
- [x] T011 [US3] Add environment variables renames array to llm-migration.json
- [x] T012 [US3] Add plugin path mapping to llm-migration.json
- [x] T013 [US3] Add migration steps array to llm-migration.json
- [x] T014 [US3] Add Go SDK module path update note to llm-migration.json

**Checkpoint**: At this point, User Stories 1, 2 AND 3 should all work independently

---

## Phase 6: User Story 4 - Changelog and README Updates (Priority: P2)

**Goal**: Integrate migration information into existing documentation (CHANGELOG, README) for discoverability.

**Independent Test**: Can be fully tested by verifying that users can find migration
guidance from the CHANGELOG entry for the rename release and from the README support
section.

### Implementation for User Story 4

- [x] T015 [US4] Add migration section to CHANGELOG.md v0.5.0 release notes
- [x] T016 [US4] Add migration guide link to README.md Support section
- [x] T017 [US4] Update README.md overview to reference migration guide

**Checkpoint**: All user stories should now be independently functional

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T018 [P] Run markdown linting on all modified files
- [x] T019 [P] Validate llm-migration.json syntax
- [x] T020 [P] Verify all migration documentation is discoverable within 3 navigation steps
- [x] T021 [P] Test that migration can be completed in under 5 minutes following steps
- [x] T022 [P] Validate JSON manifest against schema
- [x] T023 [P] Test migration guide completeness by following steps manually
- [x] T024 [P] Verify all links in documentation are functional

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: N/A for docs - no blocking prerequisites
- **User Stories (Phase 3-6)**: All depend on tool setup completion
  - US1 and US2 (P1) can proceed in parallel
  - US3 and US4 (P2) can proceed in parallel after P1 completion
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after tool setup - No dependencies on other stories
- **User Story 2 (P1)**: Can start after tool setup - No dependencies on other stories (parallel to US1)
- **User Story 3 (P2)**: Can start after P1 stories complete - No dependencies on other P2 stories
- **User Story 4 (P2)**: Can start after P1 stories complete - No dependencies on other P2 stories (parallel to US3)

### Within Each User Story

- Documentation creation tasks are sequential within each story
- No cross-story dependencies that break independence

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- US1 and US2 can be worked on in parallel by different team members
- US3 and US4 can be worked on in parallel by different team members
- All Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Stories 1 & 2

```bash
# Launch US1 and US2 in parallel:
Task: "Create MIGRATION.md with overview section"
Task: "Add plugin discovery paths section to MIGRATION.md"

# Continue with remaining tasks for each story independently
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 3: User Story 1 (Environment Variable Migration Guide)
3. **STOP and VALIDATE**: Test User Story 1 independently (verify variable mapping documentation)
4. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo (Complete P1 features)
4. Add User Stories 3 & 4 → Test independently → Deploy/Demo (Full feature set)
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup together
2. Once setup is done:
   - Developer A: User Story 1 (Environment variables)
   - Developer B: User Story 2 (Plugin directories)
3. Once P1 stories complete:
   - Developer A: User Story 3 (JSON manifest)
   - Developer B: User Story 4 (Documentation integration)
4. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
