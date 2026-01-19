# Research & Design: Capability Discovery Polish

**Feature**: Capability Discovery Polish
**Status**: Complete
**Date**: 2026-01-18

## Overview

This feature is a consolidation of well-understood quality improvements. Research primarily
focused on verifying the current state of the codebase to confirm assumptions made in the
consolidated development plan.

## Findings

### 1. DryRunHandler Interface (Issue #294)

- **Status**: Confirmed existing.
- **Location**: `sdk/go/pluginsdk/dry_run.go`
- **Definition**:

  ```go
  type DryRunHandler interface {
      HandleDryRun(req *pbc.DryRunRequest) (*pbc.DryRunResponse, error)
  }
  ```

- **Conclusion**: Issue #294 can be closed immediately as "Already Fixed".

### 2. Proto Field Inconsistency (Issue #295)

- **Status**: Confirmed.
- **File**: `proto/finfocus/v1/costsource.proto`
- **Observation**: `GetPluginInfoResponse` has both `capabilities` (enum) and `metadata` (map).
- **Plan**: Add comments clarifying that `metadata` is legacy/deprecated and `capabilities` is the
  modern standard.

### 3. Slice Copying (Issue #301)

- **Status**: Confirmed use of `make` + `copy` in `sdk/go/pluginsdk/sdk.go`.
- **Plan**: Replace with `append([]T(nil), src...)` pattern.
- **Reference**: Go generic slice copy optimization patterns often favor `append` for conciseness
  and similar/better performance for small slices.

## Design Decisions

### Single Source of Truth for Capabilities

We will introduce `sdk/go/pluginsdk/capability_compat.go` to hold the `legacyCapabilityNames`
map. This will be the **only** place where Enum -> String mapping happens.

```go
var legacyCapabilityNames = map[pbc.PluginCapability]string{
    pbc.PluginCapability_PLUGIN_CAPABILITY_DRY_RUN: "supports_dry_run",
    // ...
}
```

### Auto-Discovery Logic

Auto-discovery will happen in `NewPluginInfo`. It will check for:

- `DryRunHandler` -> `PLUGIN_CAPABILITY_DRY_RUN`
- `RecommendationsProvider` -> `PLUGIN_CAPABILITY_RECOMMENDATIONS`
- `BudgetsProvider` -> `PLUGIN_CAPABILITY_BUDGETS`
- `DismissProvider` -> `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS`

**Priority Rule**: If `WithCapabilities(...)` option is provided, it **overrides** all
auto-discovery. This allows plugins to explicitly disable features even if they implement the
interface (e.g., behind a feature flag).

## Agent Context Updates

No major architectural changes or new technologies. The `finfocus-senior-engineer` agent context
should be updated to know about the new "Capability Discovery Pattern" to avoid suggesting manual
string maps in the future.
