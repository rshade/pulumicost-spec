# Implementation Plan: Domain Enum Validation Performance Optimization

**Branch**: `001-domain-enum-optimization` | **Date**: 2025-11-17 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-domain-enum-optimization/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for
the execution workflow.

## Summary

Optimize registry package domain enum validation functions by replacing linear search through slices with map-based
lookups for all enum types (Provider, DiscoverySource, PluginStatus, SecurityLevel, InstallationMethod,
PluginCapability, SystemPermission, AuthMethod). The optimization aims to improve validation performance while
maintaining 100% backward compatibility and establishing consistent patterns across registry and pricing packages.

## Technical Context

**Language/Version**: Go 1.24.10 (toolchain go1.25.4)
**Primary Dependencies**: Standard library only (no external dependencies for validation)
**Storage**: N/A
**Testing**: Go standard testing package with benchmarks (`testing.B`)
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: SDK library (single project)
**Performance Goals**: Validation operations complete in under 100 nanoseconds per operation for enums with up to
50 values
**Constraints**: Must maintain backward compatibility with existing function signatures and behavior; must support
8 enum types across registry package
**Scale/Scope**: 8 enum types with 4-14 values each (Provider: 5, DiscoverySource: 4, PluginStatus: 6, SecurityLevel:
4, InstallationMethod: 4, PluginCapability: 14, SystemPermission: 9, AuthMethod: 6)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

### I. gRPC Proto Specification-First Development

**Status**: ✅ PASS - Not applicable
**Rationale**: This feature optimizes Go SDK validation functions, not gRPC protocol definitions. No proto changes
required.

### II. Multi-Provider gRPC Consistency

**Status**: ✅ PASS - Not applicable
**Rationale**: Provider enum validation is internal SDK functionality. Validation pattern consistency benefits all
providers equally.

### III. Test-First Protocol (NON-NEGOTIABLE)

**Status**: ✅ PASS - Test-first approach required
**Approach**:

1. Write performance benchmark tests comparing slice-based vs map-based validation
2. Write unit tests for all 8 enum validation functions with valid/invalid inputs
3. Tests must fail initially (benchmarks will show performance difference)
4. Implement map-based validation to make tests pass
5. Verify benchmarks show performance improvement and unit tests maintain 100% accuracy

**Gate**: TDD workflow mandatory. Benchmarks and unit tests must be written before implementation.

### IV. Protobuf Backward Compatibility

**Status**: ✅ PASS - Not applicable
**Rationale**: No protobuf changes. SDK-only optimization maintaining existing function signatures.

### V. Comprehensive Documentation

**Status**: ⚠️ CONDITIONAL PASS - Documentation updates required
**Requirements**:

- Function comments remain accurate (no changes needed - behavior unchanged)
- Benchmark results documented in commit message or PR description
- Performance characteristics compared (before/after) in documentation
- Pattern consistency between registry and pricing packages noted in CLAUDE.md

**Gate**: Performance comparison documentation required before merge.

### VI. Performance as a gRPC Requirement

**Status**: ✅ PASS - Performance measurement required
**Approach**:

- Benchmark tests for all validation functions (IsValidProvider, IsValidDiscoverySource, etc.)
- Compare slice-based vs map-based performance across different enum sizes
- Measure memory allocations per operation
- Document scalability characteristics (5, 10, 50 value enums)
- Target: < 100ns per operation for enums up to 50 values

**Gate**: Benchmark tests must demonstrate measurable improvement.

### VII. Validation at Multiple Levels

**Status**: ✅ PASS - Multi-layer testing required
**Validation Layers**:

- **Unit tests**: 100% accuracy for valid/invalid enum values (all 8 types)
- **Benchmark tests**: Performance comparison (slice vs map) with memory profiling
- **Integration tests**: Existing tests continue to pass (no behavior changes)
- **CI validation**: All existing CI checks pass (golangci-lint, buf, test suite)

**Gate**: All validation layers must pass before merge.

## Constitution Summary

**Overall Status**: ✅ CONDITIONAL PASS

**Action Items Before Implementation**:

1. Write benchmark tests first (TDD requirement)
2. Write unit tests for all 8 enum types (TDD requirement)
3. Document performance comparison after implementation
4. Verify pattern consistency with pricing package

**No gate violations requiring complexity justification.**

## Project Structure

### Documentation (this feature)

```text
specs/001-domain-enum-optimization/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
sdk/go/registry/
├── domain.go            # Existing: 8 enum types with linear search validation
├── domain_test.go       # Existing: Unit tests for enum validation
├── validate.go          # Existing: Plugin name validation (not affected)
└── validate_test.go     # Existing: Plugin name tests (not affected)

sdk/go/pricing/
├── domain.go            # Reference: BillingMode validation pattern (38+ enums)
├── domain_test.go       # Reference: Validation test patterns
└── validate.go          # Reference: Schema validation (not affected)

sdk/go/testing/
├── benchmark_test.go    # Existing: Performance benchmark patterns
└── conformance_test.go  # Existing: Conformance test patterns

tests/
├── unit/                # New: Benchmark tests for validation performance
└── integration/         # Existing: Integration tests (should pass unchanged)
```

**Structure Decision**: Single project SDK structure with modifications isolated to `sdk/go/registry/domain.go`. New
benchmark tests in `sdk/go/registry/domain_test.go`. Reference patterns from `sdk/go/pricing/` for consistency.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

_No constitution violations - this section is empty._
