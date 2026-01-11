# SDK Client Configuration Contract

**Feature**: SDK Polish v0.4.15
**Date**: 2026-01-10
**Type**: Internal API Contract (Go SDK)

## Overview

This contract documents the PulumiCost Go SDK client configuration API for timeout handling and context management.

---

## ClientConfig Configuration

### Struct Definition

```go
type ClientConfig struct {
    // BaseURL is the server's base URL (e.g., "http://localhost:8080").
    BaseURL string

    // Protocol specifies the RPC protocol to use.
    // Defaults to ProtocolConnect.
    Protocol Protocol

    // HTTPClient is the HTTP client to use for requests.
    // If nil, a default client is created using cfg.Timeout (or DefaultClientTimeout when Timeout is 0).
    // If a custom HTTPClient is provided, the caller retains ownership.
    HTTPClient *http.Client

    // Timeout is the per-client default timeout for RPC calls.
    // This field is only applied if cfg.HTTPClient is nil (i.e., when NewClient creates the HTTP client).
    // If a custom cfg.HTTPClient is provided, the caller must set HTTPClient.Timeout directly.
    // A value of 0 (default) means use the DefaultClientTimeout (30 seconds).
    // Context deadlines (if set) take precedence over this per-client timeout.
    Timeout time.Duration

    // ConnectOptions allows passing additional connect.ClientOption values.
    ConnectOptions []connect.ClientOption
}
```

### Timeout Precedence Rules

When making RPC calls, timeout is resolved in this order (highest to lowest priority):

1. **Context Deadline** (if set via `context.WithTimeout()`)
2. **Custom HTTPClient.Timeout** (if `HTTPClient` provided)
3. **ClientConfig.Timeout** (if > 0)
4. **DefaultClientTimeout** (30 seconds - fallback)

### Timeout Behavior Examples

#### Example 1: Client Timeout Only

```go
cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
cfg = cfg.WithTimeout(5 * time.Second)
client := pluginsdk.NewClient(cfg)

// RPC calls timeout after 5 seconds if no context deadline set
resp, err := client.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
if err != nil {
    // err contains context.DeadlineExceeded if timeout exceeded
}
```

#### Example 2: Context Deadline Overrides Client Timeout

```go
cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
cfg = cfg.WithTimeout(30 * time.Second)
client := pluginsdk.NewClient(cfg)

// Context deadline (1 second) takes precedence over client timeout (30 seconds)
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    // Times out after 1 second (not 30 seconds)
}
```

#### Example 3: Custom HTTPClient Timeout

```go
customClient := &http.Client{
    Timeout: 10 * time.Second,
}

cfg := pluginsdk.ClientConfig{
    BaseURL:    "http://localhost:8080",
    HTTPClient: customClient,
    Timeout:    30 * time.Second, // Ignored (HTTPClient takes precedence)
}

client := pluginsdk.NewClient(cfg)
// RPC calls timeout after 10 seconds (custom client timeout)
```

#### Example 4: Default Timeout (No Configuration)

```go
cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
client := pluginsdk.NewClient(cfg)

// Uses DefaultClientTimeout (30 seconds)
resp, err := client.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
```

---

## WithTimeout() Fluent API

### Method Signature

```go
// WithTimeout returns a copy of the configuration with the specified timeout.
// This allows for fluent configuration chaining.
func (c ClientConfig) WithTimeout(timeout time.Duration) ClientConfig
```

### Usage Pattern

```go
// Fluent configuration chaining
cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
    .WithTimeout(5 * time.Minute)
    .WithProtocol(pluginsdk.ProtocolGRPC)

client := pluginsdk.NewClient(cfg)
```

### Implementation Details

- Returns a **copy** of the configuration (does not modify original)
- Clears default `HTTPClient` to force `NewClient` to rebuild with new timeout
- Zero timeout is valid (uses `DefaultClientTimeout`)

---

## Error Handling

### wrapRPCError Function

```go
// wrapRPCError wraps an RPC error with context about the operation.
// It distinguishes context cancellation/timeout from other errors.
func wrapRPCError(ctx context.Context, operation string, err error) error
```

### Error Messages

| Error Condition           | Error Message                              | Error Type                                   |
| ------------------------- | ------------------------------------------ | -------------------------------------------- |
| Context deadline exceeded | `"<operation> RPC cancelled or timed out"` | `errors.Join(context.DeadlineExceeded, err)` |
| Context cancelled         | `"<operation> RPC cancelled or timed out"` | `errors.Join(context.Canceled, err)`         |
| Other RPC failure         | `"<operation> RPC failed: <error>"`        | `fmt.Errorf` wrapping original error         |

### Error Example

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        // err contains: "GetPluginInfo RPC cancelled or timed out"
        // with wrapped context.DeadlineExceeded
    }
}
```

---

## Constants

```go
// DefaultClientTimeout is the default HTTP client timeout for plugin requests.
const DefaultClientTimeout = 30 * time.Second
```

---

## Testing Requirements

### Unit Tests

- ✅ `TestClientConfig_WithTimeout` - Verifies WithTimeout() sets timeout correctly
- ❌ **Missing**: Timeout behavior integration test with slow server
- ❌ **Missing**: Context deadline precedence test
- ❌ **Missing**: Custom HTTPClient timeout precedence test

### Integration Tests

- ❌ **Missing**: Slow mock server test (client timeout verification)
- ❌ **Missing**: Context timeout vs client timeout test

---

## Backward Compatibility

### Breaking Changes

**None** - All changes maintain backward compatibility:

1. Existing code without timeout configuration continues to work
2. Default behavior (30-second timeout) unchanged
3. Custom HTTPClient usage unchanged

### Migration Guide

No migration required. New timeout features are opt-in via `WithTimeout()`.

---

## Security Considerations

### Timeout as DoS Protection

- Per-client timeout prevents indefinite blocking on slow/rogue servers
- Context deadlines allow per-request timeout control
- Custom HTTPClient callers responsible for their own timeout configuration

### Error Message Security

- `wrapRPCError` does not expose server stack traces or internal details
- Error messages are user-friendly and do not reveal implementation details
- Detailed errors logged server-side (not sent to clients)

---

## Performance Implications

### Timeout Overhead

- Timeout checking is O(1) - no performance impact
- Context deadline checking is built into Go's context package (negligible overhead)

### Connection Pooling

- HTTP clients with timeouts still benefit from connection pooling
- `CloseIdleConnections()` called on client close
- `HighThroughputClientConfig` provides connection pool optimization

---

## Future Enhancements

### Potential Improvements

1. **Per-Method Timeouts**: Allow different timeouts for different RPC methods
2. **Retry Configuration**: Add automatic retry with exponential backoff
3. **Circuit Breaker**: Add circuit breaker pattern for cascading failure prevention
4. **Metrics**: Add timeout metrics and monitoring

### Extension Points

- `ConnectOptions` allows passing custom interceptors
- Custom `HTTPClient` allows full control over HTTP behavior
- Context deadlines provide per-request control without SDK changes
