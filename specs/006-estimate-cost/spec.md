# Feature Specification: "What-If" Cost Estimation API

**Feature Branch**: `006-estimate-cost`
**Created**: 2025-11-24
**Status**: Draft
**Input**: GitHub Issue #79 - Add EstimateCost RPC to CostSource gRPC service for proactive cost management

## Clarifications

### Session 2025-11-24

- Q: When a plugin cannot determine pricing for given attributes (e.g., incomplete configuration,
  unsupported region, or ambiguous pricing tier), what should the API behavior be? → A: Return
  gRPC error with detailed message explaining missing/ambiguous attributes
- Q: When the plugin's pricing source is unavailable or returns errors (network timeout, API rate
  limit, service degradation), how should EstimateCost behave? → A: Return gRPC error immediately
  indicating pricing source unavailability (SDK should not handle retry; plugins and core handle
  retry logic)
- Q: When the resource type format is invalid (e.g., "invalid-format", "aws:ec2:Instance" missing
  module, or "aws/ec2/instance" wrong separator), what should the API response be? → A: Return
  gRPC InvalidArgument error immediately with message explaining expected format
- Q: What observability signals should EstimateCost operations emit for monitoring and debugging?
  → A: Structured logs (request/response/errors), metrics (latency/success rate), distributed
  tracing
- Q: How should the system handle missing or null values in the attributes field of EstimateCost
  requests? → A: Treat null/missing attributes as empty struct and let plugin decide if valid

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Basic Cost Estimation (Priority: P1)

As a developer using FinFocus, I want to get a cost estimate for a resource before
deploying it, so that I can make informed decisions about resource configurations during
development.

**Why this priority**: This is the core value proposition of the feature - enabling proactive
cost management instead of reactive analysis.

**Independent Test**: Can be fully tested by sending an EstimateCost request with valid resource
type and attributes, and verifying a cost estimate is returned. Delivers immediate value by
providing cost visibility before deployment.

**Acceptance Scenarios**:

1. **Given** a valid resource type (e.g., "aws:ec2/instance:Instance") and attributes
   (e.g., instance_type: "t3.micro"), **When** the user calls EstimateCost, **Then** the system
   returns an estimated monthly cost in USD
2. **Given** a valid resource type with multiple pricing attributes, **When** the user calls
   EstimateCost, **Then** the system calculates and returns the combined estimated cost
3. **Given** a valid resource type, **When** the user calls EstimateCost, **Then** the response
   includes both the cost amount and currency

---

### User Story 2 - Configuration Comparison (Priority: P2)

As a developer, I want to compare costs between different resource configurations, so that I can
optimize my infrastructure choices before deployment.

**Why this priority**: Builds on P1 to provide comparative analysis capability, which is
essential for cost optimization workflows.

**Independent Test**: Can be tested by calling EstimateCost multiple times with different
attribute values for the same resource type and comparing the returned costs.

**Acceptance Scenarios**:

1. **Given** the same resource type with different attribute values (e.g., "t3.micro" vs
   "t3.large"), **When** the user calls EstimateCost for each configuration, **Then** the system
   returns different cost estimates reflecting the configuration differences
2. **Given** multiple EstimateCost requests in sequence, **When** each request is processed,
   **Then** each response is consistent and deterministic for the same inputs

---

### User Story 3 - Unsupported Resource Handling (Priority: P3)

As a developer, I want clear feedback when requesting cost estimates for unsupported resources,
so that I understand the system's capabilities and limitations.

**Why this priority**: Essential for user experience but not core functionality - users need to
know when the system cannot provide estimates.

**Independent Test**: Can be tested by sending EstimateCost requests with unsupported resource
types and verifying appropriate error responses.

**Acceptance Scenarios**:

1. **Given** an unsupported resource type, **When** the user calls EstimateCost, **Then** the
   system returns an appropriate error indicating the resource is not supported
2. **Given** a resource type with invalid format (e.g., "invalid-format", "aws:ec2:Instance"
   missing module, "aws/ec2/instance" wrong separator), **When** the user calls EstimateCost,
   **Then** the system returns gRPC InvalidArgument error with message explaining the expected
   provider:module/resource:Type format
3. **Given** a request with null or missing attributes field, **When** the user calls
   EstimateCost, **Then** the system treats it as empty attributes and forwards to plugin for
   validation
4. **Given** a resource type with missing required attributes, **When** the user calls
   EstimateCost, **Then** the system returns an error describing which attributes are missing
5. **Given** a resource type with ambiguous or incomplete pricing attributes, **When** the user
   calls EstimateCost, **Then** the system returns a descriptive error explaining which attributes
   are ambiguous or invalid

---

### Edge Cases

