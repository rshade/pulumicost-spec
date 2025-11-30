# Feature Specification: FOCUS 1.2 Column Audit

**Feature Branch**: `010-focus-column-audit`
**Created**: 2025-11-28
**Status**: Draft
**Input**: User description: "Pulumicost Spec - FinOps FOCUS 1.2 Column Audit"

## Clarifications

### Session 2025-11-28

- Q: What is the correct column for billing frequency in FOCUS 1.2? → A: There is no
  `BillingFrequency` column in FOCUS 1.2. The correct column is `ChargeFrequency` with
  values: One-Time, Recurring, Usage-Based. This column is already implemented.
- Q: How many columns are in the official FOCUS 1.2 specification? → A: 57 columns total
  (not 35 as originally stated). Source: [FOCUS Column Library](https://focus.finops.org/focus-columns/)

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Schema Completeness Verification (Priority: P1)

A plugin developer wants to ensure their FOCUS 1.2 implementation covers all mandatory
columns defined by the FinOps Foundation, so they can certify their plugin's compliance.

**Why this priority**: Without complete column coverage, plugins cannot claim FOCUS 1.2
compliance, which is a core value proposition of PulumiCost.

**Independent Test**: Can be verified by comparing the proto definition against the official
FOCUS 1.2 specification and running an automated audit script.

**Acceptance Scenarios**:

1. **Given** the `focus.proto` file exists, **When** an audit script parses all field
   definitions, **Then** every mandatory FOCUS 1.2 column is present with the correct
   protobuf type.
2. **Given** a missing column is detected, **When** the audit report is generated,
   **Then** the report lists the specific column name, expected type, and FOCUS
   specification section reference.

---

### User Story 2 - Type Correctness Validation (Priority: P2)

A data engineer needs confidence that Protobuf types align with FOCUS data types so that
data transformations preserve accuracy when converting between formats.

**Why this priority**: Incorrect types can cause data loss or silent corruption during
cost aggregation and reporting.

**Independent Test**: Can be validated by mapping each proto field type to its FOCUS
equivalent and verifying semantic compatibility.

**Acceptance Scenarios**:

1. **Given** a field defined as `google.protobuf.Timestamp` in proto, **When** compared
   to FOCUS specification, **Then** the corresponding FOCUS column is documented as
   DateTime type.
2. **Given** a field defined as `double` in proto, **When** compared to FOCUS
   specification, **Then** the corresponding FOCUS column is documented as Decimal type.
3. **Given** an enum type in proto, **When** compared to FOCUS specification, **Then**
   all valid enum values from FOCUS are represented in the proto enum.

---

### User Story 3 - Builder API Completeness (Priority: P3)

A plugin developer using the Go SDK needs builder methods for all FOCUS columns so they
can construct fully compliant cost records without directly manipulating the proto struct.

**Why this priority**: Builder pattern provides a safer, more ergonomic API that reduces
errors and improves developer experience.

**Independent Test**: Can be verified by ensuring every field in `FocusCostRecord` has a
corresponding builder method in `focus_builder.go`.

**Acceptance Scenarios**:

1. **Given** a new column is added to `focus.proto`, **When** the builder is updated,
   **Then** a corresponding `With*` method exists that sets the field value.
2. **Given** a developer uses only builder methods, **When** they call `Build()`,
   **Then** all FOCUS mandatory fields can be populated without accessing the underlying
   struct directly.

---

### User Story 4 - Developer Documentation (Priority: P2)

A new plugin developer needs comprehensive documentation to understand how to use the
FOCUS 1.2 types, builder patterns, and validation functions without reading source code.

**Why this priority**: Good documentation reduces onboarding time and support burden,
enabling faster plugin development and fewer integration errors.

**Independent Test**: Can be verified by having a new developer successfully implement
a FOCUS-compliant plugin using only the documentation (no source code reading required).

**Acceptance Scenarios**:

1. **Given** a developer reads the SDK documentation, **When** they look up any FOCUS
   column, **Then** they find a clear description, valid values, and usage examples.
2. **Given** a developer wants to use the builder API, **When** they consult the
   documentation, **Then** they find complete method signatures with parameter
   descriptions and example code.
3. **Given** a developer encounters a validation error, **When** they check the
   documentation, **Then** they find troubleshooting guidance and common error resolutions.

---

### User Story 5 - User-Facing Documentation (Priority: P3)

A FinOps practitioner or cost analyst needs reference documentation to understand what
each FOCUS column means and how it maps to their cloud provider's billing data.

**Why this priority**: End users need to understand the data model to effectively
analyze and report on cloud costs.

**Independent Test**: Can be verified by a non-developer successfully understanding
the meaning and purpose of each FOCUS column from the documentation alone.

**Acceptance Scenarios**:

1. **Given** a user reads the FOCUS column reference, **When** they look up any column,
   **Then** they find a plain-language description suitable for non-developers.
2. **Given** a user wants to understand provider mappings, **When** they consult the
   documentation, **Then** they find examples showing how AWS/Azure/GCP fields map to
   FOCUS columns.

---

### Edge Cases

- What happens when a FOCUS column is marked "conditional" rather than "mandatory"?
  - Conditional columns should still be present in the proto but documented as optional.
- How does the system handle FOCUS columns that may be added in future versions?
  - The `extended_columns` map provides forward-compatibility for new columns.
- How should documentation handle provider-specific variations?
  - Document the canonical FOCUS representation, then provide provider-specific mapping
    tables showing how each provider's field names translate to FOCUS columns.
- What if a provider doesn't support a particular FOCUS column?
  - Document which columns are universally supported vs. provider-specific, and explain
    how to handle missing data (null values, default values, or omission).

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST define all 57 FOCUS 1.2 columns in `focus.proto`, including
  14 mandatory, 1 recommended, and 42 conditional columns.
- **FR-002**: System MUST use `google.protobuf.Timestamp` for all DateTime fields.
- **FR-003**: System MUST use `double` for all Decimal/Numeric fields.
- **FR-004**: System MUST use enum types for all FOCUS enumerated columns
  (ChargeCategory, ChargeClass, ChargeFrequency, PricingCategory, ServiceCategory,
  CommitmentDiscountCategory).
- **FR-005**: System MUST add the missing mandatory column `ContractedCost` (Decimal).
- **FR-006**: System MUST add the 19 missing columns (1 mandatory, 18 conditional)
  identified in the audit.
- **FR-007**: Each proto field MUST include a comment referencing the specific FOCUS
  section (e.g., "// Section 3.14: Charge Frequency").
- **FR-008**: The `FocusRecordBuilder` MUST provide `With*` methods for all proto fields.
- **FR-009**: An audit script MUST be created to automatically verify column coverage
  against the official FOCUS 1.2 specification.

### Documentation Requirements

- **FR-010**: All exported Go functions, types, and methods MUST have godoc comments
  achieving 80%+ documentation coverage.
- **FR-011**: Each builder method MUST include a godoc comment describing its purpose,
  parameters, and the FOCUS section it implements.
- **FR-012**: A developer guide (`sdk/go/pluginsdk/README.md`) MUST be created or updated
  with:
  - Quick start example for building FOCUS records
  - Complete builder method reference with examples
  - Validation error troubleshooting guide
  - Migration guide for existing plugins
- **FR-013**: A user-facing FOCUS column reference (`docs/focus-columns.md`) MUST be
  created with:
  - Plain-language description of each FOCUS column
  - Data type and valid values for each column
  - Provider mapping examples (AWS, Azure, GCP to FOCUS)
  - Common use cases and query patterns
- **FR-014**: Proto file comments MUST serve as inline documentation, including:
  - Column description matching FOCUS specification
  - Data type constraints and valid value ranges
  - Cross-references to related columns
- **FR-015**: Code examples MUST be provided in `examples/plugins/` demonstrating:
  - Building a complete FOCUS record with all columns
  - Handling optional/conditional columns
  - Validating records before submission

### Key Entities

- **FocusCostRecord**: The primary proto message containing all FOCUS 1.2 columns
  representing a single cost line item.
- **FOCUS Enum Types**: `FocusChargeCategory`, `FocusChargeClass`, `FocusChargeFrequency`,
  `FocusPricingCategory`, `FocusServiceCategory`, `FocusCommitmentDiscountCategory`,
  plus new enums for `CommitmentDiscountStatus`, `CommitmentDiscountType`,
  `CapacityReservationStatus`.
- **FocusRecordBuilder**: Go SDK helper that provides fluent API for constructing
  valid FOCUS records.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: 100% of FOCUS 1.2 mandatory columns are present in `focus.proto`
  (14/14 mandatory columns).
- **SC-002**: All 57 FOCUS 1.2 columns are defined in `focus.proto`.
- **SC-003**: 100% of proto fields have corresponding builder methods in
  `focus_builder.go`.
- **SC-004**: Audit script completes with zero missing mandatory columns reported.
- **SC-005**: All enum types contain all valid values specified in FOCUS 1.2
  documentation.
- **SC-006**: Every proto field includes a documentation comment with FOCUS
  section reference.

### Documentation Outcomes

- **SC-007**: Go documentation coverage is 80% or higher as measured by
  `go doc` coverage tools.
- **SC-008**: Developer guide includes working code examples that compile and
  pass validation.
- **SC-009**: User-facing documentation covers all 57 FOCUS columns with
  plain-language descriptions.
- **SC-010**: Provider mapping documentation includes at least one example per
  major cloud provider (AWS, Azure, GCP).
- **SC-011**: New developers can build a valid FOCUS record using only the
  documentation within 30 minutes (validated by user testing).

## Appendix: FOCUS 1.2 Column Reference

### Audit Summary

**Source**: [FOCUS Specification v1.2](https://focus.finops.org/focus-specification/v1-2/)
and [FOCUS Column Library](https://focus.finops.org/focus-columns/)

| Metric                  | Count |
| ----------------------- | ----- |
| Total FOCUS 1.2 Columns | 57    |
| Currently Implemented   | 38    |
| Missing                 | 19    |
| Missing Mandatory       | 1     |
| Missing Conditional     | 18    |

### Implemented Columns (38)

#### Mandatory Columns (13 of 14 implemented)

| Column             | Type     | Proto Field              |
| ------------------ | -------- | ------------------------ |
| BilledCost         | Decimal  | `billed_cost`            |
| BillingAccountId   | String   | `billing_account_id`     |
| BillingAccountName | String   | `billing_account_name`   |
| BillingCurrency    | String   | `billing_currency`       |
| BillingPeriodEnd   | DateTime | `billing_period_end`     |
| BillingPeriodStart | DateTime | `billing_period_start`   |
| ChargeCategory     | String   | `charge_category` (enum) |
| ChargeClass        | String   | `charge_class` (enum)    |
| ChargeDescription  | String   | `charge_description`     |
| ChargePeriodEnd    | DateTime | `charge_period_end`      |
| ChargePeriodStart  | DateTime | `charge_period_start`    |
| Provider           | String   | `provider_name`          |

#### Recommended Columns (1 of 1 implemented)

| Column          | Type   | Proto Field               |
| --------------- | ------ | ------------------------- |
| ChargeFrequency | String | `charge_frequency` (enum) |

#### Conditional Columns (24 of 42 implemented)

| Column                     | Type    | Proto Field                           |
| -------------------------- | ------- | ------------------------------------- |
| AvailabilityZone           | String  | `availability_zone`                   |
| CommitmentDiscountCategory | String  | `commitment_discount_category` (enum) |
| CommitmentDiscountId       | String  | `commitment_discount_id`              |
| CommitmentDiscountName     | String  | `commitment_discount_name`            |
| ConsumedQuantity           | Decimal | `consumed_quantity`                   |
| ConsumedUnit               | String  | `consumed_unit`                       |
| EffectiveCost              | Decimal | `effective_cost`                      |
| InvoiceId                  | String  | `invoice_id`                          |
| InvoiceIssuer              | String  | `invoice_issuer`                      |
| ListCost                   | Decimal | `list_cost`                           |
| ListUnitPrice              | Decimal | `list_unit_price`                     |
| PricingCategory            | String  | `pricing_category` (enum)             |
| PricingQuantity            | Decimal | `pricing_quantity`                    |
| PricingUnit                | String  | `pricing_unit`                        |
| RegionId                   | String  | `region_id`                           |
| RegionName                 | String  | `region_name`                         |
| ResourceId                 | String  | `resource_id`                         |
| ResourceName               | String  | `resource_name`                       |
| ResourceType               | String  | `resource_type`                       |
| ServiceCategory            | String  | `service_category` (enum)             |
| ServiceName                | String  | `service_name`                        |
| SkuId                      | String  | `sku_id`                              |
| SkuPriceId                 | String  | `sku_price_id`                        |
| SubAccountId               | String  | `sub_account_id`                      |
| SubAccountName             | String  | `sub_account_name`                    |
| Tags                       | String  | `tags` (map)                          |

### Missing Columns (19)

#### Missing Mandatory (1) - HIGH PRIORITY

| Column             | Type    | Category  |
| ------------------ | ------- | --------- |
| **ContractedCost** | Decimal | Financial |

#### Missing Conditional (18)

| Column                             | Type    | Category    |
| ---------------------------------- | ------- | ----------- |
| BillingAccountType                 | String  | Account     |
| SubAccountType                     | String  | Account     |
| CapacityReservationId              | String  | Capacity    |
| CapacityReservationStatus          | String  | Capacity    |
| CommitmentDiscountQuantity         | Decimal | Commitment  |
| CommitmentDiscountStatus           | String  | Commitment  |
| CommitmentDiscountType             | String  | Commitment  |
| CommitmentDiscountUnit             | String  | Commitment  |
| ContractedUnitPrice                | Decimal | Pricing     |
| PricingCurrency                    | String  | Pricing     |
| PricingCurrencyContractedUnitPrice | Decimal | Pricing     |
| PricingCurrencyEffectiveCost       | Decimal | Pricing     |
| PricingCurrencyListUnitPrice       | Decimal | Pricing     |
| Publisher                          | String  | Origination |
| ServiceSubcategory                 | String  | Service     |
| SkuMeter                           | String  | SKU         |
| SkuPriceDetails                    | String  | SKU         |

### Assumptions

- The [FOCUS Specification v1.2](https://focus.finops.org/focus-specification/v1-2/) is
  the authoritative source for column definitions.
- The `extended_columns` map provides adequate forward-compatibility for future FOCUS
  versions.
- Protobuf `double` type is acceptable for FOCUS Decimal fields (no precision
  requirements specified in FOCUS).
- Enum values in the current implementation align with FOCUS 1.2 specification values.
- Conditional columns are important for full FOCUS compliance but can be prioritized
  after mandatory columns are complete.
