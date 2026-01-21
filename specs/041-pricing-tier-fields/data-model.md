# Data Model: Pricing Tier & Probability Fields

## Message Updates

### `finfocus.v1.EstimateCostResponse`

| Field                          | Type                   | Number | Description                                                                 |
| ------------------------------ | ---------------------- | ------ | --------------------------------------------------------------------------- |
| `pricing_category`             | `FocusPricingCategory` | 3      | Categorizes the pricing model applied (Standard, Committed, Dynamic).       |
| `spot_interruption_risk_score` | `double`               | 4      | Probability-based reliability risk for Dynamic/Spot instances (0.0 to 1.0). |

### `finfocus.v1.GetProjectedCostResponse`

| Field                          | Type                   | Number | Description                                                                 |
| ------------------------------ | ---------------------- | ------ | --------------------------------------------------------------------------- |
| `pricing_category`             | `FocusPricingCategory` | 8      | Categorizes the pricing model applied (Standard, Committed, Dynamic).       |
| `spot_interruption_risk_score` | `double`               | 9      | Probability-based reliability risk for Dynamic/Spot instances (0.0 to 1.0). |

## Enum Reuse

### `FocusPricingCategory`

- `FOCUS_PRICING_CATEGORY_STANDARD`: On-demand pricing.
- `FOCUS_PRICING_CATEGORY_COMMITTED`: Reserved Instances, Savings Plans.
- `FOCUS_PRICING_CATEGORY_DYNAMIC`: Spot instances, Preemptible VMs.

## Constraints & Validation

- `spot_interruption_risk_score` MUST be between 0.0 and 1.0.
- `spot_interruption_risk_score` SHOULD only be populated when `pricing_category` is `FOCUS_PRICING_CATEGORY_DYNAMIC`.
