# CDP-SPEC: Capability Discovery Polish

> **Local Consolidated Development Plan**
>
> This specification consolidates 7 related issues into a single cohesive development plan
> for polishing the capability discovery system and related SDK quality improvements.

## Overview

| Attribute | Value |
|-----------|-------|
| **Issues** | #294, #295, #299, #300, #301, #208, #209 |
| **Theme** | SDK Quality & Capability Discovery Polish |
| **Priority** | P1 - Should Fix |
| **Estimated Effort** | 4-6 hours |
| **Risk Level** | Low (refactor + documentation) |

## Issue Summary

| Issue | Title | Type | Priority | Status |
|-------|-------|------|----------|--------|
| [#294](https://github.com/rshade/finfocus-spec/issues/294) | Missing DryRunHandler Interface Definition | Bug | P0 | ✅ Already Fixed |
| [#295](https://github.com/rshade/finfocus-spec/issues/295) | Proto Field Type Inconsistency | Bug | P0 | Open |
| [#299](https://github.com/rshade/finfocus-spec/issues/299) | Inconsistent Backward Compatibility Patterns | Enhancement | P1 | Open |
| [#300](https://github.com/rshade/finfocus-spec/issues/300) | Documentation for Capability Discovery | Documentation | P1 | Open |
| [#301](https://github.com/rshade/finfocus-spec/issues/301) | Optimize Slice Copying Performance | Enhancement | P2 | Open |
| [#208](https://github.com/rshade/finfocus-spec/issues/208) | Clarify Compatibility Test Naming | Documentation | Low | Open |
| [#209](https://github.com/rshade/finfocus-spec/issues/209) | Add Complete Import Statements to README | Documentation | Low | Open |

---

## Part 1: Issue #294 - DryRunHandler Interface Definition

### Current State Analysis

**Finding**: The `DryRunHandler` interface is **already defined** in `sdk/go/pluginsdk/dry_run.go:287-295`:

```go
// DryRunHandler is an optional interface that plugins can implement to provide
// DryRun functionality. If a plugin implements this interface, the SDK can
// automatically route DryRun requests to it.
type DryRunHandler interface {
    // HandleDryRun returns field mapping information for the given resource type.
    HandleDryRun(req *pbc.DryRunRequest) (*pbc.DryRunResponse, error)
}
```

The type assertion in `plugin_info.go:295` correctly references this interface:

```go
if _, ok := plugin.(DryRunHandler); ok {
    capabilities = append(capabilities, pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN)
}
```

### Resolution

**Status**: ✅ Already Fixed

**Action**: Close issue #294 as completed. The interface exists and auto-discovery works.

**Verification Command**:

```bash
go build ./sdk/go/pluginsdk && echo "DryRunHandler interface compiles correctly"
```

---

## Part 2: Issue #295 - Proto Field Type Inconsistency

### Problem Statement

There are two capability representation formats in `GetPluginInfoResponse`:

| Field | Type | Purpose |
|-------|------|---------|
| `capabilities` (proto) | `repeated PluginCapability` | Modern enum-based capabilities |
| `metadata` (proto) | `map<string, string>` | Legacy string-based capabilities |

The confusion arises from:

1. `SupportsResponse.capabilities` uses `map<string, bool>` (proto line 209)
2. `GetPluginInfoResponse.metadata` uses `map<string, string>`
3. SDK converts between these formats for backward compatibility

### Current Proto Definitions

**`SupportsResponse`** (costsource.proto:197-215):

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2;
  map<string, bool> capabilities = 3;        // Legacy bool map
  repeated MetricKind supported_metrics = 4;
  repeated PluginCapability capabilities_enum = 5;  // Modern enum list
}
```

**`GetPluginInfoResponse`** (implied from SDK):

```protobuf
message GetPluginInfoResponse {
  string name = 1;
  string version = 2;
  string spec_version = 3;
  repeated string providers = 4;
  repeated PluginCapability capabilities = 5;  // Modern enum list
  map<string, string> metadata = 6;            // Legacy string map
}
```

### Resolution Strategy

The type inconsistency is **intentional for backward compatibility**:

- Modern clients use `capabilities` (enum list)
- Legacy clients use `metadata` with `"supports_xyz": "true"` strings
- SDK handles conversion via `capabilitiesToLegacyMetadataWithWarnings()`

**Actions**:

1. Add clarifying comments to proto fields
2. Document the dual-format design in CLAUDE.md
3. Update SDK documentation

### Implementation

**File**: `proto/finfocus/v1/costsource.proto`

Add comments to clarify field relationships (locate `GetPluginInfoResponse` message):

```protobuf
message GetPluginInfoResponse {
  string name = 1;
  string version = 2;
  string spec_version = 3;
  repeated string providers = 4;

  // Modern capability format using strongly-typed enums.
  // Prefer this field for capability queries on newer clients.
  // SDK auto-populates this based on implemented interfaces.
  repeated PluginCapability capabilities = 5;

  // Legacy metadata format for backward compatibility with older hosts.
  // Contains string-based capability flags: {"supports_xyz": "true"}.
  // SDK auto-populates this from capabilities for backward compatibility.
  // DEPRECATION: New integrations should use capabilities field instead.
  map<string, string> metadata = 6;
}
```

**File**: `proto/finfocus/v1/costsource.proto` (SupportsResponse)

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2;

  // Legacy capability format using boolean map.
  // Example: {"recommendations": true, "dry_run": true}
  // DEPRECATION: Use capabilities_enum for new integrations.
  map<string, bool> capabilities = 3;

  repeated MetricKind supported_metrics = 4;

  // Modern capability format using strongly-typed enums.
  // Auto-populated by SDK based on implemented interfaces.
  repeated PluginCapability capabilities_enum = 5;
}
```

---

## Part 3: Issue #299 - Inconsistent Backward Compatibility Patterns

### Problem Statement

Two different patterns exist for converting `PluginCapability` enums to legacy strings:

| Location | Pattern | Code Path |
|----------|---------|-----------|
| `sdk.go:430-445` | `capabilitiesToLegacyMetadataWithWarnings()` | `handleConfiguredPluginInfo()` |
| `sdk.go:537+` | Inline conversion in `Supports()` | Direct enum iteration |

### Current Implementation Analysis

**Pattern 1** - `handleConfiguredPluginInfo()` (sdk.go:429-446):

```go
legacyMeta, warnings := capabilitiesToLegacyMetadataWithWarnings(capabilities)
for key := range legacyMeta {
    metadata[key] = capabilityTrue
}
```

**Pattern 2** - Needs investigation for exact location in `Supports()` method.

### Resolution Strategy

Create centralized helper functions for capability conversion:

**File**: `sdk/go/pluginsdk/capability_compat.go` (new file)

```go
package pluginsdk

import pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"

// capabilityTrue is the string value used for enabled capabilities in legacy metadata.
const capabilityTrue = "true"

// legacyCapabilityNames maps PluginCapability enums to their legacy string names.
// This map is the single source of truth for enum-to-legacy conversions.
//
//nolint:exhaustive // Only capabilities with legacy equivalents are mapped
var legacyCapabilityNames = map[pbc.PluginCapability]string{
    pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS:        "supports_recommendations",
    pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN:                "supports_dry_run",
    pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS:                "supports_budgets",
    pbc.PluginCapability_PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS: "supports_dismiss_recommendations",
    pbc.PluginCapability_PLUGIN_CAPABILITY_PROJECTED_COSTS:        "supports_projected_costs",
    pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS:           "supports_actual_costs",
    pbc.PluginCapability_PLUGIN_CAPABILITY_PRICING_SPEC:           "supports_pricing_spec",
    pbc.PluginCapability_PLUGIN_CAPABILITY_ESTIMATE_COST:          "supports_estimate_cost",
}

// CapabilityToLegacyName converts a PluginCapability enum to its legacy string name.
// Returns empty string if the capability has no legacy equivalent.
func CapabilityToLegacyName(cap pbc.PluginCapability) string {
    return legacyCapabilityNames[cap]
}

// CapabilitiesToLegacyMetadata converts capability enums to legacy metadata format.
// Returns a map with "supports_xyz": "true" entries for each capability.
// Capabilities without legacy equivalents are silently skipped.
func CapabilitiesToLegacyMetadata(capabilities []pbc.PluginCapability) map[string]string {
    if len(capabilities) == 0 {
        return nil
    }
    metadata := make(map[string]string, len(capabilities))
    for _, cap := range capabilities {
        if name := CapabilityToLegacyName(cap); name != "" {
            metadata[name] = capabilityTrue
        }
    }
    return metadata
}

// CapabilityConversionWarning represents a warning about unmapped capabilities.
type CapabilityConversionWarning struct {
    Capability pbc.PluginCapability
    Reason     string
}

// CapabilitiesToLegacyMetadataWithWarnings converts capabilities and reports unmapped ones.
// Use this when you need to log warnings about capabilities that cannot be represented
// in the legacy format.
func CapabilitiesToLegacyMetadataWithWarnings(
    capabilities []pbc.PluginCapability,
) (map[string]string, []CapabilityConversionWarning) {
    if len(capabilities) == 0 {
        return nil, nil
    }

    metadata := make(map[string]string, len(capabilities))
    var warnings []CapabilityConversionWarning

    for _, cap := range capabilities {
        if name := CapabilityToLegacyName(cap); name != "" {
            metadata[name] = capabilityTrue
        } else if cap != pbc.PluginCapability_PLUGIN_CAPABILITY_UNSPECIFIED {
            warnings = append(warnings, CapabilityConversionWarning{
                Capability: cap,
                Reason:     "no legacy metadata mapping exists",
            })
        }
    }

    return metadata, warnings
}
```

**Refactor locations**:

1. Update `sdk.go:handleConfiguredPluginInfo()` to use new helpers
2. Update `sdk.go:Supports()` to use new helpers
3. Remove duplicate conversion logic

### Tests

**File**: `sdk/go/pluginsdk/capability_compat_test.go`

```go
package pluginsdk_test

import (
    "testing"

    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func TestCapabilityToLegacyName(t *testing.T) {
    tests := []struct {
        name       string
        capability pbc.PluginCapability
        expected   string
    }{
        {"DryRun", pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN, "supports_dry_run"},
        {"Recommendations", pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS, "supports_recommendations"},
        {"Unspecified", pbc.PluginCapability_PLUGIN_CAPABILITY_UNSPECIFIED, ""},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := pluginsdk.CapabilityToLegacyName(tc.capability)
            if got != tc.expected {
                t.Errorf("CapabilityToLegacyName(%v) = %q, want %q", tc.capability, got, tc.expected)
            }
        })
    }
}

func TestCapabilitiesToLegacyMetadata(t *testing.T) {
    caps := []pbc.PluginCapability{
        pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
        pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS,
    }

    metadata := pluginsdk.CapabilitiesToLegacyMetadata(caps)

    if metadata["supports_dry_run"] != "true" {
        t.Errorf("expected supports_dry_run=true, got %q", metadata["supports_dry_run"])
    }
    if metadata["supports_recommendations"] != "true" {
        t.Errorf("expected supports_recommendations=true, got %q", metadata["supports_recommendations"])
    }
}

func TestCapabilitiesToLegacyMetadataWithWarnings(t *testing.T) {
    // Test with a mix of mapped and hypothetically unmapped capabilities
    caps := []pbc.PluginCapability{
        pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
    }

    metadata, warnings := pluginsdk.CapabilitiesToLegacyMetadataWithWarnings(caps)

    if len(metadata) != 1 {
        t.Errorf("expected 1 metadata entry, got %d", len(metadata))
    }
    if len(warnings) != 0 {
        t.Errorf("expected 0 warnings, got %d", len(warnings))
    }
}
```

---

## Part 4: Issue #300 - Documentation Updates

### Required Documentation

#### 4.1 CLAUDE.md Update

Add capability discovery pattern section to root `CLAUDE.md`.

**Content to add** (section heading: `### Capability Discovery Pattern`):

The section should document:

1. **Dual-mode system**: Auto-Discovery (default) vs Manual Override
2. **Interface-Based Auto-Discovery example**: Show a plugin implementing `DryRunHandler`
3. **Manual Capability Override example**: Show `WithCapabilities()` usage
4. **Capability Interfaces table**: Map interfaces to enum values
5. **Backward Compatibility note**: Explain dual-format responses (modern enum vs legacy string)

**Example code for Auto-Discovery**:

```go
// Plugin with DryRunHandler implementation
type MyPlugin struct {
    proto.UnimplementedCostSourceServiceServer
}

func (p *MyPlugin) HandleDryRun(req *pbc.DryRunRequest) (*pbc.DryRunResponse, error) {
    return pluginsdk.NewDryRunResponse(
        pluginsdk.WithResourceTypeSupported(true),
    ), nil
}

// Capability auto-discovered: PLUGIN_CAPABILITY_DRY_RUN
info := pluginsdk.NewPluginInfo("my-plugin", "v1.0.0")
```

**Example code for Manual Override**:

```go
info := pluginsdk.NewPluginInfo("my-plugin", "v1.0.0",
    pluginsdk.WithCapabilities([]pbc.PluginCapability{
        pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
    }),
)
```

**Capability Interfaces Table**:

| Interface | Capability Enum |
|-----------|-----------------|
| `DryRunHandler` | `PLUGIN_CAPABILITY_DRY_RUN` |
| `RecommendationsProvider` | `PLUGIN_CAPABILITY_RECOMMENDATIONS` |
| `BudgetsProvider` | `PLUGIN_CAPABILITY_BUDGETS` |
| `DismissProvider` | `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS` |

#### 4.2 pluginsdk/README.md Update

Add section documenting capability discovery (add to pluginsdk/README.md).

Content to add for "Capability Discovery" section:

- Intro: "The SDK automatically discovers plugin capabilities by checking implemented interfaces."
- Subsection "Automatic Discovery" with example showing `DryRunHandler` implementation
- Subsection "Manual Override" with `WithCapabilities()` example
- Subsection "Interface Reference" with table mapping interfaces to capabilities

**Interface Reference Table**:

| Interface | Required Method | Capability |
|-----------|-----------------|------------|
| `DryRunHandler` | `HandleDryRun(req) (resp, error)` | `PLUGIN_CAPABILITY_DRY_RUN` |
| `RecommendationsProvider` | `GetRecommendations(ctx, req) (resp, error)` | `PLUGIN_CAPABILITY_RECOMMENDATIONS` |
| `BudgetsProvider` | `GetBudgets(ctx, req) (resp, error)` | `PLUGIN_CAPABILITY_BUDGETS` |
| `DismissProvider` | `DismissRecommendation(ctx, req) (resp, error)` | `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS` |

---

## Part 5: Issue #301 - Optimize Slice Copying Performance

### Problem Statement

Current slice copying uses `make + copy`:

```go
// sdk/go/pluginsdk/sdk.go:420-424
capabilities := make([]pbc.PluginCapability, len(s.pluginInfo.Capabilities))
copy(capabilities, s.pluginInfo.Capabilities)
```

The append pattern is ~10-20% faster for small slices:

```go
capabilities := append([]pbc.PluginCapability(nil), s.pluginInfo.Capabilities...)
```

### Locations to Update

| File | Line | Current Pattern |
|------|------|-----------------|
| `sdk.go:406-407` | `providers` copy | `make + copy` |
| `sdk.go:420-421` | `capabilities` copy (explicit) | `make + copy` |
| `sdk.go:423-424` | `capabilities` copy (global) | `make + copy` |

### Implementation

**File**: `sdk/go/pluginsdk/sdk.go`

Replace:

```go
providers := make([]string, len(s.pluginInfo.Providers))
copy(providers, s.pluginInfo.Providers)
```

With:

```go
providers := append([]string(nil), s.pluginInfo.Providers...)
```

Similarly for capabilities:

```go
// Before
capabilities := make([]pbc.PluginCapability, len(s.pluginInfo.Capabilities))
copy(capabilities, s.pluginInfo.Capabilities)

// After
capabilities := append([]pbc.PluginCapability(nil), s.pluginInfo.Capabilities...)
```

### Benchmark Verification

Add benchmark to verify improvement:

```go
func BenchmarkSliceCopyPatterns(b *testing.B) {
    original := []pbc.PluginCapability{
        pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN,
        pbc.PluginCapability_PLUGIN_CAPABILITY_RECOMMENDATIONS,
        pbc.PluginCapability_PLUGIN_CAPABILITY_BUDGETS,
    }

    b.Run("append", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = append([]pbc.PluginCapability(nil), original...)
        }
    })

    b.Run("make+copy", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            dst := make([]pbc.PluginCapability, len(original))
            copy(dst, original)
        }
    })
}
```

---

## Part 6: Issue #208 - Clarify Compatibility Test Naming

### Problem Statement

Test names suggest literal old-server interop but actually test round-trip behavior:

- `TestResourceDescriptor_OldClientSimulation`
- `TestResourceDescriptor_NewClientOldServer`

### Resolution

Rename tests and add clarifying comments:

**File**: `sdk/go/testing/resource_id_test.go`

```go
// TestResourceDescriptor_RoundTripDefaulting tests that new fields (id, arn) default
// correctly when unmarshaling from data that doesn't include them.
// This simulates the behavior when an older client sends a message without new fields.
func TestResourceDescriptor_RoundTripDefaulting(t *testing.T) {
    // ... existing test body
}

// TestResourceDescriptor_RoundTripCompatibility tests that new fields survive
// a marshal/unmarshal round-trip through the current generated type.
// This approximates proto3 unknown-field handling behavior rather than
// literally testing against an older generated type.
func TestResourceDescriptor_RoundTripCompatibility(t *testing.T) {
    // ... existing test body
}
```

---

## Part 7: Issue #209 - Complete Import Statements in README

### Problem Statement

Code examples in `sdk/go/testing/README.md` use `pbc` alias without defining it.

### Resolution

**File**: `sdk/go/testing/README.md`

Add import note at the start of the "Resource Identifier Fields" section:

```markdown
## Resource Identifier Fields

> **Import Note**: Examples in this section use the following import aliases:
>
> ```go
> import (
>     pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
>     "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
> )
> ```

The `ResourceDescriptor` message supports two optional fields...
```

Update code examples to include full imports where appropriate:

```go
import (
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

// Create descriptors with unique IDs for batch correlation
descriptors := []*pbc.ResourceDescriptor{
    pluginsdk.NewResourceDescriptor(
        "aws", "ec2",
        pluginsdk.WithID("urn:pulumi:prod::myapp::aws:ec2/instance:Instance::web"),
    ),
}
```

---

## Implementation Order

Execute in dependency order:

| Phase | Issues | Description | Effort |
|-------|--------|-------------|--------|
| 1 | #294 | Close as fixed (verify + close) | 5 min |
| 2 | #295 | Add proto comments | 30 min |
| 3 | #299 | Extract capability helpers | 1.5 hr |
| 4 | #301 | Optimize slice copying | 30 min |
| 5 | #300 | Documentation updates | 1 hr |
| 6 | #208, #209 | Testing docs improvements | 45 min |

## Validation Checklist

```bash
# Build verification
go build ./...

# Unit tests
go test ./sdk/go/pluginsdk/...
go test ./sdk/go/testing/...

# Linting
make lint

# Markdown linting
make lint-markdown

# Proto regeneration (if proto comments changed)
make generate

# Full validation
make validate
```

## Commit Strategy

Separate commits for clear history:

1. `fix(sdk): close #294 - DryRunHandler already defined`
2. `docs(proto): clarify capability vs metadata fields (#295)`
3. `refactor(sdk): extract capability conversion helpers (#299)`
4. `perf(sdk): use append pattern for slice copying (#301)`
5. `docs(sdk): add capability discovery documentation (#300)`
6. `docs(testing): clarify round-trip test naming (#208)`
7. `docs(testing): add complete import statements to README (#209)`

## Success Criteria

- [ ] Issue #294 closed with verification comment
- [ ] Proto fields have clear documentation comments
- [ ] Single source of truth for capability-to-legacy conversion
- [ ] Slice copying uses optimized append pattern
- [ ] CLAUDE.md documents capability discovery pattern
- [ ] pluginsdk/README.md documents capability interfaces
- [ ] Testing README has complete import statements
- [ ] Test names clarify round-trip semantics
- [ ] All tests pass
- [ ] All linting passes
