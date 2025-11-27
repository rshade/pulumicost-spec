# Tasks: Trace ID Validation for TracingUnaryServerInterceptor

**Input**: Design documents from `/specs/008-trace-id-validation/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Test-first development required per constitution. Tests MUST be written before implementation and
verified to fail.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Go SDK**: `sdk/go/` at repository root
- **Primary location**: `sdk/go/pluginsdk/` package
- **Dependencies**: `sdk/go/pricing/` for validation reuse

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Verify Go 1.24+ toolchain and dependencies in go.mod
- [x] T002 [P] Confirm existing ValidateTraceID function in sdk/go/pricing/observability_validate.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T003 Create GenerateTraceID helper function in sdk/go/pluginsdk/traceid.go
- [x] T004 [P] Add crypto/rand import to pluginsdk package
- [x] T005 Verify pricing package import direction (pluginsdk → pricing is safe)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Secure Trace ID Processing (Priority: P1) 🎯 MVP

**Goal**: Validate incoming trace_id values against established format and replace invalid ones with generated IDs

**Independent Test**: Send gRPC requests with various malformed trace_id values and verify they are replaced
with valid generated IDs

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T006 [P] [US1] Add table-driven test cases for invalid trace IDs in sdk/go/pluginsdk/logging_test.go
- [x] T007 [P] [US1] Add test cases for valid trace IDs (should be preserved) in sdk/go/pluginsdk/logging_test.go
- [x] T008 [P] [US1] Add test cases for edge cases (control chars, excessive length) in sdk/go/pluginsdk/logging_test.go
- [x] T009 [US1] Verify all new tests FAIL against current implementation

### Implementation for User Story 1

- [x] T010 [US1] Modify TracingUnaryServerInterceptor in sdk/go/pluginsdk/logging.go to validate incoming trace IDs
- [x] T011 [US1] Add logic to generate replacement trace IDs when validation fails
- [x] T012 [US1] Update interceptor to store validated/generated trace ID in context
- [x] T013 [US1] Verify all User Story 1 tests now PASS

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Automatic Trace ID Generation (Priority: P2)

**Goal**: Automatically generate valid trace IDs when none is provided or when provided one is invalid

**Independent Test**: Send requests without trace_id metadata and verify a valid trace_id is generated and added to context

### Tests for User Story 2 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T014 [P] [US2] Add test cases for missing trace_id metadata in sdk/go/pluginsdk/logging_test.go
- [x] T015 [P] [US2] Add test cases for empty trace_id values in sdk/go/pluginsdk/logging_test.go
- [x] T016 [US2] Verify generation tests FAIL against current implementation
- [x] T017 [US2] Ensure interceptor generates trace ID when metadata is missing
- [x] T018 [US2] Ensure interceptor generates trace ID when trace_id value is empty
- [x] T019 [US2] Verify all User Story 2 tests now PASS

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Backward Compatible Integration (Priority: P3)

**Goal**: Enable validation by default without breaking existing plugin code

**Independent Test**: Upgrade SDK in existing plugin code and verify all existing tests pass without modification

### Tests for User Story 3 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T020 [P] [US3] Add backward compatibility tests in sdk/go/pluginsdk/logging_test.go
- [x] T021 [P] [US3] Test existing plugin scenarios with valid trace IDs
- [x] T022 [US3] Verify backward compatibility tests PASS (no breaking changes)
- [x] T023 [US3] Confirm no API signature changes to TracingUnaryServerInterceptor
- [x] T024 [US3] Verify existing plugins continue working unchanged
- [x] T025 [US3] Document backward compatibility guarantees

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T026 [P] Add benchmark tests for validation performance in sdk/go/pluginsdk/logging_test.go
- [x] T027 [P] Update package documentation in sdk/go/pluginsdk/logging.go
- [x] T028 [P] Add usage examples to quickstart.md
- [x] T029 Run make validate to ensure all tests pass
- [x] T030 [P] Update CHANGELOG.md with feature details

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Builds on US1 validation logic
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Tests backward compatibility

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Helper functions before interceptor modifications
- Validation logic before generation logic
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Add table-driven test cases for invalid trace IDs in sdk/go/pluginsdk/logging_test.go"
Task: "Add test cases for valid trace IDs (should be preserved) in sdk/go/pluginsdk/logging_test.go"
Task: "Add test cases for edge cases (control chars, excessive length) in sdk/go/pluginsdk/logging_test.go"

# Then verify they fail:
Task: "Verify all new tests FAIL against current implementation"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (validation logic)
   - Developer B: User Story 2 (generation logic)
   - Developer C: User Story 3 (compatibility testing)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
