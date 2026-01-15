# Research & Decisions - Dual-Layer Capability Discovery

**Feature Branch**: `037-capability-discovery`
**Status**: Completed

## Research Findings

### 0.1 Protobuf Compatibility

**Question**: Does modifying `SupportsRequest` to take `ResourceDescriptor` break existing clients?
**Analysis**:

- Current `SupportsRequest` definition in `proto/finfocus/v1/costsource.proto`:

  ```protobuf
  message SupportsRequest {
    // resource contains the resource descriptor to check support for
    ResourceDescriptor resource = 1;
  }
  ```

- The `SupportsRequest` ALREADY takes a `ResourceDescriptor` as field 1.
- **Finding**: No breaking change needed for the Request message. We can re-use the existing field.
- **Constraint Verified**: Backward compatibility is maintained (request message structure unchanged).

### 0.2 SDK Reflection Performance

**Question**: Impact of reflection during `Serve` initialization?
**Analysis**:

- `Serve` is called once at plugin startup.
- Reflection logic (checking if `s.plugin` implements `BudgetsProvider`, `RecommendationsProvider`) takes microseconds.
- **Finding**: Negligible performance impact. Safe to implement in `Serve` or `NewServer`.

### 0.3 ResourceDescriptor Source

**Question**: Where does the Host get the `ResourceDescriptor`?
**Analysis**:

- Host already constructs `ResourceDescriptor` for other calls like `GetActualCost` or `GetRecommendations`.
- **Finding**: The Host has this object available in its internal resource graph processing loop.

## Design Decisions

### Decision 1: Protobuf Changes

- **Request**: No change needed to `SupportsRequest` (it already has `ResourceDescriptor`).
- **Response**: Add `repeated PluginCapability capabilities_enum = 5;` to `SupportsResponse`.
  - Field 3 is `map<string, bool> capabilities` (legacy).
  - Field 4 is `repeated MetricKind supported_metrics`.
  - Field 5 is the next available ID.

### Decision 2: SDK Auto-Discovery Logic

- **Mechanism**: "Opt-Out" via `Serve` initialization.
- **Implementation**:
  - In `NewServer` (or `NewServerWithOptions`), inspect `plugin` implementation.
  - Build a `[]PluginCapability` slice based on interfaces:
    - `BudgetsProvider` -> `PLUGIN_CAPABILITY_BUDGETS`
    - `RecommendationsProvider` -> `PLUGIN_CAPABILITY_RECOMMENDATIONS`
    - `DismissProvider` -> `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS` (Need to verify enum name)
    - `GetProjectedCost` (Core Plugin interface) -> `PLUGIN_CAPABILITY_PROJECTED_COSTS`
    - `GetActualCost` (Core Plugin interface) -> `PLUGIN_CAPABILITY_ACTUAL_COSTS`
  - Store this list in the `Server` struct.
  - In `GetPluginInfo` (Global): Return this list.
  - In `Supports` (Granular):
    - If user returns a response with empty `capabilities_enum`, populate it with the stored "Global"
      list (Inherit Global).
    - If user explicitly sets `capabilities_enum` (even to empty list? No, empty implies inherit per
      clarifications), they control it.
    - **Refinement from Clarification**: "If a plugin provides an empty `capabilities_enum` list in a
      `Supports` response, should the Host assume 'None Supported' or 'Inherit Global'? → A: Inherit
      Global".
    - So SDK logic: Call user `Supports`. If `resp.CapabilitiesEnum` is empty, copy global capabilities into it.

### Decision 3: Enum Definitions

- Need to check `proto/finfocus/v1/enums.proto` to ensure all needed capabilities exist.

- Required Enums:
  - `PLUGIN_CAPABILITY_PROJECTED_COSTS`
  - `PLUGIN_CAPABILITY_ACTUAL_COSTS`
  - `PLUGIN_CAPABILITY_PRICING_SPEC`
  - `PLUGIN_CAPABILITY_ESTIMATE_COST`
  - `PLUGIN_CAPABILITY_RECOMMENDATIONS`
  - `PLUGIN_CAPABILITY_BUDGETS`
  - `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS` (or similar)
  - `PLUGIN_CAPABILITY_DRY_RUN`

### Decision 4: Legacy Compatibility

- SDK must sync `capabilities_enum` to `map<string, bool> capabilities` in both `GetPluginInfo` and `Supports`.
- Mapping:
  - `PLUGIN_CAPABILITY_RECOMMENDATIONS` <-> "recommendations"
  - `PLUGIN_CAPABILITY_BUDGETS` <-> "budgets"
  - `PLUGIN_CAPABILITY_DRY_RUN` <-> "dry_run"
