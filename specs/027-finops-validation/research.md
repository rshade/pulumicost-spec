# Research: Contextual FinOps Validation

**Feature**: 027-finops-validation
**Date**: 2025-12-25

## Research Tasks

### 1. FOCUS Specification Cost Relationships

**Decision**: Implement cost hierarchy validation: ListCost >= BilledCost >= EffectiveCost

**Rationale**: The FOCUS specification defines these cost fields with clear semantics:

- **ListCost**: The cost at list price before any discounts (highest)
- **BilledCost**: The amount actually billed to the customer
- **EffectiveCost**: The actual cost after all discounts and adjustments (lowest)

This hierarchy is implicit in the FOCUS spec definitions and widely understood in FinOps
practice. When discounts are applied, EffectiveCost should be less than or equal to BilledCost.

**Alternatives Considered**:

- Strict equality checks: Rejected - too inflexible for real-world data
- No relationship validation: Rejected - misses critical data quality issues

**Special Cases**:

- Zero costs are valid (free tier usage)
- Negative costs (credits/refunds) bypass hierarchy validation
- ChargeClass "Correction" records exempt from hierarchy rules

---

### 2. CommitmentDiscountStatus Conditional Requirements

**Decision**: Require CommitmentDiscountStatus when CommitmentDiscountId is present AND
ChargeCategory is USAGE

**Rationale**: Per FOCUS 1.2 Section 3.17:

- CommitmentDiscountStatus is CONDITIONAL - required when CommitmentDiscountId is not null
  AND ChargeCategory is "Usage"
- Status values: USED (utilizing commitment) or UNUSED (unused portion)
- Purchase charges don't require status because they represent the commitment purchase itself

**Alternatives Considered**:

- Always require status when ID present: Rejected - breaks Purchase category charges
- Never validate status: Rejected - misses data quality issues for utilization tracking

**Proto Enum Values** (from enums.proto):

```protobuf
enum FocusCommitmentDiscountStatus {
  FOCUS_COMMITMENT_DISCOUNT_STATUS_UNSPECIFIED = 0;
  FOCUS_COMMITMENT_DISCOUNT_STATUS_USED = 1;
  FOCUS_COMMITMENT_DISCOUNT_STATUS_UNUSED = 2;
}
```

---

### 3. CapacityReservationStatus Conditional Requirements

**Decision**: Require CapacityReservationStatus when CapacityReservationId is present AND
ChargeCategory is USAGE (same pattern as commitment discounts)

**Rationale**: Per FOCUS 1.2 Section 3.7:

- CapacityReservationStatus follows identical conditional logic to CommitmentDiscountStatus
- Status values: USED or UNUSED
- Enables accurate capacity utilization tracking

**Proto Enum Values** (from enums.proto):

```protobuf
enum FocusCapacityReservationStatus {
  FOCUS_CAPACITY_RESERVATION_STATUS_UNSPECIFIED = 0;
  FOCUS_CAPACITY_RESERVATION_STATUS_USED = 1;
  FOCUS_CAPACITY_RESERVATION_STATUS_UNUSED = 2;
}
```

---

### 4. PricingQuantity and PricingUnit Consistency

**Decision**: Require PricingUnit when PricingQuantity > 0

**Rationale**: A quantity without a unit is meaningless for cost analysis and chargeback.
The FOCUS spec defines these fields as paired - if you have a quantity, you must know the
unit of measure.

**Alternatives Considered**:

- Infer unit from context: Rejected - too error-prone, no standard inference rules
- Allow empty unit: Rejected - breaks downstream calculations

---

### 5. Zero-Allocation Validation Pattern

**Decision**: Use package-level sentinel errors for all validation failures

**Rationale**: The existing codebase (validation.go) demonstrates the zero-allocation pattern:

```go
var (
    ErrProjectedCostRequestNil = errors.New("request is required")
    // ... other sentinel errors
)
```

This pattern:

- Returns pre-allocated error values on failure (no allocation)
- Returns nil on success (no allocation)
- Achieves <100ns performance target

**Implementation Pattern**:

```go
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var (
    ErrCostRelationshipViolation = errors.New("effective_cost must not exceed billed_cost")
    // ... more sentinel errors
)
```

---

### 6. ValidationError Structured Type Design

**Decision**: Create a ValidationError struct implementing the error interface with four
discrete fields

**Rationale**: FR-007 requires structured error messages with FieldName, Constraint,
ActualValue, ExpectedValue. Go's error interface allows custom implementations.

**Design**:

```go
type ValidationError struct {
    FieldName     string
    Constraint    string
    ActualValue   string
    ExpectedValue string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s (actual: %s, expected: %s)",
        e.FieldName, e.Constraint, e.ActualValue, e.ExpectedValue)
}
```

**Zero-Allocation Consideration**: Use pre-allocated ValidationError instances for common
cases, or accept allocation on error path (errors are exceptional).

---

### 7. Validation Mode Implementation

**Decision**: Add ValidateFocusRecordWithOptions function with ValidationOptions parameter

**Rationale**: FR-009 requires both fail-fast (default) and aggregate modes. Adding an
options-based API preserves backward compatibility (FR-008) while enabling new functionality.

**Design**:

```go
type ValidationMode int

const (
    ValidationModeFailFast ValidationMode = iota  // Default
    ValidationModeAggregate
)

type ValidationOptions struct {
    Mode ValidationMode
}

func ValidateFocusRecordWithOptions(r *FocusCostRecord, opts ValidationOptions) []error
```

**Backward Compatibility**: Existing `ValidateFocusRecord(r)` calls internal
`ValidateFocusRecordWithOptions(r, ValidationOptions{Mode: ValidationModeFailFast})`

---

### 8. Floating-Point Comparison Best Practices

**Decision**: Reuse existing `floatEquals` function with `contractedCostTolerance = 0.0001`

**Rationale**: The codebase already has a well-tested floating-point comparison helper in
focus_conformance.go that uses relative tolerance. Cost comparisons should use the same
tolerance to ensure consistency.

**Existing Implementation**:

```go
const contractedCostTolerance = 0.0001

func floatEquals(a, b, tolerance float64) bool {
    if a == b {
        return true
    }
    diff := math.Abs(a - b)
    largest := math.Max(math.Abs(a), math.Abs(b))
    return diff <= largest*tolerance
}
```

---

## Summary

All technical decisions are resolved. Key implementation points:

1. Extend `validateBusinessRules()` with new contextual validation functions
2. Use sentinel errors for zero-allocation happy path
3. Add `ValidationError` type for structured error inspection
4. Add `ValidateFocusRecordWithOptions()` for aggregate mode
5. Reuse existing `floatEquals()` for cost comparisons
6. Exempt ChargeClass CORRECTION from cost hierarchy rules
