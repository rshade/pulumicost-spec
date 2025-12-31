# Tasks: Plugin Capability Dry Run Mode

**Input**: Design documents from `/specs/032-plugin-dry-run/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Included per constitution requirement (Test-First Protocol is NON-NEGOTIABLE).

**Organization**: Tasks are grouped by user story to enable independent implementation and
testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Based on plan.md project structure:

- **Proto definitions**: `proto/pulumicost/v1/`
- **Generated code**: `sdk/go/proto/pulumicost/v1/`
- **SDK helpers**: `sdk/go/pluginsdk/`
- **Testing framework**: `sdk/go/testing/`
- **Examples**: `examples/requests/dry_run/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create examples directory structure at examples/requests/dry_run/
- [x] T002 [P] Verify current field numbers in proto/pulumicost/v1/costsource.proto for
      safe field additions

---

## Phase 2: Foundational (Test Stubs - Expected to FAIL)

**Purpose**: Write conformance tests that define expected behavior BEFORE proto changes

**CRITICAL**: Per constitution III (Test-First Protocol), these tests MUST be written
first and MUST FAIL because the proto messages don't exist yet. This defines the contract.

### Test Stubs (Will Fail - Proto Not Yet Defined)

- [x] T003 [P] Write conformance test stub for DryRun RPC basic functionality in
      sdk/go/testing/dry_run_conformance_test.go - define expected request/response contract
      (EXPECTED TO FAIL: DryRunRequest/Response types don't exist yet)
- [x] T004 [P] Write conformance test stub for unsupported resource behavior in
      sdk/go/testing/dry_run_conformance_test.go - define resource_type_supported=false
      (EXPECTED TO FAIL: DryRunResponse type doesn't exist yet)
- [x] T005 [P] Write conformance test stub for FieldSupportStatus enum usage in
      sdk/go/testing/dry_run_conformance_test.go - define expected enum values
      (EXPECTED TO FAIL: FieldSupportStatus enum doesn't exist yet)
- [x] T006 [P] Write unit test stub for FocusFieldNames helper in
      sdk/go/pluginsdk/dry_run_test.go - define expected field list
      (EXPECTED TO FAIL: dry_run.go doesn't exist yet)
- [x] T007 Verify all test stubs fail with compilation errors (proto types missing)

**Checkpoint**: Tests written, all failing - ready to implement proto to make them pass

---

## Phase 3: Proto Definitions (Make Tests Compile)

**Purpose**: Add proto definitions that make the test stubs compile

### Proto Schema Changes

- [x] T008 Add FieldSupportStatus enum to proto/pulumicost/v1/enums.proto with values:
      UNSPECIFIED, SUPPORTED, UNSUPPORTED, CONDITIONAL, DYNAMIC
- [x] T009 [P] Add FieldMapping message to proto/pulumicost/v1/costsource.proto with
      fields: field_name, support_status, condition_description, expected_type
- [x] T010 [P] Add DryRunRequest message to proto/pulumicost/v1/costsource.proto with
      fields: resource (ResourceDescriptor), simulation_parameters (map)
- [x] T011 [P] Add DryRunResponse message to proto/pulumicost/v1/costsource.proto with
      fields: field_mappings, configuration_valid, configuration_errors, resource_type_supported
- [x] T012 Add DryRun RPC to CostSourceService in proto/pulumicost/v1/costsource.proto
- [x] T013 Add dry_run field to GetActualCostRequest in proto/pulumicost/v1/costsource.proto
- [x] T014 [P] Add dry_run_result field to GetActualCostResponse in
      proto/pulumicost/v1/costsource.proto
- [x] T015 [P] Add dry_run field to GetProjectedCostRequest in
      proto/pulumicost/v1/costsource.proto
- [x] T016 [P] Add dry_run_result field to GetProjectedCostResponse in
      proto/pulumicost/v1/costsource.proto
- [x] T017 Update SupportsResponse capabilities documentation in
      proto/pulumicost/v1/costsource.proto to include "dry_run" key

### Code Generation & Validation

- [x] T018 Run `make generate` to regenerate Go code from updated proto definitions
- [x] T019 Run `buf lint` to validate proto style compliance
- [x] T020 Run `buf breaking` to verify no unintended breaking changes
- [x] T021 Verify generated code compiles with `go build ./...`
- [x] T022 Verify test stubs from Phase 2 now COMPILE (types exist but tests still fail)

**Checkpoint**: Proto definitions complete, tests compile but fail - ready for implementation

---

## Phase 4: User Story 1 - Discover Plugin Field Mappings (Priority: P1)

**Goal**: Allow hosts to query a plugin for its FOCUS field mapping capabilities without
triggering actual data retrieval.

**Independent Test**: Send dry-run request for "aws:ec2:Instance" and verify response
contains expected field mappings, response time <100ms, no external API calls made.

### Implementation for User Story 1

- [x] T023 [US1] Create FocusFieldNames() helper function in sdk/go/pluginsdk/dry_run.go
      that returns all ~50 FocusCostRecord field names
- [x] T024 [US1] Create NewFieldMapping() helper in sdk/go/pluginsdk/dry_run.go for
      constructing FieldMapping proto messages
- [x] T025 [US1] Create DryRunHandler interface in sdk/go/pluginsdk/dry_run.go that
      plugins can implement for dry-run support
- [x] T026 [US1] Update MockPlugin in sdk/go/testing/mock_plugin.go to implement
      DryRun RPC with configurable field mappings
- [x] T027 [US1] Add DryRun capability to MockPlugin.Supports() response in
      sdk/go/testing/mock_plugin.go
- [x] T028 [US1] Run tests and verify T003-T006 now pass (RED -> GREEN)

**Checkpoint**: User Story 1 complete - basic dry-run field discovery works independently

---

## Phase 5: User Story 2 - Validate Plugin Configuration (Priority: P2)

**Goal**: Enable plugin configuration validation during dry-run requests, returning clear
error messages for misconfigured plugins.

**Independent Test**: Intentionally misconfigure plugin and verify dry-run returns
configuration_valid=false with descriptive error messages.

### Tests for User Story 2

- [x] T029 [P] [US2] Write conformance test for configuration validation in
      sdk/go/testing/dry_run_conformance_test.go - test valid config returns
      configuration_valid=true
- [x] T030 [P] [US2] Write conformance test for configuration error reporting in
      sdk/go/testing/dry_run_conformance_test.go - test invalid config returns errors list

### Implementation for User Story 2

- [x] T031 [US2] Add ConfigValidator interface in sdk/go/pluginsdk/dry_run.go for plugins
      to implement configuration validation
- [x] T032 [US2] Update DryRunHandler to call ConfigValidator if implemented in
      sdk/go/pluginsdk/dry_run.go
- [x] T033 [US2] Add configurable validation errors to MockPlugin in
      sdk/go/testing/mock_plugin.go
- [x] T034 [US2] Run tests and verify T029-T030 now pass

**Checkpoint**: User Story 2 complete - configuration validation works independently

---

## Phase 6: User Story 3 - Compare Plugin Capabilities (Priority: P3)

**Goal**: Enable comparison of field mappings across multiple plugins by supporting
conditional and dynamic field status indicators with descriptive explanations.

**Independent Test**: Query two mock plugins for same resource type and verify field
mappings can be compared (different CONDITIONAL/DYNAMIC statuses with descriptions).

### Tests for User Story 3

- [x] T035 [P] [US3] Write conformance test for CONDITIONAL field status in
      sdk/go/testing/dry_run_conformance_test.go - verify condition_description populated
- [x] T036 [P] [US3] Write conformance test for DYNAMIC field status in
      sdk/go/testing/dry_run_conformance_test.go - verify condition_description populated
- [x] T037 [P] [US3] Write conformance test for simulation_parameters in
      sdk/go/testing/dry_run_conformance_test.go - verify parameters affect field status

### Implementation for User Story 3

- [x] T038 [US3] Add WithCondition() builder method to FieldMapping helper in
      sdk/go/pluginsdk/dry_run.go for setting condition_description
- [x] T039 [US3] Add simulation parameter handling to DryRunHandler in
      sdk/go/pluginsdk/dry_run.go
- [x] T040 [US3] Update MockPlugin to support configurable CONDITIONAL/DYNAMIC fields in
      sdk/go/testing/mock_plugin.go
- [x] T041 [US3] Run tests and verify T035-T037 now pass

**Checkpoint**: User Story 3 complete - plugin comparison capabilities work independently

---

## Phase 7: Integration (dry_run Flag on Cost RPCs)

**Goal**: Add dry_run flag support to existing GetActualCost and GetProjectedCost RPCs

### Tests for Integration

- [x] T042 [P] Write conformance test for GetActualCost with dry_run=true in
      sdk/go/testing/dry_run_conformance_test.go - verify dry_run_result populated
      (EXPECTED TO FAIL: types don't exist yet)
- [x] T043 [P] Write conformance test for GetProjectedCost with dry_run=true in
      sdk/go/testing/dry_run_conformance_test.go - verify dry_run_result populated
      (EXPECTED TO FAIL: types don't exist yet)
- [x] T044 [P] Write conformance test for dry_run=false default behavior in
      sdk/go/testing/dry_run_conformance_test.go - verify normal cost retrieval
      (EXPECTED TO FAIL: tests not yet written)
- [x] T045 Update MockPlugin GetActualCost to check dry_run flag and return
      DryRunResponse in sdk/go/testing/mock_plugin.go
- [x] T046 Update MockPlugin GetProjectedCost to check dry_run flag and return
      DryRunResponse in sdk/go/testing/mock_plugin.go
- [x] T047 Run tests and verify T042-T044 now pass

**Checkpoint**: Integration complete - dry_run flag works on existing cost RPCs

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, benchmarks, and validation across all stories

### Documentation (Multi-Provider Examples per Constitution II)

- [x] T048 [P] Create examples/requests/dry_run/aws_ec2.json with sample DryRunRequest
      and DryRunResponse payloads for AWS EC2 instances
- [x] T049 [P] Create examples/requests/dry_run/azure_vm.json with sample DryRunRequest
      and DryRunResponse payloads for Azure Virtual Machines
- [x] T050 [P] Create examples/requests/dry_run/gcp_compute.json with sample DryRunRequest
      and DryRunResponse payloads for GCP Compute Engine instances
- [x] T051 [P] Create examples/requests/dry_run/k8s_pod.json with sample DryRunRequest
      and DryRunResponse payloads for Kubernetes Pod resources
- [x] T052 [P] Create examples/requests/dry_run/README.md documenting example usage
      and cross-provider patterns
- [x] T053 [P] Update sdk/go/pluginsdk/README.md with DryRun implementation guide
- [x] T054 Update CLAUDE.md with new patterns and learnings from this feature

### Performance Benchmarks (per Constitution VI)

- [x] T055 Add BenchmarkDryRun function to sdk/go/testing/benchmark_test.go measuring
      DryRun RPC latency across 1000 iterations
- [x] T056 Run performance benchmark with `go test -bench=BenchmarkDryRun -benchmem
./sdk/go/testing/` - verify <100ms p99 latency and document results

### Validation

- [x] T057 Run full test suite with `make test`
- [x] T058 Run full lint suite with `make lint`
- [x] T059 Verify all examples pass schema validation
- [x] T060 Run quickstart.md validation scenarios manually

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Test Stubs (Phase 2)**: Depends on Setup - write tests that define expected behavior
- **Proto Definitions (Phase 3)**: Depends on Test Stubs - make tests compile
- **User Stories (Phase 4-6)**: All depend on Proto phase completion
  - US1, US2, US3 can proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 then P2 then P3)
- **Integration (Phase 7)**: Can proceed in parallel with US2/US3 after US1 complete
- **Polish (Phase 8)**: Depends on all user stories being complete

### Test-First Flow (Constitution III Compliance)

```text
Phase 2: Write test stubs → Tests FAIL (types don't exist)
Phase 3: Add proto definitions → Tests COMPILE but still FAIL (no implementation)
Phase 4: Implement US1 → Tests PASS (RED → GREEN → REFACTOR)
```

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Proto (Phase 3) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Proto (Phase 3) - Builds on US1 patterns
  but independently testable
- **User Story 3 (P3)**: Can start after Proto (Phase 3) - Builds on US1 patterns
  but independently testable

### Within Each User Story

- Tests written in Phase 2 define expected behavior
- Implementation makes tests pass (RED -> GREEN)
- Story complete before moving to next priority

### Parallel Opportunities

- T003-T006 (test stubs) can run in parallel
- T009-T011 (message definitions) can run in parallel
- T013-T016 (field additions) can run in parallel
- All test tasks for US2/US3 (T029-T030, T035-T037, T042-T044) can be written in parallel
- Documentation tasks (T048-T054) can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: Test-First Flow

```bash
# Phase 2: Launch all test stubs in parallel (will fail - no types yet):
Task: "Write conformance test stub for DryRun RPC basic functionality"
Task: "Write conformance test stub for unsupported resource behavior"
Task: "Write conformance test stub for FieldSupportStatus enum usage"
Task: "Write unit test stub for FocusFieldNames helper"

# Verify tests fail (expected):
Task: "Verify all test stubs fail with compilation errors"

# Phase 3: Add proto definitions to make tests compile:
Task: "Add FieldSupportStatus enum"
Task: "Add FieldMapping message"
Task: "Add DryRunRequest message"
Task: "Add DryRunResponse message"

# Phase 4: Implement to make tests pass:
Task: "Create FocusFieldNames() helper"
Task: "Create NewFieldMapping() helper"
Task: "Run tests and verify T003-T006 now pass"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Test Stubs (tests will FAIL - this is expected)
3. Complete Phase 3: Proto Definitions (tests will COMPILE but still fail)
4. Complete Phase 4: User Story 1 Implementation (tests PASS)
5. **STOP and VALIDATE**: Test DryRun RPC independently
6. Can demo basic field mapping discovery

### Incremental Delivery

1. Complete Setup + Test Stubs + Proto -> Foundation ready, tests defined
2. Add User Story 1 -> Tests pass -> Can query field mappings (MVP!)
3. Add User Story 2 -> Tests pass -> Configuration validation works
4. Add User Story 3 -> Tests pass -> Comparison capabilities work
5. Add Integration -> dry_run flag on cost RPCs works
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Test Stubs + Proto together
2. Once Proto is done:
   - Developer A: User Story 1 (core functionality)
   - Developer B: User Story 2 (validation) after US1 patterns established
   - Developer C: Integration phase after US1 complete
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **Test-First is NON-NEGOTIABLE**: Tests written in Phase 2 must fail initially
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Run `make generate` after proto changes before implementing SDK code
- Run `buf lint` and `buf breaking` to catch proto issues early
- Multi-provider examples (AWS, Azure, GCP, K8s) required per constitution II
- Performance benchmarks required per constitution VI

---

## Completion Summary

### Phase 7 - Integration (dry_run Flag on Cost RPCs): COMPLETE

- T042-T044: Conformance tests for GetActualCost/GetProjectedCost dry_run flag
- T045-T046: MockPlugin updated to check dry_run flag and return DryRunResponse
- T047: All integration tests passing

### Phase 8 - Polish & Cross-Cutting Concerns: COMPLETE

- T048-T052: Multi-provider examples created (AWS EC2, Azure VM, GCP Compute, K8s Pod)
- T053: SDK README updated with DryRun implementation guide
- T055-T056: BenchmarkDryRun functions added and validated (<100ms p99)
- T057: Test suite passing
- T058: Linting passing
- T059: Schema validation passing for examples

### All Tasks Complete

All Phase 8 tasks are now complete:

- T053: SDK README updated with comprehensive DryRun implementation guide
- T054: CLAUDE.md updated with DryRun patterns and helper documentation
- T055-T056: BenchmarkDryRun functions added and validated (<100ms p99)
- T060: Quickstart scenarios validated (all DryRun tests pass)
