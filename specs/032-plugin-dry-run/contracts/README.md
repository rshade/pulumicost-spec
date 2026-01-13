# Contracts: Plugin Capability Dry Run Mode

This directory contains the gRPC/proto contract definitions for the dry-run feature.

## Files

### dry_run.proto

Proto snippet showing all additions required for the dry-run feature. This is not a
standalone proto file but a reference showing what will be merged into:

- `proto/finfocus/v1/enums.proto` - FieldSupportStatus enum
- `proto/finfocus/v1/costsource.proto` - Messages and RPC definition

## Implementation Notes

### Field Numbers

When implementing, use the following field numbers:

| Message | Field | Number |
|---------|-------|--------|
| GetActualCostRequest | dry_run | 6 |
| GetActualCostResponse | dry_run_result | 3 |
| GetProjectedCostRequest | dry_run | 5 |
| GetProjectedCostResponse | dry_run_result | 6 |

### Backward Compatibility

- All new fields are optional (proto3 default)
- `dry_run` flag defaults to false (existing behavior unchanged)
- DryRun RPC returns Unimplemented for legacy plugins
- `SupportsResponse.capabilities["dry_run"]` indicates support

### Validation Requirements

After proto changes:

1. Run `make generate` to regenerate Go code
2. Run `buf lint` to validate proto style
3. Run `buf breaking` to verify no unintended breaking changes
4. Add conformance tests for DryRun RPC behavior
