# Feature Specification: SDK Documentation Consolidation

**Feature Branch**: `031-sdk-docs-consolidation`
**Created**: 2025-12-31
**Status**: Draft
**Input**: User description: "Consolidate and improve SDK documentation by addressing
12 open documentation issues covering inline code comments, godoc examples,
performance tuning guides, CORS best practices, migration guides, rate limiting
patterns, and thread safety documentation for the pluginsdk package"

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Plugin Developer Learning the SDK (Priority: P1)

A new plugin developer wants to understand how to properly use the pluginsdk to
build their first cost source plugin. They need clear, comprehensive documentation
that explains patterns, best practices, and common pitfalls without having to
reverse-engineer the code.

**Why this priority**: New developers are the primary audience for documentation.
Poor documentation leads to incorrect usage, bugs, and increased support burden.
Getting developers started correctly is foundational.

**Independent Test**: Can be fully tested by having a new developer follow the
documentation to build a working plugin without external assistance.

**Acceptance Scenarios**:

1. **Given** a developer reading the pluginsdk README,
   **When** they look for client usage patterns,
   **Then** they find a complete example showing `NewClient()`, `defer client.Close()`,
   and typical operations
2. **Given** a developer reading godoc for `Client.Close()`,
   **When** they view the documentation,
   **Then** they see a clear example demonstrating the defer pattern
3. **Given** a developer reading `NewClient()` documentation,
   **When** they examine HTTP client configuration options,
   **Then** they understand ownership semantics (who manages the client lifecycle)

---

### User Story 2 - Developer Migrating from gRPC to Connect-go (Priority: P1)

An existing plugin developer has a working pure gRPC plugin and wants to migrate
to the new connect-go multi-protocol support for browser compatibility and simpler
HTTP/1.1 fallback.

**Why this priority**: Migration documentation is critical for existing users
adopting new capabilities. The connect-go integration is a major feature that
enables gRPC-Web browser support.

**Independent Test**: Can be tested by migrating an existing gRPC-only plugin to
connect-go following only the documentation.

**Acceptance Scenarios**:

1. **Given** a developer with an existing gRPC plugin,
   **When** they follow the migration guide,
   **Then** they can convert their server to support both gRPC and Connect protocols
2. **Given** a developer reading the migration guide,
   **When** they reach the protocol selection section,
   **Then** they understand when to use gRPC vs Connect vs gRPC-Web based on deployment
3. **Given** a developer migrating their client code,
   **When** they follow the client migration section,
   **Then** they can switch from grpc.Dial to the new multi-protocol Client

---

### User Story 3 - Operator Configuring CORS for Production (Priority: P2)

A DevOps engineer or operator needs to configure CORS settings for a plugin server
deployed in various environments (local development, behind an API gateway,
multi-tenant SaaS).

**Why this priority**: CORS misconfiguration is a common source of production issues.
Clear guidance prevents security vulnerabilities and debugging headaches.

**Independent Test**: Can be tested by following the guide to configure CORS for a
specific deployment scenario and verifying browser access works correctly.

**Acceptance Scenarios**:

1. **Given** an operator deploying to a single-origin frontend,
   **When** they read the CORS guide,
   **Then** they find a configuration example matching their scenario
2. **Given** an operator debugging CORS errors in the browser,
   **When** they consult the troubleshooting section,
   **Then** they can identify and fix common issues (missing headers, preflight failures)
3. **Given** an operator concerned about security,
   **When** they read the security guidelines,
   **Then** they understand the risks of wildcard origins and proper credential handling

---

### User Story 4 - Developer Optimizing Client Performance (Priority: P2)

A developer building a high-throughput application needs to tune the HTTP client
configuration for optimal performance when making many concurrent requests to
multiple plugins.

**Why this priority**: Performance tuning documentation prevents developers from
using incorrect defaults that cause connection pool exhaustion or excessive memory
usage.

