# Quickstart: Implementing GetRecommendations

**Feature**: 013-recommendations-rpc

## Overview

This guide shows how to implement the `GetRecommendations` RPC in a PulumiCost plugin.
The implementation is optional - plugins that don't support recommendations simply don't
implement the `RecommendationsProvider` interface.

## Prerequisites

- Go 1.25.4+
- Existing PulumiCost plugin implementation
- `pulumicost-spec` SDK imported

## Step 1: Implement the RecommendationsProvider Interface

```go
package myplugin

import (
    "context"

    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// MyPlugin implements the PulumiCost Plugin interface
type MyPlugin struct {
    // ... existing fields ...
}

// Ensure MyPlugin implements RecommendationsProvider
var _ pluginsdk.RecommendationsProvider = (*MyPlugin)(nil)

// GetRecommendations retrieves cost optimization recommendations.
// This method is optional - only implement if your backend supports recommendations.
func (p *MyPlugin) GetRecommendations(
    ctx context.Context,
    req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
    // 1. Fetch recommendations from your backend
    backendRecs, err := p.fetchFromBackend(ctx, req.GetFilter())
    if err != nil {
        return nil, err
    }

    // 2. Convert to proto messages
    var recommendations []*pbc.Recommendation
    for _, rec := range backendRecs {
        recommendations = append(recommendations, convertToProto(rec))
    }

    // 3. Apply filtering
    filtered := applyFilter(recommendations, req.GetFilter())

    // 4. Apply pagination
    page, nextToken := paginate(filtered, int(req.GetPageSize()), req.GetPageToken())

    // 5. Calculate summary
    summary := calculateSummary(filtered, req.GetProjectionPeriod())

    return &pbc.GetRecommendationsResponse{
        Recommendations: page,
        Summary:         summary,
        NextPageToken:   nextToken,
    }, nil
}
```

## Step 2: Declare Capability in Supports (Optional)

For deterministic capability discovery, declare recommendations support in `Supports`:

```go
func (p *MyPlugin) Supports(
    ctx context.Context,
    req *pbc.SupportsRequest,
) (*pbc.SupportsResponse, error) {
    // ... existing logic ...

    return &pbc.SupportsResponse{
        Supported: true,
        Reason:    "",
        Capabilities: map[string]bool{
            "recommendations": true,  // Declare recommendations support
        },
    }, nil
}
```

## Step 3: Create Recommendation Messages

### AWS Rightsizing Example

```go
func createAWSRightsizing(instance *aws.Instance, rec *aws.Recommendation) *pbc.Recommendation {
    return &pbc.Recommendation{
        Id:         rec.ID,
        Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
        ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
        Resource: &pbc.ResourceRecommendationInfo{
            Id:           instance.ID,
            Name:         instance.Name,
            Provider:     "aws",
            ResourceType: "ec2",
            Region:       instance.Region,
            Sku:          instance.InstanceType,
            Utilization: &pbc.ResourceUtilization{
                CpuPercent:    rec.CurrentUtilization.CPU,
                MemoryPercent: rec.CurrentUtilization.Memory,
            },
        },
        ActionDetail: &pbc.Recommendation_Rightsize{
            Rightsize: &pbc.RightsizeAction{
                CurrentInstanceType:     instance.InstanceType,
                RecommendedInstanceType: rec.TargetInstanceType,
            },
        },
        Impact: &pbc.RecommendationImpact{
            EstimatedSavings:  rec.MonthlySavings,
            Currency:          "USD",
            ProjectionPeriod:  "monthly",
            CurrentCost:       rec.CurrentMonthlyCost,
            ProjectedCost:     rec.CurrentMonthlyCost - rec.MonthlySavings,
            SavingsPercentage: (rec.MonthlySavings / rec.CurrentMonthlyCost) * 100,
        },
        Priority:    pbc.RecommendationPriority_RECOMMENDATION_PRIORITY_MEDIUM,
        Description: fmt.Sprintf("Rightsize %s from %s to %s", instance.Name,
                                 instance.InstanceType, rec.TargetInstanceType),
        Reasoning:   rec.FindingReasonCodes,
        Source:      "aws-cost-explorer",
        CreatedAt:   timestamppb.Now(),
    }
}
```

### Kubernetes Request Sizing Example

