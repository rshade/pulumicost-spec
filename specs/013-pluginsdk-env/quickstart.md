# Quickstart: Centralized Environment Variable Handling

**Branch**: `013-pluginsdk-env` | **Date**: 2025-12-07

## Overview

This feature adds centralized environment variable handling to the `pluginsdk` package,
providing a standardized way for PulumiCost plugins to read configuration from the
environment.

## Installation

The environment handling is part of the existing `pluginsdk` package. No additional
installation required:

```go
import "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
```

## Basic Usage

### Reading Port Configuration

```go
package main

import (
    "fmt"
    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    // GetPort reads PULUMICOST_PLUGIN_PORT only (no fallback)
    port := pluginsdk.GetPort()
    if port == 0 {
        fmt.Println("Error: PULUMICOST_PLUGIN_PORT must be set")
        return
    }
    fmt.Printf("Starting on port %d\n", port)
}
```

### Reading Logging Configuration

```go
// Get log level (PULUMICOST_LOG_LEVEL or LOG_LEVEL)
logLevel := pluginsdk.GetLogLevel()
if logLevel == "" {
    logLevel = "info" // plugin default
}

// Get log format (PULUMICOST_LOG_FORMAT)
logFormat := pluginsdk.GetLogFormat()
if logFormat == "" {
    logFormat = "json" // plugin default
}

// Get log file path (PULUMICOST_LOG_FILE)
logFile := pluginsdk.GetLogFile()
// Empty means stdout
```

### Reading Trace ID

```go
// Get trace ID for distributed tracing (PULUMICOST_TRACE_ID)
traceID := pluginsdk.GetTraceID()
if traceID != "" {
    // Use external trace ID for correlation
    logger = logger.With().Str("trace_id", traceID).Logger()
}
```

### Checking Test Mode

```go
// IsTestMode returns true only when PULUMICOST_TEST_MODE="true"
if pluginsdk.IsTestMode() {
    // Use mock data or bypass external services
    return mockPricingData()
}

// GetTestMode does the same but logs a warning for invalid values
// Use GetTestMode during startup to validate configuration
if pluginsdk.GetTestMode() {
    fmt.Println("Running in test mode")
}
```

## Using with Serve()

The `pluginsdk.Serve()` function automatically uses `GetPort()` internally:

```go
func main() {
    ctx := context.Background()
    plugin := &MyPlugin{}

    // Serve reads PULUMICOST_PLUGIN_PORT only (no PORT fallback)
    // Set Port: 0 to use environment variable
    err := pluginsdk.Serve(ctx, pluginsdk.ServeConfig{
        Plugin: plugin,
        Port:   0, // Uses GetPort() internally
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Environment Variable Reference

### Port Configuration

| Variable | Description |
|----------|-------------|
| `PULUMICOST_PLUGIN_PORT` | Required port variable (no fallback) |

### Logging Configuration

| Variable | Priority | Description |
|----------|----------|-------------|
| `PULUMICOST_LOG_LEVEL` | 1 (primary) | Log verbosity (debug, info, warn, error) |
| `LOG_LEVEL` | 2 (fallback) | Legacy log level variable |
| `PULUMICOST_LOG_FORMAT` | - | Log format (json, text) |
| `PULUMICOST_LOG_FILE` | - | Log file path (empty = stdout) |

### Tracing Configuration

| Variable | Description |
|----------|-------------|
| `PULUMICOST_TRACE_ID` | External trace ID for distributed tracing |

### Test Mode Configuration

| Variable | Valid Values | Description |
|----------|--------------|-------------|
| `PULUMICOST_TEST_MODE` | `"true"`, `"false"` | Enable test mode (strict matching) |

## Accessing Constants

All environment variable names are exported as constants for reference:

```go
fmt.Println(pluginsdk.EnvPort)             // "PULUMICOST_PLUGIN_PORT"
fmt.Println(pluginsdk.EnvLogLevel)         // "PULUMICOST_LOG_LEVEL"
fmt.Println(pluginsdk.EnvLogLevelFallback) // "LOG_LEVEL"
fmt.Println(pluginsdk.EnvLogFormat)        // "PULUMICOST_LOG_FORMAT"
fmt.Println(pluginsdk.EnvLogFile)          // "PULUMICOST_LOG_FILE"
fmt.Println(pluginsdk.EnvTraceID)          // "PULUMICOST_TRACE_ID"
fmt.Println(pluginsdk.EnvTestMode)         // "PULUMICOST_TEST_MODE"
```

## Migration Guide

### Migrating from os.Getenv("PORT")

**Breaking Change**: The `PORT` environment variable is no longer supported. You must use
`PULUMICOST_PLUGIN_PORT` instead.

Before:

```go
portStr := os.Getenv("PORT")
port, _ := strconv.Atoi(portStr)
```

After:

```go
// Environment: PULUMICOST_PLUGIN_PORT=8080 (required)
port := pluginsdk.GetPort()
if port == 0 {
    log.Fatal("PULUMICOST_PLUGIN_PORT must be set")
}
```

### Migrating from os.Getenv("LOG_LEVEL")

Before:

```go
logLevel := os.Getenv("LOG_LEVEL")
```

After:

```go
logLevel := pluginsdk.GetLogLevel()
```

## Testing

Test your plugin with different environment configurations:

```bash
# Test with port variable (required)
PULUMICOST_PLUGIN_PORT=8080 go test ./...

# Test without port (should fail or use explicit default)
go test ./...

# Test logging configuration
PULUMICOST_LOG_LEVEL=debug PULUMICOST_LOG_FORMAT=json go test ./...

# Test log level fallback (legacy support)
LOG_LEVEL=debug go test ./...

# Test mode
PULUMICOST_TEST_MODE=true go test ./...
```
