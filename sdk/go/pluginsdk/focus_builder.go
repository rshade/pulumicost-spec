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

// WithIdentity sets the identity and hierarchy fields per FOCUS 1.2 Section 2.1.
func (b *FocusRecordBuilder) WithIdentity(
	providerName, billingAccountID, billingAccountName string,
) *FocusRecordBuilder {
	b.record.ProviderName = providerName
	b.record.BillingAccountId = billingAccountID
	b.record.BillingAccountName = billingAccountName
	return b
}

// WithSubAccount sets the sub-account fields (e.g., AWS Account, Azure Subscription, GCP Project) per FOCUS 1.2 Section 2.1.
func (b *FocusRecordBuilder) WithSubAccount(subAccountID, subAccountName string) *FocusRecordBuilder {
	b.record.SubAccountId = subAccountID
	b.record.SubAccountName = subAccountName
	return b
}

// WithBillingPeriod sets the billing period start/end and currency per FOCUS 1.2 Section 2.2.
func (b *FocusRecordBuilder) WithBillingPeriod(start, end time.Time, currency string) *FocusRecordBuilder {
	b.record.BillingPeriodStart = timestamppb.New(start)
	b.record.BillingPeriodEnd = timestamppb.New(end)
	b.record.BillingCurrency = currency
	return b
}

// WithChargePeriod sets the start and end time of the charge per FOCUS 1.2 Section 2.3.
func (b *FocusRecordBuilder) WithChargePeriod(start, end time.Time) *FocusRecordBuilder {
	b.record.ChargePeriodStart = timestamppb.New(start)
	b.record.ChargePeriodEnd = timestamppb.New(end)
	return b
}

// WithServiceCategory sets the service category per FOCUS 1.2 Section 2.6.
func (b *FocusRecordBuilder) WithServiceCategory(category pbc.FocusServiceCategory) *FocusRecordBuilder {
	b.record.ServiceCategory = category
	return b
}

// WithChargeDetails sets the charge category and pricing category per FOCUS 1.2 Section 2.4.
func (b *FocusRecordBuilder) WithChargeDetails(
	chargeCat pbc.FocusChargeCategory,
	pricingCat pbc.FocusPricingCategory,
) *FocusRecordBuilder {
	b.record.ChargeCategory = chargeCat
	b.record.PricingCategory = pricingCat
	return b
}

// WithChargeClassification sets the charge class, description, and frequency per FOCUS 1.2 Section 2.4.
func (b *FocusRecordBuilder) WithChargeClassification(
	chargeClass pbc.FocusChargeClass,
	description string,
	frequency pbc.FocusChargeFrequency,
) *FocusRecordBuilder {
	b.record.ChargeClass = chargeClass
	b.record.ChargeDescription = description
	b.record.ChargeFrequency = frequency
	return b
}

// WithPricing sets the pricing quantity, unit, and list unit price per FOCUS 1.2 Section 2.5.
func (b *FocusRecordBuilder) WithPricing(quantity float64, unit string, listUnitPrice float64) *FocusRecordBuilder {
	b.record.PricingQuantity = quantity
	b.record.PricingUnit = unit
	b.record.ListUnitPrice = listUnitPrice
	return b
}

// WithFinancials sets the cost amounts and currency per FOCUS 1.2 Section 2.10.
func (b *FocusRecordBuilder) WithFinancials(
	billed, list, effective float64,
	currency, invoiceID string,
) *FocusRecordBuilder {
	b.record.BilledCost = billed
	b.record.ListCost = list
	b.record.EffectiveCost = effective
	b.record.BillingCurrency = currency
	b.record.InvoiceId = invoiceID
	return b
}

// WithUsage sets the consumed quantity and unit per FOCUS 1.2 Section 2.11.
func (b *FocusRecordBuilder) WithUsage(quantity float64, unit string) *FocusRecordBuilder {
	b.record.ConsumedQuantity = quantity
	b.record.ConsumedUnit = unit
	return b
}

