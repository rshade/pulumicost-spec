# Feature Specification: Plugin Capability Dry Run Mode

**Feature Branch**: `032-plugin-dry-run`
**Created**: 2025-12-31
**Status**: Draft
**Input**: User description: "Allow the host to query a plugin for its mapping logic
('What fields would you return for this resource?') without performing data retrieval.
Useful for debugging and validation."

## Clarifications

### Session 2025-12-31

- Q: Dry-run API design pattern? → A: Hybrid approach - add dedicated `DryRun` RPC for
  standalone discovery AND add `dry_run` flag to existing GetActualCost/GetProjectedCost
  RPCs for inline validation during cost calls.
- Q: Partial resource type support handling? → A: Return all known FOCUS fields with
  explicit "supported/unsupported" status per field, enabling hosts to make informed
  decisions about partial coverage.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Discover Plugin Field Mappings (Priority: P1)

As a host application operator, I want to query a plugin to discover what FOCUS fields it
would populate for a given resource type, without the plugin actually fetching cost data
from external sources. This allows me to validate plugin configurations and understand
mapping behavior before running real cost queries.

**Why this priority**: This is the core value proposition - understanding plugin
capabilities without incurring data retrieval costs or latency. Essential for debugging
and integration validation.

**Independent Test**: Can be fully tested by sending a dry-run request for a known
resource type and verifying the response contains expected field mappings without any
external API calls being made.

**Acceptance Scenarios**:

1. **Given** a configured plugin supporting AWS EC2 instances, **When** host sends a
   dry-run request for resource type "aws:ec2:Instance", **Then** plugin returns the list
   of FOCUS fields it would populate (e.g., ServiceCategory, ChargeType, PricingUnit)
   without making AWS API calls.
2. **Given** a plugin that does not support a resource type, **When** host sends a
   dry-run request for that unsupported resource type, **Then** plugin returns an
   appropriate indication that the resource is not supported.
3. **Given** a plugin in dry-run mode, **When** the request is processed, **Then** no
   external data sources are contacted and response time is under 100ms.

---

### User Story 2 - Validate Plugin Configuration (Priority: P2)

As a platform developer integrating multiple cost plugins, I want to use dry-run mode to
validate that my plugin configuration is correct before deploying to production. This
helps catch configuration errors early without waiting for actual cost data retrieval.

**Why this priority**: Configuration validation reduces integration friction and prevents
production issues. Important for developer experience but not blocking core functionality.

**Independent Test**: Can be tested by intentionally misconfiguring a plugin and
verifying dry-run mode returns clear error information about the configuration issue.

**Acceptance Scenarios**:

1. **Given** a plugin with valid configuration, **When** host sends a dry-run request,
   **Then** plugin confirms configuration is valid and returns field mapping information.
2. **Given** a plugin with missing required configuration, **When** host sends a dry-run
   request, **Then** plugin returns a clear error indicating what configuration is missing.

---

### User Story 3 - Compare Plugin Capabilities (Priority: P3)

As a platform administrator evaluating multiple plugins for the same provider, I want to
compare what fields each plugin would populate for the same resource type. This helps me
choose the most comprehensive plugin for my needs.

**Why this priority**: Comparison functionality builds on the core dry-run capability.
Valuable for plugin selection but relies on P1 being complete.

**Independent Test**: Can be tested by querying multiple plugins for the same resource
type and comparing the returned field lists.

**Acceptance Scenarios**:

1. **Given** two plugins supporting the same provider, **When** host sends dry-run
   requests to both for the same resource type, **Then** host can compare the field
   mappings returned by each plugin.
2. **Given** a plugin with optional field mappings based on configuration, **When** host
   sends a dry-run request, **Then** plugin indicates which fields are conditionally
   available and under what conditions.

---

### Edge Cases

- **Partial support**: When a plugin partially supports a resource type, it returns all
  known FOCUS fields with explicit supported/unsupported status per field.
- **Uninitialized plugin**: Plugin returns an error indicating initialization is required
  before dry-run queries can be processed.
