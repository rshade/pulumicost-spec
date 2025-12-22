# Feature Specification: Target Resources for Recommendations

**Feature Branch**: `019-target-resources`
**Created**: 2025-12-17
**Status**: Draft
**Input**: Add target_resources field to GetRecommendationsRequest for resource-scoped recommendations

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Stack-Scoped Recommendations (Priority: P1)

As a Pulumi user, I want to get cost optimization recommendations only for the resources in my
deployed stack, so that I can focus on actionable recommendations for my infrastructure without
noise from unrelated resources.

**Why this priority**: This is the primary use case that enables targeted cost analysis for
Pulumi stacks. Without this, users receive recommendations for all resources in their cloud
account, making it difficult to identify stack-specific optimizations.

**Independent Test**: Can be fully tested by passing a list of resource descriptors to the
GetRecommendations call and verifying only matching recommendations are returned. Delivers
immediate value for stack-focused cost analysis.

**Acceptance Scenarios**:

1. **Given** a user has a Pulumi stack with 10 EC2 instances, **When** they call
   GetRecommendations with those 10 resources as target_resources, **Then** only
   recommendations affecting those specific instances are returned
2. **Given** a user specifies target_resources for AWS resources, **When** recommendations
   exist for both targeted and non-targeted resources, **Then** only recommendations for
   targeted resources appear in the response
3. **Given** a user provides an empty target_resources list, **When** they call
   GetRecommendations, **Then** the system returns recommendations for all resources
   (backward compatible behavior)

---

### User Story 2 - Pre-Deployment Cost Optimization (Priority: P2)

As a DevOps engineer, I want to analyze proposed resource configurations before deployment,
so that I can optimize costs proactively and avoid deploying suboptimal infrastructure.

**Why this priority**: Enables shift-left cost optimization by analyzing resources during
`pulumi preview` phase, but requires the core targeting capability from P1 first.

**Independent Test**: Can be tested by passing proposed resource configurations (provider,
type, SKU, region) and verifying relevant recommendations are returned even though resources
don't exist yet.

**Acceptance Scenarios**:

1. **Given** a DevOps engineer has a proposed configuration for 5 new m5.xlarge instances,
   **When** they call GetRecommendations with these as target_resources, **Then**
   recommendations for rightsizing or alternative SKUs are returned
2. **Given** proposed resources span multiple providers (AWS and Azure), **When**
   target_resources includes both, **Then** recommendations from both providers are
   returned appropriately

---

### User Story 3 - Batch Resource Analysis (Priority: P3)

As a FinOps analyst, I want to query recommendations for a specific list of resource IDs
across multiple accounts, so that I can perform targeted cost audits on high-value resources.

**Why this priority**: Supports advanced enterprise workflows but builds on the same
targeting mechanism from P1/P2.

**Independent Test**: Can be tested by providing resource descriptors with specific tags or
identifiers and verifying batch filtering works correctly.

**Acceptance Scenarios**:

1. **Given** a FinOps analyst has a list of 50 high-cost resources from a cost report,
   **When** they call GetRecommendations with those resources, **Then** only
   recommendations for those specific resources are returned
2. **Given** target_resources combined with a filter for minimum savings threshold,
   **When** recommendations are retrieved, **Then** results match BOTH the target scope
   AND the filter criteria (AND logic)

---

### Edge Cases

- What happens when target_resources exceeds the maximum allowed limit (100 resources)?
  - System returns an InvalidArgument error with a clear message about the limit
- What happens when a target resource has invalid or missing required fields?
  - System validates each ResourceDescriptor and returns InvalidArgument for the first
    invalid entry
- How does the system handle when none of the target resources have recommendations?
  - System returns an empty recommendations list with successful response
- What happens when target_resources contains duplicate entries?
  - System processes duplicates without error; recommendations appear once regardless
    of duplicates

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST accept a list of resource descriptors (target_resources) in
  recommendation requests
- **FR-002**: System MUST limit target_resources to a maximum of 100 entries per request
- **FR-003**: System MUST return InvalidArgument error when target_resources exceeds
  the maximum limit
- **FR-004**: System MUST validate each resource descriptor in target_resources has
  valid provider and resource_type
- **FR-005**: System MUST return only recommendations that match at least one target
  resource when target_resources is provided
- **FR-006**: System MUST preserve existing behavior (return all recommendations) when
  target_resources is empty or not provided
- **FR-007**: System MUST apply both target_resources scope AND filter criteria using
  AND logic
- **FR-008**: System MUST match recommendations to target resources using strict equality:
  provider and resource_type must always match; if SKU, region, or tags are specified
  in the target, they must also match exactly
- **FR-009**: System MUST maintain full backward compatibility with existing clients
  that don't use target_resources

### Key Entities

- **ResourceDescriptor**: Describes a cloud resource with provider, resource_type, SKU,
  region, and tags. Used to specify which resources to analyze.
- **GetRecommendationsRequest**: The request message that will include the new
  target_resources field alongside existing filter, pagination, and exclusion parameters.
- **Recommendation**: A cost optimization suggestion that includes the affected resource
  information, used for matching against target_resources.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Users can retrieve recommendations for a specific list of up to 100
  resources in a single request
- **SC-002**: Requests with target_resources return only matching recommendations with
  100% precision (no false positives)
- **SC-003**: Empty target_resources maintains identical behavior to current
  implementation (100% backward compatibility)
- **SC-004**: Invalid target_resources requests fail fast with clear error messages
  in under 100ms
- **SC-005**: Plugin developers can implement target_resources filtering with provided
  SDK helpers

## Clarifications

### Session 2025-12-17

- Q: How should optional fields (SKU/region/tags) affect matching when specified in target?
  → A: Strict match - if target specifies optional fields, they must match exactly

## Assumptions

- Resource matching uses strict equality: provider and resource_type are always required
  to match; if SKU, region, or tags are specified in the target, they must also match
  exactly (unspecified optional fields are not checked)
- The 100-resource limit is sufficient for typical Pulumi stacks (most have 10-50 resources)
- Plugins may implement target_resources filtering at the backend level or apply
  filtering client-side
- Tags use exact key-value matching: all tags specified in the target must be present
  on the recommendation's resource with identical values (consistent with FR-008)
