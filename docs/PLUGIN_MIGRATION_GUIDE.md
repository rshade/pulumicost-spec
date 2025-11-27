# Plugin Migration Guide: FOCUS 1.2 Integration

This guide explains how to update existing plugins to support the FinOps FOCUS 1.2 specification using the new `FocusRecordBuilder`.

## Key Changes

- **New Data Model**: Cost records now strictly follow the FOCUS 1.2 schema.
- **Builder Pattern**: Direct struct initialization of cost records is discouraged. Use `NewFocusRecordBuilder`.
- **Strict Enums**: Service, Charge, and Pricing categories are now Enums.

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
