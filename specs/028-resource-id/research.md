# Research: Resource ID and ARN Fields

**Date**: 2025-12-26
**Feature**: 028-resource-id

## Research Topics

### 1. Field Number Availability

**Question**: Are field numbers 7 and 8 available in ResourceDescriptor?

**Finding**: Yes, both field numbers are available.

Current ResourceDescriptor fields:

| Field | Number | Type |
|-------|--------|------|
| provider | 1 | string |
| resource_type | 2 | string |
| sku | 3 | string |
| region | 4 | string |
| tags | 5 | map<string, string> |
| utilization_percentage | 6 | optional double |
| **id** | **7** | **string (NEW)** |
| **arn** | **8** | **string (NEW)** |

**Decision**: Use field 7 for `id` and field 8 for `arn`.

**Rationale**: Sequential numbering follows proto best practices. Fields 1-15 use
single-byte encoding, making them efficient for frequently-used fields.

---

### 2. Existing ARN Field Consistency

**Question**: How is ARN currently used in the codebase?

**Finding**: `GetActualCostRequest` already has an `arn` field (field 5):

```protobuf
message GetActualCostRequest {
  string resource_id = 1;
  google.protobuf.Timestamp start = 2;
  google.protobuf.Timestamp end = 3;
  map<string, string> tags = 4;
  string arn = 5;  // Canonical Cloud Identifier
}
```

**Decision**: Name the new field `arn` for consistency with existing usage.

**Rationale**: Using the same field name across related messages reduces cognitive
load for plugin developers and maintains API consistency.

---

### 3. Cross-Provider Identifier Formats

**Question**: What canonical identifier formats do major cloud providers use?

**Finding**: Each provider has a unique format:

| Provider | Format | Example |
|----------|--------|---------|
| AWS | ARN | `arn:aws:ec2:us-east-1:123456789012:instance/i-abc123` |
| Azure | Resource ID | `/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Compute/virtualMachines/{vm}` |
| GCP | Full Resource Name | `//compute.googleapis.com/projects/{project}/zones/{zone}/instances/{name}` |
| Kubernetes | UID or path | `{cluster}/{namespace}/{kind}/{name}` or UUID |
| Cloudflare | Zone + Resource | `{zone-id}/{resource-type}/{resource-id}` |

**Decision**: Accept all formats in the `arn` field without strict validation.

**Rationale**: Plugins understand their provider's format and can validate/parse
as needed. Protocol layer should not enforce provider-specific validation.

**Alternatives Considered**:

1. Separate fields per provider (rejected: increases proto complexity)
2. Union type with provider enum (rejected: proto3 doesn't support discriminated unions well)
3. Strict ARN format validation (rejected: would break non-AWS providers)

---

### 4. ID Field Semantics

**Question**: Should `id` have any format constraints?

**Finding**: The `id` field is for client correlation, not resource identification.

Common client ID formats:

- Pulumi URNs: `urn:pulumi:prod::myapp::aws:ec2/instance:Instance::webserver`
- UUIDs: `550e8400-e29b-41d4-a716-446655440000`
- Internal tracking IDs: `cost-analysis-batch-2024-001`

**Decision**: Treat `id` as an opaque string with no validation.

**Rationale**: Clients control the format; plugins just pass it through for
correlation. Any validation would be unnecessarily restrictive.

---

### 5. SDK Helper Patterns

**Question**: What helper patterns exist for ResourceDescriptor?

**Finding**: No existing builder pattern for ResourceDescriptor. However, the
codebase uses builder patterns extensively:

- `FocusRecordBuilder` in `pluginsdk/focus_builder.go`
- `ContractCommitmentBuilder` in `pluginsdk/contract_commitment_builder.go`
- `ResourceMatcher` helper in `pluginsdk/helpers.go`

**Decision**: ~~Do not add a ResourceDescriptorBuilder for this feature.~~
**UPDATED**: Add WithID/WithARN helper methods per FR-011 requirement.

**Rationale**: ~~ResourceDescriptor is a simple message with string fields. Direct
struct initialization is cleaner than a builder for this case.~~ User requirement
to maintain SDK consistency with other builder patterns in the codebase and fulfill
FR-011 (pluginsdk MUST provide helper functions for working with both fields).

**What we will add**:

- `WithID(id string)` helper method for setting correlation ID
- `WithARN(arn string)` helper method for setting canonical resource identifier
- Comprehensive proto comments documenting both fields
- Example code in the testing package showing proper usage
- SDK README updates with correlation patterns

---

### 6. Backward Compatibility Verification

**Question**: Will adding these fields break existing plugins?

**Finding**: No, proto3 guarantees forward/backward compatibility for new optional
fields:

1. **Old client → New server**: Server receives empty strings for id/arn, falls
   back to type/sku/region matching (current behavior)
2. **New client → Old server**: Server ignores unknown fields, processes request
   normally (proto3 behavior)
3. **Wire format**: New fields use new field numbers, no reuse of existing tags

**Decision**: Proceed with optional string fields.

**Rationale**: This is the standard proto3 evolution pattern. No migration
required for existing plugins.

---

### 7. ResourceRecommendationInfo Consideration

**Question**: Does ResourceRecommendationInfo need an `id` field for correlation?

**Finding**: `ResourceRecommendationInfo` is the response-side message in
recommendations:

```protobuf
message ResourceRecommendationInfo {
  string resource_id = 1;    // Already exists!
  string resource_type = 2;
  string provider = 3;
  string region = 4;
  string account = 5;
  // ... other fields
}
```

The `resource_id` field (field 1) already exists and can be used for correlation.

**Decision**: Plugins should copy `ResourceDescriptor.id` to
`ResourceRecommendationInfo.resource_id` for correlation.

**Rationale**: No proto changes needed to ResourceRecommendationInfo. The existing
`resource_id` field serves the correlation purpose. Documentation will clarify
this pattern.

---

## Summary of Decisions

| Topic | Decision |
|-------|----------|
| Field numbers | Use 7 for `id`, 8 for `arn` |
| Field naming | `arn` for consistency with GetActualCostRequest |
| ARN formats | Accept all provider formats without validation |
| ID semantics | Opaque pass-through, no validation |
| SDK helpers | No builder; use direct struct + documentation |
| Compatibility | Standard proto3 optional fields, no migration |
| Response correlation | Use existing `resource_id` in ResourceRecommendationInfo |

## Next Steps

Proceed to Phase 1: Design & Contracts
