feat(ci): implement automated typescript sdk publishing

## Summary

This PR enables the automated publication of the TypeScript client SDK
to GitHub Packages. It scopes the package to `@rshade/finfocus-client`,
configures independent versioning via `release-please`, and adds a
GitHub Actions workflow to handle the build-and-publish pipeline
triggered by new releases.

## Test plan

- [x] Local build verification: `npm run build -w packages/client` passed
- [x] JSON validation: `npm run validate` passed
- [x] YAML linting: `.github/workflows/publish-ts-client.yml` verified
- [x] Verified independent versioning config in `release-please-config.json`

## Changes

### New files

- `.github/workflows/publish-ts-client.yml` - Automated publishing pipeline

### Modified files

- `sdk/typescript/packages/client/package.json` - Scoped name and registry config
- `release-please-config.json` - Added TS SDK versioning strategy
- `.release-please-manifest.json` - Initialized SDK version tracking
- `specs/039-ts-publishing-infra/` - Added plan, research, and tasks

### Housekeeping

- Updated `GEMINI.md` context to reflect TS SDK publishing capabilities

Closes #311
