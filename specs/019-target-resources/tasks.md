# Tasks: Target Resources for Recommendations

**Input**: Design documents from `/specs/019-target-resources/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Included per Constitution Test-First Protocol (Section III)

**Organization**: Tasks grouped by user story for independent implementation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto**: `proto/pulumicost/v1/costsource.proto`
- **SDK Testing**: `sdk/go/testing/`
- **Generated Code**: `sdk/go/proto/pulumicost/v1/` (via `make generate`)

---

## Phase 1: Setup (Proto Definition)

**Purpose**: Add target_resources field to proto and regenerate SDK

- [x] T001 Add target_resources field to GetRecommendationsRequest in proto/pulumicost/v1/costsource.proto
- [x] T002 Run `make generate` to regenerate Go SDK code in sdk/go/proto/pulumicost/v1/
- [x] T003 Run `make lint` to validate proto with buf lint and Go linting

---

## Phase 2: Foundational (SDK Validation Infrastructure)

**Purpose**: Core validation constants and functions that ALL user stories depend on

**CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Add MaxTargetResources constant (100) in sdk/go/testing/contract.go
- [x] T005 Add ErrTargetResourcesExceedsLimit error variable in sdk/go/testing/contract.go
- [x] T006 Add unit tests for ValidateTargetResources in sdk/go/testing/contract_test.go
- [x] T007 Implement ValidateTargetResources function in sdk/go/testing/contract.go

**Checkpoint**: Validation infrastructure ready - user story implementation can begin

---

## Phase 3: User Story 1 - Stack-Scoped Recommendations (Priority: P1)

**Goal**: Filter recommendations to only return those matching specified target resources

**Independent Test**: Pass list of ResourceDescriptors to GetRecommendations, verify only
matching recommendations returned

### Tests for User Story 1

> **NOTE**: Write these tests FIRST, ensure they FAIL before implementation

- [x] T008 [P] [US1] Add test TestTargetResourcesFiltering_SingleResource in sdk/go/testing/integration_test.go
- [x] T009 [P] [US1] Add test TestTargetResourcesFiltering_MultipleResources in sdk/go/testing/integration_test.go
- [x] T010 [P] [US1] Add test TestTargetResourcesFiltering_EmptyPreservesExisting in sdk/go/testing/integration_test.go

### Implementation for User Story 1

- [x] T011 [US1] Add filterByTargetResources helper function in sdk/go/testing/mock_plugin.go
- [x] T012 [US1] Add matchesResourceDescriptor helper function in sdk/go/testing/mock_plugin.go
- [x] T013 [US1] Update GetRecommendations in MockPlugin to apply target_resources filtering in sdk/go/testing/mock_plugin.go
- [x] T014 [US1] Add integration test for backward compatibility (empty target_resources) in sdk/go/testing/integration_test.go
- [x] T015 [US1] Run `make test` to verify all US1 tests pass

**Checkpoint**: User Story 1 complete - stack-scoped filtering works independently

---

## Phase 4: User Story 2 - Pre-Deployment Cost Optimization (Priority: P2)

**Goal**: Support recommendations for proposed resources (SKU + region matching)

**Independent Test**: Pass proposed resource configs with SKU/region, verify SKU-specific
recommendations returned

### Tests for User Story 2

- [x] T016 [P] [US2] Add test TestTargetResourcesFiltering_WithSKU in sdk/go/testing/integration_test.go
- [x] T017 [P] [US2] Add test TestTargetResourcesFiltering_WithRegion in sdk/go/testing/integration_test.go
- [x] T018 [P] [US2] Add test TestTargetResourcesFiltering_MultiProvider in sdk/go/testing/integration_test.go

### Implementation for User Story 2

- [x] T019 [US2] Enhance matchesResourceDescriptor to support strict SKU matching in sdk/go/testing/mock_plugin.go
- [x] T020 [US2] Enhance matchesResourceDescriptor to support strict region matching in sdk/go/testing/mock_plugin.go
- [x] T021 [US2] Run `make test` to verify all US2 tests pass

**Checkpoint**: User Story 2 complete - SKU/region-specific filtering works

---

## Phase 5: User Story 3 - Batch Resource Analysis (Priority: P3)

**Goal**: Support batch queries with AND logic between target_resources and filter

**Independent Test**: Pass target_resources with RecommendationFilter, verify AND logic applied

### Tests for User Story 3

- [x] T022 [P] [US3] Add test TestTargetResourcesFiltering_WithTags in sdk/go/testing/integration_test.go
- [x] T023 [P] [US3] Add test TestTargetResourcesFiltering_ANDLogicWithFilter in sdk/go/testing/integration_test.go
- [x] T024 [P] [US3] Add test TestTargetResourcesFiltering_LargeList in sdk/go/testing/integration_test.go

### Implementation for User Story 3

- [x] T025 [US3] Enhance matchesResourceDescriptor to support tag subset matching in sdk/go/testing/mock_plugin.go
- [x] T026 [US3] Verify GetRecommendations applies filter after target_resources in sdk/go/testing/mock_plugin.go
- [x] T027 [US3] Run `make test` to verify all US3 tests pass

**Checkpoint**: User Story 3 complete - batch filtering with AND logic works

---

## Phase 6: Edge Cases & Error Handling

**Purpose**: Validate error conditions and edge cases from spec

- [x] T028 [P] Add test TestTargetResourcesValidation_ExceedsLimit (covered by TestValidateTargetResources)
- [x] T029 [P] Add test TestTargetResourcesValidation_InvalidResourceDescriptor (covered by TestValidateTargetResources)
- [x] T030 [P] Add test TestTargetResourcesFiltering_NoMatchReturnsEmpty in sdk/go/testing/integration_test.go
- [x] T031 [P] Add test TestTargetResourcesFiltering_DuplicatesHandled in sdk/go/testing/integration_test.go

**Checkpoint**: All edge cases covered

---

## Phase 7: Polish & Documentation

**Purpose**: Final validation and documentation updates

- [x] T032 Run `make lint` to verify code quality (Go linting passes)
- [x] T033 Run `make test` to verify all tests pass
- [x] T034 Run `make validate` for full validation pipeline (Go + tests pass, markdown errors in ISSUE_DRAFT.md only)
- [x] T035 [P] Update sdk/go/testing/README.md with target_resources examples
- [x] T036 Validate quickstart.md examples compile correctly
- [x] T037 Verify SDK helpers are usable by running quickstart.md examples against mock plugin

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 - BLOCKS all user stories
- **User Stories (Phases 3-5)**: All depend on Phase 2 completion
  - US1, US2, US3 can proceed in parallel after Phase 2
  - Or sequentially: P1 → P2 → P3
- **Edge Cases (Phase 6)**: Can run after Phase 2, parallel with user stories
- **Polish (Phase 7)**: Depends on all phases complete

### User Story Dependencies

- **User Story 1 (P1)**: After Phase 2 - No dependencies on other stories
- **User Story 2 (P2)**: After Phase 2 - Builds on US1 implementation but independently testable
- **User Story 3 (P3)**: After Phase 2 - Builds on US1/US2 but independently testable

### Within Each User Story

1. Tests MUST be written and FAIL before implementation
2. Helper functions before main implementation
3. Integration after helpers complete
4. Story checkpoint validates independence

### Parallel Opportunities

**Phase 2 (can run together)**:

- T004, T005 (constants and errors)
- T006 must complete before T007 (test-first per Constitution III)

**Each User Story (tests can run together)**:

- US1: T008, T009, T010
- US2: T016, T017, T018
- US3: T022, T023, T024

**Phase 6 (all can run together)**:

- T028, T029, T030, T031

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "T008 [P] [US1] Add test TestTargetResourcesFiltering_SingleResource"
Task: "T009 [P] [US1] Add test TestTargetResourcesFiltering_MultipleResources"
Task: "T010 [P] [US1] Add test TestTargetResourcesFiltering_EmptyPreservesExisting"

# Then implement sequentially:
Task: "T011 [US1] Add filterByTargetResources helper function"
Task: "T012 [US1] Add matchesResourceDescriptor helper function"
Task: "T013 [US1] Update GetRecommendations in MockPlugin"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (proto + generate)
2. Complete Phase 2: Foundational (validation)
3. Complete Phase 3: User Story 1 (basic filtering)
4. **STOP and VALIDATE**: Test stack-scoped filtering independently
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy (MVP!)
3. Add User Story 2 → Test independently → Deploy (SKU/region support)
4. Add User Story 3 → Test independently → Deploy (batch + AND logic)
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 + Edge Cases
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files or independent test functions
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Constitution compliance: Test-First Protocol (Section III) enforced
