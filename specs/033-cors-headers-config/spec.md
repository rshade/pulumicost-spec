# Feature Specification: Configurable CORS Headers

**Feature Branch**: `033-cors-headers-config`
**Created**: 2026-01-04
**Status**: Draft
**Input**: GitHub Issue #228 - Make CORS headers configurable in pluginsdk WebConfig

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Custom Allowed Headers (Priority: P1)

As a plugin developer deploying to a security-conscious environment, I want to configure which
HTTP headers are allowed in cross-origin requests so that I can minimize the security footprint
by excluding headers my deployment doesn't use.

**Why this priority**: This is the core value proposition. The current hardcoded headers include
`Authorization` and `X-CSRF-Token` which may not be needed in all deployments, and some organizations
require minimal header exposure for compliance.

**Independent Test**: Can be fully tested by configuring custom allowed headers in WebConfig,
making a preflight OPTIONS request, and verifying the response contains only the specified headers.

**Acceptance Scenarios**:

1. **Given** a WebConfig with custom AllowedHeaders specified, **When** the server receives a
   preflight OPTIONS request, **Then** the `Access-Control-Allow-Headers` response header
   contains exactly the specified headers.

2. **Given** a WebConfig with AllowedHeaders set to nil, **When** the server receives a
   preflight OPTIONS request, **Then** the `Access-Control-Allow-Headers` response header
   contains the sensible default headers (maintaining backward compatibility).

3. **Given** a WebConfig with an empty slice for AllowedHeaders, **When** the server receives a
   preflight OPTIONS request, **Then** the `Access-Control-Allow-Headers` response header
   is empty or omitted.

---

### User Story 2 - Custom Exposed Headers (Priority: P2)

As a plugin developer, I want to configure which response headers are exposed to browser JavaScript
so that I can add observability headers like `X-Request-ID` or restrict exposed headers for compliance.

**Why this priority**: Exposed headers enable client-side observability (tracing, request correlation)
which is important for debugging and monitoring, but is secondary to controlling what headers are
accepted in requests.

**Independent Test**: Can be fully tested by configuring custom exposed headers, making an actual
CORS request, and verifying the `Access-Control-Expose-Headers` contains the specified headers.

**Acceptance Scenarios**:

1. **Given** a WebConfig with custom ExposedHeaders specified, **When** the server returns a
   CORS response, **Then** the `Access-Control-Expose-Headers` header contains exactly the
   specified headers.

2. **Given** a WebConfig with ExposedHeaders set to nil, **When** the server returns a CORS
   response, **Then** the default gRPC/Connect headers are exposed (backward compatibility).

3. **Given** a WebConfig with custom ExposedHeaders including `X-Request-ID`, **When** the server
   sets this header on the response, **Then** client-side JavaScript can read the `X-Request-ID`
   value.

---

### User Story 3 - Builder Method Configuration (Priority: P3)

As a plugin developer, I want to configure headers using the existing builder pattern so that my
configuration code remains consistent with how I configure other WebConfig options.

**Why this priority**: Developer ergonomics and API consistency are important for adoption, but
the feature works without builder methods (direct field assignment also works).

**Independent Test**: Can be tested by using `WithAllowedHeaders()` and `WithExposedHeaders()`
builder methods and verifying the resulting WebConfig has the expected values.

**Acceptance Scenarios**:

1. **Given** a default WebConfig, **When** I call `WithAllowedHeaders([]string{"Content-Type"})`,
   **Then** the returned WebConfig has AllowedHeaders set to `["Content-Type"]` and all other
   fields remain unchanged.

2. **Given** a WebConfig, **When** I chain multiple builder methods including header configuration,
   **Then** all configurations are applied correctly.

---

### Edge Cases

- What happens when AllowedHeaders contains duplicate header names?
  - The duplicates are included as-is (no deduplication). This mirrors browser behavior and avoids
    the overhead of deduplication.

- What happens when AllowedHeaders contains an empty string?
  - The empty string is included. Configuration validation is the caller's responsibility.

- How does the system handle case sensitivity in header names?
  - Header names are used as-is. HTTP headers are case-insensitive per spec, so browsers will
    match regardless of case.

- What happens when ExposedHeaders includes a forbidden header (e.g., `Set-Cookie`)?
  - The header is included in the response. Browsers enforce their own restrictions on which
    headers can actually be read.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: WebConfig MUST support an `AllowedHeaders` field of type `[]string` that specifies
  which headers are permitted in cross-origin requests.

- **FR-002**: WebConfig MUST support an `ExposedHeaders` field of type `[]string` that specifies
  which response headers are accessible to client-side JavaScript.

- **FR-003**: When `AllowedHeaders` is nil, the system MUST use the current default headers:
  `Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token,
  X-Requested-With, Connect-Protocol-Version, Connect-Timeout-Ms, Grpc-Timeout, X-Grpc-Web,
  X-User-Agent`.

- **FR-004**: When `ExposedHeaders` is nil, the system MUST use the current default headers:
  `Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin, Connect-Content-Encoding,
  Connect-Content-Type`.

- **FR-005**: The `WithAllowedHeaders([]string)` builder method MUST return a copy of WebConfig
  with the AllowedHeaders field set to a defensive copy of the input slice.

- **FR-006**: The `WithExposedHeaders([]string)` builder method MUST return a copy of WebConfig
  with the ExposedHeaders field set to a defensive copy of the input slice.

- **FR-007**: The corsMiddleware function MUST use the configured headers when generating
  the `Access-Control-Allow-Headers` and `Access-Control-Expose-Headers` response headers.

- **FR-008**: When AllowedHeaders is an empty slice (not nil), the system MUST set an empty
  `Access-Control-Allow-Headers` header, allowing only simple headers per CORS specification.

- **FR-009**: When ExposedHeaders is an empty slice (not nil), the system MUST set an empty
  `Access-Control-Expose-Headers` header.

- **FR-010**: The default headers MUST be documented in the WebConfig struct field comments and
  in the package documentation (README.md).

### Key Entities

- **WebConfig**: Configuration struct for gRPC-Web and CORS support. Extended with two new optional
  fields for header customization: `AllowedHeaders` and `ExposedHeaders`.

- **Default Header Sets**: Two predefined lists of headers used when custom headers are not
  specified. These represent the headers required for Connect/gRPC-Web protocol compatibility.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Plugin developers can configure custom CORS headers in under 2 lines of code using
  the builder pattern.

- **SC-002**: All existing plugins continue to function without modification (100% backward
  compatibility when not using new fields).

- **SC-003**: The feature is fully documented with examples covering common use cases: security-
  minimal configuration, observability headers, and default behavior.

- **SC-004**: Test coverage for the new functionality achieves at least 90% line coverage.

- **SC-005**: Performance impact on CORS middleware is negligible (less than 1 microsecond
  additional overhead per request).

## Assumptions

- Plugin developers are familiar with CORS headers and understand the implications of modifying them.
- The existing builder pattern style (`With*` methods that return copies) is the preferred API style.
- Defensive copying of slices is required to prevent external mutation of configuration.
- The current hardcoded headers represent a sensible default for Connect/gRPC-Web compatibility.
- Header validation (e.g., checking for valid header syntax) is out of scope for this feature.

## Dependencies

- This feature has no external dependencies beyond the existing pluginsdk package.
- Implementation depends on the existing `corsMiddleware` function and `WebConfig` struct.
