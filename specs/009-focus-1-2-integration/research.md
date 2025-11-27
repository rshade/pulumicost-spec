# Research: FinOps FOCUS 1.2 Integration

**Feature**: `focus-1-2-integration`
**Status**: Research Complete

## 1. Mandatory Columns

Based on FinOps FOCUS 1.2 specification (and common FinOps usage), the following columns are generally treated as mandatory or highly critical. We will treat them as "Mandatory" in our `FocusBuilder` validation logic.

**Identity & Hierarchy**
*   `BillingAccountId` (Explicit "MUST" in spec)
*   `ProviderName` (Implicitly required for context)
*   `ChargePeriodStart`
*   `ChargePeriodEnd`

**Service & Product**
*   `ServiceCategory` (Mandatory, Enum)
*   `ServiceName`

**Charge Details**
*   `ChargeCategory` (Mandatory, Enum)
*   `PricingCategory` (Mandatory, Enum)

**Financials**
*   `BilledCost` (Primary cost metric)
*   `Currency` (Required if cost is present)

*Decision*: The `FocusBuilder.Build()` method will return an error if any of these fields are missing/empty.

## 2. Controlled Vocabularies (Enums)

We will implement strict Enums in `enums.proto` based on the search results.

### FocusServiceCategory
*(Note: Full list not found in search, will use standard list from FOCUS 1.0/1.1 + generic placeholders which are forward compatible. We will allow "Unspecified" for unknown values)*

*   `Compute`
*   `Storage`
*   `Network`
*   `Database`
*   `Analytics`
*   `MachineLearning` (AI/ML)
*   `Management` (Governance)
*   `Security` (Identity)
*   `DeveloperTools`
*   `Other`

### FocusChargeCategory
*(Source: Search Result [1])*

*   `Usage`
*   `Purchase`
*   `Credit`
*   `Tax`
*   `Refund`
*   `Adjustment` (Often seen in 1.0, mapping to Refund/Credit or separate) -> *Refinement*: Spec lists 5 specific ones. We will stick to the 5 search results + `Unspecified`.

### FocusPricingCategory
*(Source: Search Result [1])*

*   `Standard` (On-Demand)
*   `Committed` (Savings Plans/RI)
*   `Dynamic` (Spot)
*   `Other`

## 3. Data Type Decisions

*   **Financials**: Protobuf `double` (as per Clarification Q1).
*   **Timestamps**: `google.protobuf.Timestamp`.
*   **Tags/Extensions**: `map<string, string>`.

## 4. Versioning & "Backpack" Strategy

*   **Strategy**: Use a `map<string, string> extended_columns` field in `FocusCostRecord`.
*   **Usage**:
    *   Provider-specific columns (e.g., `aws_resource_tags`).
    *   Future FOCUS columns (e.g., from FOCUS 1.3).
    *   "Unknown" Enum values (if a provider sends a category not in our Enum, set Enum to `Other` and put raw value in `extended_columns["raw_service_category"]`).
