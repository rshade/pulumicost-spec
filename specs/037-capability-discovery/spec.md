# Feature Specification: Dual-Layer Capability Discovery

**Feature Branch**: `037-capability-discovery`
**Created**: 2026-01-15
**Status**: Draft
**Input**: Issue #194, PR #290, User Description: "Formalize a dual-layer capability discovery protocol..."

## Clarifications

### Session 2026-01-15

- **Q**: What specific information should the Host pass to Granular Discovery (`Supports`)?
  **A**: ResourceDescriptor (Provider, Service, Type, Region).
- **Q**: SDK auto-discovery: Opt-In or Opt-Out?
  **A**: Opt-Out (Assume support for all implemented interfaces by default).
- **Q**: How to handle granular capability requests for unsupported resources?
  **A**: Strict (Reject with clear error code).
- **Q**: Allow dynamic custom capability keys in `capabilities_enum`?
  **A**: No — Static List only (from `PluginCapability` enum).
- **Q**: Empty `capabilities_enum` in `Supports`: "None Supported" or "Inherit Global"?
  **A**: Inherit Global (supports incremental adoption).

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Global Feature Registration

As a Host System (e.g., FinFocus Core),
I want to identify which major feature sets (Budgets, Recommendations) a plugin supports at initialization
So that I can register internal service handlers and UI components only for supported features.

**Why this priority**: Essential for system stability and UI performance.

**Independent Test**:
Configure a plugin to support only "Budgets". Verify the Host initializes the Budgets service but disables Recommendations.

**Acceptance Scenarios**:

1. **Given** a plugin that implements Budgets but not Recommendations,
   **When** the plugin initializes and reports its information,
   **Then** the response includes the "Budgets" capability but NOT "Recommendations".
2. **Given** the above plugin,
   **When** the Host processes the response,
   **Then** the Host registers the Budgets handler and does NOT attempt to call Recommendations services.

---

### User Story 2 - Granular Resource Support

As a Host System,
I want to verify if a globally supported feature is valid for a specific resource
So that I don't waste resources or cause errors by calling services for unsupported resource types.

**Why this priority**: Optimizes performance and reduces error log noise.

**Independent Test**:
Simulate a resource that supports "Carbon" and one that does not, within a plugin that supports Carbon globally.

**Acceptance Scenarios**:

1. **Given** a plugin that supports "Recommendations" globally,
   **And** a specific legacy resource that does NOT support recommendations,
   **When** the Host queries support for that resource,
   **Then** the list of applicable capabilities in the response does NOT contain "Recommendations".
2. **Given** the above response,
   **When** the Host decides next steps,
   **Then** the Host skips the Recommendations service call for that specific resource.

### Edge Cases

- **Version Mismatch**: If a Host requests capabilities from an older plugin version that doesn't
  support the new discovery protocol, the system MUST fall back to legacy or assume global support.
- **Partial Support**: If a resource supports a feature "partially" (e.g., read but not write), the
  system currently treats this as "Supported" (binary) unless specific granular capabilities defined.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST implement a global discovery mechanism where plugins advertise major capabilities at initialization.
- **FR-002**: The system MUST implement granular discovery where plugins advertise
  resource-specific capabilities using `ResourceDescriptor` (Provider, Service, Type, Region).
- **FR-003**: The support response MUST include strongly-typed capabilities applicable to that
  resource, restricted to `PluginCapability` enum values.
- **FR-004**: The Developer Kit (SDK) MUST automatically detect supported features from plugin code
  structure and populate the capability list by default (Opt-Out mechanism).
- **FR-005**: The Host Application MUST NOT invoke a feature-specific service if the plugin did not
  advertise the corresponding global capability during initialization.
- **FR-006**: The Host Application MUST NOT invoke a feature-specific service for a resource if the
  granular support check indicates missing capability, even if present globally.
- **FR-007**: The Plugin MUST reject requests for specific capabilities on resources where support
  has been explicitly denied (e.g., via `Supports` check) with clear error code (e.g., Unimplemented).
- **FR-008**: If granular `Supports` response contains empty `capabilities_enum` list, the Host
  SHOULD assume resource supports all globally advertised capabilities (Backward Compatibility).

### System Constraints

- **CON-001**: Changes MUST be 100% backward compatible with existing string-based capability maps.
- **CON-002**: No manual configuration required for plugin developers to enable granular discovery
  when using standard interfaces.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Plugin developers can enable granular capability discovery solely by implementing
  standard interfaces (Zero manual config).
- **SC-002**: Host system avoids 100% of service calls for unsupported resource features (verified
  via call logs).
- **SC-003**: Existing integration tests relying on legacy string-based maps pass without
  modification.
