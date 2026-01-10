# Implementation Tasks: v0.4.14 SDK Polish Release

**Branch**: `001-sdk-polish-release`
**Date**: 2025-01-04
**Input**: Tasks generated from plan.md, spec.md, data-model.md, contracts/, research.md, quickstart.md

## Phase 1: Setup

**Goal**: Verify development environment and prepare for implementation

- [x] T001 Verify Go 1.25.5 version and gRPC/protobuf dependencies installed
- [x] T002 Run `make generate` to ensure proto generation works
- [x] T003 Run `make test` to establish baseline test pass rate
- [x] T004 Verify golangci-lint is configured and can run `golangci-lint run`
- [x] T004a [P] Add ClientConfig.Timeout precedence verification tests in sdk/go/pluginsdk/client_test.go
  (context deadline > ClientConfig.Timeout > default 30s)

**Independent Test Criteria**: All environment verification tasks pass, proto generation succeeds, linting works, and
baseline tests pass

---

## Phase 2: Foundational

**Goal**: Prepare cross-cutting infrastructure for all user stories

- [x] T005 [P] Create sdk/go/pluginsdk/health.go file with HealthChecker interface and HealthStatus struct
- [x] T006 [P] Create sdk/go/pluginsdk/context.go file with ValidateContext, ContextRemainingTime,
  ContextDeadline helpers
- [x] T007 [P] Create sdk/go/pluginsdk/arn.go file with DetectARNProvider, ValidateARNConsistency helpers and ARN
  pattern constants
- [x] T008 [P] Add ClientConfig.Timeout field to ClientConfig struct in sdk/go/pluginsdk/client.go

**Independent Test Criteria**: All new SDK files compile successfully and new constants/interfaces are accessible

---

## Phase 3: User Story 1 - Plugin Development Experience (Priority: P1)

**Goal**: As a plugin developer, I want access to health checking, context validation, and ARN format helpers so I can
build robust plugins more efficiently with better error handling

**Why this priority**: These developer experience improvements directly impact plugin quality and reduce development
time across all plugins. Custom health checking is critical for production deployments.

**Independent Test**: Can be fully tested by implementing a sample plugin with custom health checker and ARN
validation, demonstrating that context errors are caught early and health status is accurately reported.

### Tests (TDD - Write First)

- [x] T009 [US1] Write failing tests for ValidateContext in sdk/go/pluginsdk/context_test.go (nil context, cancelled
  context, valid context)
- [x] T010 [US1] Write failing tests for ContextRemainingTime and ContextDeadline in sdk/go/pluginsdk/context_test.go
- [x] T011 [US1] Write failing tests for DetectARNProvider in sdk/go/pluginsdk/arn_test.go (AWS, Azure, GCP, Kubernetes,
  unknown formats)
- [x] T012 [US1] Write failing tests for ValidateARNConsistency in sdk/go/pluginsdk/arn_test.go (valid match,
  mismatch, unrecognized format)
- [x] T013 [US1] Write failing tests for HealthChecker interface in sdk/go/pluginsdk/health_test.go (returns nil,
  returns error, times out, panics)

### Implementation

- [x] T014 [US1] Implement ValidateContext function in sdk/go/pluginsdk/context.go with nil and
  expired/cancelled context checks
- [x] T015 [US1] Implement ContextRemainingTime function in sdk/go/pluginsdk/context.go returning time until
  deadline or MaxInt64 duration
- [x] T016 [US1] Implement ContextDeadline function in sdk/go/pluginsdk/context.go returning context deadline or zero
  time
- [x] T017 [US1] Implement DetectARNProvider function in sdk/go/pluginsdk/arn.go using prefix/pattern matching (AWS,
  Azure, GCP, Kubernetes)
- [x] T018 [US1] Implement ValidateARNConsistency function in sdk/go/pluginsdk/arn.go comparing detected provider to
  expected provider
- [x] T019 [US1] Add ARN pattern constants (AWSARNPrefix, AzureARNPrefix, GCPARNPrefix, KubernetesFormat)
  in sdk/go/pluginsdk/arn.go
