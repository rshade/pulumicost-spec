# Implementation Plan: "What-If" Cost Estimation API

**Branch**: `006-estimate-cost` | **Date**: 2025-11-24 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-estimate-cost/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.
See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add a new EstimateCost RPC method to the CostSource gRPC service that enables proactive cost
estimation for Pulumi resources before deployment. The API accepts a resource type string and
structured attributes, returning an estimated monthly cost. This enables developers to compare
different resource configurations during development and make cost-informed infrastructure
decisions.

## Technical Context

**Language/Version**: Go 1.24+ (matches existing SDK, toolchain go1.25.4), Protocol Buffers v3
**Primary Dependencies**:

- google.golang.org/grpc (gRPC Go implementation)
- google.golang.org/protobuf (Go protobuf runtime)
- buf v1.32.1 (Protocol buffer toolchain)
- github.com/rs/zerolog v1.34.0+ (structured logging per spec 005-zerolog)

**Storage**: N/A (specification repository - no runtime storage)
**Testing**:

- Go testing framework (sdk/go/testing/ harness)
- bufconn for in-memory gRPC testing
- Conformance tests (Basic/Standard/Advanced levels)
- Performance benchmarks with memory profiling

**Target Platform**: Cross-platform (specification defines gRPC protocol)
**Project Type**: gRPC Protocol Specification (protobuf-first)
**Performance Goals**:

- RPC response time <500ms (per SC-002)
- Support concurrent requests per conformance level
- Zero-allocation enum validation patterns where applicable

**Constraints**:

- Backward compatibility with existing CostSource service
- buf breaking change detection must pass
- No retry logic in SDK (per FR-014, handled by plugins/core)
- Resource type format validation: provider:module/resource:Type

**Scale/Scope**:

- Single new RPC method (EstimateCost)
- 2 new protobuf messages (EstimateCostRequest, EstimateCostResponse)
- Support for 3+ major cloud providers (AWS, Azure, GCP per SC-003)
- Integration with existing testing framework (sdk/go/testing/)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

### I. gRPC Proto Specification-First Development ✅ PASS

- ✅ Proto definitions are source of truth: New EstimateCost RPC will be defined in
  `proto/finfocus/v1/costsource.proto`
- ✅ SDK code generated from proto: Go SDK will be regenerated via buf after proto updates
- ✅ Comprehensive validation: EstimateCostRequest and EstimateCostResponse will have full validation
- ✅ Breaking changes detected: buf lint/breaking checks in CI
- ✅ Proper gRPC status codes: Spec defines InvalidArgument, NotFound, Unavailable errors

### II. Multi-Provider gRPC Consistency ✅ PASS

- ✅ Cross-provider support: SC-003 requires AWS, Azure, GCP support
- ✅ Provider-agnostic messages: EstimateCostRequest uses generic resource_type and attributes fields
- ✅ No provider-specific extensions: Uses google.protobuf.Struct for flexible attribute handling
- ⚠️ **ACTION REQUIRED**: Phase 1 must create cross-provider examples for EstimateCost requests

### III. Test-First Protocol (NON-NEGOTIABLE) ✅ PASS

- ✅ TDD required: Phase 0 research will define conformance test approach
- ✅ Red-Green-Refactor: Tests written before proto changes
- ✅ gRPC error conditions: Spec defines error scenarios (FR-003, FR-008, FR-009, FR-010, FR-014)
- ✅ Integration with existing testing framework: sdk/go/testing/ harness supports new RPC

### IV. Protobuf Backward Compatibility ✅ PASS

- ✅ Non-breaking change: Adding new RPC method is backward compatible
- ✅ buf breaking check will pass: Existing RPCs unchanged
- ✅ Field numbering: Phase 1 will assign proper field numbers (1-15 for frequent fields)
- ✅ UnimplementedServer: Existing UnimplementedCostSourceServiceServer handles new method

### V. Comprehensive Documentation ✅ PASS

- ✅ Proto comments required: Phase 1 will add inline documentation
- ✅ RPC contract documentation: Request/response semantics defined in spec
- ⚠️ **ACTION REQUIRED**: Phase 1 must create example gRPC request/response payloads
- ⚠️ **ACTION REQUIRED**: Phase 1 must update README with EstimateCost usage

