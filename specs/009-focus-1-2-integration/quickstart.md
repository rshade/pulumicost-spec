# Quickstart: FinOps FOCUS 1.2 SDK

## Overview

The `FocusBuilder` provides a type-safe, validation-enforced way to create
`FocusCostRecord` objects aligned with FinOps FOCUS 1.2.

## Basic Usage

```go
package main

import (
 "fmt"
 "time"

 pb "github.com/pulumi/finfocus-spec/sdk/go/finfocus/v1"
 "github.com/pulumi/finfocus-spec/sdk/go/pluginsdk"
)

func main() {
 // Create a new record using the Builder
 record, err := pluginsdk.NewFocusRecordBuilder().
  WithIdentity("AWS", "123456789012").
  WithChargePeriod(time.Now().Add(-24*time.Hour), time.Now()).
  WithService("AmazonEC2", pb.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE).
  WithChargeDetails(
   pb.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
   pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
  ).
  WithFinancials(1.25, "USD"). // BilledCost
  WithExtension("aws_product_code", "P-12345"). // Backpack
  Build()

 if err != nil {
  // Handle validation error (e.g., missing mandatory fields)
  panic(err)
 }

 fmt.Printf("Created Record: %+v\n", record)
}
```

## Validation Rules

The `.Build()` method will return an error if any of the following are missing:

- BillingAccountId
- ChargePeriodStart / End
- ServiceCategory
- ChargeCategory
- BilledCost (and Currency)

## Extension ("Backpack")

Use `.WithExtension(key, value)` to add any data not covered by the strict schema.
These values are stored in the `extended_columns` map.

```go
// Example extension usage
builder.WithExtension("custom_field", "custom_value")
```
