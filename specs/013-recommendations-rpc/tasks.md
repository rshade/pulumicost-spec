# Tasks: GetRecommendations RPC

**Input**: Design documents from `/specs/013-recommendations-rpc/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Tests**: Conformance tests are explicitly required per SC-008 in spec.md.

**Organization**: Tasks are grouped by user story to enable independent implementation and
testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Based on plan.md project structure:

- **Proto**: `proto/pulumicost/v1/costsource.proto`
- **SDK**: `sdk/go/proto/`, `sdk/go/pluginsdk/`, `sdk/go/testing/`
- **Examples**: `examples/recommendations/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and proto generation tooling verification

- [ ] T001 Verify buf v1.32.1 is installed via `make generate`
- [ ] T002 [P] Verify zerolog v1.34.0+ and prometheus/client_golang dependencies in go.mod
- [ ] T003 [P] Create examples/recommendations/ directory structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core proto definitions that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Proto Definitions

- [ ] T010 Add RecommendationCategory enum to proto/pulumicost/v1/costsource.proto (values:
  UNSPECIFIED, COST, PERFORMANCE, SECURITY, RELIABILITY)
- [ ] T011 [P] Add RecommendationActionType enum to proto/pulumicost/v1/costsource.proto
  (values: UNSPECIFIED, RIGHTSIZE, TERMINATE, PURCHASE_COMMITMENT, ADJUST_REQUESTS, MODIFY,
  DELETE_UNUSED)
- [ ] T012 [P] Add RecommendationPriority enum to proto/pulumicost/v1/costsource.proto
  (values: UNSPECIFIED, LOW, MEDIUM, HIGH, CRITICAL)
- [ ] T013 Add ResourceUtilization message to proto/pulumicost/v1/costsource.proto (fields:
  cpu_percent, memory_percent, storage_percent, network_in_mbps, network_out_mbps,
  custom_metrics)
- [ ] T014 [P] Add ResourceRecommendationInfo message to proto/pulumicost/v1/costsource.proto
  (fields: id, name, provider, resource_type, region, sku, tags, utilization)
- [ ] T015 Add RecommendationImpact message to proto/pulumicost/v1/costsource.proto (fields:
  estimated_savings, currency, projection_period, current_cost, projected_cost,
  savings_percentage, implementation_cost, migration_effort_hours)
- [ ] T016 Add RecommendationSummary message to proto/pulumicost/v1/costsource.proto (fields:
  total_recommendations, total_estimated_savings, currency, projection_period,
  count_by_category, savings_by_category, count_by_action_type, savings_by_action_type)

### Action Detail Messages

- [ ] T017 Add RightsizeAction message to proto/pulumicost/v1/costsource.proto (fields:
  current_sku, recommended_sku, current_instance_type, recommended_instance_type,
  projected_utilization)
- [ ] T018 [P] Add TerminateAction message to proto/pulumicost/v1/costsource.proto (fields:
  termination_reason, idle_days)
- [ ] T019 [P] Add CommitmentAction message to proto/pulumicost/v1/costsource.proto (fields:
  commitment_type, term, payment_option, recommended_quantity, scope)
- [ ] T020 [P] Add KubernetesResources message to proto/pulumicost/v1/costsource.proto
  (fields: cpu, memory)
- [ ] T021 Add KubernetesAction message to proto/pulumicost/v1/costsource.proto (fields:
  cluster_id, namespace, controller_kind, controller_name, container_name, current_requests,
  recommended_requests, current_limits, recommended_limits, algorithm) - depends on T020
- [ ] T022 [P] Add ModifyAction message to proto/pulumicost/v1/costsource.proto (fields:
  modification_type, current_config, recommended_config)

### Core Recommendation Message

- [ ] T023 Add Recommendation message with oneof action_detail to
  proto/pulumicost/v1/costsource.proto (all fields from data-model.md including oneof for
  rightsize, terminate, commitment, kubernetes, modify) - depends on T010-T022

### Request/Response Messages

- [ ] T024 Add RecommendationFilter message to proto/pulumicost/v1/costsource.proto (fields:
  provider, region, resource_type, category, action_type)
- [ ] T025 Add GetRecommendationsRequest message to proto/pulumicost/v1/costsource.proto
  (fields: filter, projection_period, page_size, page_token) - depends on T024
