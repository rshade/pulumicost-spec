# Feature Specification: Contextual FinOps Validation

**Feature Branch**: `027-finops-validation`
**Created**: 2025-12-25
**Status**: Draft
**Input**: User description: "feat(pluginsdk): implement contextual FinOps validation"

## Clarifications

### Session 2025-12-25

- Q: What happens when all cost fields are zero (free tier)? → A: Valid - Zero-cost records
  pass validation (free tier is legitimate usage)
- Q: What format should structured error messages use? → A: Structured fields - Error type
  with FieldName, Constraint, ActualValue, ExpectedValue fields

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Cost Relationship Validation (Priority: P1)

Plugin developers need to validate that cost records have logically consistent relationships
between financial amounts. For example, EffectiveCost (after discounts) should never exceed
BilledCost (what the customer actually pays), and ListCost (undiscounted price) should be
greater than or equal to EffectiveCost.

**Why this priority**: Cost relationship errors are the most common data quality issues in
FinOps. Invalid cost relationships can lead to incorrect billing analysis, budget forecasts,
and FinOps reports. This is the foundational validation that all other contextual checks
depend upon.

**Independent Test**: Can be fully tested by submitting a FocusCostRecord with various cost
field combinations and verifying that invalid relationships are rejected with clear error
messages.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord with EffectiveCost greater than BilledCost,
   **When** validation is run,
   **Then** an error is returned indicating the cost relationship violation.
2. **Given** a FocusCostRecord with ListCost less than EffectiveCost,
   **When** validation is run,
   **Then** an error is returned indicating the pricing hierarchy violation.
3. **Given** a FocusCostRecord with valid cost relationships (ListCost >= BilledCost >= EffectiveCost),
   **When** validation is run,
   **Then** validation passes successfully.
4. **Given** a FocusCostRecord with negative EffectiveCost representing credits/refunds,
   **When** validation is run,
   **Then** validation handles the credit scenario appropriately.

---

### User Story 2 - Commitment Discount Field Consistency (Priority: P1)

Plugin developers need to validate that commitment discount fields are internally consistent.
When a CommitmentDiscountId is present, related fields like CommitmentDiscountStatus,
CommitmentDiscountType, and CommitmentDiscountCategory should have appropriate values.

**Why this priority**: Commitment discounts represent significant cost optimization
opportunities (Reserved Instances, Savings Plans). Inconsistent commitment data leads to
incorrect utilization analysis and missed optimization recommendations.

**Independent Test**: Can be fully tested by submitting records with various commitment field
combinations and verifying consistency rules are enforced.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord with CommitmentDiscountId set but CommitmentDiscountStatus
   unspecified,
   **When** validation is run and ChargeCategory is Usage,
   **Then** an error is returned indicating missing status field.
2. **Given** a FocusCostRecord with CommitmentDiscountStatus set to "Used" but no
   CommitmentDiscountId,
   **When** validation is run,
   **Then** an error is returned indicating orphaned status.
3. **Given** a FocusCostRecord with consistent commitment discount fields,
   **When** validation is run,
   **Then** validation passes successfully.
4. **Given** a FocusCostRecord with CommitmentDiscountCategory of "Spend" and appropriate
   spend-based fields,
   **When** validation is run,
   **Then** validation passes successfully.

---

### User Story 3 - Pricing Model Consistency (Priority: P2)

Plugin developers need to validate that pricing-related fields are consistent with the charge
category and pricing model. For example, usage-based charges should have valid pricing
quantity and unit price relationships.

**Why this priority**: Accurate pricing model validation enables reliable cost allocation and
chargeback mechanisms. This builds on the foundation of cost relationship validation.

**Independent Test**: Can be fully tested by creating records with various pricing scenarios
and verifying that inconsistent pricing relationships are detected.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord with PricingQuantity > 0 but PricingUnit empty,
   **When** validation is run,
   **Then** an error is returned indicating missing unit specification.
2. **Given** a FocusCostRecord where ContractedCost significantly differs from
   ContractedUnitPrice × PricingQuantity,
   **When** validation is run,
   **Then** an error is returned indicating calculation mismatch.
3. **Given** a FocusCostRecord with ChargeCategory "Purchase" and appropriate purchase fields,
   **When** validation is run,
   **Then** validation passes successfully.

---

### User Story 4 - Capacity Reservation Consistency (Priority: P2)

Plugin developers need to validate that capacity reservation fields follow the same
consistency patterns as commitment discounts - when CapacityReservationId is present,
CapacityReservationStatus should be appropriately set.

**Why this priority**: Capacity reservations are a key cost optimization mechanism, and
consistent data enables accurate utilization tracking.

