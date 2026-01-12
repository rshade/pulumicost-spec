# Feature Specification: Add ARN to GetActualCostRequest

**Feature Branch**: `018-proto-add-arn`  
**Created**: 2025-12-14
**Status**: Implemented  
**Input**: See [GitHub Issue #157](https://github.com/rshade/finfocus-spec/issues/157) - Add
dedicated `arn` field to `GetActualCostRequest` for canonical cloud identifier support.

## User Scenarios & Testing

### User Story 1 - Receive Canonical Identifier in Cost Requests (Priority: P1)

As a cloud provider plugin developer (e.g., for AWS), I need to receive the resource's canonical
identifier (ARN) in the cost request so that I can accurately query the cloud provider's cost API
without ambiguity or missing context (like region/account).

**Why this priority**: P1. This is the core objective. Without this, plugins cannot reliably
identify resources in systems like AWS Cost Explorer, where simple IDs (like `i-123`) are
insufficient or ambiguous.

**Independent Test**: Can be fully tested by generating a `GetActualCostRequest` with the `arn`
field populated and verifying that a consumer (mock plugin) can read and extract this value.

**Acceptance Scenarios**:

1. **Given** a `GetActualCostRequest` object, **When** the `arn` field is populated with a valid
   AWS ARN, **Then** the plugin can access this specific string distinct from the `resource_id`.
2. **Given** a legacy `GetActualCostRequest` (without `arn`), **When** processed by the plugin,
   **Then** the `arn` field is empty/null, and the plugin falls back to existing logic using
   `resource_id`.

### Edge Cases

- **Empty ARN**: If the Core does not have an ARN for a resource, the field remains empty. The
  plugin must handle this gracefully (e.g., by attempting to derive it or falling back to
  `resource_id`).
- **Mismatched IDs**: If `resource_id` and `arn` point to potentially different things (e.g.,
  logical name vs physical ARN), the plugin treats `arn` as the source of truth for cloud API
  queries.

## Requirements

### Functional Requirements

- **FR-001**: The `GetActualCostRequest` message definition MUST include a dedicated field for a
  canonical cloud identifier.
- **FR-002**: This new field (suggested name: `arn`) MUST accept string values representing global
  identifiers (e.g., AWS ARNs, Azure Resource IDs).
- **FR-003**: The new field MUST be optional to ensure backward compatibility with existing clients
  and plugins that do not populate or read it.
- **FR-004**: The field MUST be distinct from the existing `resource_id` field to allow separation
  of Pulumi-internal identifiers from cloud-provider identifiers.

### Key Entities

- **GetActualCostRequest**: The data contract passed from the Core to a Plugin to request cost data
  for a specific resource. It now includes:
  - `resource_id`: The primary identifier known to Pulumi (often a logical name or URN).
  - `arn`: The definitive, globally unique identifier used by the cloud provider.

## Success Criteria

### Measurable Outcomes

- **SC-001**: The Cost Source Interface definition contains the new `arn` field.
- **SC-002**: The interface definition compilation process completes successfully without errors,
  producing updated language-specific bindings.
- **SC-003**: Backward compatibility is maintained: existing client implementations that initialize
  `GetActualCostRequest` without the `arn` field continue to compile and run without modification.
