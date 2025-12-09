# Data Model: GetRecommendations RPC

**Feature**: 013-recommendations-rpc
**Date**: 2025-12-04

## Enums

### RecommendationCategory

Classification of recommendation types.

| Value | Number | Description |
|-------|--------|-------------|
| `RECOMMENDATION_CATEGORY_UNSPECIFIED` | 0 | Default/unknown category |
| `RECOMMENDATION_CATEGORY_COST` | 1 | Cost optimization recommendation |
| `RECOMMENDATION_CATEGORY_PERFORMANCE` | 2 | Performance improvement recommendation |
| `RECOMMENDATION_CATEGORY_SECURITY` | 3 | Security enhancement recommendation |
| `RECOMMENDATION_CATEGORY_RELIABILITY` | 4 | Reliability/availability recommendation |

### RecommendationActionType

Type of action recommended.

| Value | Number | Description |
|-------|--------|-------------|
| `RECOMMENDATION_ACTION_TYPE_UNSPECIFIED` | 0 | Default/unknown action |
| `RECOMMENDATION_ACTION_TYPE_RIGHTSIZE` | 1 | Resize resource to optimal size |
| `RECOMMENDATION_ACTION_TYPE_TERMINATE` | 2 | Delete/terminate unused resource |
| `RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT` | 3 | Purchase reserved/committed capacity |
| `RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS` | 4 | Adjust Kubernetes requests/limits |
| `RECOMMENDATION_ACTION_TYPE_MODIFY` | 5 | Modify resource configuration |
| `RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED` | 6 | Delete unused/orphaned resources |

### RecommendationPriority

Priority level for recommendations.

| Value | Number | Description |
|-------|--------|-------------|
| `RECOMMENDATION_PRIORITY_UNSPECIFIED` | 0 | Default/unknown priority |
| `RECOMMENDATION_PRIORITY_LOW` | 1 | Low priority, minor impact |
| `RECOMMENDATION_PRIORITY_MEDIUM` | 2 | Medium priority |
| `RECOMMENDATION_PRIORITY_HIGH` | 3 | High priority, significant impact |
| `RECOMMENDATION_PRIORITY_CRITICAL` | 4 | Critical priority, immediate action |

## Messages

### GetRecommendationsRequest

Request message for GetRecommendations RPC.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `filter` | RecommendationFilter | 1 | Optional filter criteria |
| `projection_period` | string | 2 | "daily", "monthly" (default), "annual" |
| `page_size` | int32 | 3 | Maximum results per page (default: 50) |
| `page_token` | string | 4 | Continuation token from previous response |

### GetRecommendationsResponse

Response message for GetRecommendations RPC.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `recommendations` | repeated Recommendation | 1 | List of recommendations |
| `summary` | RecommendationSummary | 2 | Aggregated summary statistics |
| `next_page_token` | string | 3 | Token for next page (empty if last) |

### RecommendationFilter

Filter criteria for narrowing recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `provider` | string | 1 | Filter by provider (aws, azure, gcp, kubernetes) |
| `region` | string | 2 | Filter by region |
| `resource_type` | string | 3 | Filter by resource type |
| `category` | RecommendationCategory | 4 | Filter by category |
| `action_type` | RecommendationActionType | 5 | Filter by action type |

### Recommendation

A single cost optimization recommendation.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `id` | string | 1 | Unique recommendation identifier |
| `category` | RecommendationCategory | 2 | Recommendation category |
| `action_type` | RecommendationActionType | 3 | Recommended action type |
| `resource` | ResourceInfo | 4 | Affected resource details |
| `rightsize` | RightsizeAction | 5 | oneof: Rightsizing details |
| `terminate` | TerminateAction | 6 | oneof: Termination details |
| `commitment` | CommitmentAction | 7 | oneof: Commitment purchase details |
| `kubernetes` | KubernetesAction | 8 | oneof: Kubernetes action details |
| `modify` | ModifyAction | 9 | oneof: Modification details |
| `impact` | RecommendationImpact | 10 | Financial impact assessment |
| `priority` | RecommendationPriority | 11 | Recommendation priority |
| `confidence_score` | optional double | 12 | Confidence (0.0-1.0), nil if unavailable |
| `description` | string | 13 | Human-readable description |
| `reasoning` | repeated string | 14 | List of reasons for recommendation |
| `source` | string | 15 | Data source (aws, kubecost, etc.) |
| `created_at` | google.protobuf.Timestamp | 16 | When recommendation was generated |
| `metadata` | map<string, string> | 17 | Additional provider-specific data |

### ResourceInfo

Information about the resource being recommended for action.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `id` | string | 1 | Resource identifier |
| `name` | string | 2 | Resource name |
| `provider` | string | 3 | Cloud provider |
| `resource_type` | string | 4 | Resource type |
| `region` | string | 5 | Deployment region |
| `sku` | string | 6 | SKU/instance type |
| `tags` | map<string, string> | 7 | Resource tags |
| `utilization` | ResourceUtilization | 8 | Current utilization metrics |

### ResourceUtilization

Current resource utilization metrics.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `cpu_percent` | double | 1 | CPU utilization percentage |
| `memory_percent` | double | 2 | Memory utilization percentage |
| `storage_percent` | double | 3 | Storage utilization percentage |
| `network_in_mbps` | double | 4 | Network ingress (Mbps) |
| `network_out_mbps` | double | 5 | Network egress (Mbps) |
| `custom_metrics` | map<string, double> | 6 | Provider-specific metrics |

