# Tasks: FOCUS 1.2 Column Audit

**Input**: Design documents from `/specs/010-focus-column-audit/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅, quickstart.md ✅

**Tests**: Tests ARE included as this is a spec/SDK project where conformance testing is required.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Tech Stack

- **Language**: Go 1.24+ (toolchain go1.25.4), Protocol Buffers v3
- **Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
- **Testing**: go test, conformance tests via bufconn harness, buf lint/breaking
- **Constraints**: Backward compatible proto changes, buf breaking check must pass

## Path Conventions

- **Proto definitions**: `proto/finfocus/v1/`
- **Generated code**: `sdk/go/proto/finfocus/v1/`
- **SDK code**: `sdk/go/pluginsdk/`
- **Documentation**: `docs/`, `sdk/go/pluginsdk/README.md`
- **Examples**: `examples/plugins/`
- **Scripts**: `scripts/`

---

## Phase 1: Setup

**Purpose**: Validate current state and prepare for implementation

- [x] T001 Create feature branch `010-focus-column-audit` if not exists
- [x] T002 [P] Verify buf CLI v1.32.1 is installed via `make generate`
- [x] T003 [P] Run `buf breaking` to establish baseline compatibility check

---

## Phase 2: Foundational (Proto Definitions)

**Purpose**: Add new enum types and proto fields - BLOCKS all SDK implementation

**⚠️ CRITICAL**: No SDK work can begin until proto changes are complete and generated

### New Enum Types

- [x] T004 [P] Add `FocusCommitmentDiscountStatus` enum to `proto/finfocus/v1/enums.proto`
- [x] T005 [P] Add `FocusCapacityReservationStatus` enum to `proto/finfocus/v1/enums.proto`

### Proto Field Additions (19 fields)

- [x] T006 [P] Add `contracted_cost` (field 41) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T007 [P] Add `billing_account_type` (field 42) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T008 [P] Add `sub_account_type` (field 43) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T009 [P] Add `capacity_reservation_id` (field 44) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T010 [P] Add `capacity_reservation_status` (field 45) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T011 [P] Add `commitment_discount_quantity` (field 46) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T012 [P] Add `commitment_discount_status` (field 47) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T013 [P] Add `commitment_discount_type` (field 48) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T014 [P] Add `commitment_discount_unit` (field 49) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T015 [P] Add `contracted_unit_price` (field 50) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T016 [P] Add `pricing_currency` (field 51) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T017 [P] Add `pricing_currency_contracted_unit_price` (field 52) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T018 [P] Add `pricing_currency_effective_cost` (field 53) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T019 [P] Add `pricing_currency_list_unit_price` (field 54) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T020 [P] Add `publisher` (field 55) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T021 [P] Add `service_subcategory` (field 56) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T022 [P] Add `sku_meter` (field 57) to FocusCostRecord in `proto/finfocus/v1/focus.proto`
- [x] T023 [P] Add `sku_price_details` (field 58) to FocusCostRecord in `proto/finfocus/v1/focus.proto`

### Proto Validation and Generation

- [x] T024 Run `buf lint` to validate proto syntax in `proto/finfocus/v1/`
- [x] T025 Run `buf breaking` to confirm backward compatibility
- [x] T026 Run `make generate` to regenerate Go code in `sdk/go/proto/finfocus/v1/`
- [x] T027 Run `go build ./...` to verify generated code compiles

**Checkpoint**: Proto definitions complete - SDK implementation can now begin

---

## Phase 3: User Story 1 - Schema Completeness Verification (Priority: P1) 🎯 MVP

**Goal**: Verify all 57 FOCUS 1.2 columns are present in focus.proto with correct types

**Independent Test**: Run audit script to compare proto against FOCUS 1.2 specification

### Tests for User Story 1

- [x] T028 [US1] Write conformance test for 57 FOCUS columns in `sdk/go/pluginsdk/focus_conformance_test.go`
- [x] T029 [US1] Write test for new enum types validation in `sdk/go/pluginsdk/focus_conformance_test.go`

### Implementation for User Story 1

- [x] T030 [US1] Create audit script `scripts/audit_focus_columns.go` to verify column coverage
- [x] T031 [US1] Add FOCUS section reference comments to all 19 new fields in `proto/finfocus/v1/focus.proto`
- [x] T032 [US1] Add FOCUS section reference comments to new enums in `proto/finfocus/v1/enums.proto`
- [x] T033 [US1] Run audit script and verify 57/57 columns pass

**Checkpoint**: All 57 FOCUS 1.2 columns verified in proto with correct types

---

## Phase 4: User Story 2 - Type Correctness Validation (Priority: P2)

**Goal**: Ensure proto types align with FOCUS data types (Decimal→double, DateTime→Timestamp)

**Independent Test**: Audit script validates type mappings

### Tests for User Story 2

- [x] T034 [US2] Add type mapping tests to `sdk/go/pluginsdk/focus_conformance_test.go`
- [x] T035 [US2] Add enum value completeness tests to `sdk/go/pluginsdk/focus_conformance_test.go`

### Implementation for User Story 2

- [x] T036 [US2] Extend audit script to validate type mappings in `scripts/audit_focus_columns.go`
- [x] T037 [US2] Verify all enum values match FOCUS 1.2 specification (run enum conformance tests)
- [x] T038 [US2] Update `sdk/go/pluginsdk/focus_conformance.go` with new column validation

**Checkpoint**: Type correctness verified - proto types match FOCUS specification

---

## Phase 5: User Story 3 - Builder API Completeness (Priority: P3)

**Goal**: Add builder methods for all 19 new FOCUS columns

**Independent Test**: Every field in FocusCostRecord has a corresponding builder method

### Tests for User Story 3

- [x] T039 [US3] Write tests for `WithContractedCost` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T040 [P] [US3] Write tests for `WithBillingAccountType` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T041 [P] [US3] Write tests for `WithSubAccountType` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T042 [P] [US3] Write tests for `WithCapacityReservation` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T043 [P] [US3] Write tests for `WithCommitmentDiscountDetails` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T044 [P] [US3] Write tests for `WithContractedUnitPrice` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T045 [P] [US3] Write tests for `WithPricingCurrency` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T046 [P] [US3] Write tests for `WithPricingCurrencyPrices` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T047 [P] [US3] Write tests for `WithPublisher` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T048 [P] [US3] Write tests for `WithServiceSubcategory` in `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T049 [P] [US3] Write tests for `WithSkuDetails` in `sdk/go/pluginsdk/focus_builder_test.go`

