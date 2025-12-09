# Implementation Plan: SDK Support for PULUMICOST_LOG_FILE

**Branch**: `015-log-file` | **Date**: 2025-12-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/015-log-file/spec.md`

## Summary

Enable plugins to redirect all log output to a file specified by `PULUMICOST_LOG_FILE` environment
variable. When set, the SDK's logging infrastructure writes to the file instead of stderr. This
allows the Core CLI to orchestrate plugins without log pollution in user-facing output.

**Key finding from research**: The SDK already has infrastructure in place:

- `EnvLogFile` constant and `GetLogFile()` function exist in `sdk/go/pluginsdk/env.go`
- `NewPluginLogger()` accepts an `io.Writer` parameter
- `newDefaultLogger()` currently hardcodes `os.Stderr`

The implementation requires enhancing these existing functions to:

1. Auto-configure output based on `PULUMICOST_LOG_FILE` when present
2. Handle file creation, append mode, and error fallback
3. Provide a high-level API that plugins can use without manual setup

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: zerolog v1.34.0+ (already in go.mod), stdlib only for file operations
**Storage**: File system (log file) - append mode with 0644 permissions
**Testing**: go test with existing sdk/go/pluginsdk test patterns
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: SDK library (single package addition)
**Performance Goals**: <10ms additional startup latency, no per-log-call overhead
**Constraints**: Backward compatible - no breaking changes to existing APIs
**Scale/Scope**: Single package modification (`sdk/go/pluginsdk`)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | N/A | No proto changes - SDK helper code only |
| II. Multi-Provider gRPC Consistency | N/A | Provider-agnostic logging feature |
| III. Test-First Protocol | PASS | Tests written before implementation |
| IV. Protobuf Backward Compatibility | N/A | No proto changes |
| V. Comprehensive Documentation | PASS | README update required |
| VI. Performance as gRPC Requirement | PASS | <10ms latency budget specified |
| VII. Validation at Multiple Levels | PASS | Unit tests + integration tests |

**Gate Result**: PASS - All applicable principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/015-log-file/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (via /speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── env.go              # Existing - already has EnvLogFile, GetLogFile()
├── env_test.go         # Existing - add tests for log file scenarios
├── logging.go          # MODIFY - add NewLogWriter(), update newDefaultLogger()
├── logging_test.go     # MODIFY - add log file configuration tests
└── README.md           # MODIFY - document PULUMICOST_LOG_FILE usage
```

**Structure Decision**: Single package modification - all changes confined to `sdk/go/pluginsdk`
with no new packages or directories needed.
