# Feature Specification: FOCUS 1.3 Migration

**Feature Branch**: `026-focus-1-3-migration`
**Created**: 2025-12-23
**Status**: Draft
**Input**: GitHub Issue #183 - feat(focus): migrate to FOCUS 1.3 specification

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Add FOCUS 1.3 Cost Allocation Columns (Priority: P1)

As a plugin developer implementing cost data for shared resources (Kubernetes pods, database
instances), I need to populate allocation-specific columns so that FinOps practitioners can
understand how costs are split across workloads and see the methodology, not just the final
allocated amounts.

**Why this priority**: Split cost allocation is the primary pain point FOCUS 1.3 addresses.
Without these columns, practitioners must build custom allocation logic. This is the core value
proposition of FOCUS 1.3.

**Independent Test**: Can be fully tested by building a FocusCostRecord with allocation columns
and validating all fields are correctly populated and accessible.

**Acceptance Scenarios**:

1. **Given** a plugin returning cost data for a shared Kubernetes cluster, **When** the plugin
   sets allocated resource details, **Then** the record contains AllocatedResourceId,
   AllocatedResourceName, AllocatedMethodId, AllocatedMethodDetails, and AllocatedTags fields.
2. **Given** a FocusRecordBuilder instance, **When** the developer calls allocation-related
   builder methods, **Then** all allocation fields are correctly populated in the resulting
   proto message.
3. **Given** a cost record with allocation columns populated, **When** validation runs, **Then**
   the record passes FOCUS 1.3 conformance checks.

---

### User Story 2 - Add Service/Host Provider Columns (Priority: P1)

As a plugin developer working with multi-vendor billing (resellers, marketplace purchases), I
need to distinguish between the service provider (who sells the service) and the host provider
(where it runs) so that practitioners can accurately identify billing relationships and support
contacts.

**Why this priority**: The deprecated Provider/Publisher columns caused confusion in multi-vendor
scenarios. These new columns resolve definitional conflicts and are mandatory replacements.

**Independent Test**: Can be fully tested by building a FocusCostRecord with ServiceProviderName
and HostProviderName and verifying correct population and validation.

**Acceptance Scenarios**:

1. **Given** a cost record for an Azure Marketplace service running on Azure, **When** the
   plugin populates provider columns, **Then** ServiceProviderName reflects the ISV vendor and
   HostProviderName reflects "Azure".
2. **Given** a FocusRecordBuilder instance, **When** the developer calls WithServiceProvider()
   and WithHostProvider(), **Then** both fields are correctly populated.
3. **Given** the deprecated ProviderName field is populated, **When** validation runs with
   FOCUS 1.3 mode, **Then** a deprecation warning is logged but the record still validates
   (backward compatibility).

---

### User Story 3 - Support Contract Commitment Dataset (Priority: P2)

As a FinOps practitioner, I need to query contract commitment data separately from cost/usage
data so that I can see all active commitments, their terms, remaining units, and expiration
dates in a single query without parsing cost records.

**Why this priority**: Contract commitments were previously embedded in cost rows, making
analysis difficult. This is the first supplemental dataset in FOCUS, representing a significant
architectural addition.

**Independent Test**: Can be fully tested by creating a ContractCommitment message, populating
all fields, and verifying round-trip serialization.

**Acceptance Scenarios**:

1. **Given** a new ContractCommitment proto message, **When** the developer populates all 12
   fields, **Then** all fields serialize correctly and are accessible via generated Go code.
2. **Given** a contract commitment for a 3-year reserved instance, **When** the commitment is
   created, **Then** ContractCommitmentId, ContractId, ContractPeriodStart, ContractPeriodEnd,
   ContractCommitmentCost, and ContractCommitmentQuantity are all populated.
3. **Given** a ContractCommitmentBuilder, **When** Build() is called with all required fields,
   **Then** the commitment passes validation.

---

### User Story 4 - Add Contract Applied Column (Priority: P2)

As a plugin developer returning cost records, I need to indicate whether a contract applies to
a specific charge so that practitioners can link cost records to their contract commitment data.

**Why this priority**: This column connects the Cost/Usage dataset to the new Contract
Commitment dataset, enabling cross-dataset analysis.

