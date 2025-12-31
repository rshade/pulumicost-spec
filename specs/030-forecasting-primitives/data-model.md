# Data Model: Forecasting Primitives

**Feature**: 030-forecasting-primitives
**Date**: 2025-12-30

## Entities

### GrowthType (Enum)

Represents the mathematical model used for projecting cost growth over time.

| Value | Number | Description |
|-------|--------|-------------|
| GROWTH_TYPE_UNSPECIFIED | 0 | Default/unset - treated as NONE |
| GROWTH_TYPE_NONE | 1 | No growth applied to projections |
| GROWTH_TYPE_LINEAR | 2 | Additive growth per period |
| GROWTH_TYPE_EXPONENTIAL | 3 | Compounding growth per period |

**Validation Rules**:

- UNSPECIFIED and NONE are functionally equivalent (no growth)
- LINEAR and EXPONENTIAL require a valid `growth_rate` value

**State Transitions**: N/A (enum is stateless)

### growth_rate (Field)

The rate of growth per projection period, expressed as a decimal fraction.

| Attribute | Value |
|-----------|-------|
| Type | optional double |
| Range | >= -1.0 (no upper bound) |
| Unit | Decimal fraction (0.10 = 10%) |
| Default | Unset (nil in Go) |

**Validation Rules**:

- REQUIRED when `growth_type` is LINEAR or EXPONENTIAL
- MUST be >= -1.0 (prevents negative costs)
- IGNORED when `growth_type` is NONE or UNSPECIFIED
- Value of 0.0 is valid (no growth, equivalent to NONE)

**Examples**:

| Value | Meaning |
|-------|---------|
| 0.10 | 10% growth per period |
| 0.05 | 5% growth per period |
| 0.0 | No growth |
| -0.10 | 10% decline per period |
| -1.0 | 100% decline (to zero) |
| 2.0 | 200% growth per period (hyper-growth) |

## Message Extensions

### ResourceDescriptor (Extended)

```text
Existing fields (1-8):
  1: provider (string)
  2: resource_type (string)
  3: sku (string)
  4: region (string)
  5: tags (map<string, string>)
  6: utilization_percentage (optional double)
  7: id (string)
  8: arn (string)

New fields:
  9: growth_type (GrowthType) - Default growth model for this resource
  10: growth_rate (optional double) - Default growth rate for this resource
```

**Semantics**: Growth parameters on ResourceDescriptor serve as defaults for all
projection requests involving this resource.

### GetProjectedCostRequest (Extended)

```text
Existing fields (1-2):
  1: resource (ResourceDescriptor)
  2: utilization_percentage (double)

New fields:
  3: growth_type (GrowthType) - Override growth model for this request
  4: growth_rate (optional double) - Override growth rate for this request
```

**Semantics**: Request-level growth parameters override ResourceDescriptor defaults
when set, enabling different projection scenarios for the same resource.

## Relationships

```text
┌─────────────────────────────────┐
│   GetProjectedCostRequest       │
│  ┌───────────────────────────┐  │
│  │ growth_type (override)    │──┼──┐
│  │ growth_rate (override)    │  │  │
│  └───────────────────────────┘  │  │
│              │                  │  │  Overrides if set
│              ▼                  │  │
│  ┌───────────────────────────┐  │  │
│  │    ResourceDescriptor     │  │  │
│  │  ┌─────────────────────┐  │  │  │
│  │  │ growth_type (default)│◄─┼──┼──┘
│  │  │ growth_rate (default)│  │  │
│  │  └─────────────────────┘  │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
```

## Growth Formulas

### Linear Growth

```text
cost_at_period_n = base_cost × (1 + rate × n)

Where:
  base_cost = current monthly cost
  rate = growth_rate (decimal)
  n = number of periods (months)

Example: base=$100, rate=0.10, n=3
  Period 0: $100 × (1 + 0.10 × 0) = $100.00
  Period 1: $100 × (1 + 0.10 × 1) = $110.00
  Period 2: $100 × (1 + 0.10 × 2) = $120.00
  Period 3: $100 × (1 + 0.10 × 3) = $130.00
```

### Exponential Growth

```text
cost_at_period_n = base_cost × (1 + rate)^n

Where:
  base_cost = current monthly cost
  rate = growth_rate (decimal)
  n = number of periods (months)

Example: base=$100, rate=0.10, n=3
  Period 0: $100 × (1.10)^0 = $100.00
  Period 1: $100 × (1.10)^1 = $110.00
  Period 2: $100 × (1.10)^2 = $121.00
  Period 3: $100 × (1.10)^3 = $133.10
```

## Validation Matrix

| growth_type | growth_rate | Result |
|-------------|-------------|--------|
| UNSPECIFIED | unset | Valid (no growth) |
| UNSPECIFIED | set | Valid (rate ignored, warning) |
| NONE | unset | Valid (no growth) |
| NONE | set | Valid (rate ignored, warning) |
| LINEAR | unset | Invalid (InvalidArgument) |
| LINEAR | >= -1.0 | Valid |
| LINEAR | < -1.0 | Invalid (InvalidArgument) |
| EXPONENTIAL | unset | Invalid (InvalidArgument) |
| EXPONENTIAL | >= -1.0 | Valid |
| EXPONENTIAL | < -1.0 | Invalid (InvalidArgument) |

## Backward Compatibility

| Scenario | Behavior |
|----------|----------|
| Old client, new server | Works - new fields default to unset/UNSPECIFIED |
| New client, old server | Works - server ignores unknown fields |
| Old client, old server | Unchanged |
| New client, new server | Full forecasting support |
