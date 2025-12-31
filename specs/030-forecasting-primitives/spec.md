# Feature Specification: Forecasting Primitives

**Feature Branch**: `030-forecasting-primitives`
**Created**: 2025-12-30
**Status**: Draft
**Input**: Update resource.proto to include GrowthType (Enum: NONE, LINEAR, EXPONENTIAL)
and GrowthRate (double) to enable forward-looking cost projections.

## Clarifications

### Session 2025-12-30

- Q: Where should growth fields be placed? → A: Both ResourceDescriptor (default) and
  GetProjectedCostRequest (override)
- Q: Should growth_rate have explicit bounds? → A: Lower bound only (>= -1.0), no upper bound

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Project Future Costs with Growth Assumptions (Priority: P1)

As a FinOps practitioner, I want to specify growth assumptions (linear or exponential) for
resources so that I can project future costs based on anticipated usage growth patterns.

**Why this priority**: Cost projections are the primary use case for forecasting primitives.
Without the ability to model growth, organizations cannot plan budgets or anticipate
infrastructure scaling costs.

**Independent Test**: Can be fully tested by providing a resource descriptor with growth
parameters and verifying that projected cost responses incorporate the growth model into
their calculations.

**Acceptance Scenarios**:

1. **Given** a resource descriptor with `growth_type = LINEAR` and `growth_rate = 0.10`,
   **When** requesting projected costs for 12 months ahead,
   **Then** the system returns cost projections that increase by 10% each month linearly.

2. **Given** a resource descriptor with `growth_type = EXPONENTIAL` and `growth_rate = 0.05`,
   **When** requesting projected costs for 6 months ahead,
   **Then** the system returns cost projections that compound at 5% per period.

3. **Given** a resource descriptor with `growth_type = NONE`,
   **When** requesting projected costs,
   **Then** the system returns cost projections based on current usage without any growth
   assumptions applied.

---

### User Story 2 - Default Behavior for Resources Without Growth Specifications (Priority: P2)

As a plugin developer, I want resources without explicit growth specifications to behave
predictably so that existing integrations continue to work without modification.

**Why this priority**: Backward compatibility ensures existing users are not disrupted.
New fields must have safe defaults that maintain current behavior.

**Independent Test**: Can be tested by sending requests without growth fields and verifying
responses match pre-feature behavior.

**Acceptance Scenarios**:

1. **Given** a resource descriptor with no growth fields specified,
   **When** requesting projected costs,
   **Then** the system treats this as `growth_type = NONE` with no growth applied.

2. **Given** a resource descriptor with only `growth_type` specified (no rate),
   **When** the growth type is `NONE`,
   **Then** the system ignores the missing growth rate and returns projections without growth.

---

### User Story 3 - Validate Growth Parameters (Priority: P3)

As a FinOps practitioner, I want invalid growth parameters to be rejected with clear error
messages so that I can correct configuration mistakes before relying on projections.

**Why this priority**: Input validation prevents incorrect projections from being silently
used for budget planning, which could lead to significant financial planning errors.

**Independent Test**: Can be tested by sending requests with invalid growth parameters and
verifying appropriate error responses are returned.

**Acceptance Scenarios**:

1. **Given** a resource descriptor with `growth_type = LINEAR` but no `growth_rate` specified,
   **When** requesting projected costs,
   **Then** the system returns an error indicating the growth rate is required for
   linear growth.

2. **Given** a resource descriptor with `growth_type = EXPONENTIAL` and a negative
   `growth_rate`,
   **When** requesting projected costs,
   **Then** the system accepts this as a valid declining growth model (negative growth
   represents expected usage decrease).

3. **Given** a resource descriptor with `growth_type = NONE` and a `growth_rate` value
   specified,
   **When** requesting projected costs,
   **Then** the system ignores the growth rate, returns projections without growth, and
   logs a DEBUG message noting that `growth_rate` was ignored for `GROWTH_TYPE_NONE`.

---

### Edge Cases

- What happens when growth rate exceeds 100% per period (e.g., `growth_rate = 2.0`)?
  The system accepts it as valid for rapid scaling scenarios. The SDK MUST log a warning
  via zerolog at WARN level when `growth_rate > 1.0` to alert operators of unusually high
  growth assumptions.
- How does the system handle growth rate of exactly 0.0?
  The system treats this identically to `growth_type = NONE` (no growth applied).