// WithResource sets the resource details per FOCUS 1.2 Section 2.7.
func (b *FocusRecordBuilder) WithResource(resourceID, resourceName, resourceType string) *FocusRecordBuilder {
	b.record.ResourceId = resourceID
	b.record.ResourceName = resourceName
	b.record.ResourceType = resourceType
	return b
}

// WithService sets the service details per FOCUS 1.2 Section 2.6.
func (b *FocusRecordBuilder) WithService(
	category pbc.FocusServiceCategory,
	serviceName string,
) *FocusRecordBuilder {
	b.record.ServiceCategory = category
	b.record.ServiceName = serviceName
	return b
}

// WithSKU sets the SKU details per FOCUS 1.2 Section 2.8.
func (b *FocusRecordBuilder) WithSKU(skuID, skuPriceID string) *FocusRecordBuilder {
	b.record.SkuId = skuID
	b.record.SkuPriceId = skuPriceID
	return b
}

// WithLocation sets the region and availability zone per FOCUS 1.2 Section 2.9.
func (b *FocusRecordBuilder) WithLocation(regionID, regionName, availabilityZone string) *FocusRecordBuilder {
	b.record.RegionId = regionID
	b.record.RegionName = regionName
	b.record.AvailabilityZone = availabilityZone
	return b
}

// WithCommitmentDiscount sets the commitment discount details per FOCUS 1.2 Section 2.12.
func (b *FocusRecordBuilder) WithCommitmentDiscount(
	category pbc.FocusCommitmentDiscountCategory,
	discountID, discountName string,
) *FocusRecordBuilder {
	b.record.CommitmentDiscountCategory = category
	b.record.CommitmentDiscountId = discountID
	b.record.CommitmentDiscountName = discountName
	return b
}

// WithInvoice sets the invoice details per FOCUS 1.2 Section 2.13.
func (b *FocusRecordBuilder) WithInvoice(invoiceID, invoiceIssuer string) *FocusRecordBuilder {
	b.record.InvoiceId = invoiceID
	b.record.InvoiceIssuer = invoiceIssuer
	return b
}

// WithTag adds a single tag key-value pair per FOCUS 1.2 Section 2.14.
func (b *FocusRecordBuilder) WithTag(key, value string) *FocusRecordBuilder {
	b.record.Tags[key] = value
	return b
}

// WithTags sets multiple tags at once per FOCUS 1.2 Section 2.14.
func (b *FocusRecordBuilder) WithTags(tags map[string]string) *FocusRecordBuilder {
	for k, v := range tags {
		b.record.Tags[k] = v
	}
	return b
}

// WithExtension adds a key-value pair to the extended columns (Backpack) per FOCUS 1.2 Section 2.14.
func (b *FocusRecordBuilder) WithExtension(key, value string) *FocusRecordBuilder {
	b.record.ExtendedColumns[key] = value
	return b
}

// =============================================================================
// FOCUS 1.2 New Column Builder Methods (19 new columns)
// =============================================================================

// WithContractedCost sets the contracted cost per FOCUS 1.2 Section 3.20.
// This is a MANDATORY field representing the cost calculated by multiplying
// contracted unit price and pricing quantity.
func (b *FocusRecordBuilder) WithContractedCost(cost float64) *FocusRecordBuilder {
	b.record.ContractedCost = cost
	return b
}

// WithBillingAccountType sets the billing account type per FOCUS 1.2 Section 3.3.
// This is a CONDITIONAL field representing the provider-assigned name to identify
// the type of billing account.
func (b *FocusRecordBuilder) WithBillingAccountType(accountType string) *FocusRecordBuilder {
	b.record.BillingAccountType = accountType
	return b
}

// WithSubAccountType sets the sub-account type per FOCUS 1.2 Section 3.45.
// This is a CONDITIONAL field representing the provider-assigned identifier
// for sub-account classification.
func (b *FocusRecordBuilder) WithSubAccountType(accountType string) *FocusRecordBuilder {
	b.record.SubAccountType = accountType
	return b
}

