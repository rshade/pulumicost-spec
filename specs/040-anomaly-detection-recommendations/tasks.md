# Tasks: Anomaly Detection via Recommendations

**Input**: Design documents from `/specs/040-anomaly-detection-recommendations/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Conformance tests are included per constitution (Test-First Protocol).

**Organization**: Tasks follow proto-first development pattern per constitution.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

This is a gRPC specification repository:

- **Proto**: `proto/finfocus/v1/costsource.proto`
- **Go SDK**: `sdk/go/proto/` (generated), `sdk/go/testing/` (conformance)
- **TypeScript SDK**: `sdk/typescript/packages/client/`
- **Docs**: `docs/sdk/`

---

## Phase 1: Setup (Proto Contract First)

**Purpose**: Add enum values to proto definitions per constitution (Contract First)

- [x] T001 Add RECOMMENDATION_CATEGORY_ANOMALY (=5) to RecommendationCategory enum in
  proto/finfocus/v1/costsource.proto with documentation comments
- [x] T002 Add RECOMMENDATION_ACTION_TYPE_INVESTIGATE (=12) to RecommendationActionType enum in
  proto/finfocus/v1/costsource.proto with documentation comments
- [x] T003 Run `buf lint` to validate proto changes in proto/finfocus/v1/costsource.proto
- [x] T004 Run `buf breaking` to verify backward compatibility of proto changes

**Checkpoint**: Proto contract updated - SDK regeneration can proceed

---

## Phase 2: Foundational (SDK Regeneration)

**Purpose**: Regenerate all SDK bindings from updated proto definitions

**⚠️ CRITICAL**: User story conformance tests require regenerated SDK code

- [x] T005 Run `make generate` to regenerate Go SDK bindings in sdk/go/proto/finfocus/v1/
- [x] T006 Verify generated Go enum constants in sdk/go/proto/finfocus/v1/costsource.pb.go:
  RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY and
  RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE
- [x] T007 [P] Regenerate TypeScript SDK bindings in sdk/typescript/packages/client/ via buf
- [x] T007a [P] Verify generated TypeScript enum constants include RECOMMENDATION_CATEGORY_ANOMALY
  and RECOMMENDATION_ACTION_TYPE_INVESTIGATE in sdk/typescript/packages/client/

**Checkpoint**: All SDKs regenerated - conformance testing can begin

---

## Phase 3: User Story 1 - Unified Actionable Cost Insights View (Priority: P1) 🎯 MVP

**Goal**: Enable GetRecommendations to return both optimization recommendations and cost anomalies
in a single response

**Independent Test**: Call GetRecommendations and verify ANOMALY category recommendations can be
returned alongside other categories

### Conformance Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before mock implementation**

- [x] T008 [US1] Write conformance test: plugin returns ANOMALY category recommendation in
  sdk/go/testing/anomaly_conformance_test.go - test GetRecommendations returns mixed categories
- [x] T009 [US1] Write conformance test: category filter works for ANOMALY in
  sdk/go/testing/anomaly_conformance_test.go - test category=ANOMALY returns only anomalies
- [x] T010 [US1] Write conformance test: plugins not supporting anomalies return zero ANOMALY
  recommendations in sdk/go/testing/anomaly_conformance_test.go

### Implementation for User Story 1

- [x] T011 [US1] Update MockPlugin in sdk/go/testing/mock_plugin.go to support returning
  ANOMALY category recommendations for testing
- [x] T012 [US1] Run conformance tests to verify ANOMALY category filtering works with existing
  RecommendationFilter logic

**Checkpoint**: User Story 1 complete - unified view works with ANOMALY category

---

## Phase 4: User Story 2 - Anomaly Triage by Confidence Score (Priority: P2)

**Goal**: Verify existing confidence_score filter works with ANOMALY recommendations

**Independent Test**: Return anomaly recommendations with varying confidence scores and verify
min_confidence_score filter works correctly

### Conformance Tests for User Story 2

- [x] T013 [US2] Write conformance test: confidence_score filter works for anomalies in
  sdk/go/testing/anomaly_conformance_test.go - test min_confidence_score=0.5 filters correctly

### Implementation for User Story 2

- [x] T014 [US2] Update MockPlugin to return anomaly recommendations with varying confidence
  scores (0.3, 0.6, 0.9) in sdk/go/testing/mock_plugin.go
- [x] T015 [US2] Run conformance test to verify min_confidence_score filter applies to ANOMALY
  category recommendations

**Checkpoint**: User Story 2 complete - confidence-based triage works for anomalies

---

## Phase 5: User Story 3 - Investigate Anomalous Spending (Priority: P2)

**Goal**: Verify INVESTIGATE action_type signals human investigation required

**Independent Test**: Return anomaly recommendation with INVESTIGATE action and verify it contains
sufficient context in description and metadata

### Conformance Tests for User Story 3

- [x] T016 [US3] Write conformance test: INVESTIGATE action_type returned with ANOMALY category
  in sdk/go/testing/anomaly_conformance_test.go
- [x] T017 [US3] Write conformance test: anomaly recommendation has description and metadata
  context in sdk/go/testing/anomaly_conformance_test.go

### Implementation for User Story 3

- [x] T018 [US3] Update MockPlugin to return ANOMALY recommendations with INVESTIGATE action,
  description, and metadata in sdk/go/testing/mock_plugin.go
- [x] T019 [US3] Run conformance tests to verify INVESTIGATE action semantics

**Checkpoint**: User Story 3 complete - INVESTIGATE action provides investigation context

---

## Phase 6: User Story 4 - Exclude Anomalies from Optimization Workflows (Priority: P3)

**Goal**: Verify category filtering allows excluding ANOMALY recommendations

**Independent Test**: Return mixed recommendations and verify filtering by non-ANOMALY categories
excludes anomalies

### Conformance Tests for User Story 4

- [x] T020 [US4] Write conformance test: filter by category=COST excludes ANOMALY in
  sdk/go/testing/anomaly_conformance_test.go
- [x] T021 [US4] Write conformance test: negative estimated_savings accepted for overspend anomalies
  in sdk/go/testing/anomaly_conformance_test.go

### Implementation for User Story 4

- [x] T022 [US4] Update MockPlugin to return mixed recommendations (COST + ANOMALY) with
  varying estimated_savings in sdk/go/testing/mock_plugin.go
- [x] T023 [US4] Run conformance tests to verify category exclusion and negative savings

**Checkpoint**: User Story 4 complete - automation workflows can exclude anomalies

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, validation, and final quality checks

- [x] T024 [P] Add anomaly usage section to SDK documentation in docs/sdk/README.md or
  sdk/go/pluginsdk/README.md with field mapping table from data-model.md. Include prominent
  callout: "Note: estimated_savings is NEGATIVE for overspend anomalies (cost above baseline)"
  (Note: Documentation provided via specs/040/quickstart.md)
- [x] T025 [P] Add anomaly example JSON to documentation showing all recommended fields
  (Note: Examples provided via specs/040/quickstart.md and data-model.md)
- [x] T026 [P] Run `make lint` to verify all Go code passes linting
- [x] T027 [P] Run `make test` to verify all conformance tests pass
- [x] T028 [P] Run `make validate` to run full validation pipeline
- [x] T029 Verify quickstart.md example code compiles (optional - example is for documentation)
- [x] T030 [P] Update CLAUDE.md with new enum values if not already updated by agent context script

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Proto)**: No dependencies - start immediately (CONTRACT FIRST per constitution)
- **Phase 2 (SDK Regen)**: Depends on Phase 1 - BLOCKS all conformance tests
- **User Stories (Phases 3-6)**: All depend on Phase 2 completion
  - Can proceed sequentially in priority order (P1 → P2 → P2 → P3)
  - Or in parallel if team capacity allows
- **Phase 7 (Polish)**: Depends on at least User Story 1 being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 2 - No dependencies on other stories - **MVP**
- **User Story 2 (P2)**: Can start after Phase 2 - Extends US1 test fixtures
- **User Story 3 (P2)**: Can start after Phase 2 - Extends US1 test fixtures
- **User Story 4 (P3)**: Can start after Phase 2 - Extends US1 test fixtures

### Within Each User Story

- Write conformance test FIRST - ensure it FAILS
- Update MockPlugin to support test case
- Run test to verify behavior
- Mark story complete when tests pass

### Parallel Opportunities

- T003 and T004 can run after T001+T002 complete (parallel lint + breaking check)
- T005 and T007 can run in parallel (Go SDK and TypeScript SDK regeneration)
- T024, T025, T026, T027, T028, T030 can all run in parallel (documentation + validation)

---

## Parallel Example: Phase 2 (SDK Regeneration)

```bash
# Launch SDK regeneration in parallel:
Task: "Run make generate to regenerate Go SDK"
Task: "Regenerate TypeScript SDK via buf"
```

---

## Parallel Example: Phase 7 (Polish)

```bash
# Launch all polish tasks together:
Task: "Add anomaly usage section to SDK documentation"
Task: "Add anomaly example JSON to documentation"
Task: "Run make lint"
Task: "Run make test"
Task: "Run make validate"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Proto enum additions (T001-T004)
2. Complete Phase 2: SDK regeneration (T005-T007)
3. Complete Phase 3: User Story 1 conformance (T008-T012)
4. **STOP and VALIDATE**: Verify ANOMALY category works in GetRecommendations
5. Deploy/release if ready

### Incremental Delivery

1. Proto + SDK regen → Foundation ready
2. Add User Story 1 → Test → Merge (MVP!)
3. Add User Story 2 → Test → Merge (confidence filtering)
4. Add User Story 3 → Test → Merge (INVESTIGATE action)
5. Add User Story 4 → Test → Merge (exclusion filtering)
6. Polish → Documentation complete → Final release

### Single Developer Strategy

1. T001-T007 sequentially (proto → sdk regen)
2. T008-T012 (User Story 1)
3. T013-T015 (User Story 2)
4. T016-T019 (User Story 3)
5. T020-T023 (User Story 4)
6. T024-T030 in parallel (documentation + validation)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Constitution requires test-first: conformance tests written before mock implementation
- Commit after each phase completion
- Stop at any checkpoint to validate independently
- Total tasks: 31
- Proto changes: 2 enum values
- No new messages or RPCs - minimal additive change
