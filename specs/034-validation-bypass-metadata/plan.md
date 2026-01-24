# Implementation Plan: Validation Bypass Metadata

**Branch**: `034-validation-bypass-metadata` | **Date**: 2026-01-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/034-validation-bypass-metadata/spec.md`

## Summary

Add fields to `ValidationResult` in `sdk/go/pricing/observability.go` to carry metadata about why
validation policies were bypassed (e.g., `--yolo` flag). This enables audit trails to survive
stateless service boundaries by including timestamp, reason, operator, severity, and mechanism
type in a structured `BypassMetadata` slice.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod)
**Primary Dependencies**: Standard library only (`time`, `encoding/json`)
**Storage**: N/A (stateless struct extension; retention is caller's responsibility)
**Testing**: Go standard testing (`go test`), table-driven tests, benchmarks
**Target Platform**: Linux/macOS/Windows (cross-platform Go SDK)
**Project Type**: Single project (Go SDK package extension)
**Performance Goals**: Zero-allocation for common paths, <10 ns/op for enum validation
**Constraints**: Backward compatible with existing ValidationResult consumers
**Scale/Scope**: SDK-level change affecting `sdk/go/pricing/` package

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First**: This change extends Go SDK types, not proto definitions. The
  `ValidationResult` struct is a Go-only type used for observability, not a gRPC message.
  Proto changes are NOT required.
- [x] **Spec Consumes**: No complex pricing logic. Simple metadata recording for audit purposes.
- [x] **Multi-Provider**: Provider-agnostic design. Bypass metadata applies to any validation
  regardless of cloud provider.
- [x] **FinFocus Alignment**: Uses `finfocus` naming conventions throughout.
- [x] **SDK Synchronization**: Go SDK only change. TypeScript SDK not affected as this is an
  observability helper type, not a gRPC-generated type.

**Constitution Status**: ✅ All gates pass. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/034-validation-bypass-metadata/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A for this feature)
├── checklists/          # Quality checklists
│   └── requirements.md
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
sdk/go/pricing/
├── observability.go     # Extend ValidationResult, add BypassMetadata types
├── observability_test.go # Add tests for bypass metadata functionality
├── bypass.go            # New file: BypassMetadata, BypassSeverity, BypassMechanism
└── bypass_test.go       # New file: Tests and benchmarks for bypass types

sdk/go/testing/
└── bypass_conformance_test.go  # Conformance tests for bypass metadata serialization
```

**Structure Decision**: Extend existing `sdk/go/pricing/` package with new `bypass.go` file for
bypass-specific types, keeping `observability.go` focused on its current responsibilities while
adding the `Bypasses` field to `ValidationResult`.

## Complexity Tracking

> No constitution violations. This is a straightforward SDK extension.

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| N/A       | N/A        | N/A                                  |
