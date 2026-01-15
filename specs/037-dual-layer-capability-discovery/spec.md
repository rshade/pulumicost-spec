# Feature Specification: Dual-Layer Capability Discovery

**Feature Branch**: `feat/plugin-capability-discovery`
**Created**: 2026-01-15
**Status**: Draft
**Input**: Issue #194, PR #290

## Objective

Formalize a dual-layer capability discovery protocol that enables both high-level service registration (Global) and granular, resource-specific feature negotiation (Granular). This ensures the Host only calls RPCs that a plugin explicitly supports for a given context.

## Discovery Layers

### 1. Global Discovery (`GetPluginInfo`)

*   **When**: Called once during plugin initialization.
*   **Purpose**: To determine which major feature sets (Budgets, Recommendations, etc.) the plugin supports across its entire lifecycle.
*   **Host Action**: The Host uses this to register internal service handlers. If a plugin does not advertise `PLUGIN_CAPABILITY_BUDGETS`, the Host should never attempt to call `GetBudgets`.

### 2. Granular Discovery (`Supports`)

*   **When**: Called per resource or resource type during execution.
*   **Purpose**: To determine if a globally supported feature is applicable to a specific `ResourceDescriptor` (Provider, Region, Resource Type).
*   **Host Action**: The Host uses this to decide whether to execute specific logic for a resource. For example, a plugin might support Recommendations globally, but return `false` for `PLUGIN_CAPABILITY_RECOMMENDATIONS` in the `Supports` response for a specific legacy SKU.

## Technical Implementation

### Protobuf Changes

`SupportsResponse` is enhanced with a `capabilities_enum` field:

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2;
  map<string, bool> capabilities = 3; // Legacy string map
  repeated MetricKind supported_metrics = 4;
  repeated PluginCapability capabilities_enum = 5; // New strongly-typed list
}
```

### SDK Auto-Discovery

The Go SDK automatically populates these fields by introspecting implemented interfaces:

*   `RecommendationsProvider` -> `PLUGIN_CAPABILITY_RECOMMENDATIONS`
*   `BudgetsProvider` -> `PLUGIN_CAPABILITY_BUDGETS`
*   `DryRunProvider` -> `PLUGIN_CAPABILITY_DRY_RUN` (Future)

## Host Integration Requirements (FinFocus Core)

Consuming applications (like `finfocus-core`) SHOULD follow this logic:

1.  **Init Phase**: Call `GetPluginInfo`. Store global capabilities.
2.  **Registration**: Only enable feature-specific UI/API components if the capability is present globally.
3.  **Execution Phase**:
    *   Call `Supports(resource)`.
    *   Check `capabilities_enum` in the response.
    *   If a capability (e.g., `CARBON`) is present globally but **absent** in the granular `Supports` response, the Host MUST NOT call the related RPC for that resource.

## Success Criteria

1.  Plugins can advertise granular support without extra code (auto-discovery).
2.  Hosts can avoid unnecessary RPC failures by checking support before calling.
3.  100% backward compatibility with string-based `capabilities` map.
