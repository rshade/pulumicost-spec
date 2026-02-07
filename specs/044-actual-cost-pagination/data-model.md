# Data Model: GetActualCost Pagination Support

**Feature**: 044-actual-cost-pagination
**Date**: 2026-02-04

## Entities

### GetActualCostRequest (modified)

Extends the existing `GetActualCostRequest` protobuf message with pagination fields.

| Field | Type | Number | Description | Validation |
|-------|------|--------|-------------|------------|
| `resource_id` | string | 1 | Flexible resource ID per plugin | Existing |
| `start` | google.protobuf.Timestamp | 2 | Query period start | Existing |
| `end` | google.protobuf.Timestamp | 3 | Query period end | Existing |
| `tags` | map\<string, string\> | 4 | Optional filter tags | Existing |
| `arn` | string | 5 | Canonical cloud identifier | Existing |
| `dry_run` | bool | 6 | Dry-run mode flag | Existing |
| **`page_size`** | **int32** | **7** | **Max records per page** | **Default: 50, Max: 1000, clamped** |
| **`page_token`** | **string** | **8** | **Continuation token** | **Base64-encoded offset or empty** |

**Validation Rules**:

- `page_size <= 0` -> use `DefaultPageSize` (50)
- `page_size > MaxPageSize` (1000) -> clamp to `MaxPageSize`
- `page_token` empty -> start from offset 0
- `page_token` malformed -> return `InvalidArgument` gRPC error
- `page_token` valid but offset > total -> return empty results, empty next token
- `dry_run = true` -> ignore `page_size` and `page_token`

### GetActualCostResponse (modified)

Extends the existing `GetActualCostResponse` protobuf message with pagination metadata.

| Field | Type | Number | Description | Validation |
|-------|------|--------|-------------|------------|
| `results` | repeated ActualCostResult | 1 | Cost data points | Existing |
| `fallback_hint` | FallbackHint | 2 | Fallback signal | Existing |
| `dry_run_result` | DryRunResponse | 3 | Dry-run metadata | Existing |
| **`next_page_token`** | **string** | **4** | **Token for next page** | **Non-empty if more pages; empty if last** |
| **`total_count`** | **int32** | **5** | **Total matching records** | **0 if unknown/expensive to compute** |

**Validation Rules**:

- `len(results) <= effective_page_size` (response MUST NOT exceed requested page size)
- `next_page_token` empty when on last page or single-page result
- `next_page_token` non-empty when more records exist
- `total_count >= 0` (proto3 zero-value means "unknown")
- When `dry_run_result` is populated, `results` is empty, `next_page_token` is empty,
  `total_count` is 0

### Page Token (internal encoding)

Opaque string representing position in a paginated result set.

| Property | Value |
|----------|-------|
| Encoding | Base64 of integer offset (e.g., `base64("100")` -> `"MTAw"`) |
| Stateless | All position info embedded in token itself |
| Expiration | None (format-only validation) |
| Direction | Forward only (no backward pagination) |

**Encoding/Decoding**:

```text
Encode: offset (int) -> strconv.Itoa -> base64.StdEncoding.EncodeToString -> token (string)
Decode: token (string) -> base64.StdEncoding.DecodeString -> strconv.Atoi -> offset (int)
```

**Reuses existing functions**: `EncodePageToken()` and `DecodePageToken()` from
`sdk/go/pluginsdk/helpers.go:978-997`.

### ActualCostIterator (new - client-side)

Client-side abstraction for consuming paginated actual cost responses.
NOT safe for concurrent use from multiple goroutines (same contract as `sql.Rows`).

| Field | Type | Description |
|-------|------|-------------|
| `ctx` | context.Context | Request context for cancellation |
| `fetchFn` | func(ctx, pageToken, pageSize) -> (response, error) | RPC callback |
| `pageSize` | int32 | Requested page size |
| `pageToken` | string | Current continuation token |
| `current` | []*pbc.ActualCostResult | Current page records |
| `index` | int | Position within current page |
| `totalCount` | int32 | Total count from most recent response |
| `done` | bool | True when all pages consumed |
| `err` | error | First error encountered |

**State Transitions**:

```text
Created -> Fetching -> HasRecords -> Fetching -> ... -> Done
                \-> Error (on RPC failure)
```

**API**:

| Method | Returns | Description |
|--------|---------|-------------|
| `Next()` | bool | Advances to next record; returns false when done or error |
| `Record()` | *ActualCostResult | Returns current record (valid after Next() returns true) |
| `Err()` | error | Returns first error encountered (check after Next() returns false) |
| `TotalCount()` | int32 | Returns total_count from most recent response (0 if unknown) |

## Relationships

```text
GetActualCostRequest ---[contains]--> page_size, page_token
        |
        v (RPC call)
GetActualCostResponse ---[contains]--> next_page_token, total_count
        |                                    |
        v                                    v
  ActualCostResult[]                  (fed back as page_token
        |                              in next request)
        v
ActualCostIterator ---[consumes]--> GetActualCostResponse (page by page)
        |
        v (yields one at a time)
  ActualCostResult
```

## SDK Helper: PaginateActualCosts

Server-side helper for plugin developers to paginate in-memory slices.

**Signature**:

```go
func PaginateActualCosts(
    results []*pbc.ActualCostResult,
    pageSize int32,
    pageToken string,
) ([]*pbc.ActualCostResult, string, int32, error)
//  ^page results        ^next token ^total ^error
```

**Behavior**:

1. Clamp `pageSize` to `[DefaultPageSize, MaxPageSize]`
2. Decode `pageToken` to offset (0 if empty)
3. Return error if token malformed
4. Extract `results[offset:offset+pageSize]`
5. Generate `nextToken` if more records exist
6. Return `int32(len(results))` as `totalCount`

**Differences from `PaginateRecommendations`**:

- Returns 4 values instead of 3 (adds `totalCount`)
- Operates on `[]*pbc.ActualCostResult` instead of `[]*pbc.Recommendation`
- Otherwise identical algorithm

## Constants

| Constant | Value | Notes |
|----------|-------|-------|
| `DefaultPageSize` | 50 | Already exists in `pluginsdk/helpers.go:923` |
| `MaxPageSize` | 1000 | Already exists in `pluginsdk/helpers.go:926` |

No new constants needed. Reuse existing constants for cross-RPC consistency.

## Logging Fields

| Field Constant | Value | Usage |
|----------------|-------|-------|
| `FieldPageSize` | `"page_size"` | Already exists in `pluginsdk/logging.go:67` |
| `FieldResultCount` | `"result_count"` | New: count of records in current page response |
