# Quickstart: Configurable CORS Headers

**Feature**: 033-cors-headers-config
**Date**: 2026-01-04

## Overview

This guide shows how to customize CORS headers in your FinFocus plugin for security compliance
and observability requirements.

## Default Behavior

If you don't configure custom headers, the SDK uses sensible defaults compatible with Connect/gRPC-Web:

```go
// No changes needed - existing plugins work as-is
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    Web: pluginsdk.WebConfig{
        Enabled:        true,
        AllowedOrigins: []string{"https://app.example.com"},
    },
}
```

## Use Case 1: Minimal Security Footprint

Remove headers you don't need (e.g., no Authorization if using a different auth mechanism):

```go
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    Web: pluginsdk.DefaultWebConfig().
        WithWebEnabled(true).
        WithAllowedOrigins([]string{"https://app.example.com"}).
        WithAllowedHeaders([]string{
            "Accept",
            "Content-Type",
            "Content-Length",
            "Accept-Encoding",
            // Omitting: Authorization, X-CSRF-Token
            "Connect-Protocol-Version",
            "Connect-Timeout-Ms",
            "Grpc-Timeout",
            "X-Grpc-Web",
        }),
}
```

## Use Case 2: Add Observability Headers

Expose request correlation headers for client-side tracing:

```go
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    Web: pluginsdk.DefaultWebConfig().
        WithWebEnabled(true).
        WithAllowedOrigins([]string{"https://app.example.com"}).
        WithAllowedHeaders([]string{
            "Accept",
            "Content-Type",
            "X-Request-ID",        // Custom tracing header
            "X-Correlation-ID",    // Custom tracing header
            "Connect-Protocol-Version",
            "Connect-Timeout-Ms",
            "Grpc-Timeout",
            "X-Grpc-Web",
        }).
        WithExposedHeaders([]string{
            "Grpc-Status",
            "Grpc-Message",
            "Grpc-Status-Details-Bin",
            "X-Request-ID",        // Expose for client correlation
            "X-Trace-ID",          // Expose for client tracing
        }),
}
```

## Use Case 3: Compliance (Restrict All Custom Headers)

For strict security environments, allow only CORS-safelisted headers:

```go
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    Web: pluginsdk.DefaultWebConfig().
        WithWebEnabled(true).
        WithAllowedOrigins([]string{"https://app.example.com"}).
        WithAllowedHeaders([]string{}). // Empty = simple headers only
        WithExposedHeaders([]string{}), // Empty = no custom headers exposed
}
```

## Builder Method Chaining

All `With*` methods can be chained and applied in any order:

```go
web := pluginsdk.DefaultWebConfig().
    WithWebEnabled(true).
    WithAllowedOrigins([]string{"https://app.example.com"}).
    WithAllowCredentials(true).
    WithAllowedHeaders([]string{"Content-Type", "Authorization"}).
    WithExposedHeaders([]string{"X-Request-ID"}).
    WithHealthEndpoint(true)
```

## Default Headers Reference

### DefaultAllowedHeaders

These headers are used when `AllowedHeaders` is nil:

| Header | Purpose |
| ------ | ------- |
| Accept | Standard HTTP content negotiation |
| Content-Type | Request body MIME type |
| Content-Length | Request body size |
| Accept-Encoding | Compression support |
| Authorization | Bearer tokens, API keys |
| X-CSRF-Token | Cross-site request forgery protection |
| X-Requested-With | XMLHttpRequest indicator |
| Connect-Protocol-Version | Connect protocol version |
| Connect-Timeout-Ms | Connect request timeout |
| Grpc-Timeout | gRPC request timeout |
| X-Grpc-Web | gRPC-Web protocol indicator |
| X-User-Agent | Client identification |

### DefaultExposedHeaders

These headers are used when `ExposedHeaders` is nil:

| Header | Purpose |
| ------ | ------- |
| Grpc-Status | gRPC status code |
| Grpc-Message | gRPC status message |
| Grpc-Status-Details-Bin | gRPC error details (base64) |
| Connect-Content-Encoding | Connect response encoding |
| Connect-Content-Type | Connect response content type |

## Testing Your Configuration

Verify CORS headers with curl:

```bash
# Preflight request
curl -X OPTIONS https://your-plugin:8080/grpc \
  -H "Origin: https://app.example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization" \
  -v 2>&1 | grep -i "access-control"

# Expected output should show your configured headers:
# < Access-Control-Allow-Headers: Content-Type, Authorization
# < Access-Control-Expose-Headers: Grpc-Status, X-Request-ID
```
