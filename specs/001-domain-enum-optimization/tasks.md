# Tasks: Domain Enum Validation Performance Optimization

**Input**: Design documents from `/specs/001-domain-enum-optimization/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Test tasks included per TDD requirement from constitution (Test-First Protocol)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Single project Go SDK structure
- Modifications in `sdk/go/registry/`
- Tests in `sdk/go/registry/` (co-located with implementation)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project verification and baseline establishment

- [x] T001 Verify Go 1.24.10 (toolchain go1.25.4) installation and project dependencies
- [x] T002 Run existing test suite to establish baseline in sdk/go/registry/
- [x] T003 [P] Run existing benchmarks to capture current performance in sdk/go/registry/

**Checkpoint**: Baseline established - current performance metrics documented

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core testing infrastructure that MUST be complete before implementing optimization

**⚠️ CRITICAL**: TDD requirement - tests must be written and FAIL before optimization implementation

- [ ] T004 Create benchmark test infrastructure following sdk/go/testing/benchmark_test.go patterns
- [ ] T005 Document current performance baseline from T003 in benchmark comments

**Checkpoint**: Test infrastructure ready - optimization implementation can begin

---

## Phase 3: User Story 1 - Plugin Developer Validates Domain Values (Priority: P1) 🎯 MVP

**Goal**: Optimize all 8 enum validation functions for zero-allocation, fast validation while maintaining 100%
backward compatibility

**Independent Test**: All validation functions return correct results (100% accuracy) with zero allocations and
< 30 ns/op performance

### Tests for User Story 1 (TDD Required - Write FIRST, ensure they FAIL)

> **NOTE: Write these tests FIRST per constitution Test-First Protocol requirement**

- [ ] T006 [P] [US1] Add benchmark for IsValidProvider (current implementation) in sdk/go/registry/domain_test.go
- [ ] T007 [P] [US1] Add benchmark for IsValidDiscoverySource (current) in sdk/go/registry/domain_test.go
- [ ] T008 [P] [US1] Add benchmark for IsValidPluginStatus (current) in sdk/go/registry/domain_test.go
- [ ] T009 [P] [US1] Add benchmark for IsValidSecurityLevel (current) in sdk/go/registry/domain_test.go
- [ ] T010 [P] [US1] Add benchmark for IsValidInstallationMethod (current) in sdk/go/registry/domain_test.go
- [ ] T011 [P] [US1] Add benchmark for IsValidPluginCapability (current) in sdk/go/registry/domain_test.go
- [ ] T012 [P] [US1] Add benchmark for IsValidSystemPermission (current) in sdk/go/registry/domain_test.go
- [ ] T013 [P] [US1] Add benchmark for IsValidAuthMethod (current) in sdk/go/registry/domain_test.go
- [ ] T014 [P] [US1] Add unit tests for edge cases (empty string, case mismatch, nil) in sdk/go/registry/domain_test.go
- [ ] T015 [US1] Run benchmarks with `go test -bench=. -benchmem` and verify they capture current allocation pattern

**Checkpoint**: All benchmarks written and show current performance (should show allocations > 0)

### Implementation for User Story 1

> **Implementation order**: Package-level variables first, then update functions to use them

**Provider Enum (5 values)**:

- [ ] T016 [US1] Create package-level variable `allProviders` in sdk/go/registry/domain.go
- [ ] T017 [US1] Update `AllProviders()` to return `allProviders` in sdk/go/registry/domain.go
- [ ] T018 [US1] Update `IsValidProvider()` to iterate `allProviders` in sdk/go/registry/domain.go

**DiscoverySource Enum (4 values)**:

- [ ] T019 [US1] Create package-level variable `allDiscoverySources` in sdk/go/registry/domain.go
- [ ] T020 [US1] Update `AllDiscoverySources()` to return `allDiscoverySources` in sdk/go/registry/domain.go
- [ ] T021 [US1] Update `IsValidDiscoverySource()` to iterate `allDiscoverySources` in sdk/go/registry/domain.go

**PluginStatus Enum (6 values)**:

- [ ] T022 [US1] Create package-level variable `allPluginStatuses` in sdk/go/registry/domain.go
- [ ] T023 [US1] Update `AllPluginStatuses()` to return `allPluginStatuses` in sdk/go/registry/domain.go
- [ ] T024 [US1] Update `IsValidPluginStatus()` to iterate `allPluginStatuses` in sdk/go/registry/domain.go

**SecurityLevel Enum (4 values)**:

- [ ] T025 [US1] Create package-level variable `allSecurityLevels` in sdk/go/registry/domain.go
- [ ] T026 [US1] Update `AllSecurityLevels()` to return `allSecurityLevels` in sdk/go/registry/domain.go
- [ ] T027 [US1] Update `IsValidSecurityLevel()` to iterate `allSecurityLevels` in sdk/go/registry/domain.go

**InstallationMethod Enum (4 values)**:

- [ ] T028 [US1] Create package-level variable `allInstallationMethods` in sdk/go/registry/domain.go
- [ ] T029 [US1] Update `AllInstallationMethods()` to return `allInstallationMethods` in sdk/go/registry/domain.go
- [ ] T030 [US1] Update `IsValidInstallationMethod()` to iterate `allInstallationMethods` in sdk/go/registry/domain.go

**PluginCapability Enum (14 values - largest)**:

- [ ] T031 [US1] Create package-level variable `allPluginCapabilities` in sdk/go/registry/domain.go
- [ ] T032 [US1] Update `AllPluginCapabilities()` to return `allPluginCapabilities` in sdk/go/registry/domain.go
- [ ] T033 [US1] Update `IsValidPluginCapability()` to iterate `allPluginCapabilities` in sdk/go/registry/domain.go

**SystemPermission Enum (9 values)**:

- [ ] T034 [US1] Create package-level variable `allSystemPermissions` in sdk/go/registry/domain.go
- [ ] T035 [US1] Update `AllSystemPermissions()` to return `allSystemPermissions` in sdk/go/registry/domain.go
- [ ] T036 [US1] Update `IsValidSystemPermission()` to iterate `allSystemPermissions` in sdk/go/registry/domain.go

**AuthMethod Enum (6 values)**:

- [ ] T037 [US1] Create package-level variable `allAuthMethods` in sdk/go/registry/domain.go
- [ ] T038 [US1] Update `AllAuthMethods()` to return `allAuthMethods` in sdk/go/registry/domain.go
- [ ] T039 [US1] Update `IsValidAuthMethod()` to iterate `allAuthMethods` in sdk/go/registry/domain.go

### Verification for User Story 1

- [ ] T040 [US1] Run all unit tests with `go test ./sdk/go/registry/...` - verify 100% pass (no behavior changes)
- [ ] T041 [US1] Run benchmarks with `go test -bench=. -benchmem ./sdk/go/registry/` - verify 0 allocs/op
- [ ] T042 [US1] Verify all validation functions perform < 30 ns/op per contract requirements
- [ ] T043 [US1] Run `make test` from repository root to verify no regressions
- [ ] T044 [US1] Run `make lint` from repository root to verify code quality

**Checkpoint**: User Story 1 complete - All 8 enum validations optimized with zero allocations and 4-6x
performance improvement

---

## Phase 4: User Story 2 - Performance Benchmarking and Comparison (Priority: P2)

**Goal**: Document performance characteristics with comprehensive benchmarks comparing optimized implementation
against alternative approaches

**Independent Test**: Benchmark results demonstrate measurable performance improvement and compare slice vs map
approaches across different enum sizes

### Tests for User Story 2 (Performance Benchmarks)

> **NOTE: These are comparison benchmarks to validate optimization decisions**

- [ ] T045 [P] [US2] Add map-based comparison benchmark for Provider in sdk/go/registry/domain_test.go
- [ ] T046 [P] [US2] Add map-based comparison benchmark for PluginCapability (14 values) in
  sdk/go/registry/domain_test.go
- [ ] T047 [P] [US2] Add scalability benchmark with varying enum sizes (4, 6, 9, 14 values) in
  sdk/go/registry/domain_test.go
- [ ] T048 [US2] Run benchmark comparison with `go test -bench=. -benchmem ./sdk/go/registry/` and collect results

### Documentation for User Story 2

- [ ] T049 [US2] Document benchmark results in specs/001-domain-enum-optimization/performance-results.md
- [ ] T050 [US2] Create performance comparison table (before/after, slice vs map) in performance-results.md
- [ ] T051 [US2] Document scalability characteristics for different enum sizes in performance-results.md
- [ ] T052 [US2] Update research.md with actual measured performance vs predictions

**Checkpoint**: Performance characteristics documented and validated against initial research predictions

---

## Phase 5: User Story 3 - Consistent Validation Patterns Across Packages (Priority: P3)

**Goal**: Ensure validation pattern consistency between registry and pricing packages for maintainability

**Independent Test**: Code review confirms both packages use identical validation patterns (optimized
slice-based approach)

### Implementation for User Story 3

- [ ] T053 [US3] Review pricing package validation pattern in sdk/go/pricing/domain.go
- [ ] T054 [US3] Document validation pattern consistency in sdk/go/CLAUDE.md
- [ ] T055 [US3] Update sdk/go/registry/CLAUDE.md with optimization details and pattern guidance
- [ ] T056 [US3] Add pattern documentation to sdk/go/pricing/CLAUDE.md recommending same optimization
- [ ] T057 [US3] Create validation pattern guideline in specs/001-domain-enum-optimization/validation-pattern.md

**Checkpoint**: Validation patterns documented and consistent across packages

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final documentation, cleanup, and verification

- [ ] T058 [P] Update root CLAUDE.md with optimization results and recommendations
- [ ] T059 [P] Update CHANGELOG.md with performance optimization entry
- [ ] T060 [P] Verify all markdown files pass linting with `make lint-markdown`
- [ ] T061 Run complete validation pipeline with `make validate` from repository root
- [ ] T062 Verify quickstart.md examples work correctly
- [ ] T063 Final code review for pattern consistency and code quality

**Checkpoint**: Feature complete and ready for PR

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion - BLOCKS US2 and US3
  - Must complete first as it implements the core optimization
- **User Story 2 (Phase 4)**: Depends on US1 completion (needs optimized implementation to benchmark)
- **User Story 3 (Phase 5)**: Depends on US1 completion (needs optimized pattern to document consistency)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: BLOCKS all other stories - core optimization must be implemented first
- **User Story 2 (P2)**: Can start after US1 - independently testable (benchmarks only)
- **User Story 3 (P3)**: Can start after US1 - independently testable (documentation only)

### Within User Story 1 (Critical Path)

**Test-First Order (TDD)**:

1. T006-T015: Write all benchmarks and unit tests FIRST (must FAIL initially)
2. T016-T039: Implement optimization (8 enum types, 3 tasks each = 24 implementation tasks)
3. T040-T044: Verification (tests now PASS with improved performance)

**Per-Enum Pattern** (can parallelize across enums after tests written):

1. Create package-level variable
2. Update `AllXxx()` function
3. Update `IsValidXxx()` function

### Parallel Opportunities

**Phase 1 Setup**: All tasks (T001-T003) can run in parallel

**Phase 2 Foundational**: T004 and T005 must run sequentially (T005 depends on T004)

**Phase 3 User Story 1**:

- **Tests** (T006-T014): All benchmarks can be written in parallel (different test functions)
- **Implementation** (T016-T039): After tests are written, each enum type (3 tasks) can be implemented in
  parallel
  - Provider (T016-T018) || DiscoverySource (T019-T021) || PluginStatus (T022-T024) || ...
  - Within each enum: 3 tasks must be sequential (variable → AllXxx() → IsValidXxx())

**Phase 4 User Story 2**:

- Benchmarks (T045-T047): Can run in parallel
- Documentation (T049-T052): Must run sequentially after benchmarks

**Phase 5 User Story 3**:

- All tasks (T053-T057) must run sequentially (reviewing and documenting patterns)

**Phase 6 Polish**: Tasks T058-T060 can run in parallel, T061-T063 must run sequentially

---

## Parallel Example: User Story 1 - Writing Tests

```bash
# Launch all benchmark test tasks in parallel (after infrastructure ready):
Task: "Add benchmark for IsValidProvider in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidDiscoverySource in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidPluginStatus in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidSecurityLevel in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidInstallationMethod in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidPluginCapability in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidSystemPermission in sdk/go/registry/domain_test.go"
Task: "Add benchmark for IsValidAuthMethod in sdk/go/registry/domain_test.go"
Task: "Add unit tests for edge cases in sdk/go/registry/domain_test.go"
```

## Parallel Example: User Story 1 - Implementing Enums

```bash
# After tests are written and failing, implement all 8 enum optimizations in parallel:
Task: "Create package-level variable allProviders in sdk/go/registry/domain.go"
Task: "Create package-level variable allDiscoverySources in sdk/go/registry/domain.go"
Task: "Create package-level variable allPluginStatuses in sdk/go/registry/domain.go"
Task: "Create package-level variable allSecurityLevels in sdk/go/registry/domain.go"
Task: "Create package-level variable allInstallationMethods in sdk/go/registry/domain.go"
Task: "Create package-level variable allPluginCapabilities in sdk/go/registry/domain.go"
Task: "Create package-level variable allSystemPermissions in sdk/go/registry/domain.go"
Task: "Create package-level variable allAuthMethods in sdk/go/registry/domain.go"

