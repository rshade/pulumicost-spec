# Plan: Pulumicost Spec - FinOps FOCUS 1.2 Column Audit

**Status:** Proposed
**Version:** 1.0
**Focus:** Data Integrity & Full Schema Compliance

## 1. Context

After the initial implementation of the FOCUS 1.2 integration (`pulumicost-spec-focus.md`), we must verify that the resulting Protobuf definitions strictly adhere to the official [FinOps FOCUS 1.2 Specification](https://focus.finops.org). The initial pass prioritized architecture and key columns; this pass ensures **100% column coverage**.

## 2. Objectives

1.  **Column Completeness Audit:** Verify that every mandatory and conditional column defined in FOCUS 1.2 exists in `focus.proto`.
2.  **Type Verification:** Ensure Protobuf types (e.g., `google.protobuf.Timestamp` vs `string`) align with FOCUS data types.
3.  **Missing Column Implementation:** Add any columns missed during the initial "happy path" implementation.

## 3. Verification Strategy

We will cross-reference `proto/pulumicost/v1/focus.proto` against the master list below.

### 3.1 The "Must-Have" List (FOCUS 1.2)

**Identity & Hierarchy**
*   [ ] `BillingAccountId` (String)
*   [ ] `BillingAccountName` (String)
*   [ ] `BillingCurrency` (String)
*   [ ] `ProviderName` (String)
*   [ ] `RegionId` (String)
*   [ ] `RegionName` (String)
*   [ ] `SubAccountId` (String)
*   [ ] `SubAccountName` (String)
*   [ ] `AvailabilityZone` (String) - *Conditional*

**Cost & Usage**
*   [ ] `BilledCost` (Decimal)
*   [ ] `BillingFrequency` (String)
*   [ ] `BillingPeriodStart` (Timestamp)
*   [ ] `BillingPeriodEnd` (Timestamp)
*   [ ] `ChargeCategory` (Enum)
*   [ ] `ChargeClass` (Enum) - *New in 1.0+*
*   [ ] `ChargeDescription` (String)
*   [ ] `ChargeFrequency` (Enum)
*   [ ] `ChargePeriodStart` (Timestamp)
*   [ ] `ChargePeriodEnd` (Timestamp)
*   [ ] `ConsumedQuantity` (Decimal)
*   [ ] `ConsumedUnit` (String)
*   [ ] `EffectiveCost` (Decimal)
*   [ ] `ListCost` (Decimal)
*   [ ] `ListUnitPrice` (Decimal)
*   [ ] `PricingCategory` (Enum)
*   [ ] `PricingQuantity` (Decimal)
*   [ ] `PricingUnit` (String)

**Product & Service**
*   [ ] `ResourceId` (String)
*   [ ] `ResourceName` (String)
*   [ ] `ResourceType` (String)
*   [ ] `ServiceCategory` (Enum)
*   [ ] `ServiceName` (String)
*   [ ] `SkuId` (String)
*   [ ] `SkuPriceId` (String)

**Discounts & Metadata**
*   [ ] `CommitmentDiscountCategory` (Enum)
*   [ ] `CommitmentDiscountId` (String)
*   [ ] `CommitmentDiscountName` (String)
*   [ ] `InvoiceIssuer` (String)
*   [ ] `InvoiceId` (String) - *Critical for Reconciliation*
*   [ ] `Tags` (Map)

## 4. Implementation Plan

### Phase 1: Audit Script
*   [ ] Create `scripts/audit_focus_columns.sh` (or a Go test) that parses `focus.proto` and reports missing fields.

### Phase 2: Schema Expansion
*   [ ] Add any missing fields to `proto/pulumicost/v1/focus.proto`.
*   [ ] Ensure comments in the Proto file reference the specific FOCUS section (e.g., "// Section 2.4: Charge Details").

### Phase 3: Builder Update
*   [ ] Update `sdk/go/pluginsdk/focus_builder.go` to include methods for any newly added columns (e.g., `WithCommitmentDetails(...)`).
