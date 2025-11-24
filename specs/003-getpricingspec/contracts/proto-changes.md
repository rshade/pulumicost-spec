# Proto Contract Changes: GetPricingSpec Enhancement

**Feature**: 003-getpricingspec
**Date**: 2025-11-22

## New Message: PricingTier

Add after existing PricingSpec message in `proto/pulumicost/v1/costsource.proto`:

```protobuf
// PricingTier represents one tier in a tiered pricing model.
message PricingTier {
  // min_quantity is the minimum quantity for this tier (inclusive)
  double min_quantity = 1;
  // max_quantity is the maximum quantity for this tier (exclusive, 0 = unlimited)
  double max_quantity = 2;
  // rate_per_unit is the rate charged for this tier
  double rate_per_unit = 3;
  // description provides human-readable tier description (e.g., "First 50 TB")
  string description = 4;
}
```

## Updated Message: PricingSpec

Add new fields to existing PricingSpec message:

```protobuf
message PricingSpec {
  // ... existing fields 1-11 ...

  // unit specifies the billing unit for rate_per_unit (e.g., "hour", "GB-month", "request")
  string unit = 12;
  // assumptions contains human-readable strings explaining pricing derivation
  repeated string assumptions = 13;
  // pricing_tiers contains tiered pricing breakdown (for billing_mode="tiered")
  repeated PricingTier pricing_tiers = 14;
}
```

## Wire Format Compatibility

- **Backward Compatible**: Yes - adding new optional fields
- **Forward Compatible**: Yes - old clients ignore new fields
- **buf breaking**: Expected to PASS

## Expected buf lint/breaking Results

```bash
# Should pass without errors
make generate
make lint
```

## Generated Go Code Impact

New generated types in `sdk/go/proto/`:

```go
type PricingTier struct {
    MinQuantity  float64 `protobuf:"fixed64,1,opt,name=min_quantity,json=minQuantity,proto3"`
    MaxQuantity  float64 `protobuf:"fixed64,2,opt,name=max_quantity,json=maxQuantity,proto3"`
    RatePerUnit  float64 `protobuf:"fixed64,3,opt,name=rate_per_unit,json=ratePerUnit,proto3"`
    Description  string  `protobuf:"bytes,4,opt,name=description,proto3"`
}

// PricingSpec updated with new fields
type PricingSpec struct {
    // ... existing fields ...
    Unit         string         `protobuf:"bytes,12,opt,name=unit,proto3"`
    Assumptions  []string       `protobuf:"bytes,13,rep,name=assumptions,proto3"`
    PricingTiers []*PricingTier `protobuf:"bytes,14,rep,name=pricing_tiers,json=pricingTiers,proto3"`
}
```
