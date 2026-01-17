# Actionable Tasks: Dual-Layer Capability Discovery

**Feature Branch**: `037-capability-discovery`
**Status**: Planned
**Spec**: [spec.md](spec.md)

## Implementation Strategy

- **Phase 1 (Setup)**: Prepare Protocol Buffers and generate Go code. This unlocks the data model
  for both SDK and tests.
- **Phase 2 (Foundational)**: Implement core SDK logic for Global discovery (interface reflection).
  This is the base for granular checks.
- **Phase 3 (User Story 1)**: Implement Global Discovery end-to-end, including `GetPluginInfo`
  updates.
- **Phase 4 (User Story 2)**: Implement Granular Discovery (`Supports`) with `ResourceDescriptor`
  input and "Inherit Global" fallback logic.
- **Phase 5 (Polish)**: Ensure strict error handling (FR-007) and comprehensive documentation.

## Dependencies

- Phase 1 blocks all other phases.
- Phase 2 blocks Phase 3 and 4.
- Phase 3 blocks Phase 4 (Granular needs Global context for inheritance).

## Phase 1: Setup & Protocol Definition

**Goal**: Define the data model in Protobuf and generate code bindings.

- [x] T001 Update `PluginCapability` enum in `proto/finfocus/v1/enums.proto` with all required
  values (ACTUAL, PROJECTED, ESTIMATE, RECOMMENDATIONS, BUDGETS, DISMISS, DRY_RUN)
- [x] T002 Update `SupportsResponse` in `proto/finfocus/v1/costsource.proto` to include
  `repeated PluginCapability capabilities_enum = 5`
- [x] T003 Verify `SupportsRequest` in `proto/finfocus/v1/costsource.proto` already has
  `ResourceDescriptor` (Field 1) — no action if confirmed, else update
- [x] T004 Run `buf generate` to regenerate Go SDK code from updated protos

## Phase 2: Foundational SDK Logic

**Goal**: Enable the SDK to "know itself" via reflection (Auto-Discovery).

- [x] T005 [P] Implement `inferCapabilities` helper function in `sdk/go/pluginsdk/server.go` to
  reflect on `Plugin` interface and return `[]PluginCapability`
- [x] T006 [P] Update `Server` struct in `sdk/go/pluginsdk/server.go` to store
  `globalCapabilities []PluginCapability` computed at initialization
- [x] T007 [P] Create `capabilitiesToLegacyMetadata` helper to convert `[]PluginCapability` ->
  `map[string]bool` for backward compatibility

## Phase 3: Global Discovery (User Story 1)

**Goal**: Host can identify supported features at initialization.

- [x] T008 [US1] Update `GetPluginInfo` implementation in `sdk/go/pluginsdk/server.go` to include
  `capabilities` (enum list) in the response
- [x] T009 [US1] Ensure `GetPluginInfo` populates the legacy `metadata` map using
  `capabilitiesToLegacyMetadata`
- [x] T010 [US1] Create unit test in `sdk/go/pluginsdk/plugin_info_test.go` verifying `GetPluginInfo`
  returns correct capabilities for a mock plugin implementing specific interfaces

## Phase 4: Granular Discovery (User Story 2)

**Goal**: Host can check support per resource with "Inherit Global" fallback.

- [x] T011 [US2] Update `Supports` method signature in `sdk/go/pluginsdk/server.go` to handle
  `ResourceDescriptor`
- [x] T012 [US2] Implement "Inherit Global" logic in `Supports`: if plugin returns empty
  `CapabilitiesEnum`, copy from `Server.globalCapabilities`
- [x] T013 [US2] Implement Legacy Sync in `Supports`: populate `Capabilities` (string map) from
  `CapabilitiesEnum` before returning
- [x] T014 [US2] Create conformance test in `sdk/go/testing/conformance_test.go` for Granular
  Discovery (various resource inputs)

## Phase 5: Polish & Cross-Cutting

**Goal**: Enforce strict contracts and documentation.

- [x] T015 [P] Implement strict error handling: Ensure SDK/Plugin returns `Code.Unimplemented` if a
  granular capability check explicitly fails (FR-007)
- [x] T016 [P] Update `quickstart.md` with examples of implementing `Supports` for granular logic
- [x] T017 [P] Verify `buf breaking` passes to ensure backward compatibility
