# Quickstart: FOCUS 1.2 Column Audit

**Date**: 2025-11-28
**Feature**: 010-focus-column-audit

## Overview

This quickstart guide demonstrates how to use the new FOCUS 1.2 columns added by
this feature. After implementation, the `FocusCostRecord` message and
`FocusRecordBuilder` will support all 57 FOCUS 1.2 columns.

## Prerequisites

- Go 1.24+
- pulumicost-spec SDK v0.5.0+ (after this feature is merged)

## Basic Usage

### Building a Complete FOCUS Record

```go
package main

import (
    "time"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func main() {
    now := time.Now()
    periodStart := now.AddDate(0, -1, 0) // Last month
    periodEnd := now

    record, err := pluginsdk.NewFocusRecordBuilder().
        // Identity (Mandatory)
        WithIdentity("AWS", "123456789012", "Production Account").
        WithSubAccount("111222333444", "Development").
        WithBillingAccountType("Organization").       // NEW
        WithSubAccountType("Member").                  // NEW

        // Billing Period (Mandatory)
        WithBillingPeriod(periodStart, periodEnd, "USD").

        // Charge Period (Mandatory)
        WithChargePeriod(periodStart, periodEnd).

        // Charge Details (Mandatory)
        WithChargeDetails(
            pbc.FOCUS_CHARGE_CATEGORY_USAGE,
            pbc.FOCUS_PRICING_CATEGORY_STANDARD,
        ).
        WithChargeClassification(
            pbc.FOCUS_CHARGE_CLASS_REGULAR,
            "EC2 Instance Usage",
            pbc.FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
        ).

        // Financial (Mandatory)
        WithFinancials(100.00, 120.00, 100.00, "USD", "INV-2024-001").
        WithContractedCost(95.00).                     // NEW - MANDATORY
        WithContractedUnitPrice(0.095).                // NEW

        // Service (Conditional)
        WithService(pbc.FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2").
        WithServiceSubcategory("Virtual Machines").    // NEW
        WithPublisher("AWS").                          // NEW

        // Resource (Conditional)
        WithResource("i-1234567890abcdef0", "web-server-1", "m5.large").

        // SKU (Conditional)
        WithSKU("sku-12345", "price-67890").
        WithSkuDetails("hour", "On-Demand Linux/UNIX"). // NEW

        // Location (Conditional)
        WithLocation("us-east-1", "US East (N. Virginia)", "us-east-1a").

        // Pricing (Conditional)
        WithPricing(720.0, "Hours", 0.10).
        WithPricingCurrency("EUR").                    // NEW
        WithPricingCurrencyPrices(0.085, 90.00, 0.095). // NEW

        // Usage (Conditional)
        WithUsage(720.0, "Hours").

        // Commitment Discount (Conditional)
        WithCommitmentDiscount(
            pbc.FOCUS_COMMITMENT_DISCOUNT_CATEGORY_USAGE,
            "ri-12345",
            "1-year Reserved Instance",
        ).
        WithCommitmentDiscountDetails(               // NEW
            720.0, // quantity
            pbc.FOCUS_COMMITMENT_DISCOUNT_STATUS_USED,
            "Reserved Instance",
            "Hours",
        ).

        // Capacity Reservation (Conditional) - NEW
        WithCapacityReservation(
            "cr-12345",
            pbc.FOCUS_CAPACITY_RESERVATION_STATUS_USED,
        ).

        // Invoice (Conditional)
        WithInvoice("INV-2024-001", "Amazon Web Services").

        // Tags (Conditional)
        WithTags(map[string]string{
            "Environment": "Production",
            "Team":        "Platform",
        }).

        Build()

    if err != nil {
        panic(err)
    }

    // Use the record...
    _ = record
}
```

## New Builder Methods

### Financial Methods

```go
// ContractedCost - MANDATORY
builder.WithContractedCost(95.00)

// ContractedUnitPrice - Conditional
builder.WithContractedUnitPrice(0.095)
```

### Account Type Methods

```go
// BillingAccountType - Conditional
builder.WithBillingAccountType("Organization")

// SubAccountType - Conditional
builder.WithSubAccountType("Member")
```

### Capacity Reservation Methods

```go
// CapacityReservation - Conditional (sets both ID and Status)
builder.WithCapacityReservation(
    "cr-12345",
    pbc.FOCUS_CAPACITY_RESERVATION_STATUS_USED,
)
```

### Commitment Discount Extended Methods

```go
// CommitmentDiscountDetails - Conditional (sets quantity, status, type, unit)
builder.WithCommitmentDiscountDetails(
    720.0,                                          // quantity
    pbc.FOCUS_COMMITMENT_DISCOUNT_STATUS_USED,      // status
    "Reserved Instance",                            // type
    "Hours",                                        // unit
)
```

