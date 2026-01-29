# Research: Documentation Drift Audit

**Date**: 2026-01-29
**Branch**: 043-docs-drift-audit

## Executive Summary

Verification of the current repository state against Issues #347 and #348 findings. Some issues
have been partially resolved since the audit was conducted. Issue #348 adds additional scope for
godoc coverage and code example modernization.

## Current State Verification

### 1. Version Consistency ✅ RESOLVED

All version references in README.md now show v0.5.4:

| Line | Content                                    | Status |
| ---- | ------------------------------------------ | ------ |
| 1    | `# FinFocus Specification v0.5.4`          | ✅     |
| 17   | `FinFocus Specification v0.5.4 is a...`    | ✅     |
| 832  | `Current version: **v0.5.4** (production)` | ✅     |
| 858  | `**FinFocus Specification v0.5.4** - ...`  | ✅     |

**Decision**: No action needed for version consistency.

### 2. RPC Count Mismatch ❌ STILL PRESENT

**README Claims** (2 locations):

- Line 80: "CostSourceService with 8 RPC methods"
- Line 699: "The CostSourceService provides 8 RPC methods"

**Actual RPC Count** (from proto/finfocus/v1/costsource.proto):

CostSourceService (11 RPCs):

1. Name
2. Supports
3. GetActualCost
4. GetProjectedCost
5. GetPricingSpec
6. EstimateCost
7. GetRecommendations
8. DismissRecommendation
9. GetBudgets
10. GetPluginInfo
11. DryRun

ObservabilityService (3 RPCs):

1. HealthCheck
2. GetMetrics
3. GetServiceLevelIndicators

**Decision**: Update both lines 80 and 699 to state "11 RPC methods". Consider documenting
observability RPCs separately.

### 3. Example Count Mismatch ❌ STILL PRESENT

**README Claims** (2 locations):

- Line 65: "10 comprehensive pricing examples"
- Line 660: "8 comprehensive examples"

