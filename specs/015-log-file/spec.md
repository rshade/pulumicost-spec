# Feature Specification: SDK Support for PULUMICOST_LOG_FILE

**Feature Branch**: `015-log-file`
**Created**: 2025-12-08
**Status**: Draft
**Input**: User description: "Plugins currently log to stderr, which pollutes the Core CLI
output. The SDK should check for `PULUMICOST_LOG_FILE` environment variable and configure
`zerolog` to write to that file if present. This enables a cleaner UX where Core controls
the log destination."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Core CLI Controls Log Destination (Priority: P1)

As a Core CLI developer, I want plugins to write logs to a file I specify so that the CLI
output remains clean and focused on user-facing information.

**Why this priority**: This is the primary use case - enabling Core to orchestrate multiple
plugins without log output mixing with CLI user interface. This directly solves the stated
problem of stderr pollution.

**Independent Test**: Can be fully tested by setting `PULUMICOST_LOG_FILE=/tmp/test.log`,
running a plugin, and verifying logs appear in the file while stderr remains clean.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_LOG_FILE` is set to a valid file path, **When** a plugin using the
   SDK logs messages, **Then** all log output is written to the specified file instead of
   stderr.
2. **Given** `PULUMICOST_LOG_FILE` is set to a valid file path, **When** a plugin logs at
   various levels (debug, info, warn, error), **Then** all levels are captured in the log
   file.
3. **Given** `PULUMICOST_LOG_FILE` is set, **When** the plugin starts, **Then** the log file
   is created if it does not exist or appended to if it already exists.

---

### User Story 2 - Default Behavior Without Environment Variable (Priority: P2)

As a plugin developer testing locally, I want logs to appear on stderr by default so that I
can see log output during development without additional configuration.

**Why this priority**: Maintains backward compatibility and ensures developers have a good
experience when running plugins standalone without Core orchestration.

**Independent Test**: Can be fully tested by running a plugin without setting
`PULUMICOST_LOG_FILE` and verifying logs appear on stderr as they do today.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_LOG_FILE` is not set, **When** a plugin using the SDK logs messages,
   **Then** all log output is written to stderr.
2. **Given** `PULUMICOST_LOG_FILE` is set to an empty string, **When** a plugin using the SDK
   logs messages, **Then** all log output is written to stderr (treating empty as unset).

---

### User Story 3 - Graceful Handling of Invalid Paths (Priority: P3)

As a Core CLI developer, I want plugins to handle invalid log file paths gracefully so that
plugin operation is not disrupted by configuration errors.

**Why this priority**: Error resilience is important for production deployments, but this is
less critical than the core functionality.

**Independent Test**: Can be fully tested by setting `PULUMICOST_LOG_FILE` to an invalid path
(e.g., directory that doesn't exist) and verifying the plugin logs an error and falls back
to stderr.

**Acceptance Scenarios**:

1. **Given** `PULUMICOST_LOG_FILE` points to a path in a non-existent directory, **When** a
   plugin starts, **Then** it logs a warning to stderr and falls back to stderr logging.
2. **Given** `PULUMICOST_LOG_FILE` points to a path without write permissions, **When** a
   plugin starts, **Then** it logs a warning to stderr and falls back to stderr logging.
3. **Given** `PULUMICOST_LOG_FILE` points to a directory instead of a file, **When** a
   plugin starts, **Then** it logs a warning to stderr and falls back to stderr logging.

---

### Edge Cases

- What happens when the log file path is a directory instead of a file? System should detect
  this and fall back to stderr with a warning.
- How does the system handle log file rotation? This is out of scope for MVP - file rotation
  is the responsibility of the deployer (e.g., logrotate).
- What happens when disk space is exhausted during logging? Standard OS behavior applies -
  zerolog will surface write errors naturally.
- What happens if multiple plugins write to the same log file simultaneously? The SDK should
  open the file in append mode with proper file locking semantics.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: SDK MUST check for the `PULUMICOST_LOG_FILE` environment variable at
  initialization time.
- **FR-002**: SDK MUST configure zerolog to write to the specified file path when
  `PULUMICOST_LOG_FILE` is set and non-empty.
- **FR-003**: SDK MUST configure zerolog to write to stderr when `PULUMICOST_LOG_FILE` is
  not set or is empty.
- **FR-004**: SDK MUST create the log file if it does not exist (with standard file
  permissions 0644).
- **FR-005**: SDK MUST append to the log file if it already exists (not truncate).
- **FR-006**: SDK MUST fall back to stderr logging when the specified file path is invalid
  or inaccessible, logging a warning about the fallback.
- **FR-007**: SDK MUST open log files in append mode to support multiple plugins writing to
  the same file.
- **FR-008**: SDK MUST provide a way for plugins to obtain a logger configured according to
  these rules without manual setup.

### Key Entities

- **LogConfiguration**: Represents the logging setup including output destination (file or
  stderr), derived from environment variables.
- **Logger**: The zerolog logger instance configured according to the log configuration,
  provided to plugin code.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Plugins using the SDK can redirect all log output to a file by setting a single
  environment variable.
- **SC-002**: No log output appears on stderr when `PULUMICOST_LOG_FILE` is set to a valid,
  writable path.
- **SC-003**: Plugin startup time is not noticeably impacted by log file configuration (less
  than 10ms additional latency).
- **SC-004**: Existing plugins continue to function without modification (backward
  compatibility - stderr logging when env var is unset).
- **SC-005**: Core CLI can cleanly capture plugin output streams without log interference,
  enabling structured user-facing output.

## Assumptions

- Plugins use zerolog for logging (as established in the SDK and documented in existing
  examples).
- The `PULUMICOST_LOG_FILE` environment variable name follows the established `PULUMICOST_*`
  naming convention.
- Log file rotation and cleanup are handled externally (e.g., by logrotate or the deploying
  system).
- The SDK initialization is the appropriate place to configure logging (single point of
  configuration).
- File permissions 0644 are appropriate for log files in this context.
