# Pricing Package

Package `pricing` provides domain types, validation, and growth calculation helpers for the
PulumiCost SDK.

## Features

- **Billing Mode Validation**: Validate 44+ billing modes across all cloud providers
- **Growth Projections**: Calculate linear and exponential cost growth projections
- **JSON Schema Validation**: Validate PricingSpec documents against embedded schema
- **Warning Detection**: Detect unrealistic growth assumptions

## Installation

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pricing"
```

## Growth Projection Helpers

The package provides comprehensive helpers for projecting future costs based on growth
assumptions.

### Growth Types

| Type | Description | Formula |
|------|-------------|---------|
| `GROWTH_TYPE_NONE` | No growth applied | `cost = baseCost` |
| `GROWTH_TYPE_LINEAR` | Linear (additive) growth | `cost = baseCost * (1 + rate * n)` |
| `GROWTH_TYPE_EXPONENTIAL` | Compound growth | `cost = baseCost * (1 + rate)^n` |
| `GROWTH_TYPE_UNSPECIFIED` | Default (treated as NONE) | `cost = baseCost` |

### Basic Usage

```go
import (
    "github.com/rshade/pulumicost-spec/sdk/go/pricing"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Apply linear growth: 10% per period for 12 periods
// Formula: $100 * (1 + 0.10 * 12) = $220
cost := pricing.ApplyLinearGrowth(100.0, 0.10, 12)

// Apply exponential growth: 5% per period for 12 periods
// Formula: $100 * (1.05)^12 = $179.59
cost := pricing.ApplyExponentialGrowth(100.0, 0.05, 12)

// Use the dispatcher function with growth type enum
rate := 0.10
cost := pricing.ApplyGrowth(100.0, pbc.GrowthType_GROWTH_TYPE_LINEAR, &rate, 12)
```

### Validation

Always validate growth parameters before applying them:

```go
rate := 0.10
err := pricing.ValidateGrowthParams(pbc.GrowthType_GROWTH_TYPE_LINEAR, &rate)
if err != nil {
    // Handle validation error
}
```

**Validation Rules:**

- `LINEAR` or `EXPONENTIAL` require a non-nil growth rate
- Growth rate must be >= -1.0 (rates below -100% would result in negative costs)
- `NONE` and `UNSPECIFIED` are always valid (rate is ignored)

### Parameter Resolution

When merging request-level overrides with resource-level defaults:

```go
// Request-level parameters override resource-level defaults
effectiveType, effectiveRate := pricing.ResolveGrowthParams(
    requestGrowthType, requestRate,   // Request-level (overrides)
    resourceGrowthType, resourceRate, // Resource-level (defaults)
)
```

### Warning Detection

Detect potentially unrealistic growth assumptions:

```go
warnings := pricing.CheckGrowthWarnings(
    pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
    &rate,
    periods,
)

for _, w := range warnings {
    log.Warn().
        Str("code", w.Code).
        Str("message", w.Message).
        Float64("rate", w.Rate).
        Int("periods", w.Periods).
        Msg("growth warning")
}
```

**Warning Codes:**

| Code | Condition | Description |
|------|-----------|-------------|
| `OVERFLOW_RISK` | calculation would overflow | Projection would result in +Inf (requires `CheckGrowthWarningsWithCost`) |
| `HIGH_GROWTH_RATE` | rate > 1.0 (>100%) | Hyper-growth may be unrealistic |
| `LONG_PROJECTION` | exponential + periods > 36 | Long-term projections become unreliable |

**Convenience Functions:**

```go
if pricing.IsHighGrowthRate(rate) {
    // Log warning about high growth rate
}

if pricing.IsLongProjection(pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL, periods) {
    // Log warning about long projection
}
```

## Constants

The package exports constants for threshold values:

```go
const (
    HighGrowthRateThreshold = 1.0   // 100% per period
    LongProjectionThreshold = 36    // months
    MinValidGrowthRate      = -1.0  // -100% (minimum allowed)
)
```

## Error Handling

```go
var (
    ErrMissingGrowthRate = errors.New("growth_rate required for LINEAR/EXPONENTIAL growth type")
    ErrInvalidGrowthRate = errors.New("growth_rate must be >= -1.0")
    ErrOverflow          = errors.New("growth projection resulted in overflow")
)
```

Check for specific errors:

```go
// Validation errors
err := pricing.ValidateGrowthParams(growthType, rate)
if errors.Is(err, pricing.ErrMissingGrowthRate) {
    // Handle missing rate for LINEAR/EXPONENTIAL
}
if errors.Is(err, pricing.ErrInvalidGrowthRate) {
    // Handle rate below -1.0
}

// Overflow errors (from ProjectCostSafely)
cost, warnings, err := pricing.ProjectCostSafely(baseCost, growthType, rate, periods)
if errors.Is(err, pricing.ErrOverflow) {
    // Handle overflow - projection would result in +Inf
}
```

## Safe Projections

For production use, prefer `ProjectCostSafely` which combines validation, overflow detection, and warning checks:

```go
cost, warnings, err := pricing.ProjectCostSafely(
    baseCost,
    pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
    &rate,
    periods,
)
if err != nil {
    if errors.Is(err, pricing.ErrOverflow) {
        // Handle overflow - projection would result in +Inf
    }
    // Handle other validation errors
}
for _, w := range warnings {
    log.Warn().Str("code", w.Code).Msg(w.Message)
}
```

### Overflow Detection

Check for potential overflow before calculation:

```go
if pricing.CheckOverflowRisk(baseCost, growthType, &rate, periods) {
    // Don't perform calculation - it would overflow
}

// Or use CheckGrowthWarningsWithCost for OVERFLOW_RISK warning
warnings := pricing.CheckGrowthWarningsWithCost(baseCost, growthType, &rate, periods)
for _, w := range warnings {
    if w.Code == "OVERFLOW_RISK" {
        // Handle overflow risk
    }
}
```

## Troubleshooting

### Why am I getting +Inf results?

Overflow to `+Inf` occurs when the projected cost exceeds `math.MaxFloat64` (~1.8×10^308).

**Common causes:**

1. **Extreme growth rates**: Rates like 100.0 (10000%) quickly overflow
2. **Very long projections**: Even 10% growth over 100,000 periods overflows
3. **Large base costs**: Starting near MaxFloat64 with any positive growth

**Solutions:**

- Use `CheckOverflowRisk()` to detect overflow before calculation
- Use `ProjectCostSafely()` which returns `ErrOverflow` instead of +Inf
- Cap projection periods to realistic values (e.g., 120 months = 10 years)

### Can I use negative growth rates?

Yes, negative rates model cost decline. Valid range is >= -1.0:

```go
// 10% decline per period
rate := -0.10
cost := pricing.ApplyExponentialGrowth(100.0, rate, 12)
// Result: 100 * (0.90)^12 = $28.24

// Complete decline to zero
rate := -1.0
cost := pricing.ApplyExponentialGrowth(100.0, rate, 1)
// Result: 100 * (0)^1 = $0
```

**Rates below -1.0 are invalid** (would produce negative costs) and will fail validation.

### What happens with negative periods?

Negative periods are mathematically valid and produce "backprojection" results:

- **Linear**: `100 * (1 + 0.10 * -5) = 50` (projects backward)
- **Exponential**: `100 * (1.10)^-5 = 62.09` (fractional multiplier)

This can be useful for calculating what costs were in past periods. If negative periods don't
make sense for your use case, validate `periods >= 0` in your application.

### How do I detect unrealistic projections?

Use `CheckGrowthWarnings` or `CheckGrowthWarningsWithCost`:

```go
warnings := pricing.CheckGrowthWarningsWithCost(baseCost, growthType, &rate, periods)
for _, w := range warnings {
    switch w.Code {
    case "OVERFLOW_RISK":
        // Projection would overflow to +Inf
    case "HIGH_GROWTH_RATE":
        // Rate exceeds 100% per period
    case "LONG_PROJECTION":
        // Exponential over 36 months
    }
}
```

## Billing Mode Validation

```go
if pricing.ValidBillingMode("per_hour") {
    // Valid billing mode
}

// Get all supported billing modes
modes := pricing.GetAllBillingModes()
```

**Supported Categories:**

- **Time-based**: per_hour, per_minute, per_second, per_day, per_month, per_year
- **Storage-based**: per_gb_month, per_gb_hour, per_gb_day
- **Usage-based**: per_request, per_operation, per_transaction, per_invocation
- **Compute-based**: per_cpu_hour, per_vcpu_hour, per_memory_gb_hour
- **Database-specific**: per_rcu, per_wcu, per_dtu, per_ru
- **Pricing models**: on_demand, reserved, spot, savings_plan

## PricingSpec Validation

Validate JSON documents against the embedded pricing spec schema:

```go
jsonData := []byte(`{
    "provider": "aws",
    "resource_type": "ec2",
    "billing_mode": "per_hour",
    "rate_per_unit": 0.10,
    "currency": "USD"
}`)

err := pricing.ValidatePricingSpec(jsonData)
if err != nil {
    // Handle validation error
}
```

## Performance

| Operation | Time | Allocations |
|-----------|------|-------------|
| ApplyLinearGrowth | <5 ns/op | 0 B/op |
| ApplyExponentialGrowth | <30 ns/op | 0 B/op |
| ApplyGrowth | <35 ns/op | 0 B/op |
| ValidateGrowthParams | <10 ns/op | 0 B/op |
| CheckOverflowRisk | <50 ns/op | 0 B/op |
| CheckGrowthWarnings | <200 ns/op | ~100 B/op (when warnings returned) |
| ProjectCostSafely | <300 ns/op | ~100 B/op (when warnings returned) |

## References

- [GrowthType Proto Definition](../../../proto/pulumicost/v1/enums.proto)
- [ResourceDescriptor Proto](../../../proto/pulumicost/v1/costsource.proto)
- [Forecasting Primitives Spec](../../../specs/030-forecasting-primitives/)
