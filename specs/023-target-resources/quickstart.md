# Quickstart: Target Resources for Recommendations

**Feature**: 019-target-resources

## Overview

The `target_resources` field allows clients to request recommendations for a specific set of
resources rather than all resources in scope. This enables stack-scoped recommendations,
pre-deployment analysis, and batch resource queries.

## Basic Usage

### Request Recommendations for Specific Resources

```go
import (
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Define target resources (e.g., from a Pulumi stack)
targets := []*pbc.ResourceDescriptor{
    {
        Provider:     "aws",
        ResourceType: "ec2",
        Sku:          "t3.medium",
        Region:       "us-east-1",
    },
    {
        Provider:     "aws",
        ResourceType: "rds",
        Sku:          "db.t3.micro",
        Region:       "us-east-1",
    },
}

// Create request with target resources
req := &pbc.GetRecommendationsRequest{
    TargetResources: targets,
    Filter: &pbc.RecommendationFilter{
        Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
    },
}

// Call the service
resp, err := client.GetRecommendations(ctx, req)
if err != nil {
    return fmt.Errorf("get recommendations: %w", err)
}

// Process recommendations (only for target resources)
for _, rec := range resp.Recommendations {
    fmt.Printf("Recommendation for %s: %s\n",
        rec.Resource.ResourceType, rec.Description)
}
```

### Backward Compatible Usage

```go
// Empty target_resources = existing behavior (all recommendations)
req := &pbc.GetRecommendationsRequest{
    TargetResources: nil, // or empty slice
    Filter: &pbc.RecommendationFilter{
        Provider: "aws",
    },
}
```

## Matching Behavior

### Strict Matching

When optional fields are specified, they must match exactly:

```go
// This target matches ONLY t3.medium in us-east-1
target := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
    Sku:          "t3.medium",  // Must match exactly
    Region:       "us-east-1",  // Must match exactly
}

// This target matches ANY ec2 instance (any SKU, any region)
broadTarget := &pbc.ResourceDescriptor{
    Provider:     "aws",
    ResourceType: "ec2",
    // sku and region omitted = not checked
}
```

### AND Logic with Filter

target_resources and filter are combined with AND logic:

```go
// Returns recommendations that:
// 1. Match one of the target resources, AND
// 2. Have HIGH priority
req := &pbc.GetRecommendationsRequest{
    TargetResources: []*pbc.ResourceDescriptor{
        {Provider: "aws", ResourceType: "ec2"},
        {Provider: "aws", ResourceType: "rds"},
    },
    Filter: &pbc.RecommendationFilter{
        Priority: pbc.RecommendationPriority_RECOMMENDATION_PRIORITY_HIGH,
    },
}
```

## Validation

### Maximum Resources

```go
// Maximum 100 resources per request
targets := make([]*pbc.ResourceDescriptor, 101)
// ...populate targets...

req := &pbc.GetRecommendationsRequest{TargetResources: targets}
// Returns InvalidArgument: "target_resources exceeds maximum of 100"
```

### Resource Validation

Each resource descriptor is validated:

```go
targets := []*pbc.ResourceDescriptor{
    {Provider: "invalid"}, // Missing resource_type, invalid provider
}
// Returns InvalidArgument with details about the invalid resource
```

## Plugin Implementation

Plugins implementing `GetRecommendations` should:

1. Check if `target_resources` is non-empty
2. If non-empty, filter recommendations to matching resources
3. Apply `filter` criteria after target filtering
4. Return empty list if no matches (not an error)

```go
func (p *MyPlugin) GetRecommendations(
    ctx context.Context,
    req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
    // Get all recommendations from backend
    allRecs := p.fetchRecommendations(ctx)

    // Apply target_resources filtering
    if len(req.GetTargetResources()) > 0 {
        allRecs = filterByTargetResources(allRecs, req.GetTargetResources())
    }

    // Apply filter criteria
    if req.GetFilter() != nil {
        allRecs = applyFilter(allRecs, req.GetFilter())
    }

    return &pbc.GetRecommendationsResponse{
        Recommendations: allRecs,
    }, nil
}
```
