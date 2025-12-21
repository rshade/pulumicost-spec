# Data Model: Add ARN to GetActualCostRequest

## Entities

### GetActualCostRequest

A message sent to a plugin to request actual cost data.

| Field | Type | Number | Description | New/Existing |
| :--- | :--- | :--- | :--- | :--- |
| `resource_id` | `string` | 1 | Pulumi-internal identifier (URN, logical name). | Existing |
| `start` | `Timestamp` | 2 | Start time of the window. | Existing |
| `end` | `Timestamp` | 3 | End time of the window. | Existing |
| `tags` | `map<string, string>` | 4 | Resource tags. | Existing |
| `arn` | `string` | 5 | Canonical Cloud Identifier (e.g., AWS ARN). | **New** |

## Relationships

* Used in `CostSourceService.GetActualCost` RPC.
