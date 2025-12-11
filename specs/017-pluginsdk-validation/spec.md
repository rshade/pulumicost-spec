# Feature Specification: PluginSDK Request Validation Helpers

**Feature Branch**: `017-pluginsdk-validation`
**Created**: 2025-12-10
**Status**: Draft
**Input**: User description: "Create pluginsdk/validation.go with shared request
validation helpers that both core and plugins can use"

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Core Pre-flight Validation (Priority: P1)

As a core developer, I need to validate GetProjectedCostRequest messages before
sending them to plugins so that I can catch configuration errors early and
provide actionable error messages that guide developers to fix issues.

**Why this priority**: Pre-flight validation in core is the primary use case
that prevents invalid requests from ever reaching plugins, reducing debugging
time and improving developer experience.

**Independent Test**: Can be fully tested by creating invalid requests with
missing fields and verifying validation returns specific, actionable error
messages that reference the mapping helper functions.

**Acceptance Scenarios**:

1. **Given** a GetProjectedCostRequest with nil Resource,
   **When** ValidateProjectedCostRequest is called,
   **Then** it returns an error "resource is required"
2. **Given** a GetProjectedCostRequest with empty Provider,
   **When** ValidateProjectedCostRequest is called,
   **Then** it returns an error containing "provider is required" with guidance
   on how provider is extracted from resource type
3. **Given** a GetProjectedCostRequest with empty SKU,
   **When** ValidateProjectedCostRequest is called,
   **Then** it returns an error referencing "mapping.ExtractAWSSKU()" or
   equivalent helper

---

### User Story 2 - Plugin Defense-in-Depth Validation (Priority: P2)

As a plugin developer, I need to validate incoming requests using the same
validation logic as core so that my plugin can reject invalid requests with
consistent error messages, even if core validation was bypassed.

**Why this priority**: Defense-in-depth ensures robustness even when core and
plugin versions mismatch or when core validation is disabled for testing.

**Independent Test**: Can be fully tested by calling plugin RPC methods with
invalid requests and verifying they return InvalidArgument gRPC status codes
with validation error messages.

**Acceptance Scenarios**:

1. **Given** a plugin using ValidateProjectedCostRequest,
   **When** it receives a request with empty resource_type,
   **Then** it returns InvalidArgument gRPC error with "resource_type is
   required"
2. **Given** a plugin using ValidateActualCostRequest,
   **When** it receives a request with end time before start time,
   **Then** it returns InvalidArgument gRPC error with "end time must be after
   start time"

---

### User Story 3 - Actionable Error Messages (Priority: P3)

As a developer debugging integration issues, I need validation error messages
that tell me exactly how to fix the problem so that I can quickly resolve
configuration issues without consulting documentation.

**Why this priority**: Developer experience improvement that reduces
time-to-resolution for common integration mistakes.

**Independent Test**: Can be fully tested by triggering each validation error
and verifying the error message contains specific function names or property
keys to check.

**Acceptance Scenarios**:

1. **Given** a request with missing SKU for AWS resources,
   **When** validation fails,
   **Then** the error message references "mapping.ExtractAWSSKU()" and lists
   the property keys checked (instanceType, type, volumeType)
2. **Given** a request with missing region for AWS resources,
   **When** validation fails,
   **Then** the error message references "mapping.ExtractAWSRegion()" and
   mentions availabilityZone derivation

---

### Edge Cases

- What happens when the request is nil? Returns error "request is nil"
- What happens when Resource is non-nil but all fields are empty? Validates
  each field in order, returning first error encountered
- How does validation handle zero timestamps vs nil timestamps? Nil timestamps
  fail validation; zero timestamps (Unix epoch) are technically valid but may
  warn
- What if start and end times are equal? Valid - represents an instant query

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide `ValidateProjectedCostRequest()` function
  that validates the request before sending to plugin
- **FR-002**: System MUST provide `ValidateActualCostRequest()` function that
  validates actual cost requests
- **FR-003**: Validation functions MUST return nil when the request is valid
- **FR-004**: Validation functions MUST return descriptive errors with guidance
  on how to fix the issue
- **FR-005**: Error messages MUST reference specific mapping helper function
  names when applicable (e.g., "use mapping.ExtractAWSSKU()")
- **FR-006**: Validation MUST check fields in a consistent order: nil request,
  nil nested message, required string fields, semantic validations
- **FR-007**: ValidateProjectedCostRequest MUST validate: non-nil request,
  non-nil Resource, non-empty provider, non-empty resource_type
- **FR-008**: ValidateActualCostRequest MUST validate: non-nil request,
  non-empty resource_id, non-nil start time, non-nil end time, end time after
  start time

### Key Entities _(include if feature involves data)_

- **GetProjectedCostRequest**: Proto message containing ResourceDescriptor for
  cost projection; key fields: Resource.Provider, Resource.ResourceType,
  Resource.Sku, Resource.Region
- **GetActualCostRequest**: Proto message for historical cost queries; key
  fields: resource_id, start (Timestamp), end (Timestamp), tags
- **ResourceDescriptor**: Nested message describing the cloud resource;
  contains provider, resource_type, sku, region, tags

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Validation functions execute in under 100 nanoseconds with zero
  heap allocations (matching existing pluginsdk performance standards)
- **SC-002**: All validation error messages include actionable guidance
  (function name or property keys to check)
- **SC-003**: 100% test coverage on all validation functions with table-driven
  tests covering valid and invalid cases
- **SC-004**: Integration with existing code is non-breaking - validation is
  opt-in, not automatically applied

## Assumptions

- The proto messages are imported from the existing generated proto package
- The mapping package helpers exist and are documented - error messages
  reference them by name
- Validation is stateless and does not require any external dependencies
  (network, filesystem, etc.)
- Region and SKU validation is provider-agnostic at this level -
  provider-specific validation is handled by the mapping helpers
- Simple error formatting is acceptable (no need for custom error types at
  this stage)

## Out of Scope

- Response validation (ValidateProjectedCostResponse,
  ValidateActualCostResponse) - covered in future issue
- Provider-specific validation logic (e.g., validating AWS region names) -
  handled by mapping package
- Request transformation or correction - validation only reports errors, does
  not fix them
- gRPC interceptor integration - calling code decides whether to use validation
- Custom error types or error codes - simple error messages are sufficient
