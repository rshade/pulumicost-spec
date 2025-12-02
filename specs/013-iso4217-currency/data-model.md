# Data Model: ISO 4217 Currency Package

**Feature**: 013-iso4217-currency
**Date**: 2025-11-30

## Entities

### Currency

Represents an ISO 4217 currency with complete metadata.

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| Code | string | 3-letter alphabetic code | Exactly 3 uppercase letters (A-Z) |
| Name | string | Official currency name | Non-empty |
| NumericCode | string | 3-digit numeric code | Exactly 3 digits, leading zeros preserved |
| MinorUnits | int | Decimal places for formatting | 0-4 (typically 0, 2, or 3) |

**Examples**:

| Code | Name | NumericCode | MinorUnits |
|------|------|-------------|------------|
| USD | US Dollar | 840 | 2 |
| EUR | Euro | 978 | 2 |
| JPY | Yen | 392 | 0 |
| KWD | Kuwaiti Dinar | 414 | 3 |
| XXX | No currency | 999 | 0 |
| XTS | Test currency | 963 | 0 |

### Validation Rules

1. **Code Format**: Must be exactly 3 uppercase ASCII letters (A-Z)
2. **Case Sensitivity**: Validation is case-sensitive; "usd" is invalid
3. **Whitespace**: No leading/trailing whitespace allowed; " USD" is invalid
4. **Active Only**: Only active ISO 4217 currencies are valid; historic codes rejected

### Data Storage Pattern

Following the zero-allocation pattern from `sdk/go/registry/domain.go`:

```text
Package Initialization:
┌─────────────────────────────────────────────┐
│ var allCurrencies = []Currency{             │
│   {Code: "AED", Name: "UAE Dirham", ...},   │
│   {Code: "AFN", Name: "Afghani", ...},      │
│   ... (180+ entries)                        │
│ }                                           │
└─────────────────────────────────────────────┘
         │
         ▼ (allocated once at package init)
┌─────────────────────────────────────────────┐
│ Validation: IsValid("USD")                  │
│   └─ Linear scan over allCurrencies slice   │
│   └─ Returns true/false (no allocation)     │
└─────────────────────────────────────────────┘
```

### Memory Layout

| Component | Size | Notes |
|-----------|------|-------|
| Currency struct | ~50 bytes | 4 string headers + int |
| allCurrencies slice | ~9 KB | ~180 currencies × 50 bytes |
| currencyMap (optional) | ~14 KB | For O(1) GetCurrency lookup |

**Total Package Memory**: ~23 KB (acceptable for SDK package)

## Relationships

```text
┌─────────────────────────────────────────────────────────────┐
│                    sdk/go/currency/                         │
│  ┌─────────────┐                                            │
│  │  Currency   │◄──── Exported type                         │
│  │  struct     │                                            │
│  └─────────────┘                                            │
│        │                                                    │
│        ▼                                                    │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ allCurrencies []Currency (package-level, private)   │    │
│  └─────────────────────────────────────────────────────┘    │
│        │                                                    │
│        ├─────────────────┬─────────────────┐                │
│        ▼                 ▼                 ▼                │
│  ┌──────────┐     ┌─────────────┐   ┌──────────────┐        │
│  │ IsValid  │     │ GetCurrency │   │ AllCurrencies│        │
│  │ (code)   │     │ (code)      │   │ ()           │        │
│  └──────────┘     └─────────────┘   └──────────────┘        │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼ (imports)
┌─────────────────────────────────────────────────────────────┐
│              sdk/go/pluginsdk/focus_conformance.go          │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ validateCurrency(code, field) error                   │  │
│  │   └─ Uses currency.IsValid(code)                      │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## State Transitions

N/A - Currency data is immutable and read-only. No state transitions.

## Indexes / Lookup Optimization

For `GetCurrency()` to achieve O(1) lookup (optional optimization):

```text
var currencyByCode = map[string]*Currency{
    "USD": &allCurrencies[X],
    "EUR": &allCurrencies[Y],
    ...
}
```

**Trade-off**: +5KB memory for O(1) lookup vs O(n) linear scan.
**Recommendation**: Implement map for GetCurrency() since metadata retrieval is less
performance-critical than IsValid(), and map provides cleaner API.
