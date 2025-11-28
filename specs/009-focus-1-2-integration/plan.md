# Implementation Plan - FinOps FOCUS 1.2 Integration

**Feature**: `focus-1-2-integration`  
**Version**: 1.0.0  
**Status**: Planned  
**Spec**: [specs/009-focus-1-2-integration/spec.md](./spec.md)

## Technical Context

This feature upgrades the `pulumicost-spec` to align with FinOps FOCUS 1.2 by introducing
a new `FocusCostRecord` Protobuf message, standardizing categorization via strict Enums
(Vocabulary), and implementing a "Backpack & Builder" pattern in the Go SDK.

**Key constraints:**

- **Zero-Allocation Goal**: The SDK Builder must be highly performant.
- **Backpack Pattern**: Must support arbitrary extension data without schema changes.
- **Strict Enums**: `ServiceCategory`, `ChargeCategory`, `PricingCategory` must be defined in Proto and enforced in Go.
- **Protobuf Double**: Financial fields must use `double` as per research decision.
- **Validation**: `Build()` returns error on missing mandatory fields (FOCUS 1.2 business rules).

**Unknowns & Risks:**

- [x] **Validation Rules**: Resolved. Mandatory fields identified: `BillingAccountId`,
  `ChargePeriod`, `ServiceCategory`, `ChargeCategory`, `BilledCost`, `Currency`.
- [x] **Enum Values**: Resolved. Mapped from FOCUS 1.0/1.1 and search results.
- [x] **Versioning Strategy**: Resolved. Use `extended_columns` for future fields.
- [x] **Breaking Changes**: Resolved. Additive change, safe.

## Constitution Check

| Principle | Status | Notes |
| :--- | :--- | :--- |
| **I. Proto Specification-First** | ✅ | Plan starts with `focus.proto` and `enums.proto` definition. |
| **II. Multi-Provider Consistency** | ✅ | FOCUS 1.2 is inherently multi-provider; Enums will cover all major clouds. |
| **III. Test-First Protocol** | ✅ | Tasks include "Create Conformance Tests" before Implementation. |
| **IV. Backward Compatibility** | ✅ | New message type `FocusCostRecord` is additive; existing messages unchanged. |
| **V. Comprehensive Documentation** | ✅ | Plan includes `PLUGIN_MIGRATION_GUIDE.md` and proto comments. |
| **VI. Performance Requirement** | ✅ | "Backpack" implementation must be efficient (map handling); benchmarks planned. |
| **VII. Multi-Level Validation** | ✅ | `Build()` validation + Conformance Test Suite planned. |

## Phase 0: Research & Design

**Goal**: Resolve unknowns and finalize the data model.

1. **Research FOCUS 1.2 Spec**:
    - Action: Extract the full list of mandatory vs optional columns.
    - Action: Extract the full list of allowed values for Service, Charge, and Pricing categories.
    - Output: `specs/009-focus-1-2-integration/research.md`
2. **Define Data Model**:
    - Action: Map FOCUS columns to Protobuf fields (types, field IDs).
    - Action: Define Enum values in `enums.proto`.
    - Output: `specs/009-focus-1-2-integration/data-model.md`
3. **Define Contracts**:
    - Action: Draft `proto/pulumicost/v1/focus.proto` and `proto/pulumicost/v1/enums.proto`.
    - Output: `specs/009-focus-1-2-integration/contracts/focus.proto`

## Phase 1: Specification & SDK Core

**Goal**: Implement the Protobuf spec and base SDK Builder.

1. **Protobuf Implementation**:
    - Create `proto/pulumicost/v1/enums.proto` (Vocabularies).
    - Create `proto/pulumicost/v1/focus.proto` (`FocusCostRecord` message).
    - Run `buf generate` to create Go code.
2. **SDK Builder Implementation**:
    - Create `sdk/go/pluginsdk/focus_builder.go`.
    - Implement `With...` methods for all fields.
    - Implement `WithExtension` (Backpack).
    - Implement `Build()` with validation logic (mandatory fields check).
3. **Conformance Tests**:
    - Create `sdk/go/pluginsdk/focus_builder_test.go`.
    - Test: Happy path (all fields).
    - Test: Missing mandatory field -> Error.
    - Test: Enum enforcement.
    - Test: Extension data preservation.

## Phase 2: Validation & Documentation

**Goal**: Verify compliance and document usage.

1. **Conformance Validator**:
    - Create `sdk/go/pluginsdk/focus_conformance.go` (Exported validation function).
    - Implement strict FOCUS 1.2 business rule checks (e.g., if ChargeCategory=Usage, UsageQuantity must be > 0).
2. **Benchmarks**:
    - Create `sdk/go/pluginsdk/focus_benchmark_test.go`.
    - Benchmark: Builder allocation and serialization speed.
3. **Documentation**:
    - Create `docs/PLUGIN_MIGRATION_GUIDE.md` (Focus Upgrade).
    - Update `README.md` with FOCUS 1.2 badge/info.
    - Create `examples/plugins/focus-example.go` demonstrating Builder usage.
