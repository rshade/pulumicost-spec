# Feature Specification: Centralized Environment Variable Handling

**Feature Branch**: `013-pluginsdk-env`
**Created**: 2025-12-07
**Status**: Draft
**Input**: User description: "feat(pluginsdk): Add centralized environment variable handling"

## Clarifications

### Session 2025-12-07

- Q: Should `GetLogLevel()` fall back to reading `LOG_LEVEL` if `PULUMICOST_LOG_LEVEL` is not
  set? → A: Yes, add fallback: `PULUMICOST_LOG_LEVEL` → `LOG_LEVEL` (matches port pattern)

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Developer Uses Standard Port Configuration (Priority: P1)

A plugin developer building a new PulumiCost plugin needs their plugin to automatically bind to
the correct gRPC port. They set `PULUMICOST_PLUGIN_PORT=8080` in their environment and their
plugin correctly binds to port 8080 without any code changes.

**Why this priority**: This is the core problem that caused E2E testing failures. Plugin port
configuration must work consistently across all plugins to enable reliable plugin orchestration
by pulumicost-core.

**Independent Test**: Can be fully tested by setting `PULUMICOST_PLUGIN_PORT=8080`, starting a
plugin, and verifying it binds to port 8080. Delivers immediate value by fixing the port
mismatch issue.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_PLUGIN_PORT=8080` is set, **When** a plugin starts without explicit
   port configuration, **Then** it binds to port 8080
2. **Given** `PULUMICOST_PLUGIN_PORT` is not set, **When** a plugin starts without explicit
   port, **Then** it returns a clear error indicating which environment variable to set
3. **Given** `PULUMICOST_PLUGIN_PORT` contains an invalid value, **When** a plugin starts,
   **Then** it returns 0 (caller handles error)

---

### User Story 2 - Plugin Developer Configures Logging via Environment (Priority: P2)

A plugin developer deploying to production wants to control logging behavior without code
changes. They set `PULUMICOST_LOG_LEVEL=debug`, `PULUMICOST_LOG_FORMAT=json`, and optionally
`PULUMICOST_LOG_FILE=/var/log/plugin.log` to get structured debug logs suitable for log
aggregation systems.

**Why this priority**: Logging configuration is essential for production deployments but not
critical for basic plugin functionality. Plugins can function with default logging.

**Independent Test**: Can be tested by setting log environment variables, starting a plugin,
and verifying log output format, verbosity, and file destination match expectations.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_LOG_LEVEL=debug` is set, **When** a plugin reads logging
   configuration, **Then** the plugin can access the "debug" log level value
2. **Given** `PULUMICOST_LOG_FORMAT=json` is set, **When** a plugin reads logging
   configuration, **Then** the plugin can access the "json" format value
3. **Given** `PULUMICOST_LOG_FILE=/var/log/plugin.log` is set, **When** a plugin reads logging
   configuration, **Then** the plugin can access the file path for log output
4. **Given** no logging environment variables are set, **When** a plugin reads logging
   configuration, **Then** empty strings are returned (plugin uses its defaults)

---

### User Story 3 - Operations Team Traces Requests Across Services (Priority: P2)

An operations team troubleshooting a production issue wants to trace cost calculations across
the plugin ecosystem. They inject `PULUMICOST_TRACE_ID=abc123` when invoking a plugin, and
this trace ID appears in all related logs and responses.

**Why this priority**: Distributed tracing is valuable for debugging but not essential for
basic plugin operation. Plugins work correctly without trace IDs.

**Independent Test**: Can be tested by setting `PULUMICOST_TRACE_ID`, invoking a plugin
operation, and verifying the trace ID is accessible for correlation.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_TRACE_ID=abc123` is set, **When** a plugin reads trace configuration,
   **Then** the plugin can access "abc123" as the trace ID
2. **Given** no trace ID is set, **When** a plugin reads trace configuration, **Then** an
   empty string is returned (no trace context)

---

### User Story 4 - Plugin Developer Enables Test Mode (Priority: P2)

A plugin developer writing integration tests needs to enable test mode to use mock data or
bypass external service calls. They set `PULUMICOST_TEST_MODE=true` and the plugin operates
in test mode with deterministic behavior.

**Why this priority**: Test mode is essential for reliable plugin testing but not required
for production operation. Plugins must support both modes.

**Independent Test**: Can be tested by setting `PULUMICOST_TEST_MODE=true`, invoking plugin
operations, and verifying test-specific behavior is active.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_TEST_MODE=true` is set, **When** a plugin reads test configuration,
   **Then** the plugin can detect that test mode is enabled