- [x] T020 [US1] Implement HealthChecker interface in sdk/go/pluginsdk/health.go with Check(ctx) error method
- [x] T021 [US1] Implement HealthStatus struct in sdk/go/pluginsdk/health.go with Healthy, Message, Details,
  LastChecked fields

### Integration

- [x] T022 [US1] Add type assertion check for HealthChecker in sdk/go/pluginsdk/sdk.go Serve() function
- [x] T023 [US1] Implement HealthHandler in sdk/go/pluginsdk/sdk.go to use HealthChecker.Check() if implemented
- [x] T024 [US1] Update HTTP /healthz endpoint in sdk/go/pluginsdk/sdk.go to call custom HealthChecker if present
- [x] T025 [US1] Update gRPC health service in sdk/go/pluginsdk/sdk.go to call custom HealthChecker if present
- [x] T026 [US1] Add panic recovery for HealthChecker.Check() in sdk/go/pluginsdk/sdk.go (return
  HTTP 503 / gRPC Unavailable)
- [x] T027 [US1] Add context timeout for HealthChecker.Check() in sdk/go/pluginsdk/sdk.go health handler
- [x] T028 [US1] Populate HealthStatus.LastChecked timestamp in sdk/go/pluginsdk/sdk.go health handler

**Independent Test Criteria**: Sample plugin implementing HealthChecker passes health check correctly, context
validation catches nil/expired contexts, ARN detection works for all providers, and all tests pass with -race flag

---

## Phase 4: User Story 2 - Plugin Information & Discovery (Priority: P1)

**Goal**: As a plugin operator, I want to retrieve plugin metadata quickly and reliably with clear error messages so I
can understand plugin capabilities and troubleshoot issues efficiently

**Why this priority**: GetPluginInfo is a core RPC that clients rely on for discovery. Performance and clear error
messages are critical for production operations.

**Independent Test**: Can be fully tested by calling GetPluginInfo on various plugins (new and legacy) and verifying
response time, error message clarity, and metadata accuracy.

### Tests (TDD - Write First)

- [x] T029 [US2] Write failing GetPluginInfoPerformance conformance test in sdk/go/testing/conformance_test.go
  (<100ms requirement)
- [x] T030 [US2] Write failing tests for GetPluginInfo error message mapping in sdk/go/pluginsdk/sdk_test.go (nil
  response, incomplete metadata, invalid spec_version)
- [x] T031 [US2] Write failing tests for GetPluginInfoProvider interface detection in sdk/go/pluginsdk/sdk_test.go
  (implemented vs legacy)

### Implementation