**Independent Test**: Can be fully tested by setting ContractApplied on a FocusCostRecord and
verifying the link to contract data.

**Acceptance Scenarios**:

1. **Given** a cost record covered by a reserved instance contract, **When** ContractApplied is
   set, **Then** the field contains data linking to the Contract Commitment dataset.
2. **Given** a FocusRecordBuilder, **When** WithContractApplied() is called, **Then** the
   ContractApplied field is correctly populated.

---

### User Story 5 - Maintain Backward Compatibility (Priority: P3)

As a plugin developer with existing FOCUS 1.2 implementations, I need my current code to
continue working without modification so that I can adopt FOCUS 1.3 features incrementally.

**Why this priority**: Breaking existing plugins would cause adoption friction. All new columns
should be optional for backward compatibility.

**Independent Test**: Can be fully tested by running existing FOCUS 1.2 conformance tests and
verifying they still pass.

**Acceptance Scenarios**:

1. **Given** a FocusCostRecord built using only FOCUS 1.2 fields, **When** validation runs,
   **Then** the record passes validation (new FOCUS 1.3 fields are optional).
2. **Given** existing FOCUS 1.2 builder code, **When** compiled against the updated SDK,
   **Then** no compile errors occur and existing tests pass.
3. **Given** the deprecated ProviderName field, **When** populated alongside new
   ServiceProviderName/HostProviderName, **Then** both are accepted for transition period
   compatibility.

---

### Edge Cases

- When both deprecated ProviderName and new ServiceProviderName are populated with conflicting
  values: Log deprecation warning and prefer ServiceProviderName (new field wins).
