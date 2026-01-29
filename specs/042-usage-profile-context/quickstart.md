# Quickstart: Usage Profile Context

**Feature**: 042-usage-profile-context
**Audience**: Plugin developers implementing profile-aware cost estimation

## Overview

The Usage Profile feature allows the FinFocus Core to signal workload intent (DEV, PROD, BURST)
to plugins. This enables context-aware cost estimation and recommendations.

## 5-Minute Integration

### Step 1: Check for Profile in Request

When handling `GetProjectedCost` or `GetRecommendations`, check the `usage_profile` field:

```go
import (
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
    "github.com/rs/zerolog/log"
)

func (p *MyPlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    // Get the usage profile (defaults to UNSPECIFIED if not set)
    profile := req.GetUsageProfile()

    // Log when applying profile-specific behavior
    if profile != pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED {
        log.Info().
            Str("usage_profile", profile.String()).
            Str("resource_type", req.GetResource().GetResourceType()).
            Msg("Applying profile-specific cost estimation")
    }

    // Apply profile-appropriate defaults
    hours := getMonthlyHours(profile)
    instanceClass := getInstanceClass(profile)

    // ... rest of implementation
}
```

### Step 2: Implement Profile-Specific Logic

Create helper functions for profile-aware defaults. These are plugin-specific implementations;
the examples below show common patterns:

```go
// getMonthlyHours returns the assumed monthly hours for a profile.
// Hour values are guidelines - plugins have discretion for their resource types.
func getMonthlyHours(profile pbc.UsageProfile) float64 {
    switch profile {
    case pbc.UsageProfile_USAGE_PROFILE_PROD:
        return 730 // 24/7 operation
    case pbc.UsageProfile_USAGE_PROFILE_DEV:
        return 160 // ~8 hours/day, 5 days/week
    case pbc.UsageProfile_USAGE_PROFILE_BURST:
        // Burst duration is plugin discretion - batch jobs, load tests, etc.
        // This example uses 200hr; adjust based on your resource type
        return 200
    default:
        // UNSPECIFIED or unknown - use plugin default (typically PROD behavior)
        return 730
    }
}

// getInstanceClass returns the recommended instance class for a profile.
// This is an example for compute resources; adapt for your resource type.
func getInstanceClass(profile pbc.UsageProfile) string {
    switch profile {
    case pbc.UsageProfile_USAGE_PROFILE_PROD:
        return "production" // m5, m6i, etc.
    case pbc.UsageProfile_USAGE_PROFILE_DEV:
        return "burstable" // t3, t4g, etc.
    case pbc.UsageProfile_USAGE_PROFILE_BURST:
        return "scalable" // auto-scaling groups
    default:
        return "production"
    }
}
```

### Step 3: Handle Unknown Profile Values

Use the SDK's `NormalizeUsageProfile` helper for forward compatibility:

```go
import (
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

func (p *MyPlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    // NormalizeUsageProfile returns UNSPECIFIED for unknown values
    // and logs a warning automatically
    profile := pluginsdk.NormalizeUsageProfile(req.GetUsageProfile())

    // ... use normalized profile
}
```

The SDK provides these helpers (see `sdk/go/pluginsdk/usage_profile.go`):

- `IsValidUsageProfile(profile)` - Returns true for known values
- `ParseUsageProfile(s string)` - Parses "dev", "prod", "burst" strings
- `UsageProfileString(profile)` - Returns lowercase string representation
- `NormalizeUsageProfile(profile)` - Returns UNSPECIFIED for unknown values, logs warning

## Complete Example: AWS EC2 Plugin

