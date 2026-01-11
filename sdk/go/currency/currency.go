// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package currency provides ISO 4217 currency validation and metadata.
//
// The package implements zero-allocation validation following the pattern
// established in sdk/go/registry/domain.go.
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

// String returns the currency code.
// Implements fmt.Stringer interface.
func (c Currency) String() string {
	return c.Code
}

// allCurrencies is a package-level slice containing all valid ISO 4217 currencies.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals,mnd // Intentional optimization for zero-allocation validation; minor units are ISO-defined constants
var allCurrencies = []Currency{
	{Code: "AED", Name: "UAE Dirham", NumericCode: "784", MinorUnits: 2},
	{Code: "AFN", Name: "Afghani", NumericCode: "971", MinorUnits: 2},
	{Code: "ALL", Name: "Lek", NumericCode: "008", MinorUnits: 2},
	{Code: "AMD", Name: "Armenian Dram", NumericCode: "051", MinorUnits: 2},
	{Code: "ANG", Name: "Netherlands Antillean Guilder", NumericCode: "532", MinorUnits: 2},
	{Code: "AOA", Name: "Kwanza", NumericCode: "973", MinorUnits: 2},
	{Code: "ARS", Name: "Argentine Peso", NumericCode: "032", MinorUnits: 2},
	{Code: "AUD", Name: "Australian Dollar", NumericCode: "036", MinorUnits: 2},
	{Code: "AWG", Name: "Aruban Florin", NumericCode: "533", MinorUnits: 2},
	{Code: "AZN", Name: "Azerbaijan Manat", NumericCode: "944", MinorUnits: 2},
	{Code: "BAM", Name: "Convertible Mark", NumericCode: "977", MinorUnits: 2},
	{Code: "BBD", Name: "Barbados Dollar", NumericCode: "052", MinorUnits: 2},
	{Code: "BDT", Name: "Taka", NumericCode: "050", MinorUnits: 2},
	{Code: "BGN", Name: "Bulgarian Lev", NumericCode: "975", MinorUnits: 2},
	{Code: "BHD", Name: "Bahraini Dinar", NumericCode: "048", MinorUnits: 3},
	{Code: "BIF", Name: "Burundi Franc", NumericCode: "108", MinorUnits: 0},
	{Code: "BMD", Name: "Bermudian Dollar", NumericCode: "060", MinorUnits: 2},
	{Code: "BND", Name: "Brunei Dollar", NumericCode: "096", MinorUnits: 2},
	{Code: "BOB", Name: "Boliviano", NumericCode: "068", MinorUnits: 2},
	{Code: "BOV", Name: "Mvdol", NumericCode: "984", MinorUnits: 2},
	{Code: "BRL", Name: "Brazilian Real", NumericCode: "986", MinorUnits: 2},
	{Code: "BSD", Name: "Bahamian Dollar", NumericCode: "044", MinorUnits: 2},
	{Code: "BTN", Name: "Ngultrum", NumericCode: "064", MinorUnits: 2},
	{Code: "BWP", Name: "Pula", NumericCode: "072", MinorUnits: 2},
	{Code: "BYN", Name: "Belarusian Ruble", NumericCode: "933", MinorUnits: 2},
	{Code: "BZD", Name: "Belize Dollar", NumericCode: "084", MinorUnits: 2},
	{Code: "CAD", Name: "Canadian Dollar", NumericCode: "124", MinorUnits: 2},
	{Code: "CDF", Name: "Congolese Franc", NumericCode: "976", MinorUnits: 2},
	{Code: "CHE", Name: "WIR Euro", NumericCode: "947", MinorUnits: 2},
	{Code: "CHF", Name: "Swiss Franc", NumericCode: "756", MinorUnits: 2},
	{Code: "CHW", Name: "WIR Franc", NumericCode: "948", MinorUnits: 2},
	{Code: "CLF", Name: "Unidad de Fomento", NumericCode: "990", MinorUnits: 4},
	{Code: "CLP", Name: "Chilean Peso", NumericCode: "152", MinorUnits: 0},
	{Code: "CNY", Name: "Yuan Renminbi", NumericCode: "156", MinorUnits: 2},
	{Code: "COP", Name: "Colombian Peso", NumericCode: "170", MinorUnits: 2},
	{Code: "COU", Name: "Unidad de Valor Real", NumericCode: "970", MinorUnits: 2},
	{Code: "CRC", Name: "Costa Rican Colon", NumericCode: "188", MinorUnits: 2},
	{Code: "CUC", Name: "Peso Convertible", NumericCode: "931", MinorUnits: 2},
	{Code: "CUP", Name: "Cuban Peso", NumericCode: "192", MinorUnits: 2},
	{Code: "CVE", Name: "Cabo Verde Escudo", NumericCode: "132", MinorUnits: 2},
	{Code: "CZK", Name: "Czech Koruna", NumericCode: "203", MinorUnits: 2},
	{Code: "DJF", Name: "Djibouti Franc", NumericCode: "262", MinorUnits: 0},
	{Code: "DKK", Name: "Danish Krone", NumericCode: "208", MinorUnits: 2},
	{Code: "DOP", Name: "Dominican Peso", NumericCode: "214", MinorUnits: 2},
	{Code: "DZD", Name: "Algerian Dinar", NumericCode: "012", MinorUnits: 2},
	{Code: "EGP", Name: "Egyptian Pound", NumericCode: "818", MinorUnits: 2},
	{Code: "ERN", Name: "Nakfa", NumericCode: "232", MinorUnits: 2},
	{Code: "ETB", Name: "Ethiopian Birr", NumericCode: "230", MinorUnits: 2},
	{Code: "EUR", Name: "Euro", NumericCode: "978", MinorUnits: 2},
	{Code: "FJD", Name: "Fiji Dollar", NumericCode: "242", MinorUnits: 2},
	{Code: "FKP", Name: "Falkland Islands Pound", NumericCode: "238", MinorUnits: 2},
	{Code: "GBP", Name: "Pound Sterling", NumericCode: "826", MinorUnits: 2},
	{Code: "GEL", Name: "Lari", NumericCode: "981", MinorUnits: 2},
	{Code: "GHS", Name: "Ghana Cedi", NumericCode: "936", MinorUnits: 2},
	{Code: "GIP", Name: "Gibraltar Pound", NumericCode: "292", MinorUnits: 2},
	{Code: "GMD", Name: "Dalasi", NumericCode: "270", MinorUnits: 2},
	{Code: "GNF", Name: "Guinean Franc", NumericCode: "324", MinorUnits: 0},
	{Code: "GTQ", Name: "Quetzal", NumericCode: "320", MinorUnits: 2},
	{Code: "GYD", Name: "Guyana Dollar", NumericCode: "328", MinorUnits: 2},
	{Code: "HKD", Name: "Hong Kong Dollar", NumericCode: "344", MinorUnits: 2},
	{Code: "HNL", Name: "Lempira", NumericCode: "340", MinorUnits: 2},
	{Code: "HRK", Name: "Kuna", NumericCode: "191", MinorUnits: 2},
	{Code: "HTG", Name: "Gourde", NumericCode: "332", MinorUnits: 2},
	{Code: "HUF", Name: "Forint", NumericCode: "348", MinorUnits: 2},
	{Code: "IDR", Name: "Rupiah", NumericCode: "360", MinorUnits: 2},
	{Code: "ILS", Name: "New Israeli Sheqel", NumericCode: "376", MinorUnits: 2},
	{Code: "INR", Name: "Indian Rupee", NumericCode: "356", MinorUnits: 2},
	{Code: "IQD", Name: "Iraqi Dinar", NumericCode: "368", MinorUnits: 3},
	{Code: "IRR", Name: "Iranian Rial", NumericCode: "364", MinorUnits: 2},
	{Code: "ISK", Name: "Iceland Krona", NumericCode: "352", MinorUnits: 0},
	{Code: "JMD", Name: "Jamaican Dollar", NumericCode: "388", MinorUnits: 2},
	{Code: "JOD", Name: "Jordanian Dinar", NumericCode: "400", MinorUnits: 3},
	{Code: "JPY", Name: "Yen", NumericCode: "392", MinorUnits: 0},
	{Code: "KES", Name: "Kenyan Shilling", NumericCode: "404", MinorUnits: 2},
	{Code: "KGS", Name: "Som", NumericCode: "417", MinorUnits: 2},
	{Code: "KHR", Name: "Riel", NumericCode: "116", MinorUnits: 2},
	{Code: "KMF", Name: "Comorian Franc", NumericCode: "174", MinorUnits: 0},
	{Code: "KPW", Name: "North Korean Won", NumericCode: "408", MinorUnits: 2},
	{Code: "KRW", Name: "Won", NumericCode: "410", MinorUnits: 0},
	{Code: "KWD", Name: "Kuwaiti Dinar", NumericCode: "414", MinorUnits: 3},
	{Code: "KYD", Name: "Cayman Islands Dollar", NumericCode: "136", MinorUnits: 2},
	{Code: "KZT", Name: "Tenge", NumericCode: "398", MinorUnits: 2},
	{Code: "LAK", Name: "Lao Kip", NumericCode: "418", MinorUnits: 2},
	{Code: "LBP", Name: "Lebanese Pound", NumericCode: "422", MinorUnits: 2},
	{Code: "LKR", Name: "Sri Lanka Rupee", NumericCode: "144", MinorUnits: 2},
	{Code: "LRD", Name: "Liberian Dollar", NumericCode: "430", MinorUnits: 2},
	{Code: "LSL", Name: "Loti", NumericCode: "426", MinorUnits: 2},
	{Code: "LYD", Name: "Libyan Dinar", NumericCode: "434", MinorUnits: 3},
	{Code: "MAD", Name: "Moroccan Dirham", NumericCode: "504", MinorUnits: 2},
	{Code: "MDL", Name: "Moldovan Leu", NumericCode: "498", MinorUnits: 2},
	{Code: "MGA", Name: "Malagasy Ariary", NumericCode: "969", MinorUnits: 2},
	{Code: "MKD", Name: "Denar", NumericCode: "807", MinorUnits: 2},
	{Code: "MMK", Name: "Kyat", NumericCode: "104", MinorUnits: 2},
	{Code: "MNT", Name: "Tugrik", NumericCode: "496", MinorUnits: 2},
	{Code: "MOP", Name: "Pataca", NumericCode: "446", MinorUnits: 2},
	{Code: "MRU", Name: "Ouguiya", NumericCode: "929", MinorUnits: 2},
	{Code: "MUR", Name: "Mauritius Rupee", NumericCode: "480", MinorUnits: 2},
	{Code: "MVR", Name: "Rufiyaa", NumericCode: "462", MinorUnits: 2},
	{Code: "MWK", Name: "Malawi Kwacha", NumericCode: "454", MinorUnits: 2},
	{Code: "MXN", Name: "Mexican Peso", NumericCode: "484", MinorUnits: 2},
	{Code: "MXV", Name: "Mexican Unidad de Inversion", NumericCode: "979", MinorUnits: 2},
	{Code: "MYR", Name: "Malaysian Ringgit", NumericCode: "458", MinorUnits: 2},
	{Code: "MZN", Name: "Mozambique Metical", NumericCode: "943", MinorUnits: 2},
	{Code: "NAD", Name: "Namibia Dollar", NumericCode: "516", MinorUnits: 2},
	{Code: "NGN", Name: "Naira", NumericCode: "566", MinorUnits: 2},
	{Code: "NIO", Name: "Cordoba Oro", NumericCode: "558", MinorUnits: 2},
	{Code: "NOK", Name: "Norwegian Krone", NumericCode: "578", MinorUnits: 2},
	{Code: "NPR", Name: "Nepalese Rupee", NumericCode: "524", MinorUnits: 2},
	{Code: "NZD", Name: "New Zealand Dollar", NumericCode: "554", MinorUnits: 2},
	{Code: "OMR", Name: "Rial Omani", NumericCode: "512", MinorUnits: 3},
	{Code: "PAB", Name: "Balboa", NumericCode: "590", MinorUnits: 2},
	{Code: "PEN", Name: "Sol", NumericCode: "604", MinorUnits: 2},
	{Code: "PGK", Name: "Kina", NumericCode: "598", MinorUnits: 2},
	{Code: "PHP", Name: "Philippine Peso", NumericCode: "608", MinorUnits: 2},
	{Code: "PKR", Name: "Pakistan Rupee", NumericCode: "586", MinorUnits: 2},
	{Code: "PLN", Name: "Zloty", NumericCode: "985", MinorUnits: 2},
	{Code: "PYG", Name: "Guarani", NumericCode: "600", MinorUnits: 0},
	{Code: "QAR", Name: "Qatari Rial", NumericCode: "634", MinorUnits: 2},
	{Code: "RON", Name: "Romanian Leu", NumericCode: "946", MinorUnits: 2},
	{Code: "RSD", Name: "Serbian Dinar", NumericCode: "941", MinorUnits: 2},
	{Code: "RUB", Name: "Russian Ruble", NumericCode: "643", MinorUnits: 2},
	{Code: "RWF", Name: "Rwanda Franc", NumericCode: "646", MinorUnits: 0},
	{Code: "SAR", Name: "Saudi Riyal", NumericCode: "682", MinorUnits: 2},
	{Code: "SBD", Name: "Solomon Islands Dollar", NumericCode: "090", MinorUnits: 2},
	{Code: "SCR", Name: "Seychelles Rupee", NumericCode: "690", MinorUnits: 2},
	{Code: "SDG", Name: "Sudanese Pound", NumericCode: "938", MinorUnits: 2},
	{Code: "SEK", Name: "Swedish Krona", NumericCode: "752", MinorUnits: 2},
	{Code: "SGD", Name: "Singapore Dollar", NumericCode: "702", MinorUnits: 2},
	{Code: "SHP", Name: "Saint Helena Pound", NumericCode: "654", MinorUnits: 2},
	{Code: "SLE", Name: "Leone", NumericCode: "925", MinorUnits: 2},
	{Code: "SOS", Name: "Somali Shilling", NumericCode: "706", MinorUnits: 2},
	{Code: "SRD", Name: "Surinam Dollar", NumericCode: "968", MinorUnits: 2},
	{Code: "SSP", Name: "South Sudanese Pound", NumericCode: "728", MinorUnits: 2},
	{Code: "STN", Name: "Dobra", NumericCode: "930", MinorUnits: 2},
	{Code: "SVC", Name: "El Salvador Colon", NumericCode: "222", MinorUnits: 2},
	{Code: "SYP", Name: "Syrian Pound", NumericCode: "760", MinorUnits: 2},
	{Code: "SZL", Name: "Lilangeni", NumericCode: "748", MinorUnits: 2},
	{Code: "THB", Name: "Baht", NumericCode: "764", MinorUnits: 2},
	{Code: "TJS", Name: "Somoni", NumericCode: "972", MinorUnits: 2},
	{Code: "TMT", Name: "Turkmenistan New Manat", NumericCode: "934", MinorUnits: 2},
	{Code: "TND", Name: "Tunisian Dinar", NumericCode: "788", MinorUnits: 3},
	{Code: "TOP", Name: "Pa'anga", NumericCode: "776", MinorUnits: 2},
	{Code: "TRY", Name: "Turkish Lira", NumericCode: "949", MinorUnits: 2},
	{Code: "TTD", Name: "Trinidad and Tobago Dollar", NumericCode: "780", MinorUnits: 2},
	{Code: "TWD", Name: "New Taiwan Dollar", NumericCode: "901", MinorUnits: 2},
	{Code: "TZS", Name: "Tanzanian Shilling", NumericCode: "834", MinorUnits: 2},
	{Code: "UAH", Name: "Hryvnia", NumericCode: "980", MinorUnits: 2},
	{Code: "UGX", Name: "Uganda Shilling", NumericCode: "800", MinorUnits: 0},
	{Code: "USD", Name: "US Dollar", NumericCode: "840", MinorUnits: 2},
	{Code: "USN", Name: "US Dollar (Next day)", NumericCode: "997", MinorUnits: 2},
	{Code: "UYI", Name: "Uruguay Peso en Unidades Indexadas", NumericCode: "940", MinorUnits: 0},
	{Code: "UYU", Name: "Peso Uruguayo", NumericCode: "858", MinorUnits: 2},
	{Code: "UYW", Name: "Unidad Previsional", NumericCode: "927", MinorUnits: 4},
	{Code: "UZS", Name: "Uzbekistan Sum", NumericCode: "860", MinorUnits: 2},
	{Code: "VED", Name: "Bolivar Soberano", NumericCode: "926", MinorUnits: 2},
	{Code: "VES", Name: "Bolivar Soberano", NumericCode: "928", MinorUnits: 2},
	{Code: "VND", Name: "Dong", NumericCode: "704", MinorUnits: 0},
	{Code: "VUV", Name: "Vatu", NumericCode: "548", MinorUnits: 0},
	{Code: "WST", Name: "Tala", NumericCode: "882", MinorUnits: 2},
	{Code: "XAF", Name: "CFA Franc BEAC", NumericCode: "950", MinorUnits: 0},
	{Code: "XAG", Name: "Silver", NumericCode: "961", MinorUnits: 0},
	{Code: "XAU", Name: "Gold", NumericCode: "959", MinorUnits: 0},
	{Code: "XBA", Name: "Bond Markets Unit European Composite Unit (EURCO)", NumericCode: "955", MinorUnits: 0},
	{Code: "XBB", Name: "Bond Markets Unit European Monetary Unit (E.M.U.-6)", NumericCode: "956", MinorUnits: 0},
	{Code: "XBC", Name: "Bond Markets Unit European Unit of Account 9 (E.U.A.-9)", NumericCode: "957", MinorUnits: 0},
	{Code: "XBD", Name: "Bond Markets Unit European Unit of Account 17 (E.U.A.-17)", NumericCode: "958", MinorUnits: 0},
	{Code: "XCD", Name: "East Caribbean Dollar", NumericCode: "951", MinorUnits: 2},
	{Code: "XDR", Name: "SDR (Special Drawing Right)", NumericCode: "960", MinorUnits: 0},
	{Code: "XOF", Name: "CFA Franc BCEAO", NumericCode: "952", MinorUnits: 0},
	{Code: "XPD", Name: "Palladium", NumericCode: "964", MinorUnits: 0},
	{Code: "XPF", Name: "CFP Franc", NumericCode: "953", MinorUnits: 0},
	{Code: "XPT", Name: "Platinum", NumericCode: "962", MinorUnits: 0},
	{Code: "XSU", Name: "Sucre", NumericCode: "994", MinorUnits: 0},
	{Code: "XTS", Name: "Codes specifically reserved for testing purposes", NumericCode: "963", MinorUnits: 0},
	{Code: "XUA", Name: "ADB Unit of Account", NumericCode: "965", MinorUnits: 0},
	{Code: "XXX", Name: "No currency", NumericCode: "999", MinorUnits: 0},
	{Code: "YER", Name: "Yemeni Rial", NumericCode: "886", MinorUnits: 2},
	{Code: "ZAR", Name: "Rand", NumericCode: "710", MinorUnits: 2},
	{Code: "ZMW", Name: "Zambian Kwacha", NumericCode: "967", MinorUnits: 2},
	{Code: "ZWL", Name: "Zimbabwe Dollar", NumericCode: "932", MinorUnits: 2},
}

