# Tasks: Validation Bypass Metadata

**Input**: Design documents from `/specs/034-validation-bypass-metadata/`
**Prerequisites**: plan.md, spec.md, data-model.md, research.md, quickstart.md

**Tests**: Included as this is SDK code requiring comprehensive test coverage per project standards.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing
of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Go SDK**: `sdk/go/pricing/` for bypass types, `sdk/go/testing/` for conformance tests
- Follow existing patterns in `sdk/go/registry/domain.go` for enum validation

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and file scaffolding

- [x] T001 Create `sdk/go/pricing/bypass.go` with Apache 2.0 copyright header and package
  declaration
- [x] T002 Create `sdk/go/pricing/bypass_test.go` with Apache 2.0 copyright header and package
  declaration

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core enum types that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: User story work cannot begin until enum types and validation are complete

- [x] T003 [P] Define `BypassSeverity` type and constants (warning, error, critical) in
  `sdk/go/pricing/bypass.go`
- [x] T004 [P] Define `BypassMechanism` type and constants (flag, env_var, config, programmatic) in
  `sdk/go/pricing/bypass.go`
- [x] T005 [P] Implement `allBypassSeverities` package-level slice for zero-allocation validation in
  `sdk/go/pricing/bypass.go`
- [x] T006 [P] Implement `allBypassMechanisms` package-level slice for zero-allocation validation in
  `sdk/go/pricing/bypass.go`
- [x] T007 Implement `IsValidBypassSeverity(string) bool` validation function in
  `sdk/go/pricing/bypass.go`
- [x] T008 Implement `IsValidBypassMechanism(string) bool` validation function in
  `sdk/go/pricing/bypass.go`
- [x] T009 Implement `AllBypassSeverities() []BypassSeverity` accessor function in
  `sdk/go/pricing/bypass.go`
- [x] T010 Implement `AllBypassMechanisms() []BypassMechanism` accessor function in
  `sdk/go/pricing/bypass.go`

**Checkpoint**: Enum types ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Audit Trail for Bypassed Validations (Priority: P1) 🎯 MVP

**Goal**: Enable ValidationResult to carry complete bypass metadata for audit trails

**Independent Test**: Create a ValidationResult with bypass metadata, serialize to JSON, deserialize,
and verify all fields are preserved including timestamp, reason, operator, severity, and mechanism.

### Tests for User Story 1

- [x] T011 [P] [US1] Add enum validation tests for BypassSeverity in `sdk/go/pricing/bypass_test.go`
- [x] T012 [P] [US1] Add enum validation tests for BypassMechanism in `sdk/go/pricing/bypass_test.go`
- [x] T013 [P] [US1] Add benchmark tests for enum validation (<10 ns/op, 0 allocs) in
  `sdk/go/pricing/bypass_test.go`

### Implementation for User Story 1

- [x] T014 [US1] Define `BypassMetadata` struct with all fields (Timestamp, ValidationName,
  OriginalError, Reason, Operator, Severity, Mechanism, Truncated) in `sdk/go/pricing/bypass.go`
- [x] T015 [US1] Add JSON struct tags to BypassMetadata with snake_case naming and omitempty where
  appropriate in `sdk/go/pricing/bypass.go`
- [x] T016 [US1] Implement `NewBypassMetadata` constructor with required fields and functional
  options in `sdk/go/pricing/bypass.go`
- [x] T017 [US1] Implement `WithReason(string)` option with 500-char truncation logic in
  `sdk/go/pricing/bypass.go`
- [x] T018 [US1] Implement `WithOperator(string)` option with "unknown" default in
  `sdk/go/pricing/bypass.go`
- [x] T019 [US1] Implement `WithSeverity(BypassSeverity)` option in `sdk/go/pricing/bypass.go`
- [x] T020 [US1] Implement `WithMechanism(BypassMechanism)` option in `sdk/go/pricing/bypass.go`
- [x] T021 [US1] Implement `ValidateBypassMetadata(BypassMetadata) error` validation function in
  `sdk/go/pricing/bypass.go`
- [x] T022 [US1] Add `Bypasses []BypassMetadata` field to ValidationResult struct with JSON tag
  `json:"bypasses,omitempty"` in `sdk/go/pricing/observability.go`
- [x] T023 [US1] Add BypassMetadata constructor and validation tests in
  `sdk/go/pricing/bypass_test.go`
- [x] T024 [US1] Add reason truncation tests (exactly 500 chars, over 500 chars, under 500 chars) in
  `sdk/go/pricing/bypass_test.go`
- [x] T025 [US1] Add JSON round-trip serialization tests for ValidationResult with bypasses in
  `sdk/go/pricing/bypass_test.go`

**Checkpoint**: ValidationResult can carry bypass metadata and survives JSON serialization

---

## Phase 4: User Story 2 - Display Bypass Information in CLI Output (Priority: P2)

**Goal**: Provide formatters for displaying bypass information to operators

**Independent Test**: Format a ValidationResult with multiple bypasses and verify output includes
severity levels, validation names, and clear visual distinction from passed/failed validations.

### Tests for User Story 2

- [x] T026 [P] [US2] Add tests for FormatBypassSummary output in `sdk/go/pricing/bypass_test.go`
- [x] T027 [P] [US2] Add tests for FormatBypassDetail output in `sdk/go/pricing/bypass_test.go`

### Implementation for User Story 2

