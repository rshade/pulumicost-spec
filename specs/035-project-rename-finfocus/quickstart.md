# Quickstart: Migrating to FinFocus

## For SDK Consumers

Update your `go.mod` to use the new module name:

```bash
go mod edit -module github.com/rshade/finfocus-spec
go get github.com/rshade/finfocus-spec/sdk/go/pluginsdk
```

Update your imports:

```go
import (
    // Old
    // pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
    
    // New
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)
```

## For Plugin Developers

1. Re-run `make generate` to produce the new `finfocus.v1` bindings.
2. Update your `CostSourceServiceServer` implementation to use the new package names.
3. Verify your `Name()` RPC returns the correct plugin name.

## Tagline Usage

The primary tagline for the project is:
> **FinFocus: Focusing your finances left**
