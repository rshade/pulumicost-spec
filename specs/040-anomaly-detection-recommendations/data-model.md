# Data Model: Anomaly Detection via Recommendations

**Feature**: 040-anomaly-detection-recommendations
**Date**: 2026-01-19

## Overview

This feature extends two existing enums with no new entities or structural changes.

## Enum Extensions

### RecommendationCategory

**Location**: `proto/finfocus/v1/costsource.proto`

| Value | Number | Description |
|-------|--------|-------------|
| RECOMMENDATION_CATEGORY_UNSPECIFIED | 0 | Default/unknown |
| RECOMMENDATION_CATEGORY_COST | 1 | Cost optimization |
| RECOMMENDATION_CATEGORY_PERFORMANCE | 2 | Performance improvement |
| RECOMMENDATION_CATEGORY_SECURITY | 3 | Security enhancement |
| RECOMMENDATION_CATEGORY_RELIABILITY | 4 | Reliability/availability |
| **RECOMMENDATION_CATEGORY_ANOMALY** | **5** | **NEW: Cost anomaly requiring investigation** |

**Validation Rules**:

- ANOMALY recommendations SHOULD have `action_type = INVESTIGATE` (convention, not enforced)
- ANOMALY recommendations SHOULD populate `confidence_score` when provider supports it
- ANOMALY recommendations MAY have negative `estimated_savings` for overspend anomalies

### RecommendationActionType

**Location**: `proto/finfocus/v1/costsource.proto`

| Value | Number | Description |
|-------|--------|-------------|
| RECOMMENDATION_ACTION_TYPE_UNSPECIFIED | 0 | Default/unknown |
| RECOMMENDATION_ACTION_TYPE_RIGHTSIZE | 1 | Resize resource |
| RECOMMENDATION_ACTION_TYPE_TERMINATE | 2 | Delete resource |
| RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT | 3 | Buy reserved/savings plan |
| RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS | 4 | Adjust K8s requests |
| RECOMMENDATION_ACTION_TYPE_MODIFY | 5 | Generic modification |
| RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED | 6 | Remove unused resource |
| RECOMMENDATION_ACTION_TYPE_MIGRATE | 7 | Move to different region/zone |
| RECOMMENDATION_ACTION_TYPE_CONSOLIDATE | 8 | Combine resources |
| RECOMMENDATION_ACTION_TYPE_SCHEDULE | 9 | Start/stop scheduling |
| RECOMMENDATION_ACTION_TYPE_REFACTOR | 10 | Architectural change |
| RECOMMENDATION_ACTION_TYPE_OTHER | 11 | Catch-all |
| **RECOMMENDATION_ACTION_TYPE_INVESTIGATE** | **12** | **NEW: Requires human investigation** |

**Validation Rules**:

- INVESTIGATE actions indicate no automated remediation is appropriate
- INVESTIGATE actions are semantically paired with ANOMALY category (convention)
- INVESTIGATE actions SHOULD include detailed `description` for investigation context

## Field Mapping for Anomalies

Existing `Recommendation` message fields serve anomaly use cases:

| Field | Anomaly Semantic | Example |
|-------|------------------|---------|
| `id` | Unique anomaly identifier | `"aws-anomaly-12345"` |
| `category` | Always `RECOMMENDATION_CATEGORY_ANOMALY` | `5` |
| `action_type` | Typically `RECOMMENDATION_ACTION_TYPE_INVESTIGATE` | `12` |
| `resource.provider` | Cloud provider | `"aws"` |
| `resource.resource_type` | Affected service or resource type | `"ec2"`, `"service"` |
| `resource.region` | Affected region (if applicable) | `"us-east-1"` |
| `impact.estimated_savings` | Anomaly deviation amount (may be negative) | `-1500.00` |
| `impact.currency` | Currency code | `"USD"` |
| `confidence_score` | Anomaly detection confidence (0.0-1.0) | `0.85` |
| `description` | Human-readable anomaly summary | `"150% above baseline"` |
| `metadata` | Provider-specific details | `{"baseline": "1000.00"}` |
| `created_at` | Anomaly detection timestamp | `2026-01-19T10:30:00Z` |

## State Transitions

Not applicable. Anomaly recommendations are stateless snapshots from cost management backends.
They do not have lifecycle states beyond:

- **Created**: Plugin returns anomaly from backend API
- **Filtered**: User applies category/confidence filters
- **Dismissed**: User calls `DismissRecommendation` RPC (existing functionality)

## Relationships

```text
┌─────────────────────────────────────────────────────────────────────────┐
│                          GetRecommendations RPC                         │
└─────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      RecommendationFilter                               │
│  ┌─────────────────┐  ┌──────────────────┐  ┌───────────────────────┐  │
│  │ category=ANOMALY│  │ min_confidence   │  │ action_type=INVESTIGATE│  │
│  └─────────────────┘  └──────────────────┘  └───────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         Recommendation                                   │
│  ┌──────────────┐  ┌───────────────────┐  ┌──────────────────────────┐  │
│  │ ANOMALY (5)  │  │ INVESTIGATE (12)  │  │ confidence_score: 0.85   │  │
│  │ category     │  │ action_type       │  │ estimated_savings: -1500 │  │
│  └──────────────┘  └───────────────────┘  └──────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

## Backward Compatibility

| Client Version | Server Version | Behavior |
|----------------|----------------|----------|
| Old | Old | No change (no anomalies) |
| Old | New | Old client sees ANOMALY as unknown enum (proto3 default) |
| New | Old | New client never receives ANOMALY (old server doesn't generate) |
| New | New | Full anomaly support |

Proto3 unknown enum handling ensures safe interoperability during migration.