- [x] T032 [US2] Create GetPluginInfoProvider interface in sdk/go/pluginsdk/sdk.go with GetPluginInfo() method
- [x] T033 [US2] Implement GetPluginInfo RPC handler in sdk/go/pluginsdk/sdk.go with interface type assertion
- [x] T034 [US2] Update GetPluginInfo error messages in sdk/go/pluginsdk/sdk.go (nil → "unable to retrieve plugin
  metadata", etc.)
- [x] T035 [US2] Add structured logging for GetPluginInfo errors in sdk/go/pluginsdk/sdk.go (log technical details
  server-side)
- [x] T036 [US2] Return Unimplemented status for legacy plugins in sdk/go/pluginsdk/sdk.go GetPluginInfo handler
- [x] T037 [US2] Add GetPluginInfoPerformance conformance test to Basic conformance level in
  sdk/go/testing/conformance_test.go

### Documentation

- [x] T038 [US2] Add "Migrating to GetPluginInfo" section to sdk/go/pluginsdk/README.md
- [x] T039 [US2] Add static metadata example (NewPluginInfo helper) to sdk/go/pluginsdk/README.md migration guide
- [x] T040 [US2] Add dynamic metadata example (GetPluginInfoProvider interface) to
  sdk/go/pluginsdk/README.md migration guide
- [x] T041 [US2] Add backward compatibility guidance (Unimplemented status handling) to
  sdk/go/pluginsdk/README.md migration guide
- [x] T042 [US2] Add code examples for WithProviders, WithDescription options to
  sdk/go/pluginsdk/README.md migration guide

**Independent Test Criteria**: GetPluginInfo completes within 100ms for all iterations in performance test, error
messages are user-friendly (no internal details), legacy plugins return Unimplemented status, and migration guide
is complete with examples

---

## Phase 5: User Story 3 - Connect Protocol Robustness (Priority: P2)

**Goal**: As a plugin developer, I want the Connect protocol to handle concurrent requests, large payloads, and
connection resets gracefully so my plugin remains stable under real-world conditions

**Why this priority**: Connect protocol is used for client communication. Stability issues directly affect plugin
reliability in production.

**Independent Test**: Can be fully tested by running concurrent requests, large payload transfers, and connection
reset scenarios against a test plugin, verifying no panics or data corruption occur.

### Tests (TDD - Write First)

- [x] T043 [US3] Write failing test for concurrent request handling (100+ requests) in sdk/go/pluginsdk/connect_test.go
- [x] T123 [US3] Write failing test for large request payloads (>1MB) in
  sdk/go/pluginsdk/connect_test.go
- [x] T045 [US3] Write failing test for payload size rejection (>1MB) in sdk/go/pluginsdk/connect_test.go
- [x] T046 [US3] Write failing test for graceful shutdown during active requests in sdk/go/pluginsdk/connect_test.go
- [x] T047 [US3] Write failing test for connection reset handling in sdk/go/pluginsdk/connect_test.go

### Implementation

- [x] T048 [US3] Add context deadline support to client RPC methods in sdk/go/pluginsdk/client.go (respect context
  deadlines)
- [x] T049 [US3] Add per-request timeout configuration via context.WithTimeout in sdk/go/pluginsdk/client.go
  documentation
- [x] T050 [US3] Add ClientConfig.Timeout option for overriding default 30-second timeout in sdk/go/pluginsdk/client.go
- [x] T051 [US3] Implement timeout error handling with appropriate gRPC status codes in
  sdk/go/pluginsdk/client.go
- [x] T052 [US3] Add 1MB payload size limit validation in sdk/go/pluginsdk/connect.go
- [x] T053 [US3] Add explicit error messages for payloads exceeding 1MB in
  sdk/go/pluginsdk/connect.go
- [x] T054 [US3] Implement graceful shutdown handling for active requests in sdk/go/pluginsdk/connect.go
- [x] T055 [US3] Add connection reset handling and cleanup in sdk/go/pluginsdk/connect.go
- [x] T056 [US3] Add concurrent request handling tests (run with -race flag) in sdk/go/pluginsdk/connect_test.go
- [x] T057 [US3] Add large payload streaming/chunking tests in sdk/go/pluginsdk/connect_test.go
- [x] T058 [US3] Add graceful shutdown scenario tests in sdk/go/pluginsdk/connect_test.go
- [x] T059 [US3] Add connection reset scenario tests in sdk/go/pluginsdk/connect_test.go

### CORS Complexity Reduction

- [x] T060 [US3] Verify Serve() function cognitive complexity in sdk/go/pluginsdk/sdk.go is <20
- [x] T061 [US3] Add unit tests for validateCORSConfig() edge cases in sdk/go/pluginsdk/sdk_test.go (nil config,
  valid config, invalid config)
- [x] T062 [US3] Verify CORS behavior remains functionally unchanged after refactoring in sdk/go/pluginsdk/sdk.go

**Independent Test Criteria**: 100+ concurrent requests handled safely without race conditions, large payloads (>1MB)
processed correctly, graceful shutdown completes in-flight requests, connection resets don't cause panics,
Serve() complexity <20, and all Connect tests pass with -race flag

---

## Phase 6: User Story 4 - Testing Infrastructure & Quality (Priority: P3)

**Goal**: As a plugin maintainer, I want stable CI benchmarks, comprehensive edge case coverage, and fuzz testing tools
so I can ensure plugin quality and catch regressions early

**Why this priority**: Improves long-term maintainability and quality of the SDK and plugins using it. Less critical
for immediate plugin functionality but valuable for ecosystem health.

**Independent Test**: Can be fully tested by running the test suite and verifying benchmarks don't fail spuriously,
extreme value tests catch edge cases, and fuzz tests discover potential issues.

### CI Benchmark Stability

- [x] T063 [US4] Update benchmark alert threshold to 150% in .github/workflows/benchmarks.yml
- [x] T064 [US4] Set fail-on-alert to false in .github/workflows/benchmarks.yml
- [x] T065 [US4] Enable comment-on-alert in .github/workflows/benchmarks.yml
- [x] T066 [US4] Add documentation for expected CI variance in README.md or AGENTS.md

### Extreme Value Testing

- [x] T067 [US4] Write failing tests for IEEE 754 special values (infinity, NaN) in
  sdk/go/pluginsdk/focus_conformance_test.go
- [x] T068 [US4] Write failing tests for max/min valid float64 values in sdk/go/pluginsdk/focus_conformance_test.go
- [x] T069 [US4] Implement infinity/NaN validation and rejection in cost validation functions in sdk/go/pluginsdk/
- [x] T070 [US4] Add clear error messages for infinity/NaN rejection in sdk/go/pluginsdk/ cost validation

### Fuzz Testing

- [x] T071 [US4] Write FuzzResourceDescriptorID fuzz test with diverse seed corpus in
  sdk/go/pluginsdk/helpers_test.go
- [x] T073 [US4] Implement fuzz test logic (no panics, ID round-trips correctly) in
  sdk/go/pluginsdk/helpers_test.go
- [x] T074 [US4] Add CI job for short fuzz test runs in .github/workflows/benchmarks.yml or test.yml

### Test Coverage & Quality

- [x] T075 [US4] Run code coverage analysis and ensure >80% coverage maintained in sdk/go/pluginsdk/
- [x] T076 [US4] Run golangci-lint and ensure all linters pass in sdk/go/pluginsdk/
- [x] T077 [US4] Run all tests with -race flag to verify no data races in sdk/go/pluginsdk/

**Independent Test Criteria**: Benchmarks generate alerts but don't block PRs, extreme value tests reject infinity/NaN
with clear errors, fuzz test runs 60+ minutes without panics, code coverage >80%, all linters pass, and -race
tests pass

---

## Phase 7: Polish & Cross-Cutting Concerns

**Goal**: Final validation, documentation updates, and release readiness

- [x] T078 Add context validation examples to sdk/go/pluginsdk/README.md
- [x] T079 Add ARN detection/validation examples to sdk/go/pluginsdk/README.md
- [x] T080 Add HealthChecker interface documentation to sdk/go/pluginsdk/README.md
- [x] T081 Add ClientConfig.Timeout usage examples to sdk/go/pluginsdk/README.md
- [x] T082 Run `make lint` and ensure no linting errors
- [x] T083 Run `make test` and ensure all tests pass
- [x] T084 Run `make validate` if available to ensure all validation layers pass
- [x] T085 Verify all 12 issues have corresponding tests that pass
- [x] T086 Verify code coverage >80% across all SDK packages
- [x] T087 Verify Serve() cognitive complexity <20
- [x] T088 Run `markdownlint *.md` on all modified markdown files
- [x] T089 Update CHANGELOG.md with v0.4.14 entry summarizing 12 issues across 4 themes

**Independent Test Criteria**: All validation passes, documentation is complete and consistent, all tests pass, and
release is ready for merge

---

## Dependencies

### User Story Dependencies

```text
Phase 3 (US1 - Plugin Development Experience, P1)
  ├─ Foundational (T005-T008) - New files and ClientConfig field
  └─ Independent (no other story dependencies)

Phase 4 (US2 - Plugin Information & Discovery, P1)
  ├─ Foundational (T005-T008) - SDK structure ready
  └─ Independent (no other story dependencies)

Phase 5 (US3 - Connect Protocol Robustness, P2)
  ├─ Foundational (T008) - ClientConfig.Timeout
  └─ Phase 3 (US1) - Context validation helpers (T014-T016) useful for timeout handling

Phase 6 (US4 - Testing Infrastructure & Quality, P3)
- Phase 3 (US1) - Context/ARN helpers have tests
- Phase 4 (US2) - GetPluginInfo performance test
- Phase 5 (US3) - Connect protocol tests
- Dependent on all implementation phases (tests implementation of US1-US3)

Phase 7 (Polish)
  └─ Dependent on Phases 3-6 (polishes all implemented features)
```

### Parallel Execution Opportunities

**Phase 2 (Foundational)**:

- T005, T006, T007 can run in parallel [P] (different files: health.go, context.go, arn.go)

**Phase 3 (US1 Tests)**:

- T009, T010 can run in parallel [P] (different test functions in same file, no dependencies)
- T011, T012 can run in parallel [P] (different test functions in same file, no dependencies)

**Phase 3 (US1 Implementation)**:

- T014, T015, T016 can run in parallel [P] (different functions in context.go)
- T017, T018, T019 can run in parallel [P] (different functions in arn.go)

**Phase 4 (US2 Tests)**:

- T029, T030, T031 can run in parallel [P] (different test files: conformance_test.go, sdk_test.go)

**Phase 4 (US2 Documentation)**:

- T038-T042 can run in parallel [P] (different sections of README.md)

**Phase 6 (US4)**:

- T063-T066 (CI benchmarks) can run in parallel [P] (different YAML/README changes)
- T067-T070 (extreme values) can run in parallel [P] (different tests)
- T071-T074 (fuzz testing) can run in parallel [P] (different test functions)

---

## Implementation Strategy

### MVP Scope (Minimum Viable Product)

**Recommended MVP**: Complete Phase 3 (US1 - Plugin Development Experience) + Phase 4 (US2 - Plugin Information &
Discovery)

**Rationale**:

- Both P1 priority stories
- Directly improve plugin developer experience
- Enable new features (health checking, context validation, ARN helpers, GetPluginInfo)
- Independent and testable
- ~40% of total tasks (T001-T042)

**MVP Task Count**: 42 tasks (T001-T042)

### Incremental Delivery Phases

1. **Phase 3 (US1)**: Implement health checking, context validation, ARN helpers - Plugin developers get new tools
   immediately
2. **Phase 4 (US2)**: Implement GetPluginInfo polish - Plugin operators get better error messages and migration guide
3. **Phase 5 (US3)**: Improve Connect protocol robustness - Better stability under load
4. **Phase 6 (US4)**: Enhance testing infrastructure - Better long-term quality
5. **Phase 7 (Polish)**: Final validation and documentation - Release readiness

### Test-First Approach (TDD)

**Constitution Requirement** (Principle III - Test-First Protocol - NON-NEGOTIABLE):

All test additions MUST follow TDD:

1. Write failing tests FIRST (T009-T013, T029-T031, T043-T047, T067-T068, T071)
2. Implement features to make tests PASS
3. Run `make test` after each implementation
4. All tests MUST pass with -race flag

**TDD Enforcement**:

- Failing tests written before implementation (marked as [USx] tests)
- Implementation tasks follow tests (marked as [USx] implementation)
- Linter prevents merging untested code (golangci-lint includes coverage check)

---

## Summary

**Total Task Count**: 89 tasks (T001-T089)

**Tasks Per User Story**:

- Setup & Foundational: 8 tasks (T001-T008)
- US1 (Plugin Development Experience, P1): 20 tasks (T009-T028)
- US2 (Plugin Information & Discovery, P1): 14 tasks (T029-T042)
- US3 (Connect Protocol Robustness, P2): 20 tasks (T043-T062)
- US4 (Testing Infrastructure & Quality, P3): 15 tasks (T063-T077)
- Polish & Cross-Cutting: 12 tasks (T078-T089)

**Parallel Opportunities**: 15 parallel task groups identified (see Dependencies section)

**Independent Test Criteria Per Story**:

- US1: Sample plugin with health checker, context validation, and ARN helpers works correctly
- US2: GetPluginInfo <100ms, user-friendly errors, legacy plugins return Unimplemented, migration guide complete
- US3: Concurrent requests, large payloads, graceful shutdown, connection resets all work correctly
- US4: Benchmarks don't block, extreme values rejected, fuzz test stable, coverage >80%

**Format Validation**:
✅ ALL tasks follow checklist format:

- ✅ Start with `- [ ]` checkbox
- ✅ Sequential Task ID (T001-T089)
- ✅ [P] marker for parallelizable tasks
- ✅ [USx] label for user story phases
- ✅ Clear action description with file path

**Recommended MVP**: Phase 3 + Phase 4 (42 tasks, ~40% of total) - Plugin developers get immediate value with
health checking, context validation, ARN helpers, and GetPluginInfo polish
