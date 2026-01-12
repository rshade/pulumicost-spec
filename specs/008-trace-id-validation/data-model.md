# Data Model: Trace ID Validation

**Feature**: 008-trace-id-validation
**Date**: 2025-11-26
**Purpose**: Define entities, attributes, and validation rules

## Entity Overview

This feature operates on a single value type with no persistent storage.
The trace ID flows through the request context during gRPC call processing.

## Entities

### Trace ID

A distributed tracing correlation identifier following OpenTelemetry W3C format.

| Attribute | Type   | Constraints          | Description                         |
| --------- | ------ | -------------------- | ----------------------------------- |
| value     | string | See validation rules | 32-character hexadecimal identifier |

**Validation Rules**:

| Rule          | Pattern/Value                           | Error Condition                                        |
| ------------- | --------------------------------------- | ------------------------------------------------------ |
| Format        | `^[0-9a-f]{32}$`                        | Non-hex characters or wrong length                     |
| Not empty     | `len > 0`                               | Empty string (for validation; generation handles this) |
| Not all zeros | `!= "00000000000000000000000000000000"` | Reserved invalid value                                 |

**Lifecycle States**:

```text
┌─────────────┐
│  Incoming   │  From gRPC metadata header
└──────┬──────┘
       │
       ▼
┌─────────────┐     Valid?
│  Validate   │──────────────┐
└──────┬──────┘              │
       │ Invalid             │ Yes
       ▼                     │
┌─────────────┐              │
│  Generate   │              │
└──────┬──────┘              │
       │                     │
       ▼                     ▼
┌─────────────────────────────┐
│      Store in Context       │
└─────────────────────────────┘
```

### Request Context

Go context carrying the validated trace ID through the request handler.

| Attribute  | Type       | Description                                |
| ---------- | ---------- | ------------------------------------------ |
| traceIDKey | contextKey | Private key type for context value storage |
| value      | string     | Validated or generated trace ID            |

**Context Key Definition**:

```go
type contextKey string
const traceIDKey contextKey = "finfocus-trace-id"
```

### gRPC Metadata

Incoming request metadata containing trace ID header.

| Attribute | Type     | Description                                                |
| --------- | -------- | ---------------------------------------------------------- |
| key       | string   | `"x-finfocus-trace-id"` (constant: `TraceIDMetadataKey`) |
| values    | []string | Header values (first value used if multiple)               |

## Data Flow

```text
Client Request
     │
     ▼
┌────────────────────────────────────────┐
│         gRPC Metadata Headers          │
│  x-finfocus-trace-id: <trace_id>     │
└────────────────┬───────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────┐
│     TracingUnaryServerInterceptor      │
│  1. Extract from metadata              │
│  2. Validate format                    │
│  3. Generate if invalid/missing        │
│  4. Store in context                   │
└────────────────┬───────────────────────┘
                 │
                 ▼
┌────────────────────────────────────────┐
│           Request Handler              │
│  TraceIDFromContext(ctx) → trace_id    │
└────────────────────────────────────────┘
```

## Validation Decision Table

| Incoming Value        | Action         | Result                |
| --------------------- | -------------- | --------------------- |
| Valid 32 hex chars    | Preserve       | Original trace ID     |
| Empty string `""`     | Generate       | New trace ID          |
| No metadata header    | Generate       | New trace ID          |
| Too short (<32 chars) | Generate       | New trace ID          |
| Too long (>32 chars)  | Generate       | New trace ID          |
| Non-hex characters    | Generate       | New trace ID          |
| All zeros             | Generate       | New trace ID          |
| Control characters    | Generate       | New trace ID          |
| Unicode characters    | Generate       | New trace ID          |
| Multiple values       | Validate first | First value or new ID |

## Generation Specification

New trace IDs are generated using cryptographically secure randomness:

```text
Input:  16 bytes from crypto/rand
Output: 32 lowercase hexadecimal characters

Example: "a1b2c3d4e5f6789012345678abcdef01"
```

**Properties**:

- Collision probability: < 1 in 2^128 (≈ 3.4 × 10^38)
- Generation time: ~150 nanoseconds
- Memory allocation: 48 bytes (16-byte buffer + 32-byte string)

## API Surface

### New Function

```go
// GenerateTraceID creates a cryptographically random trace ID.
// Returns a 32-character lowercase hexadecimal string.
func GenerateTraceID() string
```

### Modified Function

```go
// TracingUnaryServerInterceptor returns a gRPC server interceptor that:
// 1. Extracts trace_id from incoming request metadata
// 2. Validates the trace_id format (32 hex chars, not all zeros)
// 3. Generates a new trace_id if invalid or missing
// 4. Stores the validated/generated trace_id in the request context
//
// The interceptor ensures every request has a valid trace_id for
// observability, preventing log injection attacks from malicious input.
func TracingUnaryServerInterceptor() grpc.UnaryServerInterceptor
```

### Existing Functions (unchanged signature)

```go
// TraceIDFromContext extracts the trace ID from the given context.
func TraceIDFromContext(ctx context.Context) string

// ContextWithTraceID returns a new context with the trace ID stored.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context

// ValidateTraceID validates OpenTelemetry trace identifiers.
// (in pricing package - reused by interceptor)
func ValidateTraceID(traceID string) error
```

## Constants

| Constant             | Value                     | Package   | Description              |
| -------------------- | ------------------------- | --------- | ------------------------ |
| `TraceIDMetadataKey` | `"x-finfocus-trace-id"` | pluginsdk | gRPC metadata header key |
| `FieldTraceID`       | `"trace_id"`              | pluginsdk | Logging field name       |
| `traceIDKey`         | `"finfocus-trace-id"`   | pluginsdk | Context key (private)    |
