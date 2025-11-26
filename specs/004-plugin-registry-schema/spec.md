# Feature Specification: Plugin Registry Index JSON Schema

**Feature Branch**: `004-plugin-registry-schema`
**Created**: 2025-11-23
**Status**: Draft
**Input**: GitHub Issue #68 - Add plugin registry index JSON Schema

## Clarifications

### Session 2025-11-23

- Q: Should the schema enforce that the plugin `name` field matches its registry key?
  → A: Schema does NOT enforce match; documented as invalid for consumer validation
  (pulumicost-core test file)
- Q: Should a deprecated plugin without a `deprecation_message` be allowed?
  → A: No, require `deprecation_message` when `deprecated: true` (no existing entries)
- Q: Should schema validate that `max_spec_version` >= `min_spec_version`?
  → A: No, document as invalid for consumer validation (JSON Schema limitation)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Validate Registry Index File (Priority: P1)

A registry contributor creates or modifies a registry index file (`registry.json`) and needs
to validate that it conforms to the schema before submitting changes.

**Why this priority**: This is the core purpose of the schema - enabling validation of
registry entries to ensure consistency and correctness across the ecosystem.

**Independent Test**: Can be fully tested by validating example registry.json files against
the schema and verifying that valid entries pass and invalid entries fail with clear error
messages.

**Acceptance Scenarios**:

1. **Given** a registry.json with all required fields populated correctly,
   **When** validated against the schema,
   **Then** validation passes with no errors
2. **Given** a registry.json with a missing required field (e.g., `min_spec_version`),
   **When** validated against the schema,
   **Then** validation fails with a clear error message indicating the missing field
3. **Given** a registry.json with an invalid field value (e.g., invalid semver pattern),
   **When** validated against the schema,
   **Then** validation fails with a clear error message indicating the invalid value

---

### User Story 2 - Discover Plugin Metadata (Priority: P1)

A PulumiCost user runs `pulumicost plugin install kubecost` and the CLI fetches the registry
index to locate the plugin's repository, supported providers, and minimum spec version.

**Why this priority**: This represents the primary consumer use case that drives the need
for a formal schema - the plugin installation system in pulumicost-core.

**Independent Test**: Can be tested by verifying that the schema captures all metadata
fields required by the plugin installation system to discover and install plugins.

**Acceptance Scenarios**:

1. **Given** a registry index with multiple plugins,
   **When** the user searches for a plugin by name,
   **Then** the system can retrieve the plugin's repository, description, and supported
   providers
2. **Given** a registry entry with `security_level: "official"`,
   **When** the user views plugin details,
   **Then** the trust level is clearly communicated
3. **Given** a plugin with `deprecated: true`,
   **When** listed in search results,
   **Then** the deprecation status and migration message are visible

---

### User Story 3 - Contribute Plugin to Registry (Priority: P2)

A third-party plugin developer wants to add their plugin to the official registry and needs
clear guidance on required fields and format.

**Why this priority**: Enables ecosystem growth by providing a clear specification for
plugin contributions.

**Independent Test**: Can be tested by providing an example registry entry and validating
that it can be created from scratch following only the schema documentation.

**Acceptance Scenarios**:

1. **Given** a new plugin with valid metadata,
   **When** the developer creates a registry entry following the schema,
   **Then** the entry validates successfully
2. **Given** a plugin name with invalid characters (e.g., uppercase, special chars),
   **When** validated,
   **Then** validation fails with pattern error explaining naming requirements
3. **Given** a description shorter than 10 characters,
   **When** validated,
   **Then** validation fails indicating minimum length requirement

---

### User Story 4 - Filter Plugins by Capability (Priority: P3)

A user wants to find plugins that support specific capabilities (e.g., `cost_projection`,
`real_time_data`) or providers (e.g., `aws`, `kubernetes`).

**Why this priority**: Enhances discoverability but depends on basic search working first.

**Independent Test**: Can be tested by creating registry entries with various capabilities
and providers and verifying filtering operations work correctly.

**Acceptance Scenarios**:

1. **Given** registry entries with various `supported_providers`,
   **When** filtering by provider,
   **Then** only matching plugins are returned
2. **Given** registry entries with various `capabilities`,
   **When** filtering by capability,
   **Then** only plugins with that capability are returned
