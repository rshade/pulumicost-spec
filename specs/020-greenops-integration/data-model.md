# Data Model: GreenOps Integration

## New Entities

### MetricKind (Enum)

Represents the type of sustainability metric supported by a plugin.

| Value | Name | Description | Unit |
|-------|------|-------------|------|
| 0 | UNSPECIFIED | Default unspecified value | N/A |
| 1 | CARBON_FOOTPRINT | GreenHouse Gas emissions | gCO2e |
| 2 | ENERGY_CONSUMPTION | Electrical energy consumed | kWh |
| 3 | WATER_USAGE | Water consumed | L |

## Updated Entities

### SupportsResponse

Added field to advertise GreenOps capabilities.

| Field | Type | Description |
|-------|------|-------------|
| supported_metrics | repeated MetricKind | List of sustainability metrics supported for the resource |

### GetProjectedCostRequest

Added global utilization assumption.

| Field | Type | Description |
|-------|------|-------------|
| utilization_percentage | double | Global default utilization (0.0 to 1.0) for all resources in request |

### ResourceDescriptor

Added per-resource utilization override.

| Field | Type | Description |
|-------|------|-------------|
| utilization_percentage | optional double | Per-resource utilization override (0.0 to 1.0) |

## Validation Rules

- **Utilization Range**: `0.0 <= utilization_percentage <= 1.0`.
- **Clamping**: Plugins MUST clamp values outside this range to [0.0, 1.0].
- **Precedence**: Resource-level `utilization_percentage` > Global-level `utilization_percentage` > Default (0.5).
