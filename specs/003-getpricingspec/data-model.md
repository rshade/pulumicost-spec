# Data Model: GetPricingSpec RPC Enhancement

**Feature**: 003-getpricingspec
**Date**: 2025-11-22

## Entities

### PricingSpec (Updated)

Existing message with new fields added.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| provider | string | Yes | Cloud provider identifier |
| resource_type | string | Yes | Resource type (e.g., "ec2", "s3") |
| sku | string | No | Provider SKU or instance size |
| region | string | No | Deployment region |
| billing_mode | string | Yes | How resource is billed (per_hour, per_gb_month, tiered, not_implemented) |
| rate_per_unit | double | Yes | Price per billing unit |
| currency | string | Yes | Currency code (e.g., "USD") |
| description | string | No | Human-readable pricing description |
| metric_hints | repeated UsageMetricHint | No | Guidance on usage metrics |
| plugin_metadata | map<string, string> | No | Plugin-specific metadata |
| source | string | No | Pricing data source |
| **unit** | string | **NEW** | Billing unit for rate_per_unit (e.g., "hour", "GB-month") |
| **assumptions** | repeated string | **NEW** | List of pricing assumptions |
| **pricing_tiers** | repeated PricingTier | **NEW** | Tiered pricing breakdown |

### PricingTier (New)

New message for tiered pricing support.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| min_quantity | double | Yes | Minimum quantity for tier (inclusive) |
| max_quantity | double | Yes | Maximum quantity for tier (0 = unlimited) |
| rate_per_unit | double | Yes | Rate for this tier |
| description | string | No | Tier description (e.g., "First 50 TB") |

### Relationships

```text
GetPricingSpecRequest
└── ResourceDescriptor (1:1)

GetPricingSpecResponse
└── PricingSpec (1:1)
    └── PricingTier (1:N)
```

## Validation Rules

### PricingSpec Validation

- `billing_mode` MUST be one of: per_hour, per_gb_month, per_request, flat, per_day, per_cpu_hour, tiered, not_implemented
- `currency` SHOULD be valid ISO 4217 code (informational only, no strict validation)
- `rate_per_unit` MUST be >= 0
- If `billing_mode` is "not_implemented", then `rate_per_unit` MUST be 0
- `assumptions` SHOULD contain at least one entry when `billing_mode` is not "not_implemented"

### PricingTier Validation

- `min_quantity` MUST be >= 0
- `max_quantity` MUST be >= min_quantity OR 0 (unlimited)
- `rate_per_unit` MUST be >= 0
- Tiers MUST be contiguous (no gaps between max of tier N and min of tier N+1)
- Tiers MUST be ordered by min_quantity ascending

### Error Conditions

| Condition | gRPC Status | ErrorCode |
|-----------|-------------|-----------|
| Missing provider or resource_type | InvalidArgument | ERROR_CODE_INVALID_RESOURCE |
| Unknown region/SKU combination | NotFound | ERROR_CODE_RESOURCE_NOT_FOUND |
| Unsupported provider | InvalidArgument | ERROR_CODE_INVALID_PROVIDER |

## State Transitions

N/A - PricingSpec is stateless. Each request returns current pricing information.

## Protobuf Field Numbers

| Message | Field | Number | Wire Type |
|---------|-------|--------|-----------|
| PricingSpec | unit | 12 | string |
| PricingSpec | assumptions | 13 | repeated string |
| PricingSpec | pricing_tiers | 14 | repeated message |
| PricingTier | min_quantity | 1 | double |
| PricingTier | max_quantity | 2 | double |
| PricingTier | rate_per_unit | 3 | double |
| PricingTier | description | 4 | string |
