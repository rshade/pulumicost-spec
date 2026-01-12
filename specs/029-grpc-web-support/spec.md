# Feature Specification: Multi-Protocol Plugin Access (gRPC-Web and HTTP)

**Feature Branch**: `029-grpc-web-support`
**Created**: 2025-12-29
**Status**: Draft
**Input**: User description: "gRPC-Web support for browser-based plugin access enabling
full lifecycle management, plus Go-based client support for batch resource queries"
**Related Issue**: [#189](https://github.com/rshade/finfocus-spec/issues/189)

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Web Dashboard Queries Plugin for Cost Estimate (Priority: P1)

A developer using the Pulumi Insights web dashboard wants to get a cost estimate for a
resource they're about to deploy. The browser-based application needs to communicate
directly with a running FinFocus plugin to retrieve pricing information without
requiring a separate backend proxy.

**Why this priority**: This is the core value proposition - enabling browser-based cost
queries is the primary use case that unlocks all web-based tooling. Without this, no web
integration is possible.

**Independent Test**: Can be fully tested by loading a web page that calls the
EstimateCost RPC via gRPC-Web protocol and displays the returned cost estimate.

**Acceptance Scenarios**:

1. **Given** a FinFocus plugin is running and accessible,
   **When** a browser client sends an EstimateCost request via gRPC-Web protocol,
   **Then** the plugin returns cost estimate data in a format the browser can parse.
2. **Given** a browser client on a different origin than the plugin,
   **When** the client attempts to call plugin RPCs,
   **Then** CORS headers allow the cross-origin request to succeed.
3. **Given** a browser client sends a malformed request,
   **When** the plugin receives it,
   **Then** appropriate error information is returned that the browser can display.

---

### User Story 2 - Backend Service Queries Costs for Multiple Resources (Priority: P1)

The Pulumi Insights backend service needs to retrieve cost data for a batch of resources
in a single operation. When a user views their infrastructure costs, the backend should
efficiently query one or more plugins to get cost estimates for all resources in their
stack without making individual sequential requests.

**Why this priority**: Batch operations are critical for production use cases where
querying resources one-at-a-time would be too slow. This enables real-world integration
with platforms like Pulumi Insights.

**Independent Test**: Can be fully tested by writing a Go program that sends a list of
10+ resource descriptors and receives cost data for all of them efficiently.

**Acceptance Scenarios**:

1. **Given** a Go-based backend service with a list of 50 resources,
   **When** the service queries a plugin for cost estimates,
   **Then** the service receives cost data for all resources without timeout.
2. **Given** a batch request containing resources from different providers,
   **When** the backend queries appropriate plugins,
   **Then** costs are aggregated from multiple plugins correctly.
3. **Given** a batch request where some resources are unsupported,
   **When** the plugin processes the request,
   **Then** supported resources return cost data and unsupported ones return clear
   error indicators without failing the entire batch.
4. **Given** a large batch of 500+ resources,
   **When** the backend queries the plugin,
   **Then** the response is returned within acceptable time limits (based on plugin
   capacity).

---

### User Story 3 - Multi-Tenant Platform Manages Per-Org Plugin Instances (Priority: P1)

Pulumi Insights serves multiple organizations, each with their own cloud credentials
(AWS accounts, Azure subscriptions, GCP projects). The platform needs to launch and
manage separate plugin instances for each organization to ensure credential isolation
and prevent data leakage between tenants.

**Why this priority**: Multi-tenant credential isolation is a security requirement.
Without proper isolation, one organization's credentials could be exposed to another,
creating compliance and security risks.

**Independent Test**: Can be fully tested by launching plugins for two different
organizations and verifying each plugin uses only its designated credentials.

**Acceptance Scenarios**:

1. **Given** a platform serving Organization A and Organization B,
   **When** the orchestrator launches AWS plugins for each organization,
   **Then** each plugin instance receives only its organization's AWS credentials.
2. **Given** Organization A's plugin is running,
   **When** Organization B requests cost data,
   **Then** the request is routed to Organization B's plugin instance (not A's).
