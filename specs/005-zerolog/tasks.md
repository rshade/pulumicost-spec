# Tasks: Zerolog SDK Logging Utilities

**Input**: Design documents from `/specs/005-zerolog/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED per spec (FR-010: 90%+ coverage)

**Organization**: Tasks grouped by user story for independent implementation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **SDK Package**: `sdk/go/pluginsdk/`
- **Tests**: `sdk/go/pluginsdk/` (same package, `_test.go` files)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and package structure

- [x] T001 Create pluginsdk package directory at sdk/go/pluginsdk/
- [x] T002 Add zerolog v1.34.0+ dependency to go.mod
- [x] T003 Run go mod tidy to update go.sum

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and constants that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Define context key type and trace ID context key in
      sdk/go/pluginsdk/logging.go
- [x] T005 [P] Define TraceIDMetadataKey constant in sdk/go/pluginsdk/logging.go
- [x] T006 [P] Define all 11 field name constants (FieldTraceID, FieldComponent,
      etc.) in sdk/go/pluginsdk/logging.go

**Checkpoint**: Foundation ready - user story implementation can begin

---

## Phase 3: User Story 1 - Plugin Developer Creates Standardized Logger (P1) 🎯 MVP

**Goal**: Plugin developers can create a configured zerolog logger with plugin
metadata in a single function call

**Independent Test**: Create logger, log message, verify JSON output contains
plugin_name, plugin_version, timestamp fields

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T007 [P] [US1] Unit test for NewPluginLogger with default stderr in
      sdk/go/pluginsdk/logging_test.go
- [x] T008 [P] [US1] Unit test for NewPluginLogger with custom io.Writer in
      sdk/go/pluginsdk/logging_test.go
- [x] T009 [P] [US1] Unit test for log level filtering (Debug vs Info) in
      sdk/go/pluginsdk/logging_test.go
- [x] T010 [P] [US1] Unit test for empty plugin name/version handling in
      sdk/go/pluginsdk/logging_test.go

### Implementation for User Story 1

- [x] T011 [US1] Implement NewPluginLogger function with plugin name, version,
      level, and io.Writer parameters in sdk/go/pluginsdk/logging.go
- [x] T012 [US1] Configure logger to output structured JSON (FR-009) with
      plugin_name and plugin_version base fields in sdk/go/pluginsdk/logging.go
- [x] T013 [US1] Default to os.Stderr when io.Writer is nil in
      sdk/go/pluginsdk/logging.go
- [x] T014 [US1] Implement file output support with --logfile flag in
      sdk/go/pluginsdk/logging.go
- [x] T015 [US1] Implement directory output support with --logdir flag and
      default filenames in sdk/go/pluginsdk/logging.go
- [x] T016 [P] [US1] Unit test for file output configuration in
      sdk/go/pluginsdk/logging_test.go

**Checkpoint**: User Story 1 complete - plugin logger creation works independently

---

## Phase 4: User Story 2 - Core System Traces Requests (P1)

**Goal**: System operators can trace requests through gRPC plugin calls using
trace_id propagation

**Independent Test**: Configure gRPC server with interceptor, send request with
trace_id metadata, verify handler can retrieve trace_id from context

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T017 [P] [US2] Unit test for ContextWithTraceID/TraceIDFromContext in
      sdk/go/pluginsdk/logging_test.go
- [x] T018 [P] [US2] Unit test for TraceIDFromContext with empty context in
      sdk/go/pluginsdk/logging_test.go
- [x] T019 [P] [US2] Integration test for TracingUnaryServerInterceptor with
      bufconn in sdk/go/pluginsdk/logging_test.go
- [x] T020 [P] [US2] Integration test for interceptor with missing metadata in
      sdk/go/pluginsdk/logging_test.go
- [x] T021 [P] [US2] Integration test for concurrent requests with different
      trace_ids in sdk/go/pluginsdk/logging_test.go

### Implementation for User Story 2

- [x] T022 [US2] Implement ContextWithTraceID function in
      sdk/go/pluginsdk/logging.go
- [x] T023 [US2] Implement TraceIDFromContext function in
      sdk/go/pluginsdk/logging.go
- [x] T024 [US2] Implement TracingUnaryServerInterceptor that extracts trace_id
      from gRPC metadata in sdk/go/pluginsdk/logging.go
- [x] T025 [US2] Handle multiple trace_id values in metadata (use first) in
      sdk/go/pluginsdk/logging.go

**Checkpoint**: User Story 2 complete - trace ID propagation works independently

---

## Phase 5: User Story 3 - Plugin Developer Logs Operation Timing (P2)

**Goal**: Plugin developers can measure and log operation duration using a
defer-friendly helper function

**Depends on**: User Story 1 (requires logger instance for timing output)

**Independent Test**: Call LogOperation, perform task, call returned function,
verify duration_ms appears in log output

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T026 [P] [US3] Unit test for LogOperation timing accuracy in
      sdk/go/pluginsdk/logging_test.go
- [x] T027 [P] [US3] Unit test for LogOperation log output format in
      sdk/go/pluginsdk/logging_test.go

### Implementation for User Story 3

- [x] T028 [US3] Implement LogOperation function returning closure in
      sdk/go/pluginsdk/logging.go
- [x] T029 [US3] Log operation name and duration_ms using standard field
      constants in sdk/go/pluginsdk/logging.go

**Checkpoint**: User Story 3 complete - operation timing works independently

---

## Phase 6: User Story 4 - Plugin Developer Uses Standard Field Names (P2)

**Goal**: Plugin developers use consistent field names for ecosystem-wide log
analysis compatibility

**Independent Test**: Use field constants in log statements, verify resulting
JSON uses exact field names

### Tests for User Story 4

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T030 [P] [US4] Unit test verifying all 11 field constants have correct
      string values in sdk/go/pluginsdk/logging_test.go

### Implementation for User Story 4

- [x] T031 [US4] Code review: verify all 11 field constants from T006 are
      exported, have godoc comments, and match data-model.md values in
      sdk/go/pluginsdk/logging.go

**Checkpoint**: User Story 4 complete - field constants verified

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, benchmarks, and final validation

- [x] T032 [P] Add package documentation, function godoc comments, and logging
      best practices (SC-005, SC-006) in sdk/go/pluginsdk/logging.go
- [x] T033 [P] Create benchmark tests for logger construction in
      sdk/go/pluginsdk/logging_test.go
- [x] T034 [P] Create benchmark tests for log call overhead in
      sdk/go/pluginsdk/logging_test.go
- [x] T035 [P] Create benchmark tests for interceptor performance in
      sdk/go/pluginsdk/logging_test.go
- [x] T036 Verify 90%+ test coverage with go test -cover
- [x] T037 Run all tests and benchmarks, fix any issues
- [x] T038 Validate quickstart.md examples compile and run correctly

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase
  - US1 and US2 are both P1, can run in parallel
  - US3 and US4 are both P2, can run in parallel after P1
- **Polish (Phase 7)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies on other stories
- **User Story 2 (P1)**: No dependencies on other stories (context helpers
  independent of logger)
- **User Story 3 (P2)**: Depends on US1 (needs logger to log timing)
- **User Story 4 (P2)**: No implementation - validation only (constants in
  Phase 2)

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Implementation follows test requirements
- Story complete before moving to next priority

### Parallel Opportunities

**Phase 2 (Foundational)**:

- T005 and T006 can run in parallel

**Phase 3 (US1 Tests)**:

- T007, T008, T009, T010, T016 can all run in parallel

**Phase 4 (US2 Tests)**:

- T017, T018, T019, T020, T021 can all run in parallel

**Phase 5-6 (US3-4 Tests)**:

- All tests within each phase can run in parallel

**Phase 7 (Polish)**:

- T032, T033, T034, T035 can all run in parallel

---

## Parallel Example: User Story 2 Tests

```bash
# Launch all US2 tests together:
Task: "Unit test for ContextWithTraceID/TraceIDFromContext"
Task: "Unit test for TraceIDFromContext with empty context"
Task: "Integration test for TracingUnaryServerInterceptor with bufconn"
Task: "Integration test for interceptor with missing metadata"
Task: "Integration test for concurrent requests with different trace_ids"
```

---

## Implementation Strategy

### MVP First (User Stories 1 & 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL)
3. Complete Phase 3: User Story 1 (logger creation)
4. Complete Phase 4: User Story 2 (trace propagation)
5. **STOP and VALIDATE**: Both P1 stories work independently
6. Can deploy/demo basic logging + tracing

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add US1 → Logger creation works → Can demo
3. Add US2 → Trace propagation works → Can demo
4. Add US3 → Operation timing works → Can demo
5. Add US4 (validation) + Polish → Full feature complete

### Single Developer Strategy

Execute phases sequentially:

1. Phase 1 → Phase 2 (blocking)
2. Phase 3 (US1) → Phase 4 (US2)
3. Phase 5 (US3) → Phase 6 (US4)
4. Phase 7 (Polish)

---

## Notes

- [P] tasks = different files or independent code, no dependencies
- [Story] label maps task to specific user story
- Each user story is independently testable
- Verify tests fail before implementing
- Commit after each task or logical group
- All code in single file (logging.go) - coordinate to avoid conflicts
- Test file (logging_test.go) can have parallel test function writing
