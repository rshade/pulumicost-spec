# GetPluginInfo Error Message Contract

**Feature**: SDK Polish v0.4.15
**Date**: 2026-01-10
**Type**: gRPC Error Contract

## Overview

This contract documents the user-friendly error messages returned by the `GetPluginInfo` RPC call.

---

## RPC Definition

### Method Signature

```go
// GetPluginInfo returns metadata about the plugin including name, version,
// spec version, supported providers, and optional metadata.
func (s *Server) GetPluginInfo(
    ctx context.Context,
    req *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error)
```

### Request

```protobuf
message GetPluginInfoRequest {
    // Request body is empty (no parameters needed)
}
```

### Response

```protobuf
message GetPluginInfoResponse {
    string name = 1;            // Plugin name
    string version = 2;         // Plugin version
    string spec_version = 3;      // Specification version (e.g., "v1.2.0")
    repeated string providers = 4;  // Supported cloud providers
    map<string, string> metadata = 5; // Optional metadata key-value pairs
}
```

---

## Error Messages

### Error Hierarchy

```
Server.GetPluginInfo()
├─ Plugin implements PluginInfoProvider?
│  ├─ Yes: Delegate to plugin.GetPluginInfo()
│  │  ├─ Error? → Log + Return "unable to retrieve plugin metadata"
│  │  ├─ Nil response? → Return "unable to retrieve plugin metadata"
│  │  ├─ Incomplete metadata? → Return "plugin metadata is incomplete"
│  │  └─ Invalid spec_version? → Return "plugin reported an invalid specification version"
│  └─ No: Static PluginInfo configured?
│     ├─ Yes: Return configured PluginInfo
│     └─ No: Return Unimplemented (legacy plugin)
```

### Error Message Matrix

| Error Code    | Message                                            | Condition                                                                 | gRPC Status           |
| ------------- | -------------------------------------------------- | ------------------------------------------------------------------------- | --------------------- |
| Internal      | "unable to retrieve plugin metadata"               | Plugin returns nil or error when implementing PluginInfoProvider          | `codes.Internal`      |
| Internal      | "plugin metadata is incomplete"                    | Required fields (name, version, spec_version) are empty                   | `codes.Internal`      |
| Internal      | "plugin reported an invalid specification version" | spec_version fails ValidateSpecVersion()                                  | `codes.Internal`      |
| Unimplemented | "GetPluginInfo not implemented"                    | Plugin does not implement PluginInfoProvider and no PluginInfo configured | `codes.Unimplemented` |

---

## Error Conditions

### 1. Plugin Returns Nil Response

**FR-006**: The `GetPluginInfo` RPC MUST return "unable to retrieve plugin metadata" when the plugin returns nil.

**Trigger Conditions**:

- Plugin implements `PluginInfoProvider` interface
- Plugin's `GetPluginInfo()` method returns `nil, nil` (or nil, error)

**Server Behavior**:

```go
resp, err := provider.GetPluginInfo(ctx, req)
if err != nil {
    s.logger.Error().Err(err).Msg("GetPluginInfo handler error")
    return nil, status.Error(codes.Internal, "unable to retrieve plugin metadata")
}
if resp == nil {
    s.logger.Error().Msg("GetPluginInfo returned nil response")
    return nil, status.Error(codes.Internal, "unable to retrieve plugin metadata")
}
```

**Client Behavior**:

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    st, ok := status.FromError(err)
    if ok && st.Code() == codes.Internal {
        fmt.Printf("Error: %s\n", st.Message())
        // Output: "unable to retrieve plugin metadata"
    }
}
```

---

### 2. Plugin Returns Incomplete Metadata

**FR-007**: The `GetPluginInfo` RPC MUST return "plugin metadata is incomplete" when required fields are empty.

**Required Fields**:

- `name` - Plugin name identifier
- `version` - Plugin version string
- `spec_version` - Specification version (must be valid)

**Trigger Conditions**:

- Plugin implements `PluginInfoProvider` interface
- Plugin's `GetPluginInfo()` returns response with one or more empty required fields

**Server Behavior**:

```go
if resp.GetName() == "" || resp.GetVersion() == "" || resp.GetSpecVersion() == "" {
    s.logger.Error().
        Str("name", resp.GetName()).
        Str("version", resp.GetVersion()).
        Str("spec_version", resp.GetSpecVersion()).
        Msg("GetPluginInfo returned incomplete response")
    return nil, status.Error(codes.Internal, "plugin metadata is incomplete")
}
```

**Client Behavior**:

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    st, ok := status.FromError(err)
    if ok && st.Code() == codes.Internal {
        fmt.Printf("Error: %s\n", st.Message())
        // Output: "plugin metadata is incomplete"
    }
}
```

