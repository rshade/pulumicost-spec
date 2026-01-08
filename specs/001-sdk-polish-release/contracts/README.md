# SDK Interface Contracts: v0.4.14

**Feature**: v0.4.14 SDK Polish Release
**Date**: 2025-01-04
**Purpose**: Document new and modified SDK interfaces introduced by this feature

## Overview

This document defines the Go SDK interface contracts that plugin developers interact with. These contracts are stable,
documented, and versioned commitments.

---

## HealthChecker Interface

**Package**: `sdk/go/pluginsdk`
**Version**: v0.4.14+
**Stability**: Stable (backward compatible)

### Interface Definition

```go
type HealthChecker interface {
    Check(ctx context.Context) error
}
```

### Contract Guarantees

**SDK Guarantees**:

- SDK automatically detects if a plugin implements `HealthChecker` via type assertion
- SDK calls `Check()` on health endpoint requests (`/healthz` HTTP and gRPC health service)
- SDK sets a context timeout for `Check()` (configurable via health endpoint settings)
- SDK catches panics in `Check()` and treats them as errors
- SDK returns HTTP 503 / gRPC `Unavailable` when `Check()` returns non-nil error
- SDK returns HTTP 200 / gRPC `OK` when `Check()` returns nil
- SDK populates `HealthStatus.LastChecked` timestamp automatically
- SDK never stores or persists `HealthStatus` (stateless per-request)

**Plugin Implementation Requirements**:

- `Check()` must respect context cancellation (return `context.Canceled` if context expired)
- `Check()` should be fast (<1 second ideally) to avoid health endpoint timeouts
- `Check()` can return any error message; SDK wraps in HTTP 503 response
- `Check()` implementation is optional; plugins not implementing it get default "always healthy"
- `Check()` should not modify plugin state (idempotent reads only recommended)

**Error Handling**:

| Scenario | SDK Behavior | HTTP Status | gRPC Status |
|----------|--------------|-------------|---------------|
| `Check()` returns nil | Healthy | 200 | OK |
| `Check()` returns error | Unhealthy | 503 | Unavailable |
| `Check()` times out | Unhealthy | 503 | Unavailable |
| `Check()` panics | Unhealthy | 503 | Unavailable |
| Plugin doesn't implement `HealthChecker` | Healthy (default) | 200 | OK |

### Example Implementation

```go
package main

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/rshade/pulumicost-sdk/pluginsdk"
)

type MyPlugin struct {
    db *sql.DB
}

func (p *MyPlugin) Check(ctx context.Context) error {
    // Check database connectivity with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    if err := p.db.PingContext(ctx); err != nil {
        return fmt.Errorf("database unavailable: %w", err)
    }
    return nil
}

func main() {
    plugin := &MyPlugin{db: dbPool}
    pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: plugin,
        Web:    pluginsdk.DefaultWebConfig().WithHealthEndpoint(true),
    })
}
```

### Breaking Changes

**None expected**. This is a new optional interface.

---

## Context Validation Helpers

**Package**: `sdk/go/pluginsdk`
**Version**: v0.4.14+
**Stability**: Stable (backward compatible)

### Function Contracts

#### ValidateContext

```go
func ValidateContext(ctx context.Context) error
```

**Contract Guarantees**:

- Returns non-nil error if `ctx == nil`: message "context cannot be nil"
- Returns non-nil error if `ctx.Err() != nil`: message "context already cancelled or expired: {err}"
- Returns `nil` if context is valid (not nil, not cancelled, not expired)
- Does not modify the context
- Thread-safe (no side effects)

**Example Usage**:

```go
if err := pluginsdk.ValidateContext(ctx); err != nil {
    return fmt.Errorf("invalid context: %w", err)
}
```

#### ContextRemainingTime

```go
func ContextRemainingTime(ctx context.Context) time.Duration
```

**Contract Guarantees**:

- Returns `time.Until(deadline)` if context has deadline
- Returns `time.Duration(math.MaxInt64)` if context has no deadline
- Returns negative duration if deadline already expired
- Does not modify the context
- Thread-safe (read-only context inspection)

**Example Usage**:

```go
remaining := pluginsdk.ContextRemainingTime(ctx)
if remaining < 5*time.Second {
    log.Warn("request timeout imminent", "remaining", remaining)
}
```

#### ContextDeadline

```go
func ContextDeadline(ctx context.Context) (time.Time, bool)
```

**Contract Guarantees**:

- Returns `(deadline, true)` if context has deadline set
- Returns `(time.Time{}, false)` if context has no deadline
- Does not modify the context
- Thread-safe (read-only context inspection)

**Example Usage**:

```go
deadline, hasDeadline := pluginsdk.ContextDeadline(ctx)
if hasDeadline {
    log.Info("request deadline", "deadline", deadline)
}
```

### Breaking Changes

**None expected**. These are new helper functions.

---

## ARN Format Helpers

**Package**: `sdk/go/pluginsdk`
**Version**: v0.4.14+
**Stability**: Stable (backward compatible)

### Function Contracts

#### DetectARNProvider

```go
func DetectARNProvider(arn string) string
```

**Contract Guarantees**:

- Returns `"aws"` if ARN starts with `"arn:aws:"`
- Returns `"azure"` if ARN starts with `"/subscriptions/"`
- Returns `"gcp"` if ARN starts with `"//"`
- Returns `"kubernetes"` if ARN contains `"{cluster}/{namespace}/"` pattern
- Returns `""` (empty string) for unrecognized formats (not an error)
- Returns error if ARN format is ambiguous (could match multiple providers)
- Does not modify input string
- Thread-safe (pure function, no side effects)
- Uses simple prefix matching (not regex)

**Examples**:

| Input | Output |
|-----------------------------------------------------|---------|
| `"arn:aws:ec2:us-east-1:123:instance/i-abc"` | `"aws"` |
| `"/subscriptions/sub-123/resourceGroups/rg/..."` | `"azure"`|
| `"//compute.googleapis.com/projects/proj/instances/i"` | `"gcp"` |
| `"my-cluster/namespace/pod/my-pod"` | `"kubernetes"` |
| `"custom:format"` | `""` |

#### ValidateARNConsistency

```go
func ValidateARNConsistency(arn, expectedProvider string) error
```

**Contract Guarantees**:

- Calls `DetectARNProvider(arn)` internally
- Returns `nil` if detected provider matches `expectedProvider`
- Returns `nil` if detected provider is `""` (unrecognized ARN)
- Returns non-nil error if detected provider doesn't match expected
- Error message format: "ARN format {arn} detected as {detected} but expected {expected}"
- Does not modify input strings
- Thread-safe (pure function, no side effects)

**Examples**:

```go
// Valid match
err := ValidateARNConsistency("arn:aws:ec2:...", "aws")
// err == nil

// Mismatch
err := ValidateARNConsistency("arn:aws:ec2:...", "azure")
// err.Error() == "ARN format arn:aws:ec2:... detected as aws but expected azure"

// Unrecognized ARN
err := ValidateARNConsistency("custom:format", "aws")
// err == nil (caller should handle unrecognized separately)
```

### Constants

```go
const (
    AWSARNPrefix     = "arn:aws:"
    AzureARNPrefix   = "/subscriptions/"
    GCPARNPrefix     = "//"
    KubernetesFormat = "{cluster}/{namespace}/"
)
```

**Contract Guarantees**:

- Constants are exported for documentation and testing purposes
- Values match detection patterns used by `DetectARNProvider()`
- Constants are immutable (go const)

### Breaking Changes

**None expected**. These are new helper functions and constants.

---

## Client Configuration

**Package**: `sdk/go/pluginsdk`
**Version**: v0.4.14+
**Stability**: Stable (backward compatible)

