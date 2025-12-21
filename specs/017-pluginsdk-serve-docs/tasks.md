# Tasks: Document pluginsdk.Serve() Behavior

**Input**: Design documents from `/specs/017-pluginsdk-serve-docs/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, quickstart.md

**Tests**: No tests explicitly requested in the feature specification. This is a documentation
feature, so validation is via markdown linting and example code compilation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Documentation target**: `sdk/go/pluginsdk/README.md`
- **Source reference**: `sdk/go/pluginsdk/sdk.go`, `sdk/go/pluginsdk/env.go`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Verify existing structure and prepare for documentation updates

- [ ] T001 Verify existing README.md exists at sdk/go/pluginsdk/README.md (or create if missing)
- [ ] T002 [P] Review existing Go doc comments in sdk/go/pluginsdk/sdk.go for Serve() function
- [ ] T003 [P] Review existing Go doc comments in sdk/go/pluginsdk/env.go for environment functions

---

## Phase 2: Foundational (Documentation Structure)

**Purpose**: Create base documentation structure that all user stories will build on

**⚠️ CRITICAL**: User story sections depend on this structure being in place

- [ ] T004 Create documentation header and overview section in sdk/go/pluginsdk/README.md
- [ ] T005 Create table of contents structure in sdk/go/pluginsdk/README.md
- [ ] T006 [P] Create section placeholder for Serve() documentation in sdk/go/pluginsdk/README.md
- [ ] T007 [P] Create section placeholder for environment variables in sdk/go/pluginsdk/README.md

**Checkpoint**: Foundation ready - user story documentation can now be added

---

## Phase 3: User Story 1 - Plugin Developer Learns Serve() Usage (Priority: P1)

**Goal**: Document Serve() function signature, ServeConfig struct, and basic usage so developers
can write a minimal plugin main() function.

**Independent Test**: A developer reading this section can implement a working plugin that starts
and accepts gRPC connections within 15 minutes.

### Implementation for User Story 1

- [ ] T008 [US1] Document Serve() function signature and return behavior in sdk/go/pluginsdk/README.md
- [ ] T009 [US1] Document ServeConfig struct with all 5 fields (Plugin, Port, Registry, Logger,
  UnaryInterceptors) in sdk/go/pluginsdk/README.md
- [ ] T010 [US1] Add minimal plugin example code (from quickstart.md) to sdk/go/pluginsdk/README.md
- [ ] T011 [US1] Document Plugin interface requirements in sdk/go/pluginsdk/README.md
- [ ] T012 [US1] Document optional interfaces (SupportsProvider, RecommendationsProvider) in
  sdk/go/pluginsdk/README.md

**Checkpoint**: User Story 1 complete - developers can now understand and use Serve() for basic
plugin implementation

---

## Phase 4: User Story 2 - Plugin Developer Understands Port Resolution (Priority: P1)

**Goal**: Document port resolution priority and configuration methods so developers understand
how to configure ports for their plugins.

**Independent Test**: A developer can correctly configure ports using each method (flag, env var,
ephemeral) based on the documentation.

### Implementation for User Story 2

- [ ] T013 [US2] Document port resolution priority (--port > PULUMICOST_PLUGIN_PORT > ephemeral)
  in sdk/go/pluginsdk/README.md
- [ ] T014 [US2] Document ParsePortFlag() function and flag.Parse() requirement in
  sdk/go/pluginsdk/README.md
- [ ] T015 [US2] Document `PORT=<port>` stdout announcement format in sdk/go/pluginsdk/README.md
- [ ] T016 [P] [US2] Add port configuration examples (flag, env var, ephemeral) to
  sdk/go/pluginsdk/README.md

**Checkpoint**: User Story 2 complete - developers understand all port configuration options

---

## Phase 5: User Story 3 - Plugin Developer Configures Logging and Tracing (Priority: P2)

**Goal**: Document all environment variables for logging configuration and distributed tracing.

**Independent Test**: A developer can configure logging levels, formats, and trace IDs by setting
the documented environment variables.

### Implementation for User Story 3

- [ ] T017 [US3] Create environment variables reference table in sdk/go/pluginsdk/README.md
- [ ] T018 [P] [US3] Document PULUMICOST_LOG_LEVEL with supported values in
  sdk/go/pluginsdk/README.md
- [ ] T019 [P] [US3] Document PULUMICOST_LOG_FORMAT (json/text) in sdk/go/pluginsdk/README.md
- [ ] T020 [P] [US3] Document PULUMICOST_LOG_FILE in sdk/go/pluginsdk/README.md
- [ ] T021 [P] [US3] Document PULUMICOST_TRACE_ID for distributed tracing in
  sdk/go/pluginsdk/README.md
- [ ] T022 [P] [US3] Document PULUMICOST_TEST_MODE behavior in sdk/go/pluginsdk/README.md
- [ ] T023 [US3] Document getter functions (GetPort, GetLogLevel, etc.) in
  sdk/go/pluginsdk/README.md

**Checkpoint**: User Story 3 complete - developers can configure all logging and tracing options

---

## Phase 6: User Story 4 - DevOps Engineer Deploys Multiple Plugins (Priority: P2)

**Goal**: Document why generic PORT is not supported and how to orchestrate multiple plugins.

**Independent Test**: A DevOps engineer understands why PORT fallback doesn't exist and can deploy
multiple plugins with unique ports.

### Implementation for User Story 4

- [ ] T024 [US4] Document multi-plugin orchestration scenario in sdk/go/pluginsdk/README.md
- [ ] T025 [US4] Explain why generic PORT env var is not supported (multi-plugin conflicts) in
  sdk/go/pluginsdk/README.md
- [ ] T026 [P] [US4] Add deployment example showing two plugins with distinct --port flags in
  sdk/go/pluginsdk/README.md

**Checkpoint**: User Story 4 complete - DevOps engineers understand multi-plugin deployment

---

## Phase 7: User Story 5 - Plugin Developer Implements Graceful Shutdown (Priority: P3)

**Goal**: Document graceful shutdown behavior when context is canceled.

**Independent Test**: A developer understands how context cancellation triggers GracefulStop()
and that in-flight requests complete.

### Implementation for User Story 5

- [ ] T027 [US5] Document graceful shutdown sequence in sdk/go/pluginsdk/README.md
- [ ] T028 [US5] Explain GracefulStop() vs Stop() behavior in sdk/go/pluginsdk/README.md
- [ ] T029 [P] [US5] Add shutdown example code showing context cancellation handling in
  sdk/go/pluginsdk/README.md

**Checkpoint**: User Story 5 complete - developers understand shutdown behavior

---

## Phase 8: Edge Cases and Error Handling

**Goal**: Document error conditions and edge cases from spec.md.

### Implementation for Edge Cases

- [ ] T030 [P] Document error: "failed to listen: address already in use" (port conflict) in
  sdk/go/pluginsdk/README.md
- [ ] T031 [P] Document error: "failed to listen: permission denied" in sdk/go/pluginsdk/README.md
- [ ] T032 [P] Document behavior when ParsePortFlag() called before flag.Parse() in
  sdk/go/pluginsdk/README.md
- [ ] T033 [P] Document context.Canceled error handling in sdk/go/pluginsdk/README.md

**Checkpoint**: All edge cases documented

---

## Phase 9: Polish & Validation

**Purpose**: Final validation and cross-cutting improvements

- [ ] T034 Run markdown linting on sdk/go/pluginsdk/README.md
- [ ] T035 Verify all code examples compile (go build check)
- [ ] T036 [P] Update sdk/go/CLAUDE.md with reference to new README documentation
- [ ] T037 Cross-reference with quickstart.md to ensure consistency
- [ ] T038 Final review of documentation against FR-001 through FR-011 requirements

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
- **Edge Cases (Phase 8)**: Can run after Foundational, parallel with User Stories
- **Polish (Phase 9)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational - Independent of US1
- **User Story 3 (P2)**: Can start after Foundational - Independent
- **User Story 4 (P2)**: Can start after Foundational - May reference US2 port content
- **User Story 5 (P3)**: Can start after Foundational - Independent

### Parallel Opportunities

- T002, T003 can run in parallel (different files)
- T006, T007 can run in parallel (different sections)
- T018, T019, T020, T021, T022 can run in parallel (different env vars)
- T030, T031, T032, T033 can run in parallel (different error conditions)
- User Story 3 tasks (T018-T022) can run in parallel with User Story 4 tasks (T024-T026)

---

## Parallel Example: Environment Variables (User Story 3)

```bash
# Launch all environment variable documentation tasks together:
Task: "Document PULUMICOST_LOG_LEVEL with supported values in sdk/go/pluginsdk/README.md"
Task: "Document PULUMICOST_LOG_FORMAT (json/text) in sdk/go/pluginsdk/README.md"
Task: "Document PULUMICOST_LOG_FILE in sdk/go/pluginsdk/README.md"
Task: "Document PULUMICOST_TRACE_ID for distributed tracing in sdk/go/pluginsdk/README.md"
Task: "Document PULUMICOST_TEST_MODE behavior in sdk/go/pluginsdk/README.md"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1 (Serve() basics)
4. Complete Phase 4: User Story 2 (Port resolution)
5. **STOP and VALIDATE**: Developers can implement a basic plugin
6. Review/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Documentation structure ready
2. Add User Story 1 → Test: developer can use Serve() (MVP!)
3. Add User Story 2 → Test: developer understands port configuration
4. Add User Story 3 → Test: developer can configure logging
5. Add User Story 4 → Test: DevOps can deploy multiple plugins
6. Add User Story 5 → Test: developer understands graceful shutdown
7. Add Edge Cases → Complete error documentation
8. Polish → Final validation

---

## Notes

- [P] tasks = different files or sections, no dependencies
- [Story] label maps task to specific user story for traceability
- All documentation goes to single file: sdk/go/pluginsdk/README.md
- Example code must be copy-paste-ready and compile successfully
- Run `npx markdownlint-cli2 sdk/go/pluginsdk/README.md` for validation
- Commit after each user story completion for incremental delivery
