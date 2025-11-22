# Feature Specification: GetPricingSpec RPC

**Feature Branch**: `003-getpricingspec`
**Created**: 2025-11-22
**Status**: Draft
**Input**: GitHub Issue #62 - Add GetPricingSpec() RPC to CostSourceService for pricing transparency

## User Scenarios & Testing

### User Story 1 - Get Flat-Rate Pricing Information (Priority: P1)

A cost estimation user wants to retrieve detailed pricing information for a resource with
flat-rate billing (like EC2 instances) to understand the exact rate, unit, and pricing
assumptions being used.

**Why this priority**: Core functionality - enables users to see transparent pricing breakdowns
for the most common billing model. This is the fundamental use case that unblocks the
aws-public plugin and provides immediate value for cost verification.

**Independent Test**: Can be fully tested by calling GetPricingSpec for an EC2 instance and
verifying the response contains rate_per_unit, currency, billing_mode, unit, description,
and assumptions fields. Delivers transparent pricing visibility.

**Acceptance Scenarios**:

1. **Given** a CostSource plugin implements GetPricingSpec, **When** the user requests
   pricing for an EC2 t3.micro in us-east-1, **Then** the response contains
   rate_per_unit=0.0104, currency=USD, billing_mode=per_hour, unit=hour, and a list of
   pricing assumptions.
2. **Given** a CostSource plugin implements GetPricingSpec, **When** the user requests
   pricing for an EBS gp3 volume, **Then** the response contains rate_per_unit=0.08,
   currency=USD, billing_mode=per_gb_month, unit=GB-month.
3. **Given** a CostSource plugin implements GetPricingSpec, **When** the response is
   returned, **Then** all mandatory fields (rate_per_unit, currency, billing_mode, unit,
   description) are populated.

---

### User Story 2 - Get Tiered Pricing Information (Priority: P2)

A cost estimation user wants to retrieve detailed pricing information for resources with
tiered billing (like S3 storage) to understand how costs are calculated at different usage
levels.

**Why this priority**: Extends transparency to more complex pricing models. S3 and similar
services use tiered pricing, which is common but requires additional response structure.

**Independent Test**: Can be fully tested by calling GetPricingSpec for an S3 bucket and
verifying the response contains pricing_tiers with min_quantity, max_quantity, rate_per_unit,
and description for each tier.

**Acceptance Scenarios**:

1. **Given** a CostSource plugin implements tiered pricing, **When** the user requests
   pricing for S3 Standard storage, **Then** the response contains billing_mode=tiered
   and a pricing_tiers array with at least two tiers.
2. **Given** a CostSource plugin returns tiered pricing, **When** the response contains
   pricing_tiers, **Then** each tier has min_quantity, max_quantity (0 for unlimited),
   rate_per_unit, and description fields.
3. **Given** pricing_tiers is populated, **When** the user examines the tiers, **Then**
   they are ordered by min_quantity ascending and tiers are contiguous without gaps.

---

### User Story 3 - Handle Not-Implemented Resources (Priority: P3)

A cost estimation user calls GetPricingSpec for a resource type that the plugin does not
support, and receives a graceful response indicating the limitation.

**Why this priority**: Ensures backwards compatibility and graceful degradation. Plugins
should not fail when pricing info is unavailable; they should communicate limitations
clearly.

**Independent Test**: Can be fully tested by calling GetPricingSpec for an unsupported
resource type and verifying the response contains billing_mode=not_implemented and helpful
assumptions explaining the limitation.

**Acceptance Scenarios**:

1. **Given** a CostSource plugin does not support pricing for Lambda, **When** the user
   requests GetPricingSpec for a Lambda function, **Then** the response contains
   billing_mode=not_implemented, unit=unknown, and rate_per_unit=0.
2. **Given** pricing is not implemented for a resource, **When** the response is returned,
   **Then** the assumptions field contains at least one entry explaining why pricing is
   unavailable.

---

### Edge Cases

- What happens when ResourceDescriptor is missing required fields (provider, resource_type)?
  - Return gRPC InvalidArgument status with descriptive error message