### ClientConfig.Timeout

```go
type ClientConfig struct {
    // ... existing fields ...
    Timeout time.Duration // 0 = use default 30s timeout
}
```

**Contract Guarantees**:

- `Timeout` field is optional; defaults to `0` if not set
- `0` value means use default 30-second timeout
- Positive values override default timeout
- Negative values are invalid (SDK returns error)
- Applies to all RPC methods unless overridden by context deadline
- Context deadline takes precedence over `ClientConfig.Timeout`
- Thread-safe (read-only after client initialization)

**Precedence Rules** (highest to lowest):

1. Context deadline (if set via `context.WithTimeout`)
2. `ClientConfig.Timeout` (if set to non-zero value)
3. Default 30-second timeout

**Example Usage**:

```go
// Configure 2-minute default timeout
client, err := pluginsdk.NewClient(ctx, pluginsdk.ClientConfig{
    Address: "localhost:50051",
    Timeout: 2 * time.Minute,
})

// All requests use 2-minute timeout unless context has deadline
resp, err := client.GetActualCost(ctx, req)

// Override with context deadline (takes precedence)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
resp, err = client.GetActualCost(ctx, req) // Uses 5-minute timeout
```

### Breaking Changes

**None expected**. This is a new optional field with safe default behavior.

---

## GetPluginInfo RPC Behavior

**Package**: `sdk/go/pluginsdk`
**Version**: v0.4.14+
**Stability**: Stable (backward compatible)

### Interface Detection

**SDK Guarantees**:

- SDK checks if plugin implements `GetPluginInfoProvider` interface via type assertion
- If implemented: calls plugin's `GetPluginInfo()` method
- If not implemented: returns gRPC `Unimplemented` status
- Error messages are user-friendly (no internal implementation details)
- Technical details logged server-side for debugging

```go
type GetPluginInfoProvider interface {
    GetPluginInfo(ctx context.Context, req *pbc.GetPluginInfoRequest) (*pbc.GetPluginInfoResponse, error)
}
```

### Error Message Mapping

| Internal Error                                | Client-Facing Message                              | gRPC Status Code        |
| --------------------------------------------- | -------------------------------------------------- | ----------------------- |
| `plugin returned nil response`                | `unable to retrieve plugin metadata`               | `codes.Internal`        |
| `plugin returned incomplete metadata`         | `plugin metadata is incomplete`                    | `codes.InvalidArgument` |
| `plugin returned invalid spec_version format` | `plugin reported an invalid specification version` | `codes.InvalidArgument` |
| Legacy plugin (not implementing interface)    | `GetPluginInfo not implemented by this plugin`     | `codes.Unimplemented`   |

**Contract Guarantees**:

- Client-facing errors never include internal implementation details
- Server-side logging includes technical details for debugging
- gRPC status codes follow standard semantics
- Backward compatible with legacy plugins (they return `Unimplemented`)

### Breaking Changes

**None expected**. This is a new RPC with backward-compatible fallback behavior.

---

## Summary of Contracts

| Contract Type    | Count | Stability | Breaking Changes |
| ---------------- | ----- | --------- | ---------------- |
| Interfaces       | 1     | Stable    | None             |
| Helper Functions | 5     | Stable    | None             |
| Config Options   | 1     | Stable    | None             |
| Constants        | 4     | Stable    | None             |
| RPC Behavior     | 1     | Stable    | None             |

**Total**: 12 stable contract additions with **0 breaking changes**.

All new contracts are **backward compatible** and **opt-in**:

- `HealthChecker`: Optional interface (default "always healthy" if not implemented)
- Context helpers: New functions (existing code unchanged)
- ARN helpers: New functions (existing code unchanged)
- `ClientConfig.Timeout`: Optional field (defaults to 30s if not set)
- `GetPluginInfo`: New RPC with fallback to `Unimplemented` for legacy plugins
