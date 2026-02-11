# Feature Specification: GetActualCost Pagination Support

**Feature Branch**: `044-actual-cost-pagination`
**Created**: 2026-02-04
**Status**: Draft
**Input**: GitHub Issue #353 - Add pagination support to GetActualCost RPC

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Paginated Cost Retrieval for Large Datasets (Priority: P1)

As a FinFocus host application, I need to retrieve large volumes of actual cost records
in manageable pages so that I can process datasets with 10,000+ records without exceeding
message size limits or exhausting memory.

**Why this priority**: This is the core capability. Without pagination, plugins returning
large datasets (e.g., 30-day queries with hourly granularity across hundreds of resources)
risk exceeding the 4MB gRPC message size limit, causing complete request failures. This
directly addresses the reliability and scalability gap in the current single-response model.

**Independent Test**: Can be fully tested by requesting actual cost data with a specified
page size and verifying that the response contains at most that many records, along with
a continuation token when more records exist.

**Acceptance Scenarios**:

1. **Given** a plugin with 500 cost records available for a query, **When** the host
   requests actual costs with a page size of 100, **Then** the response contains exactly
   100 records and a non-empty continuation token.
2. **Given** a plugin with 500 cost records and the host has retrieved the first 4 pages
   (400 records), **When** the host requests the 5th page using the continuation token,
   **Then** the response contains 100 records and an empty continuation token indicating
   no more pages.
3. **Given** a plugin with 50 cost records, **When** the host requests actual costs with
   a page size of 100, **Then** the response contains all 50 records and an empty
   continuation token.

---

### User Story 2 - Backward-Compatible Default Behavior (Priority: P1)

As a plugin developer who has not yet implemented pagination, I need the system to continue
working without any changes to my existing implementation so that existing plugins are not
broken by the protocol update.

**Why this priority**: Backward compatibility is equally critical as the core feature.
Breaking existing plugins would undermine the ecosystem's stability and trust.

**Independent Test**: Can be tested by running an existing (non-paginated) plugin against
a host that sends pagination parameters and verifying the plugin returns all records in a
single response with an empty continuation token.

**Acceptance Scenarios**:

1. **Given** a plugin that does not implement pagination logic, **When** a host sends a
   request without pagination parameters (page size = 0, no token), **Then** the plugin
   returns all available records in a single response with an empty continuation token.
2. **Given** a plugin that does not implement pagination logic, **When** a host sends a
   request with a page size of 100, **Then** the plugin returns all available records
   (regardless of page size) with an empty continuation token, since it has no awareness
   of the pagination fields.
3. **Given** a host that does not use pagination, **When** it calls GetActualCost without
   setting page size or token, **Then** it receives all records as it does today with no
   behavioral change.

---

### User Story 3 - SDK Pagination Helpers for Plugin Developers (Priority: P2)

As a plugin developer, I need ready-to-use pagination helper functions so that I can
implement paginated responses without building pagination logic from scratch.

**Why this priority**: While plugins can implement pagination manually, providing SDK
helpers reduces implementation effort and ensures consistency across the plugin ecosystem.
This follows the established pattern from GetRecommendations.

**Independent Test**: Can be tested by calling the pagination helper with a slice of cost
records, a page size, and a page token, and verifying it returns the correct page of
results with the appropriate next token.

**Acceptance Scenarios**:

1. **Given** a plugin developer with 1,000 cost records in memory, **When** they call the
   pagination helper with page size 100 and an empty token, **Then** they receive the first
   100 records and a token pointing to record 101.
2. **Given** a continuation token from a previous response, **When** the plugin developer
   calls the pagination helper with that token, **Then** they receive the next page of
   records starting from the correct offset.
3. **Given** an invalid or corrupted continuation token, **When** the plugin developer
   calls the pagination helper, **Then** they receive a clear error indicating the token
   is malformed.

---

### User Story 4 - Client-Side Page Iteration (Priority: P2)

As a host developer consuming paginated cost data, I need a convenient way to iterate
through all pages of cost records so that I can collect the full dataset without manually
managing tokens.

**Why this priority**: While hosts can manually loop through pages, a client-side iterator
reduces boilerplate and prevents common mistakes like forgetting to check for the last
page or mishandling empty tokens.

**Independent Test**: Can be tested by creating an iterator over a multi-page dataset and
verifying it yields all records across all pages in order, then signals completion.

**Acceptance Scenarios**:

