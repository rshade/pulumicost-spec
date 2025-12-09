# Tasks: SDK Support for PULUMICOST_LOG_FILE

**Input**: Design documents from `/specs/015-log-file/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Included per Constitution principle III (Test-First Protocol).

**Organization**: Tasks are grouped by user story to enable independent implementation and
testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `sdk/go/pluginsdk/` - all changes in pluginsdk package

---

## Phase 1: Setup

**Purpose**: Project initialization and constants

- [x] T001 Define file permission constant `LogFilePermissions = 0644` in
  sdk/go/pluginsdk/logging.go
- [x] T002 Define file flags constant `LogFileFlags = os.O_APPEND|os.O_CREATE|os.O_WRONLY` in
  sdk/go/pluginsdk/logging.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core NewLogWriter() function that all user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Implement `NewLogWriter() io.Writer` function skeleton in sdk/go/pluginsdk/logging.go
  that returns os.Stderr (minimal implementation)
- [x] T004 Add `NewLogWriter()` function documentation with example usage in
  sdk/go/pluginsdk/logging.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Core CLI Controls Log Destination (Priority: P1)

**Goal**: Enable plugins to write logs to a file specified by `PULUMICOST_LOG_FILE`

**Independent Test**: Set `PULUMICOST_LOG_FILE=/tmp/test.log`, run plugin, verify logs appear in
file while stderr is clean

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T005 [P] [US1] Write test `TestNewLogWriter_ValidPath` in sdk/go/pluginsdk/logging_test.go -
  verify file writer returned when env var set to valid path
- [x] T006 [P] [US1] Write test `TestNewLogWriter_FileCreated` in sdk/go/pluginsdk/logging_test.go -
  verify file is created if not exists
- [x] T007 [P] [US1] Write test `TestNewLogWriter_FileAppended` in sdk/go/pluginsdk/logging_test.go -
  verify existing file is appended (not truncated)
- [x] T008 [P] [US1] Write test `TestNewLogWriter_AllLogLevels` in sdk/go/pluginsdk/logging_test.go -
  verify debug/info/warn/error all captured in file

### Implementation for User Story 1

- [x] T009 [US1] Implement path validation in `NewLogWriter()` - check `GetLogFile()` return value
  in sdk/go/pluginsdk/logging.go
- [x] T010 [US1] Implement file opening with `os.OpenFile(path, LogFileFlags, LogFilePermissions)`
  in sdk/go/pluginsdk/logging.go
- [x] T011 [US1] Return opened file as `io.Writer` when successful in
  sdk/go/pluginsdk/logging.go
- [x] T012 [US1] Update `newDefaultLogger()` to use `NewLogWriter()` instead of hardcoded
  os.Stderr in sdk/go/pluginsdk/logging.go

**Checkpoint**: User Story 1 complete - plugins can redirect logs to file via env var

---

## Phase 4: User Story 2 - Default Behavior Without Environment Variable (Priority: P2)

**Goal**: Ensure logs appear on stderr by default (backward compatibility)

**Independent Test**: Run plugin without setting `PULUMICOST_LOG_FILE`, verify logs appear on
stderr

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T013 [P] [US2] Write test `TestNewLogWriter_EnvNotSet` in sdk/go/pluginsdk/logging_test.go -
  verify os.Stderr returned when env var not set
- [x] T014 [P] [US2] Write test `TestNewLogWriter_EmptyString` in sdk/go/pluginsdk/logging_test.go -
  verify os.Stderr returned when env var is empty string

### Implementation for User Story 2

- [x] T015 [US2] Add empty string check in `NewLogWriter()` - return os.Stderr when
  `GetLogFile() == ""` in sdk/go/pluginsdk/logging.go
- [x] T016 [US2] Verify backward compatibility - existing plugins work without modification
  (integration test with mock plugin in sdk/go/pluginsdk/logging_test.go)

**Checkpoint**: User Story 2 complete - default behavior preserved

---

## Phase 5: User Story 3 - Graceful Handling of Invalid Paths (Priority: P3)

**Goal**: Handle invalid log file paths gracefully with fallback to stderr

**Independent Test**: Set `PULUMICOST_LOG_FILE=/nonexistent/dir/test.log`, verify warning logged
to stderr and subsequent logs go to stderr

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T017 [P] [US3] Write test `TestNewLogWriter_DirectoryPath` in
  sdk/go/pluginsdk/logging_test.go - verify stderr + warning when path is a directory
- [x] T018 [P] [US3] Write test `TestNewLogWriter_NonexistentParent` in
  sdk/go/pluginsdk/logging_test.go - verify stderr + warning when parent dir doesn't exist
- [x] T019 [P] [US3] Write test `TestNewLogWriter_PermissionDenied` in
  sdk/go/pluginsdk/logging_test.go - verify stderr + warning on permission error (if testable)

### Implementation for User Story 3

- [x] T020 [US3] Add directory detection with `os.Stat()` in `NewLogWriter()` - check
  `FileInfo.IsDir()` in sdk/go/pluginsdk/logging.go
- [x] T021 [US3] Add warning logging to stderr on file open failure in `NewLogWriter()` in
  sdk/go/pluginsdk/logging.go
- [x] T022 [US3] Return os.Stderr as fallback on any open error in `NewLogWriter()` in
  sdk/go/pluginsdk/logging.go

**Checkpoint**: User Story 3 complete - invalid paths handled gracefully

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final validation

- [x] T023 [P] Update sdk/go/pluginsdk/README.md with PULUMICOST_LOG_FILE documentation
- [x] T024 [P] Add example usage to NewLogWriter() godoc in sdk/go/pluginsdk/logging.go
- [x] T025 Run `make lint` and fix any issues
- [x] T026 Run `make test` and ensure all tests pass
- [x] T027 Verify quickstart.md scenarios work correctly

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in priority order (P1 → P2 → P3)
  - Or in parallel if desired (each is independently testable)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other
  stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Tests backward compatibility
  of US1 implementation
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Adds error handling to US1
  implementation

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Implementation builds on foundational `NewLogWriter()` skeleton
- Story complete before moving to next priority

### Parallel Opportunities

- T005-T008 (US1 tests) can run in parallel
- T013-T014 (US2 tests) can run in parallel
- T017-T019 (US3 tests) can run in parallel
- T023-T024 (documentation) can run in parallel

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all tests for User Story 1 together:
Task: T005 "Write test TestNewLogWriter_ValidPath in sdk/go/pluginsdk/logging_test.go"
Task: T006 "Write test TestNewLogWriter_FileCreated in sdk/go/pluginsdk/logging_test.go"
Task: T007 "Write test TestNewLogWriter_FileAppended in sdk/go/pluginsdk/logging_test.go"
Task: T008 "Write test TestNewLogWriter_AllLogLevels in sdk/go/pluginsdk/logging_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T004)
3. Complete Phase 3: User Story 1 (T005-T012)
4. **STOP and VALIDATE**: Test `PULUMICOST_LOG_FILE=/tmp/test.log` works
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Core feature complete
3. Add User Story 2 → Test independently → Backward compatibility verified
4. Add User Story 3 → Test independently → Error handling complete
5. Complete Polish → Documentation and validation complete

### Single Developer Flow

```text
T001 → T002 → T003 → T004 → [T005-T008 parallel] → T009 → T010 → T011 → T012
→ [T013-T014 parallel] → T015 → T016 → [T017-T019 parallel] → T020 → T021 → T022
→ [T023-T024 parallel] → T025 → T026 → T027
```

---

## Notes

- [P] tasks = different files or test functions, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Run `make lint` and `make test` after each logical group
- Stop at any checkpoint to validate story independently