- [ ] T026 Add GetRecommendationsResponse message to proto/pulumicost/v1/costsource.proto
  (fields: recommendations, summary, next_page_token) - depends on T023, T016

### Service Method

- [ ] T027 Add GetRecommendations RPC method to CostSourceService in
  proto/pulumicost/v1/costsource.proto - depends on T025, T026
- [ ] T028 Add capabilities field (map<string, bool>) to SupportsResponse message in
  proto/pulumicost/v1/costsource.proto

### Code Generation

- [ ] T029 Run `make generate` to regenerate sdk/go/proto/ from updated costsource.proto -
  depends on T027, T028
- [ ] T030 Run `buf lint` and `buf breaking` to verify proto changes pass validation -
  depends on T029

**Checkpoint**: Proto foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Retrieve All Recommendations (Priority: P1) 🎯 MVP

**Goal**: Enable plugins to return all available recommendations with valid fields and summary

**Independent Test**: Call GetRecommendations with empty filter, verify recommendations
returned with required fields populated and summary included

### PluginSDK Interface

- [ ] T040 [US1] Define RecommendationsProvider interface in sdk/go/pluginsdk/sdk.go with
  GetRecommendations method signature
- [ ] T041 [US1] Add type assertion in sdk/go/pluginsdk/sdk.go to auto-detect
  RecommendationsProvider implementation (follow SupportsProvider pattern)
- [ ] T042 [US1] Add GetRecommendations dispatch method in sdk/go/pluginsdk/sdk.go that
  routes to plugin or returns empty list

### Validation Helpers

- [ ] T043 [US1] Add ValidateRecommendation function in sdk/go/pluginsdk/helpers.go
  (validates required fields: id, category, action_type, resource, impact)
- [ ] T044 [US1] Add ValidateRecommendationImpact function in sdk/go/pluginsdk/helpers.go
  (validates currency via sdk/go/currency, estimated_savings)
- [ ] T045 [US1] Add ValidateConfidenceScore function in sdk/go/pluginsdk/helpers.go
  (validates 0.0-1.0 range when present)

### Mock Plugin Implementation

- [ ] T046 [US1] Add GetRecommendations method to MockPlugin in sdk/go/testing/mock_plugin.go
- [ ] T047 [US1] Add SetRecommendationsResponse configuration method in
  sdk/go/testing/mock_plugin.go for test customization
- [ ] T048 [US1] Add GenerateSampleRecommendations helper in sdk/go/testing/mock_plugin.go
  to create realistic test data

### Basic Conformance Tests

- [ ] T049 [US1] Add TestGetRecommendations_Basic in sdk/go/testing/conformance_test.go
  (validates response has correct message types)
- [ ] T050 [US1] Add TestGetRecommendations_EmptyPlugin in sdk/go/testing/conformance_test.go
  (validates empty list returned for non-implementing plugins)
- [ ] T051 [US1] Add TestGetRecommendations_SummaryCalculation in
  sdk/go/testing/conformance_test.go (validates summary totals match individual items)

### Integration Tests

- [ ] T052 [US1] Add TestGetRecommendationsIntegration in
  sdk/go/testing/integration_test.go using TestHarness

**Checkpoint**: User Story 1 complete - basic recommendation retrieval works

---

## Phase 4: User Story 2 - Filter by Category (Priority: P2)

**Goal**: Enable filtering recommendations by category (cost, performance, security,
reliability)

**Independent Test**: Create mock recommendations across categories, verify filter returns
only matching items

### Filtering Implementation

- [ ] T054 [US2] Add applyFilter helper function in sdk/go/pluginsdk/helpers.go for
  filtering recommendations by RecommendationFilter criteria
- [ ] T055 [US2] Implement category filtering logic in applyFilter (match
  RecommendationCategory enum values)

### Conformance Tests for Category Filtering

- [ ] T056 [US2] Add TestGetRecommendations_FilterByCategory in
  sdk/go/testing/conformance_test.go (validates COST category filter)
- [ ] T057 [US2] Add TestGetRecommendations_FilterByCategory_Performance in
  sdk/go/testing/conformance_test.go (validates PERFORMANCE category filter)
- [ ] T058 [US2] Add TestGetRecommendations_FilterNoMatches in
  sdk/go/testing/conformance_test.go (validates empty result when no recommendations match)
- [ ] T059 [US2] Add TestGetRecommendations_FilterByProvider in
  sdk/go/testing/conformance_test.go (validates provider filter returns only matching)

