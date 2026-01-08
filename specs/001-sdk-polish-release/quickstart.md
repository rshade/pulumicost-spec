# Quickstart: v0.4.14 SDK Polish Features

**Feature**: v0.4.14 SDK Polish Release
**Date**: 2025-01-04
**Audience**: Plugin developers
**Purpose**: Get started quickly with new SDK features

---

## Overview

v0.4.14 introduces developer experience improvements to help you build robust plugins:

1. **Custom Health Checking** - Implement real health checks for your plugin
2. **Context Validation** - Catch context errors early with clear messages
3. **ARN Detection & Validation** - Identify and validate cloud resource identifiers
4. **Configurable Timeouts** - Override default 30-second timeout for your clients
5. **GetPluginInfo** - Provide plugin metadata for discovery

All features are **opt-in** and **backward compatible** with existing plugins.

---

## Quick Examples

### 1. Custom Health Checking

Add a health check to verify your plugin's dependencies (database, APIs, etc.):

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/rshade/pulumicost-sdk/pluginsdk"
)

type MyPlugin struct {
    db *sql.DB
}

// Implement HealthChecker interface (opt-in)
func (p *MyPlugin) Check(ctx context.Context) error {
    // Verify database connectivity with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    if err := p.db.PingContext(ctx); err != nil {
        return fmt.Errorf("database unavailable: %w", err)
    }
    return nil
}

func main() {
    plugin := &MyPlugin{db: dbPool}

    // SDK automatically detects HealthChecker implementation
    pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: plugin,
        Web:    pluginsdk.DefaultWebConfig().WithHealthEndpoint(true),
    })
}
```

**Result**: Health endpoint (`/healthz`) returns database status instead of always "healthy".

---

### 2. Context Validation

Validate contexts before RPC calls to catch errors early:

```go
import "github.com/rshade/pulumicost-sdk/pluginsdk"

// Before: Could panic or return cryptic errors
resp, err := client.GetActualCost(ctx, req)

// After: Clear error messages
if err := pluginsdk.ValidateContext(ctx); err != nil {
    return fmt.Errorf("invalid context: %w", err)
}
resp, err := client.GetActualCost(ctx, req)

// Check remaining time for logging
remaining := pluginsdk.ContextRemainingTime(ctx)
if remaining < 5*time.Second {
    log.Warn("request timeout imminent", "remaining", remaining)
}
```

**Result**: Clear errors like "context cannot be nil" or "context already cancelled" instead of panics.

---

### 3. ARN Detection & Validation

Identify cloud provider from resource identifiers and validate consistency:

```go
import "github.com/rshade/pulumicost-sdk/pluginsdk"

// Detect provider from ARN
arn := "arn:aws:ec2:us-east-1:123:instance/i-abc"
provider := pluginsdk.DetectARNProvider(arn)
// provider == "aws"

// Validate ARN matches expected provider
expected := "aws"
if err := pluginsdk.ValidateARNConsistency(arn, expected); err != nil {
    log.Error("ARN provider mismatch", "error", err)
    // Error: "ARN format arn:aws:... detected as aws but expected azure"
}

// Works for all major providers
providers := []string{
    "arn:aws:ec2:...",                        // aws
    "/subscriptions/sub-123/resourceGroups/...",    // azure
    "//compute.googleapis.com/projects/proj/...",   // gcp
    "my-cluster/namespace/pod/my-pod",           // kubernetes
}
for _, arn := range providers {
    provider = pluginsdk.DetectARNProvider(arn)
    log.Info("Detected provider", "arn", arn, "provider", provider)
}
```

**Result**: Automatic provider detection and validation across AWS, Azure, GCP, and Kubernetes.

---

### 4. Configurable Timeouts

Override default 30-second timeout for operations that take longer:

```go
import "github.com/rshade/pulumicost-sdk/pluginsdk"

// Configure client with 2-minute default timeout
client, err := pluginsdk.NewClient(ctx, pluginsdk.ClientConfig{
    Address: "localhost:50051",
    Timeout: 2 * time.Minute, // Override default 30s
})
if err != nil {
    log.Fatal(err)
}

// All RPC methods use 2-minute timeout unless overridden
resp, err := client.GetActualCost(ctx, req)

// Override per-request with context (takes precedence)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
resp, err = client.GetActualCost(ctx, req) // Uses 5-minute timeout
```

**Result**: Flexible timeout configuration per-client and per-request.

---

### 5. GetPluginInfo (Metadata Discovery)

Provide plugin metadata for clients to discover capabilities:

```go
import "github.com/rshade/pulumicost-sdk/pluginsdk"

// Option 1: Static metadata (recommended for most plugins)
info := pluginsdk.NewPluginInfo("my-cost-plugin", "v1.0.0",
    pluginsdk.WithProviders("aws", "azure"),
    pluginsdk.WithDescription("Cost analysis for AWS and Azure"),
)

pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
    Plugin:     &MyPlugin{},
    PluginInfo: info, // Provide metadata
})

// Option 2: Dynamic metadata (advanced)
type MyPlugin struct{}

func (p *MyPlugin) GetPluginInfo(ctx context.Context, req *pbc.GetPluginInfoRequest) (
    *pbc.GetPluginInfoResponse, error) {
    // Compute metadata dynamically (e.g., discover supported providers)
    providers := discoverSupportedProviders()
    return &pbc.GetPluginInfoResponse{
        Name:        "my-cost-plugin",
        Version:      "v1.0.0",
        Providers:    providers,
        SpecVersion:  "v0.4.14",
    }, nil
}
```

**Result**: Clients can discover plugin metadata (name, version, providers, spec_version) via GetPluginInfo RPC.

---

## Migration from Legacy Plugins

Existing plugins work without changes. To enable new features:

### Step 1: Enable Health Checking (Optional)

Add `HealthChecker` interface to your plugin struct:

```go
type MyPlugin struct {
    // ... existing fields ...
}