- **Runtime-dependent mappings**: Fields that depend on runtime data (e.g., region-specific
  pricing tiers) are marked as "conditional" with a description of the dependency.
- **Dynamic/computed fields**: Marked as "dynamic" in the field status, indicating the
  value cannot be determined without actual data retrieval.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide a mechanism for hosts to request field mapping
  information without triggering data retrieval operations.
- **FR-002**: System MUST clearly indicate which FOCUS fields a plugin would populate
  for a given resource type in dry-run mode.
- **FR-003**: System MUST return dry-run responses significantly faster than actual cost
  queries (no external API calls).
- **FR-004**: System MUST indicate when a requested resource type is not supported by
  the plugin.
- **FR-005**: System MUST report field support status using a defined enum: SUPPORTED
  (always populated), UNSUPPORTED (not available), CONDITIONAL (depends on resource
  configuration), DYNAMIC (requires runtime data).
- **FR-006**: System MUST validate plugin configuration as part of dry-run processing
  and report configuration errors.
- **FR-007**: System MUST NOT make any external network calls to cost data sources when
  operating in dry-run mode.
- **FR-008**: System MUST provide a dedicated `DryRun` RPC for standalone field mapping
  discovery.
- **FR-009**: System MUST support a `dry_run` flag on existing GetActualCost and
  GetProjectedCost RPCs to enable inline validation during cost calls.

### Key Entities

- **DryRunRequest**: Contains `ResourceDescriptor` (reuses existing message) to identify
  the resource type being queried, plus optional simulation parameters.
- **DryRunResponse**: Contains list of `FieldMapping` entries for all known FOCUS fields,
  configuration validation status, and any errors encountered.
- **FieldMapping**: Represents a single FOCUS field with:
  - `field_name`: FOCUS field identifier (e.g., "ServiceCategory", "ChargeType")
  - `support_status`: Enum of SUPPORTED, UNSUPPORTED, CONDITIONAL, DYNAMIC
  - `condition_description`: Human-readable explanation when status is CONDITIONAL/DYNAMIC
  - `expected_type`: Data type the field would contain (string, double, timestamp, etc.)

### Integration with Existing RPCs

- **Supports RPC**: The existing `SupportsResponse.capabilities` map can indicate dry-run
  support via `{"dry_run": true}`. This allows hosts to check capability before calling.
- **GetPluginInfo RPC**: Can report dry-run support in `metadata` for diagnostic tools.
- **GetActualCost/GetProjectedCost RPCs**: Gain optional `dry_run` bool field; when true,
  return `DryRunResponse` data instead of cost data (via response union or separate field).

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Dry-run requests complete in under 100ms since no external data retrieval
  occurs.
- **SC-002**: 100% of supported resource types return accurate field mapping information
  matching actual behavior.
- **SC-003**: Plugin developers can validate configurations without provisioning test
  cloud resources.
- **SC-004**: Platform operators can compare plugin capabilities in under 5 minutes for
  any resource type.
- **SC-005**: Configuration errors are detected and reported clearly, reducing
  integration debugging time by 50% compared to discovering errors during actual cost
  queries.

## Assumptions

- Plugins have static knowledge of their field mapping logic that can be introspected
  without making external calls.
- Field mappings are deterministic for a given resource type (though some fields may be
  conditionally populated).
- The dry-run mode is a transport-layer feature that does not require changes to cost
  calculation logic.
- Plugins can distinguish between "supported with full mapping" and "supported with
  partial mapping" for resource types.

## Constraints

- Dry-run mode explicitly stays within the transport layer and does not calculate costs.
- Existing GetActualCost and GetProjectedCost RPCs gain an optional `dry_run` flag; when
  set, they return field mappings instead of performing data retrieval. The flag defaults
  to false for backward compatibility.
- A new dedicated `DryRun` RPC provides standalone discovery without modifying cost call
  semantics.
- Must be backward compatible with existing plugins (plugins not implementing dry-run
  should gracefully indicate lack of support).

## Dependencies

- Existing CostSource gRPC service definition (proto/pulumicost/v1/costsource.proto)
- FOCUS field definitions in the SDK
- Plugin registry for resource type identification
