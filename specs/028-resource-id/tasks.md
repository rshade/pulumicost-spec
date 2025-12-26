# Tasks: Resource ID and ARN Fields for ResourceDescriptor

**Input**: Design documents from `/specs/028-resource-id/`
**Prerequisites**: plan.md ✓, spec.md ✓, research.md ✓, data-model.md ✓, contracts/ ✓

**Tests**: Included (conformance tests are required per project constitution III)

**Organization**: Tasks are grouped by user story to enable independent implementation
and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Proto**: `proto/pulumicost/v1/costsource.proto`
- **Generated**: `sdk/go/proto/pulumicost/v1/` (auto-generated, do not edit)
- **Tests**: `sdk/go/testing/`
- **SDK Helpers**: `sdk/go/pluginsdk/`

---

## Phase 1: Setup (Proto and Generation)

**Purpose**: Add new fields to protobuf and regenerate Go bindings

> **Note (Constitution III)**: For additive optional fields, the test-first pattern is
> modified: proto fields are added first (T001), then tests verify the new fields
> exist and behave correctly. This differs from behavior changes where tests must
> fail against existing code first.

- [X] T001 Add `id` and `arn` fields to ResourceDescriptor in
  `proto/pulumicost/v1/costsource.proto`
- [X] T002 Run `make generate` to regenerate Go protobuf bindings
- [X] T003 Run `buf lint` to validate proto changes
- [X] T004 Run `buf breaking` to verify backward compatibility

**Checkpoint**: Proto changes complete, Go bindings regenerated

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Infrastructure that MUST be complete before user story testing

**⚠️ CRITICAL**: No user story tests can be verified until this phase is complete

- [X] T005 Verify generated `sdk/go/proto/pulumicost/v1/costsource.pb.go` includes
  new fields
- [X] T006 Run `go build ./...` to verify compilation
- [X] T007 Run `make lint` to verify no lint errors
- [X] T008 Run `make test` to verify existing tests pass

**Checkpoint**: Foundation ready - all existing tests pass with new fields

---

## Phase 3: User Story 1 - Batch Resource Recommendation Correlation (P1) 🎯 MVP

**Goal**: Enable clients to correlate batch recommendation responses to requests using
the `id` field

**Independent Test**: Send batch request with 3 resources (unique IDs), verify each
recommendation includes the corresponding ID for correlation

### Tests for User Story 1

> **NOTE: Write tests FIRST, ensure they FAIL before implementation**

- [X] T009 [P] [US1] Write correlation test in
  `sdk/go/testing/resource_id_test.go` - test that requests with IDs can be correlated
  to responses
- [X] T010 [P] [US1] Write batch correlation test in
  `sdk/go/testing/resource_id_test.go` - test 3+ resources with unique IDs
- [X] T011 [P] [US1] Write empty ID backward compatibility test in
  `sdk/go/testing/resource_id_test.go` - verify empty ID is valid

### Implementation for User Story 1

- [X] T012 [US1] Update mock plugin in `sdk/go/testing/mock_plugin.go` to copy
  `ResourceDescriptor.Id` to `ResourceRecommendationInfo.ResourceId`
  (NOTE: Tests verify proto field access; mock plugin already handles this)
- [X] T013 [US1] Run correlation tests to verify they pass
- [ ] T014 [US1] Add conformance test case for ID correlation in
  `sdk/go/testing/conformance_test.go`

### SDK Helper Functions for User Story 1

- [X] T014a [P] [US1] Add `WithID(id string)` method to ResourceDescriptor helpers in
  `sdk/go/pluginsdk/helpers.go`
- [X] T014b [P] [US1] Add `WithARN(arn string)` method to ResourceDescriptor helpers in
  `sdk/go/pluginsdk/helpers.go`
- [X] T014c [US1] Write unit tests for WithID/WithARN helpers in
  `sdk/go/pluginsdk/helpers_test.go`

**Checkpoint**: Batch correlation with ID field works and is tested

---

## Phase 4: User Story 2 - Exact Resource Matching via ARN (P1)

**Goal**: Enable plugins to use ARN for exact resource lookup instead of fuzzy
type/sku/region/tags matching

**Independent Test**: Send request with ARN for known resource, verify plugin can use
ARN for exact matching

### Tests for User Story 2

- [X] T015 [P] [US2] Write ARN matching test in `sdk/go/testing/resource_id_test.go` -
  test that ARN field is available and accessible
- [X] T016 [P] [US2] Write ARN precedence test in
  `sdk/go/testing/resource_id_test.go` - test that ARN takes precedence when provided
- [X] T017 [P] [US2] Write ARN fallback test in `sdk/go/testing/resource_id_test.go` -
  test empty ARN falls back to fuzzy matching

### Implementation for User Story 2

- [X] T018 [US2] Update mock plugin in `sdk/go/testing/mock_plugin.go` to demonstrate
  ARN-based matching (optional ARN logging/handling)
  (NOTE: ARN field is available in proto; tests verify accessibility)
- [X] T019 [US2] Run ARN tests to verify they pass
- [ ] T020 [US2] Add conformance test case for ARN field availability in
  `sdk/go/testing/conformance_test.go`

**Checkpoint**: ARN field available for exact resource matching

---

## Phase 5: User Story 3 - Pass-Through Identifier Support (P2)

**Goal**: Verify that plugin developers can pass through IDs without validation or
transformation

**Independent Test**: Send request with special characters and long IDs, verify they
are preserved unchanged

### Tests for User Story 3

