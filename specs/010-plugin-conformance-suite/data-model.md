# Data Model: Plugin Conformance Test Suite

**Feature**: 011-plugin-conformance-suite
**Date**: 2025-11-28

## Entity Relationship Overview

```text
┌─────────────────────┐
│  ConformanceSuite   │
│  (main entry point) │
└──────────┬──────────┘
           │ contains
           ▼
┌─────────────────────┐      ┌─────────────────────┐
│   TestCategory      │──────│  ConformanceLevel   │
│  (SpecValidation,   │      │  (Basic, Standard,  │
│   RPCCorrectness,   │      │   Advanced)         │
│   Performance,      │      └─────────────────────┘
│   Concurrency)      │
└──────────┬──────────┘
           │ produces
           ▼
┌─────────────────────┐
│ ConformanceResult   │
│ (JSON serializable) │
└──────────┬──────────┘
           │ contains
           ▼
┌─────────────────────┐      ┌─────────────────────┐
│   CategoryResult    │──────│    TestResult       │
│ (per-category       │      │ (per-test details)  │
│  summary)           │      └─────────────────────┘
└─────────────────────┘
```

## Entities

### ConformanceLevel

Enumeration of certification levels with associated requirements.

| Field | Type   | Description                                  |
| ----- | ------ | -------------------------------------------- |
| Value | int    | Enum value (0=Basic, 1=Standard, 2=Advanced) |
| Name  | string | Human-readable name                          |

**Values**:

- `ConformanceLevelBasic` (0) - Core functionality, required for all plugins
- `ConformanceLevelStandard` (1) - Production readiness, recommended for deployment
- `ConformanceLevelAdvanced` (2) - High performance, for demanding environments

**Validation Rules**:

- Value must be 0, 1, or 2
- Higher levels include all tests from lower levels

### TestCategory

Grouping of related conformance tests.

| Field       | Type   | Description                        |
| ----------- | ------ | ---------------------------------- |
| Name        | string | Category identifier (unique)       |
| Description | string | Human-readable description         |
| MinLevel    | int    | Minimum conformance level required |

**Values**:

- `SpecValidation` - Schema compliance and enum validation (Basic+)
- `RPCCorrectness` - Valid/invalid input handling (Basic+)
- `Performance` - Latency and allocation benchmarks (Standard+)
- `Concurrency` - Parallel request handling (Standard+)

### ConformanceSuite

Main entry point for running conformance tests.

| Field      | Type                   | Description              |
| ---------- | ---------------------- | ------------------------ |
| tests      | []ConformanceTest      | All registered tests     |
| config     | SuiteConfig            | Runtime configuration    |
| categories | map string to int list | Test indices by category |

**Relationships**:

- Contains multiple ConformanceTest instances
- Produces ConformanceResult on execution

### SuiteConfig

Configuration options for suite execution.

| Field             | Type          | Default | Description                  |
| ----------------- | ------------- | ------- | ---------------------------- |
| TargetLevel       | int           | 1       | Target conformance level     |
| Timeout           | time.Duration | 60s     | Per-test timeout             |
| ParallelRequests  | int           | 10      | Concurrency test parallelism |
| EnableBenchmarks  | bool          | true    | Run performance benchmarks   |
| BenchmarkDuration | time.Duration | 5s      | Duration per benchmark       |

### ConformanceTest

Single conformance test definition.

| Field       | Type                           | Description                 |
| ----------- | ------------------------------ | --------------------------- |
| Name        | string                         | Unique test identifier      |
| Description | string                         | What the test validates     |
| Category    | TestCategory                   | Parent category             |
| MinLevel    | ConformanceLevel               | Minimum level for this test |
| TestFunc    | func(\*TestHarness) TestResult | Test implementation         |

**Validation Rules**:

- Name must be unique within suite
- Name should follow Go test naming conventions (PascalCase)

### TestResult

Result of a single test execution (existing type, extended).

| Field    | Type          | Description                     |
| -------- | ------------- | ------------------------------- |
| Name     | string        | Test name                       |
| Method   | string        | RPC method tested               |
| Success  | bool          | Pass/fail status                |
| Error    | error         | Error if failed (nil if passed) |
| Duration | time.Duration | Execution time                  |
| Details  | string        | Human-readable details          |
| Category | string        | Parent category name            |

