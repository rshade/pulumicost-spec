---
description: "Task list for SDK Polish v0.4.15 feature implementation"
---

# Tasks: SDK Polish v0.4.15

**Input**: Design documents from `/specs/034-sdk-polish/`
**Feature Branch**: `034-sdk-polish`
**Tech Stack**: Go 1.25.5, gRPC, protobuf, buf v1.32.1

**Summary**: SDK Polish v0.4.15 is primarily a verification and testing task. All three features are already
implemented in the codebase. The work focuses on:

1. Verifying existing implementations match requirements
2. Adding missing integration/conformance tests
3. Updating documentation

## Phase 1: Setup (Verification Preparation)

**Purpose**: Verify existing implementation state and prepare test infrastructure

- [x] T001 Review existing ClientConfig.Timeout implementation in sdk/go/pluginsdk/client.go:82-132
- [x] T002 Review existing wrapRPCError function in sdk/go/pluginsdk/client.go:163-176
- [x] T003 Review existing GetPluginInfo error handling in sdk/go/pluginsdk/sdk.go:323-399
- [x] T004 Review existing Performance_GetPluginInfoLatency test in sdk/go/testing/performance.go:211-216
- [x] T005 Review existing unit tests in sdk/go/pluginsdk/client_test.go
- [x] T006 Review existing SDK tests in sdk/go/pluginsdk/sdk_test.go:1014-1105

---

## Phase 2: Foundational (Missing Integration Tests)

**Purpose**: Add missing integration tests that were identified in research as gaps

**⚠️ CRITICAL**: These tests verify behavior that currently has no coverage

### User Story 1 Integration Tests (US1 - P1)

**Goal**: Verify timeout configuration works correctly with real RPC calls

**Independent Test**: Create a slow mock server, configure client with timeout, verify timeout triggers

#### Timeout Integration Tests

- [x] T007 [P] [US1] Create timeout test helper with slow mock server in sdk/go/pluginsdk/timeout_test.go
- [x] T008 [US1] Add TestClientTimeout_ExceedsConfiguredTimeout in sdk/go/pluginsdk/timeout_test.go
- [x] T009 [US1] Add TestClientTimeout_ContextDeadlinePrecedence in sdk/go/pluginsdk/timeout_test.go
- [x] T010 [US1] Add TestClientTimeout_DefaultValue in sdk/go/pluginsdk/timeout_test.go
- [x] T011 [P] [US1] Add TestClientTimeout_CustomHTTPClientPrecedence in sdk/go/pluginsdk/timeout_test.go

### User Story 2 Integration Tests (US2 - P2)

**Goal**: Verify error messages match user-friendly format requirements

**Independent Test**: Configure mock plugin to return nil/incomplete/invalid responses,
verify client receives correct error messages

#### Error Message Conformance Tests

- [x] T012 [P] [US2] Create error message test helper in sdk/go/pluginsdk/error_test.go
- [x] T013 [US2] Add TestGetPluginInfoError_NilResponse in sdk/go/pluginsdk/error_test.go
- [x] T014 [US2] Add TestGetPluginInfoError_IncompleteMetadata in sdk/go/pluginsdk/error_test.go
- [x] T015 [US2] Add TestGetPluginInfoError_InvalidSpecVersion in sdk/go/pluginsdk/error_test.go

### User Story 3 Integration Tests (US3 - P3)

**Goal**: Verify performance test handles legacy plugins gracefully (FR-011)

**Independent Test**: Run performance test against legacy plugin (no GetPluginInfo),
verify Unimplemented handled gracefully

#### Legacy Plugin Performance Test

- [x] T016 [P] [US3] Create mock legacy plugin for performance testing in sdk/go/testing/mock_legacy_plugin.go
- [x] T017 [US3] Add TestGetPluginInfoPerformance_LegacyPlugin in sdk/go/testing/performance_test.go

---

## Phase 3: User Story 1 - Client Timeout Configuration (Priority: P1) 🎯 MVP

**Goal**: Verify client timeout configuration works as specified in spec.md

**Independent Test**: All timeout integration tests pass (T007-T011)

### Verification Tests (Required for Verification Task)

> **NOTE: Tests verify existing implementation behavior - ensure they PASS**

- [x] T018 [P] [US1] Run existing TestClientConfig_WithTimeout in sdk/go/pluginsdk/client_test.go:234-241
- [x] T019 [US1] Run new TestClientTimeout_ExceedsConfiguredTimeout and verify timeout triggers
- [x] T020 [US1] Run new TestClientTimeout_ContextDeadlinePrecedence and verify context takes precedence
- [x] T021 [US1] Run new TestClientTimeout_DefaultValue and verify 30-second default
- [x] T022 [US1] Run new TestClientTimeout_CustomHTTPClientPrecedence and verify HTTPClient timeout precedence

