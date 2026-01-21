# Migration Guide: Pricing Tier & Spot Risk Fields

This guide helps existing plugin developers adopt the new `pricing_category` and
`spot_interruption_risk_score` fields introduced in spec 041.

## Backward Compatibility Guarantee

**Existing plugins will continue to work without modification.** The new fields have
proto3 default values that pass validation:

| Field | Default Value | Validation Status |
|-------|---------------|-------------------|
| `pricing_category` | `UNSPECIFIED` (0) | Valid |
| `spot_interruption_risk_score` | `0.0` | Valid |

The combination of `UNSPECIFIED` + `0.0` is explicitly supported for backward
compatibility and produces no validation errors or warnings.

## Migration Timeline

| Phase | Target Date | Action |
|-------|-------------|--------|
| Phase 1 | Immediate | No action required - existing plugins work as-is |
| Phase 2 | v1.0 Release | Plugins SHOULD populate `pricing_category` with meaningful values |
| Phase 3 | v2.0 Release | `UNSPECIFIED` may trigger advisory warnings (not errors) |

## Quick Start: Adding Pricing Category

### Before (Legacy Plugin)

```go
func (p *MyPlugin) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (
    *pb.GetProjectedCostResponse, error) {
    return &pb.GetProjectedCostResponse{
        UnitPrice:    0.05,
        Currency:     "USD",
        CostPerMonth: 36.50,
        // No pricing_category or spot_interruption_risk_score
    }, nil
}
```

### After (Updated Plugin)

```go
func (p *MyPlugin) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (
    *pb.GetProjectedCostResponse, error) {
    return &pb.GetProjectedCostResponse{
        UnitPrice:                 0.05,
        Currency:                  "USD",
        CostPerMonth:              36.50,
        PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
        SpotInterruptionRiskScore: 0.0, // Not applicable for standard pricing
    }, nil
}
```

## Validation Rules

### Zero Score Is Always Valid

The `spot_interruption_risk_score` of `0.0` is valid for **all** pricing categories:

```go
// All of these pass validation
resp1 := &pb.EstimateCostResponse{PricingCategory: pb.UNSPECIFIED, SpotInterruptionRiskScore: 0.0}  // OK
resp2 := &pb.EstimateCostResponse{PricingCategory: pb.STANDARD, SpotInterruptionRiskScore: 0.0}    // OK
resp3 := &pb.EstimateCostResponse{PricingCategory: pb.COMMITTED, SpotInterruptionRiskScore: 0.0}   // OK
resp4 := &pb.EstimateCostResponse{PricingCategory: pb.DYNAMIC, SpotInterruptionRiskScore: 0.0}     // OK (unusual)
```

### Non-Zero Score Requires DYNAMIC Category

Non-zero risk scores are **only valid** when `pricing_category` is `DYNAMIC`:

```go
// VALID: DYNAMIC category with non-zero risk
resp := &pb.EstimateCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
    SpotInterruptionRiskScore: 0.8,
}

// INVALID: Non-DYNAMIC category with non-zero risk - fails validation
resp := &pb.EstimateCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
    SpotInterruptionRiskScore: 0.5, // ERROR: risk score requires DYNAMIC category
}
```

### Epsilon Tolerance for Floating-Point

Validation uses epsilon tolerance (1e-9) for floating-point comparison:

- Values within `1e-9` of zero are treated as zero
- Values within `1e-9` of 1.0 are treated as 1.0
- This handles floating-point arithmetic representation errors

## Best Practices by Pricing Category

### STANDARD (On-Demand)

Use for pay-as-you-go resources with no commitment discount:

```go
resp := &pb.GetProjectedCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
    SpotInterruptionRiskScore: 0.0, // Not applicable
}
```

### COMMITTED (Reserved/Savings Plans)

Use for resources covered by reservations, savings plans, or committed use discounts:

```go
resp := &pb.GetProjectedCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_COMMITTED,
    SpotInterruptionRiskScore: 0.0, // Committed resources have no interruption risk
}
```

### DYNAMIC (Spot/Preemptible)

Use for spot instances, preemptible VMs, or other interruptible resources:

```go
// Risk data available
resp := &pb.GetProjectedCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
    SpotInterruptionRiskScore: 0.75, // 75% interruption probability
}

// Risk data unavailable (use CheckSpotRiskConsistency for warnings)
resp := &pb.GetProjectedCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
    SpotInterruptionRiskScore: 0.0, // Valid but unusual - log a warning
}
```

### UNSPECIFIED (Legacy/Unknown)

Use only when pricing information is genuinely unavailable. Avoid for new plugins:

```go
// Legacy plugin behavior - still valid
resp := &pb.GetProjectedCostResponse{
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_UNSPECIFIED,
    SpotInterruptionRiskScore: 0.0,
}
```

