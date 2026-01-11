# GetPluginInfo Performance Conformance Test Contract

**Feature**: SDK Polish v0.4.15
**Date**: 2026-01-10
**Type**: Test Contract

## Overview

This contract documents the `Performance_GetPluginInfoLatency` conformance test that validates GetPluginInfo RPC performance.

---

## Test Definition

### Test Metadata

```go
{
    Name:        "Performance_GetPluginInfoLatency",
    Description: "Validates GetPluginInfo RPC latency within thresholds",
    Category:    CategoryPerformance,
    MinLevel:    ConformanceLevelStandard,
    TestFunc:    createGetPluginInfoLatencyTest(),
}
```

### Test Function

```go
func createGetPluginInfoLatencyTest() func(*TestHarness) TestResult {
    return func(harness *TestHarness) TestResult {
        baseline := GetBaseline(MethodGetPluginInfo)
        result := measureLatency(MethodGetPluginInfo, LatencyTestIterations, func() error {
            _, callErr := harness.Client().GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
            return callErr
        })
        compareToBaseline(result, baseline)
        return buildLatencyTestResult(MethodGetPluginInfo, result, baseline)
    }
}
```

---

## Performance Baseline

### Method: GetPluginInfo

```go
{
    Method:          MethodGetPluginInfo,
    StandardLatency: GetPluginInfoStandardLatencyMs * time.Millisecond,  // 100ms
    AdvancedLatency: GetPluginInfoAdvancedLatencyMs * time.Millisecond,  // 50ms
}
```

### Threshold Constants

```go
// GetPluginInfoStandardLatencyMs is the GetPluginInfo RPC standard latency threshold in milliseconds.
GetPluginInfoStandardLatencyMs = 100

// GetPluginInfoAdvancedLatencyMs is the GetPluginInfo RPC advanced latency threshold in milliseconds.
GetPluginInfoAdvancedLatencyMs = 50
```

---

## Test Execution

### Iteration Count

**FR-010**: The performance conformance test MUST run 10 iterations and fail if any exceeds 100ms.

```go
// LatencyTestIterations = 10
result := measureLatency(MethodGetPluginInfo, LatencyTestIterations, func() error {
    _, callErr := harness.Client().GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
    return callErr
})
```

### Latency Measurement

```go
func measureLatency(name string, iterations int, fn func() error) *PerformanceResult {
    result := &PerformanceResult{
        Method:     name,
        Iterations: iterations,
        MinLatency: time.Hour, // Start with high value
    }

    var totalDuration time.Duration
    for range iterations {
        start := time.Now()
        _ = fn()
        duration := time.Since(start)

        totalDuration += duration
        if duration < result.MinLatency {
            result.MinLatency = duration
        }
        if duration > result.MaxLatency {
            result.MaxLatency = duration
        }
    }

    if iterations > 0 {
        result.AvgLatency = totalDuration / time.Duration(iterations)
    }

    return result
}
```

### Metrics Collected

| Metric     | Description                                | Use                             |
| ---------- | ------------------------------------------ | ------------------------------- |
| MinLatency | Minimum observed latency across iterations | Identify best-case performance  |
| AvgLatency | Average latency across iterations          | **Used for pass/fail decision** |
| MaxLatency | Maximum observed latency across iterations | Identify worst-case performance |
| Iterations | Number of test iterations                  | Should be 10 (FR-010)           |

---

## Pass/Fail Criteria

### Standard Conformance (FR-010)

**Condition**: Average latency ≤ 100ms

```go
func compareToBaseline(result *PerformanceResult, baseline *PerformanceBaseline) {
    if baseline.StandardLatency > 0 {
        result.PassedStandard = result.AvgLatency <= baseline.StandardLatency

        if baseline.StandardLatency > 0 {
            variance := float64(result.AvgLatency-baseline.StandardLatency) /
                float64(baseline.StandardLatency) * PercentageCalculationFactor
            result.VariancePercent = variance
        }
    } else {
        result.PassedStandard = true // No Standard requirement for this method
    }
}
```

**Success Example**:

```
Method: GetPluginInfo
Iterations: 10
MinLatency: 45ms
AvgLatency: 78ms  ← Pass (≤ 100ms)
MaxLatency: 120ms
Status: PASSED
```

**Failure Example**:

```
Method: GetPluginInfo
Iterations: 10
MinLatency: 85ms
AvgLatency: 112ms  ← Fail (> 100ms)
MaxLatency: 150ms
Status: FAILED
Error: "latency 112ms exceeds threshold 100ms"
```

