# ISO 4217 Currency Package

Package `currency` provides ISO 4217 currency validation and metadata for the PulumiCost SDK.

## Features

- **Zero-allocation validation**: `IsValid()` validates currency codes with 0 allocations
- **Complete metadata**: Access currency name, numeric code, and decimal places
- **180+ currencies**: All active ISO 4217 currencies included
- **High performance**: <15 ns/op for validation operations

## Installation

```go
import "github.com/rshade/pulumicost-spec/sdk/go/currency"
```

## Usage

### Validate a Currency Code

```go
if currency.IsValid("USD") {
    fmt.Println("USD is a valid currency")
}

// Case-sensitive - lowercase is invalid
currency.IsValid("usd") // false
```

### Get Currency Metadata

```go
c, err := currency.GetCurrency("USD")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s (%s) uses %d decimal places\n", c.Name, c.Code, c.MinorUnits)
// Output: US Dollar (USD) uses 2 decimal places

// GetCurrency returns a defensive copy; modifying the returned struct will not
// affect future lookups or validation.
```

### List All Currencies

```go
for _, c := range currency.AllCurrencies() {
    fmt.Printf("%s: %s\n", c.Code, c.Name)
}
```

## Currency Struct

```go
type Currency struct {
    Code        string // 3-letter ISO 4217 code (e.g., "USD")
    Name        string // Official currency name (e.g., "US Dollar")
    NumericCode string // 3-digit numeric code (e.g., "840")
    MinorUnits  int    // Decimal places (e.g., 2 for USD, 0 for JPY)
}
```

## Performance

The package uses zero-allocation validation patterns for high-performance currency validation:

| Operation | Time | Allocations |
|-----------|------|-------------|
| IsValid() | <15 ns/op | 0 B/op, 0 allocs/op |
| GetCurrency() | O(1) lookup with map | Defensive copy (safe to mutate result) |
| AllCurrencies() | <5 ns/op | 0 B/op, 0 allocs/op |

## References

- [ISO 4217 Currency Codes](https://www.iso.org/iso-4217-currency-codes.html)
- [IBAN Currency Codes](https://www.iban.com/currency-codes)
