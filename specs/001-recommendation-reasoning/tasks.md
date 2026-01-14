# Actionable Tasks: Standardized Recommendation Reasoning Metadata

**Feature**: Standardized Recommendation Reasoning Metadata
**Branch**: `001-recommendation-reasoning`
**Status**: Completed

## Implementation Strategy
- **MVP**: Add the `RecommendationReason` enum and update the `Recommendation` message to support standardized reasoning.
- **Incremental**:
  1. Define the Test (Validation - Test First). [X]
  2. Define the Enum (Schema). [X]
  3. Update the Message (Schema). [X]
  4. Generate SDK (Code). [X]
  5. Verify (Pass Test). [X]

## Dependencies
- **Story Dependency**: US2 depends on US1 (Schema must exist before consumption).
- **Critical Path**: T001 -> T002 -> T003 -> T004 -> T005 -> T006

## Parallel Execution Examples
- **US1**: T003 (Enum) and T004 (Message update) can be drafted in parallel, but must be committed together for `buf generate`.

---

## Phase 1: Setup
**Goal**: Verify environment and tools.

- [X] T001 Verify `buf` installation and linting configuration
  - Action: Run `buf lint` to ensure clean baseline.
  - File: `proto/`

---

## Phase 2: Foundational (Prerequisites)
**Goal**: Ensure core schema files are ready for updates.

*No specific foundational tasks for this feature beyond Setup.*

---

## Phase 3: User Story 1 - Plugin Developer Maps Upstream Reasons (Priority P1)
**Goal**: Enable plugins to set standardized reason codes in recommendations.
**Independent Test**: Can create a `Recommendation` object with `PrimaryReason` set to `RECOMMENDATION_REASON_IDLE` and serialize it.

### Verification (Test First)
- [X] T002 [US1] Create serialization test in `sdk/go/finfocus/v1/recommendation_test.go`
  - Action: Create a new test file (or update existing) to verify that `RecommendationReason` can be set, serialized, and deserialized correctly.
  - Note: This test will fail compilation until T005 is complete. This satisfies Constitution Principle V.
  - File: `sdk/go/finfocus/v1/recommendation_test.go`

### Schema Definitions
- [X] T003 [US1] Define `RecommendationReason` enum in `proto/finfocus/v1/enums.proto`
  - Action: Add the enum definition with values: UNSPECIFIED, OVER_PROVISIONED, UNDER_PROVISIONED, IDLE, REDUNDANT, OBSOLETE_GENERATION.
  - File: `proto/finfocus/v1/enums.proto`

- [X] T004 [US1] Update `Recommendation` message in `proto/finfocus/v1/costsource.proto`
  - Action: Add `primary_reason` (type: `RecommendationReason`) and `secondary_reasons` (type: `repeated RecommendationReason`) fields.
  - File: `proto/finfocus/v1/costsource.proto`

### SDK Generation
- [X] T005 [US1] Generate Go SDK
  - Action: Run `buf generate` to update the Go bindings and make T002 compile/pass.
  - File: `sdk/go/`

---

## Phase 4: User Story 2 - Consumer Tool Displays Reasoning (Priority P2)
**Goal**: Ensure consumers can switch/case on the new enum for display logic.
**Independent Test**: Test demonstrates type-safe switching on the enum.

### Integration Validation
- [X] T006 [US2] Verify enum consumption in test/example
  - Action: Extend `recommendation_test.go` to include a test case that switches on `PrimaryReason` and asserts the correct string output (mimicking a UI/CLI display logic).
  - File: `sdk/go/finfocus/v1/recommendation_test.go`

---



## Phase 5: Polish & Cross-Cutting

**Goal**: Final cleanup and documentation.



- [X] T007 Run final linting and formatting

  - Action: Run `buf lint` and `go fmt ./sdk/...`

  - File: `(Repo root)`
