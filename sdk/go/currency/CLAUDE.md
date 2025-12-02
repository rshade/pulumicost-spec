# CLAUDE.md - Currency Package

This file provides guidance for Claude Code when working with the `sdk/go/currency` package.

## Package Overview

The `currency` package provides ISO 4217 currency validation and metadata. It follows the
zero-allocation validation pattern established in `sdk/go/registry/domain.go`.

## Key Files

- `currency.go` - Currency struct, allCurrencies slice, currencyByCode map, GetCurrency(),
  AllCurrencies()
- `validate.go` - IsValid() function using linear scan for zero-allocation validation
- `currency_test.go` - Table-driven unit tests
- `benchmark_test.go` - Performance benchmarks
- `doc.go` - Package documentation

## Design Patterns

### Zero-Allocation Validation

The `IsValid()` function uses a linear scan over a package-level slice to avoid allocations:

```go
var allCurrencies = []Currency{...} // Allocated once at package init

func IsValid(code string) bool {
    for _, c := range allCurrencies {
        if c.Code == code {
            return true
        }
    }
    return false
}
```

### Map Lookup for Metadata

The `GetCurrency()` function uses a map for O(1) lookup:

```go
var currencyByCode = map[string]*Currency{...} // Built at init

func GetCurrency(code string) (*Currency, error) {
    if c, ok := currencyByCode[code]; ok {
        return c, nil
    }
    return nil, ErrCurrencyNotFound
}
```

## Performance Requirements

- `IsValid()`: <15 ns/op, 0 B/op, 0 allocs/op
- `GetCurrency()`: <50 ns/op, 0 B/op, 0 allocs/op
- `AllCurrencies()`: <5 ns/op, 0 B/op, 0 allocs/op

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

1. `allCurrencies` slice in `currency.go`
2. The map is built automatically from `allCurrencies` at init

### Case Sensitivity

Currency codes are case-sensitive per ISO 4217. Always use uppercase (e.g., "USD" not "usd").
