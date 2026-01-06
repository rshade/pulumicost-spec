# Research: v0.4.14 SDK Polish Release

**Date**: 2025-01-04
**Feature**: v0.4.14 SDK Polish Release
**Purpose**: Document technical decisions and best practices for SDK implementation improvements

## Research Topics

### 1. Context-Based Timeout Handling in Go gRPC Clients

**Decision**: Use Go's `context.WithTimeout()` pattern for per-request timeouts, with `ClientConfig.Timeout` option for default override.

**Rationale**:

- Go's standard library uses context deadlines for cancellation and timeouts (idiomatic pattern)
- gRPC/Go SDK already supports context-based request handling
- Per-request configuration is more flexible than only supporting global client timeout
- Backward compatible: existing code without explicit timeouts continues to work with default 30s timeout

**Alternatives Considered**:

- **Option A**: Only support context-based timeouts (no ClientConfig.Timeout)
  - Rejected because: Some users prefer configuration over code changes for timeout defaults
- **Option B**: Separate timeout parameter on each client method
  - Rejected because: Duplicates Go's built-in context mechanism, less flexible
- **Option C**: Global timeout variable package-level (no client config)
  - Rejected because: Doesn't support per-client or per-request flexibility

**Implementation Pattern**:

```go
// Client accepts context for all RPC methods
func (c *Client) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    // Context deadline is automatically respected by gRPC
    return c.client.GetActualCost(ctx, req)
}

// Optional: ClientConfig.Timeout provides default context if nil
func (c *Client) GetActualCost(ctx context.Context, req *pb.GetActualCostRequest) (*pb.GetActualCostResponse, error) {
    if ctx == nil {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(context.Background(), c.config.Timeout)
        defer cancel()
    }
    return c.client.GetActualCost(ctx, req)
}
```

**References**:

- Go context package documentation: https://pkg.go.dev/context
- gRPC-Go context handling: https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-go.md

---

### 2. Context Validation Helper Design

**Decision**: Provide `ValidateContext()`, `ContextRemainingTime()`, and `ContextDeadline()` as standalone helper functions in new `context.go` file.

**Rationale**:

- Prevents nil context panics before they occur (defensive programming)
- Provides clear error messages for expired/cancelled contexts (better UX)
- ContextRemainingTime() is useful for logging and metrics (observability)
- Separation of concerns: validation logic independent of client methods

**Alternatives Considered**:

- **Option A**: Embed validation in each client method directly
  - Rejected because: Code duplication, harder to test, inconsistent error messages
- **Option B**: Require plugins to always provide valid contexts
  - Rejected because: Increases cognitive load, doesn't catch errors early
- **Option C**: Use panic-based validation (RequireContext only)
  - Rejected because: Panics are not user-friendly, hard to debug in production

**Implementation Pattern**:

```go
// ValidateContext checks that a context is usable for RPC calls
func ValidateContext(ctx context.Context) error {
    if ctx == nil {
        return errors.New("context cannot be nil")
    }
    if err := ctx.Err(); err != nil {
        return fmt.Errorf("context already cancelled or expired: %w", err)
    }
    return nil
}

// ContextRemainingTime returns time until deadline
func ContextRemainingTime(ctx context.Context) time.Duration {
    deadline, ok := ctx.Deadline()
    if !ok {
        return time.Duration(math.MaxInt64)
    }
    return time.Until(deadline)
}

// ContextDeadline returns the context deadline or zero time if none set
func ContextDeadline(ctx context.Context) (time.Time, bool) {
    return ctx.Deadline()
}
```

**References**:

- Go context error handling: https://pkg.go.dev/context#Deadline
- Defensive programming in Go: https://go.dev/doc/effective_go#errors

---

### 3. HealthChecker Interface Design

**Decision**: Define `HealthChecker` interface with `Check(ctx context.Context) error` method. Return `HTTPStatus 503 / gRPC Unavailable` on timeout/panic.

**Rationale**:

- Single method interface (Go idiom: "The bigger the interface, the weaker the abstraction")
- Context parameter allows timeout configuration for health checks
- Error return allows plugins to define unhealthy conditions
- HTTP 503 / gRPC Unavailable indicates temporary unavailability (retryable)
- Backward compatible: plugins implementing interface get custom logic, others get default "always healthy"

