# Quickstart: Forecasting Primitives

This guide shows how to use the new forecasting primitives in cost projections.

## Overview

Forecasting primitives allow you to model expected usage growth when projecting future
costs. Two growth models are supported:

- **Linear**: Costs increase by a fixed percentage of the base cost each period
- **Exponential**: Costs compound at a fixed rate each period

## Basic Usage

### 1. Linear Growth (10% per month)

```go
import (
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// Create a resource with 10% monthly linear growth
resource := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
    Sku:          "t3.medium",
    Region:       "us-east-1",
    GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
    GrowthRate:   proto.Float64(0.10), // 10% growth
}

// Request projected costs
req := &pbc.GetProjectedCostRequest{
    Resource: resource,
}
```

**Result**: If base cost is $100/month:

- Month 1: $110 (+10%)
- Month 2: $120 (+20%)
- Month 3: $130 (+30%)

### 2. Exponential Growth (5% per month)

```go
resource := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "s3",
    Sku:          "standard",
    Region:       "us-east-1",
    GrowthType:   pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
    GrowthRate:   proto.Float64(0.05), // 5% compounding
}
```

**Result**: If base cost is $100/month:

- Month 1: $105.00
- Month 2: $110.25
- Month 3: $115.76

### 3. Declining Usage (-10% per month)

```go
resource := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
    Sku:          "t3.large",
    Region:       "us-east-1",
    GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
    GrowthRate:   proto.Float64(-0.10), // 10% decline
}
```

**Result**: If base cost is $100/month:

- Month 1: $90
- Month 2: $80
- Month 3: $70

## Override at Request Level

You can override resource defaults for specific projection scenarios:

```go
// Resource has 10% linear growth as default
resource := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
    Sku:          "t3.medium",
    Region:       "us-east-1",
    GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
    GrowthRate:   proto.Float64(0.10),
}

// But request uses 5% exponential (override)
req := &pbc.GetProjectedCostRequest{
    Resource:   resource,
    GrowthType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
    GrowthRate: proto.Float64(0.05),
}
```

## Validation

The SDK validates growth parameters automatically:

```go
import "github.com/rshade/finfocus-spec/sdk/go/pricing"

// Validate before sending request
err := pricing.ValidateGrowthParams(
    pbc.GrowthType_GROWTH_TYPE_LINEAR,
    proto.Float64(0.10),
)
if err != nil {
    // Handle validation error
}
```

### Validation Rules

| Scenario | Valid? |
|----------|--------|
| LINEAR with rate | Yes |
| LINEAR without rate | No (error) |
| EXPONENTIAL with rate | Yes |
| EXPONENTIAL without rate | No (error) |
| NONE with rate | Yes (rate ignored) |
| NONE without rate | Yes |
| Rate < -1.0 | No (error) |
| Rate > 1.0 (e.g., 200%) | Yes (hyper-growth) |

## Common Patterns

### Budget Planning (12-month projection)

```go
// Conservative estimate with 5% annual growth (compound monthly)
resource.GrowthType = pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL
resource.GrowthRate = proto.Float64(0.004) // ~5% annual
```

### Scaling Scenarios

```go
// Aggressive growth scenario
aggressive := &pbc.GetProjectedCostRequest{
    Resource:   resource,
    GrowthType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
    GrowthRate: proto.Float64(0.20), // 20% monthly
}

// Conservative scenario
conservative := &pbc.GetProjectedCostRequest{
    Resource:   resource,
    GrowthType: pbc.GrowthType_GROWTH_TYPE_LINEAR,
    GrowthRate: proto.Float64(0.05), // 5% monthly
}
```

### Sunsetting Resources

```go
// Resource being phased out over 10 months
resource.GrowthType = pbc.GrowthType_GROWTH_TYPE_LINEAR
resource.GrowthRate = proto.Float64(-0.10) // 10% decline
```

## Error Handling

```go
import "github.com/rs/zerolog/log"

resp, err := client.GetProjectedCost(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok && st.Code() == codes.InvalidArgument {
        // Handle validation error (e.g., missing growth_rate)
        log.Error().Str("message", st.Message()).Msg("Invalid request")
    }
}
```

## Backward Compatibility

Existing code continues to work unchanged:

```go
// This still works - no growth applied
resource := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
    Sku:          "t3.medium",
    Region:       "us-east-1",
    // No growth fields set - defaults to GROWTH_TYPE_UNSPECIFIED (no growth)
}
```
