# Data Model: Centralized Environment Variable Handling

**Branch**: `013-pluginsdk-env` | **Date**: 2025-12-07

## Entity Overview

This feature introduces no persistent data storage. All entities are in-memory constants and
functions that read environment variables at runtime.

## Constants (Environment Variable Names)

### EnvPort

- **Value**: `"PULUMICOST_PLUGIN_PORT"`
- **Type**: `string` (constant)
- **Purpose**: Canonical environment variable for plugin gRPC port
- **Used by**: `GetPort()` function

### EnvLogLevel

- **Value**: `"PULUMICOST_LOG_LEVEL"`
- **Type**: `string` (constant)
- **Purpose**: Canonical environment variable for log verbosity
- **Used by**: `GetLogLevel()` function

### EnvLogLevelFallback

- **Value**: `"LOG_LEVEL"`
- **Type**: `string` (constant)
- **Purpose**: Legacy fallback for log level configuration
- **Used by**: `GetLogLevel()` function when `EnvLogLevel` is not set

### EnvLogFormat

- **Value**: `"PULUMICOST_LOG_FORMAT"`
- **Type**: `string` (constant)
- **Purpose**: Environment variable for log output format
- **Valid values**: `"json"`, `"text"` (or empty for plugin default)

### EnvLogFile

- **Value**: `"PULUMICOST_LOG_FILE"`
- **Type**: `string` (constant)
- **Purpose**: Environment variable for log file path

### EnvTraceID

- **Value**: `"PULUMICOST_TRACE_ID"`
- **Type**: `string` (constant)
- **Purpose**: Environment variable for injecting external trace IDs

### EnvTestMode

- **Value**: `"PULUMICOST_TEST_MODE"`
- **Type**: `string` (constant)
- **Purpose**: Environment variable for enabling plugin test mode
- **Valid values**: `"true"`, `"false"` (strict matching)

## Functions

### GetPort() int

- **Inputs**: None (reads from environment)
- **Outputs**: `int` - Port number or 0 if not configured/invalid
- **Validation**:
  - Reads `PULUMICOST_PLUGIN_PORT` only (no fallback)
  - Returns 0 if not set or invalid
  - Only positive integers are valid

### GetLogLevel() string

- **Inputs**: None (reads from environment)
- **Outputs**: `string` - Log level or empty string
- **Validation**:
  - Reads `PULUMICOST_LOG_LEVEL` first
  - Falls back to `LOG_LEVEL` if canonical not set
  - No value validation (plugin responsibility)

### GetLogFormat() string

- **Inputs**: None (reads from environment)
- **Outputs**: `string` - Log format or empty string
- **Validation**: None (pass-through)

### GetLogFile() string

- **Inputs**: None (reads from environment)
- **Outputs**: `string` - File path or empty string
- **Validation**: None (file system validation is caller responsibility)

### GetTraceID() string

- **Inputs**: None (reads from environment)
- **Outputs**: `string` - Trace ID or empty string
- **Validation**: None (pass-through)

### GetTestMode() bool

- **Inputs**: None (reads from environment)
- **Outputs**: `bool` - true only if `PULUMICOST_TEST_MODE` == `"true"`
- **Validation**:
  - Returns true only for exact string `"true"`
  - Returns false for `"false"`, empty, or any other value
  - Logs warning via zerolog when value is not `"true"` or `"false"`

### IsTestMode() bool

- **Inputs**: None (reads from environment)
- **Outputs**: `bool` - Same as `GetTestMode()` but without warning logging
- **Purpose**: Convenience function for repeated checks without log spam

## State Transitions

Not applicable. All functions are pure reads with no state modifications.

## Relationships

```text
┌─────────────────────────────────────────────────────────────┐
│                    Environment Variables                     │
├─────────────────────────────────────────────────────────────┤
│ PULUMICOST_PLUGIN_PORT ─►│──► GetPort() ──► int             │
│                          │                                  │
│ PULUMICOST_LOG_LEVEL ───┐│                                  │
│ LOG_LEVEL ──────────────►│──► GetLogLevel() ──► string      │
│                          │                                  │
│ PULUMICOST_LOG_FORMAT ──►│──► GetLogFormat() ──► string     │
│ PULUMICOST_LOG_FILE ────►│──► GetLogFile() ──► string       │
│ PULUMICOST_TRACE_ID ────►│──► GetTraceID() ──► string       │
│ PULUMICOST_TEST_MODE ───►│──► GetTestMode() ──► bool        │
│                          │                     (logs warn)  │
│                          │──► IsTestMode() ──► bool         │
│                          │                     (no logging) │
└─────────────────────────────────────────────────────────────┘
```

## Validation Rules

| Variable | Validation | Invalid Behavior |
|----------|------------|------------------|
| Port | Must be positive integer | Return 0 (no fallback) |
| Log Level | None | Pass through as-is |
| Log Format | None | Pass through as-is |
| Log File | None | Pass through as-is |
| Trace ID | None | Pass through as-is |
| Test Mode | Must be `"true"` or `"false"` | Return false, log warning |

## Data Volume / Scale

- Constants: 7 string constants (~180 bytes total)
- Functions: 7 functions, no allocations per call
- Performance: O(1) for all operations
- Memory: No dynamic allocation (stdlib only)
