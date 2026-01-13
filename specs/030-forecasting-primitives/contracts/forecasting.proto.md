# Proto Contract: Forecasting Primitives

This document defines the protocol buffer changes for the forecasting primitives feature.

## New Enum: GrowthType

**File**: `proto/finfocus/v1/enums.proto`

```protobuf
// GrowthType represents the mathematical model used for projecting cost growth.
// Used in ResourceDescriptor and GetProjectedCostRequest for forward-looking projections.
enum GrowthType {
  // Default value - treated identically to GROWTH_TYPE_NONE.
  // When unset, no growth is applied to projections.
  GROWTH_TYPE_UNSPECIFIED = 0;

  // No growth applied to projections.
  // Cost projections remain constant at the base cost.
  GROWTH_TYPE_NONE = 1;

  // Linear (additive) growth per projection period.
  // Formula: cost_at_n = base_cost * (1 + rate * n)
  // Requires growth_rate to be set.
  GROWTH_TYPE_LINEAR = 2;

  // Exponential (compounding) growth per projection period.
  // Formula: cost_at_n = base_cost * (1 + rate)^n
  // Requires growth_rate to be set.
  GROWTH_TYPE_EXPONENTIAL = 3;
}
```

## Extended Message: ResourceDescriptor

**File**: `proto/finfocus/v1/costsource.proto`

Add the following fields after field 8 (arn):

```protobuf
message ResourceDescriptor {
  // ... existing fields 1-8 ...

  // growth_type specifies the default growth model for cost projections.
  // OPTIONAL. When set, defines how projected costs should grow over time.
  // Can be overridden by GetProjectedCostRequest.growth_type.
  //
  // Values:
  //   - GROWTH_TYPE_UNSPECIFIED/NONE: No growth (constant projections)
  //   - GROWTH_TYPE_LINEAR: Additive growth (cost * (1 + rate * periods))
  //   - GROWTH_TYPE_EXPONENTIAL: Compounding growth (cost * (1 + rate)^periods)
  //
  // When LINEAR or EXPONENTIAL, growth_rate MUST also be provided.
  GrowthType growth_type = 9;

  // growth_rate specifies the default growth rate per projection period.
  // OPTIONAL. Required when growth_type is LINEAR or EXPONENTIAL.
  //
  // Valid range: >= -1.0 (no upper bound)
  //   - Positive values: growth (e.g., 0.10 = 10% growth per period)
  //   - Zero: no growth (equivalent to GROWTH_TYPE_NONE)
  //   - Negative values: decline (e.g., -0.10 = 10% decline per period)
  //   - -1.0: complete decline to zero cost
  //
  // Values below -1.0 are invalid (would produce negative costs).
  // Can be overridden by GetProjectedCostRequest.growth_rate.
  optional double growth_rate = 10;
}
```

## Extended Message: GetProjectedCostRequest

**File**: `proto/finfocus/v1/costsource.proto`

Add the following fields after field 2 (utilization_percentage):

```protobuf
message GetProjectedCostRequest {
  // ... existing fields 1-2 ...

  // growth_type overrides ResourceDescriptor.growth_type for this request.
  // OPTIONAL. When set, takes precedence over the resource-level default.
  //
  // Use case: Project different growth scenarios for the same resource
  // without modifying the resource descriptor.
  //
  // When LINEAR or EXPONENTIAL, growth_rate MUST also be provided
  // (either here or in ResourceDescriptor).
  GrowthType growth_type = 3;

  // growth_rate overrides ResourceDescriptor.growth_rate for this request.
  // OPTIONAL. When set, takes precedence over the resource-level default.
  //
  // Valid range: >= -1.0 (no upper bound)
  //
  // Override semantics: If this field is set, it fully replaces
  // ResourceDescriptor.growth_rate for this request.
  optional double growth_rate = 4;
}
```

## Validation Specification

### SDK Validation Function

```go
// ValidateGrowthParams validates growth_type and growth_rate combination.
// Returns InvalidArgument gRPC status on validation failure.
//
// Rules:
//   - LINEAR/EXPONENTIAL require growth_rate to be set
//   - growth_rate must be >= -1.0
//   - NONE/UNSPECIFIED ignore growth_rate (optional warning)
func ValidateGrowthParams(growthType GrowthType, growthRate *float64) error
```

### Error Codes

| Condition | gRPC Status | Message |
|-----------|-------------|---------|
| LINEAR without rate | InvalidArgument | "growth_rate required for LINEAR growth type" |
| EXPONENTIAL without rate | InvalidArgument | "growth_rate required for EXPONENTIAL growth type" |
| Rate < -1.0 | InvalidArgument | "growth_rate must be >= -1.0" |

## Wire Format Examples

### Request with Linear Growth

```json
{
  "resource": {
    "provider": "aws",
    "resource_type": "ec2",
    "sku": "t3.medium",
    "region": "us-east-1",
    "growth_type": "GROWTH_TYPE_LINEAR",
    "growth_rate": 0.10
  },
  "utilization_percentage": 0.5
}
```

### Request with Override

```json
{
  "resource": {
    "provider": "aws",
    "resource_type": "ec2",
    "sku": "t3.medium",
    "region": "us-east-1",
    "growth_type": "GROWTH_TYPE_LINEAR",
    "growth_rate": 0.10
  },
  "utilization_percentage": 0.5,
  "growth_type": "GROWTH_TYPE_EXPONENTIAL",
  "growth_rate": 0.05
}
```

In this example, the request uses EXPONENTIAL at 5% (override), not LINEAR at 10%
(resource default).

## Breaking Change Analysis

| Change | Breaking? | Reason |
|--------|-----------|--------|
| Add GrowthType enum | No | New type, no existing references |
| Add growth_type to ResourceDescriptor | No | Optional field, defaults to unset |
| Add growth_rate to ResourceDescriptor | No | Optional field, defaults to unset |
| Add growth_type to GetProjectedCostRequest | No | Optional field, defaults to unset |
| Add growth_rate to GetProjectedCostRequest | No | Optional field, defaults to unset |

**buf breaking check**: Expected to PASS (no breaking changes)