// currencyByCode provides O(1) lookup for GetCurrency().
// Built at package initialization from allCurrencies.
//
//nolint:gochecknoglobals // Intentional optimization for O(1) lookup
var currencyByCode map[string]*Currency

// init builds the currencyByCode map from allCurrencies slice.
//
//nolint:gochecknoinits // Required for package initialization
func init() {
	currencyByCode = make(map[string]*Currency, len(allCurrencies))
	for i := range allCurrencies {
		currencyByCode[allCurrencies[i].Code] = &allCurrencies[i]
	}
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
//	fmt.Printf("%s has %d decimal places
", c.Name, c.MinorUnits)
//	// Output: US Dollar has 2 decimal places
func GetCurrency(code string) (*Currency, error) {
	if c, ok := currencyByCode[code]; ok {
		cpy := *c
		return &cpy, nil
	}
	return nil, ErrCurrencyNotFound
}

// AllCurrencies returns a slice of all valid ISO 4217 currencies.
// The returned slice is a reference to package-level data and MUST NOT be modified.
//
// The slice contains 180+ active currencies sorted alphabetically by code.
//
// Example:
//
//	for _, c := range currency.AllCurrencies() {
//	    fmt.Printf("%s: %s (%d decimals)
", c.Code, c.Name, c.MinorUnits)
//	}
func AllCurrencies() []Currency {
	return allCurrencies
}