**Checkpoint**: User Story 2 complete - category filtering works independently

---

## Phase 5: User Story 3 - Filter by Action Type (Priority: P2)

**Goal**: Enable filtering recommendations by action type (rightsize, terminate,
purchase_commitment, etc.)

**Independent Test**: Create mock recommendations with different action types, verify filter
accuracy

### Action Type Filtering Implementation

- [ ] T060 [US3] Add action_type filtering logic to applyFilter in
  sdk/go/pluginsdk/helpers.go (match RecommendationActionType enum values)

### Conformance Tests for Action Type Filtering

- [ ] T061 [US3] Add TestGetRecommendations_FilterByActionType_Rightsize in
  sdk/go/testing/conformance_test.go
- [ ] T062 [US3] Add TestGetRecommendations_FilterByActionType_AdjustRequests in
  sdk/go/testing/conformance_test.go
- [ ] T063 [US3] Add TestGetRecommendations_FilterCombined in
  sdk/go/testing/conformance_test.go (validates multiple criteria applied together)
- [ ] T064 [US3] Add TestGetRecommendations_FilterByRegion in
  sdk/go/testing/conformance_test.go (validates region filter returns only matching)
- [ ] T065 [US3] Add TestGetRecommendations_FilterByResourceType in
  sdk/go/testing/conformance_test.go (validates resource_type filter)

**Checkpoint**: User Story 3 complete - action type filtering works independently

---

## Phase 6: User Story 4 - Paginate Large Result Sets (Priority: P3)

**Goal**: Enable pagination through large recommendation sets with page tokens

**Independent Test**: Create 500 mock recommendations, verify page_size and page_token work
correctly

### Pagination Implementation

- [ ] T070 [US4] Add paginate helper function in sdk/go/pluginsdk/helpers.go (accepts
  recommendations, page_size, page_token; returns page and next_token)
- [ ] T071 [US4] Implement opaque base64 token generation in paginate helper (encode offset
  or cursor)
- [ ] T072 [US4] Add validatePageToken function in sdk/go/pluginsdk/helpers.go (validates
  token format, returns decoded offset)

### Conformance Tests for Pagination

- [ ] T073 [US4] Add TestGetRecommendations_Pagination_FirstPage in
  sdk/go/testing/conformance_test.go (validates page_size respected)
- [ ] T074 [US4] Add TestGetRecommendations_Pagination_NextPage in
  sdk/go/testing/conformance_test.go (validates next_page_token works)
- [ ] T075 [US4] Add TestGetRecommendations_Pagination_LastPage in
  sdk/go/testing/conformance_test.go (validates empty token on last page)
- [ ] T076 [US4] Add TestGetRecommendations_Pagination_InvalidToken in
  sdk/go/testing/conformance_test.go (validates error on invalid token)

### Performance Tests

- [ ] T077 [US4] Add BenchmarkGetRecommendations_LargeResultSet in
  sdk/go/testing/benchmark_test.go (10,000 recommendations, verify <500ms per SC-005)

**Checkpoint**: User Story 4 complete - pagination works for large result sets

---

## Phase 7: User Story 5 - Provider-Specific Details (Priority: P3)

**Goal**: Enable viewing provider-specific action details (AWS instance types, K8s container
resources, etc.)

**Independent Test**: Verify action-specific fields populated correctly for each
recommendation type

### Example Recommendations

- [ ] T080 [P] [US5] Create examples/recommendations/aws_rightsizing.json with EC2
  instance rightsizing recommendation
- [ ] T081 [P] [US5] Create examples/recommendations/kubernetes_request_sizing.json with
  container request adjustment recommendation
- [ ] T082 [P] [US5] Create examples/recommendations/azure_advisor.json with Azure
  VM rightsizing recommendation
- [ ] T083 [P] [US5] Create examples/recommendations/gcp_recommender.json with GCE
  instance recommendation
- [ ] T084 [P] [US5] Create examples/recommendations/commitment_purchase.json with
  reserved instance recommendation

### Conformance Tests for Action Details

- [ ] T085 [US5] Add TestGetRecommendations_RightsizeAction in
  sdk/go/testing/conformance_test.go (validates current/recommended instance types)
- [ ] T086 [US5] Add TestGetRecommendations_KubernetesAction in
  sdk/go/testing/conformance_test.go (validates cluster, namespace, container details)
