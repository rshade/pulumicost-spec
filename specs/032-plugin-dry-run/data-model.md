# Data Model: Plugin Capability Dry Run Mode

**Feature**: 032-plugin-dry-run
**Date**: 2025-12-31

## Entity Overview

```text
┌─────────────────────┐     ┌──────────────────────┐
│   DryRunRequest     │────▶│   DryRunResponse     │
├─────────────────────┤     ├──────────────────────┤
│ resource            │     │ field_mappings[]     │
│ simulation_params   │     │ configuration_valid  │
└─────────────────────┘     │ configuration_errors │
                            │ resource_supported   │
                            └──────────────────────┘
                                      │
                                      ▼ (repeated)
                            ┌──────────────────────┐
                            │    FieldMapping      │
                            ├──────────────────────┤
                            │ field_name           │
                            │ support_status       │
                            │ condition_description│
                            │ expected_type        │
                            └──────────────────────┘
                                      │
                                      ▼ (enum)
                            ┌──────────────────────┐
                            │ FieldSupportStatus   │
                            ├──────────────────────┤
                            │ UNSPECIFIED = 0      │
                            │ SUPPORTED = 1        │
                            │ UNSUPPORTED = 2      │
                            │ CONDITIONAL = 3      │
                            │ DYNAMIC = 4          │
                            └──────────────────────┘
```

## Entities

### FieldSupportStatus (Enum)

Represents the support status of a FOCUS field for a given resource type.

| Value | Numeric | Description |
|-------|---------|-------------|
| FIELD_SUPPORT_STATUS_UNSPECIFIED | 0 | Unknown/default status |
| FIELD_SUPPORT_STATUS_SUPPORTED | 1 | Field is always populated for this resource type |
| FIELD_SUPPORT_STATUS_UNSUPPORTED | 2 | Field is never populated for this resource type |
| FIELD_SUPPORT_STATUS_CONDITIONAL | 3 | Field is populated based on resource configuration |
| FIELD_SUPPORT_STATUS_DYNAMIC | 4 | Field requires runtime data to determine value |

**Validation Rules**:

- Must be one of the defined values (0-4)
- UNSPECIFIED (0) should only be used when status cannot be determined

### FieldMapping (Message)

Represents the support status for a single FOCUS field.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| field_name | string | Yes | FOCUS field identifier (e.g., "service_category", "billed_cost") |
| support_status | FieldSupportStatus | Yes | The support status enum value |
| condition_description | string | No | Human-readable explanation when status is CONDITIONAL or DYNAMIC |
| expected_type | string | No | Expected data type ("string", "double", "timestamp", "enum") |

**Validation Rules**:

- `field_name` must match a valid FocusCostRecord field name
- `condition_description` should be non-empty when `support_status` is CONDITIONAL or DYNAMIC
- `expected_type` should use proto3 type names or "enum" for enumerated fields

**Examples**:

```json
{
  "field_name": "service_category",
  "support_status": "FIELD_SUPPORT_STATUS_SUPPORTED",
  "condition_description": "",
  "expected_type": "enum"
}
```

```json
{
  "field_name": "availability_zone",
  "support_status": "FIELD_SUPPORT_STATUS_CONDITIONAL",
  "condition_description": "Only populated for regional resources in multi-AZ providers",
  "expected_type": "string"
}
```

### DryRunRequest (Message)

Request message for the DryRun RPC and dry_run flag on cost RPCs.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| resource | ResourceDescriptor | Yes | Resource type to query field mappings for |
| simulation_parameters | map<string, string> | No | Optional parameters to simulate different scenarios |

**Validation Rules**:

- `resource` must have valid `provider` and `resource_type` fields
- `simulation_parameters` is optional; unknown keys are ignored

**Reuses**: Existing `ResourceDescriptor` message from costsource.proto

### DryRunResponse (Message)

Response message containing field mapping information.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| field_mappings | repeated FieldMapping | Yes | Status for each known FOCUS field |
| configuration_valid | bool | Yes | Whether plugin configuration is valid |
| configuration_errors | repeated string | No | List of configuration error messages |
| resource_type_supported | bool | Yes | Whether the resource type is supported |

**Validation Rules**:

- `field_mappings` should contain entries for all known FOCUS fields (~50 fields)
- `configuration_errors` should be populated when `configuration_valid` is false
- When `resource_type_supported` is false, `field_mappings` may be empty

**State Transitions**: N/A (stateless introspection)

## Existing Entity Modifications

### GetActualCostRequest (Modification)

Add optional dry_run flag:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| dry_run | bool | No | When true, return DryRunResponse instead of cost data |

**Backward Compatibility**: Defaults to false, preserving existing behavior.

### GetProjectedCostRequest (Modification)

Add optional dry_run flag:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| dry_run | bool | No | When true, return DryRunResponse instead of cost data |

**Backward Compatibility**: Defaults to false, preserving existing behavior.

### GetActualCostResponse (Modification)

Add optional dry_run_result field:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| dry_run_result | DryRunResponse | No | Populated when request.dry_run is true |

### GetProjectedCostResponse (Modification)

Add optional dry_run_result field:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| dry_run_result | DryRunResponse | No | Populated when request.dry_run is true |

### SupportsResponse (Modification)

Document expected capability key:

```protobuf
// capabilities may include:
//   "dry_run": true  - Plugin supports DryRun RPC
```

## Relationships

```text
ResourceDescriptor ──(input)──▶ DryRunRequest
                                      │
                                      ▼
                              DryRunResponse
                                      │
                                      ▼ (1:many)
                              FieldMapping
                                      │
                                      ▼ (has-a)
                              FieldSupportStatus

FocusCostRecord ──(defines)──▶ Valid field_name values
```

## Data Volume Assumptions

- **Field count**: ~50-70 FOCUS fields per response (based on FocusCostRecord)
- **Response size**: ~5-10 KB per DryRunResponse (small, cacheable)
- **Request frequency**: Expected to be low (debugging/validation use case)
- **Caching**: Responses are cacheable by resource type (deterministic)
