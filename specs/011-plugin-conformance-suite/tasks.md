# Tasks: Plugin Conformance Test Suite

**Input**: Design documents from `/specs/011-plugin-conformance-suite/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: This feature IS a testing framework, so tests validate the suite itself. Tests are
included to validate suite behavior.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Package**: `sdk/go/testing/` (extending existing package)
- Tests: `sdk/go/testing/*_test.go`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Core types and structures shared by all user stories

- [ ] T001 Define ConformanceLevel enum type in sdk/go/testing/conformance.go
- [ ] T002 Define TestCategory enum type in sdk/go/testing/conformance.go
- [ ] T003 [P] Define SuiteConfig struct in sdk/go/testing/conformance.go
- [ ] T004 [P] Define ValidationError struct in sdk/go/testing/conformance.go
- [ ] T005 Extend TestResult struct with Category field in sdk/go/testing/harness.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core suite infrastructure that MUST be complete before ANY user story

**CRITICAL**: No user story work can begin until this phase is complete

- [ ] T006 Define ConformanceSuite struct in sdk/go/testing/conformance.go
- [ ] T007 Implement NewConformanceSuite() constructor in sdk/go/testing/conformance.go
- [ ] T008 Implement NewConformanceSuiteWithConfig() in sdk/go/testing/conformance.go
- [ ] T009 [P] Define CategoryResult struct in sdk/go/testing/conformance.go
- [ ] T010 [P] Define ResultSummary struct in sdk/go/testing/conformance.go
- [ ] T011 Define ConformanceResult struct in sdk/go/testing/conformance.go
- [ ] T012 Implement ConformanceResult.ToJSON() method in sdk/go/testing/report.go
- [ ] T013 Implement ConformanceResult.Passed() method in sdk/go/testing/conformance.go
- [ ] T014 [P] Define PerformanceBaseline struct in sdk/go/testing/performance.go
- [ ] T015 [P] Implement DefaultBaselines() function in sdk/go/testing/performance.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Validate Plugin Specification Compliance (Priority: P1)

**Goal**: Run automated tests that verify `GetPricingSpec` responses are schema-compliant

**Independent Test**: Import conformance library, call spec validation, receive pass/fail with
specific error messages

### Tests for User Story 1

- [ ] T016 [P] [US1] Test spec validation passes for valid PricingSpec in sdk/go/testing/spec_validation_test.go
- [ ] T017 [P] [US1] Test spec validation fails for invalid billing mode in sdk/go/testing/spec_validation_test.go
- [ ] T018 [P] [US1] Test spec validation fails for missing required fields in sdk/go/testing/spec_validation_test.go

### Implementation for User Story 1

- [ ] T019 [US1] Create validatePricingSpecSchema() in sdk/go/testing/spec_validation.go
- [ ] T020 [US1] Create validateBillingModeEnum() in sdk/go/testing/spec_validation.go
- [ ] T021 [US1] Create validateRequiredFields() in sdk/go/testing/spec_validation.go
- [ ] T022 [US1] Implement ValidationError formatting with field-level details in sdk/go/testing/spec_validation.go
- [ ] T023 [US1] Create SpecValidationTest conformance test in sdk/go/testing/spec_validation.go
- [ ] T024 [US1] Register spec validation tests in ConformanceSuite for Basic level in sdk/go/testing/conformance.go
- [ ] T025 [US1] Implement RunSpecValidation() convenience function in sdk/go/testing/spec_validation.go

**Checkpoint**: User Story 1 complete - spec validation works independently

---

## Phase 4: User Story 2 - Verify RPC Method Correctness (Priority: P1)

**Goal**: Automated tests exercise all RPC methods with valid/invalid inputs

**Independent Test**: Run RPC correctness suite, receive detailed results for each method

### Tests for User Story 2

- [ ] T026 [P] [US2] Test RPC correctness for valid Name request in sdk/go/testing/rpc_correctness_test.go
- [ ] T027 [P] [US2] Test RPC correctness for valid Supports request in sdk/go/testing/rpc_correctness_test.go
- [ ] T028 [P] [US2] Test RPC returns error for nil resource descriptor in sdk/go/testing/rpc_correctness_test.go
- [ ] T029 [P] [US2] Test RPC returns InvalidArgument for invalid time range in sdk/go/testing/rpc_correctness_test.go

### Implementation for User Story 2

- [ ] T030 [US2] Create testNameRPC() function in sdk/go/testing/rpc_correctness.go
- [ ] T031 [US2] Create testSupportsRPC() function in sdk/go/testing/rpc_correctness.go
- [ ] T032 [US2] Create testGetActualCostRPC() function in sdk/go/testing/rpc_correctness.go
- [ ] T033 [US2] Create testGetProjectedCostRPC() function in sdk/go/testing/rpc_correctness.go
- [ ] T034 [US2] Create testGetPricingSpecRPC() function in sdk/go/testing/rpc_correctness.go
- [ ] T035 [US2] Create testNilResourceHandling() function in sdk/go/testing/rpc_correctness.go
- [ ] T036 [US2] Create testInvalidTimeRangeHandling() function in sdk/go/testing/rpc_correctness.go
- [ ] T037 [US2] Register RPC correctness tests in ConformanceSuite for Basic level in sdk/go/testing/conformance.go
- [ ] T038 [US2] Implement RunRPCCorrectness() convenience function in sdk/go/testing/rpc_correctness.go

**Checkpoint**: User Story 2 complete - RPC correctness works independently

---

## Phase 5: User Story 3 - Measure Plugin Performance (Priority: P2)

**Goal**: Standardized benchmarks measuring latency and memory allocations

**Independent Test**: Run benchmark suite, receive min/avg/max latency and allocation counts

### Tests for User Story 3

- [ ] T039 [P] [US3] Test benchmark returns latency metrics in sdk/go/testing/performance_test.go
- [ ] T040 [P] [US3] Test benchmark compares against baseline thresholds with <10% variance (SC-003) in sdk/go/testing/performance_test.go
- [ ] T041 [P] [US3] Test benchmark warns on excessive allocations in sdk/go/testing/performance_test.go

### Implementation for User Story 3

- [ ] T042 [US3] Create measureLatency() function in sdk/go/testing/performance.go
- [ ] T043 [US3] Create measureAllocations() function in sdk/go/testing/performance.go
- [ ] T044 [US3] Create compareToBaseline() function in sdk/go/testing/performance.go
- [ ] T045 [US3] Create PerformanceResult struct in sdk/go/testing/performance.go
- [ ] T046 [US3] Create performanceTests for each RPC method in sdk/go/testing/performance.go
- [ ] T047 [US3] Register performance tests in ConformanceSuite for Standard level in sdk/go/testing/conformance.go
- [ ] T048 [US3] Implement RunPerformanceBenchmarks() convenience function in sdk/go/testing/performance.go

**Checkpoint**: User Story 3 complete - performance benchmarks work independently

---

## Phase 6: User Story 4 - Detect Concurrency Issues (Priority: P2)

**Goal**: Tests exercise plugin under concurrent load, detect race conditions

**Independent Test**: Run concurrency suite with `-race` flag, receive pass/fail with stack traces

### Tests for User Story 4

- [ ] T049 [P] [US4] Test concurrent requests complete successfully in sdk/go/testing/concurrency_test.go
- [ ] T050 [P] [US4] Test race detection integration works in sdk/go/testing/concurrency_test.go
- [ ] T051 [P] [US4] Test response consistency under load in sdk/go/testing/concurrency_test.go

### Implementation for User Story 4

- [ ] T052 [US4] Create ConcurrencyConfig struct in sdk/go/testing/concurrency.go
- [ ] T053 [US4] Create runParallelRequests() function in sdk/go/testing/concurrency.go
- [ ] T054 [US4] Create validateConsistentResponses() function in sdk/go/testing/concurrency.go
- [ ] T055 [US4] Create concurrencyTest for Standard level (10 requests) in sdk/go/testing/concurrency.go
- [ ] T056 [US4] Create concurrencyTest for Advanced level (50 requests) in sdk/go/testing/concurrency.go
- [ ] T057 [US4] Register concurrency tests in ConformanceSuite for Standard/Advanced levels in sdk/go/testing/conformance.go
- [ ] T058 [US4] Implement RunConcurrencyTests() convenience function in sdk/go/testing/concurrency.go

**Checkpoint**: User Story 4 complete - concurrency tests work independently

---

## Phase 7: User Story 5 - Run Complete Conformance Suite (Priority: P3)

**Goal**: Single command runs full suite, returns comprehensive report

**Independent Test**: Import library, run full suite, receive JSON report with certification level

### Tests for User Story 5

- [ ] T059 [P] [US5] Test full suite returns consolidated report in sdk/go/testing/conformance_test.go
- [ ] T060 [P] [US5] Test suite determines correct certification level in sdk/go/testing/conformance_test.go
- [ ] T061 [P] [US5] Test suite provides actionable failure feedback in sdk/go/testing/conformance_test.go

### Implementation for User Story 5

- [ ] T062 [US5] Implement ConformanceSuite.Run() method in sdk/go/testing/conformance.go
- [ ] T063 [US5] Implement ConformanceSuite.RunCategory() method in sdk/go/testing/conformance.go
- [ ] T064 [US5] Implement determineLevelAchieved() logic in sdk/go/testing/conformance.go
- [ ] T065 [US5] Implement aggregateResults() function in sdk/go/testing/conformance.go
- [ ] T066 [US5] Implement PrintReport() function in sdk/go/testing/report.go
- [ ] T067 [US5] Implement RunBasicConformance() convenience function in sdk/go/testing/conformance.go
- [ ] T068 [US5] Implement RunStandardConformance() convenience function in sdk/go/testing/conformance.go
- [ ] T069 [US5] Implement RunAdvancedConformance() convenience function in sdk/go/testing/conformance.go

**Checkpoint**: User Story 5 complete - full suite works with JSON output

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, cleanup, and final validation

- [ ] T070 [P] Update sdk/go/testing/README.md with conformance suite documentation
- [ ] T071 [P] Add conformance examples to sdk/go/testing/README.md
- [ ] T072 Run make lint to verify code style in sdk/go/testing/
- [ ] T073 Run make test to verify all tests pass
- [ ] T074 [P] Verify JSON report schema matches contracts/conformance_api.go
- [ ] T075 Run quickstart.md validation against MockPlugin

### Edge Case Coverage (from spec.md)

- [ ] T076 [P] Test suite handles nil plugin implementation gracefully in sdk/go/testing/conformance_test.go
- [ ] T077 [P] Test suite recovers from plugin panics in sdk/go/testing/rpc_correctness_test.go
- [ ] T078 [P] Test suite enforces timeout for slow plugin responses in sdk/go/testing/concurrency_test.go
- [ ] T079 [P] Test suite handles unimplemented plugin (empty responses) in sdk/go/testing/conformance_test.go

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - US1 and US2 are both P1 - can run in parallel
  - US3 and US4 are both P2 - can run in parallel (after P1s or independently)
  - US5 depends on US1-US4 (aggregates all categories)
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 3 (P2)**: Can start after Foundational - Independent of US1/US2
- **User Story 4 (P2)**: Can start after Foundational - Independent of US1/US2/US3
- **User Story 5 (P3)**: Depends on US1-US4 (integrates all test categories)

### Within Each User Story

- Tests written first, must FAIL before implementation
- Validation functions before test registration
- Individual tests before convenience functions
- Story complete before moving to next priority

### Parallel Opportunities

- T003, T004 can run in parallel (different structs)
- T009, T010 can run in parallel (different structs)
- T014, T015 can run in parallel (different functions)
- All [P] tests within a story can run in parallel
- US1 and US2 can run in parallel (both P1, different files)
- US3 and US4 can run in parallel (both P2, different files)

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "T016 Test spec validation passes for valid PricingSpec"
Task: "T017 Test spec validation fails for invalid billing mode"
Task: "T018 Test spec validation fails for missing required fields"

# After tests written, implement in order:
Task: "T019 Create validatePricingSpecSchema()"
Task: "T020 Create validateBillingModeEnum()"
# ... continue sequentially
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Spec Validation)
4. **STOP and VALIDATE**: Test spec validation independently
5. Plugin developers can use spec validation immediately

### Incremental Delivery

1. Setup + Foundational → Framework ready
2. Add US1 (Spec Validation) → MVP! Plugin devs can validate specs
3. Add US2 (RPC Correctness) → Expanded coverage
4. Add US3 (Performance) → Production readiness metrics
5. Add US4 (Concurrency) → Thread safety validation
6. Add US5 (Full Suite) → Complete certification experience

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (spec_validation.go)
   - Developer B: User Story 2 (rpc_correctness.go)
3. Then:
   - Developer A: User Story 3 (performance.go)
   - Developer B: User Story 4 (concurrency.go)
4. Finally: Either developer completes User Story 5 (full suite integration)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Run `make lint` and `make test` after each phase

## Implicitly Covered Requirements

The following requirements are satisfied by the package structure and test-first design rather than
explicit tasks:

- **FR-011** (importable library): Satisfied by implementing in `sdk/go/testing/` package with
  public exports
- **FR-015** (idempotent): Satisfied by test-first approach - tests verify consistent results
  across multiple runs
