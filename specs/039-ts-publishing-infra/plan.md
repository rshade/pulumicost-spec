# Implementation Plan - TypeScript SDK Publishing Infrastructure

**Feature**: TypeScript SDK Publishing Infrastructure
**Branch**: `039-ts-publishing-infra`
**Status**: Planning

## Technical Context

We are enabling the automated publication of the `finfocus-client` TypeScript SDK to GitHub Packages. This transforms the SDK from a source-only artifact into a consumable NPM package (`@rshade/finfocus-client`).

### Architecture

- **Registry**: GitHub Packages (`npm.pkg.github.com`).
- **Versioning**: Automated via `release-please` (Google's Release Please Action).
- **CI/CD**: GitHub Actions workflow triggering on release publication.
- **Workspace**: The TS SDK is part of a nested NPM workspace rooted at `sdk/typescript`.

### Existing Components

- `sdk/typescript/package.json`: Workspace root.
- `sdk/typescript/packages/client/package.json`: Target package (currently named `finfocus-client`).
- `release-please-config.json`: Configuration for versioning strategies.
- `.github/workflows/`: CI pipeline location.

## Constitution Check

| Principle | Compliance | Notes |
| :--- | :--- | :--- |
| **III. Spec Consumes, It Does Not Calculate** | N/A | Infrastructure task. |
| **IV. Separation of Concerns** | ✅ | Creating a dedicated publishing pipeline separate from the Go SDK/Spec. |
| **VII. Documentation & Identity** | ✅ | Scoping the package to `@rshade` establishes clear identity. |

## Phased Implementation

### Phase 1: Configuration (Package & Versioning)

**Goal**: Prepare the package metadata and release automation configuration.

1.  **Package Scoping**: Rename `finfocus-client` to `@rshade/finfocus-client` in `sdk/typescript/packages/client/package.json`.
2.  **Publish Config**: Add `publishConfig` pointing to GitHub Packages.
3.  **Release Please**: Configure `release-please-config.json` to track the SDK path and `.release-please-manifest.json` to initialize the version.

### Phase 2: CI/CD Pipeline

**Goal**: Create the workflow to build and publish the package.

1.  **Workflow File**: Create `.github/workflows/publish-ts-client.yml`.
2.  **Build Logic**: Ensure the workflow installs dependencies at the workspace root (`sdk/typescript`) and builds the specific package before publishing.

### Phase 3: Verification

**Goal**: Verify configurations (dry run where possible).

1.  **Lint**: Ensure JSON/YAML files are valid.
2.  **Build**: Verify `npm run build` works locally with the new name.

## Unknowns & Risks

- **Token Permissions**: GitHub Packages usually requires `packages: write` permission in the workflow. We must ensure the token used has this scope.
- **Workspace Isolation**: Publishing a package from within a workspace can sometimes be tricky with `npm publish` if it relies on hoisted dependencies. `npm publish -w` or navigating to the directory is usually required.

## Verification Plan

### Automated Tests
- `npm run build` in `sdk/typescript/packages/client`.
- `yamllint` on the new workflow file.
- `validate-json` on modified config files.

### Manual Verification
- Review diffs for `package.json` and release configs.
