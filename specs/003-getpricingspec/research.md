# Research: GetPricingSpec RPC Enhancement

**Feature**: 003-getpricingspec
**Date**: 2025-11-22

## Research Tasks

### 1. Current Proto State Analysis

**Question**: What fields already exist in PricingSpec and what needs to be added?

**Finding**: Current PricingSpec message (lines 136-162) has:

- `provider`, `resource_type`, `sku`, `region` (duplicating ResourceDescriptor)
- `billing_mode`, `rate_per_unit`, `currency`, `description` (matches spec)
- `metric_hints`, `plugin_metadata`, `source` (not in spec requirements)

**Missing per FR-003**:

- `unit` (string) - billing unit (e.g., "hour", "GB-month")
- `assumptions` (repeated string) - pricing assumptions list
- `pricing_tiers` (repeated PricingTier) - tiered pricing support

**Decision**: Add three new fields to PricingSpec; create new PricingTier message
**Rationale**: Backward compatible (new optional fields), aligns with spec requirements
**Alternatives**: Considered creating new response message, rejected due to breaking change

### 2. Protobuf Field Numbering Strategy

**Question**: What field numbers to use for new PricingSpec fields?

**Finding**: Current highest field number in PricingSpec is 11 (source)

**Decision**: Use field numbers 12-14 for new fields

- `unit` = 12
- `assumptions` = 13
- `pricing_tiers` = 14

**Rationale**: Sequential numbering, fields 1-15 are single-byte encoded (efficient)
**Alternatives**: None required - standard practice

### 3. PricingTier Message Design

**Question**: What fields does PricingTier need per FR-004?

**Finding**: FR-004 requires: min_quantity, max_quantity, rate_per_unit, description

**Decision**: Create PricingTier message with:

```protobuf
message PricingTier {
  double min_quantity = 1;
  double max_quantity = 2;  // 0 = unlimited
  double rate_per_unit = 3;
  string description = 4;
}
```

**Rationale**: Matches spec exactly, double type for quantities supports fractional values
**Alternatives**: Considered using int64 for quantities, rejected for flexibility

### 4. Unit Field vs Existing metric_hints

**Question**: Does new `unit` field duplicate `metric_hints.unit`?

**Finding**: metric_hints is for usage guidance, unit is the billing unit for rate_per_unit

**Decision**: Keep both - they serve different purposes

- `unit`: what rate_per_unit is priced in (e.g., "hour" for $0.01/hour)
- `metric_hints`: guidance on what usage metrics to collect

**Rationale**: Different semantics; unit describes the rate, metric_hints describes data collection
**Alternatives**: Remove metric_hints - rejected, would break existing usage

### 5. gRPC Error Handling Pattern

**Question**: How to implement FR-011/FR-012 error responses?

**Finding**: Proto already has ErrorCode enum with:

- `ERROR_CODE_INVALID_RESOURCE` (6) - for InvalidArgument
- `ERROR_CODE_RESOURCE_NOT_FOUND` (7) - for NotFound

**Decision**: Use standard gRPC status codes with ErrorDetail in metadata

```go
// Missing required fields
return nil, status.Error(codes.InvalidArgument, "provider and resource_type are required")

// Unknown region/SKU
return nil, status.Error(codes.NotFound, "unknown region or SKU combination")
```

**Rationale**: Standard gRPC pattern, ErrorDetail can be added to metadata for structured errors
**Alternatives**: Return valid response with error flag - rejected per clarification answer

### 6. JSON Schema Updates

**Question**: What schema changes needed for new fields?

**Finding**: schemas/pricing_spec.schema.json needs updates

**Decision**: Add to schema:

- `unit`: string type
- `assumptions`: array of strings
- `pricing_tiers`: array of objects with min_quantity, max_quantity, rate_per_unit, description

**Rationale**: Schema must match proto for JSON serialization compatibility
**Alternatives**: None - schema and proto must align

### 7. Backward Compatibility Verification

**Question**: Will adding fields break existing plugins?

**Finding**: Adding optional fields to protobuf messages is backward compatible

- Existing plugins sending old format will have empty/default values for new fields
- New plugins can populate new fields without breaking old consumers

**Decision**: Safe to add fields; no buf breaking changes expected

**Rationale**: Protobuf wire format guarantees forward/backward compatibility for new fields
**Alternatives**: None needed - standard protobuf behavior

## Summary

All technical decisions resolved. Implementation can proceed with:

1. Add PricingTier message to proto
2. Add unit, assumptions, pricing_tiers fields to PricingSpec
3. Update JSON schema with matching fields
4. Implement gRPC error handling per standard patterns
5. Update examples with new fields
6. Regenerate SDK via `make generate`
