# Tasks: FallbackHint Enum for Plugin Orchestration

**Input**: Design documents from `/specs/001-fallback-hint/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅, quickstart.md ✅

**Tests**: Tests are included as this feature modifies core protocol and SDK components.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto definition**: `proto/pulumicost/v1/costsource.proto`
- **SDK helpers**: `sdk/go/pluginsdk/`
- **Generated code**: `sdk/go/proto/` (regenerated, do not edit manually)
- **Spec contracts**: `specs/001-fallback-hint/contracts/costsource.proto`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Validate contracts and prepare for proto modification

- [x] T001 Validate contract proto matches expected changes in
      specs/001-fallback-hint/contracts/costsource.proto
- [x] T002 [P] Review existing GetActualCostResponse usage patterns in sdk/go/pluginsdk/helpers.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Proto definition changes that MUST be complete before any SDK work

**⚠️ CRITICAL**: No SDK or user story work can begin until proto changes are complete and regenerated

- [x] T003 Add FallbackHint enum to proto/pulumicost/v1/costsource.proto with four values
- [x] T004 Add fallback_hint field to GetActualCostResponse message in
      proto/pulumicost/v1/costsource.proto
- [x] T005 Run `make generate` to regenerate Go bindings in sdk/go/proto/
- [x] T006 Run `buf lint` to validate proto changes
- [x] T007 Verify generated code compiles with `go build ./...`

**Checkpoint**: Proto changes complete - SDK implementation can now begin

---

## Phase 3: User Story 5 - Backwards Compatibility (Priority: P1) 🎯 MVP

**Goal**: Ensure existing plugins without hint field continue to work - default/unspecified treated as "no fallback"

**Independent Test**: Call legacy plugin (without hint), verify core treats as no-fallback

### Tests for User Story 5

- [x] T008 [P] [US5] Add test for default hint value (unspecified = 0) behavior in
      sdk/go/pluginsdk/helpers_test.go
- [x] T009 [P] [US5] Add test verifying CreateActualCostResponse returns unspecified hint by default in
      sdk/go/pluginsdk/helpers_test.go
- [x] T009a [P] [US5] Add test verifying wire format (JSON marshaling) omits fallback_hint when
      zero/unspecified (SC-004) in sdk/go/pluginsdk/helpers_test.go

### Implementation for User Story 5

- [x] T010 [US5] Verify CreateActualCostResponse continues to work without modification in
      sdk/go/pluginsdk/helpers.go
- [x] T011 [US5] Add documentation comment to CreateActualCostResponse explaining default hint behavior
      in sdk/go/pluginsdk/helpers.go

**Checkpoint**: Backwards compatibility verified - existing code works without changes

---

## Phase 4: User Story 1 - Plugin Returns Data (Priority: P1)

**Goal**: Plugin can return data with explicit "no fallback" signal using functional options

**Independent Test**: Call GetActualCost with data, verify response has FALLBACK_HINT_NONE

### Tests for User Story 1

- [x] T012 [P] [US1] Add test for WithFallbackHint(NONE) option in sdk/go/pluginsdk/helpers_test.go
- [x] T013 [P] [US1] Add test for NewActualCostResponse with results and explicit NONE hint in
      sdk/go/pluginsdk/helpers_test.go

### Implementation for User Story 1

- [x] T014 [US1] Define ActualCostResponseOption type for functional options in
      sdk/go/pluginsdk/helpers.go
- [x] T015 [US1] Implement WithFallbackHint option function in sdk/go/pluginsdk/helpers.go
- [x] T016 [US1] Implement WithResults option function in sdk/go/pluginsdk/helpers.go
- [x] T017 [US1] Implement NewActualCostResponse constructor using functional options in
      sdk/go/pluginsdk/helpers.go
- [x] T018 [US1] Add documentation comments for functional options pattern in
      sdk/go/pluginsdk/helpers.go (FR-008)

**Checkpoint**: Plugins can explicitly signal no-fallback with data

---

## Phase 5: User Story 2 - No Data, Recommend Fallback (Priority: P1)

**Goal**: Plugin can return empty response with FALLBACK_RECOMMENDED signal

**Independent Test**: Call GetActualCost with no data, verify FALLBACK_RECOMMENDED hint

### Tests for User Story 2

- [x] T019 [P] [US2] Add test for WithFallbackHint(RECOMMENDED) with empty results in
      sdk/go/pluginsdk/helpers_test.go
- [x] T020 [P] [US2] Add test for NewActualCostResponse with nil results and RECOMMENDED hint in
      sdk/go/pluginsdk/helpers_test.go

### Implementation for User Story 2

- [x] T021 [US2] Verify FALLBACK_HINT_RECOMMENDED correctly set via WithFallbackHint in
      sdk/go/pluginsdk/helpers.go
- [x] T022 [US2] Add example usage for "no data" scenario in documentation comment in
      sdk/go/pluginsdk/helpers.go

**Checkpoint**: Plugins can signal "no data, try others"

---

## Phase 6: User Story 3 - Cannot Handle Request (Priority: P2)

**Goal**: Plugin can signal FALLBACK_REQUIRED for unsupported resource types

**Independent Test**: Call GetActualCost with unsupported type, verify FALLBACK_REQUIRED hint

### Tests for User Story 3

- [x] T023 [P] [US3] Add test for WithFallbackHint(REQUIRED) in sdk/go/pluginsdk/helpers_test.go
- [x] T024 [P] [US3] Add test for NewActualCostResponse with REQUIRED hint for unsupported type in
      sdk/go/pluginsdk/helpers_test.go

### Implementation for User Story 3

- [x] T025 [US3] Verify FALLBACK_HINT_REQUIRED correctly set via WithFallbackHint in
      sdk/go/pluginsdk/helpers.go
- [x] T026 [US3] Add example usage for "cannot handle" scenario in documentation comment in
      sdk/go/pluginsdk/helpers.go

**Checkpoint**: Plugins can signal "not my job, must try others"

---

## Phase 7: User Story 4 - Error Handling (Priority: P2)

**Goal**: Ensure errors are returned as gRPC errors, not hint values - clarify this in documentation

**Independent Test**: Simulate API failure, verify gRPC error returned (not hint)

### Tests for User Story 4

- [x] T027 [P] [US4] Add test verifying GetActualCost returns error for API failures (not hint) in
      sdk/go/pluginsdk/helpers_test.go

### Implementation for User Story 4

- [x] T028 [US4] Add documentation clarifying error vs hint distinction in
      sdk/go/pluginsdk/helpers.go (FR-008)
- [x] T029 [US4] Add example in documentation showing when to return error vs hint in
      sdk/go/pluginsdk/helpers.go (FR-008)

**Checkpoint**: Error semantics clearly documented

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, validation, and final polish

- [x] T030 [P] Update quickstart.md with new functional options pattern examples (FR-009)
- [x] T031 [P] Add FallbackHint decision matrix to quickstart.md in
      specs/001-fallback-hint/quickstart.md (FR-009, SC-005)
- [x] T032 [P] Document zero-cost vs no-data distinction in quickstart.md in
      specs/001-fallback-hint/quickstart.md (FR-010)
- [x] T033 [P] Document edge case: data + hint conflict (data wins) in
      specs/001-fallback-hint/quickstart.md
- [x] T034 [P] Document edge case: unrecognized hint values (treat as UNSPECIFIED) in
      specs/001-fallback-hint/quickstart.md
- [x] T035 [P] Document edge case: fallback chain termination in
      specs/001-fallback-hint/quickstart.md
- [x] T036 Run `make lint` to verify all code passes linting
- [x] T037 Run `make test` to verify all tests pass
- [x] T038 Run `make validate` for full validation pipeline
- [x] T039 Verify contract proto in specs/001-fallback-hint/contracts/costsource.proto matches final
      implementation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories (proto changes required first)
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - US5 (Backwards Compat) is independent - can start immediately after proto
  - US1 (Data + NONE) is independent - can start immediately after proto
  - US2 (No Data + RECOMMENDED) depends on US1 functional options implementation
  - US3 (Cannot Handle + REQUIRED) depends on US1 functional options implementation
  - US4 (Error Handling) is independent - documentation only
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 5 (P1)**: No implementation dependencies - verification only
- **User Story 1 (P1)**: Core implementation - functional options pattern
- **User Story 2 (P1)**: Depends on US1 (uses WithFallbackHint from US1)
- **User Story 3 (P2)**: Depends on US1 (uses WithFallbackHint from US1)
- **User Story 4 (P2)**: Independent documentation task

### Within Each User Story

- Tests written first to define expected behavior
- Implementation follows tests
- Documentation updated after implementation

### Parallel Opportunities

- T001, T002 can run in parallel (Setup phase)
- T008, T009 can run in parallel (US5 tests)
- T012, T013 can run in parallel (US1 tests)
- T019, T020 can run in parallel (US2 tests)
- T023, T024 can run in parallel (US3 tests)
- T030, T031, T032, T033, T034, T035 can run in parallel (Polish documentation)

---

## Parallel Example: User Story 1

```bash
# Launch tests in parallel:
Task: "Add test for WithFallbackHint(NONE) option in sdk/go/pluginsdk/helpers_test.go"
Task: "Add test for NewActualCostResponse with results and explicit NONE hint in sdk/go/pluginsdk/helpers_test.go"

