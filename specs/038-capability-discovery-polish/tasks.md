# Tasks: Capability Discovery Polish

**Feature**: Capability Discovery Polish
**Status**: In Progress
**Branch**: `038-capability-discovery-polish`

## Phase 1: Setup

**Goal**: Prepare project for development and verify environment.

- [x] T001 Verify branch is active and up to date with `main` in `specs/038-capability-discovery-polish`
- [x] T002 [P] Verify current state of `sdk/go/pluginsdk/dry_run.go` for DryRunHandler existence
- [x] T003 [P] Verify current state of `sdk/go/pluginsdk/sdk.go` to identify copy/make patterns
- [x] T004 [P] Verify current state of `proto/finfocus/v1/costsource.proto` for GetPluginInfoResponse definition

## Phase 2: Foundational

**Goal**: Fix non-functional prerequisites and "low-hanging fruit" issues (#295, #208, #209).

- [x] T005 Add clarification comments to `GetPluginInfoResponse` in `proto/finfocus/v1/costsource.proto`
- [x] T006 Add clarification comments to `SupportsResponse` in `proto/finfocus/v1/costsource.proto`
- [x] T007 Run `make generate` to ensure proto comments don't break generation in `proto/`
- [x] T008 [P] Rename compatibility tests in `sdk/go/testing/resource_id_test.go`
- [x] T009 [P] Add comments explaining round-trip nature of tests in `sdk/go/testing/resource_id_test.go`
- [x] T010 [P] Add import notes to `sdk/go/testing/README.md`
- [x] T011 [P] Fix incomplete code examples in `sdk/go/testing/README.md`

## Phase 3: Backward Compatibility (US3)

**Goal**: Implement the centralized capability conversion logic to support legacy clients.

**Independent Test**: `go test ./sdk/go/pluginsdk/capability_compat_test.go`

- [x] T012 [US3] Create new file `sdk/go/pluginsdk/capability_compat.go`
- [x] T013 [US3] Implement `legacyCapabilityNames` map in `sdk/go/pluginsdk/capability_compat.go`
- [x] T014 [US3] Implement `CapabilitiesToLegacyMetadata` function in `sdk/go/pluginsdk/capability_compat.go`
- [x] T015 [US3] Implement `CapabilitiesToLegacyMetadataWithWarnings` function in `sdk/go/pluginsdk/capability_compat.go`
- [x] T016 [US3] Create tests for helper functions in `sdk/go/pluginsdk/capability_compat_test.go`

## Phase 4: Automatic & Manual Discovery (US1 & US2)

**Goal**: Update SDK core logic to auto-discover capabilities from interfaces OR use manual overrides.

**Independent Test**: `go test ./sdk/go/pluginsdk/sdk_test.go` (ensure PluginInfo creation logic is correct)

- [x] T017 [US1] [US2] Refactor `NewPluginInfo` in `sdk/go/pluginsdk/plugin_info.go` to use `capability_compat` helpers
- [x] T018 [US1] [US2] Update `NewPluginInfo` to prioritize `WithCapabilities` override if provided in `sdk/go/pluginsdk/plugin_info.go`
- [x] T019 [US1] Ensure `NewPluginInfo` checks for `DryRunHandler` interface in `sdk/go/pluginsdk/plugin_info.go`
- [x] T020 [US1] Ensure `NewPluginInfo` checks for `RecommendationsProvider` interface in `sdk/go/pluginsdk/plugin_info.go`
- [x] T021 [US1] Ensure `NewPluginInfo` checks for `BudgetsProvider` interface in `sdk/go/pluginsdk/plugin_info.go`
- [x] T022 [US1] Ensure `NewPluginInfo` checks for `DismissProvider` interface in `sdk/go/pluginsdk/plugin_info.go`
- [x] T023 [US3] Refactor `Supports` method in `sdk/go/pluginsdk/sdk.go` to use
      `capability_compat` helpers (removing duplicate logic)
- [x] T024 [US1] [US2] Add unit tests for auto-discovery and override logic in
      `sdk/go/pluginsdk/plugin_info_test.go`

## Phase 5: SDK Performance Optimization (US4)

**Goal**: Optimize memory usage for slice copying.

**Independent Test**: `go test -bench=. ./sdk/go/pluginsdk/...`

- [x] T025 [US4] Create benchmark for slice copying in `sdk/go/pluginsdk/slice_benchmark_test.go`
- [x] T026 [US4] Replace `make+copy` with `append` pattern for `providers` in `sdk/go/pluginsdk/sdk.go`
- [x] T027 [US4] Replace `make+copy` with `append` pattern for `capabilities` in `sdk/go/pluginsdk/sdk.go`
- [x] T028 [US4] Run benchmarks again to verify improvement

## Phase 6: Polish & Documentation

**Goal**: Finalize documentation and verify system integrity.

- [x] T029 Add "Capability Discovery Pattern" section to `CLAUDE.md`
- [x] T030 Add "Capability Discovery" section with examples to `sdk/go/pluginsdk/README.md`
- [x] T031 Run `make lint` to ensure all new code passes linting
- [x] T032 Run `make validate` for full project validation
- [x] T033 [P] Verify TypeScript SDK client prefers `capabilities` enum over `metadata`
      in `sdk/typescript/`

## Dependencies

1. Phase 3 (Back Compat) MUST be completed before Phase 4 (Discovery Logic) because the
   refactoring in Phase 4 relies on the helpers created in Phase 3.
2. Phase 1 & 2 can be done in parallel with Phase 3 & 4 (mostly).
3. Phase 5 is independent.

## Implementation Strategy

1. **Start with Foundational fixes (Phase 2)** to get the "easy wins" and doc fixes out of the
   way.
2. **Implement the Compatibility Helper (Phase 3)** to create the single source of truth.
3. **Refactor the Core Logic (Phase 4)** to use the helper and implement the auto-discovery rules.
4. **Optimize Performance (Phase 5)** as a final code polish.
5. **Verify & Document (Phase 6)**.
