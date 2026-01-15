# Research: Plugin Capability Discovery

## Unknowns & Research Tasks

| Unknown                | Research Task                                                     | Finding                                                                                |
| ---------------------- | ----------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| Enum Placement         | Verify if `enums.proto` is the right place for `PluginCapability` | Yes, `enums.proto` already contains domain-specific enums.                             |
| Auto-Discovery Logic   | How to detect implemented interfaces in Go SDK?                   | Use Go type assertions (e.g., `plugin.(SupportsProvider)`) in `GetPluginInfo` handler. |
| Backward Compatibility | How to map new enums back to legacy string map?                   | Maintain a static map of `PluginCapability` to `string` keys in the SDK.               |

## Decisions

### Decision 1: PluginCapability Enum Placement

- **Decision**: Add `PluginCapability` to `proto/finfocus/v1/enums.proto`.
- **Rationale**: Keeps `costsource.proto` clean and follows the pattern of segregating enums into a dedicated file.
- **Alternatives considered**: Inline in `costsource.proto` (rejected to maintain consistency).

### Decision 2: Go SDK Auto-Discovery

- **Decision**: The SDK's `GetPluginInfo` implementation will use type assertions on the `Plugin` instance.
- **Rationale**: Go interfaces are checked at runtime. By asserting against `SupportsProvider`,
  `RecommendationsProvider`, etc., the SDK can determine capabilities without user intervention.
- **Alternatives considered**: Manual registration (rejected for poor DX).

### Decision 3: Legacy Map Support

- **Decision**: Continue populating `SupportsResponse.capabilities` map[string]bool and add it to
  `GetPluginInfoResponse.metadata` if needed, but primarily keep it where it exists.
- **Rationale**: Existing clients expect it. The new enum field is additive.
