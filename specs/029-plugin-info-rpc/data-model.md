<!-- markdownlint-disable MD013 -->
# Data Model: Add GetPluginInfo RPC

## Entities

### Plugin Info

Represents the metadata of a loaded cost source plugin.

| Field          | Type                | Description                                                                                                                                    |
| :------------- | :------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`         | `string`            | The display name of the plugin (e.g., "aws-cost-plugin").                                                                                      |
| `version`      | `string`            | The semantic version of the plugin implementation (e.g., "v1.2.0").                                                                            |
| `spec_version` | `string`            | The version of the `finfocus-spec` protocol the plugin was compiled against (e.g., "v0.4.11").                                               |
| `providers`    | `[]string`          | List of cloud providers supported by this plugin (e.g., `["aws"]`).                                                                            |
| `metadata`     | `map[string]string` | Optional key-value pairs for additional metadata (e.g., build hash, commit ID). Free-form; no key restrictions or size limits enforced by SDK. |

## Validation Rules

1. **Spec Version**: MUST be a valid Semantic Version (vX.Y.Z).
2. **Name**: MUST be non-empty.
3. **Version**: MUST be non-empty.
4. **Providers**: MUST contain at least one provider if functional.

## RPC Methods

### `GetPluginInfo`

Retrieves the `PluginInfo` from the cost source.

- **Request**: `GetPluginInfoRequest` (Empty)
- **Response**: `GetPluginInfoResponse`
- **Error Handling**:
  - `Unimplemented`: Plugin is too old to support this RPC. Client should treat `spec_version` as "Unknown" and proceed with warning.
  - `Internal`: Plugin failed to retrieve its own metadata (unexpected).
