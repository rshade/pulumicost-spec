# Plugin Migration Guide: FOCUS 1.2 Integration

This guide explains how to update existing plugins to support the FinOps FOCUS 1.2 specification using the new `FocusRecordBuilder`.

## Key Changes

- **New Data Model**: Cost records now strictly follow the FOCUS 1.2 schema.
- **Builder Pattern**: Direct struct initialization of cost records is discouraged. Use `NewFocusRecordBuilder`.
- **Strict Enums**: Service, Charge, and Pricing categories are now Enums.

## Breaking Changes (v0.2.0)

The following proto field renames align with FOCUS 1.2 naming conventions:

| Old Field Name   | New Field Name      | Proto Field # | Notes                              |
| ---------------- | ------------------- | ------------- | ---------------------------------- |
| `currency`       | `billing_currency`  | 18            | Aligns with FOCUS BillingCurrency  |
| `usage_quantity` | `consumed_quantity` | 20            | Aligns with FOCUS ConsumedQuantity |
| `usage_unit`     | `consumed_unit`     | 21            | Aligns with FOCUS ConsumedUnit     |

### Wire Format Impact

- **Binary protobuf**: Field numbers unchanged; binary compatibility preserved.
- **JSON serialization**: Field names changed (e.g., `currency` → `billingCurrency`).
- **Go SDK**: Use `FocusRecordBuilder` methods which handle naming automatically.

### Migration Example

**Before (v0.1.x):**

```go
record.Currency = "USD"
record.UsageQuantity = 100.0
record.UsageUnit = "Hours"
```

**After (v0.2.0):**

```go
// Option 1: Use the builder (recommended)
builder.WithFinancials(billed, list, effective, "USD", invoiceID)
builder.WithUsage(100.0, "Hours")

// Option 2: Direct field access (if needed)
record.BillingCurrency = "USD"
record.ConsumedQuantity = 100.0
record.ConsumedUnit = "Hours"
```

## Migration Steps

### 1. Update Imports

Add the `pluginsdk` and `proto` imports:

```go
import (
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)
```

### 2. Use FocusRecordBuilder

Replace direct struct initialization with the Builder:

**Before:**

```go
// Old style (hypothetical)
result := &pbc.ActualCostResult{
    Cost: 10.0,
    ...
}
```

**After:**

```go
builder := pluginsdk.NewFocusRecordBuilder()
builder.WithIdentity("aws", "account-id", "Account Name")
builder.WithChargePeriod(start, end)
builder.WithServiceCategory(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE)
builder.WithChargeDetails(pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE, pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD)
builder.WithFinancials(10.5, 10.5, 10.5, "USD", "inv-123")

record, err := builder.Build()
if err != nil {
    // Handle validation error
}
```

### 3. Handle Extensions (The Backpack)

If you have custom fields that don't map to FOCUS 1.2 columns, use `WithExtension`:

```go
builder.WithExtension("MyCustomField", "Value")
```
