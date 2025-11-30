# Feature Specification: Structured Logging Example for EstimateCost

**Feature Branch**: `007-zerolog-logging-example`
**Created**: 2025-11-26
**Status**: Draft
**Input**: GitHub Issue #83 - Add structured logging example for EstimateCost (T040)

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Developer Learns Logging Patterns (Priority: P1)

A plugin developer wants to understand how to properly integrate zerolog structured
logging with PulumiCost plugin operations, specifically for the EstimateCost RPC.

**Why this priority**: Documentation through working examples is the most effective way
to ensure consistent logging practices across the plugin ecosystem. This example serves
as the canonical reference for NFR-001 compliance.

**Independent Test**: Can be tested by running the example code and verifying the
structured JSON output contains all expected fields (trace_id, operation, resource_type,
duration_ms, cost details).

**Acceptance Scenarios**:

1. **Given** a developer reads the integration test example, **When** they examine the
   logging code, **Then** they can understand how to create a configured logger with
   plugin name and version
2. **Given** a developer reads the example, **When** they examine the request logging
   pattern, **Then** they can see how to log incoming requests with relevant context
   fields (resource_type, attributes count)
3. **Given** a developer reads the example, **When** they examine the response logging
   pattern, **Then** they can see how to log successful responses with cost details
   (estimated_cost, currency)

---

### User Story 2 - Plugin Developer Implements Error Logging (Priority: P1)

A plugin developer needs to understand the proper pattern for logging errors in a way
that supports debugging and monitoring, including correlation IDs for distributed tracing.

**Why this priority**: Error logging is equally critical to success logging - operators
need clear error context to diagnose production issues. Correlation ID handling enables
end-to-end tracing.

**Independent Test**: Can be tested by running the error scenario example and verifying
the error log contains trace_id, error code, and descriptive message.

**Acceptance Scenarios**:

1. **Given** an EstimateCost call fails with an unsupported resource, **When** the
   developer logs the error, **Then** the log includes error code, error message,
   resource_type, and trace_id
2. **Given** a request arrives with a correlation ID header, **When** the developer
   logs any operation, **Then** the trace_id appears in all related log entries
3. **Given** an operation times out or encounters an internal error, **When** the
   developer logs the failure, **Then** the log includes duration_ms to help identify
   performance issues

---

### User Story 3 - Operator Monitors EstimateCost Health (Priority: P2)

An operator wants to be able to query logs to understand EstimateCost operation health,
including latency percentiles and error patterns.

**Why this priority**: Operational monitoring requires consistent log structure. While
not required for plugin development, this enables ecosystem-wide observability.

**Independent Test**: Can be tested by verifying log output follows the documented field
naming conventions and can be parsed by standard JSON tools.

**Acceptance Scenarios**:

1. **Given** logs from multiple EstimateCost calls, **When** an operator queries by
   operation field, **Then** they can filter to only EstimateCost operations
2. **Given** logs with duration_ms fields, **When** an operator aggregates by
   resource_type, **Then** they can calculate average latency per resource type
3. **Given** logs with error_code fields, **When** an operator groups errors, **Then**
   they can identify the most common failure modes

---

### Edge Cases

- What happens when trace_id is not present in the request context?
  (Log without trace_id field rather than failing; allows graceful degradation)
- How should very large attribute maps be logged?
  (Log attribute count rather than full attributes to avoid log bloat)
- What happens when the estimated cost is zero?
  (Log normally with cost=0; this is a valid business case)
- How should sensitive attribute values be handled?
  (Do not log attribute values; only log count and keys if needed)

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: Example MUST demonstrate creating a zerolog.Logger with plugin name and
  version using the SDK's NewPluginLogger function
- **FR-002**: Example MUST demonstrate logging a request with resource_type and
  operation context fields
- **FR-003**: Example MUST demonstrate logging a successful response with estimated
  cost, currency, and duration_ms fields
- **FR-004**: Example MUST demonstrate logging an error with error code, error message,
  and original request context
- **FR-005**: Example MUST demonstrate correlation ID (trace_id) propagation through
  request logging
- **FR-006**: Example MUST use the standard field name constants (FieldTraceID,
  FieldOperation, FieldDurationMs, FieldResourceType, etc.)
- **FR-007**: Example MUST include code comments explaining logging best practices
- **FR-008**: Example MUST be placed in sdk/go/testing/integration_test.go as a test
  function that can be run and verified
- **FR-009**: Example MUST demonstrate the LogOperation timing helper for measuring
  operation duration
- **FR-010**: Example MUST follow existing integration test patterns in the file

### Key Entities

- **Logger**: Configured zerolog instance with plugin metadata
- **TraceID**: Correlation identifier propagated via context for distributed tracing
- **LogFields**: Standard field constants for consistent naming (from 005-zerolog spec)
- **Operation**: EstimateCost RPC being logged
- **CostResult**: The estimated cost value and currency to be logged on success

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Developers can run the example test and see structured JSON log output
- **SC-002**: Example demonstrates all three logging scenarios (request, success response,
  error response)
- **SC-003**: All log output uses standard field names from the zerolog SDK utilities
- **SC-004**: Example code compiles and passes go test without modification
- **SC-005**: Code comments provide clear guidance that developers can apply to their
  own plugins

## Assumptions

- The zerolog SDK logging utilities from spec 005-zerolog are implemented and available
- EstimateCost RPC from spec 006-estimate-cost is implemented and available
- Plugin developers are familiar with zerolog's builder pattern API
- The integration test file already exists and follows established patterns
- Structured JSON output is the standard format (not console/pretty printing)

## Dependencies

- 005-zerolog: Zerolog SDK Logging Utilities (provides NewPluginLogger, field constants,
  LogOperation)
- 006-estimate-cost: EstimateCost RPC (provides the operation being logged)
- sdk/go/testing package: Existing test infrastructure

## Out of Scope

- Creating new zerolog utilities (use existing from 005-zerolog)
- Implementing actual EstimateCost functionality (use mock or existing implementation)
- Log aggregation or centralized logging infrastructure
- Console/pretty-print formatters
- Metric emission (covered separately in NFR-002 of 006-estimate-cost)
