# Validation Pattern Guidelines: Zero-Allocation Enum Validation

**Pattern Name**: Optimized Slice-Based Enum Validation
**Created**: 2025-11-17
**Status**: Recommended for all SDK packages

## Pattern Overview

This document defines the **zero-allocation enum validation pattern** established in the registry package and recommended
for consistent implementation across the FinFocus SDK.

## Pattern Definition

### Core Pattern Components

#### 1. Package-Level Slice Variable

Define a package-level slice containing all valid enum values:

```go
// allEnumName is a package-level slice containing all valid EnumType values.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allEnumName = []EnumType{
    EnumValue1, EnumValue2, EnumValue3,
}
```

#### 2. Accessor Function

Provide a public accessor function that returns the package-level slice:

```go
// AllEnumName returns a slice of all supported enum values.
func AllEnumName() []EnumType {
    return allEnumName  // Returns reference to package variable
}
```

#### 3. Validation Function

Implement validation using direct iteration over the package-level slice:

```go
// IsValidEnumName checks if the given string represents a valid enum value.
func IsValidEnumName(value string) bool {
    enumValue := EnumType(value)
    for _, valid := range allEnumName {  // Direct reference to package variable
        if enumValue == valid {
            return true
        }
    }
    return false
}
```

### Pattern Benefits

1. **Zero Allocation**: No per-call memory allocation (0 B/op, 0 allocs/op)
2. **High Performance**: 5-12 ns/op for enums with 4-14 values
3. **Cache Friendly**: Small slices fit in CPU L1/L2 cache
4. **Simple**: Minimal code complexity
5. **Maintainable**: Single source of truth for enum values
6. **Type Safe**: Strong typing with Go enum pattern

## When to Use This Pattern

### Recommended For

Use this pattern when:

- ✅ Enum has **< 20 values** (optimal performance range)
- ✅ Enum size is **unlikely to exceed 40-50 values**
- ✅ Validation performance target is **< 100 ns/op**
- ✅ Zero allocation is desired (GC pressure minimization)
- ✅ Memory footprint is a concern (small enums)

### Consider Alternatives When

Evaluate map-based validation if:

- ⚠️ Enum has **> 40 values** (approaching crossover point)
- ⚠️ Enum frequently grows (many new values added regularly)
- ⚠️ Constant-time lookup is required regardless of size
- ⚠️ Memory footprint is not a concern

### Decision Tree

```text
Enum size?
├─ < 20 values  → ✅ Use optimized slice pattern (this pattern)
├─ 20-40 values → Benchmark both approaches
│                  - Likely slice-based still optimal
│                  - Map-based competitive at ~40 values
└─ > 50 values  → Consider map-based validation
                   - Constant ~16 ns/op vs linear scaling
```

## Implementation Guide

### Step 1: Define Enum Type and Constants

```go
package mypackage

// MyEnum represents valid enum values.
type MyEnum string

const (
    MyEnumValue1 MyEnum = "value1"
    MyEnumValue2 MyEnum = "value2"
    MyEnumValue3 MyEnum = "value3"
)
```

### Step 2: Create Package-Level Variable

```go
// allMyEnums is a package-level slice containing all valid MyEnum values.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allMyEnums = []MyEnum{
    MyEnumValue1,
    MyEnumValue2,
    MyEnumValue3,
}
```

**Important Notes**:

- Use `//nolint:gochecknoglobals` directive to suppress linter warnings
- Add documentation explaining the optimization intent
- Initialize with all enum constants for completeness

### Step 3: Implement Accessor Function

```go
// AllMyEnums returns all valid MyEnum values.
func AllMyEnums() []MyEnum {
    return allMyEnums
}
```

### Step 4: Implement Validation Function

```go
// IsValidMyEnum checks if the given string represents a valid MyEnum value.
func IsValidMyEnum(value string) bool {
    myEnum := MyEnum(value)
    for _, valid := range allMyEnums {  // Direct iteration over package variable
        if myEnum == valid {
            return true
        }
    }
    return false
}
```

