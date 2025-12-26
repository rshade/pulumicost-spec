# Tasks: Contextual FinOps Validation

**Input**: Design documents from `/specs/027-finops-validation/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Included per Constitution Principle III (Test-First Protocol) and SC-002 (100%
code coverage requirement).

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **SDK extension**: `sdk/go/pluginsdk/` at repository root
- **Testing framework**: `sdk/go/testing/`

---

## Phase 1: Setup

**Purpose**: Create new files and verify existing structure

- [x] T001 Verify existing focus_conformance.go structure in sdk/go/pluginsdk/focus_conformance.go
- [x] T002 [P] Create validation_error.go skeleton in sdk/go/pluginsdk/validation_error.go
- [x] T003 [P] Create validation_options.go skeleton in sdk/go/pluginsdk/validation_options.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types that ALL user stories depend on

**CRITICAL**: No user story work can begin until these types exist

- [x] T004 Define ValidationError struct with four fields in sdk/go/pluginsdk/validation_error.go
- [x] T005 Implement Error() method for ValidationError in sdk/go/pluginsdk/validation_error.go
- [x] T006 [P] Write tests for ValidationError in sdk/go/pluginsdk/validation_error_test.go
- [x] T007 Define ValidationMode enum (FailFast, Aggregate) in sdk/go/pluginsdk/validation_options.go
- [x] T008 Define ValidationOptions struct in sdk/go/pluginsdk/validation_options.go
- [x] T009 [P] Write tests for ValidationMode and ValidationOptions in sdk/go/pluginsdk/validation_options_test.go
- [x] T010 Define sentinel error variables for all validation rules in sdk/go/pluginsdk/focus_conformance.go
- [x] T011 Create ValidateFocusRecordWithOptions function signature in sdk/go/pluginsdk/focus_conformance.go
- [x] T012 Wire ValidateFocusRecord to call ValidateFocusRecordWithOptions with FailFast mode in sdk/go/pluginsdk/focus_conformance.go

**Checkpoint**: Foundation ready - ValidationError, ValidationMode, and function signatures exist

---

## Phase 3: User Story 1 - Cost Relationship Validation (Priority: P1)

**Goal**: Validate EffectiveCost <= BilledCost and ListCost >= EffectiveCost

**Independent Test**: Submit FocusCostRecord with various cost field combinations and verify
correct error/success responses

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T013 [P] [US1] Test EffectiveCost > BilledCost returns error in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T014 [P] [US1] Test ListCost < EffectiveCost returns error in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T015 [P] [US1] Test valid cost hierarchy passes in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T016 [P] [US1] Test zero costs (free tier) pass validation in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T017 [P] [US1] Test negative costs (credits) are exempt in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T018 [P] [US1] Test ChargeClass CORRECTION exempts cost hierarchy rules in sdk/go/pluginsdk/focus_conformance_test.go

### Implementation for User Story 1

- [x] T019 [US1] Implement validateCostHierarchy helper function in sdk/go/pluginsdk/focus_conformance.go
- [x] T020 [US1] Add validateCostHierarchy call to validateBusinessRules in sdk/go/pluginsdk/focus_conformance.go
- [x] T021 [US1] Run tests and verify all US1 tests pass with go test ./sdk/go/pluginsdk/...

**Checkpoint**: Cost relationship validation works independently

---

## Phase 4: User Story 2 - Commitment Discount Consistency (Priority: P1)

**Goal**: Validate CommitmentDiscountId and CommitmentDiscountStatus field dependencies

**Independent Test**: Submit records with various commitment field combinations and verify
consistency rules enforced

### Tests for User Story 2

- [x] T022 [P] [US2] Test CommitmentDiscountId set + Usage category requires status in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T023 [P] [US2] Test CommitmentDiscountStatus set requires CommitmentDiscountId in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T024 [P] [US2] Test consistent commitment fields pass validation in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T025 [P] [US2] Test Purchase category does not require status in sdk/go/pluginsdk/focus_conformance_test.go

### Implementation for User Story 2

- [x] T026 [US2] Implement validateCommitmentDiscountConsistency helper in sdk/go/pluginsdk/focus_conformance.go
- [x] T027 [US2] Add validateCommitmentDiscountConsistency call to validateBusinessRules in sdk/go/pluginsdk/focus_conformance.go
- [x] T028 [US2] Run tests and verify all US2 tests pass with go test ./sdk/go/pluginsdk/...

**Checkpoint**: Commitment discount consistency validation works independently

---

## Phase 5: User Story 3 - Pricing Model Consistency (Priority: P2)

**Goal**: Validate PricingQuantity requires PricingUnit

**Independent Test**: Submit records with pricing quantity > 0 and verify unit is required

### Tests for User Story 3

- [x] T029 [P] [US3] Test PricingQuantity > 0 without PricingUnit returns error in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T030 [P] [US3] Test PricingQuantity > 0 with PricingUnit passes in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T031 [P] [US3] Test PricingQuantity = 0 does not require unit in sdk/go/pluginsdk/focus_conformance_test.go

### Implementation for User Story 3

- [x] T032 [US3] Implement validatePricingConsistency helper in sdk/go/pluginsdk/focus_conformance.go
- [x] T033 [US3] Add validatePricingConsistency call to validateBusinessRules in sdk/go/pluginsdk/focus_conformance.go
- [x] T034 [US3] Run tests and verify all US3 tests pass with go test ./sdk/go/pluginsdk/...

**Checkpoint**: Pricing model consistency validation works independently

---

## Phase 6: User Story 4 - Capacity Reservation Consistency (Priority: P2)

**Goal**: Validate CapacityReservationId requires CapacityReservationStatus for Usage charges

**Independent Test**: Submit records with capacity reservation fields and verify consistency

### Tests for User Story 4

- [x] T035 [P] [US4] Test CapacityReservationId set + Usage requires status in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T036 [P] [US4] Test consistent capacity reservation fields pass in sdk/go/pluginsdk/focus_conformance_test.go

### Implementation for User Story 4

- [x] T037 [US4] Implement validateCapacityReservationConsistency helper in sdk/go/pluginsdk/focus_conformance.go
- [x] T038 [US4] Add validateCapacityReservationConsistency call to validateBusinessRules in sdk/go/pluginsdk/focus_conformance.go
- [x] T039 [US4] Run tests and verify all US4 tests pass with go test ./sdk/go/pluginsdk/...

**Checkpoint**: Capacity reservation consistency validation works independently

---

## Phase 7: User Story 5 - Validation Error Aggregation (Priority: P3)

**Goal**: Support aggregate mode that collects all errors instead of fail-fast

**Independent Test**: Submit record with multiple errors and verify all are returned in
aggregate mode

### Tests for User Story 5

- [x] T040 [P] [US5] Test aggregate mode returns all errors in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T041 [P] [US5] Test fail-fast mode returns first error only in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T042 [P] [US5] Test aggregate mode with valid record returns empty slice in sdk/go/pluginsdk/focus_conformance_test.go

### Implementation for User Story 5

- [x] T043 [US5] Implement aggregate error collection in ValidateFocusRecordWithOptions in sdk/go/pluginsdk/focus_conformance.go
- [x] T044 [US5] Ensure fail-fast mode still works as before in sdk/go/pluginsdk/focus_conformance.go
- [x] T045 [US5] Run tests and verify all US5 tests pass with go test ./sdk/go/pluginsdk/...

**Checkpoint**: Both validation modes work correctly

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, performance, and integration

- [x] T046 [P] Add benchmark tests for zero-allocation validation in sdk/go/pluginsdk/focus_benchmark_test.go
- [x] T047 [P] Verify SC-001: <100ns, 0 allocs on valid records with go test -bench=. ./sdk/go/pluginsdk/...
- [x] T048 [P] Update sdk/go/testing/README.md with new validation capabilities
- [x] T049 [P] Add contextual validation conformance tests in sdk/go/testing/focus13_conformance_test.go
- [x] T050 Run full validation suite: make lint && make test
- [x] T051 Verify backward compatibility: existing focus_conformance_test.go tests still pass
- [x] T052 Run quickstart.md examples manually to verify documentation accuracy

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phases 3-7)**: All depend on Foundational completion
  - US1 and US2 are both P1 - can run in parallel
  - US3 and US4 are both P2 - can run in parallel after P1
  - US5 is P3 - run after P2
- **Polish (Phase 8)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Foundational only - independent
- **User Story 2 (P1)**: Foundational only - independent, can parallel with US1
- **User Story 3 (P2)**: Foundational only - independent
- **User Story 4 (P2)**: Foundational only - independent, can parallel with US3
- **User Story 5 (P3)**: Depends on Foundational + aggregation logic affects all validators

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Helper function before integration into validateBusinessRules
- Run tests after each implementation step

### Parallel Opportunities

- T002, T003: Create skeleton files in parallel
- T006, T009: Test files can be written in parallel
- All US1 tests (T013-T018) can run in parallel
- All US2 tests (T022-T025) can run in parallel
- US1 and US2 implementation can proceed in parallel
- All Phase 8 tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for US1 together:
Task: "Test EffectiveCost > BilledCost returns error in focus_conformance_test.go"
Task: "Test ListCost < EffectiveCost returns error in focus_conformance_test.go"
Task: "Test valid cost hierarchy passes in focus_conformance_test.go"
Task: "Test zero costs (free tier) pass validation in focus_conformance_test.go"
Task: "Test negative costs (credits) are exempt in focus_conformance_test.go"
Task: "Test ChargeClass CORRECTION exempts cost hierarchy in focus_conformance_test.go"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Cost Hierarchy)
4. Complete Phase 4: User Story 2 (Commitment Discounts)
5. **STOP and VALIDATE**: Both P1 stories are independently testable
6. Run `make lint && make test` to verify

### Incremental Delivery

1. Setup + Foundational → Core types ready
2. US1 → Cost hierarchy validation works → Test independently
3. US2 → Commitment consistency works → Test independently
4. US3 + US4 → P2 features work → Test independently
5. US5 → Aggregate mode works → Full feature complete
6. Polish → Documentation, benchmarks, conformance tests

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational done:
   - Developer A: User Story 1 (Cost Hierarchy)
   - Developer B: User Story 2 (Commitment Discounts)
3. After P1 complete:
   - Developer A: User Story 3 (Pricing)
   - Developer B: User Story 4 (Capacity Reservation)
4. Any developer: User Story 5 (Aggregation)
5. Team: Polish phase in parallel

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All validation functions must achieve zero-allocation on happy path (SC-001)
