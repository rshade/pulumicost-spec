# Tasks: Configurable CORS Headers

**Input**: Design documents from `/specs/033-cors-headers-config/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Included per Constitution III (Test-First is NON-NEGOTIABLE)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **SDK Package**: `sdk/go/pluginsdk/`
- **Tests**: Same package with `_test.go` suffix

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Define default header constants used by all user stories

- [x] T001 Define DefaultAllowedHeaders constant in sdk/go/pluginsdk/options.go
- [x] T002 Define DefaultExposedHeaders constant in sdk/go/pluginsdk/options.go

---

## Phase 2: User Story 1 - Custom Allowed Headers (Priority: P1)

**Goal**: Enable plugin developers to customize Access-Control-Allow-Headers

**Independent Test**: Configure custom AllowedHeaders, send preflight OPTIONS, verify response

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T003 [P] [US1] Add test for nil AllowedHeaders (uses defaults) in sdk/go/pluginsdk/sdk_test.go
- [x] T004 [P] [US1] Add test for custom AllowedHeaders in sdk/go/pluginsdk/sdk_test.go
- [x] T005 [P] [US1] Add test for empty AllowedHeaders slice (FR-008) in sdk/go/pluginsdk/sdk_test.go

### Implementation for User Story 1

- [x] T006 [US1] Add AllowedHeaders field to WebConfig struct in sdk/go/pluginsdk/options.go
- [x] T007 [US1] Update corsMiddleware to use AllowedHeaders (nil → default, empty → empty, custom → join) in sdk/go/pluginsdk/sdk.go
- [x] T008 [US1] Verify tests T003-T005 pass

**Checkpoint**: Custom AllowedHeaders is functional and tested independently ✅

---

## Phase 3: User Story 2 - Custom Exposed Headers (Priority: P2)

**Goal**: Enable plugin developers to customize Access-Control-Expose-Headers

**Independent Test**: Configure custom ExposedHeaders, send CORS request, verify response

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T009 [P] [US2] Add test for nil ExposedHeaders (uses defaults) in sdk/go/pluginsdk/sdk_test.go
- [x] T010 [P] [US2] Add test for custom ExposedHeaders in sdk/go/pluginsdk/sdk_test.go
- [x] T011 [P] [US2] Add test for empty ExposedHeaders slice (FR-009) in sdk/go/pluginsdk/sdk_test.go

### Implementation for User Story 2

- [x] T012 [US2] Add ExposedHeaders field to WebConfig struct in sdk/go/pluginsdk/options.go
- [x] T013 [US2] Update corsMiddleware to use ExposedHeaders (nil → default, empty → empty, custom → join) in sdk/go/pluginsdk/sdk.go
- [x] T014 [US2] Verify tests T009-T011 pass

**Checkpoint**: Custom ExposedHeaders is functional and tested independently ✅

---

## Phase 4: User Story 3 - Builder Method Configuration (Priority: P3)

**Goal**: Provide fluent builder methods consistent with existing WebConfig API

**Independent Test**: Use WithAllowedHeaders/WithExposedHeaders, verify config values

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T015 [P] [US3] Add test for WithAllowedHeaders builder method in sdk/go/pluginsdk/options_test.go
- [x] T016 [P] [US3] Add test for WithExposedHeaders builder method in sdk/go/pluginsdk/options_test.go
- [x] T017 [P] [US3] Add test for defensive slice copying in builder methods in sdk/go/pluginsdk/options_test.go
- [x] T018 [P] [US3] Add test for builder method chaining in sdk/go/pluginsdk/options_test.go

### Implementation for User Story 3

- [x] T019 [US3] Implement WithAllowedHeaders builder method in sdk/go/pluginsdk/options.go
- [x] T020 [US3] Implement WithExposedHeaders builder method in sdk/go/pluginsdk/options.go
- [x] T021 [US3] Verify tests T015-T018 pass

**Checkpoint**: Builder methods are functional and tested independently ✅

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, performance verification, and final validation

- [x] T022 [P] Add benchmark test for corsMiddleware overhead in sdk/go/pluginsdk/sdk_test.go
- [x] T023 [P] Update WebConfig struct godoc comments with default header values in sdk/go/pluginsdk/options.go
- [x] T024 [P] Update sdk/go/pluginsdk/README.md with CORS header configuration examples
- [x] T025 Run make lint and fix any issues
- [x] T026 Run make test and verify all tests pass
- [x] T027 Verify benchmark shows <1μs overhead per request (SC-005)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - constants needed by US1 and US2
- **User Story 1 (Phase 2)**: Depends on Setup (T001)
- **User Story 2 (Phase 3)**: Depends on Setup (T002), can run parallel to US1
- **User Story 3 (Phase 4)**: Depends on US1 and US2 (fields must exist for builder methods)
- **Polish (Phase 5)**: Depends on all user stories complete

### User Story Dependencies

```text
Setup (T001-T002)
    ├── US1: AllowedHeaders (T003-T008) ─┐
    │                                    │
    └── US2: ExposedHeaders (T009-T014) ─┼── US3: Builder Methods (T015-T021)
                                         │
                                         └── Polish (T022-T027)
```

### Within Each User Story

- Tests MUST be written and FAIL before implementation (Constitution III)
- Implementation after tests
- Verify tests pass after implementation

### Parallel Opportunities

**Phase 2 (US1) Tests**:

- T003, T004, T005 can run in parallel (different test functions)

**Phase 3 (US2) Tests**:

- T009, T010, T011 can run in parallel (different test functions)

**Phase 4 (US3) Tests**:

- T015, T016, T017, T018 can run in parallel (different test functions)

**Phase 5 (Polish)**:

- T022, T023, T024 can run in parallel (different files)

**Cross-Story Parallelism**:

- US1 and US2 can be implemented in parallel (different fields, different middleware sections)

---

## Parallel Example: User Story 1

```bash
# Launch all US1 tests together (TDD - must fail first):
Task: "T003 [P] [US1] Add test for nil AllowedHeaders"
Task: "T004 [P] [US1] Add test for custom AllowedHeaders"
Task: "T005 [P] [US1] Add test for empty AllowedHeaders slice"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: User Story 1 (T003-T008)
3. **STOP and VALIDATE**: Test AllowedHeaders independently
4. Deploy/demo if ready (partial feature)

### Incremental Delivery

1. Setup → Constants ready
2. US1 → AllowedHeaders works → Test independently
3. US2 → ExposedHeaders works → Test independently
4. US3 → Builder methods work → Full API complete
5. Polish → Docs, benchmarks, final validation

### Parallel Team Strategy

With two developers:

1. Developer A: Setup → US1
2. Developer B: (wait for Setup) → US2
3. Both: US3 (simple, can split builder methods)
4. Both: Polish tasks in parallel

---

## Notes

- [P] tasks = different files or test functions, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Constitution III: Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