1. **Given** a paginated cost response spanning 5 pages, **When** the host uses the
   iterator to consume all pages, **Then** it receives all records in order without gaps
   or duplicates.
2. **Given** a single-page response (total records fewer than page size), **When** the
   host uses the iterator, **Then** it yields all records and signals completion after one
   iteration.
3. **Given** a network error occurs while fetching page 3 of 5, **When** the iterator
   encounters the error, **Then** it surfaces the error to the caller with the records
   from pages 1-2 already delivered.

---

### User Story 5 - Conformance Validation for Paginated Plugins (Priority: P3)

As a plugin certification process, I need conformance tests that validate pagination
behavior so that certified plugins handle pagination correctly and consistently.

**Why this priority**: Conformance tests ensure ecosystem quality but are not required
for the feature to function. They build on the core pagination capability.

**Independent Test**: Can be tested by running the conformance suite against a mock plugin
that implements pagination and verifying all pagination-specific test cases pass.

**Acceptance Scenarios**:

1. **Given** a plugin that correctly implements pagination, **When** the conformance suite
   runs, **Then** all pagination tests pass (page boundaries, token handling, last page).
2. **Given** a plugin that returns more records than the requested page size, **When** the
   conformance suite runs, **Then** the over-sized page test fails with a descriptive error.
3. **Given** a plugin that does not implement pagination (legacy), **When** the conformance
   suite runs basic tests, **Then** the basic tests pass (pagination is optional for basic
   conformance level).

---

### Edge Cases

- What happens when a host requests page size of 0? The system uses the default page size
  (50 records).
- What happens when a host requests a page size exceeding the maximum (1,000)? The system
  caps the page size at the maximum rather than returning an error.
- What happens when a host submits a malformed or expired continuation token? The system
  returns a clear error indicating the token is invalid.
- What happens when the underlying data changes between page requests (records added or
  removed)? The pagination contract does not guarantee consistency across pages; plugins
  may return duplicate or missing records if the dataset changes. This is documented
  as expected behavior for offset-based pagination.
- What happens when a host requests a page beyond the available data (valid token but
  offset exceeds total records)? The system returns an empty result set with an empty
  continuation token.
- What happens when GetActualCost is called with dry_run=true and pagination parameters?
  The dry_run flag takes precedence; pagination parameters are ignored and the dry run
  response is returned as normal.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The GetActualCost request MUST accept a page size parameter specifying the
  maximum number of cost records to return per page.
- **FR-002**: The GetActualCost request MUST accept a continuation token parameter for
  retrieving subsequent pages of results.
- **FR-003**: The GetActualCost response MUST include a continuation token field that is
  non-empty when additional pages are available and empty when no more pages remain.
- **FR-004**: The GetActualCost response MUST include an optional total count field
  representing the total number of matching records (may be omitted or set to 0 if
  expensive to compute).
- **FR-005**: When page size is <= 0 (including unset, which defaults to 0 in proto3),
  the system MUST use a default page size of 50 records (matching the existing
  `DefaultPageSize` constant used by GetRecommendations).
- **FR-006**: When page size exceeds the maximum of 1,000, the system MUST cap the
  effective page size at 1,000 without returning an error.
- **FR-007**: A response MUST NOT contain more records than the effective page size.
- **FR-008**: When a malformed continuation token is provided, the system MUST return
  a gRPC `InvalidArgument` status error indicating the token is invalid.
- **FR-009**: When a valid token points beyond available data, the system MUST return
  an empty result set with an empty continuation token.
- **FR-010**: Plugins that do not implement pagination MUST continue to function by
  returning all records in a single response with an empty continuation token.
- **FR-011**: Hosts that do not use pagination MUST continue to receive all records by
  omitting the page size and token parameters.
- **FR-012**: The SDK MUST provide a pagination helper function for plugin developers to
  paginate in-memory slices of actual cost records. The helper MUST auto-populate the
  total count from the slice length.
- **FR-013**: The SDK MUST provide a client-side iterator for hosts to consume all pages
  of a paginated GetActualCost response, using the standard Go `Next()`/`Record()`/`Err()`
  pattern (similar to `sql.Rows`).
- **FR-014**: The pagination pattern MUST follow the same conventions established by
  GetRecommendations (offset-based tokens, same encoding, same default/max constants).
- **FR-015**: When dry_run is true on the request, pagination parameters MUST be ignored
  and the dry run response returned as normal.
- **FR-016**: Conformance tests MUST validate pagination behavior at the Standard
  conformance level (optional for Basic level).
