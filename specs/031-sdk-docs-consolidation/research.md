# Research: SDK Documentation Consolidation

**Date**: 2025-12-31
**Branch**: `031-sdk-docs-consolidation`
**Status**: Complete

## Overview

This document captures research findings for all documentation gaps identified in the 12 GitHub
issues. Each section addresses a specific NEEDS CLARIFICATION item from the plan.

## 1. Connect-go Migration Patterns

**Decision**: Document migration from pure gRPC to connect-go with dual-mode serving
**Rationale**: Existing implementation in sdk.go provides a clean pattern that enables web
support without breaking native gRPC clients
**Alternatives Considered**: Single-protocol enforcement (rejected: loses backward compatibility)

### Key Findings

**Server Migration Steps**:

1. Add connect-go dependencies to go.mod
2. Create `ConnectHandler` adapter type that wraps existing `Server`
3. Implement `serveConnect()` function with h2c + health + CORS
4. Add `WebConfig` struct to `ServeConfig`
5. Branch serving mode: `Web.Enabled=true` → connect, `Web.Enabled=false` → pure gRPC

**Client Migration Steps**:

1. Define Protocol enumeration (Connect, gRPC, gRPC-Web)
2. Create configuration-based client factory
3. Add convenience constructors per protocol
4. Implement simplified RPC methods wrapping connect client
5. Provide resource cleanup via `Close()` method

**Protocol Selection Matrix**:

| Protocol | Transport | Browser Support | Performance | Use Case |
|----------|-----------|-----------------|-------------|----------|
| gRPC | HTTP/2 only | No | Best (binary) | Server-to-server, native clients |
| Connect | HTTP/1.1+ | Yes (fetch API) | Good (JSON) | Web dashboards, REST clients |
| gRPC-Web | HTTP/1.1+ | Yes | Good (protobuf) | Web clients needing binary protocol |

**Testing Considerations**:

- Unit tests via `ConnectHandler` adapter
- Integration tests with actual HTTP client
- CORS preflight testing with explicit origin headers
- Backward compatibility tests with pure gRPC clients

## 2. CORS Deployment Scenarios

**Decision**: Document 5 distinct deployment scenarios with security guidelines
**Rationale**: CORS misconfiguration is a common source of production issues; clear guidance prevents security vulnerabilities
**Alternatives Considered**: Generic CORS guide (rejected: not actionable enough)

### 5 Deployment Scenarios

#### 1. Local Development

```go
cfg := pluginsdk.DefaultWebConfig().
    WithWebEnabled(true).
    WithAllowedOrigins([]string{"http://localhost:3000", "http://127.0.0.1:3000"}).
    WithHealthEndpoint(true)
```

- Multiple localhost aliases handled
- No credentials needed for local CORS
- Health endpoint useful for local testing

#### 2. Single-Origin Production

```go
cfg := pluginsdk.DefaultWebConfig().
    WithWebEnabled(true).
    WithAllowedOrigins([]string{"https://app.example.com"}).
    WithAllowCredentials(true)
```

- Specific origin prevents unauthorized access
- HTTPS required before sending credentials
- Most restrictive and secure option

#### 3. Multi-Origin (Trusted Partners)

```go
cfg := pluginsdk.DefaultWebConfig().
    WithWebEnabled(true).
    WithAllowedOrigins([]string{
        "https://app.example.com",
        "https://dashboard.example.com",
        "https://partner.trusted.com",
    }).
    WithAllowCredentials(true)
```

- Explicit whitelist for partner integrations
- If origin list > 10, consider API gateway pattern

#### 4. API Gateway Pattern

- Plugins behind gateway don't need CORS
- Gateway handles all CORS for consistency
- Centralized authentication/authorization
- Better observability and logging

#### 5. Multi-Tenant SaaS

- Each tenant gets own plugin server instance
- Complete isolation between tenant origins
- Per-tenant rate limiting possible

### Security Guidelines

**Wildcard Origin (`*`)** - When to Avoid:

- Cannot send credentials (browser blocks)
- No protection against CSRF-like attacks
- Legitimate only for public APIs with no sensitive data

**Credentials Handling**:

- HTTPS required (browser enforces)
- Credentials enabled ONLY for trusted origins
- Authorization headers preferred over cookies for cross-origin

### Debugging Approaches

