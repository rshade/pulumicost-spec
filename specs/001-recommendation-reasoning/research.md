# Research: Standardized Recommendation Reasoning Metadata

**Feature**: Standardized Recommendation Reasoning Metadata
**Date**: 2026-01-13
**Status**: Completed

## Decisions

### 1. Protobuf Enum Naming Convention

- **Decision**: Use `RecommendationReason` for the enum type and `RECOMMENDATION_REASON_VALUE_NAME` for values.
- **Rationale**: Consistent with existing enums in `proto/finfocus/v1/enums.proto`
  (e.g., `FocusServiceCategory`, `FOCUS_SERVICE_CATEGORY_COMPUTE`).
- **Alternatives Considered**: Short names (e.g., `REASON_IDLE`), but this risks collisions
  and violates the established pattern.

### 2. Validation Strategy

- **Decision**: Rely on standard Protobuf validation for enum values (unknown values
  become 0/UNSPECIFIED in proto3).
- **Rationale**: `protoc-gen-validate` is not currently integrated into the build
  pipeline. The standard proto3 behavior is sufficient for transport.
- **Context**: The "Standard Domain Enum Pattern" (string-based) is used for registry domain
  types, but this feature is adding a core Protobuf enum. We should use the generated int32
  enum types directly for zero-allocation performance, which is even better than string iteration.

### 3. SDK Implementation

- **Decision**: Use generated Go code for the enum type.
- **Rationale**: `buf generate` produces high-quality, performant Go code. No custom domain
  wrapper is needed unless we need string representation logic different from the generated `String()` method.

## Open Questions Resolved

- **Naming**: Confirmed standard UpperSnakeCase pattern.
- **Validation**: Confirmed reliance on proto3 defaults + custom logic if needed, avoiding new dependencies.
- **Pattern**: Confirmed that the "Standard Domain Enum Pattern" applies to string-based
  enums, while this feature introduces a Protobuf int32 enum.
