# Tasks: "What-If" Cost Estimation API

**Input**: Design documents from `/specs/006-estimate-cost/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: This feature follows Test-First Protocol (Constitution Principle III). Conformance tests
are REQUIRED per the TDD approach.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and verify existing structure

- [x] T001 Verify Go 1.24+ toolchain and Protocol Buffers v3 installation
- [x] T002 Verify buf v1.32.1 is installed in bin/ (run `make generate` if needed)
- [x] T003 [P] Verify zerolog v1.34.0+ dependency in go.mod
- [x] T004 [P] Review existing proto/pulumicost/v1/costsource.proto to identify cost field type
      pattern (examine GetActualCost and GetProjectedCost response messages - likely double, string,
      or custom decimal type - document findings in comment for T013)

---

## Phase 2: Foundational (Blocking Prerequisites) - TDD RED Phase

**Purpose**: Write conformance tests FIRST (must FAIL) before any implementation

**CRITICAL**: Tests must be written and FAILING before proto changes. This follows Constitution
Principle III (Test-First Protocol).

- [x] T005 [P] Create test helper for EstimateCost in sdk/go/testing/harness.go
      (add method to TestHarness struct)
- [x] T006 [P] Add EstimateCost mock implementation in sdk/go/testing/mock_plugin.go
      (ConfigurableEstimateCost field and handler)
- [x] T007 Create Basic conformance tests for EstimateCost in sdk/go/testing/conformance_test.go
      (TestEstimateCostBasicConformance function)
- [x] T008 [P] Create Standard conformance tests for EstimateCost in sdk/go/testing/conformance_test.go
      (TestEstimateCostStandardConformance function)
- [x] T009 [P] Create Advanced conformance tests for EstimateCost in sdk/go/testing/conformance_test.go
      (TestEstimateCostAdvancedConformance function)
- [x] T010 Add EstimateCost benchmarks in sdk/go/testing/benchmark_test.go
      (BenchmarkEstimateCost function with multiple scenarios)
- [x] T011 Run `go test -v ./sdk/go/testing/` and verify ALL EstimateCost tests FAIL (RED phase complete)

**Checkpoint**: All tests written and failing - ready for protocol implementation (GREEN phase)

---

## Phase 3: User Story 1 - Basic Cost Estimation (Priority: P1) 🎯 MVP

**Goal**: Enable developers to get cost estimates for resources before deployment by implementing
the core EstimateCost RPC with request/response messages.

**Independent Test**: Send EstimateCost request with valid resource type and attributes, verify
cost estimate is returned. Tests: TestEstimateCostBasicConformance must pass.

### Protocol Definition for User Story 1

- [x] T012 [US1] Add EstimateCostRequest message to proto/pulumicost/v1/costsource.proto
      (fields: resource_type string=1, attributes google.protobuf.Struct=2)
- [x] T013 [US1] Add EstimateCostResponse message to proto/pulumicost/v1/costsource.proto
      (fields: currency string=1, cost_monthly [USE TYPE FROM T004 - double/string/custom decimal]=2)
- [x] T014 [US1] Add EstimateCost RPC method to CostSource service in
      proto/pulumicost/v1/costsource.proto (method signature with inline documentation)
- [x] T015 [US1] Add comprehensive proto comments for EstimateCost RPC, request, and response
      messages per Constitution Principle V
- [x] T016 [US1] Run `make generate` to regenerate Go SDK from proto definitions in sdk/go/proto/

### Validation Implementation for User Story 1

- [x] T017 [US1] Implement resource type format validation function in sdk/go/pricing/validate.go
      (ValidateResourceType with regex pattern provider:module/resource:Type)
- [x] T018 [US1] Implement null/missing attributes normalization in sdk/go/pricing/validate.go
      (NormalizeAttributes function)
- [x] T019 [US1] Add unit tests for validation functions in sdk/go/pricing/validate_test.go

### Testing Framework Integration for User Story 1

- [x] T020 [US1] Update mock plugin EstimateCost handler to use generated proto messages in
      sdk/go/testing/mock_plugin.go
- [x] T021 [US1] Update TestHarness EstimateCost helper to use generated proto messages in
      sdk/go/testing/harness.go
- [x] T022 [US1] Run `go test -v ./sdk/go/testing/` and verify Basic conformance tests PASS (GREEN phase for US1)

**Checkpoint**: Basic EstimateCost RPC functional - can estimate cost for supported resource types

---

## Phase 4: User Story 2 - Configuration Comparison (Priority: P2)

**Goal**: Enable developers to compare costs between different resource configurations by ensuring
deterministic, consistent responses.

**Independent Test**: Call EstimateCost multiple times with different attribute values for same
resource type, verify different costs. Tests: TestEstimateCostStandardConformance must pass.

### Deterministic Behavior for User Story 2

- [x] T023 [P] [US2] Add determinism tests to sdk/go/testing/conformance_test.go
      (verify identical inputs produce identical outputs per FR-011)
- [x] T024 [P] [US2] Add concurrent request tests to sdk/go/testing/conformance_test.go
      (test 10+ concurrent EstimateCost calls per Standard conformance)
- [x] T025 [US2] Update mock plugin to support multiple resource types in
      sdk/go/testing/mock_plugin.go (add ResourceTypes map[string]EstimateCostConfig)
- [x] T026 [US2] Run `go test -v ./sdk/go/testing/` and verify Standard conformance tests PASS

### Cross-Provider Examples for User Story 2

- [x] T027 [P] [US2] Create examples/requests/ directory
- [x] T028 [P] [US2] Create AWS example in examples/requests/estimate_cost_aws.json
      (aws:ec2/instance:Instance with t3.micro attributes)
- [x] T029 [P] [US2] Create Azure example in examples/requests/estimate_cost_azure.json
      (azure:compute/virtualMachine:VirtualMachine with Standard_B1s)
- [x] T030 [P] [US2] Create GCP example in examples/requests/estimate_cost_gcp.json
      (gcp:compute/instance:Instance with e2-micro)
- [x] T031 [US2] Verify all example JSON files match protobuf Struct format

**Checkpoint**: Configuration comparison working - can compare costs across different configs and providers

---

## Phase 5: User Story 3 - Unsupported Resource Handling (Priority: P3)

**Goal**: Provide clear feedback for unsupported resources, invalid formats, and missing attributes
to guide developers.

**Independent Test**: Send EstimateCost requests with invalid resource types, verify appropriate
error responses. Tests: TestEstimateCostAdvancedConformance must pass.

### Error Handling for User Story 3

- [x] T032 [P] [US3] Add invalid format error test to sdk/go/testing/conformance_test.go
      (test gRPC InvalidArgument for "invalid-format" per FR-003)
- [x] T033 [P] [US3] Add unsupported resource error test to sdk/go/testing/conformance_test.go
      (test gRPC NotFound for unsupported resource per FR-008)
- [x] T034 [P] [US3] Add missing attributes error test to sdk/go/testing/conformance_test.go
      (test gRPC InvalidArgument for empty attributes per FR-009)
- [x] T035 [P] [US3] Add ambiguous attributes error test to sdk/go/testing/conformance_test.go
      (test descriptive error messages per FR-010)
- [x] T036 [P] [US3] Add pricing source unavailable error test to sdk/go/testing/conformance_test.go
      (test gRPC Unavailable per FR-014)
- [x] T037 [P] [US3] Add zero cost handling test to sdk/go/testing/conformance_test.go
      (test valid response with cost=0 per FR-013)
- [x] T038 [US3] Update mock plugin with error scenario support in sdk/go/testing/mock_plugin.go
      (add ForceError field and error injection)
- [x] T039 [US3] Run `go test -v ./sdk/go/testing/` and verify Advanced conformance tests PASS

### Observability Integration for User Story 3

- [ ] T040 [P] [US3] Add structured logging example in sdk/go/testing/integration_test.go
      (demonstrate zerolog integration per NFR-001) -
      [Issue #83](https://github.com/rshade/pulumicost-spec/issues/83)
- [ ] T041 [P] [US3] Add metrics tracking example in sdk/go/testing/integration_test.go
      (demonstrate latency/success rate tracking per NFR-002) -
      [Issue #84](https://github.com/rshade/pulumicost-spec/issues/84)
- [ ] T042 [P] [US3] Add tracing support example in sdk/go/testing/integration_test.go
      (demonstrate correlation ID handling per NFR-003) -
      [Issue #85](https://github.com/rshade/pulumicost-spec/issues/85)

**Checkpoint**: Error handling complete - all error scenarios return appropriate gRPC status codes and messages

---

## Phase 6: Performance & Benchmarking

**Purpose**: Verify performance requirements and optimize as needed

- [ ] T043 [P] Run benchmarks with `go test -bench=BenchmarkEstimateCost -benchmem ./sdk/go/testing/`
      and verify <500ms response time (SC-002)
      [Issue #86](https://github.com/rshade/pulumicost-spec/issues/86)
- [ ] T044 [P] Run concurrent benchmark with 50+ requests per Advanced conformance and verify
      <500ms under load
      [Issue #87](https://github.com/rshade/pulumicost-spec/issues/87)
- [ ] T045 Add performance regression tests to CI in .github/workflows/ (benchmark comparison between commits)
      [Issue #88](https://github.com/rshade/pulumicost-spec/issues/88)

---

## Phase 7: Documentation & Polish

**Purpose**: Complete documentation and cross-cutting concerns

- [ ] T046 [P] Update README.md with EstimateCost RPC usage and examples per Constitution Principle V
- [ ] T047 [P] Update examples/README.md with EstimateCost cross-provider coverage matrix
- [ ] T048 [P] Create migration guide in docs/ if cost field type differs from existing RPCs
- [ ] T049 [P] Update CHANGELOG.md with EstimateCost feature entry (version 0.3.0 or next)
- [ ] T050 Run `make lint` to verify all linting passes (buf lint, golangci-lint, markdownlint)
- [ ] T051 Run `make test` to verify all tests pass
- [ ] T052 Run `make validate` to verify complete validation pipeline (tests + linting + npm validations)
- [ ] T053 Verify buf breaking check passes in CI (non-breaking change per Constitution Principle IV)
- [ ] T054 Update data-model.md line 56 to replace "[decimal type TBD]" with actual type determined in T004
      [Issue #89](https://github.com/rshade/pulumicost-spec/issues/89)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup (T001-T004) - BLOCKS all user stories until tests written
- **User Story 1 (Phase 3)**: Depends on Foundational (T005-T011) - Core RPC implementation
- **User Story 2 (Phase 4)**: Depends on US1 complete (T012-T022) - Builds on basic functionality
- **User Story 3 (Phase 5)**: Depends on US1 complete (T012-T022) - Independent error handling tests
- **Performance (Phase 6)**: Depends on US1 complete (T012-T022) - Can run in parallel with US2/US3
- **Documentation (Phase 7)**: Depends on all user stories complete (T012-T042)

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational tests written - No dependencies on other stories
- **User Story 2 (P2)**: Can start after US1 basic functionality works (T012-T022) - Extends with comparison capability
- **User Story 3 (P3)**: Can start after US1 basic functionality works (T012-T022) - Independent error testing

### Within Each User Story

- Tests FIRST (Foundational Phase T005-T011) - MUST FAIL before implementation
- Protocol definition (proto messages) before code generation (T012-T016)
- Code generation before validation implementation (T017-T019)
- Validation before testing framework integration (T020-T022)
- All tasks within phase must complete before moving to next phase

### Parallel Opportunities

**Setup Phase**:

- T003 and T004 can run in parallel

**Foundational Phase**:

- T005 and T006 can run in parallel (different files)
- T008 and T009 can run in parallel (different test functions)

**User Story 1**:

- T017 and T018 can run in parallel after T016 completes (different validation functions)

**User Story 2**:

- T023 and T024 can run in parallel (different test functions)
- T028, T029, T030 can run in parallel (different example files)

**User Story 3**:

- T032-T037 can all run in parallel (different test functions in same file)
- T040, T041, T042 can run in parallel (different observability examples)

**Performance Phase**:

- T043 and T044 can run in parallel (different benchmark scenarios)

**Documentation Phase**:

- T046, T047, T048, T049 can all run in parallel (different doc files)

---

## Parallel Example: Foundational Phase (TDD RED)

```bash
# Launch all test scaffolding tasks together:
Task: "Create test helper for EstimateCost in sdk/go/testing/harness.go"
Task: "Add EstimateCost mock implementation in sdk/go/testing/mock_plugin.go"

