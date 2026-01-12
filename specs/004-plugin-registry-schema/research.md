# Research: Plugin Registry Index JSON Schema

**Date**: 2025-11-23
**Feature**: 004-plugin-registry-schema

## Research Topics

### 1. JSON Schema draft 2020-12 Features for Conditional Validation

**Decision**: Use `dependentRequired` for deprecation_message requirement

**Rationale**: JSON Schema draft 2020-12 provides `dependentRequired` as the cleaner syntax
for conditional field requirements (when field A exists, field B is required). This is
simpler than `if/then/else` for this use case.

**Alternatives considered**:

- `if/then/else` construct - More verbose, harder to read for simple conditional
- Custom validation in consumer - Defeats purpose of schema validation

**Implementation**:

```json
"dependentRequired": {
  "deprecated": ["deprecation_message"]
}
```

### 2. Alignment with registry.proto Enums

**Decision**: Map proto enum values to JSON Schema string enums

**Rationale**: The registry.proto defines several enums that must be represented in the
JSON Schema. Proto enum values are SCREAMING_SNAKE_CASE but JSON convention is snake_case.

**Mappings**:

| Proto Enum    | Proto Values                   | JSON Schema Values                               |
| ------------- | ------------------------------ | ------------------------------------------------ |
| SecurityLevel | SECURITY_LEVEL_UNTRUSTED, etc. | `untrusted`, `community`, `verified`, `official` |
| Capabilities  | (string array in proto)        | `cost_retrieval`, `cost_projection`, etc.        |

**Note**: PluginInfo.capabilities is already `repeated string` in proto, so the schema
uses lowercase snake_case values matching the proto field values.

### 3. Plugin Name Pattern Validation

**Decision**: Use regex pattern `^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`

**Rationale**: Plugin names must be:

- Lowercase alphanumeric with hyphens
- Cannot start or end with hyphen
- Single character names allowed

**Alternatives considered**:

- Simpler `^[a-z0-9-]+$` - Would allow leading/trailing hyphens
- Require minimum 2 characters - Too restrictive for short names like `k8s`

### 4. Repository Format Validation

**Decision**: Use pattern `^[a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+$`

**Rationale**: GitHub repository format is `owner/repo` where both can contain alphanumeric,
hyphens, underscores, and repo can also contain dots.

**Examples**:

- `rshade/finfocus-plugin-kubecost` - Valid
- `my-org/my.repo.name` - Valid
- `invalid` - Invalid (no slash)

### 5. NPM Validation Script Integration

**Decision**: Add new `validate:registry` npm script

**Rationale**: Existing validation uses `validate_examples.js` with AJV. Registry validation
follows the same pattern but validates `examples/registry.json` against the registry schema.

**Implementation approach**:

1. Add `validate:registry` script to package.json
2. Update `validate` script to include registry validation
3. Follow existing pattern of using `--strict=false` for AJV

### 6. Semantic Version Pattern

**Decision**: Use pattern `^\d+\.\d+\.\d+$`

**Rationale**: Standard semver format without pre-release or build metadata. Keeps pattern
simple while covering all version fields (schema_version, min_spec_version, max_spec_version).

**Note**: Pre-release versions (e.g., `1.0.0-alpha`) not supported in initial version.
Can be added later if needed.

## Best Practices Applied

### JSON Schema Best Practices

1. **Use `$defs` for reusable components**: RegistryEntry defined once, referenced in
   additionalProperties
2. **Provide examples**: Each field with pattern validation includes examples array
3. **Set `additionalProperties: false`**: Prevents schema drift and catches typos
4. **Include descriptions**: All properties have clear descriptions

### Validation Integration Best Practices

1. **CI integration**: Validation runs in existing `make validate` pipeline
2. **Example-driven testing**: Example registry.json serves as both documentation and test
3. **Schema versioning**: `schema_version` field allows future schema evolution

## Unresolved Items

None - all technical decisions resolved through research.
