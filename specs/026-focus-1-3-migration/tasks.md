# Tasks: FOCUS 1.3 Migration

**Input**: Design documents from `/specs/026-focus-1-3-migration/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Test tasks are included as specified in FR-013 through FR-017 of the spec.

**Organization**: Tasks are grouped by user story to enable independent implementation
and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto definitions**: `proto/finfocus/v1/`
- **Go SDK**: `sdk/go/`
- **Testing framework**: `sdk/go/testing/`
- **Plugin SDK**: `sdk/go/pluginsdk/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and verify existing infrastructure

- [x] T001 Verify branch 026-focus-1-3-migration exists and is checked out
- [x] T002 [P] Verify buf v1.32.1 is installed (`make generate` installs if needed)
- [x] T003 [P] Verify Go 1.25.5+ and protobuf dependencies are available
- [x] T004 Review existing `proto/finfocus/v1/focus.proto` for field number conflicts

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core proto changes that MUST be complete before ANY user story can
be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Add `FocusContractCommitmentCategory` enum to `proto/finfocus/v1/focus.proto`
  with values: UNSPECIFIED(0), SPEND(1), USAGE(2)
- [x] T006 Add `ContractCommitment` message to `proto/finfocus/v1/focus.proto`
  with 12 fields per contracts/contract_commitment.proto
- [x] T006a Add FOCUS 1.3 specification reference comments to all new proto fields
  in `proto/finfocus/v1/focus.proto` (reference FOCUS 1.3 section numbers for
  each column: AllocatedMethodId, ServiceProviderName, etc.)
- [x] T007 Mark `provider_name` field (1) as deprecated with `[deprecated = true]`
  option in FocusCostRecord
- [x] T008 Mark `publisher` field (55) as deprecated with `[deprecated = true]`
  option in FocusCostRecord
- [x] T009 Run `make generate` to regenerate Go bindings from proto changes
- [x] T010 Verify generated code compiles: `go build ./...`
- [x] T011 Update `sdk/go/registry/` enums if ContractCommitmentCategory
  needs registry exposure

**Checkpoint**: Foundation ready - proto changes compiled, user story
implementation can now begin in parallel

---

## Phase 3: User Story 1 - Cost Allocation Columns (Priority: P1) 🎯 MVP

**Goal**: As a FinOps Practitioner, I can track split cost allocations to analyze
how shared resource costs are distributed across workloads.

**Independent Test**: Create FocusRecordBuilder with allocation fields, verify
proto serialization and validation

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T012 [P] [US1] Unit test for allocation field validation in
  `sdk/go/pluginsdk/focus_builder_test.go` - test AllocatedMethodId requires
  AllocatedResourceId
- [x] T013 [P] [US1] Unit test for allocation builder methods in
  `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T014 [P] [US1] Benchmark test for allocation operations in
  `sdk/go/testing/benchmark_test.go` - target <100ns/op, 0 allocs

### Implementation for User Story 1

- [x] T015 [US1] Add `allocated_method_id` field (61) string to FocusCostRecord
  in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T016 [P] [US1] Add `allocated_method_details` field (62) string to
  FocusCostRecord in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T017 [P] [US1] Add `allocated_resource_id` field (63) string to
  FocusCostRecord in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T018 [P] [US1] Add `allocated_resource_name` field (64) string to
  FocusCostRecord in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T019 [US1] Add `allocated_tags` field (65) map<string,string> to
  FocusCostRecord in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T020 [US1] Run `make generate` after proto field additions (completed in Phase 2)
- [x] T021 [US1] Implement `WithAllocation(methodId, methodDetails string)`
  builder method in `sdk/go/pluginsdk/focus_builder.go`
- [x] T022 [US1] Implement `WithAllocatedResource(resourceId, resourceName string)`
  builder method in `sdk/go/pluginsdk/focus_builder.go`
- [x] T023 [US1] Implement `WithAllocatedTags(tags map[string]string)`
  builder method in `sdk/go/pluginsdk/focus_builder.go`
- [x] T024 [US1] Add validation in `Build()` method: if allocated_method_id
  is set, allocated_resource_id MUST also be set
- [x] T025 [US1] Verify tests pass: `go test -v ./sdk/go/pluginsdk/...`

**Checkpoint**: User Story 1 complete - allocation fields functional and
testable independently

---

## Phase 4: User Story 2 - Service/Host Provider Columns (Priority: P1)

**Goal**: As a FinOps Practitioner, I can distinguish between service provider
and host provider to accurately attribute costs in marketplace scenarios.

**Independent Test**: Create FocusRecordBuilder with provider fields, verify
deprecation warning when both old and new fields set

### Tests for User Story 2

- [x] T026 [P] [US2] Unit test for service/host provider builder methods in
  `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T027 [P] [US2] Unit test for deprecation warning logging when provider_name
  and service_provider_name both set in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T028 [P] [US2] Benchmark test for provider operations in
  `sdk/go/testing/benchmark_test.go` - target <100ns/op, 0 allocs

