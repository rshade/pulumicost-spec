# Research: GetRecommendations RPC

**Feature**: 013-recommendations-rpc
**Date**: 2025-12-04

## Research Questions

### 1. Proto File Organization

**Question**: Should recommendation messages be in `costsource.proto` or a new `recommendations.proto`?

**Decision**: Add to existing `costsource.proto`

**Rationale**:

- GetRecommendations is an RPC method on `CostSourceService`, keeping it with the service
- Existing pattern: all `CostSourceService` RPCs and messages are in `costsource.proto`
- Enum types (ErrorCategory, ErrorCode, MetricType) are already in `costsource.proto`
- Splitting would fragment the service definition

**Alternatives Considered**:

- Separate `recommendations.proto`: Would require additional import, fragments service definition
- Split messages only: Inconsistent with existing patterns in the repo

### 2. Enum Placement

**Question**: Where should `RecommendationCategory`, `RecommendationActionType`, and
`RecommendationPriority` enums be defined?

**Decision**: Define in `costsource.proto` alongside existing enums (ErrorCategory, MetricType)

**Rationale**:

- Existing enums (ErrorCategory, ErrorCode, MetricType, SLIStatus) are in `costsource.proto`
- No separate `enums.proto` file exists in the current structure
- `registry.proto` and `enums.proto` exist but are for different domains
- Consistency with existing patterns

**Alternatives Considered**:

- New `recommendations_enums.proto`: Unnecessary fragmentation
- `enums.proto` file: Exists but contains BillingMode, different domain

### 3. Optional Interface Pattern

**Question**: How to implement the optional `RecommendationsProvider` interface in PluginSDK?

**Decision**: Follow existing `SupportsProvider` pattern in `sdk/go/pluginsdk/sdk.go`

**Rationale**:

- `SupportsProvider` is an established pattern (line 40-43 of sdk.go)
- Type assertion pattern: `provider, ok := s.plugin.(RecommendationsProvider)`
- Returns empty list for plugins not implementing the interface
- Maintains backward compatibility

**Implementation Pattern**:

```go
// RecommendationsProvider is an optional interface that plugins can implement
// to provide cost optimization recommendations.
type RecommendationsProvider interface {
    GetRecommendations(ctx context.Context, req *pbc.GetRecommendationsRequest) (
        *pbc.GetRecommendationsResponse, error)
}
```

**Alternatives Considered**:

- Required interface: Would break backward compatibility
- Capability flags only: Less type-safe than interface assertion

### 4. Metrics Integration

**Question**: How to add Prometheus metrics for GetRecommendations?

**Decision**: Extend existing `PluginMetrics` structure in `metrics.go`

**Rationale**:

- Existing `MetricsUnaryServerInterceptor` automatically captures request count and duration
- Add recommendation-specific counter for items returned per response
- Follow existing metric naming: `finfocus_plugin_*`

**New Metrics**:

```go
// In addition to automatic request_total and request_duration:
finfocus_plugin_recommendations_returned_total  // Counter: total recommendations
finfocus_plugin_recommendations_per_response    // Histogram: recommendations per call
```

**Alternatives Considered**:

- No additional metrics: Would miss recommendation volume tracking
- Separate metrics package: Unnecessary complexity

### 5. Capability Declaration in Supports

**Question**: How should plugins declare recommendations capability in `Supports` metadata?

**Decision**: Extend `SupportsResponse` with optional capabilities field

**Rationale**:

- Spec FR-018 requires deterministic capability discovery via `Supports`
- Adding `map<string, bool> capabilities` to `SupportsResponse` is backward compatible
- Example: `{"recommendations": true}` in capabilities map

**Proto Addition**:

```protobuf
message SupportsResponse {
  bool supported = 1;
  string reason = 2;
  // New field for capability declaration
  map<string, bool> capabilities = 3;  // e.g., {"recommendations": true}
}
```

**Alternatives Considered**:

- Separate Capabilities RPC: Over-engineering for single capability
- Boolean field: Less extensible than map

### 6. Pagination Token Format

**Question**: What format should pagination tokens use?