# Launch all conformance test writing tasks together (after helpers done):
Task: "Create Basic conformance tests in sdk/go/testing/conformance_test.go"
Task: "Create Standard conformance tests in sdk/go/testing/conformance_test.go"
Task: "Create Advanced conformance tests in sdk/go/testing/conformance_test.go"
```

## Parallel Example: User Story 2

```bash
# Launch all example file creation tasks together:
Task: "Create AWS example in examples/requests/estimate_cost_aws.json"
Task: "Create Azure example in examples/requests/estimate_cost_azure.json"
Task: "Create GCP example in examples/requests/estimate_cost_gcp.json"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational - Write ALL tests FIRST (T005-T011) - RED phase
3. Complete Phase 3: User Story 1 (T012-T022) - GREEN phase for basic functionality
4. **STOP and VALIDATE**: Run `go test -v ./sdk/go/testing/` - Basic conformance must pass
5. **STOP and VALIDATE**: Run `make validate` - All linting and validation must pass
6. Deploy/demo EstimateCost RPC with basic functionality

### Incremental Delivery

1. Complete Setup + Foundational (TDD RED) → All tests written and failing
2. Add User Story 1 (TDD GREEN) → Basic functionality works → Test Basic conformance → Deploy/Demo (MVP!)
3. Add User Story 2 → Configuration comparison works → Test Standard conformance → Deploy/Demo
4. Add User Story 3 → Error handling complete → Test Advanced conformance → Deploy/Demo
5. Add Performance validation → Verify <500ms response time → Deploy/Demo
6. Each phase adds value without breaking previous functionality

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T011) - Everyone writes tests
2. Once Foundational tests are failing (RED):
   - Developer A: User Story 1 protocol + validation (T012-T022) - Get to GREEN
   - Developer B: User Story 2 examples (T027-T031) - Prepare examples while US1 implements
   - Developer C: User Story 3 error tests (T032-T039) - Add error scenario coverage
3. After US1 GREEN, US2 and US3 can integrate and test in parallel

---

## Task Summary

- **Total Tasks**: 54
- **Setup Phase**: 4 tasks
- **Foundational Phase (TDD RED)**: 7 tasks (T005-T011)
- **User Story 1 (TDD GREEN for MVP)**: 11 tasks (T012-T022)
- **User Story 2 (Configuration Comparison)**: 9 tasks (T023-T031)
- **User Story 3 (Error Handling)**: 11 tasks (T032-T042)
- **Performance Phase**: 3 tasks (T043-T045)
- **Documentation Phase**: 9 tasks (T046-T054)
- **Parallel Opportunities**: 18 tasks marked [P]

**MVP Scope**: Phases 1-3 (T001-T022) = 22 tasks to basic functionality

---

## Notes

- [P] tasks = different files, no dependencies - can run in parallel
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **CRITICAL**: Verify tests FAIL (T011) before implementing proto changes (T012-T016)
- Constitution Principle III (Test-First Protocol) enforced via Foundational phase
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Run `make validate` frequently to catch issues early
- buf breaking check must pass (non-breaking change)
