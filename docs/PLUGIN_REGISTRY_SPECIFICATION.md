# PulumiCost Plugin Registry Specification v1.0

## Overview

The PulumiCost Plugin Registry Specification defines a comprehensive framework for plugin discovery, registration, and management within the PulumiCost ecosystem. This specification enables dynamic plugin discovery from multiple sources and provides centralized plugin lifecycle management.

## Table of Contents

- [Plugin Discovery Mechanisms](#plugin-discovery-mechanisms)
- [Plugin Manifest Format](#plugin-manifest-format)
- [Plugin Registry Format](#plugin-registry-format)
- [Versioning and Compatibility](#versioning-and-compatibility)
- [Security Model](#security-model)
- [Lifecycle Management](#lifecycle-management)
- [Reference Implementation](#reference-implementation)

## Plugin Discovery Mechanisms

### 1. Filesystem Discovery

Plugins can be discovered from local filesystem locations using directory scanning:

```
/usr/local/lib/pulumicost/plugins/
├── kubecost/
│   ├── 1.2.3/
│   │   ├── manifest.json
│   │   └── kubecost-plugin
│   └── 1.2.2/
│       ├── manifest.json
│       └── kubecost-plugin
└── aws-cost-explorer/
    └── 2.1.0/
        ├── manifest.json
        └── aws-cost-explorer-plugin
```

**Discovery Process:**
1. Scan configured plugin directories
2. Look for `manifest.json` files in subdirectories
3. Parse and validate manifest schemas
4. Register plugins with the local registry cache

**Configuration:**
```yaml
discovery:
  filesystem:
    enabled: true
    paths:
      - "/usr/local/lib/pulumicost/plugins"
      - "~/.pulumicost/plugins"
      - "./plugins"
    recursive: true
    watch_for_changes: true
```

### 2. Registry Discovery

Plugins can be discovered from remote registries using HTTP/HTTPS APIs:

**Registry Endpoint Structure:**
- `GET /plugins` - List all available plugins
- `GET /plugins/{name}` - Get specific plugin details
- `GET /plugins/{name}/{version}/manifest` - Get plugin manifest
- `GET /search?q={query}&provider={provider}` - Search plugins

**Example Registry Discovery:**
```yaml
discovery:
  registries:
    - name: "community"
      url: "https://registry.pulumicost.dev"
      priority: 1
      timeout: "30s"
      auth:
        type: "none"
    - name: "enterprise"
      url: "https://enterprise-registry.internal"
      priority: 2
      timeout: "60s"
      auth:
        type: "api_key"
        api_key: "${REGISTRY_API_KEY}"
```

### 3. URL-Based Discovery

Individual plugins can be installed directly from URLs:

```bash
# Install from direct manifest URL
pulumicost plugin install https://releases.example.com/plugins/my-plugin/1.0.0/manifest.json

# Install from GitHub releases
pulumicost plugin install github://pulumicost/kubecost-plugin@v1.2.3

# Install from OCI registry
pulumicost plugin install oci://registry.example.com/pulumicost/plugins/kubecost:1.2.3
```

**Supported URL Schemes:**
- `https://` - Direct HTTP/HTTPS URLs
- `github://` - GitHub repository releases
- `oci://` - OCI-compliant container registries
- `s3://` - Amazon S3 buckets
- `gs://` - Google Cloud Storage
- `file://` - Local file paths

## Plugin Manifest Format

The plugin manifest uses JSON Schema validation and follows semantic versioning. See [`schemas/plugin_manifest.schema.json`](../schemas/plugin_manifest.schema.json) for the complete specification.

### Core Manifest Structure

```json
{
  "name": "plugin-name",
  "display_name": "Human Readable Name", 
  "description": "Brief description of plugin functionality",
  "version": "1.2.3",
  "api_version": "v1",
  "plugin_type": "cost_source",
  "capabilities": ["actual_cost", "projected_cost"],
  "supported_providers": ["aws", "azure"],
  "requirements": { ... },
  "authentication": { ... },
  "installation": { ... },
  "configuration": { ... },
  "contacts": { ... },
  "metadata": { ... }
}
```

### Plugin Types

- **`cost_source`** - Provides cost data from various sources
- **`observability`** - Provides monitoring and metrics capabilities  
- **`registry`** - Provides plugin registry functionality
- **`aggregator`** - Aggregates and processes cost data

### Capabilities

Standard capabilities include:
- `actual_cost` - Historical cost data retrieval
- `projected_cost` - Future cost projections
- `pricing_spec` - Pricing specification data
- `real_time_cost` - Real-time cost monitoring
- `cost_forecasting` - Cost forecasting algorithms
- `cost_optimization` - Cost optimization recommendations
- `budget_alerts` - Budget threshold alerting
- `custom_metrics` - Custom metrics collection
- `multi_cloud` - Multi-cloud support
- `kubernetes_native` - Kubernetes-native functionality
- `tagging_support` - Resource tagging support
- `drill_down` - Detailed cost breakdowns
- `cost_allocation` - Cost allocation algorithms
- `showback` - Cost showback reporting
- `chargeback` - Cost chargeback functionality

## Plugin Registry Format

The plugin registry format aggregates multiple plugin entries with metadata. See [`schemas/plugin_registry.schema.json`](../schemas/plugin_registry.schema.json) for the complete specification.

### Registry Structure

```json
{
  "registry_version": "1.0.0",
  "registry_metadata": {
    "name": "Community Registry",
    "url": "https://registry.pulumicost.dev",
    "supported_api_versions": ["v1"]
  },
  "plugins": [
    {
      "name": "plugin-name",
      "latest_version": "1.2.3",
      "versions": [
        {
          "version": "1.2.3",
          "manifest_url": "https://registry.example.com/plugins/plugin-name/1.2.3/manifest.json",
          "published": "2024-08-31T12:00:00Z",
          "deprecated": false,
          "yanked": false
        }
      ]
    }
  ]
}
```

## Versioning and Compatibility

### Semantic Versioning

All plugins and registries MUST follow [Semantic Versioning 2.0.0](https://semver.org/):

- **Major version** (`X.y.z`) - Breaking changes to plugin API or manifest format
- **Minor version** (`x.Y.z`) - Backward-compatible feature additions
- **Patch version** (`x.y.Z`) - Backward-compatible bug fixes

### API Version Compatibility

Plugins specify API version compatibility using `min_api_version` and `max_api_version`:

```json
{
  "api_version": "v1",
  "requirements": {
    "min_api_version": "v1",
    "max_api_version": "v1"
  }
}
```

### Version Resolution Algorithm

1. **Exact Version** - `plugin@1.2.3` installs exactly version 1.2.3
2. **Version Range** - `plugin@^1.2.0` installs highest compatible 1.x.x version
3. **Latest** - `plugin@latest` installs the highest non-prerelease version
4. **Channel** - `plugin@beta` installs the latest beta version

### Breaking Change Policy

- **Major API changes** require major version bump
- **Manifest schema changes** require registry version bump
- **Deprecation warnings** must be provided for 1 major version before removal
- **Security vulnerabilities** may warrant immediate version yanking

## Security Model

### Plugin Verification

Plugins can be verified through multiple mechanisms:

1. **Cryptographic Signatures** - RSA/ECDSA signatures using registry keys
2. **Checksum Validation** - SHA-256 checksums for binary integrity
3. **Registry Verification** - Official verification badges from trusted registries
4. **Security Scanning** - Automated vulnerability scanning of plugin binaries

### Authentication Methods

Supported authentication methods for plugin APIs:

- `api_key` - API key authentication
- `oauth2` - OAuth 2.0 flows
- `service_account` - Service account credentials
- `iam_role` - Cloud provider IAM roles
- `mutual_tls` - Mutual TLS authentication
- `basic_auth` - HTTP Basic authentication
- `bearer_token` - Bearer token authentication

### Security Scanning

Plugin registries SHOULD implement security scanning:

```json
{
  "security_scan_results": {
    "status": "secure",
    "vulnerabilities": [],
    "last_scan": "2024-08-31T12:00:00Z",
    "scan_version": "v1.5.0"
  }
}
```

### Sandboxing and Isolation

Plugin execution environments SHOULD implement:

- **Process isolation** - Separate processes for plugin execution
- **Network restrictions** - Limited network access based on declared requirements
- **Filesystem restrictions** - Limited filesystem access to necessary directories
- **Resource limits** - CPU and memory usage limits
- **Capability-based security** - Plugins can only access declared capabilities

## Lifecycle Management

### Installation

```bash
# Install from registry
pulumicost plugin install kubecost

# Install specific version
pulumicost plugin install kubecost@1.2.3

# Install from URL
pulumicost plugin install https://example.com/plugin.json

# Install with configuration
pulumicost plugin install kubecost --config config.json
```

### Updates

```bash
# Update all plugins
pulumicost plugin update

# Update specific plugin
pulumicost plugin update kubecost

# Update to specific version
pulumicost plugin update kubecost@1.3.0

# Check for updates
pulumicost plugin list --outdated
```

### Removal

```bash
# Remove plugin
pulumicost plugin remove kubecost

# Remove specific version
pulumicost plugin remove kubecost@1.2.3

# Remove and clean cache
pulumicost plugin remove kubecost --clean
```

### Status and Management

```bash
# List installed plugins
pulumicost plugin list

# Show plugin details
pulumicost plugin show kubecost

# Validate plugin installation
pulumicost plugin validate kubecost

# Enable/disable plugins
pulumicost plugin enable kubecost
pulumicost plugin disable kubecost
```

## Reference Implementation

### Go SDK Components

The Go SDK provides reference implementations in the following packages:

- **`sdk/go/registry`** - Plugin registry client and server implementations
- **`sdk/go/discovery`** - Plugin discovery mechanisms
- **`sdk/go/installer`** - Plugin installation and lifecycle management
- **`sdk/go/validation`** - Manifest and binary validation

### Core Interfaces

```go
// PluginRegistry defines the interface for plugin registries
type PluginRegistry interface {
    ListPlugins(ctx context.Context, filter *PluginFilter) ([]*PluginInfo, error)
    GetPlugin(ctx context.Context, name, version string) (*PluginManifest, error)
    SearchPlugins(ctx context.Context, query string) ([]*PluginInfo, error)
    RegisterPlugin(ctx context.Context, manifest *PluginManifest) error
}

// PluginDiscovery defines the interface for plugin discovery
type PluginDiscovery interface {
    Discover(ctx context.Context) ([]*PluginManifest, error)
    Watch(ctx context.Context, callback func(*PluginManifest)) error
}

// PluginInstaller defines the interface for plugin installation
type PluginInstaller interface {
    Install(ctx context.Context, source string, opts *InstallOptions) error
    Update(ctx context.Context, name string, opts *UpdateOptions) error
    Remove(ctx context.Context, name string) error
    Validate(ctx context.Context, name string) (*ValidationResult, error)
}
```

### Example Usage

```go
// Create registry client
registry := NewRegistryClient("https://registry.pulumicost.dev")

// Discover plugins
plugins, err := registry.ListPlugins(ctx, &PluginFilter{
    PluginType: PluginTypeCostSource,
    Providers:  []string{"aws", "kubernetes"},
})

// Install plugin
installer := NewPluginInstaller()
err = installer.Install(ctx, "kubecost@1.2.3", &InstallOptions{
    RegistryURL: "https://registry.pulumicost.dev",
    VerifySignature: true,
})
```

## Compliance and Standards

### Schema Validation

All JSON documents MUST validate against their respective JSON schemas:
- Plugin manifests: `schemas/plugin_manifest.schema.json`
- Plugin registries: `schemas/plugin_registry.schema.json`

### Protocol Buffer Definitions

gRPC service definitions are provided in `proto/pulumicost/v1/plugin_registry.proto` for programmatic access to plugin registries.

### Testing Requirements

Plugin implementations SHOULD provide:
- Unit tests for core functionality
- Integration tests with the plugin framework
- Conformance tests using the standard test suite
- Security vulnerability scanning
- Performance benchmarking

### Documentation Standards

Plugin documentation MUST include:
- Installation and configuration instructions
- API reference documentation
- Usage examples and tutorials
- Troubleshooting guides
- Security considerations
- Performance characteristics

## Migration and Evolution

### Backward Compatibility

- Registry format changes require major version increments
- Plugin manifest changes must maintain backward compatibility within major versions
- Deprecated fields must be supported for at least one major version

### Forward Compatibility

- New optional fields can be added in minor versions
- Unknown fields in manifests should be preserved but ignored
- Clients should gracefully handle unknown plugin types and capabilities

---

This specification provides a comprehensive framework for plugin management in the PulumiCost ecosystem while maintaining flexibility for future enhancements and extensions.