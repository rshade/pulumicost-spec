# Tasks: GetBudgets RPC for Plugin-Provided Budget Information

**Input**: Design documents from `/specs/001-get-budgets-rpc/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Following test-first protocol per constitution, basic conformance tests are included for gRPC behavior validation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto definitions**: `proto/pulumicost/v1/`
- **Generated SDK**: `sdk/go/proto/`
- **Plugin SDK**: `sdk/go/pluginsdk/`
- **Testing**: `sdk/go/testing/`
- **Examples**: `examples/`
- **Schemas**: `schemas/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Proto definitions and basic project structure for budget functionality

- [x] T001 Create budget.proto with core message definitions in proto/pulumicost/v1/budget.proto
- [x] T002 Add GetBudgets RPC to CostSource service in proto/pulumicost/v1/costsource.proto
- [x] T003 [P] Validate proto syntax with buf lint
- [x] T004 [P] Create budget_spec.schema.json for JSON validation in schemas/budget_spec.schema.json

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: SDK generation and testing infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Generate Go SDK from proto definitions with make generate
- [x] T006 Add BudgetsProvider interface to pluginsdk in sdk/go/pluginsdk/sdk.go
- [x] T007 Implement optional RPC detection in Server wrapper for GetBudgets
- [x] T008 Add ValidateBudgetsResponse to test harness in sdk/go/testing/harness.go
- [x] T009 Create basic conformance test structure for GetBudgets RPC (empty, single,
  multiple budgets, error conditions)
- [x] T010 [P] Update mock plugin with GetBudgets implementation in sdk/go/testing/mock_plugin.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Unified Budget Visibility Across Providers (Priority: P1) 🎯 MVP

**Goal**: Enable plugins to provide budget information from cloud providers (AWS, GCP, Azure) in a unified format

**Independent Test**: Configure AWS and GCP plugins, verify budget data appears in unified
output with correct provider attribution

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T011 [P] [US1] Basic GetBudgets RPC conformance test in sdk/go/testing/
- [x] T012 [P] [US1] Cross-provider budget mapping validation test

### Implementation for User Story 1

