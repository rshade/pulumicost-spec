# Feature Specification: Resource ID and ARN Fields for ResourceDescriptor

**Feature Branch**: `028-resource-id`
**Created**: 2025-12-26
**Status**: Draft
**Input**: User description: "Add resource ID and ARN fields to ResourceDescriptor
for recommendation correlation and exact resource matching"

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Batch Resource Recommendation Correlation (Priority: P1)

A cost management system sends multiple resources in a single GetRecommendations
request. When recommendations are returned, the system can match each
recommendation back to the original resource in the request using the resource
ID.

**Why this priority**: This is the core problem being solved. Without this
capability, batch recommendation requests cannot properly correlate responses
to requests, breaking multi-resource analysis workflows.

**Independent Test**: Can be fully tested by sending a batch request with 3
resources (each with unique IDs), receiving recommendations, and verifying that
each recommendation's resource can be traced back to the correct input resource.

**Acceptance Scenarios**:

1. **Given** a GetRecommendationsRequest with 3 ResourceDescriptors
   (IDs: "res-001", "res-002", "res-003"),
   **When** the plugin returns recommendations,
   **Then** recommendations for each resource include the corresponding ID
   enabling correlation.

2. **Given** a ResourceDescriptor with ID "my-resource-id",
   **When** the plugin processes the request and returns a recommendation,
   **Then** the recommendation response includes "my-resource-id" for matching.

3. **Given** a GetRecommendationsRequest with a ResourceDescriptor that has an
   empty ID,
   **When** the plugin returns recommendations,
   **Then** the response is valid (ID is optional for backward compatibility).

---

### User Story 2 - Exact Resource Matching via ARN (Priority: P1)

A plugin receives a ResourceDescriptor with a canonical cloud resource identifier
(ARN). The plugin uses this identifier to look up the exact resource in the
cloud provider's API or cost data, rather than relying on fuzzy matching by
type/sku/region/tags.

**Why this priority**: Exact matching eliminates ambiguity when multiple
resources share the same type, SKU, and region. Critical for accurate cost
attribution and recommendations.

**Independent Test**: Can be tested by sending a request with an ARN for a
known resource and verifying the plugin returns data for that exact resource,
not a similar one.

**Acceptance Scenarios**:

1. **Given** a ResourceDescriptor with ARN
   `arn:aws:ec2:us-east-1:123456789012:instance/i-abc123`,
   **When** the AWS plugin processes the request,
   **Then** the plugin queries cost data for instance `i-abc123` specifically.

2. **Given** a ResourceDescriptor with Azure Resource ID
   `/subscriptions/sub-1/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1`,
   **When** the Azure plugin processes the request,
   **Then** the plugin queries cost data for `vm-1` specifically.

3. **Given** a ResourceDescriptor with both ARN and type/sku/region fields,
   **When** the plugin processes the request,
   **Then** the plugin prefers the ARN for exact matching (ARN takes precedence).

4. **Given** a ResourceDescriptor with an empty ARN,
   **When** the plugin processes the request,
   **Then** the plugin falls back to matching by type/sku/region/tags.

---

### User Story 3 - Pass-Through Identifier Support (Priority: P2)

Plugin developers can receive the resource ID from requests and include it in
responses without needing to understand or validate the ID format. The ID is
treated as an opaque pass-through value.

**Why this priority**: Enables plugin implementation flexibility while ensuring
consistent behavior across all plugins.

**Independent Test**: Can be tested by implementing a mock plugin that receives
requests with various ID formats and verifying it can pass them through
unchanged.

**Acceptance Scenarios**:

1. **Given** a ResourceDescriptor with ID containing special characters
   (e.g., "urn:pulumi:stack::project::aws:ec2/instance:Instance::webserver"),
   **When** the plugin processes the request,
   **Then** the ID is preserved exactly as provided.

2. **Given** a ResourceDescriptor with a very long ID (256+ characters),
   **When** the plugin processes the request,
   **Then** the ID is preserved without truncation.

---

### User Story 4 - Backward Compatible Protocol Evolution (Priority: P2)

Existing plugins and clients continue to work without modification when the new
fields are added. Both fields are optional and have empty defaults.

**Why this priority**: Ensures smooth adoption without breaking existing
deployments.

**Independent Test**: Can be tested by running existing plugin implementations
against updated proto definitions and verifying all RPCs continue to function.

**Acceptance Scenarios**:

1. **Given** an existing plugin compiled against the old proto (without new
   fields),
   **When** a new client sends a request with id and arn fields set,
   **Then** the plugin ignores the unknown fields and processes the request
   normally.

2. **Given** a new plugin compiled with the new fields,
   **When** an old client sends a request without id/arn fields,
   **Then** both fields default to empty string and the request is processed
   normally using type/sku/region/tags matching.

---

### Edge Cases

- What happens when duplicate IDs are provided in the same request?
  (Handled by caller - no server-side validation required)
- How does the system handle empty string vs. unset ID/ARN?
  (Both treated as "not provided")
- What if ARN format is invalid for the provider?
  (Plugin logs warning, falls back to type/sku/region matching)
