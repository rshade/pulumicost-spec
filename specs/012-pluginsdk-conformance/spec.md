# Feature Specification: PluginSDK Conformance Testing Adapters

**Feature Branch**: `012-pluginsdk-conformance`
**Created**: 2025-11-30
**Status**: Draft
**Input**: GitHub Issue #98 - feat(pluginsdk): Add conformance testing support for Plugin interface

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Run Basic Conformance Tests on Plugin (Priority: P1)

As a plugin developer using the `pluginsdk.Plugin` interface, I want to run basic conformance tests
against my plugin implementation with a single function call, so that I can validate my plugin meets
minimum PulumiCost specification requirements without manually converting to the raw gRPC interface.

**Why this priority**: This is the core value proposition. Plugin developers using the high-level
`Plugin` interface should not need to understand the underlying gRPC conversion to run conformance
tests. This removes the main friction point in the current developer experience.

**Independent Test**: Can be fully tested by creating a simple plugin implementation, calling the
adapter function, and verifying conformance results are returned with pass/fail status.

**Acceptance Scenarios**:

1. **Given** a plugin implementing `pluginsdk.Plugin`, **When** I call
   `pluginsdk.RunBasicConformance(plugin)`, **Then** I receive conformance results showing
   pass/fail status for all basic-level tests
2. **Given** a plugin with invalid PricingSpec responses, **When** I run basic conformance,
   **Then** I receive specific error messages identifying schema violations
3. **Given** a plugin that properly implements all methods, **When** I run basic conformance,
   **Then** all tests pass and the result indicates successful basic conformance

---

### User Story 2 - Run Standard Conformance Tests on Plugin (Priority: P1)

As a plugin developer preparing for production deployment, I want to run standard conformance tests
against my plugin, so that I can verify my plugin meets production-readiness requirements including
error handling and consistency guarantees.

**Why this priority**: Standard conformance is the target for most production plugins. This adapter
is equally important as basic conformance since most developers will aim for this level.

**Independent Test**: Can be tested by running standard conformance against both passing and failing
plugin implementations, verifying appropriate results and error details are returned.

**Acceptance Scenarios**:

1. **Given** a plugin implementing `pluginsdk.Plugin`, **When** I call
   `pluginsdk.RunStandardConformance(plugin)`, **Then** I receive conformance results including
   all basic and standard-level tests
2. **Given** a plugin that fails concurrency tests, **When** I run standard conformance, **Then**
   the result clearly indicates the concurrency category failed with specific details
3. **Given** a plugin that passes standard conformance, **When** I check the result, **Then** the
   conformance level achieved is clearly indicated as "Standard"

---

### User Story 3 - Run Advanced Conformance Tests on Plugin (Priority: P2)

As a plugin developer building high-performance plugins, I want to run advanced conformance tests
that include strict performance and scalability requirements, so that I can certify my plugin for
demanding production environments.

**Why this priority**: Advanced conformance is optional and targeted at specialized use cases. Most
plugins will not need this level, but the adapter should be available for those that do.

**Independent Test**: Can be tested by running advanced conformance against plugins with varying
performance characteristics, verifying the strict thresholds are properly evaluated.

**Acceptance Scenarios**:

1. **Given** a plugin implementing `pluginsdk.Plugin`, **When** I call
   `pluginsdk.RunAdvancedConformance(plugin)`, **Then** I receive conformance results including
   performance benchmarks with strict thresholds
2. **Given** a plugin that meets standard but not advanced requirements, **When** I run advanced
   conformance, **Then** the result shows standard level achieved with specific failures on
   advanced tests
3. **Given** a plugin that exceeds all thresholds, **When** I run advanced conformance, **Then**
   the result indicates "Advanced" conformance level achieved

---

### User Story 4 - Print Formatted Conformance Report (Priority: P2)

As a plugin developer reviewing test results, I want a formatted conformance report printed to my
test output, so that I can quickly understand which tests passed, which failed, and what actions
to take.

**Why this priority**: The report formatting improves developer experience but is secondary to the
core testing functionality. Developers can work with raw results if needed.

**Independent Test**: Can be tested by capturing test output after calling the print function and
verifying the format includes all expected sections and information.

**Acceptance Scenarios**:

1. **Given** conformance results from any level, **When** I call
   `pluginsdk.PrintConformanceReport(t, result)`, **Then** the test log shows a formatted report
   with pass/fail counts and details
2. **Given** results with some failed tests, **When** I print the report, **Then** failed tests
   are clearly highlighted with actionable error messages
3. **Given** results showing Standard conformance achieved, **When** I print the report, **Then**
   the conformance level is prominently displayed in the output

---

### Edge Cases

- **Nil plugin**: Functions return descriptive error immediately (FR-007)
- **Plugin panics**: Recovered by conformance harness; included in test failure results (FR-008)
- **gRPC server failure**: Error propagated from `NewServer()` call
- **Partial interface**: N/A - Go compiler enforces full interface implementation at compile time
- **Outside test context**: Functions work without `*testing.T`; use `PrintConformanceReport` for
  test output

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: Package MUST provide `RunBasicConformance(plugin)` function that accepts a
  `pluginsdk.Plugin` and returns conformance results
- **FR-002**: Package MUST provide `RunStandardConformance(plugin)` function for standard-level
  conformance testing
- **FR-003**: Package MUST provide `RunAdvancedConformance(plugin)` function for advanced-level
  conformance testing
- **FR-004**: All adapter functions MUST internally wrap the Plugin using `NewServer(plugin)` to
  create a `CostSourceServiceServer`
- **FR-005**: All adapter functions MUST delegate to the existing conformance functions in
  `sdk/go/testing` without duplicating test logic
- **FR-006**: Package MUST provide `PrintConformanceReport(t, result)` function for formatted
  output
- **FR-007**: Functions MUST return appropriate errors when passed nil plugins
- **FR-008**: Functions MUST handle plugin panics gracefully by delegating to the underlying
  conformance harness which uses `defer recover()` for panic recovery
- **FR-009**: The `ConformanceResult` type from `sdk/go/testing` MUST be re-exported or aliased
  for convenience
- **FR-010**: All adapter functions MUST be usable in standard test files without additional setup

### Key Entities

- **Plugin**: The high-level interface defined in `pluginsdk` that developers implement
- **Server**: The gRPC server wrapper that converts `Plugin` to `CostSourceServiceServer`
- **ConformanceResult**: The result structure from `sdk/go/testing` containing pass/fail status,
  test details, and conformance level achieved
- **ConformanceLevel**: Enumeration of Basic, Standard, Advanced certification levels

## Assumptions

- The existing `sdk/go/testing` package provides comprehensive conformance testing for
  `CostSourceServiceServer`
- The `pluginsdk.Server` type correctly implements `CostSourceServiceServer` by wrapping a `Plugin`
- Plugin developers are familiar with standard testing patterns
- Import cycles between `pluginsdk` and `sdk/go/testing` can be avoided using import aliases
- The existing conformance result types provide sufficient detail for plugin developers

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Plugin developers can run conformance tests with a single function call and no
  manual type conversion
- **SC-002**: Integration requires adding less than 10 lines of test code to an existing plugin
  test file
- **SC-003**: Test output clearly indicates conformance level achieved (Basic, Standard, or
  Advanced)
- **SC-004**: Failed tests include specific, actionable error messages identifying the violation
- **SC-005**: 100% of tests from `sdk/go/testing` conformance suite are exercised through the
  adapters
- **SC-006**: Adapter functions have zero impact on conformance test accuracy (same results as
  calling `sdk/go/testing` directly)
- **SC-007**: Documentation includes complete usage examples for all three conformance levels
