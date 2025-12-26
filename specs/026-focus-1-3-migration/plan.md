# Implementation Plan: FOCUS 1.3 Migration

**Branch**: `026-focus-1-3-migration` | **Date**: 2025-12-23 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/026-focus-1-3-migration/spec.md`

## Summary

Migrate the PulumiCost FOCUS implementation from version 1.2 to 1.3 by:

1. Adding 8 new columns to FocusCostRecord proto message for split cost allocation and
   provider identification
2. Creating a new ContractCommitment proto message (12 fields) for the supplemental dataset
3. Extending the Go SDK builder API with corresponding methods
4. Adding deprecation handling for ProviderName and Publisher fields
5. Maintaining full backward compatibility with existing FOCUS 1.2 implementations

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless proto definitions and SDK)
**Testing**: go test (conformance tests, integration tests, benchmarks)
**Target Platform**: Multi-platform Go library (Linux, macOS, Windows)
**Project Type**: gRPC proto specification + Go SDK library
**Performance Goals**: Builder operations <100ns/op, 0 allocs/op for validation
**Constraints**: Proto field numbers 59+ for backward compatibility, no breaking changes
**Scale/Scope**: 8 new FocusCostRecord columns, 12-field ContractCommitment message

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. gRPC Proto Specification-First | ✅ PASS | Proto changes defined before SDK builder methods |
| II. Multi-Provider gRPC Consistency | ✅ PASS | FOCUS 1.3 is provider-agnostic by design |
| III. Test-First Protocol | ✅ PASS | Conformance tests will be written first (FR-013-016) |
| IV. Protobuf Backward Compatibility | ✅ PASS | Field numbers 59+, no field removals |
| V. Comprehensive Documentation | ✅ PASS | FR-017-019 require documentation updates |
| VI. Performance as gRPC Requirement | ✅ PASS | SC-005 requires <100ns builder operations |
| VII. Validation at Multiple Levels | ✅ PASS | FR-008-010a define validation requirements |

**Gate Status**: ALL PASS - Proceed to Phase 0

## Project Structure

### Documentation (this feature)

```text
specs/026-focus-1-3-migration/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto contract definitions)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
proto/pulumicost/v1/
├── focus.proto          # FocusCostRecord message (add 8 new fields)
├── contract.proto       # NEW: ContractCommitment message (12 fields)
└── enums.proto          # Add new enums if needed

sdk/go/
├── proto/               # Generated code (via buf generate)
│   └── pulumicost/v1/
│       ├── focus.pb.go
│       └── contract.pb.go  # NEW
├── pluginsdk/
│   ├── focus_builder.go          # Add new builder methods
│   ├── focus_conformance.go      # Add FOCUS 1.3 validation
│   ├── contract_builder.go       # NEW: ContractCommitment builder
│   └── contract_conformance.go   # NEW: Contract validation
└── testing/
    ├── focus_conformance_test.go # Add FOCUS 1.3 tests
    ├── contract_test.go          # NEW: Contract commitment tests
    └── benchmark_test.go         # Add new benchmarks

docs/
└── focus-columns.md     # Update with FOCUS 1.3 columns
```

**Structure Decision**: Follows existing repository structure with proto definitions in
`proto/pulumicost/v1/`, generated code in `sdk/go/proto/`, and SDK implementation in
`sdk/go/pluginsdk/`. New ContractCommitment proto may be added to existing focus.proto
or as a separate contract.proto file.

## Complexity Tracking

> No constitutional violations - table not needed.