3. **Given** an organization has not made any requests recently,
   **When** the idle timeout period expires,
   **Then** the orchestrator terminates the idle plugin to conserve resources.
4. **Given** a plugin instance crashes or becomes unhealthy,
   **When** the orchestrator detects the failure,
   **Then** a new plugin instance is launched with the same configuration.
5. **Given** Organization A's plugin is compromised,
   **When** an attacker attempts to access other organizations' data,
   **Then** process isolation prevents access to other organizations' credentials.

---

### User Story 4 - Web Client Discovers Available Plugin Services (Priority: P2)

A web-based development tool needs to discover what services and methods a running
plugin supports. Using service reflection, the browser can dynamically understand
plugin capabilities without prior knowledge of the plugin's implementation.

**Why this priority**: Service discovery enables dynamic UI generation and plugin
introspection, which is essential for building flexible web tools that work with
any plugin.

**Independent Test**: Can be fully tested by having a browser client call the reflection
service and displaying the list of available RPC methods.

**Acceptance Scenarios**:

1. **Given** a plugin is running with reflection enabled,
   **When** a browser client queries the reflection endpoint via gRPC-Web,
   **Then** the client receives a list of available services and methods.
2. **Given** multiple plugins are running,
   **When** a browser client queries each plugin's reflection,
   **Then** each plugin returns its own service definitions independently.

---

### User Story 5 - Web Client Monitors Plugin Health (Priority: P2)

An operations dashboard needs to monitor the health status of running plugins. The web
interface should be able to check if plugins are healthy and responding without
requiring a backend intermediary.

**Why this priority**: Health monitoring is critical for operational visibility but
depends on basic connectivity (P1) being established first.

**Independent Test**: Can be fully tested by having a browser periodically call a health
check endpoint and displaying status indicators.

**Acceptance Scenarios**:

1. **Given** a healthy running plugin,
   **When** a browser client calls the health check endpoint,
   **Then** the client receives a positive health status.
2. **Given** a plugin that is unhealthy or overloaded,
   **When** a browser client calls the health check endpoint,
   **Then** the client receives appropriate status information indicating the issue.
3. **Given** a plugin that is not responding,
   **When** a browser client attempts to connect,
   **Then** the client receives a timeout or connection error within a reasonable time.

---

### User Story 6 - Go Client Library for Plugin Integration (Priority: P2)

A developer building a Go-based service that needs to integrate with FinFocus plugins
wants a simple client library that handles connection management, protocol negotiation,
and error handling. The SDK should provide a typed Go client that makes plugin
integration straightforward.

**Why this priority**: Go clients enable server-to-server integration which is essential
for platforms like Pulumi Insights. A well-designed client library reduces integration
friction.

**Independent Test**: Can be fully tested by using the client library to call each RPC
method and verifying type-safe responses.

**Acceptance Scenarios**:

1. **Given** a Go service importing the SDK client library,
   **When** the developer creates a new plugin client,
   **Then** the client connects to the plugin with minimal configuration.
2. **Given** a configured plugin client,
   **When** the developer calls EstimateCost with typed request parameters,
   **Then** a typed response is returned without manual serialization.
3. **Given** a plugin that returns an error,
   **When** the Go client receives the error,
   **Then** the error is properly typed and includes actionable information.

---

### User Story 7 - Web Application Executes All Cost RPCs (Priority: P3)

A full-featured cost management web application needs to execute all available
cost-related operations: getting actual costs, projected costs, pricing specs,
recommendations, and budget information. All 8 CostSourceService RPCs should be
accessible from the browser.

**Why this priority**: Full RPC access builds on basic connectivity and enables complete
feature parity with native gRPC clients, but is not required for initial value delivery.