- [ ] T013 [US1] Implement AWS budget mapping in aws-cost-plugin repository
  - Issue: [rshade/pulumicost-plugin-aws-ce#4](https://github.com/rshade/pulumicost-plugin-aws-ce/issues/4)
- [ ] T014 [US1] Implement GCP budget mapping in gcp-cost-plugin repository
- [ ] T015 [US1] Implement Azure budget mapping in azure-cost-plugin repository
- [ ] T016 [US1] Add provider filtering support in GetBudgets RPC
  - Issue: [rshade/pulumicost-core#263](https://github.com/rshade/pulumicost-core/issues/263)
- [ ] T017 [US1] Add currency handling for multi-provider budgets
  - Issue: [rshade/pulumicost-core#263](https://github.com/rshade/pulumicost-core/issues/263)
- [ ] T018 [US1] Update budget summary calculation for cross-provider aggregation
  - Issue: [rshade/pulumicost-core#263](https://github.com/rshade/pulumicost-core/issues/263)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Kubernetes Budget Tracking (Priority: P2)

**Goal**: Enable Kubecost plugin to expose namespace-level budget information

**Independent Test**: Configure Kubecost plugin, verify namespace budgets appear with proper filtering and status

### Tests for User Story 2 ⚠️

- [ ] T019 [P] [US2] Kubecost budget mapping validation test
  - Issue: [rshade/pulumicost-core#264](https://github.com/rshade/pulumicost-core/issues/264)
- [ ] T020 [P] [US2] Namespace-level budget filtering test
  - Issue: [rshade/pulumicost-core#264](https://github.com/rshade/pulumicost-core/issues/264)

### Implementation for User Story 2

- [ ] T021 [US2] Implement Kubecost budget mapping in kubecost-cost-plugin repository
  - Issue: [rshade/pulumicost-plugin-kubecost#38](https://github.com/rshade/pulumicost-plugin-kubecost/issues/38)
- [ ] T022 [US2] Add namespace filtering support to BudgetFilter
  - Issue: [rshade/pulumicost-core#266](https://github.com/rshade/pulumicost-core/issues/266)
- [ ] T023 [US2] Implement namespace-specific budget status tracking
  - Issue: [rshade/pulumicost-core#266](https://github.com/rshade/pulumicost-core/issues/266)
- [ ] T024 [US2] Add Kubecost-specific metadata handling
  - Issue: [rshade/pulumicost-core#266](https://github.com/rshade/pulumicost-core/issues/266)
- [ ] T025 [US2] Update budget summary for namespace aggregation
  - Issue: [rshade/pulumicost-core#266](https://github.com/rshade/pulumicost-core/issues/266)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Multi-Cloud Budget Health Overview (Priority: P3)

**Goal**: Provide aggregated budget health status across all configured providers

**Independent Test**: Configure multiple providers, verify aggregated summary shows correct health breakdown

### Tests for User Story 3 ⚠️

- [ ] T026 [P] [US3] Budget health aggregation test
  - Issue: [rshade/pulumicost-core#265](https://github.com/rshade/pulumicost-core/issues/265)
- [ ] T027 [P] [US3] Multi-provider summary calculation test
  - Issue: [rshade/pulumicost-core#265](https://github.com/rshade/pulumicost-core/issues/265)

### Implementation for User Story 3

- [ ] T028 [US3] Implement budget health status calculation logic
  - Issue: [rshade/pulumicost-core#267](https://github.com/rshade/pulumicost-core/issues/267)
- [ ] T029 [US3] Add threshold-based alerting for budget warnings
  - Issue: [rshade/pulumicost-core#267](https://github.com/rshade/pulumicost-core/issues/267)
- [ ] T030 [US3] Implement forecasted spending calculations
  - Issue: [rshade/pulumicost-core#267](https://github.com/rshade/pulumicost-core/issues/267)
- [ ] T031 [US3] Add budget health aggregation across providers
  - Issue: [rshade/pulumicost-core#267](https://github.com/rshade/pulumicost-core/issues/267)
- [ ] T032 [US3] Update summary statistics for health status breakdown
  - Issue: [rshade/pulumicost-core#267](https://github.com/rshade/pulumicost-core/issues/267)

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, and validation improvements

- [x] T033 [P] Create AWS budget example in examples/budgets/aws-budget.json
- [x] T034 [P] Create GCP budget example in examples/budgets/gcp-budget.json
- [x] T035 [P] Create Kubecost budget example in examples/budgets/kubecost-budget.json
- [x] T036 [P] Add GetBudgets request examples in examples/requests/
- [x] T037 Update PLUGIN_DEVELOPER_GUIDE.md with budget implementation guidance
- [x] T038 Add budget RPC documentation to proto comments
- [x] T039 Run full validation suite with make validate
- [x] T040 Performance test with 100-1000 budget scale (in sdk/go/testing/benchmark_test.go)
- [x] T041 Update CHANGELOG.md (managed by release-please)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of other stories
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent of other stories

### Within Each User Story

- Tests MUST be written and FAIL before implementation (test-first protocol)
- Proto/SDK changes before plugin implementations
- Core functionality before advanced features
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Provider-specific implementations can run in parallel
- All Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Basic GetBudgets RPC conformance test in sdk/go/testing/"
Task: "Cross-provider budget mapping validation test"

# Launch provider implementations in parallel:
Task: "Implement AWS budget mapping in plugin (provider-specific)"
Task: "Implement GCP budget mapping in plugin (provider-specific)"
Task: "Implement Azure budget mapping in plugin (provider-specific)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (proto definitions)
2. Complete Phase 2: Foundational (SDK generation, interfaces)
3. Complete Phase 3: User Story 1 (AWS/GCP/Azure unified visibility)
4. **STOP and VALIDATE**: Test cross-provider budget visibility independently
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
   - Developer A: User Story 1 (cross-provider visibility)
   - Developer B: User Story 2 (Kubernetes integration)
   - Developer C: User Story 3 (health aggregation)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (test-first protocol)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Follow proto-first development: proto changes → SDK generation → implementation
