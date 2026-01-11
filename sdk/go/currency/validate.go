package currency

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
	// Map lookup for O(1) validation with zero allocations
	// The currencyByCode map is built once at package init
	_, ok := currencyByCode[code]
	return ok
}
