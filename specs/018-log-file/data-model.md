# Data Model: SDK Support for PULUMICOST_LOG_FILE

**Date**: 2025-12-08
**Feature**: 015-log-file

## Overview

This feature adds log file configuration to the SDK. The data model is minimal as this is
primarily a configuration and I/O feature rather than a domain model change.

## Entities

### LogConfiguration (Conceptual)

The log configuration is not a struct but a set of values derived from environment variables
at runtime. The configuration is resolved once at startup and used to create the logger.

| Attribute | Type | Source | Default |
|-----------|------|--------|---------|
| Output Destination | `io.Writer` | `PULUMICOST_LOG_FILE` | `os.Stderr` |
| Log Level | `zerolog.Level` | `PULUMICOST_LOG_LEVEL` | `zerolog.InfoLevel` |
| Log Format | `string` | `PULUMICOST_LOG_FORMAT` | `json` |

**Resolution Rules**:

1. If `PULUMICOST_LOG_FILE` is empty or unset → use `os.Stderr`
2. If `PULUMICOST_LOG_FILE` is set but invalid → warn to stderr, use `os.Stderr`
3. If `PULUMICOST_LOG_FILE` is valid path → open file in append mode

### File Handle State

The log file handle is managed internally by the SDK:

```text
┌─────────────────────────────────────────────────────────────┐
│                    Log File Lifecycle                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Plugin Startup                                              │
│       │                                                      │
│       ▼                                                      │
│  ┌─────────────┐     env not set      ┌────────────────┐    │
│  │ GetLogFile()├─────────────────────►│ Use os.Stderr  │    │
│  └─────┬───────┘                      └────────────────┘    │
│        │ env set                                             │
│        ▼                                                     │
│  ┌─────────────┐     is directory     ┌────────────────┐    │
│  │ Stat path   ├─────────────────────►│ Warn + Stderr  │    │
│  └─────┬───────┘                      └────────────────┘    │
│        │ is file or not exists                               │
│        ▼                                                     │
│  ┌─────────────┐     open error       ┌────────────────┐    │
│  │ OpenFile()  ├─────────────────────►│ Warn + Stderr  │    │
│  └─────┬───────┘                      └────────────────┘    │
│        │ success                                             │
│        ▼                                                     │
│  ┌─────────────────┐                                        │
│  │ Return *os.File │                                        │
│  └─────────────────┘                                        │
│                                                              │
│  Plugin Shutdown                                             │
│       │                                                      │
│       ▼                                                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │ OS closes file handle on process exit (no explicit  │    │
│  │ close required)                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## Relationships

### Environment Variables → Logger Configuration

```text
┌──────────────────────────┐
│   Environment Variables   │
├──────────────────────────┤
│ PULUMICOST_LOG_FILE      │──┐
│ PULUMICOST_LOG_LEVEL     │  │
│ PULUMICOST_LOG_FORMAT    │  │
└──────────────────────────┘  │
                              │
                              ▼
                    ┌─────────────────┐
                    │  NewLogWriter() │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   io.Writer     │
                    │ (file or stderr)│
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ NewPluginLogger │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ zerolog.Logger  │
                    └─────────────────┘
```

## Validation Rules

| Rule | Validation | Error Handling |
|------|------------|----------------|
| Path is directory | `os.Stat()` check | Fallback to stderr |
| Path not writable | `os.OpenFile()` fails | Fallback to stderr |
| Parent dir missing | `os.OpenFile()` fails | Fallback to stderr |
| Empty path | Length check | Use stderr (expected) |

## State Transitions

This feature has no explicit state machine. The log writer is determined once at startup
and remains constant for the process lifetime.

## Constants

| Name | Value | Purpose |
|------|-------|---------|
| `LogFilePermissions` | `0644` | File permissions for created log files |
| `LogFileFlags` | `O_APPEND\|O_CREATE\|O_WRONLY` | File open flags |

## Constraints

1. **Immutable after creation**: Log writer cannot be changed after logger initialization
2. **Single destination**: Log output goes to exactly one writer (file or stderr)
3. **No rotation**: File rotation is external responsibility
4. **Process-scoped**: File handle lifetime matches process lifetime
