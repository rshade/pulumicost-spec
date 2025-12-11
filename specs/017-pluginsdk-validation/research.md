# Research: PluginSDK Request Validation

**Feature**: `017-pluginsdk-validation`
**Status**: Complete

## Decisions and Findings

### 1. Mapping Helper References

**Unknown**: "The mapping package helpers exist and are documented" (from Spec Assumptions).
**Finding**: Verified `sdk/go/pluginsdk/mapping/aws.go`.

- `ExtractAWSSKU(properties)` exists.
- `ExtractAWSRegion(properties)` exists.
- `ExtractAWSRegionFromAZ(az)` exists.

**Decision**: The error messages will explicitly reference `mapping.ExtractAWSSKU()` and `mapping.ExtractAWSRegion()`
when validation fails for SKU or Region fields respectively. This guides the user to the correct solution.

**Rationale**: Meets FR-005 and US-03 ("Actionable Error Messages").

### 2. Performance & Allocations

**Requirement**: <100ns, zero allocations.
**Research**: Go's string validation (checking `len(s) > 0`) is cheap and zero-allocation. Accessing Proto fields
(e.g., `req.GetResource().GetProvider()`) is also cheap and zero-allocation if just reading pointers/strings.
**Risk**: Error formatting (`fmt.Errorf`) causes allocations.
**Mitigation**:

- The happy path (valid request) MUST NOT allocate.
- The error path (invalid request) is allowed to allocate to create the error message.
- Validation logic should return `nil` immediately upon success.

### 3. Validation Logic Order

**Decision**:

1. Check `req == nil`.
2. Check `req.Resource == nil` (for Projected).
3. Check required string fields (`Provider`, `ResourceType`, `Sku`, `Region`) using `len() > 0`.

**Rationale**: Fail fast. Checking nil pointers first prevents panics. Checking string lengths is cheapest.

### 4. Error Message Format

**Decision**: Simple string errors. "resource.provider is required", "resource.sku is required (use
mapping.ExtractAWSSKU)".
**Rationale**: No need for custom error types (Assumption: "Simple error formatting is acceptable").

## Alternatives Considered

- **Custom Error Types**: defining `ValidationError` struct.
  - _Rejected_: Premature optimization. `status.Error(codes.InvalidArgument, ...)` is the standard gRPC pattern,
    but these are helper functions returning `error`. The caller (Plugin/Core) wraps them in gRPC status if needed.
- **Validation Library (e.g., go-playground/validator)**:
  - _Rejected_: Adds heavy dependency and reflection overhead. We need <100ns performance.