---

### 3. Plugin Reports Invalid Spec Version

**FR-008**: The `GetPluginInfo` RPC MUST return "plugin reported an invalid specification version" for malformed spec_version.

**Valid Spec Version Formats**:

- `vX.Y.Z` (e.g., "v1.2.0")
- Must start with "v"
- Must be parsable semantic version

**Trigger Conditions**:

- Plugin implements `PluginInfoProvider` interface
- Plugin's `GetPluginInfo()` returns response with invalid `spec_version` format

**Server Behavior**:

```go
if specErr := ValidateSpecVersion(resp.GetSpecVersion()); specErr != nil {
    s.logger.Error().
        Err(specErr).
        Msg("GetPluginInfo returned invalid spec_version")
    return nil, status.Error(codes.Internal, "plugin reported an invalid specification version")
}
```

**Client Behavior**:

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    st, ok := status.FromError(err)
    if ok && st.Code() == codes.Internal {
        fmt.Printf("Error: %s\n", st.Message())
        // Output: "plugin reported an invalid specification version"
    }
}
```

---

### 4. Plugin Does Not Implement GetPluginInfo (Legacy)

**Graceful Degradation**: Legacy plugins that don't implement `GetPluginInfo` return `Unimplemented` error.

**Trigger Conditions**:

- Plugin does NOT implement `PluginInfoProvider` interface
- No `PluginInfo` configured in `ServeConfig`

**Server Behavior**:

```go
s.logger.Debug().
    Str("plugin", s.plugin.Name()).
    Msg("GetPluginInfo not implemented (legacy plugin)")
return nil, status.Error(codes.Unimplemented, "GetPluginInfo not implemented")
```

**Client Behavior** (Recommended):

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    if status.Code(err) == codes.Unimplemented {
        // Legacy plugin - use fallback values
        log.Info("Plugin does not implement GetPluginInfo")
        return &PluginMetadata{
            Name:    "unknown",
            Version: "unknown",
        }
    }
    return nil, fmt.Errorf("GetPluginInfo failed: %w", err)
}
```

---

## Server-Side Logging

### Error Logging Pattern

All errors are logged server-side with detailed information for debugging:

```go
// Pattern 1: Handler error
s.logger.Error().
    Err(err).  // Original error with stack trace
    Msg("GetPluginInfo handler error")

// Pattern 2: Nil response
s.logger.Error().
    Msg("GetPluginInfo returned nil response")

// Pattern 3: Incomplete metadata
s.logger.Error().
    Str("name", resp.GetName()).
    Str("version", resp.GetVersion()).
    Str("spec_version", resp.GetSpecVersion()).
    Msg("GetPluginInfo returned incomplete response")

// Pattern 4: Invalid spec version
s.logger.Error().
    Err(specErr).  // Validation error details
    Msg("GetPluginInfo returned invalid spec_version")

// Pattern 5: Legacy plugin (debug level)
s.logger.Debug().
    Str("plugin", s.plugin.Name()).
    Msg("GetPluginInfo not implemented (legacy plugin)")
```

### Logging Levels

| Error Type                      | Level | Rationale                             |
| ------------------------------- | ----- | ------------------------------------- |
| Plugin handler error            | Error | Unexpected error, needs investigation |
| Nil/incomplete/invalid response | Error | Bug in plugin implementation          |
| Legacy plugin                   | Debug | Expected for older plugins            |

### Client vs Server Messages

| Server Log                                                            | Client Message                                     | Security                     |
| --------------------------------------------------------------------- | -------------------------------------------------- | ---------------------------- |
| "GetPluginInfo returned nil response"                                 | "unable to retrieve plugin metadata"               | ✅ No implementation details |
| "GetPluginInfo handler error: <stack trace>"                          | "unable to retrieve plugin metadata"               | ✅ No stack trace leaked     |
| "GetPluginInfo returned invalid spec_version: expected format vX.Y.Z" | "plugin reported an invalid specification version" | ✅ No validation details     |

---

## Testing Requirements

### Unit Tests