- How does the system handle invalid region or SKU combinations?
  - Return gRPC NotFound status with descriptive error message
- What happens when a plugin returns both pricing_tiers and a non-zero rate_per_unit?
  - Both are allowed: tiered billing uses pricing_tiers, flat billing uses rate_per_unit
- How does the system handle unknown currencies or currency conversion requirements?
  - Out of scope; currency field is informational only (no validation)

## Requirements

### Functional Requirements

- **FR-001**: System MUST provide a GetPricingSpec RPC method in the CostSourceService
  interface
- **FR-002**: GetPricingSpec MUST accept a GetPricingSpecRequest containing a
  ResourceDescriptor
- **FR-003**: GetPricingSpec MUST return a GetPricingSpecResponse with rate_per_unit
  (double), currency (string), billing_mode (string), unit (string), description (string),
  assumptions (repeated string), and pricing_tiers (repeated PricingTier)
- **FR-004**: PricingTier MUST contain min_quantity (double), max_quantity (double),
  rate_per_unit (double), and description (string)
- **FR-005**: System MUST support billing_mode values including: per_hour, per_gb_month,
  tiered, not_implemented
- **FR-006**: System MUST support unit values including: hour, GB-month, request, unknown
- **FR-007**: Plugins MAY choose not to implement GetPricingSpec; the method is optional
- **FR-008**: When pricing is unavailable, system MUST return billing_mode=not_implemented
  and rate_per_unit=0
- **FR-009**: assumptions field MUST be populated with human-readable strings explaining
  pricing derivation
- **FR-010**: GetPricingSpec MUST NOT break existing GetProjectedCost functionality
  (backward compatible)
- **FR-011**: GetPricingSpec MUST return gRPC InvalidArgument status when ResourceDescriptor
  is missing required fields (provider, resource_type)
- **FR-012**: GetPricingSpec MUST return gRPC NotFound status when region or SKU combination
  is invalid or unknown

### Key Entities

- **GetPricingSpecRequest**: Request containing a ResourceDescriptor to query pricing for
- **GetPricingSpecResponse**: Response containing transparent pricing breakdown with rate,
  currency, billing mode, unit, description, assumptions, and optional pricing tiers
- **PricingTier**: Represents one tier in a tiered pricing model with quantity ranges and
  per-unit rate
- **ResourceDescriptor**: Existing entity containing provider, resource_type, sku, region,
  and tags

## Success Criteria

### Measurable Outcomes

- **SC-001**: Cost estimation tools can retrieve structured pricing information for any
  supported resource
- **SC-002**: Users can verify cost calculations by examining the returned rate_per_unit,
  billing_mode, and assumptions
- **SC-003**: Downstream integration tools can programmatically parse all fields in the
  GetPricingSpecResponse
- **SC-004**: Plugin developers can implement GetPricingSpec or gracefully decline by
  returning not_implemented
- **SC-005**: All existing GetProjectedCost calls continue to function without modification
  after adding GetPricingSpec

## Assumptions

- Currency will initially be USD; multi-currency support may be added later
- billing_mode values follow existing patterns in sdk/go/pricing and can be extended
- Plugins use the same ResourceDescriptor structure already defined for GetProjectedCost
- Pricing tiers are exclusive (a quantity falls into exactly one tier)
- max_quantity=0 indicates unlimited (no upper bound for the tier)
- The RPC is synchronous and stateless, consistent with other CostSourceService methods

## Out of Scope

- Currency conversion between different currencies
- Historical pricing queries (this returns current pricing only)
- Real-time pricing updates or pricing change notifications
- Cost optimization recommendations based on pricing data
- Aggregation of pricing across multiple resources

## Clarifications

### Session 2025-11-22

- Q: How should GetPricingSpec handle invalid input (missing fields, invalid region/SKU)?
  → A: Return gRPC status codes (InvalidArgument, NotFound) with descriptive error messages
- Q: How handle both pricing_tiers and non-zero rate_per_unit in response?
  → A: Allow both; tiered billing uses tiers, flat billing uses rate_per_unit
