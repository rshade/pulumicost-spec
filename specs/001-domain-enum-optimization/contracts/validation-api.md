# Validation API Contract

**Feature**: Domain Enum Validation Performance Optimization
**Package**: `github.com/rshade/pulumicost-spec/sdk/go/registry`
**Date**: 2025-11-17

## API Contract

This document defines the public API contract for registry package enum validation. All functions maintain backward
compatibility with existing implementations.

## Enum Accessor Functions

### AllProviders

```go
func AllProviders() []Provider
```

**Description**: Returns all supported cloud provider enum values

**Returns**: Slice containing all 5 Provider enum constants

**Performance**: Zero allocation (returns reference to package-level variable)

**Contract**:

- MUST return all Provider enum constants
- MUST return values in consistent order
- MUST NOT allocate memory on repeated calls
- Return value MAY be mutated by caller (defensive copy not required)

**Example**:

```go
providers := registry.AllProviders()
// Returns: [aws, azure, gcp, kubernetes, custom]
```

### AllDiscoverySources

```go
func AllDiscoverySources() []DiscoverySource
```

**Description**: Returns all plugin discovery source enum values

**Returns**: Slice containing all 4 DiscoverySource enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for DiscoverySource type

### AllPluginStatuses

```go
func AllPluginStatuses() []PluginStatus
```

**Description**: Returns all plugin operational status enum values

**Returns**: Slice containing all 6 PluginStatus enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for PluginStatus type

### AllSecurityLevels

```go
func AllSecurityLevels() []SecurityLevel
```

**Description**: Returns all plugin security level enum values

**Returns**: Slice containing all 4 SecurityLevel enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for SecurityLevel type

### AllInstallationMethods

```go
func AllInstallationMethods() []InstallationMethod
```

**Description**: Returns all plugin installation method enum values

**Returns**: Slice containing all 4 InstallationMethod enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for InstallationMethod type

### AllPluginCapabilities

```go
func AllPluginCapabilities() []PluginCapability
```

**Description**: Returns all plugin capability enum values

**Returns**: Slice containing all 14 PluginCapability enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for PluginCapability type

### AllSystemPermissions

```go
func AllSystemPermissions() []SystemPermission
```

**Description**: Returns all system permission enum values

**Returns**: Slice containing all 9 SystemPermission enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for SystemPermission type

### AllAuthMethods

```go
func AllAuthMethods() []AuthMethod
```

**Description**: Returns all authentication method enum values

**Returns**: Slice containing all 6 AuthMethod enum constants

**Performance**: Zero allocation

**Contract**: Same as AllProviders for AuthMethod type

## Validation Functions

### IsValidProvider

```go
func IsValidProvider(p string) bool
```

**Description**: Validates whether a string represents a valid Provider enum value

**Parameters**:

- `p` (string): Provider name to validate

**Returns**: `true` if `p` matches a valid Provider enum constant, `false` otherwise

**Performance**: < 30 nanoseconds per call, zero allocation

**Behavior**:

- Case-sensitive comparison (only lowercase matches)
- Exact match required
- Empty string returns `false`
- Invalid values return `false`

**Examples**:

```go
IsValidProvider("aws")              // true
IsValidProvider("AWS")              // false (case mismatch)
IsValidProvider("invalid")          // false
IsValidProvider("")                 // false
```

**Contract**:

- MUST return `true` for all values in `AllProviders()`
- MUST return `false` for all other input values
- MUST perform case-sensitive comparison
- MUST NOT allocate memory
- MUST complete in < 100 nanoseconds

### IsValidDiscoverySource

```go
func IsValidDiscoverySource(source string) bool
```

**Description**: Validates whether a string represents a valid DiscoverySource enum value

**Parameters**:

- `source` (string): Discovery source to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for DiscoverySource values

**Examples**:

```go
IsValidDiscoverySource("filesystem")  // true
IsValidDiscoverySource("registry")    // true
IsValidDiscoverySource("invalid")     // false
```

### IsValidPluginStatus

```go
func IsValidPluginStatus(status string) bool
```

**Description**: Validates whether a string represents a valid PluginStatus enum value

**Parameters**:

- `status` (string): Plugin status to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for PluginStatus values

**Examples**:

```go
IsValidPluginStatus("active")         // true
IsValidPluginStatus("installed")      // true
IsValidPluginStatus("unknown")        // false
```

### IsValidSecurityLevel

```go
func IsValidSecurityLevel(level string) bool
```

