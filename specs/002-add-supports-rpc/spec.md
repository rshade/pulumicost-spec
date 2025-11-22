# Feature Specification: Add Supports() RPC Method to CostSourceService

**Feature Branch**: `002-add-supports-rpc`
**Created**: 2025-11-20
**Status**: Already Implemented (Verification Only)
**Input**: GitHub Issue #64 - Add Supports() RPC method to CostSourceService

**Note**: Analysis revealed Supports() RPC was fully implemented in v0.1.0. This spec
documents requirements for verification purposes. Actual work needed is in pulumicost-core.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Query Plugin Capabilities (Priority: P1)

As a PulumiCost client, I need to query a plugin to determine if it supports a specific
resource type and region before making cost requests, so I can provide meaningful feedback
to users and avoid unnecessary API calls.

**Why this priority**: This is the core functionality that enables the entire Supports() RPC
capability. Without this, clients cannot determine plugin capabilities via gRPC.

**Independent Test**: Can be fully tested by calling the Supports() RPC method with a
ResourceDescriptor and verifying the response indicates support status and optional reason.

**Acceptance Scenarios**:

1. **Given** a plugin that supports AWS EC2 in us-east-1,
   **When** client calls Supports() with a matching ResourceDescriptor,
   **Then** response returns supported=true
2. **Given** a plugin that does not support a specific resource type,
   **When** client calls Supports() with that ResourceDescriptor,
   **Then** response returns supported=false with a reason explaining why
3. **Given** a plugin that does not support a specific region,
   **When** client calls Supports() with that ResourceDescriptor,
   **Then** response returns supported=false with a reason indicating the unsupported region

---

### User Story 2 - Graceful Capability Discovery (Priority: P2)

As a plugin developer, I need to implement the Supports() RPC method in my plugin so that
clients can discover what resources my plugin can price, enabling better plugin selection
and error messages.

**Why this priority**: Plugin developers need a clear contract to implement. This enables
the ecosystem to grow with capability-aware plugins.

**Independent Test**: Can be tested by implementing a plugin with Supports() method and
verifying it responds correctly via gRPC test harness.

**Acceptance Scenarios**:

1. **Given** a plugin implementation with Supports() method,
   **When** the plugin uses the generated CostSourceServiceServer interface,
   **Then** the Supports() method signature is available for implementation
2. **Given** a plugin with partial resource support,
   **When** Supports() is called with various ResourceDescriptors,
   **Then** each returns appropriate supported status and reason

---

### Edge Cases

- What happens when Supports() is called with a nil or empty ResourceDescriptor?
- How does the system handle network errors during Supports() RPC calls?
- What happens when a plugin returns an invalid response (missing fields)?
- How does the system behave when Supports() takes too long to respond?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: CostSourceService proto definition MUST include a Supports RPC method
- **FR-002**: Supports RPC MUST accept a SupportsRequest message containing a
  ResourceDescriptor
- **FR-003**: Supports RPC MUST return a SupportsResponse message with supported (boolean)
  and reason (string) fields
- **FR-004**: Generated Go SDK code MUST include Supports method in CostSourceServiceClient
  interface
- **FR-005**: Generated Go SDK code MUST include Supports method signature in
  CostSourceServiceServer interface
- **FR-006**: UnimplementedCostSourceServiceServer MUST include a default Supports
  implementation returning Unimplemented error
- **FR-007**: Proto code generation MUST complete without errors for all supported
  languages (Go)
- **FR-008**: Build process MUST succeed with the regenerated proto code
- **FR-009**: CHANGELOG MUST be updated to document the new RPC method

### Key Entities

- **ResourceDescriptor**: Describes a cloud resource for capability checking (resource type,
  provider, region, properties)
- **SupportsRequest**: Request message containing a ResourceDescriptor to check for support
- **SupportsResponse**: Response message indicating whether the resource is supported and
  optionally why not
- **CostSourceService**: The gRPC service definition that will include the new Supports RPC
  method

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Proto linting passes with no errors after adding Supports RPC
- **SC-002**: Generated Go code compiles successfully without errors
- **SC-003**: All existing tests continue to pass after proto regeneration
- **SC-004**: Supports method is available in CostSourceServiceClient and
  CostSourceServiceServer interfaces
- **SC-005**: Plugin developers can implement Supports() and have it called via gRPC by
  clients

## Assumptions

- SupportsRequest and SupportsResponse message types already exist in the proto definition
  (per issue notes)
- ResourceDescriptor message type already exists in the proto definition
- Only Go SDK needs to be regenerated (no other language SDKs currently maintained)
- This is a backward-compatible change (MINOR version bump appropriate)
- No changes to existing message definitions required

## Dependencies

- Proto compiler (buf) must be available
- Go toolchain must be installed
- Existing SupportsRequest and SupportsResponse messages must be correctly defined

## Out of Scope

- Changes to SupportsRequest or SupportsResponse message definitions
- Implementation of the handler in pulumicost-core pluginsdk (separate issue)
- Plugin-specific implementation of the Supports() method
- Performance optimization of the Supports() RPC
