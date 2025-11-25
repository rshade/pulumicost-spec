# Phase 0: Research & Technical Decisions

**Feature**: "What-If" Cost Estimation API
**Date**: 2025-11-24

## Research Tasks

### 1. Protobuf Message Design for EstimateCost

**Decision**: Use `google.protobuf.Struct` for attributes field

**Rationale**:
- Existing GetActualCost and GetProjectedCost RPCs already use Struct for flexibility
- Supports arbitrary nested key-value pairs matching Pulumi resource input properties
- Aligns with assumption: "The `google.protobuf.Struct` type is sufficient for representing all
  resource attribute types"
- Avoids need for provider-specific message types
- Allows plugins to interpret attributes based on their pricing models

**Alternatives Considered**:
- **Map<string, string>**: Rejected - Cannot represent nested structures or typed values
- **Custom AttributeValue oneof**: Rejected - Adds unnecessary complexity and proto message bloat
- **Provider-specific messages**: Rejected - Violates Constitution Principle II (Multi-Provider
  Consistency)

**Implementation Notes**:
- EstimateCostRequest will have fields: `string resource_type` (field 1), `google.protobuf.Struct
  attributes` (field 2)
- EstimateCostResponse will have fields: `string currency` (field 1), protobuf decimal type for
  `cost_monthly` (field 2)
- Field numbering follows best practice: most frequently accessed fields use numbers 1-15 for
  single-byte encoding

### 2. Decimal Type for Cost Representation

**Decision**: Research existing cost field types in GetActualCost/GetProjectedCost responses

**Findings**:
- Existing RPCs likely use `double` or `string` for cost values
- Financial precision requires avoiding floating-point rounding errors
- gRPC/Protobuf does not have native decimal type

**Decision**: Use approach consistent with existing GetActualCost and GetProjectedCost RPCs

**Rationale**:
- Maintains consistency across all cost-related RPCs
- Assumption states: "Decimal precision for costs follows existing patterns in GetActualCost and
  GetProjectedCost responses"
- Avoids introducing incompatible precision handling between different cost APIs

**Action Required**: Examine proto/pulumicost/v1/costsource.proto to determine exact type used in
existing cost response messages

### 3. Resource Type Validation Strategy

**Decision**: Implement format validation in SDK, not protobuf

**Rationale**:
- Protobuf string fields cannot enforce regex patterns
- FR-003 requires validation and specific gRPC InvalidArgument error responses
- Validation logic belongs in SDK layer, not proto layer
- Allows for clear, actionable error messages explaining format expectations

**Pattern**: `provider:module/resource:Type`
- **provider**: Cloud provider or platform (e.g., "aws", "azure", "gcp", "kubernetes")
- **module**: Provider module (e.g., "ec2", "compute", "k8s")
- **resource**: Resource name (e.g., "instance", "virtualMachine", "pod")
- **Type**: Pascal case resource type name

