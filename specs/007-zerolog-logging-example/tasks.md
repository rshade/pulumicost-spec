# Tasks: Structured Logging Example for EstimateCost

**Input**: Design documents from `/specs/007-zerolog-logging-example/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: This feature is a documentation example that IS a test. No separate test tasks needed.

**Organization**: Tasks grouped by user story to enable independent verification of each logging pattern.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Target file**: `sdk/go/testing/integration_test.go`
- **Dependency**: `sdk/go/logging/` (from 005-zerolog)

---

## Phase 1: Setup (Verification)

**Purpose**: Verify prerequisites are met before implementation

- [x] T001 Verify 005-zerolog utilities exist in sdk/go/pluginsdk/ (NewPluginLogger, field constants)
- [x] T002 Verify EstimateCost RPC is available in MockPlugin in sdk/go/testing/mock_plugin.go
- [x] T003 Review existing test patterns in sdk/go/testing/integration_test.go

---

## Phase 2: Foundational (Test Function Scaffold)

**Purpose**: Create the test function structure that all user stories will use

**CRITICAL**: This scaffold must be complete before implementing logging examples

- [x] T004 Add TestStructuredLoggingExample function scaffold in sdk/go/testing/integration_test.go
- [x] T005 Add import statements for zerolog, bytes, encoding/json, sdk/go/pluginsdk in sdk/go/testing/integration_test.go
- [x] T006 Add bytes.Buffer setup for log capture in TestStructuredLoggingExample
- [x] T007 Add helper function parseMultipleLogEntries, assertLogContains, assertLogNotContains in sdk/go/testing/integration_test.go

**Checkpoint**: Scaffold ready - logging example subtests can now be added

---

## Phase 3: User Story 1 - Plugin Developer Learns Logging Patterns (Priority: P1)

**Goal**: Demonstrate creating a configured logger and logging request/response patterns

**Independent Test**: Run `go test -v ./sdk/go/testing/ -run TestStructuredLoggingExample/RequestLogging`
and verify JSON output contains trace_id, operation, resource_type fields

### Implementation for User Story 1

- [x] T008 [US1] Add t.Run("RequestLogging") subtest demonstrating request logging in sdk/go/testing/integration_test.go
- [x] T009 [US1] Demonstrate NewPluginLogger with plugin name and version in RequestLogging subtest
- [x] T010 [US1] Demonstrate logging resource_type and attribute_count (not values) in RequestLogging subtest
- [x] T011 [US1] Add t.Run("SuccessResponseLogging") subtest in sdk/go/testing/integration_test.go
- [x] T012 [US1] Demonstrate logging cost_monthly, currency, and duration_ms in SuccessResponseLogging subtest
- [x] T013 [US1] Demonstrate LogOperation timing helper in SuccessResponseLogging subtest
- [x] T014 [US1] Add code comments explaining logging best practices for FR-007 compliance

**Checkpoint**: User Story 1 complete - developers can see request/response logging patterns

---

## Phase 4: User Story 2 - Plugin Developer Implements Error Logging (Priority: P1)

**Goal**: Demonstrate error logging with correlation IDs and error context

**Independent Test**: Run `go test -v ./sdk/go/testing/ -run TestStructuredLoggingExample/ErrorLogging`
and verify JSON output contains trace_id, error_code, error message fields

### Implementation for User Story 2

- [x] T015 [US2] Add t.Run("ErrorLogging") subtest using ConfigurableErrorMockPlugin in sdk/go/testing/integration_test.go
- [x] T016 [US2] Demonstrate error logging with FieldErrorCode and gRPC status code in ErrorLogging subtest
- [x] T017 [US2] Demonstrate including original request context (resource_type) in error logs
- [x] T018 [US2] Add t.Run("CorrelationIDPropagation") subtest in sdk/go/testing/integration_test.go
- [x] T019 [US2] Demonstrate ContextWithTraceID and TraceIDFromContext usage in CorrelationIDPropagation subtest
- [x] T020 [US2] Add assertions verifying trace_id appears in all log entries
- [x] T021 [US2] Add code comments explaining correlation ID best practices

**Checkpoint**: User Story 2 complete - developers can see error logging and tracing patterns

---

## Phase 5: User Story 3 - Operator Monitors EstimateCost Health (Priority: P2)

**Goal**: Demonstrate consistent log structure for operational monitoring

**Independent Test**: Verify all log outputs use standard field names parseable by JSON tools

### Implementation for User Story 3

- [x] T022 [US3] Add log structure assertions verifying all entries use FieldTraceID, FieldOperation constants in sdk/go/testing/integration_test.go
- [x] T023 [US3] Add t.Run("LogStructureValidation") subtest verifying JSON parseable output
- [x] T024 [US3] Demonstrate filterability by operation field across multiple log entries
- [x] T025 [US3] Add code comments explaining operational monitoring considerations

**Checkpoint**: User Story 3 complete - operators can understand log query patterns

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

- [x] T026 Run go test -v ./sdk/go/testing/ -run TestStructuredLoggingExample to verify all subtests pass
- [x] T027 Run make lint to ensure code quality
- [x] T028 [P] Update sdk/go/testing/CLAUDE.md with logging example documentation
- [x] T029 Validate implementation matches quickstart.md patterns

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verify prerequisites first
- **Foundational (Phase 2)**: Depends on Setup - creates test scaffold
- **User Story 1 (Phase 3)**: Depends on Foundational - adds request/response logging
- **User Story 2 (Phase 4)**: Depends on Foundational - can run parallel to US1 (different subtests)
- **User Story 3 (Phase 5)**: Depends on US1 and US2 (validates their output)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Phase 2 scaffold - demonstrates normal flow logging
- **User Story 2 (P1)**: Depends on Phase 2 scaffold - demonstrates error flow logging
- **User Story 3 (P2)**: Depends on US1 and US2 logs to validate structure

### Within Each User Story

- Subtest creation before implementation details
- Logging demonstrations before assertions
- Code comments after implementation

### Parallel Opportunities

- T008, T011 (US1 subtests) can be created in parallel (different subtests)
- T015, T018 (US2 subtests) can be created in parallel (different subtests)
- US1 and US2 can be worked on in parallel after Phase 2 (no dependencies between them)

---

## Parallel Example: Phase 2 Setup

```bash
# Import statements and buffer setup can be done in sequence (same file):
Task: "Add import statements for zerolog, bytes, encoding/json"
Task: "Add bytes.Buffer setup for log capture"
Task: "Add helper function parseLogEntry"
```

## Parallel Example: User Story 1 + User Story 2

```bash
# After Phase 2 scaffold is complete, these can run in parallel:
# Developer A: User Story 1 (RequestLogging, SuccessResponseLogging subtests)
# Developer B: User Story 2 (ErrorLogging, CorrelationIDPropagation subtests)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup verification
2. Complete Phase 2: Test scaffold
3. Complete Phase 3: User Story 1 (request/response logging)
4. **STOP and VALIDATE**: Run `go test -v ./sdk/go/testing/ -run TestStructuredLoggingExample`
5. Merge if US1 is sufficient for initial documentation

### Full Implementation

1. Complete Setup + Foundational
2. Add User Story 1 → Verify independently
3. Add User Story 2 → Verify independently
4. Add User Story 3 → Validate log structure
5. Polish and final validation

---

## Notes

- This feature adds to existing file - no new file creation needed
- All logging examples use the 005-zerolog SDK utilities
- Example code serves as documentation (FR-007: code comments required)
- Focus on demonstrating patterns, not exhaustive test coverage
- Avoid logging sensitive attribute values - log count only