```go
func createK8sRequestSizing(rec *kubecost.Recommendation) *pbc.Recommendation {
    return &pbc.Recommendation{
        Id:         rec.ID,
        Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
        ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS,
        Resource: &pbc.ResourceRecommendationInfo{
            Id:           fmt.Sprintf("%s/%s/%s", rec.ClusterID, rec.Namespace, rec.Controller),
            Name:         rec.Controller,
            Provider:     "kubernetes",
            ResourceType: "container",
            Region:       rec.ClusterID,
        },
        ActionDetail: &pbc.Recommendation_Kubernetes{
            Kubernetes: &pbc.KubernetesAction{
                ClusterId:      rec.ClusterID,
                Namespace:      rec.Namespace,
                ControllerKind: rec.ControllerKind,
                ControllerName: rec.Controller,
                ContainerName:  rec.Container,
                CurrentRequests: &pbc.KubernetesResources{
                    Cpu:    rec.CurrentRequest.CPU,
                    Memory: rec.CurrentRequest.Memory,
                },
                RecommendedRequests: &pbc.KubernetesResources{
                    Cpu:    rec.RecommendedRequest.CPU,
                    Memory: rec.RecommendedRequest.Memory,
                },
                Algorithm: rec.Algorithm,
            },
        },
        Impact: &pbc.RecommendationImpact{
            EstimatedSavings: rec.MonthlySavings,
            Currency:         "USD",
            ProjectionPeriod: "monthly",
        },
        Priority:    pbc.RecommendationPriority_RECOMMENDATION_PRIORITY_HIGH,
        Description: fmt.Sprintf("Adjust requests for %s/%s", rec.Namespace, rec.Container),
        Source:      "kubecost",
        CreatedAt:   timestamppb.Now(),
    }
}
```

## Step 4: Add Metrics (Optional)

Enable Prometheus metrics for recommendations:

```go
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    UnaryInterceptors: []grpc.UnaryServerInterceptor{
        pluginsdk.MetricsUnaryServerInterceptor("my-plugin"),
    },
}
```

The interceptor automatically records:

- `pulumicost_plugin_requests_total{grpc_method="...",grpc_code="...",plugin_name="..."}`
- `pulumicost_plugin_request_duration_seconds{grpc_method="...",plugin_name="..."}`

## Step 5: Test with Conformance Suite

```go
func TestGetRecommendations(t *testing.T) {
    plugin := NewMyPlugin()
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    // Test basic retrieval
    resp, err := harness.Client().GetRecommendations(context.Background(),
        &pbc.GetRecommendationsRequest{})
    require.NoError(t, err)
    assert.NotNil(t, resp.Summary)

    // Test filtering
    resp, err = harness.Client().GetRecommendations(context.Background(),
        &pbc.GetRecommendationsRequest{
            Filter: &pbc.RecommendationFilter{
                Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
            },
        })
    require.NoError(t, err)
    for _, rec := range resp.Recommendations {
        assert.Equal(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST, rec.Category)
    }

    // Test pagination
    resp, err = harness.Client().GetRecommendations(context.Background(),
        &pbc.GetRecommendationsRequest{PageSize: 10})
    require.NoError(t, err)
    assert.LessOrEqual(t, len(resp.Recommendations), 10)
}
```

## Summary Calculation

Helper function for calculating the response summary:

```go
func calculateSummary(recs []*pbc.Recommendation, period string) *pbc.RecommendationSummary {
    summary := &pbc.RecommendationSummary{
        TotalRecommendations: int32(len(recs)),
        ProjectionPeriod:     period,
        CountByCategory:      make(map[string]int32),
        SavingsByCategory:    make(map[string]float64),
        CountByActionType:    make(map[string]int32),
        SavingsByActionType:  make(map[string]float64),
    }

    var totalSavings float64
    var currency string

    for _, rec := range recs {
        cat := rec.Category.String()
        action := rec.ActionType.String()

        summary.CountByCategory[cat]++
        summary.CountByActionType[action]++

        if rec.Impact != nil {
            savings := rec.Impact.EstimatedSavings
            totalSavings += savings
            summary.SavingsByCategory[cat] += savings
            summary.SavingsByActionType[action] += savings
            currency = rec.Impact.Currency
        }
    }

    summary.TotalEstimatedSavings = totalSavings
    summary.Currency = currency

    return summary
}
```

## Error Handling

Return appropriate gRPC status codes:

```go
import "google.golang.org/grpc/status"
import "google.golang.org/grpc/codes"

func (p *MyPlugin) GetRecommendations(ctx context.Context, req *pbc.GetRecommendationsRequest) (
    *pbc.GetRecommendationsResponse, error) {

    // Invalid pagination token
    if req.GetPageToken() != "" && !isValidToken(req.GetPageToken()) {
        return nil, status.Error(codes.InvalidArgument, "invalid page_token")
    }

    // Backend unavailable
    recs, err := p.fetchFromBackend(ctx)
    if err != nil {
        return nil, status.Error(codes.Unavailable, "recommendation service unavailable")
    }

    // ...
}
```

## Next Steps

1. Implement `GetRecommendations` in your plugin
2. Run conformance tests to validate
3. Enable metrics for production monitoring
4. Document supported recommendation types in your plugin README