### Step 5: Add String() Method (Optional)

```go
// String returns the enum value as a string (e.g., "value1").
func (e MyEnum) String() string {
    return string(e)
}
```

### Step 6: Add Benchmark Tests

```go
func BenchmarkIsValidMyEnum(b *testing.B) {
    testCases := []string{"value1", "invalid", "value2", ""}
    b.ResetTimer()
    for i := range b.N {
        _ = IsValidMyEnum(testCases[i%len(testCases)])
    }
}
```

### Step 7: Add Completeness Tests

```go
func TestAllMyEnumsCompleteness(t *testing.T) {
    expected := 3  // Update when adding new values
    if len(AllMyEnums()) != expected {
        t.Errorf("Expected %d enum values, got %d", expected, len(AllMyEnums()))
    }
}
```

## Anti-Patterns to Avoid

### ❌ Anti-Pattern 1: Function-Returned Slices

**Don't do this**:

```go
func AllMyEnums() []MyEnum {
    return []MyEnum{  // ❌ Allocates new slice every call!
        MyEnumValue1, MyEnumValue2, MyEnumValue3,
    }
}
```

**Why it's wrong**:

- Allocates memory on every call
- Adds 20-30 ns overhead
- Creates GC pressure

**Do this instead**:

```go
var allMyEnums = []MyEnum{MyEnumValue1, MyEnumValue2, MyEnumValue3}

func AllMyEnums() []MyEnum {
    return allMyEnums  // ✅ Zero allocation
}
```

### ❌ Anti-Pattern 2: Calling Function in Validation

**Don't do this**:

```go
func IsValidMyEnum(value string) bool {
    myEnum := MyEnum(value)
    for _, valid := range AllMyEnums() {  // ❌ Calls function, allocates slice
        if myEnum == valid {
            return true
        }
    }
    return false
}
```

**Why it's wrong**:

- Calls accessor function on every validation
- Allocates new slice if function returns new slice

**Do this instead**:

```go
func IsValidMyEnum(value string) bool {
    myEnum := MyEnum(value)
    for _, valid := range allMyEnums {  // ✅ Direct slice access
        if myEnum == valid {
            return true
        }
    }
    return false
}
```

### ❌ Anti-Pattern 3: Missing nolint Directive

**Don't do this**:

```go
var allMyEnums = []MyEnum{...}  // ❌ Linter will complain about global variable
```

**Why it's wrong**:

- Linter warnings about global variables
- No explanation for intentional pattern

**Do this instead**:

```go
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allMyEnums = []MyEnum{...}  // ✅ Documented intentional global
```

### ❌ Anti-Pattern 4: Map for Small Enums

**Don't do this** (for < 20 values):

```go
var myEnumMap = map[MyEnum]struct{}{  // ❌ Slower than slice for small enums
    MyEnumValue1: {}, MyEnumValue2: {}, MyEnumValue3: {},
}

func IsValidMyEnum(value string) bool {
    _, ok := myEnumMap[MyEnum(value)]
    return ok
}
```

**Why it's wrong**:

- 2x slower than optimized slice for small enums
- Higher memory footprint (~6x for 8 enums)
- Unnecessary complexity

**Do this instead**:

```go
// Use optimized slice pattern for enums with < 20 values
```

## Package-Specific Applications

### Registry Package (Implemented ✅)

**Status**: Fully implemented as of 2025-11-17

**Enums**: 8 types (Provider, DiscoverySource, PluginStatus, SecurityLevel, InstallationMethod, PluginCapability,
SystemPermission, AuthMethod)

**Performance**:

- 5-12 ns/op across all enum types
- 0 B/op, 0 allocs/op (zero allocation)
- 2x faster than map-based alternatives

**Location**: `sdk/go/registry/domain.go`

### Pricing Package (Recommended)

**Status**: Not yet optimized (uses function-returned slices)

