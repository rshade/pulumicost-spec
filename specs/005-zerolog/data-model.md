# Data Model: Zerolog SDK Logging Utilities

**Date**: 2025-11-24
**Feature**: 005-zerolog

## Entities

### 1. Logger Configuration

The logger is configured at plugin startup and carries plugin metadata.

**Fields**:

| Field         | Type          | Description       | Validation                           |
| ------------- | ------------- | ----------------- | ------------------------------------ |
| pluginName    | string        | Plugin identifier | Non-empty for meaningful logs        |
| pluginVersion | string        | Semantic version  | Non-empty for meaningful logs        |
| level         | zerolog.Level | Minimum log level | Valid zerolog level (Trace-Disabled) |
| output        | io.Writer     | Log destination   | Defaults to os.Stderr if nil         |

**Relationships**:

- Logger → produces → Log entries
- Logger → includes → Plugin metadata on all entries

### 2. Trace Context

Trace ID propagates through gRPC calls for distributed tracing.

**Fields**:

| Field   | Type   | Description               | Validation                   |
| ------- | ------ | ------------------------- | ---------------------------- |
| traceID | string | Unique request identifier | Any string, empty if not set |

**Relationships**:

- gRPC metadata → extracted by → Interceptor → stored in → Context
- Context → read by → Handler → included in → Log entries

### 3. Field Constants

Standardized field names for ecosystem-wide consistency.

**Constants**:

| Constant           | Value            | Usage                            |
| ------------------ | ---------------- | -------------------------------- |
| FieldTraceID       | "trace_id"       | Distributed tracing correlation  |
| FieldComponent     | "component"      | System component identifier      |
| FieldOperation     | "operation"      | RPC method or operation name     |
| FieldDurationMs    | "duration_ms"    | Operation timing in milliseconds |
| FieldResourceURN   | "resource_urn"   | Pulumi resource identifier       |
| FieldResourceType  | "resource_type"  | Cloud resource type              |
| FieldPluginName    | "plugin_name"    | Plugin identifier                |
| FieldPluginVersion | "plugin_version" | Plugin version                   |
| FieldCostMonthly   | "cost_monthly"   | Monthly cost value               |
| FieldAdapter       | "adapter"        | Cost data source adapter         |
| FieldErrorCode     | "error_code"     | Error classification             |

### 4. Metadata Keys

gRPC metadata keys for trace propagation.

**Constants**:

| Constant           | Value                   | Description               |
| ------------------ | ----------------------- | ------------------------- |
| TraceIDMetadataKey | "x-pulumicost-trace-id" | gRPC metadata header name |

### 5. Context Keys

Internal context keys for value storage.

**Type**: Unexported `contextKey string` type

**Values**:

| Key        | Value                 | Description              |
| ---------- | --------------------- | ------------------------ |
| traceIDKey | "pulumicost-trace-id" | Context key for trace ID |

## State Transitions

### Logger Lifecycle

```text
[Not Created] → NewPluginLogger() → [Active]
                                        ↓
                                    Log calls (Info, Debug, Error, etc.)
                                        ↓
                                    [Active] (loggers are immutable)
```

### Trace ID Lifecycle (per request)

```text
[gRPC Request] → Interceptor extracts from metadata
                      ↓
               [Context with TraceID]
                      ↓
               Handler reads via TraceIDFromContext
                      ↓
               [Included in log entries]
                      ↓
               [Request completes]
```

## Data Flow

### Plugin Startup

```text
Plugin main() → NewPluginLogger(name, version, level, writer)
                    ↓
              Logger with base fields (plugin_name, plugin_version)
                    ↓
              Store in plugin struct for handler use
```

### Request Handling

```text
gRPC Request with x-pulumicost-trace-id header
          ↓
TracingUnaryServerInterceptor extracts trace_id
          ↓
Context enriched with trace_id
          ↓
Handler receives enriched context
          ↓
TraceIDFromContext(ctx) returns trace_id
          ↓
logger.Str(FieldTraceID, traceID).Msg("...")
          ↓
JSON log entry with trace_id field
```

### Operation Timing

```text
done := LogOperation(logger, "GetProjectedCost")
          ↓
Start time captured
          ↓
... perform operation ...
          ↓
defer done() or done() explicitly
          ↓
Duration calculated and logged with operation field
```

## JSON Log Entry Structure

### Standard Entry

```json
{
  "level": "info",
  "time": "2025-11-24T10:30:00Z",
  "plugin_name": "aws-public",
  "plugin_version": "v1.0.0",
  "message": "Request processed"
}
```

### Entry with Trace ID

```json
{
  "level": "info",
  "time": "2025-11-24T10:30:00Z",
  "plugin_name": "aws-public",
  "plugin_version": "v1.0.0",
  "trace_id": "abc123",
  "operation": "GetProjectedCost",
  "message": "Processing cost request"
}
```

### Entry with Operation Timing

```json
{
  "level": "info",
  "time": "2025-11-24T10:30:00Z",
  "plugin_name": "aws-public",
  "plugin_version": "v1.0.0",
  "operation": "GetProjectedCost",
  "duration_ms": 45,
  "message": "operation completed"
}
```

### Error Entry

```json
{
  "level": "error",
  "time": "2025-11-24T10:30:00Z",
  "plugin_name": "aws-public",
  "plugin_version": "v1.0.0",
  "trace_id": "abc123",
  "error": "connection refused",
  "error_code": "CONN_REFUSED",
  "message": "Failed to fetch pricing data"
}
```
