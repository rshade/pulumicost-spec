# Data Model: Configurable CORS Headers

**Feature**: 033-cors-headers-config
**Date**: 2026-01-04

## Overview

This feature extends the existing `WebConfig` struct with two new optional fields. No new entities
are created; this is a pure extension of existing data structures.

## Entity: WebConfig (Extended)

### Current Fields (unchanged)

| Field | Type | Default | Description |
| ----- | ---- | ------- | ----------- |
| Enabled | bool | false | Enables gRPC-Web support |
| AllowedOrigins | []string | nil | Origins permitted for CORS |
| AllowCredentials | bool | false | Allow credentials in CORS requests |
| EnableHealthEndpoint | bool | false | Expose /healthz endpoint |

### New Fields

| Field | Type | Default | Description |
| ----- | ---- | ------- | ----------- |
| AllowedHeaders | []string | nil | Headers permitted in CORS requests. nil = use defaults |
| ExposedHeaders | []string | nil | Headers exposed to JavaScript. nil = use defaults |

### Field Semantics

#### AllowedHeaders

- **nil**: Use `DefaultAllowedHeaders` (Connect/gRPC-Web compatible set)
- **empty []string{}**: Set empty `Access-Control-Allow-Headers` (simple headers only)
- **populated []string**: Use exactly these headers, joined by ", "

#### ExposedHeaders

- **nil**: Use `DefaultExposedHeaders` (gRPC status headers)
- **empty []string{}**: Set empty `Access-Control-Expose-Headers`
- **populated []string**: Use exactly these headers, joined by ", "

## Constants

### DefaultAllowedHeaders

```text
Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,
X-CSRF-Token, X-Requested-With, Connect-Protocol-Version, Connect-Timeout-Ms,
Grpc-Timeout, X-Grpc-Web, X-User-Agent
```

### DefaultExposedHeaders

```text
Grpc-Status, Grpc-Message, Grpc-Status-Details-Bin,
Connect-Content-Encoding, Connect-Content-Type
```

## Validation Rules

| Rule | Scope | Enforcement |
| ---- | ----- | ----------- |
| No validation | AllowedHeaders | Caller responsibility |
| No validation | ExposedHeaders | Caller responsibility |
| Defensive copy | Both fields | At config time via builder methods |

## State Transitions

Not applicable. WebConfig is an immutable configuration struct with no state transitions.

## Relationships

```text
ServeConfig
    └── Web: WebConfig
            ├── AllowedHeaders []string (NEW)
            └── ExposedHeaders []string (NEW)
```

The WebConfig is embedded in ServeConfig and passed to corsMiddleware at server startup.
Configuration is read-only after server initialization.

## Migration Notes

### Backward Compatibility

Existing code continues to work without changes:

```go
// Before (still works - nil fields use defaults)
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    Web: pluginsdk.WebConfig{
        Enabled:        true,
        AllowedOrigins: []string{"https://app.example.com"},
    },
}

// After (new optional customization)
config := pluginsdk.ServeConfig{
    Plugin: myPlugin,
    Web: pluginsdk.DefaultWebConfig().
        WithWebEnabled(true).
        WithAllowedOrigins([]string{"https://app.example.com"}).
        WithAllowedHeaders([]string{"Content-Type", "X-Request-ID"}), // NEW
}
```

No breaking changes. Zero-value (nil) for new fields preserves existing behavior.
