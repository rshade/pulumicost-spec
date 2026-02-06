# Implementation Plan: GetActualCost Pagination Support

**Branch**: `044-actual-cost-pagination` | **Date**: 2026-02-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/044-actual-cost-pagination/spec.md`

## Summary

Add pagination support to the `GetActualCost` RPC to enable retrieval of large cost
datasets (10,000+ records) in manageable pages. The implementation follows the established
`GetRecommendations` pagination pattern: offset-based tokens with base64 encoding,
`DefaultPageSize` of 50, `MaxPageSize` of 1000. Proto-first approach adds `page_size`,
`page_token`, `next_page_token`, and `total_count` fields. SDK provides a
`PaginateActualCosts()` helper for plugins and an `ActualCostIterator` for hosts.
TypeScript SDK gets a parallel `actualCostIterator()` async generator. Backward
compatibility is preserved: existing plugins and hosts continue to work without changes.

## Technical Context

**Language/Version**: Go 1.25.6 (per go.mod) + Protocol Buffers v3, TypeScript (SDK)
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1,
zerolog (logging), connectrpc (TypeScript)
**Storage**: N/A (stateless pagination with offset-based tokens)
**Testing**: `go test` (unit + integration), conformance suite in `sdk/go/testing/`,
benchmarks with `-bench -benchmem`
**Target Platform**: Cross-platform gRPC library (Linux, macOS, Windows)
**Project Type**: Single project (gRPC specification + multi-language SDK)
**Performance Goals**: <100ms p99 for pagination helper operations; 10,000 records
retrievable within 10 seconds across all pages; zero-allocation token encoding/decoding
**Constraints**: <100MB memory per maximum page (1,000 records); responses must stay
under 4MB gRPC message size limit; backward compatible with existing plugins/hosts
**Scale/Scope**: Handles datasets of 10,000+ records; pages of up to 1,000 records;
affects proto, Go SDK (pluginsdk, testing), and TypeScript SDK (client)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First** (I): Proto definitions updated first (`costsource.proto`
  fields 7-8 on request, 4-5 on response), then SDK code follows.
- [x] **Spec Consumes** (III): Pagination is data transport mechanics, not pricing
  logic. No calculators or tiered pricing embedded.
- [x] **Multi-Provider** (II): Pagination is provider-agnostic. No provider-specific
  pagination logic. Works identically for AWS, Azure, GCP, Kubernetes plugins.
- [x] **FinFocus Alignment** (VII): All naming uses FinFocus conventions. No legacy
  PulumiCost references introduced.
- [x] **SDK Synchronization** (XIII): Plan includes both Go SDK helpers and TypeScript
  SDK `actualCostIterator()`. Proto changes regenerate both SDKs via buf.
- [x] **Documentation Integrity** (XIV): Plan includes README updates for pluginsdk,
  godoc comments on all exported functions, and quickstart guide.
- [x] **Test-First Protocol** (V): Conformance tests written first to define expected
  pagination behavior, then implementation to make them pass.
- [x] **Backward Compatibility** (VI): New fields are additive only. `page_size = 0`
  and empty `page_token` preserve existing behavior. buf breaking check passes.
- [x] **Performance** (VIII): Zero-allocation token encoding, benchmarks required for
  `PaginateActualCosts()` and iterator.
- [x] **Follow Established Patterns** (X): Replicates exact `GetRecommendations`
  pagination pattern for consistency.

## Project Structure

### Documentation (this feature)

```text
specs/044-actual-cost-pagination/
├── plan.md              # This file
├── research.md          # Phase 0 output (completed)
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto diff)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
└── costsource.proto                    # Add pagination fields to GetActualCost messages

sdk/go/pluginsdk/
├── helpers.go                          # Add PaginateActualCosts() function
├── helpers_test.go                     # Add pagination unit tests
├── logging.go                          # Add FieldResultCount constant
├── actual_cost_iterator.go             # New: ActualCostIterator (Next/Record/Err)
├── actual_cost_iterator_test.go        # New: Iterator unit tests
└── README.md                           # Update with pagination documentation

sdk/go/proto/finfocus/v1/
└── costsource.pb.go                    # Regenerated from proto (make generate)

sdk/go/testing/
├── mock_plugin.go                      # Add pagination to mock GetActualCost
├── pagination_conformance_test.go      # New: Pagination conformance tests
├── integration_test.go                 # Add paginated GetActualCost integration tests
└── benchmark_test.go                   # Add pagination benchmarks

sdk/typescript/packages/client/src/
├── utils/pagination.ts                 # Add actualCostIterator() async generator
└── utils/pagination.test.ts            # Add iterator tests
```

**Structure Decision**: This feature extends the existing repository structure. No new
directories are needed. The proto-first approach means `costsource.proto` changes drive
all generated code. New Go files are limited to the iterator (`actual_cost_iterator.go`)
and conformance tests (`pagination_conformance_test.go`). All other changes are additions
to existing files, following the established patterns.

## Complexity Tracking

> No constitution violations. All gates pass without justification needed.
