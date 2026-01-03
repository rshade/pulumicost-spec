# Implementation Plan: JSON-LD / Schema.org Serialization

**Branch**: `032-jsonld-serialization` | **Date**: 2025-12-31 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/032-jsonld-serialization/spec.md`

## Summary

High-performance JSON-LD serialization for FOCUS cost data enabling enterprise knowledge graph
indexing. Transforms existing FocusCostRecord and ContractCommitment protobuf messages into
JSON-LD 1.1 format with Schema.org vocabulary mappings and custom FOCUS namespace support.
This is a presentation/transport layer feature - no proto changes required.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: encoding/json (stdlib), crypto/sha256 (stdlib), no external JSON-LD
library needed for serialization-only use case
**Storage**: N/A (stateless serialization library)
**Testing**: go test with table-driven tests, benchmarks for performance validation
**Target Platform**: Go SDK library (cross-platform)
**Project Type**: Single library package within existing SDK
**Performance Goals**: <1ms single record, <5s for 10k records, bounded memory via streaming
**Constraints**: Zero external dependencies preferred, JSON-LD 1.1 output format
**Scale/Scope**: 66-field FocusCostRecord, 12-field ContractCommitment, streaming batch support

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Proto-First | **PASS** | No proto changes - serializing existing messages |
| II. Multi-Provider | **PASS** | JSON-LD output is provider-agnostic |
| III. Test-First | **REQUIRED** | Conformance tests for JSON-LD output structure |
| IV. Backward Compat | **PASS** | New package, no breaking changes |
| V. Documentation | **REQUIRED** | README, examples, godoc comments |
| VI. Performance | **REQUIRED** | Benchmarks for single/batch serialization |
| VII. Validation | **REQUIRED** | JSON-LD output validation tests |

**Gate Result**: PASS - No violations. Test-first, documentation, performance, and validation
requirements will be addressed in implementation.

## Project Structure

### Documentation (this feature)

```text
specs/032-jsonld-serialization/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (JSON-LD context definitions)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/jsonld/                    # NEW PACKAGE
├── serializer.go                 # Core serialization logic
├── serializer_test.go            # Unit tests with table-driven cases
├── serializer_benchmark_test.go  # Performance benchmarks
├── context.go                    # JSON-LD context management
├── context_test.go               # Context configuration tests
├── vocabulary.go                 # FOCUS vocabulary definitions
├── schema_org.go                 # Schema.org mapping definitions
├── id_generator.go               # @id generation (user-provided + fallback)
├── id_generator_test.go          # ID generation tests
├── streaming.go                  # Batch streaming serialization
├── streaming_test.go             # Streaming tests
├── doc.go                        # Package documentation
└── README.md                     # Usage guide and examples

examples/jsonld/                  # Example outputs
├── focus_cost_record.jsonld      # Single record example
├── contract_commitment.jsonld    # ContractCommitment example
├── batch_output.jsonld           # Batch serialization example
└── README.md                     # Example documentation
```

**Structure Decision**: New `sdk/go/jsonld/` package following existing SDK patterns
(separate from pluginsdk which handles gRPC concerns). This maintains separation of
concerns - pluginsdk for gRPC, jsonld for linked data serialization.

## Complexity Tracking

> No Constitution Check violations requiring justification.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| N/A       | N/A        | N/A                                  |