- ContractApplied references to non-existent ContractCommitmentId: Accept as opaque reference
  (no cross-dataset validation; referential integrity is consumer's responsibility).
- AllocatedMethodId without AllocatedResourceId: Validation error (method requires resource).
- FOCUS 1.2 records processed by FOCUS 1.3 code: Pass validation (all new fields optional).

## Requirements _(mandatory)_

### Functional Requirements

#### Proto Definition Updates

- **FR-001**: System MUST add eight new columns to FocusCostRecord: AllocatedMethodId,
  AllocatedMethodDetails, AllocatedResourceId, AllocatedResourceName, AllocatedTags,
  ContractApplied, HostProviderName, ServiceProviderName.
- **FR-002**: System MUST create a new ContractCommitment proto message with 12 fields:
  ContractCommitmentId, ContractId, ContractCommitmentCategory, ContractCommitmentType,
  ContractCommitmentPeriodStart, ContractCommitmentPeriodEnd, ContractPeriodStart,
  ContractPeriodEnd, ContractCommitmentCost, ContractCommitmentQuantity, ContractCommitmentUnit,
  BillingCurrency.
- **FR-003**: System MUST add deprecation comments to ProviderName and Publisher fields in
  FocusCostRecord proto.
- **FR-004**: System MUST use proto field numbers 59+ for new FocusCostRecord columns to
  maintain backward compatibility.

#### Go SDK Builder Updates

- **FR-005**: System MUST add builder methods for all new FocusCostRecord columns:
  WithAllocation(), WithAllocatedResource(), WithAllocatedTags(), WithContractApplied(),
  WithServiceProvider(), WithHostProvider().
- **FR-006**: System MUST create a ContractCommitmentBuilder for constructing contract
  commitment records.
- **FR-007**: Builder methods MUST follow existing naming conventions and fluent API patterns.

#### Validation Updates

- **FR-008**: System MUST validate new FOCUS 1.3 columns according to specification
  requirements (mandatory, recommended, conditional classifications).
- **FR-009**: System MUST log deprecation warnings when deprecated fields (ProviderName,
  Publisher) are populated.
- **FR-010**: System MUST maintain validation compatibility for FOCUS 1.2 records (new columns
  optional).
- **FR-010a**: System MUST require AllocatedResourceId when AllocatedMethodId is populated
  (allocation method requires an allocated resource target).

#### Enum Additions

- **FR-011**: System MUST add FocusContractCommitmentCategory enum if required by
  specification.
- **FR-012**: System MUST add FocusContractCommitmentType enum if required by specification.

#### Testing Updates

- **FR-013**: System MUST add conformance tests for all new FOCUS 1.3 columns.
- **FR-014**: System MUST add integration tests for ContractCommitment dataset.
- **FR-015**: System MUST verify existing FOCUS 1.2 conformance tests continue to pass.
- **FR-016**: System MUST add benchmarks for new builder methods.

#### Documentation Updates

- **FR-017**: System MUST update focus-columns.md with FOCUS 1.3 column documentation.
- **FR-018**: System MUST update FocusCostRecord proto comments to reference FOCUS 1.3
  specification.
- **FR-019**: System MUST document the new ContractCommitment dataset and its relationship
  to FocusCostRecord.

### Key Entities

- **FocusCostRecord**: Extended with 8 new columns for allocation and provider identification.
  Represents a single cost line item in the Cost and Usage dataset.
- **ContractCommitment**: New entity representing contractual commitment terms. Contains 12
  fields covering commitment identification, temporal bounds, and financial metrics. Links to
  FocusCostRecord via ContractApplied column.
- **AllocatedResource**: Conceptual entity representing the resource receiving an allocated
  charge. Identified by AllocatedResourceId and AllocatedResourceName, with associated tags.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All 8 new FocusCostRecord columns are accessible through the Go SDK builder API.
- **SC-002**: ContractCommitment proto message compiles and generates valid Go code with all
  12 fields.
- **SC-003**: 100% of existing FOCUS 1.2 conformance tests pass without modification.
- **SC-004**: New FOCUS 1.3 conformance tests achieve 100% coverage of new columns and
  validation rules.
- **SC-005**: Builder operations for new columns complete in under 100 nanoseconds per
  operation (matching existing performance baselines).
- **SC-006**: Documentation covers all new columns with examples matching the quality of
  existing FOCUS 1.2 documentation.
- **SC-007**: Deprecation warnings appear in logs when deprecated columns (ProviderName,
  Publisher) are used.
- **SC-008**: Plugin developers can incrementally adopt FOCUS 1.3 features without rewriting
  existing FOCUS 1.2 code.

## Assumptions

- FOCUS 1.3 specification is finalized (ratified December 5, 2025) and column definitions are
  stable.
- New columns follow the same conditional/mandatory classification patterns as FOCUS 1.2.
- The ContractCommitment dataset is queried separately (not embedded in cost RPC responses).
- AllocatedTags follows the same map<string, string> pattern as existing Tags field.
- ContractApplied column contains structured data (likely JSON) linking to ContractCommitmentId.
- Deprecation period for ProviderName and Publisher extends through FOCUS 1.3 (removed in 1.4).

## Out of Scope

- Recency and completeness metadata (dataset-level, not record-level).
- Changes to existing RPC definitions (GetActualCost, GetProjectedCost, etc.).
- New RPC methods for querying ContractCommitment data (may be addressed in future spec).
- Automatic migration of existing FOCUS 1.2 data to FOCUS 1.3 format.
- UI/CLI changes for visualizing contract commitments.

## Clarifications

### Session 2025-12-23

- Q: When both deprecated ProviderName and new ServiceProviderName are populated with conflicting
  values, how should validation behave? → A: Warn and prefer ServiceProviderName (new field wins)
- Q: How should the system handle a ContractApplied field that references a ContractCommitmentId
  that doesn't exist? → A: Accept as opaque reference (no cross-dataset validation)
- Q: When AllocatedMethodId is set but AllocatedResourceId is NOT set, how should validation
  behave? → A: Require AllocatedResourceId when AllocatedMethodId is set (validation error)

## References

- [FOCUS 1.3 Announcement](https://www.finops.org/insights/introducing-focus-1-3/)
- [FOCUS Specification](https://focus.finops.org/)
- [FOCUS GitHub Repository](https://github.com/FinOps-Open-Cost-and-Usage-Spec/FOCUS_Spec)
- GitHub Issue #183: feat(focus): migrate to FOCUS 1.3 specification
- Related PR: rshade/finfocus-spec#99 - feat(focus): Implement FOCUS 1.2 integration [merged]
