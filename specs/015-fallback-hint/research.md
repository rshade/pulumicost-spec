# Research: FallbackHint Enum

**Feature**: `001-fallback-hint`
**Date**: 2025-12-06

## 1. Unknowns & Clarifications

### naming Convention

- **Question**: What naming convention should be used for the FallbackHint enum values?
- **Decision**: Use `FALLBACK_HINT_` prefix (e.g., `FALLBACK_HINT_UNSPECIFIED`).
- **Rationale**: Standard Protobuf practice to avoid scoping collisions in generated code
  (especially C++/Go). Confirmed in spec clarifications.

### Zero Cost vs No Data

- **Question**: How to distinguish "no data found" vs "zero cost"?
- **Decision**:
  - Empty array `[]` = `FALLBACK_HINT_RECOMMENDED` (No data, try elsewhere).
  - Array with `[{cost: 0.00}]` = `FALLBACK_HINT_NONE` (Data exists, it's just free).
- **Rationale**: Preserves the semantic difference between "I don't know" and "I know it costs
  nothing".

### SDK Implementation Pattern

- **Question**: How to expose this in the Go SDK?
- **Decision**: Use Functional Options Pattern
  (e.g., `NewActualCostResponse(..., WithFallbackHint(hint))`).
- **Rationale**: idiomatic Go, extensible, and keeps the constructor signature clean.

## 2. Technology Choices

### Protobuf Enum

- **Choice**: standard `enum` in `costsource.proto`.
- **Values**:
  0. `FALLBACK_HINT_UNSPECIFIED` (Default)
  1. `FALLBACK_HINT_NONE`
  2. `FALLBACK_HINT_RECOMMENDED`
  3. `FALLBACK_HINT_REQUIRED`

### Integration Strategy

- **Approach**: Add field `fallback_hint` to `GetActualCostResponse`.
- **Compatibility**: Field number 2 (next available is likely 2, need to check
  `GetActualCostResponse` definition).
  - `GetActualCostResponse` has `repeated ActualCostResult results = 1;`. Next is 2.

## 3. Best Practices

- **Proto3 Defaults**: The default value (0) must be the "safe" or "backward compatible"
  behavior. `UNSPECIFIED` mapping to "No Fallback" ensures existing plugins (which send 0
  implicitly) don't trigger unwanted fallbacks.
