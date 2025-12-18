# Data Model: Target Resources for Recommendations

**Feature**: 019-target-resources
**Date**: 2025-12-17

## Entity Overview

This feature extends an existing message rather than introducing new entities.

## Modified Messages

### GetRecommendationsRequest (Extended)

**File**: `proto/pulumicost/v1/costsource.proto`
**Lines**: 652-667 (current), adding field 6

```protobuf
message GetRecommendationsRequest {
  // Existing fields (1-5)
  RecommendationFilter filter = 1;
  string projection_period = 2;
  int32 page_size = 3;
  string page_token = 4;
  repeated string excluded_recommendation_ids = 5;

  // NEW FIELD
  repeated ResourceDescriptor target_resources = 6;
}
```

### ResourceDescriptor (Reused)

**File**: `proto/pulumicost/v1/costsource.proto`
**Lines**: 181-218 (unchanged)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| provider | string | Yes | Cloud provider: aws, azure, gcp, kubernetes, custom |
| resource_type | string | Yes | Resource type (max 256 chars, alphanumeric + -:/) |
| sku | string | No | Provider-specific SKU (e.g., t3.micro, Standard_B1s) |
| region | string | No | Deployment region (e.g., us-east-1, eastus) |
| tags | map<string,string> | No | Label/tag hints for matching |

## Validation Rules

### target_resources Field

| Rule | Value | Error |
|------|-------|-------|
| Max entries | 100 | InvalidArgument: "target_resources exceeds maximum of 100" |
| Item validation | Each ResourceDescriptor validated | InvalidArgument: "target_resources[N]: {error}" |
| Empty list | Allowed | Preserves existing behavior (all recommendations) |
| Duplicates | Allowed | Processed without error, results deduplicated |

### ResourceDescriptor Validation (Existing)

| Field | Rule | Error |
|-------|------|-------|
| provider | Required, must be valid enum | ErrEmptyProvider, ErrInvalidProvider |
| resource_type | Required, max 256 chars, format regex | ErrEmptyResourceType, ErrEmptyResourceTypeFmt |
| tags | Max 50 entries, key max 128, value max 256 | ErrTooManyTags, ErrTagKeyTooLong, ErrTagValueTooLong |

## Matching Logic

### Resource-to-Recommendation Matching

A recommendation matches a target resource when:

1. **provider** equals (required)
2. **resource_type** equals (required)
3. **sku** equals IF specified in target (optional)
4. **region** equals IF specified in target (optional)
5. **tags** all specified tags present on recommendation's resource (subset match)

### Request Processing Logic

```text
INPUT: GetRecommendationsRequest with target_resources and filter

1. Validate target_resources (length, each item)
2. IF target_resources is empty:
     Return all recommendations matching filter (existing behavior)
3. ELSE:
     a. Get all recommendations from backend
     b. Filter to recommendations matching ANY target resource
     c. Apply RecommendationFilter criteria (AND logic)
     d. Return filtered results
```

## State Transitions

N/A - This is a stateless request/response extension. No persistent state changes.

## Relationships

```text
GetRecommendationsRequest
  └── target_resources: repeated ResourceDescriptor (NEW)
        └── (matches against) Recommendation.resource: ResourceRecommendationInfo
```

**Note**: `ResourceDescriptor` (input) matches against `ResourceRecommendationInfo` (output).
These are different message types but share common fields for matching purposes.
