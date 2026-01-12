# Research: Centralized Environment Variable Handling

**Branch**: `013-pluginsdk-env` | **Date**: 2025-12-07

## Research Summary

This feature adds centralized environment variable handling to the pluginsdk package. Research
focused on understanding existing patterns in the codebase and Go best practices.

## Decisions

### Decision 1: File Location

**Decision**: Create `sdk/go/pluginsdk/env.go`

**Rationale**: The pluginsdk package already contains the `Serve()` function that needs
modification. Adding env.go to the same package avoids cross-package dependencies and follows
the existing SDK structure pattern.

**Alternatives considered**:

- `sdk/go/config/` - Would require new package and import cycle risks
- `sdk/go/env/` - Unnecessary package fragmentation for 7 constants and 7 functions

### Decision 2: Port Implementation Pattern (No Fallback)

**Decision**: Use canonical variable only for port configuration

```go
func GetPort() int {
    if v := os.Getenv(EnvPort); v != "" {
        if port, err := strconv.Atoi(v); err == nil && port > 0 {
            return port
        }
    }
    return 0
}
```

**Rationale**: The `PORT` environment variable has too much potential for conflict with other
services. Using `PULUMICOST_PLUGIN_PORT` exclusively ensures clear ownership and prevents
accidental port conflicts. This is a breaking change but provides clearer semantics.

**Alternatives considered**:

- Fallback to PORT - Rejected; too much conflict potential with other services
- Return error for invalid values - Rejected; complicates API and existing code uses return 0

### Decision 2b: Log Level Fallback Pattern

**Decision**: Use canonical-first, fallback-second for logging configuration only

```go
func GetLogLevel() string {
    if v := os.Getenv(EnvLogLevel); v != "" {
        return v
    }
    return os.Getenv(EnvLogLevelFallback)
}
```

**Rationale**: LOG_LEVEL is a common convention that doesn't have the same conflict potential
as PORT. Supporting fallback enables gradual migration of existing plugins.

### Decision 3: Test Mode Validation

**Decision**: Use strict `"true"`/`"false"` string matching with warning logging

**Rationale**: Follows existing pattern in `finfocus-plugin-aws-public/internal/plugin
/testmode.go`. Strict matching prevents accidental test mode activation from typos like
`"yes"` or `"1"`.

**Alternatives considered**:

- Case-insensitive matching - Rejected; too permissive for safety-critical flag
- Accept `"1"`/`"0"` - Rejected; inconsistent with existing plugin pattern

### Decision 4: Logging for Invalid Values

**Decision**: Only log warnings for invalid `PULUMICOST_TEST_MODE` values (per FR-013)

**Rationale**: Test mode is boolean with discrete valid values; other variables are either
string pass-through (log level, format) or numeric with silent fallback (port).

**Alternatives considered**:

- Log warnings for all invalid values - Rejected; port fallback is expected behavior
- No logging at all - Rejected; FR-013 explicitly requires warning

### Decision 5: No Dependencies Beyond stdlib

**Decision**: Use only `os`, `strconv`, `strings` from Go stdlib

**Rationale**: Environment variable reading is trivial and doesn't benefit from external
dependencies. Avoids adding complexity to the SDK.

**Alternatives considered**:

- Use `envconfig` library - Rejected; overkill for 8 variables
- Use `viper` - Rejected; massive dependency for simple env reading

## Existing Code Analysis

### Current Port Resolution (`sdk/go/pluginsdk/sdk.go:210-231`)

```go
func resolvePort(requested int) (int, error) {
    if requested > 0 {
        return requested, nil
    }
    portEnv := os.Getenv("PORT")
    if portEnv == "" {
        return 0, nil
    }
    value, err := strconv.Atoi(portEnv)
    if err != nil {
        // ... error handling
        return 0, nil
    }
    return value, nil
}
```

**Analysis**:

- Only reads `PORT`, not `PULUMICOST_PLUGIN_PORT`
- Returns 0 when not set (caller handles)
- Logs to stderr on parse error
- Will be modified to use `GetPort()` which adds canonical variable

### Existing Test Mode Pattern (`finfocus-plugin-aws-public`)

```go
const testModeEnvVar = "PULUMICOST_TEST_MODE"

func IsTestMode() bool {
    return os.Getenv(testModeEnvVar) == "true"
}

func ValidateTestModeEnv() {
    value := os.Getenv(testModeEnvVar)
    if value != "" && value != "true" && value != "false" {
        // log warning
    }
}
```

**Analysis**: This pattern will be adopted in the SDK with minor modifications.

### Logging Pattern (`finfocus-core`)

Uses `zerolog` for structured logging. The SDK already imports zerolog (see `sdk.go:12`).

## Technical Notes

### Environment Variable Constants

All constants follow the `PULUMICOST_` prefix convention established by finfocus-core:

| Constant | Value | Type |
|----------|-------|------|
| EnvPort | `PULUMICOST_PLUGIN_PORT` | Canonical (no fallback) |
| EnvLogLevel | `PULUMICOST_LOG_LEVEL` | Canonical |
| EnvLogLevelFallback | `LOG_LEVEL` | Fallback (for logging only) |
| EnvLogFormat | `PULUMICOST_LOG_FORMAT` | Canonical |
| EnvLogFile | `PULUMICOST_LOG_FILE` | Canonical |
| EnvTraceID | `PULUMICOST_TRACE_ID` | Canonical |
| EnvTestMode | `PULUMICOST_TEST_MODE` | Canonical |

### Function Signatures

Based on existing patterns and requirements:

```go
// Port - returns 0 when not configured
func GetPort() int

// Logging - returns empty string when not configured
func GetLogLevel() string
func GetLogFormat() string
func GetLogFile() string

// Tracing - returns empty string when not configured
func GetTraceID() string

// Test mode - returns false when not configured or invalid
func GetTestMode() bool      // logs warning for invalid values
func IsTestMode() bool       // convenience alias, no logging
```

## No Outstanding Research Items

All NEEDS CLARIFICATION items from Technical Context have been resolved:

- Language/Version: Go 1.25.5 (confirmed from go.mod)
- Dependencies: stdlib only (`os`, `strconv`, `strings`)
- Testing: Go testing framework (existing pattern in pluginsdk)
- Platform: Cross-platform Go library

## Next Phase

Proceed to Phase 1: Design & Contracts with `data-model.md` and `quickstart.md`.