**Independent Test**: Can be fully tested by creating a web page with buttons for each
RPC type and verifying each returns expected data.

**Acceptance Scenarios**:

1. **Given** a running plugin,
   **When** a browser client calls Name RPC via gRPC-Web,
   **Then** the plugin name is returned.
2. **Given** a running plugin,
   **When** a browser client calls GetActualCost with valid parameters,
   **Then** historical cost data is returned.
3. **Given** a running plugin,
   **When** a browser client calls GetRecommendations,
   **Then** optimization recommendations are returned (or empty list if not supported).
4. **Given** a running plugin,
   **When** a browser client calls GetBudgets,
   **Then** budget information is returned (or appropriate error if not supported).

---

### User Story 8 - Plugin Lifecycle Management from Web UI (Priority: P3)

A platform administrator using a web-based management console wants to discover running
plugins, view their status, and manage their lifecycle. The web UI should provide
visibility into which plugins are available and their operational state.

**Why this priority**: Full lifecycle management is an advanced feature that builds on
all previous capabilities and provides administrative control.

**Independent Test**: Can be fully tested by loading a web page that enumerates running
plugins and displays their status information.

**Acceptance Scenarios**:

1. **Given** multiple plugins are running,
   **When** a browser client queries the plugin registry,
   **Then** a list of available plugins with their endpoints is returned.
2. **Given** a running plugin,
   **When** a browser client queries its status,
   **Then** operational metrics and health information are displayed.

---

### Edge Cases

- What happens when a browser client attempts to connect to a plugin that requires
  authentication but provides no credentials?
- How does the system handle browser clients with slow or unstable network connections?
- What happens when a browser client sends a request larger than the maximum allowed
  message size?
- How does the system behave when the browser has disabled certain security features
  (mixed content, CORS preflight)?
- What happens when a plugin is restarted while a browser client has an active
  connection?
- What happens when a batch request includes duplicate resource identifiers?
- How does the system handle partial failures in batch operations (some resources
  succeed, others fail)?
- What is the maximum batch size supported, and how does the system respond when
  exceeded?
- What happens when two requests for the same organization arrive simultaneously and
  no plugin instance exists yet?
- How does the system handle credential rotation for a running plugin instance?
- What happens when the maximum number of plugin instances per node is reached?

## Requirements _(mandatory)_

### Functional Requirements

#### Protocol Support

- **FR-001**: Plugins MUST serve all 8 CostSourceService RPCs via HTTP-based protocols
  accessible from both browsers and Go clients.
- **FR-002**: Plugins MUST handle Cross-Origin Resource Sharing (CORS) to allow requests
  from browser applications on different origins.
- **FR-003**: Plugins MUST support service reflection accessible from all client types.
- **FR-004**: Plugins MUST announce their HTTP endpoint on startup (in addition to the
  existing PORT announcement).
- **FR-005**: Plugins MUST continue to support existing native gRPC clients (backward
  compatibility).

#### Batch Operations

- **FR-006**: The SDK MUST support efficient batch queries for cost data across multiple
  resources in a single operation.
- **FR-007**: Batch operations MUST return partial results when some resources succeed
  and others fail, rather than failing the entire batch.
- **FR-008**: Batch operations MUST support configurable concurrency limits to prevent
  overwhelming plugins.
- **FR-009**: The SDK MUST provide clear error reporting for each resource in a batch
  that indicates success, failure, or unsupported status.

#### Client Libraries

- **FR-010**: The SDK MUST provide a Go client library with typed methods for all
  CostSourceService RPCs.
- **FR-011**: The Go client library MUST handle connection pooling and reconnection
  automatically.
- **FR-012**: The Go client library MUST support configurable timeouts and retry
  policies.

#### Health and Monitoring

- **FR-013**: Plugins MUST provide a health check endpoint accessible via HTTP for
  browser-based and programmatic monitoring.
- **FR-014**: Plugins MUST handle graceful shutdown for HTTP connections the same as for
  native gRPC connections.

