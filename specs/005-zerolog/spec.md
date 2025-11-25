# Feature Specification: Zerolog SDK Logging Utilities

**Feature Branch**: `005-zerolog`
**Created**: 2025-11-24
**Status**: Draft
**Input**: GitHub Issue #75 - Add zerolog v1.34.0+ SDK logging utilities

## Clarifications

### Session 2025-11-24

- Q: Where does the logger write output? → A: Default to os.Stderr; accept
  io.Writer parameter for flexibility; support --logfile/--logdir flags with
  default filenames (core.log, plugin.log, etc.)
- Q: Should SDK handle log rotation for file output? → A: Delegate to external
  tools (logrotate, systemd, etc.)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Plugin Developer Creates Standardized Logger (Priority: P1)

A plugin developer wants to create a consistent, well-configured logger for their
PulumiCost plugin that includes standard metadata fields and follows ecosystem
conventions.

**Why this priority**: This is the foundation - without a standard logger
constructor, plugins cannot produce consistent logs. Every other feature depends
on this.

**Independent Test**: Can be tested by creating a logger instance with plugin
name/version and verifying it outputs structured JSON with required base fields.

**Acceptance Scenarios**:

1. **Given** a plugin developer calls NewPluginLogger with name "aws-public" and
   version "v1.0.0", **When** logging at Info level, **Then** the output includes
   structured JSON with plugin_name, plugin_version, and timestamp fields
2. **Given** a plugin developer specifies Debug log level, **When** Debug messages
   are logged, **Then** they appear in output; when Info level is specified, Debug
   messages are suppressed
3. **Given** a plugin developer creates multiple logger instances, **When** each
   uses different names/versions, **Then** each logger correctly identifies its
   source plugin

---

### User Story 2 - Core System Traces Requests Through Plugin Calls (Priority: P1)

A system operator needs to trace a single request through the entire PulumiCost
system, including into plugin RPC calls, to debug issues or analyze performance.

**Why this priority**: Distributed tracing is essential for debugging production
issues in a plugin-based architecture. Equal priority with logging since both are
required for operational visibility.

**Independent Test**: Can be tested by configuring a gRPC server with the
interceptor, making a call with trace_id metadata, and verifying the trace_id is
available in the handler context.

**Acceptance Scenarios**:

1. **Given** a gRPC server configured with TracingUnaryServerInterceptor, **When**
   a client sends a request with `x-pulumicost-trace-id` metadata header, **Then**
   the handler can retrieve the trace_id from context using TraceIDFromContext
2. **Given** a gRPC request arrives without trace_id metadata, **When** the
   interceptor processes it, **Then** TraceIDFromContext returns an empty string
   (no error thrown)
3. **Given** multiple concurrent requests with different trace_ids, **When** each
   is processed, **Then** each handler sees its own correct trace_id without
   cross-contamination

---

### User Story 3 - Plugin Developer Logs Operation Timing (Priority: P2)

A plugin developer wants to easily measure and log how long specific operations
take without writing boilerplate timing code.

**Why this priority**: Operation timing is important for performance monitoring
but not strictly required for basic logging functionality.

**Independent Test**: Can be tested by calling LogOperation, performing a task,
calling the returned function, and verifying duration_ms appears in log output.

**Acceptance Scenarios**:

1. **Given** a developer calls LogOperation with operation name "GetProjectedCost",
   **When** they call the returned function after completing the operation,
   **Then** the log includes operation name and duration_ms field
2. **Given** an operation takes 250ms, **When** the timing function is called,
   **Then** duration_ms is recorded as approximately 250 (within reasonable
   precision)

---

### User Story 4 - Plugin Developer Uses Standard Field Names (Priority: P2)

A plugin developer needs to use consistent field names when logging so that log
aggregation and analysis tools can easily parse logs from any PulumiCost plugin.

**Why this priority**: Standardization enables ecosystem-wide log analysis but
individual plugins can function without it.

**Independent Test**: Can be tested by using the exported field constants in log
statements and verifying the resulting JSON uses those exact field names.

