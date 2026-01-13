# Data Model: FOCUS 1.2 Column Audit

**Date**: 2025-11-28
**Feature**: 010-focus-column-audit

## Entity Overview

This feature extends the existing `FocusCostRecord` proto message with 19 additional
columns and adds 2 new enum types to achieve full FOCUS 1.2 compliance.

## Primary Entity: FocusCostRecord

### Current State (38 fields)

The existing `FocusCostRecord` message in `proto/finfocus/v1/focus.proto` contains
38 fields covering 13 of 14 mandatory columns and 24 of 42 conditional columns.

### Target State (57 fields)

After implementation, `FocusCostRecord` will contain all 57 FOCUS 1.2 columns:

- 14 mandatory columns (100%)
- 1 recommended column (100%)
- 42 conditional columns (100%)

## New Fields

### Financial Fields (Mandatory + Conditional)

| Field                 | Proto Type | FOCUS Type | Level       | Description                                  |
| --------------------- | ---------- | ---------- | ----------- | -------------------------------------------- |
| contracted_cost       | double     | Decimal    | Mandatory   | Cost = ContractedUnitPrice × PricingQuantity |
| contracted_unit_price | double     | Decimal    | Conditional | Agreed-upon unit price per pricing unit      |

### Account Fields (Conditional)

| Field                | Proto Type | FOCUS Type | Level       | Description             |
| -------------------- | ---------- | ---------- | ----------- | ----------------------- |
| billing_account_type | string     | String     | Conditional | Type of billing account |
| sub_account_type     | string     | String     | Conditional | Type of sub-account     |

### Capacity Reservation Fields (Conditional)

| Field                       | Proto Type                     | FOCUS Type | Level       | Description                     |
| --------------------------- | ------------------------------ | ---------- | ----------- | ------------------------------- |
| capacity_reservation_id     | string                         | String     | Conditional | Capacity reservation identifier |
| capacity_reservation_status | FocusCapacityReservationStatus | String     | Conditional | Used/Unused                     |

### Commitment Discount Fields (Conditional)

| Field                        | Proto Type                    | FOCUS Type | Level       | Description                 |
| ---------------------------- | ----------------------------- | ---------- | ----------- | --------------------------- |
| commitment_discount_quantity | double                        | Decimal    | Conditional | Amount of discount          |
| commitment_discount_status   | FocusCommitmentDiscountStatus | String     | Conditional | Used/Unused                 |
| commitment_discount_type     | string                        | String     | Conditional | Type of commitment discount |
| commitment_discount_unit     | string                        | String     | Conditional | Unit of measurement         |

### Pricing Currency Fields (Conditional)

| Field                                  | Proto Type | FOCUS Type | Level       | Description                    |
| -------------------------------------- | ---------- | ---------- | ----------- | ------------------------------ |
| pricing_currency                       | string     | String     | Conditional | ISO 4217 currency code         |
| pricing_currency_contracted_unit_price | double     | Decimal    | Conditional | Price in pricing currency      |
| pricing_currency_effective_cost        | double     | Decimal    | Conditional | Cost in pricing currency       |
| pricing_currency_list_unit_price       | double     | Decimal    | Conditional | List price in pricing currency |

### Origination Fields (Conditional)

| Field     | Proto Type | FOCUS Type | Level       | Description                       |
| --------- | ---------- | ---------- | ----------- | --------------------------------- |
| publisher | string     | String     | Conditional | Entity that published the service |

### Service Fields (Conditional)

| Field               | Proto Type | FOCUS Type | Level       | Description                     |
| ------------------- | ---------- | ---------- | ----------- | ------------------------------- |
| service_subcategory | string     | String     | Conditional | Granular service classification |

### SKU Fields (Conditional)

| Field             | Proto Type | FOCUS Type | Level       | Description                |
| ----------------- | ---------- | ---------- | ----------- | -------------------------- |
| sku_meter         | string     | String     | Conditional | Meter identifier           |
| sku_price_details | string     | String     | Conditional | Additional pricing details |

## New Enum Types

### FocusCommitmentDiscountStatus

```text
FOCUS_COMMITMENT_DISCOUNT_STATUS_UNSPECIFIED = 0
FOCUS_COMMITMENT_DISCOUNT_STATUS_USED = 1
FOCUS_COMMITMENT_DISCOUNT_STATUS_UNUSED = 2
```