#### Configuration and Security

- **FR-015**: The SDK MUST provide configuration options for CORS origins via YAML file
  and environment variables, allowing operators to restrict which web origins can access
  plugins. Default: deny all cross-origin requests (secure default).
- **FR-016**: Plugins MUST use ambient credentials from environment variables as the
  primary authentication model for cloud API access.
- **FR-017**: Plugins MUST return appropriate error responses that all client types can
  parse and handle.
- **FR-018**: The SDK MUST support configurable network binding (loopback-only vs.
  network-accessible) with secure defaults.
- **FR-019**: Plugins MUST handle concurrent requests from browser, Go, and native gRPC
  clients without interference.
- **FR-019a**: The SDK MUST support plugin configuration via YAML file and environment
  variables for all configurable options (timeouts, CORS, network binding, etc.).
  Environment variables take precedence over YAML values.

#### Multi-Tenant Plugin Lifecycle (Orchestrator)

- **FR-020**: The orchestrator MUST launch separate plugin instances for each tenant
  (organization) to ensure credential isolation.
- **FR-021**: The orchestrator MUST pass tenant credentials to plugins via environment
  variables at launch time (not via request headers).
- **FR-022**: The orchestrator MUST route requests to the correct plugin instance based
  on tenant identifier.
- **FR-023**: The orchestrator MUST support lazy instantiation of plugins (launch on
  first request for a tenant).
- **FR-024**: The orchestrator MUST support configurable idle timeout to terminate
  inactive plugin instances and conserve resources. Default: 5 minutes.
- **FR-025**: The orchestrator MUST automatically restart failed plugin instances with
  the same configuration.
- **FR-026**: The orchestrator MUST enforce maximum plugin instance limits per node to
  prevent resource exhaustion.
- **FR-027**: The orchestrator MUST provide an API to query running plugin instances
  and their tenant assignments.
- **FR-028**: The SDK MUST include tenant context in request metadata for audit logging
  purposes (but NOT credentials).

### Key Entities

- **Browser Client**: A web application running in a user's browser that needs to
  communicate with FinFocus plugins. Communicates via HTTP-based protocols.
- **Go Client**: A Go-based service or application (e.g., Pulumi Insights backend) that
  queries plugins programmatically. Uses the SDK-provided client library.
- **Plugin Endpoint**: The network address and port where a plugin serves requests. May
  support multiple protocols simultaneously.
- **Batch Request**: A single operation that queries cost data for multiple resources,
  returning aggregated results with per-resource status.
- **CORS Configuration**: Settings that define which web origins are permitted to access
  plugin services.
- **Service Reflection**: A mechanism for clients to discover available services and
  methods at runtime.
- **Orchestrator**: The component responsible for launching, routing, and managing
  plugin instances. Handles multi-tenant isolation and lifecycle management.
- **Tenant**: An organization or account that has its own cloud credentials and isolated
  plugin instances. Identified by a unique tenant ID.
- **Plugin Instance**: A running plugin process serving a specific tenant, launched with
  that tenant's credentials in environment variables.
- **Plugin Pool**: The set of all running plugin instances managed by an orchestrator,
  organized by tenant and plugin type.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Browser clients can successfully call all 8 CostSourceService RPCs and
  receive valid responses.
- **SC-002**: Go clients can successfully call all 8 CostSourceService RPCs using the
  SDK client library.
- **SC-003**: Batch queries for 100 resources complete within 5 seconds under normal
  plugin load.
- **SC-004**: Existing native gRPC clients continue to work without modification after
  multi-protocol support is added.
- **SC-005**: Browser clients on different origins can access plugins when CORS is
  properly configured.
- **SC-006**: Service discovery via reflection works from all client types, returning
  accurate service definitions.
- **SC-007**: Plugin health checks are accessible from all client types within 1 second
  response time.