- [ ] T087 [US5] Add TestGetRecommendations_CommitmentAction in
  sdk/go/testing/conformance_test.go (validates commitment_type, term, quantity)

**Checkpoint**: User Story 5 complete - provider-specific details work correctly

---

## Phase 8: Observability (Cross-Cutting)

**Purpose**: Add logging, metrics, and tracing per FR-020 through FR-022

### Logging (zerolog)

- [ ] T090 [P] Add logging constants in sdk/go/pluginsdk/logging.go (FieldRecommendationCount,
  FieldFilterCategory, FieldFilterActionType, FieldPageSize, FieldTotalSavings)
- [ ] T091 Add GetRecommendations request/response logging in sdk/go/pluginsdk/sdk.go using
  zerolog patterns

### Prometheus Metrics

- [ ] T092 Add recommendations-specific metrics in sdk/go/pluginsdk/metrics.go
  (pulumicost_plugin_recommendations_returned_total counter,
  pulumicost_plugin_recommendations_per_response histogram)
- [ ] T093 Add GetRecommendations metrics recording in MetricsUnaryServerInterceptor

### Tracing

- [ ] T094 Add trace context propagation in GetRecommendations dispatch method
  (extract/inject trace headers)

### Observability Tests

- [ ] T095 Add TestGetRecommendations_MetricsRecorded in sdk/go/testing/conformance_test.go
  (validates metrics emitted)
- [ ] T096 Add TestGetRecommendations_LoggingFormat in sdk/go/testing/conformance_test.go
  (validates log structure)

**Checkpoint**: Observability complete - metrics and logs consistent with other RPC methods

---

## Phase 9: Polish & Documentation

**Purpose**: Final documentation, examples, and cleanup

### Documentation

- [ ] T100 [P] Update sdk/go/testing/README.md with GetRecommendations testing examples
- [ ] T101 [P] Update examples/README.md with recommendation example descriptions
- [ ] T102 [P] Add GetRecommendations section to sdk/go/CLAUDE.md
- [ ] T103 Update sdk/go/pluginsdk/README.md with RecommendationsProvider interface docs

### Conformance Suite Updates

- [ ] T104 Update RunBasicConformanceTests in sdk/go/testing/conformance.go to include
  recommendation tests
- [ ] T105 Update RunStandardConformanceTests in sdk/go/testing/conformance.go to include
  filtering and pagination tests
- [ ] T106 Update RunAdvancedConformanceTests in sdk/go/testing/conformance.go to include
  performance tests

### Final Validation

- [ ] T107 Run full conformance suite: `go test -v ./sdk/go/testing/ -run TestConformance`
- [ ] T108 Run benchmarks: `go test -bench=. -benchmem ./sdk/go/testing/`
- [ ] T109 Run `make lint` to verify all code passes linting
- [ ] T110 Run `make validate` for full validation pipeline
- [ ] T111 Verify quickstart.md steps work with generated code

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - US1 (P1): Can start after Phase 2
  - US2 (P2): Can start after Phase 2, independent of US1
  - US3 (P2): Can start after Phase 2, independent of US1/US2
  - US4 (P3): Can start after Phase 2, independent of other stories
  - US5 (P3): Can start after Phase 2, independent of other stories
- **Observability (Phase 8)**: Can start after Phase 3 (needs base SDK implementation)
- **Polish (Phase 9)**: Depends on all phases complete

### Within Each User Story

- Validation helpers before conformance tests
- Mock plugin updates before integration tests
- Core implementation before tests that use it

### Parallel Opportunities

- T011, T012 can run in parallel (different enum definitions)
- T017-T022 action messages can run in parallel (independent messages)
- T080-T084 example files can run in parallel (different files)
- T090 logging constants can run in parallel with T092 metrics definitions
- User Stories 2-5 can all run in parallel after Phase 2 completes

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - proto definitions)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test basic recommendation retrieval
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Proto foundation ready
2. User Story 1 → Basic retrieval works (MVP!)
3. User Stories 2 + 3 → Filtering works
4. User Story 4 → Pagination works
5. User Story 5 → Provider details work
6. Observability → Production-ready
7. Polish → Documentation complete

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Conformance tests required per SC-008
- Proto changes must pass `buf lint` and `buf breaking`
- Currency validation uses existing sdk/go/currency package per research.md decision
- Follow SupportsProvider pattern for optional interface per research.md decision
