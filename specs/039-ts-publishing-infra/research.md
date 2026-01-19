# Research: TypeScript SDK Publishing Infrastructure

**Feature**: TypeScript SDK Publishing Infrastructure
**Branch**: `039-ts-publishing-infra`
**Status**: Complete

## Strategic Choice: GitHub Packages vs. npmjs.com

**Decision**: GitHub Packages (`npm.pkg.github.com`)

**Rationale**:
- **Ecosystem Integration**: Keeps code, issues, and packages tightly coupled in the same ecosystem.
- **Permission Management**: Allows for strict permission scoping (same tokens for code and packages).
- **Enterprise Alignment**: Common choice for internal/spec-related packages in enterprise environments.

**Alternatives Considered**:
- **npmjs.com**: The public registry. Rejected because it requires separate authentication tokens and account management, increasing friction for a specification project.

## Implementation Details

### Package Naming
- **Scope**: `@rshade` (Owner of the repo)
- **Package Name**: `@rshade/finfocus-client` (Scoped packages are required for GitHub Packages)

### Publishing Workflow
- **Trigger**: `release` event with type `published`.
- **Authentication**: `GITHUB_TOKEN` (Must have `packages: write` permission).
- **Registry**: Configured in `.npmrc` during the workflow run or via `publishConfig` in `package.json`.

### Versioning Strategy
- **Tool**: `release-please`
- **Config**: Add a separate entry in `release-please-config.json` for `sdk/typescript/packages/client`.
- **Type**: `node` (Standard semver for Node.js packages).
