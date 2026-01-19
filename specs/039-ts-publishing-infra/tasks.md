# Tasks: TypeScript SDK Publishing Infrastructure

**Feature**: TypeScript SDK Publishing Infrastructure
**Status**: In Progress
**Branch**: `039-ts-publishing-infra`

## Phase 1: Configuration

**Goal**: Update package metadata and release tool configuration.

- [ ] T001 Rename package to `@rshade/finfocus-client` and add `publishConfig` in `sdk/typescript/packages/client/package.json`
- [ ] T002 Add SDK package to `release-please-config.json`
- [ ] T003 Initialize SDK version in `.release-please-manifest.json`

## Phase 2: CI/CD

**Goal**: Create the publishing workflow.

- [ ] T004 Create `.github/workflows/publish-ts-client.yml`

## Phase 3: Verification

**Goal**: Ensure build stability.

- [ ] T005 Verify build passes with new package name (`cd sdk/typescript && npm install && npm run build`)
- [ ] T006 Lint `.github/workflows/publish-ts-client.yml` (using `yamllint` or visual inspection)
