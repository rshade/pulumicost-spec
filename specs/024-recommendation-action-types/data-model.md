# Data Model: Extend RecommendationActionType Enum

**Feature**: 019-recommendation-action-types
**Date**: 2025-12-17

## Enums

### RecommendationActionType (Extended)

Type of action recommended. Extended from 7 values to 12 values.

| Value | Number | Description | Source Platforms |
|-------|--------|-------------|------------------|
| `RECOMMENDATION_ACTION_TYPE_UNSPECIFIED` | 0 | Default/unknown action | All |
| `RECOMMENDATION_ACTION_TYPE_RIGHTSIZE` | 1 | Resize resource to optimal size | AWS, Azure, GCP, Kubecost |
| `RECOMMENDATION_ACTION_TYPE_TERMINATE` | 2 | Delete/terminate unused resource | AWS, Azure, GCP |
| `RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT` | 3 | Purchase reserved/committed capacity | AWS, Azure, GCP |
| `RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS` | 4 | Adjust Kubernetes requests/limits | Kubecost |
| `RECOMMENDATION_ACTION_TYPE_MODIFY` | 5 | Modify resource configuration | All |
| `RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED` | 6 | Delete unused/orphaned resources | AWS, Azure, GCP |
| `RECOMMENDATION_ACTION_TYPE_MIGRATE` | 7 | **NEW** Move workloads to different regions/zones/SKUs | Azure Advisor, GCP Recommender |
| `RECOMMENDATION_ACTION_TYPE_CONSOLIDATE` | 8 | **NEW** Combine multiple resources into fewer, larger ones | Azure Advisor, Kubecost |
| `RECOMMENDATION_ACTION_TYPE_SCHEDULE` | 9 | **NEW** Start/stop resources on schedule (dev/test) | AWS Instance Scheduler, Azure Automation |
| `RECOMMENDATION_ACTION_TYPE_REFACTOR` | 10 | **NEW** Architectural changes (e.g., move to serverless) | GCP Recommender |
| `RECOMMENDATION_ACTION_TYPE_OTHER` | 11 | **NEW** Provider-specific recommendations not fitting other categories | All |

## Proto Definition Changes

### Before (Current)

```protobuf
// RecommendationActionType indicates the type of action recommended.
enum RecommendationActionType {
  RECOMMENDATION_ACTION_TYPE_UNSPECIFIED = 0;
  RECOMMENDATION_ACTION_TYPE_RIGHTSIZE = 1;
  RECOMMENDATION_ACTION_TYPE_TERMINATE = 2;
  RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT = 3;
  RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS = 4;
  RECOMMENDATION_ACTION_TYPE_MODIFY = 5;
  RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED = 6;
}
```

### After (Extended)

```protobuf
// RecommendationActionType indicates the type of action recommended.
enum RecommendationActionType {
  RECOMMENDATION_ACTION_TYPE_UNSPECIFIED = 0;
  RECOMMENDATION_ACTION_TYPE_RIGHTSIZE = 1;
  RECOMMENDATION_ACTION_TYPE_TERMINATE = 2;
  RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT = 3;
  RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS = 4;
  RECOMMENDATION_ACTION_TYPE_MODIFY = 5;
  RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED = 6;
  // Move workloads to different regions, zones, or SKUs for cost optimization.
  // Common in Azure Advisor (region migration) and GCP Recommender (zone migration).
  RECOMMENDATION_ACTION_TYPE_MIGRATE = 7;
  // Combine multiple smaller resources into fewer, larger ones for efficiency.
  // Common in Azure Advisor and Kubecost for node consolidation.
  RECOMMENDATION_ACTION_TYPE_CONSOLIDATE = 8;
  // Start/stop resources on a schedule (e.g., dev/test environments).
  // Common in AWS Instance Scheduler and Azure Automation.
  RECOMMENDATION_ACTION_TYPE_SCHEDULE = 9;
  // Architectural changes such as moving to serverless or managed services.
  // Common in GCP Recommender for App Engine/Cloud Run suggestions.
  RECOMMENDATION_ACTION_TYPE_REFACTOR = 10;
  // Catch-all for provider-specific recommendations not fitting other categories.
  // Use when no other action type accurately describes the recommendation.
  RECOMMENDATION_ACTION_TYPE_OTHER = 11;
}
```

## Messages (No Changes)

No new message types are required. The existing `ModifyAction` message handles all new action
types through its `modification_type` field and `current_config`/`recommended_config` maps.

### ModifyAction Usage for New Action Types

| Action Type | modification_type | current_config | recommended_config |
|-------------|-------------------|----------------|-------------------|
| MIGRATE | "region_migration" | region: "us-east-1" | region: "us-west-2" |
| MIGRATE | "sku_migration" | sku: "Standard_D4s" | sku: "Standard_D2s" |
| CONSOLIDATE | "node_consolidation" | node_count: "5" | node_count: "3" |
| SCHEDULE | "start_stop_schedule" | schedule: "none" | schedule: "weekdays-only" |
| REFACTOR | "serverless_migration" | compute_type: "vm" | compute_type: "cloud_run" |
| OTHER | (provider-specific) | (varies) | (varies) |

## Validation Rules (No Changes)

Existing validation rules apply to new enum values:

| Entity | Field | Rule |
|--------|-------|------|
| Recommendation | action_type | Required, valid enum value (now includes 7-11) |
| RecommendationFilter | action_type | If present, valid enum value (now includes 7-11) |

## Relationships (No Changes)

The enum extension does not affect message relationships. `RecommendationActionType` continues
to be used in:

- `Recommendation.action_type` (field 3)
- `RecommendationFilter.action_type` (field 5)
- `RecommendationSummary.count_by_action_type` (map keys)
- `RecommendationSummary.savings_by_action_type` (map keys)

## Backward Compatibility

| Scenario | Behavior |
|----------|----------|
| Old plugin, new core | Plugin returns 0-6, core accepts normally |
| New plugin, old core | Plugin returns 7-11, core sees as unknown numeric value |
| Old client filtering by 7-11 | Filter ignored or returns empty (graceful degradation) |
| Round-trip serialization | All values 0-11 preserved correctly |

## Generated Code Impact

After `make generate`:

### Go Constants (sdk/go/proto/finfocus/v1/costsource.pb.go)

```go
const (
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED        RecommendationActionType = 0
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE          RecommendationActionType = 1
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE          RecommendationActionType = 2
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT RecommendationActionType = 3
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS    RecommendationActionType = 4
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY             RecommendationActionType = 5
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED      RecommendationActionType = 6
    // NEW
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE            RecommendationActionType = 7
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE        RecommendationActionType = 8
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE           RecommendationActionType = 9
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR           RecommendationActionType = 10
    RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER              RecommendationActionType = 11
)
```

### String Methods

The generated `.String()` method will automatically include string representations for all
new values, enabling proper logging and JSON serialization.
