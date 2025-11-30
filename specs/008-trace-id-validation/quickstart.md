# Quickstart: Trace ID Validation

**Feature**: 008-trace-id-validation
**Date**: 2025-11-26
**Purpose**: Usage guide for plugin developers

## Overview

The `TracingUnaryServerInterceptor` now validates incoming trace IDs and
automatically generates valid replacements for malformed or missing values.
This ensures all requests have valid trace IDs for distributed tracing
while preventing log injection attacks.

## What Changed

### Before (v0.4.x and earlier)

```go
// Trace IDs passed through without validation
// Malicious input could corrupt logs
ctx = ContextWithTraceID(ctx, values[0])  // Any string accepted
```

### After (this release)

```go
// Trace IDs are validated against format rules
// Invalid values are replaced with generated ones
traceID = ValidateAndEnsureTraceID(incoming)  // Always valid output
ctx = ContextWithTraceID(ctx, traceID)
```

## Usage

### No Code Changes Required

If you're already using `TracingUnaryServerInterceptor()`, validation is
automatically enabled when you upgrade the SDK:

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"

// Your existing code continues to work unchanged
server := grpc.NewServer(
    grpc.UnaryInterceptor(pluginsdk.TracingUnaryServerInterceptor()),
)
```

### Accessing the Trace ID

Retrieve the validated trace ID in your handler:

```go
func (s *myPlugin) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    // Always returns a valid trace ID (never empty, never malformed)
    traceID := pluginsdk.TraceIDFromContext(ctx)

    // Safe to use in logs without sanitization
    log.Info().Str("trace_id", traceID).Msg("processing request")

    // ... handler logic
}
```

### Generating Trace IDs (for clients)

If you need to generate a trace ID on the client side:

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"

// Generate a valid trace ID
traceID := pluginsdk.GenerateTraceID()
// Result: "a1b2c3d4e5f6789012345678abcdef01" (32 hex chars)

// Include in gRPC metadata
md := metadata.New(map[string]string{
    pluginsdk.TraceIDMetadataKey: traceID,
})
ctx := metadata.NewOutgoingContext(ctx, md)
```

## Validation Rules

The following trace ID formats are accepted:

| Format             | Example                            | Status   |
| ------------------ | ---------------------------------- | -------- |
| Valid 32 hex chars | `abcdef1234567890abcdef1234567890` | Accepted |
| Valid with numbers | `12345678901234567890123456789012` | Accepted |

The following formats are rejected (replacement generated):

| Format        | Example                            | Reason                   |
| ------------- | ---------------------------------- | ------------------------ |
| Empty         | `""`                               | Missing trace ID         |
| Too short     | `abcdef`                           | Less than 32 characters  |
| Too long      | `abcdef...01` (33+ chars)          | More than 32 characters  |
| Non-hex       | `ghijklmnop...`                    | Invalid characters (g-z) |
| All zeros     | `00000000000000000000000000000000` | Reserved invalid value   |
| Control chars | `abc\ndef...`                      | Security risk            |
| Unicode       | `ąbcdef...`                        | Invalid encoding         |

## Security Benefits

### Log Injection Prevention

Without validation, attackers could inject malicious trace IDs:

```text
# Malicious input
x-pulumicost-trace-id: fake\nERROR: System compromised\ntrace_id=

# Would appear in logs as:
INFO trace_id=fake
ERROR: System compromised
trace_id= processing request
```

With validation, the malicious input is replaced:

```text
# Same malicious input
x-pulumicost-trace-id: fake\nERROR: System compromised\ntrace_id=

# Now appears in logs as:
INFO trace_id=a1b2c3d4e5f6789012345678abcdef01 processing request
```

### Buffer Overflow Prevention

Large trace IDs (e.g., 10KB) are rejected and replaced, preventing
memory exhaustion or buffer issues in downstream systems.

## Backward Compatibility

| Scenario          | Behavior                              |
| ----------------- | ------------------------------------- |
| Valid trace IDs   | Preserved unchanged                   |
| Existing plugins  | Continue working without code changes |
| Missing trace IDs | Now automatically generated           |
| Invalid trace IDs | Now replaced (was passed through)     |

The only observable change is that invalid trace IDs are no longer
propagated to handlers. This is a security improvement that does not
affect correct usage.

## Performance

| Operation                               | Time    | Memory |
| --------------------------------------- | ------- | ------ |
| Validation (valid input)                | ~50 ns  | 0 B    |
| Validation (invalid input) + generation | ~200 ns | 48 B   |

Both scenarios are well under the 1ms latency budget.

## Troubleshooting

### Trace ID changed unexpectedly

If you notice trace IDs changing between client and server:

1. **Check client format**: Ensure client sends 32 lowercase hex characters
2. **Check for all-zeros**: All-zero trace IDs are invalid per OpenTelemetry spec
3. **Check for typos**: Uppercase letters (A-F) are invalid (use a-f)

### How to verify validation is working

```go
// Test with invalid input
md := metadata.New(map[string]string{
    pluginsdk.TraceIDMetadataKey: "invalid!",
})
ctx := metadata.NewOutgoingContext(context.Background(), md)

// Handler receives a valid generated trace ID
traceID := pluginsdk.TraceIDFromContext(ctx)
// traceID will be a valid 32-char hex string, not "invalid!"
```

## Related Documentation

- [Spec: 008-trace-id-validation](./spec.md)
- [Data Model](./data-model.md)
- [Research](./research.md)
- [OpenTelemetry Trace Context](https://www.w3.org/TR/trace-context/)