# Then sequential implementation:
Task: "Define ActualCostResponseOption type for functional options in sdk/go/pluginsdk/helpers.go"
Task: "Implement WithFallbackHint option function in sdk/go/pluginsdk/helpers.go"
Task: "Implement WithResults option function in sdk/go/pluginsdk/helpers.go"
Task: "Implement NewActualCostResponse constructor using functional options in sdk/go/pluginsdk/helpers.go"
```

---

## Implementation Strategy

### MVP First (Proto + Backwards Compatibility + Core Options)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - proto changes)
3. Complete Phase 3: User Story 5 (Backwards Compatibility)
4. Complete Phase 4: User Story 1 (Functional Options)
5. **STOP and VALIDATE**: Test proto changes and SDK options independently
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Proto changes ready
2. Add US5 → Test backwards compatibility → Verify existing code works
3. Add US1 → Test functional options → Core SDK usable (MVP!)
4. Add US2 + US3 → Additional hint values documented and tested
5. Add US4 → Error handling documented
6. Polish → All documentation complete

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (proto changes)
2. Once proto regenerated:
   - Developer A: User Story 5 (backwards compat verification)
   - Developer B: User Story 1 (functional options implementation)
3. After US1 complete:
   - Developer A: User Story 2 + US3 (uses US1 options)
   - Developer B: User Story 4 (error docs - independent)
4. All: Polish phase documentation

---

## Summary

| Metric | Value |
|--------|-------|
| **Total Tasks** | 40 |
| **Setup Tasks** | 2 |
| **Foundational Tasks** | 5 |
| **US5 (Backwards Compat) Tasks** | 5 |
| **US1 (Data + NONE) Tasks** | 7 |
| **US2 (No Data + RECOMMENDED) Tasks** | 4 |
| **US3 (Cannot Handle + REQUIRED) Tasks** | 4 |
| **US4 (Error Handling) Tasks** | 3 |
| **Polish Tasks** | 10 |
| **Parallelizable Tasks** | 17 |
| **MVP Scope** | Phases 1-4 (US5 + US1) |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- The functional options pattern (WithFallbackHint, WithResults) follows Go idioms from research.md
- Default hint value (0 = UNSPECIFIED) preserves backwards compatibility per spec
