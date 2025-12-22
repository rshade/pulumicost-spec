# Implementation Plan: GetRecommendations RPC

**Branch**: `013-recommendations-rpc` | **Date**: 2025-12-04 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/014-recommendations-rpc/spec.md`

## Summary

Add `GetRecommendations` RPC to `CostSourceService` enabling plugins to surface FinOps
optimization recommendations from various cost management platforms (AWS Cost Explorer,
Kubecost, Azure Advisor, GCP Recommender). Implementation follows proto-first development:
define protobuf messages and service method, generate Go SDK, add PluginSDK support with
optional interface pattern, and comprehensive conformance testing.

## Technical Context

**Language/Version**: Go 1.25.4 (toolchain as specified in go.mod)
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1,
zerolog v1.34.0+, prometheus/client_golang
**Storage**: N/A (stateless RPC, recommendations fetched from backend services)
**Testing**: go test, bufconn (in-memory gRPC), conformance suite (Basic/Standard/Advanced)
**Target Platform**: gRPC service specification (proto definitions + Go SDK)
**Project Type**: gRPC protocol specification repository
**Performance Goals**: <500ms response for <100 items, support 10,000 recommendations
**Constraints**: Backward compatible proto changes, buf breaking check must pass
**Scale/Scope**: 3 enums, 15+ proto messages, 1 new RPC method, PluginSDK interface

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                           | Status | Evidence                                             |
| ----------------------------------- | ------ | ---------------------------------------------------- |
| I. Proto Specification-First        | PASS   | Proto messages defined before SDK implementation     |
| II. Multi-Provider Consistency      | PASS   | Spec includes AWS, Azure, GCP, Kubernetes examples   |
| III. Test-First Protocol            | PASS   | Conformance tests (Basic/Standard/Advanced) required |
| IV. Protobuf Backward Compatibility | PASS   | New RPC method, no breaking changes to existing      |
| V. Comprehensive Documentation      | PASS   | Proto comments, examples, README updates planned     |
| VI. Performance as Requirement      | PASS   | SC-002: <500ms, SC-005: 10,000 items supported       |
| VII. Multi-Layer Validation         | PASS   | buf lint, conformance tests, integration tests       |

**Gate Result**: PASS - All constitutional principles satisfied. No violations requiring
justification.

## Project Structure

### Documentation (this feature)

```text
specs/013-recommendations-rpc/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto excerpts)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
proto/pulumicost/v1/
└── costsource.proto     # Add GetRecommendations RPC + all recommendation messages
                         # (single-file approach per research.md decision)

sdk/go/
├── proto/               # Generated protobuf code (buf generate)
├── pluginsdk/
│   ├── sdk.go           # Add RecommendationsProvider interface
│   ├── helpers.go       # Add recommendation validation helpers
│   └── metrics.go       # Add GetRecommendations metrics
├── testing/
│   ├── harness.go       # Add ValidateRecommendationsResponse
│   ├── mock_plugin.go   # Add GetRecommendations mock implementation
│   ├── conformance.go   # Add recommendation conformance tests
│   └── integration_test.go # Add GetRecommendations integration tests
└── registry/            # May need new enum types for recommendations

examples/
├── recommendations/     # New directory for recommendation examples
│   ├── aws_rightsizing.json
│   ├── kubernetes_request_sizing.json
│   ├── azure_advisor.json
│   └── gcp_recommender.json
└── README.md            # Update with recommendation examples
```

**Structure Decision**: Single project structure. This is a proto specification repository
with Go SDK generation. No frontend/backend split needed.

## Complexity Tracking

> No constitutional violations requiring justification.
