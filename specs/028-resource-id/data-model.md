# Data Model: Resource ID and ARN Fields

**Date**: 2025-12-26
**Feature**: 028-resource-id

## Entity Changes

### ResourceDescriptor (Modified)

The `ResourceDescriptor` protobuf message is extended with two new optional fields.

**Current Fields** (unchanged):

| Field | Number | Type | Required | Description |
|-------|--------|------|----------|-------------|
| provider | 1 | string | Yes | Cloud provider ("aws", "azure", "gcp", "kubernetes", "custom") |
| resource_type | 2 | string | Yes | Resource type identifier |
| sku | 3 | string | No | Provider-specific SKU |
| region | 4 | string | No | Deployment region |
| tags | 5 | map<string, string> | No | Resource tags |
| utilization_percentage | 6 | optional double | No | Override for utilization (0.0-1.0) |

**New Fields**:

| Field | Number | Type | Required | Description |
|-------|--------|------|----------|-------------|
| id | 7 | string | No | Client correlation identifier (opaque pass-through) |
| arn | 8 | string | No | Canonical cloud resource identifier |

### Field Semantics

#### `id` (field 7)

- **Purpose**: Request/response correlation in batch operations
- **Set by**: Client (finfocus-core)
- **Validated by**: None (opaque)
- **Pass-through**: Plugins MUST copy to response unchanged
- **Default**: Empty string (no correlation)

**Valid formats** (examples, not exhaustive):

```text
urn:pulumi:prod::myapp::aws:ec2/instance:Instance::webserver
550e8400-e29b-41d4-a716-446655440000
cost-batch-2024-12-26-001
my-custom-tracking-id
```

#### `arn` (field 8)

- **Purpose**: Exact resource identification for precise cost lookups
- **Set by**: Client (from cloud provider outputs)
- **Validated by**: Plugin (format check, optional)
- **Precedence**: When provided, takes precedence over type/sku/region matching
- **Default**: Empty string (use fuzzy matching)

**Valid formats by provider**:

| Provider | Format Pattern | Example |
|----------|----------------|---------|
| AWS | `arn:aws:{service}:{region}:{account}:{resource}` | `arn:aws:ec2:us-east-1:123456789012:instance/i-abc123` |
| Azure | `/subscriptions/{sub}/resourceGroups/{rg}/providers/{provider}/{type}/{name}` | `/subscriptions/sub-1/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1` |
| GCP | `//{service}.googleapis.com/projects/{project}/{scope}/{name}` | `//compute.googleapis.com/projects/my-project/zones/us-central1-a/instances/vm-1` |
| Kubernetes | `{cluster}/{namespace}/{kind}/{name}` or UID | `prod-cluster/default/Deployment/nginx` |
| Cloudflare | `{zone-id}/{resource-type}/{resource-id}` | `abc123/dns_record/xyz789` |

## Relationship to Response Messages

### GetRecommendationsRequest → GetRecommendationsResponse Correlation

```text
┌─────────────────────────────────────────────────────────────────┐
│ GetRecommendationsRequest                                       │
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ target_resources[]:ResourceDescriptor                       │ │
│ │   ├─ id: "res-001"  ←── correlation ID                      │ │
│ │   ├─ arn: "arn:aws:ec2:..."  ←── exact match                │ │
│ │   ├─ provider: "aws"                                        │ │
│ │   ├─ resource_type: "ec2"                                   │ │
│ │   └─ ...                                                    │ │
│ └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ GetRecommendationsResponse                                      │
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ recommendations[]:Recommendation                            │ │
│ │   └─ resource:ResourceRecommendationInfo                    │ │
│ │       └─ resource_id: "res-001"  ←── copied from id         │ │
│ └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Matching Logic (Plugin Side)

```text
┌────────────────────────────────────────────┐
│ ResourceDescriptor received                │
├────────────────────────────────────────────┤
│ arn provided?                              │
│   ├─ Yes → Parse ARN, exact resource match │
│   │        (use ARN for API queries)       │
│   │                                        │
│   └─ No → Fuzzy match by:                  │
│           provider + resource_type +       │
│           sku + region + tags              │
│                                            │
│ id provided?                               │
│   ├─ Yes → Include in response for         │
│   │        client correlation              │
│   │                                        │
│   └─ No → Omit from response               │
│           (backward compatible)            │
└────────────────────────────────────────────┘
```

## State Transitions

N/A - Both fields are stateless request parameters. No lifecycle or state
transitions apply.

## Validation Rules

### Protocol Layer (buf lint)

- Field names follow snake_case convention ✓
- Field numbers are in valid range (1-536870911) ✓
- No reserved field numbers used ✓

### Runtime Layer (Plugin)

| Field | Validation | On Invalid |
|-------|------------|------------|
| id | None | N/A (always valid) |
| arn | Optional format check | Log warning, fall back to fuzzy match |

### Client Layer (finfocus-core)

| Field | Validation | On Invalid |
|-------|------------|------------|
| id | Must be unique per batch request | Client responsibility |
| arn | Should match provider format | Client responsibility |

## Data Volume Considerations

- **Field size**: id typically <500 chars, arn typically <2048 chars
- **Batch size**: Up to 100 ResourceDescriptors per request (existing limit)
- **Memory overhead**: ~2.5 KB additional per request at max batch size
- **Wire format impact**: Negligible (optional fields, empty by default)