**Alternatives Considered**:

- **Option A**: Return boolean instead of error
  - Rejected because: Error provides diagnostic information, boolean doesn't explain what's wrong
- **Option B**: Multiple methods (Check(), IsHealthy(), GetStatus())
  - Rejected because: Unnecessarily complex, Check() is sufficient
- **Option C**: Require plugins to implement HealthStatus struct directly
  - Rejected because: Tightly couples to our HealthStatus format, less flexible

**Implementation Pattern**:

```go
// HealthChecker allows plugins to provide custom health check logic
type HealthChecker interface {
    Check(ctx context.Context) error
}

// HealthStatus provides detailed health information
type HealthStatus struct {
    Healthy     bool              `json:"healthy"`
    Message     string            `json:"message,omitempty"`
    Details     map[string]string `json:"details,omitempty"`
    LastChecked time.Time         `json:"last_checked"`
}

// SDK detects and uses HealthChecker automatically
type ServeConfig struct {
    Plugin interface{}
    // ... other fields
}

// In Serve():
if hc, ok := config.Plugin.(HealthChecker); ok {
    // Use custom health check
    err := hc.Check(ctx)
    if err != nil {
        return HTTP 503 / gRPC Unavailable
    }
} else {
    // Default: always healthy
    return HealthStatus{Healthy: true}
}
```

**References**:

- Go health check patterns: https://github.com/grpc-ecosystem/grpc-health-probe
- HTTP 503 semantics: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/503

---

### 4. ARN Format Detection and Validation

**Decision**: `DetectARNProvider()` uses prefix/pattern matching. `ValidateARNConsistency()` checks ARN format matches expected provider. Ambiguous formats return explicit error.

**Rationale**:

- Simple string prefix matching is fast and reliable for provider detection
- Explicit error on ambiguity prevents silent misattribution (fail-fast principle)
- Empty string for unrecognized formats (not error) allows flexible identifier formats
- Exported pattern constants enable testing and documentation

**Alternatives Considered**:

- **Option A**: Use regex for all pattern matching
  - Rejected because: Overkill for simple prefix detection, harder to maintain
- **Option B**: Guess provider on ambiguous formats (first match wins)
  - Rejected because: Silent misattribution is dangerous, violates fail-fast
- **Option C**: Require explicit provider parameter instead of detection
  - Rejected because: Increases boilerplate, detection is more convenient

**Implementation Pattern**:

```go
const (
    AWSARNPrefix     = "arn:aws:"
    AzureARNPrefix   = "/subscriptions/"
    GCPARNPrefix     = "//"
    KubernetesFormat = "{cluster}/{namespace}/"
)

// DetectARNProvider returns the cloud provider inferred from ARN format
func DetectARNProvider(arn string) string {
    switch {
    case strings.HasPrefix(arn, AWSARNPrefix):
        return "aws"
    case strings.HasPrefix(arn, AzureARNPrefix):
        return "azure"
    case strings.HasPrefix(arn, GCPARNPrefix):
        return "gcp"
    case strings.Contains(arn, "{") && strings.Contains(arn, "}/"):
        return "kubernetes"
    default:
        return "" // Unrecognized
    }
}

// ValidateARNConsistency checks if ARN format matches expected provider
func ValidateARNConsistency(arn, expectedProvider string) error {
    detected := DetectARNProvider(arn)
    if detected == "" {
        return fmt.Errorf("ARN format unrecognized")
    }
    if detected != expectedProvider {
        return fmt.Errorf("ARN format %q detected as %q but expected %q", arn, detected, expectedProvider)
    }
    return nil
}

// Handle ambiguous formats (e.g., string starts with multiple prefixes)
// Decision: Explicit error "ARN format ambiguous, could be multiple providers"
```

**References**:

- AWS ARN format: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
- Azure Resource ID format: https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/
- GCP Full Resource Names: https://cloud.google.com/apis/design/resource_names

---

### 5. GetPluginInfo Error Message Strategy

**Decision**: Replace internal error details with user-friendly messages while preserving technical logs server-side.

**Rationale**:

