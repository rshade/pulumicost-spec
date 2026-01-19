feat(sdk): polish capability discovery with auto-detection

## Summary

This PR consolidates the capability discovery logic into a robust,
performant, and auto-detecting system. It introduces a centralized
compatibility layer for legacy string-map metadata, enables the SDK to
automatically infer capabilities from implemented Go interfaces (like
`DryRunHandler`), and optimizes memory allocation in high-traffic paths.

## Test plan

- [x] `go test ./sdk/go/pluginsdk/...` passes (verified auto-discovery)
- [x] `go test -bench=. ./sdk/go/pluginsdk/...` shows zero-alloc improvements
- [x] `golangci-lint run` passes with zero issues
- [x] Validated backward compatibility with legacy metadata formats

## Changes

### New files

- `sdk/go/pluginsdk/capability_compat.go` - Centralized Enum/Metadata conversion

### Modified files

- `sdk/go/pluginsdk/plugin_info.go` - Implemented interface-based auto-discovery
- `sdk/go/pluginsdk/sdk.go` - Integrated new discovery logic and optimized slices
- `proto/finfocus/v1/costsource.proto` - Added clarifying comments to messages
- `sdk/go/pluginsdk/README.md` - Documented the Capability Discovery Pattern

### Housekeeping

- Renamed compatibility tests in `sdk/go/testing` for clarity

Closes #208
Closes #209
Closes #294
Closes #295
Closes #299
Closes #300
Closes #301
