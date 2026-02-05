package currency_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/currency"
)

// TestGetSymbol_VariousCases tests GetSymbol with various currency codes.
func TestGetSymbol_VariousCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		// Major currencies with symbols
		{name: "USD", code: "USD", expected: "$"},
		{name: "EUR", code: "EUR", expected: "€"},
		{name: "GBP", code: "GBP", expected: "£"},
		{name: "JPY", code: "JPY", expected: "¥"},
		{name: "CNY", code: "CNY", expected: "¥"},
		{name: "INR", code: "INR", expected: "₹"},
		{name: "KRW", code: "KRW", expected: "₩"},
		{name: "THB", code: "THB", expected: "฿"},
		{name: "RUB", code: "RUB", expected: "₽"},
		{name: "UAH", code: "UAH", expected: "₴"},
		{name: "ILS", code: "ILS", expected: "₪"},
		{name: "PHP", code: "PHP", expected: "₱"},
		{name: "TRY", code: "TRY", expected: "₺"},
		{name: "PLN", code: "PLN", expected: "zł"},
		{name: "BRL", code: "BRL", expected: "R$"},
		{name: "VND", code: "VND", expected: "₫"},

		// Dollar currencies with prefixes
		{name: "CAD", code: "CAD", expected: "C$"},
		{name: "AUD", code: "AUD", expected: "A$"},
		{name: "NZD", code: "NZD", expected: "NZ$"},
		{name: "HKD", code: "HKD", expected: "HK$"},
		{name: "SGD", code: "SGD", expected: "S$"},
		{name: "TWD", code: "TWD", expected: "NT$"},

		// Currencies without symbols (fallback to code)
		{name: "CHF fallback", code: "CHF", expected: "CHF"},
		{name: "XAU fallback", code: "XAU", expected: "XAU"},
		{name: "XDR fallback", code: "XDR", expected: "XDR"},

		// Invalid codes (fallback to code)
		{name: "Invalid XYZ", code: "XYZ", expected: "XYZ"},
		{name: "Invalid ABC", code: "ABC", expected: "ABC"},
		{name: "Empty", code: "", expected: ""},
		{name: "Lowercase", code: "usd", expected: "usd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := currency.GetSymbol(tt.code)
			if got != tt.expected {
				t.Errorf("GetSymbol(%q) = %q, want %q", tt.code, got, tt.expected)
			}
		})
	}
}