**Independent Test**: Can be fully tested by submitting records with capacity reservation
field combinations.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord with CapacityReservationId set but CapacityReservationStatus
   unspecified,
   **When** validation is run and ChargeCategory is Usage,
   **Then** an error is returned.
2. **Given** a FocusCostRecord with consistent capacity reservation fields,
   **When** validation is run,
   **Then** validation passes successfully.

---

### User Story 5 - Validation Error Aggregation (Priority: P3)

Plugin developers need the option to receive all validation errors at once rather than
fail-fast behavior. This enables comprehensive data quality reports and batch correction
workflows.

**Why this priority**: For large datasets, collecting all errors in a single pass is more
efficient than fixing errors one at a time. This enhances the developer experience for data
quality workflows.

**Independent Test**: Can be fully tested by submitting a record with multiple validation
errors and verifying all errors are returned.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord with multiple validation errors,
   **When** validation is run in aggregate mode,
   **Then** all errors are returned in a structured collection.
2. **Given** a FocusCostRecord with multiple errors,
   **When** validation is run in fail-fast mode (default),
   **Then** only the first error is returned.

---

### Edge Cases

- Zero-cost records (free tier): Valid and pass validation without error.
- How does validation handle negative costs (credits, refunds, corrections)?
- What happens when optional fields are missing vs explicitly set to zero/empty?
- How does validation behave with ChargeClass "Correction" which has different rules?
- What happens with mixed currency scenarios (BillingCurrency vs PricingCurrency)?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST validate that EffectiveCost does not exceed BilledCost when both
  are positive (excluding correction charges).
- **FR-002**: System MUST validate that ListCost is greater than or equal to EffectiveCost
  when both are present and positive.
- **FR-003**: System MUST validate that CommitmentDiscountStatus is set when
  CommitmentDiscountId is present and ChargeCategory is "Usage".
- **FR-004**: System MUST validate that CommitmentDiscountId is present when
  CommitmentDiscountStatus indicates a commitment relationship.
- **FR-005**: System MUST validate that CapacityReservationStatus is set when
  CapacityReservationId is present and ChargeCategory is "Usage".
- **FR-006**: System MUST validate that PricingUnit is specified when PricingQuantity is
  greater than zero.
- **FR-007**: System MUST provide structured error messages with discrete fields (FieldName,
  Constraint, ActualValue, ExpectedValue) that callers can inspect programmatically.
- **FR-008**: System MUST maintain backward compatibility with existing ValidateFocusRecord
  function.
- **FR-009**: System MUST support both fail-fast (default) and aggregate-all-errors
  validation modes.
- **FR-010**: System MUST exclude ChargeClass "Correction" from certain validation rules per
  FOCUS specification.
- **FR-011**: System MUST handle negative cost values appropriately (credits, refunds do not
  violate cost hierarchies).
- **FR-012**: System MUST achieve zero-allocation validation on the happy path (valid
  records).

### Key Entities

- **FocusCostRecord**: The primary record type containing all cost and commitment fields to
  be validated.
- **ValidationError**: A structured error type with four discrete fields: FieldName (the
  violating field), Constraint (the rule violated), ActualValue (what was found),
  ExpectedValue (what was expected). Enables programmatic error inspection.
- **ValidationMode**: An enumeration or option to control fail-fast vs aggregate validation
  behavior.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Validation of a valid FocusCostRecord completes in under 100 nanoseconds with
  zero memory allocations.
- **SC-002**: All 12 functional requirements have corresponding test cases with 100% code
  coverage.
- **SC-003**: ValidationError instances expose FieldName, Constraint, ActualValue, and
  ExpectedValue as discrete fields accessible programmatically.
- **SC-004**: Existing tests using ValidateFocusRecord continue to pass without modification.
- **SC-005**: Plugin developers can identify and fix data quality issues in their FinOps
  datasets using the validation output.
- **SC-006**: Validation rules align with FOCUS 1.2/1.3 specification requirements documented
  at focus.finops.org.

## Assumptions

- Cost values are represented as float64 following the existing FocusCostRecord protobuf
  definition.
- Floating-point comparison will use appropriate tolerance (existing
  `contractedCostTolerance = 0.0001`).
- The existing `ValidateFocusRecord` function in `focus_conformance.go` will be extended
  rather than replaced.
- Credit/refund scenarios are identified by ChargeClass "Correction" or negative cost values.
- The FOCUS specification rules at focus.finops.org are the authoritative source for business
  logic constraints.

## Out of Scope

- Cross-record validation (validating relationships between multiple cost records).
- Time-series validation (validating cost trends over time).
- Provider-specific validation rules beyond FOCUS specification.
- Currency conversion or multi-currency validation.
- Real-time streaming validation (batch validation only).
