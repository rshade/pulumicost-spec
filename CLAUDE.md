# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **pulumicost-spec**, a repository that provides the canonical protocol and schemas for PulumiCost plugins. It defines:
- gRPC service definitions for cost source plugins
- JSON schemas for pricing specifications
- Go SDK with generated protobuf code and helper types

## Build Commands

- `make generate` - Generate Go code from protobuf definitions (installs buf locally in bin/)
- `make tidy` - Run `go mod tidy` to clean up dependencies  
- `make test` - Run all Go tests including integration tests
- `make lint` - Run golangci-lint and buf lint
- `make validate` - Run tests and linting together
- `make clean` - Remove generated proto files
- `make clean-all` - Remove generated files and local tools (bin/)
- `go build ./...` - Build all Go packages
- `go test -bench=. -benchmem ./sdk/go/testing/` - Run performance benchmarks

## Architecture

### Core Components

**Proto Definition (`proto/pulumicost/costsource.proto`)**
- Defines `CostSource` gRPC service with RPCs for: Name, Supports, GetActualCost, GetProjectedCost, GetPricingSpec
- Contains message definitions for requests/responses
- Uses Google protobuf types (Empty, Timestamp)

**JSON Schema (`schemas/pricing_spec.schema.json`)**
- Validates PricingSpec documents
- Defines required fields: provider, resource_type, billing_mode, rate_per_unit, currency
- Enforces billing_mode enum values and data types

**Go SDK (`sdk/go/`)**
- `sdk/go/proto/` - Generated protobuf Go code (do not edit manually)
- `sdk/go/types/domain.go` - BillingMode enum constants and validation helpers
- `sdk/go/types/validate.go` - JSON schema validation for PricingSpec documents
- `sdk/go/testing/` - Comprehensive plugin testing framework

**Testing Framework (`sdk/go/testing/`)**
- `harness.go` - In-memory gRPC test harness with bufconn
- `mock_plugin.go` - Configurable mock plugin implementation
- `integration_test.go` - Comprehensive integration tests for all RPC methods
- `benchmark_test.go` - Performance benchmarks with memory profiling
- `conformance_test.go` - Multi-level plugin conformance testing (Basic/Standard/Advanced)
- `README.md` - Complete testing guide for plugin developers

**Examples (`examples/`)**
- `examples/specs/` - 8 comprehensive cross-vendor JSON examples
- `examples/README.md` - Documentation of all billing models and examples

### Generated Code

The `sdk/go/proto/` directory contains generated Go protobuf code. To regenerate:
1. Run `make generate` (automatically installs buf v1.32.1 locally in bin/)
2. Generated code is automatically validated in CI

### Code Generation Dependencies

- **buf** - Protocol buffer toolchain (installed locally in bin/ via make generate)
- **google.golang.org/protobuf** - Go protobuf runtime
- **google.golang.org/grpc** - gRPC Go implementation
- **golangci-lint** - Go linting (installed via make depend)

### Local Tool Management

The project uses local tool installation to avoid version conflicts:
- `bin/buf` - buf CLI v1.32.1 installed automatically
- Tools are excluded from git via `.gitignore`
- `make clean-all` removes all local tools

## Development Workflow

1. **Modify Proto**: Edit `proto/pulumicost/costsource.proto`
2. **Update Schema**: Edit `schemas/pricing_spec.schema.json` if PricingSpec message changes
3. **Regenerate**: Run `make generate` to update Go bindings
4. **Update Types**: Modify helper code in `sdk/go/types/` as needed
5. **Test**: Run `make validate` to run tests and linting
6. **Verify**: Run integration tests with `go test -v ./sdk/go/testing/`

## Testing Workflow

### Integration Testing
- Use `TestHarness` for in-memory gRPC testing with bufconn
- Create mock plugins with `NewMockPlugin()` for configurable behavior
- Run comprehensive tests: `go test -v ./sdk/go/testing/`

### Performance Testing
- Run benchmarks: `go test -bench=. -benchmem ./sdk/go/testing/`
- Measure all RPC methods with memory profiling
- Test different data sizes and concurrent requests

### Conformance Testing
- **Basic**: Required for all plugins - core functionality
- **Standard**: Recommended for production - reliability and consistency  
- **Advanced**: High-performance requirements - scalability and performance

