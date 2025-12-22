# Quickstart: Extended RecommendationActionType Enum

**Feature**: 019-recommendation-action-types
**Date**: 2025-12-17

## Overview

This feature extends the `RecommendationActionType` enum with 5 new values to support
additional FinOps platform recommendation categories.

## New Action Types

| Action Type | Value | Use Case |
|-------------|-------|----------|
| MIGRATE | 7 | Move workloads to different regions/zones/SKUs |
| CONSOLIDATE | 8 | Combine resources into fewer, larger ones |
| SCHEDULE | 9 | Start/stop resources on schedule |
| REFACTOR | 10 | Architectural changes (e.g., serverless) |
| OTHER | 11 | Provider-specific catch-all |

## Usage Examples

### Plugin Developer: Return Migration Recommendation

```go
import (
    pb "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func (p *MyPlugin) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsRequest) (
    *pb.GetRecommendationsResponse, error) {

    recommendations := []*pb.Recommendation{
        {
            Id:         "rec-migrate-001",
            Category:   pb.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
            ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
            Resource: &pb.ResourceInfo{
                Id:           "i-1234567890abcdef0",
                Name:         "web-server-1",
                Provider:     "aws",
                ResourceType: "ec2_instance",
                Region:       "us-east-1",
            },
            Modify: &pb.ModifyAction{
                ModificationType: "region_migration",
                CurrentConfig: map[string]string{
                    "region": "us-east-1",
                },
                RecommendedConfig: map[string]string{
                    "region": "us-west-2",
                },
            },
            Impact: &pb.RecommendationImpact{
                EstimatedSavings: 150.00,
                Currency:         "USD",
                ProjectionPeriod: "monthly",
            },
            Description: "Migrate to us-west-2 for 15% cost reduction",
        },
    }

    return &pb.GetRecommendationsResponse{
        Recommendations: recommendations,
    }, nil
}
```

### Plugin Developer: Return Schedule Recommendation

```go
{
    Id:         "rec-schedule-001",
    ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
    Resource: &pb.ResourceInfo{
        Id:           "dev-cluster",
        Name:         "dev-environment",
        Provider:     "azure",
        ResourceType: "aks_cluster",
    },
    Modify: &pb.ModifyAction{
        ModificationType: "start_stop_schedule",
        CurrentConfig: map[string]string{
            "running": "24x7",
        },
        RecommendedConfig: map[string]string{
            "schedule": "weekdays-8am-6pm",
            "timezone": "America/New_York",
        },
    },
    Impact: &pb.RecommendationImpact{
        EstimatedSavings: 500.00,
        Currency:         "USD",
        ProjectionPeriod: "monthly",
        SavingsPercentage: 60.0,
    },
    Description: "Schedule dev cluster for business hours only",
}
```

### Core CLI: Filter by New Action Types

```go
// Filter to show only migration recommendations
filter := &pb.RecommendationFilter{
    ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
}

// Filter for scheduling opportunities (dev/test environments)
filter := &pb.RecommendationFilter{
    ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
}

// Filter for consolidation recommendations (node/resource merging)
filter := &pb.RecommendationFilter{
    ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
}

// Filter for architectural refactoring (serverless migration)
filter := &pb.RecommendationFilter{
    ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
}

// Filter for provider-specific recommendations
filter := &pb.RecommendationFilter{
    ActionType: pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
}
```

### Handling Unknown Action Types (Backward Compatibility)

```go
func displayActionType(actionType pb.RecommendationActionType) string {
    switch actionType {
    case pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE:
        return "Migrate"
    case pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE:
        return "Consolidate"
    case pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE:
        return "Schedule"
    case pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR:
        return "Refactor"
    case pb.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER:
        return "Other"
    // ... existing cases ...
    default:
        // Handle unknown values gracefully (backward compatibility)
        return fmt.Sprintf("Unknown (%d)", actionType)
    }
}
```

## Testing the Changes

### Verify Proto Compilation

```bash
make generate
make lint
```

### Run Conformance Tests

```bash
go test -v ./sdk/go/testing/ -run TestConformance
```

### Verify Serialization Round-Trip

```bash
go test -v ./sdk/go/testing/ -run TestActionTypeSerialization
```

## Migration Notes

- **Existing plugins**: No changes required. Continue using action types 0-6.
- **New recommendations**: Use action types 7-11 for better categorization.
- **Filtering**: Core and CLI can filter by new action types immediately after SDK update.
- **Display**: Add human-readable labels for new action types in UI/CLI.

## Related Documentation

- [Data Model](data-model.md) - Complete enum definition and usage patterns
- [Research](research.md) - Decision rationale for new action types
- [Spec](spec.md) - Full feature specification
