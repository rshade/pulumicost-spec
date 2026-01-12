# Implementation Plan: Plugin Capability Dry Run Mode

**Branch**: `032-plugin-dry-run` | **Date**: 2025-12-31 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/032-plugin-dry-run/spec.md`

## Summary

Add a dry-run capability to the CostSource gRPC service that allows hosts to query plugins
for their FOCUS field mapping logic without performing actual cost data retrieval. This
enables debugging, validation, and plugin capability comparison. Implementation uses a
hybrid approach: a dedicated `DryRun` RPC for standalone discovery plus an optional
`dry_run` flag on existing GetActualCost/GetProjectedCost RPCs for inline validation.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod) + Protocol Buffers v3
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless RPC introspection, no data persistence)
**Testing**: Go testing + sdk/go/testing harness (bufconn), conformance tests
**Target Platform**: gRPC service (cross-platform, plugin architecture)
**Project Type**: gRPC specification + Go SDK
**Performance Goals**: <100ms response time for dry-run requests (no external calls)
**Constraints**: Backward compatible with existing plugins, transport-layer only (no cost
calculation), must integrate with existing Supports/GetPluginInfo RPCs
**Scale/Scope**: ~50 FOCUS fields to report status for, integrates with existing CostSource
service (12 RPCs currently)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. gRPC Proto Specification-First | PASS | Proto changes to costsource.proto are primary deliverable |
| II. Multi-Provider Consistency | PASS | DryRunRequest uses existing ResourceDescriptor (provider-agnostic) |
| III. Test-First Protocol | PASS | Conformance tests will define expected DryRun RPC behavior first |
| IV. Protobuf Backward Compatibility | PASS | New RPC + optional field additions only, no breaking changes |
| V. Comprehensive Documentation | PASS | Proto comments, SDK docs, examples required before merge |
| VI. Performance as Requirement | PASS | 100ms target in spec, tested at Basic conformance level |
| VII. Validation at Multiple Levels | PASS | buf lint, conformance tests, SDK integration tests planned |

**Gate Result**: PASS - All constitution principles satisfied. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/032-plugin-dry-run/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto snippets)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
├── costsource.proto     # Add DryRun RPC, DryRunRequest/Response messages,
│                        # FieldMapping message, FieldSupportStatus enum,
│                        # dry_run flag to GetActualCostRequest/GetProjectedCostRequest
└── enums.proto          # Add FieldSupportStatus enum (if separate)

sdk/go/
├── proto/               # Generated code (buf generate)
│   └── finfocus/v1/
│       └── costsource.pb.go
├── pluginsdk/           # SDK helpers for dry-run implementation
│   ├── dry_run.go       # DryRunHandler helper, field mapping utilities
│   └── dry_run_test.go  # Unit tests
└── testing/             # Conformance tests
    ├── dry_run_conformance_test.go  # New conformance tests for DryRun RPC
    └── mock_plugin.go               # Update to support dry-run behavior

examples/
└── requests/
    └── dry_run/         # Example DryRun request/response payloads
        ├── aws_ec2.json
        ├── azure_vm.json
        └── README.md
```

**Structure Decision**: Single project structure following existing repository patterns.
Proto changes are the primary deliverable with SDK helpers and conformance tests.

## Complexity Tracking

> No violations requiring justification. All changes follow existing patterns.
