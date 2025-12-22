# Tasks: PluginSDK Conformance Testing Adapters

**Input**: Design documents from `/specs/011-pluginsdk-conformance/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Tests**: Included (Constitution Principle III: Test-First Protocol is NON-NEGOTIABLE)

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Go SDK**: `sdk/go/pluginsdk/` for new conformance adapter code
- **Tests**: `sdk/go/pluginsdk/conformance_test.go`
- **Documentation**: `sdk/go/pluginsdk/README.md`

---

## Phase 1: Setup

**Purpose**: Create conformance.go file structure with package imports and type aliases

- [x] T001 Create sdk/go/pluginsdk/conformance.go with package declaration and imports
- [x] T002 Add type aliases for ConformanceResult, ConformanceLevel, ResultSummary from
      sdk/go/testing in sdk/go/pluginsdk/conformance.go
- [x] T003 [P] Create sdk/go/pluginsdk/conformance_test.go with package declaration and test
      imports
- [x] T004 [P] Verify no import cycle by running `go build ./sdk/go/pluginsdk/...`

---

## Phase 2: Foundational (Nil Plugin Validation)

**Purpose**: Implement core nil validation that ALL adapter functions depend on

**CRITICAL**: This validation must work before any conformance adapter can be used safely

- [x] T005 Implement internal `validatePlugin(plugin Plugin) error` function returning
      descriptive error for nil in sdk/go/pluginsdk/conformance.go
- [x] T006 Write test `TestValidatePluginNil` verifying nil returns error in
      sdk/go/pluginsdk/conformance_test.go
- [x] T007 Write test `TestValidatePluginValid` verifying non-nil passes in
      sdk/go/pluginsdk/conformance_test.go
- [x] T008 Run tests to confirm T006 and T007 pass with `go test -v ./sdk/go/pluginsdk/...`

**Checkpoint**: Nil validation foundation ready - user story implementation can begin

---

## Phase 3: User Story 1 - Run Basic Conformance (Priority: P1) MVP

**Goal**: Plugin developers can run basic conformance tests with single function call

**Independent Test**: Call `RunBasicConformance(plugin)` and verify ConformanceResult returned
with pass/fail status

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T009 [US1] Write test `TestRunBasicConformanceNilPlugin` expecting error for nil plugin in
      sdk/go/pluginsdk/conformance_test.go
- [x] T010 [US1] Write test `TestRunBasicConformanceValidPlugin` using mock plugin expecting
      ConformanceResult in sdk/go/pluginsdk/conformance_test.go
- [x] T011 [US1] Run tests to confirm T009 and T010 FAIL (function not implemented yet)

### Implementation for User Story 1

- [x] T012 [US1] Implement `RunBasicConformance(plugin Plugin) (*ConformanceResult, error)` in
      sdk/go/pluginsdk/conformance.go
- [x] T013 [US1] Add godoc comment documenting function purpose, parameters, return values
- [x] T014 [US1] Run tests to confirm T009 and T010 now PASS

**Checkpoint**: User Story 1 complete - basic conformance adapter functional

---

## Phase 4: User Story 2 - Run Standard Conformance (Priority: P1)

**Goal**: Plugin developers can run standard (production-ready) conformance tests

**Independent Test**: Call `RunStandardConformance(plugin)` and verify ConformanceResult shows
Standard level tests included

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T015 [US2] Write test `TestRunStandardConformanceNilPlugin` expecting error for nil plugin
      in sdk/go/pluginsdk/conformance_test.go
- [x] T016 [US2] Write test `TestRunStandardConformanceValidPlugin` using mock plugin expecting
      ConformanceResult in sdk/go/pluginsdk/conformance_test.go
- [x] T017 [US2] Run tests to confirm T015 and T016 FAIL (function not implemented yet)

### Implementation for User Story 2

- [x] T018 [US2] Implement `RunStandardConformance(plugin Plugin) (*ConformanceResult, error)` in
      sdk/go/pluginsdk/conformance.go
- [x] T019 [US2] Add godoc comment documenting function purpose, parameters, return values
- [x] T020 [US2] Run tests to confirm T015 and T016 now PASS

**Checkpoint**: User Stories 1 AND 2 complete - basic and standard conformance adapters functional

---

## Phase 5: User Story 3 - Run Advanced Conformance (Priority: P2)

**Goal**: Plugin developers can run advanced (high-performance) conformance tests

**Independent Test**: Call `RunAdvancedConformance(plugin)` and verify ConformanceResult includes
performance benchmarks

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T021 [US3] Write test `TestRunAdvancedConformanceNilPlugin` expecting error for nil plugin
      in sdk/go/pluginsdk/conformance_test.go
- [x] T022 [US3] Write test `TestRunAdvancedConformanceValidPlugin` using mock plugin expecting
      ConformanceResult in sdk/go/pluginsdk/conformance_test.go
- [x] T023 [US3] Run tests to confirm T021 and T022 FAIL (function not implemented yet)

### Implementation for User Story 3

- [x] T024 [US3] Implement `RunAdvancedConformance(plugin Plugin) (*ConformanceResult, error)` in
      sdk/go/pluginsdk/conformance.go
- [x] T025 [US3] Add godoc comment documenting function purpose, parameters, return values
- [x] T026 [US3] Run tests to confirm T021 and T022 now PASS

**Checkpoint**: All three conformance level adapters complete

---

## Phase 6: User Story 4 - Print Conformance Report (Priority: P2)

**Goal**: Plugin developers can print formatted conformance report to test output

**Independent Test**: Call `PrintConformanceReport(t, result)` and verify formatted output
includes pass/fail counts and conformance level

### Tests for User Story 4

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T027 [US4] Write test `TestPrintConformanceReportNilResult` verifying no panic for nil
      result in sdk/go/pluginsdk/conformance_test.go
- [x] T028 [US4] Write test `TestPrintConformanceReportValidResult` verifying formatted output in
      sdk/go/pluginsdk/conformance_test.go
- [x] T029 [US4] Run tests to confirm T027 and T028 FAIL (function not implemented yet)

### Implementation for User Story 4

- [x] T030 [US4] Implement `PrintConformanceReport(t *testing.T, result *ConformanceResult)` in
      sdk/go/pluginsdk/conformance.go
- [x] T031 [US4] Add godoc comment documenting function purpose, parameters
- [x] T032 [US4] Run tests to confirm T027 and T028 now PASS

**Checkpoint**: All four adapter functions complete - full feature implemented

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, linting, and validation

- [x] T033 [P] Add Conformance Testing section to sdk/go/pluginsdk/README.md with usage examples
- [x] T034 [P] Run `make lint` to verify code passes all linting rules
- [x] T035 [P] Run `make test` to verify all tests pass including new conformance tests
- [x] T036 [P] Update sdk/go/pluginsdk/CLAUDE.md with conformance testing patterns if file exists,
      otherwise skip
- [x] T037 Verify integration by running full conformance suite through adapters manually
- [x] T038 Run quickstart.md validation - ensure examples compile and execute correctly

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on T001-T002 completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - US1 and US2 are both P1 priority - can proceed in parallel or sequentially
  - US3 and US4 are both P2 priority - can proceed after P1 stories
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Independent of US1
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Independent of US1/US2
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Independent but uses
  ConformanceResult from US1/US2/US3

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Implementation follows test
- Godoc documentation with implementation
- Tests PASS to complete story

### Parallel Opportunities

- T003 and T004 can run in parallel within Setup
- US1 (T009-T014) and US2 (T015-T020) can run in parallel (both P1, different functions)
- US3 (T021-T026) and US4 (T027-T032) can run in parallel (both P2, different functions)
- T033, T034, T035 can all run in parallel in Polish phase

---

## Parallel Example: User Story 1 and 2 Together

```bash
# Launch US1 and US2 tests in parallel (different test functions):
Task: "Write test TestRunBasicConformanceNilPlugin in conformance_test.go"
Task: "Write test TestRunStandardConformanceNilPlugin in conformance_test.go"