3. **Given** registry entries with `keywords`,
   **When** searching by keyword,
   **Then** matching plugins are discoverable

---

### Edge Cases

- **Name/key mismatch**: When a plugin `name` doesn't match its registry key, the schema
  does not enforce this constraint (JSON Schema limitation). Consuming applications
  (e.g., pulumicost-core) MUST validate this match and reject mismatched entries.
- **Deprecated without message**: Schema requires `deprecation_message` when `deprecated`
  is `true`; validation fails if message is missing.
- **Version constraint violation**: When `max_spec_version` < `min_spec_version`, the
  schema does not enforce this (JSON Schema limitation). Consuming applications MUST
  validate version ordering and reject invalid entries.
- **Multiple providers display**: How plugins with multiple `supported_providers` are
  displayed is a consumer UX concern (CLI formatting); schema validates array correctness.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Schema MUST validate registry index files containing plugin metadata
- **FR-002**: Schema MUST require `schema_version` and `plugins` at the root level
- **FR-003**: Schema MUST require `name`, `description`, `repository`, `author`,
  `supported_providers`, and `min_spec_version` for each registry entry
- **FR-004**: Schema MUST validate plugin name pattern as lowercase alphanumeric with
  hyphens (pattern: `^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)
- **FR-005**: Schema MUST validate repository format as `owner/repo` pattern
- **FR-006**: Schema MUST validate version strings as semantic versioning format
  (`^\d+\.\d+\.\d+$`)
- **FR-007**: Schema MUST restrict `supported_providers` to enum values: `aws`, `azure`,
  `gcp`, `kubernetes`, `custom`
- **FR-008**: Schema MUST restrict `capabilities` to enum values matching registry.proto:
  `cost_retrieval`, `cost_projection`, `pricing_specs`, `real_time_data`, `historical_data`,
  `filtering`, `aggregation`, `tagging`, `recommendations`, `anomaly_detection`,
  `forecasting`, `budgets`, `alerts`, `custom`
- **FR-009**: Schema MUST restrict `security_level` to enum values matching registry.proto:
  `untrusted`, `community`, `verified`, `official`
- **FR-010**: Schema MUST validate description length between 10 and 500 characters
- **FR-011**: Schema MUST validate keywords array with max 10 items, each max 30 characters
- **FR-012**: Schema MUST support optional fields: `license`, `homepage`,
  `max_spec_version`, `keywords`, `deprecated`, `deprecation_message`
- **FR-013**: Schema MUST disallow additional properties to prevent schema drift
- **FR-014**: Schema MUST require `deprecation_message` when `deprecated` is `true`
  (using `dependentRequired` or `if/then`)

### Key Entities

- **RegistryIndex**: Top-level container with `schema_version` and `plugins` map
- **RegistryEntry**: Individual plugin metadata including name, description, repository,
  author, providers, capabilities, security level, and version requirements
- **Supported Providers**: Enumerated cloud platforms the plugin supports
- **Capabilities**: Plugin features aligned with gRPC service capabilities
- **Security Level**: Trust tier indicating verification status

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Valid registry entries pass schema validation with no errors
- **SC-002**: Invalid registry entries fail validation with clear, actionable error
  messages identifying the specific field and issue
- **SC-003**: Schema field patterns and enums match existing registry.proto definitions
  exactly
- **SC-004**: Example registry file demonstrates all supported field types and validates
  successfully
- **SC-005**: Validation can be performed using standard JSON Schema tools (AJV, etc.)
  without custom code

## Assumptions

- The registry index is consumed by pulumicost-core's plugin installation system
  (rshade/pulumicost-core#163)
- Plugin names must be unique within the registry
- The schema follows JSON Schema draft 2020-12 specification
- Third-party registries will use this same schema for compatibility
- The `capabilities` enum aligns with PluginInfo.capabilities string values in
  registry.proto
- SPDX license identifiers are used but not validated by the schema (validation would
  require external reference)
- Homepage URLs use URI format validation

## Out of Scope

- Version constraint syntax documentation (separate documentation file)
- Go SDK validation helpers for the registry schema
- Automated registry submission workflow
- Plugin signature verification logic
- Registry synchronization between instances
