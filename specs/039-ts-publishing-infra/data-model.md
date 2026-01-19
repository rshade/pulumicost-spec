# Data Model: TypeScript SDK Publishing Infrastructure

**Feature**: TypeScript SDK Publishing Infrastructure
**Branch**: `039-ts-publishing-infra`

## Entities

### Package Configuration (`package.json`)

**Description**: The manifest file for the `@rshade/finfocus-client` package.

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | MUST be `@rshade/finfocus-client` (scoped). |
| `version` | `string` | Semantic version (e.g., `0.1.0`). Managed by `release-please`. |
| `publishConfig` | `object` | Configuration for publishing to the registry. |
| `publishConfig.registry` | `string` | URL: `https://npm.pkg.github.com`. |

### Release Configuration (`release-please-config.json`)

**Description**: Configuration for the release automation tool.

| Field | Type | Description |
| :--- | :--- | :--- |
| `packages` | `map` | Map of path to package configuration. |
| `packages["sdk/typescript/packages/client"]` | `object` | Configuration for the client SDK. |
| `packages[...].release-type` | `string` | MUST be `node`. |
| `packages[...].package-name` | `string` | MUST be `@rshade/finfocus-client`. |

### GitHub Workflow (`publish-ts-client.yml`)

**Description**: The CI/CD pipeline definition.

| Trigger | Event |
| :--- | :--- |
| `on.release` | `types: [published]` |

| Job | Steps |
| :--- | :--- |
| `publish` | Checkout -> Setup Node -> Install -> Build -> Publish |
