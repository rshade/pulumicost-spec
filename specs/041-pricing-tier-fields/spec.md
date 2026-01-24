# Feature Specification: Add Pricing Tier & Probability Fields

**Feature Branch**: `041-pricing-tier-fields`
**Created**: 2026-01-20
**Status**: Draft
**Input**: User description: "Add Pricing Tier & Probability Fields... We need a way to tell the Core why the price is low..."

## Session 2026-01-20

- Q: How should we name the new pricing tier enum to avoid collision with the existing volume-discount
  `PricingTier` message? → A: Reuse the existing `FocusPricingCategory` (Standard, Committed, Dynamic) to
  maintain alignment with FOCUS 1.2.
- Q: Which RPC responses should include these new fields? → A: Both `EstimateCostResponse` and
  `GetProjectedCostResponse`.
- Q: Should the risk enum use the `Focus` prefix? → A: No, use `SpotInterruptionRisk` as it is a custom
  extension not part of the official FOCUS spec.
- Q: Should we use an enumeration or a numeric score for risk? → A: Use a `double spot_interruption_risk_score`
  (0.0 - 1.0) to allow for precise probability data from provider APIs.

## User Scenarios & Testing

### User Story 1 - Identify Spot Instance Risk (Priority: P1)

A user viewing cost data needs to immediately understand if a low price is due to using Spot instances and
what the associated risk is, so they can make informed decisions about workload placement.

**Why this priority**: Critical for "green" or "efficient" ops where cost savings must be balanced against reliability.

**Independent Test**: Can be tested by creating a mock cost record with `pricing_category = DYNAMIC` and
`spot_interruption_risk_score = 0.8` and verifying the data is correctly serialized and accessible in the SDK.

**Acceptance Scenarios**:

1. **Given** a cost record for a Spot instance with 80% interruption probability, **When** `EstimateCost` or
   `GetProjectedCost` is called, **Then** the response contains `pricing_category = DYNAMIC` and
   `spot_interruption_risk_score = 0.8`.
2. **Given** a standard on-demand resource, **When** processed, **Then** `pricing_category` is `STANDARD`
   and `spot_interruption_risk_score` is 0.0 (or omitted).

---

### User Story 2 - Explain Cost Basis (Priority: P2)

A user investigating unexpected cost variances needs to know if a resource is covered by a Reserved Instance
(RI) or Savings Plan, explaining why its effective cost is lower than on-demand rates.

**Why this priority**: Helps users trust the data by explaining "why" a cost is what it is, reducing confusion.

**Independent Test**: Create a record with `pricing_category = COMMITTED` and verify it can be distinguished
from `STANDARD` pricing.

**Acceptance Scenarios**:

1. **Given** a resource covered by a reservation, **When** `EstimateCost` or `GetProjectedCost` is called,
   **Then** the response contains `pricing_category = COMMITTED`.
2. **Given** the same resource without reservation, **When** processed, **Then** `pricing_category` is `STANDARD`.

### Edge Cases

- **Risk without Spot**: If `spot_interruption_risk_score` is greater than 0.0 but `pricing_category` is NOT
  `DYNAMIC`, the risk score should be ignored by consumers.
- **Unspecified Values**: If `pricing_category` is `UNSPECIFIED` (0), consumers should default to `STANDARD` behavior.

## Requirements

### Functional Requirements

- **FR-001**: The implementation MUST reuse the existing `FocusPricingCategory` enumeration: `UNSPECIFIED`
  (0), `STANDARD` (1), `COMMITTED` (2), `DYNAMIC` (3).
- **FR-002**: The specification MUST define a `double spot_interruption_risk_score` field, where 0.0
  represents no risk and 1.0 represents 100% probability of interruption.
- **FR-003**: `EstimateCostResponse` MUST include `pricing_category` and `spot_interruption_risk_score` fields.
- **FR-004**: `GetProjectedCostResponse` MUST include `pricing_category` and `spot_interruption_risk_score` fields.
- **FR-005**: The `spot_interruption_risk_score` field SHOULD be considered relevant only when `pricing_category` is `DYNAMIC`.

### Key Entities

- **FocusPricingCategory**: (Existing) Categorized pricing model (Standard, Committed/Reserved, Dynamic/Spot).
- **SpotInterruptionRiskScore**: Probability-based reliability risk for Dynamic/Spot instances (0.0 to 1.0).

## Success Criteria

### Measurable Outcomes

- **SC-001**: The data model explicitly defines the four pricing tiers and three risk levels, verified by schema inspection.
- **SC-002**: Client libraries can programmatically set and retrieve pricing tier and risk information without data loss.
- **SC-003**: A "High Risk Spot" scenario can be fully represented in a single data object and validated against the schema.
