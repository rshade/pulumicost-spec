# Implementation Plan: Multi-Protocol Plugin Access

**Branch**: `029-grpc-web-support` | **Date**: 2025-12-29 | **Spec**: [specs/029-grpc-web-support/spec.md](./spec.md)
**Input**: Feature specification from `/specs/029-grpc-web-support/spec.md`
**Updated**: 2025-12-29 (Switched to connect-go)

## Summary

Enable multi-protocol support for FinFocus plugins using connect-go, allowing:

- **gRPC**: Existing clients continue to work unchanged
- **gRPC-Web**: Direct browser connectivity
- **Connect**: JSON/curl-friendly HTTP/1.1 access

Also includes a robust Go client SDK for batch operations and connection management.

## Technical Context

**Language/Version**: Go 1.25.5+
**Primary Dependencies**:

- `connectrpc.com/connect` - Multi-protocol RPC framework
- `connectrpc.com/grpchealth` - Standard health check support
- `golang.org/x/net/http2/h2c` - HTTP/2 cleartext for gRPC compatibility
- `golang.org/x/sync/errgroup` - Batch operation concurrency

**Storage**: N/A
**Testing**: Go `testing` package, `httptest` for integration tests
**Target Platform**: Linux/macOS/Windows (Plugin host), Web Browsers (Clients)
**Performance Goals**: Minimal overhead; <5s for 100-resource batch query
**Constraints**: Must support existing gRPC clients; wire-compatible with google.golang.org/grpc

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Proto-First**: Adds connect-go plugin to buf.gen.yaml (non-breaking)
- [x] **Performance**: Connect-go has minimal overhead, benchmarks required
- [x] **Documentation**: New client usage documented in README and examples
- [x] **Test-First**: Conformance tests for all three protocols
- [x] **Backward Compatible**: Existing gRPC clients work unchanged

## Project Structure

### Documentation (this feature)

```text
specs/029-grpc-web-support/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── tasks.md
```

### Source Code Changes

```text
# Modified files
buf.gen.yaml                    # Add connect-go plugin

# Generated files (new)
sdk/go/proto/finfocus/v1/
└── finfocusv1connect/
    └── costsource.connect.go   # Connect handlers (generated)

# SDK changes
sdk/go/pluginsdk/
├── sdk.go                      # Update Serve() for connect-go
├── options.go                  # WebConfig for CORS/origins
├── health.go                   # /healthz endpoint handler
├── connect.go                  # NEW: Connect handler setup
└── client/                     # NEW: Go client library
    ├── client.go               # Connection management
    ├── methods.go              # Typed RPC methods
    └── batch.go                # Batch operations

# Examples
examples/plugins/
└── web-client/
    ├── index.html              # Browser example
    └── README.md               # Usage instructions
```

## Architecture

### Server Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Server                            │
│  ┌───────────────┐  ┌─────────────────────────────────────┐ │
│  │   /healthz    │  │    Connect Handler (h2c wrapped)    │ │
│  │  (HTTP GET)   │  │  ┌─────────┐ ┌────────┐ ┌────────┐  │ │
│  └───────────────┘  │  │  gRPC   │ │gRPC-Web│ │Connect │  │ │
│                     │  │ (HTTP/2)│ │(HTTP/1)│ │ (JSON) │  │ │
│                     │  └─────────┘ └────────┘ └────────┘  │ │
│                     └─────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │  Plugin Server  │
                    │  (implements    │
                    │   Plugin iface) │
                    └─────────────────┘
```

### Client Architecture

```text
┌─────────────────────────────────────────┐
│           pluginsdk/client              │
│  ┌─────────────────────────────────────┐│
│  │           Client                    ││
│  │  - Connect client (HTTP)            ││
│  │  - Connection pooling               ││
│  │  - Batch operations (errgroup)      ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

## Key Implementation Details

### 1. buf.gen.yaml Changes

```yaml
plugins:
  # Existing
  - plugin: buf.build/protocolbuffers/go
    out: sdk/go/proto
    opt: paths=source_relative
  - plugin: buf.build/grpc/go
    out: sdk/go/proto
    opt: paths=source_relative
  # New
  - plugin: buf.build/connectrpc/go
    out: sdk/go/proto
    opt: paths=source_relative
```

### 2. Serve() Function Changes

The `Serve()` function will:

1. Create a connect handler from the plugin
2. Wrap with h2c for HTTP/2 cleartext (gRPC compatibility)
3. Add /healthz endpoint to mux
4. Apply CORS if WebConfig.Enabled
5. Start HTTP server

### 3. Backward Compatibility

- Existing `Plugin` interface unchanged
- Existing `ServeConfig` fields work as before
- New `WebConfig` field is optional (zero value = disabled)
- Existing gRPC clients connect to same port, same protocol

## Migration Notes

### For Plugin Developers

No changes required. Existing plugins work unchanged.

### For Client Developers

- **Existing gRPC clients**: Work unchanged
- **New browser clients**: Use grpc-web or Connect protocol
- **New Go clients**: Can use `pluginsdk/client` package or raw grpc-go
