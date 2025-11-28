package pluginsdk

import (
	"errors"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// ValidateFocusRecord checks if a record complies with FOCUS 1.2 mandatory fields and business rules.
func ValidateFocusRecord(r *pbc.FocusCostRecord) error {
	if r == nil {
		return errors.New("record is nil")
	}

	// Mandatory fields
	if r.GetBillingAccountId() == "" {
		return errors.New("billing_account_id is required")
	}
	if r.GetChargePeriodStart() == nil || r.GetChargePeriodEnd() == nil {
		return errors.New("charge_period (start/end) is required")
	}
	if r.GetServiceCategory() == pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_UNSPECIFIED {
		return errors.New("service_category is required")
	}
	if r.GetChargeCategory() == pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_UNSPECIFIED {
		return errors.New("charge_category is required")
	}
	if r.GetCurrency() == "" {
		return errors.New("currency is required")
	}

	// Business Rules
	//
	// Rule: Usage records should have usage quantity (unless it's a fixed fee, but usually usage implies quantity)
	// The spec example said "if ChargeCategory=Usage, UsageQuantity must be > 0".
	if r.GetChargeCategory() == pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE {
		if r.GetUsageQuantity() == 0 {
			// 0 usage for usage charge seems wrong.
			return errors.New("usage_quantity must be non-zero for usage charge category")
		}
	}

	return nil
}