### CategoryResult

Aggregated results for a test category.

| Field   | Type         | Description             |
| ------- | ------------ | ----------------------- |
| Name    | string       | Category name           |
| Passed  | int          | Number of passed tests  |
| Failed  | int          | Number of failed tests  |
| Skipped | int          | Number of skipped tests |
| Results | []TestResult | Individual test results |

### ConformanceResult

Complete result of suite execution (JSON output).

| Field         | Type                      | Description                     |
| ------------- | ------------------------- | ------------------------------- |
| Version       | string                    | Report schema version ("1.0.0") |
| Timestamp     | time.Time                 | Execution timestamp             |
| PluginName    | string                    | Name from plugin's Name() RPC   |
| LevelAchieved | ConformanceLevel          | Highest level passed            |
| Summary       | ResultSummary             | Aggregate counts                |
| Categories    | map[string]CategoryResult | Results by category             |
| Duration      | time.Duration             | Total execution time            |

**JSON Serialization**:

```go
type ConformanceResult struct {
    Version       string                    `json:"version"`
    Timestamp     time.Time                 `json:"timestamp"`
    PluginName    string                    `json:"plugin_name"`
    LevelAchieved string                    `json:"level_achieved"`
    Summary       ResultSummary             `json:"summary"`
    Categories    map[string]CategoryResult `json:"categories"`
    Duration      string                    `json:"duration"`
}
```

### ResultSummary

Aggregate test counts.

| Field   | Type | Description           |
| ------- | ---- | --------------------- |
| Total   | int  | Total tests executed  |
| Passed  | int  | Tests that passed     |
| Failed  | int  | Tests that failed     |
| Skipped | int  | Tests skipped (level) |

### PerformanceBaseline

Threshold values for performance conformance.

| Field           | Type          | Description                 |
| --------------- | ------------- | --------------------------- |
| Method          | string        | RPC method name             |
| StandardLatency | time.Duration | Standard level threshold    |
| AdvancedLatency | time.Duration | Advanced level threshold    |
| MaxAllocBytes   | int64         | Maximum allocation per call |

**Canonical Values** (from sdk/go/testing/README.md):

| Method              | Standard | Advanced |
| ------------------- | -------- | -------- |
| Name                | 100ms    | 50ms     |
| Supports            | 50ms     | 25ms     |
| GetProjectedCost    | 200ms    | 100ms    |
| GetPricingSpec      | 200ms    | 100ms    |
| GetActualCost (24h) | 2s       | 1s       |
| GetActualCost (30d) | N/A      | 10s      |

### ValidationError

Structured error for field-level validation failures.

| Field    | Type        | Description                  |
| -------- | ----------- | ---------------------------- |
| Field    | string      | Field name that failed       |
| Value    | interface{} | Actual value received        |
| Expected | string      | Expected value or constraint |
| Message  | string      | Human-readable error message |

**Implements**: `error` interface

## State Transitions

### ConformanceResult.LevelAchieved

```text
[Not Run] → Basic → Standard → Advanced
    ↓         ↓         ↓
  (fail)   (fail)    (fail)
    ↓         ↓         ↓
  None     Basic    Standard
```

Determination logic:

1. Run all Basic tests → if any fail, LevelAchieved = None
2. Run all Standard tests → if any fail, LevelAchieved = Basic
3. Run all Advanced tests → if any fail, LevelAchieved = Standard
4. All pass → LevelAchieved = Advanced

## Relationships Summary

| From              | To                | Cardinality | Description                 |
| ----------------- | ----------------- | ----------- | --------------------------- |
| ConformanceSuite  | ConformanceTest   | 1:N         | Suite contains tests        |
| ConformanceTest   | TestCategory      | N:1         | Test belongs to category    |
| ConformanceTest   | ConformanceLevel  | N:1         | Test has minimum level      |
| ConformanceSuite  | ConformanceResult | 1:1         | Execution produces result   |
| ConformanceResult | CategoryResult    | 1:N         | Result has category results |
| CategoryResult    | TestResult        | 1:N         | Category has test results   |