func (p *MyPlugin) Check(ctx context.Context) error {
    // Your health check logic
    return nil // Healthy
}
```

SDK automatically detects this interface. No configuration needed.

### Step 2: Add GetPluginInfo (Optional)

Provide static metadata:

```go
pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
    Plugin:     &MyPlugin{},
    PluginInfo: pluginsdk.NewPluginInfo("my-plugin", "v1.0.0"),
})
```

### Step 3: Use Context Validation (Optional)

Wrap existing RPC calls:

```go
if err := pluginsdk.ValidateContext(ctx); err != nil {
    return fmt.Errorf("invalid context: %w", err)
}
// ... existing code ...
```

### Step 4: Use ARN Helpers (Optional)

Replace manual ARN parsing:

```go
// Before
provider := "aws"
if strings.HasPrefix(arn, "arn:aws:") {
    provider = "aws"
} else if strings.HasPrefix(arn, "/subscriptions/") {
    provider = "azure"
}

// After
provider := pluginsdk.DetectARNProvider(arn)
```

---

## Common Patterns

### Validating Contexts Before RPC Calls

```go
func MyRPCMethod(ctx context.Context) error {
    // Validate early
    if err := pluginsdk.ValidateContext(ctx); err != nil {
        return fmt.Errorf("invalid context: %w", err)
    }

    // ... rest of method ...
}
```

### Logging Health Check Details

```go
func (p *MyPlugin) Check(ctx context.Context) error {
    start := time.Now()
    defer func() {
        log.Info("health check completed", "duration", time.Since(start))
    }()

    // Your checks ...
    return nil
}
```

### Handling ARN Validation Errors

```go
func ValidateResourceARN(arn string) error {
    provider := pluginsdk.DetectARNProvider(arn)
    if provider == "" {
        return fmt.Errorf("unrecognized ARN format: %s", arn)
    }

    expected := "aws" // or detect from config
    if err := pluginsdk.ValidateARNConsistency(arn, expected); err != nil {
        return fmt.Errorf("ARN provider mismatch: %w", err)
    }
    return nil
}
```

### Configuring Timeouts per Client

```go
// Client for fast operations (default 30s)
fastClient, _ := pluginsdk.NewClient(ctx, pluginsdk.ClientConfig{
    Address: "localhost:50051",
    Timeout: 10 * time.Second,
})

// Client for slow operations
slowClient, _ := pluginsdk.NewClient(ctx, pluginsdk.ClientConfig{
    Address: "localhost:50051",
    Timeout: 5 * time.Minute,
})
```

---

## Testing Your Plugin

### Test Health Checking

```go
func TestMyPlugin_Check(t *testing.T) {
    plugin := &MyPlugin{db: mockDB}

    // Healthy database
    mockDB.ExpectPing().WillReturn(nil)
    err := plugin.Check(context.Background())
    if err != nil {
        t.Errorf("expected healthy, got error: %v", err)
    }

    // Unhealthy database
    mockDB.ExpectPing().WillReturn(errors.New("connection refused"))
    err = plugin.Check(context.Background())
    if err == nil {
        t.Error("expected unhealthy, got nil error")
    }
}
```

### Test Context Validation

```go
func TestValidateContext(t *testing.T) {
    // Nil context
    err := pluginsdk.ValidateContext(nil)
    if err == nil {
        t.Error("expected error for nil context")
    }

    // Cancelled context
    ctx, cancel := context.WithCancel(context.Background())
    cancel()
    err = pluginsdk.ValidateContext(ctx)
    if err == nil {
        t.Error("expected error for cancelled context")
    }

    // Valid context
    err = pluginsdk.ValidateContext(context.Background())
    if err != nil {
        t.Errorf("expected nil error for valid context, got: %v", err)
    }
}
```

### Test ARN Detection

```go
func TestDetectARNProvider(t *testing.T) {
    tests := []struct {
        arn      string
        provider string
    }{
        {"arn:aws:ec2:...", "aws"},
        {"/subscriptions/sub-123/...", "azure"},
        {"//compute.googleapis.com/...", "gcp"},
        {"my-cluster/namespace/pod/...", "kubernetes"},
        {"custom:format", ""},
    }

    for _, tt := range tests {
        got := pluginsdk.DetectARNProvider(tt.arn)
        if got != tt.provider {
            t.Errorf("DetectARNProvider(%q) = %q, want %q", tt.arn, got, tt.provider)
        }
    }
}
```

---

## Getting Help

- **Full Documentation**: See [contracts/README.md](./contracts/README.md) for detailed interface contracts
- **Migration Guide**: See `sdk/go/pluginsdk/README.md` → "Migrating to GetPluginInfo" section
- **Examples**: See `examples/` directory for complete plugin examples
- **Issues**: Report bugs or request features at <https://github.com/rshade/pulumicost-spec/issues>

---

## What's Next?

1. **Add Health Checking** to your plugin (1-2 hours)
2. **Use Context Validation** in your RPC methods (30 minutes)
3. **Validate ARN Formats** with ARN helpers (30 minutes)
4. **Configure Timeouts** for long-running operations (15 minutes)
5. **Provide GetPluginInfo** metadata (15 minutes)

Total estimated time: **3-4 hours** for typical plugin.

All features are **backward compatible** - you can adopt them incrementally without breaking existing functionality.