```go
package awsplugin

import (
    "context"
    "strings"

    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    "github.com/rs/zerolog/log"
)

type AWSComputePlugin struct {
    pbc.UnimplementedCostSourceServiceServer
}

func (p *AWSComputePlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    resource := req.GetResource()
    // Use SDK helper for forward compatibility
    profile := pluginsdk.NormalizeUsageProfile(req.GetUsageProfile())

    // Log profile application (FR-008)
    log.Info().
        Str("usage_profile", profile.String()).
        Str("sku", resource.GetSku()).
        Str("region", resource.GetRegion()).
        Msg("Calculating projected cost")

    // Get profile-appropriate pricing
    hourlyRate := p.getHourlyRate(resource.GetSku(), resource.GetRegion())
    monthlyHours := getMonthlyHours(profile) // Reuse helper from Step 2

    // Apply profile-specific instance recommendation
    recommendedSku := p.getRecommendedSku(resource.GetSku(), profile)

    return &pbc.GetProjectedCostResponse{
        UnitPrice:     hourlyRate,
        Currency:      "USD",
        CostPerMonth:  hourlyRate * monthlyHours,
        BillingDetail: formatBillingDetail(profile, recommendedSku),
    }, nil
}

// getHourlyRate is plugin-specific - fetch from your pricing source
func (p *AWSComputePlugin) getHourlyRate(sku, region string) float64 {
    // Implementation: query AWS pricing API or local cache
    return 0.0 // Placeholder
}

func (p *AWSComputePlugin) getRecommendedSku(
    currentSku string,
    profile pbc.UsageProfile,
) string {
    // For DEV profile, suggest burstable alternatives
    if profile == pbc.UsageProfile_USAGE_PROFILE_DEV {
        if isM5Family(currentSku) {
            return toBurstable(currentSku) // m5.large -> t3.large
        }
    }
    return currentSku
}

// isM5Family checks if SKU is in the M5 instance family (plugin-specific)
func isM5Family(sku string) bool {
    return strings.HasPrefix(strings.ToLower(sku), "m5")
}

// toBurstable converts M5 SKU to equivalent T3 burstable (plugin-specific)
func toBurstable(sku string) string {
    // m5.large -> t3.large, m5.xlarge -> t3.xlarge, etc.
    return strings.Replace(strings.ToLower(sku), "m5", "t3", 1)
}

func formatBillingDetail(profile pbc.UsageProfile, sku string) string {
    switch profile {
    case pbc.UsageProfile_USAGE_PROFILE_DEV:
        return "on-demand (dev workload, 160hr/mo assumed)"
    case pbc.UsageProfile_USAGE_PROFILE_BURST:
        return "on-demand (burst workload, plugin-discretion hours)"
    default:
        return "on-demand (production, 730hr/mo)"
    }
}
```

## Recommendations Integration

For `GetRecommendations`, use the profile to adjust recommendation priorities:

```go
func (p *AWSComputePlugin) GetRecommendations(
    ctx context.Context,
    req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
    profile := normalizeProfile(req.GetUsageProfile())

    recommendations := p.generateRecommendations(req.GetFilter())

    // Adjust priorities based on profile
    for _, rec := range recommendations {
        adjustPriorityForProfile(rec, profile)
    }

    return &pbc.GetRecommendationsResponse{
        Recommendations: recommendations,
    }, nil
}

func adjustPriorityForProfile(rec *pbc.Recommendation, profile pbc.UsageProfile) {
    switch profile {
    case pbc.UsageProfile_USAGE_PROFILE_DEV:
        // DEV: Boost cost savings, lower reliability recommendations
        if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST {
            boostPriority(rec)
        }
    case pbc.UsageProfile_USAGE_PROFILE_PROD:
        // PROD: Balance cost and reliability
        if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_RELIABILITY {
            boostPriority(rec)
        }
    case pbc.UsageProfile_USAGE_PROFILE_BURST:
        // BURST: Focus on scale-out recommendations
        if rec.GetActionType() == pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE {
            boostPriority(rec)
        }
    }
}
```

## Testing Your Implementation

Write tests to verify profile handling:

```go
func TestProfileHandling(t *testing.T) {
    plugin := &AWSComputePlugin{}

    tests := []struct {
        name     string
        profile  pbc.UsageProfile
        wantHrs  float64
    }{
        {"unspecified uses prod hours", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, 730},
        {"prod uses full month", pbc.UsageProfile_USAGE_PROFILE_PROD, 730},
        {"dev uses business hours", pbc.UsageProfile_USAGE_PROFILE_DEV, 160},
        {"burst uses short term", pbc.UsageProfile_USAGE_PROFILE_BURST, 200},
        {"unknown treated as unspecified", pbc.UsageProfile(999), 730},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := &pbc.GetProjectedCostRequest{
                Resource: &pbc.ResourceDescriptor{
                    Provider:     "aws",
                    ResourceType: "ec2",
                    Sku:          "t3.micro",
                    Region:       "us-east-1",
                },
                UsageProfile: tt.profile,
            }

            resp, err := plugin.GetProjectedCost(context.Background(), req)
            require.NoError(t, err)

            // Verify hours assumption is reflected in cost
            expectedCost := resp.GetUnitPrice() * tt.wantHrs
            assert.InDelta(t, expectedCost, resp.GetCostPerMonth(), 0.01)
        })
    }
}
```

## Next Steps

1. **Document your profile behavior**: Add profile-specific assumptions to your plugin's README
2. **Consider DryRun support**: Report profile-dependent fields in DryRun responses
3. **Add benchmarks**: Verify profile handling doesn't impact performance