**Current Implementation**:

```go
// Current pattern in pricing package (NOT optimized)
func getAllBillingModes() []BillingMode {
    return []BillingMode{  // ❌ Allocates on every call
        PerHour, PerMinute, /* ... 38+ values ... */
    }
}
```

**Recommended Optimization**:

```go
// Recommended pattern for pricing package
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allBillingModes = []BillingMode{
    PerHour, PerMinute, PerSecond, /* ... all 44 values ... */
}

func getAllBillingModes() []BillingMode {
    return allBillingModes  // ✅ Zero allocation
}

func ValidBillingMode(mode string) bool {
    billingMode := BillingMode(mode)
    for _, valid := range allBillingModes {  // ✅ Direct slice access
        if billingMode == valid {
            return true
        }
    }
    return false
}
```

**Expected Performance** (based on registry package scaling):

- Current: ~40-60 ns/op with allocation overhead
- Optimized: ~20-25 ns/op (44 values × 0.5 ns/value + 5 ns base)
- Memory: 0 allocs/op (vs current 1 alloc/op)

**Recommendation**: Apply same pattern as registry package in future PR

### Future Packages

All future SDK packages should adopt this pattern for enum validation when:

- Enum size < 20 values (optimal)
- Enum size 20-40 values (likely still optimal, benchmark to confirm)
- Zero allocation is desired
- Performance target < 100 ns/op

## Testing Requirements

### Required Tests

1. **Benchmark Tests** - Verify zero allocation and performance:

   ```go
   func BenchmarkIsValidEnum(b *testing.B) {
       // Test with valid and invalid inputs
       // Verify 0 B/op, 0 allocs/op
   }
   ```

2. **Completeness Tests** - Ensure all enum values included:

   ```go
   func TestAllEnumsCompleteness(t *testing.T) {
       // Verify expected count matches actual
   }
   ```

3. **Edge Case Tests** - Test boundary conditions:

   ```go
   func TestValidationEdgeCases(t *testing.T) {
       // Empty strings, case sensitivity, invalid values
   }
   ```

### Performance Targets

- **Allocation**: 0 B/op, 0 allocs/op (strict requirement)
- **Performance**: < 30 ns/op for enums with < 15 values
- **Scalability**: ~0.5 ns per additional enum value (linear)

## Pattern Evolution

### Version History

- **v1.0** (2025-11-17): Initial pattern established in registry package
  - Optimized slice-based validation
  - Zero allocation achieved
  - Benchmark validation included

### Future Considerations

**Potential Enhancements**:

1. **Code Generation**: Generate enum validation code from schema definitions
2. **Hybrid Approach**: Automatic selection between slice/map based on enum size
3. **Compile-Time Validation**: Use Go generics for type-safe enum validation

**When to Revisit**:

- If enum sizes consistently exceed 50 values
- If performance requirements change (e.g., < 10 ns/op target)
- If new Go language features enable better patterns

## References

### Implementation Examples

- **Registry Package**: `sdk/go/registry/domain.go` (reference implementation)
- **Research Analysis**: `specs/001-domain-enum-optimization/research.md`
- **Performance Results**: `specs/001-domain-enum-optimization/performance-results.md`

### Related Documentation

- **Registry CLAUDE.md**: Pattern usage in registry package
- **Pricing CLAUDE.md**: Current implementation and recommendations
- **Quickstart Guide**: `specs/001-domain-enum-optimization/quickstart.md`

## Summary

The **optimized slice-based enum validation pattern** provides:

- ✅ Zero allocation (0 B/op, 0 allocs/op)
- ✅ High performance (5-12 ns/op for 4-14 values)
- ✅ Simple implementation (minimal complexity)
- ✅ Type safety (Go enum pattern)
- ✅ Maintainability (single source of truth)

**Recommendation**: Use this pattern as the standard approach for all enum validation in the FinFocus SDK when
enum size is < 40 values.
