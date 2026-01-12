# Research: Zerolog SDK Logging Utilities

**Date**: 2025-11-24
**Feature**: 005-zerolog

## Technology Decisions

### 1. Zerolog Library Selection

**Decision**: Use zerolog v1.34.0+ as the logging library

**Rationale**:

- Zero-allocation JSON logging (critical for high-throughput plugin workloads)
- Fluent builder API matches Go idioms
- Wide adoption in production Go systems
- Active maintenance and security updates
- Native support for structured JSON output

**Alternatives Considered**:

- **zap**: Slightly higher allocations, more complex configuration
- **logrus**: Higher memory overhead, slower performance
- **slog (stdlib)**: Go 1.21+, less mature ecosystem, fewer features

### 2. Context Key Pattern

**Decision**: Use custom type for context key to avoid collisions

**Rationale**:

- Go best practice to use unexported custom types for context keys
- Prevents accidental key collisions with other packages
- Type-safe retrieval from context

**Implementation**:

```go
type contextKey string
const traceIDKey contextKey = "finfocus-trace-id"
```

### 3. gRPC Interceptor Pattern

**Decision**: Implement UnaryServerInterceptor using grpc/metadata package

**Rationale**:

- Standard gRPC pattern for cross-cutting concerns
- Metadata extraction is well-documented
- Supports both client and server-side propagation
- Works with existing FinFocus gRPC infrastructure

**Pattern**:

```go
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{},
        info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        if md, ok := metadata.FromIncomingContext(ctx); ok {
            if values := md.Get(TraceIDMetadataKey); len(values) > 0 {
                ctx = ContextWithTraceID(ctx, values[0])
            }
        }
        return handler(ctx, req)
    }
}
```

### 4. Logger Output Configuration

**Decision**: Default to os.Stderr with io.Writer flexibility

**Rationale**:

- stderr is Unix standard for logs (stdout for data)
- io.Writer interface allows testing with buffers
- File output via --logfile/--logdir flags for production
- External tools handle rotation (logrotate, systemd)

**Implementation**:

```go
func NewPluginLogger(name, version string, level zerolog.Level, w io.Writer) zerolog.Logger
// w defaults to os.Stderr when nil
```

### 5. Field Constants Design

**Decision**: Export string constants for all standard field names

**Rationale**:

- Compile-time safety for field names
- IDE autocomplete support
- Prevents typos in distributed team
- Easy to add new fields without breaking existing code

**Constants**:

```go
const (
    FieldTraceID       = "trace_id"
    FieldComponent     = "component"
    FieldOperation     = "operation"
    FieldDurationMs    = "duration_ms"
    FieldResourceURN   = "resource_urn"
    FieldResourceType  = "resource_type"
    FieldPluginName    = "plugin_name"
    FieldPluginVersion = "plugin_version"
    FieldCostMonthly   = "cost_monthly"
    FieldAdapter       = "adapter"
    FieldErrorCode     = "error_code"
)
```

### 6. LogOperation Timing Pattern

**Decision**: Return a closure that logs duration when called

**Rationale**:

- Defer-friendly pattern for automatic timing
- No need for explicit start/stop calls
- Captures operation context at call time
- Idiomatic Go pattern (similar to trace spans)

**Implementation**:

```go
func LogOperation(logger zerolog.Logger, operation string) func() {
    start := time.Now()
    return func() {
        logger.Info().
            Str(FieldOperation, operation).
            Int64(FieldDurationMs, time.Since(start).Milliseconds()).
            Msg("operation completed")
    }
}
```

### 7. Package Location

**Decision**: Create new `sdk/go/pluginsdk` package

**Rationale**:

- Clear separation from existing packages (pricing, proto, testing, registry)
- Name indicates purpose: utilities for plugin developers
- Follows existing SDK package structure
- Room for future plugin utilities (metrics, config, etc.)

## Best Practices Research

### Zerolog Best Practices

1. **Logger construction**: Create once at startup, pass by value (loggers are
   immutable)
2. **Contextual logging**: Use `With()` to add persistent fields
3. **Level filtering**: Set at construction, use `zerolog.GlobalLevel()` sparingly
4. **Error handling**: Use `.Err(err)` method for error fields
5. **Sampling**: Available but not needed for SDK (plugin decision)

### gRPC Interceptor Best Practices

1. **Metadata key naming**: Use lowercase with hyphens (gRPC normalizes to
   lowercase)
2. **Context enrichment**: Always pass enriched context to handler
3. **Error handling**: Don't swallow errors, propagate status codes
4. **Performance**: Interceptors should be lightweight (<1μs overhead)

### Structured Logging Best Practices

1. **Field naming**: snake_case for JSON compatibility
2. **Required fields**: timestamp, level, message always present
3. **Contextual fields**: trace_id, component, operation for correlation
4. **Numeric fields**: Use appropriate types (int64 for ms, float64 for costs)

## Testing Strategy

### Unit Tests

- Logger construction with various configurations
- Context key get/set operations
- Field constant values
- LogOperation timing accuracy

### Integration Tests (bufconn)

- Interceptor extracts trace_id from metadata
- Interceptor handles missing metadata gracefully
- Concurrent requests maintain isolation
- Full plugin example with all utilities

### Benchmarks

- Logger construction performance
- Log call overhead (Info, Debug, Error)
- Interceptor overhead per request
- Memory allocation tracking

## Dependencies to Add

```go
require (
    github.com/rs/zerolog v1.34.0
)
```

No indirect dependencies beyond existing gRPC packages.

## Risk Assessment

| Risk                         | Likelihood | Impact | Mitigation                             |
| ---------------------------- | ---------- | ------ | -------------------------------------- |
| zerolog API breaking changes | Low        | Medium | Pin to v1.34.x, monitor releases       |
| Performance regression       | Low        | High   | Benchmarks in CI, compare to baseline  |
| Field name collisions        | Low        | Low    | Use constants, document reserved names |
| Context key collisions       | Low        | Low    | Use custom unexported type             |