- ✅ `TestGetPluginInfo/provider_implements_GetPluginInfoProvider_returns_metadata` - Basic implementation test
- ✅ `TestGetPluginInfo/provider_does_not_implement_GetPluginInfoProvider_returns_Unimplemented` - Legacy plugin test
- ✅ `TestGetPluginInfo/plugin_returns_nil_returns_error` - Nil response test
- ✅ `TestGetPluginInfo/plugin_returns_incomplete_returns_error` - Incomplete metadata test
- ✅ `TestGetPluginInfo/plugin_returns_invalid_spec_version_returns_error` - Invalid spec version test

### Conformance Tests

- ❌ **Missing**: Verify error message format matches FR-006, FR-007, FR-008
- ❌ **Missing**: Verify server-side logging occurs with detailed info
- ❌ **Missing**: Verify legacy plugin Unimplemented error is handled gracefully

### Integration Tests

- ❌ **Missing**: Mock plugin returns nil → client receives "unable to retrieve plugin metadata"
- ❌ **Missing**: Mock plugin returns incomplete → client receives "plugin metadata is incomplete"
- ❌ **Missing**: Mock plugin returns invalid spec_version → client receives "plugin reported an invalid specification version"
- ❌ **Missing**: Legacy plugin → client handles Unimplemented gracefully

---

## Validation Rules

### Response Validation

```go
func (s *Server) validateGetPluginInfoResponse(resp *pbc.GetPluginInfoResponse) error {
    if resp == nil {
        return errors.New("unable to retrieve plugin metadata")
    }
    if resp.GetName() == "" || resp.GetVersion() == "" || resp.GetSpecVersion() == "" {
        return errors.New("plugin metadata is incomplete")
    }
    if err := ValidateSpecVersion(resp.GetSpecVersion()); err != nil {
        return errors.New("plugin reported an invalid specification version")
    }
    return nil
}
```

### Spec Version Validation

```go
// ValidateSpecVersion validates the spec_version format (e.g., "v1.2.0").
func ValidateSpecVersion(version string) error {
    if version == "" {
        return errors.New("spec_version is required")
    }
    if !strings.HasPrefix(version, "v") {
        return fmt.Errorf("spec_version must start with 'v', got: %s", version)
    }
    // Additional semantic version validation...
    return nil
}
```

---

## Backward Compatibility

### Legacy Plugin Support

Plugins that don't implement `GetPluginInfo`:

- Return `Unimplemented` error (not a failure)
- Client should handle gracefully (use fallback values)
- Server logs at Debug level (not error)

### Error Message Stability

Error messages are **contractual obligations**:

- Changes require major version bump
- Existing plugins depend on these messages
- Client code parses these messages

---

## Security Considerations

### Information Disclosure Prevention

**What Clients See**:

- Generic error messages ("unable to retrieve plugin metadata")
- No stack traces
- No validation details
- No internal state information

**What Server Logs** (admin access only):

- Full error messages with stack traces
- Validation failure details
- Plugin response state

### DoS Protection

- `GetPluginInfo` is fast (no external API calls)
- No timeout protection needed (sub-millisecond operation)
- Server implements graceful degradation for legacy plugins

---

## Migration Guide

### For Plugin Developers

**Old Behavior** (Legacy Plugin):

```go
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

// No GetPluginInfo() method
```

**New Behavior** (With GetPluginInfo):

```go
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) GetPluginInfo(
    ctx context.Context,
    req *pbc.GetPluginInfoRequest,
) (*pbc.GetPluginInfoResponse, error) {
    // MUST return non-nil response with required fields
    return &pbc.GetPluginInfoResponse{
        Name:        "my-plugin",
        Version:     "1.0.0",
        SpecVersion: "v1.2.0",
        Providers:   []string{"aws", "azure"},
    }, nil
}
```

### For Client Developers

**Handling Legacy Plugins**:

```go
resp, err := client.GetPluginInfo(ctx, &pbc.GetPluginInfoRequest{})
if err != nil {
    if status.Code(err) == codes.Unimplemented {
        // Legacy plugin - use fallback values
        log.Info("Plugin does not implement GetPluginInfo")
        return &PluginMetadata{
            Name:    "unknown",
            Version: "unknown",
        }
    }
    return nil, fmt.Errorf("GetPluginInfo failed: %w", err)
}
```

---

## Performance Implications

### Error Handling Overhead

- Validation is O(1) - negligible performance impact
- Logging is async in zerolog - minimal blocking
- gRPC status code creation is O(1) - no overhead

### GetPluginInfo Performance

- **Standard Conformance**: 100ms threshold
- **Advanced Conformance**: 50ms threshold
- Should not make external API calls (per spec assumptions)
- Local metadata retrieval only