2. **Given** `PULUMICOST_TEST_MODE=false` or unset, **When** a plugin reads test
   configuration, **Then** the plugin operates in normal production mode
3. **Given** `PULUMICOST_TEST_MODE` contains an invalid value (not "true"/"false"), **When**
   a plugin reads test configuration, **Then** a warning is logged and test mode is disabled

---

### User Story 5 - Plugin Developer Migrates Logging Configuration (Priority: P3)

A developer with an existing plugin that uses direct `os.Getenv("LOG_LEVEL")` calls wants to
migrate to the standardized environment handling. They update their imports to use
`pluginsdk.GetLogLevel()` and their plugin continues to work with both old (`LOG_LEVEL`) and
new (`PULUMICOST_LOG_LEVEL`) environment variable names.

**Why this priority**: Migration support ensures backward compatibility for logging
configuration. Port configuration does NOT have fallback - plugins must use
`PULUMICOST_PLUGIN_PORT` exclusively.

**Independent Test**: Can be tested by setting `LOG_LEVEL=debug` (without PULUMICOST_ prefix),
calling `pluginsdk.GetLogLevel()`, and verifying it returns "debug".

**Acceptance Scenarios**:

1. **Given** an existing plugin using `os.Getenv("LOG_LEVEL")`, **When** migrated to use
   `pluginsdk.GetLogLevel()`, **Then** the plugin works identically with `LOG_LEVEL` set
2. **Given** a migrated plugin, **When** operator switches to `PULUMICOST_LOG_LEVEL`,
   **Then** the plugin works without code changes
3. **Given** both `PULUMICOST_LOG_LEVEL` and `LOG_LEVEL` are set, **When** plugin reads
   config, **Then** `PULUMICOST_LOG_LEVEL` takes precedence

---

### Edge Cases

- What happens when `PULUMICOST_PLUGIN_PORT` contains a non-numeric value like "abc"? The
  value is treated as invalid, and `GetPort()` returns 0 (no fallback).
- What happens when port value is negative or zero? Non-positive values are treated as
  invalid; `GetPort()` returns 0.
- What happens when port value exceeds valid port range (65535)? Value is accepted as-is;
  network layer will reject during binding (consistent with standard Go behavior).
- What happens when environment variable contains leading/trailing whitespace? Standard
  `os.Getenv` returns the value as-is including whitespace, which will cause parse failures
  for numeric values.
- What happens when `PULUMICOST_TEST_MODE` is set to an unexpected value like "yes"? A
  warning is logged and test mode defaults to disabled for safety.
- What happens when `PULUMICOST_LOG_FILE` path is not writable? Plugin should fail gracefully
  with a clear error message (implementation responsibility).

## Requirements _(mandatory)_

### Functional Requirements

#### Port Configuration

- **FR-001**: Plugin SDK MUST provide exported constants for all standard PulumiCost
  environment variable names
- **FR-002**: Plugin SDK MUST provide a `GetPort()` function that reads
  `PULUMICOST_PLUGIN_PORT` only (no fallback to PORT)
- **FR-003**: `GetPort()` MUST return 0 when no valid port is configured (caller handles
  error)
- **FR-004**: `GetPort()` MUST reject non-positive port values as invalid
- **FR-005**: `pluginsdk.Serve()` MUST use the centralized `GetPort()` function instead of
  direct environment variable access
- **FR-006**: Error messages for missing port configuration MUST explicitly name
  `PULUMICOST_PLUGIN_PORT` as the required variable

#### Logging Configuration

- **FR-007**: Plugin SDK MUST provide `GetLogLevel()` function that reads
  `PULUMICOST_LOG_LEVEL` first, then falls back to `LOG_LEVEL`, returning empty string if
  neither is set
- **FR-008**: Plugin SDK MUST provide `GetLogFormat()` function returning the configured log
  format or empty string