- What if ARN doesn't match any known resource?
  (Plugin returns empty results or NotFound, depending on RPC)
- What is the maximum reasonable ID/ARN length?
  (No hard limit in proto; AWS ARNs can be ~2048 chars)

## Requirements _(mandatory)_

### Functional Requirements

#### ID Field (Correlation)

- **FR-001**: Protocol MUST include an `id` field in the `ResourceDescriptor`
  message for request/response correlation.
- **FR-002**: The `id` field MUST be optional with an empty string default for
  backward compatibility.
- **FR-003**: The `id` field MUST be passed through unchanged by plugins
  (no validation or transformation required).
- **FR-004**: Plugins MUST include the `id` value in any responses or
  recommendations related to the resource.

#### ARN Field (Exact Matching)

- **FR-005**: Protocol MUST include an `arn` field in the `ResourceDescriptor`
  message for canonical cloud resource identification.
- **FR-006**: The `arn` field MUST be optional with an empty string default for
  backward compatibility.
- **FR-007**: When `arn` is provided, plugins SHOULD use it for exact resource
  matching instead of type/sku/region/tags.
- **FR-008**: Plugins MUST gracefully handle invalid or unrecognized ARN formats
  by falling back to type/sku/region/tags matching.

#### Documentation and SDK

- **FR-009**: Protocol documentation MUST clearly describe both fields with
  usage guidance and provider-specific examples.
- **FR-010**: Go SDK MUST regenerate bindings to include the new fields.
- **FR-011**: pluginsdk MUST provide helper functions for working with both
  fields.

### Key Entities

- **ResourceDescriptor**: Extended with two new fields:
  - `id` (field 7): Client correlation identifier (opaque pass-through)
  - `arn` (field 8): Canonical cloud resource identifier (for exact matching)
- **Recommendation**: Existing message - plugins include the `id` for
  correlation in returned recommendations.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Batch recommendation requests with N resources can be correlated
  with N recommendations in O(1) lookup time using the ID field.
- **SC-002**: Plugins can perform exact resource matching using ARN when
  provided, eliminating ambiguous matches.
- **SC-003**: 100% of existing plugins and clients continue to function without
  modification after the protocol update.
- **SC-004**: Plugin developers can implement pass-through ID correlation with
  zero additional validation code.
- **SC-005**: Protocol change follows semantic versioning (minor version bump
  for backward-compatible addition).

## Technical Context _(informational)_

### Background

The `GetRecommendations` RPC accepts a `target_resources` field containing
multiple `ResourceDescriptor` messages. Currently, when recommendations are
returned, there is no standardized way to match a recommendation back to its
corresponding input resource in batch operations.

Additionally, `GetActualCostRequest` has an `arn` field (field 5), but
`ResourceDescriptor` does not, creating inconsistency for exact resource
identification.

**Current Workarounds**:

1. pulumicost-core sets the `ResourceID` on returned recommendations when
   there's exactly one resource in the request. This fails for batch requests.
2. Plugins rely on type/sku/region/tags matching, which can be ambiguous when
   multiple similar resources exist.

### Protocol Impact

The change adds fields 7 and 8 to `ResourceDescriptor`:

```protobuf
message ResourceDescriptor {
  // ... existing fields 1-6 ...

  // id is a client-specified identifier for request/response correlation.
  // OPTIONAL. When provided, plugins MUST include this ID in any
  // recommendations or responses related to this resource, enabling
  // clients to match responses to their original requests in batch operations.
  //
  // The ID is treated as an opaque string - plugins MUST NOT validate or
  // transform this value. Common formats include Pulumi URNs, UUIDs, or
  // application-specific identifiers.
  //
  // Example: "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::webserver"
  string id = 7;

  // arn is the canonical cloud resource identifier for exact matching.
  // OPTIONAL. When provided, plugins SHOULD use this for precise resource
  // lookup instead of matching by type/sku/region/tags.
  //
  // This field uses "arn" as the name for consistency with GetActualCostRequest,
  // but accepts canonical identifiers from any cloud provider:
  //
  // AWS ARN:
  //   arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0
  //
  // Azure Resource ID:
  //   /subscriptions/{sub-id}/resourceGroups/{rg}/providers/
  //   Microsoft.Compute/virtualMachines/{vm-name}
  //
  // GCP Full Resource Name:
  //   //compute.googleapis.com/projects/{project}/zones/{zone}/instances/{name}
  //
  // Kubernetes Resource:
  //   {cluster}/{namespace}/{kind}/{name} or UID
  //
  // Cloudflare:
  //   {zone-id}/{resource-type}/{resource-id}
  //
  // If the ARN format is unrecognized or the resource is not found, plugins
  // SHOULD fall back to type/sku/region/tags matching and MAY log a warning.
  string arn = 8;
}
```

### Field Comparison

| Field | Purpose | Set By | Validated By Plugin | Used For |
|-------|---------|--------|---------------------|----------|
| `id` | Correlation | Client | No (pass-through) | Matching responses to requests |
| `arn` | Exact lookup | Client | Yes (format check) | Precise resource identification |

### Cross-Repository Impact

This section provides implementation guidance for dependent repositories.

#### pulumicost-core

