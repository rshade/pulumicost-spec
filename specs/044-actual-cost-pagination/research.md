# Research: GetActualCost Pagination Support

**Feature**: 044-actual-cost-pagination
**Date**: 2026-02-04
**Status**: Complete

## Research Tasks

### R1: Existing Pagination Pattern (GetRecommendations)

**Task**: Understand the established pagination pattern to ensure GetActualCost pagination
is consistent.

**Decision**: Replicate the GetRecommendations offset-based pagination pattern exactly.

**Rationale**: The codebase already has a well-tested pagination implementation in
GetRecommendations. Reusing the same pattern ensures consistency, reduces cognitive load
for plugin developers, and allows sharing of encoding/decoding utilities.

**Findings**:

- **Proto fields**: `page_size` (int32), `page_token` (string) on request;
  `next_page_token` (string) on response
- **Constants**: `DefaultPageSize = 50`, `MaxPageSize = 1000` (in `pluginsdk/helpers.go:922-926`)
- **Token encoding**: Base64-encoded integer offset (`EncodePageToken`/`DecodePageToken`
  in `pluginsdk/helpers.go:978-997`)
- **Helper function**: `PaginateRecommendations()` handles page size clamping, token
  decoding, slice extraction, and next token generation (`pluginsdk/helpers.go:930-976`)
- **Response builder**: Functional options pattern via `NewActualCostResponse()` with
  `WithResults()`, `WithFallbackHint()` options
- **Test coverage**: 7 dedicated pagination tests in `pluginsdk/helpers_test.go:615-737`
- **Mock implementation**: `paginateMockRecommendations()` in `testing/mock_plugin.go:1502-1543`

**Alternatives Considered**:

- **Cursor-based pagination**: More complex, unnecessary for fixed time-range queries.
  Offset-based is simpler and matches the existing pattern.
- **Keyset pagination**: Requires a stable sort key. Cost records don't have a natural
  unique sort key across providers.

### R2: Proto Field Number Assignment

**Task**: Determine correct field numbers for new pagination fields in
GetActualCostRequest and GetActualCostResponse.

**Decision**: Use field numbers 7-8 for request, field number 4-5 for response.

**Rationale**: Proto3 field numbers must be unique within a message and cannot reuse
previously assigned numbers. Current GetActualCostRequest uses fields 1-6, and
GetActualCostResponse uses fields 1-3.

**Findings**:

- **GetActualCostRequest** current fields: `resource_id` (1), `start` (2), `end` (3),
  `tags` (4), `arn` (5), `dry_run` (6)
- **GetActualCostResponse** current fields: `results` (1), `fallback_hint` (2),
  `dry_run_result` (3)
- New request fields: `page_size` = 7, `page_token` = 8
- New response fields: `next_page_token` = 4, `total_count` = 5

**Alternatives Considered**:

- Using a nested `PaginationParams` message: Adds unnecessary nesting. Flat fields match
  GetRecommendations pattern and are simpler.

### R3: Total Count Field Design

**Task**: Determine how `total_count` should be represented and when it's populated.

**Decision**: Use `int32 total_count` (optional, 0 means unknown/not computed).

**Rationale**: Some plugins can efficiently compute total count (in-memory datasets),
while others cannot (streaming from external APIs). Making it optional with 0 as
"unknown" follows proto3 zero-value semantics naturally.

**Findings**:

- GetRecommendations does NOT have a `total_count` field (only `RecommendationSummary`)
- The spec requires `total_count` as an optional field (FR-004)
- The SDK `PaginateActualCosts` helper will auto-populate from slice length (O(1))
- External API-backed plugins can omit it (return 0)

**Alternatives Considered**:

- `google.protobuf.Int32Value` wrapper: Distinguishes "0 records" from "unknown". However,
  cost queries with 0 records return empty results anyway, so the distinction is unnecessary.
- `optional int32`: Proto3 optional adds `has_*` methods but complicates SDK code for
  minimal benefit.

### R4: Client-Side Iterator Design

**Task**: Determine the Go API pattern for consuming paginated actual cost responses.

**Decision**: Use the standard Go `Next()`/`Record()`/`Err()` pattern (similar to
`sql.Rows`), as specified in the feature spec.

