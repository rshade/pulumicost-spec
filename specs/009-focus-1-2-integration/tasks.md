# Tasks: FinOps FOCUS 1.2 Integration

**Feature Branch**: `009-focus-1-2-integration`
**Spec**: [specs/009-focus-1-2-integration/spec.md](./spec.md)

## Phase 1: Setup & Infrastructure

**Goal**: Initialize new Protobuf contracts and prepare the project for FOCUS 1.2 integration.

- [x] T001 Create `proto/pulumicost/v1/enums.proto` with Focus Service, Charge, and Pricing enums `specs/009-focus-1-2-integration/contracts/enums.proto`
- [x] T002 Create `proto/pulumicost/v1/focus.proto` with `FocusCostRecord` message definition `specs/009-focus-1-2-integration/contracts/focus.proto`
- [x] T003 Update `proto/pulumicost/v1/costsource.proto` to import `enums.proto`
      and `focus.proto` `proto/pulumicost/v1/costsource.proto`
- [x] T004 Run `buf generate` to generate Go SDK code for new protos `Makefile`

## Phase 2: Foundational

**Goal**: Establish core test harness and interfaces required for all user stories.

- [x] T005 Create `sdk/go/pluginsdk/focus_builder.go` with initial `FocusBuilder` interface definition `sdk/go/pluginsdk/focus_builder.go`
- [x] T006 Create empty `sdk/go/pluginsdk/focus_builder_test.go` test file `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T007 [P] Create `sdk/go/pluginsdk/focus_conformance.go` and implement validation logic `sdk/go/pluginsdk/focus_conformance.go`

## Phase 3: User Story 1 - Create Standardized Cost Records

**Goal**: Enable developers to create FOCUS-compliant records with mandatory fields validation.
**Priority**: P1

- [x] T008 [US1] Add tests: Verify happy path, missing mandatory fields, and enum enforcement `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T009 [US1] Implement `NewFocusRecordBuilder()` constructor `sdk/go/pluginsdk/focus_builder.go`
- [x] T010 [US1] Implement `WithServiceCategory` and standard field setters in Builder `sdk/go/pluginsdk/focus_builder.go`
- [x] T011 [US1] Implement `WithChargeDetails` and `WithPricingCategory` setters `sdk/go/pluginsdk/focus_builder.go`
- [x] T012 [US1] Implement `WithFinancials` and `WithIdentity` setters `sdk/go/pluginsdk/focus_builder.go`
- [x] T013 [US1] Implement `Build()` method with mandatory field validation
      (AccountID, Dates, Service/Charge/Pricing, BilledCost) `sdk/go/pluginsdk/focus_builder.go`

## Phase 4: User Story 2 - Future-Proof Extension ("Backpack")

**Goal**: Allow arbitrary extension data via `extended_columns` without schema changes.
**Priority**: P2

- [x] T014 [US2] Add tests: Verify extension data is preserved and multiple extensions handled `sdk/go/pluginsdk/focus_builder_test.go`
- [x] T015 [US2] Implement `WithExtension(key, value)` method in Builder `sdk/go/pluginsdk/focus_builder.go`
- [x] T016 [US2] Update `Build()` to populate `extended_columns` map in the record `sdk/go/pluginsdk/focus_builder.go`

## Phase 5: User Story 3 - Stable Upgrade Path ("Shield")

**Goal**: Ensure internal refactoring doesn't break plugin code (Integration/Pattern verification).
**Priority**: P3

- [x] T017 [US3] [P] Create `examples/plugins/focus-example.go` demonstrating full Builder usage `examples/plugins/focus-example.go`
- [x] T018 [US3] Verify example compiles against generated SDK `examples/plugins/focus-example.go`
- [x] T019 [US3] Document usage of Builder vs direct Struct in `docs/PLUGIN_MIGRATION_GUIDE.md` `docs/PLUGIN_MIGRATION_GUIDE.md`

## Phase 6: Polish & Cross-Cutting

**Goal**: Finalize documentation, benchmarks, and cleanup.

- [x] T020 Create `sdk/go/pluginsdk/focus_benchmark_test.go` for allocation benchmarks `sdk/go/pluginsdk/focus_benchmark_test.go`
- [x] T021 Run benchmarks; if map allocation exceeds X bytes/record or Y ns/op, optimize `sdk/go/pluginsdk/focus_builder.go`
- [x] T022 Update `README.md` with FOCUS 1.2 support badge and quick links `README.md`
- [x] T023 Final `buf lint` and `go vet` check `Makefile`

## Dependencies

1. **T001-T004 (Proto)** -> MUST be done first to generate Go types.
2. **T005-T007 (Foundational)** -> Unlocks parallel work on Builder methods.
3. **T008 (US1 Tests)** -> Defines requirements for T009-T013.
4. **T009-T013 (US1 Impl)** -> Core implementation.
5. **T014 (US2 Tests)** -> Defines requirements for T015-T016.
6. **T017 (Example)** -> Depends on full Builder implementation (US1+US2).

## Implementation Strategy

- **MVP Scope**: T001-T013 (Proto + US1). This delivers the core mandatory FOCUS record creation.
- **Full Scope**: All tasks. Adds the critical "Backpack" feature and documentation.
- **Parallel Execution**:
  - T007 (Conformance) can be written alongside T009-T013.
  - T017 (Example) can be started once T009 is defined, refining as methods are added.
