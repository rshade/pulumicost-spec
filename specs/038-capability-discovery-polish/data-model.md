# Data Model: Capability Discovery

## Core Entities

### PluginInfo (SDK)

The central struct holding plugin configuration.

| Field | Type | Description |
| :--- | :--- | :--- |
| `Name` | `string` | Plugin name |
| `Version` | `string` | Plugin version |
| `Capabilities` | `[]PluginCapability` | List of supported capabilities (Enums) |
| `Metadata` | `map[string]string` | **Legacy** string map for backward compatibility |

### Capability Mapping (New)

A static mapping defining the relationship between Enums and Legacy Strings.

**Location**: `sdk/go/pluginsdk/capability_compat.go`

| Enum | Legacy String Key | Value |
| :--- | :--- | :--- |
| `PLUGIN_CAPABILITY_DRY_RUN` | `supports_dry_run` | `"true"` |
| `PLUGIN_CAPABILITY_RECOMMENDATIONS` | `supports_recommendations` | `"true"` |
| `PLUGIN_CAPABILITY_BUDGETS` | `supports_budgets` | `"true"` |
| `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS` | `supports_dismiss_recommendations` | `"true"` |
| `PLUGIN_CAPABILITY_PROJECTED_COSTS` | `supports_projected_costs` | `"true"` |
| `PLUGIN_CAPABILITY_ACTUAL_COSTS` | `supports_actual_costs` | `"true"` |
| `PLUGIN_CAPABILITY_PRICING_SPEC` | `supports_pricing_spec` | `"true"` |
| `PLUGIN_CAPABILITY_ESTIMATE_COST` | `supports_estimate_cost` | `"true"` |

## Interfaces

### Auto-Discovery Interfaces

| Interface Name | Method Signature | Maps To Capability |
| :--- | :--- | :--- |
| `DryRunHandler` | `HandleDryRun(*pbc.DryRunRequest) (*pbc.DryRunResponse, error)` | `PLUGIN_CAPABILITY_DRY_RUN` |
| `RecommendationsProvider` | `GetRecommendations(...)` | `PLUGIN_CAPABILITY_RECOMMENDATIONS` |
| `BudgetsProvider` | `GetBudgets(...)` | `PLUGIN_CAPABILITY_BUDGETS` |
| `DismissProvider` | `DismissRecommendation(...)` | `PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS` |