**Examples**:
- Valid: `aws:ec2/instance:Instance`
- Valid: `azure:compute/virtualMachine:VirtualMachine`
- Invalid: `aws:ec2:Instance` (missing module separator)
- Invalid: `aws/ec2/instance` (wrong separators)
- Invalid: `invalid-format` (doesn't match pattern)

**Implementation Location**: sdk/go/pricing/validate.go (new validation function)

### 4. Error Handling and gRPC Status Codes

**Decision**: Map requirement categories to gRPC status codes

**Mapping**:

| Scenario | Requirement | gRPC Status Code | Error Message Pattern |
|----------|-------------|------------------|----------------------|
| Invalid format | FR-003 | InvalidArgument | "resource type must follow provider:module/resource:Type format" |
| Unsupported resource | FR-008 | NotFound | "resource type {type} is not supported by this plugin" |
| Missing attributes | FR-009 | InvalidArgument | "missing required attributes: [{list}]" |
| Ambiguous attributes | FR-010 | InvalidArgument | "ambiguous or invalid attributes: {details}" |
| Pricing source unavailable | FR-014 | Unavailable | "pricing source unavailable: {reason}" |
| Zero cost (valid) | FR-013 | OK | Return successful response with cost=0 |

**Rationale**:
- Follows gRPC best practices for error handling
- InvalidArgument for client errors (bad input)
- NotFound for valid requests to unsupported resources
- Unavailable for transient server/external dependency failures
- Consistent with existing gRPC service error patterns

**Implementation Notes**:
- Use `status.Error()` or `status.Errorf()` in Go SDK
- Include structured error details using `google.rpc.ErrorInfo` when appropriate
- Error messages must be actionable (tell user what to fix)

### 5. Observability Integration

**Decision**: Leverage existing zerolog integration from spec 005-zerolog

**Requirements** (NFR-001 to NFR-003):
- Structured logs: request/response/errors with context
- Metrics: latency, success rate, error rate
- Distributed tracing: end-to-end visibility

**Implementation Approach**:
- Use `github.com/rs/zerolog` v1.34.0+ for structured logging
- Log at appropriate levels:
  - `Info`: Successful EstimateCost calls with resource_type and latency
  - `Warn`: Validation errors, unsupported resources
  - `Error`: Pricing source failures, unexpected errors
  - `Debug`: Full request/response details (attributes, cost breakdown)
- Include correlation IDs from gRPC metadata for tracing
- Metrics via existing registry (if present) or defer to plugin implementations

**Log Fields**:
- `rpc_method`: "EstimateCost"
- `resource_type`: The requested resource type
- `latency_ms`: Request duration
- `status`: "success" | "error"
- `error_code`: gRPC status code if error
- `error_message`: Error details if applicable

### 6. Testing Strategy

**Decision**: Extend existing testing framework with EstimateCost support

**Conformance Levels**:

**Basic (Required for all plugins)**:
- Successfully estimate cost for at least one supported resource type
- Return proper error for unsupported resource type
- Return proper error for invalid resource type format
- Handle null/missing attributes per FR-005

**Standard (Recommended for production)**:
- All Basic requirements plus:
- Estimate cost for multiple resource types (3+)
- Return descriptive errors for missing/ambiguous attributes
- Consistent/deterministic results for identical inputs (FR-011)
- Response time <500ms for standard resource types
- Handle 10+ concurrent requests

**Advanced (High-performance requirements)**:
- All Standard requirements plus:
- Support 50+ concurrent requests
- Response time <500ms even under load
- Proper error handling when pricing source unavailable
- Comprehensive observability signals (logs, metrics, traces)

**Test Implementation**:
- Extend `sdk/go/testing/harness.go` to support EstimateCost RPC
- Add mock implementations in `sdk/go/testing/mock_plugin.go`
- Add conformance tests in `sdk/go/testing/conformance_test.go`
- Add benchmarks in `sdk/go/testing/benchmark_test.go`
- Use bufconn for in-memory gRPC testing (existing pattern)

### 7. Example Payload Design

**Decision**: Create cross-provider request examples for Phase 1

**Required Examples** (per Constitution Principle II):

**AWS Example** (`examples/requests/estimate_cost_aws.json`):

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {
    "instanceType": "t3.micro",
    "region": "us-east-1",
    "tenancy": "default"
  }
}
```

**Azure Example** (`examples/requests/estimate_cost_azure.json`):

```json
{
  "resource_type": "azure:compute/virtualMachine:VirtualMachine",
  "attributes": {
    "vmSize": "Standard_B1s",
    "location": "eastus",
    "osType": "Linux"
  }
}
```

**GCP Example** (`examples/requests/estimate_cost_gcp.json`):

```json
{
  "resource_type": "gcp:compute/instance:Instance",
  "attributes": {
    "machineType": "e2-micro",
    "zone": "us-central1-a"
  }
}
```

**Response Example** (shared across providers):

```json
{
  "cost_monthly": "7.30",
  "currency": "USD"
}
```

**Rationale**:
- Demonstrates real-world usage patterns
- Shows attribute structure for each major provider
- Validates that `google.protobuf.Struct` handles provider-specific attributes
- Provides documentation through examples

## Summary

All technical decisions resolved. Key findings:

1. **Protobuf Design**: Use `google.protobuf.Struct` for attributes (consistent with existing RPCs)
2. **Decimal Type**: Follow existing GetActualCost/GetProjectedCost patterns (needs proto review)
3. **Validation**: Implement resource type format validation in SDK layer
4. **Error Handling**: Clear gRPC status code mapping with actionable error messages
5. **Observability**: Leverage zerolog integration with structured logging
6. **Testing**: Extend existing conformance framework with 3-level testing
7. **Examples**: Cross-provider request/response payloads for AWS, Azure, GCP

Ready to proceed to Phase 1 design.
