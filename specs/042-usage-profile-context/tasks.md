# Tasks: Usage Profile Context

**Input**: Design documents from `/specs/042-usage-profile-context/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Conformance tests included as required by spec (FR validation, SDK pattern conformance)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto definitions**: `proto/finfocus/v1/`
- **Go SDK**: `sdk/go/` (generated: `proto/`, helpers: `pluginsdk/`, tests: `testing/`)
- **TypeScript SDK**: `sdk/typescript/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Proto definitions and code generation - foundation for all SDK work

- [x] T001 Add UsageProfile enum to proto/finfocus/v1/enums.proto
- [x] T002 Add usage_profile field (field 6) to GetProjectedCostRequest in proto/finfocus/v1/costsource.proto
- [x] T003 Add usage_profile field (field 7) to GetRecommendationsRequest in proto/finfocus/v1/costsource.proto
- [x] T004 Regenerate Go SDK bindings with `make generate`
- [x] T005 Verify generated code compiles with `go build ./...`

**Checkpoint**: Proto definitions complete, generated code compiles

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core SDK helpers that ALL user stories depend on

**⚠️ CRITICAL**: No user story implementation can begin until this phase is complete

- [x] T006 Create sdk/go/pluginsdk/usage_profile.go with package declaration and imports
- [x] T007 Implement IsValidUsageProfile() validation function with zero-allocation pattern in sdk/go/pluginsdk/usage_profile.go
- [x] T008 Implement ParseUsageProfile() string-to-enum parser in sdk/go/pluginsdk/usage_profile.go
- [x] T009 Implement UsageProfileString() enum-to-string helper in sdk/go/pluginsdk/usage_profile.go
- [x] T010 Implement NormalizeUsageProfile() unknown value handler in sdk/go/pluginsdk/usage_profile.go
- [x] T011 [P] Create sdk/go/pluginsdk/usage_profile_test.go with table-driven tests for all helpers
- [x] T012 [P] Add benchmark tests for validation functions (target: <15 ns/op, 0 allocs) in sdk/go/pluginsdk/usage_profile_test.go

**Checkpoint**: Foundation ready - SDK helpers implemented and tested

---

## Phase 3: User Story 1 - Development Environment Cost Estimation (Priority: P1) 🎯 MVP

**Goal**: Enable Core to signal DEV intent to plugins for development workload cost estimation

**Independent Test**: Run cost estimate with `usage_profile=DEV` and verify plugins receive the profile context

### Conformance Tests for User Story 1

- [x] T013 [P] [US1] Add DEV profile conformance test case in sdk/go/testing/usage_profile_conformance_test.go
- [x] T014 [P] [US1] Add UNSPECIFIED default behavior test case in sdk/go/testing/usage_profile_conformance_test.go

### Implementation for User Story 1

- [x] T015 [US1] Update MockPlugin to handle usage_profile in GetProjectedCost in sdk/go/testing/mock_plugin.go
- [x] T016 [US1] Add profile-aware logging example to quickstart.md (DEV profile scenario)
- [x] T017 [US1] Implement WithProfileDefaults() builder method for DEV profile in sdk/go/pluginsdk/focus_builder.go

**Checkpoint**: DEV profile handling works end-to-end, UNSPECIFIED falls back to default behavior

---

## Phase 4: User Story 2 - Production Environment Cost Estimation (Priority: P1)

**Goal**: Enable Core to signal PROD intent to plugins for production workload cost estimation

**Independent Test**: Run cost estimate with `usage_profile=PROD` and verify plugins receive
the profile context and apply production assumptions

### Conformance Tests for User Story 2

- [x] T018 [P] [US2] Add PROD profile conformance test case in sdk/go/testing/usage_profile_conformance_test.go

### Implementation for User Story 2

- [x] T019 [US2] Extend WithProfileDefaults() builder method for PROD profile in sdk/go/pluginsdk/focus_builder.go
- [x] T020 [US2] Update MockPlugin to apply PROD defaults (730hr usage) in sdk/go/testing/mock_plugin.go
- [x] T021 [US2] Add PROD profile scenario to quickstart.md examples

**Checkpoint**: Both DEV and PROD profiles work correctly, with distinct behavior for each

---

## Phase 5: User Story 3 - Burst Workload Cost Estimation (Priority: P2)

**Goal**: Enable Core to signal BURST intent to plugins for temporary/high-load workload cost estimation

**Independent Test**: Run cost estimate with `usage_profile=BURST` and verify plugins receive
the profile context and apply burst assumptions

### Conformance Tests for User Story 3

- [x] T022 [P] [US3] Add BURST profile conformance test case in sdk/go/testing/usage_profile_conformance_test.go
- [x] T023 [P] [US3] Add unknown profile value forward-compatibility test in sdk/go/testing/usage_profile_conformance_test.go

### Implementation for User Story 3

- [x] T024 [US3] Extend WithProfileDefaults() builder method for BURST profile in sdk/go/pluginsdk/focus_builder.go
- [x] T025 [US3] Update MockPlugin to apply BURST defaults (high data transfer) in sdk/go/testing/mock_plugin.go
- [x] T026 [US3] Add BURST profile scenario to quickstart.md examples

**Checkpoint**: All three profiles (DEV, PROD, BURST) work correctly with distinct behavior

