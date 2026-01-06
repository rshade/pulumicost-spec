# Data Model: v0.4.14 SDK Polish Release

**Feature**: v0.4.14 SDK Polish Release
**Date**: 2025-01-04
**Purpose**: Define entities, interfaces, and data structures introduced by this feature

## Overview

This feature introduces new interfaces, helper functions, and validation patterns to the Go SDK. The data model consists of:

- **Interfaces**: HealthChecker for custom health checking
- **Structs**: HealthStatus for detailed health information
- **Helper Functions**: Context validation, ARN format detection/validation
- **Configuration**: ClientConfig.Timeout option
- **Patterns**: Error handling, timeout configuration, backward compatibility

No new persistent data storage is introduced (this is SDK-only changes).

---

## Interface Definitions

### HealthChecker

**Package**: `sdk/go/pluginsdk`
**Purpose**: Allows plugins to provide custom health check logic

```go
type HealthChecker interface {
    // Check performs a health check and returns nil if healthy,
    // or an error describing the health issue.
    Check(ctx context.Context) error
}
```

**Implementation Requirements**:

- Context parameter allows timeout configuration for health checks
- Error return indicates unhealthy state; nil indicates healthy
- SDK automatically detects and uses this interface if implemented by plugin
- Timeout/panic in Check() is caught and returns HTTP 503 / gRPC Unavailable

**Usage Pattern**:

```go
type MyPlugin struct {
    dbPool *sql.DB
}

func (p *MyPlugin) Check(ctx context.Context) error {
    // Verify database connectivity
    if err := p.dbPool.PingContext(ctx); err != nil {
        return fmt.Errorf("database unavailable: %w", err)
    }
    return nil
}

// SDK automatically detects HealthChecker implementation
pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
    Plugin: &MyPlugin{dbPool: pool},
    Web:    pluginsdk.DefaultWebConfig().WithHealthEndpoint(true),
})
```

**Validation Rules**:

- Check() timeout is configured by caller (SDK or manual health endpoint)
- Panics in Check() are caught and treated as errors
- Context cancellation is respected (Check() should return context.Canceled)

**State Transitions**:

- N/A (HealthChecker is stateless, Check() is called on each health request)

---

## Struct Definitions

### HealthStatus

**Package**: `sdk/go/pluginsdk`
**Purpose**: Detailed health information returned by health endpoints

```go
type HealthStatus struct {
    Healthy     bool              `json:"healthy"`
    Message     string            `json:"message,omitempty"`
    Details     map[string]string `json:"details,omitempty"`
    LastChecked time.Time         `json:"last_checked"`
}
```

**Field Descriptions**:
| Field | Type | Description |
|--------------|---------------------|----------------------------------------------------------|
| Healthy | bool | `true` if system is healthy, `false` otherwise |
| Message | string | Human-readable status message (optional) |
| Details | map[string]string | Additional diagnostic information (optional) |
| LastChecked | time.Time | Timestamp of last health check (populated by SDK) |

**Validation Rules**:

- `Message` is optional; if empty, defaults to "healthy" or "unhealthy"
- `Details` is optional; map can be empty or nil
- `LastChecked` is populated automatically by SDK (not set by plugin)

**Usage Pattern**:

```go
// Plugin can return custom details
func (p *MyPlugin) Check(ctx context.Context) error {
    start := time.Now()
    if err := p.dbPool.PingContext(ctx); err != nil {
        return err
    }
    latency := time.Since(start)

    // SDK converts error to HealthStatus with details
    // (SDK populates LastChecked, Message based on error)
    return nil
}

// SDK constructs response:
// HealthStatus{
//     Healthy: true,
//     Message: "",
//     Details: map[string]string{"latency": latency.String()},
//     LastChecked: time.Now(),
// }
```

**JSON Representation**:

```json
{
  "healthy": true,
  "message": "",
  "details": {
    "latency": "12ms"
  },
  "last_checked": "2025-01-04T10:30:00Z"
}
```

---

## Configuration Structures

### ClientConfig.Timeout

**Package**: `sdk/go/pluginsdk`
**Purpose**: Override default 30-second timeout for all client RPC methods

```go
type ClientConfig struct {
    // ... existing fields ...
    Timeout time.Duration // 0 = use default 30s timeout
}
```

**Field Description**:
| Field | Type | Description |
|---------|--------------|-------------------------------------------------------|
| Timeout | time.Duration | Per-client default timeout (0 = use default 30s) |

**Validation Rules**:

- `0` value means use default 30-second timeout
- Negative values are invalid (should error)
- Very large values (>1 hour) should be rejected or warned

**Usage Pattern**:

```go
// Option 1: Configure default timeout for all requests
client, err := pluginsdk.NewClient(ctx, pluginsdk.ClientConfig{
    Address: "localhost:50051",
    Timeout: 2 * time.Minute, // Override default to 2 minutes
})
if err != nil {
    log.Fatal(err)
}

// All RPC methods use 2-minute timeout unless overridden by context
resp, err := client.GetActualCost(ctx, req)

// Option 2: Per-request timeout via context (takes precedence)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
resp, err := client.GetActualCost(ctx, req) // Uses 5-minute timeout
```

