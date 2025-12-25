# Research: FOCUS 1.3 Migration

**Branch**: `026-focus-1-3-migration` | **Date**: 2025-12-23

## Research Tasks

### 1. FOCUS 1.3 Column Specifications

**Task**: Research exact FOCUS 1.3 column definitions, data types, and requirements.

**Decision**: Use official FOCUS 1.3 specification column definitions with these mappings:

| FOCUS Column | Proto Type | FOCUS Requirement |
|--------------|-----------|-------------------|
| AllocatedMethodId | string | Conditional |
| AllocatedMethodDetails | string | Recommended |
| AllocatedResourceId | string | Conditional |
| AllocatedResourceName | string | Conditional |
| AllocatedTags | map<string, string> | Conditional |
| ContractApplied | string | Conditional |
| ServiceProviderName | string | Conditional |
| HostProviderName | string | Conditional |

**Rationale**: These match the official FOCUS 1.3 specification released December 2025. All new
columns are conditional (not mandatory) to maintain backward compatibility.

**Alternatives considered**:

- Make some columns mandatory: Rejected - would break FOCUS 1.2 backward compatibility
- Use different data types: Rejected - must follow FOCUS specification exactly

### 2. ContractCommitment Dataset Schema

**Task**: Research the 12-field ContractCommitment dataset structure.

**Decision**: Implement as a separate proto message with these fields:

| Field | Proto Type | Purpose |
|-------|-----------|---------|
| ContractCommitmentId | string | Unique identifier for commitment |
| ContractId | string | Parent contract identifier |
| ContractCommitmentCategory | enum | Category (e.g., Spend, Usage) |
| ContractCommitmentType | string | Type of commitment |
| ContractCommitmentPeriodStart | google.protobuf.Timestamp | Commitment period start |
| ContractCommitmentPeriodEnd | google.protobuf.Timestamp | Commitment period end |
| ContractPeriodStart | google.protobuf.Timestamp | Contract period start |
| ContractPeriodEnd | google.protobuf.Timestamp | Contract period end |
| ContractCommitmentCost | double | Financial commitment amount |
| ContractCommitmentQuantity | double | Committed quantity |
| ContractCommitmentUnit | string | Unit for quantity |
| BillingCurrency | string | ISO 4217 currency code |

**Rationale**: This structure matches the FOCUS 1.3 Contract Commitment dataset specification.
Using a separate proto message (not embedded in FocusCostRecord) reflects FOCUS's design of
decoupled datasets.

**Alternatives considered**:

- Embed in FocusCostRecord: Rejected - FOCUS 1.3 explicitly separates these datasets
- Use nested messages: Rejected - simpler flat structure aligns with FOCUS patterns

### 3. Deprecation Handling Pattern

**Task**: Research best practices for protobuf field deprecation.

**Decision**: Use protobuf `deprecated` option and application-level warnings:

```protobuf
// Deprecated: Use service_provider_name instead. Will be removed in FOCUS 1.4.
string provider_name = 1 [deprecated = true];

// Deprecated: Use host_provider_name instead. Will be removed in FOCUS 1.4.
string publisher = 55 [deprecated = true];
```

**Rationale**: The `deprecated` option generates compiler warnings in generated code. Combined
with application-level logging (FR-009), this gives developers clear migration signals.

**Alternatives considered**:

- Remove fields immediately: Rejected - breaks wire compatibility
- Only add comments: Rejected - generated code wouldn't warn developers

### 4. Proto Field Numbering Strategy

**Task**: Determine field numbers for new FocusCostRecord columns.

**Decision**: Use field numbers 59-66 for new columns:

| Field | Number | Rationale |
|-------|--------|-----------|
| service_provider_name | 59 | First new field |
| host_provider_name | 60 | Provider pair |
| allocated_method_id | 61 | Allocation group start |
| allocated_method_details | 62 | Allocation group |
| allocated_resource_id | 63 | Allocation group |
| allocated_resource_name | 64 | Allocation group |
| allocated_tags | 65 | Allocation metadata |
| contract_applied | 66 | Contract link |

**Rationale**: Current highest field number is 58 (sku_price_details). Using 59+ maintains
backward compatibility - existing serialized messages remain valid.

**Alternatives considered**:

- Use gaps in existing numbering: Rejected - existing gaps may have semantic meaning
- Reserve large range: Not needed - FOCUS spec is stable

### 5. AllocatedTags Data Structure

**Task**: Research whether AllocatedTags should match existing Tags structure.

**Decision**: Use `map<string, string>` matching the existing Tags field pattern.

**Rationale**: Consistency with existing `tags` field (field 22) in FocusCostRecord. Same
pattern used throughout the codebase for key-value metadata.

**Alternatives considered**:

- Use repeated TagEntry message: Rejected - inconsistent with existing pattern
- Use JSON string: Rejected - loses type safety

### 6. ContractApplied Data Format

**Task**: Research the ContractApplied column format for linking datasets.

**Decision**: Use a simple string containing the ContractCommitmentId value.

**Rationale**: Per clarification session, ContractApplied is treated as an opaque reference
with no cross-dataset validation. A string ID is sufficient for linking and follows FOCUS's
approach of decoupled datasets where referential integrity is the consumer's responsibility.

**Alternatives considered**:

- JSON object with multiple IDs: Rejected - overcomplicated for opaque reference
- Repeated field for multiple contracts: Rejected - FOCUS spec uses singular

### 7. Enum Requirements

**Task**: Research whether new enums are needed for FOCUS 1.3.

**Decision**: Add one new enum for ContractCommitmentCategory:

```protobuf
enum FocusContractCommitmentCategory {
  FOCUS_CONTRACT_COMMITMENT_CATEGORY_UNSPECIFIED = 0;
  FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND = 1;
  FOCUS_CONTRACT_COMMITMENT_CATEGORY_USAGE = 2;
}
```

**Rationale**: ContractCommitmentCategory is explicitly enumerated in FOCUS 1.3 with Spend
and Usage as the defined values. ContractCommitmentType is a free-form string per spec.

**Alternatives considered**:

- Use string for category: Rejected - FOCUS explicitly enumerates values
- Add ContractCommitmentType enum: Rejected - spec allows free-form strings

### 8. Validation Rule Dependencies

**Task**: Research validation dependencies between new fields.

**Decision**: Implement one validation dependency per clarification session:

- **AllocatedMethodId requires AllocatedResourceId**: When AllocatedMethodId is populated,
  AllocatedResourceId MUST also be populated. Rationale: An allocation method without an
  allocation target is meaningless.

No other field dependencies identified - all other new fields are independently optional.

**Rationale**: This was clarified during spec development. The dependency is logical -
you can't describe HOW costs are allocated without identifying WHAT they're allocated to.

**Alternatives considered**:

- Require all allocation fields together: Rejected - too restrictive
- No dependencies: Rejected - allows nonsensical data

## Summary

All technical unknowns have been resolved:

1. **Column definitions**: 8 new FocusCostRecord columns mapped with types and requirements
2. **ContractCommitment**: 12-field separate proto message structure defined
3. **Deprecation**: Use protobuf `deprecated` option + application warnings
4. **Field numbers**: 59-66 for new columns, preserving backward compatibility
5. **AllocatedTags**: map<string, string> matching existing Tags pattern
6. **ContractApplied**: Simple string reference (opaque, no validation)
7. **Enums**: One new enum (FocusContractCommitmentCategory)
8. **Validation**: AllocatedMethodId requires AllocatedResourceId

Ready to proceed to Phase 1: Design & Contracts.
