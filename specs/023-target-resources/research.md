# Research: Target Resources for Recommendations

**Feature**: 019-target-resources
**Date**: 2025-12-17

## Research Summary

No critical unknowns requiring external research. All technical decisions are based on existing
proto patterns and SDK conventions established in this repository.

## Decision Log

### D1: Proto Field Number

**Decision**: Use field number 6 for `target_resources`
**Rationale**: Sequential numbering after existing fields 1-5 in `GetRecommendationsRequest`
**Alternatives Considered**:

- Higher field number (e.g., 10+): Rejected - breaks sequential convention
- Reserved field: Not applicable - no prior removals in this message

### D2: Reuse ResourceDescriptor

**Decision**: Reuse existing `ResourceDescriptor` message type for target resources
**Rationale**: Maintains consistency with other RPCs (Supports, GetProjectedCost, GetPricingSpec)
**Alternatives Considered**:

- New `TargetResource` message: Rejected - would duplicate fields with no benefit
- Inline fields: Rejected - less structured, harder to validate

### D3: Maximum Resource Limit

**Decision**: 100 resources maximum per request
**Rationale**:

- Typical Pulumi stacks contain 10-50 resources
- 100 provides headroom for large deployments
- Prevents unbounded memory/processing costs
- Consistent with other paginated limits in the proto

**Alternatives Considered**:

- 50 resources: Too restrictive for larger stacks
- 500 resources: Potentially expensive validation/matching
- No limit: Unbounded resource consumption risk

### D4: Matching Semantics

**Decision**: Strict equality matching when optional fields are specified
**Rationale**:

- Predictable behavior: specifying more fields = more precise filtering
- Consistent with existing `RecommendationFilter` field behavior
- Users control precision by choosing which fields to include

**Alternatives Considered**:

- Lenient matching (hints only): Less predictable, harder to test
- Configurable match mode: Over-engineered for initial implementation

### D5: AND Logic with Filter

**Decision**: target_resources AND filter applied together (scope then select)
**Rationale**:

- target_resources defines SCOPE (which resources to analyze)
- filter defines SELECTION (which recommendations to return within scope)
- Consistent with composition patterns in other APIs

**Alternatives Considered**:

- OR logic: Would expand results unpredictably
- target_resources replaces filter: Loses existing filter functionality

### D6: Validation Location

**Decision**: Add validation in `sdk/go/testing/contract.go`
**Rationale**:

- Consistent with existing validation patterns (MaxPageSize, etc.)
- Reusable by both Core and Plugin implementations
- Enables conformance testing

**Alternatives Considered**:

- Proto-level validation (protoc-gen-validate): Not currently used in project
- Per-plugin validation: Inconsistent enforcement

## Best Practices Applied

### Proto Design

- Add new field at end of message (field 6)
- Use `repeated` for list fields
- Reuse existing message types when appropriate
- Add comprehensive proto comments with use cases and validation rules

### SDK Validation

- Define constant for max limit (`MaxTargetResources = 100`)
- Validate each item in repeated field
- Return early on first validation failure
- Use existing `ValidateResourceDescriptor` function

### Testing

- Write integration tests for all acceptance scenarios
- Test edge cases (empty list, max limit, invalid items, duplicates)
- Test AND logic with filter combinations
- Add to conformance suite for plugin developers

## External References

- Existing ResourceDescriptor usage: `proto/pulumicost/v1/costsource.proto:181-218`
- Existing validation patterns: `sdk/go/testing/contract.go`
- Mock plugin patterns: `sdk/go/testing/mock_plugin.go`
- Constitution requirements: `.specify/memory/constitution.md`