- **FR-009**: Plugin SDK MUST provide `GetLogFile()` function returning the configured log
  file path or empty string

#### Tracing Configuration

- **FR-010**: Plugin SDK MUST provide `GetTraceID()` function returning the configured trace
  ID or empty string

#### Test Mode Configuration

- **FR-011**: Plugin SDK MUST provide `GetTestMode()` function returning true only when
  `PULUMICOST_TEST_MODE` is explicitly set to "true"
- **FR-012**: Plugin SDK MUST provide `IsTestMode()` convenience function returning boolean
- **FR-013**: `GetTestMode()` MUST log a warning when `PULUMICOST_TEST_MODE` contains an
  invalid value (not "true" or "false")

### Key Entities

- **Environment Variable Constants**: Named constants providing the canonical variable names:
  - `EnvPort` = `PULUMICOST_PLUGIN_PORT`
  - `EnvLogLevel` = `PULUMICOST_LOG_LEVEL`
  - `EnvLogLevelFallback` = `LOG_LEVEL`
  - `EnvLogFormat` = `PULUMICOST_LOG_FORMAT`
  - `EnvLogFile` = `PULUMICOST_LOG_FILE`
  - `EnvTraceID` = `PULUMICOST_TRACE_ID`
  - `EnvTestMode` = `PULUMICOST_TEST_MODE`
- **Port Configuration**: Numeric port value with validation (no fallback)
- **Logging Configuration**: String-based log level (with fallback), format, and file path
- **Trace Context**: Optional trace ID string for distributed tracing
- **Test Mode**: Boolean flag for enabling test-specific plugin behavior

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All plugins using `pluginsdk.GetPort()` correctly bind to the port set via
  `PULUMICOST_PLUGIN_PORT` on first attempt
- **SC-002**: Plugin developers can configure all standard environment variables by
  referencing a single source of truth in `pluginsdk/env.go`
- **SC-003**: E2E tests between pulumicost-core and plugins pass when core sets
  `PULUMICOST_PLUGIN_PORT` (fixing the original issue)
- **SC-004**: Documentation clearly shows plugin developers how to use environment
  configuration functions
- **SC-005**: Existing plugins using `LOG_LEVEL` can migrate to `PULUMICOST_LOG_LEVEL` with
  a clear deprecation path (logging fallback only)

## Assumptions

- The `pluginsdk` package already exists and is the appropriate location for this
  functionality
- Plugins are responsible for their own default values when environment variables are not set
- Port validation beyond positive integer check is delegated to the network layer
- Log level and format values are passed through as strings; semantic validation is the
  plugin's responsibility
- Test mode validation (strict "true"/"false") follows the pattern in pulumicost-plugin-
  aws-public

## Cross-Repository Consistency Notes

Based on analysis of pulumicost-core, pulumicost-plugin-aws-public, and pulumicost-plugin-
aws-ce:

### Variables Used Consistently

- `PULUMICOST_PLUGIN_PORT` - Core sets this, plugins should read via `GetPort()`
- `PULUMICOST_TRACE_ID` - Core and plugins use for distributed tracing
- `PULUMICOST_LOG_LEVEL` - Core uses this; plugins should adopt
- `PULUMICOST_LOG_FORMAT` - Core uses this; plugins should adopt

### Variables Requiring Migration

- `LOG_LEVEL` in pulumicost-plugin-aws-public should migrate to `PULUMICOST_LOG_LEVEL`
  (fallback supported for backward compatibility)
- `PORT` usage must be replaced with `PULUMICOST_PLUGIN_PORT` (no fallback - breaking change)

### Core-Only Variables (Not in Plugin SDK)

The following are used by pulumicost-core only and should NOT be in the plugin SDK:

- `PULUMICOST_OUTPUT_FORMAT` - CLI output formatting
- `PULUMICOST_OUTPUT_PRECISION` - Decimal precision for CLI output
- `PULUMICOST_CONFIG_STRICT` - Config file error handling
- `PULUMICOST_CONFIG` - Config file path
- `PULUMICOST_PLUGIN_<NAME>_<KEY>` - Plugin-specific config (passed by core)
- `GITHUB_TOKEN` - Registry authentication