1. **Browser DevTools**: Check Network tab for CORS headers
2. **curl Testing**: Simulate browser with Origin header
3. **Server-side Logging**: Log origin, method, allowed status
4. **Integration Tests**: Test preflight and actual requests

## 3. Rate Limiting Patterns

**Decision**: Document token bucket pattern using golang.org/x/time/rate
**Rationale**: Standard Go library provides efficient, well-tested rate limiting
**Alternatives Considered**: Manual implementation (rejected: error-prone)

### Token Bucket Implementation

```go
import "golang.org/x/time/rate"

// Create rate limiter: 100 requests per second, burst of 200
limiter := rate.NewLimiter(rate.Limit(100), 200)

func (p *MyPlugin) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
    // Check rate limit
    if !limiter.Allow() {
        return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
    }
    // Proceed with request...
}
```

### Cloud Provider Rate Limits (Reference)

| Provider | Service | Default Limit | Recommended Local Limit |
|----------|---------|---------------|------------------------|
| AWS | Cost Explorer | 5 req/sec | 3 req/sec |
| Azure | Cost Management | 100 req/5min | 15 req/min |
| GCP | Billing API | 1000 req/min | 800 req/min |
| Kubernetes | Metrics API | Varies | 50 req/sec |

### Backoff Strategies

**Exponential Backoff with Jitter**:

```go
func backoff(attempt int) time.Duration {
    base := 100 * time.Millisecond
    max := 30 * time.Second

    // Exponential: 100ms, 200ms, 400ms, 800ms, ...
    delay := base * time.Duration(1<<attempt)
    if delay > max {
        delay = max
    }

    // Add jitter (±25%)
    jitter := time.Duration(rand.Float64() * 0.5 * float64(delay))
    return delay - jitter/2 + jitter
}
```

### Proper gRPC Status Codes

| Situation | Status Code | Description |
|-----------|-------------|-------------|
| Local rate limit | `ResourceExhausted` | Plugin's internal limit reached |
| Upstream throttling | `Unavailable` | Backend API returned 429 |
| Retry recommended | `Unavailable` | Include Retry-After header |

## 4. Thread Safety Documentation

**Decision**: Document thread safety guarantees for all major SDK components
**Rationale**: Thread safety documentation prevents race conditions that are difficult to diagnose
**Alternatives Considered**: Runtime checks (rejected: performance overhead)

### Component Thread Safety Summary

| Component | Thread-Safe | Notes |
|-----------|-------------|-------|
| **Client** | ✅ YES | Safe for concurrent RPC calls |
| **Server** | ✅ YES | Assumes Plugin implementation is safe |
| **WebConfig** | ✅ YES | Read-only after construction |
| **PluginMetrics** | ✅ YES | Prometheus internal atomics |
| **ResourceMatcher** | ❌ NO | Configure before Serve(), then read-only |
| **FocusRecordBuilder** | ❌ NO | Single-threaded builder pattern |

### Client Thread Safety

The `Client` struct wraps `http.Client` which is explicitly designed for concurrent use.
All methods are stateless request/response operations.

**Safe Usage**:

```go
// Create once, use from multiple goroutines
client := pluginsdk.NewConnectClient("http://localhost:8080")

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        name, err := client.Name(ctx) // Safe concurrent call
        // ...
    }()
}
wg.Wait()
client.Close()
```

### Server Thread Safety

The `Server` struct delegates to the `Plugin` interface which must be implemented in a
thread-safe manner by plugin developers. The gRPC framework handles concurrent request dispatch.

### ResourceMatcher Thread Safety

**NOT thread-safe** - Must be configured before `Serve()` is called:

```go
// CORRECT: Configure during initialization
matcher := pluginsdk.NewResourceMatcher()
matcher.AddProvider("aws")
matcher.AddResourceType("aws:ec2/instance:Instance")
// ... complete all configuration ...

// After Serve() is called, matcher becomes effectively read-only
pluginsdk.Serve(ctx, config)
```

### FocusRecordBuilder Thread Safety

**NOT thread-safe** - Each goroutine should create its own builder:

```go
// CORRECT: One builder per goroutine
go func() {
    builder := pluginsdk.NewFocusRecordBuilder()
    builder.WithIdentity("AWS", "123456789012", "Production")
    // ...
    record, err := builder.Build()
}()
```

