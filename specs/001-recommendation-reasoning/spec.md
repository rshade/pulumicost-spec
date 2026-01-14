# Feature Specification: Standardized Recommendation Reasoning Metadata

**Feature Branch**: `001-recommendation-reasoning`
**Created**: 2026-01-13
**Status**: Draft
**Input**: User description: "Define machine-readable reason codes (e.g., OVER_PROVISIONED_CPU) in recommendation metadata to bridge raw data and user understanding. Boundary Check: Transport-only; logic remains in the plugin/upstream provider."

## Clarifications

### Session 2026-01-13

- Q: How should multiple reasons be handled (e.g., Idle + Old Gen)? → A: Use a structured approach with a specific `primary_reason` field and a `secondary_reasons` list to distinguish the main driver from contributing factors.
- Q: Which categories should the initial set of reason codes include? → A: Use a comprehensive set including `UNSPECIFIED`, `OVER_PROVISIONED`, `UNDER_PROVISIONED` (Performance), `IDLE`, `REDUNDANT`, and `OBSOLETE_GENERATION`.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Developer Maps Upstream Reasons (Priority: P1)

A plugin developer implementing the FinFocus spec for a cloud provider (e.g., AWS, GCP) needs to map the provider-specific recommendation reason (e.g., "Low CPU Utilization") to a standardized FinFocus reason code. This ensures that downstream tools can understand the "why" of a recommendation regardless of the source.

**Why this priority**: This is the core function of the feature—enabling the transport of standardized metadata. Without this, the feature has no value.

**Independent Test**: Can be tested by creating a mock plugin response that includes the new reason code and verifying it is correctly serialized and accessible.

**Acceptance Scenarios**:

1. **Given** a plugin receives a recommendation from an upstream API with a provider-specific reason (e.g., "Idle"), **When** the plugin constructs the FinFocus Recommendation object, **Then** it can set a `Reason` field to a standardized enum value (e.g., `REASON_IDLE`).
2. **Given** a plugin encounters a provider reason that has no direct standard mapping, **When** it constructs the Recommendation object, **Then** it can set the `Reason` field to `REASON_UNSPECIFIED` or similar fallback.

---

### User Story 2 - Consumer Tool Displays Reasoning (Priority: P2)

A developer building a dashboard or CLI tool using the FinFocus SDK wants to display the rationale behind a cost-saving recommendation. They want to show a consistent label (e.g., "Over-provisioned") regardless of whether the recommendation came from AWS, Azure, or Kubernetes.

**Why this priority**: This validates the end-user value of the standardization.

**Independent Test**: Can be tested by consuming a recommendation object with the reason code set and verifying the consumer code can switch/case on the enum to display a UI string.

**Acceptance Scenarios**:

1. **Given** a Recommendation object with `Reason: REASON_OVER_PROVISIONED`, **When** the consumer application processes it, **Then** it can deterministically identify the reason type without parsing arbitrary text strings.

### Edge Cases

- What happens when a new reason code is added to the spec but the consumer SDK is outdated? (Protobuf handling of unknown enum values).
- How does the system handle complex reasons that might fit multiple categories (e.g., both "Idle" and "Old Generation")? (Resolved: Use `primary_reason` and `secondary_reasons`).

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST define a standardized set of machine-readable reason codes (e.g., enumerated type) for recommendations.
- **FR-002**: The set of codes MUST include: `UNSPECIFIED`, `OVER_PROVISIONED`, `UNDER_PROVISIONED`, `IDLE`, `REDUNDANT`, and `OBSOLETE_GENERATION`.
- **FR-003**: The Recommendation message structure MUST include a `primary_reason` field and a `secondary_reasons` list to transport these codes.
- **FR-004**: The system MUST allow for an "Other" or "Unspecified" reason code for cases not covered by the standard set.
- **FR-005**: The SDK MUST provide language-specific bindings for these reason codes.

### Key Entities

- **RecommendationReason**: An enumeration of standard reason codes (e.g., `OVER_PROVISIONED`, `IDLE`).
- **Recommendation**: The existing message structure, updated to include:
  - `primary_reason` (RecommendationReason): The main driver for the recommendation.
  - `secondary_reasons` (List<RecommendationReason>): Contributing factors.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: The schema definition successfully compiles/validates with the new `RecommendationReason` structure.
- **SC-002**: The generated SDKs/libraries include the `RecommendationReason` type and constants.
- **SC-003**: A dummy plugin response containing a specific reason code (e.g., `REASON_IDLE`) can be serialized and deserialized back to the exact same value.
- **SC-004**: Existing plugins or code that do not populate this field still function (backward compatibility), defaulting to an Unspecified value.