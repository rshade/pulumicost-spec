# Quickstart: GetActualCost Pagination

**Feature**: 044-actual-cost-pagination

## Plugin Developer: Paginated Responses

Use the `PaginateActualCosts` SDK helper to paginate in-memory cost record slices.

### Basic Usage

```go
import (
    pluginsdk "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func (p *MyPlugin) GetActualCost(
    ctx context.Context,
    req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
    // 1. Handle dry_run (pagination ignored)
    if req.DryRun {
        return p.handleDryRun(ctx, req)
    }

    // 2. Fetch all cost records from your data source
    allResults, err := p.fetchCostData(ctx, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "fetch failed: %v", err)
    }

    // 3. Apply pagination (handles page_size defaults, token decoding, slicing)
    page, nextToken, totalCount, err := pluginsdk.PaginateActualCosts(
        allResults,
        req.PageSize,
        req.PageToken,
    )
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "%v", err)
    }

    // 4. Return paginated response
    return &pbc.GetActualCostResponse{
        Results:       page,
        FallbackHint:  pbc.FallbackHint_FALLBACK_HINT_NONE,
        NextPageToken: nextToken,
        TotalCount:    totalCount,
    }, nil
}
```

### No Changes Needed for Legacy Plugins

Existing plugins that return all results in a single response continue to work.
When `page_size` and `page_token` are not set (proto3 defaults of 0 and ""),
hosts receive all records with an empty `next_page_token`.

## Host Developer: Consuming Paginated Results

### Using the ActualCostIterator

```go
import (
    pluginsdk "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// Create an iterator that automatically handles page tokens
iter := pluginsdk.NewActualCostIterator(ctx,
    func(ctx context.Context, pageToken string, pageSize int32) (*pbc.GetActualCostResponse, error) {
        return client.GetActualCost(ctx, &pbc.GetActualCostRequest{
            ResourceId: "i-abc123",
            Start:      startTimestamp,
            End:        endTimestamp,
            PageSize:   pageSize,
            PageToken:  pageToken,
        })
    },
    100, // page size
)

// Iterate through all records across all pages
var allCosts []*pbc.ActualCostResult
for iter.Next() {
    allCosts = append(allCosts, iter.Record())
}
if err := iter.Err(); err != nil {
    return fmt.Errorf("iteration failed: %w", err)
}

fmt.Printf("Retrieved %d records (total: %d)\n", len(allCosts), iter.TotalCount())
```

### Manual Page Iteration

```go
var allResults []*pbc.ActualCostResult
pageToken := ""

for {
    resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
        ResourceId: "i-abc123",
        Start:      startTimestamp,
        End:        endTimestamp,
        PageSize:   100,
        PageToken:  pageToken,
    })
    if err != nil {
        return fmt.Errorf("page request failed: %w", err)
    }

    allResults = append(allResults, resp.Results...)

    if resp.NextPageToken == "" {
        break // Last page
    }
    pageToken = resp.NextPageToken
}
```

## TypeScript: Client-Side Iterator

```typescript
import { actualCostIterator } from "@finfocus/client/utils/pagination";

const request = {
  resourceId: "i-abc123",
  start: startTimestamp,
  end: endTimestamp,
  pageSize: 100,
};

// Async generator automatically handles page tokens
for await (const record of actualCostIterator(client, request)) {
  console.log(`Cost: ${record.cost} ${record.currency}`);
}
```

## Key Constants

| Constant | Value | Description |
|----------|-------|-------------|
| DefaultPageSize | 50 | Used when `page_size` is 0 or unset |
| MaxPageSize | 1000 | Maximum allowed; larger values are clamped |

## Edge Cases

| Scenario | Behavior |
|----------|----------|
| `page_size = 0` | Uses DefaultPageSize (50) |
| `page_size > 1000` | Clamped to MaxPageSize (1000) |
| Malformed `page_token` | Returns `InvalidArgument` gRPC error |
| Token beyond data | Empty results, empty next token |
| `dry_run = true` with pagination | Pagination ignored, dry run response returned |
| Dataset changes between pages | Possible duplicates/gaps (documented behavior) |
