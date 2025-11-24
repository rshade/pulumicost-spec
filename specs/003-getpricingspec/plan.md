# Implementation Plan: GetPricingSpec RPC Enhancement

**Branch**: `001-getpricingspec` | **Date**: 2025-11-22 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-getpricingspec/spec.md`

## Summary

Enhance the existing GetPricingSpec RPC to support transparent pricing breakdowns with
assumptions and tiered pricing. The RPC method already exists in the proto; this feature
adds missing fields to PricingSpec message: `unit`, `assumptions`, and `pricing_tiers`
with a new PricingTier message type.

## Technical Context

**Language/Version**: Go 1.24+ with protobuf
**Primary Dependencies**: buf v1.32.1, google.golang.org/protobuf, google.golang.org/grpc
**Storage**: N/A (stateless RPC)
**Testing**: go test, conformance tests via sdk/go/testing harness
**Target Platform**: gRPC server/client (cross-platform)
**Project Type**: Single (protobuf specification repository)
**Performance Goals**: <100ms response time, 0 allocations for validation
**Constraints**: Backward compatible with existing PricingSpec usage
**Scale/Scope**: Plugin developers implementing CostSourceService

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Proto Specification-First | ✅ PASS | Adding fields to existing proto messages |
| II. Multi-Provider Consistency | ✅ PASS | No provider-specific fields added |
| III. Test-First Protocol | ⏳ PENDING | Conformance tests to be written in Phase 1 |
| IV. Backward Compatibility | ✅ PASS | Adding new optional fields only; no removals |
| V. Comprehensive Documentation | ⏳ PENDING | Proto comments to be added |
| VI. Performance Requirements | ✅ PASS | Zero-allocation validation pattern exists |
| VII. Multi-layer Validation | ✅ PASS | buf lint, JSON schema, conformance tests |

## Project Structure

### Documentation (this feature)

```text
specs/003-getpricingspec/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
proto/pulumicost/v1/
└── costsource.proto     # Add PricingTier message, update PricingSpec

sdk/go/
├── proto/               # Generated code (via buf generate)
├── pricing/             # Domain types - add PricingTier, Unit enums
│   ├── domain.go        # BillingMode constants
│   └── validate.go      # JSON schema validation
└── testing/             # Test framework
    ├── harness.go       # In-memory gRPC harness
    └── mock_plugin.go   # Mock with GetPricingSpec support

schemas/
└── pricing_spec.schema.json  # Update with assumptions, pricing_tiers

examples/specs/
└── *.json               # Update examples with new fields
```

**Structure Decision**: Single project structure - this is a protobuf specification repository with generated SDK code.

## Complexity Tracking

> No violations requiring justification. All gates pass or pending test-first workflow.
