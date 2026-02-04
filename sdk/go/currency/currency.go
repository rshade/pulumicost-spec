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

	// Symbol is the currency symbol (e.g., "$", "€", "£").
	// Empty string for currencies without a commonly used symbol.
	Symbol string
}

// String returns the currency code.
// Implements fmt.Stringer interface.
func (c Currency) String() string {
	return c.Code
}

// allCurrencies is a package-level slice containing all valid ISO 4217 currencies.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals,mnd,golines // Intentional optimization; minor units/names are ISO-defined constants
var allCurrencies = []Currency{
	{Code: "AED", Name: "UAE Dirham", NumericCode: "784", MinorUnits: 2, Symbol: "د.إ"},
	{Code: "AFN", Name: "Afghani", NumericCode: "971", MinorUnits: 2, Symbol: "؋"},
	{Code: "ALL", Name: "Lek", NumericCode: "008", MinorUnits: 2, Symbol: "L"},
	{Code: "AMD", Name: "Armenian Dram", NumericCode: "051", MinorUnits: 2, Symbol: "֏"},
	{Code: "ANG", Name: "Netherlands Antillean Guilder", NumericCode: "532", MinorUnits: 2, Symbol: "ƒ"},
	{Code: "AOA", Name: "Kwanza", NumericCode: "973", MinorUnits: 2, Symbol: "Kz"},
	{Code: "ARS", Name: "Argentine Peso", NumericCode: "032", MinorUnits: 2, Symbol: "$"},
	{Code: "AUD", Name: "Australian Dollar", NumericCode: "036", MinorUnits: 2, Symbol: "A$"},
	{Code: "AWG", Name: "Aruban Florin", NumericCode: "533", MinorUnits: 2, Symbol: "ƒ"},
	{Code: "AZN", Name: "Azerbaijan Manat", NumericCode: "944", MinorUnits: 2, Symbol: "₼"},
	{Code: "BAM", Name: "Convertible Mark", NumericCode: "977", MinorUnits: 2, Symbol: "KM"},
	{Code: "BBD", Name: "Barbados Dollar", NumericCode: "052", MinorUnits: 2, Symbol: "$"},
	{Code: "BDT", Name: "Taka", NumericCode: "050", MinorUnits: 2, Symbol: "৳"},
	{Code: "BGN", Name: "Bulgarian Lev", NumericCode: "975", MinorUnits: 2, Symbol: "лв"},
	{Code: "BHD", Name: "Bahraini Dinar", NumericCode: "048", MinorUnits: 3, Symbol: ".د.ب"},
	{Code: "BIF", Name: "Burundi Franc", NumericCode: "108", MinorUnits: 0, Symbol: "FBu"},
	{Code: "BMD", Name: "Bermudian Dollar", NumericCode: "060", MinorUnits: 2, Symbol: "$"},
	{Code: "BND", Name: "Brunei Dollar", NumericCode: "096", MinorUnits: 2, Symbol: "$"},
	{Code: "BOB", Name: "Boliviano", NumericCode: "068", MinorUnits: 2, Symbol: "Bs."},
	{Code: "BOV", Name: "Mvdol", NumericCode: "984", MinorUnits: 2, Symbol: ""},
	{Code: "BRL", Name: "Brazilian Real", NumericCode: "986", MinorUnits: 2, Symbol: "R$"},
	{Code: "BSD", Name: "Bahamian Dollar", NumericCode: "044", MinorUnits: 2, Symbol: "$"},
	{Code: "BTN", Name: "Ngultrum", NumericCode: "064", MinorUnits: 2, Symbol: "Nu."},
	{Code: "BWP", Name: "Pula", NumericCode: "072", MinorUnits: 2, Symbol: "P"},
	{Code: "BYN", Name: "Belarusian Ruble", NumericCode: "933", MinorUnits: 2, Symbol: "Br"},
	{Code: "BZD", Name: "Belize Dollar", NumericCode: "084", MinorUnits: 2, Symbol: "BZ$"},
	{Code: "CAD", Name: "Canadian Dollar", NumericCode: "124", MinorUnits: 2, Symbol: "C$"},
	{Code: "CDF", Name: "Congolese Franc", NumericCode: "976", MinorUnits: 2, Symbol: "FC"},
	{Code: "CHE", Name: "WIR Euro", NumericCode: "947", MinorUnits: 2, Symbol: ""},
	{Code: "CHF", Name: "Swiss Franc", NumericCode: "756", MinorUnits: 2, Symbol: ""},
	{Code: "CHW", Name: "WIR Franc", NumericCode: "948", MinorUnits: 2, Symbol: ""},
	{Code: "CLF", Name: "Unidad de Fomento", NumericCode: "990", MinorUnits: 4, Symbol: "UF"},
	{Code: "CLP", Name: "Chilean Peso", NumericCode: "152", MinorUnits: 0, Symbol: "$"},
	{Code: "CNY", Name: "Yuan Renminbi", NumericCode: "156", MinorUnits: 2, Symbol: "¥"},
	{Code: "COP", Name: "Colombian Peso", NumericCode: "170", MinorUnits: 2, Symbol: "$"},
	{Code: "COU", Name: "Unidad de Valor Real", NumericCode: "970", MinorUnits: 2, Symbol: ""},
	{Code: "CRC", Name: "Costa Rican Colon", NumericCode: "188", MinorUnits: 2, Symbol: "₡"},
	{Code: "CUC", Name: "Peso Convertible", NumericCode: "931", MinorUnits: 2, Symbol: "$"},
	{Code: "CUP", Name: "Cuban Peso", NumericCode: "192", MinorUnits: 2, Symbol: "₱"},
	{Code: "CVE", Name: "Cabo Verde Escudo", NumericCode: "132", MinorUnits: 2, Symbol: "$"},
	{Code: "CZK", Name: "Czech Koruna", NumericCode: "203", MinorUnits: 2, Symbol: "Kč"},
	{Code: "DJF", Name: "Djibouti Franc", NumericCode: "262", MinorUnits: 0, Symbol: "Fdj"},
	{Code: "DKK", Name: "Danish Krone", NumericCode: "208", MinorUnits: 2, Symbol: "kr"},
	{Code: "DOP", Name: "Dominican Peso", NumericCode: "214", MinorUnits: 2, Symbol: "RD$"},
	{Code: "DZD", Name: "Algerian Dinar", NumericCode: "012", MinorUnits: 2, Symbol: "د.ج"},
	{Code: "EGP", Name: "Egyptian Pound", NumericCode: "818", MinorUnits: 2, Symbol: "£"},
	{Code: "ERN", Name: "Nakfa", NumericCode: "232", MinorUnits: 2, Symbol: "Nfk"},
	{Code: "ETB", Name: "Ethiopian Birr", NumericCode: "230", MinorUnits: 2, Symbol: "Br"},
	{Code: "EUR", Name: "Euro", NumericCode: "978", MinorUnits: 2, Symbol: "€"},
	{Code: "FJD", Name: "Fiji Dollar", NumericCode: "242", MinorUnits: 2, Symbol: "$"},
	{Code: "FKP", Name: "Falkland Islands Pound", NumericCode: "238", MinorUnits: 2, Symbol: "£"},
	{Code: "GBP", Name: "Pound Sterling", NumericCode: "826", MinorUnits: 2, Symbol: "£"},
	{Code: "GEL", Name: "Lari", NumericCode: "981", MinorUnits: 2, Symbol: "₾"},
	{Code: "GHS", Name: "Ghana Cedi", NumericCode: "936", MinorUnits: 2, Symbol: "₵"},
	{Code: "GIP", Name: "Gibraltar Pound", NumericCode: "292", MinorUnits: 2, Symbol: "£"},
	{Code: "GMD", Name: "Dalasi", NumericCode: "270", MinorUnits: 2, Symbol: "D"},
	{Code: "GNF", Name: "Guinean Franc", NumericCode: "324", MinorUnits: 0, Symbol: "FG"},
	{Code: "GTQ", Name: "Quetzal", NumericCode: "320", MinorUnits: 2, Symbol: "Q"},
	{Code: "GYD", Name: "Guyana Dollar", NumericCode: "328", MinorUnits: 2, Symbol: "$"},
	{Code: "HKD", Name: "Hong Kong Dollar", NumericCode: "344", MinorUnits: 2, Symbol: "HK$"},
	{Code: "HNL", Name: "Lempira", NumericCode: "340", MinorUnits: 2, Symbol: "L"},
	{Code: "HRK", Name: "Kuna", NumericCode: "191", MinorUnits: 2, Symbol: "kn"},
	{Code: "HTG", Name: "Gourde", NumericCode: "332", MinorUnits: 2, Symbol: "G"},
	{Code: "HUF", Name: "Forint", NumericCode: "348", MinorUnits: 2, Symbol: "Ft"},
	{Code: "IDR", Name: "Rupiah", NumericCode: "360", MinorUnits: 2, Symbol: "Rp"},
	{Code: "ILS", Name: "New Israeli Sheqel", NumericCode: "376", MinorUnits: 2, Symbol: "₪"},
	{Code: "INR", Name: "Indian Rupee", NumericCode: "356", MinorUnits: 2, Symbol: "₹"},
	{Code: "IQD", Name: "Iraqi Dinar", NumericCode: "368", MinorUnits: 3, Symbol: "ع.د"},
	{Code: "IRR", Name: "Iranian Rial", NumericCode: "364", MinorUnits: 2, Symbol: "﷼"},
	{Code: "ISK", Name: "Iceland Krona", NumericCode: "352", MinorUnits: 0, Symbol: "kr"},
	{Code: "JMD", Name: "Jamaican Dollar", NumericCode: "388", MinorUnits: 2, Symbol: "J$"},
	{Code: "JOD", Name: "Jordanian Dinar", NumericCode: "400", MinorUnits: 3, Symbol: "د.ا"},
	{Code: "JPY", Name: "Yen", NumericCode: "392", MinorUnits: 0, Symbol: "¥"},
	{Code: "KES", Name: "Kenyan Shilling", NumericCode: "404", MinorUnits: 2, Symbol: "KSh"},
	{Code: "KGS", Name: "Som", NumericCode: "417", MinorUnits: 2, Symbol: "KGS"},
	{Code: "KHR", Name: "Riel", NumericCode: "116", MinorUnits: 2, Symbol: "៛"},
	{Code: "KMF", Name: "Comorian Franc", NumericCode: "174", MinorUnits: 0, Symbol: "CF"},
	{Code: "KPW", Name: "North Korean Won", NumericCode: "408", MinorUnits: 2, Symbol: "₩"},
	{Code: "KRW", Name: "Won", NumericCode: "410", MinorUnits: 0, Symbol: "₩"},
	{Code: "KWD", Name: "Kuwaiti Dinar", NumericCode: "414", MinorUnits: 3, Symbol: "د.ك"},
	{Code: "KYD", Name: "Cayman Islands Dollar", NumericCode: "136", MinorUnits: 2, Symbol: "$"},
	{Code: "KZT", Name: "Tenge", NumericCode: "398", MinorUnits: 2, Symbol: "₸"},
	{Code: "LAK", Name: "Lao Kip", NumericCode: "418", MinorUnits: 2, Symbol: "₭"},
	{Code: "LBP", Name: "Lebanese Pound", NumericCode: "422", MinorUnits: 2, Symbol: "ل.ل"},
	{Code: "LKR", Name: "Sri Lanka Rupee", NumericCode: "144", MinorUnits: 2, Symbol: "Rs"},
	{Code: "LRD", Name: "Liberian Dollar", NumericCode: "430", MinorUnits: 2, Symbol: "$"},
	{Code: "LSL", Name: "Loti", NumericCode: "426", MinorUnits: 2, Symbol: "M"},
	{Code: "LYD", Name: "Libyan Dinar", NumericCode: "434", MinorUnits: 3, Symbol: "ل.د"},
	{Code: "MAD", Name: "Moroccan Dirham", NumericCode: "504", MinorUnits: 2, Symbol: "د.م."},
	{Code: "MDL", Name: "Moldovan Leu", NumericCode: "498", MinorUnits: 2, Symbol: "L"},
	{Code: "MGA", Name: "Malagasy Ariary", NumericCode: "969", MinorUnits: 2, Symbol: "Ar"},
	{Code: "MKD", Name: "Denar", NumericCode: "807", MinorUnits: 2, Symbol: "ден"},
	{Code: "MMK", Name: "Kyat", NumericCode: "104", MinorUnits: 2, Symbol: "K"},
	{Code: "MNT", Name: "Tugrik", NumericCode: "496", MinorUnits: 2, Symbol: "₮"},
	{Code: "MOP", Name: "Pataca", NumericCode: "446", MinorUnits: 2, Symbol: "MOP$"},
	{Code: "MRU", Name: "Ouguiya", NumericCode: "929", MinorUnits: 2, Symbol: "UM"},
	{Code: "MUR", Name: "Mauritius Rupee", NumericCode: "480", MinorUnits: 2, Symbol: "₨"},
	{Code: "MVR", Name: "Rufiyaa", NumericCode: "462", MinorUnits: 2, Symbol: "Rf"},
	{Code: "MWK", Name: "Malawi Kwacha", NumericCode: "454", MinorUnits: 2, Symbol: "MK"},
	{Code: "MXN", Name: "Mexican Peso", NumericCode: "484", MinorUnits: 2, Symbol: "$"},
	{Code: "MXV", Name: "Mexican Unidad de Inversion", NumericCode: "979", MinorUnits: 2, Symbol: ""},
	{Code: "MYR", Name: "Malaysian Ringgit", NumericCode: "458", MinorUnits: 2, Symbol: "RM"},
	{Code: "MZN", Name: "Mozambique Metical", NumericCode: "943", MinorUnits: 2, Symbol: "MT"},
	{Code: "NAD", Name: "Namibia Dollar", NumericCode: "516", MinorUnits: 2, Symbol: "$"},
	{Code: "NGN", Name: "Naira", NumericCode: "566", MinorUnits: 2, Symbol: "₦"},
	{Code: "NIO", Name: "Cordoba Oro", NumericCode: "558", MinorUnits: 2, Symbol: "C$"},
	{Code: "NOK", Name: "Norwegian Krone", NumericCode: "578", MinorUnits: 2, Symbol: "kr"},
	{Code: "NPR", Name: "Nepalese Rupee", NumericCode: "524", MinorUnits: 2, Symbol: "₨"},
	{Code: "NZD", Name: "New Zealand Dollar", NumericCode: "554", MinorUnits: 2, Symbol: "NZ$"},
	{Code: "OMR", Name: "Rial Omani", NumericCode: "512", MinorUnits: 3, Symbol: "ر.ع."},
	{Code: "PAB", Name: "Balboa", NumericCode: "590", MinorUnits: 2, Symbol: "B/."},
	{Code: "PEN", Name: "Sol", NumericCode: "604", MinorUnits: 2, Symbol: "S/"},
	{Code: "PGK", Name: "Kina", NumericCode: "598", MinorUnits: 2, Symbol: "K"},
	{Code: "PHP", Name: "Philippine Peso", NumericCode: "608", MinorUnits: 2, Symbol: "₱"},
	{Code: "PKR", Name: "Pakistan Rupee", NumericCode: "586", MinorUnits: 2, Symbol: "₨"},
	{Code: "PLN", Name: "Zloty", NumericCode: "985", MinorUnits: 2, Symbol: "zł"},
	{Code: "PYG", Name: "Guarani", NumericCode: "600", MinorUnits: 0, Symbol: "₲"},
	{Code: "QAR", Name: "Qatari Rial", NumericCode: "634", MinorUnits: 2, Symbol: "ر.ق"},
	{Code: "RON", Name: "Romanian Leu", NumericCode: "946", MinorUnits: 2, Symbol: "lei"},
	{Code: "RSD", Name: "Serbian Dinar", NumericCode: "941", MinorUnits: 2, Symbol: "дин."},
	{Code: "RUB", Name: "Russian Ruble", NumericCode: "643", MinorUnits: 2, Symbol: "₽"},
	{Code: "RWF", Name: "Rwanda Franc", NumericCode: "646", MinorUnits: 0, Symbol: "FRw"},
	{Code: "SAR", Name: "Saudi Riyal", NumericCode: "682", MinorUnits: 2, Symbol: "ر.س"},
	{Code: "SBD", Name: "Solomon Islands Dollar", NumericCode: "090", MinorUnits: 2, Symbol: "$"},
	{Code: "SCR", Name: "Seychelles Rupee", NumericCode: "690", MinorUnits: 2, Symbol: "₨"},
	{Code: "SDG", Name: "Sudanese Pound", NumericCode: "938", MinorUnits: 2, Symbol: "ج.س."},
	{Code: "SEK", Name: "Swedish Krona", NumericCode: "752", MinorUnits: 2, Symbol: "kr"},
	{Code: "SGD", Name: "Singapore Dollar", NumericCode: "702", MinorUnits: 2, Symbol: "S$"},
	{Code: "SHP", Name: "Saint Helena Pound", NumericCode: "654", MinorUnits: 2, Symbol: "£"},
	{Code: "SLE", Name: "Leone", NumericCode: "925", MinorUnits: 2, Symbol: "Le"},
	{Code: "SOS", Name: "Somali Shilling", NumericCode: "706", MinorUnits: 2, Symbol: "S"},
	{Code: "SRD", Name: "Surinam Dollar", NumericCode: "968", MinorUnits: 2, Symbol: "$"},
	{Code: "SSP", Name: "South Sudanese Pound", NumericCode: "728", MinorUnits: 2, Symbol: "£"},
	{Code: "STN", Name: "Dobra", NumericCode: "930", MinorUnits: 2, Symbol: "Db"},
	{Code: "SVC", Name: "El Salvador Colon", NumericCode: "222", MinorUnits: 2, Symbol: "$"},
	{Code: "SYP", Name: "Syrian Pound", NumericCode: "760", MinorUnits: 2, Symbol: "£"},
	{Code: "SZL", Name: "Lilangeni", NumericCode: "748", MinorUnits: 2, Symbol: "E"},
	{Code: "THB", Name: "Baht", NumericCode: "764", MinorUnits: 2, Symbol: "฿"},
	{Code: "TJS", Name: "Somoni", NumericCode: "972", MinorUnits: 2, Symbol: "SM"},
	{Code: "TMT", Name: "Turkmenistan New Manat", NumericCode: "934", MinorUnits: 2, Symbol: "T"},
	{Code: "TND", Name: "Tunisian Dinar", NumericCode: "788", MinorUnits: 3, Symbol: "د.ت"},
	{Code: "TOP", Name: "Pa'anga", NumericCode: "776", MinorUnits: 2, Symbol: "T$"},
	{Code: "TRY", Name: "Turkish Lira", NumericCode: "949", MinorUnits: 2, Symbol: "₺"},
	{Code: "TTD", Name: "Trinidad and Tobago Dollar", NumericCode: "780", MinorUnits: 2, Symbol: "TT$"},
	{Code: "TWD", Name: "New Taiwan Dollar", NumericCode: "901", MinorUnits: 2, Symbol: "NT$"},
	{Code: "TZS", Name: "Tanzanian Shilling", NumericCode: "834", MinorUnits: 2, Symbol: "TSh"},
	{Code: "UAH", Name: "Hryvnia", NumericCode: "980", MinorUnits: 2, Symbol: "₴"},
	{Code: "UGX", Name: "Uganda Shilling", NumericCode: "800", MinorUnits: 0, Symbol: "USh"},
	{Code: "USD", Name: "US Dollar", NumericCode: "840", MinorUnits: 2, Symbol: "$"},
	{Code: "USN", Name: "US Dollar (Next day)", NumericCode: "997", MinorUnits: 2, Symbol: "$"},
	{Code: "UYI", Name: "Uruguay Peso en Unidades Indexadas", NumericCode: "940", MinorUnits: 0, Symbol: ""},
	{Code: "UYU", Name: "Peso Uruguayo", NumericCode: "858", MinorUnits: 2, Symbol: "$U"},
	{Code: "UYW", Name: "Unidad Previsional", NumericCode: "927", MinorUnits: 4, Symbol: ""},
	{Code: "UZS", Name: "Uzbekistan Sum", NumericCode: "860", MinorUnits: 2, Symbol: "UZS"},
	{Code: "VED", Name: "Bolivar Soberano", NumericCode: "926", MinorUnits: 2, Symbol: "Bs.D"},
	{Code: "VES", Name: "Bolivar Soberano", NumericCode: "928", MinorUnits: 2, Symbol: "Bs.S"},
	{Code: "VND", Name: "Dong", NumericCode: "704", MinorUnits: 0, Symbol: "₫"},
	{Code: "VUV", Name: "Vatu", NumericCode: "548", MinorUnits: 0, Symbol: "VT"},
	{Code: "WST", Name: "Tala", NumericCode: "882", MinorUnits: 2, Symbol: "WS$"},
	{Code: "XAF", Name: "CFA Franc BEAC", NumericCode: "950", MinorUnits: 0, Symbol: "FCFA"},
	{Code: "XAG", Name: "Silver", NumericCode: "961", MinorUnits: 0, Symbol: ""},
	{Code: "XAU", Name: "Gold", NumericCode: "959", MinorUnits: 0, Symbol: ""},
	{Code: "XBA", Name: "Bond Markets Unit European Composite Unit (EURCO)", NumericCode: "955", MinorUnits: 0, Symbol: ""},
	{Code: "XBB", Name: "Bond Markets Unit European Monetary Unit (E.M.U.-6)", NumericCode: "956", MinorUnits: 0, Symbol: ""},
	{Code: "XBC", Name: "Bond Markets Unit European Unit of Account 9 (E.U.A.-9)", NumericCode: "957", MinorUnits: 0, Symbol: ""},
	{Code: "XBD", Name: "Bond Markets Unit European Unit of Account 17 (E.U.A.-17)", NumericCode: "958", MinorUnits: 0, Symbol: ""},
	{Code: "XCD", Name: "East Caribbean Dollar", NumericCode: "951", MinorUnits: 2, Symbol: "$"},
	{Code: "XDR", Name: "SDR (Special Drawing Right)", NumericCode: "960", MinorUnits: 0, Symbol: ""},
	{Code: "XOF", Name: "CFA Franc BCEAO", NumericCode: "952", MinorUnits: 0, Symbol: "CFA"},
	{Code: "XPD", Name: "Palladium", NumericCode: "964", MinorUnits: 0, Symbol: ""},
	{Code: "XPF", Name: "CFP Franc", NumericCode: "953", MinorUnits: 0, Symbol: "₣"},
	{Code: "XPT", Name: "Platinum", NumericCode: "962", MinorUnits: 0, Symbol: ""},
	{Code: "XSU", Name: "Sucre", NumericCode: "994", MinorUnits: 0, Symbol: ""},
	{Code: "XTS", Name: "Codes specifically reserved for testing purposes", NumericCode: "963", MinorUnits: 0, Symbol: ""},
	{Code: "XUA", Name: "ADB Unit of Account", NumericCode: "965", MinorUnits: 0, Symbol: ""},
	{Code: "XXX", Name: "No currency", NumericCode: "999", MinorUnits: 0, Symbol: ""},
	{Code: "YER", Name: "Yemeni Rial", NumericCode: "886", MinorUnits: 2, Symbol: "﷼"},
	{Code: "ZAR", Name: "Rand", NumericCode: "710", MinorUnits: 2, Symbol: "R"},
	{Code: "ZMW", Name: "Zambian Kwacha", NumericCode: "967", MinorUnits: 2, Symbol: "ZK"},
	{Code: "ZWL", Name: "Zimbabwe Dollar", NumericCode: "932", MinorUnits: 2, Symbol: "Z$"},
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
//	fmt.Printf("%s has %d decimal places\n", c.Name, c.MinorUnits)
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
//	    fmt.Printf("%s: %s (%d decimals)\n", c.Code, c.Name, c.MinorUnits)
//	}
func AllCurrencies() []Currency {
	return allCurrencies
}
