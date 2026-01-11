# Quickstart Guide: SDK Polish v0.4.15

**Feature**: SDK Polish v0.4.15
**Date**: 2026-01-10
**Audience**: Plugin Developers

## Overview

SDK Polish v0.4.15 adds three enhancements to the PulumiCost Go SDK:

1. **Configurable Client Timeouts** - Per-client timeout configuration with context deadline support
2. **User-Friendly GetPluginInfo Errors** - Improved error messages for better developer experience
3. **GetPluginInfo Performance Conformance** - Performance validation for GetPluginInfo RPC

This guide shows how to use these features and verify they work correctly.

---

## Feature 1: Configurable Client Timeouts

### Use Case

You need to configure timeout behavior for a client that may call long-running RPC operations
(e.g., `GetActualCost` with large date ranges).

### Basic Timeout Configuration

```go
package main

import (
    "context"
    "log"
    "time"

    pluginsdk "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    // Create client with 5-second timeout
    cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
    cfg = cfg.WithTimeout(5 * time.Second)

    client := pluginsdk.NewClient(cfg)
    defer client.Close()

    // Make RPC call (times out after 5 seconds if not completed)
    resp, err := client.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
    if err != nil {
        log.Printf("RPC failed: %v", err)
        return
    }

    log.Printf("Plugin name: %s, version: %s", resp.GetName(), resp.GetVersion())
}
```

### Context Deadline Override

Context deadlines take precedence over client-level timeouts:

```go
cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
cfg = cfg.WithTimeout(30 * time.Second)
client := pluginsdk.NewClient(cfg)

// Context deadline (1 second) takes precedence over client timeout (30 seconds)
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    // Error: "GetPluginInfo RPC cancelled or timed out"
    // Times out after 1 second (not 30 seconds)
    log.Printf("RPC timed out: %v", err)
    return
}
```

### Default Timeout Behavior

If no timeout is configured, the default 30-second timeout is used:

```go
cfg := pluginsdk.DefaultClientConfig("http://localhost:8080")
// No WithTimeout() call
client := pluginsdk.NewClient(cfg)

// Uses DefaultClientTimeout (30 seconds)
resp, err := client.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
```

### Custom HTTPClient Timeouts

If you provide a custom `HTTPClient`, its timeout takes precedence:

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

### Timeout Precedence Rules

When making RPC calls, timeout is resolved in this order (highest to lowest priority):

1. **Context Deadline** (if set via `context.WithTimeout()`)
2. **Custom HTTPClient.Timeout** (if `HTTPClient` provided)
3. **ClientConfig.Timeout** (if > 0)
4. **DefaultClientTimeout** (30 seconds - fallback)

### Error Messages

Timeout errors are wrapped with user-friendly messages:

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        // Error message: "GetPluginInfo RPC cancelled or timed out"
        // Wrapped error contains context.DeadlineExceeded
        log.Printf("RPC timed out: %v", err)
    }
}
```

---

## Feature 2: User-Friendly GetPluginInfo Errors

### Use Case

Your plugin returns invalid metadata (nil, incomplete, or invalid spec_version), and you want to
provide actionable error messages to developers.

### Implementing GetPluginInfo

Add the `PluginInfoProvider` interface to your plugin:

```go
package main

import (
    "context"
    "log"

    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
    pluginsdk "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

// MyPlugin implements Plugin interface and optional PluginInfoProvider
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-cost-plugin"
}

func (p *MyPlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    // Implementation...
    return &pbc.GetProjectedCostResponse{}, nil
}

func (p *MyPlugin) GetActualCost(
    ctx context.Context,
    req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
    // Implementation...
    return &pbc.GetActualCostResponse{}, nil
}

func (p *MyPlugin) GetPricingSpec(
    ctx context.Context,
    req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
    // Implementation...
    return &pbc.GetPricingSpecResponse{}, nil
}

