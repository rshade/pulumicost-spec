# Feature Specification: Enhancing Support Discovery & Request Options for GreenOps Integration

**Feature Branch**: `020-greenops-integration`  
**Created**: 2025-12-18  
**Status**: Draft  
**Input**: User description: "Enhancing Support Discovery & Request Options for GreenOps Integration: Update SupportsResponse to advertise supported metrics and GetProjectedCostRequest to include utilization_percentage."

## Clarifications

### Session 2025-12-18

- Q: Which specific GreenOps metrics should be included in the initial MetricKind enumeration? → A: Carbon Footprint, Energy Consumption, and Water Usage
- Q: What is the preferred Protobuf type for the utilization_percentage field? → A: Protobuf double (0.0 to 1.0)
- Q: What should be the standard unit for the new Water Usage metric? → A: Liters (L)
- Q: How should a plugin handle an advertised metric that cannot be calculated for a specific request? → A: Omit the metric from response
- Q: What should be the scope of the utilization_percentage field? → A: Global default with per-resource override

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Capability Discovery for GreenOps (Priority: P1)

As a cloud cost engine, I want to know which sustainability metrics a plugin supports for a specific resource, so that I can automatically configure the user interface with the appropriate data columns.

**Why this priority**: This is foundational for the user experience. Without knowing what metrics are available, the engine cannot display GreenOps data effectively.

**Independent Test**: Can be tested by calling the `Supports` RPC on a plugin and verifying that the `supported_metrics` field in the response contains the expected `MetricKind` values.

**Acceptance Scenarios**:

1. **Given** a plugin that supports Carbon Footprint and Energy Consumption, **When** the engine calls `Supports`, **Then** the response MUST include `METRIC_KIND_CARBON_FOOTPRINT` and `METRIC_KIND_ENERGY_CONSUMPTION` in the `supported_metrics` list.
2. **Given** a plugin that does not support any sustainability metrics, **When** the engine calls `Supports`, **Then** the `supported_metrics` list MUST be empty.

---

### User Story 2 - Accurate Impact Modeling via Utilization (Priority: P2)

As a user/engine, I want to provide a specific utilization assumption when requesting projected costs, so that the sustainability calculations (like CCF - Cloud Carbon Footprint) are more accurate for my specific workload.

**Why this priority**: Sustainability metrics are highly sensitive to utilization. Providing this parameter allows for much more precise impact modeling than generic averages.

**Independent Test**: Can be tested by sending a `GetProjectedCostRequest` with a specific `utilization_percentage` and verifying that the plugin receives this value and incorporates it into its response (e.g., higher utilization leading to higher carbon estimates).

**Acceptance Scenarios**:

1. **Given** a `GetProjectedCostRequest` with `utilization_percentage` set to 0.75, **When** the plugin processes the request, **Then** it MUST use 75% utilization as the basis for its sustainability calculations.
2. **Given** a `GetProjectedCostRequest` with no `utilization_percentage` specified, **When** the plugin processes the request, **Then** it SHOULD use a sensible default (e.g., 0.5 or 50%) and document this behavior.

---

### Edge Cases

- **Invalid Utilization**: Values for `utilization_percentage` outside the range of 0.0 to 1.0 (e.g., -0.1 or 1.5) MUST be clamped to the nearest valid boundary (0.0 or 1.0) by the plugin.
- **Metric Unavailability**: If a plugin advertises support for a metric in `SupportsResponse` but cannot calculate it for a specific `GetProjectedCostRequest`, it MUST omit that metric from the response rather than returning a default value.

## Assumptions

- **Default Utilization**: If `utilization_percentage` is not provided, plugins will use a default value of 0.5 (50%) for impact calculations, unless they have more specific data.
- **Metric Availability**: The `MetricKind` enumeration will be extended to specifically include Carbon Footprint, Energy Consumption, and Water Usage.
- **Range Enforcement**: Values for `utilization_percentage` outside [0.0, 1.0] will be clamped to the nearest valid boundary (0.0 or 1.0) by the plugin.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: `SupportsResponse` MUST include a `repeated MetricKind supported_metrics` field to advertise GreenOps capabilities.
- **FR-002**: `GetProjectedCostRequest` MUST include a global `double utilization_percentage` field to serve as a default for all resources in the request.
- **FR-003**: Individual resource specifications within `GetProjectedCostRequest` MUST include an optional `double utilization_percentage` field to override the global default.
- **FR-004**: Plugins MUST use the resource-level `utilization_percentage` if provided; otherwise, they MUST use the global default. If neither is provided, the default assumption (0.5) applies.
- **FR-005**: The engine MUST use the `supported_metrics` list to determine which UI columns or data fields to activate for a given resource.

### Key Entities

- **SupportsResponse**: The message returned by plugins to describe their capabilities for a specific resource type.
- **GetProjectedCostRequest**: The message sent by the engine to request cost and impact estimates.
- **MetricKind**: An enumeration of supported sustainability/impact metrics. Initial supported metrics are Carbon Footprint (gCO2e), Energy Consumption (kWh), and Water Usage (L).
- **Utilization Percentage**: A Protobuf `double` value (0.0 to 1.0). It can be specified globally at the request level (default) or per-resource (override) to scale impact calculations.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: 100% of plugins implementing this spec can successfully advertise their supported sustainability metrics via the `Supports` RPC.
- **SC-002**: The engine can dynamically generate a column for "Carbon Footprint" only when the plugin explicitly advertises it via `supported_metrics`.
- **SC-003**: Changes in the `utilization_percentage` parameter in `GetProjectedCostRequest` result in mathematically consistent changes in the impact metrics returned by the plugin (assuming a linear or documented non-linear model).