- [x] T028 [US2] Implement `FormatBypassSummary([]BypassMetadata) string` for CLI summary output in
  `sdk/go/pricing/bypass.go`
- [x] T029 [US2] Implement `FormatBypassDetail(BypassMetadata) string` for detailed single bypass
  output in `sdk/go/pricing/bypass.go`
- [x] T030 [US2] Implement `HasBypasses(ValidationResult) bool` helper function in
  `sdk/go/pricing/bypass.go`
- [x] T031 [US2] Implement `CountBypassesBySeverity([]BypassMetadata) map[BypassSeverity]int`
  aggregation helper in `sdk/go/pricing/bypass.go`

**Checkpoint**: CLI can display bypass information clearly to operators

---

## Phase 5: User Story 3 - Query Historical Bypass Events (Priority: P3)

**Goal**: Enable filtering and querying of bypass metadata for compliance analysis

**Independent Test**: Create multiple ValidationResults with different timestamps, operators, and
severities, then filter by each criterion and verify correct results.

### Tests for User Story 3

- [x] T032 [P] [US3] Add tests for FilterByTimeRange in `sdk/go/pricing/bypass_test.go`
- [x] T033 [P] [US3] Add tests for FilterByOperator in `sdk/go/pricing/bypass_test.go`
- [x] T034 [P] [US3] Add tests for FilterBySeverity in `sdk/go/pricing/bypass_test.go`

### Implementation for User Story 3

- [x] T035 [US3] Implement `FilterByTimeRange([]BypassMetadata, start, end time.Time)
  []BypassMetadata` in `sdk/go/pricing/bypass.go`
- [x] T036 [US3] Implement `FilterByOperator([]BypassMetadata, operator string) []BypassMetadata` in
  `sdk/go/pricing/bypass.go`
- [x] T037 [US3] Implement `FilterBySeverity([]BypassMetadata, severity BypassSeverity)
  []BypassMetadata` in `sdk/go/pricing/bypass.go`
- [x] T038 [US3] Implement `FilterByMechanism([]BypassMetadata, mechanism BypassMechanism)
  []BypassMetadata` in `sdk/go/pricing/bypass.go`

**Checkpoint**: Bypass metadata can be queried and filtered for compliance reporting

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Conformance tests, documentation, and final validation

- [x] T039 [P] Create `sdk/go/testing/bypass_conformance_test.go` with Apache 2.0 header
- [x] T040 [P] Add conformance test for bypass metadata JSON round-trip in
  `sdk/go/testing/bypass_conformance_test.go`
- [x] T041 [P] Add conformance test for ValidationResult backward compatibility (empty Bypasses) in
  `sdk/go/testing/bypass_conformance_test.go`
- [x] T042 Run `make test` and verify all tests pass
- [x] T043 Run `make lint` and fix any linting issues
- [x] T044 Verify quickstart.md examples compile and work correctly
- [x] T045 [P] Update `sdk/go/pricing/CLAUDE.md` with bypass metadata documentation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can proceed in priority order (P1 → P2 → P3)
  - US2 and US3 can technically start in parallel after US1 core types are done
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after US1 T022 (ValidationResult extension) - Uses BypassMetadata
  types
- **User Story 3 (P3)**: Can start after US1 T014 (BypassMetadata struct) - Uses BypassMetadata
  types

### Within Each User Story

- Tests written first (TDD approach per project standards)
- Types/structs before functions
- Validation functions before helper functions
- Unit tests before integration/conformance tests

### Parallel Opportunities

**Phase 2 (Foundational):**

```text
T003 + T004 + T005 + T006 can run in parallel (different constants/slices)
```

**Phase 3 (US1):**

```text
T011 + T012 + T013 can run in parallel (different test functions)
```

**Phase 4 (US2):**

```text
T026 + T027 can run in parallel (different test functions)
```

**Phase 5 (US3):**

```text
T032 + T033 + T034 can run in parallel (different test functions)
```

**Phase 6 (Polish):**

```text
T039 + T040 + T041 + T045 can run in parallel (different files)
```

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all tests for User Story 1 together:
Task: "Add enum validation tests for BypassSeverity in sdk/go/pricing/bypass_test.go"
Task: "Add enum validation tests for BypassMechanism in sdk/go/pricing/bypass_test.go"
Task: "Add benchmark tests for enum validation in sdk/go/pricing/bypass_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T010)
3. Complete Phase 3: User Story 1 (T011-T025)
4. **STOP and VALIDATE**: Run `go test ./sdk/go/pricing/...` and verify JSON serialization
5. Deploy/demo if ready - ValidationResult now carries bypass metadata

### Incremental Delivery

1. Complete Setup + Foundational → Enum types ready
2. Add User Story 1 → Test JSON round-trip → MVP complete (audit trail works)
3. Add User Story 2 → Test formatting → CLI can display bypasses
4. Add User Story 3 → Test filtering → Compliance queries work
5. Each story adds value without breaking previous stories

---

## Notes

- [P] tasks = different files or independent test functions
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Follow existing patterns in `sdk/go/registry/domain.go` for enum implementation
- Zero-allocation validation is required per performance goals (<10 ns/op)
- JSON struct tags use snake_case per existing project conventions
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- **Retention Policy (FR-012)**: The 90-day retention requirement is a caller responsibility, not
  SDK implementation. The SDK provides the data structures; callers must implement their own
  retention policies for stored bypass metadata.
