# Feature Specification: SDK Polish v0.4.15

**Feature Branch**: `034-sdk-polish`
**Created**: 2026-01-10
**Status**: Draft
**Input**: User description: "SDK Polish: configurable timeouts, user-friendly errors,
performance tests"

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Client Timeout Configuration (Priority: P1)

A plugin host application needs to configure per-request timeouts when calling plugins
that may have long-running operations (e.g., `GetActualCost` with large date ranges
spanning months of data).

**Why this priority**: Long-running cost queries can block indefinitely without proper
timeout handling. This directly impacts host application reliability and user experience
when querying historical cost data.

**Independent Test**: Can be fully tested by creating a client with a custom timeout,
making a request to a slow mock server, and verifying the request times out as expected.

**Acceptance Scenarios**:

1. **Given** a `ClientConfig` with `Timeout` set to 5 seconds,
   **When** a plugin RPC takes longer than 5 seconds,
   **Then** the request fails with a context deadline exceeded error.
2. **Given** a `ClientConfig` using `WithTimeout(10 * time.Second)`,
   **When** a plugin RPC completes within 10 seconds,
   **Then** the request succeeds normally.
3. **Given** a context with a deadline shorter than `ClientConfig.Timeout`,
   **When** a plugin RPC is called,
   **Then** the context deadline takes precedence.
4. **Given** the default `ClientConfig` with no explicit timeout,
   **When** a client is created,
   **Then** the default 30-second timeout is applied.

---

### User Story 2 - User-Friendly GetPluginInfo Error Messages (Priority: P2)

A host application developer encounters error messages when plugins return invalid
metadata. The error messages should be user-friendly and actionable rather than exposing
internal implementation details.

**Why this priority**: Clear error messages reduce debugging time and improve the
developer experience. Technical error messages like "plugin returned nil response" expose
implementation details and provide poor guidance.

**Independent Test**: Can be fully tested by configuring a mock plugin to return
nil/incomplete/invalid responses and verifying the returned error messages match the
user-friendly format.

**Acceptance Scenarios**:

1. **Given** a plugin that returns a nil response from `GetPluginInfo`,
   **When** a host calls `GetPluginInfo`,
   **Then** the error message is "unable to retrieve plugin metadata".
2. **Given** a plugin that returns incomplete metadata (missing name, version, or
   spec_version),
   **When** a host calls `GetPluginInfo`,
   **Then** the error message is "plugin metadata is incomplete".
3. **Given** a plugin that returns an invalid spec_version format,
   **When** a host calls `GetPluginInfo`,
   **Then** the error message is "plugin reported an invalid specification version".

---

### User Story 3 - GetPluginInfo Performance Conformance (Priority: P3)

A plugin developer wants to ensure their `GetPluginInfo` implementation meets performance
requirements. The conformance test suite should validate that `GetPluginInfo` responds
within acceptable latency bounds.

**Why this priority**: Performance conformance ensures plugins meet operational
requirements. `GetPluginInfo` is called during plugin discovery and should be fast
(under 100ms) since it requires no external API calls.

**Independent Test**: Can be fully tested by running the performance conformance test
against any plugin implementation and verifying all iterations complete within 100ms.

**Acceptance Scenarios**:

1. **Given** a plugin implementation,
   **When** the `GetPluginInfo` conformance test runs 10 iterations,
   **Then** all iterations complete within 100ms each.
2. **Given** a slow `GetPluginInfo` implementation taking 150ms,
   **When** the conformance test runs,
   **Then** the test fails with a clear message indicating which iteration exceeded
   the threshold.
3. **Given** a plugin that does not implement `GetPluginInfo`,
   **When** the conformance test runs,
   **Then** the test handles the `Unimplemented` error gracefully (legacy plugin support).

---

### Edge Cases

- **Zero Timeout Behavior**: When ClientConfig.Timeout is set to 0, the default 30-second timeout is applied.
  **Test Criteria**: Client times out after exactly 30 seconds when no custom timeout is set.
- **Custom HTTPClient Precedence**: When both HTTPClient and Timeout are set in ClientConfig, HTTPClient.Timeout
  takes precedence. **Test Criteria**: RPC calls respect HTTPClient.Timeout value over ClientConfig.Timeout.
- **Empty Strings in Metadata**: When GetPluginInfo returns empty strings for required fields, it's treated as
  incomplete metadata. **Test Criteria**: Empty strings trigger "plugin metadata is incomplete" error.
- **Connection Errors in Performance Tests**: When performance test cannot connect to the plugin, it fails with
  connection error (not performance failure). **Test Criteria**: Test exits with connection error before
  measuring latency.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The SDK MUST support configurable per-client timeouts via
  `ClientConfig.Timeout` field.
- **FR-002**: The SDK MUST provide a `WithTimeout(duration)` method on `ClientConfig`
  for fluent configuration.
- **FR-003**: The SDK MUST respect context deadlines over client-level timeouts when
  both are set.
- **FR-004**: The SDK MUST use the default timeout (30 seconds) when `Timeout` is zero
  and no custom `HTTPClient` is provided.
- **FR-005**: The SDK MUST properly identify and wrap context timeout errors via
  `wrapRPCError`.
- **FR-006**: The `GetPluginInfo` RPC MUST return "unable to retrieve plugin metadata"
  when the plugin returns nil.
- **FR-007**: The `GetPluginInfo` RPC MUST return "plugin metadata is incomplete" when
  required fields are empty.
- **FR-008**: The `GetPluginInfo` RPC MUST return "plugin reported an invalid
  specification version" for malformed versions.
- **FR-009**: The testing framework MUST include a `GetPluginInfoPerformance`
  conformance test.
- **FR-010**: The performance conformance test MUST run 10 iterations and fail if any
  exceeds 100ms.
- **FR-011**: The performance conformance test MUST handle legacy plugins that return
  `Unimplemented` gracefully.

### Key Entities

- **ClientConfig**: Configuration struct for creating SDK clients, including timeout
  settings and HTTP client options.
- **Server**: gRPC server wrapper that handles `GetPluginInfo` validation and error
  message formatting.
- **ConformanceTest**: Test case structure in the testing framework that validates
  plugin behavior and performance.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All timeout-related tests pass, demonstrating correct timeout behavior
  at client and context levels.
- **SC-002**: Error messages returned by `GetPluginInfo` match the user-friendly format
  specified in FR-006, FR-007, and FR-008.
- **SC-003**: The `GetPluginInfoPerformance` conformance test passes for all mock
  plugin implementations.
- **SC-004**: No regressions in existing client or server functionality (all existing
  tests continue to pass).
- **SC-005**: `GetPluginInfo` conformance test completes in under 2 seconds total
  (10 iterations × 100ms + overhead).

## Assumptions

- The existing `ClientConfig.Timeout` field and `WithTimeout()` method are correctly
  implemented but may need verification/documentation.
- The existing `wrapRPCError` function correctly identifies `context.DeadlineExceeded`
  errors.
- Plugin developers understand that `GetPluginInfo` should not make external API calls
  (hence the 100ms requirement).
- The conformance test framework already supports adding new test cases without
  architectural changes.

## Dependencies

- No new external dependencies required.
- Changes are isolated to `sdk/go/pluginsdk/client.go`, `sdk/go/pluginsdk/sdk.go`,
  and `sdk/go/testing/conformance_test.go`.
- Related GitHub issues: #226 (timeouts), #244 (performance test), #245 (error messages).
