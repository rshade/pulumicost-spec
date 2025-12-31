# Implementation Plan: SDK Documentation Consolidation

**Branch**: `031-sdk-docs-consolidation` | **Date**: 2025-12-31 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/031-sdk-docs-consolidation/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.
See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Consolidate and improve SDK documentation by addressing 12 open documentation issues covering
inline code comments, godoc examples, performance tuning guides, CORS best practices, migration
guides, rate limiting patterns, and thread safety documentation for the pluginsdk package.
This is a **documentation-only** feature with no proto or SDK code changes required.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod) + markdown documentation
**Primary Dependencies**: pluginsdk package, testing package, zerolog (logging examples)
**Storage**: N/A (pure documentation)
**Testing**: Validation via `make lint-markdown`, godoc example compilation, copy-paste verification
**Target Platform**: Documentation consumers (plugin developers, operators, contributors)
**Project Type**: Documentation consolidation across existing SDK packages
**Performance Goals**: N/A (documentation only)
**Constraints**: All examples must compile without modification (copy-paste ready)
**Scale/Scope**: 12 GitHub issues, ~15-20 documentation files affected

### Affected Files

**Inline Code Documentation (Quick Wins)**:

- `sdk/go/pluginsdk/client.go` - NewClient() ownership semantics comment (Issue #240)
- `sdk/go/pluginsdk/focus_conformance.go` - contractedCostTolerance explanation (Issue #211)
- `sdk/go/pluginsdk/focus_builder.go` - WithTags shared-map semantics (Issue #207)
- `sdk/go/pluginsdk/focus_conformance.go` - Correlation pattern comment fix (Issue #206)
- `sdk/go/testing/*_test.go` - Test naming clarification (Issue #208)

**Godoc Examples**:

- `sdk/go/pluginsdk/example_test.go` - Client.Close() example (Issue #238)
- `sdk/go/testing/README.md` - Complete import statements with pbc alias (Issue #209)

**Comprehensive Guides (README sections or new docs)**:

- `sdk/go/pluginsdk/README.md` - Migration Guide: gRPC to connect-go (Issue #235)
- `sdk/go/pluginsdk/README.md` - Performance Tuning section (Issue #237)
- `sdk/go/pluginsdk/README.md` - CORS Best Practices section (Issue #236)
- `sdk/go/pluginsdk/README.md` - Rate Limiting documentation (Issue #233)
- `sdk/go/pluginsdk/README.md` - Thread Safety documentation (Issue #231)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
| --------- | ------ | ----- |
| I. gRPC Proto Specification-First | ✅ PASS | No proto changes required - documentation only |
| II. Multi-Provider gRPC Consistency | ✅ PASS | Examples will include cross-provider patterns |
| III. Test-First Protocol | ⚠️ N/A | Documentation feature - validation via lint + compile |
| IV. Protobuf Backward Compatibility | ✅ PASS | No breaking changes - documentation only |
| V. Comprehensive Documentation | ✅ PASS | This feature directly addresses documentation currency |
| VI. Performance as a gRPC Requirement | ✅ PASS | Performance tuning guide will document benchmarks |
| VII. Validation at Multiple Levels | ✅ PASS | Markdown lint + godoc compile validation |

**Documentation Currency Gate** (NON-NEGOTIABLE):

- ✅ This PR specifically addresses 12 open documentation issues
- ✅ All SDK README files will be updated in this PR
- ✅ Example files will be added/updated as needed

**Example Quality Standards Gate** (NON-NEGOTIABLE):

- ✅ All examples will be copy-paste ready (compile without modification)
- ✅ Examples will use realistic values from existing codebase patterns
- ✅ No placeholder values ("TODO", "example-value", "xxx")
- ✅ Cross-provider examples where applicable (AWS, Azure, GCP)

## Project Structure

### Documentation (this feature)

```text
specs/031-sdk-docs-consolidation/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Documentation locations (primary targets)
sdk/go/pluginsdk/
├── README.md            # Main SDK documentation (major updates)
├── client.go            # Inline comment additions
├── focus_builder.go     # Inline comment additions
├── focus_conformance.go # Inline comment additions
└── example_test.go      # New godoc example file

sdk/go/testing/
├── README.md            # Testing framework docs (import fixes)
└── *_test.go            # Test naming clarifications
```

**Structure Decision**: Single project structure - all changes are inline documentation
improvements to existing SDK packages. No new packages or major restructuring required.

## Complexity Tracking

> No constitutional violations requiring justification. This is a documentation-only feature
> that directly supports Constitution Principle V (Comprehensive Documentation).

## Research Required (Phase 0)

The following items require investigation before implementation:

1. **Connect-go Migration Patterns** - Research best practices for gRPC to connect-go migration
2. **CORS Best Practices** - Research standard CORS configurations for 5+ deployment scenarios
3. **Rate Limiting Patterns** - Research golang.org/x/time/rate token bucket implementation
4. **Thread Safety Documentation** - Audit existing SDK components for concurrency guarantees
5. **HTTP Client Ownership Patterns** - Research standard Go idioms for client ownership semantics
