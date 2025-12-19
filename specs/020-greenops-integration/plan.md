# Implementation Plan: GreenOps Integration

**Feature Branch**: `020-greenops-integration`  
**Created**: 2025-12-18  
**Status**: Planning  
**Feature Spec**: [specs/020-greenops-integration/spec.md](spec.md)

## Technical Context

### Existing Infrastructure

- **Proto**: `proto/pulumicost/v1/costsource.proto` defines the core gRPC service and messages.
- **SDK**: Go SDK is generated from the proto definitions.
- **Validation**: `buf` is used for linting and breaking change detection.

### Unknowns & Research

- **MetricKind Location**: Decided to place it in `costsource.proto` as a top-level enum.
- **Utilization Scope**: Implementation will follow a global default with per-resource override model.
- **Units**: Standardized on gCO2e, kWh, and Liters.

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Spec-First | ✅ | Proto changes planned before implementation. |
| II. Multi-Provider | ✅ | Standardized metrics (Carbon, Energy, Water) are provider-agnostic. |
| III. Test-First | ✅ | Conformance tests will be updated/added. |
| IV. Backward Compatibility | ✅ | New fields added with appropriate numbering; no breaking changes. |
| V. Documentation | ✅ | Proto comments and `quickstart.md` included. |
| VI. Performance | ✅ | Metric list and double fields are lightweight. |
| VII. Multi-Layer Validation | ✅ | `buf lint` will be run. |

## Implementation Phases

### Phase 0: Research (Completed)

- Artifact: [research.md](research.md)
- Decisions on enum values, units, and field scope are finalized.

### Phase 1: Design & Contracts (In Progress)

- [x] Data Model: [data-model.md](data-model.md)
- [x] Quickstart: [quickstart.md](quickstart.md)
- [ ] Protobuf Updates: Modify `proto/pulumicost/v1/costsource.proto`
- [ ] SDK Generation: Run `make generate`

### Phase 2: Testing & Validation

- [ ] Conformance Tests: Update `sdk/go/conformance` to validate `SupportsResponse` metrics.
- [ ] Integration Tests: Add test cases for `utilization_percentage` override logic.
- [ ] Examples: Add a new GreenOps example in `examples/plugins/greenops-plugin.json` (mock).

## Gates & Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Protobuf Versioning | Low | Use `buf` to ensure no breaking changes occur. |
| Plugin Adoption | Medium | Provide clear documentation and a `quickstart.md` guide. |
| Unit Mismatch | Low | Explicitly document units (gCO2e, kWh, L) in proto comments. |
