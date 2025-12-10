# Implementation Plan: GetBudgets RPC for Plugin-Provided Budget Information

**Branch**: `001-get-budgets-rpc` | **Date**: 2025-12-09 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-get-budgets-rpc/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.
See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add a GetBudgets RPC to the CostSource gRPC service to allow plugins to provide budget
information from cloud cost management services. Implementation follows proto-first development
with generated Go SDK, comprehensive testing, and cross-provider examples.

## Technical Context

**Language/Version**: Go 1.24.10 (toolchain 1.25.4)
**Primary Dependencies**: gRPC, protobuf, buf v1.32.1
**Storage**: N/A (protocol specification repository)
**Testing**: Go testing framework, gRPC conformance tests (Basic/Standard/Advanced levels)
**Target Platform**: Cross-platform (gRPC clients for AWS, Azure, GCP, Kubernetes plugins)
**Project Type**: Protocol specification library with generated SDK
**Performance Goals**: gRPC response time <5 seconds for budget queries, support 100-1000 budgets per department
**Constraints**: Protobuf backward compatibility, buf validation passing, no breaking changes without MAJOR version bump
**Scale/Scope**: Multi-provider gRPC service supporting 4+ cloud providers with unified budget interface

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

### I. gRPC Proto Specification-First Development

✅ **PASS**: Feature begins with protobuf specification updates to `costsource.proto` before implementation.

### II. Multi-Provider gRPC Consistency

✅ **PASS**: Budget RPC supports all major providers (AWS, Azure, GCP, Kubernetes) with provider-agnostic message fields.

### III. Test-First Protocol

✅ **PASS**: Conformance tests will be written defining expected gRPC behavior before proto implementation.

### IV. Protobuf Backward Compatibility

✅ **PASS**: New optional RPC method (GetBudgets) maintains backward compatibility - plugins can return Unimplemented.

### V. Comprehensive Documentation

✅ **PASS**: Proto messages and RPC methods will include comprehensive inline comments for generated documentation.

### VI. Performance as a gRPC Requirement

✅ **PASS**: Performance goals defined (<5s response time, 100-1000 budget scale) with conformance test requirements.

### VII. Validation at Multiple Levels

✅ **PASS**: Implementation will include buf validation, JSON schema validation, and gRPC conformance testing.

**Overall**: All constitutional gates PASS. Feature aligns with gRPC proto-first development principles.

### Post-Design Re-evaluation

✅ **CONFIRMED**: Design phase complete. All constitutional requirements met:

- Proto definitions designed with cross-provider consistency
- Test-first approach confirmed with conformance test plan
- Backward compatibility maintained through optional RPC
- Comprehensive documentation planned with examples
- Performance requirements specified and testable
- Multi-layer validation approach confirmed (buf, schema, conformance)

## Project Structure

### Documentation (this feature)

```text
specs/001-get-budgets-rpc/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── README.md
│   ├── get-budgets-request-aws.json
│   ├── get-budgets-response-aws.json
│   ├── get-budgets-request-gcp.json
│   ├── get-budgets-response-gcp.json
│   ├── get-budgets-request-kubecost.json
│   └── get-budgets-response-kubecost.json
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
proto/pulumicost/v1/
├── budget.proto         # NEW: Budget message definitions
└── costsource.proto     # UPDATED: Add GetBudgets RPC

sdk/go/
├── proto/               # GENERATED: gRPC service and message types
│   ├── pulumicost/
│   │   └── v1/
│   │       ├── budget.pb.go
│   │       └── costsource.pb.go
│   └── pulumicost/
│       └── v1/
│           ├── budget_grpc.pb.go
│           └── costsource_grpc.pb.go
├── pluginsdk/           # UPDATED: Add BudgetsProvider interface
│   └── sdk.go
└── testing/             # UPDATED: Add budget validation functions
    ├── harness.go
    └── mock_plugin.go

examples/
├── specs/               # UPDATED: Add budget examples
│   ├── aws-budget.json
│   ├── gcp-budget.json
│   └── kubecost-budget.json
└── requests/            # UPDATED: Add GetBudgets request examples
    ├── get-budgets-aws.json
    ├── get-budgets-gcp.json
    └── get-budgets-kubecost.json

schemas/                 # UPDATED: Add budget schema
└── budget_spec.schema.json
```

**Structure Decision**: Protocol specification repository structure maintained. New budget
functionality adds proto definitions, generated SDK code, examples, and schemas while preserving
existing patterns. No architectural changes required - follows established gRPC proto-first
development workflow.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |
