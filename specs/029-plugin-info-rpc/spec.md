<!-- markdownlint-disable MD013 -->
# Feature Specification: Add GetPluginInfo RPC

**Feature Branch**: `029-plugin-info-rpc`  
**Created**: 2025-12-30  
**Status**: Draft  
**Input**: User description: "Add GetPluginInfo() RPC for spec version compatibility checking"

## Clarifications

### Session 2025-12-30

- Q: What is the timeout for the `GetPluginInfo` call? → A: 5 seconds (Option B).
- Q: What format must `spec_version` follow? → A: SemVer (vX.Y.Z) (Option A).
- Q: How should malformed metadata be handled? → A: Strict (Option A) - Fail initialization if Name, Version, or SpecVersion are missing/invalid.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Compatibility Verification (Priority: P1)

As the Core System, I want to request metadata from a loaded plugin so that I can verify if its spec version is compatible with my current version.

**Why this priority**: Essential for preventing runtime errors caused by mismatched protocol versions as the ecosystem grows.

**Independent Test**: Can be fully tested by creating a mock plugin implementing the RPC and verifying the Core receives the correct data.

**Acceptance Scenarios**:

1. **Given** a plugin that implements `GetPluginInfo`, **When** the Core calls this RPC, **Then** it receives the plugin's name, version, and `spec_version`.
2. **Given** a plugin returns a `spec_version`, **When** the Core checks it, **Then** it can determine if the version is supported.

---

### User Story 2 - Diagnostic Visibility (Priority: P2)

As a Developer or Operator, I want to see the spec version of installed plugins when listing them, so that I can debug compatibility issues or identify outdated plugins.

**Why this priority**: Provides visibility and easier troubleshooting for users managing multiple plugins.

**Independent Test**: Can be tested by running the plugin list command and verifying output contains spec version info.

**Acceptance Scenarios**:

1. **Given** multiple installed plugins, **When** a consumer calls `GetPluginInfo` on each,
   **Then** the response includes `spec_version` that can be displayed in diagnostic tools.

> **Note**: The actual CLI/list command is a consumer-side implementation detail, not part of this SDK.

---

### User Story 3 - Graceful Degradation (Priority: P3)

As the Core System, I want to handle plugins that do not implement `GetPluginInfo` gracefully, so that existing plugins continue to work without modification.

**Why this priority**: Ensures backward compatibility and prevents breaking changes for existing deployments.

**Independent Test**: Can be tested by attempting to call the RPC on a plugin that does not implement it and asserting the system catches the error.

**Acceptance Scenarios**:

1. **Given** a legacy plugin that does not implement `GetPluginInfo`, **When** the Core calls the RPC, **Then** it receives an "Unimplemented" error (or similar) and treats the spec version as "Unknown" or logs a warning, but does not crash.

---

### Edge Cases

- **Malformed/empty metadata**: Handled by FR-007 (SDK validates SemVer format) and FR-008
  _(Consumer)_ (Core validates required fields). SDK returns validation error; consumer fails init.
- **Plugin hangs during call**: Handled by FR-005 _(Consumer)_ with 5-second timeout. SDK has no
  timeout responsibility; consumers must implement.
- **spec_version format changes**: Out of scope for this feature. Future spec versions would define
  migration path. Current implementation uses SemVer (vX.Y.Z) per FR-007.

## Requirements _(mandatory)_

### Scope Clarification

This feature implements the **plugin SDK side** of GetPluginInfo. The following are **consumer-side
responsibilities** that plugin SDK users (e.g., pulumicost-core) must implement:

- Calling GetPluginInfo with appropriate timeout
- Handling Unimplemented errors from legacy plugins
- Validating response fields before use
- Displaying plugin info in CLI tools

> **Note**: Consumer-side requirements are marked with _(Consumer)_ below. A tracking issue will be
> created in pulumicost-core to implement these once this spec version is released.

### Functional Requirements

- **FR-001**: The `CostSource` service definition MUST include a `GetPluginInfo` operation.
- **FR-002**: The `GetPluginInfo` response MUST include the plugin's `name`, `version`, `spec_version`, and a list of supported `providers`.
- **FR-003**: The plugin SDK MUST provide a default implementation of `GetPluginInfo` that automatically returns the SDK's compiled Spec Version.
- **FR-004**: The SDK MUST define a constant for the current Spec Version to ensure consistency.
- **FR-005** _(Consumer)_: The Core System MUST attempt to call `GetPluginInfo` upon plugin
  initialization with a timeout of **5 seconds**.
- **FR-006** _(Consumer)_: The Core System MUST handle RPC errors (e.g., Unimplemented) from
  `GetPluginInfo` by defaulting to a safe state (e.g., assuming compatibility or warning user)
  rather than terminating.
- **FR-007**: The `spec_version` returned by the plugin MUST be a valid Semantic Version (vX.Y.Z).
- **FR-008** _(Consumer)_: The Core System MUST validate that `name`, `version`, and `spec_version`
  are present and valid in the `GetPluginInfo` response. If any are missing or invalid, plugin
  initialization MUST fail.

### Key Entities _(include if feature involves data)_

- **Plugin Metadata**: Contains identity and versioning information (Name, Version, Spec Version, Providers).
- **Spec Version**: A version string identifier (e.g., "v0.4.11") indicating the protocol version used at build time.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Core system successfully retrieves metadata (including spec version) from 100% of compliant plugins.
- **SC-002** _(Consumer)_: Consumer tools (e.g., `pulumicost plugin list`) can display the spec
  version for all compliant plugins by calling `GetPluginInfo`.
- **SC-003**: System successfully initializes 100% of legacy plugins (non-compliant) without crashing, logging a compatibility warning instead.
