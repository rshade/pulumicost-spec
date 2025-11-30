# Data Model: Plugin Registry Index

**Date**: 2025-11-23
**Feature**: 004-plugin-registry-schema

## Entity Definitions

### RegistryIndex (Root)

Top-level container for the plugin registry index.

| Field          | Type   | Required | Validation                 | Description                       |
| -------------- | ------ | -------- | -------------------------- | --------------------------------- |
| schema_version | string | Yes      | Pattern: `^\d+\.\d+\.\d+$` | Version of registry schema format |
| plugins        | object | Yes      | Map of RegistryEntry       | Plugin name → metadata map        |

**Constraints**:

- `additionalProperties: false` - No extra fields allowed
- `plugins` values must conform to RegistryEntry schema

### RegistryEntry

Metadata for a single plugin in the registry.

| Field               | Type    | Required | Validation                                                      | Description                    |
| ------------------- | ------- | -------- | --------------------------------------------------------------- | ------------------------------ |
| name                | string  | Yes      | Pattern: `^[a-z0-9][a-z0-9-]*[a-z0-9]$\|^[a-z0-9]$`, 1-50 chars | Unique plugin identifier       |
| description         | string  | Yes      | 10-500 chars                                                    | Human-readable description     |
| repository          | string  | Yes      | Pattern: `^[a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+$`                     | GitHub owner/repo              |
| author              | string  | Yes      | 1-100 chars                                                     | Plugin author/organization     |
| supported_providers | array   | Yes      | Enum items, min 1, unique                                       | Cloud providers supported      |
| min_spec_version    | string  | Yes      | Pattern: `^\d+\.\d+\.\d+$`                                      | Minimum spec version required  |
| license             | string  | No       | -                                                               | SPDX license identifier        |
| homepage            | string  | No       | URI format                                                      | Plugin homepage URL            |
| capabilities        | array   | No       | Enum items, unique                                              | Plugin capabilities            |
| security_level      | string  | No       | Enum, default: "community"                                      | Trust level                    |
| max_spec_version    | string  | No       | Pattern: `^\d+\.\d+\.\d+$`                                      | Maximum spec version supported |
| keywords            | array   | No       | Max 10 items, each max 30 chars, unique                         | Searchable keywords            |
| deprecated          | boolean | No       | Default: false                                                  | Deprecation status             |
| deprecation_message | string  | No       | -                                                               | Migration guidance             |

**Constraints**:

- `additionalProperties: false` - No extra fields allowed
- `dependentRequired`: When `deprecated` is true, `deprecation_message` is required

## Enumerations

### supported_providers

Cloud platforms the plugin supports.

| Value        | Description           | Proto Reference                         |
| ------------ | --------------------- | --------------------------------------- |
| `aws`        | Amazon Web Services   | PluginSpecification.supported_providers |
| `azure`      | Microsoft Azure       | PluginSpecification.supported_providers |
| `gcp`        | Google Cloud Platform | PluginSpecification.supported_providers |
| `kubernetes` | Kubernetes clusters   | PluginSpecification.supported_providers |
| `custom`     | Custom/other sources  | PluginSpecification.supported_providers |

### capabilities

Plugin features aligned with gRPC service capabilities.

| Value               | Description                       | Proto Reference         |
| ------------------- | --------------------------------- | ----------------------- |
| `cost_retrieval`    | Retrieve actual cost data         | PluginInfo.capabilities |
| `cost_projection`   | Project future costs              | PluginInfo.capabilities |
| `pricing_specs`     | Provide pricing specifications    | PluginInfo.capabilities |
| `real_time_data`    | Real-time cost updates            | PluginInfo.capabilities |
| `historical_data`   | Historical cost queries           | PluginInfo.capabilities |
| `filtering`         | Filter cost data                  | PluginInfo.capabilities |
| `aggregation`       | Aggregate cost data               | PluginInfo.capabilities |
| `tagging`           | Tag-based cost allocation         | PluginInfo.capabilities |
| `recommendations`   | Cost optimization recommendations | PluginInfo.capabilities |
| `anomaly_detection` | Detect cost anomalies             | PluginInfo.capabilities |
| `forecasting`       | Cost forecasting                  | PluginInfo.capabilities |
| `budgets`           | Budget tracking                   | PluginInfo.capabilities |
| `alerts`            | Cost alerts                       | PluginInfo.capabilities |
| `custom`            | Custom capabilities               | PluginInfo.capabilities |

### security_level

Plugin security trust levels.

| Value       | Description                           | Proto Reference          |
| ----------- | ------------------------------------- | ------------------------ |
| `untrusted` | Untrusted, requires explicit approval | SECURITY_LEVEL_UNTRUSTED |
| `community` | Community verified (default)          | SECURITY_LEVEL_COMMUNITY |
| `verified`  | Officially verified                   | SECURITY_LEVEL_VERIFIED  |
| `official`  | Official PulumiCost plugin            | SECURITY_LEVEL_OFFICIAL  |

## Relationships

```text
RegistryIndex
    │
    └── plugins (map)
            │
            └── RegistryEntry
                    ├── supported_providers[] ──► Provider enum
                    ├── capabilities[] ──────────► Capability enum
                    └── security_level ─────────► SecurityLevel enum
```

## State Transitions

### Plugin Lifecycle in Registry

```text
[Not Listed] ──► [Active] ──► [Deprecated] ──► [Removed]
                    │              │
                    └──────────────┘
                    (can revert deprecation)
```

**State Indicators**:

- **Active**: `deprecated: false` or field absent
- **Deprecated**: `deprecated: true` with `deprecation_message`
- **Removed**: Entry deleted from registry (not represented in schema)

## Validation Rules (Consumer Responsibility)

These validations cannot be enforced by JSON Schema and must be validated by consuming
applications (e.g., pulumicost-core):

1. **Name/key match**: Plugin `name` must match its key in the `plugins` map
2. **Version ordering**: `max_spec_version` >= `min_spec_version` when both specified
3. **Repository existence**: Repository URL resolves to valid GitHub repo
4. **License validity**: License string is valid SPDX identifier
