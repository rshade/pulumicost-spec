# Quickstart: Implementing Anomaly Recommendations

**Feature**: 040-anomaly-detection-recommendations
**Audience**: Plugin developers implementing anomaly detection

## Overview

This guide explains how to return cost anomalies as recommendations using the new
`RECOMMENDATION_CATEGORY_ANOMALY` and `RECOMMENDATION_ACTION_TYPE_INVESTIGATE` enum values.

## Prerequisites

- FinFocus SDK version with anomaly support (v0.X.X+)
- Access to a cost management backend that provides anomaly data (AWS, Azure, or similar)

## Step 1: Check Provider Anomaly API Support

| Provider | API | SDK Support |
|----------|-----|-------------|
| AWS | Cost Anomaly Detection | `aws-sdk-go-v2/service/costexplorer` |
| Azure | Cost Management Anomalies | `azure-sdk-for-go/sdk/resourcemanager/costmanagement` |
| GCP | Budget Alerts | `cloud.google.com/go/billing/budgets` |

## Step 2: Map Provider Data to Recommendation

### AWS Cost Anomaly Detection Example

```go
import (
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func mapAWSAnomalyToRecommendation(anomaly *costexplorer.Anomaly) *pbc.Recommendation {
    // Normalize AWS anomaly score (0-100) to confidence score (0.0-1.0)
    confidence := float64(*anomaly.AnomalyScore) / 100.0

    // Calculate impact (may be negative for overspend)
    impact := -(*anomaly.Impact.TotalActualSpend - *anomaly.Impact.TotalExpectedSpend)

    return &pbc.Recommendation{
        Id:         *anomaly.AnomalyId,
        Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
        ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE,
        Resource: &pbc.ResourceRecommendationInfo{
            Provider:     "aws",
            ResourceType: *anomaly.DimensionValue, // e.g., "SERVICE" or "LINKED_ACCOUNT"
            Region:       extractRegion(anomaly),
        },
        Impact: &pbc.RecommendationImpact{
            EstimatedSavings: impact,
            Currency:         "USD",
        },
        ConfidenceScore: &confidence,
        Description:     formatAnomalyDescription(anomaly),
        Source:          "aws-cost-anomaly-detection",
        Metadata: map[string]string{
            "baseline_amount":    fmt.Sprintf("%.2f", *anomaly.Impact.TotalExpectedSpend),
            "actual_amount":      fmt.Sprintf("%.2f", *anomaly.Impact.TotalActualSpend),
            "anomaly_start_date": anomaly.AnomalyStartDate.String(),
            "anomaly_end_date":   anomaly.AnomalyEndDate.String(),
        },
    }
}
```

### Azure Cost Management Example

```go
func mapAzureAnomalyToRecommendation(anomaly *armcostmanagement.Anomaly) *pbc.Recommendation {
    confidence := *anomaly.Properties.ConfidenceScore

    return &pbc.Recommendation{
        Id:         *anomaly.ID,
        Category:   pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
        ActionType: pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE,
        Resource: &pbc.ResourceRecommendationInfo{
            Provider:     "azure",
            ResourceType: *anomaly.Properties.ResourceType,
            Region:       *anomaly.Properties.Location,
        },
        Impact: &pbc.RecommendationImpact{
            EstimatedSavings: -*anomaly.Properties.AnomalyAmount, // Negative for overspend
            Currency:         "USD",
        },
        ConfidenceScore: &confidence,
        Description:     *anomaly.Properties.Message,
        Source:          "azure-cost-management",
    }
}
```

## Step 3: Implement GetRecommendations Handler

```go
func (p *MyPlugin) GetRecommendations(
    ctx context.Context,
    req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
    var recommendations []*pbc.Recommendation

    // Check if anomalies should be included (no filter or ANOMALY filter)
    includeAnomalies := req.Filter == nil ||
        req.Filter.Category == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_UNSPECIFIED ||
        req.Filter.Category == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY

    if includeAnomalies {
        // Fetch anomalies from backend
        anomalies, err := p.fetchAnomalies(ctx)
        if err != nil {
            // Log but don't fail - anomalies are supplementary
            p.logger.Warn().Err(err).Msg("Failed to fetch anomalies")
        } else {
            recommendations = append(recommendations, anomalies...)
        }
    }

    // Fetch traditional recommendations
    optimizations, err := p.fetchOptimizations(ctx, req.Filter)
    if err != nil {
        return nil, err
    }
    recommendations = append(recommendations, optimizations...)

    // Apply confidence score filter if specified
    if req.Filter != nil && req.Filter.MinConfidenceScore > 0 {
        recommendations = filterByConfidence(recommendations, req.Filter.MinConfidenceScore)
    }

    return &pbc.GetRecommendationsResponse{
        Recommendations: recommendations,
        Summary:         computeSummary(recommendations),
    }, nil
}
```

## Step 4: Handle Negative Savings

Anomalies representing overspend have negative `estimated_savings`:

```go
// For overspend anomalies, estimated_savings is negative
if anomaly.ActualCost > anomaly.ExpectedCost {
    impact.EstimatedSavings = -(anomaly.ActualCost - anomaly.ExpectedCost)
}

// For underspend anomalies (rare), estimated_savings could be positive
if anomaly.ActualCost < anomaly.ExpectedCost {
    impact.EstimatedSavings = anomaly.ExpectedCost - anomaly.ActualCost
}
```

## Step 5: Testing

Write a conformance test for your anomaly implementation:

```go
func TestAnomalyRecommendations(t *testing.T) {
    plugin := NewMyPlugin()
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    // Test: Anomaly recommendations are returned
    resp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{})
    require.NoError(t, err)

    anomalies := filterByCategory(resp.Recommendations,
        pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY)

    for _, rec := range anomalies {
        assert.Equal(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY, rec.Category)
        assert.Equal(t, pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_INVESTIGATE, rec.ActionType)
        assert.NotEmpty(t, rec.Description)
    }

    // Test: Category filter works
    filteredResp, err := harness.Client().GetRecommendations(ctx, &pbc.GetRecommendationsRequest{
        Filter: &pbc.RecommendationFilter{
            Category: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY,
        },
    })
    require.NoError(t, err)

    for _, rec := range filteredResp.Recommendations {
        assert.Equal(t, pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_ANOMALY, rec.Category)
    }
}
```

## Best Practices

1. **Always set confidence score** when your provider supports it
2. **Use descriptive metadata keys**: `baseline_amount`, `deviation_percent`, `detection_time`
3. **Include investigation context** in the `description` field
4. **Handle missing anomaly data gracefully** - don't fail the entire request
5. **Log anomaly fetch errors** but return partial results

## Troubleshooting

| Issue | Solution |
|-------|----------|
| No anomalies returned | Check provider API credentials and anomaly detection configuration |
| Confidence score is nil | Provider may not support confidence; leave nil (optional field) |
| Filter not working | Ensure filter logic checks for UNSPECIFIED as "include all" |
