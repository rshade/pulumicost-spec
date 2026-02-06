# Proto Contract Changes: GetActualCost Pagination

**File**: `proto/finfocus/v1/costsource.proto`

## GetActualCostRequest Changes

```protobuf
// GetActualCostRequest contains parameters for retrieving historical cost data.
message GetActualCostRequest {
  // resource_id is a flexible ID per plugin (e.g., "i-abc123", "namespace/default")
  string resource_id = 1;
  // start timestamp for the cost query period
  google.protobuf.Timestamp start = 2;
  // end timestamp for the cost query period
  google.protobuf.Timestamp end = 3;
  // tags provide optional extra filters for cost retrieval
  map<string, string> tags = 4;
  // New field: Canonical Cloud Identifier (e.g. AWS ARN, Azure Resource ID, GCP Full Resource Name)
  string arn = 5;
  // dry_run when true, returns DryRunResponse in dry_run_result field
  // instead of performing actual cost data retrieval.
  // Default: false (normal cost retrieval behavior).
  // When true, the response will contain dry_run_result instead of results.
  bool dry_run = 6;

+ // page_size is the maximum number of cost records to return per page.
+ // Default: 50 (matches DefaultPageSize). Maximum: 1000 (matches MaxPageSize).
+ // Values <= 0 use the default. Values > 1000 are clamped to 1000.
+ // Ignored when dry_run is true.
+ int32 page_size = 7;
+ // page_token is the continuation token from a previous GetActualCost response.
+ // Empty string requests the first page of results.
+ // Ignored when dry_run is true.
+ string page_token = 8;
}
```

## GetActualCostResponse Changes

```protobuf
// GetActualCostResponse contains the list of actual cost results.
message GetActualCostResponse {
  // results contains the actual cost data points for the requested period
  repeated ActualCostResult results = 1;
  // fallback_hint indicates whether the core should attempt to query other plugins
  FallbackHint fallback_hint = 2;
  // dry_run_result contains field mapping information when request.dry_run
  // was true. Empty/nil when dry_run was false or not set.
  // When populated, results field will be empty.
  DryRunResponse dry_run_result = 3;

+ // next_page_token is the token for retrieving the next page of results.
+ // Non-empty when additional pages are available. Empty when this is the
+ // last page or when all results fit in a single response.
+ string next_page_token = 4;
+ // total_count is the total number of matching cost records across all pages.
+ // Optional: may be 0 if the total is expensive to compute.
+ // When populated by the SDK PaginateActualCosts helper, this is automatically
+ // set to the slice length.
+ int32 total_count = 5;
}
```

## Wire Compatibility

| Aspect | Status | Notes |
|--------|--------|-------|
| Field number reuse | Safe | Fields 7-8 (request), 4-5 (response) are unused |
| Default values | Safe | `page_size=0` and `page_token=""` are proto3 defaults |
| Old client -> new server | Safe | Server receives defaults, returns all records |
| New client -> old server | Safe | Server ignores unknown fields, returns all records |
| buf breaking check | Passes | Additive changes only (new fields) |

## Backward Compatibility Matrix

| Client Version | Server Version | Behavior |
|----------------|----------------|----------|
| Old (no pagination) | Old (no pagination) | Unchanged: all records in one response |
| Old (no pagination) | New (pagination-aware) | Server sees `page_size=0`, returns all records |
| New (pagination-aware) | Old (no pagination) | Server ignores new fields, returns all records |
| New (pagination-aware) | New (pagination-aware) | Full pagination support |