### Implementation for User Story 2

- [x] T029 [US2] Add `service_provider_name` field (59) string to FocusCostRecord
  in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T030 [P] [US2] Add `host_provider_name` field (60) string to FocusCostRecord
  in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T031 [US2] Run `make generate` after proto field additions (completed in Phase 2)
- [x] T032 [US2] Implement `WithServiceProvider(name string)` builder method
  in `sdk/go/pluginsdk/focus_builder.go`
- [x] T033 [US2] Implement `WithHostProvider(name string)` builder method
  in `sdk/go/pluginsdk/focus_builder.go`
- [x] T034 [US2] Add deprecation warning logic in `Build()` method: log warning
  if provider_name is set when service_provider_name is also set, prefer
  service_provider_name value
- [x] T035 [US2] Add deprecation warning logic: log warning if publisher is set
  when host_provider_name is also set, prefer host_provider_name value
- [x] T036 [US2] Verify tests pass: `go test -v ./sdk/go/pluginsdk/...`

**Checkpoint**: User Story 2 complete - provider fields functional with
deprecation handling

---

## Phase 5: User Story 3 - Contract Commitment Dataset (Priority: P2)

**Goal**: As a Cloud Finance Manager, I can create ContractCommitment records
to track contract obligations separately from cost line items.

**Independent Test**: Create ContractCommitmentBuilder, verify all 12 fields
serialize correctly

### Tests for User Story 3

- [x] T037 [P] [US3] Unit test for ContractCommitmentBuilder in
  `sdk/go/pluginsdk/contract_commitment_builder_test.go`
- [x] T038 [P] [US3] Unit test for ContractCommitment validation (required
  fields, period consistency, non-negative values) in
  `sdk/go/pluginsdk/contract_commitment_builder_test.go`
- [x] T039 [P] [US3] Benchmark test for ContractCommitment operations in
  `sdk/go/testing/benchmark_test.go`

### Implementation for User Story 3

- [x] T040 [US3] Create `sdk/go/pluginsdk/contract_commitment_builder.go`
  with `ContractCommitmentBuilder` struct
- [x] T041 [US3] Implement `NewContractCommitmentBuilder()` constructor
- [x] T042 [US3] Implement `WithIdentity(commitmentId, contractId string)`
  builder method
- [x] T043 [US3] Implement `WithCategory(category FocusContractCommitmentCategory)`
  builder method
- [x] T044 [US3] Implement `WithType(commitmentType string)` builder method
- [x] T045 [US3] Implement `WithCommitmentPeriod(start, end time.Time)`
  builder method
- [x] T046 [US3] Implement `WithContractPeriod(start, end time.Time)`
  builder method
- [x] T047 [US3] Implement `WithFinancials(cost, quantity float64, unit,
  currency string)` builder method
- [x] T048 [US3] Implement `Build() (*pbc.ContractCommitment, error)` method
  with validation
- [x] T049 [US3] Add validation: contract_commitment_id, contract_id,
  billing_currency are required
- [x] T050 [US3] Add validation: period_end >= period_start
- [x] T051 [US3] Add validation: cost >= 0, quantity >= 0
- [x] T052 [US3] Reuse existing ISO 4217 currency validation from
  `sdk/go/pricing/` for billing_currency
- [x] T053 [US3] Verify tests pass: `go test -v ./sdk/go/pluginsdk/...`
- [x] T053a [US3] Document ContractCommitment dataset and its relationship to
  FocusCostRecord in `sdk/go/pluginsdk/README.md` (include ContractApplied
  linking pattern, 12-field overview, usage examples)

**Checkpoint**: User Story 3 complete - ContractCommitment dataset functional

---

## Phase 6: User Story 4 - Contract Applied Column (Priority: P2)

**Goal**: As a Cloud Finance Manager, I can link cost records to contract
commitments to correlate usage with contractual obligations.

**Independent Test**: Create FocusRecordBuilder with ContractApplied field,
verify it accepts any string value (opaque reference)

### Tests for User Story 4

