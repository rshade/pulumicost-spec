# Research: GreenOps Integration

## Decisions

### 1. Protobuf Enum for MetricKind

- **Decision**: Define `MetricKind` as a top-level enum in `proto/pulumicost/v1/costsource.proto`.
- **Rationale**: It is a core part of the capability discovery and cost projection logic.
- **Values**:
  - `METRIC_KIND_UNSPECIFIED = 0`
  - `METRIC_KIND_CARBON_FOOTPRINT = 1` (Unit: gCO2e)
  - `METRIC_KIND_ENERGY_CONSUMPTION = 2` (Unit: kWh)
  - `METRIC_KIND_WATER_USAGE = 3` (Unit: L)

### 2. Utilization Percentage Field Type and Scope

- **Decision**: Use `double` for `utilization_percentage` (0.0 to 1.0).
- **Placement**:
  - Global: `GetProjectedCostRequest.utilization_percentage`
  - Override: `ResourceDescriptor.utilization_percentage` (as `optional double`)
- **Rationale**: Protobuf `double` is consistent with existing cost fields. The tiered approach
  (global + override) provides maximum flexibility as requested.

### 3. Error Handling for Unsupported/Unavailable Metrics

- **Decision**: Omit metrics from the response if they cannot be calculated, even if advertised as supported.
- **Rationale**: Idiomatic Protobuf approach to avoid ambiguity with zero values.

## Alternatives Considered

### 1. Separate Service for GreenOps

- **Evaluation**: Rejected.
- **Reason**: GreenOps metrics are intrinsically linked to cost projections and resource discovery.
  Adding a new service would increase complexity without clear benefit.

### 2. Use google.protobuf.FloatValue for Nullability

- **Evaluation**: Rejected for global field, but considered for override.
- **Reason**: `optional double` in proto3 provides the necessary presence detection without the overhead of wrapper types.

## Dependencies & Best Practices

- **Buf**: Use `buf lint` and `buf generate` to ensure proto quality and SDK consistency.
- **Units**: Adhere to SI units (Liters, kWh) and standard sustainability reporting units (gCO2e).
