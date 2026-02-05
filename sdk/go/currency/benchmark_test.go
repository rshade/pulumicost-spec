package currency_test

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/currency"
)

// T014: Benchmark for currency.IsValid() targeting <15 ns/op.
func BenchmarkIsValid(b *testing.B) {
	// Benchmark with a common currency (early in list)
	b.Run("USD", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			currency.IsValid("USD")
		}
	})

	// Benchmark with a currency in middle of list
	b.Run("MXN", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			currency.IsValid("MXN")
		}
	})

	// Benchmark with a currency at end of list
	b.Run("ZWL", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			currency.IsValid("ZWL")
		}
	})

	// Benchmark with invalid code (worst case - full scan)
	b.Run("Invalid", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			currency.IsValid("XYZ")
		}
	})
}

// T021: Benchmark for GetCurrency().
func BenchmarkGetCurrency(b *testing.B) {
	b.Run("USD", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, _ = currency.GetCurrency("USD")
		}
	})

	b.Run("JPY", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, _ = currency.GetCurrency("JPY")
		}
	})

	b.Run("Invalid", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_, _ = currency.GetCurrency("XYZ")
		}
	})
}

// T027: Benchmark for AllCurrencies().
func BenchmarkAllCurrencies(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = currency.AllCurrencies()
	}
}

// Benchmark String() method.
func BenchmarkCurrency_String(b *testing.B) {
	c := currency.Currency{Code: "USD", Name: "US Dollar", NumericCode: "840", MinorUnits: 2, Symbol: "$"}
	b.ReportAllocs()
	for range b.N {
		_ = c.String()
	}
}

// Benchmark GetSymbol() - targeting <20 ns/op, 0 allocs/op.
func BenchmarkGetSymbol(b *testing.B) {
	b.Run("USD", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.GetSymbol("USD")
		}
	})

	b.Run("EUR", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.GetSymbol("EUR")
		}
	})

	b.Run("CHF_NoSymbol", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.GetSymbol("CHF")
		}
	})

	b.Run("Invalid", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.GetSymbol("XYZ")
		}
	})
}

// Benchmark FormatAmount() - targeting <500 ns/op.
func BenchmarkFormatAmount(b *testing.B) {
	b.Run("USD_Simple", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.FormatAmount(1234.56, "USD")
		}
	})

	b.Run("USD_Large", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.FormatAmount(1234567890.12, "USD")
		}
	})

	b.Run("JPY_NoDecimals", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.FormatAmount(123456789, "JPY")
		}
	})

	b.Run("KWD_ThreeDecimals", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.FormatAmount(1234.567, "KWD")
		}
	})
}

// Benchmark FormatAmountNoSymbol().
func BenchmarkFormatAmountNoSymbol(b *testing.B) {
	b.Run("USD_Simple", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.FormatAmountNoSymbol(1234.56, "USD")
		}
	})

	b.Run("USD_Large", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			_ = currency.FormatAmountNoSymbol(1234567890.12, "USD")
		}
	})
}