- Client-facing errors should be actionable (users can't fix internal SDK issues)
- Server-side technical logs enable debugging for SDK maintainers
- Separation of concerns: UX for users, diagnostics for developers

**Current → Mapped Examples**:
| Current Error | Proposed Error |
|----------------------------------------|-----------------------------------------------|
| "plugin returned nil response" | "unable to retrieve plugin metadata" |
| "plugin returned incomplete metadata" | "plugin metadata is incomplete" |
| "plugin returned invalid spec_version format" | "plugin reported an invalid specification version" |

**Alternatives Considered**:

- **Option A**: Include all technical details in client errors
  - Rejected because: Exposes implementation details, not actionable for users
- **Option B**: Log only, return generic errors to clients
  - Rejected because: Too generic, users can't distinguish different failure modes
- **Option C**: Stack trace in client error message
  - Rejected because: Violates principle, exposes internal paths

**Implementation Pattern**:

```go
// GetPluginInfo handler
func (s *Server) GetPluginInfo(ctx context.Context, req *pb.GetPluginInfoRequest) (*pb.GetPluginInfoResponse, error) {
    info, err := s.getPluginInfo()
    if err == nil {
        return info, nil
    }

    // Log technical details server-side
    logger.Error("GetPluginInfo failed", "error", err, "details", s.debugInfo(err))

    // Return user-friendly error to client
    switch {
    case errors.Is(err, errNilResponse):
        return nil, status.Error(codes.Internal, "unable to retrieve plugin metadata")
    case errors.Is(err, errIncompleteMetadata):
        return nil, status.Error(codes.InvalidArgument, "plugin metadata is incomplete")
    case errors.Is(err, errInvalidSpecVersion):
        return nil, status.Error(codes.InvalidArgument, "plugin reported an invalid specification version")
    default:
        return nil, status.Error(codes.Unknown, "failed to retrieve plugin metadata")
    }
}
```

**References**:

- gRPC error handling best practices: https://grpc.github.io/grpc/core/md_doc_statuscodes.html
- User-friendly error messages: https://blog.codinghorror.com/the-art-of-error-messages/

---

### 6. Fuzz Testing in Go

**Decision**: Use Go's native fuzz testing (`testing.F`) with diverse seed corpus for `ResourceDescriptor.ID` field.

**Rationale**:

- Go 1.18+ native fuzzing is built-in and well-integrated
- Seed corpus provides coverage of known edge cases (URLs, ARNs, Unicode, null bytes)
- Fuzz testing discovers unexpected panics and edge cases
- CI runs short fuzz tests on PRs, extended runs locally

**Seed Corpus Examples**:

```go
f.Add("urn:pulumi:prod::myapp::resource::name")
f.Add("https://example.com/resource?id=12345#anchor")
f.Add("arn:aws:ec2:us-east-1:123456789012:instance/i-abc")
f.Add("/subscriptions/sub-123/resourceGroups/rg/providers/...")
f.Add("//compute.googleapis.com/projects/proj/zones/z/instances/i")
f.Add("") // Empty string
f.Add(strings.Repeat("a", 10000)) // Very long string
f.Add("null\x00byte") // Embedded null
f.Add("emoji\U0001F600test") // Unicode
```

**Alternatives Considered**:

- **Option A**: External fuzzing framework (libFuzzer, AFL)
  - Rejected because: Go native fuzzing is sufficient and better integrated
- **Option B**: Property-based testing (quickcheck)
  - Rejected because: Fuzz testing better for arbitrary input, property testing for invariants
- **Option C**: Manual edge case tests only
  - Rejected because: Cannot discover unexpected input patterns fuzzing finds

**Implementation Pattern**:

```go
func FuzzResourceDescriptorID(f *testing.F) {
    // Seed corpus
    f.Add("urn:pulumi:prod::myapp::resource::name")
    f.Add("arn:aws:ec2:us-east-1:123:instance/i-abc")
    // ... more seed cases

    f.Fuzz(func(t *testing.T, id string) {
        // Should never panic
        desc := pluginsdk.NewResourceDescriptor().WithID(id)

        // ID should round-trip correctly
        if desc.ID != id {
            t.Errorf("ID mismatch: got %q, want %q", desc.ID, id)
        }
    })
}

// Running locally: go test -fuzz=FuzzResourceDescriptorID -fuzztime=60s
// CI: go test -fuzz=FuzzResourceDescriptorID -fuzztime=10s
```

**References**:

- Go fuzz testing tutorial: https://go.dev/doc/fuzz
- Fuzz testing best practices: https://go.dev/security/fuzz/

---

### 7. CI Benchmark Stability Configuration

**Decision**: Set benchmark alert threshold to 150%, disable fail-on-alert, enable comment-on-alert.

**Rationale**:

- GitHub Actions CI has inherent infrastructure variability (noisy neighbors, CPU throttling)
- 50% tolerance reduces spurious failures while still catching real regressions
- Alerts provide visibility without blocking legitimate code changes
- Comments notify maintainers of potential performance issues

**Configuration**:

```yaml
# .github/workflows/benchmarks.yml
alert-threshold: "150%" # 50% tolerance for CI variance
fail-on-alert: false # Don't block PRs
comment-on-alert: true # Post comment for visibility
```

**Alternatives Considered**:

- **Option A**: Keep 110% threshold (10% tolerance)
  - Rejected because: Too strict for shared CI infrastructure, causes spurious failures
- **Option B**: Disable benchmarks in CI entirely
  - Rejected because: Loses performance regression detection
- **Option C**: Use dedicated benchmarking infrastructure
  - Rejected because: Overkill for this project, adds maintenance cost

**Documentation Requirement**:

- Document expected CI variance in README
- Explain that benchmark alerts are informational, not blocking
- Encourage manual benchmark runs for performance validation

**References**:

- GitHub Actions performance variability: https://github.com/actions/runner-images/issues
- Benchmarking best practices: https://golang.org/pkg/testing/#hdr.Benchmark

---

### 8. Extreme Value Testing (IEEE 754)

**Decision**: Add test cases for `math.Inf(1)`, `math.Inf(-1)`, `math.NaN()` in cost validation. Reject with clear errors.

**Rationale**:

- IEEE 754 special values (infinity, NaN) are valid float64 but invalid for monetary values
- Silent acceptance leads to calculation errors downstream
- Clear error messages prevent confusion
- Test max/min float64 to ensure valid extremes work correctly

**Test Cases**:

```go
tests := []struct {
    name        string
    billedCost  float64
    expectError bool
}{
    {"positive infinity", math.Inf(1), true},
    {"negative infinity", math.Inf(-1), true},
    {"NaN", math.NaN(), true},
    {"max float64", math.MaxFloat64, false},
    {"min positive float64", math.SmallestNonzeroFloat64, false},
    {"zero cost", 0.0, false},
}
```

**Alternatives Considered**:

- **Option A**: Accept infinity/NaN as valid
  - Rejected because: Invalid for cost calculations, causes downstream errors
- **Option B**: Silently clip to max/min values
  - Rejected because: Hides errors from users, unexpected behavior
- **Option C**: Don't test extreme values
  - Rejected because: Misses edge cases that could cause panics or crashes

**Implementation Pattern**:

```go
func ValidateCost(cost float64) error {
    if math.IsInf(cost, 1) || math.IsInf(cost, -1) {
        return errors.New("cost cannot be infinity")
    }
    if math.IsNaN(cost) {
        return errors.New("cost cannot be NaN")
    }
    if cost < 0 {
        return errors.New("cost cannot be negative")
    }
    return nil
}
```

**References**:

- IEEE 754 floating-point standard: https://ieeexplore.ieee.org/document/4610935
- Go math package specials: https://pkg.go.dev/math#IsInf

---

### 9. Code Complexity Reduction Techniques

**Decision**: Extract `validateCORSConfig()` function to reduce `Serve()` cognitive complexity below 20. Add unit tests for extracted function.

**Rationale**:

- `Serve()` currently at cognitive complexity 21 (exceeds threshold)
- Extracting validation logic reduces complexity and improves testability
- Unit tests for `validateCORSConfig()` ensure correctness
- No functional changes to CORS behavior (backward compatible)

**Refactoring Pattern**:

```go
// Before (Serve() with inline validation):
func Serve(ctx context.Context, config ServeConfig) error {
    // ... setup code ...
    if config.CORS != nil {
        if config.CORS.Origins != nil {
            // validate origins ...
        }
        if config.CORS.Methods != nil {
            // validate methods ...
        }
        // ... more inline validation ...
        // Complexity accumulates ...
    }
    // ... rest of Serve() ...
}

// After (Extracted function):
func validateCORSConfig(cors *CORSConfig) error {
    if cors == nil {
        return nil
    }
    // Validation logic isolated here
    if cors.Origins != nil {
        // validate origins ...
    }
    return nil
}

// Serve() complexity reduced:
func Serve(ctx context.Context, config ServeConfig) error {
    if err := validateCORSConfig(config.CORS); err != nil {
        return err
    }
    // ... rest of Serve() ...
}
```

**Alternatives Considered**:

- **Option A**: Ignore complexity threshold
  - Rejected because: Complexity indicates maintenance burden, should be addressed
- **Option B**: Split Serve() into multiple functions
  - Rejected because: More invasive refactoring, higher risk
- **Option C**: Increase threshold to 25
  - Rejected because: Threshold exists for maintainability, should be respected

**Verification**:

```go
// Unit tests for extracted function
func TestValidateCORSConfig(t *testing.T) {
    tests := []struct {
        name    string
        cors    *CORSConfig
        wantErr bool
    }{
        {"valid config", &CORSConfig{...}, false},
        {"nil config", nil, false},
        {"invalid origins", &CORSConfig{...}, true},
    }
    // ... test table ...
}
```

**References**:

- Cognitive complexity definition: https://en.wikipedia.org/wiki/Cyclomatic_complexity
- Go refactoring patterns: https://github.com/golang/go/wiki/CodeReviewComments#refactoring

---

### 10. Backward Compatibility Patterns

**Decision**: Maintain backward compatibility for legacy plugins (no HealthChecker implementation → default "always healthy", GetPluginInfo returns Unimplemented).

**Rationale**:

- Existing plugins cannot be updated immediately
- Graceful degradation ensures ecosystem stability
- Documentation guides migration path
- Default behaviors are safe (always healthy, Unimplemented status)

**Patterns Applied**:

1. **Interface Detection** (HealthChecker):

```go
// Type assertion for optional interface
if hc, ok := config.Plugin.(HealthChecker); ok {
    // Use custom implementation
    err := hc.Check(ctx)
    // ...
} else {
    // Default behavior
    return HealthStatus{Healthy: true}
}
```

2. **Fallback for Legacy RPCs** (GetPluginInfo):

```go
// Implement GetPluginInfo RPC
func (s *Server) GetPluginInfo(ctx context.Context, req *pb.GetPluginInfoRequest) (*pb.GetPluginInfoResponse, error) {
    // Check if plugin implements GetPluginInfo
    if infoProvider, ok := s.plugin.(GetPluginInfoProvider); ok {
        return infoProvider.GetPluginInfo(ctx, req)
    }
    // Legacy plugin: return Unimplemented
    return nil, status.Error(codes.Unimplemented, "GetPluginInfo not implemented by this plugin")
}
```

3. **Documentation** (Migration Guide):

```markdown
## Migrating to GetPluginInfo

### Backward Compatibility

- Legacy plugins (not implementing GetPluginInfo) return Unimplemented status
- Clients should handle Unimplemented gracefully
- Migration is optional for existing plugins
```

**Alternatives Considered**:

- **Option A**: Require all plugins to implement new features
  - Rejected because: Breaking change, would require all plugins to update simultaneously
- **Option B**: Use semantic versioning to indicate breaking change
  - Rejected because: Unnecessary if backward compatibility can be maintained
- **Option C**: Provide shim/wrapper for legacy plugins
  - Rejected because: Adds complexity, interface detection is simpler

**References**:

- Go interface assertion patterns: https://go.dev/tour/methods-and-interfaces#type-assertions
- gRPC backward compatibility: https://grpc.github.io/grpc/core/md_doc_grpc-well-known-types.html

---

## Summary of Decisions

All research topics have been resolved with clear technical decisions. Key patterns applied:

1. **Go idiomatic patterns**: context-based timeouts, interface detection, simple helper functions
2. **User-focused design**: User-friendly error messages, backward compatibility, fail-fast on ambiguity
3. **Test-First approach**: TDD required, comprehensive test coverage (unit, integration, fuzz, conformance)
4. **Performance awareness**: CI benchmarking, timeout configuration, concurrency testing
5. **Maintainability**: Code complexity reduction, separation of concerns, defensive programming

**No implementation blockers**: All technical decisions are clear and can proceed to Phase 1 (design and contracts).