### Implementation Verification

- [x] T023 [US1] Verify ClientConfig.Timeout field in sdk/go/pluginsdk/client.go:99
- [x] T024 [US1] Verify WithTimeout() method in sdk/go/pluginsdk/client.go:116-132
- [x] T025 [US1] Verify wrapRPCError context deadline handling in sdk/go/pluginsdk/client.go:163-176
- [x] T026 [US1] Verify NewClient timeout application in sdk/go/pluginsdk/client.go:204-213

**Checkpoint**: User Story 1 verified - timeout configuration works correctly

---

## Phase 4: User Story 2 - User-Friendly GetPluginInfo Error Messages (Priority: P2)

**Goal**: Verify error messages match user-friendly format specified in FR-006, FR-007, FR-008

**Independent Test**: All error message tests pass (T013-T015)

### Verification Tests (Required for Verification Task)

- [x] T027 [P] [US2] Run existing TestGetPluginInfo tests in sdk/go/pluginsdk/sdk_test.go:1014-1105
- [x] T028 [US2] Run new TestGetPluginInfoError_NilResponse and verify "unable to retrieve plugin metadata" message
- [x] T029 [US2] Run new TestGetPluginInfoError_IncompleteMetadata and verify "plugin metadata is incomplete" message
- [x] T030 [US2] Run new TestGetPluginInfoError_InvalidSpecVersion and verify
  "plugin reported an invalid specification version" message

### Implementation Verification

- [x] T031 [US2] Verify nil response error message in sdk/go/pluginsdk/sdk.go:342
- [x] T032 [US2] Verify incomplete metadata error message in sdk/go/pluginsdk/sdk.go:350
- [x] T033 [US2] Verify invalid spec version error message in sdk/go/pluginsdk/sdk.go:361
- [x] T034 [US2] Verify server-side logging captures detailed error info

**Checkpoint**: User Story 2 verified - error messages match requirements

---

## Phase 5: User Story 3 - GetPluginInfo Performance Conformance (Priority: P3)

**Goal**: Verify performance test exists and handles legacy plugins gracefully (FR-009, FR-010, FR-011)

**Independent Test**: All performance tests pass, including legacy plugin test

### Verification Tests (Required for Verification Task)

- [x] T035 [P] [US3] Run existing Performance_GetPluginInfoLatency test
- [x] T036 [US3] Verify test runs 10 iterations per FR-010
- [x] T037 [US3] Verify test uses 100ms threshold for Standard conformance
- [x] T038 [US3] Run new TestGetPluginInfoPerformance_LegacyPlugin and verify graceful handling

### Implementation Verification

- [x] T039 [US3] Verify GetPluginInfoStandardLatencyMs constant in sdk/go/testing/harness.go:64
- [x] T040 [US3] Verify GetPluginInfoAdvancedLatencyMs constant in sdk/go/testing/harness.go:66
- [x] T041 [US3] Verify createGetPluginInfoLatencyTest function in sdk/go/testing/performance.go:297-307
- [x] T042 [US3] Verify measureLatency function runs correct iteration count

**Checkpoint**: User Story 3 verified - performance conformance works correctly

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation updates and final validation

### Documentation Updates (Required per Constitution Check)

- [x] T043 [P] Update SDK README.md with timeout configuration examples
  (maps to FR-001, FR-002, FR-003, FR-004, Constitution V)
- [x] T044 [P] Update SDK README.md with GetPluginInfo error message documentation
  (maps to FR-006, FR-007, FR-008, Constitution V)
- [x] T045 [P] Update SDK README.md with performance conformance test documentation
  (maps to FR-009, FR-010, FR-011, Constitution V)
- [x] T046 Add conformance testing documentation to docs/conformance.md (maps to FR-009, FR-010, FR-011, Constitution V)

### Final Validation

- [x] T047 Run all timeout-related tests and verify SC-001
- [x] T048 Run all error message tests and verify SC-002
- [x] T049 Run all performance tests and verify SC-003, SC-005
- [x] T050 Run existing test suite and verify no regressions (SC-004)
- [x] T051 Run `make test` to validate all tests pass
- [x] T052 Run `make lint` to validate code quality

---

## Dependencies & Execution Order

### Phase Dependencies

