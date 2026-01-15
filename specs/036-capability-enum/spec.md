# Feature Specification: Add PluginCapability Enum

**Feature Branch**: `036-capability-enum`
**Created**: 2026-01-14
**Status**: Draft
**Input**: User description: "Add PluginCapability Enum for Feature Discovery..."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Capability Discovery (Priority: P1)

As a Consumer Application Developer (e.g., finfocus-core),
I want to discover plugin capabilities using a strongly-typed enum
So that I can reliably route requests to plugins that support specific features (like Carbon metrics or Budgets)
without relying on error-prone string keys.

**Why this priority**: Core value proposition. Enables safe and reliable feature negotiation.

**Independent Test**: Can be tested by creating a mock plugin that implements specific interfaces
and verifying the consumer receives the correct enum values.

**Acceptance Scenarios**:

1. **Given** a plugin that implements `GetProjectedCost` and `GetRecommendations`,
   **When** `GetPluginInfo` is called,
   **Then** the response includes `PLUGIN_CAPABILITY_PROJECTED_COSTS` and
   `PLUGIN_CAPABILITY_RECOMMENDATIONS` in the `capabilities` field.

2. **Given** a plugin that supports Carbon metrics,
   **When** `GetPluginInfo` is called,
   **Then** the response includes `PLUGIN_CAPABILITY_CARBON`.

---

### User Story 2 - SDK Auto-Discovery for Plugin Developers (Priority: P1)

As a Plugin Developer,
I want the SDK to automatically detect my plugin's capabilities based on the RPCs I implement
So that I don't have to manually maintain a list of capabilities and risk it getting out of sync with my code.

**Why this priority**: Improves Developer Experience (DX) and reduces bugs.

**Independent Test**: Implement a plugin using the Go SDK with only `GetActualCost`
and verify the SDK automatically reports `PLUGIN_CAPABILITY_ACTUAL_COSTS`.

**Acceptance Scenarios**:

1. **Given** a plugin implementation using the Go SDK,
   **When** the SDK initializes the plugin server,
   **Then** it introspects the implemented methods and populates the `capabilities` list automatically.

---

### User Story 3 - Backward Compatibility (Priority: P2)

As a Consumer Application,
I want existing string-based capability checks to continue working
So that I am not forced to refactor all my existing integration code immediately.

**Why this priority**: Prevents breaking changes for existing integrations.

**Independent Test**: Verify that the `metadata` or `capabilities` map (legacy) still contains the expected string keys.

**Acceptance Scenarios**:

1. **Given** a plugin that implements `GetRecommendations`,
   **When** `GetPluginInfo` or `Supports` is called,
   **Then** the legacy string-based `capabilities` map still contains `"recommendations": true`.

### Edge Cases

- **New/Unsupported Capabilities**: If a plugin offers a feature not yet defined in the
  `PluginCapability` enum, it cannot be advertised via the typed field. It MUST rely on the
  legacy string map until the enum is updated.
- **Version Mismatch**: If a Consumer uses an older version of the SDK/Proto than the Plugin,
  it may receive unknown enum values (which Protobuf handles gracefully as integers). The
  Consumer MUST handle unknown enum values gracefully.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST define a `PluginCapability` enum in `proto/finfocus/v1/enums.proto`
  with values for all currently supported features (Projected Costs, Actual Costs,
  Carbon, Recommendations, Dry Run, Budgets, Energy, Water).
- **FR-002**: System MUST add a `repeated PluginCapability capabilities` field to the
  `GetPluginInfoResponse` message in `costsource.proto`.
- **FR-003**: The Go SDK MUST automatically detect implemented RPCs via runtime interface
  type assertions (e.g., `plugin.(SupportsProvider)`) and populate the `PluginCapability`
  list in `GetPluginInfo`.
- **FR-004**: The system MUST maintain the existing `map<string, bool>` capabilities field
  for backward compatibility.
- **FR-005**: The project Constitution MUST be updated to require a corresponding
  `PluginCapability` enum value for every new RPC added to `CostSourceService`.

### Key Entities

- **PluginCapability**: A comprehensive enum representing all functional features a plugin can offer.
- **GetPluginInfoResponse**: The response message containing the plugin's metadata and now its typed capabilities.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All 8 defined capabilities (Projected, Actual, Carbon, Recommendations, Dry Run,
  Budgets, Energy, Water) are representable via the new enum.
- **SC-002**: A plugin implementing `GetProjectedCost` automatically reports
  `PLUGIN_CAPABILITY_PROJECTED_COSTS` without manual configuration in the plugin code.
- **SC-003**: `GetPluginInfo` calls return both the new `capabilities` enum list and the
  legacy string map, ensuring 100% backward compatibility for existing clients.
