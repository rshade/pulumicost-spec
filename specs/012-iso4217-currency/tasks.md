# Tasks: ISO 4217 Currency Validation Package

**Input**: Design documents from `/specs/012-iso4217-currency/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests ARE required per FR-007 (>90% coverage) and constitution (Test-First Protocol).

**Organization**: Tasks grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Exact file paths included in descriptions

## Path Conventions

Per plan.md structure:

```text
sdk/go/currency/           # New package
sdk/go/pluginsdk/          # Migration target
```

---

## Phase 1: Setup (Package Infrastructure)

**Purpose**: Create package structure and basic scaffolding

- [x] T001 Create package directory at sdk/go/currency/
- [x] T002 Create package documentation at sdk/go/currency/README.md
- [x] T003 [P] Create package-specific CLAUDE.md at sdk/go/currency/CLAUDE.md

---

## Phase 2: Foundational (Core Data Structure)

**Purpose**: Currency type and data that ALL user stories depend on

**⚠️ CRITICAL**: All user stories require this phase to complete first

- [x] T004 Define Currency struct in sdk/go/currency/currency.go
- [x] T005 Create allCurrencies package-level slice with complete ISO 4217 data
      (180+ currencies) in sdk/go/currency/currency.go
- [x] T006 Create currencyByCode map for O(1) lookup in sdk/go/currency/currency.go
- [x] T007 Add ErrCurrencyNotFound error variable in sdk/go/currency/currency.go
- [x] T008 Implement String() method for Currency type in sdk/go/currency/currency.go

**Checkpoint**: Core data structure ready - user story implementation can begin

---

## Phase 3: User Story 1 - Validate Currency Codes (Priority: P1) 🎯 MVP

**Goal**: Provide zero-allocation currency code validation function

**Independent Test**: Call IsValid() with various codes and verify boolean responses

### Tests for User Story 1

> **NOTE: Write tests FIRST, ensure they FAIL before implementation (TDD)**

- [x] T009 [P] [US1] Write table-driven tests for IsValid() function covering valid codes
      (USD, EUR, JPY, XXX, XTS) in sdk/go/currency/currency_test.go
- [x] T010 [P] [US1] Write tests for invalid codes (XYZ, usd, empty, whitespace)
      in sdk/go/currency/currency_test.go
- [x] T011 [P] [US1] Write edge case tests (historic codes DEM, supranational XDR)
      in sdk/go/currency/currency_test.go

### Implementation for User Story 1

- [x] T012 [US1] Implement IsValid(code string) bool function using linear scan
      in sdk/go/currency/validate.go
- [x] T013 [US1] Verify all IsValid tests pass with `go test -v ./sdk/go/currency/`
- [x] T014 [P] [US1] Write benchmark for IsValid() targeting <15 ns/op
      in sdk/go/currency/benchmark_test.go
- [x] T015 [US1] Run benchmarks and verify 0 allocs/op with
      `go test -bench=. -benchmem ./sdk/go/currency/`

**Checkpoint**: IsValid() works independently - core validation complete (MVP)

---

## Phase 4: User Story 2 - Retrieve Currency Metadata (Priority: P2)

**Goal**: Provide function to retrieve full currency metadata by code

**Independent Test**: Call GetCurrency() with valid/invalid codes and verify metadata/error

### Tests for User Story 2

- [x] T016 [P] [US2] Write tests for GetCurrency() with valid codes returning metadata
      (USD, JPY, KWD) in sdk/go/currency/currency_test.go
- [x] T017 [P] [US2] Write tests for GetCurrency() with invalid codes returning
      ErrCurrencyNotFound in sdk/go/currency/currency_test.go
- [x] T018 [P] [US2] Write tests verifying metadata values (name, numeric code, minor units)
      in sdk/go/currency/currency_test.go

### Implementation for User Story 2

- [x] T019 [US2] Implement GetCurrency(code string) (\*Currency, error) using map lookup
      in sdk/go/currency/currency.go
- [x] T020 [US2] Verify all GetCurrency tests pass
- [x] T021 [P] [US2] Write benchmark for GetCurrency() in sdk/go/currency/benchmark_test.go

**Checkpoint**: GetCurrency() works independently - metadata retrieval complete

---

## Phase 5: User Story 3 - List All Currencies (Priority: P3)

**Goal**: Provide function to list all valid ISO 4217 currencies

**Independent Test**: Call AllCurrencies() and verify list contains 180+ entries with metadata

### Tests for User Story 3

- [x] T022 [P] [US3] Write test verifying AllCurrencies() returns 180+ currencies
      in sdk/go/currency/currency_test.go
- [x] T023 [P] [US3] Write test verifying each currency in list has non-empty fields
      in sdk/go/currency/currency_test.go
- [x] T024 [P] [US3] Write test verifying list is sorted alphabetically by code
      in sdk/go/currency/currency_test.go

### Implementation for User Story 3

- [x] T025 [US3] Implement AllCurrencies() []Currency returning package-level slice
      in sdk/go/currency/currency.go
- [x] T026 [US3] Verify all AllCurrencies tests pass
- [x] T027 [P] [US3] Write benchmark for AllCurrencies() in sdk/go/currency/benchmark_test.go

**Checkpoint**: AllCurrencies() works independently - listing complete

---

## Phase 6: User Story 4 - Migrate Existing Validation (Priority: P4)

**Goal**: Update focus_conformance.go to use new currency package

**Independent Test**: Run existing FOCUS conformance tests and verify they pass

### Tests for User Story 4

- [ ] T028 [US4] Run existing pluginsdk tests to establish baseline with
      `go test -v ./sdk/go/pluginsdk/`

### Implementation for User Story 4

- [ ] T029 [US4] Add import for currency package in sdk/go/pluginsdk/focus_conformance.go
- [ ] T030 [US4] Update validateCurrency() to use currency.IsValid()
      in sdk/go/pluginsdk/focus_conformance.go
- [ ] T031 [US4] Remove inline iso4217Currencies map from sdk/go/pluginsdk/focus_conformance.go
- [ ] T032 [US4] Run pluginsdk tests to verify no regression with
      `go test -v ./sdk/go/pluginsdk/`
- [ ] T033 [US4] Run full test suite with `make test` to verify integration

**Checkpoint**: Migration complete - all existing tests pass with new package

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, coverage verification, and final validation

- [ ] T034 [P] Run test coverage report with `go test -cover ./sdk/go/currency/`
      and verify >90%
- [ ] T035 [P] Add package documentation comments in sdk/go/currency/doc.go
- [ ] T036 [P] Run linting with `make lint` and fix any issues
- [ ] T037 Verify quickstart.md examples work by creating test file
- [ ] T038 Update sdk/go/CLAUDE.md to reference new currency package
- [ ] T039 Run full validation with `make validate`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - US1-US3 can proceed in parallel (different functions)
  - US4 depends on US1 completion (needs IsValid())
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: After Foundational - No dependencies on other stories
- **User Story 2 (P2)**: After Foundational - No dependencies on other stories
- **User Story 3 (P3)**: After Foundational - No dependencies on other stories
- **User Story 4 (P4)**: After Foundational + US1 (uses IsValid())

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD per constitution)
- Implementation follows tests
- Benchmarks can be parallel with implementation verification

### Parallel Opportunities

- T002, T003 can run in parallel (different files)
- T009, T010, T011 can run in parallel (all US1 test files)
- T016, T017, T018 can run in parallel (all US2 test files)
- T022, T023, T024 can run in parallel (all US3 test files)
- US1, US2, US3 can be worked in parallel after Foundational
- T034, T035, T036 can run in parallel (different concerns)

---

## Parallel Example: User Story 1

```bash
# Launch all US1 tests together:
Task: "Write table-driven tests for IsValid() in sdk/go/currency/currency_test.go"
Task: "Write tests for invalid codes in sdk/go/currency/currency_test.go"
Task: "Write edge case tests in sdk/go/currency/currency_test.go"

# After implementation, launch benchmark in parallel with verification:
Task: "Write benchmark for IsValid() in sdk/go/currency/benchmark_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - Currency struct + data)
3. Complete Phase 3: User Story 1 (IsValid function)
4. **STOP and VALIDATE**: Test IsValid() independently
5. Core validation is now usable

### Incremental Delivery

1. Setup + Foundational → Data structure ready
2. User Story 1 → IsValid() works → MVP complete!
3. User Story 2 → GetCurrency() works → Metadata retrieval
4. User Story 3 → AllCurrencies() works → Listing capability
5. User Story 4 → Migration complete → Ecosystem integrated

### Full Implementation (Recommended)

Since US1-US3 are all in sdk/go/currency/:

1. Complete Setup + Foundational
2. Complete US1, US2, US3 (all in same package, efficient to do together)
3. Complete US4 (migration, separate package)
4. Complete Polish phase

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- TDD required: write tests first, verify they fail, then implement
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Performance targets: <15 ns/op, 0 allocs/op for IsValid()
