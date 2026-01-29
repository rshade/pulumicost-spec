# Data Model: Usage Profile Context

**Feature**: 042-usage-profile-context
**Date**: 2026-01-27
**Status**: Complete

## Entity Definitions

### UsageProfile Enum

**Purpose**: Represents the intended workload context for cost estimation and recommendations.

| Value | Proto Number | Description |
|-------|--------------|-------------|
| `USAGE_PROFILE_UNSPECIFIED` | 0 | Default/no preference - plugin applies standard behavior |
| `USAGE_PROFILE_PROD` | 1 | Production workloads - full utilization assumptions (730hr/month) |
| `USAGE_PROFILE_DEV` | 2 | Development workloads - reduced utilization (160hr/month typical) |
| `USAGE_PROFILE_BURST` | 3 | Burst/temporary workloads - high-intensity, scale-out assumptions |

**Proto Definition**:

```protobuf
// UsageProfile represents the intended workload context for cost estimation.
// Plugins use this to apply profile-appropriate defaults to cost calculations
// and recommendations.
//
// Usage:
//   - Core sets usage_profile in requests based on CLI flags (--profile=dev)
//   - Plugins apply profile-specific behavior (e.g., reduced utilization for DEV)
//   - Unknown values are treated as UNSPECIFIED for forward compatibility
enum UsageProfile {
  // Default value - no preference stated.
  // Plugins apply their standard estimation behavior.
  // Backward compatible with requests that don't specify a profile.
  USAGE_PROFILE_UNSPECIFIED = 0;

  // Production workload intent.
  // Plugins should assume:
  //   - Full-time utilization (730 hours/month for compute)
  //   - Production-grade instance types
  //   - High availability and redundancy considerations
  //   - Longer retention periods for storage
  USAGE_PROFILE_PROD = 1;

  // Development workload intent.
  // Plugins should assume:
  //   - Reduced utilization (e.g., 160 hours/month for business hours)
  //   - Burstable/cost-efficient instance types preferred
  //   - Minimal redundancy requirements
  //   - Shorter retention periods acceptable
  USAGE_PROFILE_DEV = 2;

  // Burst/temporary workload intent.
  // Plugins should assume:
  //   - High-intensity, short-duration usage
  //   - Scale-out architectures
  //   - Elevated data transfer/network costs
  //   - Batch processing or load testing scenarios
  USAGE_PROFILE_BURST = 3;
}
```

**Validation Rules**:

- Proto3 default is 0 (`UNSPECIFIED`) when field is not set
- Unknown values (e.g., from future spec versions) MUST be treated as `UNSPECIFIED`
- Plugins SHOULD log warning when encountering unknown values

### Extended Messages

#### GetProjectedCostRequest (Extended)

**New Field**: `usage_profile` (field number 6)

```protobuf
message GetProjectedCostRequest {
  ResourceDescriptor resource = 1;
  double utilization_percentage = 2;
  GrowthType growth_type = 3;
  optional double growth_rate = 4;
  bool dry_run = 5;

  // usage_profile signals the intended workload context.
  // Plugins use this to apply profile-appropriate defaults.
  // Examples:
  //   - DEV: Assume 160 hours/month, prefer burstable instances
  //   - PROD: Assume 730 hours/month, use production instance types
  //   - BURST: Assume high data transfer, scale-out patterns
  //
  // When UNSPECIFIED (default), plugins apply their standard behavior.
  // Unknown values are treated as UNSPECIFIED for forward compatibility.
  UsageProfile usage_profile = 6;
}
```

#### GetRecommendationsRequest (Extended)

**New Field**: `usage_profile` (field number 7)

```protobuf
message GetRecommendationsRequest {
  RecommendationFilter filter = 1;
  string projection_period = 2;
  int32 page_size = 3;
  string page_token = 4;
  repeated string excluded_recommendation_ids = 5;
  repeated ResourceDescriptor target_resources = 6;

  // usage_profile provides context for recommendation generation.
  // Plugins may adjust recommendation priorities based on profile:
  //   - DEV: Prioritize cost savings over performance
  //   - PROD: Balance reliability with cost optimization
  //   - BURST: Focus on scale-out and resource efficiency
  //
  // When UNSPECIFIED (default), plugins use their standard prioritization.
  // Unknown values are treated as UNSPECIFIED for forward compatibility.
  UsageProfile usage_profile = 7;
}
```

## Relationships

```text
┌─────────────────────────────────────┐
│        UsageProfile (enum)          │
│  UNSPECIFIED | PROD | DEV | BURST   │
└─────────────────────────────────────┘
           │
           │ included in
           ▼
┌─────────────────────────────────────┐     ┌─────────────────────────────────────┐
│    GetProjectedCostRequest          │     │    GetRecommendationsRequest        │
│  ─────────────────────────────────  │     │  ─────────────────────────────────  │
│  resource: ResourceDescriptor       │     │  filter: RecommendationFilter       │
│  utilization_percentage: double     │     │  projection_period: string          │
│  growth_type: GrowthType            │     │  page_size: int32                   │
│  growth_rate: optional double       │     │  page_token: string                 │
│  dry_run: bool                      │     │  excluded_recommendation_ids: []    │
│  usage_profile: UsageProfile ←NEW   │     │  target_resources: []               │
└─────────────────────────────────────┘     │  usage_profile: UsageProfile ←NEW   │
                                            └─────────────────────────────────────┘
```

## State Transitions

Not applicable - `UsageProfile` is a stateless enum representing intent at request time.
There are no state transitions; each request independently specifies its profile.

## SDK Type Mappings

### Go SDK

```go
// Generated proto type
type UsageProfile int32

const (
    UsageProfile_USAGE_PROFILE_UNSPECIFIED UsageProfile = 0
    UsageProfile_USAGE_PROFILE_PROD        UsageProfile = 1
    UsageProfile_USAGE_PROFILE_DEV         UsageProfile = 2
    UsageProfile_USAGE_PROFILE_BURST       UsageProfile = 3
)

// Helper functions (sdk/go/pluginsdk/usage_profile.go)
func IsValidUsageProfile(profile pbc.UsageProfile) bool
func ParseUsageProfile(s string) (pbc.UsageProfile, error)
func UsageProfileString(profile pbc.UsageProfile) string
func NormalizeUsageProfile(profile pbc.UsageProfile) pbc.UsageProfile // Returns UNSPECIFIED for unknown
```

### TypeScript SDK

```typescript
// Generated from proto
enum UsageProfile {
    USAGE_PROFILE_UNSPECIFIED = 0,
    USAGE_PROFILE_PROD = 1,
    USAGE_PROFILE_DEV = 2,
    USAGE_PROFILE_BURST = 3,
}

// Helper functions (sdk/typescript/packages/client/src/usage-profile.ts)
function isValidUsageProfile(profile: UsageProfile): boolean;
function parseUsageProfile(s: string): UsageProfile;
function usageProfileString(profile: UsageProfile): string;
```

## Backward Compatibility

| Scenario | Behavior |
|----------|----------|
| Old plugin, new request with profile | Plugin ignores unknown field (proto3 semantics) |
| New plugin, old request without profile | Field defaults to UNSPECIFIED (0), plugin uses standard behavior |
| New plugin, unknown profile value | Plugin treats as UNSPECIFIED, logs warning |

## Performance Considerations

- Enum validation: O(1) time, 0 allocations (package-level slice pattern)
- Profile field: 1-byte wire encoding (field numbers 1-15)
- No database impact (stateless)
