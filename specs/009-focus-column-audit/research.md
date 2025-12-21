# Research: FOCUS 1.2 Column Audit

**Date**: 2025-11-28
**Feature**: 010-focus-column-audit
**Source**: [FOCUS Specification v1.2](https://focus.finops.org/focus-specification/v1-2/)

## Research Summary

All NEEDS CLARIFICATION items have been resolved through research of the official
FOCUS 1.2 specification. This document captures the decisions and rationale for
implementing the 19 missing columns.

## Missing Column Specifications

### 1. ContractedCost (MANDATORY)

**Decision**: Add as `double contracted_cost` field
**Rationale**: Required for FOCUS compliance; calculated as ContractedUnitPrice ×
PricingQuantity
**Alternatives Considered**: None - mandatory column

| Attribute     | Value                                                                     |
| ------------- | ------------------------------------------------------------------------- |
| Data Type     | Decimal → `double` in proto                                               |
| Feature Level | Mandatory                                                                 |
| Allows Nulls  | No                                                                        |
| Description   | Cost calculated by multiplying contracted unit price and pricing quantity |

### 2. CommitmentDiscountStatus (CONDITIONAL)

**Decision**: Add as enum `FocusCommitmentDiscountStatus` with values: Used, Unused
**Rationale**: Limited set of allowed values suits enum pattern
**Alternatives Considered**: String field - rejected for type safety

| Attribute      | Value                                                                   |
| -------------- | ----------------------------------------------------------------------- |
| Data Type      | String (enum) → `FocusCommitmentDiscountStatus`                         |
| Feature Level  | Conditional                                                             |
| Allowed Values | Used, Unused                                                            |
| Constraint     | Required when CommitmentDiscountId not null and ChargeCategory is Usage |

### 3. CommitmentDiscountType (CONDITIONAL)

**Decision**: Add as `string commitment_discount_type` field
**Rationale**: Provider-assigned identifier with no fixed values
**Alternatives Considered**: None - string is appropriate

| Attribute     | Value                                                            |
| ------------- | ---------------------------------------------------------------- |
| Data Type     | String                                                           |
| Feature Level | Conditional                                                      |
| Description   | Provider-assigned identifier for the type of commitment discount |

### 4. CommitmentDiscountQuantity (CONDITIONAL)

**Decision**: Add as `double commitment_discount_quantity` field
**Rationale**: Decimal type for quantity values
**Alternatives Considered**: None

| Attribute     | Value                                                    |
| ------------- | -------------------------------------------------------- |
| Data Type     | Decimal → `double`                                       |
| Feature Level | Conditional                                              |
| Description   | Amount of commitment discount purchased or accounted for |

### 5. CommitmentDiscountUnit (CONDITIONAL)

**Decision**: Add as `string commitment_discount_unit` field
**Rationale**: Provider-specified measurement unit
**Alternatives Considered**: None

| Attribute     | Value                                                                |
| ------------- | -------------------------------------------------------------------- |
| Data Type     | String                                                               |
| Feature Level | Conditional                                                          |
| Description   | Provider-specified measurement unit for commitment discount quantity |

### 6. CapacityReservationId (CONDITIONAL)

**Decision**: Add as `string capacity_reservation_id` field
**Rationale**: Provider-assigned identifier
**Alternatives Considered**: None

| Attribute     | Value                                                         |
| ------------- | ------------------------------------------------------------- |
| Data Type     | String                                                        |
| Feature Level | Conditional                                                   |
| Description   | Identifier assigned to a capacity reservation by the provider |

### 7. CapacityReservationStatus (CONDITIONAL)

**Decision**: Add as enum `FocusCapacityReservationStatus` with values: Used, Unused
**Rationale**: Limited set of allowed values suits enum pattern
**Alternatives Considered**: String field - rejected for type safety

| Attribute      | Value                                                                    |
| -------------- | ------------------------------------------------------------------------ |
| Data Type      | String (enum) → `FocusCapacityReservationStatus`                         |
| Feature Level  | Conditional                                                              |
| Allowed Values | Used, Unused                                                             |
| Constraint     | Required when CapacityReservationId not null and ChargeCategory is Usage |

### 8. BillingAccountType (CONDITIONAL)

**Decision**: Add as `string billing_account_type` field
**Rationale**: Provider-assigned name with no fixed values
**Alternatives Considered**: None

| Attribute     | Value                                                          |
| ------------- | -------------------------------------------------------------- |
| Data Type     | String                                                         |
| Feature Level | Conditional                                                    |
| Description   | Provider-assigned name to identify the type of billing account |

### 9. SubAccountType (CONDITIONAL)

**Decision**: Add as `string sub_account_type` field
**Rationale**: Provider-assigned identifier with no fixed values
**Alternatives Considered**: None

| Attribute     | Value                                                       |
| ------------- | ----------------------------------------------------------- |
| Data Type     | String                                                      |
| Feature Level | Conditional                                                 |
| Description   | Provider-assigned identifier for sub-account classification |

### 10. ContractedUnitPrice (CONDITIONAL)

**Decision**: Add as `double contracted_unit_price` field
**Rationale**: Decimal type for price values
**Alternatives Considered**: None

| Attribute     | Value                                                      |
| ------------- | ---------------------------------------------------------- |
| Data Type     | Decimal → `double`                                         |
| Feature Level | Conditional                                                |
| Description   | Agreed-upon unit price per pricing unit for associated SKU |

### 11. PricingCurrency (CONDITIONAL)

**Decision**: Add as `string pricing_currency` field
**Rationale**: ISO 4217 currency code
**Alternatives Considered**: None

| Attribute     | Value                                                                     |
| ------------- | ------------------------------------------------------------------------- |
| Data Type     | String                                                                    |
| Feature Level | Conditional                                                               |
| Description   | Currency for pricing-related columns when different from billing currency |
| Format        | ISO 4217 currency code                                                    |

### 12. PricingCurrencyContractedUnitPrice (CONDITIONAL)

**Decision**: Add as `double pricing_currency_contracted_unit_price` field
**Rationale**: Decimal type for price values in pricing currency
**Alternatives Considered**: None

| Attribute     | Value                                                 |
| ------------- | ----------------------------------------------------- |
| Data Type     | Decimal → `double`                                    |
| Feature Level | Conditional                                           |
| Description   | Contracted unit price denominated in pricing currency |

### 13. PricingCurrencyEffectiveCost (CONDITIONAL)

**Decision**: Add as `double pricing_currency_effective_cost` field
**Rationale**: Decimal type for cost values in pricing currency
**Alternatives Considered**: None

| Attribute     | Value                                          |
| ------------- | ---------------------------------------------- |
| Data Type     | Decimal → `double`                             |
| Feature Level | Conditional                                    |
| Description   | Effective cost denominated in pricing currency |

### 14. PricingCurrencyListUnitPrice (CONDITIONAL)

**Decision**: Add as `double pricing_currency_list_unit_price` field
**Rationale**: Decimal type for price values in pricing currency
**Alternatives Considered**: None

| Attribute     | Value                                           |
| ------------- | ----------------------------------------------- |
| Data Type     | Decimal → `double`                              |
| Feature Level | Conditional                                     |
| Description   | List unit price denominated in pricing currency |

### 15. Publisher (CONDITIONAL)

**Decision**: Add as `string publisher` field
**Rationale**: Provider-assigned identifier for the entity that published the service
**Alternatives Considered**: None

| Attribute     | Value                                        |
| ------------- | -------------------------------------------- |
| Data Type     | String                                       |
| Feature Level | Conditional                                  |
| Description   | Entity that published the service or product |

### 16. ServiceSubcategory (CONDITIONAL)

**Decision**: Add as `string service_subcategory` field
**Rationale**: Granular service classification
**Alternatives Considered**: Enum - rejected due to variable provider values

| Attribute     | Value                                                                |
| ------------- | -------------------------------------------------------------------- |
| Data Type     | String                                                               |
| Feature Level | Conditional                                                          |
| Description   | Granular service classification supporting functional categorization |

### 17. SkuMeter (CONDITIONAL)

**Decision**: Add as `string sku_meter` field
**Rationale**: Provider-assigned meter identifier
**Alternatives Considered**: None

| Attribute     | Value                                          |
| ------------- | ---------------------------------------------- |
| Data Type     | String                                         |
| Feature Level | Conditional                                    |
| Description   | Provider-assigned meter identifier for the SKU |

### 18. SkuPriceDetails (CONDITIONAL)

**Decision**: Add as `string sku_price_details` field
**Rationale**: Provider-specific pricing information
**Alternatives Considered**: Structured type - rejected for flexibility

| Attribute     | Value                                                    |
| ------------- | -------------------------------------------------------- |
| Data Type     | String                                                   |
| Feature Level | Conditional                                              |
| Description   | Additional provider-specific pricing details for the SKU |

## New Enum Types Required

### FocusCommitmentDiscountStatus

```protobuf
enum FocusCommitmentDiscountStatus {
  FOCUS_COMMITMENT_DISCOUNT_STATUS_UNSPECIFIED = 0;
  FOCUS_COMMITMENT_DISCOUNT_STATUS_USED = 1;
  FOCUS_COMMITMENT_DISCOUNT_STATUS_UNUSED = 2;
}
```

### FocusCapacityReservationStatus

```protobuf
enum FocusCapacityReservationStatus {
  FOCUS_CAPACITY_RESERVATION_STATUS_UNSPECIFIED = 0;
  FOCUS_CAPACITY_RESERVATION_STATUS_USED = 1;
  FOCUS_CAPACITY_RESERVATION_STATUS_UNUSED = 2;
}
```

## Proto Field Number Assignment

Starting from field number 41 (last used is 40 for `invoice_issuer`):

| Field Number | Column Name                            |
| ------------ | -------------------------------------- |
| 41           | contracted_cost                        |
| 42           | billing_account_type                   |
| 43           | sub_account_type                       |
| 44           | capacity_reservation_id                |
| 45           | capacity_reservation_status            |
| 46           | commitment_discount_quantity           |
| 47           | commitment_discount_status             |
| 48           | commitment_discount_type               |
| 49           | commitment_discount_unit               |
| 50           | contracted_unit_price                  |
| 51           | pricing_currency                       |
| 52           | pricing_currency_contracted_unit_price |
| 53           | pricing_currency_effective_cost        |
| 54           | pricing_currency_list_unit_price       |
| 55           | publisher                              |
| 56           | service_subcategory                    |
| 57           | sku_meter                              |
| 58           | sku_price_details                      |

## Best Practices Applied

### Protobuf Design

- **Field numbering**: Sequential from 41-58, avoiding reserved ranges
- **Naming convention**: snake_case matching existing pattern
- **Type mapping**: FOCUS Decimal → proto double, FOCUS DateTime → Timestamp
- **Enums**: Used for columns with fixed allowed values (CommitmentDiscountStatus,
  CapacityReservationStatus)

### Documentation

- **Proto comments**: Include FOCUS section reference (e.g., "// Section 3.XX")
- **Godoc comments**: Describe purpose, parameters, return values
- **Examples**: Show real-world usage for each new column

### Testing

- **Conformance tests**: Validate new column presence and types
- **Builder tests**: Test each new With\* method
- **Audit script**: Verify 57/57 columns present

## Sources

- [FOCUS Specification v1.2](https://focus.finops.org/focus-specification/v1-2/)
- [FOCUS Column Library](https://focus.finops.org/focus-columns/)
