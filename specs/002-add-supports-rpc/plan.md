# Implementation Plan: Add Supports() RPC Method to CostSourceService

**Branch**: `002-add-supports-rpc` | **Date**: 2025-11-20 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-add-supports-rpc/spec.md`

## Summary

**FINDING: The Supports() RPC method is fully implemented in finfocus-spec v0.1.0.**

### Verified Components

The proto file `proto/finfocus/v1/costsource.proto` contains:

- `rpc Supports(SupportsRequest) returns (SupportsResponse);` (line 17)
- `SupportsRequest` message with `ResourceDescriptor resource` field (lines 38-42)
- `SupportsResponse` message with `supported` and `reason` fields (lines 44-50)

The generated Go SDK in `sdk/go/proto/finfocus/v1/` includes:

- `CostSourceServiceClient.Supports()` method (line 67 in costsource_grpc.pb.go)
- `CostSourceServiceServer.Supports()` interface method (line 118)
- `UnimplementedCostSourceServiceServer.Supports()` default implementation (line 138)
- `_CostSourceService_Supports_Handler` gRPC handler (line 189)
- Service descriptor entry in `CostSourceService_ServiceDesc.Methods` (lines 272-275)

The testing framework includes:

- `ValidateSupportsResponse()` validator in harness.go
- `MockPlugin.Supports()` implementation in mock_plugin.go
- Integration tests for Supports in integration_test.go
- Error handling tests including Supports method

### Root Cause Analysis

The issue **is NOT in finfocus-spec**. The Supports RPC was included in the v0.1.0 release.

Per the issue's "Related Issues" section:

> "After this is implemented, finfocus-core#TBD needs to be completed to actually register
> and implement the handler in `pluginsdk`."

**The actual problem is in finfocus-core's pluginsdk**, which likely doesn't expose the
Supports method to plugin implementations. When a plugin calls `client.Supports()`, the
pluginsdk needs to forward this to the plugin's implementation.

### Conclusion

GitHub Issue #64 should be closed as "Already Complete" for finfocus-spec. A new issue
should be created in finfocus-core to update the pluginsdk to expose the Supports method.

## Technical Context

**Language/Version**: Go 1.25.4
**Primary Dependencies**: buf v1.32.1, google.golang.org/protobuf, google.golang.org/grpc
**Storage**: N/A
**Testing**: go test, buf lint, buf breaking
**Target Platform**: Linux/macOS (gRPC services)
**Project Type**: Single - gRPC protocol specification repository
**Performance Goals**: N/A for proto definition
**Constraints**: Backward compatibility with existing proto consumers
**Scale/Scope**: Single RPC method addition (already complete)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                           | Status  | Notes                                   |
| ----------------------------------- | ------- | --------------------------------------- |
| I. gRPC Proto Specification-First   | ✅ Pass | Proto already defines Supports RPC      |
| II. Multi-Provider Consistency      | ✅ Pass | ResourceDescriptor is provider-agnostic |
| III. Test-First Protocol            | ⚠️ N/A  | No changes needed - already implemented |
| IV. Protobuf Backward Compatibility | ✅ Pass | No breaking changes                     |
| V. Comprehensive Documentation      | ✅ Pass | Proto comments exist for all messages   |
| VI. Performance as Requirement      | ✅ Pass | Standard RPC method                     |
| VII. Validation at Multiple Levels  | ✅ Pass | buf lint passes                         |

## Project Structure

### Documentation (this feature)

```text
specs/002-add-supports-rpc/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Not needed - feature already exists
├── checklists/
│   └── requirements.md  # Quality checklist
└── tasks.md             # Not needed - no implementation required
```

### Source Code (repository root)

```text
proto/finfocus/v1/
└── costsource.proto     # Already contains Supports RPC (lines 15-17, 38-50)

sdk/go/proto/finfocus/v1/
├── costsource.pb.go     # Already contains SupportsRequest/Response types
└── costsource_grpc.pb.go # Already contains Supports client/server methods
```

**Structure Decision**: No changes required - feature already implemented.

## Complexity Tracking

No violations - feature already exists in codebase.

## Recommended Actions

1. **Close GitHub Issue #64** with comment explaining:
   - The Supports RPC is fully implemented in finfocus-spec v0.1.0
   - All components verified: proto, messages, handlers, service descriptor, testing framework
   - The actual issue is in finfocus-core pluginsdk

2. **Unblock finfocus-core#160** by commenting:
   - Issue: <https://github.com/rshade/finfocus-core/issues/160>
   - The dependency (finfocus-spec#64) is already complete
   - Supports RPC has been in v0.1.0 since release
   - Issue #160 can now proceed with pluginsdk implementation
   - Suggested comment:

   ```text
   This issue is now unblocked. finfocus-spec#64 is already complete - the Supports()
   RPC method has been in finfocus-spec since v0.1.0 release:

   - Proto: `rpc Supports(SupportsRequest) returns (SupportsResponse);`
   - Messages: SupportsRequest, SupportsResponse
   - Generated code: Client/Server interfaces, handlers, service descriptor
   - Testing: ValidateSupportsResponse, MockPlugin.Supports, integration tests

   Update finfocus-core to use v0.1.0 or later and proceed with pluginsdk implementation.
   ```

3. **Verify finfocus-plugin-aws-public** is using:
   - `github.com/rshade/finfocus-spec v0.1.0` or later
   - Run `go get github.com/rshade/finfocus-spec@v0.1.0` if needed

## Verification Commands

```bash
# Verify proto contains Supports RPC
grep -n "rpc Supports" proto/finfocus/v1/costsource.proto

# Verify generated code contains Supports method
grep "func.*Supports" sdk/go/proto/finfocus/v1/*.go

# Build and test
make test
make lint
```
