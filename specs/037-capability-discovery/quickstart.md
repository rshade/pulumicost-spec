# Quickstart Guide: Dual-Layer Capability Discovery

## Overview

The FinFocus SDK now supports "Dual-Layer Capability Discovery" to help Hosts understand what
your plugin supports both Globally (at startup) and Granularly (per resource).

- **Global Discovery**: "Does this plugin support Recommendations?"
- **Granular Discovery**: "Does this plugin support Recommendations *for this specific RDS instance*?"

## SDK Auto-Discovery (Opt-Out)

By default, the Go SDK automatically detects which features you support based on the interfaces you implement.

**Example**:
If you implement the `RecommendationsProvider` interface:

```go
func (p *MyPlugin) GetRecommendations(...) (...) { ... }
```

The SDK will automatically advertise `PLUGIN_CAPABILITY_RECOMMENDATIONS` globally.

## Granular Support

To explicitly enable/disable features for specific resources, implement the `Supports` method. You
can receive a `ResourceDescriptor` to make decisions.

```go
func (p *MyPlugin) Supports(ctx context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
    res := req.GetResource()
    
    // Default: Return empty capabilities list to inherit Global capabilities automatically
    // This reduces boilerplate if you only want to check basic "Supported" status.
    caps := []pbc.PluginCapability{}

    // Granular Logic Example:
    // Legacy resources don't support Recommendations
    if res.GetTags()["legacy"] == "true" {
        // Explicitly define subset of capabilities (overrides Global inheritance)
        caps = []pbc.PluginCapability{
            pbc.PluginCapability_PLUGIN_CAPABILITY_ACTUAL_COSTS,
            // Recommendations removed
        }
    }

    return &pbc.SupportsResponse{
        Supported:        true,
        CapabilitiesEnum: caps,
    }, nil
}
```

## Enforcing Granular Support (Strict Error Handling)

If you advertise that a resource does NOT support a capability (e.g., Recommendations), you should
enforce this in your RPC implementation:

```go
func (p *MyPlugin) GetRecommendations(ctx context.Context, req *pbc.GetRecommendationsRequest) (*pbc.GetRecommendationsResponse, error) {
    // Check if target resources support recommendations
    for _, res := range req.GetTargetResources() {
        if res.GetTags()["legacy"] == "true" {
             return nil, status.Error(codes.Unimplemented, "recommendations not supported for legacy resource")
        }
    }
    // ... implementation ...
}
```

## Backward Compatibility

The SDK automatically syncs your typed `CapabilitiesEnum` to the legacy `Capabilities` string map,
so older Hosts will continue to work without changes.
