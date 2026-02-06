# Tasks: GetActualCost Pagination Support

**Input**: Design documents from `/specs/044-actual-cost-pagination/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Included per Constitution V (Test-First Protocol) and spec FR-016.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Proto-First Foundation)

**Purpose**: Update protobuf definitions and regenerate SDK code. This is the contract-first
foundation that all other work depends on.

- [x] T001 Add `page_size` (field 7) and `page_token` (field 8) to `GetActualCostRequest`
  message in `proto/finfocus/v1/costsource.proto`. Add `next_page_token` (field 4) and
  `total_count` (field 5) to `GetActualCostResponse` message. Include proto comments per
  contracts/proto-diff.md.
- [x] T002 Run `make generate` to regenerate Go protobuf bindings in
  `sdk/go/proto/finfocus/v1/costsource.pb.go` and TypeScript bindings.
- [x] T003 Run `buf breaking` to verify no breaking changes to existing proto contract.
  Validate that `make lint` passes for buf lint rules.

**Checkpoint**: Proto contract updated, all SDKs regenerated, backward compatibility
verified.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add logging field and response builder options that all user stories depend on.

- [x] T004 Add `FieldResultCount = "result_count"` constant to
  `sdk/go/pluginsdk/logging.go` following the existing `FieldPageSize` pattern at line 67.
- [x] T005 [P] Add `WithNextPageToken(token string)` option function for
  `ActualCostResponseOption` in `sdk/go/pluginsdk/helpers.go`. Follows the existing
  `WithFallbackHint` and `WithResults` pattern.
- [x] T006 [P] Add `WithTotalCount(count int32)` option function for
  `ActualCostResponseOption` in `sdk/go/pluginsdk/helpers.go`. Follows the existing
  functional options pattern.

**Checkpoint**: Foundation ready. Response builders and logging constants available for
all user stories.

---

## Phase 3: User Story 1 - Paginated Cost Retrieval (Priority: P1) MVP

**Goal**: Hosts can retrieve large cost datasets (10,000+ records) in manageable pages
via the GetActualCost RPC without exceeding gRPC message size limits.

**Independent Test**: Request actual cost data with a specified page size and verify the
response contains at most that many records with a continuation token when more exist.

**FRs**: FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-008, FR-009,
FR-010, FR-011, FR-014, FR-015

### Tests for User Story 1

> Write tests FIRST, ensure they FAIL before implementation (Constitution V).

- [x] T007 [P] [US1] Write `TestPaginateActualCosts_FirstPageDefaultSize` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify that calling with page_size=0 returns
  DefaultPageSize (50) records and a non-empty next token when more records exist.
  Mirror `TestPaginateRecommendations_FirstPageDefaultSize` at line 616.
- [x] T008 [P] [US1] Write `TestPaginateActualCosts_FirstPageCustomSize` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify page_size=100 returns exactly 100 records
  from a 500-record dataset with a valid continuation token.
- [x] T009 [P] [US1] Write `TestPaginateActualCosts_LastPage` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify the last page returns remaining records
  and an empty next_page_token.
- [x] T010 [P] [US1] Write `TestPaginateActualCosts_InvalidToken` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify malformed token returns an error (FR-008).
- [x] T011 [P] [US1] Write `TestPaginateActualCosts_OffsetBeyondRange` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify valid token pointing beyond data returns
  empty results and empty next token (FR-009).
- [x] T012 [P] [US1] Write `TestPaginateActualCosts_MaxPageSizeClamping` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify page_size > 1000 is clamped to
  MaxPageSize (FR-006).
- [x] T013 [P] [US1] Write `TestPaginateActualCosts_TotalCount` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify total_count is auto-populated from
  slice length.
- [x] T014 [P] [US1] Write `TestPaginateActualCosts_SinglePageFitsAll` test in
  `sdk/go/pluginsdk/helpers_test.go`. Verify that when total records < page_size,
  all records returned with empty next token (FR-007).

### Implementation for User Story 1

- [x] T015 [US1] Implement `PaginateActualCosts()` function in
  `sdk/go/pluginsdk/helpers.go`. Accepts `[]*pbc.ActualCostResult`, `pageSize int32`,
  `pageToken string`. Returns `([]*pbc.ActualCostResult, string, int32, error)`.
  Reuses existing `EncodePageToken()`/`DecodePageToken()`. Clamps page_size to
  `[DefaultPageSize, MaxPageSize]`. Auto-populates total_count from `len(results)`.
  Algorithm mirrors `PaginateRecommendations()` at line 930.
- [x] T016 [US1] Add pagination support to mock plugin `GetActualCost` in
  `sdk/go/testing/mock_plugin.go`. Implement `paginateMockActualCosts()` function
  mirroring `paginateMockRecommendations()` at line 1502. Generate configurable
  number of mock `ActualCostResult` records. Wire pagination into mock's
  `GetActualCost` handler to use `page_size` and `page_token` from request.
- [x] T017 [US1] Add paginated GetActualCost integration test in
  `sdk/go/testing/integration_test.go`. Test full gRPC round-trip through
  TestHarness: request with page_size, verify response has correct record count
  and next_page_token. Iterate through all pages and verify total record count
  matches expected. Verify dry_run=true ignores pagination (FR-015).

**Checkpoint**: Core pagination works end-to-end. Hosts can paginate through large
actual cost datasets. Existing plugins return all records in one page (backward
compatible per FR-010, FR-011).

---

## Phase 4: User Story 2 - Backward-Compatible Default Behavior (Priority: P1)

**Goal**: Existing plugins and hosts continue working without any changes after the
protocol update.

**Independent Test**: Run an existing (non-paginated) plugin against a host that sends
pagination parameters and verify the plugin returns all records in a single response.

**FRs**: FR-010, FR-011

### Tests for User Story 2

- [x] T018 [P] [US2] Write `TestBackwardCompat_LegacyPluginNoPagination` test in
  `sdk/go/testing/integration_test.go`. Create a mock plugin that does NOT implement
  pagination logic. Send request with page_size=100 and verify plugin returns all
  records with empty next_page_token.
- [x] T019 [P] [US2] Write `TestBackwardCompat_LegacyHostNoPaginationParams` test in
  `sdk/go/testing/integration_test.go`. Send request with page_size=0 and empty
  page_token to a paginated plugin and verify all records returned (default behavior).

### Implementation for User Story 2

- [x] T020 [US2] Verify that the existing mock plugin `GetActualCost` handler returns
  all results when page_size=0 and page_token="" (proto3 defaults). Ensure no behavioral
  change when pagination fields are absent. Update mock plugin if needed to handle both
  paginated and non-paginated modes in `sdk/go/testing/mock_plugin.go`.
- [x] T021 [US2] Run all existing conformance tests (`go test -v ./sdk/go/testing/
  -run TestConformance`) and integration tests to verify no regressions from proto
  field additions.

**Checkpoint**: All existing tests pass. Legacy plugins and hosts work unchanged.

---

## Phase 5: User Story 3 - SDK Pagination Polish & Observability (Priority: P2)

**Goal**: Plugin developers can implement paginated GetActualCost responses using
ready-to-use SDK helper functions in under 10 lines of pagination code.

**Independent Test**: Call the pagination helper with a slice of cost records, a page
size, and a page token, and verify correct page extraction and next token.

**FRs**: FR-012, FR-014, FR-017

### Tests for User Story 3

- [x] T022 [P] [US3] Write `TestPaginateActualCosts_FullIteration` test in
  `sdk/go/pluginsdk/helpers_test.go`. Iterate through all pages of a 500-record
  dataset with page_size=100 and verify: 5 pages total, each page has correct
  count, no duplicates, no gaps, total_count consistent across pages.
- [x] T023 [P] [US3] Write pagination benchmark `BenchmarkPaginateActualCosts` in
  `sdk/go/pluginsdk/helpers_test.go`. Benchmark with 1000 records, page_size=100.
  Verify <100ms p99 and measure allocations. Target: match
  `PaginateRecommendations` performance.

### Implementation for User Story 3

- [x] T024 [US3] Add structured logging for paginated GetActualCost calls. In
  `sdk/go/pluginsdk/logging.go`, ensure `FieldResultCount` is documented. In a new
  logging example or existing logging helpers, demonstrate logging `page_size` and
  `result_count` fields using zerolog per FR-017 pattern.
- [x] T025 [US3] Verify `PaginateActualCosts` godoc comments include complete
  usage example showing the quickstart pattern from `quickstart.md`. Ensure all
  exported functions have godoc comments in `sdk/go/pluginsdk/helpers.go`.

**Checkpoint**: Plugin developers have a complete, documented, performant pagination
helper. Logging observability matches GetRecommendations pattern.

---

## Phase 6: User Story 4 - Client-Side Page Iteration (Priority: P2)

**Goal**: Host developers can iterate through all pages of paginated actual cost data
using a convenient iterator without manually managing continuation tokens.

**Independent Test**: Create an iterator over a multi-page dataset and verify it yields
all records across all pages in order, then signals completion.

**FRs**: FR-013

### Tests for User Story 4

- [x] T026 [P] [US4] Write `TestActualCostIterator_MultiPage` test in
  `sdk/go/pluginsdk/actual_cost_iterator_test.go`. Create iterator over 5-page
  dataset (500 records, page_size=100). Verify all records yielded in order,
  `Err()` is nil, `TotalCount()` returns expected value.
- [x] T027 [P] [US4] Write `TestActualCostIterator_SinglePage` test in
  `sdk/go/pluginsdk/actual_cost_iterator_test.go`. Verify iterator with fewer
  records than page_size yields all records and signals completion after one fetch.
- [x] T028 [P] [US4] Write `TestActualCostIterator_FetchError` test in
  `sdk/go/pluginsdk/actual_cost_iterator_test.go`. Verify that when fetchFn returns
  an error on page 3, `Next()` returns false and `Err()` surfaces the error.
  Records from pages 1-2 should have been delivered before the error.
- [x] T029 [P] [US4] Write `TestActualCostIterator_EmptyDataset` test in
  `sdk/go/pluginsdk/actual_cost_iterator_test.go`. Verify iterator with 0 records
  returns false on first `Next()` call with nil `Err()`.
- [x] T030 [P] [US4] Write `TestActualCostIterator_ContextCancellation` test in
  `sdk/go/pluginsdk/actual_cost_iterator_test.go`. Verify iterator respects context
  cancellation during fetch.

### Implementation for User Story 4

- [x] T031 [US4] Implement `ActualCostIterator` struct and `NewActualCostIterator()`
  constructor in new file `sdk/go/pluginsdk/actual_cost_iterator.go`. Include
  Apache 2.0 copyright header. Struct fields: ctx, fetchFn, pageSize, pageToken,
  current, index, totalCount, done, err. Constructor accepts `context.Context`,
  fetch function callback, and `int32` page size.
- [x] T032 [US4] Implement `Next()`, `Record()`, `Err()`, `TotalCount()` methods on
  `ActualCostIterator` in `sdk/go/pluginsdk/actual_cost_iterator.go`. `Next()` lazily
  fetches pages on demand. `Record()` returns current `*pbc.ActualCostResult`.
  `Err()` returns first error. `TotalCount()` returns total_count from last response.
- [x] T033 [US4] Add `actualCostIterator()` async generator function to
  `sdk/typescript/packages/client/src/utils/pagination.ts`. Mirror the existing
  `recommendationsIterator()` pattern (lines 6-25). Accept `CostSourceClient` and
  `GetActualCostRequest`. Yield `ActualCostResult` records. Handle `nextPageToken`
  continuation automatically.
- [x] T034 [US4] Add TypeScript iterator tests in
  `sdk/typescript/packages/client/src/utils/pagination.test.ts`. Test multi-page
  iteration, single-page, and error handling for `actualCostIterator()`.

**Checkpoint**: Host developers in both Go and TypeScript can iterate through paginated
actual cost data with a clean, idiomatic API.

---

## Phase 7: User Story 5 - Conformance Validation (Priority: P3)

**Goal**: Conformance test suite validates pagination behavior at Standard level,
ensuring certified plugins handle pagination correctly and consistently.

**Independent Test**: Run conformance suite against a mock plugin that implements
pagination and verify all pagination-specific test cases pass.

**FRs**: FR-016

### Tests for User Story 5

- [x] T035 [P] [US5] Write `TestPaginationConformance_FirstPage` test in new file
  `sdk/go/testing/pagination_conformance_test.go`. Include Apache 2.0 copyright header.
  Verify first page returns correct record count and non-empty next token.
- [x] T036 [P] [US5] Write `TestPaginationConformance_MiddlePage` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify middle page returns correct
  records using continuation token from previous page.
- [x] T037 [P] [US5] Write `TestPaginationConformance_LastPage` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify last page returns remaining
  records and empty next token.
- [x] T038 [P] [US5] Write `TestPaginationConformance_EmptyResult` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify offset beyond data returns
  empty results with empty next token.
- [x] T039 [P] [US5] Write `TestPaginationConformance_InvalidToken` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify malformed token returns
  gRPC InvalidArgument error.
- [x] T040 [P] [US5] Write `TestPaginationConformance_OversizedPage` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify response never contains
  more records than effective page size.
- [x] T041 [P] [US5] Write `TestPaginationConformance_DefaultPageSize` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify page_size=0 uses
  DefaultPageSize (50).
- [x] T042 [P] [US5] Write `TestPaginationConformance_BackwardCompat` test in
  `sdk/go/testing/pagination_conformance_test.go`. Verify legacy plugin (no pagination
  awareness) passes basic conformance with pagination fields present in request.

### Implementation for User Story 5

- [x] T043 [US5] Add pagination benchmarks to `sdk/go/testing/benchmark_test.go`.
  Benchmark `GetActualCost` with pagination through TestHarness: measure per-page
  latency, memory per page with 1000 records, and full iteration through 10,000
  records. Verify <100ms p99 per page.
- [x] T044 [US5] Ensure all 8 conformance tests (T035-T042) pass against the mock
  plugin from T016. Run full conformance suite: `go test -v ./sdk/go/testing/
  -run TestPaginationConformance`. Verify SC-008 (at least 5 distinct behaviors).

**Checkpoint**: Conformance suite validates 8 distinct pagination behaviors. Mock plugin
passes all tests. SC-008 satisfied.

---

## Phase 8: Polish and Cross-Cutting Concerns

**Purpose**: Documentation, validation, and final quality checks.

- [x] T045 [P] Update `sdk/go/pluginsdk/README.md` with pagination documentation.
  Add `PaginateActualCosts` section under the existing pagination helpers. Add
  `ActualCostIterator` section with usage examples from quickstart.md. Document
  `WithNextPageToken` and `WithTotalCount` response options.
- [x] T046 [P] Verify godoc coverage is >80% for all modified packages. Run
  `go doc ./sdk/go/pluginsdk/` and verify all new exported functions
  (`PaginateActualCosts`, `NewActualCostIterator`, `ActualCostIterator`,
  `WithNextPageToken`, `WithTotalCount`, `FieldResultCount`) have documentation.
- [x] T047 [P] Run `make validate` to execute full validation pipeline: tests,
  linting (Go, buf, markdown, YAML), and npm validations. Fix any issues.
- [x] T048 [P] Run `make lint` with extended timeout to verify golangci-lint passes
  on all modified Go files.
- [x] T049 [P] Run performance benchmarks: `go test -bench=. -benchmem
  ./sdk/go/pluginsdk/` and `go test -bench=. -benchmem ./sdk/go/testing/`.
  Verify PaginateActualCosts performance matches PaginateRecommendations baseline.
  Verify SC-002 (10k records in <10s) and SC-003 (<100MB per max page).
- [x] T050 Run quickstart.md validation: verify all Go code examples in
  `specs/044-actual-cost-pagination/quickstart.md` compile against the updated SDK.
- [x] T051 Verify TypeScript SDK builds: `cd sdk/typescript && npm run build` and
  `npm test` to ensure `actualCostIterator()` compiles and tests pass.

---

## Dependencies and Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies. Start immediately.
- **Phase 2 (Foundational)**: Depends on Phase 1 (needs regenerated proto).
- **Phase 3 (US1 - P1)**: Depends on Phase 2. Core MVP.
- **Phase 4 (US2 - P1)**: Depends on Phase 3 (needs pagination implementation to
  verify backward compat).
- **Phase 5 (US3 - P2)**: Depends on Phase 3 (extends PaginateActualCosts with
  docs and benchmarks).
- **Phase 6 (US4 - P2)**: Depends on Phase 2 only (iterator is independent of
  server-side helper). Can run in parallel with Phase 3-5.
- **Phase 7 (US5 - P3)**: Depends on Phase 3 and Phase 4 (needs working mock with
  pagination).
- **Phase 8 (Polish)**: Depends on all user story phases being complete.

### User Story Dependencies

```text
Phase 1 (Setup)
    └── Phase 2 (Foundational)
         ├── Phase 3 (US1: Core Pagination) ── MVP
         │    ├── Phase 4 (US2: Backward Compat)
         │    ├── Phase 5 (US3: SDK Helpers Polish)
         │    └── Phase 7 (US5: Conformance)
         └── Phase 6 (US4: Client Iterator) ── parallel with US1
              └── Phase 8 (Polish) ── after all stories