### Advanced Conformance

**Condition**: Average latency ≤ 50ms

```go
if baseline.AdvancedLatency > 0 {
    result.PassedAdvanced = result.AvgLatency <= baseline.AdvancedLatency
} else {
    result.PassedAdvanced = true // No Advanced requirement for this method
}
```

**Success Example**:

```
Method: GetPluginInfo
Iterations: 10
MinLatency: 20ms
AvgLatency: 38ms  ← Pass (≤ 50ms)
MaxLatency: 55ms
Status: PASSED (Advanced)
```

**Failure Example**:

```
Method: GetPluginInfo
Iterations: 10
MinLatency: 40ms
AvgLatency: 62ms  ← Fail (> 50ms)
MaxLatency: 85ms
Status: FAILED (Advanced)
Error: "latency 62ms exceeds threshold 50ms"
```

---

## Test Result Structure

### Success Result

```go
TestResult{
    Method:   "GetPluginInfo",
    Category: CategoryPerformance,
    Success:  true,
    Duration: 78ms,  // AvgLatency
    Details:  "Avg: 78ms (threshold: 100ms)",
}
```

### Failure Result

```go
TestResult{
    Method:   "GetPluginInfo",
    Category: CategoryPerformance,
    Success:  false,
    Error: fmt.Errorf("latency %.2fms exceeds threshold %.2fms",
        float64(perfResult.AvgLatency.Milliseconds()),
        float64(baseline.StandardLatency.Milliseconds())),
    Duration: 112ms,  // AvgLatency
    Details:  "Avg: 112ms, Min: 85ms, Max: 150ms",
}
```

---

## Legacy Plugin Handling (FR-011)

### Unimplemented Error

**FR-011**: The performance conformance test MUST handle legacy plugins that return `Unimplemented` gracefully.

**Expected Behavior**:

```go
// Test should NOT fail if plugin returns Unimplemented
_, err := harness.Client().GetPluginInfo(context.Background(), &pbc.GetPluginInfoRequest{})
if status.Code(err) == codes.Unimplemented {
    // Legacy plugin - handle gracefully
    return TestResult{
        Method:   "GetPluginInfo",
        Category: CategoryPerformance,
        Success:  true,  // Legacy plugin is acceptable
        Duration: 0,
        Details:  "Legacy plugin (GetPluginInfo not implemented)",
    }
}
```

**Test Flow**:

```
1. Call GetPluginInfo()
   ↓
2. Check error code
   ↓
3. [Branch]
   ├─ Unimplemented → Return success (legacy plugin)
   ├─ Other error → Return failure
   └─ No error → Measure latency
```

---

## Test Scenarios

### Scenario 1: Fast Plugin (Passes Standard)

**Setup**: Plugin responds in ~50ms

**Expected Result**:

```
Min: 45ms
Avg: 52ms  ← Pass (≤ 100ms)
Max: 60ms
Status: PASSED
```

### Scenario 2: Slow Plugin (Fails Standard)

**Setup**: Plugin responds in ~150ms

**Expected Result**:

```
Min: 140ms
Avg: 152ms  ← Fail (> 100ms)
Max: 165ms
Status: FAILED
Error: "latency 152ms exceeds threshold 100ms"
```

### Scenario 3: Variable Performance (Avg Passes, Max Exceeds)

**Setup**: Plugin response varies (30-130ms)

**Expected Result**:

```
Min: 30ms
Avg: 78ms  ← Pass (≤ 100ms)
Max: 130ms  ← Exceeds threshold, but not used for decision
Status: PASSED
Note: Only AvgLatency is used for pass/fail decision
```

### Scenario 4: Legacy Plugin (Unimplemented)

**Setup**: Plugin does not implement GetPluginInfo

**Expected Result**:

```
Error: Unimplemented
Status: PASSED (legacy plugin handling)
Details: "Legacy plugin (GetPluginInfo not implemented)"
Note: Test should not fail on Unimplemented error
```

---

## Integration with Conformance Suite

### Test Registration

```go
suite := plugintesting.NewPluginConformanceSuite()
suite.AddTest(plugintesting.ConformanceTest{
    Name:        "Performance_GetPluginInfoLatency",
    Description: "Validates GetPluginInfo RPC latency within thresholds",
    Category:    CategoryPerformance,
    MinLevel:    ConformanceLevelStandard,
    TestFunc:    createGetPluginInfoLatencyTest(),
})
```

### Running the Test

