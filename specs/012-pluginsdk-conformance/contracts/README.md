# Contracts: PluginSDK Conformance Testing Adapters

**Date**: 2025-11-30
**Feature**: 012-pluginsdk-conformance

## Overview

This feature does not introduce new external API contracts. The adapter functions are internal
Go SDK utilities that wrap existing functionality.

## Why No API Contracts

1. **Internal Package**: `pluginsdk` is a Go package, not an HTTP/gRPC service
2. **Type Aliases**: Uses existing types from `sdk/go/testing` via aliases
3. **No Protocol Changes**: No new protobuf messages or gRPC methods
4. **No External Integration**: No external services or data formats

## Existing Contracts Referenced

The adapter functions delegate to existing contracts:

- **Proto Contract**: `proto/pulumicost/v1/costsource.proto` (unchanged)
- **Testing Contract**: `sdk/go/testing/conformance.go` function signatures (consumed)

## Go API Contract (Documentation)

For Go consumers, the API contract is defined by:

1. **Function Signatures**: Documented in `data-model.md`
2. **Type Definitions**: Re-exported via type aliases
3. **Error Behavior**: Documented in function godoc comments
4. **Package Documentation**: `pluginsdk/conformance.go` package comment

See `quickstart.md` for usage examples.
