# Data Model: Plugin Capability

## Entities

### PluginCapability (Enum)

Represents a discrete functional capability of a FinFocus plugin.

| Value | Description |
|-------|-------------|
| `PLUGIN_CAPABILITY_UNSPECIFIED` | Default value. |
| `PLUGIN_CAPABILITY_PROJECTED_COSTS` | Plugin implements `GetProjectedCost`. |
| `PLUGIN_CAPABILITY_ACTUAL_COSTS` | Plugin implements `GetActualCost`. |
| `PLUGIN_CAPABILITY_CARBON` | Plugin supports carbon/sustainability metrics. |
| `PLUGIN_CAPABILITY_RECOMMENDATIONS` | Plugin implements `GetRecommendations`. |
| `PLUGIN_CAPABILITY_DRY_RUN` | Plugin implements `DryRun`. |
| `PLUGIN_CAPABILITY_BUDGETS` | Plugin implements `GetBudgets`. |
| `PLUGIN_CAPABILITY_ENERGY` | Plugin supports energy consumption metrics. |
| `PLUGIN_CAPABILITY_WATER` | Plugin supports water usage metrics. |

## Relationships

- **GetPluginInfoResponse** `1` ---- `*` **PluginCapability**
- **Plugin** (Implementation) ---- `autodetects` ---- **PluginCapability**