```go
result := suite.RunTests(t, impl)
for _, testResult := range result {
    if testResult.Method == "GetPluginInfo" && testResult.Category == CategoryPerformance {
        if !testResult.Success {
            t.Logf("❌ GetPluginInfo performance test failed: %v", testResult.Error)
        }
    }
}
```

### Conformance Level Requirement

| Conformance Level | MinLevel Requirement                 | Pass Condition     |
| ----------------- | ------------------------------------ | ------------------ |
| Basic             | N/A (No performance requirement)     | N/A                |
| Standard          | `MinLevel: ConformanceLevelStandard` | AvgLatency ≤ 100ms |
| Advanced          | `MinLevel: ConformanceLevelStandard` | AvgLatency ≤ 50ms  |

---

## Variance Calculation (SC-003)

### Purpose

Advanced conformance includes variance tracking to ensure performance consistency across multiple runs.

### Calculation

```go
variance := float64(result.AvgLatency - baseline.StandardLatency) /
    float64(baseline.StandardLatency) * PercentageCalculationFactor
result.VariancePercent = variance
```

### Example

```
Baseline: 100ms
AvgLatency: 95ms
Variance: (95 - 100) / 100 * 100 = -5%

Baseline: 100ms
AvgLatency: 108ms
Variance: (108 - 100) / 100 * 100 = 8%
Status: PASSED (variance ≤ 10%)

Baseline: 100ms
AvgLatency: 115ms
Variance: (115 - 100) / 100 * 100 = 15%
Status: FAILED (variance > 10%)
```

### Variance Threshold

```go
const MaxVariancePercent = 10.0  // SC-003 requirement
```

---

## Testing Requirements

### Unit Tests

- ❌ **Missing**: `createGetPluginInfoLatencyTest` unit test (verify logic)
- ❌ **Missing**: Legacy plugin Unimplemented handling test

### Integration Tests

- ❌ **Missing**: Fast mock plugin test (verify 100ms pass)
- ❌ **Missing**: Slow mock plugin test (verify 100ms fail)
- ❌ **Missing**: Variable performance test (verify avg vs max behavior)
- ❌ **Missing**: Legacy plugin test (verify Unimplemented handling)

---

## Performance Implications

### Test Overhead

- **Measurement overhead**: `time.Since()` is ~10-20ns per call (negligible)
- **Total test time**: 10 iterations × average latency (e.g., 10 × 50ms = 500ms)
- **Result calculation**: O(1) for min/avg/max computation

### Resource Usage

- **Memory**: Minimal (stores 10 duration values)
- **CPU**: Negligible (just time measurements)
- **Network**: Depends on plugin implementation (should be local mock for testing)

---

## Best Practices

### Running Performance Tests

1. **Use Local Mocks**: Avoid network latency affecting results
2. **Multiple Runs**: Run tests 3-5 times to account for system variance
3. **Consistent Environment**: Run on same hardware/machine for comparisons
4. **Profile Slow Tests**: If test fails, profile plugin implementation for bottlenecks

### Interpreting Results

- **AvgLatency matters most**: Used for pass/fail decision
- **MinLatency**: Best-case performance (not for decision)
- **MaxLatency**: Worst-case performance (identify outliers, not for decision)
- **Variance**: Consistency metric (Advanced conformance only)

### Plugin Development Guidance

**Target Performance**:

- **GetPluginInfo**: < 100ms (Standard), < 50ms (Advanced)
- **Should not**: Make external API calls
- **Should**: Return cached/static metadata
- **Should**: Validate responses quickly (local validation only)

**Common Performance Issues**:

- External API calls in GetPluginInfo → Cache metadata
- Expensive validation in GetPluginInfo → Move to plugin initialization
- Blocking I/O in GetPluginInfo → Use in-memory data structures

---

## Future Enhancements

### Potential Improvements

1. **Percentile Tracking**: Track P50/P95/P99 latencies (not just min/avg/max)
2. **Memory Allocation Tracking**: Measure bytes allocated per call (for memory leaks)
3. **Custom Iteration Count**: Allow configurable iteration count (not fixed at 10)
4. **Warmup Iterations**: Add warmup phase before measurement (to account for JIT compilation)
5. **Time Series Recording**: Store results over time for performance regression detection

### Extension Points

- `measureLatency()` can be extended with custom metrics
- `PerformanceResult` can add new fields (Percentiles, AllocsPerOp)
- `PerformanceBaseline` can be configured per-environment (dev vs prod thresholds)