---

## Phase 6: GetRecommendationsRequest Integration

**Goal**: Extend profile support to recommendations endpoint

**Purpose**: Complete the usage_profile integration for recommendation generation

- [x] T027 Update MockPlugin to handle usage_profile in GetRecommendations in sdk/go/testing/mock_plugin.go
- [x] T028 [P] Add recommendations profile conformance tests in sdk/go/testing/usage_profile_conformance_test.go
- [x] T029 Add recommendations example to quickstart.md showing profile-aware priority adjustment

**Checkpoint**: Profile context available in both GetProjectedCost and GetRecommendations

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: SDK synchronization, validation, and documentation

- [x] T030 [P] Run `make lint` and fix any linting issues
- [x] T031 [P] Run `make test` and verify all tests pass
- [x] T032 [P] Verify benchmark results meet <15 ns/op target for validation functions
- [x] T033 [P] Regenerate TypeScript proto bindings: `cd sdk/typescript && npm run generate`
- [x] T034 [P] Add TypeScript helpers in sdk/typescript/src/usage-profile.ts (Constitution XIII)
- [x] T035 [P] Add TypeScript profile tests in sdk/typescript/src/usage-profile.test.ts
- [x] T036 [P] Run markdown linting with `make lint-markdown`
- [x] T037 Validate quickstart.md examples compile and run correctly
- [x] T038 Update CLAUDE.md with any new patterns discovered during implementation
- [x] T039 [P] Document DryRun profile fields in sdk/go/pluginsdk/README.md (FR-007)
- [x] T040 [P] Add DryRun profile conformance test in sdk/go/testing/ (FR-007)
- [x] T041 [P] Add backward compat conformance test for old plugins in sdk/go/testing/ (FR-005)
- [x] T042 Verify SDK ergonomics: profile handling <5 lines per SC-004

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (proto generation must complete first)
- **User Stories (Phase 3-5)**: All depend on Phase 2 (SDK helpers must exist)
  - US1 and US2 are both P1 priority and can proceed in parallel
  - US3 (P2) can proceed after Foundational, independent of US1/US2
- **GetRecommendations (Phase 6)**: Can proceed after any user story phase, independent work
- **Polish (Phase 7)**: Depends on all previous phases

### User Story Dependencies

- **User Story 1 (P1 - DEV)**: No dependencies on other stories - can test DEV profile independently
- **User Story 2 (P1 - PROD)**: No dependencies on US1 - can test PROD profile independently
- **User Story 3 (P2 - BURST)**: No dependencies on US1/US2 - can test BURST profile independently

### Within Each Phase

- Proto changes (T001-T003) MUST complete before code generation (T004)
- Code generation (T004) MUST complete before compilation check (T005)
- SDK helpers (T006-T010) can be implemented incrementally but all needed before tests
- Tests and benchmarks (T011-T012) can run in parallel
- Conformance tests per story can run in parallel within their story

### Parallel Opportunities

**Phase 1 Parallel**: T002 and T003 (different proto sections)
**Phase 2 Parallel**: T011 and T012 (tests and benchmarks)
**Phase 3-5 Parallel**: All conformance test tasks marked [P] within a phase
**Phase 6 Parallel**: T028 runs parallel with T027/T029
**Phase 7 Parallel**: All tasks marked [P] can run concurrently

---

## Parallel Example: Foundational Phase

```bash
# After T006-T010 complete, launch all tests together:
Task: "Create sdk/go/pluginsdk/usage_profile_test.go with table-driven tests"
Task: "Add benchmark tests for validation functions"
```

## Parallel Example: User Story 1

```bash
# Launch conformance tests in parallel:
Task: "Add DEV profile conformance test case"
Task: "Add UNSPECIFIED default behavior test case"
```

---

## Implementation Strategy

### MVP First (Setup + Foundational + User Story 1)

1. Complete Phase 1: Setup (proto definitions)
2. Complete Phase 2: Foundational (SDK helpers)
3. Complete Phase 3: User Story 1 (DEV profile)
4. **STOP and VALIDATE**: Test DEV profile works, UNSPECIFIED defaults work
5. Merge to main if ready

### Incremental Delivery

1. **MVP**: Setup → Foundational → US1 (DEV) → Validate
2. **+PROD**: Add User Story 2 → Test independently
3. **+BURST**: Add User Story 3 → Test independently
4. **+Recommendations**: Add Phase 6 → Test independently
5. **+TypeScript**: Add Phase 7 TypeScript tasks → Full SDK parity

### Task Count Summary

| Phase | Task Count | Story |
|-------|------------|-------|
| Phase 1: Setup | 5 | - |
| Phase 2: Foundational | 7 | - |
| Phase 3: User Story 1 | 5 | US1 (P1) |
| Phase 4: User Story 2 | 4 | US2 (P1) |
| Phase 5: User Story 3 | 5 | US3 (P2) |
| Phase 6: Recommendations | 3 | - |
| Phase 7: Polish | 13 | - |
| **Total** | **42** | |

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Performance target: <15 ns/op, 0 allocs/op for validation functions
- Commit after each task or logical group
- Stop at any checkpoint to validate independently
- TypeScript SDK update (Constitution XIII) is required but can be done after Go SDK is complete