- When the resource type format is invalid (e.g., "invalid-format", "aws:ec2:Instance" missing
  module, "aws/ec2/instance" wrong separator), the system returns a gRPC InvalidArgument error
  immediately with a message explaining the expected format (provider:module/resource:Type)
- When attributes field is null or missing, the system treats it as an empty struct and forwards
  the request to the plugin, which decides if the empty attributes are valid for estimation
- When a plugin cannot determine pricing for the given attributes (e.g., incomplete configuration,
  unsupported region, ambiguous pricing tier), the system returns a gRPC error with a detailed
  message explaining which attributes are missing or ambiguous
- When the plugin's pricing source is unavailable or returns errors (network timeout, API rate
  limit, service degradation), the system returns a gRPC error immediately indicating pricing
  source unavailability. Retry logic is handled by plugins and core, not the SDK
- How does system handle currency conversion if the pricing source uses a different currency?
- What happens when estimated cost is zero (e.g., free tier resources)?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide an EstimateCost RPC method on the CostSource gRPC service
- **FR-002**: System MUST accept a resource type string identifying the Pulumi resource (e.g., "aws:ec2/instance:Instance")
- **FR-003**: System MUST validate resource type format and return gRPC InvalidArgument error with
  descriptive message if format does not match provider:module/resource:Type pattern
- **FR-004**: System MUST accept structured attributes representing resource input properties
- **FR-005**: System MUST treat null or missing attributes field as empty struct and forward to
  plugin for validation (plugin determines if empty attributes are valid)
- **FR-006**: System MUST return an estimated monthly cost for the specified resource configuration
- **FR-007**: System MUST return the currency of the estimated cost (e.g., "USD")
- **FR-008**: System MUST return appropriate gRPC errors for unsupported resource types
- **FR-009**: System MUST return appropriate gRPC errors when required attributes are missing
- **FR-010**: System MUST return descriptive gRPC errors when pricing cannot be determined for
  given attributes, explaining which attributes are missing, ambiguous, or invalid
- **FR-011**: System MUST return deterministic results for identical inputs
- **FR-012**: System MUST support the same resource types as reported by the Supports RPC
- **FR-013**: System MUST handle zero-cost estimates (e.g., free tier) by returning a valid response with zero cost
- **FR-014**: System MUST return gRPC errors immediately when pricing source is unavailable or
  returns errors, without implementing retry logic in the SDK (retry handled by plugins/core)

### Non-Functional Requirements

#### Observability

- **NFR-001**: System MUST emit structured logs for EstimateCost requests, responses, and errors
  with relevant context (resource type, error details)
- **NFR-002**: System MUST track metrics including request latency, success rate, and error rate
  for EstimateCost operations
- **NFR-003**: System MUST support distributed tracing to enable end-to-end visibility through
  cost estimation flows

### Key Entities

- **EstimateCostRequest**: Represents a request to estimate costs, containing the resource type
  string and structured attributes map
- **EstimateCostResponse**: Represents the cost estimation result, containing the monthly cost
  decimal value and currency string
- **Resource Type**: String identifier for Pulumi resources in format "provider:module/resource:Type"
- **Attributes**: Structured key-value data representing resource input properties (mirrors Pulumi resource declaration)

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Users can obtain a cost estimate for any supported resource in a single API call
- **SC-002**: Cost estimates are returned within 500ms for standard resource types
- **SC-003**: System correctly estimates costs for at least 3 major cloud providers (AWS, Azure, GCP)
- **SC-004**: 100% of EstimateCost requests return either a valid estimate or a descriptive error
- **SC-005**: Users can compare at least 5 different configurations per resource type to find optimal cost
- **SC-006**: Operations emit observable signals (logs, metrics, traces) enabling operators to
  monitor API health and debug failures

## Assumptions

- The pricing data used for estimates is available in the plugin's pricing source
- Estimated costs are based on list prices and do not include discounts, reserved instances, or negotiated rates
- Monthly cost assumes standard 730 hours/month for hourly-billed resources
- Currency returned matches the plugin's configured pricing source (typically USD)
- The `google.protobuf.Struct` type is sufficient for representing all resource attribute types
- Decimal precision for costs follows existing patterns in GetActualCost and GetProjectedCost responses
- Retry logic for pricing source failures is the responsibility of plugins and core components, not
  the SDK

## Scope Boundaries

### In Scope

- Single resource cost estimation
- Monthly cost calculation
- Support for resources already supported by the Supports RPC
- Standard attribute handling via protobuf Struct

### Out of Scope

- Multi-resource or stack-level cost estimation
- Reserved instance or savings plan pricing
- Real-time pricing updates
- Currency conversion between different currencies
- Cost estimation for custom/unsupported resource types
- Historical pricing or price trend analysis
