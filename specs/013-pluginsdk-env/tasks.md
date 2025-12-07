# Tasks: Centralized Environment Variable Handling

**Input**: Design documents from `/specs/013-pluginsdk-env/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Tests**: TDD is mandatory per Constitution. Tests are included for each user story.

**Organization**: Tasks are grouped by user story to enable independent implementation and
testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **SDK package**: `sdk/go/pluginsdk/`
- Tests co-located with implementation: `sdk/go/pluginsdk/env_test.go`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create env.go file structure and define all environment variable constants

- [x] T001 Create env.go with package declaration and imports in sdk/go/pluginsdk/env.go
- [x] T002 [P] Define EnvPort constant as "PULUMICOST_PLUGIN_PORT" in sdk/go/pluginsdk/env.go
- [x] T003 [P] Define EnvLogLevel constant as "PULUMICOST_LOG_LEVEL" in sdk/go/pluginsdk/env.go
- [x] T004 [P] Define EnvLogLevelFallback constant as "LOG_LEVEL" in sdk/go/pluginsdk/env.go
- [x] T005 [P] Define EnvLogFormat constant as "PULUMICOST_LOG_FORMAT" in sdk/go/pluginsdk/env.go
- [x] T006 [P] Define EnvLogFile constant as "PULUMICOST_LOG_FILE" in sdk/go/pluginsdk/env.go
- [x] T007 [P] Define EnvTraceID constant as "PULUMICOST_TRACE_ID" in sdk/go/pluginsdk/env.go
- [x] T008 [P] Define EnvTestMode constant as "PULUMICOST_TEST_MODE" in sdk/go/pluginsdk/env.go
- [x] T009 Create env_test.go with package declaration in sdk/go/pluginsdk/env_test.go

**Checkpoint**: All constants defined. Ready for function implementation.

---

## Phase 2: User Story 1 - Standard Port Configuration (Priority: P1) MVP

**Goal**: Plugin developers use `PULUMICOST_PLUGIN_PORT` exclusively (no fallback)

**Independent Test**: Set `PULUMICOST_PLUGIN_PORT=8080`, call `GetPort()`, verify returns 8080

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T010 [US1] Write test TestGetPort_Set in sdk/go/pluginsdk/env_test.go
- [x] T011 [US1] Write test TestGetPort_NotSet_ReturnsZero in sdk/go/pluginsdk/env_test.go
- [x] T012 [US1] Write test TestGetPort_InvalidValue_ReturnsZero in sdk/go/pluginsdk/env_test.go
- [x] T013 [US1] Write test TestGetPort_NonPositive_ReturnsZero in sdk/go/pluginsdk/env_test.go

### Implementation for User Story 1

- [x] T014 [US1] Implement GetPort() function (no fallback) in sdk/go/pluginsdk/env.go
- [x] T015 [US1] Run tests to verify GetPort() passes all test cases
- [x] T016 [US1] Modify resolvePort() in sdk/go/pluginsdk/sdk.go to use GetPort()
- [x] T017 [US1] Update error message in sdk.go to mention PULUMICOST_PLUGIN_PORT only

**Checkpoint**: User Story 1 complete. Port configuration works with canonical variable only.

---

## Phase 3: User Story 2 - Logging Configuration (Priority: P2)

**Goal**: Plugin developers can configure logging via PULUMICOST_LOG_LEVEL, LOG_FORMAT, LOG_FILE

**Independent Test**: Set `PULUMICOST_LOG_LEVEL=debug`, call `GetLogLevel()`, verify returns
"debug"

### Tests for User Story 2

- [x] T018 [P] [US2] Write test TestGetLogLevel_CanonicalVariable in sdk/go/pluginsdk/env_test.go
- [x] T019 [P] [US2] Write test TestGetLogLevel_FallbackVariable in sdk/go/pluginsdk/env_test.go
- [x] T020 [P] [US2] Write test TestGetLogLevel_CanonicalTakesPrecedence in sdk/go/pluginsdk/env_test.go
- [x] T021 [P] [US2] Write test TestGetLogLevel_NeitherSet_ReturnsEmpty in sdk/go/pluginsdk/env_test.go
- [x] T022 [P] [US2] Write test TestGetLogFormat in sdk/go/pluginsdk/env_test.go
- [x] T023 [P] [US2] Write test TestGetLogFile in sdk/go/pluginsdk/env_test.go

### Implementation for User Story 2

- [x] T024 [US2] Implement GetLogLevel() with canonical-first fallback in sdk/go/pluginsdk/env.go
- [x] T025 [P] [US2] Implement GetLogFormat() in sdk/go/pluginsdk/env.go
- [x] T026 [P] [US2] Implement GetLogFile() in sdk/go/pluginsdk/env.go
- [x] T027 [US2] Run tests to verify all logging functions pass

**Checkpoint**: User Story 2 complete. Logging configuration works independently.

---

## Phase 4: User Story 3 - Trace ID Configuration (Priority: P2)

**Goal**: Operations team can inject PULUMICOST_TRACE_ID for distributed tracing

**Independent Test**: Set `PULUMICOST_TRACE_ID=abc123`, call `GetTraceID()`, verify returns
"abc123"

### Tests for User Story 3

- [x] T028 [P] [US3] Write test TestGetTraceID_Set in sdk/go/pluginsdk/env_test.go
- [x] T029 [P] [US3] Write test TestGetTraceID_NotSet_ReturnsEmpty in sdk/go/pluginsdk/env_test.go

### Implementation for User Story 3

- [x] T030 [US3] Implement GetTraceID() in sdk/go/pluginsdk/env.go
- [x] T031 [US3] Run tests to verify GetTraceID() passes

**Checkpoint**: User Story 3 complete. Trace ID configuration works independently.

---

## Phase 5: User Story 4 - Test Mode Configuration (Priority: P2)

**Goal**: Plugin developers can enable test mode via PULUMICOST_TEST_MODE=true

**Independent Test**: Set `PULUMICOST_TEST_MODE=true`, call `IsTestMode()`, verify returns true

### Tests for User Story 4

- [x] T032 [P] [US4] Write test TestGetTestMode_True in sdk/go/pluginsdk/env_test.go
- [x] T033 [P] [US4] Write test TestGetTestMode_False in sdk/go/pluginsdk/env_test.go
- [x] T034 [P] [US4] Write test TestGetTestMode_NotSet_ReturnsFalse in sdk/go/pluginsdk/env_test.go
- [x] T035 [P] [US4] Write test TestGetTestMode_InvalidValue_ReturnsFalse in sdk/go/pluginsdk/env_test.go
- [x] T036 [P] [US4] Write test TestIsTestMode_NoWarning in sdk/go/pluginsdk/env_test.go

### Implementation for User Story 4

- [x] T037 [US4] Implement GetTestMode() with warning logging in sdk/go/pluginsdk/env.go
- [x] T038 [US4] Implement IsTestMode() without warning in sdk/go/pluginsdk/env.go
- [x] T039 [US4] Run tests to verify test mode functions pass

**Checkpoint**: User Story 4 complete. Test mode configuration works independently.

---

## Phase 6: User Story 5 - Migration Support (Priority: P3)

**Goal**: Verify Serve() integration and LOG_LEVEL backward compatibility

**Independent Test**: Set `LOG_LEVEL=debug`, call `GetLogLevel()`, verify returns "debug"

### Tests for User Story 5

- [x] T040 [P] [US5] Write integration test verifying Serve() uses GetPort() in sdk/go/pluginsdk/sdk_test.go
- [x] T041 [P] [US5] Write test verifying LOG_LEVEL fallback works in sdk/go/pluginsdk/env_test.go

### Implementation for User Story 5

- [x] T042 [US5] Verify resolvePort() modification from T016 properly integrates with Serve()
- [x] T043 [US5] Run integration tests to verify LOG_LEVEL backward compatibility

**Checkpoint**: User Story 5 complete. Serve() integration verified, LOG_LEVEL fallback works.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, validation, and cleanup

- [x] T044 [P] Add godoc comments for all constants in sdk/go/pluginsdk/env.go
- [x] T045 [P] Add godoc comments for all functions in sdk/go/pluginsdk/env.go
- [x] T046 Run `make lint` to verify code quality
- [x] T047 Run `make test` to verify all tests pass
- [x] T048 Verify quickstart.md examples work by manual testing
- [x] T049 Update sdk/go/CLAUDE.md with env.go documentation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **US1 (Phase 2)**: Depends on Setup - MVP priority
- **US2 (Phase 3)**: Depends on Setup only - can run parallel to US1
- **US3 (Phase 4)**: Depends on Setup only - can run parallel to US1/US2
- **US4 (Phase 5)**: Depends on Setup only - can run parallel to US1/US2/US3
- **US5 (Phase 6)**: Depends on US1 and US2 (needs GetPort() and GetLogLevel() for integration)
- **Polish (Phase 7)**: Depends on all user stories complete

### User Story Dependencies

```text
Setup ──┬──► US1 (P1) ───┬──► US5 (P3) ──► Polish
        │                │
        ├──► US2 (P2) ───┘
        │
        ├──► US3 (P2) ──────────────────► Polish
        │
        └──► US4 (P2) ──────────────────► Polish