# After tests fail, implement in parallel (different functions in same file):
Task: "Implement RunBasicConformance in conformance.go"
Task: "Implement RunStandardConformance in conformance.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (nil validation)
3. Complete Phase 3: User Story 1 (RunBasicConformance)
4. **STOP and VALIDATE**: Test US1 independently
5. Deploy/demo if ready - basic conformance adapter available

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → MVP ready (basic conformance)
3. Add User Story 2 → Test independently → Standard conformance available
4. Add User Story 3 → Test independently → Advanced conformance available
5. Add User Story 4 → Test independently → Report printing available
6. Polish phase → Documentation and validation complete

### Recommended Execution Order

Since this is a small feature with thin wrappers:

1. Phase 1-2: Setup and Foundational (T001-T008)
2. Phase 3: User Story 1 - RunBasicConformance (T009-T014)
3. Phase 4: User Story 2 - RunStandardConformance (T015-T020)
4. Phase 5: User Story 3 - RunAdvancedConformance (T021-T026)
5. Phase 6: User Story 4 - PrintConformanceReport (T027-T032)
6. Phase 7: Polish (T033-T038)

---

## Notes

- [P] tasks = different files or independent functions, no dependencies
- [Story] label maps task to specific user story for traceability
- Each adapter function (US1-US3) follows identical pattern: nil check → NewServer → delegate
- Test-First Protocol strictly followed per Constitution Principle III
- All functions go in single file: sdk/go/pluginsdk/conformance.go
- All tests go in single file: sdk/go/pluginsdk/conformance_test.go
- Verify tests FAIL before implementation, then PASS after
- Commit after each user story completion
