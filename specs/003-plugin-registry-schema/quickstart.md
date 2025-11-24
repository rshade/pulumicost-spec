# Quickstart: Plugin Registry Index Schema

**Date**: 2025-11-23
**Feature**: 003-plugin-registry-schema

## Overview

This guide shows how to create a valid plugin registry entry and validate it against the
schema.

## Creating a Registry Entry

### Minimal Example

```json
{
  "schema_version": "1.0.0",
  "plugins": {
    "my-plugin": {
      "name": "my-plugin",
      "description": "A cost data plugin for my custom cost source",
      "repository": "myorg/pulumicost-plugin-custom",
      "author": "My Organization",
      "supported_providers": ["custom"],
      "min_spec_version": "0.1.0"
    }
  }
}
```

### Full Example

```json
{
  "schema_version": "1.0.0",
  "plugins": {
    "kubecost": {
      "name": "kubecost",
      "description": "Kubernetes cost analysis via Kubecost API integration",
      "repository": "rshade/pulumicost-plugin-kubecost",
      "author": "PulumiCost Team",
      "license": "Apache-2.0",
      "homepage": "https://github.com/rshade/pulumicost-plugin-kubecost",
      "supported_providers": ["kubernetes"],
      "capabilities": ["cost_retrieval", "cost_projection", "real_time_data"],
      "security_level": "official",
      "min_spec_version": "0.1.0",
      "keywords": ["kubernetes", "kubecost", "k8s", "container"]
    }
  }
}
```

## Field Requirements

### Required Fields

| Field | Format | Example |
|-------|--------|---------|
| name | Lowercase alphanumeric with hyphens | `my-plugin` |
| description | 10-500 characters | `A cost data plugin...` |
| repository | owner/repo format | `myorg/my-plugin` |
| author | 1-100 characters | `My Organization` |
| supported_providers | Array of valid providers | `["aws", "azure"]` |
| min_spec_version | Semantic version | `0.1.0` |

### Optional Fields

| Field | Format | Default | Example |
|-------|--------|---------|---------|
| license | SPDX identifier | - | `Apache-2.0` |
| homepage | Valid URI | - | `https://example.com` |
| capabilities | Array of valid capabilities | - | `["cost_retrieval"]` |
| security_level | Valid level | `community` | `verified` |
| max_spec_version | Semantic version | - | `1.0.0` |
| keywords | Array (max 10, each max 30 chars) | - | `["aws", "ec2"]` |
| deprecated | Boolean | `false` | `true` |
| deprecation_message | String (required if deprecated) | - | `Use aws-v2 instead` |

## Validation

### Using npm Scripts

```bash
# Validate registry against schema
npm run validate:registry

# Validate all examples (including registry)
npm run validate
```

### Using AJV Directly

```bash
npx ajv validate -s schemas/plugin_registry.schema.json \
  -d examples/registry.json --strict=false
```

### Common Validation Errors

**Missing required field**:

```text
Error: must have required property 'min_spec_version'
```

**Invalid name pattern**:

```text
Error: must match pattern "^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$"
```

**Description too short**:

```text
Error: must NOT have fewer than 10 characters
```

**Invalid provider**:

```text
Error: must be equal to one of the allowed values
Allowed: aws, azure, gcp, kubernetes, custom
```

**Deprecated without message**:

```text
Error: must have required property 'deprecation_message'
```

## Adding Your Plugin to the Registry

1. **Create your entry**: Follow the format above with all required fields
2. **Validate locally**: Run `npm run validate:registry`
3. **Submit PR**: Add your entry to `examples/registry.json`
4. **Ensure name matches key**: The plugin `name` field must match its key in `plugins`

## Consumer Validation

The schema validates structure and patterns. Consuming applications (pulumicost-core)
additionally validate:

- Plugin name matches registry key
- Version ordering (`max_spec_version >= min_spec_version`)
- Repository URL is accessible

## Version Compatibility

- **schema_version**: Format version of the registry file
- **min_spec_version**: Minimum PulumiCost spec version required by plugin
- **max_spec_version**: Maximum PulumiCost spec version supported (optional)

When installing a plugin, pulumicost-core checks that the current spec version falls
within the plugin's supported range.
