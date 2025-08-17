# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **pulumicost-spec**, a repository that provides the canonical protocol and schemas for PulumiCost plugins. It defines:

- gRPC service definitions for cost source plugins
- JSON schemas for pricing specifications
- Go SDK with generated protobuf code and helper types

## Build Commands

### Core Build Commands

- `make generate` - Generate Go code from protobuf definitions (installs buf locally in bin/)
- `make tidy` - Run `go mod tidy` to clean up dependencies  
- `make test` - Run all Go tests including integration tests
- `make validate` - Run tests, linting, and npm validations together
- `make clean` - Remove generated proto files
- `make clean-all` - Remove generated files and local tools (bin/)
- `go build ./...` - Build all Go packages
- `go test -bench=. -benchmem ./sdk/go/testing/` - Run performance benchmarks

### Linting Commands

- `make lint` - Run all linting (Go, buf, markdown, and YAML)
- `make lint-go` - Run Go linting (golangci-lint and buf lint)
- `make lint-markdown` - Run markdown linting with markdownlint-cli2
- `make lint-markdown-fix` - Auto-fix markdown linting issues
- `make lint-yaml` - Run YAML linting on GitHub workflows
- `make lint-yaml-fix` - Auto-fix YAML linting issues

### NPM/Schema Validation Commands

- `make validate-schema` - Validate JSON schema syntax
- `make validate-examples` - Validate example files against schema
- `make validate-npm` - Run all npm validations (schema + examples)
- `npm run lint:markdown` - Direct npm markdown linting
- `npm run validate` - Direct npm validation command

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
- `sdk/go/pricing/domain.go` - BillingMode enum constants and validation helpers
- `sdk/go/pricing/validate.go` - JSON schema validation for PricingSpec documents
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
- **Node.js >=22** - Required for npm commands and markdown linting
- **markdownlint-cli2** - Markdown linting tool
- **ajv** - JSON schema validation
- **yamllint** - YAML linting tool (install with `pip install yamllint` or `brew install yamllint`)

### Local Tool Management

The project uses local tool installation to avoid version conflicts:

- `bin/buf` - buf CLI v1.32.1 installed automatically
- Tools are excluded from git via `.gitignore`
- `make clean-all` removes all local tools

## Development Workflow

1. **Setup**: Ensure Node.js >=22 and npm >=10 are installed, then run `npm install`
2. **Modify Proto**: Edit `proto/pulumicost/costsource.proto`
3. **Update Schema**: Edit `schemas/pricing_spec.schema.json` if PricingSpec message changes
4. **Regenerate**: Run `make generate` to update Go bindings
5. **Update Types**: Modify helper code in `sdk/go/pricing/` as needed
6. **Test**: Run `make validate` to run tests, linting, and npm validations
7. **Verify**: Run integration tests with `go test -v ./sdk/go/testing/`
8. **Format**: Use `make lint-markdown-fix` to auto-fix markdown issues

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

```text
github.com/rshade/pulumicost-spec/sdk/go/proto   # Generated protobuf code
github.com/rshade/pulumicost-spec/sdk/go/pricing # Domain types and validation (formerly 'types')
github.com/rshade/pulumicost-spec/sdk/go/testing # Plugin testing framework
```

## Schema Validation

The pricing package embeds the JSON schema and provides `ValidatePricingSpec(doc []byte) error`
for validating PricingSpec JSON documents against the schema.

## Versioning

Follow semantic versioning for proto changes:

- MAJOR: Breaking proto changes
- MINOR: Backward-compatible additions  
- PATCH: Bug fixes, documentation

Tag releases as `v0.1.0`, `v1.0.0`, etc.

- Project board: <https://github.com/users/rshade/projects/3>

## Common Issues & Solutions

### YAML Linting Configuration

- Issue: yamllint errors with default configuration
- Solution: Created `.yamllint` configuration file with sensible defaults
- Configuration disables document-start rule and sets line length to 120
- Use `make lint-yaml` to check YAML files and `make lint-yaml-fix` to auto-fix issues

### Dot Import Linting Issues

- Issue: golangci-lint flags dot imports (`. "package"`) as style violations
- Solution: Replace with explicit imports and use package prefixes
- Pattern: Change `import . "pkg"` to `import "pkg"` and update all function calls to use `pkg.Function()`
- Special case: When importing custom package with same name as stdlib (e.g., `testing`),
  use import alias: `import plugintesting "custom/testing"`

### Package Naming Conflicts  

- Issue: Import name conflicts between stdlib and custom packages (e.g., `testing` vs custom `testing` package)
- Solution: Use import aliases to disambiguate: `import plugintesting "github.com/repo/testing"`
- Pattern: Rename one of the imports with a descriptive alias, typically the custom package

### Package Renaming Process

- Issue: Need to rename package for better naming conventions
- Solution: Systematic approach to avoid breaking changes:
  1. `mv old_package new_package` (rename directory)
  2. Update `package` declarations in all `.go` files
  3. Update import paths in all files
  4. Update package references in code (e.g., `old.Function()` to `new.Function()`)
  5. Update test package names (`package old_test` → `package new_test`)
- Verification: Run `go build ./...`, `make test`, and `make lint` to ensure no breakage

### Mock Plugin Implementation

- Issue: gRPC method name conflicts in mock plugins
- Solution: Use `PluginName` field instead of `Name` to avoid conflicts with RPC method names
- Pattern: Separate data fields from method names in struct design

### Integration Testing Setup

- Issue: Network-based testing complexity and flakiness
- Solution: Use `bufconn` for in-memory gRPC testing in `TestHarness`
- Pattern: Always prefer in-memory testing for unit/integration tests

### Tool Management Issues

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
npm install            # Install npm dependencies
make generate          # Install buf locally and generate code
make validate          # Run tests, linting, and npm validations

# Testing
go test -v ./sdk/go/testing/                    # Integration tests
go test -bench=. -benchmem ./sdk/go/testing/    # Performance benchmarks
go test -v -run TestConformance ./sdk/go/testing/  # Conformance tests

# Linting
make lint              # Run all linting (Go, buf, markdown, YAML)
make lint-markdown     # Run markdown linting only
make lint-markdown-fix # Auto-fix markdown issues
make lint-yaml         # Run YAML linting
make lint-yaml-fix     # Auto-fix YAML issues

# Schema Validation
make validate-schema   # Validate JSON schema syntax
make validate-examples # Validate example files against schema
make validate-npm      # Run all npm validations

# Cleanup
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
