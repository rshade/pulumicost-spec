# Research: gRPC-Web Support

**Branch**: `029-grpc-web-support`
**Date**: 2025-12-29
**Updated**: 2025-12-29 (Decision changed to connect-go)

## 1. Multi-Protocol Implementation

**Decision**: Use `connectrpc.com/connect` (connect-go).

**Rationale**:

- **Multi-Protocol**: Connect-go supports three protocols on the same handler:
  - **gRPC**: Full compatibility with existing `google.golang.org/grpc` clients
  - **gRPC-Web**: Browser clients using grpc-web protocol
  - **Connect**: Native HTTP/1.1 with JSON support (curl-friendly)
- **Actively Maintained**: Developed by Buf, the same team behind the buf CLI we already use
- **No Breaking Changes**: Existing gRPC clients continue to work unchanged
- **Built-in Features**: CORS support, health checks, and interceptors are first-class
- **JSON Support**: Enables debugging with curl and browser fetch API without special tooling
- **Pre-1.0 Flexibility**: At v0.4.11, we can adopt modern tooling without stability concerns

**Alternatives Considered**:

- **`improbable-eng/grpc-web`**: Wraps existing `*grpc.Server` without proto changes, but:
  - Less actively maintained
  - Requires manual CORS setup with `rs/cors`
  - No JSON/curl support
  - Only supports gRPC-Web protocol (not Connect)
- **`grpc-gateway`**: Generates a RESTful JSON API but requires proto annotations
  (`google.api.http`) which we don't have. Overkill for our needs.

**Migration Impact**:

- Add `buf.build/connectrpc/go` plugin to `buf.gen.yaml`
- Generated code includes new `*connect` package alongside existing `*grpc` code
- `Serve()` function internals change, but `Plugin` interface unchanged
- Existing gRPC clients work without modification (wire-compatible)

## 2. CORS Handling

**Decision**: Use connect-go's built-in CORS support via `connect.WithCORS()`.

**Rationale**:

- **Integrated**: No need for external CORS middleware
- **Correct Defaults**: Pre-configured for gRPC-Web requirements
- **Configuration**: Allowed origins passed to handler options

**Headers Exposed**:

- `grpc-status`, `grpc-message` (for gRPC error handling)
- `connect-*` headers (for Connect protocol)

## 3. Health Check

**Decision**: Use `connectrpc.com/grpchealth` for standard health checks plus `/healthz` HTTP endpoint.

**Rationale**:

- **grpchealth**: Provides `grpc.health.v1.Health/Check` compatible with both gRPC and Connect
- **HTTP /healthz**: Simple endpoint for load balancers and Kubernetes liveness probes
- **Unified Handler**: Both served from the same HTTP mux

## 4. Go Client Batching

**Decision**: Client-side concurrency using `golang.org/x/sync/errgroup`.

**Rationale**:

- **Simplicity**: No need for a complex worker pool library.
- **Control**: `errgroup` allows setting a concurrency limit (using `SetLimit` in newer Go
  versions or a semaphore channel) and cancels all requests on the first error (configurable).

## 5. Multi-Tenant Orchestration Support

**Decision**: Use Environment Variables for Credentials.

**Rationale**:

- **Isolation**: Each plugin instance is a process.
- **Standard**: Cloud SDKs (AWS SDK, Azure SDK) automatically pick up credentials from env vars
  (`AWS_ACCESS_KEY_ID`, etc.).
- **Implementation**: The SDK's `Serve` function doesn't need to do anything special here; it
  just ensures it doesn't cache credentials globally in a way that prevents process-level
  isolation (which is default behavior for Go).

## 6. Connect-go Protocol Details

### Wire Compatibility

| Protocol   | Content-Type                         | HTTP Version | Client Types         |
| ---------- | ------------------------------------ | ------------ | -------------------- |
| gRPC       | `application/grpc`                   | HTTP/2       | grpc-go, grpcurl     |
| gRPC-Web   | `application/grpc-web`               | HTTP/1.1+    | grpc-web, browser    |
| Connect    | `application/proto`, `application/json` | HTTP/1.1+    | curl, fetch, connect |

### Example curl Request

```bash
# JSON request (Connect protocol)
curl -X POST http://localhost:8080/pulumicost.v1.CostSourceService/Name \
  -H "Content-Type: application/json" \
  -d '{}'

# Binary protobuf request
curl -X POST http://localhost:8080/pulumicost.v1.CostSourceService/Name \
  -H "Content-Type: application/proto" \
  --data-binary @request.bin
```
