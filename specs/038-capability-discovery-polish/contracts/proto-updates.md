# Contracts: Capability Discovery Polish

**Type**: gRPC Protobuf
**File**: `proto/finfocus/v1/costsource.proto`

## Changes

**Status**: Documentation Only. No wire format changes.

### GetPluginInfoResponse

Added comments to clarify deprecation status.

```protobuf
message GetPluginInfoResponse {
  // ...
  
  // Modern capability format using strongly-typed enums.
  // Prefer this field for capability queries on newer clients.
  // SDK auto-populates this based on implemented interfaces.
  repeated PluginCapability capabilities = 5; 

  // Legacy metadata format for backward compatibility with older hosts.
  // Contains string-based capability flags: {"supports_xyz": "true"}.
  // SDK auto-populates this from capabilities for backward compatibility.
  // DEPRECATION: New integrations should use capabilities field instead.
  map<string, string> metadata = 6;
}
```

### SupportsResponse

Added comments to clarify deprecation status.

```protobuf
message SupportsResponse {
  // ...
  
  // Legacy capability format using boolean map.
  // Example: {"recommendations": true, "dry_run": true}
  // DEPRECATION: Use capabilities_enum for new integrations.
  map<string, bool> capabilities = 3;

  // Modern capability format using strongly-typed enums.
  // Auto-populated by SDK based on implemented interfaces.
  repeated PluginCapability capabilities_enum = 5;
}
```
