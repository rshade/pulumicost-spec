# Feature Specification: Document pluginsdk.Serve() Behavior

**Feature Branch**: `001-pluginsdk-serve-docs`
**Created**: 2025-12-08
**Status**: Draft
**Input**: User description: "Provide comprehensive documentation for the pluginsdk.Serve() function,
detailing its startup behavior, environment variable usage, and expected flags."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Developer Learns Serve() Usage (Priority: P1)

A plugin developer building a new PulumiCost plugin needs to understand how to configure and start
the gRPC server using `pluginsdk.Serve()`. They need clear documentation explaining the function
signature, configuration options, and expected behavior.

**Why this priority**: Without clear documentation, developers cannot successfully implement plugins.
This is the core value proposition of the documentation feature.

**Independent Test**: Can be fully tested by a developer reading the documentation and successfully
implementing a minimal plugin that starts and accepts gRPC connections.

**Acceptance Scenarios**:

1. **Given** a developer new to PulumiCost, **When** they read the Serve() documentation,
   **Then** they can write a minimal main() function that starts a plugin server without errors.
2. **Given** a developer reading the documentation, **When** they look for configuration options,
   **Then** they find a complete list of ServeConfig fields with descriptions.

---

### User Story 2 - Plugin Developer Understands Port Resolution (Priority: P1)

A plugin developer needs to understand how the server determines which port to listen on.
They need documentation that clearly explains the priority order:
command-line flag > environment variable > ephemeral port.

**Why this priority**: Port configuration is critical for multi-plugin orchestration and debugging.
Incorrect understanding leads to port conflicts and failed deployments.

**Independent Test**: Can be fully tested by a developer configuring ports using each method
(flag, env var, default) and observing the documented behavior.

**Acceptance Scenarios**:

1. **Given** a developer passing `--port 50051`, **When** they start the plugin,
   **Then** the server listens on port 50051 as documented.
2. **Given** a developer with `PULUMICOST_PLUGIN_PORT=50052` set and no flag,
   **When** they start the plugin, **Then** the server listens on port 50052 as documented.
3. **Given** a developer with no port configuration, **When** they start the plugin,
   **Then** the server uses an ephemeral port and outputs `PORT=<assigned>` to stdout as documented.

---

### User Story 3 - Plugin Developer Configures Logging and Tracing (Priority: P2)

A plugin developer needs to configure logging levels, formats, and distributed tracing for
debugging and production monitoring. They need documentation covering all relevant environment
variables.

**Why this priority**: Proper logging and tracing are essential for production deployments
but not required for initial development.

**Independent Test**: Can be fully tested by setting environment variables and observing log output
matches documented behavior.

**Acceptance Scenarios**:

1. **Given** a developer sets `PULUMICOST_LOG_LEVEL=debug`, **When** they start the plugin,
   **Then** debug logs are emitted as documented.
2. **Given** a developer sets `PULUMICOST_TRACE_ID=abc123`, **When** requests are processed,
   **Then** the trace ID appears in logs as documented.

---

### User Story 4 - DevOps Engineer Deploys Multiple Plugins (Priority: P2)

A DevOps engineer deploying pulumicost-core needs to understand how to orchestrate multiple plugins
(e.g., aws-public + aws-ce) with unique ports. They need documentation explaining why
PORT fallback is not supported.

**Why this priority**: Multi-plugin orchestration is a production deployment concern.
The documentation prevents common deployment errors.

**Independent Test**: Can be fully tested by deploying two plugins with different --port flags
and verifying they start on distinct ports.

**Acceptance Scenarios**:

1. **Given** a DevOps engineer reading the documentation, **When** they look for why PORT is
   not supported, **Then** they find an explanation about multi-plugin conflict prevention.
2. **Given** a DevOps engineer deploying two plugins, **When** they use `--port` flags with
   distinct values, **Then** both plugins start successfully on their assigned ports.

---

### User Story 5 - Plugin Developer Implements Graceful Shutdown (Priority: P3)

A plugin developer needs to understand how context cancellation triggers graceful shutdown.
They need documentation explaining the shutdown sequence.

**Why this priority**: Graceful shutdown is important for production reliability
but not critical for initial development.

**Independent Test**: Can be fully tested by canceling the context and observing the server
stops gracefully without dropping active requests.

**Acceptance Scenarios**:

1. **Given** a developer reading the documentation, **When** they look for shutdown behavior,
   **Then** they find that context cancellation triggers GracefulStop().
2. **Given** a running plugin with active requests, **When** the context is canceled,
   **Then** existing requests complete before the server stops.

---

### Edge Cases

- What happens when an invalid port number is configured (e.g., negative, > 65535)?
- How does the system handle port conflicts (port already in use)?
- What happens if ParsePortFlag() is called before flag.Parse()?
- How does the server behave when the listener cannot be created (permission denied)?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: Documentation MUST explain the `Serve()` function signature and return behavior
- **FR-002**: Documentation MUST describe the `ServeConfig` struct with all fields and their purposes
- **FR-003**: Documentation MUST explain port resolution priority:
  `--port` flag > `PULUMICOST_PLUGIN_PORT` > ephemeral (0)
- **FR-004**: Documentation MUST explain why generic `PORT` environment variable is not supported
  (multi-plugin conflicts)
- **FR-005**: Documentation MUST describe the `PORT=<port>` stdout announcement format
- **FR-006**: Documentation MUST list all supported environment variables with their purposes
- **FR-007**: Documentation MUST explain the `ParsePortFlag()` function and requirement to call
  `flag.Parse()` first
- **FR-008**: Documentation MUST describe graceful shutdown behavior when context is canceled
- **FR-009**: Documentation MUST include a complete working example of a minimal plugin main()
  function
- **FR-010**: Documentation MUST describe the UnaryInterceptors configuration for custom middleware
- **FR-011**: Documentation MUST explain error return conditions (listener failure, gRPC server
  failure)

### Key Entities

- **Serve()**: Main entry point function that starts the gRPC server
- **ServeConfig**: Configuration struct with Plugin, Port, Registry, Logger, and UnaryInterceptors
  fields
- **ParsePortFlag()**: Helper function to retrieve --port command-line flag value
- **Environment Variables**: PULUMICOST_PLUGIN_PORT, PULUMICOST_LOG_LEVEL, PULUMICOST_LOG_FORMAT,
  PULUMICOST_LOG_FILE, PULUMICOST_TRACE_ID, PULUMICOST_TEST_MODE

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: A developer with Go experience can implement a working plugin main() function
  within 15 minutes using only the documentation
- **SC-002**: 100% of ServeConfig fields are documented with descriptions and valid value ranges
- **SC-003**: 100% of supported environment variables are documented with purposes and default
  behaviors
- **SC-004**: Documentation includes at least one complete, copy-paste-ready code example
- **SC-005**: All edge cases (invalid port, port conflict, permission errors) are documented
  with expected error behaviors
- **SC-006**: The documentation passes markdown linting without errors