func (p *MyPlugin) EstimateCost(
    ctx context.Context,
    req *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
    // Implementation...
    return &pbc.EstimateCostResponse{}, nil
}

// PluginInfoProvider implementation
func (p *MyPlugin) GetPluginInfo(
    ctx context.Context,
    req *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error) {
    // Return plugin metadata
    return &pbc.GetPluginInfoResponse{
        Name:        "my-cost-plugin",
        Version:     "1.0.0",
        SpecVersion: "v1.2.0",
        Providers:   []string{"aws", "azure", "gcp"},
        Metadata: map[string]string{
            "description":  "My custom cost plugin",
            "author":       "Your Name",
        },
    }, nil
}

func main() {
    plugin := &MyPlugin{}

    // Server automatically validates GetPluginInfo responses
    // and returns user-friendly error messages
    err := pluginsdk.Serve(context.Background(), pluginsdk.ServeConfig{
        Plugin: plugin,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

### Error Message Examples

The SDK automatically validates your `GetPluginInfo` response and returns user-friendly errors:

| Error Condition       | Client Receives                                    | Server Logs                                                           |
| --------------------- | -------------------------------------------------- | --------------------------------------------------------------------- |
| Plugin returns `nil`  | "unable to retrieve plugin metadata"               | "GetPluginInfo returned nil response"                                 |
| Required fields empty | "plugin metadata is incomplete"                    | "GetPluginInfo returned incomplete response" (with field values)      |
| Invalid spec_version  | "plugin reported an invalid specification version" | "GetPluginInfo returned invalid spec_version" (with validation error) |

### Handling Legacy Plugins

Clients should gracefully handle plugins that don't implement `GetPluginInfo`:

```go
resp, err := client.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
if err != nil {
    if status.Code(err) == codes.Unimplemented {
        // Legacy plugin - use fallback values
        log.Info("Plugin does not implement GetPluginInfo")
        return &PluginMetadata{
            Name:    "unknown",
            Version: "unknown",
        }
    }
    return nil, fmt.Errorf("GetPluginInfo failed: %w", err)
}

return &PluginMetadata{
    Name:        resp.GetName(),
    Version:     resp.GetVersion(),
    SpecVersion: resp.GetSpecVersion(),
}, nil
```

---

## Feature 3: GetPluginInfo Performance Conformance

### Use Case

You want to verify your `GetPluginInfo` implementation meets performance requirements (< 100ms per call).

### Running Performance Conformance Tests

Create a conformance test file for your plugin:

```go
package main

import (
    "os"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
    "github.com/rshade/pulumicost-spec/sdk/go/testing"
    pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-cost-plugin"
}

func (p *MyPlugin) GetProjectedCost(
    ctx context.Context,
    req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
    // Implementation...
    return &pbc.GetProjectedCostResponse{}, nil
}

func (p *MyPlugin) GetActualCost(
    ctx context.Context,
    req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
    // Implementation...
    return &pbc.GetActualCostResponse{}, nil
}

func (p *MyPlugin) GetPricingSpec(
    ctx context.Context,
    req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
    // Implementation...
    return &pbc.GetPricingSpecResponse{}, nil
}

func (p *MyPlugin) EstimateCost(
    ctx context.Context,
    req *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
    // Implementation...
    return &pbc.EstimateCostResponse{}, nil
}

func main() {
    plugin := &MyPlugin{}

    // Create test harness
    harness := testing.NewTestHarness(plugin)
    defer harness.Stop()

    // Get client connection
    conn, err := harness.CreateClientConnection()
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // Create client
    client := pbc.NewCostSourceServiceClient(conn)

    // Run GetPluginInfo performance test
    t := &testing.T{}
    harness.SetClient(client)

    testFunc := testing.createGetPluginInfoLatencyTest()
    result := testFunc(&harness)

    if !result.Success {
        t.Errorf("GetPluginInfo performance test failed: %v", result.Error)
    }

    t.Logf("GetPluginInfo performance: %s", result.Details)

    if result.Error != nil {
        os.Exit(1)
    }
}
```

### Performance Thresholds

| Conformance Level | Threshold | Description                                      |
| ----------------- | --------- | ------------------------------------------------ |
| Standard          | 100ms     | Minimum requirement for production-ready plugins |
| Advanced          | 50ms      | High-performance requirement                     |

### Interpreting Test Results

**Pass Example**:

```text
Performance_GetPluginInfoLatency: PASSED
Details: Avg: 78ms (threshold: 100ms)
Min: 45ms, Max: 95ms
```

**Fail Example**:

```text
Performance_GetPluginInfoLatency: FAILED
Error: latency 112ms exceeds threshold 100ms
Details: Avg: 112ms, Min: 85ms, Max: 150ms
```

### Performance Optimization Tips

1. **Cache Metadata**: Return static/cached metadata (no external API calls)
2. **Avoid Blocking I/O**: Don't make database or network calls in `GetPluginInfo`
3. **Use In-Memory Data**: Store metadata in struct fields for fast access
4. **Validate Early**: Validate metadata at plugin initialization, not per-request

**Good Implementation**:

```go
type MyPlugin struct {
    name     string
    version  string
    specVer  string
    providers []string
}

func (p *MyPlugin) GetPluginInfo(
    ctx context.Context,
    req *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error) {
    // Fast: Return pre-computed values (no external API calls)
    return &pbc.GetPluginInfoResponse{
        Name:        p.name,
        Version:     p.version,
        SpecVersion: p.specVer,
        Providers:   p.providers,
    }, nil
}
```

**Bad Implementation**:

```go
func (p *MyPlugin) GetPluginInfo(
    ctx context.Context,
    req *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error) {
    // SLOW: Makes external API call on every request
    resp, err := http.Get("https://api.example.com/metadata")
    if err != nil {
        return nil, err
    }

    var metadata Metadata
    json.Unmarshal(resp.Body, &metadata)
    return &metadata, nil
}
```

---

## Running All Conformance Tests

### Standard Conformance

```bash
cd your-plugin

# Run standard conformance tests (includes GetPluginInfo performance)
go test -v ./... -run Conformance
```

### Advanced Conformance

```bash
# Run advanced conformance tests (50ms threshold)
CONFORMANCE_LEVEL=advanced go test -v ./... -run Conformance
```

### Specific Performance Test

```bash
# Run only GetPluginInfo performance test
go test -v ./... -run Performance_GetPluginInfoLatency
```

---

## Testing Your Implementation

### 1. Test Timeout Configuration

Create a slow mock server and verify timeout behavior:

```go
func TestClientTimeout(t *testing.T) {
    // Create slow server that sleeps for 10 seconds
    slowPlugin := &SlowPlugin{}

    harness := testing.NewTestHarness(slowPlugin)
    defer harness.Stop()

    conn, _ := harness.CreateClientConnection()
    defer conn.Close()

    client := pluginsdk.DefaultClientConfig(harness.BaseURL()).
        WithTimeout(1 * time.Second).  // 1-second timeout
        NewClient()

    // Call should timeout after 1 second
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    _, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
    if err == nil {
        t.Error("Expected timeout error, got nil")
    }
}
```

### 2. Test Error Messages

Create a plugin that returns invalid metadata:

```go
type InvalidPlugin struct{}

func (p *InvalidPlugin) GetPluginInfo(
    ctx context.Context,
    req *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error) {
    // Return incomplete metadata (missing version)
    return &pbc.GetPluginInfoResponse{
        Name:        "invalid-plugin",
        Version:     "",  // Missing!
        SpecVersion: "v1.2.0",
    }, nil
}

func TestInvalidMetadataError(t *testing.T) {
    plugin := &InvalidPlugin{}

    harness := testing.NewTestHarness(plugin)
    defer harness.Stop()

    conn, _ := harness.CreateClientConnection()
    defer conn.Close()

    client := pbc.NewCostSourceServiceClient(conn)

    _, err := client.GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
    if err == nil {
        t.Error("Expected error for incomplete metadata, got nil")
        return
    }

    st, ok := status.FromError(err)
    if !ok || st.Code() != codes.Internal {
        t.Errorf("Expected Internal error, got: %v", err)
        return
    }

    if st.Message() != "plugin metadata is incomplete" {
        t.Errorf("Expected 'plugin metadata is incomplete', got: %s", st.Message())
    }
}
```

### 3. Test Performance

Verify your `GetPluginInfo` implementation meets performance requirements:

```bash
# Run performance conformance test
go test -v -run Performance_GetPluginInfoLatency

# Expected output:
# Performance_GetPluginInfoLatency: PASSED
# Details: Avg: 45ms (threshold: 100ms)
```

---

## Common Issues and Solutions

### Issue 1: Timeout Not Working

**Symptom**: RPC calls don't timeout after configured duration.

**Solution**: Check if custom `HTTPClient` is overriding timeout:

```go
// WRONG: Timeout field ignored when HTTPClient is set
customClient := &http.Client{Timeout: 10 * time.Second}
cfg := pluginsdk.ClientConfig{
    HTTPClient: customClient,
    Timeout:    30 * time.Second,  // Ignored!
}

// CORRECT: Use WithTimeout() or set HTTPClient.Timeout
cfg := pluginsdk.DefaultClientConfig(url).WithTimeout(30 * time.Second)
// OR
customClient := &http.Client{Timeout: 30 * time.Second}
cfg := pluginsdk.ClientConfig{HTTPClient: customClient}
```

### Issue 2: GetPluginInfo Performance Test Fails

**Symptom**: `Performance_GetPluginInfoLatency` test fails with "latency exceeds threshold".

**Solution**: Profile your `GetPluginInfo` implementation:

```bash
# Run with CPU profiling
go test -cpuprofile=cpu.prof -run Performance_GetPluginInfoLatency

# Analyze profile
go tool pprof cpu.prof
```

**Common causes**:

- External API calls in `GetPluginInfo`
- Expensive validation logic
- Blocking I/O operations

**Fix**: Move initialization to plugin constructor, return static data.

### Issue 3: Error Messages Not User-Friendly

**Symptom**: Client receives technical error messages instead of user-friendly ones.

**Solution**: Verify server-side validation is enabled (automatic in SDK v0.4.15+):

```go
// Server automatically validates responses in sdk.go:323-399
// No manual validation needed in plugin code
func (s *Server) GetPluginInfo(ctx context.Context, req *pbc.GetPluginInfoRequest) (*pbc.GetPluginInfoResponse, error) {
    // SDK validates: nil response, incomplete metadata, invalid spec_version
    // You just return the response, SDK handles validation
    return provider.GetPluginInfo(ctx, req)
}
```

---

## Next Steps

1. **Update Your Plugin**: Add `PluginInfoProvider` implementation if not present
2. **Configure Timeouts**: Add timeout configuration to client initialization
3. **Run Tests**: Execute conformance tests to verify performance
4. **Deploy**: Deploy updated plugin with user-friendly errors

---

## Additional Resources

- **SDK Documentation**: `sdk/go/README.md`
- **Conformance Test Suite**: `sdk/go/testing/conformance_test.go`
- **Performance Testing**: `sdk/go/testing/performance.go`
- **Error Handling**: `sdk/go/pluginsdk/sdk.go:323-399`
- **Client Configuration**: `sdk/go/pluginsdk/client.go:82-132`

---

## Summary

SDK Polish v0.4.15 enhances the PulumiCost Go SDK with:

- ✅ Configurable per-client timeouts with context deadline support
- ✅ User-friendly GetPluginInfo error messages (automatic server-side validation)
- ✅ GetPluginInfo performance conformance tests (100ms Standard / 50ms Advanced)

All features are backward compatible and require no migration for existing plugins.
