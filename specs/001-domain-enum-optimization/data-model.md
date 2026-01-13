# Data Model: Domain Enum Validation Performance Optimization

**Feature**: Domain Enum Validation Performance Optimization
**Date**: 2025-11-17
**Branch**: `001-domain-enum-optimization`

## Overview

This document describes the data structures and validation models for optimizing registry package enum validation. The
optimization involves converting function-returned slices to package-level variables while maintaining existing API
surfaces and behavior.

## Enum Types (Existing - No Schema Changes)

### Provider Enum

**Purpose**: Represents supported cloud providers in the FinFocus ecosystem

**Values**:

- `aws` - Amazon Web Services
- `azure` - Microsoft Azure
- `gcp` - Google Cloud Platform
- `kubernetes` - Kubernetes clusters
- `custom` - Custom provider implementations

**Validation**: Case-sensitive exact match against valid values

**Size**: 5 values (smallest enum in registry package)

### DiscoverySource Enum

**Purpose**: Represents plugin discovery mechanisms for locating cost source plugins

**Values**:

- `filesystem` - Local filesystem-based discovery
- `registry` - Remote registry-based discovery
- `url` - URL-based plugin discovery
- `git` - Git repository-based discovery

**Validation**: Case-sensitive exact match against valid values

**Size**: 4 values (tied for smallest)

### PluginStatus Enum

**Purpose**: Represents operational state of installed plugins

**Values**:

- `available` - Plugin available for installation
- `installed` - Plugin installed but not active
- `active` - Plugin installed and running
- `inactive` - Plugin installed but stopped
- `error` - Plugin in error state
- `updating` - Plugin currently being updated

**Validation**: Case-sensitive exact match against valid values

**Size**: 6 values

### SecurityLevel Enum

**Purpose**: Represents trust level for plugin security verification

**Values**:

- `untrusted` - Unverified plugin requiring explicit approval
- `community` - Community-verified plugin
- `verified` - Officially verified plugin
- `official` - Official FinFocus plugin

**Validation**: Case-sensitive exact match against valid values

**Size**: 4 values (tied for smallest)

### InstallationMethod Enum

**Purpose**: Represents deployment mechanisms for plugins

**Values**:

- `binary` - Direct binary installation
- `container` - Container image deployment
- `script` - Script-based installation
- `package` - System package manager installation

**Validation**: Case-sensitive exact match against valid values

**Size**: 4 values (tied for smallest)

### PluginCapability Enum

**Purpose**: Represents feature capabilities supported by plugins

**Values**:

- `cost_retrieval` - Cost data retrieval capability
- `cost_projection` - Cost projection capability
- `pricing_specs` - Pricing specification capability
- `historical_data` - Historical data support
- `real_time_data` - Real-time data support
- `batch_processing` - Batch processing support
- `rate_limiting` - Rate limiting support
- `caching` - Caching support
- `encryption` - Encryption support
- `compression` - Compression support
- `filtering` - Filtering support
- `aggregation` - Aggregation support
- `multi_tenancy` - Multi-tenancy support
- `audit_logging` - Audit logging support

**Validation**: Case-sensitive exact match against valid values

**Size**: 14 values (largest enum in registry package)

### SystemPermission Enum

**Purpose**: Represents system-level permissions required by plugins

**Values**:

- `network_access` - Outbound network connection permission
- `filesystem_read` - Filesystem read permission
- `filesystem_write` - Filesystem write permission
- `environment_read` - Environment variable read permission
- `process_spawn` - Process spawn permission
- `system_info` - System information access permission
- `temp_files` - Temporary file creation permission
- `config_read` - Configuration file read permission
- `metrics_collect` - Metrics collection permission

**Validation**: Case-sensitive exact match against valid values

**Size**: 9 values

### AuthMethod Enum

**Purpose**: Represents authentication methods supported for plugin access

**Values**:

- `none` - No authentication required
- `api_key` - API key authentication
- `jwt` - JWT token authentication
- `oauth2` - OAuth2 authentication
- `mtls` - Mutual TLS authentication
- `basic_auth` - Basic HTTP authentication

**Validation**: Case-sensitive exact match against valid values

**Size**: 6 values

## Validation Data Structures

### Current Implementation (Before Optimization)

```go
// Function-based slice generation (allocates on every call)
func AllProviders() []Provider {
    return []Provider{
        ProviderAWS,
        ProviderAzure,
        ProviderGCP,
        ProviderKubernetes,
        ProviderCustom,
    }
}

// Linear search through function-returned slice
func IsValidProvider(p string) bool {
    provider := Provider(p)
    for _, validProvider := range AllProviders() {
        if provider == validProvider {
            return true
        }
    }
    return false
}
```

**Memory characteristics:**

- **Allocation**: New slice allocated on each `AllProviders()` call
- **Size**: 5 elements × 8 bytes + 24-byte header = 64 bytes per call
- **GC impact**: Each validation creates garbage

### Optimized Implementation (After Optimization)

```go
// Package-level slice (allocated once at startup)
var allProviders = []Provider{
    ProviderAWS,
    ProviderAzure,
    ProviderGCP,
    ProviderKubernetes,
    ProviderCustom,
}

// Returns reference to package-level slice (zero allocation)
func AllProviders() []Provider {
    return allProviders
}

// Linear search through package-level slice (zero allocation)
func IsValidProvider(p string) bool {
    provider := Provider(p)
    for _, validProvider := range allProviders {
        if provider == validProvider {
            return true
        }
    }
    return false
}
```

**Memory characteristics:**

- **Allocation**: Zero per-call allocation (package variable allocated once)
- **Size**: 64 bytes total (fixed for program lifetime)
- **GC impact**: No per-call garbage generation