| Phase                 | Description                    | Dependencies                         |
| --------------------- | ------------------------------ | ------------------------------------ |
| Phase 1: Setup        | Review existing implementation | None - can start immediately         |
| Phase 2: Foundational | Add missing integration tests  | Depends on Phase 1 completion        |
| Phase 3: US1 (P1)     | Timeout verification           | Depends on Phase 2 completion        |
| Phase 4: US2 (P2)     | Error message verification     | Depends on Phase 2 completion        |
| Phase 5: US3 (P3)     | Performance verification       | Depends on Phase 2 completion        |
| Phase 6: Polish       | Documentation and validation   | Depends on Phases 3, 4, 5 completion |

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 2 - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Phase 2 - Independent from US1
- **User Story 3 (P3)**: Can start after Phase 2 - Independent from US1/US2

### Recommended Execution Order

1. Complete Phase 1 (Setup) to understand existing implementation
2. Complete Phase 2 (Foundational) to add missing tests
3. Complete US1 (Phase 3) - **This is the MVP**
4. Deploy/demo if US1 is successful
5. Complete US2 (Phase 4) when ready
6. Complete US3 (Phase 5) when ready
7. Complete Phase 6 (Polish) for documentation and final validation

---

## Parallel Opportunities

### Within Phase 1 (Setup)

All tasks (T001-T006) can run in parallel - they are independent file reviews.

### Within Phase 2 (Foundational)

- T007, T012, T016 can run in parallel (different test files)
- All error message tests (T013-T015) can run in parallel after T012
- All timeout tests (T008-T011) can run in parallel after T007
- All user stories can be worked on in parallel by different developers

### Within User Stories

- All verification tests within a user story can run in parallel
- All implementation verification tasks within a story can run in parallel

### Parallel Example Commands

```bash
# Run all Phase 1 setup reviews in parallel:
Task: "Review ClientConfig.Timeout implementation"
Task: "Review wrapRPCError function"
Task: "Review GetPluginInfo error handling"

# Run all US1 timeout tests in parallel:
Task: "TestClientTimeout_ExceedsConfiguredTimeout"
Task: "TestClientTimeout_ContextDeadlinePrecedence"
Task: "TestClientTimeout_DefaultValue"
Task: "TestClientTimeout_CustomHTTPClientPrecedence"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

For a minimum viable delivery:

1. Complete Phase 1: Setup (review existing implementation)
2. Complete Phase 2: Foundational (add US1 timeout tests - T007-T011)
3. Complete Phase 3: User Story 1 (verify timeout configuration works)
4. **STOP and VALIDATE**: Run timeout tests to verify behavior
5. Deploy/demo if successful

### Incremental Delivery

1. Setup + Foundational → Test infrastructure ready
2. Add User Story 1 → Test independently → Deploy (MVP!)
3. Add User Story 2 → Test independently → Deploy
4. Add User Story 3 → Test independently → Deploy
5. Polish phase → Documentation updates and final validation

### Parallel Team Strategy

With multiple developers:

1. **Developer A**: Complete Phase 1 + Phase 2 (test infrastructure)
2. **Developer B**: Complete User Story 1 verification (Phase 3)
3. **Developer C**: Complete User Story 2 verification (Phase 4)
4. **Developer A/D**: Complete User Story 3 verification (Phase 5)
5. **Anyone**: Complete Phase 6 (Polish)

---

## Task Summary

| Category              | Count  | Description                         |
| --------------------- | ------ | ----------------------------------- |
| Phase 1: Setup        | 6      | Review existing implementation      |
| Phase 2: Foundational | 11     | Add missing integration tests       |
| Phase 3: US1 (P1)     | 10     | Timeout verification                |
| Phase 4: US2 (P2)     | 8      | Error message verification          |
| Phase 5: US3 (P3)     | 8      | Performance verification            |
| Phase 6: Polish       | 10     | Documentation and validation        |
| **Total**             | **53** | Tasks for complete feature delivery |

### Tasks per User Story

- **User Story 1 (P1)**: 10 tasks (includes Phase 2 timeout tests)
- **User Story 2 (P2)**: 8 tasks (includes Phase 2 error tests)
- **User Story 3 (P3)**: 8 tasks (includes Phase 2 performance tests)

### Independent Test Criteria

**US1 Independent Test**: All timeout integration tests pass (T018-T022)
**US2 Independent Test**: All error message tests pass (T027-T030)
**US3 Independent Test**: All performance tests pass (T035-T038)

---

## Notes

- **[P]** tasks = parallelizable (different files, no dependencies)
- **[US1], [US2], [US3]** labels map tasks to specific user stories for traceability
- All features already implemented - this is a verification and testing task
- Focus on adding integration tests that are currently missing (identified in research.md)
- Commit after each task or logical group
- Run `make test` after completing each phase to validate progress
- Stop at any checkpoint to validate story independently
