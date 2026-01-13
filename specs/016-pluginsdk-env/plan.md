# Implementation Plan: Centralized Environment Variable Handling

**Branch**: `013-pluginsdk-env` | **Date**: 2025-12-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/016-pluginsdk-env/spec.md`

## Summary

Add centralized environment variable handling to the pluginsdk package to standardize how
plugins read configuration from environment variables. The implementation provides exported
constants for canonical variable names and getter functions with fallback support for backward
compatibility.

**Primary Goal**: Fix the port mismatch issue where core sets `PULUMICOST_PLUGIN_PORT` but
plugins read `PORT`.

**Technical Approach**: Create `sdk/go/pluginsdk/env.go` with constants and getter functions
that implement canonical-first, fallback-second reading pattern using Go stdlib only.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod toolchain)
**Primary Dependencies**: Go stdlib only (`os`, `strconv`, `strings`)
**Storage**: N/A (reads environment variables at runtime)
**Testing**: Go testing framework with table-driven tests
**Target Platform**: Cross-platform Go library (SDK)
**Project Type**: Single Go module (SDK library)
**Performance Goals**: N/A (simple env var reads, ~nanosecond operations)
**Constraints**: Must maintain 100% backward compatibility with `LOG_LEVEL` (no PORT fallback)
**Scale/Scope**: Used by all FinFocus plugins (3+ repos currently)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Proto-First Development | N/A | SDK helper code, no proto changes |
| II. Multi-Provider Consistency | PASS | Environment variables are provider-agnostic |
| III. Test-First Protocol | PASS | Will write unit tests before implementation |
| IV. Backward Compatibility | PASS | Fallback variables maintain compatibility |
| V. Documentation | PASS | Will document all env vars and functions |
| VI. Performance | N/A | Simple env var reads, no performance concerns |
| VII. Validation | PASS | Tests will run in CI |

**Gate Result**: PASS - All applicable principles satisfied. This is SDK helper code that
enhances the plugin development experience without modifying the gRPC protocol.

## Project Structure

### Documentation (this feature)

```text
specs/013-pluginsdk-env/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A for this feature - no API contracts)
├── checklists/          # Quality checklists
│   └── requirements.md  # Specification validation checklist
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/pluginsdk/
├── env.go               # NEW: Environment variable constants and getters
├── env_test.go          # NEW: Unit tests for env.go
├── sdk.go               # MODIFY: Use GetPort() instead of os.Getenv("PORT")
└── ...                  # Existing pluginsdk files unchanged
```

**Structure Decision**: Single file addition to existing `sdk/go/pluginsdk/` package. No new
packages or directories required. This follows the existing SDK structure pattern.

## Complexity Tracking

No constitution violations. All gates pass. No complexity justification required.
