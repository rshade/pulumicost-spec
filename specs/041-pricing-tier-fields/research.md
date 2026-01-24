# Phase 0 Research: Pricing Tier & Probability Fields

## Existing Proto Definitions

### `FocusPricingCategory` (in `proto/finfocus/v1/enums.proto`)

```protobuf
enum FocusPricingCategory {
  FOCUS_PRICING_CATEGORY_UNSPECIFIED = 0;
  FOCUS_PRICING_CATEGORY_STANDARD = 1;
  FOCUS_PRICING_CATEGORY_COMMITTED = 2;
  FOCUS_PRICING_CATEGORY_DYNAMIC = 3;
  FOCUS_PRICING_CATEGORY_OTHER = 4;
}
```

- **Assessment**: Suitable for the requirement. `COMMITTED` maps to "Reserved", `DYNAMIC` maps to "Spot".

### Target RPC Responses (in `proto/finfocus/v1/costsource.proto`)

#### `EstimateCostResponse`

Current fields:

1. `string currency`
2. `double cost_monthly`

#### `GetProjectedCostResponse`

Current fields:

1. `double unit_price`
2. `string currency`
3. `double cost_per_month`
4. `string billing_detail`
5. `repeated ImpactMetric impact_metrics`
6. `GrowthType growth_type`
7. `DryRunResponse dry_run_result`

## Technical Approach

1. **Protobuf Update**:
   - Inject `FocusPricingCategory pricing_category` and `double spot_interruption_risk_score` into both response messages.
   - Maintain backward compatibility by using new field numbers.
2. **SDK Generation**: Use `make generate` to update Go and TypeScript SDKs.
3. **SDK Polish**: Ensure the new fields are documented in the SDK and examples are updated if necessary.
