# Feature Specification: FinOps FOCUS 1.2 Integration

**Feature Branch**: `009-focus-1-2-integration`
**Created**: Thursday, November 27, 2025
**Status**: Draft
**Input**: User description:

> The `pulumicost-spec` currently serves as a high-level query API for cloud costs.
> To achieve the project's vision of becoming the universal, open-source standard
> for cloud cost observability, it must align with the industry-standard FinOps
> FOCUS 1.2 Specification. Adopting FOCUS 1.2 will transform `pulumicost` from
> a simple cost summarizer into a forensic-grade cost analysis tool. To ensure
> this is sustainable and upgradeable (e.g., to FOCUS 1.3), we will employ a
> "Backpack & Builder" strategy to insulate plugin developers from schema complexity.

## Clarifications

### Session 2025-11-27

- Q: Preferred data type for financial fields?
  → A: Protobuf `double` (for performance and broad tool compatibility).
- Q: How should the Builder handle validation failures (e.g., missing mandatory fields)?
  → A: Return an explicit error (Go idiom: `(*Record, error)`).
- Q: Should `extended_columns` keys use a naming convention/namespace to prevent
  future collisions? → A: No namespace (allow any string for simplicity and
  plugin flexibility).
- **Contextual Clarification**: Plugins are understood to query external data sources
  and then construct `FocusCostRecord` objects for downstream consumption; they do
  not persist these records themselves.
- Q: How comprehensive should the validation of `FocusCostRecord` be?
  → A: FOCUS 1.2 Compliance (validate against FOCUS 1.2 business rules beyond
  just structural checks).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create Standardized Cost Records (Priority: P1)

As a Plugin Developer, I want to construct cost records using a standardized interface
that aligns with FinOps FOCUS 1.2, so that my plugin produces data compatible with
the broader FinOps ecosystem.

**Why this priority**: This is the core function of the integration. Without the
ability to create standard records, the spec update has no value.

**Independent Test**: Can be fully tested by creating a `FocusBuilder`, populating
it with standard FOCUS fields (e.g., ServiceCategory, BilledCost), and verifying
the output cost record contains the correct values.

**Acceptance Scenarios**:

1. **Given** a new cost plugin, **When** I use the `FocusBuilder` to set
   `ServiceCategory` to `COMPUTE`, **Then** the resulting record's `service_category`
   field equals the corresponding Standard Vocabulary value.
2. **Given** a cost plugin, **When** I set financial fields like `BilledCost` and
   `ListCost`, **Then** the resulting record contains these precise numeric values.
3. **Given** a cost plugin, **When** I attempt to set a category, **Then** the
   development tools restrict me to the defined Vocabulary values
   (e.g., `FocusServiceCategory`).
4. **Given** a builder instance, **When** I call `Build()` without setting mandatory
   fields (e.g., `ListCost`), **Then** it returns a non-nil error and a nil record.

---

### User Story 2 - Future-Proof Extension (The "Backpack") (Priority: P2)

As a Plugin Developer, I want to attach provider-specific or future-standard attributes
to a record that are not yet in the strict schema, so that I do not lose valuable
context (like custom project IDs or new FOCUS 1.3 fields).

**Why this priority**: Essential for real-world usage where data often exceeds the
current standard spec. Prevents the spec from becoming a blocker.

**Independent Test**: Verify that calling `.WithExtension("my-key", "my-value")`
results in the data being present in the serialized record's `extended_columns`
collection.

**Acceptance Scenarios**:

1. **Given** a cost record builder, **When** I add a custom extension
   "TeamOwner"="DevOps", **Then** the final record's `extended_columns` collection
   contains this key-value pair.
2. **Given** a record with extensions, **When** it is serialized and deserialized,
   **Then** the extension data is preserved exactly.

---

### User Story 3 - Stable Upgrade Path (The "Shield") (Priority: P3)