**Validation Rules**:

- Required when CommitmentDiscountId is not null AND ChargeCategory is Usage
- Must be USED or UNUSED (not UNSPECIFIED) when applicable

### FocusCapacityReservationStatus

```text
FOCUS_CAPACITY_RESERVATION_STATUS_UNSPECIFIED = 0
FOCUS_CAPACITY_RESERVATION_STATUS_USED = 1
FOCUS_CAPACITY_RESERVATION_STATUS_UNUSED = 2
```

**Validation Rules**:

- Required when CapacityReservationId is not null AND ChargeCategory is Usage
- Must be USED or UNUSED (not UNSPECIFIED) when applicable

## Field Relationships

### ContractedCost Calculation

```text
IF ContractedUnitPrice != null AND PricingQuantity != null AND ChargeClass != "Correction":
    ContractedCost MUST EQUAL ContractedUnitPrice × PricingQuantity
```

### Commitment Discount Dependencies

```text
IF CommitmentDiscountId != null:
    CommitmentDiscountType MUST be present
    CommitmentDiscountUnit SHOULD be present
    IF ChargeCategory == "Usage":
        CommitmentDiscountStatus MUST be present
        CommitmentDiscountQuantity SHOULD be present
```

### Capacity Reservation Dependencies

```text
IF CapacityReservationId != null AND ChargeCategory == "Usage":
    CapacityReservationStatus MUST be present
```

### Pricing Currency Dependencies

```text
IF PricingCurrency != null AND PricingCurrency != BillingCurrency:
    PricingCurrencyContractedUnitPrice MAY be present
    PricingCurrencyEffectiveCost MAY be present
    PricingCurrencyListUnitPrice MAY be present
```

## State Transitions

Not applicable - FocusCostRecord is an immutable value object representing a point-in-time
cost record. There are no state transitions.

## Validation Rules Summary

### Field-Level Validation

| Field                       | Rule                                                       |
| --------------------------- | ---------------------------------------------------------- |
| contracted_cost             | Must be >= 0 when present                                  |
| pricing_currency            | Must be valid ISO 4217 code when present                   |
| capacity_reservation_status | Must not be UNSPECIFIED when CapacityReservationId present |
| commitment_discount_status  | Must not be UNSPECIFIED when CommitmentDiscountId present  |

### Cross-Field Validation

| Rule                              | Fields Involved                                          |
| --------------------------------- | -------------------------------------------------------- |
| ContractedCost calculation        | contracted_cost, contracted_unit_price, pricing_quantity |
| Commitment discount completeness  | commitment*discount*\* fields                            |
| Capacity reservation completeness | capacity*reservation*\* fields                           |
| Pricing currency consistency      | pricing*currency, pricing_currency*\* fields             |

## Builder Method Mapping

Each new field requires a corresponding builder method:

| Field                                  | Builder Method                                                 |
| -------------------------------------- | -------------------------------------------------------------- |
| contracted_cost                        | WithContractedCost(cost float64)                               |
| billing_account_type                   | WithBillingAccountType(accountType string)                     |
| sub_account_type                       | WithSubAccountType(accountType string)                         |
| capacity_reservation_id                | WithCapacityReservation(id string, status)                     |
| capacity_reservation_status            | (included in WithCapacityReservation)                          |
| commitment_discount_quantity           | WithCommitmentDiscountDetails(qty, status, discountType, unit) |
| commitment_discount_status             | (included in WithCommitmentDiscountDetails)                    |
| commitment_discount_type               | (included in WithCommitmentDiscountDetails)                    |
| commitment_discount_unit               | (included in WithCommitmentDiscountDetails)                    |
| contracted_unit_price                  | WithContractedUnitPrice(price float64)                         |
| pricing_currency                       | WithPricingCurrency(currency string)                           |
| pricing_currency_contracted_unit_price | WithPricingCurrencyPrices(contracted, effective, list)         |
| pricing_currency_effective_cost        | (included in WithPricingCurrencyPrices)                        |
| pricing_currency_list_unit_price       | (included in WithPricingCurrencyPrices)                        |
| publisher                              | WithPublisher(publisher string)                                |
| service_subcategory                    | WithServiceSubcategory(subcategory string)                     |
| sku_meter                              | WithSkuDetails(meter, priceDetails string)                     |
| sku_price_details                      | (included in WithSkuDetails)                                   |
