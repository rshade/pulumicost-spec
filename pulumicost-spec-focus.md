# Plan: Pulumicost Spec - FinOps FOCUS 1.2 Integration

**Status:** Proposed
**Version:** 1.5
**Focus:** Schema Definition, Enums, SDK Abstraction, and Future-Proofing

## 1. Context & Vision

The `pulumicost-spec` currently serves as a high-level query API for cloud costs. To achieve
the project's vision of becoming the **universal, open-source standard for cloud cost
observability**, it must align with the industry-standard **FinOps FOCUS 1.2 Specification**.

Adopting FOCUS 1.2 will transform `pulumicost` from a simple cost summarizer into a
forensic-grade cost analysis tool. To ensure this is sustainable and upgradeable
(e.g., to FOCUS 1.3), we will employ a **"Backpack & Builder"** strategy to insulate
plugin developers from schema complexity.

## 2. Objectives

1. **Define `FocusCostRecord`:** Create a comprehensive Protobuf message implementing
   the ~40 mandatory and optional columns of FOCUS 1.2.
2. **Future-Proofing (The "Backpack"):** Include an `extended_columns` map to support
   version 1.3+ features before strict schema support exists.
3. **SDK Abstraction (The "Shield"):** Mandate an **Opaque Builder Pattern** so plugin
   developers never touch the raw Protobuf struct, allowing us to refactor the schema
   without breaking plugin code.
4. **Conformance Validation:** Provide a test harness to verify compliance.

## 3. User Benefits (Enabled by this Spec)

- **Universal Compatibility:** Data exported from `pulumicost` is native to FinOps
  platforms (Cloudability, Vantage).
- **Forensic Precision:** Distinguish between "List Cost", "Billed Cost", and "Effective Cost".
- **Stability:** Plugins written today will compile against future Spec versions thanks
  to the Builder abstraction.

## 4. Proposed Specification Changes

### 4.1 New Protobuf Message: `FocusCostRecord`

_File:_ `proto/pulumicost/v1/focus.proto`

```protobuf
// FocusCostRecord implements FinOps FOCUS 1.2 columns
message FocusCostRecord {
  // --- Identity & Hierarchy ---
  string provider_name = 1;
  string billing_account_id = 2;
  string billing_account_name = 3;
  string sub_account_id = 4;
  string resource_id = 5;
  string resource_name = 6;

  // --- Service & Product ---
  FocusServiceCategory service_category = 7; // Use strict enum
  string service_name = 8;
  string sku_id = 9;
  string sku_price_id = 10;
  string region_id = 11;
  string region_name = 12;
  string availability_zone = 13;

  // --- Charge Details ---
  FocusChargeCategory charge_category = 14;  // Use strict enum
  string charge_description = 15;
  FocusPricingCategory pricing_category = 16; // Use strict enum
  string unit = 17;
  double usage_quantity = 18;

  // --- Financials ---
  string currency = 19;
  double billed_cost = 20;
  double list_cost = 21;
  double effective_cost = 22;
  string invoice_id = 23;

  // --- Time ---
  google.protobuf.Timestamp charge_period_start = 24;
  google.protobuf.Timestamp charge_period_end = 25;

  // --- Metadata & Extension ---
  map<string, string> tags = 26;

  // The "Backpack": Supports future FOCUS columns (1.3+) or provider-specific
  // extensions without requiring a schema schema update.
  map<string, string> extended_columns = 27;
}
```

### 4.2 Enum Standardization and Examples

_File:_ `proto/pulumicost/v1/enums.proto`

We will define strict Enums for categories with controlled vocabularies.

