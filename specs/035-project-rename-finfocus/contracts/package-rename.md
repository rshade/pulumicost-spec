# API Contracts: Package Rename

## Overview

This document describes the API contract changes resulting from the PulumiCost to
FinFocus rename. Since this is a mechanical rename rather than new functionality,
the contracts focus on package and module path changes.

## Protobuf Package Contract

### Package Declaration Change

**Old Package**: `pulumicost.v1`
**New Package**: `finfocus.v1`

### Proto File Structure

```protobuf
// Before
package pulumicost.v1;
option go_package = "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1;pbc";

// After
package finfocus.v1;
option go_package = "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1;pbc";
```

### Directory Structure Change

```text
Before:
proto/pulumicost/
├── v1/
│   ├── costsource.proto
│   └── budget.proto

After:
proto/finfocus/
├── v1/
│   ├── costsource.proto
│   └── budget.proto
```

### Service Names

Service names remain **unchanged** (generic):

- `CostSourceService`
- `BudgetService`

Message definitions remain **unchanged**:

- `CostSource`
- `PricingSpec`
- `Budget`
- `ResourceDescriptor`

## Go Module Contract

### Module Path Change

**Old Module**: `github.com/rshade/pulumicost-spec`
**New Module**: `github.com/rshade/finfocus-spec`

### Import Path Changes

| Purpose      | Old Import                                                     | New Import                                                 |
| ------------ | -------------------------------------------------------------- | ---------------------------------------------------------- |
| Plugin SDK   | `github.com/rshade/pulumicost-spec/sdk/go/pluginsdk`           | `github.com/rshade/finfocus-spec/sdk/go/pluginsdk`         |
| Pulumi API   | `github.com/rshade/pulumicost-spec/sdk/go/pulumiapi`           | `github.com/rshade/finfocus-spec/sdk/go/pulumiapi`         |
| Proto Client | `github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1` | `github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1` |

### Package Alias

**Unchanged**: The `pbc` alias for the generated proto package is preserved.

## Release Contract

### Version Bump

**Previous Version**: v0.1.0
**New Version**: v0.5.0

**Rationale**: This is a breaking change at the module/package level, justifying a
MINOR version bump in pre-1.0 semantic versioning.

### Release-Please Configuration

```json
// Before
{
  "package-name": "pulumicost-spec",
  "extra-files": [
    ...
  ]
}

// After
{
  "package-name": "finfocus-spec",
  "extra-files": [
    ...
  ]
}
```

## Backward Compatibility

### Wire Protocol

**Compatible**: The gRPC wire format remains unchanged. Only the package namespace is updated.

### Generated Code

**Breaking Change**: Generated Go code will have different package names, requiring import updates in downstream projects.

### Service Contracts

**Compatible**: Service method signatures and message types remain unchanged at the wire level.

## Migration Guide

### For SDK Consumers

1. Update `go.mod`:

   ```bash
   go get github.com/rshade/finfocus-spec/sdk/go/pluginsdk
   ```

2. Update imports:

   ```go
   // Before
   import pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"

   // After
   import pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
   ```

### For Plugin Developers

1. Regenerate proto code:

   ```bash
   make clean
   make generate
   ```

2. Update import paths in plugin implementations.

3. Verify service implementation interfaces still match.

## Breaking Changes Summary

| Change Type       | Impact                         | Migration Required          |
| ----------------- | ------------------------------ | --------------------------- |
| Protobuf package  | Wire-level package name change | Proto regeneration required |
| Go module path    | Import path changes            | `go get` and import updates |
| Generated code    | Package name changes           | Import updates in consumers |
| Service contracts | None                           | No migration required       |

## Validation

All changes must pass the following validation:

1. `make generate` produces clean proto code in `sdk/go/proto/finfocus/v1/`
2. `make test` passes with zero `pulumicost-spec` imports
3. `make validate` passes all linting and schema checks
4. `rg pulumicost` returns only historical references

## Compliance with Constitution

This rename complies with all constitutional principles:

- **Contract First**: Proto files are updated before SDK code
- **Backward Compatibility**: Breaking changes are justified and properly versioned
- **Multi-Provider**: Provider-agnostic patterns are preserved
- **FinFocus Alignment**: This is the primary implementation of the FinFocus transition