## Using Validation Helpers

### ValidateEstimateCostResponse / ValidateGetProjectedCostResponse

Always validate responses before returning them:

```go
import "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"

resp := &pb.EstimateCostResponse{
    Currency:                  "USD",
    CostMonthly:               50.0,
    PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
    SpotInterruptionRiskScore: 0.8,
}

if err := pluginsdk.ValidateEstimateCostResponse(resp); err != nil {
    return nil, status.Errorf(codes.Internal, "invalid response: %v", err)
}
return resp, nil
```

### CheckSpotRiskConsistency (Advisory Warnings)

Use to detect potentially missing or inconsistent data:

```go
warnings := pluginsdk.CheckSpotRiskConsistency(
    resp.GetPricingCategory(),
    resp.GetSpotInterruptionRiskScore(),
)
for _, w := range warnings {
    log.Warn().Str("warning", w).Msg("pricing data consistency check")
}
```

Common warnings:

- "spot_interruption_risk_score > 0.0 but pricing_category is not DYNAMIC"
- "pricing_category is DYNAMIC but spot_interruption_risk_score is 0.0 (risk data may be unavailable)"

## Determining Pricing Category

### AWS

| AWS Pricing Model | pricing_category |
|-------------------|------------------|
| On-Demand | STANDARD |
| Reserved Instance (RI) | COMMITTED |
| Savings Plan | COMMITTED |
| Spot Instance | DYNAMIC |
| Free Tier | STANDARD |

### Azure

| Azure Pricing Model | pricing_category |
|---------------------|------------------|
| Pay-As-You-Go | STANDARD |
| Reserved VM Instance | COMMITTED |
| Azure Hybrid Benefit | COMMITTED |
| Spot VM | DYNAMIC |

### GCP

| GCP Pricing Model | pricing_category |
|-------------------|------------------|
| On-Demand | STANDARD |
| Committed Use Discount (CUD) | COMMITTED |
| Sustained Use Discount (SUD) | STANDARD (automatic, no commitment) |
| Preemptible VM | DYNAMIC |
| Spot VM | DYNAMIC |

## Getting Spot Interruption Risk Data

### AWS Spot Instance Advisor

AWS provides interruption frequency data via Spot Instance Advisor:

```go
// Map AWS interruption frequency to risk score
func awsInterruptionToRiskScore(frequency string) float64 {
    switch frequency {
    case "<5%":
        return 0.025
    case "5-10%":
        return 0.075
    case "10-15%":
        return 0.125
    case "15-20%":
        return 0.175
    case ">20%":
        return 0.25 // Conservative estimate
    default:
        return 0.0 // Unknown
    }
}
```

### Azure / GCP

Azure and GCP do not publish detailed interruption statistics. Consider:

- Using historical data from your own workloads
- Setting a conservative default (e.g., 0.5 for unknown risk)
- Setting 0.0 and logging a warning about unavailable data

## Testing Your Migration

### Verify Backward Compatibility

```go
func TestBackwardCompatibility(t *testing.T) {
    // Legacy response (no new fields)
    resp := &pb.EstimateCostResponse{
        Currency:    "USD",
        CostMonthly: 50.0,
    }
    err := pluginsdk.ValidateEstimateCostResponse(resp)
    assert.NoError(t, err, "Legacy responses must remain valid")
}
```

### Verify New Field Validation

```go
func TestPricingCategoryValidation(t *testing.T) {
    // Valid: DYNAMIC with risk score
    resp := &pb.EstimateCostResponse{
        Currency:                  "USD",
        CostMonthly:               50.0,
        PricingCategory:           pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_DYNAMIC,
        SpotInterruptionRiskScore: 0.8,
    }
    err := pluginsdk.ValidateEstimateCostResponse(resp)
    assert.NoError(t, err)

    // Invalid: STANDARD with non-zero risk
    resp.PricingCategory = pb.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD
    err = pluginsdk.ValidateEstimateCostResponse(resp)
    assert.Error(t, err)
}
```

## FAQ

### Q: Do I need to update my plugin immediately?

No. Existing plugins will continue to work. The default values (UNSPECIFIED + 0.0)
pass validation. Update at your convenience before v2.0.

### Q: What if I don't know the pricing category?

Use `UNSPECIFIED` with `0.0` risk score. This is valid but discouraged for new
development. Core systems will treat this as "pricing tier unknown."

### Q: What if I have DYNAMIC pricing but don't know the risk?

Set `pricing_category = DYNAMIC` and `spot_interruption_risk_score = 0.0`. Use
`CheckSpotRiskConsistency()` to log an advisory warning, but the response will
pass validation.

### Q: Will my plugin break if I don't update?

No. Backward compatibility is guaranteed. The UNSPECIFIED + 0.0 combination is
explicitly tested and supported.
