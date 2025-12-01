# Quickstart: ISO 4217 Currency Package

**Package**: `github.com/rshade/pulumicost-spec/sdk/go/currency`

## Installation

```go
import "github.com/rshade/pulumicost-spec/sdk/go/currency"
```

## Basic Usage

### Validate a Currency Code

```go
package main

import (
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/currency"
)

func main() {
    // Valid currency
    if currency.IsValid("USD") {
        fmt.Println("USD is valid")
    }

    // Invalid currency
    if !currency.IsValid("XYZ") {
        fmt.Println("XYZ is not valid")
    }

    // Case-sensitive (lowercase is invalid)
    if !currency.IsValid("usd") {
        fmt.Println("usd is not valid (must be uppercase)")
    }
}
```

### Get Currency Metadata

```go
package main

import (
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/currency"
)

func main() {
    // Get full currency information
    usd, err := currency.GetCurrency("USD")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Code: %s\n", usd.Code)           // USD
    fmt.Printf("Name: %s\n", usd.Name)           // US Dollar
    fmt.Printf("Numeric: %s\n", usd.NumericCode) // 840
    fmt.Printf("Decimals: %d\n", usd.MinorUnits) // 2

    // Handle currencies with different decimal places
    jpy, _ := currency.GetCurrency("JPY")
    fmt.Printf("%s uses %d decimal places\n", jpy.Name, jpy.MinorUnits) // 0

    kwd, _ := currency.GetCurrency("KWD")
    fmt.Printf("%s uses %d decimal places\n", kwd.Name, kwd.MinorUnits) // 3
}
```

### List All Currencies

```go
package main

import (
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/currency"
)

func main() {
    // Get all valid currencies
    all := currency.AllCurrencies()
    fmt.Printf("Total currencies: %d\n", len(all))

    // Iterate and display
    for _, c := range all {
        fmt.Printf("%s - %s (numeric: %s, decimals: %d)\n",
            c.Code, c.Name, c.NumericCode, c.MinorUnits)
    }
}
```

## Integration with FOCUS Conformance

The currency package is used by `sdk/go/pluginsdk/focus_conformance.go` for
validating billing and pricing currency fields in FOCUS cost records.

```go
package pluginsdk

import (
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/currency"
)

// validateCurrency checks if a currency code is a valid ISO 4217 code.
func validateCurrency(code string, fieldName string) error {
    if !currency.IsValid(code) {
        return fmt.Errorf("%s must be a valid ISO 4217 currency code, got %q",
            fieldName, code)
    }
    return nil
}
```

## Formatting Monetary Values

Use `MinorUnits` to properly format currency amounts:

```go
package main

import (
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/currency"
)

func formatAmount(amount float64, currencyCode string) string {
    c, err := currency.GetCurrency(currencyCode)
    if err != nil {
        return fmt.Sprintf("%.2f %s", amount, currencyCode)
    }

    format := fmt.Sprintf("%%.%df %%s", c.MinorUnits)
    return fmt.Sprintf(format, amount, c.Code)
}

func main() {
    fmt.Println(formatAmount(1234.567, "USD")) // 1234.57 USD
    fmt.Println(formatAmount(1234.567, "JPY")) // 1235 JPY
    fmt.Println(formatAmount(1234.567, "KWD")) // 1234.567 KWD
}
```

## Error Handling

```go
package main

import (
    "errors"
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/currency"
)

func main() {
    // Check for not found error
    _, err := currency.GetCurrency("INVALID")
    if errors.Is(err, currency.ErrCurrencyNotFound) {
        fmt.Println("Currency not found")
    }
}
```

## Performance Characteristics

| Operation | Time | Allocations |
|-----------|------|-------------|
| IsValid() | <15 ns/op | 0 B/op, 0 allocs/op |
| GetCurrency() | <50 ns/op | 0 B/op, 0 allocs/op |
| AllCurrencies() | <5 ns/op | 0 B/op, 0 allocs/op |

The package uses zero-allocation validation patterns for high-performance
currency validation in hot paths.