# Then update all accessor functions in parallel:
Task: "Update AllProviders() to return allProviders in sdk/go/registry/domain.go"
Task: "Update AllDiscoverySources() to return allDiscoverySources in sdk/go/registry/domain.go"
# ... etc for all 8 enums
```

---

## Implementation Strategy

### MVP First (User Story 1 Only - Core Optimization)

1. Complete Phase 1: Setup → Baseline established
2. Complete Phase 2: Foundational → Test infrastructure ready
3. Complete Phase 3: User Story 1 → Core optimization complete
4. **STOP and VALIDATE**: Run benchmarks, verify 0 allocs/op and < 30 ns/op
5. Ready for performance validation (US2) or merge as-is

**MVP Scope**: 44 tasks (T001-T044) - Core optimization with TDD workflow

### Incremental Delivery

1. Complete Setup + Foundational (T001-T005) → Test infrastructure ready
2. Add User Story 1 (T006-T044) → Test independently → **Core optimization complete** (MVP!)
3. Add User Story 2 (T045-T052) → Document performance → **Benchmarks validated**
4. Add User Story 3 (T053-T057) → Document patterns → **Consistency documented**
5. Add Polish (T058-T063) → Final cleanup → **Feature complete**

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T005)
2. For User Story 1 tests (T006-T015):
   - All developers can write benchmarks in parallel (different functions)
3. For User Story 1 implementation (T016-T039):
   - Developer A: Provider, DiscoverySource, PluginStatus (T016-T024)
   - Developer B: SecurityLevel, InstallationMethod, PluginCapability (T025-T033)
   - Developer C: SystemPermission, AuthMethod (T034-T039)
4. Team verifies together (T040-T044)
5. User Story 2 and 3 can be done in parallel by different developers after US1 completes

---

## Task Count Summary

- **Phase 1 (Setup)**: 3 tasks
- **Phase 2 (Foundational)**: 2 tasks
- **Phase 3 (User Story 1)**: 39 tasks (10 tests + 24 implementation + 5 verification)
- **Phase 4 (User Story 2)**: 8 tasks (4 benchmarks + 4 documentation)
- **Phase 5 (User Story 3)**: 5 tasks (pattern documentation)
- **Phase 6 (Polish)**: 6 tasks (final cleanup)

**Total**: 63 tasks

**MVP Scope** (User Story 1): 44 tasks (T001-T044)
**Full Feature**: 63 tasks

---

## Notes

- [P] tasks = different files or test functions, no dependencies
- [Story] label maps task to specific user story for traceability
- TDD workflow enforced: Write tests (T006-T015) → Run and FAIL → Implement (T016-T039) → Verify PASS (T040-T044)
- Each enum type follows identical 3-step pattern: variable → AllXxx() → IsValidXxx()
- All 8 enum types modified in same file (sdk/go/registry/domain.go) but different sections
- Constitution requirement met: Test-First Protocol with benchmark tests before optimization
- Performance target: < 100 ns/op (contract), < 30 ns/op (realistic for 4-14 value enums)
- Memory target: 0 allocs/op (zero allocation validation)
- Backward compatibility: 100% (no API changes, same behavior)
