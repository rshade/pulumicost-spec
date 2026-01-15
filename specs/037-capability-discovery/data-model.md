# Data Model - Dual-Layer Capability Discovery

## Protobuf Schema Changes

### `proto/finfocus/v1/costsource.proto`

#### `SupportsResponse`

Current:

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2;
  map<string, bool> capabilities = 3;
  repeated MetricKind supported_metrics = 4;
}
```

Proposed:

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2;
  map<string, bool> capabilities = 3;
  repeated MetricKind supported_metrics = 4;
  // capabilities_enum declares optional capabilities using strongly-typed enums.
  // This field is automatically populated by the SDK based on implemented interfaces
  // or inherited from global capabilities if left empty.
  repeated PluginCapability capabilities_enum = 5;
}
```

### `proto/finfocus/v1/enums.proto`

(Verify if `PluginCapability` enum exists and has required values. If not, add them.)

Assuming `PluginCapability` needs to be defined or updated:

```protobuf
enum PluginCapability {
  PLUGIN_CAPABILITY_UNSPECIFIED = 0;
  PLUGIN_CAPABILITY_ACTUAL_COSTS = 1;
  PLUGIN_CAPABILITY_PROJECTED_COSTS = 2;
  PLUGIN_CAPABILITY_PRICING_SPEC = 3;
  PLUGIN_CAPABILITY_ESTIMATE_COST = 4;
  PLUGIN_CAPABILITY_RECOMMENDATIONS = 5;
  PLUGIN_CAPABILITY_BUDGETS = 6;
  PLUGIN_CAPABILITY_DISMISS_RECOMMENDATIONS = 7;
  PLUGIN_CAPABILITY_DRY_RUN = 8;
  // Add others as needed (Carbon/Energy/Water are MetricKinds, not functional capabilities usually, but check spec)
}
```

## Go SDK Data Structures

### `Server` struct update

Add a field to store globally discovered capabilities.

```go
type Server struct {
    // ... existing fields
    globalCapabilities []pbc.PluginCapability
}
```

### `Supports` Method Logic

```go
func (s *Server) Supports(...) {
    // ... existing logic ...
    resp, err := s.plugin.Supports(...)
    // ...
    
    // Auto-population logic
    if len(resp.CapabilitiesEnum) == 0 {
        // Inherit Global
        resp.CapabilitiesEnum = s.globalCapabilities
    }
    
    // Sync to legacy map
    if len(resp.Capabilities) == 0 {
        resp.Capabilities = convertEnumToMap(resp.CapabilitiesEnum)
    }
}
```
