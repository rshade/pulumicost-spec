# Tasks: GreenOps Integration

**Input**: Design documents from `specs/020-greenops-integration/`
**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [data-model.md](data-model.md), [research.md](research.md)

**Tests**: Test-First Protocol is MANDATORY per Constitution Section III.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Environment verification and branch readiness

- [X] T001 Verify `buf` CLI version and project dependencies in repository root
- [X] T002 [P] Review `proto/pulumicost/v1/costsource.proto` for insertion points
- [X] T003 [P] Ensure `make generate` environment is functional

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Protobuf definition updates and SDK regeneration

**⚠️ CRITICAL**: No user story work can begin until the protocol is updated and SDK is regenerated.

- [X] T004 Define `MetricKind` enum in `proto/pulumicost/v1/costsource.proto`
- [X] T005 Update `SupportsResponse` message with `supported_metrics` in `proto/pulumicost/v1/costsource.proto`
- [X] T006 Update `GetProjectedCostRequest` message with `utilization_percentage` in `proto/pulumicost/v1/costsource.proto`
- [X] T007 Update `ResourceDescriptor` message with `utilization_percentage` override in `proto/pulumicost/v1/costsource.proto`
- [X] T008 [P] Run `buf lint` on `proto/` directory
- [X] T009 [P] Run `buf breaking --against .git#branch=main`
- [X] T010 Regenerate Go SDK using `make generate`

**Checkpoint**: Foundation ready - Protobuf contracts are updated and Go SDK is regenerated.

---

## Phase 3: User Story 1 - Capability Discovery (Priority: P1) 🎯 MVP

**Goal**: Enable plugins to advertise supported GreenOps metrics.

**Independent Test**: Call `Supports` RPC on a plugin and verify `supported_metrics` contains the expected `MetricKind` values.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T011 [P] [US1] Create conformance test for GreenOps metrics discovery in `sdk/go/testing/greenops_discovery_test.go`
- [X] T012 [P] [US1] Add unit tests for `MetricKind` validation in `sdk/go/pluginsdk/validation_test.go`

### Implementation for User Story 1

- [X] T013 [US1] Add validation logic for `MetricKind` in `sdk/go/pluginsdk/validation.go`
- [X] T014 [US1] Update mock plugin for testing to return `supported_metrics` in `sdk/go/testing/mock_plugin.go`
- [X] T015 [US1] Verify US1 conformance tests pass

**Checkpoint**: User Story 1 is fully functional and testable independently.

---

## Phase 4: User Story 2 - Accurate Impact Modeling (Priority: P2)

**Goal**: Implement utilization precedence logic (Override > Global > Default).

**Independent Test**: Call `GetProjectedCost` with various combinations of global and per-resource utilization and verify the correct value is used.

### Tests for User Story 2

- [X] T016 [P] [US2] Create integration test for utilization precedence logic in `sdk/go/testing/utilization_test.go`
- [X] T017 [P] [US2] Add edge case tests for utilization clamping in `sdk/go/testing/utilization_clamping_test.go`
- [X] T018 [P] [US2] Add test for metric omission when data is unavailable in `sdk/go/testing/metric_omission_test.go`

### Implementation for User Story 2

- [X] T019 [US2] Implement `GetUtilization` helper in `sdk/go/pluginsdk/utilization.go` (extracts value following precedence rules)
- [X] T020 [US2] Update mock plugin to use the `GetUtilization` helper in `sdk/go/testing/mock_plugin.go`
- [X] T021 [US2] Verify US2 integration and edge case tests pass

**Checkpoint**: User Story 2 is fully functional and handles utilization precedence correctly.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final validation

- [X] T022 [P] Update GreenOps examples in `examples/plugins/greenops-plugin.json`
- [X] T023 [P] Update `examples/requests/projected-cost-with-utilization.json`
- [X] T024 [P] Add documentation comments to `proto/pulumicost/v1/costsource.proto` for units (gCO2e, kWh, L)
- [X] T025 [P] Run `make benchmarks` to ensure no regression in cost projection performance
- [X] T026 Run `quickstart.md` validation tasks

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Start immediately.
- **Foundational (Phase 2)**: Depends on T001-T003. BLOCKS all user stories.
- **User Stories (Phases 3 & 4)**: Depend on SDK regeneration (T010). Can proceed in parallel after T010.
- **Polish (Phase 5)**: Depends on completion of Phase 3 and 4.

### Parallel Opportunities

- T002 and T003 can run in parallel.
- T008 and T009 can run in parallel after T004-T007 are drafted.
- User Story 1 (Phase 3) and User Story 2 (Phase 4) can run in parallel.
- T021-T023 can run in parallel.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 & 2 (Foundation).
2. Complete Phase 3 (Capability Discovery).
3. **STOP and VALIDATE**: Verify plugins can advertise metrics.

### Incremental Delivery

1. Deliver Foundation + US1 (Discovery).
2. Deliver US2 (Utilization Logic).
3. Finalize with Polish and Examples.
