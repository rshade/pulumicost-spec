# Data Model: FinFocus Package Mapping

## Protobuf Package Mapping

| Component | Old Package | New Package |
|-----------|-------------|-------------|
| Core | `pulumicost.v1` | `finfocus.v1` |

## Go Module Mapping

| Component | Old Path | New Path |
|-----------|----------|----------|
| Root Module | `github.com/rshade/pulumicost-spec` | `github.com/rshade/finfocus-spec` |
| Proto SDK | `.../sdk/go/proto/pulumicost/v1` | `.../sdk/go/proto/finfocus/v1` |
| Go Alias | `pbc` | `pbc` (Unchanged) |

## Implementation Constants

- **Version**: `v0.5.0`
- **Branding**: `FinFocus`
- **Tagline**: `Focusing your finances left`
