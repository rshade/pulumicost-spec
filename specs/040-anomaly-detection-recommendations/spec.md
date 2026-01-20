# Feature Specification: Anomaly Detection via Recommendations

**Feature Branch**: `040-anomaly-detection-recommendations`
**Created**: 2026-01-19
**Status**: Draft
**Input**: GitHub Issue #315 - Add ANOMALY category and INVESTIGATE action for cost anomaly detection

## User Scenarios & Testing

### User Story 1 - Unified Actionable Cost Insights View (Priority: P1)

A FinOps practitioner needs to see all actionable cost items in a single view, including both
optimization recommendations (rightsize, terminate, purchase commitment) and cost anomalies
(unusual spending patterns). Currently, they would need to query multiple endpoints to get a
complete picture of items requiring attention.

**Why this priority**: This is the core value proposition - unifying anomalies and recommendations
into a single query enables simpler dashboards, reduces integration complexity, and matches how
AWS presents anomalies in their Cost Management console.

**Independent Test**: Can be tested by calling `GetRecommendations` and verifying that both
traditional optimization recommendations and anomaly-type recommendations can be returned in
the same response, delivering a complete "action items" view.

**Acceptance Scenarios**:

1. **Given** a plugin connected to a cost management service that detects anomalies,
   **When** a user calls `GetRecommendations` without category filtering,
   **Then** both optimization recommendations (COST, PERFORMANCE, etc.) and anomaly
   recommendations (ANOMALY category) are returned in the same response.

2. **Given** a plugin that supports anomaly detection,
   **When** a user calls `GetRecommendations` with `category=ANOMALY`,
   **Then** only recommendations with `RECOMMENDATION_CATEGORY_ANOMALY` are returned.

3. **Given** a plugin that does NOT support anomaly detection,
   **When** a user calls `GetRecommendations`,
   **Then** zero recommendations with `RECOMMENDATION_CATEGORY_ANOMALY` are returned
   (no error is raised - the plugin simply never generates anomaly recommendations).

---

### User Story 2 - Anomaly Triage by Confidence Score (Priority: P2)

A FinOps engineer wants to prioritize which anomalies to investigate first. They want to filter
anomalies by confidence score to focus on high-confidence detections and avoid spending time on
false positives.

**Why this priority**: Anomaly detection inherently involves statistical confidence levels.
Leveraging the existing `confidence_score` field and `min_confidence_score` filter enables
automated triage workflows without requiring new filtering infrastructure.

**Independent Test**: Can be tested by returning anomaly recommendations with varying confidence
scores and using the existing `min_confidence_score` filter to verify correct filtering behavior.

**Acceptance Scenarios**:

1. **Given** multiple anomaly recommendations with confidence scores of 0.3, 0.6, and 0.9,
   **When** a user calls `GetRecommendations` with `min_confidence_score=0.5`,
   **Then** only the anomaly recommendations with scores 0.6 and 0.9 are returned.

2. **Given** an anomaly recommendation from a provider that supports confidence scoring,
   **When** the recommendation is returned,
   **Then** the `confidence_score` field contains the provider's anomaly detection confidence
   (0.0-1.0 scale).

---

### User Story 3 - Investigate Anomalous Spending (Priority: P2)

A cloud engineer receives an anomaly recommendation and needs to understand what action to take.
The recommendation should clearly indicate that investigation is required rather than suggesting
a specific remediation like rightsizing or termination.

**Why this priority**: Traditional recommendations have clear actions (rightsize, terminate, buy
commitment). Anomalies are different - they indicate something unusual happened that requires
human investigation. The `INVESTIGATE` action communicates this distinction clearly.

**Independent Test**: Can be tested by verifying that anomaly recommendations return with
`action_type=RECOMMENDATION_ACTION_TYPE_INVESTIGATE` and that the description field provides
sufficient context for investigation.

**Acceptance Scenarios**:

1. **Given** a cost anomaly is detected (e.g., 150% above baseline),
   **When** the anomaly is returned as a recommendation,
   **Then** the `action_type` is `RECOMMENDATION_ACTION_TYPE_INVESTIGATE`.

