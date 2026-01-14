# Data Model: Standardized Recommendation Reasoning Metadata

## Entities

### RecommendationReason (Enum)
Represents the standardized reason code for a recommendation.

| Value | Name | Description |
|-------|------|-------------|
| 0 | `RECOMMENDATION_REASON_UNSPECIFIED` | Default value. Used when the reason is unknown or not mapped. |
| 1 | `RECOMMENDATION_REASON_OVER_PROVISIONED` | Resource capacity significantly exceeds utilization. |
| 2 | `RECOMMENDATION_REASON_UNDER_PROVISIONED` | Resource capacity is insufficient for the workload (Performance). |
| 3 | `RECOMMENDATION_REASON_IDLE` | Resource is active but has no significant utilization. |
| 4 | `RECOMMENDATION_REASON_REDUNDANT` | Resource provides duplicate functionality (e.g., unused IP). |
| 5 | `RECOMMENDATION_REASON_OBSOLETE_GENERATION` | Resource uses an older hardware generation (e.g., AWS t2 vs t3). |

### Recommendation (Message Update)
The existing `Recommendation` message will be updated with new fields.

| Field | Type | Description |
|-------|------|-------------|
| `primary_reason` | `RecommendationReason` | The main driver for the recommendation. |
| `secondary_reasons` | `repeated RecommendationReason` | Contributing factors. |
