package pluginsdk

import (
	"errors"
	"fmt"
	"math"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// iso4217Currencies contains valid ISO 4217 currency codes.
// This is a comprehensive set of active ISO 4217 currency codes commonly
// encountered in cloud billing across major providers (AWS, Azure, GCP).
//
//nolint:gochecknoglobals // Package-level constant data for currency validation.
var iso4217Currencies = map[string]bool{
	"AED": true, "AFN": true, "ALL": true, "AMD": true, "ANG": true,
	"AOA": true, "ARS": true, "AUD": true, "AWG": true, "AZN": true,
	"BAM": true, "BBD": true, "BDT": true, "BGN": true, "BHD": true,
	"BIF": true, "BMD": true, "BND": true, "BOB": true, "BRL": true,
	"BSD": true, "BTN": true, "BWP": true, "BYN": true, "BZD": true,
	"CAD": true, "CDF": true, "CHF": true, "CLP": true, "CNY": true,
	"COP": true, "CRC": true, "CUP": true, "CVE": true, "CZK": true,
	"DJF": true, "DKK": true, "DOP": true, "DZD": true, "EGP": true,
	"ERN": true, "ETB": true, "EUR": true, "FJD": true, "FKP": true,
	"GBP": true, "GEL": true, "GHS": true, "GIP": true, "GMD": true,
	"GNF": true, "GTQ": true, "GYD": true, "HKD": true, "HNL": true,
	"HRK": true, "HTG": true, "HUF": true, "IDR": true, "ILS": true,
	"INR": true, "IQD": true, "IRR": true, "ISK": true, "JMD": true,
	"JOD": true, "JPY": true, "KES": true, "KGS": true, "KHR": true,
	"KMF": true, "KPW": true, "KRW": true, "KWD": true, "KYD": true,
	"KZT": true, "LAK": true, "LBP": true, "LKR": true, "LRD": true,
	"LSL": true, "LYD": true, "MAD": true, "MDL": true, "MGA": true,
	"MKD": true, "MMK": true, "MNT": true, "MOP": true, "MRU": true,
	"MUR": true, "MVR": true, "MWK": true, "MXN": true, "MYR": true,
	"MZN": true, "NAD": true, "NGN": true, "NIO": true, "NOK": true,
	"NPR": true, "NZD": true, "OMR": true, "PAB": true, "PEN": true,
	"PGK": true, "PHP": true, "PKR": true, "PLN": true, "PYG": true,
	"QAR": true, "RON": true, "RSD": true, "RUB": true, "RWF": true,
	"SAR": true, "SBD": true, "SCR": true, "SDG": true, "SEK": true,
	"SGD": true, "SHP": true, "SLE": true, "SOS": true, "SRD": true,
	"SSP": true, "STN": true, "SVC": true, "SYP": true, "SZL": true,
	"THB": true, "TJS": true, "TMT": true, "TND": true, "TOP": true,
	"TRY": true, "TTD": true, "TWD": true, "TZS": true, "UAH": true,
	"UGX": true, "USD": true, "UYU": true, "UZS": true, "VES": true,
	"VND": true, "VUV": true, "WST": true, "XAF": true, "XCD": true,
	"XOF": true, "XPF": true, "YER": true, "ZAR": true, "ZMW": true,
	"ZWL": true,
}

// contractedCostTolerance is the relative tolerance for ContractedCost validation.
// Allows for floating-point precision differences up to 0.01%.
const contractedCostTolerance = 0.0001

// ValidateFocusRecord checks if a record complies with FOCUS 1.2 mandatory fields and business rules.
// Reference: https://focus.finops.org
func ValidateFocusRecord(r *pbc.FocusCostRecord) error {
	if r == nil {
		return errors.New("record is nil")
	}

	// Validate mandatory fields (FOCUS 1.2).
	if err := validateMandatoryFields(r); err != nil {
		return err
	}

	// Validate currency codes (ISO 4217).
	if err := validateCurrencyFields(r); err != nil {
		return err
	}

	// Validate business rules.
	if err := validateBusinessRules(r); err != nil {
		return err
	}

	return nil
}

// validateMandatoryFields checks all 14 mandatory FOCUS 1.2 fields.
func validateMandatoryFields(r *pbc.FocusCostRecord) error {
	// Identity fields (FOCUS 1.2 Section 2.1).
	if r.GetProviderName() == "" {
		return errors.New("provider_name is required")
	}
	if r.GetBillingAccountId() == "" {
		return errors.New("billing_account_id is required")
	}

	// Billing period (FOCUS 1.2 Section 2.2).
	if r.GetBillingPeriodStart() == nil || r.GetBillingPeriodEnd() == nil {
		return errors.New("billing_period (start/end) is required")
	}
	if r.GetBillingCurrency() == "" {
		return errors.New("billing_currency is required")
	}

	// Charge period (FOCUS 1.2 Section 2.3).
	if r.GetChargePeriodStart() == nil || r.GetChargePeriodEnd() == nil {
		return errors.New("charge_period (start/end) is required")
	}

	// Charge details (FOCUS 1.2 Section 2.4).
	if r.GetChargeCategory() == pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_UNSPECIFIED {
		return errors.New("charge_category is required")
	}
	if r.GetChargeClass() == pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_UNSPECIFIED {
		return errors.New("charge_class is required")
	}
	if r.GetChargeDescription() == "" {
		return errors.New("charge_description is required")
	}

	// Service details (FOCUS 1.2 Section 2.6).
	if r.GetServiceCategory() == pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_UNSPECIFIED {
		return errors.New("service_category is required")
	}
	if r.GetServiceName() == "" {
		return errors.New("service_name is required")
	}

	// Note: BilledCost and ContractedCost are mandatory but can be zero or negative
	// (e.g., credits, refunds per FOCUS 1.2 spec).
	//
	// Note: BillingAccountName is mandatory per FOCUS 1.2 but validation is intentionally
	// relaxed to accommodate cloud providers that may not include account names in all
	// billing data exports. Plugins should populate this field when available.
	return nil
}

// validateCurrencyFields validates ISO 4217 currency codes.
func validateCurrencyFields(r *pbc.FocusCostRecord) error {
	if err := validateCurrency(r.GetBillingCurrency(), "billing_currency"); err != nil {
		return err
	}

	// PricingCurrency is optional, but if present, must be valid ISO 4217.
	if r.GetPricingCurrency() != "" {
		if err := validateCurrency(r.GetPricingCurrency(), "pricing_currency"); err != nil {
			return err
		}
	}
	return nil
}

// validateBusinessRules validates FOCUS 1.2 business rules.
func validateBusinessRules(r *pbc.FocusCostRecord) error {
	// Rule: Usage records must have positive consumed quantity.
	// FOCUS 1.2: If ChargeCategory=Usage, ConsumedQuantity must be > 0.
	if r.GetChargeCategory() == pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE {
		if r.GetConsumedQuantity() <= 0 {
			return errors.New("consumed_quantity must be positive for usage charge category")
		}
	}

	// Rule: ContractedCost must equal ContractedUnitPrice × PricingQuantity.
	// FOCUS 1.2 Section 3.20: When ContractedUnitPrice and PricingQuantity are present
	// and ChargeClass is not "Correction", this relationship must hold.
	if err := validateContractedCostRule(r); err != nil {
		return err
	}

	return nil
}

// validateCurrency checks if a currency code is a valid ISO 4217 code.
func validateCurrency(code string, fieldName string) error {
	if !iso4217Currencies[code] {
		return fmt.Errorf("%s must be a valid ISO 4217 currency code, got %q", fieldName, code)
	}
	return nil
}

// validateContractedCostRule verifies ContractedCost = ContractedUnitPrice × PricingQuantity.
// This rule applies when both ContractedUnitPrice and PricingQuantity are present (non-zero)
// and ChargeClass is not Correction.
func validateContractedCostRule(r *pbc.FocusCostRecord) error {
	// Skip validation for correction charges.
	if r.GetChargeClass() == pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_CORRECTION {
		return nil
	}

	contractedUnitPrice := r.GetContractedUnitPrice()
	pricingQuantity := r.GetPricingQuantity()

	// Only validate when both values are present (non-zero).
	if contractedUnitPrice == 0 || pricingQuantity == 0 {
		return nil
	}

	contractedCost := r.GetContractedCost()
	expectedCost := contractedUnitPrice * pricingQuantity

	// Use relative tolerance for floating-point comparison.
	if !floatEquals(contractedCost, expectedCost, contractedCostTolerance) {
		return fmt.Errorf(
			"contracted_cost (%f) must equal contracted_unit_price (%f) × pricing_quantity (%f) = %f",
			contractedCost, contractedUnitPrice, pricingQuantity, expectedCost,
		)
	}

	return nil
}

// floatEquals compares two floats with relative tolerance.
func floatEquals(a, b, tolerance float64) bool {
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	largest := math.Max(math.Abs(a), math.Abs(b))
	return diff <= largest*tolerance
}
