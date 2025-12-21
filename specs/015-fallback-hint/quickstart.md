# Quickstart: Using Fallback Hints

## Overview

The `FallbackHint` enum allows your plugin to tell the Core system whether it should look elsewhere
for cost data. This is useful when your plugin is specialized or when it knows it doesn't have data
for a specific resource.

## Updating Your Plugin

### 1. Regenerate SDK

Ensure your plugin is using the latest `pulumicost` SDK (v1.x.x).

```bash
go get github.com/rshade/pulumicost-spec/sdk/go@latest
```

### 2. Using Functional Options

When constructing a `GetActualCostResponse`, use the new `WithFallbackHint` option.

#### Scenario A: Data Found (Default)

If your plugin found data, you don't need to do anything. The default is `FALLBACK_HINT_UNSPECIFIED` which acts as "No Fallback".

```go
// Default behavior: Core uses your results and stops.
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithResults(results),
)
```

#### Scenario B: No Data, Try Others

If your plugin checked its source but found nothing (e.g., resource too new), tell Core to try other plugins.

```go
// "I found nothing, but maybe someone else knows."
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithResults(nil), // or empty slice
    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
)
```

#### Scenario C: Not My Job

If your plugin receives a request for a resource type it doesn't handle, strictly require fallback.

```go
// "I definitely don't handle this type."
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_REQUIRED),
)
```

#### Scenario D: It's Free (Zero Cost)

If you know the resource exists but costs $0.00 (free tier), return a result with 0 cost and
explicit `NONE` (optional, but clear).

```go
// "I found it, and it's free. Don't ask anyone else."
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithResults([]*pbc.ActualCostResult{
        {Cost: 0.0, Timestamp: ...},
    }),
    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
)
```

## FallbackHint Decision Matrix

Use this matrix to determine the correct hint for your scenario:

| Scenario | Has Results? | Hint Value | Core Behavior |
|----------|--------------|------------|---------------|
| Found actual cost data | Yes | `UNSPECIFIED` or `NONE` | Use data, stop |
| Found zero-cost data (free tier) | Yes (cost: 0) | `NONE` | Use data, stop |
| No billing data exists yet | No | `RECOMMENDED` | Try other plugins |
| Resource type not supported | No | `REQUIRED` | Must try others |
| API error / network failure | N/A | Return gRPC error | Error handling path |

## Zero-Cost vs No-Data Distinction

It's critical to distinguish between "resource costs nothing" and "no billing data found":

### Zero-Cost (Free Tier)

The resource exists and genuinely costs $0.00 (e.g., AWS free tier, always-free resources).

```go
// Resource is free - return result with 0 cost and NONE hint
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithResults([]*pbc.ActualCostResult{
        {Cost: 0.0, Source: "aws-ce"},
    }),
    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
)
```

### No Data Found

No billing records exist (e.g., resource too new, not yet in billing system).

```go
// No billing data - return empty results with RECOMMENDED hint
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
)
```

### SaaS Vendor Translation

Some SaaS vendors return `0.00` for "not found". Your plugin must translate:

```go
// If vendor returns 0.00 but resource doesn't exist in their system
if vendorResponse.Cost == 0.0 && !vendorResponse.ResourceExists {
    // Translate to proper "no data" semantics
    return pluginsdk.NewActualCostResponse(
        pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
    ), nil
}
```

## Edge Cases

### Data + Hint Conflict

If a plugin returns results AND sets `FALLBACK_RECOMMENDED` or `FALLBACK_REQUIRED`,
the Core system prioritizes data:

- **Behavior**: Core uses the returned results and logs a warning
- **Rationale**: Data presence indicates the plugin did its job

### Unrecognized Hint Values

If Core receives an unrecognized hint value (future extensibility):

- **Behavior**: Treat as `FALLBACK_HINT_UNSPECIFIED` (no fallback)
- **Rationale**: Safe default for forward compatibility

### Fallback Chain Termination

If all plugins in a fallback chain return `FALLBACK_RECOMMENDED`:

- **Behavior**: Core stops and returns empty result
- **Rationale**: Prevents infinite delegation loops

## Best Practices

- **Do not** set `RECOMMENDED` if you are returning valid cost data (unless you are partial source, which is rare).
- **Do not** use `REQUIRED` for API errors. Return a gRPC error instead.
- **Always** use the functional options pattern provided by the SDK.