**Acceptance Scenarios**:

1. **Given** a developer imports the logging package, **When** they use
   FieldTraceID, FieldOperation, FieldResourceType constants, **Then** logs
   contain fields named exactly "trace_id", "operation", "resource_type"
2. **Given** all PulumiCost plugins use the standard field constants, **When**
   logs are aggregated, **Then** they can be queried consistently across all
   plugins

---

### Edge Cases

- What happens when NewPluginLogger is called with empty plugin name or version?
  (Logs will include empty strings for those fields - not an error)
- What happens when multiple trace_ids are passed in metadata?
  (Use first value only)
- How does the system handle invalid log level values?
  (Use zerolog's behavior - invalid becomes Info)
- What happens if LogOperation's returned function is never called?
  (No duration logged - developer responsibility)
- How are malformed metadata keys handled?
  (gRPC metadata is case-insensitive, interceptor handles normalization)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide NewPluginLogger function that creates a
  zerolog.Logger configured with plugin name, version, log level, and output
  destination (default: os.Stderr; accepts io.Writer for flexibility)
- **FR-011**: System MUST support --logfile and --logdir flags for file-based
  logging with default filenames (e.g., core.log, plugin.log based on component)
- **FR-002**: System MUST export TracingUnaryServerInterceptor function that
  returns a gRPC UnaryServerInterceptor
- **FR-003**: Interceptor MUST extract trace_id from gRPC metadata key
  "x-pulumicost-trace-id" and add it to request context
- **FR-004**: System MUST provide TraceIDFromContext function that extracts
  trace_id string from context
- **FR-005**: System MUST provide ContextWithTraceID function that adds trace_id
  to context
- **FR-006**: System MUST export standard field name constants for consistent
  logging (FieldTraceID, FieldComponent, FieldOperation, FieldDurationMs,
  FieldResourceURN, FieldResourceType, FieldPluginName, FieldPluginVersion,
  FieldCostMonthly, FieldAdapter, FieldErrorCode)
- **FR-007**: System MUST provide LogOperation function that returns a function
  to log operation duration
- **FR-008**: System MUST use zerolog v1.34.0 or higher as the logging library
- **FR-009**: Logger output MUST be structured JSON for machine parsing
- **FR-010**: All utilities MUST have unit tests with 90%+ code coverage

### Key Entities

- **Logger**: Configured zerolog instance with plugin metadata (name, version,
  level)
- **TraceID**: String identifier for distributed request tracing, propagated via
  gRPC metadata
- **Field Constants**: Standardized string constants for consistent log field
  naming across ecosystem
- **Context Keys**: Context key types for storing/retrieving trace_id in request
  context

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Plugin developers can create a configured logger in a single
  function call
- **SC-002**: Trace IDs flow correctly from gRPC client metadata through to
  handler context in 100% of requests
- **SC-003**: All 11 standard field name constants are exported and documented
- **SC-004**: Unit test coverage reaches 90% or higher for all logging utilities
- **SC-005**: Example usage demonstrates complete pattern for plugin
  implementation
- **SC-006**: Documentation includes logging best practices for plugin developers

## Assumptions

- zerolog v1.34.0+ is compatible with existing Go SDK dependencies
- Plugin developers are familiar with zerolog's builder pattern API
- gRPC metadata propagation is already established in core system
- Plugins use standard Go context.Context for request handling
- JSON log output is the required format (not console/pretty printing in
  production)
- Performance overhead of interceptor and logging is acceptable (zerolog is
  already highly optimized)

## Dependencies

- zerolog v1.34.0+ library
- google.golang.org/grpc (already a project dependency)
- google.golang.org/grpc/metadata (already a project dependency)

## Out of Scope

- Log aggregation or centralized logging infrastructure
- Log rotation (delegate to external tools like logrotate/systemd)
- Console/pretty-print formatters for development (can be added later)
- Streaming RPC interceptors (only unary for now)
- Automatic trace_id generation in SDK (responsibility of core system)
