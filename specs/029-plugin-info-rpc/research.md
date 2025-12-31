<!-- markdownlint-disable MD013 -->
# Research: Add GetPluginInfo RPC

**Branch**: `029-plugin-info-rpc` | **Date**: 2025-12-30

## 1. Proto Definition

**Decision**: Add `GetPluginInfo` to `CostSource` service in `proto/pulumicost/v1/costsource.proto`.

**Rationale**:

- `CostSource` is the primary service implemented by plugins.
- Centralizes metadata retrieval alongside functional RPCs.
- Avoids creating a separate "MetadataService" for a single method, reducing complexity.

**Proto Definition**:

```protobuf
message GetPluginInfoRequest {}

message GetPluginInfoResponse {
  string name = 1;          // Plugin name (e.g., "aws-public")
  string version = 2;       // Plugin version (e.g., "v1.0.0")
  string spec_version = 3;  // Spec version plugin was built against (e.g., "v0.4.11")
  repeated string providers = 4; // Supported providers (e.g., ["aws"])
  map<string, string> metadata = 5; // Optional additional metadata
}

service CostSourceService {
  // ... existing RPCs ...
  rpc GetPluginInfo(GetPluginInfoRequest) returns (GetPluginInfoResponse);
}
```

## 2. SDK Implementation

**Decision**: Implement a default `GetPluginInfo` handler in `pluginsdk.BasePlugin` (or equivalent base struct) that plugins embed.

**Rationale**:

- Reduces boilerplate for plugin developers.
- Ensures consistent `spec_version` reporting (the SDK knows its own version).
- Allows plugins to override if dynamic metadata is needed, though the default should suffice for 99% of cases.

**Implementation Details**:

- `pluginsdk` package will export a constant `SpecVersion`.
- The base implementation will return this constant.
- `Name` and `Version` will be populated from the plugin's configuration/initialization struct.

## 3. Versioning Strategy

**Decision**: Use SemVer (vX.Y.Z) for `spec_version`.

**Rationale**:

- Standard in the Go ecosystem.
- Allows for compatibility ranges (e.g., core supports `^1.0.0`).
- Consistent with the project's existing versioning.

## 4. Compatibility

**Decision**: Handle `Unimplemented` error in Core.

**Rationale**:

- Existing plugins built against older protos will not implement `GetPluginInfo`.
- gRPC servers return `Unimplemented` (status code 12) for unknown methods.
- Core must treat this as "Unknown/Legacy Version" rather than a fatal error, logging a warning but proceeding (graceful degradation).

## 5. Alternatives Considered

- **Alternative**: CLI flag for version (e.g., `--version`).
  - Rejected: Inconsistent with gRPC-first architecture. Requires executing the binary just to check version, which is slower/heavier than an RPC if the connection is already established (or prevents establishing a connection if the binary interface changed).
- **Alternative**: Metadata in `Register` call (if one existed).
  - Rejected: We don't have a dynamic registration RPC; plugins are often loaded by path or discovered. An explicit query RPC is more robust.