### RightsizeAction

Details for rightsizing recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `current_sku` | string | 1 | Current SKU/size |
| `recommended_sku` | string | 2 | Recommended SKU/size |
| `current_instance_type` | string | 3 | Current instance type |
| `recommended_instance_type` | string | 4 | Recommended instance type |
| `projected_utilization` | ResourceUtilization | 5 | Expected utilization after resize |

### TerminateAction

Details for termination recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `termination_reason` | string | 1 | Reason for termination |
| `idle_days` | int32 | 2 | Days resource has been idle |

### CommitmentAction

Details for commitment purchase recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `commitment_type` | string | 1 | Type: reserved_instance, savings_plan, cud |
| `term` | string | 2 | Term: 1_year, 3_year |
| `payment_option` | string | 3 | Payment option |
| `recommended_quantity` | double | 4 | Recommended purchase quantity |
| `scope` | string | 5 | Scope (account, region, etc.) |

### KubernetesAction

Details for Kubernetes resource adjustment recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `cluster_id` | string | 1 | Kubernetes cluster identifier |
| `namespace` | string | 2 | Kubernetes namespace |
| `controller_kind` | string | 3 | Controller type (Deployment, StatefulSet) |
| `controller_name` | string | 4 | Controller name |
| `container_name` | string | 5 | Container name |
| `current_requests` | KubernetesResources | 6 | Current resource requests |
| `recommended_requests` | KubernetesResources | 7 | Recommended requests |
| `current_limits` | KubernetesResources | 8 | Current resource limits |
| `recommended_limits` | KubernetesResources | 9 | Recommended limits |
| `algorithm` | string | 10 | Recommendation algorithm used |

### KubernetesResources

Kubernetes CPU and memory resource specifications.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `cpu` | string | 1 | CPU specification (e.g., "100m", "2") |
| `memory` | string | 2 | Memory specification (e.g., "256Mi", "2Gi") |

### ModifyAction

Details for generic modification recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `modification_type` | string | 1 | Type of modification |
| `current_config` | map<string, string> | 2 | Current configuration |
| `recommended_config` | map<string, string> | 3 | Recommended configuration |

### RecommendationImpact

Financial impact of implementing a recommendation.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `estimated_savings` | double | 1 | Estimated savings amount |
| `currency` | string | 2 | ISO 4217 currency code |
| `projection_period` | string | 3 | Period for projection |
| `current_cost` | double | 4 | Current cost |
| `projected_cost` | double | 5 | Projected cost after action |
| `savings_percentage` | double | 6 | Savings as percentage |
| `implementation_cost` | optional double | 7 | One-time implementation cost |
| `migration_effort_hours` | optional double | 8 | Estimated migration effort |

### RecommendationSummary

Aggregated summary of recommendations.

| Field | Type | Number | Description |
|-------|------|--------|-------------|
| `total_recommendations` | int32 | 1 | Total count of recommendations |
| `total_estimated_savings` | double | 2 | Total estimated savings |
| `currency` | string | 3 | ISO 4217 currency code |
| `projection_period` | string | 4 | Period for projection |
| `count_by_category` | map<string, int32> | 5 | Count per category |
| `savings_by_category` | map<string, double> | 6 | Savings per category |
| `count_by_action_type` | map<string, int32> | 7 | Count per action type |
| `savings_by_action_type` | map<string, double> | 8 | Savings per action type |

## Relationships

```text
GetRecommendationsRequest
    └── RecommendationFilter (optional)

GetRecommendationsResponse
    ├── Recommendation (repeated)
    │   ├── ResourceInfo
    │   │   └── ResourceUtilization (optional)
    │   ├── [oneof action_detail]
    │   │   ├── RightsizeAction
    │   │   │   └── ResourceUtilization (projected)
    │   │   ├── TerminateAction
    │   │   ├── CommitmentAction
    │   │   ├── KubernetesAction
    │   │   │   ├── KubernetesResources (current_requests)
    │   │   │   ├── KubernetesResources (recommended_requests)
    │   │   │   ├── KubernetesResources (current_limits)
    │   │   │   └── KubernetesResources (recommended_limits)
    │   │   └── ModifyAction
    │   └── RecommendationImpact
    └── RecommendationSummary
```

## Validation Rules

| Entity | Field | Rule |
|--------|-------|------|
| Recommendation | id | Required, non-empty |
| Recommendation | category | Required, valid enum value |
| Recommendation | action_type | Required, valid enum value |
| Recommendation | resource | Required |
| Recommendation | impact | Required |
| Recommendation | confidence_score | If present, 0.0 <= value <= 1.0 |
| RecommendationImpact | currency | Required, valid ISO 4217 code |
| RecommendationImpact | estimated_savings | Required |
| RecommendationSummary | currency | Required, valid ISO 4217 code |
| ResourceInfo | id | Required, non-empty |
| ResourceInfo | provider | Required, non-empty |
| GetRecommendationsRequest | projection_period | "daily", "monthly", "annual" |
| GetRecommendationsRequest | page_size | 0 < value <= 1000 (0 uses default 50) |
