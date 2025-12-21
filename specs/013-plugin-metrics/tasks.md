# Tasks: Standardized Plugin Metrics

**Input**: Design documents from `/specs/013-plugin-metrics/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Constitution requires Test-First Protocol - tests are included per spec requirements.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

All source code in `sdk/go/pluginsdk/` per plan.md structure decision.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add Prometheus dependency and prepare metrics module structure

- [x] T001 Add prometheus/client_golang dependency to go.mod via `go get github.com/prometheus/client_golang`
- [x] T002 Run `go mod tidy` to update go.sum
- [x] T003 [P] Create metrics.go file skeleton with package declaration and imports in
      sdk/go/pluginsdk/metrics.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define constants and types that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Define MetricNamespace, MetricSubsystem, DefaultMetricsPort, DefaultMetricsPath constants
      in sdk/go/pluginsdk/metrics.go
- [x] T005 Define DefaultHistogramBuckets variable (5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s,
      2.5s, 5s) in sdk/go/pluginsdk/metrics.go
- [x] T006 Define PluginMetrics struct with RequestsTotal, RequestDuration, and Registry fields in
      sdk/go/pluginsdk/metrics.go
- [x] T007 Define MetricsServerConfig struct with Port, Path, and Registry fields in
      sdk/go/pluginsdk/metrics.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Enable Metrics Collection for Plugin (Priority: P1) 🎯 MVP

**Goal**: Plugin maintainers can enable metrics collection via a single interceptor configuration

**Independent Test**: Create plugin with metrics interceptor, make gRPC requests, verify counter
increments and histogram records latency

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T008 [P] [US1] Create metrics_test.go with TestNewPluginMetrics verifying registry creation
      in sdk/go/pluginsdk/metrics_test.go
- [x] T009 [P] [US1] Add TestMetricsUnaryServerInterceptor_CounterIncrement verifying counter
      increments on request in sdk/go/pluginsdk/metrics_test.go
- [x] T010 [P] [US1] Add TestMetricsUnaryServerInterceptor_HistogramObservation verifying duration
      recorded in sdk/go/pluginsdk/metrics_test.go
- [x] T011 [P] [US1] Add TestMetricsUnaryServerInterceptor_ErrorHandling verifying grpc_code label
      for errors in sdk/go/pluginsdk/metrics_test.go
- [x] T012 [P] [US1] Add integration test using bufconn harness verifying MetricsInterceptor
      chains correctly with TracingUnaryServerInterceptor (both interceptors active, trace_id and
      metrics both recorded) in sdk/go/pluginsdk/metrics_test.go
- [x] T012a [P] [US1] Add TestNoMetricsOverhead_InterceptorNotConfigured verifying no metrics
      code executes when interceptor is not added to server chain in sdk/go/pluginsdk/metrics_test.go

### Implementation for User Story 1

- [x] T013 [US1] Implement NewPluginMetrics(pluginName string) creating registry, counter, histogram
      in sdk/go/pluginsdk/metrics.go
- [x] T014 [US1] Implement MetricsUnaryServerInterceptor(pluginName string) creating default metrics
      and returning interceptor in sdk/go/pluginsdk/metrics.go
- [x] T015 [US1] Implement MetricsInterceptorWithRegistry(metrics \*PluginMetrics) for custom registry
      usage in sdk/go/pluginsdk/metrics.go
- [x] T016 [US1] Add method name extraction using path.Base(info.FullMethod) for grpc_method label
      in sdk/go/pluginsdk/metrics.go
- [x] T017 [US1] Add gRPC status code extraction using status.Code(err).String() for grpc_code label
      in sdk/go/pluginsdk/metrics.go
- [x] T018 [US1] Verify all tests pass with `go test -v ./sdk/go/pluginsdk/ -run TestMetrics`

**Checkpoint**: User Story 1 complete - metrics interceptor functional, counter and histogram work

---

## Phase 4: User Story 2 - Query Metrics via Standard Endpoint (Priority: P2)

**Goal**: Operations engineers can scrape metrics from HTTP endpoint for monitoring integration

**Independent Test**: Start plugin with metrics, make requests, query /metrics endpoint, verify
Prometheus format output with accurate counts

### Tests for User Story 2

- [x] T019 [P] [US2] Add TestStartMetricsServer_DefaultConfig verifying server starts on default
      port in sdk/go/pluginsdk/metrics_test.go
- [x] T020 [P] [US2] Add TestStartMetricsServer_CustomPort verifying configurable port in
      sdk/go/pluginsdk/metrics_test.go
- [x] T021 [P] [US2] Add TestStartMetricsServer_MetricsEndpoint verifying /metrics returns
      Prometheus format in sdk/go/pluginsdk/metrics_test.go
- [x] T022 [P] [US2] Add TestStartMetricsServer_Shutdown verifying graceful shutdown in
      sdk/go/pluginsdk/metrics_test.go

### Implementation for User Story 2

- [x] T023 [US2] Implement StartMetricsServer(config MetricsServerConfig) returning \*http.Server in
      sdk/go/pluginsdk/metrics.go
- [x] T024 [US2] Add default port handling (9090) and path handling (/metrics) in StartMetricsServer
      in sdk/go/pluginsdk/metrics.go
- [x] T025 [US2] Use promhttp.HandlerFor(registry, opts) for metrics endpoint handler in
      sdk/go/pluginsdk/metrics.go
- [x] T026 [US2] Add server startup in goroutine with error channel for startup errors in
      sdk/go/pluginsdk/metrics.go
- [x] T027 [US2] Verify all tests pass with `go test -v ./sdk/go/pluginsdk/ -run TestStartMetrics`

**Checkpoint**: User Stories 1 AND 2 complete - full metrics pipeline working

---

## Phase 5: User Story 3 - Identify Plugin Performance Issues (Priority: P3)

**Goal**: Plugin maintainers can see per-method latency distributions for performance analysis

**Independent Test**: Make requests to different gRPC methods (GetProjectedCost, GetActualCost),
verify histogram has separate buckets per grpc_method label

### Tests for User Story 3

- [x] T028 [P] [US3] Add TestMetrics_PerMethodLabels verifying grpc_method label distinguishes
      methods in sdk/go/pluginsdk/metrics_test.go
- [x] T029 [P] [US3] Add TestMetrics_AllGRPCMethods verifying Name, Supports, GetProjectedCost,
      GetActualCost, GetPricingSpec, EstimateCost tracked in sdk/go/pluginsdk/metrics_test.go

### Implementation for User Story 3

- [x] T030 [US3] Verify histogram correctly uses grpc_method label from interceptor (already
      implemented in T016) - add method-specific integration test in sdk/go/pluginsdk/metrics_test.go
- [x] T031 [US3] Add comprehensive integration test calling all 6 gRPC methods and verifying
      distinct histogram entries in sdk/go/pluginsdk/metrics_test.go
- [x] T031a [US3] Add TestMetrics_CountAccuracy sending exactly 1000 requests and verifying
      counter shows 1000 (within 1% = 990-1010) in sdk/go/pluginsdk/metrics_test.go

**Checkpoint**: All user stories complete - full metrics with per-method breakdown

---

## Phase 6: Performance & Benchmarks

**Purpose**: Verify <5% overhead requirement (SC-004 from spec)

- [x] T032 [P] Create metrics_benchmark_test.go with BenchmarkMetricsInterceptor_Overhead in
      sdk/go/pluginsdk/metrics_benchmark_test.go
- [x] T033 [P] Add BenchmarkMetricsInterceptor_NoMetrics baseline (interceptor disabled) in
      sdk/go/pluginsdk/metrics_benchmark_test.go
- [x] T034 Compare benchmark results to verify <5% overhead with
      `go test -bench=. -benchmem ./sdk/go/pluginsdk/`
- [x] T034a Add TestMetrics_LatencyAccuracy with controlled 100ms sleep, verifying recorded
      duration is 100ms ± 1ms in sdk/go/pluginsdk/metrics_test.go
- [x] T035 Document benchmark results in README.md metrics section

---

## Phase 7: Polish & Documentation

**Purpose**: Documentation updates and final validation

- [x] T036 [P] Add Metrics section to sdk/go/pluginsdk/README.md with usage examples from
      quickstart.md
- [x] T037 [P] Add inline godoc comments to all exported functions in
      sdk/go/pluginsdk/metrics.go
- [x] T038 Run `make lint` to verify code style compliance
- [x] T039 Run `make test` to verify all tests pass
- [x] T040 Validate quickstart.md examples compile and work correctly
- [x] T041 Update CLAUDE.md Active Technologies section if needed

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
- **Performance (Phase 6)**: Depends on User Story 1 (metrics interceptor must exist)
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Uses PluginMetrics.Registry from
  US1 but independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Validates US1 implementation but
  independently testable

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD per constitution)
- Types/structs before functions
- Core implementation before integration tests
- Story complete before moving to next priority

### Parallel Opportunities

- T003 (metrics.go skeleton) can run in parallel with T001-T002
- All T008-T012 tests can run in parallel (different test functions)
- All T019-T022 tests can run in parallel
- T028-T029 tests can run in parallel
- T032-T033 benchmarks can run in parallel
- T036-T037 documentation can run in parallel

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all tests for User Story 1 together:
Task: "Create metrics_test.go with TestNewPluginMetrics in sdk/go/pluginsdk/metrics_test.go"
Task: "Add TestMetricsUnaryServerInterceptor_CounterIncrement in sdk/go/pluginsdk/metrics_test.go"
Task: "Add TestMetricsUnaryServerInterceptor_HistogramObservation in sdk/go/pluginsdk/metrics_test.go"
Task: "Add TestMetricsUnaryServerInterceptor_ErrorHandling in sdk/go/pluginsdk/metrics_test.go"
Task: "Add integration test using bufconn harness in sdk/go/pluginsdk/metrics_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (add dependency)
2. Complete Phase 2: Foundational (constants, types)
3. Complete Phase 3: User Story 1 (core interceptor)
4. **STOP and VALIDATE**: Test interceptor independently
5. Plugin maintainers can now use `MetricsUnaryServerInterceptor("my-plugin")`

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Metrics collection works (MVP!)
3. Add User Story 2 → Test independently → HTTP endpoint works
4. Add User Story 3 → Test independently → Per-method analysis works
5. Add Performance → Verify <5% overhead
6. Polish → Documentation complete

### Single Developer Strategy

Execute phases sequentially: 1 → 2 → 3 → 4 → 5 → 6 → 7

Each phase delivers testable value. Stop at any checkpoint to validate.

---

## Notes

- [P] tasks = different files or test functions, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- Verify tests fail before implementing (TDD requirement)
- Commit after each task or logical group
- New dependency: github.com/prometheus/client_golang
- All code in existing sdk/go/pluginsdk/ package (no new packages)