- **SC-008**: Plugins handle at least 100 concurrent connections (mixed client types)
  without degradation.
- **SC-009**: Plugin documentation and examples enable developers to set up Go client
  integration within 15 minutes.
- **SC-010**: Plugin documentation and examples enable developers to set up browser
  access within 30 minutes.
- **SC-011**: Error messages returned to clients are clear enough for operators to
  diagnose issues.
- **SC-012**: Batch operations with partial failures return results for successful
  resources and clear errors for failed ones.
- **SC-013**: Orchestrator can manage at least 100 concurrent plugin instances without
  degradation.
- **SC-014**: Plugin launch latency (cold start) is under 5 seconds for first request
  to a new tenant.
- **SC-015**: Idle plugin cleanup correctly terminates instances after configured
  timeout without affecting other tenants.
- **SC-016**: No credentials are logged or exposed in error messages, traces, or
  metrics.

## Assumptions

- Plugins will be accessed from modern browsers (Chrome, Firefox, Safari, Edge) that
  support standard web APIs.
- Go clients will use Go 1.21+ with standard library HTTP support.
- Network infrastructure between clients and plugins allows HTTP/1.1 and/or HTTP/2
  traffic.
- Operators deploying plugins in production will configure appropriate network security
  (firewalls, TLS termination) at the infrastructure level.
- Browser clients will use standard JavaScript/TypeScript libraries for protocol
  communication.
- Initial implementation focuses on unary (request-response) RPCs; streaming support is
  out of scope for this specification.
- Batch operations will be implemented client-side in the SDK using parallel requests
  to existing RPCs. A dedicated batch RPC is out of scope (see future enhancement).
- **Multi-tenant platforms will run separate plugin instances per tenant (organization)
  for credential isolation. This is the recommended and default model.**
- Plugin credentials are passed via environment variables at launch time, not via
  request headers. This ensures compatibility with standard cloud SDK credential chains.
- The orchestrator is a web platform concern (e.g., Pulumi Insights), not part of
  finfocus-core. finfocus-core is a CLI that launches plugins locally via command
  line; multi-tenant orchestration is implemented by the consuming web service.

## Scope Boundaries

### In Scope

- HTTP-based protocol support for all existing CostSourceService RPCs
- CORS configuration for cross-origin browser access
- Go client library with typed methods and connection management
- Batch query support for multiple resources
- Service reflection accessible from all client types
- HTTP health check endpoint
- SDK configuration options for HTTP/browser/Go client access
- Multi-tenant plugin lifecycle management (orchestrator requirements)
- Per-tenant plugin instance isolation via environment variable credentials
- Idle plugin cleanup and automatic restart on failure
- Documentation and examples for browser and Go client setup

### Out of Scope

- Specific authentication implementations (JWT, OAuth, etc.) - left to operators
- Per-request credential passing (see future enhancement issue #220)
- Streaming RPC support (all current RPCs are unary)
- TLS certificate management (handled at infrastructure level)
- Specific browser client library implementations (TypeScript types may be provided)
- WebSocket transport (may be considered for future streaming support)
- Cross-plugin query aggregation (batch queries target single plugin)
- Credential rotation for running plugins (requires restart)
- Dedicated batch RPC (see future enhancement issue #221 - use client-side parallelism)

## Clarifications

### Session 2025-12-29

- Q: What should be the default idle timeout for terminating inactive plugin instances?
  A: 5 minutes (balanced approach for intermittent queries)
- Q: Should browsers call plugins directly or through backend proxy? → A: Direct browser access enabled (AJAX-like calls)
- Q: How should plugin options (timeout, lifetime, CORS) be configured? → A: Via YAML file and environment variables

## Dependencies

- Existing CostSourceService proto definitions (no changes required for this spec)
- Current pluginsdk.Serve() implementation (will be extended)
- gRPC reflection support (already implemented via #181)
- Orchestrator implementation (web platform responsibility, not finfocus-core)
