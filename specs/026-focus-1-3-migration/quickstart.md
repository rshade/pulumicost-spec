# Quickstart: FOCUS 1.3 Migration

**Branch**: `026-focus-1-3-migration` | **Date**: 2025-12-23

## Overview

This guide covers implementing FOCUS 1.3 features in your PulumiCost plugin.

## New Capabilities

### 1. Split Cost Allocation

Track how shared resource costs are allocated to specific workloads:

```go
record := pluginsdk.NewFocusRecordBuilder().
    // ... existing FOCUS 1.2 fields ...
    WithAllocation(
        "proportional-cpu",           // AllocatedMethodId
        "CPU-weighted proportional",  // AllocatedMethodDetails
    ).
    WithAllocatedResource(
        "pod-frontend-abc123",        // AllocatedResourceId
        "frontend-service",           // AllocatedResourceName
    ).
    WithAllocatedTags(map[string]string{
        "team":        "platform",
        "environment": "production",
    }).
    Build()
```

### 2. Service/Host Provider Distinction

Clearly identify billing relationships in multi-vendor scenarios:

```go
record := pluginsdk.NewFocusRecordBuilder().
    // For an Azure Marketplace ISV product:
    WithServiceProvider("Datadog Inc").      // Who sells the service
    WithHostProvider("Microsoft Azure").     // Where it runs
    // ... other fields ...
    Build()
```

### 3. Contract Commitment Linking

Link cost records to contract commitment data:

```go
record := pluginsdk.NewFocusRecordBuilder().
    // ... existing fields ...
    WithContractApplied("commitment-ri-123456").  // Links to ContractCommitment
    Build()
```

### 4. Contract Commitment Dataset

Create standalone contract commitment records:

```go
commitment := pluginsdk.NewContractCommitmentBuilder().
    WithIdentity(
        "commitment-ri-123456",     // ContractCommitmentId
        "contract-enterprise-2024", // ContractId
    ).
    WithCategory(pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND).
    WithType("3-Year Reserved Instance").
    WithCommitmentPeriod(startTime, endTime).
    WithContractPeriod(contractStart, contractEnd).
    WithFinancials(
        120000.00,  // ContractCommitmentCost ($120k commitment)
        0,          // ContractCommitmentQuantity (not applicable for SPEND)
        "",         // ContractCommitmentUnit (not applicable for SPEND)
        "USD",      // BillingCurrency
    ).
    Build()
```

## Migration from FOCUS 1.2

### Backward Compatibility

Your existing FOCUS 1.2 code continues to work unchanged:

```go
// This still works - all new fields are optional
record := pluginsdk.NewFocusRecordBuilder().
    WithIdentity("AWS", "123456789012", "Production").
    WithBillingPeriod(start, end, "USD").
    WithChargePeriod(chargeStart, chargeEnd).
    WithService(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE, "Amazon EC2").
    WithChargeDetails(pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
                      pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD).
    WithChargeClassification(pbc.FocusChargeClass_FOCUS_CHARGE_CLASS_REGULAR,
                             "EC2 m5.large usage",
                             pbc.FocusChargeFrequency_FOCUS_CHARGE_FREQUENCY_USAGE_BASED).
    WithFinancials(100.00, 100.00, 95.00, "USD", "INV-2024-001").
    WithContractedCost(95.00).
    Build()
```

### Deprecated Fields

The following fields are deprecated but still work in FOCUS 1.3:

| Deprecated | Use Instead | Warning |
|------------|-------------|---------|
| `WithIdentity()` sets `provider_name` | `WithServiceProvider()` | Deprecation warning logged |
| `WithPublisher()` | `WithHostProvider()` | Deprecation warning logged |

### Deprecation Handling

When both old and new fields are set:

```go
record := pluginsdk.NewFocusRecordBuilder().
    WithIdentity("AWS", "123", "Prod").            // Sets provider_name
    WithServiceProvider("Amazon Web Services").    // Sets service_provider_name
    Build()
// Warning logged: "provider_name is deprecated, using service_provider_name"
// service_provider_name value takes precedence
```

## Validation Rules

### Allocation Field Dependency

```go
// VALID: Both method and resource set
record.WithAllocation("method-id", "details").
       WithAllocatedResource("resource-id", "resource-name")

// INVALID: Method without resource (validation error)
record.WithAllocation("method-id", "details")
// Error: "allocated_method_id requires allocated_resource_id"
```

### Contract Reference (No Validation)

```go
// Contract references are not validated against ContractCommitment dataset
record.WithContractApplied("any-commitment-id")  // Always passes
```

## Testing

### Conformance Tests

```go
// Run FOCUS 1.3 conformance tests
func TestFocus13Conformance(t *testing.T) {
    plugin := &MyPlugin{}
    result := plugintesting.RunFocus13ConformanceTests(t, plugin)
    if result.FailedTests > 0 {
        t.Errorf("FOCUS 1.3 conformance failed: %s", result.Summary)
    }
}
```

### Backward Compatibility Tests

```go
// Verify FOCUS 1.2 records still work
func TestFocus12BackwardCompat(t *testing.T) {
    // Build FOCUS 1.2-only record
    record, err := pluginsdk.NewFocusRecordBuilder().
        // Only FOCUS 1.2 fields...
        Build()

    require.NoError(t, err, "FOCUS 1.2 record should validate")
}
```

## Performance

Builder operations for new fields maintain existing performance:

- Target: <100 nanoseconds per operation
- Memory: 0 allocations for validation
- Benchmark: `go test -bench=BenchmarkFocus13 -benchmem ./sdk/go/testing/`

## Next Steps

1. Update your plugin to use new provider columns (`WithServiceProvider`, `WithHostProvider`)
2. Add allocation data for shared resources if applicable
3. Implement ContractCommitment support if you have contract data
4. Run conformance tests to verify FOCUS 1.3 compliance
