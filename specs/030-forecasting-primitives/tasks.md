# Tasks: Forecasting Primitives

**Input**: Design documents from `/specs/030-forecasting-primitives/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Constitution mandates Test-First Protocol (TDD) - conformance tests must fail before
proto implementation.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Exact file paths included in descriptions

## Path Conventions

- **Proto definitions**: `proto/finfocus/v1/`
- **Generated code**: `sdk/go/proto/` (via buf generate)
- **SDK helpers**: `sdk/go/pricing/`
- **Conformance tests**: `sdk/go/testing/`
- **Examples**: `examples/requests/`

---

## Phase 1: Setup

**Purpose**: Project initialization and proto tooling verification

- [X] T001 Verify buf v1.32.1 is installed via `make generate` in repository root
- [X] T002 [P] Run `buf lint` to confirm proto files pass linting before changes
- [X] T003 [P] Run `buf breaking` against main branch to establish baseline

**Checkpoint**: Proto tooling verified, ready for foundational work

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core proto changes that ALL user stories depend on

**CRITICAL**: No user story work can begin until GrowthType enum exists in proto

**TDD Deviation Note**: Per Constitution III (Test-First Protocol), tests should fail before
implementation. However, Phase 3 tests (T009-T011) reference the GrowthType enum which must
exist for tests to compile. This is an acceptable deviation: the enum is a foundational type
required for test compilation, not business logic. Tests will fail at runtime (missing fields)
not compile time, satisfying the TDD spirit.

- [X] T004 Add GrowthType enum to `proto/finfocus/v1/enums.proto` with values:
  GROWTH_TYPE_UNSPECIFIED (0), GROWTH_TYPE_NONE (1), GROWTH_TYPE_LINEAR (2),
  GROWTH_TYPE_EXPONENTIAL (3)
- [X] T005 Add import for enums.proto in `proto/finfocus/v1/costsource.proto` if not present
- [X] T006 Run `buf lint` to validate new enum follows style conventions
- [X] T007 Run `buf breaking` to confirm no breaking changes introduced
- [X] T008 Run `make generate` to regenerate Go SDK code in `sdk/go/proto/`

**Checkpoint**: GrowthType enum exists, SDK regenerated, user stories can begin

---

## Phase 3: User Story 1 - Project Future Costs (Priority: P1) MVP

**Goal**: Enable users to specify growth assumptions (linear/exponential) for cost projections

**Independent Test**: Provide resource descriptor with growth params, verify projected cost
response incorporates growth model

### Tests for User Story 1 (TDD - Write First, Must Fail)

- [X] T009 [P] [US1] Write conformance test for LINEAR growth in
  `sdk/go/testing/forecasting_test.go`: Given growth_type=LINEAR, rate=0.10, verify
  cost increases 10% per period linearly
- [X] T010 [P] [US1] Write conformance test for EXPONENTIAL growth in
  `sdk/go/testing/forecasting_test.go`: Given growth_type=EXPONENTIAL, rate=0.05, verify
  cost compounds at 5% per period
- [X] T011 [P] [US1] Write conformance test for NONE growth in
  `sdk/go/testing/forecasting_test.go`: Given growth_type=NONE, verify no growth applied
- [X] T012 [P] [US1] Write conformance test verifying GetActualCost ignores growth params in
  `sdk/go/testing/forecasting_test.go`: Given growth params on resource, verify GetActualCost
  response is unchanged (FR-008)
- [X] T013 [US1] Run tests, confirm all four FAIL (proto fields don't exist yet)

### Implementation for User Story 1

- [X] T014 [US1] Add growth_type (field 9) and growth_rate (field 10) to ResourceDescriptor
  message in `proto/finfocus/v1/costsource.proto` per contracts/forecasting.proto.md
- [X] T015 [US1] Add growth_type (field 3) and growth_rate (field 4) to GetProjectedCostRequest
  message in `proto/finfocus/v1/costsource.proto` per contracts/forecasting.proto.md
- [X] T016 [US1] Add comprehensive proto comments to growth_type and growth_rate fields
  documenting semantics, valid ranges, and examples per Constitution V
- [X] T017 [US1] Run `buf lint` to validate proto changes
- [X] T018 [US1] Run `buf breaking` to confirm backward compatibility
- [X] T019 [US1] Run `make generate` to regenerate Go SDK
- [X] T020 [P] [US1] Implement ApplyLinearGrowth helper function in `sdk/go/pricing/growth.go`:
  `cost_at_n = base_cost * (1 + rate * n)`
- [X] T021 [P] [US1] Implement ApplyExponentialGrowth helper function in
  `sdk/go/pricing/growth.go`: `cost_at_n = base_cost * (1 + rate)^n`
- [X] T022 [US1] Run conformance tests, verify T009-T012 now PASS

**Checkpoint**: User Story 1 complete - growth projections work with LINEAR, EXPONENTIAL, NONE

---

## Phase 4: User Story 2 - Default Behavior (Priority: P2)

**Goal**: Ensure backward compatibility - resources without growth specs behave as before

**Independent Test**: Send requests without growth fields, verify identical to pre-feature behavior

### Tests for User Story 2 (TDD - Write First, Must Fail)

- [X] T023 [P] [US2] Write backward compatibility test in
  `sdk/go/testing/forecasting_backward_test.go`: Request without growth fields returns
  projection without growth (same as current behavior)
- [X] T024 [P] [US2] Write test for UNSPECIFIED equals NONE in
  `sdk/go/testing/forecasting_backward_test.go`: GROWTH_TYPE_UNSPECIFIED treated as
  GROWTH_TYPE_NONE
- [X] T025 [US2] Run tests, confirm they PASS (default behavior already works from Phase 3)

### Implementation for User Story 2

- [X] T026 [US2] Implement ResolveGrowthType helper in `sdk/go/pricing/growth.go` that treats
  UNSPECIFIED as NONE
- [X] T027 [US2] Implement ResolveGrowthParams helper in `sdk/go/pricing/growth.go` that merges
  request-level overrides with resource-level defaults
- [X] T028 [US2] Add unit tests for ResolveGrowthType and ResolveGrowthParams in
  `sdk/go/pricing/growth_test.go`
- [X] T029 [US2] Run all tests, verify backward compatibility maintained

**Checkpoint**: User Story 2 complete - existing clients unaffected by new fields

---

## Phase 5: User Story 3 - Validate Growth Parameters (Priority: P3)

**Goal**: Reject invalid growth parameters with clear error messages

**Independent Test**: Send invalid growth params, verify InvalidArgument errors returned

### Tests for User Story 3 (TDD - Write First, Must Fail)

- [X] T030 [P] [US3] Write validation test in `sdk/go/testing/forecasting_validation_test.go`:
  LINEAR without growth_rate returns InvalidArgument error
- [X] T031 [P] [US3] Write validation test in `sdk/go/testing/forecasting_validation_test.go`:
  EXPONENTIAL without growth_rate returns InvalidArgument error
- [X] T032 [P] [US3] Write validation test in `sdk/go/testing/forecasting_validation_test.go`:
  growth_rate < -1.0 returns InvalidArgument error
- [X] T033 [P] [US3] Write validation test in `sdk/go/testing/forecasting_validation_test.go`:
  negative growth_rate >= -1.0 is accepted (valid decline)
- [X] T034 [US3] Run validation tests, confirm they FAIL (validation not implemented)

### Implementation for User Story 3

- [X] T035 [US3] Implement ValidateGrowthParams function in `sdk/go/pricing/growth.go` per
  contracts/forecasting.proto.md validation specification
- [X] T036 [US3] Add error constants for growth validation in `sdk/go/pricing/growth.go`:
  ErrMissingGrowthRate, ErrInvalidGrowthRate
- [X] T037 [US3] Add unit tests for ValidateGrowthParams edge cases in
  `sdk/go/pricing/growth_test.go`
- [X] T038 [US3] Run validation tests, verify T030-T033 now PASS

**Checkpoint**: User Story 3 complete - invalid parameters rejected with clear messages

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, edge case handling, and final validation

- [X] T039 [P] Create example gRPC request with linear growth in
  `examples/requests/projected_cost_linear_growth.json`
- [X] T040 [P] Create example gRPC request with exponential growth in
  `examples/requests/projected_cost_exponential_growth.json`
- [X] T041 [P] Create example gRPC request with override in
  `examples/requests/projected_cost_override.json`
- [X] T042 Update sdk/go/pricing/README.md with growth helper documentation
- [X] T043 [P] Implement high-rate warning log in `sdk/go/pricing/growth.go`: Log warning via
  zerolog when growth_rate > 1.0 (>100% per period) per edge case spec
- [X] T044 [P] Implement long-projection confidence log in `sdk/go/pricing/growth.go`: Log
  info when projection period > 36 months with exponential growth per edge case spec
- [X] T045 Run `make validate` to execute full validation suite
- [X] T046 Run `make test` to verify all tests pass
- [X] T047 Run quickstart.md validation - verify code samples compile and measure
  time-to-completion (target: <30 minutes per SC-006)
- [X] T048 Update CLAUDE.md with forecasting primitives patterns if needed

**Checkpoint**: Feature complete, documented, validated

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (GrowthType enum must exist)
- **User Story 2 (Phase 4)**: Depends on User Story 1 (proto fields must exist)
- **User Story 3 (Phase 5)**: Depends on User Story 1 (proto fields must exist)
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Primary story - implements core proto changes
- **User Story 2 (P2)**: Depends on US1 proto fields for default behavior testing
- **User Story 3 (P3)**: Depends on US1 proto fields for validation implementation

Note: US2 and US3 can run in parallel after US1 completes (different test files, no conflicts)

### Within Each User Story (TDD Cycle)

1. Write tests FIRST - they MUST FAIL
2. Implement proto changes (if needed)
3. Regenerate SDK
4. Implement helper functions
5. Run tests - verify they PASS

### Parallel Opportunities

**Phase 1 (Setup)**:

```text
- T002 and T003 can run in parallel (independent linting/breaking checks)
```

**Phase 2 (Foundational)**:

```text
- T006 and T007 can run in parallel after T004/T005 (different buf commands)
```

**Phase 3 (User Story 1)**:

```text
- T009, T010, T011, T012 can run in parallel (different test cases, same file)
- T020, T021 can run in parallel (different helper functions)
```

**Phase 4 (User Story 2)**:

```text
- T023, T024 can run in parallel (different test cases)
```

**Phase 5 (User Story 3)**:

```text
- T030, T031, T032, T033 can run in parallel (different validation test cases)
```

**Phase 6 (Polish)**:

```text
- T039, T040, T041 can run in parallel (different example files)
- T043, T044 can run in parallel (different logging implementations)
```

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all conformance tests for User Story 1 together:
Task: "Write conformance test for LINEAR growth in sdk/go/testing/forecasting_test.go"
Task: "Write conformance test for EXPONENTIAL growth in sdk/go/testing/forecasting_test.go"
Task: "Write conformance test for NONE growth in sdk/go/testing/forecasting_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (verify tooling)
2. Complete Phase 2: Foundational (add GrowthType enum)
3. Complete Phase 3: User Story 1 (proto fields + growth helpers)
4. **STOP and VALIDATE**: Test growth projections work
5. Deploy/demo if ready - core forecasting is functional

### Incremental Delivery

1. Setup + Foundational → Proto infrastructure ready
2. Add User Story 1 → Core growth projections work → **MVP!**
3. Add User Story 2 → Backward compatibility verified
4. Add User Story 3 → Validation prevents misuse
5. Add Polish → Documentation and examples complete

### Constitution Compliance

- **Test-First**: All user story phases begin with failing tests
- **Proto-First**: Proto changes in Phase 2/3 before SDK helpers
- **Backward Compatible**: All new fields are optional (buf breaking passes)
- **Documented**: Proto comments required in implementation tasks

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to user story for traceability
- TDD cycle: Write failing test → Implement → Verify test passes
- Run `buf lint` and `buf breaking` after every proto change
- Commit after each task or logical group
- Constitution requires conformance tests before proto changes
