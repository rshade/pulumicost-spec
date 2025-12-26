# Data Model: Contextual FinOps Validation

**Feature**: 027-finops-validation
**Date**: 2025-12-25

## Overview

This feature extends the validation layer for FocusCostRecord without modifying the proto
schema. All entities below are Go SDK types that wrap or validate existing protobuf messages.

## Entities

### 1. FocusCostRecord (Existing - Proto Message)

The primary record type being validated. No changes to the proto definition.

**Relevant Fields for Contextual Validation**:

| Field                       | Type                           | Validation Context            |
| --------------------------- | ------------------------------ | ----------------------------- |
| billed_cost                 | double                         | Cost hierarchy validation     |
| effective_cost              | double                         | Cost hierarchy validation     |
| list_cost                   | double                         | Cost hierarchy validation     |
| charge_class                | FocusChargeClass               | Exemption from certain rules  |
| charge_category             | FocusChargeCategory            | Conditional field requirements|
| commitment_discount_id      | string                         | Commitment consistency        |
| commitment_discount_status  | FocusCommitmentDiscountStatus  | Commitment consistency        |
| capacity_reservation_id     | string                         | Capacity consistency          |
| capacity_reservation_status | FocusCapacityReservationStatus | Capacity consistency          |
| pricing_quantity            | double                         | Pricing consistency           |
| pricing_unit                | string                         | Pricing consistency           |

---

### 2. ValidationError (New - Go Type)

A structured error type for programmatic error inspection.

**Fields**:

| Field         | Type   | Description                              |
| ------------- | ------ | ---------------------------------------- |
| FieldName     | string | Name of the field that failed validation |
| Constraint    | string | Description of the validation rule       |
| ActualValue   | string | String representation of actual value    |
| ExpectedValue | string | String representation of expected value  |

**Behavior**:

- Implements `error` interface via `Error() string` method
- Used when callers need structured error information
- Error message format: `{FieldName}: {Constraint} (actual: {ActualValue}, expected: {ExpectedValue})`

**Example**:

```go
err := &ValidationError{
    FieldName:     "effective_cost",
    Constraint:    "must not exceed billed_cost",
    ActualValue:   "150.00",
    ExpectedValue: "<= 100.00",
}
// err.Error() returns: "effective_cost: must not exceed billed_cost (actual: 150.00, expected: <= 100.00)"
```

---

### 3. ValidationMode (New - Go Type)

An enumeration controlling validation behavior.

**Values**:

| Constant                 | Value | Description                              |
| ------------------------ | ----- | ---------------------------------------- |
| ValidationModeFailFast   | 0     | Stop on first error (default)            |
| ValidationModeAggregate  | 1     | Collect all errors before returning      |

---

### 4. ValidationOptions (New - Go Type)

Configuration options for validation functions.

**Fields**:

| Field | Type           | Default              | Description                   |
| ----- | -------------- | -------------------- | ----------------------------- |
| Mode  | ValidationMode | ValidationModeFailFast | Controls fail-fast vs aggregate |

---

## Validation Rules

### Cost Hierarchy Rules (FR-001, FR-002)

- Applies when: BilledCost > 0 AND EffectiveCost > 0 AND ChargeClass != CORRECTION
- Constraint: EffectiveCost <= BilledCost
- Applies when: ListCost > 0 AND EffectiveCost > 0 AND ChargeClass != CORRECTION
- Constraint: ListCost >= EffectiveCost

**Exemptions**:

- Zero costs (free tier) - valid, no violation
- Negative costs (credits/refunds) - exempt from hierarchy
- ChargeClass CORRECTION - exempt per FOCUS spec

---

### Commitment Discount Rules (FR-003, FR-004)

| Condition                                  | Required Field            | Error If                    |
| ------------------------------------------ | ------------------------- | --------------------------- |
| CommitmentDiscountId != "" AND ChargeCategory == USAGE | CommitmentDiscountStatus | Status == UNSPECIFIED       |
| CommitmentDiscountStatus != UNSPECIFIED    | CommitmentDiscountId     | Id == ""                    |

---

### Capacity Reservation Rules (FR-005)

| Condition                                      | Required Field               | Error If                  |
| ---------------------------------------------- | ---------------------------- | ------------------------- |
| CapacityReservationId != "" AND ChargeCategory == USAGE | CapacityReservationStatus | Status == UNSPECIFIED     |

---

### Pricing Consistency Rules (FR-006)

| Condition            | Required Field | Error If       |
| -------------------- | -------------- | -------------- |
| PricingQuantity > 0  | PricingUnit    | Unit == ""     |

---

## State Transitions

Not applicable - validation is stateless. Each `FocusCostRecord` is validated independently.

---

## Relationships

```text
FocusCostRecord (proto message)
    │
    └──▶ ValidateFocusRecord() ─────────────────────┐
    │                                                │
    └──▶ ValidateFocusRecordWithOptions() ──────────┼──▶ []error / []ValidationError
                     │                               │
                     ▼                               │
              ValidationOptions ─────────────────────┘
                     │
                     ▼
              ValidationMode (FailFast | Aggregate)
```

---

## Zero-Allocation Considerations

For fail-fast mode (default), use package-level sentinel errors:

```go
var (
    ErrEffectiveCostExceedsBilledCost    = errors.New("effective_cost must not exceed billed_cost")
    ErrListCostLessThanEffectiveCost     = errors.New("list_cost must be >= effective_cost")
    ErrCommitmentStatusMissing           = errors.New("commitment_discount_status required when commitment_discount_id set for usage charges")
    ErrCommitmentIdMissingForStatus      = errors.New("commitment_discount_id required when commitment_discount_status is set")
    ErrCapacityReservationStatusMissing  = errors.New("capacity_reservation_status required when capacity_reservation_id set for usage charges")
    ErrPricingUnitMissing                = errors.New("pricing_unit required when pricing_quantity > 0")
)
```

For aggregate mode, ValidationError structs are allocated per error (acceptable since
aggregate mode is used for batch data quality workflows where allocation is expected).
