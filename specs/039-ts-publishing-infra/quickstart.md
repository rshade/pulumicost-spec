# Quickstart: Publishing the TypeScript SDK

This guide explains how the publishing infrastructure works for maintainers.

## How it Works

1. **Develop**: Make changes to `sdk/typescript/packages/client`.
2. **Merge**: Merge your PR to `main`.
3. **Automate**: `release-please` will automatically create a "Release PR" proposing a version bump
   (e.g., `v0.1.0` -> `v0.1.1`) and updating `CHANGELOG.md` in the client directory.
4. **Release**: Merge the "Release PR". `release-please` will create a GitHub Release tag.
5. **Publish**: The `publish-ts-client` workflow triggers on the new release, builds the SDK, and
   publishes it to GitHub Packages.

## Installing the SDK

To install the published package, consumers must configure their npm client to look at the GitHub
Registry for the `@rshade` scope.

### 1. Create/Update `.npmrc`

Create a `.npmrc` file in your project root (or user home):

```ini
@rshade:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=${GITHUB_TOKEN}
```

### 2. Install

```bash
npm install @rshade/finfocus-client
```

## Manual Verification (Dry Run)

To verify the build locally before pushing:

```bash
cd sdk/typescript
npm install
npm run build -w packages/client
```
