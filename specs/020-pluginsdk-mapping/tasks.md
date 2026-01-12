# Tasks: PluginSDK Mapping Package

**Input**: Design documents from `/specs/020-pluginsdk-mapping/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Included as per constitution requirement (Test-First Protocol)

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

```text
sdk/go/pluginsdk/mapping/    # Main package directory
├── aws.go                   # AWS extraction functions
├── azure.go                 # Azure extraction functions
├── gcp.go                   # GCP extraction functions
├── common.go                # Generic extraction functions
├── doc.go                   # Package documentation
├── mapping_test.go          # Unit tests
└── benchmark_test.go        # Performance benchmarks
```

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create package directory and documentation structure

- [x] T001 Create sdk/go/pluginsdk/mapping/ directory
- [x] T002 [P] Create doc.go with package documentation in sdk/go/pluginsdk/mapping/doc.go
- [x] T003 [P] Create property key constants in sdk/go/pluginsdk/mapping/keys.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: GCP regions list needed by US3 (GCP extraction)

**Note**: No true blocking prerequisites for US1/US2, but GCP regions validation is
needed before US3 can be implemented.

- [x] T004 [P] Create GCP regions list and IsValidGCPRegion in sdk/go/pluginsdk/mapping/gcp.go
- [x] T005 [P] Create helper function extractFromKeys in sdk/go/pluginsdk/mapping/common.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - AWS Plugin Developer Extracts Properties (Priority: P1)

**Goal**: AWS plugin developers can extract SKU and region from EC2/EBS/RDS properties

**Independent Test**: Pass AWS property maps and verify correct SKU and region extraction

### Tests for User Story 1

- [x] T006 [P] [US1] Write failing tests for ExtractAWSSKU in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T007 [P] [US1] Write failing tests for ExtractAWSRegion in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T008 [P] [US1] Write failing tests for ExtractAWSRegionFromAZ in sdk/go/pluginsdk/mapping/mapping_test.go

### Implementation for User Story 1

- [x] T009 [US1] Implement ExtractAWSRegionFromAZ in sdk/go/pluginsdk/mapping/aws.go
- [x] T010 [US1] Implement ExtractAWSSKU in sdk/go/pluginsdk/mapping/aws.go
- [x] T011 [US1] Implement ExtractAWSRegion in sdk/go/pluginsdk/mapping/aws.go
- [x] T012 [US1] Verify all US1 tests pass and add edge case tests

**Checkpoint**: User Story 1 is complete - AWS extraction fully functional

---

## Phase 4: User Story 2 - Azure Plugin Developer Extracts Properties (Priority: P2)

**Goal**: Azure plugin developers can extract SKU and region from VM properties

**Independent Test**: Pass Azure property maps and verify correct SKU and region extraction

### Tests for User Story 2

- [x] T013 [P] [US2] Write failing tests for ExtractAzureSKU in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T014 [P] [US2] Write failing tests for ExtractAzureRegion in sdk/go/pluginsdk/mapping/mapping_test.go

### Implementation for User Story 2

- [x] T015 [US2] Implement ExtractAzureSKU in sdk/go/pluginsdk/mapping/azure.go
- [x] T016 [US2] Implement ExtractAzureRegion in sdk/go/pluginsdk/mapping/azure.go
- [x] T017 [US2] Verify all US2 tests pass and add edge case tests

**Checkpoint**: User Stories 1 AND 2 complete - AWS and Azure extraction functional

---

## Phase 5: User Story 3 - GCP Plugin Developer Extracts Properties (Priority: P2)

**Goal**: GCP plugin developers can extract SKU and region from Compute properties

**Independent Test**: Pass GCP property maps and verify correct SKU and region extraction
with zone-to-region validation

### Tests for User Story 3

- [x] T018 [P] [US3] Write failing tests for ExtractGCPSKU in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T019 [P] [US3] Write failing tests for ExtractGCPRegion in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T020 [P] [US3] Write failing tests for ExtractGCPRegionFromZone in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T021 [P] [US3] Write failing tests for IsValidGCPRegion in sdk/go/pluginsdk/mapping/mapping_test.go

### Implementation for User Story 3

- [x] T022 [US3] Implement ExtractGCPRegionFromZone in sdk/go/pluginsdk/mapping/gcp.go
- [x] T023 [US3] Implement ExtractGCPSKU in sdk/go/pluginsdk/mapping/gcp.go
- [x] T024 [US3] Implement ExtractGCPRegion in sdk/go/pluginsdk/mapping/gcp.go
- [x] T025 [US3] Verify all US3 tests pass and add edge case tests

**Checkpoint**: User Stories 1, 2, AND 3 complete - All major cloud providers supported

---

## Phase 6: User Story 4 - FinOps Developer Uses Generic Extractors (Priority: P3)

**Goal**: FinOps plugin developers can extract properties using custom key lists

**Independent Test**: Pass custom property maps and verify generic extraction with fallback keys

### Tests for User Story 4

- [x] T026 [P] [US4] Write failing tests for ExtractSKU in sdk/go/pluginsdk/mapping/mapping_test.go
- [x] T027 [P] [US4] Write failing tests for ExtractRegion in sdk/go/pluginsdk/mapping/mapping_test.go

### Implementation for User Story 4

- [x] T028 [US4] Implement ExtractSKU in sdk/go/pluginsdk/mapping/common.go
- [x] T029 [US4] Implement ExtractRegion in sdk/go/pluginsdk/mapping/common.go
- [x] T030 [US4] Verify all US4 tests pass and add edge case tests

**Checkpoint**: All user stories complete - Full multi-cloud and generic support

---

## Phase 7: User Story 5 - Core System Decoupling (Priority: P1)

**Goal**: Documentation and integration guidance for finfocus-core adoption

**Independent Test**: Verify package is importable and documented for core migration

### Implementation for User Story 5

- [x] T031 [US5] Update sdk/go/pluginsdk/README.md with mapping package documentation
- [x] T032 [US5] Add migration guide section for finfocus-core adapter replacement
- [x] T033 [US5] Add usage examples for all extraction functions in doc.go

**Checkpoint**: Package ready for finfocus-core adoption

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Performance benchmarks, documentation, and validation

- [x] T034 [P] Create benchmark tests for all functions in sdk/go/pluginsdk/mapping/benchmark_test.go
- [x] T035 [P] Run benchmarks and verify <50 ns/op, 0 allocs/op targets
- [x] T036 Add nil/empty input edge case tests for all functions
- [x] T037 Run `make lint` and fix any issues
- [x] T038 Run `make test` and verify all tests pass
- [x] T038a Verify package is importable at github.com/rshade/finfocus-spec/sdk/go/pluginsdk/mapping
- [x] T039 Run quickstart.md validation scenarios manually
- [x] T040 Update sdk/go/CLAUDE.md with mapping package documentation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **Foundational (Phase 2)**: Depends on Setup; creates shared helpers
- **User Story 1 (Phase 3)**: Depends on Phase 2; AWS extraction
- **User Story 2 (Phase 4)**: Depends on Phase 2; Azure extraction (can parallel with US1)
- **User Story 3 (Phase 5)**: Depends on Phase 2 (GCP regions list); GCP extraction
- **User Story 4 (Phase 6)**: Depends on Phase 2; Generic extraction (can parallel with US1-3)
- **User Story 5 (Phase 7)**: Depends on US1-4 completion; Documentation
- **Polish (Phase 8)**: Depends on all user stories; Benchmarks and final validation

### User Story Independence

- **US1 (AWS)**: No dependencies on other stories - fully independent
- **US2 (Azure)**: No dependencies on other stories - fully independent
- **US3 (GCP)**: No dependencies on other stories - fully independent (uses shared GCP regions)
- **US4 (Generic)**: No dependencies on other stories - fully independent
- **US5 (Docs)**: Depends on US1-4 for complete documentation

### Within Each User Story

1. Write failing tests FIRST
2. Implement functions to make tests pass
3. Add edge case tests
4. Mark checkpoint complete

### Parallel Opportunities

**Setup Phase**:

- T002, T003 can run in parallel

**Foundational Phase**:

- T004, T005 can run in parallel

**User Stories** (once Foundational complete):

- US1, US2, US3, US4 can all be worked in parallel by different developers
- Within each story: test tasks marked [P] can run in parallel

**Polish Phase**:

- T034, T035 can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task T006: "Write failing tests for ExtractAWSSKU"
Task T007: "Write failing tests for ExtractAWSRegion"
Task T008: "Write failing tests for ExtractAWSRegionFromAZ"

# Then implement sequentially (simpler function first):
Task T009: "Implement ExtractAWSRegionFromAZ"  # No dependencies
Task T010: "Implement ExtractAWSSKU"           # No dependencies
Task T011: "Implement ExtractAWSRegion"        # Uses T009
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T005)
3. Complete Phase 3: User Story 1 (T006-T012)
4. **STOP and VALIDATE**: Test AWS extraction independently
5. Package usable for AWS-only plugins

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add US1 (AWS) → Test → First major cloud supported
3. Add US2 (Azure) → Test → Two clouds supported
4. Add US3 (GCP) → Test → All major clouds supported
5. Add US4 (Generic) → Test → Custom use cases supported
6. Add US5 (Docs) → Complete package ready for adoption
7. Polish → Production-ready with benchmarks

### Suggested MVP Scope

**Minimum**: User Story 1 (AWS extraction only)

- Covers most common use case
- Validates package structure and patterns
- Enables finfocus-core migration for AWS resources

---

## Notes

- All functions must handle nil/empty input without panic
- Use table-driven tests for comprehensive coverage
- Follow existing pluginsdk patterns (see env.go for reference)
- Performance target: <50 ns/op, 0 allocs/op
- GCP regions list last updated: 2025-12
