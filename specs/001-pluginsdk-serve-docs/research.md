# Research: pluginsdk.Serve() Documentation

**Feature**: 001-pluginsdk-serve-docs
**Date**: 2025-12-08
**Status**: Complete

## Overview

This document captures the research findings from analyzing the existing `pluginsdk` implementation
to inform documentation creation.

## Serve() Function Analysis

### Function Signature

```go
func Serve(ctx context.Context, config ServeConfig) error
```

**Source**: `sdk/go/pluginsdk/sdk.go:384`

### ServeConfig Struct

```go
type ServeConfig struct {
    Plugin            Plugin
    Port              int                           // If 0, uses env var or random port
    Registry          RegistryLookup                // Optional; if nil, DefaultRegistryLookup
    Logger            *zerolog.Logger               // Optional; if nil, default logger
    UnaryInterceptors []grpc.UnaryServerInterceptor // Optional; chained after tracing
}
```

**Source**: `sdk/go/pluginsdk/sdk.go:308-320`

### Port Resolution Priority

The port is resolved via `resolvePort()` function with this priority:

1. **ServeConfig.Port** - If > 0, use this value (set from --port flag via ParsePortFlag() or
   explicitly)
2. **PULUMICOST_PLUGIN_PORT** - Environment variable (via GetPort())
3. **Ephemeral (0)** - OS assigns available port

**Decision**: PORT environment variable is NOT supported to avoid multi-plugin conflicts.
**Rationale**: When pulumicost-core spawns multiple plugins (e.g., aws-public + aws-ce), each needs
a unique port. Using --port flag allows the core to allocate distinct ports.
**Source**: `sdk/go/pluginsdk/sdk.go:322-336`

### Port Announcement

The server outputs `PORT=<port>` to stdout immediately after binding:

```go
fmt.Fprintf(os.Stdout, "PORT=%d\n", addr.Port)
```

**Source**: `sdk/go/pluginsdk/sdk.go:363-375`

### Graceful Shutdown

Context cancellation triggers graceful shutdown:

```go
go func() {
    <-ctx.Done()
    grpcServer.GracefulStop()
}()
```

**Decision**: Use GracefulStop() not Stop() to allow in-flight requests to complete.
**Source**: `sdk/go/pluginsdk/sdk.go:407-411`

### Interceptor Chain

Interceptors are chained in this order:

1. Built-in `TracingUnaryServerInterceptor()` (always first)
2. User-provided interceptors from `config.UnaryInterceptors` (in order)

**Source**: `sdk/go/pluginsdk/sdk.go:395-403`

### Loopback Binding

Server binds to loopback address only (127.0.0.1):

```go
address := "127.0.0.1:0"
if port > 0 {
    address = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
}
```

**Decision**: Bind to loopback only for security (plugins should not accept external connections).
**Source**: `sdk/go/pluginsdk/sdk.go:338-361`

## ParsePortFlag() Analysis

### Function

```go
func ParsePortFlag() int
```

**Requirement**: Caller MUST call `flag.Parse()` before calling this function.
**Returns**: 0 if flag not specified or flag.Parse() not called.
**Source**: `sdk/go/pluginsdk/sdk.go:26-49`

### Flag Registration

The --port flag is registered at package initialization:

```go
var portFlag = flag.Int("port", 0, "TCP port for gRPC server (overrides PULUMICOST_PLUGIN_PORT)")
```

**Source**: `sdk/go/pluginsdk/sdk.go:24`

## Environment Variables

### Complete List

| Variable | Constant | Purpose | Default |
|----------|----------|---------|---------|
| PULUMICOST_PLUGIN_PORT | EnvPort | gRPC server port | 0 (ephemeral) |
| PULUMICOST_LOG_LEVEL | EnvLogLevel | Log verbosity | (none) |
| LOG_LEVEL | EnvLogLevelFallback | Legacy log level fallback | (none) |
| PULUMICOST_LOG_FORMAT | EnvLogFormat | Log format (json/text) | (none) |
| PULUMICOST_LOG_FILE | EnvLogFile | Log file path | stdout |
| PULUMICOST_TRACE_ID | EnvTraceID | Distributed trace ID | (none) |
| PULUMICOST_TEST_MODE | EnvTestMode | Test mode flag | false |

**Source**: `sdk/go/pluginsdk/env.go:11-41`

### Getter Functions

| Function | Returns | Fallback |
|----------|---------|----------|
| GetPort() | int | 0 if not set/invalid |
| GetLogLevel() | string | LOG_LEVEL if PULUMICOST_LOG_LEVEL not set |
| GetLogFormat() | string | empty if not set |
| GetLogFile() | string | empty (stdout) |
| GetTraceID() | string | empty if not set |
| GetTestMode() | bool | false (logs warning for invalid values) |
| IsTestMode() | bool | false (no warning) |

**Source**: `sdk/go/pluginsdk/env.go:43-111`

## Error Conditions

### Listener Creation Failure

Returns wrapped error: `"failed to listen: %w"`

**Causes**:

- Port already in use (EADDRINUSE)
- Permission denied (EPERM)
- Invalid port number (handled by resolvePort returning 0)

**Source**: `sdk/go/pluginsdk/sdk.go:344-345`

### Port Announcement Failure

Returns wrapped error: `"writing port: %w"`

**Cause**: stdout write failure (extremely rare)

**Source**: `sdk/go/pluginsdk/sdk.go:372`

### Context Cancellation

Returns `ctx.Err()` when context is canceled before or during serve.

**Source**: `sdk/go/pluginsdk/sdk.go:416-417`

### gRPC Server Failure

Returns the underlying gRPC error if server fails for reasons other than context cancellation.

**Source**: `sdk/go/pluginsdk/sdk.go:419`

## Documentation Locations

### Existing Documentation

1. **Go doc comments** in `sdk.go` - Basic function descriptions
2. **CLAUDE.md** in `sdk/go/` - Developer context file (mentions env vars)
3. **CLAUDE.md** in `sdk/go/pluginsdk/` - Package-specific notes (if exists)

### Target Documentation

1. **sdk/go/pluginsdk/README.md** - Primary target for comprehensive Serve() documentation
2. **Go doc comments** - Enhance existing comments with more detail
3. **Example code** - Verified, copy-paste-ready main() function

## Alternatives Considered

### PORT Fallback (Rejected)

**Alternative**: Support generic PORT environment variable as fallback
**Rejected Because**: Multi-plugin orchestration requires unique ports per plugin. Generic PORT
would cause conflicts when pulumicost-core spawns multiple plugins.

### External Binding (Rejected)

**Alternative**: Bind to 0.0.0.0 for external access
**Rejected Because**: Plugins are subprocess components, not standalone services. External access
is a security risk.

### Immediate Stop (Rejected)

**Alternative**: Use grpcServer.Stop() instead of GracefulStop()
**Rejected Because**: Would drop in-flight requests. GracefulStop() allows existing requests to
complete before shutdown.
