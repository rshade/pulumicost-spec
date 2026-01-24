# Data Model: Validation Bypass Metadata

**Feature**: 034-validation-bypass-metadata
**Date**: 2026-01-24
**Phase**: 1 - Design

## Entity Overview

```text
┌─────────────────────────────────────────────────────────────────┐
│                      ValidationResult                           │
├─────────────────────────────────────────────────────────────────┤
│ Valid      bool                                                 │
│ Errors     []string                                             │
│ Warnings   []string                                             │
│ Bypasses   []BypassMetadata  ◄── NEW                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ 0..* contains
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       BypassMetadata                            │
├─────────────────────────────────────────────────────────────────┤
│ Timestamp       time.Time        (when bypass occurred)         │
│ ValidationName  string           (which validation bypassed)    │
│ OriginalError   string           (error that would have shown)  │
│ Reason          string           (why bypassed, max 500 chars)  │
│ Operator        string           (who triggered bypass)         │
│ Severity        BypassSeverity   (risk level)                   │
│ Mechanism       BypassMechanism  (how triggered)                │
│ Truncated       bool             (reason was truncated)         │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│     BypassSeverity      │     │     BypassMechanism     │
├─────────────────────────┤     ├─────────────────────────┤
│ "warning"               │     │ "flag"                  │
│ "error"                 │     │ "env_var"               │
│ "critical"              │     │ "config"                │
└─────────────────────────┘     │ "programmatic"          │
                                └─────────────────────────┘
```

## Entity Definitions

### ValidationResult (Extended)

**Location**: `sdk/go/pricing/observability.go`
**Purpose**: Represents the complete result of a validation operation including any policy
overrides.

| Field    | Type             | Required | Description                            |
| -------- | ---------------- | -------- | -------------------------------------- |
| Valid    | bool             | Yes      | Overall validation passed              |
| Errors   | []string         | No       | Validation error messages              |
| Warnings | []string         | No       | Validation warning messages            |
| Bypasses | []BypassMetadata | No       | **NEW**: Bypassed validation records   |

**Validation Rules**:

- `Valid` can be `true` even with bypasses (bypassed validations don't fail the result)
- `Bypasses` is nil/empty when no validations were bypassed
- `Bypasses` is non-empty only when at least one validation was actually bypassed

**JSON Representation**:

```json
{
  "valid": true,
  "errors": [],
  "warnings": ["Budget threshold approaching"],
  "bypasses": [
    {
      "timestamp": "2026-01-24T10:30:00Z",
      "validation_name": "budget_limit",
      "original_error": "Cost exceeds budget by $500",
      "reason": "Emergency deployment approved by manager",
      "operator": "user@example.com",
      "severity": "error",
      "mechanism": "flag"
    }
  ]
}
```

### BypassMetadata

**Location**: `sdk/go/pricing/bypass.go` (new file)
**Purpose**: Individual bypass event record for audit trail.

| Field          | Type            | Required | Description                              |
| -------------- | --------------- | -------- | ---------------------------------------- |
| Timestamp      | time.Time       | Yes      | When the bypass occurred (UTC)           |
| ValidationName | string          | Yes      | Identifier of bypassed validation        |
| OriginalError  | string          | Yes      | Error message that would have shown      |
| Reason         | string          | Yes      | Human-readable bypass reason (max 500)   |
| Operator       | string          | No       | Who triggered bypass ("unknown" if N/A)  |
| Severity       | BypassSeverity  | Yes      | Risk level of bypassed validation        |
| Mechanism      | BypassMechanism | Yes      | How the bypass was triggered             |
| Truncated      | bool            | No       | True if reason was truncated             |

**Validation Rules**:

- `Timestamp` must not be zero-value (use `time.Now().UTC()` as default)
- `ValidationName` must be non-empty
- `OriginalError` must be non-empty
- `Reason` truncated to 500 chars with "..." suffix if exceeded
- `Operator` defaults to "unknown" if not provided
- `Severity` must be valid enum value
- `Mechanism` must be valid enum value

### BypassSeverity (Enum)

**Location**: `sdk/go/pricing/bypass.go`
**Purpose**: Risk level classification for bypassed validations.

| Value      | Description                                      |
| ---------- | ------------------------------------------------ |
| `warning`  | Low-risk bypass, informational alert             |
| `error`    | Medium-risk bypass, would have blocked operation |
| `critical` | High-risk bypass, security or compliance impact  |

**Go Definition**:

```go
type BypassSeverity string

const (
    BypassSeverityWarning  BypassSeverity = "warning"
    BypassSeverityError    BypassSeverity = "error"
    BypassSeverityCritical BypassSeverity = "critical"
)
```

### BypassMechanism (Enum)

**Location**: `sdk/go/pricing/bypass.go`
**Purpose**: Classification of how the bypass was triggered.

| Value         | Description                                  |
| ------------- | -------------------------------------------- |
| `flag`        | Command-line flag (e.g., `--yolo`, `--force`) |
| `env_var`     | Environment variable override                |
| `config`      | Configuration file setting                   |
| `programmatic`| Code-level API call                          |

**Go Definition**:

```go
type BypassMechanism string

const (
    BypassMechanismFlag        BypassMechanism = "flag"
    BypassMechanismEnvVar      BypassMechanism = "env_var"
    BypassMechanismConfig      BypassMechanism = "config"
    BypassMechanismProgrammatic BypassMechanism = "programmatic"
)
```

## State Transitions

```text
Validation Check Flow
─────────────────────

                    ┌──────────────────┐
                    │   Start Check    │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │ Run Validation   │
                    └────────┬─────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
        ┌──────────┐  ┌──────────┐  ┌──────────────┐
        │  PASS    │  │  FAIL    │  │ FAIL+BYPASS  │
        └──────────┘  └────┬─────┘  └──────┬───────┘
                           │               │
                           │               ▼
                           │        ┌──────────────┐
                           │        │ Record       │
                           │        │ BypassMeta   │
                           │        └──────┬───────┘
                           │               │
                           ▼               ▼
                    ┌──────────────────────────┐
                    │   Return ValidationResult │
                    │   Valid=false or true    │
                    │   Bypasses=[...] or nil  │
                    └──────────────────────────┘
```

**State Rules**:

1. PASS → `Valid=true`, `Errors=[]`, `Bypasses=nil`
2. FAIL → `Valid=false`, `Errors=[msg]`, `Bypasses=nil`
3. FAIL+BYPASS → `Valid=true`, `Errors=[]`, `Bypasses=[metadata]`

## Relationships

| From             | To              | Cardinality | Description                    |
| ---------------- | --------------- | ----------- | ------------------------------ |
| ValidationResult | BypassMetadata  | 1:0..*      | Result contains bypass records |
| BypassMetadata   | BypassSeverity  | N:1         | Each bypass has one severity   |
| BypassMetadata   | BypassMechanism | N:1         | Each bypass has one mechanism  |

## Constraints Summary

| Constraint                    | Type       | Description                             |
| ----------------------------- | ---------- | --------------------------------------- |
| Reason max length             | Data       | 500 characters, truncate with "..."     |
| Timestamp not zero            | Validation | Must be set (warn if zero)              |
| ValidationName required       | Validation | Must be non-empty string                |
| Severity valid enum           | Validation | Must be warning/error/critical          |
| Mechanism valid enum          | Validation | Must be flag/env_var/config/programmatic|
| No bypass without failure     | Business   | Can only bypass actual validation fails |
| Retention minimum             | Policy     | 90 days (caller responsibility)         |
