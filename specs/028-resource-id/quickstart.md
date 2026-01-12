# Quickstart: Resource ID and ARN Fields

**Date**: 2025-12-26
**Feature**: 028-resource-id

## Overview

This feature adds two new optional fields to `ResourceDescriptor`:

- **`id`** (field 7): Client correlation identifier for batch request/response
  matching
- **`arn`** (field 8): Canonical cloud resource identifier for exact resource
  matching

## For Plugin Developers

### Handling the `id` Field

The `id` field is an opaque pass-through value. Copy it to responses unchanged:

```go
func (p *Plugin) GetRecommendations(
    ctx context.Context,
    req *pb.GetRecommendationsRequest,
) (*pb.GetRecommendationsResponse, error) {
    var recommendations []*pb.Recommendation

    for _, target := range req.TargetResources {
        // Get recommendations for this resource...
        recs := p.findRecommendations(target)

        for _, rec := range recs {
            // Copy the client's correlation ID to the response
            rec.Resource.ResourceId = target.Id
            recommendations = append(recommendations, rec)
        }
    }

    return &pb.GetRecommendationsResponse{
        Recommendations: recommendations,
    }, nil
}
```

### Using the `arn` Field for Exact Matching

When `arn` is provided, use it for precise resource lookup:

```go
func (p *Plugin) matchResource(target *pb.ResourceDescriptor) (*Resource, error) {
    // Prefer ARN for exact matching
    if target.Arn != "" {
        resource, err := p.lookupByARN(target.Arn)
        if err == nil {
            return resource, nil
        }
        // Log warning and fall back
        p.logger.Warn().
            Str("arn", target.Arn).
            Err(err).
            Msg("ARN lookup failed, falling back to fuzzy match")
    }

    // Fall back to type/sku/region/tags matching
    return p.fuzzyMatch(target)
}
```

### AWS-Specific ARN Parsing

For AWS plugins, use the standard library:

```go
import "github.com/aws/aws-sdk-go-v2/aws/arn"

func (p *AWSPlugin) lookupByARN(arnStr string) (*Resource, error) {
    parsed, err := arn.Parse(arnStr)
    if err != nil {
        return nil, fmt.Errorf("invalid ARN: %w", err)
    }

    // Use parsed.Service, parsed.Region, parsed.Resource
    return p.client.GetResource(parsed)
}
```

## For finfocus-core Developers

### Setting Correlation IDs

When building batch requests, set unique IDs for correlation:

```go
func buildRecommendationsRequest(resources []*Resource) *pb.GetRecommendationsRequest {
    var descriptors []*pb.ResourceDescriptor

    for _, r := range resources {
        descriptors = append(descriptors, &pb.ResourceDescriptor{
            Provider:     r.Provider,
            ResourceType: r.Type,
            Sku:          r.SKU,
            Region:       r.Region,
            Id:           r.URN,                    // Use Pulumi URN for correlation
            Arn:          r.Outputs["arn"].(string), // Use cloud ARN if available
        })
    }

    return &pb.GetRecommendationsRequest{
        TargetResources: descriptors,
    }
}
```

### Correlating Responses

Use the ID to match recommendations to resources:

```go
func processRecommendations(
    req *pb.GetRecommendationsRequest,
    resp *pb.GetRecommendationsResponse,
    resources map[string]*Resource,
) {
    // Build lookup map from request
    resourceByID := make(map[string]*Resource)
    for _, desc := range req.TargetResources {
        if desc.Id != "" {
            resourceByID[desc.Id] = resources[desc.Id]
        }
    }

    // Match recommendations to resources
    for _, rec := range resp.Recommendations {
        id := rec.Resource.ResourceId
        if resource, ok := resourceByID[id]; ok {
            // Recommendation correlates to this resource
            resource.AddRecommendation(rec)
        }
    }
}
```

## Testing

### Unit Test Example

```go
func TestResourceDescriptorWithIDAndARN(t *testing.T) {
    descriptor := &pb.ResourceDescriptor{
        Provider:     "aws",
        ResourceType: "ec2",
        Sku:          "t3.micro",
        Region:       "us-east-1",
        Id:           "urn:pulumi:prod::app::aws:ec2/instance:Instance::web",
        Arn:          "arn:aws:ec2:us-east-1:123456789012:instance/i-abc123",
    }

    // Verify fields are set correctly
    assert.Equal(t, "urn:pulumi:prod::app::aws:ec2/instance:Instance::web", descriptor.Id)
    assert.Equal(t, "arn:aws:ec2:us-east-1:123456789012:instance/i-abc123", descriptor.Arn)

    // Verify empty defaults for backward compatibility
    emptyDesc := &pb.ResourceDescriptor{Provider: "aws", ResourceType: "ec2"}
    assert.Empty(t, emptyDesc.Id)
    assert.Empty(t, emptyDesc.Arn)
}
```

### Integration Test Example

```go
func TestBatchCorrelation(t *testing.T) {
    harness := testing.NewTestHarness(myPlugin)
    harness.Start(t)
    defer harness.Stop()

    // Create batch request with IDs
    req := &pb.GetRecommendationsRequest{
        TargetResources: []*pb.ResourceDescriptor{
            {Provider: "aws", ResourceType: "ec2", Id: "res-001"},
            {Provider: "aws", ResourceType: "ec2", Id: "res-002"},
            {Provider: "aws", ResourceType: "ec2", Id: "res-003"},
        },
    }

    resp, err := harness.Client().GetRecommendations(ctx, req)
    require.NoError(t, err)

    // Verify each recommendation has matching ID
    ids := make(map[string]bool)
    for _, rec := range resp.Recommendations {
        ids[rec.Resource.ResourceId] = true
    }

    // All input IDs should appear in responses
    assert.Contains(t, ids, "res-001")
    assert.Contains(t, ids, "res-002")
    assert.Contains(t, ids, "res-003")
}
```

## Migration Guide

### Existing Plugins

No changes required. Both fields default to empty strings, and existing
type/sku/region/tags matching continues to work.

### Opting In

1. Update to latest finfocus-spec
2. Run `make generate` to regenerate proto bindings
3. Add correlation ID pass-through in recommendation handlers
4. (Optional) Add ARN-based exact matching for improved accuracy

## Troubleshooting

### ID Not Appearing in Response

- Ensure plugin copies `target.Id` to `rec.Resource.ResourceId`
- Check that client is setting `Id` on request descriptors

### ARN Lookup Failing

- Verify ARN format matches provider expectations
- Check plugin logs for format validation warnings
- Ensure fallback to fuzzy matching is implemented

### Batch Correlation Mismatches

- Verify unique IDs per request (duplicates cause ambiguous matches)
- Check that all target resources have IDs set