**Independent Test**: Can be tested by applying the documented configuration for a
high-throughput scenario and measuring improved performance.

**Acceptance Scenarios**:

1. **Given** a developer unsure which config to use,
   **When** they consult the decision matrix,
   **Then** they can select between DefaultClientConfig and HighThroughputClientConfig
2. **Given** a developer needing custom tuning,
   **When** they read the parameter explanations,
   **Then** they understand MaxIdleConns, MaxIdleConnsPerHost, and IdleConnTimeout
3. **Given** a developer experiencing connection issues,
   **When** they follow monitoring tips,
   **Then** they can diagnose connection pool problems

---

### User Story 5 - Developer Implementing Rate Limiting (Priority: P2)

A plugin developer needs to implement rate limiting to avoid overwhelming cloud
provider APIs that have strict rate limits (AWS, Azure, GCP).

**Why this priority**: Rate limiting is essential for production plugins. Poor rate
limiting leads to API throttling, degraded performance, and potential service bans.

**Independent Test**: Can be tested by implementing rate limiting following the
documented patterns and verifying proper backoff behavior under load.

**Acceptance Scenarios**:

1. **Given** a developer building an AWS plugin,
   **When** they read the rate limiting section,
   **Then** they find the token bucket pattern with golang.org/x/time/rate
2. **Given** a developer handling rate limit errors,
   **When** they follow the error handling guide,
   **Then** they correctly return gRPC ResourceExhausted status codes
3. **Given** a developer implementing retry logic,
   **When** they read the backoff section,
   **Then** they implement exponential backoff with jitter

---

### User Story 6 - Developer Understanding Concurrency Guarantees (Priority: P3)

A developer building a multi-threaded application needs to understand which SDK
components are safe for concurrent use and which require synchronization.

**Why this priority**: Thread safety documentation prevents race conditions and data
corruption bugs that are difficult to diagnose.

**Independent Test**: Can be tested by running concurrent access tests with the
`-race` flag and verifying documented thread-safe components work correctly.

**Acceptance Scenarios**:

1. **Given** a developer using Client from multiple goroutines,
   **When** they check the thread safety docs,
   **Then** they confirm Client is safe for concurrent use
2. **Given** a developer configuring WebConfig,
   **When** they read the concurrency notes,
   **Then** they understand it's safe only if used immutably with the builder pattern
3. **Given** a plugin implementer,
   **When** they read the thread safety section,
   **Then** they understand their implementation must handle its own synchronization

---

### User Story 7 - Developer Understanding Code Semantics (Priority: P3)

A developer needs to understand the precise semantics of SDK functions to avoid
subtle bugs (ownership transfer, reference semantics, tolerance constants).

**Why this priority**: Semantic documentation prevents subtle bugs that are difficult
to diagnose. These are "quick wins" that improve code quality.

**Independent Test**: Can be tested by reading the inline comments and understanding
the behavior without consulting additional resources.

**Acceptance Scenarios**:

1. **Given** a developer using `WithTags()`,
   **When** they read the function comment,
   **Then** they understand tags are copied into the internal map (zero-allocation pattern)
2. **Given** a developer reading `contractedCostTolerance`,
   **When** they see the comment,
   **Then** they understand why 0.0001 (1 basis point) is used for IEEE 754 comparisons
3. **Given** a developer reading correlation pattern comments,
   **When** they see `ResourceRecommendationInfo.id`,
   **Then** the field name matches the actual proto definition

---

### Edge Cases

- What happens when a developer uses documentation for a different SDK version?
  - **Resolution**: Documentation targets current SDK version; version-specific notes added where
    behavior differs. CHANGELOG.md provides migration guidance for version transitions.
- How does the documentation handle deprecated APIs (e.g., provider_name vs
  service_provider_name)?
  - **Resolution**: Deprecated APIs documented with migration paths per FOCUS 1.3 patterns in
    research.md Section 6. Both old and new field names shown in examples.