- What happens with extremely long projection periods (e.g., 120 months) with exponential
  growth? The system calculates the projection. The SDK SHOULD log an INFO message when
  projection periods exceed 36 months with exponential growth, noting that accuracy
  decreases for distant projections. No calculation is blocked.
- How are growth fields handled for actual cost queries (historical data)?
  Growth parameters are only relevant for projected costs; they are ignored for actual
  cost queries.
- What happens when growth_rate is exactly -1.0? The system projects costs declining to
  zero over the projection period (100% decline). Values below -1.0 are rejected as invalid.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide a `GrowthType` enumeration with values:
  `GROWTH_TYPE_UNSPECIFIED`, `GROWTH_TYPE_NONE`, `GROWTH_TYPE_LINEAR`,
  `GROWTH_TYPE_EXPONENTIAL`.

- **FR-002**: System MUST provide a `growth_rate` field as an optional double value
  representing the growth rate per projection period.

- **FR-003**: System MUST treat `GROWTH_TYPE_UNSPECIFIED` and `GROWTH_TYPE_NONE`
  identically, applying no growth to projections.

- **FR-004**: For `GROWTH_TYPE_LINEAR`, the system MUST apply the growth rate as an
  additive percentage per period (e.g., rate=0.10 adds 10% of the base cost each period).

- **FR-005**: For `GROWTH_TYPE_EXPONENTIAL`, the system MUST apply the growth rate as a
  compounding factor per period (e.g., rate=0.10 multiplies by 1.10 each period).

- **FR-006**: System MUST return an `InvalidArgument` error when `growth_type` is `LINEAR`
  or `EXPONENTIAL` but `growth_rate` is not provided.

- **FR-007**: System MUST accept negative `growth_rate` values for modeling declining
  usage scenarios, with a lower bound of -1.0 (representing 100% decline to zero cost).

- **FR-007a**: System MUST return an `InvalidArgument` error when `growth_rate` is less
  than -1.0, as this would produce mathematically nonsensical negative costs.

- **FR-008**: System MUST ignore growth parameters for actual cost queries
  (`GetActualCost` RPC) as growth is only applicable to forward-looking projections.

- **FR-009**: System MUST maintain backward compatibility by treating requests without
  growth fields identically to current behavior.

- **FR-010**: When growth parameters are specified on both `ResourceDescriptor` and
  `GetProjectedCostRequest`, the request-level parameters MUST override the resource-level
  defaults.

### Key Entities

- **GrowthType (Enum)**: Represents the mathematical model used for projecting cost growth.
  Values indicate no growth, linear (additive) growth, or exponential (compounding) growth.

- **growth_rate (double)**: The rate of growth per projection period expressed as a decimal
  (0.10 = 10%). Valid range: >= -1.0 (no upper bound). Positive values indicate growth,
  zero indicates no change, negative values (down to -1.0) indicate decline.

- **ResourceDescriptor (Extended)**: The existing resource descriptor message extended with
  optional growth parameters (`growth_type`, `growth_rate`) as default values for the resource.

- **GetProjectedCostRequest (Extended)**: The existing request message extended with optional
  growth parameters that override ResourceDescriptor defaults when specified. This enables
  different projection scenarios for the same resource in a single session.

## Assumptions

- Growth rate is expressed as a decimal fraction, not a percentage (0.10 means 10%, not
  0.1%).
- The projection period for growth calculations aligns with the standard monthly projection
  period used by `GetProjectedCost`.
- Growth parameters are optional fields to maintain backward compatibility with existing
  clients.
- Plugin implementations that do not support forecasting may ignore these fields and return
  standard projections.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Users can specify growth assumptions for 100% of resources that support
  projected cost queries.

- **SC-002**: Cost projections with linear growth show predictable, additive increases that
  match the specified growth rate within 0.01% accuracy.

- **SC-003**: Cost projections with exponential growth show compounding behavior that
  matches the specified growth rate within 0.01% accuracy.

- **SC-004**: Existing clients without growth parameters continue to receive identical
  responses (100% backward compatibility).

- **SC-005**: Invalid growth parameter combinations are rejected with clear error messages
  in less than 100ms.

- **SC-006**: Plugin developers can implement forecasting support using the SDK in under
  30 minutes with provided documentation and examples.
