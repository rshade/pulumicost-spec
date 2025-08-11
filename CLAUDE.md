# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **pulumicost-spec**, a repository that provides the canonical protocol and schemas for PulumiCost plugins. It defines:
- gRPC service definitions for cost source plugins
- JSON schemas for pricing specifications
- Go SDK with generated protobuf code and helper types

## Build Commands

- `make generate` - Generate Go code from protobuf definitions (requires buf CLI)
- `make tidy` - Run `go mod tidy` to clean up dependencies  
- `make test` - Run all Go tests
- `make clean` - Remove generated proto files
- `go build ./...` - Build all Go packages
- `go test ./...` - Run tests (currently no test files exist)

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

### Generated Code

The `sdk/go/proto/` directory contains generated Go protobuf code. To regenerate:
1. Install buf CLI: https://docs.buf.build/installation
2. Run `make generate`

### Code Generation Dependencies

- **buf** - Protocol buffer toolchain for generation
- **google.golang.org/protobuf** - Go protobuf runtime
- **google.golang.org/grpc** - gRPC Go implementation

## Development Workflow

1. **Modify Proto**: Edit `proto/pulumicost/costsource.proto`
2. **Update Schema**: Edit `schemas/pricing_spec.schema.json` if PricingSpec message changes
3. **Regenerate**: Run `make generate` to update Go bindings
4. **Update Types**: Modify helper code in `sdk/go/types/` as needed
5. **Test**: Run `make test` and `go build ./...` to verify compilation

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