**Precedence Rules**:

1. Context deadline (if set) takes highest precedence
2. ClientConfig.Timeout used if context has no deadline
3. Default 30-second timeout used if neither is set

---

## Helper Functions

### Context Validation

**Package**: `sdk/go/pluginsdk`
**Purpose**: Validate context usability before RPC calls

```go
// ValidateContext checks that a context is usable for RPC calls.
func ValidateContext(ctx context.Context) error

// ContextRemainingTime returns time until deadline, negative if expired,
// or math.MaxInt64 duration if no deadline set.
func ContextRemainingTime(ctx context.Context) time.Duration

// ContextDeadline returns the context deadline or zero time if none set.
func ContextDeadline(ctx context.Context) (time.Time, bool)
```

**Validation Rules**:

**ValidateContext**:

- Returns error if `ctx == nil`: "context cannot be nil"
- Returns error if `ctx.Err() != nil`: "context already cancelled or expired: {err}"
- Returns `nil` if context is valid

**ContextRemainingTime**:

- Returns `time.Until(deadline)` if deadline is set
- Returns `time.Duration(math.MaxInt64)` if no deadline set
- Returns negative value if deadline already expired

**ContextDeadline**:

- Returns `(deadline, true)` if deadline is set
- Returns `(time.Time{}, false)` if no deadline set

**Usage Patterns**:

```go
// Validate before RPC call
if err := pluginsdk.ValidateContext(ctx); err != nil {
    return fmt.Errorf("invalid context: %w", err)
}
resp, err := client.GetActualCost(ctx, req)

// Check remaining time for logging
remaining := pluginsdk.ContextRemainingTime(ctx)
if remaining < 5*time.Second {
    log.Warn("request timeout imminent", "remaining", remaining)
}

// Check if deadline is set
deadline, hasDeadline := pluginsdk.ContextDeadline(ctx)
if hasDeadline {
    log.Info("request has deadline", "deadline", deadline)
}
```

---

### ARN Format Detection and Validation

**Package**: `sdk/go/pluginsdk`
**Purpose**: Detect cloud provider from resource identifier format and validate consistency

```go
// DetectARNProvider returns the cloud provider inferred from ARN format.
// Returns empty string if format is unrecognized.
//
// Examples:
//   - "arn:aws:..." → "aws"
//   - "/subscriptions/..." → "azure"
//   - "//compute.googleapis.com/..." → "gcp"
//   - "{cluster}/{namespace}/..." → "kubernetes"
func DetectARNProvider(arn string) string

// ValidateARNConsistency checks if ARN format matches expected provider.
// Returns nil if consistent, error describing mismatch otherwise.
func ValidateARNConsistency(arn, expectedProvider string) error

// ARN format patterns (exported for documentation)
const (
    AWSARNPrefix     = "arn:aws:"
    AzureARNPrefix   = "/subscriptions/"
    GCPARNPrefix     = "//"
    KubernetesFormat = "{cluster}/{namespace}/"
)
```

**Validation Rules**:

**DetectARNProvider**:

- Uses prefix/pattern matching (not regex)
- Returns exact string: `"aws"`, `"azure"`, `"gcp"`, `"kubernetes"`, or `""`
- Empty string returned for unrecognized formats (not an error)
- Ambiguous formats return explicit error (e.g., "ARN format ambiguous, could be multiple providers")

**ValidateARNConsistency**:

- Calls `DetectARNProvider(arn)` to detect actual provider
- Returns error if detected provider doesn't match expected
- Error format: "ARN format {arn} detected as {detected} but expected {expected}"
- Returns nil if providers match or ARN is unrecognized (handled by caller)

**ARN Format Patterns**:

| Provider   | Prefix/Pattern             | Example                                                |
| ---------- | -------------------------- | ------------------------------------------------------ |
| AWS        | `"arn:aws:"`               | `"arn:aws:ec2:us-east-1:123:instance/i-abc"`           |
| Azure      | `"/subscriptions/"`        | `"/subscriptions/sub-123/resourceGroups/rg/..."`       |
| GCP        | `"//"`                     | `"//compute.googleapis.com/projects/proj/instances/i"` |
| Kubernetes | `"{cluster}/{namespace}/"` | `"my-cluster/namespace/pod/my-pod"`                    |

**Usage Patterns**:

```go
// Detect provider from ARN
arn := "arn:aws:ec2:us-east-1:123:instance/i-abc"
provider := pluginsdk.DetectARNProvider(arn)
// provider == "aws"

// Validate consistency
expected := "aws"
if err := pluginsdk.ValidateARNConsistency(arn, expected); err != nil {
    log.Error("ARN provider mismatch", "arn", arn, "error", err)
}

// Handle unrecognized formats
unknownArn := "custom:format"
provider = pluginsdk.DetectARNProvider(unknownArn)
// provider == "" (empty string)
```

---

## Error Handling Patterns