```

### Within Each User Story

1. Tests MUST be written and FAIL before implementation
2. Implementation to make tests pass
3. Verify all tests pass before checkpoint

### Parallel Opportunities

- All constant definitions (T002-T008) can run in parallel
- All US2 tests (T018-T023) can run in parallel
- All US2 implementations (T025-T026 for GetLogFormat/GetLogFile) can run in parallel
- All US3 tests (T028-T029) can run in parallel
- All US4 tests (T032-T036) can run in parallel
- All US5 tests (T040-T041) can run in parallel
- All Polish tasks (T044-T049) can run in parallel (except T047 depends on code completion)

---

## Parallel Example: Setup Phase

```bash
# Launch all constant definitions in parallel:
Task: "Define EnvPort constant as 'PULUMICOST_PLUGIN_PORT'"
Task: "Define EnvLogLevel constant as 'PULUMICOST_LOG_LEVEL'"
Task: "Define EnvLogLevelFallback constant as 'LOG_LEVEL'"
Task: "Define EnvLogFormat constant as 'PULUMICOST_LOG_FORMAT'"
Task: "Define EnvLogFile constant as 'PULUMICOST_LOG_FILE'"
Task: "Define EnvTraceID constant as 'PULUMICOST_TRACE_ID'"
Task: "Define EnvTestMode constant as 'PULUMICOST_TEST_MODE'"
```

## Parallel Example: User Story 2 Tests

```bash
# Launch all logging tests in parallel:
Task: "Write test TestGetLogLevel_CanonicalVariable"
Task: "Write test TestGetLogLevel_FallbackVariable"
Task: "Write test TestGetLogLevel_CanonicalTakesPrecedence"
Task: "Write test TestGetLogLevel_NeitherSet_ReturnsEmpty"
Task: "Write test TestGetLogFormat"
Task: "Write test TestGetLogFile"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (constants)
2. Complete Phase 2: User Story 1 (GetPort + Serve integration)
3. **STOP and VALIDATE**: Test port configuration independently
4. Deploy/demo if ready - fixes the original E2E testing failure

### Incremental Delivery

1. Setup → Constants ready
2. US1 → Port configuration works → **MVP COMPLETE**
3. US2 → Logging configuration works
4. US3 → Trace ID works
5. US4 → Test mode works
6. US5 → Migration verified → Full backward compatibility
7. Polish → Documentation complete

### Solo Developer Strategy

1. Complete Setup
2. Complete US1 (MVP) - validates core functionality
3. Complete US2, US3, US4 in any order (all P2)
4. Complete US5 (integration verification)
5. Complete Polish

---

## Notes

- [P] tasks = different files or independent code sections
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- TDD required: verify tests fail before implementing
- Commit after each phase completion
- Stop at any checkpoint to validate story independently
- MVP (US1) can be deployed standalone to fix the original port mismatch issue