**Decision**: Opaque base64-encoded string (plugin-defined internal format)

**Rationale**:

- Plugin controls pagination implementation
- Opaque tokens prevent client-side manipulation
- Base64 encoding is URL-safe and consistent
- Existing patterns in cloud provider APIs (AWS, GCP)

**Validation**:

- Empty token = first page
- Invalid token = gRPC InvalidArgument error
- Plugin responsible for token generation/validation

**Alternatives Considered**:

- Offset-based: Exposes internal implementation, less flexible
- Cursor with defined schema: Over-specifies plugin implementation

### 7. Oneof vs Separate Action Fields

**Question**: Should action details use `oneof` or separate optional fields?

**Decision**: Use `oneof action_detail` for type-safe action details

**Rationale**:

- Ensures exactly one action type per recommendation
- Type-safe in generated Go code
- Matches spec FR-011 requirement for distinct action types
- Standard protobuf pattern for variants

**Proto Pattern**:

```protobuf
message Recommendation {
  // ... other fields ...
  oneof action_detail {
    RightsizeAction rightsize = 5;
    TerminateAction terminate = 6;
    CommitmentAction commitment = 7;
    KubernetesAction kubernetes = 8;
    ModifyAction modify = 9;
  }
}
```

**Alternatives Considered**:

- Separate optional fields: Allows multiple actions (incorrect semantics)
- Generic `google.protobuf.Any`: Loses type safety

### 8. Confidence Score Representation

**Question**: How to represent optional confidence scores?

**Decision**: Use `optional double confidence_score` with validation in SDK

**Rationale**:

- Proto3 `optional` keyword clearly indicates optional field
- SDK validates 0.0-1.0 range
- Nil/missing = confidence not available
- Zero is a valid confidence score (low confidence)

**Validation**:

```go
func ValidateConfidenceScore(score *float64) error {
    if score != nil && (*score < 0.0 || *score > 1.0) {
        return errors.New("confidence_score must be between 0.0 and 1.0")
    }
    return nil
}
```

**Alternatives Considered**:

- Wrapper message: Over-complex for single field
- Sentinel value (-1): Less clear semantics

### 9. Currency Validation Integration

**Question**: How to integrate with existing `sdk/go/currency` package?

**Decision**: Reuse existing `currency.IsValid()` for all currency fields in recommendations

**Rationale**:

- `currency` package already provides zero-allocation ISO 4217 validation
- Consistent validation across all currency fields in the SDK
- Performance: <15 ns/op validation

**Usage Points**:

- `RecommendationImpact.currency` validation
- `RecommendationSummary.currency` validation
- Error: InvalidArgument for non-ISO 4217 codes

**Alternatives Considered**:

- Duplicate validation: Violates DRY principle
- No validation: Allows invalid data

### 10. Logging Integration

**Question**: How to integrate zerolog logging for GetRecommendations?

**Decision**: Follow existing `logging.go` patterns with recommendation-specific fields

**Rationale**:

- Existing logging patterns in `pluginsdk/logging.go`
- Use existing field constants where applicable
- Add recommendation-specific fields

**Log Fields**:

```go
const (
    FieldRecommendationCount = "recommendation_count"
    FieldFilterCategory      = "filter_category"
    FieldFilterActionType    = "filter_action_type"
    FieldPageSize            = "page_size"
    FieldTotalSavings        = "total_savings"
)
```

**Log Events**:

- Request received (with filter details)
- Response complete (with count and total savings)
- Validation errors (with specific field)

## Summary

All research questions resolved. Key decisions:

| Decision | Choice |
|----------|--------|
| Proto location | `costsource.proto` (single file) |
| Enum location | `costsource.proto` |
| Optional interface | Follow `SupportsProvider` pattern |
| Metrics | Extend existing interceptor |
| Capability declaration | `capabilities` map in `SupportsResponse` |
| Pagination | Opaque base64 tokens |
| Action details | `oneof action_detail` |
| Confidence score | `optional double` with SDK validation |
| Currency validation | Reuse `sdk/go/currency` |
| Logging | Extend existing patterns |