// TestFormatAmount_VariousCases tests FormatAmount with various amounts and currencies.
func TestFormatAmount_VariousCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   float64
		code     string
		expected string
	}{
		// Basic formatting with USD
		{name: "USD simple", amount: 1234.56, code: "USD", expected: "$1,234.56"},
		{name: "USD zero", amount: 0, code: "USD", expected: "$0.00"},
		{name: "USD small", amount: 0.01, code: "USD", expected: "$0.01"},
		{name: "USD round up", amount: 1234.567, code: "USD", expected: "$1,234.57"},
		{name: "USD round down", amount: 1234.564, code: "USD", expected: "$1,234.56"},
		{name: "USD large", amount: 1234567.89, code: "USD", expected: "$1,234,567.89"},
		{name: "USD negative", amount: -1234.56, code: "USD", expected: "-$1,234.56"},

		// EUR with Euro symbol
		{name: "EUR simple", amount: 1234.56, code: "EUR", expected: "€1,234.56"},

		// JPY with 0 decimals
		{name: "JPY simple", amount: 1234.56, code: "JPY", expected: "¥1,235"},
		{name: "JPY large", amount: 123456789, code: "JPY", expected: "¥123,456,789"},
		{name: "JPY round", amount: 1234.4, code: "JPY", expected: "¥1,234"},

		// KWD with 3 decimals
		{name: "KWD simple", amount: 1234.567, code: "KWD", expected: "د.ك1,234.567"},
		{name: "KWD round", amount: 1234.5678, code: "KWD", expected: "د.ك1,234.568"},

		// CLF with 4 decimals
		{name: "CLF simple", amount: 1234.5678, code: "CLF", expected: "UF1,234.5678"},
		{name: "CLF round", amount: 1234.56789, code: "CLF", expected: "UF1,234.5679"},

		// CHF (no symbol, uses code)
		{name: "CHF fallback", amount: 1234.56, code: "CHF", expected: "CHF1,234.56"},

		// Invalid currency (uses code, default 2 decimals)
		{name: "Invalid currency", amount: 1234.56, code: "XYZ", expected: "XYZ1,234.56"},

		// Edge cases
		{name: "Very small", amount: 0.001, code: "USD", expected: "$0.00"},
		{name: "No commas needed", amount: 999.99, code: "USD", expected: "$999.99"},
		{name: "Single comma", amount: 1000, code: "USD", expected: "$1,000.00"},

		// Floating-point edge cases
		{name: "NaN", amount: math.NaN(), code: "USD", expected: "N/A"},
		{name: "Positive Inf", amount: math.Inf(1), code: "USD", expected: "N/A"},
		{name: "Negative Inf", amount: math.Inf(-1), code: "USD", expected: "N/A"},

		// Additional edge cases
		{name: "Very large negative", amount: -1234567890.12, code: "USD", expected: "-$1,234,567,890.12"},
		{name: "Negative JPY", amount: -123456, code: "JPY", expected: "-¥123,456"},
		{name: "Tiny negative rounds to zero", amount: -0.001, code: "USD", expected: "$0.00"},
		{name: "Tiny positive rounds to zero", amount: 0.001, code: "USD", expected: "$0.00"},

		// Arithmetic overflow after rounding
		{name: "MaxFloat64 overflow", amount: math.MaxFloat64, code: "USD", expected: "N/A"},
		{name: "Neg MaxFloat64 overflow", amount: -math.MaxFloat64, code: "USD", expected: "N/A"},

		// Negative rounds to zero vs nonzero
		{name: "Neg rounds to zero USD", amount: -0.004, code: "USD", expected: "$0.00"},
		{name: "Neg rounds to nonzero USD", amount: -0.005, code: "USD", expected: "-$0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := currency.FormatAmount(tt.amount, tt.code)
			if got != tt.expected {
				t.Errorf("FormatAmount(%v, %q) = %q, want %q", tt.amount, tt.code, got, tt.expected)
			}
		})
	}
}

// TestFormatAmountNoSymbol_VariousCases tests FormatAmountNoSymbol with various amounts.
func TestFormatAmountNoSymbol_VariousCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   float64
		code     string
		expected string
	}{
		// Basic formatting
		{name: "USD simple", amount: 1234.56, code: "USD", expected: "1,234.56"},
		{name: "USD large", amount: 1234567.89, code: "USD", expected: "1,234,567.89"},
		{name: "USD negative", amount: -1234.56, code: "USD", expected: "-1,234.56"},

		// JPY (0 decimals)
		{name: "JPY simple", amount: 1234.56, code: "JPY", expected: "1,235"},

		// KWD (3 decimals)
		{name: "KWD simple", amount: 1234.567, code: "KWD", expected: "1,234.567"},

		// Invalid currency (default 2 decimals)
		{name: "Invalid", amount: 1234.56, code: "XYZ", expected: "1,234.56"},

		// Edge cases
		{name: "Zero", amount: 0, code: "USD", expected: "0.00"},
		{name: "Small negative", amount: -0.01, code: "USD", expected: "-0.01"},

		// Floating-point edge cases
		{name: "NaN", amount: math.NaN(), code: "USD", expected: "N/A"},
		{name: "Positive Inf", amount: math.Inf(1), code: "USD", expected: "N/A"},
		{name: "Negative Inf", amount: math.Inf(-1), code: "USD", expected: "N/A"},

		// Arithmetic overflow after rounding
		{name: "MaxFloat64 overflow", amount: math.MaxFloat64, code: "USD", expected: "N/A"},

		// Negative zero
		{name: "Tiny negative rounds to zero", amount: -0.001, code: "USD", expected: "0.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := currency.FormatAmountNoSymbol(tt.amount, tt.code)
			if got != tt.expected {
				t.Errorf("FormatAmountNoSymbol(%v, %q) = %q, want %q", tt.amount, tt.code, got, tt.expected)
			}
		})
	}
}