### Implementation for User Story 3

- [x] T050 [US3] Implement `WithContractedCost(cost float64)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T051 [P] [US3] Implement `WithBillingAccountType(accountType string)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T052 [P] [US3] Implement `WithSubAccountType(accountType string)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T053 [P] [US3] Implement `WithCapacityReservation(id string, status)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T054 [P] [US3] Implement `WithCommitmentDiscountDetails(qty, status, discountType, unit)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T055 [P] [US3] Implement `WithContractedUnitPrice(price float64)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T056 [P] [US3] Implement `WithPricingCurrency(currency string)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T057 [P] [US3] Implement `WithPricingCurrencyPrices(contracted, effective, list float64)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T058 [P] [US3] Implement `WithPublisher(publisher string)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T059 [P] [US3] Implement `WithServiceSubcategory(subcategory string)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T060 [P] [US3] Implement `WithSkuDetails(meter, priceDetails string)` in `sdk/go/pluginsdk/focus_builder.go`
- [x] T061 [US3] Update `ValidateFocusRecord` to validate ContractedCost (mandatory) in `sdk/go/pluginsdk/focus_conformance.go`
- [x] T062 [US3] Run all builder tests and verify 100% pass rate

**Checkpoint**: All 19 new columns have builder methods with tests

---

## Phase 6: User Story 4 - Developer Documentation (Priority: P2)

**Goal**: Achieve 80%+ godoc coverage with comprehensive developer guide

**Independent Test**: New developer can build FOCUS record using only documentation

### Tests for User Story 4

- [x] T063 [US4] Verify godoc coverage meets 80%+ threshold

### Implementation for User Story 4

- [x] T064 [P] [US4] Add godoc comments to all new builder methods in `sdk/go/pluginsdk/focus_builder.go`
- [x] T065 [P] [US4] Add godoc comments to new enum types in generated code comments (proto source)
- [x] T066 [US4] Update `sdk/go/pluginsdk/README.md` with quick start example for new columns
- [x] T067 [US4] Add migration guide section to `sdk/go/pluginsdk/README.md`
- [x] T068 [US4] Add troubleshooting guide for validation errors to `sdk/go/pluginsdk/README.md`

