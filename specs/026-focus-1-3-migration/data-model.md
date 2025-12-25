# Data Model: FOCUS 1.3 Migration

**Branch**: `026-focus-1-3-migration` | **Date**: 2025-12-23

## Entity Overview

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FocusCostRecord                                    │
│  (Extended with 8 new FOCUS 1.3 columns)                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ EXISTING FIELDS (1-58)                                                       │
│ ├── provider_name [1] ⚠️ DEPRECATED                                          │
│ ├── publisher [55] ⚠️ DEPRECATED                                             │
│ └── ... (56 other existing fields)                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│ NEW FOCUS 1.3 FIELDS (59-66)                                                │
│ ├── service_provider_name [59] - string                                     │
│ ├── host_provider_name [60] - string                                        │
│ ├── allocated_method_id [61] - string                                       │
│ ├── allocated_method_details [62] - string                                  │
│ ├── allocated_resource_id [63] - string                                     │
│ ├── allocated_resource_name [64] - string                                   │
│ ├── allocated_tags [65] - map<string, string>                               │
│ └── contract_applied [66] - string                                          │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Links via contract_applied
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ContractCommitment                                   │
│  (NEW supplemental dataset - 12 fields)                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│ IDENTITY                                                                     │
│ ├── contract_commitment_id [1] - string (primary identifier)                │
│ └── contract_id [2] - string (parent contract)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ CLASSIFICATION                                                               │
│ ├── contract_commitment_category [3] - enum                                 │
│ └── contract_commitment_type [4] - string                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│ COMMITMENT PERIOD                                                            │
│ ├── contract_commitment_period_start [5] - timestamp                        │
│ └── contract_commitment_period_end [6] - timestamp                          │
├─────────────────────────────────────────────────────────────────────────────┤
│ CONTRACT PERIOD                                                              │
│ ├── contract_period_start [7] - timestamp                                   │
│ └── contract_period_end [8] - timestamp                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ FINANCIAL                                                                    │
│ ├── contract_commitment_cost [9] - double                                   │
│ ├── contract_commitment_quantity [10] - double                              │
│ ├── contract_commitment_unit [11] - string                                  │
│ └── billing_currency [12] - string (ISO 4217)                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

## FocusCostRecord Extensions

### New Fields (FOCUS 1.3)

| Field | Number | Type | Requirement | Description |
|-------|--------|------|-------------|-------------|
| service_provider_name | 59 | string | Conditional | Provider making service available for purchase |
| host_provider_name | 60 | string | Conditional | Provider hosting the underlying resource |
| allocated_method_id | 61 | string | Conditional | Identifier for allocation methodology |
| allocated_method_details | 62 | string | Recommended | Description of allocation methodology |
| allocated_resource_id | 63 | string | Conditional | ID of resource receiving allocated cost |
| allocated_resource_name | 64 | string | Conditional | Name of resource receiving allocated cost |
| allocated_tags | 65 | map | Conditional | Tags associated with allocated resource |
| contract_applied | 66 | string | Conditional | Reference to ContractCommitmentId |

### Deprecated Fields

| Field | Number | Replacement | Removal Version |
|-------|--------|-------------|-----------------|
| provider_name | 1 | service_provider_name | FOCUS 1.4 |
| publisher | 55 | host_provider_name | FOCUS 1.4 |

### Validation Rules

1. **Allocation Field Dependency**: If `allocated_method_id` is populated, `allocated_resource_id`
   MUST also be populated. Validation error otherwise.

2. **Deprecated Field Handling**: If `provider_name` is populated AND `service_provider_name` is
   populated with a different value, log deprecation warning and prefer `service_provider_name`.

3. **Backward Compatibility**: All new fields are optional. FOCUS 1.2 records (without new
   fields) MUST pass validation.

## ContractCommitment Entity

### Fields

| Field | Number | Type | Description |
|-------|--------|------|-------------|
| contract_commitment_id | 1 | string | Unique identifier for this commitment |
| contract_id | 2 | string | Parent contract identifier |
| contract_commitment_category | 3 | enum | SPEND or USAGE |
| contract_commitment_type | 4 | string | Provider-specific commitment type |
| contract_commitment_period_start | 5 | timestamp | When commitment period begins |
| contract_commitment_period_end | 6 | timestamp | When commitment period ends |
| contract_period_start | 7 | timestamp | When contract begins |
| contract_period_end | 8 | timestamp | When contract ends |
| contract_commitment_cost | 9 | double | Monetary commitment amount |
| contract_commitment_quantity | 10 | double | Committed quantity |
| contract_commitment_unit | 11 | string | Unit for quantity measurement |
| billing_currency | 12 | string | ISO 4217 currency code |

### Enum: FocusContractCommitmentCategory

| Value | Number | Description |
|-------|--------|-------------|
| UNSPECIFIED | 0 | Default/unknown |
| SPEND | 1 | Monetary spend commitment |
| USAGE | 2 | Usage quantity commitment |

### Validation Rules

1. **Required Fields**: contract_commitment_id, contract_id, billing_currency
2. **Period Consistency**: contract_commitment_period_end >= contract_commitment_period_start
3. **Currency Format**: billing_currency must be valid ISO 4217 code (reuse existing validation)
4. **Non-negative Values**: contract_commitment_cost >= 0, contract_commitment_quantity >= 0

## Relationships

### FocusCostRecord → ContractCommitment

- **Link Field**: `contract_applied` in FocusCostRecord
- **Link Type**: Opaque reference (string containing ContractCommitmentId)
- **Validation**: No cross-dataset validation (referential integrity is consumer's responsibility)
- **Cardinality**: Many FocusCostRecords can reference one ContractCommitment

### AllocatedResource (Conceptual)

The allocated resource is not a separate entity but represented by fields within FocusCostRecord:

- `allocated_resource_id` - Identifier
- `allocated_resource_name` - Display name
- `allocated_tags` - Associated metadata

This follows FOCUS 1.3's design of embedded allocation data rather than a separate dataset.

## State Transitions

### FocusCostRecord Lifecycle

```text
┌──────────────┐
│   Created    │ ─── All mandatory FOCUS 1.2 fields set
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Extended   │ ─── Optional FOCUS 1.3 fields added
└──────┬───────┘     (allocation, provider, contract)
       │
       ▼
┌──────────────┐
│  Validated   │ ─── Passes conformance checks
└──────────────┘
```

### ContractCommitment Lifecycle

```text
┌──────────────┐
│   Created    │ ─── Identity fields set (ids, category)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Bounded    │ ─── Period fields set (start/end dates)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Quantified  │ ─── Financial fields set (cost, quantity)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Validated   │ ─── Passes validation rules
└──────────────┘
```

## Proto Wire Format Compatibility

### Backward Compatibility Guarantee

- Field numbers 59-66 are new and don't conflict with existing fields
- Old clients reading new messages will ignore unknown field numbers
- New clients reading old messages will see empty/default values for new fields
- No `reserved` statements needed (no fields being removed)

### Forward Compatibility

- New messages can be read by old clients (unknown fields ignored)
- Deprecated fields (1, 55) remain in proto for transition period
- `deprecated = true` option warns at compile time