### GetPluginInfo Error Messages

**Package**: `sdk/go/pluginsdk`
**Purpose**: User-friendly error messages for GetPluginInfo RPC

**Error Mapping**:

| Internal Error                                | Client-Facing Error                                | gRPC Status Code        |
| --------------------------------------------- | -------------------------------------------------- | ----------------------- |
| `plugin returned nil response`                | `unable to retrieve plugin metadata`               | `codes.Internal`        |
| `plugin returned incomplete metadata`         | `plugin metadata is incomplete`                    | `codes.InvalidArgument` |
| `plugin returned invalid spec_version format` | `plugin reported an invalid specification version` | `codes.InvalidArgument` |

**Server-Side Logging**:

- Technical error details logged with `logger.Error()` for debugging
- Includes structured fields: error, debug info, stack trace if available
- Log entry examples:
  ```go
  logger.Error("GetPluginInfo failed",
      "error", err,
      "debug", map[string]any{
          "plugin_name": pluginName,
          "stack": debug.Stack(),
      })
  ```

**Client-Side Error**:

- gRPC status code indicates error category
- Status message is user-friendly (no internal details)
- Status details can include additional metadata if needed

---

### Health Check Error Handling

**Package**: `sdk/go/pluginsdk`
**Purpose**: Handle health check timeouts and panics

**Error Handling Rules**:

- If `HealthChecker.Check()` times out (context deadline exceeded): return HTTP 503 / gRPC `Unavailable`
- If `HealthChecker.Check()` panics: recover and return HTTP 503 / gRPC `Unavailable`
- If `HealthChecker.Check()` returns error: return HTTP 503 / gRPC `Unavailable` with error message
- If plugin doesn't implement `HealthChecker`: return HTTP 200 / gRPC `OK` with `"healthy": true`

**Response Examples**:

**Healthy (Custom HealthChecker)**:

```json
{
  "healthy": true,
  "message": "",
  "details": {
    "database": "connected",
    "latency": "12ms"
  },
  "last_checked": "2025-01-04T10:30:00Z"
}
```

**Unhealthy (Custom HealthChecker)**:

```json
{
  "healthy": false,
  "message": "database unavailable: connection timeout",
  "details": {
    "database": "disconnected"
  },
  "last_checked": "2025-01-04T10:30:00Z"
}
```

**Timeout/Panic (Custom HealthChecker)**:

```json
{
  "healthy": false,
  "message": "health check timed out",
  "last_checked": "2025-01-04T10:30:00Z"
}
```

HTTP Status: 503 / gRPC Status: `Unavailable`

**Default (No HealthChecker)**:

```json
{
  "healthy": true,
  "message": "",
  "last_checked": "2025-01-04T10:30:00Z"
}
```

HTTP Status: 200 / gRPC Status: `OK`

---

## Backward Compatibility

### GetPluginInfo RPC

**Package**: `sdk/go/pluginsdk`
**Purpose**: Support legacy plugins not implementing GetPluginInfo

**Interface Detection Pattern**:

```go
// Check if plugin implements GetPluginInfo
type GetPluginInfoProvider interface {
    GetPluginInfo(ctx context.Context, req *pbc.GetPluginInfoRequest) (*pbc.GetPluginInfoResponse, error)
}

// In Server implementation
func (s *Server) GetPluginInfo(ctx context.Context, req *pbc.GetPluginInfoRequest) (*pbc.GetPluginInfoResponse, error) {
    if infoProvider, ok := s.plugin.(GetPluginInfoProvider); ok {
        return infoProvider.GetPluginInfo(ctx, req)
    }
    // Legacy plugin: return Unimplemented
    return nil, status.Error(codes.Unimplemented, "GetPluginInfo not implemented by this plugin")
}
```

**Behavior**:

- New plugins implementing `GetPluginInfoProvider`: return actual metadata
- Legacy plugins (not implementing interface): return `Unimplemented` status
- Client should handle `Unimplemented` gracefully
- No breaking changes to existing plugins

---

## Summary

**New Interfaces**: 1 (`HealthChecker`)
**New Structs**: 1 (`HealthStatus`)
**New Helper Functions**: 5 (`ValidateContext`, `ContextRemainingTime`, `ContextDeadline`, `DetectARNProvider`, `ValidateARNConsistency`)
**New Configuration Options**: 1 (`ClientConfig.Timeout`)
**New Constants**: 4 (`AWSARNPrefix`, `AzureARNPrefix`, `GCPARNPrefix`, `KubernetesFormat`)

**No Persistent Data Storage**: This feature is SDK-only changes (no database, files, or storage).

**Validation Rules**:

- Context must not be nil or expired before RPC calls
- ARN format must match expected provider
- Health check timeouts/panics return HTTP 503 / gRPC `Unavailable`
- Infinity/NaN values rejected in cost validation
- Payloads >1MB rejected in Connect protocol

**Backward Compatibility**:

- Optional interfaces detected via type assertion
- Default behaviors safe for non-implementing plugins
- GetPluginInfo returns `Unimplemented` for legacy plugins