**Location**: Recommendation analyzer and cost query components

**Implementation Requirements**:

1. **Populate `id` field**: When calling `GetRecommendations`, set the `id`
   field to the Pulumi resource URN or internal tracking ID:

   ```go
   descriptor := &pb.ResourceDescriptor{
       Provider:     "aws",
       ResourceType: "ec2",
       Sku:          "t3.micro",
       Region:       "us-east-1",
       Id:           resource.URN, // Pulumi URN for correlation
       Arn:          resource.Outputs["arn"], // AWS ARN if available
   }
   ```

2. **Populate `arn` field**: When the cloud resource ARN/ID is known (from
   Pulumi outputs or state), include it for exact matching:
   - AWS: Use `arn` output from resource
   - Azure: Use `id` output (Azure Resource ID format)
   - GCP: Use `selfLink` or construct full resource name
   - Kubernetes: Use `metadata.uid` or `namespace/name`

3. **Correlate responses**: When processing `GetRecommendationsResponse`, use
   the returned `id` to match recommendations back to input resources:

   ```go
   // Build lookup map from request
   resourceByID := make(map[string]*Resource)
   for _, r := range request.TargetResources {
       resourceByID[r.Id] = originalResources[r.Id]
   }

   // Match recommendations to resources
   for _, rec := range response.Recommendations {
       if resource, ok := resourceByID[rec.Resource.Id]; ok {
           // Recommendation correlates to this resource
       }
   }
   ```

#### pulumicost-plugin-aws-public

**Location**: Recommendations provider implementation

**Implementation Requirements**:

1. **Pass through `id` field**: Copy the `id` from input ResourceDescriptor to
   output ResourceRecommendationInfo:

   ```go
   func (p *Plugin) GetRecommendations(
       ctx context.Context,
       req *pb.GetRecommendationsRequest,
   ) (*pb.GetRecommendationsResponse, error) {
       // For each target resource
       for _, target := range req.TargetResources {
           // ... get recommendations from AWS Cost Explorer ...

           // Include the client's ID in the response for correlation
           rec := &pb.Recommendation{
               Resource: &pb.ResourceRecommendationInfo{
                   // ... other fields ...
               },
           }
           // Preserve correlation ID
           // Note: ResourceRecommendationInfo may need id field added
       }
   }
   ```

2. **Use `arn` for exact matching**: When ARN is provided, use it to filter
   AWS Cost Explorer recommendations:

   ```go
   func (p *Plugin) matchResource(target *pb.ResourceDescriptor) string {
       if target.Arn != "" {
           // Parse ARN to extract resource ID
           arn, err := arn.Parse(target.Arn)
           if err == nil {
               // Use ARN for exact Cost Explorer filtering
               return extractResourceID(arn)
           }
           // Log warning and fall back to fuzzy matching
           p.logger.Warn().
               Str("arn", target.Arn).
               Msg("invalid ARN format, falling back to type/region matching")
       }
       // Fall back to type/sku/region/tags matching
       return ""
   }
   ```

3. **Handle missing ARN gracefully**: Not all resources have ARNs at request
   time (e.g., during preview). Fall back to existing matching logic.

#### pluginsdk Updates

**Location**: `sdk/go/pluginsdk/`

**New Helper Functions**:

1. **ResourceDescriptor builder methods**:

   ```go
   // WithID sets the correlation ID for request/response matching.
   // The ID is passed through unchanged to responses.
   func (b *ResourceDescriptorBuilder) WithID(id string) *ResourceDescriptorBuilder

   // WithARN sets the canonical cloud resource identifier for exact matching.
   // Supported formats: AWS ARN, Azure Resource ID, GCP Full Resource Name,
   // Kubernetes UID, Cloudflare zone/resource ID.
   func (b *ResourceDescriptorBuilder) WithARN(arn string) *ResourceDescriptorBuilder
   ```

2. **ARN parsing helpers** (optional, for plugin convenience):

   ```go
   // ParseARN attempts to parse a canonical resource identifier and returns
   // the provider type and resource components. Returns an error if the
   // format is not recognized.
   func ParseARN(arn string) (*ParsedARN, error)

   // ParsedARN contains the parsed components of a canonical resource ID.
   type ParsedARN struct {
       Provider     string // "aws", "azure", "gcp", "kubernetes", "cloudflare"
       ResourceType string // Provider-specific resource type
       ResourceID   string // The unique resource identifier
       Region       string // Region/location if available
       Account      string // Account/subscription/project if available
   }
   ```

3. **Documentation**: Add comprehensive godoc with examples for each provider
   format.

## Assumptions

- Field numbers 7 and 8 are available in ResourceDescriptor
  (verified: field 6 is utilization_percentage)
- The `id` is a string to accommodate various identifier formats
  (URNs, UUIDs, Pulumi resource names)
- The `arn` uses the name "arn" for consistency with `GetActualCostRequest.arn`
  but accepts any provider's canonical format
- No server-side ID uniqueness validation is required - callers are responsible
  for unique IDs if needed
- Empty string is semantically equivalent to "not provided" for both fields
- Plugins should prefer ARN when available but must handle cases where only
  type/sku/region/tags are provided