### VI. Performance as a gRPC Requirement ✅ PASS

- ✅ Performance target defined: SC-002 specifies <500ms response time
- ✅ Conformance testing: Integration with existing Basic/Standard/Advanced levels
- ✅ Benchmark requirements: Phase 1 will add EstimateCost benchmarks
- ✅ Concurrent request support: Spec defines deterministic behavior (FR-011)

### VII. Validation at Multiple Levels ✅ PASS

- ✅ Protobuf layer: buf validates proto syntax and breaking changes
- ✅ Data layer: Validation in FR-003 (format), FR-005 (null handling), FR-008-010 (attributes)
- ✅ Service layer: Conformance tests will validate EstimateCost behavior
- ✅ SDK layer: Integration tests via bufconn harness
- ✅ CI layer: All validation runs in GitHub Actions

**GATE STATUS**: ✅ **PASS** - All constitutional requirements met. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
└── costsource.proto          # Add EstimateCost RPC and messages here

sdk/go/
├── proto/                     # Generated code (regenerated via buf)
│   └── finfocus/v1/
│       ├── costsource.pb.go
│       └── costsource_grpc.pb.go
├── pricing/                   # Domain types and validation
│   ├── domain.go
│   └── validate.go
└── testing/                   # Testing framework
    ├── harness.go             # Add EstimateCost test support
    ├── mock_plugin.go         # Add EstimateCost mock methods
    ├── conformance_test.go    # Add EstimateCost conformance tests
    └── benchmark_test.go      # Add EstimateCost benchmarks

examples/
└── requests/                  # New directory for RPC examples
    ├── estimate_cost_aws.json
    ├── estimate_cost_azure.json
    └── estimate_cost_gcp.json
```

**Structure Decision**: This is a gRPC protocol specification repository. Changes are limited
to proto definitions, generated SDK code, testing framework extensions, and example payloads.
No application code or services are implemented here - those are in plugin repositories that
implement the CostSource service interface.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A - All constitutional requirements passed. No violations to justify.

## Phase Completion Summary

### Phase 0: Research ✅ COMPLETE

**Output**: `research.md`

**Key Decisions**:

1. Use `google.protobuf.Struct` for attributes (consistent with existing RPCs)
2. Resource type validation in SDK layer (provider:module/resource:Type pattern)
3. Clear gRPC status code mapping for all error scenarios
4. Leverage zerolog integration for observability
5. Extend existing conformance framework with 3-level testing
6. Cross-provider examples for AWS, Azure, GCP

### Phase 1: Design & Contracts ✅ COMPLETE

**Outputs**:

- `data-model.md` - Protobuf message definitions and validation rules
- `contracts/estimate_cost.proto` - RPC contract definition with inline documentation
- `contracts/examples.md` - Request/response examples for all scenarios
- `quickstart.md` - Developer guide with code examples and best practices
- `CLAUDE.md` - Updated agent context with new technologies

**Artifacts Created**:

- EstimateCostRequest protobuf message specification
- EstimateCostResponse protobuf message specification
- Resource type format validation pattern
- Error handling patterns with gRPC status codes
- Cross-provider request examples (AWS, Azure, GCP)
- Error scenario examples
- Configuration comparison use cases

**Constitution Re-Check** ✅ PASS:

- All Phase 1 artifacts align with constitutional principles
- Cross-provider examples created as required
- Proto contract defined with comprehensive documentation
- Testing framework integration planned
- Performance requirements maintained

## Next Steps

**Ready for** `/speckit.tasks`:

- Convert functional requirements into actionable tasks
- Generate dependency-ordered implementation checklist
- Create test-first task sequence per TDD principle

**Implementation Order** (from tasks.md):

1. Write conformance tests for EstimateCost (RED phase)
2. Update proto/finfocus/v1/costsource.proto
3. Run `make generate` to regenerate Go SDK
4. Implement validation functions in sdk/go/pricing/
5. Extend testing framework in sdk/go/testing/
6. Create cross-provider examples
7. Run tests to verify (GREEN phase)
8. Update documentation
