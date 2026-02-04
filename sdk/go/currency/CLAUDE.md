# CLAUDE.md - Currency Package

This file provides guidance for Claude Code when working with the `sdk/go/currency` package.

## Package Overview

The `currency` package provides ISO 4217 currency validation, metadata, and formatting. It follows the
zero-allocation validation pattern established in `sdk/go/registry/domain.go`.

## Key Files

- `currency.go` - Currency struct (with Symbol field), allCurrencies slice, currencyByCode map,
  GetCurrency(), AllCurrencies()
- `symbol.go` - GetSymbol(), FormatAmount(), FormatAmountNoSymbol() helper functions
- `validate.go` - IsValid() function using map lookup for O(1) zero-allocation validation
- `currency_test.go` - Table-driven unit tests
- `symbol_test.go` - Tests for symbol and formatting functions
- `benchmark_test.go` - Performance benchmarks
- `doc.go` - Package documentation

## Currency Struct

```go
type Currency struct {
    Code        string  // 3-letter ISO code (e.g., "USD")
    Name        string  // Full name (e.g., "US Dollar")
    NumericCode string  // 3-digit numeric code (e.g., "840")
    MinorUnits  int     // Decimal places (0-4)
    Symbol      string  // Currency symbol (e.g., "$", "€"), empty if none
}
```

## Design Patterns

### Zero-Allocation Validation

The `IsValid()` function uses the `currencyByCode` map for O(1) lookup with zero allocations:

```go
var currencyByCode = map[string]*Currency{...} // Built at init

func IsValid(code string) bool {
    _, ok := currencyByCode[code]
    return ok
}
```

Note: An earlier implementation used linear scan over `allCurrencies`, but map lookup provides
better performance (~15 ns/op vs ~800 ns/op) while maintaining zero allocations.

### Symbol Lookup with Fallback

The `GetSymbol()` function returns the currency symbol if defined, otherwise falls back to the
currency code:

```go
func GetSymbol(code string) string {
    if c, ok := currencyByCode[code]; ok && c.Symbol != "" {
        return c.Symbol
    }
    return code  // Fallback to code (e.g., "CHF" for Swiss Franc)
}
```

### Amount Formatting

The `FormatAmount()` function formats monetary amounts with:

- Proper decimal precision from MinorUnits (0-4 decimals)
- Thousands separators (commas)
- Currency symbol prefix

```go
currency.FormatAmount(1234.56, "USD")   // "$1,234.56"
currency.FormatAmount(1234.56, "JPY")   // "¥1,235" (0 decimals)
currency.FormatAmount(1234.567, "KWD")  // "د.ك1,234.567" (3 decimals)
```

## Performance Requirements

- `IsValid()`: <15 ns/op, 0 B/op, 0 allocs/op
- `GetCurrency()`: <50 ns/op, 0 B/op, 0 allocs/op
- `AllCurrencies()`: <5 ns/op, 0 B/op, 0 allocs/op
- `GetSymbol()`: <20 ns/op, 0 B/op, 0 allocs/op
- `FormatAmount()`: <500 ns/op, 2-3 allocs/op (acceptable for display)

## Testing Commands

```bash
# Run all tests
go test -v ./sdk/go/currency/

# Run benchmarks
go test -bench=. -benchmem ./sdk/go/currency/

# Run with coverage
go test -cover ./sdk/go/currency/
```

## Common Issues

### Adding New Currencies

When ISO 4217 is updated, add new currencies to:

1. `allCurrencies` slice in `currency.go` (include Symbol field)
2. The map is built automatically from `allCurrencies` at init

### Currency Symbols

- Some currencies have no widely-used symbol (e.g., CHF, XAU, XDR)
- For these, `Symbol` is empty string and `GetSymbol()` returns the code
- Symbols are sourced from Unicode CLDR and common usage

### Case Sensitivity

Currency codes are case-sensitive per ISO 4217. Always use uppercase (e.g., "USD" not "usd").

## Design Decisions

### Why Custom Implementation vs External Libraries

We evaluated `golang.org/x/text/currency` and `github.com/bojanz/currency` but chose
a custom implementation for zero-allocation validation and minimal dependencies.
See [Issue #358](https://github.com/rshade/finfocus-spec/issues/358) for full analysis.
