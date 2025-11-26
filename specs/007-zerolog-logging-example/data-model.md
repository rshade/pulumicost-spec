# Data Model: Structured Logging Example for EstimateCost

**Feature**: 007-zerolog-logging-example
**Date**: 2025-11-26

## Overview

This feature is a documentation/example feature. The data model consists of:

1. **Existing entities** from dependencies (used, not created)
2. **Log entry structure** (output format, not persisted data)

## Existing Entities (Dependencies)

### From 005-zerolog

| Entity | Purpose | Key Fields |
|--------|---------|------------|
| `Logger` | Configured zerolog instance | plugin_name, plugin_version, level |
| `TraceID` | Correlation identifier | string value from context |

### From 006-estimate-cost

| Entity | Purpose | Key Fields |
|--------|---------|------------|
| `EstimateCostRequest` | RPC input | resource_type (string), attributes (Struct) |
| `EstimateCostResponse` | RPC output | currency (string), cost_monthly (double) |

### From sdk/go/testing

| Entity | Purpose | Key Fields |
|--------|---------|------------|
| `MockPlugin` | Test plugin implementation | ShouldErrorOnEstimateCost, EstimateCostDelay |
| `TestHarness` | In-memory gRPC harness | client, server |

## Log Entry Structure

The example demonstrates structured JSON log entries with these field patterns:

### Request Log Entry

```json
{
  "level": "info",
  "time": "2025-11-26T10:30:00Z",
  "trace_id": "abc123",
  "operation": "EstimateCost",
  "resource_type": "aws:ec2/instance:Instance",
  "attribute_count": 3,
  "message": "Processing cost estimation request"
}
```

### Success Response Log Entry

```json
{
  "level": "info",
  "time": "2025-11-26T10:30:00Z",
  "trace_id": "abc123",
  "operation": "EstimateCost",
  "resource_type": "aws:ec2/instance:Instance",
  "cost_monthly": 8.76,
  "currency": "USD",
  "duration_ms": 45,
  "message": "Cost estimation completed"
}
```

### Error Log Entry

```json
{
  "level": "error",
  "time": "2025-11-26T10:30:00Z",
  "trace_id": "abc123",
  "operation": "EstimateCost",
  "resource_type": "invalid:resource",
  "error_code": "INVALID_ARGUMENT",
  "error": "Invalid resource_type format",
  "duration_ms": 12,
  "message": "Cost estimation failed"
}
```

## Field Constants Mapping

| Constant | JSON Field | Type | Description |
|----------|------------|------|-------------|
| `FieldTraceID` | `trace_id` | string | Correlation ID for distributed tracing |
| `FieldOperation` | `operation` | string | RPC method name ("EstimateCost") |
| `FieldResourceType` | `resource_type` | string | Pulumi resource type |
| `FieldCostMonthly` | `cost_monthly` | float64 | Estimated monthly cost |
| `FieldDurationMs` | `duration_ms` | int64 | Operation duration in milliseconds |
| `FieldErrorCode` | `error_code` | string | gRPC/custom error code |
| `FieldPluginName` | `plugin_name` | string | Plugin identifier |
| `FieldPluginVersion` | `plugin_version` | string | Plugin version |

## Relationships

```text
TestHarness (1) ----> (1) MockPlugin
     |
     v
CostSourceServiceClient ----> EstimateCost RPC
     |
     v
Logger ----> Log Entries (request, response, error)
     |
     +----> TraceID (from context)
```

## State Transitions

Not applicable - this is a stateless example demonstrating logging patterns.

## Validation Rules

| Field | Rule | Enforced By |
|-------|------|-------------|
| trace_id | May be empty (graceful degradation) | Application logic |
| resource_type | Must be non-empty for logging | Proto validation |
| cost_monthly | Must be non-negative | Proto field semantics |
| duration_ms | Auto-calculated by LogOperation | LogOperation helper |

## Data Volume Considerations

Not applicable - example produces minimal log output for demonstration purposes only.
