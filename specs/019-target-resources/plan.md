# Implementation Plan: Target Resources for Recommendations

**Branch**: `019-target-resources` | **Date**: 2025-12-17 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/019-target-resources/spec.md`

## Summary

Add `repeated ResourceDescriptor target_resources` field (field 6) to `GetRecommendationsRequest`
proto message, enabling resource-scoped recommendations for Pulumi stacks. Includes SDK validation
(max 100 resources), mock plugin filtering support, and comprehensive conformance tests.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod) + Protocol Buffers v3
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless RPC extension)
**Testing**: go test + bufconn harness + conformance suite
**Target Platform**: gRPC service specification
**Project Type**: Protocol specification repository with Go SDK
**Performance Goals**: Validation < 100ms, no impact on recommendation RPC latency
**Constraints**: Full backward compatibility (empty target_resources = existing behavior)
**Scale/Scope**: Max 100 resources per request, typical stacks 10-50 resources

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | PASS | Proto change before SDK implementation |
| II. Multi-Provider gRPC Consistency | PASS | ResourceDescriptor is provider-agnostic |
| III. Test-First Protocol | PASS | Conformance tests written before proto change |
| IV. Protobuf Backward Compatibility | PASS | New field 6, empty = existing behavior |
| V. Comprehensive Documentation | PASS | Proto comments + SDK README update |
| VI. Performance as gRPC Requirement | PASS | Validation < 100ms, no latency impact |
| VII. Validation at Multiple Levels | PASS | buf + contract + conformance + integration |

## Project Structure

### Documentation (this feature)

```text
specs/019-target-resources/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
proto/pulumicost/v1/
└── costsource.proto          # Add target_resources field to GetRecommendationsRequest

sdk/go/
├── proto/pulumicost/v1/      # Generated code (make generate)
│   └── costsource.pb.go      # Updated with target_resources field
└── testing/
    ├── contract.go           # Add MaxTargetResources + validation
    ├── contract_test.go      # Add validation tests
    ├── mock_plugin.go        # Add target_resources filtering
    └── integration_test.go   # Add target_resources tests
```

**Structure Decision**: Single project structure - this is a protocol specification repository
with generated SDK code. All changes center on proto definitions and SDK testing utilities.

## Complexity Tracking

> No Constitution Check violations. All changes follow established patterns.