### Pricing Currency Methods

```go
// PricingCurrency - Conditional (ISO 4217 code)
builder.WithPricingCurrency("EUR")

// PricingCurrencyPrices - Conditional (sets all three pricing currency fields)
builder.WithPricingCurrencyPrices(
    0.085,  // contracted unit price in pricing currency
    90.00,  // effective cost in pricing currency
    0.095,  // list unit price in pricing currency
)
```

### Service/SKU Extended Methods

```go
// Publisher - Conditional
builder.WithPublisher("AWS")

// ServiceSubcategory - Conditional
builder.WithServiceSubcategory("Virtual Machines")

// SkuDetails - Conditional (sets both meter and price details)
builder.WithSkuDetails("hour", "On-Demand Linux/UNIX")
```

## New Enum Types

### FocusCommitmentDiscountStatus

```go
pbc.FOCUS_COMMITMENT_DISCOUNT_STATUS_UNSPECIFIED // 0
pbc.FOCUS_COMMITMENT_DISCOUNT_STATUS_USED        // 1
pbc.FOCUS_COMMITMENT_DISCOUNT_STATUS_UNUSED      // 2
```

### FocusCapacityReservationStatus

```go
pbc.FOCUS_CAPACITY_RESERVATION_STATUS_UNSPECIFIED // 0
pbc.FOCUS_CAPACITY_RESERVATION_STATUS_USED        // 1
pbc.FOCUS_CAPACITY_RESERVATION_STATUS_UNUSED      // 2
```

## Validation

The builder validates all fields when `Build()` is called. Validation includes:

1. **Mandatory field presence**: ContractedCost must be set
2. **Cross-field consistency**: ContractedCost = ContractedUnitPrice × PricingQuantity
3. **Enum value validity**: Status enums must not be UNSPECIFIED when IDs are present
4. **ISO 4217 compliance**: PricingCurrency must be valid currency code

## Running the Audit Script

After implementation, verify column completeness:

```bash
go run scripts/audit_focus_columns.go
```

Expected output:

```text
FOCUS 1.2 Column Audit
======================
Total columns: 57
Implemented: 57
Missing: 0

✅ All FOCUS 1.2 columns are implemented
```

## Migration Guide

If upgrading from a previous version:

1. **Field renames**: Some fields were renamed to align with FOCUS 1.2 naming
2. **Optional fields**: All new conditional fields default to zero values
3. **Mandatory field**: `ContractedCost` is now mandatory - set it to match
   your billed cost if contracted pricing is not applicable

### Field Rename Migration

The following fields were renamed for FOCUS 1.2 compliance:

| Old Name | New Name | Migration |
|----------|----------|-----------|
| `Currency` | `BillingCurrency` | Use `WithFinancials()` or set directly |
| `UsageQuantity` | `ConsumedQuantity` | Use `WithUsage()` or set directly |
| `UsageUnit` | `ConsumedUnit` | Use `WithUsage()` or set directly |

**Before (v0.1.x):**

```go
record.Currency = "USD"
record.UsageQuantity = 720.0
record.UsageUnit = "Hours"
```

**After (v0.2.0):**

```go
// Option 1: Use builder methods (recommended)
builder.WithFinancials(billed, list, effective, "USD", invoiceID)
builder.WithUsage(720.0, "Hours")

// Option 2: Direct field access
record.BillingCurrency = "USD"
record.ConsumedQuantity = 720.0
record.ConsumedUnit = "Hours"
```

### Adding ContractedCost

```go
// Minimal migration - use billed cost as fallback:
builder.WithContractedCost(billedCost)

// If you have contracted pricing:
builder.WithContractedCost(contractedUnitPrice * pricingQuantity)
builder.WithContractedUnitPrice(contractedUnitPrice)
```

### Complete Migration Example

```go
// Before: Old plugin returning raw cost data
func (p *OldPlugin) GetCost() float64 {
    return 100.0
}

// After: FOCUS 1.2 compliant record
func (p *NewPlugin) GetCostRecord() (*pbc.FocusCostRecord, error) {
    cost := 100.0
    return pluginsdk.NewFocusRecordBuilder().
        WithIdentity("AWS", p.AccountID, p.AccountName).
        WithBillingPeriod(p.PeriodStart, p.PeriodEnd, "USD").
        WithChargePeriod(p.ChargeStart, p.ChargeEnd).
        WithChargeDetails(
            pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
            pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
        ).
        WithChargeClassification(
            pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
            "Usage charge",
            pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED,
        ).
        WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "EC2").
        WithFinancials(cost, cost, cost, "USD", "").
        WithContractedCost(cost). // Required for FOCUS 1.2
        WithUsage(1.0, "Hour").
        Build()
}
```
