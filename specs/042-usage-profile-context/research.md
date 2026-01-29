# Research: Usage Profile Context

**Feature**: 042-usage-profile-context
**Date**: 2026-01-27
**Status**: Complete

## Research Questions

### 1. Where should UsageProfile enum be defined?

**Decision**: `proto/finfocus/v1/enums.proto`

**Rationale**: The existing `enums.proto` contains all shared enums used across messages:

- `GrowthType` - mathematical models for cost projections
- `RecommendationCategory` - recommendation classifications
- `FieldSupportStatus` - dry-run field support levels
- `PluginCapability` - feature capabilities

`UsageProfile` follows the same pattern: a shared enum used in multiple request messages.

**Alternatives Considered**:

- `costsource.proto`: Rejected - would clutter the main service file; enums belong in `enums.proto`
- New `usage.proto`: Rejected - overkill for a single enum; follows existing consolidation pattern

### 2. Which messages should include usage_profile?

**Decision**: `GetProjectedCostRequest` and `GetRecommendationsRequest` only

**Rationale**: Per spec clarification:

- `GetProjectedCostRequest`: Cost projections benefit from workload intent (DEV=160hr, PROD=730hr)
- `GetRecommendationsRequest`: Recommendations should be context-aware (DEV=cost-saving, PROD=reliability)
- `GetActualCostRequest`: Excluded - actual costs reflect what happened, not intent

**Alternatives Considered**:

- Include in all requests: Rejected - actual costs are historical, profile doesn't apply
- Include in `EstimateCostRequest`: Deferred - could be added later if needed; current scope focuses on projections

### 3. What enum values should UsageProfile have?

**Decision**: Four values following proto3 conventions:

```protobuf
enum UsageProfile {
  USAGE_PROFILE_UNSPECIFIED = 0;  // Default, plugin applies standard behavior
  USAGE_PROFILE_PROD = 1;         // Production workloads (full utilization)
  USAGE_PROFILE_DEV = 2;          // Development workloads (reduced utilization)
  USAGE_PROFILE_BURST = 3;        // Burst/temporary workloads (scale-out)
}
```

**Rationale**:

- `UNSPECIFIED = 0` is proto3 requirement and provides backward compatibility
- Ordering (PROD=1, DEV=2) reflects priority - production is the most common enterprise use case
- `BURST` added per user story 3 for batch/load-testing scenarios

**Alternatives Considered**:

- String-based profile: Rejected - enums provide type safety and code generation benefits
- More granular profiles (STAGING, TEST, etc.): Deferred - can be added later without breaking changes

### 4. How should plugins handle unknown/future profile values?

**Decision**: Treat unknown as `UNSPECIFIED`, log warning at INFO level

**Rationale**: Per spec requirement FR-009:

- Enables forward compatibility when new profiles are added
- Graceful degradation prevents plugin failures
- Logging provides visibility for debugging

**Implementation Pattern**:

```go
func handleUsageProfile(profile pbc.UsageProfile) {
    switch profile {
    case pbc.UsageProfile_USAGE_PROFILE_PROD:
        // Apply production defaults
    case pbc.UsageProfile_USAGE_PROFILE_DEV:
        // Apply development defaults
    case pbc.UsageProfile_USAGE_PROFILE_BURST:
        // Apply burst defaults
    case pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED:
        // Apply plugin default behavior
    default:
        // Unknown value - treat as UNSPECIFIED
        log.Warn().Int32("usage_profile", int32(profile)).Msg("Unknown usage profile, treating as UNSPECIFIED")
    }
}
```

### 5. What SDK helpers should be provided?

**Decision**: Profile-aware builder methods on existing builders

**Rationale**: Per spec requirement FR-010 and existing SDK patterns:

- `FocusRecordBuilder` has `WithProfileDefaults(profile)` method
- Follow functional options pattern used throughout SDK
- Zero-allocation implementation for performance

**SDK Helpers to Implement**:

```go
// Profile detection and validation
func IsValidUsageProfile(profile pbc.UsageProfile) bool
func ParseUsageProfile(s string) (pbc.UsageProfile, error)
func UsageProfileString(profile pbc.UsageProfile) string

// Builder integration (on FocusRecordBuilder or similar)
func (b *FocusRecordBuilder) WithProfileDefaults(profile pbc.UsageProfile) *FocusRecordBuilder
```

**Alternatives Considered**:

- New `UsageProfileBuilder`: Rejected - profile is a single field, not complex enough for builder
- Extension methods on proto: Rejected - Go doesn't support extension methods

### 6. How should logging work for profile-specific behavior?

**Decision**: Structured logging with `usage_profile` field at INFO level

**Rationale**: Per spec requirement FR-008:

- INFO level ensures visibility without cluttering debug logs
- Structured field enables log aggregation and filtering
- Consistent with existing zerolog patterns in SDK

**Logging Pattern**:

```go
log.Info().
    Str("usage_profile", "DEV").
    Str("resource_type", "ec2").
    Msg("Applying development profile defaults")
```

### 7. Should TypeScript SDK be updated?

**Decision**: Yes, required by Constitution XIII (Multi-Language SDK Synchronization)

**Rationale**:

- TypeScript SDK exists at `sdk/typescript/`
- Proto regeneration via buf will update generated code
- Client wrappers may need manual updates for profile handling

**Implementation Steps**:

1. Regenerate proto bindings: `npm run generate` in `sdk/typescript/`
2. Update client wrappers if they expose profile-related methods
3. Add TypeScript tests for profile handling

### 8. What field number should usage_profile use?

**Decision**: Next available field number in each message

**Rationale**:

- `GetProjectedCostRequest`: Field 6 (after `dry_run` at 5)
- `GetRecommendationsRequest`: Field 7 (after `target_resources` at 6)

Proto3 reserves field numbers 1-15 for frequently used fields (1-byte encoding).
Both messages have space in the 1-15 range, but usage_profile is optional,
so higher numbers are acceptable.

### 9. Performance requirements for SDK helpers?

**Decision**: <15 ns/op, 0 allocs/op for validation functions

**Rationale**: Per Constitution VIII (Performance is Paramount) and existing patterns:

- Currency validation: <15 ns/op, 0 allocs/op
- Registry validation: 5-12 ns/op, 0 allocs/op
- UsageProfile validation should match these benchmarks

**Implementation Pattern**:

```go
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allUsageProfiles = []pbc.UsageProfile{
    pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
    pbc.UsageProfile_USAGE_PROFILE_PROD,
    pbc.UsageProfile_USAGE_PROFILE_DEV,
    pbc.UsageProfile_USAGE_PROFILE_BURST,
}

func IsValidUsageProfile(profile pbc.UsageProfile) bool {
    for _, valid := range allUsageProfiles {
        if profile == valid {
            return true
        }
    }
    return false
}
```

## Technology Decisions Summary

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Enum location | `enums.proto` | Follows existing enum consolidation pattern |
| Messages extended | `GetProjectedCostRequest`, `GetRecommendationsRequest` | Per spec clarification - actual costs excluded |
| Default behavior | UNSPECIFIED=0 | Proto3 convention, backward compatible |
| Unknown handling | Treat as UNSPECIFIED + warn | Forward compatibility |
| SDK helpers | Validation + builder methods | Matches existing SDK patterns |
| Logging | INFO + structured fields | Per FR-008 requirement |
| TypeScript SDK | Required update | Constitution XIII |
| Performance | <15 ns/op, 0 allocs | Constitution VIII |

## Dependencies

- buf v1.32.1 (proto generation)
- google.golang.org/protobuf (Go SDK)
- zerolog (structured logging)
- TypeScript SDK regeneration

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Plugin compatibility | UNSPECIFIED default preserves existing behavior |
| Performance regression | Benchmark tests verify <15 ns/op |
| SDK version skew | TypeScript SDK updated in same PR |
| Profile misuse | Documentation clarifies plugin responsibility |
