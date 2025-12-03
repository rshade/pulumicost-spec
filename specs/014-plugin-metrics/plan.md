# Implementation Plan: Standardized Plugin Metrics

**Branch**: `014-plugin-metrics` | **Date**: 2025-12-02 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/014-plugin-metrics/spec.md`

## Summary

Implement a standardized Prometheus metrics interceptor for the Go SDK (`sdk/go/pluginsdk`) that
plugin authors can opt into for monitoring request volume and latency. The feature includes:

- A gRPC unary server interceptor (`MetricsUnaryServerInterceptor`) that records request counters
  and latency histograms
- An optional lightweight HTTP server helper for metrics exposure
- Integration with the existing interceptor chaining mechanism

## Technical Context

**Language/Version**: Go 1.25.4 (as per go.mod)
**Primary Dependencies**: google.golang.org/grpc, prometheus/client_golang (new)
**Storage**: N/A (in-memory metrics only)
**Testing**: go test with testify, bufconn for gRPC testing
**Target Platform**: Linux server, cross-platform Go
**Project Type**: SDK library extension
**Performance Goals**: <5% overhead on request processing (SC-004)
**Constraints**: Zero overhead when disabled (FR-004), fixed histogram buckets
**Scale/Scope**: Single package addition to existing pluginsdk

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Proto Specification-First | ✅ PASS | No proto changes required - SDK helper only |
| II. Multi-Provider Consistency | ✅ PASS | Provider-agnostic metrics by design |
| III. Test-First Protocol | ✅ REQUIRED | Must write interceptor tests before implementation |
| IV. Protobuf Backward Compatibility | ✅ PASS | No proto changes |
| V. Comprehensive Documentation | ✅ REQUIRED | README updates, inline docs, example usage |
| VI. Performance as Requirement | ✅ REQUIRED | Benchmark tests for <5% overhead |
| VII. Validation at Multiple Levels | ✅ REQUIRED | Unit tests, integration tests with bufconn |

**Gate Result**: PASS - No constitutional violations. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/014-plugin-metrics/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (Go interfaces)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── metrics.go           # NEW: MetricsUnaryServerInterceptor, MetricsServer helper
├── metrics_test.go      # NEW: Unit and integration tests
├── metrics_benchmark_test.go  # NEW: Performance benchmarks
├── logging.go           # EXISTING: TracingUnaryServerInterceptor (reference pattern)
├── sdk.go               # EXISTING: May need UnaryInterceptors docs update
└── README.md            # UPDATE: Add metrics section
```

**Structure Decision**: Extends existing `sdk/go/pluginsdk` package with new metrics.go file
following the established pattern from logging.go (TracingUnaryServerInterceptor).

## Complexity Tracking

> No constitutional violations requiring justification.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| N/A | N/A | N/A |

## Post-Design Constitution Re-Check

_Verified after Phase 1 design completion._

| Principle | Status | Verification |
|-----------|--------|--------------|
| I. Proto Specification-First | ✅ PASS | Design adds SDK helper only, no proto changes |
| II. Multi-Provider Consistency | ✅ PASS | Metrics are provider-agnostic (method/code/plugin labels) |
| III. Test-First Protocol | ✅ READY | Test contracts defined in data-model.md |
| IV. Protobuf Backward Compatibility | ✅ PASS | No proto changes |
| V. Comprehensive Documentation | ✅ READY | quickstart.md, inline docs in contract |
| VI. Performance as Requirement | ✅ READY | Benchmark requirements defined (<5% overhead) |
| VII. Validation at Multiple Levels | ✅ READY | Unit + integration test strategy defined |

**Post-Design Gate Result**: PASS - Ready for `/speckit.tasks`

## Generated Artifacts

| Artifact | Path | Description |
|----------|------|-------------|
| Research | [research.md](research.md) | Prometheus patterns, gRPC interceptor research |
| Data Model | [data-model.md](data-model.md) | Metric entities, labels, cardinality analysis |
| API Contract | [contracts/metrics.go.design](contracts/metrics.go.design) | Go interface signatures |
| Quickstart | [quickstart.md](quickstart.md) | Usage examples and PromQL queries |
