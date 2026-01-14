# Implementation Plan - Standardized Recommendation Reasoning Metadata

**Feature**: Standardized Recommendation Reasoning Metadata
**Status**: Draft
**Spec**: [specs/001-recommendation-reasoning/spec.md](spec.md)

## Technical Context

### Architecture
- **Protocol**: gRPC (Protobuf v3)
- **Component**: `finfocus-spec` (Core Protocol Definition & Go SDK)
- **Touchpoints**:
  - `proto/finfocus/v1/recommendation.proto`: New `Reason` enum and message fields.
  - `sdk/go/finfocus/v1/recommendation.pb.go`: Generated Go code.
  - `sdk/go/finfocus/v1/recommendation_custom.go`: Helper methods for reason codes (if needed).

### Dependencies
- **Internal**: `finfocus-spec` core protos.
- **External**: `google.golang.org/protobuf`, `google.golang.org/grpc`.

### Existing Patterns
- **Enums**: Follows standard Protobuf enum patterns (e.g., `RECOMMENDATION_REASON_UNSPECIFIED`).
- **Validation**: Needs to integrate with existing validation layers (e.g., `protoc-gen-validate` or custom SDK validation).
- **Naming**: Upper snake case for enum values (e.g., `RECOMMENDATION_REASON_IDLE`).

### Integration Points
- **Upstream**: Plugins (AWS, GCP, etc.) will populate this field.
- **Downstream**: Dashboards/CLIs will consume this field for display/filtering.

## Constitution Check

| Principle | Check | Context |
|-----------|-------|---------|
| **I. Proto First** | [x] | Changes starting in `proto/finfocus/v1/` (modeled in contracts). |
| **II. Multi-Provider** | [x] | Categories (Idle, Over-provisioned) are generic and apply to all clouds. |
| **III. Spec Consumes** | [x] | Feature transports data; does not calculate it. |
| **IV. Separation** | [x] | Logic remains in plugins; spec only defines the schema. |
| **V. Test-First** | [ ] | Conformance tests will be written in implementation phase. |
| **VI. Backward Compat** | [x] | Adding fields is non-breaking. |
| **VII. Docs & Identity** | [x] | New enum documented in `contracts/`. |
| **VIII. Performance** | [x] | Enum validation relies on zero-allocation generated code. |
| **IX. Observability** | [x] | N/A (Standard fields). |
| **X. Patterns** | [x] | Follows standard domain enum pattern (Protobuf variant). |
| **XI. Copyright** | [ ] | Header required on new files (implementation phase). |

## Unknowns & Risks (NEEDS CLARIFICATION)
- **Research 1**: Resolved. Standard Protobuf naming applies.
- **Research 2**: Resolved. No `protoc-gen-validate` dependency.
- **Research 3**: Resolved. Standard Domain Enum Pattern is for string-based domains; Protobuf enums use generated code.

## Plan Phases

### Phase 1: Research & Validation
- [x] Research Protobuf enum conventions.
- [x] Research validation rules for repeated fields.
- [x] Research "Standard Domain Enum Pattern" implementation details.
- [x] Produce `research.md`.

### Phase 2: Specification & Protocol
- [x] Define `RecommendationReason` enum in `contracts/`.
- [x] Add `primary_reason` and `secondary_reasons` to `Recommendation` message in `contracts/`.
- [ ] Run `buf lint` and `buf generate` (Implementation).
- [x] Produce `contracts/` artifacts (generated Go code preview).

### Phase 3: SDK Implementation
- [ ] Implement Go SDK helpers/validation for new reason codes.
- [ ] Ensure zero-allocation validation if applicable.
- [x] Update `data-model.md`.

### Phase 4: Documentation & Examples
- [x] Update `quickstart.md` with usage examples.
- [x] Update agent context.
