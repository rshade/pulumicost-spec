# Tasks: PluginSDK Request Validation

**Feature**: `017-pluginsdk-validation`
**Status**: Complete

## Phase 1: Setup

**Goal**: Initialize project structure for validation helpers.

- [X] T001 Create `sdk/go/pluginsdk/validation.go` with package declaration in `sdk/go/pluginsdk/validation.go`
- [X] T002 Create `sdk/go/pluginsdk/validation_test.go` with test skeleton and imports in `sdk/go/pluginsdk/validation_test.go`

## Phase 2: Foundation

**Goal**: Establish common dependencies and testing infrastructure.

- [X] T003 Import `github.com/rshade/pulumicost-spec/proto/pulumicost/v1` in `sdk/go/pluginsdk/validation.go`
- [X] T004 [P] Define `TestValidateProjectedCostRequest` function in `sdk/go/pluginsdk/validation_test.go`
- [X] T005 [P] Define `TestValidateActualCostRequest` function in `sdk/go/pluginsdk/validation_test.go`

## Phase 3: Core Pre-flight Validation (P1)

**Goal**: Implement validation for `GetProjectedCostRequest` to catch configuration errors early.
**Tests**: `TestValidateProjectedCostRequest` covering all edge cases.

- [X] T006 [P] [US1] Add test case: Nil request returns error in `sdk/go/pluginsdk/validation_test.go`
- [X] T007 [P] [US1] Add test case: Nil Resource returns error "resource is required" in `sdk/go/pluginsdk/validation_test.go`
- [X] T008 [P] [US1] Add test case: Empty Provider returns error "resource.provider is required" in `sdk/go/pluginsdk/validation_test.go`
- [X] T009 [P] [US1] Add test case: Empty ResourceType returns error "resource.resource_type is required" in `sdk/go/pluginsdk/validation_test.go`
- [X] T010 [US1] Implement `ValidateProjectedCostRequest` signature and nil checks in `sdk/go/pluginsdk/validation.go`
- [X] T011 [US1] Implement validation for `Resource` and `Provider` fields in `sdk/go/pluginsdk/validation.go`
- [X] T012 [US1] Implement validation for `ResourceType` field in `sdk/go/pluginsdk/validation.go`

## Phase 4: Plugin Defense-in-Depth Validation (P2)

**Goal**: Implement validation for `GetActualCostRequest` for robustness.
**Tests**: `TestValidateActualCostRequest` covering time ranges and required fields.

- [X] T013 [P] [US2] Add test case: Nil request returns error in `sdk/go/pluginsdk/validation_test.go`
- [X] T014 [P] [US2] Add test case: Empty ResourceId returns error in `sdk/go/pluginsdk/validation_test.go`
- [X] T015 [P] [US2] Add test case: Nil StartTime/EndTime returns error in `sdk/go/pluginsdk/validation_test.go`
- [X] T016 [P] [US2] Add test case: EndTime before StartTime returns error "end time must be after start time" in `sdk/go/pluginsdk/validation_test.go`
- [X] T017 [US2] Implement `ValidateActualCostRequest` signature and nil checks in `sdk/go/pluginsdk/validation.go`
- [X] T018 [US2] Implement `ResourceId` and Timestamp presence checks in `sdk/go/pluginsdk/validation.go`
- [X] T019 [US2] Implement TimeRange validation logic (End > Start) in `sdk/go/pluginsdk/validation.go`

## Phase 5: Actionable Error Messages (P3)

**Goal**: Ensure error messages provide specific guidance referencing mapping helpers.
**Tests**: Verify error strings contain "mapping.Extract..."

- [X] T020 [P] [US3] Add test case: Empty SKU expects error containing "use mapping.ExtractAWSSKU" in `sdk/go/pluginsdk/validation_test.go`
- [X] T021 [P] [US3] Add test case: Empty Region expects error containing "use mapping.ExtractAWSRegion" in `sdk/go/pluginsdk/validation_test.go`
- [X] T022 [US3] Implement validation for `Sku` field with specific error message in `sdk/go/pluginsdk/validation.go`
- [X] T023 [US3] Implement validation for `Region` field with specific error message in `sdk/go/pluginsdk/validation.go`

## Phase 6: Polish & Cross-Cutting

**Goal**: Code quality and standard adherence.

- [X] T024 Run `go fmt ./sdk/go/pluginsdk/...` to ensure formatting
- [X] T025 Run `go test ./sdk/go/pluginsdk/...` to verify all tests pass
- [X] T026 Verify zero allocations (benchmarks optional but recommended) in `sdk/go/pluginsdk/validation_test.go`

## Dependencies

1. Phase 1 & 2 (Setup/Foundation) -> MUST complete before Phase 3, 4, 5
2. Phase 3 (US1) -> Independent of US2
3. Phase 4 (US2) -> Independent of US1
4. Phase 5 (US3) -> Depends on basic structure of US1 (extends it)

## Implementation Strategy

1. **MVP**: Complete Phase 1, 2, and 3 (Projected Cost Validation).
2. **Increment 1**: Complete Phase 4 (Actual Cost Validation).
3. **Increment 2**: Complete Phase 5 (Refine Error Messages).
