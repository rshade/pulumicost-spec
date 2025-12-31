# Tasks: SDK Documentation Consolidation

**Input**: Design documents from `/specs/031-sdk-docs-consolidation/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, quickstart.md

**Tests**: Not explicitly requested for this documentation feature. Validation via `make lint-markdown`
and `go build ./...` to ensure examples compile.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **SDK Documentation**: `sdk/go/pluginsdk/` at repository root
- **Testing Documentation**: `sdk/go/testing/` at repository root
- **Spec Documentation**: `specs/031-sdk-docs-consolidation/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and validation infrastructure

- [ ] T001 Review existing `sdk/go/pluginsdk/README.md` structure for section placement
- [ ] T002 Review existing `sdk/go/testing/README.md` structure for section placement
- [ ] T003 [P] Verify research.md content is complete for all documentation topics

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core documentation standards that ALL user stories depend on

**⚠️ CRITICAL**: Establishes documentation patterns before detailed writing begins

- [ ] T004 [P] Create `sdk/go/pluginsdk/example_test.go` file scaffold for godoc examples
- [ ] T005 [P] Establish copy-paste verification process for code examples
- [ ] T006 [P] Confirm inline comment style guidelines from existing codebase

**Checkpoint**: Foundation ready - user story documentation can now begin in parallel

---

## Phase 3: User Story 1 - Plugin Developer Learning the SDK (Priority: P1) 🎯 MVP

**Goal**: Enable new plugin developers to understand client lifecycle and basic usage patterns

**Independent Test**: New developer can follow documentation to build a working plugin client

