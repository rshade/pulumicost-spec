package pluginsdk

import (
	"time"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FocusRecordBuilder handles the construction of FOCUS 1.2 cost records.
type FocusRecordBuilder struct {
	record *pbc.FocusCostRecord
}

// NewFocusRecordBuilder creates a new builder instance.
func NewFocusRecordBuilder() *FocusRecordBuilder {
	return &FocusRecordBuilder{
		record: &pbc.FocusCostRecord{
			ExtendedColumns: make(map[string]string),
			Tags:            make(map[string]string),
		},
	}
}

// WithIdentity sets the identity and hierarchy fields.
func (b *FocusRecordBuilder) WithIdentity(
	providerName, billingAccountID, billingAccountName string,
) *FocusRecordBuilder {
	b.record.ProviderName = providerName
	b.record.BillingAccountId = billingAccountID
	b.record.BillingAccountName = billingAccountName
	return b
}

// WithChargePeriod sets the start and end time of the charge.
func (b *FocusRecordBuilder) WithChargePeriod(start, end time.Time) *FocusRecordBuilder {
	b.record.ChargePeriodStart = timestamppb.New(start)
	b.record.ChargePeriodEnd = timestamppb.New(end)
	return b
}

// WithServiceCategory sets the service category.
func (b *FocusRecordBuilder) WithServiceCategory(category pbc.FocusServiceCategory) *FocusRecordBuilder {
	b.record.ServiceCategory = category
	return b
}

// WithChargeDetails sets the charge and pricing categories.
func (b *FocusRecordBuilder) WithChargeDetails(
	chargeCat pbc.FocusChargeCategory,
	pricingCat pbc.FocusPricingCategory,
) *FocusRecordBuilder {
	b.record.ChargeCategory = chargeCat
	b.record.PricingCategory = pricingCat
	return b
}

// WithFinancials sets the cost amounts and currency.
func (b *FocusRecordBuilder) WithFinancials(
	billed, list, effective float64,
	currency, invoiceID string,
) *FocusRecordBuilder {
	b.record.BilledCost = billed
	b.record.ListCost = list
	b.record.EffectiveCost = effective
	b.record.Currency = currency
	b.record.InvoiceId = invoiceID
	return b
}

// WithUsage sets the usage quantity and unit.
func (b *FocusRecordBuilder) WithUsage(quantity float64, unit string) *FocusRecordBuilder {
	b.record.UsageQuantity = quantity
	b.record.UsageUnit = unit
	return b
}

// WithExtension adds a key-value pair to the extended columns (Backpack).
func (b *FocusRecordBuilder) WithExtension(key, value string) *FocusRecordBuilder {
	b.record.ExtendedColumns[key] = value
	return b
}

// Build validates and returns the constructed FocusCostRecord.
func (b *FocusRecordBuilder) Build() (*pbc.FocusCostRecord, error) {
	if err := ValidateFocusRecord(b.record); err != nil {
		return nil, err
	}
	return b.record, nil
}
