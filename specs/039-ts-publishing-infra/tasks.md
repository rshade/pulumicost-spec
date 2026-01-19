# Tasks: TypeScript SDK Publishing Infrastructure

**Feature**: TypeScript SDK Publishing Infrastructure
**Status**: In Progress
**Branch**: `039-ts-publishing-infra`

## Phase 1: Setup

**Goal**: Prepare the project environment and configurations.

- [ ] T001 [P] Verify current state of `sdk/typescript/packages/client/package.json`
- [ ] T002 [P] Verify current state of `release-please-config.json`
- [ ] T003 [P] Verify current state of `.release-please-manifest.json`

## Phase 2: Foundational

**Goal**: Establish the core configurations required for publishing.

- [ ] T004 Rename package to `@rshade/finfocus-client` in `sdk/typescript/packages/client/package.json`
- [ ] T005 Add `publishConfig` with registry URL to `sdk/typescript/packages/client/package.json`
- [ ] T006 Add SDK package configuration to `release-please-config.json`
- [ ] T007 Initialize SDK version entry in `.release-please-manifest.json`

## Phase 3: Consumable NPM Package (US1)

**Goal**: Configure the package to be installable via NPM from GitHub Packages.

**Independent Test**: `cd sdk/typescript && npm install && npm run build -w packages/client` (verifies local build with new name)

- [ ] T008 [US1] Verify `sdk/typescript/packages/client/package.json` has correct name and publishConfig
- [ ] T009 [US1] Verify `npm run build` works with the new package name locally

## Phase 4: Automated Release Pipeline (US2)

**Goal**: Implement the CI/CD workflow for automated publishing.

**Independent Test**: Lint workflow file (visual inspection or tools).

- [ ] T010 [US2] Create `.github/workflows/publish-ts-client.yml` workflow file
- [ ] T011 [US2] Configure workflow trigger on `release` event (published)
- [ ] T012 [US2] Add job steps: checkout, setup-node, install dependencies, build, publish
- [ ] T013 [US2] Configure token permissions in workflow file

## Phase 5: Independent Versioning (US3)

**Goal**: Ensure independent versioning for the TS SDK.

**Independent Test**: Review `release-please-config.json` structure.

- [ ] T014 [US3] Verify `release-please-config.json` has separate entry for `sdk/typescript/packages/client`
- [ ] T015 [US3] Verify correct release type (`node`) is set for the SDK package

## Phase 6: Polish & Verification

**Goal**: Finalize and verify all configurations.

- [ ] T016 Lint JSON files (`package.json`, release configs)
- [ ] T017 Lint YAML workflow file
- [ ] T018 Run local build verification one last time

## Dependencies

- Phase 2 tasks are prerequisites for Phases 3, 4, and 5.
- Phase 4 depends on Phase 2 configuration being correct.

## Implementation Strategy

1. **Configure Package (Phase 2 & 3)**: Rename and configure `package.json` first.
2. **Configure Versioning (Phase 2 & 5)**: Update release-please configs.
3. **Build Pipeline (Phase 4)**: Create the GitHub Action.
4. **Verify (Phase 6)**: Ensure everything builds and lints correctly.
