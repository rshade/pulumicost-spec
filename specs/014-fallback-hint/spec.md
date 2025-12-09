# Feature Specification: FallbackHint Enum for Plugin Orchestration

**Feature Branch**: `001-fallback-hint`
**Created**: 2025-12-05
**Status**: Draft
**Input**: GitHub Issue #124 - Add FallbackHint enum to GetActualCostResponse

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Returns Actual Cost Data (Priority: P1)

A cost source plugin (e.g., aws-ce) queries its data source and successfully retrieves
actual cost data for a resource. The plugin returns the data along with a signal
indicating no fallback is needed.

**Why this priority**: This is the primary happy path - most cost queries will return
data successfully. Core must correctly interpret that no fallback is needed when data
is present.

**Independent Test**: Can be fully tested by calling GetActualCost on a plugin with
available data and verifying core does not invoke fallback plugins.

**Acceptance Scenarios**:

1. **Given** a plugin with actual cost data for a resource, **When** GetActualCost is
   called, **Then** response contains results and indicates no fallback needed
2. **Given** core receives a response with data and no-fallback signal, **When**
   processing the response, **Then** core returns data to user without trying
   fallback plugins

---

### User Story 2 - Plugin Has No Data, Recommends Fallback (Priority: P1)

A cost source plugin queries its data source but finds no data (e.g., new resource
not yet in billing system). The plugin returns an empty response with a signal
recommending fallback to another plugin that may provide estimates.

**Why this priority**: This enables the core use case of graceful degradation from
actual costs to estimated costs without conflating "no data" with "error".

**Independent Test**: Can be fully tested by calling GetActualCost on a plugin with
no data and verifying core attempts fallback plugin.

**Acceptance Scenarios**:

1. **Given** a plugin with no actual cost data for a resource, **When** GetActualCost
   is called, **Then** response is empty but includes fallback-recommended signal
2. **Given** core receives empty response with fallback-recommended signal, **When**
   processing the response, **Then** core attempts to call a fallback plugin
   (e.g., aws-public for estimates)

---

### User Story 3 - Plugin Cannot Handle Request Type (Priority: P2)

A plugin receives a request for a resource type or operation outside its domain.
The plugin returns an empty response with a signal requiring fallback, indicating
the request must be handled by another plugin.

**Why this priority**: Important for proper plugin specialization - allows
cost-explorer plugins to cleanly delegate to public-pricing plugins for
unsupported resource types.

**Independent Test**: Can be fully tested by calling GetActualCost on a plugin with
an unsupported resource type and verifying core must attempt fallback.

**Acceptance Scenarios**:

1. **Given** a plugin that does not support a resource type, **When** GetActualCost
   is called for that resource, **Then** response is empty with fallback-required
   signal
2. **Given** core receives fallback-required signal, **When** processing the response,
   **Then** core must try fallback plugins (not optional)

---

### User Story 4 - Plugin Returns Error (Priority: P2)

A plugin encounters a true failure (API error, credentials invalid, network timeout).
The plugin returns an error, which core handles through its error path, not the
fallback path.

**Why this priority**: Critical to distinguish between "no data available, try
fallback" and "system failure" to maintain proper error handling semantics.

**Independent Test**: Can be fully tested by simulating an API failure and verifying
core routes to error handling, not fallback.

**Acceptance Scenarios**:

1. **Given** a plugin that encounters an API failure, **When** GetActualCost is
   called, **Then** response is a gRPC error (not a success response with hint)
2. **Given** core receives an error response, **When** processing the response,
   **Then** core uses error handling path (not fallback path)

---

### User Story 5 - Backwards Compatibility with Legacy Plugins (Priority: P1)

Existing plugins that do not set the fallback hint field continue to work correctly.
The default/unspecified value is treated as "no fallback" to maintain current
behavior.

**Why this priority**: Essential for non-breaking deployment - existing plugins must
work without modification.

**Independent Test**: Can be fully tested by calling a legacy plugin (without hint
field) and verifying core treats response as no-fallback.

**Acceptance Scenarios**:

1. **Given** a legacy plugin that does not set the fallback hint, **When**
   GetActualCost is called, **Then** response has unspecified/default hint value
2. **Given** core receives response with unspecified hint, **When** processing the
   response, **Then** core treats it as "no fallback" (backwards compatible)

---

### Edge Cases

- What happens when a plugin returns data AND sets fallback-recommended? Core should
  use the data and ignore the hint (data takes precedence).
