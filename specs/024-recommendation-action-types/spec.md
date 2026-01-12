# Feature Specification: Extend RecommendationActionType Enum

**Feature Branch**: `019-recommendation-action-types`
**Created**: 2025-12-17
**Status**: Draft
**Input**: GitHub Issue #170 - Extend RecommendationActionType enum with additional action
types for comprehensive FinOps platform coverage

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Developer Uses Extended Action Types (Priority: P1)

As a cost source plugin developer, I want to categorize recommendations with more specific
action types so that users receive actionable, well-categorized cost optimization suggestions.

**Why this priority**: This is the core value proposition - enabling plugins to provide more
precise recommendation categorization that maps directly to what major FinOps platforms
(Azure Advisor, GCP Recommender, AWS Cost Explorer) provide.

**Independent Test**: Can be fully tested by implementing a mock plugin that returns
recommendations with the new action types (MIGRATE, CONSOLIDATE, SCHEDULE, REFACTOR, OTHER)
and verifying the SDK correctly serializes/deserializes them.

**Acceptance Scenarios**:

1. **Given** a plugin returns a recommendation with `RECOMMENDATION_ACTION_TYPE_MIGRATE`,
   **When** the core processes this recommendation,
   **Then** the action type is correctly preserved and displayed to the user
2. **Given** a plugin returns a recommendation with `RECOMMENDATION_ACTION_TYPE_CONSOLIDATE`,
   **When** the recommendation is serialized to JSON,
   **Then** the enum value maps to "CONSOLIDATE"
3. **Given** a plugin returns a recommendation with `RECOMMENDATION_ACTION_TYPE_SCHEDULE`,
   **When** the core filters recommendations by action type,
   **Then** scheduled recommendations are correctly identified
4. **Given** a plugin returns a recommendation with `RECOMMENDATION_ACTION_TYPE_REFACTOR`,
   **When** the core categorizes recommendations,
   **Then** architectural change recommendations are grouped separately
5. **Given** a plugin cannot categorize a recommendation into existing types,
   **When** it uses `RECOMMENDATION_ACTION_TYPE_OTHER`,
   **Then** the recommendation is accepted and processed with a generic category

---

### User Story 2 - Backward Compatibility with Existing Plugins (Priority: P1)

As an existing plugin maintainer, I want the enum extension to be backward compatible so that
my plugin continues to work without immediate updates.

**Why this priority**: Breaking existing plugins would cause immediate disruption to the
ecosystem. Proto3 enums are inherently backward compatible, but this must be explicitly
tested and documented.

**Independent Test**: Can be tested by running existing plugin implementations against the
updated SDK and verifying all existing functionality works unchanged.

**Acceptance Scenarios**:

1. **Given** an existing plugin using only the original 6 action types,
   **When** the SDK is updated,
   **Then** the plugin compiles and functions without modification
2. **Given** an existing plugin receives a request from a newer core using extended types,
   **When** the plugin doesn't recognize a new type,
   **Then** it can safely ignore or handle it as unknown
3. **Given** a plugin binary compiled with the old SDK,
   **When** it communicates with a core using the new SDK,
   **Then** gRPC communication succeeds for all existing operations

---

### User Story 3 - Core CLI Categorization (Priority: P2)

As a finfocus CLI user running `finfocus cost recommendations`, I want recommendations
categorized by the full range of action types so that I can prioritize and filter
recommendations effectively.

**Why this priority**: This is the downstream consumer of the enum extension. While critical
for end-user value, it depends on the proto changes being available first.

**Independent Test**: Can be tested by mocking plugin responses with various action types and
verifying CLI output correctly groups and labels recommendations.

**Acceptance Scenarios**:

1. **Given** the CLI receives recommendations with mixed action types including MIGRATE,
   **When** displaying results,
   **Then** migration recommendations are labeled clearly as workload migration suggestions
2. **Given** the CLI receives recommendations with SCHEDULE action type,
   **When** filtering by action type,
   **Then** users can filter to see only schedule-based cost savings (dev/test environments)
3. **Given** the CLI receives recommendations with OTHER action type,
   **When** displaying results,
   **Then** these are shown as "Other Recommendations" with provider-specific details preserved

---

### Edge Cases

- What happens when an older plugin receives a recommendation request with a filter for new
  action types?
  - The filter should be ignored or return empty results gracefully
- How does the system handle unknown enum values in proto3?
  - Proto3 preserves unknown enum values as their numeric representation; the SDK should
    handle this gracefully
- What happens if a plugin returns action type 0 (UNSPECIFIED)?
  - This indicates the plugin couldn't categorize the recommendation; core should handle it
    as "uncategorized"

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The `RecommendationActionType` enum MUST include 5 new values:
  MIGRATE (7), CONSOLIDATE (8), SCHEDULE (9), REFACTOR (10), OTHER (11)
- **FR-002**: The enum numbering MUST continue sequentially from the existing values (7-11)
  to maintain proto3 compatibility
- **FR-003**: The Go SDK MUST be regenerated to include the new enum constants
- **FR-004**: The generated Go code MUST include string representation methods for all new
  enum values
- **FR-005**: Existing enum values (UNSPECIFIED through DELETE_UNUSED, values 0-6) MUST
  remain unchanged
- **FR-006**: The proto definition MUST include documentation comments describing each new
  action type's intended use case

### Key Entities

- **RecommendationActionType**: Enumeration representing the category of action a cost
  recommendation suggests. Extended from 7 values (0-6) to 12 values (0-11).
- **Recommendation**: Message type that uses RecommendationActionType to categorize the
  suggested optimization action.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All 5 new enum values are available in the generated Go SDK and can be
  referenced in code
- **SC-002**: Existing plugins compile successfully against the updated SDK without code
  changes
- **SC-003**: Round-trip serialization of all 12 enum values (including new ones) preserves
  values correctly
- **SC-004**: Plugin conformance tests pass for plugins using any of the 12 action types
- **SC-005**: The proto definition is accepted by `buf lint` without warnings or errors

## Assumptions

- The enum extension follows Proto3 semantics where unknown values are preserved as their
  numeric representation
- The `OTHER` type serves as a catch-all for provider-specific recommendations that don't
  fit other categories
- No changes are needed to the `Recommendation` message structure itself - only the enum
  is extended
- Release versioning (v0.4.9 or v0.5.0) will be determined by release-please based on
  conventional commits

## Out of Scope

- CLI implementation changes (handled in finfocus-core)
- Plugin implementation updates (each plugin repository handles their own updates)
- Changes to other enums or message types in the proto definition
- Changes to JSON schema (this feature is proto-only)
