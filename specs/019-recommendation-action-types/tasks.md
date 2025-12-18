# Tasks: Extend RecommendationActionType Enum

**Input**: Design documents from `/specs/019-recommendation-action-types/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Constitution requires Test-First Protocol (TDD). Tests written before implementation.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This is a proto-first specification repository:

- **Proto definitions**: `proto/pulumicost/v1/`
- **Generated Go SDK**: `sdk/go/proto/pulumicost/v1/`
- **Testing framework**: `sdk/go/testing/`
- **Documentation**: `PLUGIN_DEVELOPER_GUIDE.md`, `CHANGELOG.md`

---

## Phase 1: Setup (Verification)

**Purpose**: Verify current state and prepare for changes

- [ ] T001 Verify buf CLI is available via `make generate` dry-run
- [ ] T002 Run `buf lint` to confirm current proto is valid in proto/pulumicost/v1/
- [ ] T003 Run `buf breaking` baseline to establish compatibility reference
- [ ] T004 Run existing tests via `make test` to confirm green baseline

---

## Phase 2: Foundational (Test-First Protocol)

**Purpose**: Write conformance tests that MUST FAIL before proto changes

**⚠️ CRITICAL**: Constitution III requires tests written and failing before implementation

- [ ] T005 [P] Write test for MIGRATE action type serialization in
  sdk/go/testing/action_type_test.go
- [ ] T006 [P] Write test for CONSOLIDATE action type serialization in
  sdk/go/testing/action_type_test.go
- [ ] T007 [P] Write test for SCHEDULE action type serialization in
  sdk/go/testing/action_type_test.go
- [ ] T008 [P] Write test for REFACTOR action type serialization in
  sdk/go/testing/action_type_test.go
- [ ] T009 [P] Write test for OTHER action type serialization in
  sdk/go/testing/action_type_test.go
- [ ] T010 Write test for round-trip serialization of all 12 enum values in
  sdk/go/testing/action_type_test.go
- [ ] T011 Verify tests T005-T010 FAIL (enum values 7-11 don't exist yet)

**Checkpoint**: Tests written and failing - ready for proto implementation

---

## Phase 3: User Story 1 - Plugin Developer Uses Extended Action Types (Priority: P1) 🎯 MVP

**Goal**: Enable plugins to return recommendations with new action types (MIGRATE,
CONSOLIDATE, SCHEDULE, REFACTOR, OTHER)

**Independent Test**: Mock plugin returns recommendations with new action types, SDK correctly
serializes/deserializes them

### Implementation for User Story 1

- [ ] T012 [US1] Add RECOMMENDATION_ACTION_TYPE_MIGRATE = 7 with documentation comments in
  proto/pulumicost/v1/costsource.proto (line ~621)
- [ ] T013 [US1] Add RECOMMENDATION_ACTION_TYPE_CONSOLIDATE = 8 with documentation comments in
  proto/pulumicost/v1/costsource.proto
- [ ] T014 [US1] Add RECOMMENDATION_ACTION_TYPE_SCHEDULE = 9 with documentation comments in
  proto/pulumicost/v1/costsource.proto
- [ ] T015 [US1] Add RECOMMENDATION_ACTION_TYPE_REFACTOR = 10 with documentation comments in
  proto/pulumicost/v1/costsource.proto
- [ ] T016 [US1] Add RECOMMENDATION_ACTION_TYPE_OTHER = 11 with documentation comments in
  proto/pulumicost/v1/costsource.proto
- [ ] T017 [US1] Run `buf lint` to validate proto syntax in proto/pulumicost/v1/
- [ ] T018 [US1] Run `make generate` to regenerate Go SDK in sdk/go/proto/pulumicost/v1/
- [ ] T019 [US1] Verify generated Go constants exist for all 5 new values in
  sdk/go/proto/pulumicost/v1/costsource.pb.go
- [ ] T020 [US1] Run tests T005-T010 - verify they now PASS
- [ ] T021 [US1] Update mock plugin to support new action types in sdk/go/testing/mock_plugin.go

**Checkpoint**: User Story 1 complete - plugins can use all 12 action types

---

## Phase 4: User Story 2 - Backward Compatibility with Existing Plugins (Priority: P1)

**Goal**: Ensure existing plugins continue to work without modification after SDK update

**Independent Test**: Run existing plugin implementations against updated SDK, verify all
existing functionality works unchanged

### Tests for User Story 2

- [ ] T022 [P] [US2] Write backward compatibility test for plugins using only original 6 types
  in sdk/go/testing/backward_compat_test.go
- [ ] T023 [P] [US2] Write test for unknown enum value handling (core receives value 7-11 from
  new plugin, old core behavior) in sdk/go/testing/backward_compat_test.go
- [ ] T024 [P] [US2] Write test for gRPC communication between old plugin binary and new SDK in
  sdk/go/testing/backward_compat_test.go

### Implementation for User Story 2

- [ ] T025 [US2] Run `buf breaking` against main branch to confirm no breaking changes
- [ ] T026 [US2] Verify existing conformance tests still pass with `go test ./sdk/go/testing/`
- [ ] T027 [US2] Run integration tests to verify gRPC communication in
  sdk/go/testing/integration_test.go
- [ ] T028 [US2] Document backward compatibility guarantees in PLUGIN_DEVELOPER_GUIDE.md

**Checkpoint**: User Story 2 complete - backward compatibility verified and documented

---

## Phase 5: User Story 3 - Core CLI Categorization (Priority: P2)

**Goal**: Enable CLI to filter and display recommendations by new action types

**Independent Test**: Mock plugin responses with various action types, CLI correctly groups
and labels recommendations

**Note**: CLI implementation is out of scope for this repository (handled in pulumicost-core).
This phase focuses on SDK support for CLI integration.

### Implementation for User Story 3

- [ ] T029 [P] [US3] Verify existing helpers use generated .String() method for new action types
  in sdk/go/pluginsdk/helpers.go (no changes needed - protobuf generates String() automatically)
- [ ] T030 [P] [US3] Add action type filtering examples to quickstart.md in
  specs/019-recommendation-action-types/quickstart.md
- [ ] T031 [US3] Verify RecommendationSummary.count_by_action_type handles new values correctly
  via test in sdk/go/testing/summary_test.go
- [ ] T032 [US3] Verify RecommendationFilter.action_type accepts new values via test in
  sdk/go/testing/filter_test.go

**Checkpoint**: User Story 3 complete - SDK supports CLI integration for new action types

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, validation, and release preparation

- [ ] T033 [P] Update PLUGIN_DEVELOPER_GUIDE.md with new action type usage examples
- [ ] T034 [P] Add cross-provider examples for new action types:
  - examples/recommendations/migrate_region.json (Azure/GCP region migration)
  - examples/recommendations/consolidate_nodes.json (Kubecost/Azure node consolidation)
  - examples/recommendations/schedule_dev_env.json (AWS/Azure dev environment scheduling)
  - examples/recommendations/refactor_serverless.json (GCP serverless migration)
  - examples/recommendations/other_provider_specific.json (catch-all example)
- [ ] T035 Run full validation suite with `make validate`
- [ ] T036 Run `make lint` to verify all linting passes
- [ ] T037 Verify all success criteria from spec.md are met (SC-001 through SC-005)
- [ ] T038 Run quickstart.md validation scenarios manually

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - MUST complete before implementation
- **User Story 1 (Phase 3)**: Depends on Foundational tests failing
- **User Story 2 (Phase 4)**: Depends on User Story 1 (needs proto changes to test)
- **User Story 3 (Phase 5)**: Depends on User Story 1 (needs proto changes for SDK support)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Core implementation - no dependencies on other stories
- **User Story 2 (P1)**: Tests compatibility OF User Story 1 changes - depends on US1
- **User Story 3 (P2)**: Uses SDK changes FROM User Story 1 - depends on US1

### Within Each User Story

- Proto changes before SDK regeneration
- SDK regeneration before tests can pass
- Tests passing before documentation updates

### Parallel Opportunities

**Phase 2 (Foundational Tests)**:

```text
T005, T006, T007, T008, T009 can all run in parallel (different test functions)
```

**Phase 4 (Backward Compatibility Tests)**:

```text
T022, T023, T024 can all run in parallel (different test files)
```

**Phase 5 (SDK Support)**:

```text
T029, T030 can run in parallel (T029 is verification, T030 is documentation)
```

**Phase 6 (Polish)**:

```text
T033, T034 can run in parallel (T033 updates guide, T034 creates 5 example files)
```

---

## Parallel Example: Foundational Tests

```bash
# Launch all action type tests together (Phase 2):
Task: "Write test for MIGRATE action type serialization in sdk/go/testing/action_type_test.go"
Task: "Write test for CONSOLIDATE action type serialization in sdk/go/testing/action_type_test.go"
Task: "Write test for SCHEDULE action type serialization in sdk/go/testing/action_type_test.go"
Task: "Write test for REFACTOR action type serialization in sdk/go/testing/action_type_test.go"
Task: "Write test for OTHER action type serialization in sdk/go/testing/action_type_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational tests that FAIL (T005-T011)
3. Complete Phase 3: User Story 1 proto changes (T012-T021)
4. **STOP and VALIDATE**: Tests T005-T010 should now PASS
5. MVP complete - plugins can use new action types

### Incremental Delivery

1. Setup + Foundational → Tests failing (TDD ready)
2. User Story 1 → Tests passing → Proto changes deployed
3. User Story 2 → Backward compatibility verified → Safe for existing plugins
4. User Story 3 → SDK support complete → CLI integration ready
5. Polish → Documentation and validation → Release ready

### Release Checklist

Before merging:

- [ ] All tests pass (`make test`)
- [ ] All linting passes (`make lint`)
- [ ] No breaking changes (`buf breaking`)
- [ ] Documentation updated
- [ ] Success criteria SC-001 through SC-005 verified

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Constitution requires Test-First Protocol - tests MUST fail before implementation
- Proto changes automatically regenerate Go SDK via `make generate`
- Commit after each phase completion
- Stop at any checkpoint to validate independently