**Checkpoint**: 80%+ godoc coverage achieved, developer guide complete

---

## Phase 7: User Story 5 - User-Facing Documentation (Priority: P3)

**Goal**: Create FOCUS column reference with provider mappings

**Independent Test**: Non-developer can understand each column from documentation alone

### Implementation for User Story 5

- [x] T069 [P] [US5] Create `docs/focus-columns.md` with plain-language column descriptions
- [x] T070 [P] [US5] Add provider mapping table (AWS, Azure, GCP) to `docs/focus-columns.md`
- [x] T071 [US5] Add common use cases and query patterns to `docs/focus-columns.md`
- [x] T072 [US5] Update `examples/plugins/focus_example.go` with complete 57-column example

**Checkpoint**: User documentation complete with provider mappings

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T073 [P] Run `make lint` to ensure all code passes linting
- [x] T074 [P] Run `make test` to verify all tests pass
- [x] T075 [P] Run `buf breaking` to confirm final backward compatibility
- [x] T076 Run `scripts/audit_focus_columns.go` and verify 57/57 columns pass
- [x] T077 Run quickstart.md validation (compile and run example code)
- [ ] T078 Update CHANGELOG.md with feature changes (via release-please PR)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all SDK implementation
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - US1 (Schema Completeness) must complete before US2 (Type Correctness)
  - US2 and US3 can proceed in parallel after US1
  - US4 (Developer Docs) and US5 (User Docs) can proceed in parallel after US3
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

```text
Phase 1: Setup
    ↓
Phase 2: Foundational (Proto changes + generation)
    ↓
Phase 3: US1 - Schema Completeness ────────────────┐
    ↓                                              │
Phase 4: US2 - Type Correctness                    │
    ↓                                              │
Phase 5: US3 - Builder API ←───────────────────────┘
    ↓
    ├── Phase 6: US4 - Developer Docs ──┐
    │                                   ├── Phase 8: Polish
    └── Phase 7: US5 - User Docs ───────┘
```

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Proto fields before builder methods
- Builder methods before validation updates
- Implementation before documentation

### Parallel Opportunities

**Phase 2 (Proto Changes)**: T004-T023 can all run in parallel (different lines in proto files)

**Phase 5 (Builder Methods)**: T040-T049 tests and T051-T060 implementations can run in parallel

**Phase 6-7 (Documentation)**: US4 and US5 can run in parallel after US3 completes

---

## Parallel Example: Proto Field Additions

```bash
# Launch all proto field additions together (Phase 2):
Task: "Add contracted_cost (field 41) to FocusCostRecord"
Task: "Add billing_account_type (field 42) to FocusCostRecord"
Task: "Add sub_account_type (field 43) to FocusCostRecord"
# ... all 19 fields in parallel
```

## Parallel Example: Builder Method Tests

```bash
# Launch all builder tests together (Phase 5):
Task: "Write tests for WithBillingAccountType"
Task: "Write tests for WithSubAccountType"
Task: "Write tests for WithCapacityReservation"
# ... all 11 method tests in parallel
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (Proto changes)
3. Complete Phase 3: User Story 1 (Schema Completeness)
4. **STOP and VALIDATE**: Run audit script to verify 57/57 columns
5. This provides immediate value: FOCUS 1.2 compliance verified

### Incremental Delivery

1. Setup + Foundational → Proto changes complete
2. US1 (Schema Completeness) → Audit script validates all columns
3. US2 (Type Correctness) → Type mapping verified
4. US3 (Builder API) → Full SDK support for new columns
5. US4 + US5 (Documentation) → Developer and user guides complete
6. Polish → Final validation and release preparation

### Single Developer Strategy

With one developer, complete phases sequentially:

1. Phase 1 → Phase 2 → Phase 3 (MVP complete)
2. Phase 4 → Phase 5 → Phase 6/7 → Phase 8

Estimated effort: ~40-60 tasks, ~2-3 hours for proto changes, ~4-6 hours for builder + tests

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently testable via audit script and conformance tests
- Commit after each phase completion
- Stop at any checkpoint to validate incrementally
- Proto field numbers 41-58 are reserved for these additions