As a Project Maintainer, I want to be able to refactor the underlying data schema
(e.g., moving a field from "Extension" to "First Class") without breaking existing
plugin code.

**Why this priority**: Ensures long-term maintainability and encourages plugin
ecosystem growth by promising stability.

**Independent Test**: (Conceptual/Integration) Rename a field in the internal schema
but keep the Builder method signature the same. Existing plugin code should still
build and run.

**Acceptance Scenarios**:

1. **Given** plugin code using `FocusBuilder`, **When** the underlying schema
   definition is updated (e.g., adding a new field), **Then** the plugin code
   still builds without changes.
2. **Given** the Builder interface, **When** I implement it, **Then** I do not
   interact with the raw data struct directly.

### Edge Cases

- **Unknown/Unmapped Categories**: If a provider introduces a new Service Category
  not yet in the Enum, the Plugin Developer must map it to "Other" or "Unspecified"
  and store the raw value in `extended_columns` to preserve fidelity.
- **Missing Mandatory Fields**: If a Plugin Developer fails to populate a mandatory
  FOCUS 1.2 field (like `ListCost`), the `Build()` method MUST return an explicit
  error.
- **Precision Loss**: Financial values MUST use Protobuf `double` to balance precision
  with performance and ecosystem compatibility.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST define a `FocusCostRecord` data schema that implements
  the mandatory and optional columns specified in FinOps FOCUS 1.2.
- **FR-002**: System MUST define strict Controlled Vocabularies (Enums) for:
  `FocusServiceCategory`, `FocusChargeCategory`, and `FocusPricingCategory`.
- **FR-003**: System MUST support an `extended_columns` collection (key-value map)
  within the record to store attributes not defined in the strict schema.
- **FR-004**: System MUST provide an SDK `FocusBuilder` interface that encapsulates
  the construction of `FocusCostRecord` messages.
- **FR-005**: The `FocusBuilder` interface MUST provide strongly-typed methods for
  setting Vocabulary-based fields (e.g., `WithServiceCategory`).
- **FR-006**: The `FocusBuilder` interface MUST provide a method (e.g., `WithExtension`)
  to add arbitrary key-value pairs to the `extended_columns` collection.
- **FR-007**: The SDK MUST prevent (or strongly discourage via documentation and
  pattern) direct instantiation of the raw data struct by consumers.
- **FR-008**: System MUST include a conformance validation utility to verify that
  a generated record meets **FinOps FOCUS 1.2 compliance, including structural
  validity and defined business rules.**
- **FR-009**: The `Build()` method MUST return a tuple of `(*FocusCostRecord, error)`
  to enforce handling of validation failures (e.g., missing mandatory fields) at
  runtime.

### Key Entities

- **FocusCostRecord**: The core data structure representing a single line item of
  cost, aligned with the FinOps FOCUS 1.2 specification (contains Identity, Service,
  Charge, Financial [`double`], and Time details).
- **FocusServiceCategory**: Vocabulary representing the type of service
  (e.g., Compute, Storage, Database).
- **FocusChargeCategory**: Vocabulary representing the nature of the charge
  (e.g., Usage, Tax, Adjustment).
- **FocusPricingCategory**: Vocabulary representing the pricing model
  (e.g., On-Demand, Commitment, Spot).
- **FocusBuilder**: The opaque interface used by developers to construct records
  without coupling to the raw data layout.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of the mandatory FOCUS 1.2 columns are represented in the
  `FocusCostRecord` schema, and records conform to **FinOps FOCUS 1.2 business rules.**
- **SC-002**: Developers define records with 100% type safety for Service, Charge,
  and Pricing categories (build system enforces Vocabulary usage).
- **SC-003**: Arbitrary extension data added via the Builder is successfully
  serialized and retrievable from the `extended_columns` collection 100% of the time.
- **SC-004**: Existing plugin code requires Zero changes to build when the underlying
  schema is updated (simulated by interface adherence).
