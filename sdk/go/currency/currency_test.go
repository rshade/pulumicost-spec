package currency_test

import (
	"errors"
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/currency"
)

// T009: Table-driven tests for currency.IsValid() covering valid codes.
func TestIsValid_ValidCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want bool
	}{
		{name: "USD", code: "USD", want: true},
		{name: "EUR", code: "EUR", want: true},
		{name: "JPY", code: "JPY", want: true},
		{name: "GBP", code: "GBP", want: true},
		{name: "CHF", code: "CHF", want: true},
		{name: "CAD", code: "CAD", want: true},
		{name: "AUD", code: "AUD", want: true},
		{name: "CNY", code: "CNY", want: true},
		{name: "INR", code: "INR", want: true},
		{name: "BRL", code: "BRL", want: true},
		// Special codes
		{name: "XXX (no currency)", code: "XXX", want: true},
		{name: "XTS (test currency)", code: "XTS", want: true},
		{name: "XDR (SDR)", code: "XDR", want: true},
		{name: "XAU (gold)", code: "XAU", want: true},
		{name: "XAG (silver)", code: "XAG", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := currency.IsValid(tt.code); got != tt.want {
				t.Errorf("currency.IsValid(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// T010: Tests for invalid codes.
func TestIsValid_InvalidCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want bool
	}{
		{name: "XYZ (non-existent)", code: "XYZ", want: false},
		{name: "ABC (non-existent)", code: "ABC", want: false},
		{name: "lowercase usd", code: "usd", want: false},
		{name: "lowercase eur", code: "eur", want: false},
		{name: "mixed case Usd", code: "Usd", want: false},
		{name: "empty string", code: "", want: false},
		{name: "single char", code: "U", want: false},
		{name: "two chars", code: "US", want: false},
		{name: "four chars", code: "USDD", want: false},
		{name: "with leading space", code: " USD", want: false},
		{name: "with trailing space", code: "USD ", want: false},
		{name: "with spaces", code: " USD ", want: false},
		{name: "numeric", code: "123", want: false},
		{name: "special chars", code: "US$", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := currency.IsValid(tt.code); got != tt.want {
				t.Errorf("currency.IsValid(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// T011: Edge case tests (historic codes, supranational).
func TestIsValid_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code string
		want bool
	}{
		// Historic codes should be invalid
		{name: "DEM (historic)", code: "DEM", want: false},
		{name: "FRF (historic)", code: "FRF", want: false},
		{name: "ITL (historic)", code: "ITL", want: false},
		{name: "ESP (historic)", code: "ESP", want: false},
		// Supranational should be valid
		{name: "XDR (SDR)", code: "XDR", want: true},
		{name: "XSU (Sucre)", code: "XSU", want: true},
		{name: "XUA (ADB)", code: "XUA", want: true},
		// Precious metals should be valid
		{name: "XAU (gold)", code: "XAU", want: true},
		{name: "XAG (silver)", code: "XAG", want: true},
		{name: "XPT (platinum)", code: "XPT", want: true},
		{name: "XPD (palladium)", code: "XPD", want: true},
		// Test currency should be valid
		{name: "XTS (test)", code: "XTS", want: true},
		// No currency should be valid
		{name: "XXX (no currency)", code: "XXX", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := currency.IsValid(tt.code); got != tt.want {
				t.Errorf("currency.IsValid(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

// T016: Tests for GetCurrency() with valid codes.
func TestGetCurrency_ValidCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		code        string
		wantName    string
		wantNumeric string
		wantMinor   int
	}{
		{
			name:        "USD",
			code:        "USD",
			wantName:    "US Dollar",
			wantNumeric: "840",
			wantMinor:   2,
		},
		{
			name:        "JPY",
			code:        "JPY",
			wantName:    "Yen",
			wantNumeric: "392",
			wantMinor:   0,
		},
		{
			name:        "KWD",
			code:        "KWD",
			wantName:    "Kuwaiti Dinar",
			wantNumeric: "414",
			wantMinor:   3,
		},
		{
			name:        "EUR",
			code:        "EUR",
			wantName:    "Euro",
			wantNumeric: "978",
			wantMinor:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, err := currency.GetCurrency(tt.code)
			if err != nil {
				t.Fatalf("GetCurrency(%q) returned error: %v", tt.code, err)
			}
			if c.Code != tt.code {
				t.Errorf("Code = %q, want %q", c.Code, tt.code)
			}
			if c.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", c.Name, tt.wantName)
			}
			if c.NumericCode != tt.wantNumeric {
				t.Errorf("NumericCode = %q, want %q", c.NumericCode, tt.wantNumeric)
			}
			if c.MinorUnits != tt.wantMinor {
				t.Errorf("MinorUnits = %d, want %d", c.MinorUnits, tt.wantMinor)
			}
		})
	}
}

// T017: Tests for GetCurrency() with invalid codes.
func TestGetCurrency_InvalidCodes(t *testing.T) {
	t.Parallel()

	invalidCodes := []string{
		"XYZ",
		"ABC",
		"usd",
		"",
		"US",
		"USDD",
		" USD",
		"DEM", // historic
	}

	for _, code := range invalidCodes {
		t.Run(code, func(t *testing.T) {
			t.Parallel()
			c, err := currency.GetCurrency(code)
			if !errors.Is(err, currency.ErrCurrencyNotFound) {
				t.Errorf("GetCurrency(%q) error = %v, want %v", code, err, currency.ErrCurrencyNotFound)
			}
			if c != nil {
				t.Errorf("GetCurrency(%q) returned non-nil currency", code)
			}
		})
	}
}

// T018: Tests verifying metadata values.
func TestGetCurrency_MetadataValues(t *testing.T) {
	t.Parallel()

	// Test currencies with different minor units
	tests := []struct {
		code       string
		minorUnits int
	}{
		{"JPY", 0}, // No decimals
		{"USD", 2}, // Standard 2 decimals
		{"EUR", 2}, // Standard 2 decimals
		{"BHD", 3}, // 3 decimals (dinar)
		{"KWD", 3}, // 3 decimals (dinar)
		{"OMR", 3}, // 3 decimals (rial)
		{"CLF", 4}, // 4 decimals (UF)
		{"UYW", 4}, // 4 decimals
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			t.Parallel()
			c, err := currency.GetCurrency(tt.code)
			if err != nil {
				t.Fatalf("GetCurrency(%q) returned error: %v", tt.code, err)
			}
			if c.MinorUnits != tt.minorUnits {
				t.Errorf("MinorUnits for %s = %d, want %d", tt.code, c.MinorUnits, tt.minorUnits)
			}
		})
	}
}

// Regression: returned currency should be a copy and safe to modify without
// corrupting package-level data.
func TestGetCurrency_ReturnsCopy(t *testing.T) {
	t.Parallel()

	c, err := currency.GetCurrency("USD")
	if err != nil {
		t.Fatalf("GetCurrency(%q) returned error: %v", "USD", err)
	}

	// Mutate the returned currency
	c.Name = "Tampered"
	c.Code = "HAX"
	c.NumericCode = "000"
	c.MinorUnits = 7

	// Fetch again and ensure canonical data is unchanged
	canonical, err := currency.GetCurrency("USD")
	if err != nil {
		t.Fatalf("GetCurrency(%q) returned error: %v", "USD", err)
	}

	if canonical.Name == "Tampered" || canonical.Code != "USD" ||
		canonical.NumericCode != "840" || canonical.MinorUnits != 2 {
		t.Fatalf("GetCurrency returned mutable data: %+v", canonical)
	}
}

// T022: Test AllCurrencies returns 180+ currencies.
func TestAllCurrencies_Count(t *testing.T) {
	t.Parallel()

	currencies := currency.AllCurrencies()
	if len(currencies) < 180 {
		t.Errorf("AllCurrencies() returned %d currencies, want >= 180", len(currencies))
	}
}

// T023: Test each currency has non-empty fields.
func TestAllCurrencies_NonEmptyFields(t *testing.T) {
	t.Parallel()

	currencies := currency.AllCurrencies()
	for _, c := range currencies {
		if c.Code == "" {
			t.Error("Found currency with empty Code")
		}
		if c.Name == "" {
			t.Errorf("Currency %s has empty Name", c.Code)
		}
		if c.NumericCode == "" {
			t.Errorf("Currency %s has empty NumericCode", c.Code)
		}
		if len(c.Code) != 3 {
			t.Errorf("Currency %s has Code length %d, want 3", c.Code, len(c.Code))
		}
		if len(c.NumericCode) != 3 {
			t.Errorf("Currency %s has NumericCode length %d, want 3", c.Code, len(c.NumericCode))
		}
		if c.MinorUnits < 0 || c.MinorUnits > 4 {
			t.Errorf("Currency %s has MinorUnits %d, want 0-4", c.Code, c.MinorUnits)
		}
	}
}

// T024: Test list is sorted alphabetically by code.
func TestAllCurrencies_Sorted(t *testing.T) {
	t.Parallel()

	currencies := currency.AllCurrencies()
	for i := 1; i < len(currencies); i++ {
		if currencies[i-1].Code >= currencies[i].Code {
			t.Errorf("Currencies not sorted: %s >= %s at index %d",
				currencies[i-1].Code, currencies[i].Code, i)
		}
	}
}

// Test String() method.
func TestCurrency_String(t *testing.T) {
	t.Parallel()

	c := currency.Currency{Code: "USD", Name: "US Dollar", NumericCode: "840", MinorUnits: 2}
	if got := c.String(); got != "USD" {
		t.Errorf("Currency.String() = %q, want %q", got, "USD")
	}
}

// Test error type.
func TestErrCurrencyNotFound(t *testing.T) {
	t.Parallel()

	if currency.ErrCurrencyNotFound.Error() != "currency not found" {
		t.Errorf("ErrCurrencyNotFound.Error() = %q, want %q",
			currency.ErrCurrencyNotFound.Error(), "currency not found")
	}
}