- [x] T054 [P] [US4] Unit test for ContractApplied builder method in
  `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T055 [P] [US4] Unit test verifying no cross-dataset validation
  (ContractApplied accepts any string) in
  `sdk/go/pluginsdk/focus_builder_test.go`

### Implementation for User Story 4

- [x] T056 [US4] Add `contract_applied` field (66) string to FocusCostRecord
  in `proto/finfocus/v1/focus.proto` (completed in Phase 2)
- [x] T057 [US4] Run `make generate` after proto field addition (completed in Phase 2)
- [x] T058 [US4] Implement `WithContractApplied(commitmentId string)` builder
  method in `sdk/go/pluginsdk/focus_builder.go` (completed in Phase 3)
- [x] T059 [US4] Verify tests pass: `go test -v ./sdk/go/pluginsdk/...`

**Checkpoint**: User Story 4 complete - contract linking functional

---

## Phase 7: User Story 5 - Backward Compatibility (Priority: P3)

**Goal**: As a Plugin Developer, I can continue using existing FOCUS 1.2 code
while incrementally adopting FOCUS 1.3 features.

**Independent Test**: Build a FOCUS 1.2-only record (no new fields), verify
it passes validation and serializes correctly

### Tests for User Story 5

- [x] T060 [P] [US5] Unit test that FOCUS 1.2 records (no new fields) still
  validate successfully in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T061 [P] [US5] Conformance test for backward compatibility in
  `sdk/go/testing/focus13_conformance_test.go`
- [x] T062 [P] [US5] Test that deprecated fields still work (provider_name,
  publisher) in `sdk/go/pluginsdk/focus_builder_test.go`

### Implementation for User Story 5

- [x] T063 [US5] Verify all new FOCUS 1.3 fields have default/zero values that
  don't affect existing behavior
- [x] T064 [US5] Ensure `WithIdentity()` continues to set provider_name for
  backward compatibility (existing behavior)
- [x] T065 [US5] Document migration path in code comments for deprecated fields
- [x] T066 [US5] Add FOCUS 1.3 conformance test suite in
  `sdk/go/testing/focus13_conformance_test.go`
- [x] T067 [US5] Implement `RunFocus13ConformanceTests(t, plugin)` function
  in `sdk/go/testing/focus13_conformance_test.go` (tests directly callable)
- [x] T068 [US5] Verify tests pass: `go test -v ./sdk/go/testing/...`

**Checkpoint**: User Story 5 complete - backward compatibility verified

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T069 [P] Update `specs/026-focus-1-3-migration/quickstart.md` with
  actual API usage after implementation
- [x] T070 [P] Add FOCUS 1.3 section to `sdk/go/pluginsdk/README.md`
- [x] T077 [P] Update `docs/focus-columns.md` with FOCUS 1.3 column documentation
  (AllocatedMethodId, AllocatedMethodDetails, AllocatedResourceId, AllocatedResourceName,
  AllocatedTags, ContractApplied, ServiceProviderName, HostProviderName)
- [x] T071 [P] Update root CLAUDE.md with FOCUS 1.3 patterns and learnings
- [x] T072 Run full test suite: `make test`
- [x] T073 Run full linting: `make lint` (deprecated field warnings expected for
  backward compatibility tests)
- [x] T074 Run benchmarks and verify performance targets:
  `go test -bench=BenchmarkFocus -benchmem ./sdk/go/pluginsdk/`
  (All FOCUS 1.3 methods: < 2 ns/op, 0 allocs except tags at ~130 ns/op)
- [x] T075 Validate quickstart.md examples compile and run correctly
- [x] T076 Code review for consistency with existing SDK patterns

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies
  on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies
  on other stories, can run parallel with US1
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - No dependencies
  on US1/US2
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Conceptually
  links to US3 but implementation is independent
- **User Story 5 (P3)**: Should wait until US1-US4 complete for comprehensive
  backward compatibility testing

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Proto changes before SDK implementation
- Builder methods before validation logic
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks are sequential (proto changes must be ordered)
- Once Foundational phase completes:
  - US1 and US2 can run in parallel (both P1, different field groups)
  - US3 and US4 can run in parallel (both P2, different features)
- All tests for a user story marked [P] can run in parallel
- Proto field additions within a story marked [P] can run in parallel

---

## Implementation Strategy

### MVP First (User Stories 1 & 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Cost Allocation)
4. Complete Phase 4: User Story 2 (Service/Host Provider)
5. **STOP and VALIDATE**: Test MVP with all P1 features
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add US1 + US2 → Test independently → Deploy/Demo (MVP with P1 features!)
3. Add US3 + US4 → Test independently → Deploy/Demo (P2 features)
4. Add US5 → Test independently → Deploy/Demo (P3 backward compat)
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Allocation)
   - Developer B: User Story 2 (Provider)
3. After P1 complete:
   - Developer A: User Story 3 (ContractCommitment)
   - Developer B: User Story 4 (ContractApplied)
4. All: User Story 5 (Backward Compatibility verification)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Performance target: <100ns/op, 0 allocs for builder operations
- All new proto fields use numbers 59-66 (no conflicts with existing 1-58)