- **FR-017**: The SDK MUST log `page_size` and `result_count` structured fields for
  each paginated GetActualCost request, matching the existing GetRecommendations
  observability pattern.

### Key Entities

- **Page Token**: An opaque string representing the position in a paginated result set.
  Encoded by the server, passed back by the client to retrieve subsequent pages. Stateless
  (all position information is embedded in the token itself). Tokens do not expire;
  validation is format-only.
- **Page Size**: An integer specifying the maximum number of records per page. Subject to
  default (50) and maximum (1,000) constraints.
- **Total Count**: An optional integer representing the total number of matching records
  across all pages. A value of 0 means "unknown or not computed" (proto3 zero-value
  semantics); when combined with a non-empty result set this indicates the plugin chose
  not to compute the total. A value of 0 with an empty result set indicates zero matching
  records.
- **Paginated Cost Iterator**: A client-side abstraction that consumes pages of actual cost
  records by automatically managing continuation tokens and yielding records to the caller.
  Uses the standard Go `Next()`/`Record()`/`Err()` pattern (similar to `sql.Rows`). Like
  `sql.Rows`, the iterator is NOT safe for concurrent use from multiple goroutines.

### Assumptions

- **Offset-based pagination**: The existing GetRecommendations pattern uses simple
  offset-based pagination with base64-encoded tokens. This feature follows the same
  approach for consistency. Cursor-based pagination is not needed since cost records
  are typically queried over fixed time ranges.
- **No cross-page consistency guarantee**: Offset-based pagination does not guarantee
  snapshot isolation. If the underlying dataset changes between page requests, callers
  may see duplicates or gaps. This is acceptable for cost data queries where the dataset
  is typically stable for a given time range.
- **Default page size of 50**: Matches the existing `DefaultPageSize` constant (50)
  used by GetRecommendations, ensuring consistency across all paginated RPCs. At
  approximately 0.5-1KB per cost record, 50 records fits comfortably within gRPC
  default limits while keeping response times low.
- **Maximum page size of 1,000**: Aligns with the existing GetRecommendations maximum.
  At 1,000 records, responses remain well under the 4MB gRPC limit.
- **No server-side aggregation or sorting**: Per the anti-guess boundary, the SDK provides
  pagination mechanics but does not perform aggregation, sorting, or deduplication.
  That logic remains with the consumer.
- **GetProjectedCost not paginated in this feature**: The issue specifically targets
  GetActualCost. GetProjectedCost pagination can be addressed separately if needed.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Hosts can retrieve 10,000+ cost records from a plugin via paginated requests
  without any single response exceeding standard message size limits.
- **SC-002**: Plugins returning large datasets (10,000 records) complete the full paginated
  retrieval within 10 seconds for a 30-day query window.
- **SC-003**: Memory usage for processing a maximum-size page (1,000 records) remains under
  100MB on both the plugin and host side.
- **SC-004**: Existing plugins that do not implement pagination continue to pass all
  existing conformance tests without modification.
- **SC-005**: Existing hosts that do not use pagination parameters continue to receive
  complete results without behavioral changes.
- **SC-006**: Plugin developers can implement paginated GetActualCost responses using
  the SDK helper with a single `PaginateActualCosts()` call, as demonstrated in the
  quickstart pattern (`specs/044-actual-cost-pagination/quickstart.md`).
- **SC-007**: The pagination protocol passes backward-compatibility validation with no
  breaking changes to the existing contract.
- **SC-008**: Conformance tests validate at least 5 distinct pagination behaviors (first
  page, middle page, last page, empty result, invalid token).

## Clarifications

### Session 2026-02-04

- Q: Default page size: 100 (spec draft) vs 50 (existing GetRecommendations)?
  → A: Use 50 to match existing `DefaultPageSize` constant for cross-RPC
  consistency.
- Q: Should GetActualCost pagination emit structured log fields?
  → A: Match GetRecommendations pattern: log `page_size` and `result_count`
  per request.
- Q: Should page tokens have expiration/TTL?
  → A: No expiration. Tokens are stateless offsets with format-only
  validation, consistent with GetRecommendations.
- Q: Should PaginateActualCosts helper auto-populate total_count?
  → A: Yes, auto-populate from in-memory slice length since len() is O(1).
- Q: What Go API pattern for the client-side page iterator?
  → A: Standard Go `Next()`/`Record()`/`Err()` pattern, matching
  `sql.Rows` style.
