# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.1](https://github.com/rshade/pulumicost-spec/compare/v0.4.0...v0.4.1) (2025-11-29)

### Added

- add trace ID validation to TracingUnaryServerInterceptor ([#96](https://github.com/rshade/pulumicost-spec/issues/96)) ([dd410cd](https://github.com/rshade/pulumicost-spec/commit/dd410cdbc2ca88ecc87dc3bef3e4fa3488efd714)), closes [#94](https://github.com/rshade/pulumicost-spec/issues/94)
- **focus:** add complete FOCUS 1.2 column coverage with builder API ([#100](https://github.com/rshade/pulumicost-spec/issues/100)) ([2355acf](https://github.com/rshade/pulumicost-spec/commit/2355acf206f5d2ed70d9cd47fd9033525f8fe552))
- **focus:** Implement FOCUS 1.2 integration ([#99](https://github.com/rshade/pulumicost-spec/issues/99)) ([913b6ef](https://github.com/rshade/pulumicost-spec/commit/913b6ef9d9a9ca277058ded46dcf5f7cfadc7aab))
- **sdk:** migrate pluginsdk from core to spec ([#97](https://github.com/rshade/pulumicost-spec/issues/97)) ([2e35cbf](https://github.com/rshade/pulumicost-spec/commit/2e35cbf548f91f0151901795227dc07d578f3220))
- **testing:** add structured logging example for EstimateCost RPC ([#93](https://github.com/rshade/pulumicost-spec/issues/93)) ([4c583c0](https://github.com/rshade/pulumicost-spec/commit/4c583c0fc0cc179b57fb919eef6ba349d4cf7187)), closes [#83](https://github.com/rshade/pulumicost-spec/issues/83)

## [Unreleased]

### Added

- **ci:** add performance regression tests with benchmark comparison
- **docs:** add EstimateCost cross-provider coverage matrix to examples/README.md
- **sdk:** add trace ID validation to TracingUnaryServerInterceptor for security ([#94](https://github.com/rshade/pulumicost-spec/issues/94))

### Security

- **sdk:** prevent log injection attacks through malformed trace IDs by validating and replacing invalid values

## [0.4.0](https://github.com/rshade/pulumicost-spec/compare/v0.3.0...v0.4.0) (2025-11-26)

### Added

- **rpc:** implement EstimateCost RPC for what-if cost analysis ([#90](https://github.com/rshade/pulumicost-spec/issues/90)) ([d6f3c95](https://github.com/rshade/pulumicost-spec/commit/d6f3c9566da8d28550923edfe4ffe34d3c143e0e)), closes [#79](https://github.com/rshade/pulumicost-spec/issues/79)

## [0.3.0](https://github.com/rshade/pulumicost-spec/compare/v0.2.0...v0.3.0) (2025-11-24)

### Added

- **sdk:** add zerolog logging utilities for plugin standardization ([#76](https://github.com/rshade/pulumicost-spec/issues/76)) ([6d5b5ac](https://github.com/rshade/pulumicost-spec/commit/6d5b5ac06329dce03a99b595e41d5ca1273b7c40)), closes [#75](https://github.com/rshade/pulumicost-spec/issues/75)

### Documentation

- udpate sdk/go/registry/CLAUDE.md for enums ([#78](https://github.com/rshade/pulumicost-spec/issues/78)) ([f15ef76](https://github.com/rshade/pulumicost-spec/commit/f15ef769ac491056504f7a0376b413727f7969e0)), closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)

## [0.2.0](https://github.com/rshade/pulumicost-spec/compare/v0.1.0...v0.2.0) (2025-11-24)

### ⚠ BREAKING CHANGES

- **proto:** None - 100% backward compatible (additive proto changes only)
- **registry:** None - 100% backward compatible

### Added

- **proto:** enhance GetPricingSpec with transparent pricing breakdown ([#67](https://github.com/rshade/pulumicost-spec/issues/67)) ([336144e](https://github.com/rshade/pulumicost-spec/commit/336144e45e2334a677af4a0f5ccb3994126cf22a)), closes [#62](https://github.com/rshade/pulumicost-spec/issues/62)
- **schema:** add plugin registry index JSON Schema ([#70](https://github.com/rshade/pulumicost-spec/issues/70)) ([79938a7](https://github.com/rshade/pulumicost-spec/commit/79938a7fc473cd997465f3a71c734b0b2b3b692b)), closes [#68](https://github.com/rshade/pulumicost-spec/issues/68)

### Fixed

- **release:** remove release-as constraint to allow version bumps ([#73](https://github.com/rshade/pulumicost-spec/issues/73)) ([706ec65](https://github.com/rshade/pulumicost-spec/commit/706ec65dbe779242674aae94c6bc89de0f8c252a))

### Performance

- **registry:** optimize enum validation for zero-allocation performance ([#63](https://github.com/rshade/pulumicost-spec/issues/63)) ([6d3c124](https://github.com/rshade/pulumicost-spec/commit/6d3c124b4230485ee27288051181a965a31daf50)), closes [#33](https://github.com/rshade/pulumicost-spec/issues/33)

### Documentation

- **spec:** document Supports() RPC verification for issue [#64](https://github.com/rshade/pulumicost-spec/issues/64) ([#66](https://github.com/rshade/pulumicost-spec/issues/66)) ([3a17c4f](https://github.com/rshade/pulumicost-spec/commit/3a17c4f516488d033549c26c4bb9ad8ed957e3b0))

## [0.1.0](https://github.com/rshade/pulumicost-spec/compare/v0.1.0...v0.1.0) (2025-11-24)

### ⚠ BREAKING CHANGES

- **proto:** None - 100% backward compatible (additive proto changes only)
- **registry:** None - 100% backward compatible

### Added

- **proto:** enhance GetPricingSpec with transparent pricing breakdown ([#67](https://github.com/rshade/pulumicost-spec/issues/67)) ([336144e](https://github.com/rshade/pulumicost-spec/commit/336144e45e2334a677af4a0f5ccb3994126cf22a)), closes [#62](https://github.com/rshade/pulumicost-spec/issues/62)
- **schema:** add plugin registry index JSON Schema ([#70](https://github.com/rshade/pulumicost-spec/issues/70)) ([79938a7](https://github.com/rshade/pulumicost-spec/commit/79938a7fc473cd997465f3a71c734b0b2b3b692b)), closes [#68](https://github.com/rshade/pulumicost-spec/issues/68)

### Performance

- **registry:** optimize enum validation for zero-allocation performance ([#63](https://github.com/rshade/pulumicost-spec/issues/63)) ([6d3c124](https://github.com/rshade/pulumicost-spec/commit/6d3c124b4230485ee27288051181a965a31daf50)), closes [#33](https://github.com/rshade/pulumicost-spec/issues/33)

### Documentation

- **spec:** document Supports() RPC verification for issue [#64](https://github.com/rshade/pulumicost-spec/issues/64) ([#66](https://github.com/rshade/pulumicost-spec/issues/66)) ([3a17c4f](https://github.com/rshade/pulumicost-spec/commit/3a17c4f516488d033549c26c4bb9ad8ed957e3b0))

## [Unreleased]

### Added

- **Schema**: Add plugin registry index JSON Schema (`schemas/plugin_registry.schema.json`)
  - Validates registry.json files for `pulumicost plugin install` discovery
  - Aligns with registry.proto definitions (SecurityLevel, capabilities, providers)
  - Includes `dependentRequired` for deprecation_message when deprecated is true
  - npm validation scripts: `validate:registry`, `validate:registry-schema`
  - Example registry with kubecost and aws-public plugins
  - Closes [#68](https://github.com/rshade/pulumicost-spec/issues/68)

### Changed

- **Performance**: Optimized registry package enum validation for zero-allocation performance
  - Converted all 8 enum types (Provider, DiscoverySource, PluginStatus, SecurityLevel, InstallationMethod,
    PluginCapability, SystemPermission, AuthMethod) from function-returned slices to package-level variables
  - Achieved 0 B/op, 0 allocs/op across all validation functions (previously 1 alloc/op)
  - Performance improved to 5-12 ns/op (2x faster than map-based alternatives)
  - Memory footprint reduced to ~608 bytes total for all enums (vs ~3.5 KB for maps)
  - Established validation pattern for future SDK enums (see `specs/001-domain-enum-optimization/`)

## [0.1.0](https://github.com/rshade/pulumicost-spec/compare/v0.1.0...v0.1.0) (2025-11-18)

### Added

- Add comprehensive Plugin Registry Specification ([#29](https://github.com/rshade/pulumicost-spec/issues/29)) ([5825eaa](https://github.com/rshade/pulumicost-spec/commit/5825eaab1f343b1f9fcb966c38effb6a743d3494)), closes [#8](https://github.com/rshade/pulumicost-spec/issues/8)
- comprehensive testing framework and enterprise CI/CD pipeline ([#19](https://github.com/rshade/pulumicost-spec/issues/19)) ([3a235ef](https://github.com/rshade/pulumicost-spec/commit/3a235eff0ce4a172a7b920326066227540c6c8b8))
- enhance Provider enum with String() method and improved error m… ([#39](https://github.com/rshade/pulumicost-spec/issues/39)) ([acbaf0c](https://github.com/rshade/pulumicost-spec/commit/acbaf0c6db29d986211d3ca8ef19a66e3162e2c2)), closes [#4](https://github.com/rshade/pulumicost-spec/issues/4)
- freeze costsource.proto v0.1.0 specification ([#17](https://github.com/rshade/pulumicost-spec/issues/17)) ([3b485b9](https://github.com/rshade/pulumicost-spec/commit/3b485b96fef7dc2992c166980f324907e4ff06bd)), closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)
- freeze costsource.proto v0.1.0 specification ([#18](https://github.com/rshade/pulumicost-spec/issues/18)) ([a085bd2](https://github.com/rshade/pulumicost-spec/commit/a085bd202e266189efba92a1832fa8f90c0931e6)), closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)

### Documentation

- add comprehensive plugin developer guide ([#16](https://github.com/rshade/pulumicost-spec/issues/16)) ([b0a5eb3](https://github.com/rshade/pulumicost-spec/commit/b0a5eb3396b5c49839d742234666667e5e7b1ee7)), closes [#2](https://github.com/rshade/pulumicost-spec/issues/2)
- establish constitution v1.0.0 for gRPC proto specification governance ([#57](https://github.com/rshade/pulumicost-spec/issues/57)) ([54578aa](https://github.com/rshade/pulumicost-spec/commit/54578aa259cb8907f989196ce2b73e37e57f906f))

## [0.1.0](https://github.com/rshade/pulumicost-spec/compare/v0.1.0...v0.1.0) (2025-11-18)

### Added

- Add comprehensive Plugin Registry Specification ([#29](https://github.com/rshade/pulumicost-spec/issues/29)) ([5825eaa](https://github.com/rshade/pulumicost-spec/commit/5825eaab1f343b1f9fcb966c38effb6a743d3494)), closes [#8](https://github.com/rshade/pulumicost-spec/issues/8)
- comprehensive testing framework and enterprise CI/CD pipeline ([#19](https://github.com/rshade/pulumicost-spec/issues/19)) ([3a235ef](https://github.com/rshade/pulumicost-spec/commit/3a235eff0ce4a172a7b920326066227540c6c8b8))
- enhance Provider enum with String() method and improved error m… ([#39](https://github.com/rshade/pulumicost-spec/issues/39)) ([acbaf0c](https://github.com/rshade/pulumicost-spec/commit/acbaf0c6db29d986211d3ca8ef19a66e3162e2c2)), closes [#4](https://github.com/rshade/pulumicost-spec/issues/4)
- freeze costsource.proto v0.1.0 specification ([#17](https://github.com/rshade/pulumicost-spec/issues/17)) ([3b485b9](https://github.com/rshade/pulumicost-spec/commit/3b485b96fef7dc2992c166980f324907e4ff06bd)), closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)
- freeze costsource.proto v0.1.0 specification ([#18](https://github.com/rshade/pulumicost-spec/issues/18)) ([a085bd2](https://github.com/rshade/pulumicost-spec/commit/a085bd202e266189efba92a1832fa8f90c0931e6)), closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)

### Documentation

- add comprehensive plugin developer guide ([#16](https://github.com/rshade/pulumicost-spec/issues/16)) ([b0a5eb3](https://github.com/rshade/pulumicost-spec/commit/b0a5eb3396b5c49839d742234666667e5e7b1ee7)), closes [#2](https://github.com/rshade/pulumicost-spec/issues/2)
- establish constitution v1.0.0 for gRPC proto specification governance ([#57](https://github.com/rshade/pulumicost-spec/issues/57)) ([54578aa](https://github.com/rshade/pulumicost-spec/commit/54578aa259cb8907f989196ce2b73e37e57f906f))

## [Unreleased]

### Added

### Changed

### Fixed

## [0.1.0] - 2025-11-18

### Added

- Add comprehensive Plugin Registry Specification
  ([#29](https://github.com/rshade/pulumicost-spec/issues/29))
  ([5825eaa](https://github.com/rshade/pulumicost-spec/commit/5825eaab1f343b1f9fcb966c38effb6a743d3494)),
  closes [#8](https://github.com/rshade/pulumicost-spec/issues/8)
- Comprehensive testing framework and enterprise CI/CD pipeline
  ([#19](https://github.com/rshade/pulumicost-spec/issues/19))
  ([3a235ef](https://github.com/rshade/pulumicost-spec/commit/3a235eff0ce4a172a7b920326066227540c6c8b8))
- Enhance Provider enum with String() method and improved error handling
  ([#39](https://github.com/rshade/pulumicost-spec/issues/39))
  ([acbaf0c](https://github.com/rshade/pulumicost-spec/commit/acbaf0c6db29d986211d3ca8ef19a66e3162e2c2)),
  closes [#4](https://github.com/rshade/pulumicost-spec/issues/4)
- Freeze costsource.proto v0.1.0 specification
  ([#17](https://github.com/rshade/pulumicost-spec/issues/17))
  ([3b485b9](https://github.com/rshade/pulumicost-spec/commit/3b485b96fef7dc2992c166980f324907e4ff06bd)),
  closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)
- Freeze costsource.proto v0.1.0 specification
  ([#18](https://github.com/rshade/pulumicost-spec/issues/18))
  ([a085bd2](https://github.com/rshade/pulumicost-spec/commit/a085bd202e266189efba92a1832fa8f90c0931e6)),
  closes [#3](https://github.com/rshade/pulumicost-spec/issues/3)

### Documentation

- Add comprehensive plugin developer guide
  ([#16](https://github.com/rshade/pulumicost-spec/issues/16))
  ([b0a5eb3](https://github.com/rshade/pulumicost-spec/commit/b0a5eb3396b5c49839d742234666667e5e7b1ee7)),
  closes [#2](https://github.com/rshade/pulumicost-spec/issues/2)
- Establish constitution v1.0.0 for gRPC proto specification governance
  ([#57](https://github.com/rshade/pulumicost-spec/issues/57))
  ([54578aa](https://github.com/rshade/pulumicost-spec/commit/54578aa259cb8907f989196ce2b73e37e57f906f))

[0.1.0]: https://github.com/rshade/pulumicost-spec/releases/tag/v0.1.0
