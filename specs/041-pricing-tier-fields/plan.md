# Implementation Plan: Add Pricing Tier & Probability Fields

**Branch**: `041-pricing-tier-fields` | **Date**: 2026-01-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/041-pricing-tier-fields/spec.md`

## Summary

Add `pricing_category` and `spot_interruption_risk_score` fields to `EstimateCostResponse` and
`GetProjectedCostResponse` to provide better UX feedback regarding cost drivers and reliability
risks.

## Technical Context

**Language/Version**: Go 1.25.5+, TypeScript 5.0+  
**Primary Dependencies**: Protobuf 3, gRPC  
**Testing**: Go unit tests, Buf linting  
**Target Platform**: Multi-cloud (AWS, Azure, GCP)  
**Project Type**: SDK/Spec  
**Constitution Check**: Passed (Contract-first, Spec consumes, FinFocus aligned)

## Project Structure

### Documentation

```text
specs/041-pricing-tier-fields/
├── plan.md              # This file
├── research.md          # Research on existing enums and messages
├── data-model.md        # Updated response messages
└── spec.md              # Original feature specification
```

### Source Code Updates

```text
proto/finfocus/v1/
└── costsource.proto     # Update response messages
```

## Implementation Phases

### Phase 1: Protobuf Design (Contracts)

- Update `proto/finfocus/v1/costsource.proto` with the new fields.
- Field numbers: `EstimateCostResponse` (3, 4), `GetProjectedCostResponse` (8, 9).

### Phase 2: SDK Generation

- Run `make generate` to update Go and TypeScript client libraries.
- Verify generated code includes the new fields.

### Phase 3: Verification

- Add a Go unit test to `sdk/go/` (locate existing cost estimation tests) to verify the fields can be set and retrieved.
- Run `make test` and `make buf-lint`.
