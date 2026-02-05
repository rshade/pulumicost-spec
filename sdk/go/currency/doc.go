// Package currency provides ISO 4217 currency validation, metadata, and formatting.
//
// The currency package implements zero-allocation validation for ISO 4217 currency codes,
// following the established pattern from sdk/go/registry/domain.go. It provides fast,
// memory-efficient validation, metadata retrieval, and amount formatting for all active
// ISO 4217 currencies.
//
// # Features
//
//   - Zero-allocation validation: IsValid() achieves <15 ns/op with 0 allocations
//   - Complete ISO 4217 coverage: 180+ active currencies with full metadata
//   - Currency symbols: GetSymbol() returns the currency symbol (e.g., "$", "€", "£")
//   - Amount formatting: FormatAmount() formats monetary amounts with proper decimals
//   - Case-sensitive validation: Only uppercase codes are accepted
//   - Metadata retrieval: GetCurrency() provides name, numeric code, decimal places, and symbol
//   - List operations: AllCurrencies() returns all valid currencies
//
// # Performance Characteristics
//
//   - IsValid(): <15 ns/op, 0 B/op, 0 allocs/op
//   - GetCurrency(): O(1) lookup with a defensive copy for safe mutation
//   - AllCurrencies(): <5 ns/op, 0 B/op, 0 allocs/op
//   - GetSymbol(): <15 ns/op, 0 B/op, 0 allocs/op
//   - FormatAmount(): ~300 ns/op, 2-3 allocs/op (acceptable for display)
//
// # Usage Examples
//
// Basic validation:
//
//	if currency.IsValid("USD") {
//	    fmt.Println("USD is valid")
//	}
//
// Metadata retrieval:
//
//	usd, err := currency.GetCurrency("USD")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("%s (%s) uses %d decimal places\n", usd.Name, usd.Symbol, usd.MinorUnits)
//
// Currency symbols:
//
//	symbol := currency.GetSymbol("USD") // "$"
//	symbol = currency.GetSymbol("CHF")  // "CHF" (fallback when no symbol defined)
//
// Amount formatting:
//
//	formatted := currency.FormatAmount(1234.56, "USD")  // "$1,234.56"
//	formatted = currency.FormatAmount(1234.56, "JPY")   // "¥1,235" (0 decimals)
//	formatted = currency.FormatAmount(1234.567, "KWD")  // "د.ك1,234.567" (3 decimals)
//
// # Formatting Conventions
//
// The formatting functions use US English conventions:
//   - Thousands separator: comma (e.g., "1,234,567")
//   - Decimal separator: period (e.g., "1,234.56")
//   - Symbol placement: prefix position (e.g., "$1,234.56")
//   - Negative format: minus sign before symbol (e.g., "-$1,234.56")
//   - Invalid floats (NaN, Inf): returns "N/A" for graceful degradation
//
// For locale-specific formatting (e.g., European "1.234,56"), consider using
// golang.org/x/text/message or a dedicated i18n library.
//
// Listing all currencies:
//
//	for _, c := range currency.AllCurrencies() {
//	    fmt.Printf("%s: %s (%s)\n", c.Code, c.Name, c.Symbol)
//	}
//
// # Integration
//
// This package is used by sdk/go/pluginsdk/focus_conformance.go for validating
// billing and pricing currency fields in FOCUS cost records. Plugins can use
// GetSymbol() and FormatAmount() for consistent currency display across the
// FinFocus ecosystem.
//
// # References
//
//   - ISO 4217: https://www.iso.org/iso-4217-currency-codes.html
//   - FOCUS 1.2 Specification: https://focus.finops.org
//   - Unicode CLDR Currency Data: https://github.com/unicode-org/cldr
package currency
