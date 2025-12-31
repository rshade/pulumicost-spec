# Research: Forecasting Primitives

**Feature**: 030-forecasting-primitives
**Date**: 2025-12-30

## Research Questions

### 1. Field Number Assignment for New Proto Fields

**Decision**: Use field numbers 9-10 for ResourceDescriptor, 3-4 for GetProjectedCostRequest

**Rationale**:

- ResourceDescriptor currently uses fields 1-8 (provider through arn)
- Fields 9 (`growth_type`) and 10 (`growth_rate`) are next available
- GetProjectedCostRequest uses fields 1-2 (resource, utilization_percentage)
- Fields 3 (`growth_type`) and 4 (`growth_rate`) follow sequentially
- Per protobuf best practices, fields 1-15 use single-byte encoding (efficient for
  frequently used fields)

**Alternatives considered**:

- Higher field numbers (>15): Rejected - growth fields are likely to be used frequently
- Shared field numbers across messages: N/A - each message has independent numbering

### 2. GrowthType Enum Placement

**Decision**: Add GrowthType enum to `enums.proto` (existing enum file)

**Rationale**:

- `enums.proto` already contains FOCUS-related enums (FocusServiceCategory,
  FocusChargeCategory, etc.)
- Centralizing enums in one file follows existing codebase pattern
- Avoids circular imports between proto files

**Alternatives considered**:

- Add to `costsource.proto`: Rejected - would make the already large file even longer
- Create new `forecasting.proto`: Rejected - over-engineering for a single enum

### 3. Optional Field Pattern for growth_rate

**Decision**: Use `optional double growth_rate` (proto3 optional)

**Rationale**:

- Proto3 optional allows distinguishing "not set" from explicit 0.0
- Matches existing pattern used for `utilization_percentage` in ResourceDescriptor
- Required for FR-006: detect when LINEAR/EXPONENTIAL is set but rate is missing
- Go generated code will use `*float64` pointer type

**Alternatives considered**:

- Regular `double` field: Rejected - cannot distinguish 0.0 from unset
- Wrapper type (DoubleValue): Rejected - more verbose, optional keyword preferred

### 4. Growth Calculation Formula Precision

**Decision**: Document formulas in proto comments, implement in SDK helper

**Rationale**:

- LINEAR: `cost_at_period_n = base_cost * (1 + rate * n)`
- EXPONENTIAL: `cost_at_period_n = base_cost * (1 + rate)^n`
- Proto defines the contract; SDK provides helper functions
- 0.01% accuracy (SC-002/SC-003) achievable with standard IEEE 754 double precision

**Alternatives considered**:

- Use decimal types: Rejected - proto3 doesn't have native decimal; doubles sufficient
- Store pre-computed projections: Rejected - against stateless design principle

### 5. Override Semantics for Request-Level Fields

**Decision**: Request-level fields fully override resource-level when set

**Rationale**:

- Simpler mental model: "if set, use it; otherwise, use default"
- Matches existing `utilization_percentage` override pattern in GetProjectedCostRequest
- No partial override complexity (e.g., override type but use default rate)
- Both `growth_type` and `growth_rate` on request override both on resource

**Alternatives considered**:

- Partial override (type-only or rate-only): Rejected - adds complexity, unclear semantics
- Merge semantics: Rejected - harder to reason about, potential for subtle bugs

### 6. Validation Placement

**Decision**: Validation in SDK helper functions, not proto layer

**Rationale**:

- Proto definitions are descriptive, not prescriptive (no proto-level validators)
- SDK provides `ValidateGrowthParams(type, rate)` returning error
- Consistent with existing validation pattern in `sdk/go/pricing/validate.go`
- gRPC status codes (InvalidArgument) returned by SDK validation

**Alternatives considered**:

- Proto constraints (buf validate): Not mature enough for complex cross-field rules
- Per-plugin validation: Rejected - inconsistent behavior across implementations

## Existing Codebase Patterns

### Enum Definition Pattern (from enums.proto)

```protobuf
enum GrowthType {
  GROWTH_TYPE_UNSPECIFIED = 0;
  GROWTH_TYPE_NONE = 1;
  GROWTH_TYPE_LINEAR = 2;
  GROWTH_TYPE_EXPONENTIAL = 3;
}
```

### Optional Double Pattern (from costsource.proto:256)

```protobuf
// utilization_percentage is a per-resource utilization override (0.0 to 1.0).
// OPTIONAL. If provided, overrides the global request default.
optional double utilization_percentage = 6;
```

### Message Extension Pattern (ResourceDescriptor)

New fields added at end with next available numbers, comprehensive proto comments
documenting semantics, validation rules, and examples.

## Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| buf | v1.32.1 | Proto compilation, lint, breaking detection |
| google.golang.org/protobuf | latest | Go proto runtime |
| google.golang.org/grpc | latest | gRPC runtime |

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking change detection false positive | Low | Medium | Review buf output carefully |
| Field number collision | Very Low | High | Verify current max field numbers |
| SDK validation edge cases | Medium | Low | Comprehensive test coverage |

## Conclusion

All research questions resolved. No NEEDS CLARIFICATION items remain. Ready for Phase 1
design work.
