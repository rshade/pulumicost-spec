<!-- markdownlint-disable MD013 -->
# Tasks: Add GetPluginInfo RPC

**Input**: Design documents from `/specs/029-plugin-info-rpc/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Conformance tests are included as they are part of the standard protocol validation pattern
for this SDK.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of
each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto definitions**: `proto/pulumicost/v1/`
- **Generated code**: `sdk/go/proto/`
- **SDK implementation**: `sdk/go/pluginsdk/`
- **Testing framework**: `sdk/go/testing/`

---

## Phase 1: Setup (Proto & Code Generation)

**Purpose**: Update protocol buffer definitions and regenerate Go SDK

- [x] T001 Add GetPluginInfoRequest message to proto/pulumicost/v1/costsource.proto
- [x] T002 Add GetPluginInfoResponse message to proto/pulumicost/v1/costsource.proto
  - Fields: name, version, spec_version, providers, metadata
- [x] T003 Add GetPluginInfo RPC to CostSourceService in proto/pulumicost/v1/costsource.proto
- [x] T004 Run `make generate` to regenerate sdk/go/proto/ from updated proto definitions
- [x] T005 Verify generated code compiles with `go build ./sdk/go/proto/...`

**Checkpoint**: Proto updated and Go SDK regenerated successfully

---

## Phase 2: Foundational (SDK Version Constants)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Create SpecVersion constant in sdk/go/pluginsdk/version.go with current spec version value
- [x] T007 [P] Add SemVer validation function ValidateSpecVersion(string) error in sdk/go/pluginsdk/version.go
- [x] T008 [P] Add unit tests for ValidateSpecVersion in sdk/go/pluginsdk/version_test.go
- [x] T009 Verify all tests pass with `make test`

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Compatibility Verification (Priority: P1) MVP

**Goal**: Enable the Core System to request and verify plugin metadata including spec version

**Independent Test**: Create a mock plugin implementing GetPluginInfo RPC and verify the Core
receives correct data (name, version, spec_version, providers)

### Tests for User Story 1 (Write First, Must FAIL)

> **Constitution III**: Tests MUST be written and FAIL before implementation begins.

- [x] T010 [US1] Add conformance test for GetPluginInfo RPC in sdk/go/testing/conformance_test.go
- [x] T011 [US1] Add integration test for GetPluginInfo metadata retrieval in sdk/go/testing/integration_test.go
- [x] T012 [US1] Verify tests FAIL with `make test` (expected: GetPluginInfo not implemented)

### Implementation for User Story 1

- [x] T013 [US1] Update MockPlugin to implement GetPluginInfo in sdk/go/testing/mock_plugin.go
- [x] T014 [US1] Add PluginInfo struct to hold plugin metadata in sdk/go/pluginsdk/plugin_info.go
- [x] T015 [US1] Add PluginInfoOption functional options for configuring PluginInfo in sdk/go/pluginsdk/plugin_info.go
- [x] T016 [US1] Add providers field to PluginInfo struct in sdk/go/pluginsdk/plugin_info.go
- [x] T017 [US1] Add GetPluginInfo method signature to Plugin interface in sdk/go/pluginsdk/plugin.go
- [x] T018 [US1] Implement default GetPluginInfo handler in sdk/go/pluginsdk/base.go that returns PluginInfo
- [x] T019 [US1] Update Serve function to accept PluginInfo configuration in sdk/go/pluginsdk/serve.go
- [x] T020 [US1] Add server-side validation for required fields (name, version, spec_version) in PluginInfo configuration
- [x] T021 [US1] Verify tests PASS with `make validate`

**Checkpoint**: User Story 1 complete - Core can retrieve and validate plugin metadata

---

## Phase 4: User Story 2 - Diagnostic Visibility (Priority: P2)

**Goal**: Enable developers/operators to see spec version when listing plugins

**Independent Test**: Call GetPluginInfo on a plugin and verify response includes all diagnostic fields

### Implementation for User Story 2

- [x] T022 [US2] Add optional metadata map field to PluginInfo in sdk/go/pluginsdk/plugin_info.go
- [x] T023 [US2] Add WithMetadata functional option for additional key-value metadata in sdk/go/pluginsdk/plugin_info.go
- [x] T024 [US2] Add unit tests for PluginInfo with metadata in sdk/go/pluginsdk/plugin_info_test.go
- [x] T025 [US2] Add conformance test verifying metadata field in response in sdk/go/testing/conformance_test.go
- [x] T026 [US2] Run `make validate` to ensure all tests and linting pass

**Checkpoint**: User Story 2 complete - All diagnostic info available via GetPluginInfo

---

## Phase 5: User Story 3 - Graceful Degradation (Priority: P3)

**Goal**: Handle legacy plugins that don't implement GetPluginInfo without crashing

**Independent Test**: Call GetPluginInfo on a plugin that doesn't implement it and verify
Unimplemented error is handled

### Implementation for User Story 3

- [x] T027 [US3] Document Unimplemented error handling pattern in sdk/go/pluginsdk/README.md or inline comments
- [x] T028 [US3] Add example code for handling Unimplemented error in client code comments in sdk/go/pluginsdk/base.go
- [x] T029 [US3] Add conformance test for Unimplemented error scenario in sdk/go/testing/conformance_test.go
- [x] T030 [US3] Update testing framework to support legacy plugin simulation in sdk/go/testing/mock_plugin.go
- [x] T031 [US3] Add integration test for graceful degradation with legacy plugin in sdk/go/testing/integration_test.go
- [x] T032 [US3] Run `make validate` to ensure all tests and linting pass

**Checkpoint**: User Story 3 complete - Legacy plugins handled gracefully

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, documentation, and cleanup

- [x] T033 [P] Update sdk/go/pluginsdk/README.md with GetPluginInfo usage documentation (include 5s timeout recommendation per FR-005)
- [x] T034 [P] Add benchmark test for GetPluginInfo in sdk/go/testing/benchmark_test.go
- [x] T035 [P] Update CLAUDE.md with new patterns and files created
- [x] T036 Run full validation with `make validate`
- [x] T037 Run quickstart.md validation steps to verify end-to-end flow

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 -> P2 -> P3)
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Extends PluginInfo from US1 but
  is independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent error handling, no
  US1/US2 dependencies

### Within Each User Story

- **Tests FIRST** (Constitution III - NON-NEGOTIABLE)
- Tests must FAIL before implementation begins
- Models/types before handlers
- Handlers before integration
- Verify tests PASS after implementation
- Story complete before moving to next priority

### Parallel Opportunities

- T007 and T008 can run in parallel (different files)
- T033, T034, T035 can run in parallel (different files)
- After Phase 2, all three user stories could theoretically run in parallel with separate developers

---

## Parallel Example: Phase 2 Foundation

```bash
# Launch parallel tasks in Phase 2:
Task: "Add SemVer validation function in sdk/go/pluginsdk/version.go"
Task: "Add unit tests for ValidateSpecVersion in sdk/go/pluginsdk/version_test.go"
```

---

## Parallel Example: Polish Phase

```bash
# Launch all Polish tasks together:
Task: "Update sdk/go/pluginsdk/README.md with GetPluginInfo usage"
Task: "Add benchmark test in sdk/go/testing/benchmark_test.go"
Task: "Update CLAUDE.md with new patterns"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (proto changes)
2. Complete Phase 2: Foundational (version constants)
3. Complete Phase 3: User Story 1 (core GetPluginInfo)
   - Tests first (T010-T012) - must FAIL
   - Implementation (T013-T020)
   - Verify tests PASS (T021)
4. **STOP and VALIDATE**: Test GetPluginInfo independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational -> Proto and SDK foundation ready
2. Add User Story 1 -> Test independently -> Deploy (MVP!)
3. Add User Story 2 -> Test independently -> Deploy (enhanced diagnostics)
4. Add User Story 3 -> Test independently -> Deploy (full backward compatibility)
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (core functionality)
   - Developer B: User Story 2 (diagnostics) - can start in parallel
   - Developer C: User Story 3 (graceful degradation) - can start in parallel
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **Tests must be written and FAIL before implementation** (Constitution III)
- Run `make validate` at each checkpoint to catch issues early
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Proto changes (Phase 1) must complete before any SDK work
