# Phase 1: Data Model

**Feature**: "What-If" Cost Estimation API
**Date**: 2025-11-24

## Protobuf Messages

### EstimateCostRequest

**Purpose**: Request message for estimating the cost of a Pulumi resource before deployment.

**Fields**:

| Field         | Type                   | Number | Required | Description                     | Validation                                          |
| ------------- | ---------------------- | ------ | -------- | ------------------------------- | --------------------------------------------------- |
| resource_type | string                 | 1      | Yes      | Pulumi resource type identifier | Must match pattern: `provider:module/resource:Type` |
| attributes    | google.protobuf.Struct | 2      | No       | Resource input properties       | Treated as empty struct if null/missing             |

**Validation Rules**:

- `resource_type` MUST NOT be empty
- `resource_type` format MUST match `provider:module/resource:Type` pattern
  - **provider**: lowercase identifier (e.g., "aws", "azure", "gcp")
  - **module**: lowercase module name (e.g., "ec2", "compute", "storage")
  - **resource**: lowercase resource name (e.g., "instance", "bucket")
  - **Type**: Pascal case type name (e.g., "Instance", "Bucket")
- `attributes` MAY be null or missing (interpreted as empty struct per FR-005)
- `attributes` structure is plugin-specific and validated by the plugin implementation

**Example** (AWS EC2 Instance):

```protobuf
EstimateCostRequest {
  resource_type: "aws:ec2/instance:Instance"
  attributes: {
    fields: {
      key: "instanceType"
      value: { string_value: "t3.micro" }
    }
    fields: {
      key: "region"
      value: { string_value: "us-east-1" }
    }
  }
}
```

### EstimateCostResponse

**Purpose**: Response message containing the estimated monthly cost for a resource.

**Fields**:

| Field        | Type               | Number | Required | Description            | Constraints                                    |
| ------------ | ------------------ | ------ | -------- | ---------------------- | ---------------------------------------------- |
| currency     | string             | 1      | Yes      | ISO 4217 currency code | Typically "USD", must be uppercase             |
| cost_monthly | [decimal type TBD] | 2      | Yes      | Estimated monthly cost | Non-negative, precision per existing cost RPCs |

**Validation Rules**:

- `currency` MUST be a valid ISO 4217 currency code (3 uppercase letters)
- `cost_monthly` MUST be non-negative (≥0)
- `cost_monthly` precision follows existing GetActualCost/GetProjectedCost patterns
- Zero cost is valid (e.g., free tier resources per FR-013)

**Example** (Successful Estimation):

```protobuf
EstimateCostResponse {
  currency: "USD"
  cost_monthly: "7.30"  // Type depends on existing RPC pattern
}
```

**Example** (Zero Cost - Free Tier):

```protobuf
EstimateCostResponse {
  currency: "USD"
  cost_monthly: "0.00"
}
```

## RPC Service Method

### EstimateCost

**Service**: `pulumicost.v1.CostSource`

**Signature**:

```protobuf
rpc EstimateCost(EstimateCostRequest) returns (EstimateCostResponse);
```

**Semantics**:

- **Idempotent**: YES - same inputs always produce same outputs (FR-011)
- **Streaming**: NO - simple unary RPC
- **Timeout**: Implementations should complete within 500ms (SC-002)

**Error Responses**:

| Scenario                    | gRPC Status     | Error Message Pattern                                                          |
| --------------------------- | --------------- | ------------------------------------------------------------------------------ |
| Empty resource_type         | InvalidArgument | "resource_type is required"                                                    |
| Invalid format              | InvalidArgument | "resource_type must follow provider:module/resource:Type format, got: {input}" |
| Unsupported resource        | NotFound        | "resource type {type} is not supported by this plugin"                         |
| Missing required attributes | InvalidArgument | "missing required attributes for {resource_type}: [{attribute_names}]"         |
| Ambiguous attributes        | InvalidArgument | "ambiguous or invalid attributes for {resource_type}: {details}"               |
| Pricing source unavailable  | Unavailable     | "pricing source unavailable: {reason}"                                         |
| Internal error              | Internal        | "internal error during cost estimation: {details}"                             |

## Data Relationships

```text
┌─────────────────────────────────┐
│    EstimateCostRequest          │
│                                 │
│  + resource_type: string        │───────┐
│  + attributes: Struct           │       │ Input to
│                                 │       │
└─────────────────────────────────┘       │
                                           │
                                           ▼
                                  ┌────────────────┐
                                  │   Plugin       │
                                  │ Implementation │
                                  └────────────────┘
                                           │
                                           │ Returns
                                           ▼
┌─────────────────────────────────┐
│   EstimateCostResponse          │
│                                 │
│  + currency: string             │
│  + cost_monthly: decimal        │
│                                 │
└─────────────────────────────────┘
```

## Resource Type Format

**Pattern**: `provider:module/resource:Type`

**Components**:

```text
aws:ec2/instance:Instance
└┬┘ └┬┘ └──┬───┘ └──┬───┘
 │   │     │        └─ Type (PascalCase)
 │   │     └─────────── resource (lowercase)
 │   └───────────────── module (lowercase)
 └───────────────────── provider (lowercase)
```

**Examples by Provider**:

| Provider   | Resource Type                                 | Description           |
| ---------- | --------------------------------------------- | --------------------- |
| AWS        | `aws:ec2/instance:Instance`                   | EC2 compute instance  |
| AWS        | `aws:s3/bucket:Bucket`                        | S3 storage bucket     |
| Azure      | `azure:compute/virtualMachine:VirtualMachine` | Azure VM              |
| Azure      | `azure:storage/account:Account`               | Azure storage account |
| GCP        | `gcp:compute/instance:Instance`               | GCE compute instance  |
| GCP        | `gcp:storage/bucket:Bucket`                   | Cloud Storage bucket  |
| Kubernetes | `kubernetes:core/v1/pod:Pod`                  | Kubernetes pod        |

## Attribute Structure

**Type**: `google.protobuf.Struct` (arbitrary key-value pairs)

**Common Patterns**:

### Compute Resources

```json
{
  "instanceType": "t3.micro", // AWS
  "machineType": "e2-micro", // GCP
  "vmSize": "Standard_B1s", // Azure
  "region": "us-east-1",
  "zone": "us-central1-a",
  "osType": "Linux"
}
```

### Storage Resources

```json
{
  "storageClass": "STANDARD",
  "region": "us-east-1",
  "sizeGB": 100,
  "replication": "GRS"
}
```

**Plugin Responsibility**:

- Plugins interpret attributes based on their pricing models
- Plugins validate required attributes and return InvalidArgument if missing
- Plugins handle provider-specific naming conventions (instanceType vs machineType vs vmSize)

## State Transitions

N/A - EstimateCost is a stateless, idempotent query operation. No state is persisted.

## Validation Summary

**SDK Layer (Go)**:

- Resource type format validation (regex match)
- Empty resource type check
- Null/missing attributes normalization to empty struct

**Plugin Layer (Implementation-specific)**:

- Resource type support check (Supports RPC)
- Required attribute presence check
- Attribute value validation (types, ranges, enums)
- Pricing model applicability check

**Proto Layer (buf)**:

- Message structure validation
- Field type validation
- Breaking change detection

## Future Considerations

**Not in Scope** (defer to future specs):

- Batch estimation (multiple resources in one RPC)
- Cost breakdown by component (compute, storage, network)
- Historical cost estimates (time-series data)
- Currency conversion
- Savings plan / reserved instance pricing