- What if a developer's deployment scenario doesn't match any documented CORS pattern?
  - **Resolution**: CORS guide includes "Custom Configuration" subsection with builder pattern
    examples and security principles that apply to any scenario.

## Requirements _(mandatory)_

### Functional Requirements

#### Inline Code Documentation (Quick Wins)

- **FR-001**: System MUST include a comment in `NewClient()` explaining HTTP client
  ownership semantics (Issue #240)
- **FR-002**: System MUST include a comprehensive comment for `contractedCostTolerance`
  explaining the 1 basis point (0.0001) rationale, IEEE 754 considerations, and
  concrete examples (Issue #211)
- **FR-003**: System MUST include a comment for `WithTags` documenting shared-map
  semantics (zero-allocation, reference assignment) (Issue #207)
- **FR-004**: System MUST fix correlation pattern comment to reference
  `ResourceRecommendationInfo.id` instead of `resource_id` (Issue #206)
- **FR-005**: System MUST clarify test naming in `resource_id_test.go` to indicate
  round-trip semantics vs literal old-server interoperability (Issue #208)

#### Godoc Examples

- **FR-006**: System MUST include a godoc example for `Client.Close()` showing the
  `defer client.Close()` pattern with typical operations (Issue #238)
- **FR-007**: System MUST include complete import statements in
  `sdk/go/testing/README.md` examples, particularly the `pbc` package alias (Issue #209)

#### Comprehensive Guides

- **FR-008**: System MUST provide a Migration Guide from pure gRPC to connect-go
  covering server migration, client migration, protocol selection, testing, and
  backward compatibility (Issue #235)
- **FR-009**: System MUST provide a Performance Tuning Guide with configuration
  selection matrix, parameter explanations, custom examples, and monitoring tips
  (Issue #237)
- **FR-010**: System MUST provide a CORS Best Practices Guide covering 5+ deployment
  scenarios, security guidelines, debugging, and header reference (Issue #236)
- **FR-011**: System MUST provide Rate Limiting documentation covering token bucket
  pattern, cloud provider limits, backoff strategies, and proper gRPC status codes
  (Issue #233)
- **FR-012**: System MUST document thread safety guarantees for Server, Client,
  WebConfig, and PluginMetrics structs (Issue #231)

### Key Entities

- **pluginsdk README**: Primary documentation file for SDK usage
  (`sdk/go/pluginsdk/README.md`)
- **testing README**: Testing framework documentation (`sdk/go/testing/README.md`)
- **Inline Comments**: Code-level documentation in Go source files
- **Proto Comments**: Documentation in protobuf definitions

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All 12 documentation issues (#206, #207, #208, #209, #211, #231, #233,
  #235, #236, #237, #238, #240) are resolved and closed
- **SC-002**: All code examples in README files compile without modification
  (copy-paste ready)
- **SC-003**: All markdown files pass `make lint-markdown` validation
- **SC-004**: All Go code changes pass `make lint` and `make test`
- **SC-005**: Migration guide enables developers to migrate a gRPC plugin to
  connect-go without external assistance
- **SC-006**: New plugin developers can understand client lifecycle management from
  documentation alone
- **SC-007**: CORS guide covers at least 5 distinct deployment scenarios with
  working examples
- **SC-008**: Performance tuning guide includes a decision matrix for configuration
  selection
- **SC-009**: Thread safety documentation includes concurrent access tests with
  `-race` flag verification

## Dependencies & Assumptions

### Dependencies

- Issue #228 (CORS headers configurable) and #229 (CORS max-age configurable)
  should be implemented before CORS Best Practices Guide is finalized
- Existing pluginsdk implementation is stable and documented behavior matches
  actual behavior

### Assumptions

- Documentation will be written in markdown and Go godoc format
- Examples will use the current SDK API (no upcoming breaking changes)
- Target audience has intermediate Go programming knowledge
- Documentation follows existing repository conventions for structure and formatting
