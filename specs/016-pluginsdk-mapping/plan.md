# Implementation Plan: PluginSDK Mapping Package

**Branch**: `016-pluginsdk-mapping` | **Date**: 2025-12-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/016-pluginsdk-mapping/spec.md`

## Summary

Create a `pluginsdk/mapping/` subpackage providing shared helper functions for extracting SKU,
region, and other pricing-relevant fields from Pulumi resource properties. The package will
support AWS, Azure, and GCP cloud providers with provider-specific extraction logic, plus
generic fallback extractors for custom use cases. This enables pulumicost-core to delegate
cloud-specific property mapping to the canonical spec repository, keeping the core cloud-agnostic.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: None (stdlib only - no external dependencies required)
**Storage**: N/A (stateless helper functions)
**Testing**: Go testing package with table-driven tests, benchmarks
**Target Platform**: Cross-platform Go library
**Project Type**: Single library package
**Performance Goals**: <50 ns/op for extraction functions, 0 allocs/op
**Constraints**: Must not panic on nil/empty input; consistent with existing pluginsdk patterns
**Scale/Scope**: 10 public functions across 4 files (aws.go, azure.go, gcp.go, common.go)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Proto-First Development | N/A | No proto changes required; helper package |
| II. Multi-Provider Consistency | PASS | AWS, Azure, GCP all supported equally |
| III. Test-First Protocol | PASS | Unit tests with table-driven approach planned |
| IV. Backward Compatibility | PASS | New package; no breaking changes |
| V. Comprehensive Documentation | PASS | Package docs, README, examples planned |
| VI. Performance as Requirement | PASS | Benchmark targets defined |
| VII. Validation at Multiple Levels | PASS | Unit tests + integration with existing harness |

**Gate Status**: PASS - No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/016-pluginsdk-mapping/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (Go interfaces)
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── mapping/
│   ├── aws.go           # ExtractAWSSKU, ExtractAWSRegion, ExtractAWSRegionFromAZ
│   ├── azure.go         # ExtractAzureSKU, ExtractAzureRegion
│   ├── gcp.go           # ExtractGCPSKU, ExtractGCPRegion, GCP regions list
│   ├── common.go        # ExtractSKU, ExtractRegion (generic)
│   ├── doc.go           # Package documentation
│   ├── mapping_test.go  # Comprehensive unit tests
│   └── benchmark_test.go # Performance benchmarks
├── env.go               # Existing environment helpers
├── sdk.go               # Existing SDK entry point
└── README.md            # Updated with mapping package docs
```

**Structure Decision**: Single subpackage under existing `pluginsdk/` directory following
established SDK patterns. No separate test directory needed as Go tests co-locate with source.

## Complexity Tracking

No complexity violations requiring justification.