// TestThousandsSeparators_EdgeCases tests thousands separator handling for various digit counts.
func TestThousandsSeparators_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   float64
		expected string
	}{
		{name: "1 digit", amount: 1, expected: "1.00"},
		{name: "2 digits", amount: 12, expected: "12.00"},
		{name: "3 digits", amount: 123, expected: "123.00"},
		{name: "4 digits", amount: 1234, expected: "1,234.00"},
		{name: "5 digits", amount: 12345, expected: "12,345.00"},
		{name: "6 digits", amount: 123456, expected: "123,456.00"},
		{name: "7 digits", amount: 1234567, expected: "1,234,567.00"},
		{name: "8 digits", amount: 12345678, expected: "12,345,678.00"},
		{name: "9 digits", amount: 123456789, expected: "123,456,789.00"},
		{name: "10 digits", amount: 1234567890, expected: "1,234,567,890.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := currency.FormatAmountNoSymbol(tt.amount, "USD")
			if got != tt.expected {
				t.Errorf("FormatAmountNoSymbol(%v, \"USD\") = %q, want %q", tt.amount, got, tt.expected)
			}
		})
	}
}

// TestFormatAmount_AllCurrencies iterates all currencies to validate formatting.
func TestFormatAmount_AllCurrencies(t *testing.T) {
	t.Parallel()

	for _, c := range currency.AllCurrencies() {
		t.Run(c.Code, func(t *testing.T) {
			t.Parallel()

			// Test with zero amount
			result := currency.FormatAmount(0, c.Code)
			validateFormat(t, c, result, 0)

			// Test with representative amount
			result = currency.FormatAmount(1.23456789, c.Code)
			validateFormat(t, c, result, 1.23456789)
		})
	}
}

// validateFormat validates that the formatted result has the correct prefix and decimal places.
func validateFormat(t *testing.T, c currency.Currency, result string, amount float64) {
	t.Helper()

	// Verify symbol or code prefix
	expectedPrefix := c.Symbol
	if expectedPrefix == "" {
		expectedPrefix = c.Code
	}
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Errorf("FormatAmount(%v, %q) = %q; want prefix %q",
			amount, c.Code, result, expectedPrefix)
	}

	// Verify decimal places match MinorUnits
	if strings.Contains(result, ".") {
		parts := strings.Split(result, ".")
		decimalPart := parts[len(parts)-1]
		if len(decimalPart) != c.MinorUnits {
			t.Errorf("FormatAmount(%v, %q) = %q; has %d decimals, want %d",
				amount, c.Code, result, len(decimalPart), c.MinorUnits)
		}
	} else if c.MinorUnits > 0 {
		t.Errorf("FormatAmount(%v, %q) = %q; missing decimal point for currency with %d minor units",
			amount, c.Code, result, c.MinorUnits)
	}
}

// Example functions for go doc discoverability.

func ExampleGetSymbol() {
	fmt.Println(currency.GetSymbol("USD"))
	fmt.Println(currency.GetSymbol("EUR"))
	fmt.Println(currency.GetSymbol("CHF")) // No symbol, returns code
	// Output:
	// $
	// €
	// CHF
}

func ExampleFormatAmount() {
	fmt.Println(currency.FormatAmount(1234.56, "USD"))
	fmt.Println(currency.FormatAmount(-1234.56, "EUR"))
	fmt.Println(currency.FormatAmount(1234.56, "JPY")) // 0 decimals
	// Output:
	// $1,234.56
	// -€1,234.56
	// ¥1,235
}

func ExampleFormatAmountNoSymbol() {
	fmt.Println(currency.FormatAmountNoSymbol(1234.56, "USD"))
	fmt.Println(currency.FormatAmountNoSymbol(1234.567, "KWD")) // 3 decimals
	// Output:
	// 1,234.56
	// 1,234.567
}