## Validation State Transitions

### Empty String Validation

**Input**: `""`
**Expected behavior**: Returns `false` (invalid)
**Rationale**: Empty strings do not match any enum constant

**Example**:

```go
IsValidProvider("") // false - empty string not in valid set
```

### Case Sensitivity Validation

**Input**: `"AWS"` (uppercase)
**Expected behavior**: Returns `false` (invalid)
**Rationale**: Validation is case-sensitive; only lowercase matches constants

**Example**:

```go
IsValidProvider("AWS")  // false - case mismatch
IsValidProvider("aws")  // true - exact match
```

### Invalid Value Validation

**Input**: `"invalid-provider"`
**Expected behavior**: Returns `false` (invalid)
**Rationale**: Value not present in valid enum set

**Example**:

```go
IsValidProvider("invalid-provider") // false - not in valid set
```

### Valid Value Validation

**Input**: `"aws"`
**Expected behavior**: Returns `true` (valid)
**Rationale**: Value matches enum constant exactly

**Example**:

```go
IsValidProvider("aws") // true - exact match with ProviderAWS
```

## Enum Size Distribution

**Registry package enum distribution:**

| Size Category      | Enum Count | Enum Types                                         | Validation Complexity        |
| ------------------ | ---------- | -------------------------------------------------- | ---------------------------- |
| Small (4 values)   | 3          | DiscoverySource, SecurityLevel, InstallationMethod | O(4) - 4 comparisons max     |
| Small (5-6 values) | 3          | Provider, PluginStatus, AuthMethod                 | O(5-6) - 5-6 comparisons max |
| Medium (9 values)  | 1          | SystemPermission                                   | O(9) - 9 comparisons max     |
| Medium (14 values) | 1          | PluginCapability                                   | O(14) - 14 comparisons max   |

**Validation performance characteristics:**

- **Best case**: 1 comparison (first element match) - ~5 ns
- **Average case**: N/2 comparisons - ~10-20 ns for 4-6 values, ~20-30 ns for 9-14 values
- **Worst case**: N comparisons (not found) - ~15-30 ns for 4-6 values, ~30-50 ns for 9-14 values

All within < 100 ns target for enums up to 50 values.

## Memory Footprint Analysis

**Total registry package validation memory:**

| Enum Type          | Elements | Memory (bytes)  |
| ------------------ | -------- | --------------- |
| Provider           | 5        | 64 (5×8 + 24)   |
| DiscoverySource    | 4        | 56 (4×8 + 24)   |
| PluginStatus       | 6        | 72 (6×8 + 24)   |
| SecurityLevel      | 4        | 56 (4×8 + 24)   |
| InstallationMethod | 4        | 56 (4×8 + 24)   |
| PluginCapability   | 14       | 136 (14×8 + 24) |
| SystemPermission   | 9        | 96 (9×8 + 24)   |
| AuthMethod         | 6        | 72 (6×8 + 24)   |
| **Total**          | **52**   | **608 bytes**   |

**Comparison with map-based approach:**

- **Slice total**: 608 bytes (fixed)
- **Map total**: ~3,500 bytes (52 entries × ~65 bytes/entry with hash overhead)
- **Memory savings**: 5.7x less memory with slice approach

## Validation Rules

### General Validation Contract

All `IsValidXxx()` functions follow this contract:

**Input**: String value to validate
**Output**: Boolean indicating validity

**Behavior**:

1. Convert input string to enum type (type assertion)
2. Iterate through package-level slice of valid values
3. Compare each valid value against input using equality
4. Return `true` if match found, `false` otherwise

**Guarantees**:

- Case-sensitive comparison (no normalization)
- Exact match required (no partial matching)
- Empty strings always return `false`
- Nil/uninitialized values treated as empty strings (return `false`)
- Zero allocation per validation call

### Edge Case Handling

**Empty String**:

- Input: `""`
- Behavior: Iterates through all values, finds no match
- Result: `false`

**Case Mismatch**:

- Input: `"AWS"` (expected: `"aws"`)
- Behavior: String comparison fails on case difference
- Result: `false`

**Invalid Characters**:

- Input: `"aws-invalid!@#"`
- Behavior: No match found in valid set
- Result: `false`

**Tight Loop Performance**:

- Scenario: 10,000 validations in tight loop
- Behavior: Zero allocation per iteration, cache-friendly access pattern
- Expected: ~50-300 microseconds total (5-30 ns per operation)

## Relationship to Pricing Package

**Shared enum**: Provider enum is used by both registry and pricing packages

**Consistency requirement**: Validation pattern must match pricing package approach

**Pricing package pattern** (`sdk/go/pricing/domain.go`):

```go
func getAllBillingModes() []BillingMode {
    return []BillingMode{/* 38 modes */}
}

func ValidBillingMode(s string) bool {
    mode := BillingMode(s)
    for _, validMode := range getAllBillingModes() {
        if mode == validMode {
            return true
        }
    }
    return false
}
```

**Optimization preserves consistency**: Both packages will use slice-based validation with package-level variables
after optimization.

## Data Invariants

**Invariant 1: Completeness**
All enum constants defined in the package MUST be included in corresponding package-level slice variable.

**Invariant 2: Uniqueness**
Each enum value MUST appear exactly once in its package-level slice variable.

**Invariant 3: Immutability**
Package-level slice variables MUST NOT be modified after initialization (though they return mutable references,
callers should not mutate).

**Invariant 4: Synchronization**
`AllXxx()` function return value MUST match package-level slice variable content.

**Invariant 5: Backward Compatibility**
Validation function signatures and behavior MUST remain unchanged (zero breaking changes).
