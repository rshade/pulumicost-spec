# Quickstart: Contextual FinOps Validation

**Feature**: 027-finops-validation
**Date**: 2025-12-25

## Overview

This guide shows how to use the contextual FinOps validation features in the pluginsdk.

## Basic Usage (Fail-Fast Mode)

The simplest way to validate a FocusCostRecord:

```go
import (
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func processRecord(record *pbc.FocusCostRecord) error {
    // Validates all FOCUS rules including contextual business logic
    if err := pluginsdk.ValidateFocusRecord(record); err != nil {
        return fmt.Errorf("invalid record: %w", err)
    }
    // Process valid record...
    return nil
}
```

## Aggregate Mode (Collect All Errors)

For batch processing and data quality reports:

```go
import (
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func validateBatch(records []*pbc.FocusCostRecord) map[int][]error {
    opts := pluginsdk.ValidationOptions{
        Mode: pluginsdk.ValidationModeAggregate,
    }

    issues := make(map[int][]error)
    for i, record := range records {
        errs := pluginsdk.ValidateFocusRecordWithOptions(record, opts)
        if len(errs) > 0 {
            issues[i] = errs
        }
    }
    return issues
}
```

## Inspecting Structured Errors

When you need detailed error information:

```go
import (
    "errors"
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

func handleValidationError(err error) {
    var valErr *pluginsdk.ValidationError
    if errors.As(err, &valErr) {
        fmt.Printf("Field: %s\n", valErr.FieldName)
        fmt.Printf("Constraint: %s\n", valErr.Constraint)
        fmt.Printf("Actual: %s\n", valErr.ActualValue)
        fmt.Printf("Expected: %s\n", valErr.ExpectedValue)
    } else {
        fmt.Printf("Error: %s\n", err.Error())
    }
}
```

## Validation Rules Summary

### Cost Hierarchy

Records must satisfy: `ListCost >= BilledCost >= EffectiveCost` (when all are positive)

**Exemptions**:

- Zero costs (free tier usage)
- Negative costs (credits, refunds)
- ChargeClass = CORRECTION

### Commitment Discounts

| Scenario                               | Requirement                     |
| -------------------------------------- | ------------------------------- |
| CommitmentDiscountId set + Usage charge | CommitmentDiscountStatus required |
| CommitmentDiscountStatus set           | CommitmentDiscountId required   |

### Capacity Reservations

| Scenario                                  | Requirement                        |
| ----------------------------------------- | ---------------------------------- |
| CapacityReservationId set + Usage charge  | CapacityReservationStatus required |

### Pricing Consistency

| Scenario             | Requirement            |
| -------------------- | ---------------------- |
| PricingQuantity > 0  | PricingUnit required   |

## Example: Complete Plugin Validation

```go
package main

import (
    "context"
    "fmt"

    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type MyPlugin struct {
    pbc.UnimplementedCostSourceServiceServer
}

func (p *MyPlugin) GetActualCost(
    ctx context.Context,
    req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
    // Validate request
    if err := pluginsdk.ValidateActualCostRequest(req); err != nil {
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }

    // Fetch records from backend...
    records := fetchRecords(req)

    // Validate each record before returning
    for i, record := range records {
        if err := pluginsdk.ValidateFocusRecord(record); err != nil {
            return nil, status.Errorf(codes.Internal,
                "invalid record %d: %v", i, err)
        }
    }

    return &pbc.GetActualCostResponse{Results: records}, nil
}
```

## Performance Characteristics

| Mode       | Valid Record | Invalid Record |
| ---------- | ------------ | -------------- |
| Fail-Fast  | <100ns, 0 allocs | ~50ns + 1 alloc (error string) |
| Aggregate  | <100ns, 0 allocs | ~100ns per error + allocations |

The fail-fast mode uses pre-allocated sentinel errors for common violations, achieving
zero allocation on both valid records and simple error cases.
