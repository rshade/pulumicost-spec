# Tasks: GetPricingSpec RPC Enhancement

**Input**: Design documents from `/specs/003-getpricingspec/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Following TDD approach per Constitution III (Test-First Protocol)

**Organization**: Tasks grouped by user story to enable independent implementation and testing

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Proto**: `proto/pulumicost/v1/costsource.proto`
- **SDK**: `sdk/go/proto/`, `sdk/go/pricing/`, `sdk/go/testing/`
- **Schemas**: `schemas/pricing_spec.schema.json`
- **Examples**: `examples/specs/*.json`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Ensure tooling is ready for proto changes

- [ ] T001 Verify buf CLI is installed via `make generate` dry run
- [ ] T002 [P] Verify JSON schema tooling via `make validate-schema` dry run
- [ ] T003 [P] Create feature branch backup point (tag current state)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Proto changes that all user stories depend on

**CRITICAL**: All user stories require these proto changes

### Tests for Foundation

- [ ] T004 Write conformance test for GetPricingSpec with flat-rate response in
      sdk/go/testing/conformance_test.go
- [ ] T005 [P] Write conformance test for GetPricingSpec with empty/default fields in
      sdk/go/testing/conformance_test.go
- [ ] T006 Run tests - confirm they FAIL (proto fields don't exist yet)
- [ ] T006a [P] Write validation helper for PricingSpec assumptions requirement in
      sdk/go/pricing/validate.go (FR-009: assumptions SHOULD be populated when not not_implemented)

### Proto Implementation

- [ ] T007 Add PricingTier message to proto/pulumicost/v1/costsource.proto after line 162
      with fields: min_quantity (1), max_quantity (2), rate_per_unit (3), description (4)
- [ ] T008 Add unit field (12) to PricingSpec message in proto/pulumicost/v1/costsource.proto
- [ ] T009 [P] Add assumptions field (13) to PricingSpec message in
      proto/pulumicost/v1/costsource.proto
- [ ] T010 [P] Add pricing_tiers field (14) to PricingSpec message in
      proto/pulumicost/v1/costsource.proto
- [ ] T011 Add proto comments for all new fields explaining purpose and usage
- [ ] T012 Run `make generate` to regenerate Go SDK from proto
- [ ] T013 Run `make lint` to verify buf lint and buf breaking pass
- [ ] T014 Run conformance tests - confirm they now PASS

**Checkpoint**: Proto changes complete, SDK regenerated, foundation tests pass

---

## Phase 3: User Story 1 - Flat-Rate Pricing (Priority: P1) MVP

**Goal**: Enable plugins to return transparent flat-rate pricing with assumptions

**Independent Test**: Call GetPricingSpec for EC2 instance, verify rate_per_unit, unit,
billing_mode, and assumptions are populated

### Tests for User Story 1

- [ ] T015 [P] [US1] Write unit test for flat-rate PricingSpec validation in
      sdk/go/pricing/validate_test.go
- [ ] T016 [P] [US1] Write integration test for GetPricingSpec with EC2 example in
      sdk/go/testing/integration_test.go
- [ ] T017 [US1] Run tests - confirm they FAIL

### Implementation for User Story 1

- [ ] T018 [P] [US1] Add Unit type constants to sdk/go/pricing/domain.go
      (Hour, GBMonth, Request, Unknown)
- [ ] T018a [P] [US1] Add BillingMode constants for Tiered and NotImplemented in
      sdk/go/pricing/domain.go
- [ ] T019 [P] [US1] Update JSON schema with unit field (string) in
      schemas/pricing_spec.schema.json
- [ ] T020 [P] [US1] Update JSON schema with assumptions field (array of strings) in
      schemas/pricing_spec.schema.json
- [ ] T021 [US1] Update mock plugin GetPricingSpec to return unit and assumptions in
      sdk/go/testing/mock_plugin.go
- [ ] T022 [US1] Create flat-rate example ec2-t3-micro.json in examples/specs/ with
      unit=hour, assumptions populated
- [ ] T023 [US1] Create flat-rate example ebs-gp3.json in examples/specs/ with
      unit=GB-month, assumptions populated
- [ ] T024 [US1] Run `make validate-examples` to verify examples pass schema validation
- [ ] T025 [US1] Run User Story 1 tests - confirm they PASS

**Checkpoint**: Flat-rate pricing with assumptions fully functional

---

## Phase 4: User Story 2 - Tiered Pricing (Priority: P2)

**Goal**: Enable plugins to return tiered pricing breakdown for volume-based resources

**Independent Test**: Call GetPricingSpec for S3 storage, verify pricing_tiers array
contains multiple tiers with min/max quantities

### Tests for User Story 2

- [ ] T026 [P] [US2] Write unit test for PricingTier validation in
      sdk/go/pricing/validate_test.go
- [ ] T027 [P] [US2] Write integration test for GetPricingSpec with tiered response in
      sdk/go/testing/integration_test.go
- [ ] T028 [US2] Run tests - confirm they FAIL

### Implementation for User Story 2

- [ ] T029 [US2] Update JSON schema with pricing_tiers array in
      schemas/pricing_spec.schema.json
- [ ] T030 [US2] Add PricingTier object schema with min_quantity, max_quantity,
      rate_per_unit, description in schemas/pricing_spec.schema.json
- [ ] T031 [US2] Update mock plugin to support tiered pricing responses in
      sdk/go/testing/mock_plugin.go
- [ ] T032 [US2] Create tiered example s3-standard.json in examples/specs/ with
      billing_mode=tiered and 3 tiers
- [ ] T033 [US2] Run `make validate-examples` to verify tiered example passes
- [ ] T034 [US2] Run User Story 2 tests - confirm they PASS

**Checkpoint**: Tiered pricing fully functional

---

## Phase 5: User Story 3 - Not-Implemented Handling (Priority: P3)

**Goal**: Enable plugins to gracefully indicate unsupported resources

**Independent Test**: Call GetPricingSpec for unsupported resource, verify
billing_mode=not_implemented and assumptions explain limitation

### Tests for User Story 3

- [ ] T035 [P] [US3] Write integration test for GetPricingSpec with not_implemented
      response in sdk/go/testing/integration_test.go
- [ ] T036 [P] [US3] Write test for gRPC InvalidArgument error when provider missing in
      sdk/go/testing/integration_test.go
- [ ] T037 [P] [US3] Write test for gRPC NotFound error for unknown SKU in
      sdk/go/testing/integration_test.go
- [ ] T038 [US3] Run tests - confirm they FAIL

### Implementation for User Story 3

- [ ] T039 [US3] Update mock plugin to return not_implemented for unknown resource types
      in sdk/go/testing/mock_plugin.go
- [ ] T040 [US3] Update mock plugin to return InvalidArgument for missing provider in
      sdk/go/testing/mock_plugin.go
- [ ] T041 [US3] Update mock plugin to return NotFound for unknown SKU in
      sdk/go/testing/mock_plugin.go
- [ ] T042 [US3] Create not-implemented example lambda-stub.json in examples/specs/
      with billing_mode=not_implemented
- [ ] T043 [US3] Run `make validate-examples` to verify not-implemented example passes
- [ ] T044 [US3] Run User Story 3 tests - confirm they PASS

**Checkpoint**: Error handling and not-implemented graceful degradation complete

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, validation, and final cleanup

- [ ] T045 [P] Update examples/README.md with new example descriptions
- [ ] T046 [P] Update sdk/go/testing/README.md with GetPricingSpec testing guidance
- [ ] T047 Add GetPricingSpec enhancement entry to CHANGELOG.md under Unreleased
- [ ] T048 Run `make validate` to verify all tests, linting, and schema validation pass
- [ ] T049 [P] Run benchmarks `go test -bench=. -benchmem ./sdk/go/testing/` to verify
      performance
- [ ] T050 Review and update quickstart.md if implementation differs from plan
- [ ] T051 Final code review for proto comments and documentation completeness

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in priority order (P1 → P2 → P3)
  - Each story builds incrementally but is independently testable
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational - Builds on schema from US1
- **User Story 3 (P3)**: Can start after Foundational - Builds on mock from US1/US2

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Schema before implementation
- Mock updates before examples
- Examples must pass validation before marking story complete

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational proto field additions (T008-T010) can run in parallel
- All test tasks within a story marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: T015 "Write unit test for flat-rate PricingSpec validation"
Task: T016 "Write integration test for GetPricingSpec with EC2 example"

# Launch all parallel implementation tasks:
Task: T018 "Add Unit type constants to sdk/go/pricing/domain.go"
Task: T019 "Update JSON schema with unit field"
Task: T020 "Update JSON schema with assumptions field"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (proto changes)
3. Complete Phase 3: User Story 1 (flat-rate pricing)
4. **STOP and VALIDATE**: Test flat-rate pricing independently
5. Plugin developers can start using GetPricingSpec for flat-rate resources

### Incremental Delivery

1. Complete Setup + Foundational → Proto enhanced
2. Add User Story 1 → Flat-rate works → Demo/Release
3. Add User Story 2 → Tiered works → Demo/Release
4. Add User Story 3 → Error handling → Demo/Release
5. Each story adds value without breaking previous functionality

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- TDD strictly enforced per Constitution III
- Run `make validate` at each checkpoint
- Commit after each phase completion
- Version bump (0.1.0 → 0.2.0) after all stories complete
