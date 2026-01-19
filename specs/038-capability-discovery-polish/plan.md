# Implementation Plan - Capability Discovery Polish

**Feature**: Capability Discovery Polish
**Branch**: `038-capability-discovery-polish`
**Status**: Planning
**Spec**: [specs/038-capability-discovery-polish/spec.md](spec.md)

## Technical Context

This feature consolidates several quality improvements and polish items for the plugin capability
discovery system. The primary goal is to standardize how capabilities are advertised (enum vs.
legacy metadata), implement auto-discovery based on interfaces, and improve SDK performance.

### Architecture

- **Capability Discovery**: A hybrid approach where plugins can either manually declare
  capabilities OR have the SDK auto-discover them by checking implemented interfaces
  (`DryRunHandler`, etc.).
- **Backward Compatibility**: Dual-reporting of capabilities as both a modern Enum list
  (`capabilities`) and a legacy string map (`metadata`).
- **Performance**: Optimization of slice copying in the SDK hot paths.

### Existing Components

- `sdk/go/pluginsdk`: The core SDK package where `PluginInfo`, capability detection, and helper functions reside.
- `proto/finfocus/v1/costsource.proto`: Defines the `GetPluginInfoResponse` message.
- `sdk/go/testing`: Contains compatibility tests that need renaming.

### Unknowns & Risks

- **Unknowns**: None. The spec is extremely detailed and verification steps have confirmed the
  existence of key files.
- **Risks**:
  - **Proto Compatibility**: Modifying proto comments is safe, but we must ensure we don't
    accidentally rename fields.
  - **Logic Duplication**: We must be careful to remove the old capability conversion logic when
    introducing the new central helper to avoid "split brain" behavior.

## Constitution Check

| Principle                              | Compliance | Notes                                                                                           |
| :------------------------------------- | :--------- | :---------------------------------------------------------------------------------------------- |
| **I. Proto Contracts are Sacred**      | ✅         | We are clarifying comments in Proto files, not changing wire format.                            |
| **II. Multi-Provider Consistency**     | ✅         | Capability discovery is provider-agnostic.                                                      |
| **III. Spec Consumes, Not Calculates** | ✅         | N/A - Infrastructure feature.                                                                   |
| **IV. Separation of Concerns**         | ✅         | Improvements are isolated to the SDK/Proto layer.                                               |
| **V. Test-First Protocol**             | ✅         | Unit tests and benchmarks are required by the spec.                                             |
| **VI. Proto Backward Compatibility**   | ✅         | No breaking changes; clarifying deprecations.                                                   |
| **VII. Documentation & Identity**      | ✅         | Docs are a major part of this feature (Issue #300).                                             |
| **VIII. Performance**                  | ✅         | Explicit performance optimization (Issue #301).                                                 |
| **IX. Observability**                  | ✅         | N/A                                                                                             |
| **XIII. Multi-Language SDK Sync**      | ✅         | We will add tasks to verify/update TS SDK if relevant (though this is mostly Go SDK internals). |

## Phased Implementation

### Phase 0: Research & Design

**Goal**: Confirm design of helper functions and finalize proto comment updates.

- [x] Verify existing `DryRunHandler` interface (Done in planning).
- [ ] Create `research.md` (Formalize the "no unknowns" status).

### Phase 1: Foundation & Contracts (Issues #294, #295, #208, #209)

**Goal**: Fix "easy" issues: Proto comments, test renaming, and doc fixes.

1. **Proto Updates**: Add clarifying comments to `costsource.proto`.
2. **Close #294**: Verify `DryRunHandler` works as expected (it exists).
3. **Test Renaming**: Rename compatibility tests in `sdk/go/testing`.
4. **Doc Fixes**: Update `sdk/go/testing/README.md` imports.

### Phase 2: Core SDK Logic (Issues #299, #301)

**Goal**: Implement the new capability conversion logic and performance fixes.

1. **Capability Helper**: Create `sdk/go/pluginsdk/capability_compat.go`.
   - Implement `legacyCapabilityNames` map.
   - Implement `CapabilitiesToLegacyMetadata` and `...WithWarnings`.
2. **Refactor SDK**: Update `NewPluginInfo` and `Supports` to use the new helper.
3. **Performance**: Optimize slice copying in `sdk.go`.
4. **Benchmarks**: Run benchmarks to prove performance gains.

### Phase 3: Documentation & Verification (Issue #300)

**Goal**: Finalize documentation and verify everything.

1. **CLAUDE.md**: Add "Capability Discovery Pattern" section.
2. **SDK README**: Update `pluginsdk/README.md` with auto-discovery examples.
3. **Validation**: Run `make validate` and `make lint`.

## Verification Plan

### Automated Tests

- `go test ./sdk/go/pluginsdk/...` - Verify capability detection and helper logic.
- `go test ./sdk/go/testing/...` - Verify renamed compatibility tests pass.
- `go test -bench=. ./sdk/go/pluginsdk/...` - Verify slice copying performance.

### Manual Verification

- Inspect generated `CLAUDE.md` and `README.md` to ensure clarity.
- Verify `make lint` passes.
