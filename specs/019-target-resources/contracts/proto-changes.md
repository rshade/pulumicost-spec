# Proto Contract Changes: Target Resources

**Feature**: 019-target-resources
**File**: `proto/pulumicost/v1/costsource.proto`

## Change Summary

Add `repeated ResourceDescriptor target_resources = 6` to `GetRecommendationsRequest` message.

## Proto Diff

```protobuf
// GetRecommendationsRequest contains parameters for retrieving recommendations.
message GetRecommendationsRequest {
  // filter narrows the recommendations returned
  RecommendationFilter filter = 1;
  // projection_period specifies the time period for savings projection
  // Valid values: "daily", "monthly" (default), "annual"
  string projection_period = 2;
  // page_size is the maximum number of recommendations to return (default: 50, max: 1000)
  int32 page_size = 3;
  // page_token is the continuation token from a previous response
  string page_token = 4;
  // excluded_recommendation_ids contains IDs of recommendations to exclude from results.
  // Use this to filter out recommendations that have been dismissed by users.
  // Plugins should not return recommendations matching these IDs.
  repeated string excluded_recommendation_ids = 5;

+ // target_resources specifies the resources to analyze for recommendations.
+ // When provided, plugins return recommendations ONLY for these resources.
+ // When empty, plugins return recommendations for all resources in scope.
+ //
+ // Use cases:
+ //   - Stack-scoped recommendations: Pass Pulumi stack resources for targeted analysis
+ //   - Pre-deployment optimization: Analyze proposed resources before creation
+ //   - Batch resource analysis: Query recommendations for a known resource list
+ //
+ // Interaction with filter:
+ //   - target_resources defines the SCOPE (which resources to analyze)
+ //   - filter defines SELECTION CRITERIA within that scope (category, priority, etc.)
+ //   - Both are applied (AND logic): recommendations must match a target resource
+ //     AND satisfy any filter criteria
+ //
+ // Matching rules:
+ //   - provider and resource_type must always match (required fields)
+ //   - sku, region, and tags are matched only when specified in the target
+ //   - If specified, optional fields must match exactly (strict matching)
+ //
+ // Validation:
+ //   - Maximum 100 resources per request (exceeding returns InvalidArgument)
+ //   - Each ResourceDescriptor must have valid provider and resource_type
+ //   - Empty target_resources is valid (analyze all resources in scope)
+ repeated ResourceDescriptor target_resources = 6;
}
```

## Wire Format Compatibility

| Aspect | Status | Details |
|--------|--------|---------|
| Field number | Safe | 6 is unused, sequential |
| Wire type | Compatible | LENGTH_DELIMITED (repeated message) |
| Default value | Safe | Empty repeated = existing behavior |
| Old client → new server | Compatible | Server ignores missing field |
| New client → old server | Compatible | Server ignores unknown field |

## SDK Impact

### Generated Go Code

After `make generate`, `GetRecommendationsRequest` will have:

```go
type GetRecommendationsRequest struct {
    // ... existing fields ...
    TargetResources []*ResourceDescriptor `protobuf:"bytes,6,rep,name=target_resources"`
}

func (x *GetRecommendationsRequest) GetTargetResources() []*ResourceDescriptor {
    if x != nil {
        return x.TargetResources
    }
    return nil
}
```

### Validation Contract (sdk/go/testing/contract.go)

Add constant and validation function:

```go
const MaxTargetResources = 100

func ValidateGetRecommendationsRequest(req *pbc.GetRecommendationsRequest) error {
    // ... existing validation ...

    // Validate target_resources
    if len(req.GetTargetResources()) > MaxTargetResources {
        return NewContractError("target_resources", len(req.GetTargetResources()),
            fmt.Errorf("exceeds maximum of %d resources", MaxTargetResources))
    }

    for i, resource := range req.GetTargetResources() {
        if err := ValidateResourceDescriptor(resource); err != nil {
            return fmt.Errorf("target_resources[%d]: %w", i, err)
        }
    }

    return nil
}
```

## buf Breaking Change Check

This change will pass `buf breaking` because:

- New field with unused number (6)
- No field removals or type changes
- No renaming of existing fields
- Wire-compatible default (empty repeated)
