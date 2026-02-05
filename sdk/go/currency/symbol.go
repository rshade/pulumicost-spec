package currency

import (
	"math"
	"strconv"
	"strings"
)

// thousandsGroupSize is the number of digits in each group when formatting with thousands separators.
const thousandsGroupSize = 3

// notAvailable is the string returned for unformattable float values (NaN, Inf, overflow).
const notAvailable = "N/A"

// defaultDecimalPlaces is the default number of decimal places for unknown currencies.
const defaultDecimalPlaces = 2

// GetSymbol returns the currency symbol for the given uppercase ISO 4217 code.
// The lookup is case-sensitive against currencyByCode; lowercase or mixed-case
// input (e.g., "usd") will not match any entry and will be returned as-is.
// If the code is invalid or has no defined symbol, returns the code itself as fallback.
//
// Performance: <15 ns/op, 0 allocs/op
//
// Examples:
//
//	currency.GetSymbol("USD") // "$"
//	currency.GetSymbol("EUR") // "€"
//	currency.GetSymbol("CHF") // "CHF" (fallback to code, no symbol defined)
//	currency.GetSymbol("XYZ") // "XYZ" (fallback to code, invalid)
func GetSymbol(code string) string {
	if c, ok := currencyByCode[code]; ok && c.Symbol != "" {
		return c.Symbol
	}
	return code
}

// FormatAmount formats a monetary amount with currency symbol and proper decimals.
// Uses MinorUnits from the currency metadata for decimal precision.
// Includes thousands separators (commas).
//
// The symbol is placed before the amount (prefix position).
// Negative amounts are formatted with the minus sign before the symbol (e.g., "-$1,234.56").
//
// Special float values (NaN, +Inf, -Inf) return notAvailable following the graceful degradation
// pattern used elsewhere in the package.
//
// Performance: ~300 ns/op, 2-3 allocs/op (acceptable for display formatting)
//
// Examples:
//
//	currency.FormatAmount(1234.56, "USD")    // "$1,234.56"
//	currency.FormatAmount(-1234.56, "USD")   // "-$1,234.56"
//	currency.FormatAmount(1234.56, "EUR")    // "€1,234.56"
//	currency.FormatAmount(1234.56, "JPY")    // "¥1,235" (0 decimals)
//	currency.FormatAmount(1234.567, "KWD")   // "د.ك1,234.567" (3 decimals)
//	currency.FormatAmount(1234.56, "CHF")    // "CHF1,234.56" (code fallback)
//	currency.FormatAmount(1234.56, "XYZ")    // "XYZ1,234.56" (invalid currency)
//	currency.FormatAmount(math.NaN(), "USD") // "N/A"
func FormatAmount(amount float64, code string) string {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return notAvailable
	}

	decimals := getDecimals(code)
	rounded := roundAmount(amount, decimals)

	// Check for arithmetic overflow after rounding (e.g., math.MaxFloat64 * 100 → +Inf)
	if math.IsNaN(rounded) || math.IsInf(rounded, 0) {
		return notAvailable
	}

	symbol := GetSymbol(code)
	if rounded < 0 {
		formatted := formatWithDecimals(math.Abs(rounded), decimals)
		return "-" + symbol + formatted
	}
	formatted := formatWithDecimals(rounded, decimals)
	return symbol + formatted
}

// FormatAmountNoSymbol formats a monetary amount without a symbol.
// Uses MinorUnits from the currency metadata for decimal precision.
// Includes thousands separators (commas).
//
// Special float values (NaN, +Inf, -Inf) return notAvailable following the graceful degradation
// pattern used elsewhere in the package.
//
// Performance: ~250 ns/op, 2-3 allocs/op
//
// Examples:
//
//	currency.FormatAmountNoSymbol(1234.56, "USD")  // "1,234.56"
//	currency.FormatAmountNoSymbol(1234.56, "JPY")  // "1,235" (0 decimals)
//	currency.FormatAmountNoSymbol(1234.567, "KWD") // "1,234.567" (3 decimals)
//	currency.FormatAmountNoSymbol(1234.56, "XYZ")  // "1,234.56" (default 2 decimals)
//	currency.FormatAmountNoSymbol(math.NaN(), "USD") // "N/A"
func FormatAmountNoSymbol(amount float64, code string) string {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return notAvailable
	}

	decimals := getDecimals(code)
	rounded := roundAmount(amount, decimals)

	// Check for arithmetic overflow after rounding (e.g., math.MaxFloat64 * 100 → +Inf)
	if math.IsNaN(rounded) || math.IsInf(rounded, 0) {
		return notAvailable
	}

	return formatWithDecimals(rounded, decimals)
}

// getDecimals returns the number of decimal places for the given currency code.
func getDecimals(code string) int {
	if c, ok := currencyByCode[code]; ok {
		return c.MinorUnits
	}
	return defaultDecimalPlaces
}

// roundAmount rounds the amount to the given number of decimal places.
// Returns the rounded value, which may be negative zero (-0.0) for small negative inputs.
func roundAmount(amount float64, decimals int) float64 {
	multiplier := math.Pow10(decimals)
	rounded := math.Round(amount*multiplier) / multiplier

	// Normalize negative zero to positive zero
	if rounded == 0 {
		return 0
	}
	return rounded
}

// formatWithDecimals formats a pre-rounded amount with the given decimal places and thousands separators.
func formatWithDecimals(amount float64, decimals int) string {
	formatted := strconv.FormatFloat(amount, 'f', decimals, 64)
	return addThousandsSeparators(formatted)
}

// addThousandsSeparators adds commas as thousands separators to a formatted number string.
func addThousandsSeparators(s string) string {
	// Split into integer and decimal parts
	parts := strings.Split(s, ".")
	intPart := parts[0]

	// Handle negative numbers
	negative := false
	if len(intPart) > 0 && intPart[0] == '-' {
		negative = true
		intPart = intPart[1:]
	}

	// Add commas from right to left
	n := len(intPart)
	if n <= thousandsGroupSize {
		// No separators needed
		if negative {
			intPart = "-" + intPart
		}
		if len(parts) > 1 {
			return intPart + "." + parts[1]
		}
		return intPart
	}

	// Calculate number of commas needed
	numCommas := (n - 1) / thousandsGroupSize

	// Build result with commas (include sign in builder to avoid extra allocation)
	var result strings.Builder
	if negative {
		result.Grow(1 + n + numCommas)
		result.WriteByte('-')
	} else {
		result.Grow(n + numCommas)
	}

	// Add digits with commas
	firstGroupLen := n - (numCommas * thousandsGroupSize)
	result.WriteString(intPart[:firstGroupLen])

	for i := firstGroupLen; i < n; i += thousandsGroupSize {
		result.WriteByte(',')
		result.WriteString(intPart[i : i+thousandsGroupSize])
	}

	s = result.String()

	if len(parts) > 1 {
		return s + "." + parts[1]
	}
	return s
}
