# Quickstart: SDK Documentation Consolidation

**Date**: 2025-12-31
**Branch**: `031-sdk-docs-consolidation`
**Status**: Ready for Implementation

## Prerequisites

- Access to the 12 GitHub issues referenced in spec.md
- Familiarity with Go documentation conventions
- Understanding of the pluginsdk package structure

## Implementation Order

The implementation is organized into 4 tiers based on scope and dependencies.

### Tier 1: Quick Wins (Inline Comments)

Start with inline code documentation changes that require minimal effort:

#### 1.1 HTTP Client Ownership Comment (Issue #240)

**File**: `sdk/go/pluginsdk/client.go`
**Location**: Above `NewClient()` function
**Action**: Add ownership semantics documentation

```go
// NewClient creates a new FinFocus client with the given configuration.
//
// HTTP Client Ownership:
//   - If HTTPClient is nil, an internal HTTP client is created and owned by
//     this Client. Call Close() when done to release connection pool resources.
//   - If HTTPClient is provided, the caller retains ownership. Close() is a
//     no-op; the caller is responsible for closing the HTTP client.
//
// Thread Safety: Client is safe for concurrent use from multiple goroutines.
func NewClient(cfg ClientConfig) *Client {
```

#### 1.2 Tolerance Constant Explanation (Issue #211)

**File**: `sdk/go/pluginsdk/focus_conformance.go`
**Location**: Above `contractedCostTolerance` constant
**Action**: Expand comment with rationale

See research.md Section 6 for the improved comment text.

#### 1.3 WithTags Semantics (Issue #207)

**File**: `sdk/go/pluginsdk/focus_builder.go`
**Location**: Above `WithTags()` method
**Action**: Document copy semantics

See research.md Section 6 for the improved comment text.

#### 1.4 Correlation Pattern Fix (Issue #206)

**File**: Varies (search for "resource_id" correlation comments)
**Action**: Update to reference `ResourceRecommendationInfo.id`

#### 1.5 Test Naming Clarification (Issue #208)

**File**: `sdk/go/testing/*_test.go` (where resource ID tests exist)
**Action**: Rename tests to clarify purpose (round-trip vs backward-compat)

### Tier 2: Godoc Examples

Create or update example files:

#### 2.1 Client.Close() Example (Issue #238)

**File**: `sdk/go/pluginsdk/example_test.go` (create if not exists)
**Action**: Add `ExampleClient_Close()` function

```go
package pluginsdk_test

import (
    "context"

    "github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

func ExampleClient_Close() {
    // Create a client with SDK-owned HTTP client
    client := pluginsdk.NewConnectClient("http://localhost:8080")

    // Use the client for requests
    ctx := context.Background()
    name, _ := client.Name(ctx)
    _ = name

    // Close releases connection pool resources
    // This is a no-op if the client was created with a user-provided HTTPClient
    client.Close()
}

func ExampleClient_Close_userProvided() {
    // When providing your own HTTP client, you manage its lifecycle
    httpClient := &http.Client{Timeout: 60 * time.Second}

    client := pluginsdk.NewClient(pluginsdk.ClientConfig{
        BaseURL:    "http://localhost:8080",
        HTTPClient: httpClient,
    })

    // Use the client...
    ctx := context.Background()
    name, _ := client.Name(ctx)
    _ = name

    // client.Close() is a no-op here - caller manages httpClient
    client.Close()

    // Caller is responsible for closing the HTTP client
    httpClient.CloseIdleConnections()
}
```

#### 2.2 Complete Import Statements (Issue #209)

**File**: `sdk/go/testing/README.md`
**Action**: Ensure all code blocks have complete imports including `pbc` alias

```go
import (
    "context"
    "testing"

    plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
    pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)
```

### Tier 3: README Sections (Medium Scope)

Add new sections to existing README files:

#### 3.1 Thread Safety Section (Issue #231)

**File**: `sdk/go/pluginsdk/README.md`
**Location**: New section "## Thread Safety"
**Content**: Component summary table from research.md Section 4

#### 3.2 Rate Limiting Section (Issue #233)

**File**: `sdk/go/pluginsdk/README.md`
**Location**: New section "## Rate Limiting"
**Content**: Token bucket pattern from research.md Section 3

### Tier 4: Comprehensive Guides (Largest Scope)

These require more substantial documentation:

#### 4.1 Migration Guide: gRPC to connect-go (Issue #235)

**File**: `sdk/go/pluginsdk/README.md`
**Location**: New section "## Migration Guide: gRPC to Connect"
**Content**: Server and client migration steps from research.md Section 1

#### 4.2 Performance Tuning Section (Issue #237)

**File**: `sdk/go/pluginsdk/README.md`
**Location**: New section "## Performance Tuning"
**Content**:

- Connection pool configuration
- Server timeouts for DoS protection
- Protocol performance trade-offs
- Benchmark reference values

#### 4.3 CORS Best Practices Section (Issue #236)

**File**: `sdk/go/pluginsdk/README.md`
**Location**: New section "## CORS Configuration"
**Content**: 5 deployment scenarios from research.md Section 2

## Validation Checklist

Before submitting:

- [ ] Run `make lint-markdown` - all markdown files pass
- [ ] Run `go build ./...` - all example code compiles
- [ ] Run `go test ./...` - all tests pass
- [ ] Verify no placeholder values in examples
- [ ] Verify cross-provider coverage (AWS, Azure, GCP)
- [ ] Link each change to its GitHub issue number

## File Summary

| File | Changes |
|------|---------|
| `sdk/go/pluginsdk/client.go` | Inline comment (Issue #240) |
| `sdk/go/pluginsdk/focus_conformance.go` | Inline comments (Issues #206, #211) |
| `sdk/go/pluginsdk/focus_builder.go` | Inline comment (Issue #207) |
| `sdk/go/pluginsdk/example_test.go` | New file with godoc examples (Issue #238) |
| `sdk/go/pluginsdk/README.md` | Major updates (Issues #231, #233, #235, #236, #237) |
| `sdk/go/testing/README.md` | Import fixes (Issue #209) |
| `sdk/go/testing/*_test.go` | Test naming (Issue #208) |

## Issue Mapping

| GitHub Issue | Tier | Description |
|--------------|------|-------------|
| #206 | 1 | Correlation pattern comment fix |
| #207 | 1 | WithTags shared-map semantics |
| #208 | 1 | Test naming clarification |
| #209 | 2 | Complete import statements |
| #211 | 1 | contractedCostTolerance explanation |
| #231 | 3 | Thread safety documentation |
| #233 | 3 | Rate limiting patterns |
| #235 | 4 | Migration guide (gRPC to connect) |
| #236 | 4 | CORS best practices |
| #237 | 4 | Performance tuning guide |
| #238 | 2 | Client.Close() godoc example |
| #240 | 1 | HTTP client ownership semantics |