- How does system handle unrecognized hint values (future extensibility)? Core should
  treat unrecognized values as UNSPECIFIED for forward compatibility.
- What happens if all fallback plugins also recommend fallback? Core should stop and
  return empty result (avoid infinite delegation).
- How to distinguish "no data found" from "zero cost data"? Empty results array `[]`
  indicates no billing records exist (use `FALLBACK_RECOMMENDED`). Results containing
  entries with `cost: 0.00` indicate legitimate zero-cost resources like free tier
  (use `FALLBACK_NONE`). Plugins wrapping SaaS vendors that return 0.00 for "not
  found" must translate to proper empty-array semantics.

## Requirements _(mandatory)_

### Functional Requirements

- FR-001: System MUST define a FallbackHint enumeration with four values:
  FALLBACK_HINT_UNSPECIFIED (0), FALLBACK_HINT_NONE (1), FALLBACK_HINT_RECOMMENDED (2),
  FALLBACK_HINT_REQUIRED (3)
- FR-002: System MUST add a fallback_hint field to GetActualCostResponse message
  in the protocol definition
- FR-003: Plugins MUST be able to set the fallback hint when constructing
  responses via SDK helpers
- FR-004: System MUST preserve backwards compatibility - FALLBACK_HINT_UNSPECIFIED must
  be the default value (proto3 default = 0)
- **FR-005**: System MUST generate updated language bindings (Go SDK) from the
  modified protocol definition
- FR-006: System MUST provide SDK helper functions, preferably using the functional
  options pattern, for creating responses with explicit fallback hints
- **FR-007**: Existing SDK helper functions MUST continue to work without modification
  (use default hint)
- **FR-008**: System MUST document the semantics of each hint value for plugin
  developers
- **FR-009**: System MUST include workflow documentation describing core's expected
  behavior for each hint
- **FR-010**: Documentation MUST clarify the distinction between empty results (no
  billing data, use `FALLBACK_RECOMMENDED`) and zero-cost results (free tier, use
  `FALLBACK_NONE`), including guidance for plugins wrapping SaaS vendors

### Key Entities

- **FallbackHint**: Enumeration representing the plugin's recommendation for core's
  fallback behavior. Contains four values with distinct semantics for orchestration.
- **GetActualCostResponse**: The gRPC response message returned by plugins for actual
  cost queries. Now includes an optional fallback_hint field in addition to the
  existing results field.
- **CostCalculator (SDK)**: Helper type in the Go SDK that provides convenience
  methods for constructing protocol responses. Will include new method for responses
  with explicit hints.

## Clarifications

### Session 2025-12-06

- Q: What naming convention should be used for the FallbackHint enum values to avoid
  scope collisions?
  - A: Use `FALLBACK_HINT_` prefix (e.g., `FALLBACK_HINT_UNSPECIFIED`) as standard
    Protobuf practice.

- Q: How should plugins distinguish between "resource not found" (no billing data)
  vs "resource found with $0.00 cost" (free tier)?
  - A: Empty array `[]` means "no records found" → use `FALLBACK_RECOMMENDED`
  - Array with zero-cost entry `[{cost: 0.00}]` means "genuinely zero cost" →
    use `FALLBACK_NONE`
  - Plugins wrapping SaaS vendors that conflate these must translate to proper
    semantics

- Q: How should Core treat `FALLBACK_HINT_UNSPECIFIED` vs `FALLBACK_HINT_NONE` logic?
  - A: Identical logic. Both result in "no fallback". The distinction is only
    useful for observability/debugging.
  - Note: "We don't have legacy yet" context provided, but treating them identically
    ensures future safety if default behavior changes.

- Q: What should happen if a plugin returns results (data) AND sets
  `FALLBACK_HINT_REQUIRED`?
  - A: Data wins. Core should use the returned results, ignore the hint, and log a
    warning about the inconsistency.

- Q: What is the preferred Go SDK implementation pattern for helper functions
  creating responses with fallback hints?
  - A: Functional Options Pattern (e.g., `NewActualCostResponse(..., WithFallbackHint(hint))`).

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All existing plugin implementations compile and pass tests without
  modification after the change
- **SC-002**: New SDK helper enables plugins to set fallback hints with a single
  method call
- **SC-003**: Protocol documentation clearly describes when to use each hint value
  with concrete examples
- **SC-004**: Wire format overhead is zero when fallback hint is not set (proto3
  default behavior)
- **SC-005**: Plugin developers can determine correct hint value using decision
  matrix in documentation