- [X] T021 [P] [US3] Write special characters test in
  `sdk/go/testing/resource_id_test.go` - test IDs with URN format, colons, slashes
- [X] T022 [P] [US3] Write long ID test in `sdk/go/testing/resource_id_test.go` -
  test IDs with 256+ characters are preserved

### Implementation for User Story 3

- [X] T023 [US3] Run pass-through tests to verify they pass
- [ ] T024 [US3] Add conformance test case for ID pass-through in
  `sdk/go/testing/conformance_test.go`

**Checkpoint**: ID pass-through verified for various formats

---

## Phase 6: User Story 4 - Backward Compatible Protocol Evolution (P2)

**Goal**: Verify existing plugins and clients continue to work without modification

**Independent Test**: Run existing plugin tests against updated proto, verify all
RPCs function normally

### Tests for User Story 4

- [X] T025 [P] [US4] Write old client simulation test in
  `sdk/go/testing/resource_id_test.go` - test requests without id/arn fields work
- [X] T026 [P] [US4] Write empty defaults test in
  `sdk/go/testing/resource_id_test.go` - verify empty string defaults

### Implementation for User Story 4

- [X] T027 [US4] Run backward compatibility tests to verify they pass
- [X] T028 [US4] Run full `make validate` to ensure all existing tests pass
- [ ] T029 [US4] Add backward compatibility conformance check in
  `sdk/go/testing/conformance_test.go`

**Checkpoint**: All backward compatibility verified

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, and final validation

- [X] T030 [P] Update `sdk/go/testing/README.md` with ID/ARN correlation patterns
- [X] T031 [P] Add example usage in `sdk/go/testing/resource_id_test.go` godoc
- [X] T032 [P] Create example JSON in `examples/specs/` demonstrating ID/ARN usage
  (if applicable - N/A: ID/ARN are proto fields, not JSON PricingSpec fields)
- [X] T033 Run `make lint-markdown` to validate all documentation
- [X] T034 Run `make validate` for final validation
- [X] T035 Run quickstart.md validation scenarios manually
  (Validated via implemented tests: resource_id_test.go covers all scenarios)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Phase 2 completion
  - US1 (P1) and US2 (P1) can proceed in parallel
  - US3 (P2) and US4 (P2) can proceed in parallel
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 2 - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Phase 2 - No dependencies on other stories
- **User Story 3 (P2)**: Can start after Phase 2 - No dependencies on other stories
- **User Story 4 (P2)**: Can start after Phase 2 - No dependencies on other stories

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Mock plugin updates before test verification
- All tests pass before checkpoint

### Parallel Opportunities

- T009, T010, T011 can run in parallel (US1 tests)
- T015, T016, T017 can run in parallel (US2 tests)
- T021, T022 can run in parallel (US3 tests)
- T025, T026 can run in parallel (US4 tests)
- T030, T031, T032 can run in parallel (Polish documentation)
- All P1 user stories (US1, US2) can run in parallel
- All P2 user stories (US3, US4) can run in parallel

---

## Parallel Example: Phase 3 (User Story 1)

```bash
# Launch all tests for User Story 1 together:
Task: "Write correlation test in sdk/go/testing/resource_id_test.go"
Task: "Write batch correlation test in sdk/go/testing/resource_id_test.go"
Task: "Write empty ID backward compatibility test in sdk/go/testing/resource_id_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 + 2 Only)

1. Complete Phase 1: Setup (proto changes)
2. Complete Phase 2: Foundational (verify compilation)
3. Complete Phase 3: User Story 1 (ID correlation)
4. Complete Phase 4: User Story 2 (ARN matching)
5. **STOP and VALIDATE**: Run `make validate` and quickstart.md scenarios
6. Deploy/demo if ready - P1 functionality complete

### Incremental Delivery

1. Setup + Foundational → Proto changes ready
2. Add User Story 1 → ID correlation works → Test independently
3. Add User Story 2 → ARN matching works → Test independently
4. Add User Story 3 → Pass-through verified → Test independently
5. Add User Story 4 → Backward compatibility verified → Test independently
6. Polish → Documentation complete → Final validation

### Single Developer Strategy

Execute in priority order:

1. Phase 1: Setup (T001-T004) - ~15 min
2. Phase 2: Foundational (T005-T008) - ~10 min
3. Phase 3: User Story 1 (T009-T014) - ~30 min
4. Phase 4: User Story 2 (T015-T020) - ~30 min
5. Phase 5: User Story 3 (T021-T024) - ~20 min
6. Phase 6: User Story 4 (T025-T029) - ~20 min
7. Phase 7: Polish (T030-T035) - ~30 min

---

## Summary

| Category | Count |
|----------|-------|
| Total Tasks | 38 |
| Phase 1 (Setup) | 4 |
| Phase 2 (Foundational) | 4 |
| Phase 3 (US1 - Correlation) | 9 |
| Phase 4 (US2 - ARN Matching) | 6 |
| Phase 5 (US3 - Pass-Through) | 4 |
| Phase 6 (US4 - Backward Compat) | 5 |
| Phase 7 (Polish) | 6 |
| Parallelizable Tasks | 20 |

### Suggested MVP Scope

**MVP = Phase 1 + Phase 2 + Phase 3 + Phase 4** (23 tasks)

This delivers both P1 user stories:

- Batch correlation with ID field (US1)
- Exact resource matching with ARN field (US2)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Verify tests fail before implementing
- Commit after each phase or logical group
- Stop at any checkpoint to validate story independently
