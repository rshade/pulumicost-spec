# Data Model: Standardized Plugin Metrics

**Feature**: 014-plugin-metrics
**Date**: 2025-12-02

## Overview

This feature introduces Prometheus metrics instrumentation to the pluginsdk package. The data model
consists of metric types, their labels, and configuration structures.

## Entities

### 1. Request Counter (Counter Metric)

Tracks total gRPC requests received by the plugin.

**Metric Name**: `pulumicost_plugin_requests_total`

**Fields/Labels**:

| Label | Type | Description | Cardinality |
|-------|------|-------------|-------------|
| `grpc_method` | string | gRPC method name (e.g., "GetProjectedCost") | 6 (fixed) |
| `grpc_code` | string | gRPC status code (e.g., "OK", "Internal") | ~17 (bounded) |
| `plugin_name` | string | Plugin identifier provided at creation | 1 per instance |

**Validation Rules**:

- `grpc_method`: Extracted from `grpc.UnaryServerInfo.FullMethod`, base name only
- `grpc_code`: Derived from `google.golang.org/grpc/status.Code(err).String()`
- `plugin_name`: Provided at interceptor creation, must be non-empty

**Lifecycle**:

- Created: When `MetricsUnaryServerInterceptor` is instantiated
- Updated: Incremented by 1 after each RPC completes (success or failure)
- Collected: When `/metrics` endpoint is scraped

### 2. Duration Histogram (Histogram Metric)

Tracks request latency distribution.

**Metric Name**: `pulumicost_plugin_request_duration_seconds`

**Fields/Labels**:

| Label | Type | Description | Cardinality |
|-------|------|-------------|-------------|
| `grpc_method` | string | gRPC method name | 6 (fixed) |
| `plugin_name` | string | Plugin identifier | 1 per instance |

**Buckets** (fixed, in seconds):

```text
0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0
```

Equivalent to: 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s

**Validation Rules**:

- Duration must be non-negative (enforced by `time.Since()`)
- Recorded in seconds (float64) for Prometheus compatibility

**Lifecycle**:

- Created: When `MetricsUnaryServerInterceptor` is instantiated
- Updated: `Observe()` called with duration after each RPC completes
- Collected: When `/metrics` endpoint is scraped

### 3. Metrics Registry

Container for all plugin metrics.

**Type**: `*prometheus.Registry`

**Relationships**:

- Contains: Request Counter, Duration Histogram
- Used by: Metrics HTTP Server helper

**Lifecycle**:

- Created: When `NewPluginMetrics()` is called
- Destroyed: When plugin process terminates

### 4. Metrics Server Configuration

Configuration for optional HTTP metrics server.

**Fields**:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Port` | int | 9090 | HTTP port for metrics endpoint |
| `Registry` | *prometheus.Registry | nil | Custom registry (uses metrics registry if nil) |
| `Path` | string | "/metrics" | URL path for metrics endpoint |

**Validation Rules**:

- `Port`: Must be valid port number (1-65535)
- `Path`: Must start with "/"

## State Transitions

### Interceptor State

```text
┌─────────────┐
│   Created   │ ← MetricsUnaryServerInterceptor(pluginName)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Active    │ ← Interceptor added to gRPC server chain
└──────┬──────┘
       │ (each RPC)
       ▼
┌─────────────────────────────────────────────┐
│  Record: start time                          │
│  Execute: handler(ctx, req)                  │
│  Record: duration, increment counter         │
└─────────────────────────────────────────────┘
```

### Metrics Server State

```text
┌─────────────┐
│   Created   │ ← StartMetricsServer(config)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Listening  │ ← http.ListenAndServe()
└──────┬──────┘
       │ (shutdown signal)
       ▼
┌─────────────┐
│   Stopped   │ ← server.Shutdown()
└─────────────┘
```

## Cardinality Analysis

Total unique metric series (worst case):

- **Request Counter**: 6 methods × 17 codes × 1 plugin = **102 series**
- **Duration Histogram**: 6 methods × 1 plugin × (10 buckets + sum + count) = **72 series**
- **Total per plugin**: ~174 series (well within safe limits)

## Integration Points

### With Existing TracingUnaryServerInterceptor

Both interceptors can be chained:

```text
Request → TracingInterceptor → MetricsInterceptor → Handler → Response
                │                      │
                ▼                      ▼
          Context with            Metrics recorded
           trace_id
```

### With ServeConfig

```text
ServeConfig.UnaryInterceptors = [
    MetricsUnaryServerInterceptor("plugin-name"),
    // other custom interceptors...
]
```

The SDK's `Serve()` function chains TracingInterceptor first, then user interceptors.
