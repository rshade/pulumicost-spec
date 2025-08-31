# PulumiCost Plugin Registry Examples

This document provides comprehensive examples for the PulumiCost Plugin Registry system, including plugin manifests, registry formats, and usage examples.

## Plugin Manifest Examples

### Kubecost Plugin

The [kubecost-plugin.json](plugin_manifests/kubecost-plugin.json) example demonstrates a comprehensive cost source plugin:

```json
{
  "name": "kubecost",
  "display_name": "Kubecost Cost Source Plugin",
  "description": "PulumiCost plugin for retrieving cost data from Kubecost installations",
  "version": "1.2.3",
  "api_version": "v1",
  "plugin_type": "cost_source",
  "capabilities": [
    "actual_cost", "projected_cost", "pricing_spec",
    "historical_cost", "kubernetes_native", "cost_allocation", "drill_down"
  ],
  "supported_providers": ["kubernetes", "aws", "gcp", "azure"]
}
```

Key features:
- Kubernetes-native cost allocation
- Multi-cloud provider support
- Comprehensive cost analysis capabilities
- Optional authentication with API keys

### AWS Cost Explorer Plugin

The [aws-cost-explorer-plugin.json](plugin_manifests/aws-cost-explorer-plugin.json) example shows cloud provider integration:

```json
{
  "name": "aws-cost-explorer",
  "display_name": "AWS Cost Explorer Plugin",
  "version": "2.1.0",
  "plugin_type": "cost_source",
  "supported_providers": ["aws"],
  "capabilities": [
    "actual_cost", "projected_cost", "pricing_spec",
    "historical_cost", "cost_forecasting", "budget_alerts"
  ]
}
```

Key features:
- AWS-specific cost retrieval
- Cost forecasting capabilities
- Budget alert integration
- IAM role authentication

## Plugin Registry Example

The [community-registry.json](registries/community-registry.json) example demonstrates a complete registry structure:

### Registry Metadata

```json
{
  "registry_version": "1.0.0",
  "registry_metadata": {
    "name": "PulumiCost Community Registry",
    "url": "https://registry.pulumicost.dev",
    "supported_api_versions": ["v1"]
  }
}
```

### Plugin Entries

Each plugin in the registry includes:
- Multiple version entries
- Download statistics
- Security status
- Verification badges
- User ratings

### Plugin Categories

```json
{
  "categories": [
    {
      "name": "Cost Sources",
      "description": "Plugins that provide cost data from various sources",
      "plugins": ["kubecost", "aws-cost-explorer", "azure-cost-management"]
    },
    {
      "name": "Cloud Providers", 
      "description": "Plugins specific to major cloud providers",
      "plugins": ["aws-cost-explorer", "azure-cost-management"]
    }
  ]
}
```

## Usage Examples

### Installation Commands

```bash
# Install from registry
pulumicost plugin install kubecost

# Install specific version
pulumicost plugin install kubecost@1.2.3

# Install from URL
pulumicost plugin install https://releases.example.com/plugins/kubecost/manifest.json

# Install with configuration
pulumicost plugin install kubecost --config '{"endpoint": "http://kubecost:9090"}'
```

### Discovery and Search

```bash
# List all plugins
pulumicost plugin list

# Search plugins
pulumicost plugin search "kubernetes cost"

# Filter by capabilities
pulumicost plugin list --capability cost_allocation

# Filter by provider
pulumicost plugin list --provider aws
```

### Management Commands

```bash
# Update all plugins
pulumicost plugin update

# Update specific plugin
pulumicost plugin update kubecost

# Remove plugin
pulumicost plugin remove kubecost

# Validate installation
pulumicost plugin validate kubecost

# Show plugin details
pulumicost plugin show kubecost
```

## Validation Examples

All examples are validated against JSON schemas. Use the validation tools:

```bash
# Validate plugin manifest
npm run validate:examples

# Validate with Go SDK
go run validate.go plugin_manifests/kubecost-plugin.json
```

Example validation result:
```json
{
  "valid": true,
  "errors": [],
  "warnings": [
    {
      "code": "MISSING_TAGS",
      "field": "metadata.tags",
      "message": "Tags are recommended for better categorization"
    }
  ]
}
```

## Security Examples

### Plugin Verification

```json
{
  "installation": {
    "binary_url": "https://releases.pulumicost.dev/plugins/kubecost/1.2.3/kubecost-plugin-linux-amd64",
    "checksum": "sha256:a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
    "size_bytes": 15728640
  }
}
```

### Authentication Configuration

```json
{
  "authentication": {
    "required": true,
    "methods": ["api_key", "iam_role"],
    "scopes": ["ce:GetCostAndUsage", "ce:GetUsageReport"],
    "permissions": [
      "ce:GetCostAndUsage",
      "ce:GetReservationUtilization"
    ]
  }
}
```

## Integration Examples

### Configuration Schema

```json
{
  "configuration": {
    "schema_url": "https://spec.pulumicost.dev/plugins/kubecost/config-schema.json",
    "defaults": {
      "timeout": "30s",
      "currency": "USD",
      "window": "1d"
    },
    "examples": [
      {
        "endpoint": "http://kubecost-cost-analyzer:9090",
        "currency": "USD",
        "window": "1d",
        "aggregate_by": ["namespace", "deployment"]
      }
    ]
  }
}
```

### Environment Variables

```json
{
  "environment_variables": [
    {
      "name": "KUBECOST_ENDPOINT",
      "description": "Kubecost API endpoint URL",
      "required": true
    },
    {
      "name": "KUBECOST_TOKEN", 
      "description": "Authentication token for Kubecost API",
      "required": false
    },
    {
      "name": "KUBECOST_TIMEOUT",
      "description": "Request timeout in seconds",
      "required": false,
      "default_value": "30"
    }
  ]
}
```

## Best Practices

### Plugin Manifest Best Practices

1. **Complete Metadata**
   - Always provide display name and description
   - Include comprehensive tags for discoverability
   - Specify clear maintainer contact information

2. **Version Management**
   - Use semantic versioning strictly
   - Clearly specify API version compatibility
   - Provide migration guidance for breaking changes

3. **Security**
   - Specify minimum required permissions
   - Use checksums for all binary distributions
   - Implement proper error handling

4. **Dependencies**
   - Clearly specify all external dependencies
   - Include version constraints
   - Mark optional dependencies appropriately

### Registry Best Practices

1. **Quality Assurance**
   - Implement automated validation
   - Perform security scanning
   - Test plugin compatibility

2. **User Experience**
   - Maintain accurate download statistics
   - Provide clear categorization
   - Enable effective search functionality

3. **Maintenance**
   - Regular security audits
   - Timely vulnerability reporting
   - Clear deprecation processes

## Contributing

To add new plugin examples:

1. Create the plugin manifest following the schema
2. Validate the manifest using the provided tools
3. Add comprehensive documentation
4. Test with the reference implementation
5. Submit a pull request

For detailed specification information, see [PLUGIN_REGISTRY_SPECIFICATION.md](../docs/PLUGIN_REGISTRY_SPECIFICATION.md).