**Actual Example Count** (examples/specs/*.json):

1. aws-ec2-t3-micro.json
2. aws-lambda-not-implemented.json
3. aws-lambda-per-invocation.json
4. aws-s3-tiered-pricing.json
5. azure-sql-dtu.json
6. azure-vm-per-second.json
7. gcp-preemptible-spot.json
8. gcp-storage-standard.json
9. kubernetes-namespace-cpu.json

Total: **9 JSON files**

**Decision**: Update line 65 to "9 comprehensive pricing examples" and line 660 to "9 comprehensive
examples".

### 4. Testing README Function Signatures ❌ STILL PRESENT

**README Shows** (sdk/go/testing/README.md):

```go
// Line 104
result := RunStandardConformance(plugin)

// Lines 171-174
result := RunBasicConformance(plugin)
result := RunStandardConformance(plugin)
result := RunAdvancedConformance(plugin)
```

**Actual Function Signatures** (sdk/go/testing/conformance.go):

```go
func RunBasicConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)
func RunStandardConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)
func RunAdvancedConformance(impl pbc.CostSourceServiceServer) (*ConformanceResult, error)
```

**Issue**: README examples ignore the error return value. Code will not compile as shown.

**Decision**: Update all examples to show proper error handling:

```go
result, err := RunStandardConformance(plugin)
if err != nil {
    t.Fatalf("conformance tests failed to run: %v", err)
}
```

### 5. Mapping Package README ❌ MISSING

**Status**: `sdk/go/pluginsdk/mapping/` has no README.md

**Package Contents**:

- aws.go - AWS property extraction (SKU, region, ARN parsing)
- azure.go - Azure property extraction
- gcp.go - GCP property extraction
- common.go - Shared utilities
- keys.go - Property key constants
- doc.go - Package documentation (godoc only)
- mapping_test.go - Unit tests
- benchmark_test.go - Performance benchmarks

**Decision**: Create README.md documenting the package purpose and usage patterns.

### 6. SDK Helper Documentation ❌ PARTIAL

**Key Helpers in pluginsdk** (from sdk/go/pluginsdk/helpers.go):

Response Builders:

- `NewActualCostResponse(opts ...ActualCostResponseOption)`
- `NewResourceDescriptor(provider, resourceType string, opts ...)`
- `NewEstimateCostResponse(opts ...)`

Functional Options:

- `WithFallbackHint(hint pbc.FallbackHint)`
- `WithResults(results []*pbc.ActualCostResult)`
- `WithID(id string)`, `WithARN(arn string)`, `WithSKU(sku string)`, etc.

Validation:

- `ValidateActualCostResponse(resp)`
- `ValidateRecommendation(rec)`
- `ValidateResourceRecommendationInfo(res)`
- `ValidateRecommendationImpact(impact)`

**Current Documentation Status**:

- `sdk/go/CLAUDE.md` - Comprehensive documentation of FallbackHint ✅
- `sdk/go/pluginsdk/README.md` - Documents Serve() and environment variables
- Root `README.md` - No mention of helper functions

**Decision**: Add helper documentation to `sdk/go/pluginsdk/README.md` or create a section in root
README. FallbackHint is well-documented in CLAUDE.md but needs user-facing documentation.

## Scope Reduction

Based on verification, the following items are ALREADY RESOLVED:

1. ✅ Version consistency (all v0.5.4)

The following items STILL NEED WORK:

1. ❌ RPC count (2 locations: lines 80, 699)
2. ❌ Example count (2 locations: lines 65, 660)
3. ❌ Testing README function signatures (lines 104, 171-174)
4. ❌ Mapping package README (new file)
5. ❌ SDK helper user-facing documentation

## Alternatives Considered

### SDK Helper Documentation Location

| Option                         | Pros                              | Cons                           |
| ------------------------------ | --------------------------------- | ------------------------------ |
| Root README.md                 | High visibility                   | README already long            |
| sdk/go/pluginsdk/README.md     | Logical location                  | Requires users to navigate     |
| New HELPERS.md                 | Dedicated space                   | Another file to maintain       |
| Link to godoc                  | Auto-generated                    | Requires online access         |

**Decision**: Add to `sdk/go/pluginsdk/README.md` with a cross-reference from root README.

### 7. Root README Code Examples (from Issue #348) ❌ OUTDATED

**Lines 259-279** show manual struct construction:

```go
resp := &pbc.GetActualCostResponse{
    Records: []*pbc.FocusRecord{...},
}
```

**Should use SDK helper**:

```go
resp := pluginsdk.NewActualCostResponse(
    pluginsdk.WithRecords(records),
    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
)
```

**Decision**: Update examples to use modern SDK patterns.

### 8. Godoc Coverage (from Issue #348) ❌ INCOMPLETE

**pluginsdk package** - ~25 undocumented exported functions:

- Response builders: NewActualCostResponse, NewProjectedCostResponse, NewDryRunResponse, NewEstimateCostResponse
- Validation helpers: ValidateActualCostResponse, ValidateProjectedCostResponse, ValidateDryRunResponse
- Field mapping utilities: FocusFieldNames, NewFieldMapping, AllFieldsWithStatus, SetFieldStatus

**pricing package** - ~10 undocumented functions:

- ApplyGrowth, ValidateGrowthParams, CheckGrowthWarnings
- ValidatePricingSpec, ValidateBillingMode

**registry package** - Needs README.md (has 8 enum types with validators)

**Decision**: Add godoc comments to all exported functions. Create registry/ README.md.

### RPC Count Format

| Option                          | Pros                    | Cons                            |
| ------------------------------- | ----------------------- | ------------------------------- |
| "11 RPC methods"                | Simple fix              | Doesn't explain service split   |
| "11 CostSource + 3 Observability" | Complete picture      | Verbose                         |
| List all RPCs                   | Self-documenting        | Takes space                     |

**Decision**: Use "11 RPC methods" with reference to proto file for full list.

## Implementation Priority

Based on impact and effort:

| Priority | Issue                   | Impact   | Effort | Notes                            |
| -------- | ----------------------- | -------- | ------ | -------------------------------- |
| P1       | RPC count fix           | High     | Low    | 2 line changes                   |
| P1       | Example count fix       | High     | Low    | 2 line changes                   |
| P1       | Testing README sigs     | Critical | Medium | 4 code blocks to update          |
| P2       | Mapping README          | Medium   | Medium | New file ~100 lines              |
| P3       | Helper documentation    | Medium   | Medium | Add section to existing README   |
| P3       | Root README examples    | Medium   | Medium | Update ~2 code blocks            |
| P4       | Godoc pluginsdk         | Medium   | High   | ~25 functions need comments      |
| P4       | Godoc pricing           | Medium   | Medium | ~10 functions need comments      |
| P4       | Registry README         | Low      | Medium | New file ~80 lines               |
