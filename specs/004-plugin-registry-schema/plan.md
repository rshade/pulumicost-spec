# Implementation Plan: Plugin Registry Index JSON Schema

**Branch**: `004-plugin-registry-schema` | **Date**: 2025-11-23 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-plugin-registry-schema/spec.md`

## Summary

Add a JSON Schema (`schemas/plugin_registry.schema.json`) for validating plugin registry
index files (`registry.json`). The schema defines the contract for well-known plugin
registries used by `finfocus plugin install`, aligning field patterns and enums with
`registry.proto` definitions.

## Technical Context

**Language/Version**: JSON Schema draft 2020-12
**Primary Dependencies**: AJV (validation), existing `registry.proto` definitions
**Storage**: N/A (schema file in repository)
**Testing**: npm validate scripts (AJV), example registry.json validation
**Target Platform**: Any JSON Schema validator
**Project Type**: single (schema + examples)
**Performance Goals**: N/A (schema validation is inherently fast)
**Constraints**: Must align with registry.proto enums (SecurityLevel, capabilities)
**Scale/Scope**: Single schema file, 2 example plugins, npm validation integration

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                           | Status   | Notes                                         |
| ----------------------------------- | -------- | --------------------------------------------- |
| I. Proto Specification-First        | **PASS** | Schema aligns with existing registry.proto    |
| II. Multi-Provider Consistency      | **PASS** | Examples include AWS, Kubernetes plugins      |
| III. Test-First Protocol            | **PASS** | Schema validation tests before implementation |
| IV. Protobuf Backward Compatibility | **N/A**  | No proto changes (schema only)                |
| V. Comprehensive Documentation      | **PASS** | Schema has descriptions, examples provided    |
| VI. Performance as gRPC Requirement | **N/A**  | Not a gRPC feature                            |
| VII. Validation at Multiple Levels  | **PASS** | JSON Schema validates registry entries        |

**Gate Result**: PASS - All applicable principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/004-plugin-registry-schema/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # N/A (not an API feature)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
schemas/
├── plugin_registry.schema.json    # NEW: Main schema file
├── plugin_manifest.schema.json    # Existing
└── pricing_spec.schema.json       # Existing

examples/
├── registry.json                  # NEW: Example registry index
├── specs/                         # Existing pricing spec examples
└── README.md                      # Update with registry info

scripts/
└── validate_examples.js           # Update to validate registry
```

**Structure Decision**: Schema added to existing `schemas/` directory following established
patterns. Example registry placed in `examples/` at root level (not in `specs/` subdirectory)
since it represents the registry index format, not a pricing spec.

## Complexity Tracking

> No violations - schema feature follows existing patterns