**Rationale**: This is the idiomatic Go pattern for iterating over external data sources.
It's familiar to Go developers and handles error propagation naturally.

**Findings**:

- TypeScript SDK already has an async generator iterator for recommendations
  (`sdk/typescript/packages/client/src/utils/pagination.ts:6-25`)
- No Go client-side iterator exists yet (Go SDK is plugin-side focused)
- The iterator needs a callback function type for making RPC calls

**API Shape**:

```go
type ActualCostIterator struct { ... }

func NewActualCostIterator(
    ctx context.Context,
    fetchFn func(ctx context.Context, pageToken string, pageSize int32) (*pbc.GetActualCostResponse, error),
    pageSize int32,
) *ActualCostIterator

func (it *ActualCostIterator) Next() bool
func (it *ActualCostIterator) Record() *pbc.ActualCostResult
func (it *ActualCostIterator) Err() error
func (it *ActualCostIterator) TotalCount() int32  // from most recent response
```

**Alternatives Considered**:

- Channel-based iterator: More complex, harder to handle errors, not idiomatic for
  sequential consumption.
- Callback-based (`ForEach`): Less flexible, doesn't allow early termination without
  context cancellation.
- Generic iterator (`iter.Seq` from Go 1.23+): The project uses Go 1.25.6 so this is
  available, but `Next()/Record()/Err()` was explicitly specified in the feature spec.

### R5: TypeScript SDK Synchronization

**Task**: Determine what TypeScript SDK changes are required for this feature.

**Decision**: Add `actualCostIterator()` async generator in
`sdk/typescript/packages/client/src/utils/pagination.ts`, mirroring the existing
`recommendationsIterator()`.

**Rationale**: Constitution XIII (Multi-Language SDK Synchronization) requires all SDKs
to be updated when proto changes occur. The TypeScript SDK already has the pagination
pattern for recommendations.

**Findings**:

- Existing pattern in `pagination.ts` uses async generator with `do/while` loop
- Proto changes will auto-generate TypeScript bindings via buf
- Client wrapper will need the new iterator function
- No TypeScript-side pagination helper needed (token encoding is server-side)

**Alternatives Considered**:

- Deferring TypeScript to a separate PR: Violates Constitution XIII. Must be in same PR.

### R6: Conformance Test Integration

**Task**: Determine where pagination conformance tests belong in the test hierarchy.

**Decision**: Pagination conformance tests at Standard level (optional for Basic).
New file: `sdk/go/testing/pagination_conformance_test.go`.

**Rationale**: Matches FR-016 requirement. Basic plugins should work without pagination
(backward compatibility). Standard plugins should implement pagination correctly.

**Findings**:

- Existing conformance levels: Basic (required), Standard (recommended), Advanced
  (high-performance)
- Mock plugin already implements recommendation pagination in `mock_plugin.go:1495-1564`
- Need to add actual cost pagination to mock plugin
- Test cases: first page, middle page, last page, empty result, invalid token,
  over-sized page, default page size, backward compatibility

### R7: Logging and Observability

**Task**: Determine structured logging fields for paginated GetActualCost requests.

**Decision**: Log `page_size` and `result_count` fields, matching the GetRecommendations
pattern per FR-017.

**Rationale**: Consistency with existing observability patterns. The `FieldPageSize`
constant already exists in `logging.go:67`.

**Findings**:

- Existing field: `FieldPageSize = "page_size"` (reuse)
- Need new field: `FieldResultCount = "result_count"` for actual cost record count
- Existing `FieldRecommendationCount` is recommendation-specific; actual costs need
  a generic result count field

### R8: DryRun Interaction

**Task**: Clarify how pagination interacts with dry_run mode.

**Decision**: When `dry_run = true`, pagination parameters are ignored and the dry run
response is returned as normal (FR-015).

**Rationale**: DryRun is for capability introspection, not data retrieval. Pagination
is meaningless without actual data. This matches the spec's edge case documentation.

**Findings**:

- Current proto: `dry_run` field (6) on request, `dry_run_result` field (3) on response
- When dry_run is true, `results` is empty and `dry_run_result` is populated
- No pagination fields should be populated in dry_run responses
- The SDK helper and conformance tests should validate this behavior