// WithCapacityReservation sets the capacity reservation details per FOCUS 1.2 Sections 3.6, 3.7.
// This includes the capacity reservation ID and its utilization status.
// Both fields are CONDITIONAL.
func (b *FocusRecordBuilder) WithCapacityReservation(
	reservationID string,
	status pbc.FocusCapacityReservationStatus,
) *FocusRecordBuilder {
	b.record.CapacityReservationId = reservationID
	b.record.CapacityReservationStatus = status
	return b
}

// WithCommitmentDiscountDetails sets extended commitment discount details per FOCUS 1.2.
// Includes quantity (Section 3.14), status (Section 3.17), type (Section 3.18),
// and unit (Section 3.19). All fields are CONDITIONAL.
func (b *FocusRecordBuilder) WithCommitmentDiscountDetails(
	quantity float64,
	status pbc.FocusCommitmentDiscountStatus,
	discountType string,
	unit string,
) *FocusRecordBuilder {
	b.record.CommitmentDiscountQuantity = quantity
	b.record.CommitmentDiscountStatus = status
	b.record.CommitmentDiscountType = discountType
	b.record.CommitmentDiscountUnit = unit
	return b
}

// WithContractedUnitPrice sets the contracted unit price per FOCUS 1.2 Section 3.21.
// This is a CONDITIONAL field representing the agreed-upon unit price per
// pricing unit for the associated SKU.
func (b *FocusRecordBuilder) WithContractedUnitPrice(price float64) *FocusRecordBuilder {
	b.record.ContractedUnitPrice = price
	return b
}

// WithPricingCurrency sets the pricing currency per FOCUS 1.2 Section 3.34.
// This is a CONDITIONAL field used when pricing is in a different currency
// than billing. Format: ISO 4217 currency code.
func (b *FocusRecordBuilder) WithPricingCurrency(currency string) *FocusRecordBuilder {
	b.record.PricingCurrency = currency
	return b
}

// WithPricingCurrencyPrices sets pricing currency amounts per FOCUS 1.2.
// Includes contracted unit price (Section 3.35), effective cost (Section 3.36),
// and list unit price (Section 3.37). All fields are CONDITIONAL.
func (b *FocusRecordBuilder) WithPricingCurrencyPrices(
	contractedUnitPrice float64,
	effectiveCost float64,
	listUnitPrice float64,
) *FocusRecordBuilder {
	b.record.PricingCurrencyContractedUnitPrice = contractedUnitPrice
	b.record.PricingCurrencyEffectiveCost = effectiveCost
	b.record.PricingCurrencyListUnitPrice = listUnitPrice
	return b
}

// WithPublisher sets the publisher per FOCUS 1.2 Section 3.39.
// This is a CONDITIONAL field representing the entity that published the
// service or product.
func (b *FocusRecordBuilder) WithPublisher(publisher string) *FocusRecordBuilder {
	b.record.Publisher = publisher
	return b
}

// WithServiceSubcategory sets the service subcategory per FOCUS 1.2 Section 3.43.
// This is a CONDITIONAL field providing granular service classification
// supporting functional categorization.
func (b *FocusRecordBuilder) WithServiceSubcategory(subcategory string) *FocusRecordBuilder {
	b.record.ServiceSubcategory = subcategory
	return b
}

// WithSkuDetails sets SKU meter and price details per FOCUS 1.2 Sections 3.46, 3.48.
// SkuMeter is the provider-assigned meter identifier, and SkuPriceDetails contains
// additional pricing information. Both fields are CONDITIONAL.
func (b *FocusRecordBuilder) WithSkuDetails(meter, priceDetails string) *FocusRecordBuilder {
	b.record.SkuMeter = meter
	b.record.SkuPriceDetails = priceDetails
	return b
}

// Build validates and returns the constructed FocusCostRecord.
func (b *FocusRecordBuilder) Build() (*pbc.FocusCostRecord, error) {
	if err := ValidateFocusRecord(b.record); err != nil {
		return nil, err
	}
	return b.record, nil
}
