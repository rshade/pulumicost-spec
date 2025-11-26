# Quickstart: Structured Logging Example for EstimateCost

**Feature**: 007-zerolog-logging-example
**Date**: 2025-11-26

## Prerequisites

Before implementing this example:

1. **005-zerolog utilities must be implemented** in `sdk/go/logging/`
2. **006-estimate-cost RPC** is already available in proto and MockPlugin
3. Go 1.24+ with zerolog v1.34.0+ dependency

## Quick Implementation Guide

### Step 1: Add Test Function to integration_test.go

Add to `sdk/go/testing/integration_test.go`:

```go
// TestStructuredLoggingExample demonstrates structured logging patterns
// for the EstimateCost RPC, per NFR-001 of spec 006-estimate-cost.
func TestStructuredLoggingExample(t *testing.T) {
    // Setup - see implementation details in tasks
}
```

### Step 2: Demonstrate Request Logging

```go
// Log incoming request with context fields
logger.Info().
    Str(logging.FieldTraceID, traceID).
    Str(logging.FieldOperation, "EstimateCost").
    Str(logging.FieldResourceType, req.GetResourceType()).
    Int("attribute_count", len(req.GetAttributes().GetFields())).
    Msg("Processing cost estimation request")
```

### Step 3: Demonstrate Response Logging

```go
// Log successful response with cost details
done := logging.LogOperation(logger, "EstimateCost")
// ... perform RPC ...
done() // Logs duration_ms automatically

logger.Info().
    Str(logging.FieldTraceID, traceID).
    Str(logging.FieldOperation, "EstimateCost").
    Float64(logging.FieldCostMonthly, resp.GetCostMonthly()).
    Str("currency", resp.GetCurrency()).
    Msg("Cost estimation completed")
```

### Step 4: Demonstrate Error Logging

```go
// Log errors with error code and context
logger.Error().
    Err(err).
    Str(logging.FieldTraceID, traceID).
    Str(logging.FieldOperation, "EstimateCost").
    Str(logging.FieldResourceType, req.GetResourceType()).
    Str(logging.FieldErrorCode, status.Code(err).String()).
    Msg("Cost estimation failed")
```

## Running the Example

```bash
# Run the example test
go test -v ./sdk/go/testing/ -run TestStructuredLoggingExample

# Expected output shows structured JSON logs
```

## Expected Log Output

### Request Log

```json
{"level":"info","time":"...","trace_id":"abc123","operation":"EstimateCost",
 "resource_type":"aws:ec2/instance:Instance","attribute_count":2,
 "message":"Processing cost estimation request"}
```

### Success Log

```json
{"level":"info","time":"...","trace_id":"abc123","operation":"EstimateCost",
 "cost_monthly":8.76,"currency":"USD","duration_ms":45,
 "message":"Cost estimation completed"}
```

### Error Log

```json
{"level":"error","time":"...","trace_id":"abc123","operation":"EstimateCost",
 "resource_type":"invalid:resource","error_code":"INVALID_ARGUMENT",
 "error":"invalid resource_type format","duration_ms":12,
 "message":"Cost estimation failed"}
```

## Key Patterns to Follow

1. **Always include trace_id** when available (graceful degradation if missing)
2. **Use standard field constants** from `sdk/go/logging/`
3. **Never log attribute values** - log count only to prevent credential exposure
4. **Use LogOperation helper** for automatic timing measurement
5. **Include operation name** in every log entry for filterability
6. **Log at appropriate levels**: Info for normal flow, Error for failures

## File Changes Summary

| File | Change |
|------|--------|
| `sdk/go/testing/integration_test.go` | Add `TestStructuredLoggingExample` function |

## Next Steps

After implementation:

1. Run `go test -v ./sdk/go/testing/ -run TestStructuredLoggingExample`
2. Verify JSON output contains all required fields
3. Run `make lint` to ensure code quality
4. Update `sdk/go/testing/README.md` if needed