## 5. HTTP Client Ownership Patterns

**Decision**: Document explicit ownership semantics in NewClient()
**Rationale**: Go community idiom - "who creates it, closes it"
**Alternatives Considered**: Always close (rejected: breaks shared clients)

### Ownership Pattern

```go
// SDK-owned: SDK creates client, SDK closes it
client := pluginsdk.NewConnectClient("http://localhost:8080")
defer client.Close() // Closes the internal HTTP client

// User-owned: User creates client, user is responsible
httpClient := &http.Client{Timeout: 60 * time.Second}
client := pluginsdk.NewClient(pluginsdk.ClientConfig{
    BaseURL:    "http://localhost:8080",
    HTTPClient: httpClient, // User-provided
})
// client.Close() is a no-op - user must close httpClient
```

### Documentation Pattern

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
func NewClient(cfg ClientConfig) *Client
```

## 6. Inline Comment Improvements

### contractedCostTolerance Explanation

**Current** (focus_conformance.go:13-14):

```go
// contractedCostTolerance is the relative tolerance for ContractedCost validation.
// Allows for floating-point precision differences up to 0.01%.
const contractedCostTolerance = 0.0001
```

**Improved**:

```go
// contractedCostTolerance defines the relative tolerance for ContractedCost validation.
//
// The tolerance is set to 0.0001 (1 basis point, or 0.01%) to account for
// IEEE 754 floating-point precision limitations in cost calculations.
//
// Why 1 basis point?
//   - Floating-point arithmetic can introduce small errors (e.g., 0.1 + 0.2 ≠ 0.3)
//   - 1 basis point (0.01%) is the standard financial tolerance for rounding
//   - Example: $1,000,000 contracted cost allows $100 variance (0.01%)
//
// Usage in validation:
//   expected := contractedUnitPrice * pricingQuantity
//   diff := abs(contractedCost - expected)
//   valid := diff <= max(abs(contractedCost), abs(expected)) * contractedCostTolerance
//
// Reference: FOCUS 1.2 Section 3.20 (ContractedCost)
const contractedCostTolerance = 0.0001
```

### WithTags Shared-Map Semantics

**Current** (focus_builder.go:259-265):

```go
// WithTags sets multiple tags at once per FOCUS 1.2 Section 2.14.
func (b *FocusRecordBuilder) WithTags(tags map[string]string) *FocusRecordBuilder {
    for k, v := range tags {
        b.record.Tags[k] = v
    }
    return b
}
```

**Improved**:

```go
// WithTags merges the provided tags into the record's tag map per FOCUS 1.2 Section 2.14.
//
// Shared-Map Semantics (Zero-Allocation Pattern):
//   - Tags are copied into the builder's internal map, NOT assigned by reference.
//   - The input map can be safely modified after this call without affecting the record.
//   - For performance-critical code processing thousands of records, consider using
//     WithTag(key, value) to avoid map iteration overhead.
//
// Thread Safety: NOT thread-safe. Do not call from multiple goroutines.
func (b *FocusRecordBuilder) WithTags(tags map[string]string) *FocusRecordBuilder
```

### Correlation Pattern Comment Fix

**Current reference** (from Issue #206): Comment references `resource_id` but should reference `ResourceRecommendationInfo.id`

**Fix**: Update any correlation pattern comments to use the correct proto field name `ResourceRecommendationInfo.id`.

### Test Naming Clarification

**Issue #208**: Test names in `resource_id_test.go` should clarify round-trip semantics vs literal old-server interoperability.

**Pattern**:

```go
// TestResourceID_RoundTrip_PreservesValue tests that resource IDs survive
// encode/decode cycles without modification (serialization fidelity).
func TestResourceID_RoundTrip_PreservesValue(t *testing.T)

// TestResourceID_BackwardCompat_OlderServerOmitsID tests graceful handling
// when older servers return responses without the id field populated.
func TestResourceID_BackwardCompat_OlderServerOmitsID(t *testing.T)
```

## Conclusion

All NEEDS CLARIFICATION items have been resolved through research and codebase analysis.
The findings are ready for incorporation into the implementation plan.

**Next Steps**:

1. Generate data-model.md (N/A for documentation feature)
2. Generate quickstart.md (implementation guide)
3. Update agent context
4. Proceed to task generation via `/speckit.tasks`
