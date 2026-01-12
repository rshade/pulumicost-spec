# Data Model: SDK Polish v0.4.15

**Feature**: SDK Polish v0.4.15
**Date**: 2026-01-10
**Status**: Minimal (SDK Enhancement)

## Overview

SDK Polish v0.4.15 is an enhancement feature that adds verification and testing for existing SDK functionality.
No new data models or entities are introduced.

## Entity Changes

### Summary

| Entity   | Change Type | Description                                                            |
| -------- | ----------- | ---------------------------------------------------------------------- |
| None     | N/A         | No new entities introduced                                             |
| Existing | Verified    | Timeout configuration, error messages, performance tests already exist |

---

## Existing Entities (Verified)

### ClientConfig

**Location**: `sdk/go/pluginsdk/client.go:82-103`

**Fields**:

- `BaseURL string` - Server's base URL
- `Protocol Protocol` - RPC protocol (Connect/gRPC/gRPC-Web)
- `HTTPClient *http.Client` - Custom HTTP client (caller-provided)
- `Timeout time.Duration` - Per-client default timeout for RPC calls
- `ConnectOptions []connect.ClientOption` - Additional connect options

**Validation Rules**:

- Zero timeout uses `DefaultClientTimeout` (30 seconds)
- Custom `HTTPClient` takes precedence over `Timeout` field
- Context deadlines override client-level timeouts

**Relationships**:

- `ClientConfig` → Creates → `Client` (1:1)
- `Client` → Wraps → `CostSourceServiceClient` (1:1)

---

### Server

**Location**: `sdk/go/pluginsdk/sdk.go:230-244`

**Fields**:

- `plugin Plugin` - Plugin implementation
- `registry RegistryLookup` - Plugin registry for validation
- `logger zerolog.Logger` - Logging instance
- `pluginInfo *PluginInfo` - Optional plugin metadata for GetPluginInfo

**Methods**:

- `GetPluginInfo(ctx, req) (*GetPluginInfoResponse, error)` - Returns plugin metadata

**Validation Rules**:

- Returns `Unimplemented` for legacy plugins (graceful degradation)
- Validates response metadata completeness (name, version, spec_version)
- Validates spec version format using `ValidateSpecVersion()`

**Error Messages**:

- Nil response: "unable to retrieve plugin metadata"
- Incomplete metadata: "plugin metadata is incomplete"
- Invalid spec version: "plugin reported an invalid specification version"

---

### ConformanceTest

**Location**: `sdk/go/testing/performance.go` and `sdk/go/testing/conformance_test.go`

**Fields**:

- `Name string` - Test case name
- `Description string` - Test case description
- `Category Category` - Test category (e.g., Performance)
- `MinLevel ConformanceLevel` - Minimum conformance level required
- `TestFunc func(*TestHarness) TestResult` - Test function

**Performance Test Specific**:

- `Baseline PerformanceBaseline` - Latency thresholds (Standard/Advanced)
- `Iterations int` - Number of test iterations (typically 10)
- `Latency Thresholds`:
  - GetPluginInfo Standard: 100ms
  - GetPluginInfo Advanced: 50ms

---

## State Transitions

### Client Timeout Behavior

**States**:

1. **Configured** - ClientConfig created with timeout value
2. **Active** - Client created and making requests
3. **Timed Out** - Request exceeded timeout or context deadline

**Transitions**:

```text
Configured → Active (NewClient called)
Active → Timed Out (timeout exceeded or context deadline)
Active → Active (request completes within timeout)
```

**Precedence Rules**:

1. Context deadline (if set) - highest precedence
2. Custom HTTPClient.Timeout (if HTTPClient provided)
3. ClientConfig.Timeout (if Timeout > 0)
4. DefaultClientTimeout (30 seconds) - fallback

---

### GetPluginInfo Error Handling

**States**:

1. **Called** - GetPluginInfo request received
2. **Plugin Info Provider Check** - Check if plugin implements PluginInfoProvider
3. **Static Metadata Check** - Check if PluginInfo configured in ServeConfig
4. **Validation** - Validate response (nil, incomplete, invalid spec_version)
5. **Response or Error** - Return validated response or error

**Transitions**:

```text
Called → Plugin Info Provider Check
Plugin Info Provider Check → [Branch]
  ├─ Implements → Delegate to plugin
  └─ Not Implements → Static Metadata Check

Static Metadata Check → [Branch]
  ├─ PluginInfo Set → Return configured info
  └─ Not Set → Return Unimplemented

Delegate to plugin → [Branch]
  ├─ Error → Log + Return "unable to retrieve plugin metadata"
  ├─ Nil Response → Return "unable to retrieve plugin metadata"
  ├─ Incomplete → Return "plugin metadata is incomplete"
  ├─ Invalid SpecVersion → Return "plugin reported an invalid specification version"
  └─ Valid → Return response
```

---

### Performance Conformance Test Execution

**States**:

1. **Initialized** - Test harness created
2. **Measuring** - Running iterations
3. **Completed** - All iterations done
4. **Evaluated** - Compared against baseline
5. **Passed or Failed** - Final test result

**Transitions**:

```text
Initialized → Measuring
Measuring → Completed (all iterations done)
Completed → Evaluated (compareToBaseline called)
Evaluated → [Branch]
  ├─ Passed (latency ≤ threshold)
  └─ Failed (latency > threshold)
```

**Latency Metrics**:

- MinLatency: Minimum observed across iterations
- AvgLatency: Average across iterations
- MaxLatency: Maximum observed across iterations

**Pass Criteria**:

- Standard Conformance: AvgLatency ≤ StandardLatency (100ms for GetPluginInfo)
- Advanced Conformance: AvgLatency ≤ AdvancedLatency (50ms for GetPluginInfo)

---

## Data Flow

### Timeout Configuration Flow

```text
1. User creates ClientConfig
   ↓
2. Calls WithTimeout(duration) [optional]
   ↓
3. Calls NewClient(cfg)
   ↓
4. NewClient checks if HTTPClient is nil
   ↓
5. If nil: creates http.Client{Timeout: cfg.Timeout or DefaultClientTimeout}
   ↓
6. Client makes RPC call with context
   ↓
7. If context has deadline: deadline takes precedence
   ↓
8. If client timeout exceeded: context.DeadlineExceeded
   ↓
9. wrapRPCError() returns: "RPC cancelled or timed out"
```

### GetPluginInfo Error Message Flow

```text
1. Client calls GetPluginInfo
   ↓
2. Server receives request
   ↓
3. Server checks PluginInfoProvider interface
   ↓
4. If implements: calls plugin.GetPluginInfo()
   ↓
5. Plugin returns response (or error)
   ↓
6. Server validates response:
   - Is response nil? → "unable to retrieve plugin metadata"
   - Are required fields empty? → "plugin metadata is incomplete"
   - Is spec_version invalid? → "plugin reported an invalid specification version"
   ↓
7. Server logs detailed error (for debugging)
   ↓
8. Server returns gRPC error with user-friendly message
```

### Performance Test Execution Flow

```text
1. Test harness creates client connection
   ↓
2. Performance test starts
   ↓
3. Measure latency for iteration 1
   ↓
4. Measure latency for iteration 2
   ↓
... [repeat for 10 iterations]
   ↓
5. Calculate min/avg/max latency
   ↓
6. Get baseline thresholds (100ms Standard, 50ms Advanced)
   ↓
7. Compare avg latency to threshold
   ↓
8. Return TestResult:
   - Success: avg ≤ threshold
   - Error: "latency Xms exceeds threshold Yms"
```

---

## Validation Rules

### Timeout Configuration Validation

**FR-001**: ClientConfig.Timeout field must exist ✅
**FR-002**: WithTimeout() method must exist ✅
**FR-003**: Context deadlines must take precedence ✅
**FR-004**: Zero timeout must use default (30 seconds) ✅
**FR-005**: wrapRPCError must identify context timeout errors ✅

### GetPluginInfo Error Message Validation

**FR-006**: Nil response → "unable to retrieve plugin metadata" ✅
**FR-007**: Incomplete metadata → "plugin metadata is incomplete" ✅
**FR-008**: Invalid spec_version → "plugin reported an invalid specification version" ✅

### Performance Conformance Validation

**FR-009**: GetPluginInfoPerformance conformance test must exist ✅
**FR-010**: Test must run 10 iterations and fail if any exceeds 100ms ✅
**FR-011**: Test must handle Unimplemented error gracefully ✅

---

## Notes

1. **No New Data Models**: This feature verifies existing functionality, no new entities introduced
2. **Backward Compatibility**: All changes maintain backward compatibility with existing SDK usage
3. **No Proto Changes**: No protobuf definition changes required (SDK-level only)
4. **Testing Focus**: Phase 1 artifacts focus on test design, not new data structures