**Description**: Validates whether a string represents a valid SecurityLevel enum value

**Parameters**:

- `level` (string): Security level to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for SecurityLevel values

**Examples**:

```go
IsValidSecurityLevel("verified")      // true
IsValidSecurityLevel("official")      // true
IsValidSecurityLevel("unknown")       // false
```

### IsValidInstallationMethod

```go
func IsValidInstallationMethod(method string) bool
```

**Description**: Validates whether a string represents a valid InstallationMethod enum value

**Parameters**:

- `method` (string): Installation method to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for InstallationMethod values

**Examples**:

```go
IsValidInstallationMethod("binary")   // true
IsValidInstallationMethod("container")// true
IsValidInstallationMethod("invalid")  // false
```

### IsValidPluginCapability

```go
func IsValidPluginCapability(capability string) bool
```

**Description**: Validates whether a string represents a valid PluginCapability enum value

**Parameters**:

- `capability` (string): Plugin capability to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for PluginCapability values

**Examples**:

```go
IsValidPluginCapability("cost_retrieval")  // true
IsValidPluginCapability("caching")         // true
IsValidPluginCapability("invalid")         // false
```

### IsValidSystemPermission

```go
func IsValidSystemPermission(permission string) bool
```

**Description**: Validates whether a string represents a valid SystemPermission enum value

**Parameters**:

- `permission` (string): System permission to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for SystemPermission values

**Examples**:

```go
IsValidSystemPermission("network_access") // true
IsValidSystemPermission("filesystem_read")// true
IsValidSystemPermission("invalid")        // false
```

### IsValidAuthMethod

```go
func IsValidAuthMethod(method string) bool
```

**Description**: Validates whether a string represents a valid AuthMethod enum value

**Parameters**:

- `method` (string): Authentication method to validate

**Returns**: Boolean indicating validity

**Contract**: Same as IsValidProvider for AuthMethod values

**Examples**:

```go
IsValidAuthMethod("api_key")          // true
IsValidAuthMethod("jwt")              // true
IsValidAuthMethod("invalid")          // false
```

## Performance Contract

### Validation Performance Requirements

All validation functions (`IsValidXxx`) MUST meet these performance targets:

| Enum Type          | Max Values | Target Time | Max Allocation |
| ------------------ | ---------- | ----------- | -------------- |
| Provider           | 5          | < 10 ns/op  | 0 allocs/op    |
| DiscoverySource    | 4          | < 10 ns/op  | 0 allocs/op    |
| PluginStatus       | 6          | < 15 ns/op  | 0 allocs/op    |
| SecurityLevel      | 4          | < 10 ns/op  | 0 allocs/op    |
| InstallationMethod | 4          | < 10 ns/op  | 0 allocs/op    |
| PluginCapability   | 14         | < 30 ns/op  | 0 allocs/op    |
| SystemPermission   | 9          | < 20 ns/op  | 0 allocs/op    |
| AuthMethod         | 6          | < 15 ns/op  | 0 allocs/op    |

**Overall target**: < 100 nanoseconds per operation for all enum types (as specified in success criteria SC-003)

### Accessor Performance Requirements

All accessor functions (`AllXxx`) MUST meet these performance targets:

- **Allocation**: 0 bytes per call
- **Execution time**: < 5 nanoseconds (reference return)

## Backward Compatibility Guarantees

### API Surface

**MUST NOT change**:

- Function signatures (names, parameters, return types)
- Return value content (same enum values returned)
- Validation behavior (same results for all inputs)

**MAY change**:

- Internal implementation (slice allocation strategy)
- Performance characteristics (must improve, not regress)

### Breaking Change Detection

Any PR implementing this optimization MUST verify:

1. All existing tests pass unchanged
2. No API signature changes (verified by go vet)
3. Validation behavior identical (property-based tests)
4. Performance improves (benchmark comparison)

## Testing Contract

### Unit Test Requirements

All validation functions MUST have tests covering:

- Valid values (all enum constants)
- Invalid values (non-existent strings)
- Edge cases (empty string, case mismatch)
- Boundary conditions (nil, special characters)

### Benchmark Test Requirements

All validation functions MUST have benchmarks measuring:

- Nanoseconds per operation (ns/op)
- Allocations per operation (allocs/op)
- Bytes allocated per operation (B/op)

### Performance Regression Tests

CI MUST fail if:

- Any validation function allocates memory (> 0 allocs/op)
- Any validation exceeds 100 ns/op
- Performance degrades from baseline (> 10% slower)
