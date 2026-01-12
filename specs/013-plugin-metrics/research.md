# Research: Standardized Plugin Metrics

**Feature**: 014-plugin-metrics
**Date**: 2025-12-02

## Research Topics

### 1. Prometheus Go Client Library Best Practices

**Decision**: Use `github.com/prometheus/client_golang` with custom registry pattern

**Rationale**:

- Industry standard for Go metrics instrumentation
- Supports both Counter and Histogram metric types required by spec
- Custom registry pattern (`prometheus.NewRegistry()`) avoids global state pollution
- Provides `promhttp.HandlerFor()` for flexible HTTP exposure

**Alternatives Considered**:

- OpenTelemetry metrics: More complex, overkill for this use case
- Custom metrics implementation: Would duplicate existing well-tested code

**Key Implementation Pattern**:

```go
reg := prometheus.NewRegistry()
counter := prometheus.NewCounterVec(prometheus.CounterOpts{
    Name: "finfocus_plugin_requests_total",
    Help: "Total gRPC requests",
}, []string{"grpc_method", "grpc_code", "plugin_name"})
reg.MustRegister(counter)
```

**Sources**:

- [Prometheus Client Golang](https://github.com/prometheus/client_golang)

### 2. gRPC Server Interceptor Patterns for Metrics

**Decision**: Follow `grpc-ecosystem/go-grpc-prometheus` naming conventions and label structure

**Rationale**:

- Well-established patterns from [go-grpc-prometheus](https://github.com/grpc-ecosystem/go-grpc-prometheus)
- Standard labels: `grpc_method`, `grpc_code` are industry conventions
- Interceptor pattern integrates cleanly with existing `TracingUnaryServerInterceptor`
- Unary interceptor sufficient (FinFocus uses unary RPCs only)

**Alternatives Considered**:

- Using go-grpc-prometheus directly: Would add external dependency; we want minimal, focused metrics
- go-grpc-middleware/providers/openmetrics: More complex than needed

**Key Implementation Pattern**:

```go
func MetricsUnaryServerInterceptor(pluginName string) grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler) (interface{}, error) {

        start := time.Now()
        resp, err := handler(ctx, req)
        duration := time.Since(start)

        code := status.Code(err)
        method := path.Base(info.FullMethod)

        requestsTotal.WithLabelValues(method, code.String(), pluginName).Inc()
        requestDuration.WithLabelValues(method, pluginName).Observe(duration.Seconds())

        return resp, err
    }
}
```

**Sources**:

- [go-grpc-prometheus server_metrics.go](https://github.com/grpc-ecosystem/go-grpc-prometheus/blob/master/server_metrics.go)
- [gRPC Custom Prometheus Metrics Guide](https://medium.com/@sonu.sonu75/how-to-write-grpc-custom-prometheus-metric-for-grpc-server-and-client-golang-9617ac3480aa)

### 3. Histogram Bucket Configuration

**Decision**: Use fixed buckets optimized for typical RPC latencies

**Rationale**:

- Fixed buckets: 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s (from spec clarification)
- Covers sub-millisecond to multi-second range typical for cost calculation RPCs
- Aligned with Prometheus default bucket patterns but tuned for RPC use case

**Alternatives Considered**:

- `prometheus.DefBuckets`: Generic, not optimized for RPC latencies
- `prometheus.ExponentialBuckets`: More complex, less readable
- Configurable buckets: Rejected per spec clarification - keeps API simple

**Implementation**:

```go
var defaultBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0}

requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
    Name:    "finfocus_plugin_request_duration_seconds",
    Help:    "Request duration histogram",
    Buckets: defaultBuckets,
}, []string{"grpc_method", "plugin_name"})
```

### 4. Metrics HTTP Server Helper

**Decision**: Provide optional, lightweight helper with configurable port

**Rationale**:

- Per spec clarification: documented as convenience example, not required infrastructure
- Plugin authors may already have HTTP servers for health checks
- Keeps SDK minimal while providing helpful starting point

**Implementation Approach**:

```go
type MetricsServerConfig struct {
    Port     int                  // Default: 9090
    Registry *prometheus.Registry // Uses default if nil
}

func StartMetricsServer(cfg MetricsServerConfig) (*http.Server, error) {
    // Lightweight HTTP server exposing /metrics endpoint
}
```

### 5. Existing SDK Integration Patterns

**Decision**: Follow `TracingUnaryServerInterceptor` pattern from `logging.go`

**Rationale**:

- Consistent API surface with existing SDK interceptors
- Leverages existing `UnaryInterceptors` configuration in `Serve()` function
- Plugin authors already familiar with interceptor chaining pattern

**Integration Point** (from `sdk.go` lines 293-300):

```go
// Build interceptor chain: tracing first, then user interceptors
interceptors := make([]grpc.UnaryServerInterceptor, 0, 1+len(config.UnaryInterceptors))
interceptors = append(interceptors, TracingUnaryServerInterceptor())
interceptors = append(interceptors, config.UnaryInterceptors...)
```

Plugin authors add metrics interceptor to `UnaryInterceptors` slice:

```go
config := pluginsdk.ServeConfig{
    UnaryInterceptors: []grpc.UnaryServerInterceptor{
        pluginsdk.MetricsUnaryServerInterceptor("my-plugin"),
    },
}
```

## Resolved NEEDS CLARIFICATION Items

All technical context items resolved:

| Item | Resolution |
|------|------------|
| Metrics library | prometheus/client_golang |
| Interceptor pattern | Follow go-grpc-prometheus conventions |
| Bucket configuration | Fixed buckets per spec clarification |
| HTTP exposure | Optional helper per spec clarification |

## New Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/prometheus/client_golang | v1.20+ | Metrics instrumentation |

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| Prometheus dependency size | Minimal - well-maintained, widely used |
| Performance overhead | Benchmark tests required per constitution |
| API stability | Use stable v1 Prometheus client API |
