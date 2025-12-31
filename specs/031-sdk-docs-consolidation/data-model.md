# Data Model: SDK Documentation Consolidation

**Date**: 2025-12-31
**Branch**: `031-sdk-docs-consolidation`
**Status**: Complete

## Overview

This feature is **documentation-only** and does not introduce new data structures or entities.
The data model section documents the existing structures that will be referenced in
documentation.

## Documentation Targets

### Client Struct (Existing)

**Location**: `sdk/go/pluginsdk/client.go`

```go
type Client struct {
    inner      pbcconnect.CostSourceServiceClient
    httpClient *http.Client
    ownsClient bool  // true if SDK created the HTTP client
}
```

**Documentation Focus**:

- HTTP client ownership semantics
- Thread safety guarantees
- Connection pool configuration

### ResourceMatcher Struct (Existing)

**Location**: `sdk/go/pluginsdk/helpers.go`

```go
type ResourceMatcher struct {
    supportedProviders map[string]bool
    supportedTypes     map[string]bool
}
```

**Documentation Focus**:

- Thread safety constraints (configure before Serve())
- Read-only contract after initialization

### FocusRecordBuilder Struct (Existing)

**Location**: `sdk/go/pluginsdk/focus_builder.go`

```go
type FocusRecordBuilder struct {
    record *pbc.FocusCostRecord
}
```

**Documentation Focus**:

- Single-threaded builder pattern
- Shared-map semantics for WithTags
- FOCUS 1.2/1.3 compatibility

### WebConfig Struct (Existing)

**Location**: `sdk/go/pluginsdk/options.go`

```go
type WebConfig struct {
    Enabled              bool
    AllowedOrigins       []string
    AllowCredentials     bool
    EnableHealthEndpoint bool
}
```

**Documentation Focus**:

- CORS deployment scenarios
- Security guidelines
- Value semantics and defensive copying

### ServerTimeouts Struct (Existing)

**Location**: `sdk/go/pluginsdk/sdk.go`

```go
type ServerTimeouts struct {
    ReadHeaderTimeout time.Duration
    ReadTimeout       time.Duration
    WriteTimeout      time.Duration
    IdleTimeout       time.Duration
}
```

**Documentation Focus**:

- Performance tuning
- DoS protection guidelines

## No API Contracts

This feature does not introduce new APIs. All documentation changes are inline comments,
godoc examples, and README sections for existing APIs.

## Validation Rules

Documentation must satisfy:

1. **Compile Check**: All code examples must compile without modification
2. **Lint Check**: All markdown must pass `make lint-markdown`
3. **Copy-Paste Ready**: Examples use realistic values, no placeholders
4. **Cross-Provider**: Examples include AWS, Azure, GCP where applicable
