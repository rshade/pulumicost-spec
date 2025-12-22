# Research: ISO 4217 Currency Validation Package

**Date**: 2025-11-30
**Feature**: 013-iso4217-currency

## Research Areas

### 1. ISO 4217 Currency Standard

**Decision**: Use the complete ISO 4217 active currency list with all 180+ currencies.

**Rationale**:

- ISO 4217 is the international standard for currency codes maintained by ISO
- The standard defines 3-letter alphabetic codes (e.g., USD, EUR, GBP) and 3-digit numeric codes
- Each currency has associated "minor units" (decimal places) for formatting
- The standard includes supranational currencies (XDR), test codes (XTS), and no-currency (XXX)

**Alternatives Considered**:

- Subset of common currencies only - Rejected: Would not support all cloud provider billing currencies
- Dynamic lookup via external API - Rejected: Adds latency and external dependency

**Source**: [ISO 4217 Currency Codes](https://www.iso.org/iso-4217-currency-codes.html),
[IBAN Currency Codes](https://www.iban.com/currency-codes)

### 2. Zero-Allocation Validation Pattern

**Decision**: Follow the established pattern from `sdk/go/registry/domain.go`.

**Rationale**:

- Registry package achieves 5-12 ns/op with 0 allocations
- Pattern uses package-level slice variables allocated once at initialization
- Simple linear search is faster than map lookup for small-to-medium datasets (<200 items)
- Consistent with existing SDK patterns

**Pattern Details**:

```go
// Package-level slice allocated once
var allCurrencies = []Currency{...}

// Validation iterates over pre-allocated slice
func IsValid(code string) bool {
    for _, c := range allCurrencies {
        if c.Code == code {
            return true
        }
    }
    return false
}
```

**Alternatives Considered**:

- Map-based lookup (`map[string]bool`) - Rejected: Higher memory overhead, not faster for ~180 items
- Sorted slice with binary search - Rejected: More complex, marginal benefit for dataset size

### 3. Currency Metadata Structure

**Decision**: Store complete metadata for each currency (code, name, numeric code, minor units).

**Rationale**:

- Enables proper monetary value formatting (e.g., JPY uses 0 decimal places)
- Numeric codes required for some financial protocols
- Full names useful for UI/documentation generation
- Minimal memory overhead (~50 bytes per currency × 180 = ~9KB total)

**Structure**:

```go
type Currency struct {
    Code       string // 3-letter alphabetic code (e.g., "USD")
    Name       string // Full name (e.g., "US Dollar")
    NumericCode string // 3-digit numeric code (e.g., "840")
    MinorUnits int    // Decimal places (e.g., 2 for USD, 0 for JPY)
}
```

**Alternatives Considered**:

- Code-only validation (no metadata) - Rejected: Limits usefulness, metadata needed for P2 user story
- Separate lookup maps for each field - Rejected: More complex, higher memory

### 4. Historic vs Active Currencies

**Decision**: Include only active ISO 4217 currencies; exclude historic/withdrawn codes.

**Rationale**:

- Cloud billing systems use active currencies
- Historic codes (DEM, FRF, etc.) add complexity without value
- Reduces dataset size and validation time
- Can add historic support later if needed

**Alternatives Considered**:

- Include all historic currencies - Rejected: Not needed for cloud billing use case
- Configurable historic mode - Rejected: Over-engineering for current requirements

### 5. Case Sensitivity

**Decision**: Strict uppercase validation only (case-sensitive).

**Rationale**:

- ISO 4217 standard specifies uppercase codes
- Existing `focus_conformance.go` uses uppercase
- Case-insensitive would require string conversion (allocation)
- Clear, predictable behavior

**Alternatives Considered**:

- Case-insensitive validation - Rejected: Would require allocation for ToUpper()
- Normalize input - Rejected: Changes input semantics, potential for bugs

### 6. Package Location

**Decision**: Create `sdk/go/currency/` package.

**Rationale**:

- Follows existing SDK package structure (`sdk/go/registry/`, `sdk/go/pricing/`)
- Clear separation of concerns
- Reusable by other SDK packages and external consumers
- Consistent naming convention

**Alternatives Considered**:

- Add to existing `sdk/go/pricing/` - Rejected: Currency is distinct domain concept
- Add to `sdk/go/registry/` - Rejected: Registry is for plugin management, not currency

## Implementation Dependencies

| Dependency | Type | Notes |
|------------|------|-------|
| `sdk/go/registry/domain.go` | Reference | Pattern to follow for zero-allocation validation |
| `sdk/go/pluginsdk/focus_conformance.go` | Migration | Update to use new currency package |
| ISO 4217 specification | Data source | Currency list and metadata |

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Missing currencies | Low | Medium | Verify against official ISO 4217 list |
| Performance regression | Low | High | Benchmark against registry package baseline |
| Breaking existing tests | Medium | High | Run full test suite before/after migration |
| Incorrect minor units | Low | Medium | Cross-reference multiple sources |

## Conclusion

All research items resolved. Ready to proceed with Phase 1 design:

- Data model: Currency struct with 4 fields
- API: IsValid(), GetCurrency(), AllCurrencies()
- Pattern: Zero-allocation following registry package
- Migration: Update focus_conformance.go to import new package
