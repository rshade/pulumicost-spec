# Feature Specification: Capability Discovery Polish

**Feature Branch**: `038-capability-discovery-polish`
**Created**: 2026-01-18
**Status**: Draft
**Input**: Consolidates 7 related issues (#294, #295, #299, #300, #301, #208, #209) into a
single cohesive development plan for polishing the capability discovery system, improving SDK
quality, and fixing documentation.

## User Scenarios & Testing

### User Story 1 - Automatic Capability Discovery (Priority: P1)

As a plugin developer, I want the SDK to automatically detect which capabilities my plugin
supports based on the interfaces I implement, so that I don't have to manually configure them.

**Why this priority**: Core value proposition of the SDK to reduce boilerplate and
configuration errors.

**Independent Test**: Create a struct implementing `DryRunHandler`. Initialize `PluginInfo`.
Verify `PLUGIN_CAPABILITY_DRY_RUN` is present in the resulting capabilities list.

**Acceptance Scenarios**:

1. **Given** a plugin struct that implements `DryRunHandler`
   **When** `NewPluginInfo` is called
   **Then** the returned info includes `PLUGIN_CAPABILITY_DRY_RUN`
   **And** the legacy metadata includes `"supports_dry_run": "true"`

2. **Given** a plugin struct that implements `RecommendationsProvider`
   **When** `NewPluginInfo` is called
   **Then** the returned info includes `PLUGIN_CAPABILITY_RECOMMENDATIONS`

---

### User Story 2 - Manual Capability Override (Priority: P1)

As a plugin developer, I want to explicitly define my plugin's capabilities using
`WithCapabilities()`, so that I can override auto-discovery if needed (e.g., to temporarily disable
a feature or when using a dynamic proxy).

**Why this priority**: Essential for advanced use cases and debugging.

**Independent Test**: Initialize `PluginInfo` with `WithCapabilities` passing a specific list.
Verify only those capabilities are present, ignoring implemented interfaces.

**Acceptance Scenarios**:

1. **Given** a plugin struct that implements `DryRunHandler`
   **When** `NewPluginInfo` is called with `WithCapabilities([]PluginCapability{})` (empty list)
   **Then** the returned info has NO capabilities (auto-discovery is bypassed)

---

### User Story 3 - Backward Compatibility for Legacy Clients (Priority: P1)

As a maintainer of a legacy host application, I want to see capability flags in the `metadata`
map (e.g., "supports_dry_run": "true"), so that my application continues to work with new plugins
without code changes.

**Why this priority**: strict backward compatibility requirement for the
ecosystem.

**Independent Test**: Unit test the `CapabilitiesToLegacyMetadata` conversion logic.

**Acceptance Scenarios**:

1. **Given** a plugin with `PLUGIN_CAPABILITY_ESTIMATE_COST`
   **When** `GetPluginInfo` is called by a client
   **Then** the response `metadata` map contains `"supports_estimate_cost": "true"`

---

### User Story 4 - SDK Performance Optimization (Priority: P2)

As a high-throughput system operator, I want the SDK to use efficient memory allocation
patterns (specifically for slice copying), so that garbage collection overhead is minimized.

**Why this priority**: Identified as an enhancement (#301) for SDK quality.

**Independent Test**: Benchmark `append` pattern vs `make+copy` pattern.

**Acceptance Scenarios**:

1. **Given** the SDK needs to copy a slice of providers or capabilities
   **When** the operation occurs
   **Then** it uses the optimized `append([]T(nil), src...)` pattern

### Edge Cases

- **Unmapped Capabilities**: What happens if a new capability enum has no legacy string
  equivalent?
  - _Expectation_: It is omitted from legacy metadata, potentially with a warning log (handled by
    `CapabilitiesToLegacyMetadataWithWarnings`).
- **Mixed Configuration**: User provides `WithCapabilities` but also expects some
  auto-discovery?
  - _Expectation_: `WithCapabilities` is authoritative and completely replaces auto-discovery.

## Requirements

### Functional Requirements

- **FR-001**: The SDK MUST define and recognize the `DryRunHandler` interface for auto-discovery
  of `PLUGIN_CAPABILITY_DRY_RUN` (Issue #294).
- **FR-002**: The Proto definitions (`GetPluginInfoResponse`, `SupportsResponse`) MUST include
  comments clarifying the relationship and deprecation status of `capabilities` (enum) vs
  `metadata`/`capabilities` (legacy map) (Issue #295).
- **FR-003**: The SDK MUST provide a centralized helper function `CapabilitiesToLegacyMetadata`
  to convert enums to legacy string map (Issue #299).
- **FR-004**: The SDK MUST use the `append` pattern (`dst := append([]T(nil), src...)`) for slice
  copying instead of `make` + `copy` (Issue #301).
- **FR-005**: Documentation (`CLAUDE.md`, `pluginsdk/README.md`) MUST explicitly document the
  "Capability Discovery Pattern", including Auto-Discovery vs Manual Override and the Interface
  Reference table (Issue #300).
- **FR-006**: Compatibility tests in `sdk/go/testing` MUST be renamed and commented to clarify
  they test round-trip serialization, not literal interaction with old servers (Issue #208).
- **FR-007**: Documentation code examples in `sdk/go/testing/README.md` MUST include complete and
  correct import statements (Issue #209).

### Key Entities

- **PluginCapability**: Enum defining features a plugin supports (e.g., DryRun, Recommendations).
- **PluginInfo**: Struct containing metadata and capabilities of a plugin.
- **Legacy Metadata**: `map[string]string` used for backward compatibility.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All 7 identified issues (#294, #295, #299, #300, #301, #208, #209) are resolved
  and closed.
- **SC-002**: `make validate` (including linting and tests) passes successfully.
- **SC-003**: Benchmarks confirm the slice copying optimization is at least neutral or positive in
  performance.
- **SC-004**: Capability conversion logic exists in exactly one location in the codebase
  (duplicate logic removed).
- **SC-005**: `CLAUDE.md` and `pluginsdk/README.md` contain the specific "Capability Discovery"
  sections defined in requirements.