2. **Given** an anomaly recommendation,
   **When** a user reviews the recommendation,
   **Then** the `description` field contains actionable information (e.g., "Unusual spending
   detected: 150% above baseline for service X in region Y").

3. **Given** an anomaly recommendation,
   **When** a user needs additional context,
   **Then** the `metadata` map contains anomaly-specific details (e.g., baseline amount,
   deviation percentage, detection timestamp).

---

### User Story 4 - Exclude Anomalies from Optimization Workflows (Priority: P3)

An automated cost optimization system processes recommendations to generate Terraform/IaC changes.
This system should be able to exclude anomaly recommendations since they cannot be automatically
remediated and require human investigation.

**Why this priority**: Enables clean separation between automatable optimization recommendations
and human-required investigation items in downstream workflows.

**Independent Test**: Can be tested by using category filtering to exclude ANOMALY category from
results, ensuring automation systems only receive actionable recommendations.

**Acceptance Scenarios**:

1. **Given** a mix of optimization recommendations and anomaly recommendations,
   **When** a user calls `GetRecommendations` with a filter excluding ANOMALY category,
   **Then** only non-anomaly recommendations are returned.

2. **Given** an automated remediation system,
   **When** it receives an anomaly recommendation,
   **Then** the `action_type=INVESTIGATE` signals that automated remediation is not appropriate.

---

### Edge Cases

- What happens when a plugin returns ANOMALY category but uses action_type other than INVESTIGATE?
  **Answer**: This is technically valid but discouraged. The semantic pairing of ANOMALY category
  with INVESTIGATE action is a convention, not a hard constraint. Plugins may use other actions
  if they have specific remediation suggestions for certain anomaly types.

- How does the system handle negative `estimated_savings` for cost overspend anomalies?
  **Answer**: The `estimated_savings` field in `RecommendationImpact` can be negative when
  representing the additional cost from an anomaly (i.e., how much the anomaly is costing above
  baseline). This is a valid semantic interpretation.

- What happens if both category=ANOMALY and category=COST filters are applied?
  **Answer**: The filter accepts a single category value. Multiple category filtering would
  require a future API enhancement (e.g., repeated category field or category exclusion).

- How should plugins handle anomalies that span multiple resources?
  **Answer**: The `ResourceRecommendationInfo` can identify a service-level or account-level
  scope rather than a specific resource. The `resource_type` can be set to a grouping concept
  (e.g., "account", "service").

## Requirements

### Functional Requirements

- **FR-001**: System MUST add `RECOMMENDATION_CATEGORY_ANOMALY = 5` to the `RecommendationCategory`
  enum in `costsource.proto` (note: actual value assignment based on existing enum sequence).

- **FR-002**: System MUST add `RECOMMENDATION_ACTION_TYPE_INVESTIGATE = 12` to the
  `RecommendationActionType` enum in `costsource.proto` (note: actual value assignment based on
  existing enum sequence).

- **FR-003**: Plugins returning anomaly recommendations MUST set `category` to
  `RECOMMENDATION_CATEGORY_ANOMALY`.

- **FR-004**: Plugins returning anomaly recommendations SHOULD set `action_type` to
  `RECOMMENDATION_ACTION_TYPE_INVESTIGATE` unless a specific remediation action is appropriate.

- **FR-005**: The existing `confidence_score` field MUST be populated with the provider's anomaly
  detection confidence when available (0.0-1.0 scale).

- **FR-006**: The existing `estimated_savings` field in `RecommendationImpact` MUST represent the
  anomaly amount (deviation from baseline), which MAY be negative for cost overspend anomalies.

- **FR-007**: The `description` field MUST provide human-readable context about the anomaly
  (e.g., baseline, deviation percentage, affected service/resource).

- **FR-008**: Plugins that do not support anomaly detection MUST NOT return recommendations with
  `RECOMMENDATION_CATEGORY_ANOMALY` (they simply omit such recommendations).

- **FR-009**: The existing `RecommendationFilter.category` field MUST support filtering for
  `RECOMMENDATION_CATEGORY_ANOMALY` with no code changes (existing enum filtering logic applies).

- **FR-010**: SDK MUST regenerate Go code from updated proto definitions via `make generate`.

- **FR-011**: SDK documentation MUST document the semantic mapping of existing Recommendation
  fields to anomaly use cases (see issue description table).

### Key Entities

- **RecommendationCategory Enum**: Extended with ANOMALY value to classify cost anomaly
  recommendations.

- **RecommendationActionType Enum**: Extended with INVESTIGATE value to indicate that human
  investigation is required rather than automated remediation.

- **Recommendation Message**: No structural changes - existing fields (`confidence_score`,
  `estimated_savings`, `description`, `metadata`) serve anomaly use cases.

- **RecommendationFilter Message**: No changes - existing `category` filter supports the new
  ANOMALY value automatically.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can retrieve all actionable cost items (recommendations + anomalies) via
  a single `GetRecommendations` call.

- **SC-002**: Users can filter to see only anomaly recommendations by setting `category=ANOMALY`
  in the request filter.

- **SC-003**: Users can exclude anomaly recommendations from automated workflows by filtering
  on `action_type != INVESTIGATE`.

- **SC-004**: Anomaly recommendations contain sufficient context for investigation (description,
  confidence score, impact details).

- **SC-005**: All existing recommendation workflows continue to function unchanged (backward
  compatible addition).

- **SC-006**: Plugin developers can implement anomaly support by returning recommendations with
  the new enum values without any changes to the Recommendation message structure.

## Assumptions

- **A-001**: The existing `Recommendation` message fields are sufficient to represent anomaly
  data without structural changes.

- **A-002**: Plugins have access to anomaly detection data from their underlying cost management
  services (AWS Cost Anomaly Detection, Azure Cost Management Anomalies, etc.).

- **A-003**: Anomaly confidence scores from different providers can be normalized to a 0.0-1.0
  scale for consistent filtering.

- **A-004**: Negative `estimated_savings` values are acceptable for representing cost overspend
  (the amount above baseline).

- **A-005**: The semantic convention pairing ANOMALY category with INVESTIGATE action will be
  documented but not enforced at the proto validation level.

## Out of Scope

- Creating a dedicated `GetAnomalies` RPC (rejected alternative per issue description).
- Adding an `AnomalyDetails` message to the action oneof (potential future enhancement).
- Implementing anomaly detection logic within the SDK (plugins implement this by integrating
  with backend services).
- Multi-category filtering in `RecommendationFilter` (would require separate enhancement).
- Time-series anomaly visualization or historical tracking.
