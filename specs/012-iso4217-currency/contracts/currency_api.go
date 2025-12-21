// Package currency provides ISO 4217 currency validation and metadata.
//
// This file defines the API contract for the sdk/go/currency package.
// It is a design document, not compilable code.
package currency

import "errors"

// ErrCurrencyNotFound is returned when a currency code is not found.
var ErrCurrencyNotFound = errors.New("currency not found")

// Currency represents an ISO 4217 currency with complete metadata.
type Currency struct {
	// Code is the 3-letter alphabetic currency code (e.g., "USD", "EUR").
	Code string

	// Name is the official currency name (e.g., "US Dollar", "Euro").
	Name string

	// NumericCode is the 3-digit numeric code (e.g., "840", "978").
	// Stored as string to preserve leading zeros (e.g., "008" for Albanian Lek).
	NumericCode string

	// MinorUnits is the number of decimal places for the currency.
	// Common values: 0 (JPY), 2 (USD, EUR), 3 (KWD).
	MinorUnits int
}

// IsValid checks if code is a valid ISO 4217 currency code.
// Returns true for active ISO 4217 currencies, false otherwise.
//
// Validation rules:
//   - Case-sensitive: "USD" is valid, "usd" is not
//   - No whitespace: " USD" is not valid
//   - Active only: Historic currencies (e.g., "DEM") return false
//
// Performance: <15 ns/op, 0 allocs/op
//
// Example:
//
//	currency.IsValid("USD") // true
//	currency.IsValid("XYZ") // false
//	currency.IsValid("usd") // false (case-sensitive)
func IsValid(code string) bool {
	_ = code // placeholder to satisfy lint until implemented
	// Implementation: Linear scan over allCurrencies slice
	return false // placeholder
}

// GetCurrency retrieves the Currency metadata for a valid code.
// Returns ErrCurrencyNotFound if the code is not a valid ISO 4217 currency.
//
// Example:
//
//	c, err := currency.GetCurrency("USD")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Printf("%s has %d decimal places\n", c.Name, c.MinorUnits)
//	// Output: US Dollar has 2 decimal places
func GetCurrency(code string) (*Currency, error) {
	_ = code // placeholder to satisfy lint until implemented
	// Implementation: Map lookup for O(1) access
	return nil, ErrCurrencyNotFound // placeholder
}

// AllCurrencies returns a slice of all valid ISO 4217 currencies.
// The returned slice is a reference to package-level data and MUST NOT be modified.
//
// The slice contains 180+ active currencies sorted alphabetically by code.
//
// Example:
//
//	for _, c := range currency.AllCurrencies() {
//	    fmt.Printf("%s: %s (%d decimals)\n", c.Code, c.Name, c.MinorUnits)
//	}
func AllCurrencies() []Currency {
	// Implementation: Return reference to allCurrencies package-level slice
	return nil // placeholder
}

// String returns the currency code.
// Implements fmt.Stringer interface.
func (c Currency) String() string {
	return c.Code
}
