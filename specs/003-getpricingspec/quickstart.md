# Quickstart: GetPricingSpec RPC Enhancement

**Feature**: 003-getpricingspec
**Date**: 2025-11-22

## Overview

This feature enhances the GetPricingSpec RPC to provide transparent pricing breakdowns
with assumptions and tiered pricing support. Plugin developers can now return detailed
pricing information that users can verify and downstream tools can parse.

## Prerequisites

- Go 1.24+
- buf CLI (installed via `make generate`)
- Node.js 22+ (for schema validation)

## Quick Implementation Guide

### 1. Update Proto Definitions

```bash
# The proto changes add:
# - PricingTier message (new)
# - unit, assumptions, pricing_tiers fields to PricingSpec
```

### 2. Regenerate SDK

```bash
make generate
```

### 3. Implement in Your Plugin

```go
import (
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (p *MyPlugin) GetPricingSpec(ctx context.Context, req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error) {
    resource := req.GetResource()

    // Validate required fields (FR-011)
    if resource.GetProvider() == "" || resource.GetResourceType() == "" {
        return nil, status.Error(codes.InvalidArgument,
            "provider and resource_type are required")
    }

    // Check if resource is supported (FR-012)
    if !p.supportsResource(resource) {
        return nil, status.Error(codes.NotFound,
            "unknown region or SKU combination")
    }

    // Return flat-rate pricing
    return &pbc.GetPricingSpecResponse{
        Spec: &pbc.PricingSpec{
            Provider:     resource.GetProvider(),
            ResourceType: resource.GetResourceType(),
            Sku:          resource.GetSku(),
            Region:       resource.GetRegion(),
            BillingMode:  "per_hour",
            RatePerUnit:  0.0104,
            Currency:     "USD",
            Unit:         "hour",
            Description:  "EC2 t3.micro on-demand hourly rate",
            Assumptions: []string{
                "On-demand pricing (not Reserved or Spot)",
                "Linux operating system",
                "Shared tenancy",
            },
        },
    }, nil
}
```

### 4. Return Tiered Pricing

```go
// For resources with tiered pricing (like S3)
return &pbc.GetPricingSpecResponse{
    Spec: &pbc.PricingSpec{
        Provider:     "aws",
        ResourceType: "s3",
        BillingMode:  "tiered",
        Currency:     "USD",
        Unit:         "GB-month",
        Description:  "S3 Standard storage tiered pricing",
        Assumptions:  []string{"Standard storage class"},
        PricingTiers: []*pbc.PricingTier{
            {
                MinQuantity: 0,
                MaxQuantity: 51200,  // 50 TB
                RatePerUnit: 0.023,
                Description: "First 50 TB",
            },
            {
                MinQuantity: 51200,
                MaxQuantity: 512000,  // 500 TB
                RatePerUnit: 0.022,
                Description: "Next 450 TB",
            },
            {
                MinQuantity: 512000,
                MaxQuantity: 0,  // unlimited
                RatePerUnit: 0.021,
                Description: "Over 500 TB",
            },
        },
    },
}, nil
```

### 5. Handle Not-Implemented (FR-008)

```go
// For unsupported resources
return &pbc.GetPricingSpecResponse{
    Spec: &pbc.PricingSpec{
        BillingMode: "not_implemented",
        Unit:        "unknown",
        RatePerUnit: 0,
        Currency:    "USD",
        Description: "Lambda cost estimation not implemented",
        Assumptions: []string{
            "Lambda pricing not available in this version",
            "Returns $0 estimate",
        },
    },
}, nil
```

## Testing Your Implementation

### Using the Test Harness

```go
import (
    plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto"
)

func TestGetPricingSpec(t *testing.T) {
    plugin := &MyPlugin{}
    harness := plugintesting.NewTestHarness(plugin)
    harness.Start(t)
    defer harness.Stop()

    client := harness.Client()

    resp, err := client.GetPricingSpec(context.Background(),
        &pbc.GetPricingSpecRequest{
            Resource: &pbc.ResourceDescriptor{
                Provider:     "aws",
                ResourceType: "ec2",
                Sku:          "t3.micro",
                Region:       "us-east-1",
            },
        })

    require.NoError(t, err)
    require.Equal(t, "per_hour", resp.Spec.BillingMode)
    require.Equal(t, 0.0104, resp.Spec.RatePerUnit)
    require.NotEmpty(t, resp.Spec.Assumptions)
}
```

## Common Patterns

### Both rate_per_unit and pricing_tiers

Per clarification, both can be populated simultaneously:

- Flat billing: use `rate_per_unit`, leave `pricing_tiers` empty
- Tiered billing: use `pricing_tiers`, `rate_per_unit` can be base rate or 0

### Assumptions Best Practices

Include assumptions that help users understand:

- Pricing model (on-demand, reserved, spot)
- Operating system or platform
- Tenancy (shared, dedicated)
- Time calculations (hours/month)

## Validation

```bash
# Run all validations
make validate

# Run specific tests
go test -v ./sdk/go/testing/ -run TestGetPricingSpec
```