### Example Usage
```go
// Create test harness
plugin := &MyPluginImpl{}
harness := plugintesting.NewTestHarness(plugin)
harness.Start(t)
defer harness.Stop()

// Run conformance tests
result := plugintesting.RunStandardConformanceTests(t, plugin)
if result.FailedTests > 0 {
    t.Errorf("Plugin failed conformance: %s", result.Summary)
}
```

## Package Structure

```
github.com/rshade/pulumicost-spec/sdk/go/proto  # Generated protobuf code
github.com/rshade/pulumicost-spec/sdk/go/types  # Helper types and validation
```

## Schema Validation

The types package embeds the JSON schema and provides `ValidatePricingSpec(doc []byte) error` for validating PricingSpec JSON documents against the schema.

## Versioning

Follow semantic versioning for proto changes:
- MAJOR: Breaking proto changes
- MINOR: Backward-compatible additions  
- PATCH: Bug fixes, documentation

Tag releases as `v0.1.0`, `v1.0.0`, etc.
- This is the project for this repo: https://github.com/users/rshade/projects/3

## Common Issues & Solutions

### Mock Plugin Implementation
- Issue: gRPC method name conflicts in mock plugins
- Solution: Use `PluginName` field instead of `Name` to avoid conflicts with RPC method names
- Pattern: Separate data fields from method names in struct design

### Integration Testing Setup
- Issue: Network-based testing complexity and flakiness
- Solution: Use `bufconn` for in-memory gRPC testing in `TestHarness`
- Pattern: Always prefer in-memory testing for unit/integration tests

### Local Tool Management
- Issue: buf CLI version conflicts and system installation requirements
- Solution: Install tools locally in `bin/` directory with version pinning
- Pattern: `bin/toolname` with automatic installation in Makefile

### CI Pipeline Structure
- Issue: Missing integration test coverage and performance tracking
- Solution: Separate CI jobs for unit tests, integration tests, and benchmarks
- Pattern: Parallel job execution with artifact collection for benchmarks

## Best Practices Discovered

### Testing Framework Architecture
1. **Harness Pattern**: Use in-memory gRPC with bufconn for fast, reliable testing
2. **Mock Configurability**: Support error injection, delays, and custom behavior
3. **Conformance Levels**: Implement Basic/Standard/Advanced hierarchy for certification
4. **Performance Baselines**: Establish response time and memory usage benchmarks

### SDK Development Patterns  
1. **Generated vs Helper Code**: Separate protobuf generation from helper utilities
2. **Validation Integration**: Embed JSON schema for runtime validation
3. **Example Completeness**: Provide comprehensive cross-vendor examples
4. **Documentation Strategy**: Use specialized agents for technical writing

### CI/CD Optimization
1. **Tool Installation**: Local installation avoids version conflicts
2. **Validation Gates**: Generated code must be up-to-date in CI
3. **Test Coverage**: Unit → Integration → Conformance → Performance progression  
4. **Artifact Collection**: Store benchmark results for performance tracking

### Protocol Buffer Best Practices
1. **Forward Compatibility**: Use `UnimplementedServer` embedding
2. **Validation Functions**: Create comprehensive validators for all message types
3. **Error Handling**: Use proper gRPC status codes with meaningful messages
4. **Testing Support**: Design messages to support comprehensive testing scenarios

## Development Commands Reference

### Daily Development
```bash
# Setup development environment
make generate          # Install buf locally and generate code
make validate          # Run tests and linting

# Testing
go test -v ./sdk/go/testing/                    # Integration tests
go test -bench=. -benchmem ./sdk/go/testing/    # Performance benchmarks
go test -v -run TestConformance ./sdk/go/testing/  # Conformance tests

# Validation
make lint              # Run linting only
make test              # Run unit tests only
make clean-all         # Clean all generated files and tools
```

### Cross-Vendor Example Validation
```bash
# Validate all JSON examples against schema
for file in examples/specs/*.json; do 
    echo "Validating $file..."
    go run validate_examples.go "$file"
done
```

### Plugin Development Testing
```bash
# Test your plugin implementation
go test -v -run TestBasicPluginFunctionality
go test -v -run TestConformance  
go test -bench=BenchmarkAllMethods
```