```

### Within Each User Story

1. Tests FIRST (must fail before implementation per Constitution V)
2. Core implementation
3. Integration verification
4. Story checkpoint validation

### Parallel Opportunities

- **Phase 2**: T005 and T006 can run in parallel (different option functions)
- **Phase 3 tests**: T007-T014 can all run in parallel (independent test functions)
- **Phase 4 tests**: T018 and T019 can run in parallel
- **Phase 5 tests**: T022 and T023 can run in parallel
- **Phase 6 tests**: T026-T030 can all run in parallel
- **Phase 6 impl**: T033-T034 (TypeScript) can run in parallel with T031-T032 (Go)
- **Phase 7 tests**: T035-T042 can all run in parallel
- **Phase 8**: T045-T049 can all run in parallel
- **Cross-phase**: Phase 6 (US4) can run in parallel with Phases 3-5

---

## Parallel Example: User Story 1

```text
# Launch all US1 tests in parallel (T007-T014):
Task: "Write TestPaginateActualCosts_FirstPageDefaultSize in helpers_test.go"
Task: "Write TestPaginateActualCosts_FirstPageCustomSize in helpers_test.go"
Task: "Write TestPaginateActualCosts_LastPage in helpers_test.go"
Task: "Write TestPaginateActualCosts_InvalidToken in helpers_test.go"
Task: "Write TestPaginateActualCosts_OffsetBeyondRange in helpers_test.go"
Task: "Write TestPaginateActualCosts_MaxPageSizeClamping in helpers_test.go"
Task: "Write TestPaginateActualCosts_TotalCount in helpers_test.go"
Task: "Write TestPaginateActualCosts_SinglePageFitsAll in helpers_test.go"
```

## Parallel Example: User Story 4

```text
# Launch Go and TypeScript iterator work in parallel:
Task: "Implement ActualCostIterator in actual_cost_iterator.go"      # Go
Task: "Add actualCostIterator() in pagination.ts"                    # TypeScript
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Proto contract update
2. Complete Phase 2: Response builder options and logging
3. Complete Phase 3: Core pagination (US1)
4. Complete Phase 4: Backward compatibility verification (US2)
5. **STOP and VALIDATE**: Run `make validate` and `make test`. All existing tests
   pass, pagination works end-to-end.
6. This delivers SC-001, SC-004, SC-005, SC-007.

### Incremental Delivery

1. Setup + Foundational -> Proto contract ready
2. Add US1 (Core Pagination) -> Test independently -> MVP complete
3. Add US2 (Backward Compat) -> Verify no regressions -> Deploy-ready
4. Add US3 (SDK Helpers Polish) -> Better DX for plugin devs
5. Add US4 (Client Iterator) -> Better DX for host devs (can parallel with 2-3)
6. Add US5 (Conformance) -> Ecosystem quality assurance
7. Polish -> Documentation, benchmarks, final validation
8. Each story adds value without breaking previous stories.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Tests written first per Constitution V (Test-First Protocol)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All file paths are relative to repository root