```protobuf
enum FocusServiceCategory {
  SERVICE_CATEGORY_UNSPECIFIED = 0;
  SERVICE_CATEGORY_COMPUTE = 1;       // e.g., AWS EC2, Azure VMs
  SERVICE_CATEGORY_STORAGE = 2;       // e.g., AWS S3, Azure Blob Storage
  SERVICE_CATEGORY_NETWORK = 3;       // e.g., Data Transfer, Load Balancers
  SERVICE_CATEGORY_DATABASE = 4;      // e.g., AWS RDS, Azure SQL
  SERVICE_CATEGORY_ANALYTICS = 5;
  SERVICE_CATEGORY_MACHINE_LEARNING = 6;
  SERVICE_CATEGORY_DEVELOPER_TOOLS = 7;
  SERVICE_CATEGORY_SECURITY_IDENTITY = 8;
  SERVICE_CATEGORY_MANAGEMENT_GOVERNANCE = 9;
  SERVICE_CATEGORY_INTERNET_OF_THINGS = 10;
  SERVICE_CATEGORY_OTHER = 11;
}

enum FocusChargeCategory {
  CHARGE_CATEGORY_UNSPECIFIED = 0;
  CHARGE_CATEGORY_USAGE = 1;          // Costs directly tied to consumption
  CHARGE_CATEGORY_TAX = 2;            // Taxes applied to services
  CHARGE_CATEGORY_ADJUSTMENT = 3;     // Credits, refunds, or one-time charges
  CHARGE_CATEGORY_PURCHASE = 4;       // Upfront purchase of RIs/Savings Plans
}

enum FocusPricingCategory {
  PRICING_CATEGORY_UNSPECIFIED = 0;
  PRICING_CATEGORY_ON_DEMAND = 1;     // Standard pay-as-you-go rates
  PRICING_CATEGORY_COMMITMENT = 2;    // Costs from Reserved Instances / Savings Plans
  PRICING_CATEGORY_DYNAMIC = 3;       // Spot / Preemptible instances
  PRICING_CATEGORY_FREE_TIER = 4;     // Costs covered by free tier allowances
  PRICING_CATEGORY_ESTIMATE = 5;      // Used for synthesized/fallback costs
}
```

### 4.3 SDK Builder Pattern (The "Shield")

_File:_ `sdk/go/pluginsdk/focus_builder.go`

The SDK serves as an **Anti-Corruption Layer**. The plugin developer depends on the SDK's
method signatures, not the Proto's struct layout.

```go
// Plugin Developer sees this:
type FocusBuilder interface {
    WithIdentity(provider string, accountID string) FocusBuilder
    WithResource(resourceID string, resourceType string) FocusBuilder
    WithServiceCategory(category pb.FocusServiceCategory) FocusBuilder // Uses enum
    WithChargeCategory(category pb.FocusChargeCategory) FocusBuilder   // Uses enum
    WithPricingCategory(category pb.FocusPricingCategory) FocusBuilder // Uses enum
    WithCurrency(currencyCode string) FocusBuilder
    WithBilledCost(cost float64) FocusBuilder
    WithExtension(key, value string) FocusBuilder // Fills the "Backpack"
    Build() (*pb.FocusCostRecord, error)
}

// Example Usage in plugin code
record, err := pluginsdk.NewFocusRecordBuilder().
    WithIdentity("AWS", "123456789012").
    WithResource("i-abcdef1234567890a", "EC2Instance").
    WithServiceCategory(pb.FocusServiceCategory_SERVICE_CATEGORY_COMPUTE).
    WithChargeCategory(pb.FocusChargeCategory_CHARGE_CATEGORY_USAGE).
    WithPricingCategory(pb.FocusPricingCategory_PRICING_CATEGORY_ON_DEMAND).
    WithCurrency("USD").
    WithBilledCost(0.45).
    WithExtension("custom_project_id", "project-x-alpha"). // Custom data
    Build()
```

## 5. Implementation Plan

### Phase 1: Schema Definition

- [ ] Create `proto/pulumicost/v1/focus.proto` with `extended_columns` and updated
      enum types for fields.
- [ ] Create `proto/pulumicost/v1/enums.proto` with the detailed enum definitions
      (FocusServiceCategory, FocusChargeCategory, FocusPricingCategory).
- [ ] Update `costsource.proto` to import both `focus.proto` and `enums.proto`.
- [ ] Run `buf generate`.

### Phase 2: SDK & Builders

- [ ] Create `sdk/go/pluginsdk/focus_builder.go`.
  - Implement the Builder pattern hiding the struct fields.
  - Ensure methods take enum types where appropriate.
  - Ensure `WithExtension` maps to `extended_columns`.
- [ ] Create `sdk/go/pluginsdk/focus_conformance.go` (`ValidateFocusRecord`).

### Phase 3: Documentation

- [ ] Create `PLUGIN_MIGRATION_GUIDE.md`.
  - **Explicitly Warning:** "Do not instantiate `FocusCostRecord` structs directly.
    Use the Builder. Direct struct usage is unsupported and may break in minor
    version updates."

## 6. Separation of Concerns

- **Spec:** Defines the data wire format (`extended_columns` handles the unknown)
  and strict enumerations.
- **SDK:** Defines the _Developer API_ (The Builder). This is the stability boundary
  and enforces enum usage.
- **Core:** Renders `extended_columns` dynamically, ensuring new data is visible immediately.