**Relates to**: FR-001, FR-002, FR-003, FR-006 (Issues #240, #211, #207, #238)

### Implementation for User Story 1

- [ ] T007 [P] [US1] Add HTTP client ownership comment to `NewClient()` in
  `sdk/go/pluginsdk/client.go` (Issue #240)
- [ ] T008 [P] [US1] Add `ExampleClient_Close()` godoc example in
  `sdk/go/pluginsdk/example_test.go` (Issue #238)
- [ ] T009 [P] [US1] Add `ExampleClient_Close_userProvided()` godoc example in
  `sdk/go/pluginsdk/example_test.go` (Issue #238)
- [ ] T010 [P] [US1] Expand `contractedCostTolerance` comment with IEEE 754 rationale in
  `sdk/go/pluginsdk/focus_conformance.go` (Issue #211)
- [ ] T011 [P] [US1] Document `WithTags()` copy semantics in
  `sdk/go/pluginsdk/focus_builder.go` (Issue #207)
- [ ] T012 [US1] Verify all US1 examples compile with `go build ./sdk/go/pluginsdk/...`

**Checkpoint**: Plugin developers can understand client lifecycle from documentation alone

---

## Phase 4: User Story 2 - Migrating from gRPC to Connect-go (Priority: P1)

**Goal**: Enable existing gRPC plugin developers to migrate to connect-go multi-protocol support

**Independent Test**: Developer can migrate existing gRPC plugin following only the documentation

**Relates to**: FR-008 (Issue #235)

### Implementation for User Story 2

- [ ] T013 [US2] Add "## Migration Guide: gRPC to Connect" section header to
  `sdk/go/pluginsdk/README.md` (Issue #235)
- [ ] T014 [US2] Write server migration steps (5 steps) with code examples in
  `sdk/go/pluginsdk/README.md` (Issue #235)
- [ ] T015 [US2] Write client migration steps (5 steps) with code examples in
  `sdk/go/pluginsdk/README.md` (Issue #235)
- [ ] T016 [US2] Add protocol selection matrix table (gRPC vs Connect vs gRPC-Web) in
  `sdk/go/pluginsdk/README.md` (Issue #235)
- [ ] T017 [US2] Add backward compatibility guidance section in
  `sdk/go/pluginsdk/README.md` (Issue #235)
- [ ] T018 [US2] Verify migration guide code examples compile

**Checkpoint**: Developers can migrate gRPC plugins to connect-go from documentation

---

## Phase 5: User Story 3 - Operator Configuring CORS for Production (Priority: P2)

**Goal**: Enable operators to configure CORS for various deployment scenarios

**Independent Test**: Operator can configure CORS for their scenario and verify browser access

**Relates to**: FR-010 (Issue #236)

**⚠️ Dependency**: CORS WebConfig features depend on Issues #228 (CORS headers configurable)
and #229 (CORS max-age configurable). Verify these are merged before finalizing CORS documentation.

### Implementation for User Story 3

- [ ] T019 [US3] Add "## CORS Configuration" section header to
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T020 [P] [US3] Document Scenario 1: Local Development CORS configuration in
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T021 [P] [US3] Document Scenario 2: Single-Origin Production CORS configuration in
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T022 [P] [US3] Document Scenario 3: Multi-Origin Partners CORS configuration in
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T023 [P] [US3] Document Scenario 4: API Gateway Pattern in
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T024 [P] [US3] Document Scenario 5: Multi-Tenant SaaS CORS configuration in
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T025 [US3] Add CORS Security Guidelines subsection in
  `sdk/go/pluginsdk/README.md` (Issue #236)
- [ ] T026 [US3] Add CORS Debugging/Troubleshooting subsection in
  `sdk/go/pluginsdk/README.md` (Issue #236)

**Checkpoint**: Operators can configure CORS for 5+ deployment scenarios

---

## Phase 6: User Story 4 - Optimizing Client Performance (Priority: P2)

**Goal**: Enable developers to tune HTTP client configuration for high-throughput scenarios

**Independent Test**: Developer can apply configuration and measure improved performance

**Relates to**: FR-009 (Issue #237)

### Implementation for User Story 4

- [ ] T027 [US4] Add "## Performance Tuning" section header to
  `sdk/go/pluginsdk/README.md` (Issue #237)
- [ ] T028 [US4] Add configuration selection decision matrix (Default vs HighThroughput) in
  `sdk/go/pluginsdk/README.md` (Issue #237)
- [ ] T029 [US4] Document connection pool parameters (MaxIdleConns, MaxIdleConnsPerHost,
  IdleConnTimeout) in `sdk/go/pluginsdk/README.md` (Issue #237)
- [ ] T030 [US4] Add server timeout settings section (DoS protection) in
  `sdk/go/pluginsdk/README.md` (Issue #237)
- [ ] T031 [US4] Add protocol performance trade-offs section in
  `sdk/go/pluginsdk/README.md` (Issue #237)
- [ ] T032 [US4] Add connection pool monitoring tips section in
  `sdk/go/pluginsdk/README.md` (Issue #237)

**Checkpoint**: Developers can tune client performance based on documented guidance

---

## Phase 7: User Story 5 - Implementing Rate Limiting (Priority: P2)

**Goal**: Enable plugin developers to implement proper rate limiting for cloud provider APIs

**Independent Test**: Developer can implement rate limiting and verify proper backoff behavior

**Relates to**: FR-011 (Issue #233)

### Implementation for User Story 5

- [ ] T033 [US5] Add "## Rate Limiting" section header to
  `sdk/go/pluginsdk/README.md` (Issue #233)
- [ ] T034 [US5] Document token bucket pattern with golang.org/x/time/rate example in
  `sdk/go/pluginsdk/README.md` (Issue #233)
- [ ] T035 [US5] Add cloud provider rate limits reference table (AWS, Azure, GCP, K8s) in
  `sdk/go/pluginsdk/README.md` (Issue #233)
- [ ] T036 [US5] Document exponential backoff with jitter implementation in
  `sdk/go/pluginsdk/README.md` (Issue #233)
- [ ] T037 [US5] Document proper gRPC status codes for rate limiting (ResourceExhausted,
  Unavailable) in `sdk/go/pluginsdk/README.md` (Issue #233)

**Checkpoint**: Developers can implement proper rate limiting patterns

---

## Phase 8: User Story 6 - Understanding Concurrency Guarantees (Priority: P3)

**Goal**: Enable developers to understand thread safety guarantees for SDK components

**Independent Test**: Developer can run concurrent access tests with `-race` flag

**Relates to**: FR-012 (Issue #231)

### Implementation for User Story 6

- [ ] T038 [US6] Add "## Thread Safety" section header to
  `sdk/go/pluginsdk/README.md` (Issue #231)
- [ ] T039 [US6] Add thread safety summary table (Client, Server, WebConfig, PluginMetrics,
  ResourceMatcher, FocusRecordBuilder) in `sdk/go/pluginsdk/README.md` (Issue #231)
- [ ] T040 [US6] Document Client concurrent usage pattern with example in
  `sdk/go/pluginsdk/README.md` (Issue #231)
- [ ] T041 [US6] Document Server thread safety expectations in
  `sdk/go/pluginsdk/README.md` (Issue #231)
- [ ] T042 [US6] Document ResourceMatcher "configure before Serve()" contract in
  `sdk/go/pluginsdk/README.md` (Issue #231)
- [ ] T043 [US6] Document FocusRecordBuilder single-threaded usage pattern in
  `sdk/go/pluginsdk/README.md` (Issue #231)

**Checkpoint**: Developers understand thread safety guarantees for all SDK components

---

## Phase 9: User Story 7 - Understanding Code Semantics (Priority: P3)

**Goal**: Enable developers to understand precise function semantics via inline comments

**Independent Test**: Developer can understand behavior from comments without consulting source

**Relates to**: FR-004, FR-005, FR-007 (Issues #206, #208, #209)

### Implementation for User Story 7

- [ ] T044 [P] [US7] Fix correlation pattern comment to use `ResourceRecommendationInfo.id` in
  `sdk/go/pluginsdk/` files (Issue #206)
- [ ] T045 [P] [US7] Clarify test naming in `sdk/go/testing/*_test.go` for round-trip vs
  backward-compat semantics (Issue #208)
- [ ] T046 [US7] Add complete import statements with `pbc` alias to all code examples in
  `sdk/go/testing/README.md` (Issue #209)
- [ ] T047 [US7] Verify all testing README examples compile with complete imports

**Checkpoint**: Code semantics are clear from inline comments

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and quality assurance

- [ ] T048 Run `make lint-markdown` on all modified files
- [ ] T049 Run `go build ./sdk/go/pluginsdk/...` to verify all examples compile
- [ ] T050 Run `go build ./sdk/go/testing/...` to verify all examples compile
- [ ] T051 [P] Verify no placeholder values in examples ("TODO", "example-value", "xxx")
- [ ] T052 [P] Verify cross-provider coverage in examples (AWS, Azure, GCP)
- [ ] T053 Link all changes to their GitHub issue numbers in commit messages
- [ ] T054 Run quickstart.md validation checklist

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-9)**: All depend on Foundational phase completion
  - US1 and US2 (P1 priority) should complete first
  - US3, US4, US5 (P2 priority) can run in parallel after US1/US2
  - US6, US7 (P3 priority) can run in parallel after P2 stories
- **Polish (Phase 10)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies on other stories - MVP
- **User Story 2 (P1)**: No dependencies on other stories - MVP
- **User Story 3 (P2)**: No dependencies on other stories
- **User Story 4 (P2)**: No dependencies on other stories
- **User Story 5 (P2)**: No dependencies on other stories
- **User Story 6 (P3)**: No dependencies on other stories
- **User Story 7 (P3)**: No dependencies on other stories

### Within Each User Story

- Section header task before content tasks
- Content tasks can run in parallel if marked [P]
- Verification task last in each story phase

### Parallel Opportunities

- **Phase 2**: T004, T005, T006 can run in parallel
- **Phase 3 (US1)**: T007, T008, T009, T010, T011 can run in parallel (different files)
- **Phase 5 (US3)**: T020, T021, T022, T023, T024 can run in parallel (same file, different sections)
- **Phase 9 (US7)**: T044, T045 can run in parallel (different files)
- **Phase 10**: T051, T052 can run in parallel

---

## Parallel Example: User Story 1 (MVP)

```bash
# Launch all inline documentation tasks for User Story 1 together:
Task: "T007 - Add HTTP client ownership comment to NewClient() in sdk/go/pluginsdk/client.go"
Task: "T008 - Add ExampleClient_Close() godoc example in sdk/go/pluginsdk/example_test.go"
Task: "T009 - Add ExampleClient_Close_userProvided() example in sdk/go/pluginsdk/example_test.go"
Task: "T010 - Expand contractedCostTolerance comment in sdk/go/pluginsdk/focus_conformance.go"
Task: "T011 - Document WithTags() copy semantics in sdk/go/pluginsdk/focus_builder.go"

# Then run verification:
Task: "T012 - Verify all US1 examples compile"
```

---

## Parallel Example: User Story 3 (CORS Guide)

```bash
# Launch all CORS scenario documentation in parallel:
Task: "T020 - Document Scenario 1: Local Development CORS"
Task: "T021 - Document Scenario 2: Single-Origin Production CORS"
Task: "T022 - Document Scenario 3: Multi-Origin Partners CORS"
Task: "T023 - Document Scenario 4: API Gateway Pattern"
Task: "T024 - Document Scenario 5: Multi-Tenant SaaS CORS"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1 (Plugin Developer Learning)
4. Complete Phase 4: User Story 2 (Migration Guide)
5. **STOP and VALIDATE**: Test US1+US2 independently with `make lint-markdown`
6. Create PR for MVP scope (closes Issues #207, #211, #235, #238, #240)

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Stories 1+2 → Validate → PR #1 (MVP - closes 5 issues)
3. Add User Stories 3+4+5 → Validate → PR #2 (closes Issues #233, #236, #237)
4. Add User Stories 6+7 → Validate → PR #3 (closes Issues #206, #208, #209, #231)
5. Final Polish → All 12 issues closed

### Parallel Team Strategy

With multiple contributors:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Contributor A: User Stories 1 + 2 (P1 priority - MVP)
   - Contributor B: User Stories 3 + 4 + 5 (P2 priority)
   - Contributor C: User Stories 6 + 7 (P3 priority)
3. All contributors: Phase 10 Polish together

---

## GitHub Issue Mapping

| Task Range | GitHub Issue | User Story | Description |
|------------|--------------|------------|-------------|
| T007 | #240 | US1 | HTTP client ownership semantics |
| T008-T009 | #238 | US1 | Client.Close() godoc examples |
| T010 | #211 | US1 | contractedCostTolerance explanation |
| T011 | #207 | US1 | WithTags shared-map semantics |
| T013-T018 | #235 | US2 | Migration guide (gRPC to connect) |
| T019-T026 | #236 | US3 | CORS best practices |
| T027-T032 | #237 | US4 | Performance tuning guide |
| T033-T037 | #233 | US5 | Rate limiting patterns |
| T038-T043 | #231 | US6 | Thread safety documentation |
| T044 | #206 | US7 | Correlation pattern comment fix |
| T045 | #208 | US7 | Test naming clarification |
| T046-T047 | #209 | US7 | Complete import statements |

---

## Notes

- [P] tasks = different files or independent sections, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and validatable
- Verify markdown lint passes after each phase
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All code examples must compile without modification (copy-paste ready)
