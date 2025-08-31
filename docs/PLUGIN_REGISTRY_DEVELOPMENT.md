# Plugin Registry Development Guide

This document provides guidance for developing and working with the PulumiCost Plugin Registry system.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Node.js 22 or later  
- npm 10 or later
- Protocol Buffer compiler (buf)

### Initial Setup

```bash
# Install dependencies
npm install

# Install Go dependencies and generate code
make generate
make tidy

# Validate everything is working
make validate
```

## Architecture Overview

The plugin registry system consists of several key components:

### 1. Schema Definitions

- **[plugin_manifest.schema.json](../schemas/plugin_manifest.schema.json)** - JSON schema for plugin manifests
- **[plugin_registry.schema.json](../schemas/plugin_registry.schema.json)** - JSON schema for plugin registries

### 2. Protocol Buffer Definitions

- **[plugin_registry.proto](../proto/pulumicost/v1/plugin_registry.proto)** - gRPC service definitions for plugin registry operations

### 3. Go SDK Components

- **[types.go](../sdk/go/registry/types.go)** - Core data structures and types
- **[interfaces.go](../sdk/go/registry/interfaces.go)** - Service interfaces
- **[validation.go](../sdk/go/registry/validation.go)** - Plugin validation logic
- **[validation_test.go](../sdk/go/registry/validation_test.go)** - Validation tests

### 4. Examples and Documentation

- **[Plugin manifests](../examples/plugin_manifests/)** - Example plugin definitions
- **[Registry examples](../examples/registries/)** - Example registry structures  
- **[PLUGIN_REGISTRY_SPECIFICATION.md](PLUGIN_REGISTRY_SPECIFICATION.md)** - Complete specification document

## Development Workflow

### Adding New Plugin Types

1. Update the `PluginType` enum in [types.go](../sdk/go/registry/types.go)
2. Add validation in [validation.go](../sdk/go/registry/validation.go)
3. Update the JSON schema in [plugin_manifest.schema.json](../schemas/plugin_manifest.schema.json)
4. Update the proto definition in [plugin_registry.proto](../proto/pulumicost/v1/plugin_registry.proto)
5. Add tests and examples

### Adding New Capabilities

1. Update the `Capability` enum in [types.go](../sdk/go/registry/types.go)
2. Add validation logic if needed
3. Update JSON schemas
4. Update proto definitions
5. Add examples demonstrating the new capability

### Updating Schemas

When updating JSON schemas:

1. Maintain backward compatibility where possible
2. Use semantic versioning for breaking changes
3. Update all related examples
4. Regenerate Go code if needed
5. Run validation tests

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test -v ./sdk/go/registry/

# Run with coverage
go test -cover ./sdk/go/registry/

# Run benchmarks
go test -bench=. -benchmem ./sdk/go/registry/
```

### Validation Testing

```bash
# Validate JSON schemas
make validate-schema

# Validate examples against schemas
make validate-examples

# Run full validation suite
make validate
```

### Integration Testing

```bash
# Test with real plugin manifests
go test -v -run TestValidateManifest ./sdk/go/registry/

# Test binary validation
go test -v -run TestValidateBinary ./sdk/go/registry/

# Test security scanning
go test -v -run TestScanSecurity ./sdk/go/registry/
```

## Implementation Guidelines

### Plugin Validation

The validation system implements three levels of checking:

1. **Schema Validation** - JSON schema compliance
2. **Business Logic Validation** - Domain-specific rules
3. **Security Validation** - Security scanning and checks

#### Adding Custom Validators

```go
type CustomValidator struct {
    *DefaultValidator
}

func (v *CustomValidator) ValidateManifest(manifest *PluginManifest) (*ValidationResult, error) {
    // Call base validation
    result, err := v.DefaultValidator.ValidateManifest(manifest)
    if err != nil {
        return nil, err
    }
    
    // Add custom validation logic
    if manifest.Name == "reserved-name" {
        result.addError("RESERVED_NAME", "name", "Name is reserved")
    }
    
    return result, nil
}
```

### Registry Implementation

#### Basic Registry Client

```go
client := &RegistryClient{
    baseURL: "https://registry.pulumicost.dev",
    timeout: 30 * time.Second,
}

plugins, total, err := client.ListPlugins(ctx, &PluginFilter{
    PluginType: PluginTypeCostSource,
    Providers:  []string{"aws"},
}, 50, 0)
```

#### Custom Discovery

```go
type FilesystemDiscovery struct {
    paths []string
}

func (d *FilesystemDiscovery) Discover(ctx context.Context) ([]*PluginManifest, error) {
    var manifests []*PluginManifest
    
    for _, path := range d.paths {
        // Scan filesystem for manifests
        // Parse and validate
        // Add to results
    }
    
    return manifests, nil
}
```

## Security Considerations

### Plugin Binary Validation

- Always validate checksums
- Implement signature verification
- Perform security scanning
- Check executable headers

### Authentication and Authorization

- Support multiple authentication methods
- Implement proper RBAC
- Audit access logs
- Rate limit API calls

### Registry Security

- Validate all uploads
- Scan for vulnerabilities
- Implement content policies
- Monitor for malicious activity

## Performance Optimization

### Validation Performance

- Cache validation results
- Implement parallel validation
- Use efficient JSON parsing
- Optimize regex patterns

### Registry Operations

- Implement proper pagination
- Use database indexing
- Cache frequently accessed data
- Implement search optimization

## Monitoring and Observability

### Metrics to Track

- Plugin installation rates
- Validation success/failure rates
- Security scan results
- API response times
- Registry uptime

### Logging

- Log all validation failures
- Track security events
- Monitor performance metrics
- Implement structured logging

## Troubleshooting

### Common Issues

1. **Schema Validation Failures**
   - Check JSON syntax
   - Verify required fields
   - Validate field formats

2. **Binary Validation Failures**
   - Verify checksums
   - Check file permissions
   - Validate executable format

3. **Compatibility Issues**
   - Check API version requirements
   - Verify OS/architecture support
   - Validate dependencies

### Debug Commands

```bash
# Validate specific manifest
go run validate.go examples/plugin_manifests/kubecost-plugin.json

# Check schema syntax
npm run validate:schema

# Run linting
make lint

# Generate debug output
go test -v -run TestValidation ./sdk/go/registry/ > debug.log 2>&1
```

## Contributing

### Development Process

1. Create feature branch
2. Implement changes
3. Add comprehensive tests
4. Update documentation
5. Run full validation suite
6. Submit pull request

### Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comprehensive comments
- Include error handling
- Write testable code

### Documentation Requirements

- Update specification documents
- Add usage examples
- Update API documentation
- Include migration guides
- Provide troubleshooting info

## Future Enhancements

### Planned Features

1. **Enhanced Security Scanning**
   - Static analysis integration
   - Vulnerability database updates
   - Automated security reporting

2. **Registry Federation**
   - Multi-registry support
   - Registry synchronization
   - Conflict resolution

3. **Advanced Discovery**
   - OCI registry support
   - Git-based discovery
   - Automatic updates

4. **Performance Improvements**
   - Caching layer
   - Async validation
   - Streaming APIs

### Extension Points

The system is designed for extensibility:

- Custom validators
- Plugin discovery sources
- Authentication methods
- Registry backends
- Security scanners

For more information, see the [Plugin Registry Specification](PLUGIN_REGISTRY_SPECIFICATION.md).