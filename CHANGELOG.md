# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.7](https://github.com/rshade/pulumicost-spec/compare/v0.4.6...v0.4.7) (2025-12-15)


### Added

* **proto:** add arn field to actual cost request ([#160](https://github.com/rshade/pulumicost-spec/issues/160)) ([a75b42b](https://github.com/rshade/pulumicost-spec/commit/a75b42b355395f590b7d211475de05f5083bd82b)), closes [#157](https://github.com/rshade/pulumicost-spec/issues/157)

## [0.4.6](https://github.com/rshade/pulumicost-spec/compare/v0.4.5...v0.4.6) (2025-12-11)


### Added

* **pluginsdk:** add request validation helpers ([#151](https://github.com/rshade/pulumicost-spec/issues/151)) ([ef71ba6](https://github.com/rshade/pulumicost-spec/commit/ef71ba6018169b630f3068deab45b97bd1e3522c)), closes [#130](https://github.com/rshade/pulumicost-spec/issues/130)


### Fixed

* updating small bugs in spec ([#156](https://github.com/rshade/pulumicost-spec/issues/156)) ([83df05f](https://github.com/rshade/pulumicost-spec/commit/83df05f11ec2690c9b4e594128a9ba4022c20c8b))

## [0.4.5](https://github.com/rshade/pulumicost-spec/compare/v0.4.4...v0.4.5) (2025-12-10)


### Added

* **pluginsdk:** add mapping package for property extraction ([#148](https://github.com/rshade/pulumicost-spec/issues/148)) ([8fd1524](https://github.com/rshade/pulumicost-spec/commit/8fd1524877218272a7d219239058bfee43c294bc)), closes [#128](https://github.com/rshade/pulumicost-spec/issues/128)
* **proto:** add getbudgets rpc for unified budget visibility ([#149](https://github.com/rshade/pulumicost-spec/issues/149)) ([b4018d7](https://github.com/rshade/pulumicost-spec/commit/b4018d794fd5b2c2f54541102c7625ffd900f26f)), closes [#123](https://github.com/rshade/pulumicost-spec/issues/123)
* **sdk:** Support PULUMICOST_LOG_FILE for unified logging ([#145](https://github.com/rshade/pulumicost-spec/issues/145)) ([6b9a9b3](https://github.com/rshade/pulumicost-spec/commit/6b9a9b38e24f9b49a05095c140bd2500c4e8090b)), closes [#131](https://github.com/rshade/pulumicost-spec/issues/131)


### Documentation

* Document pluginsdk.Serve() behavior and configuration ([#146](https://github.com/rshade/pulumicost-spec/issues/146)) ([30687f9](https://github.com/rshade/pulumicost-spec/commit/30687f9ca34a1c262a0ac6b8f66f98301dce1987))
* **sdk:** add core-plugin interface docs and contract tests ([#150](https://github.com/rshade/pulumicost-spec/issues/150)) ([87d4428](https://github.com/rshade/pulumicost-spec/commit/87d44289a5a63b25d8baeccb11f7eb9f56ba3128)), closes [#132](https://github.com/rshade/pulumicost-spec/issues/132) [#133](https://github.com/rshade/pulumicost-spec/issues/133) [#134](https://github.com/rshade/pulumicost-spec/issues/134) [#135](https://github.com/rshade/pulumicost-spec/issues/135)

## [0.4.4](https://github.com/rshade/pulumicost-spec/compare/v0.4.3...v0.4.4) (2025-12-09)

### Added

- **pluginsdk:** add --port flag parsing for multi-plugin orchestration ([#143](https://github.com/rshade/pulumicost-spec/issues/143)) ([c0b0528](https://github.com/rshade/pulumicost-spec/commit/c0b05288e69dc70ad0105165572a5fa3714ed27f)), closes [#129](https://github.com/rshade/pulumicost-spec/issues/129) [#137](https://github.com/rshade/pulumicost-spec/issues/137)
- **pluginsdk:** add fallback hint enum for plugin orchestration ([#126](https://github.com/rshade/pulumicost-spec/issues/126)) ([ef7aab0](https://github.com/rshade/pulumicost-spec/commit/ef7aab0576e3b4815c4d273c33800557734ebb37)), closes [#124](https://github.com/rshade/pulumicost-spec/issues/124)
- **pluginsdk:** centralize environment variable handling ([#139](https://github.com/rshade/pulumicost-spec/issues/139)) ([4c9e279](https://github.com/rshade/pulumicost-spec/commit/4c9e279ad38c58ee28178f61041c387c512654ca)), closes [#127](https://github.com/rshade/pulumicost-spec/issues/127)
- **proto:** add getbudgets rpc for unified budget visibility across providers ([#145](https://github.com/rshade/pulumicost-spec/issues/145)) ([abc123d](https://github.com/rshade/pulumicost-spec/commit/abc123def456ghi789jkl012))
- **proto:** add getrecommendations rpc for finops optimization ([#125](https://github.com/rshade/pulumicost-spec/issues/125)) ([ecf92f0](https://github.com/rshade/pulumicost-spec/commit/ecf92f0af6c1dbd1d036e92a1e999a7576debaef))

### Fixed

- adding in edge case tests, and benchmark ([#142](https://github.com/rshade/pulumicost-spec/issues/142)) ([881132b](https://github.com/rshade/pulumicost-spec/commit/881132bd87d1bb69ecd9f3abff01527da16fc08f))

### Documentation

- adding in claude speckit ([#144](https://github.com/rshade/pulumicost-spec/issues/144)) ([70a6e78](https://github.com/rshade/pulumicost-spec/commit/70a6e78fffba6ac81c518a4225a4da56a34aafdf))

## [0.4.3](https://github.com/rshade/pulumicost-spec/compare/v0.4.2...v0.4.3) (2025-12-03)

### Added

- **ci:** Add Lefthook git hooks with commitlint validation ([#120](https://github.com/rshade/pulumicost-spec/issues/120)) ([afdf8f7](https://github.com/rshade/pulumicost-spec/commit/afdf8f78afb2cac5dfdad95359acb2871c727be7)), closes [#55](https://github.com/rshade/pulumicost-spec/issues/55)
- **pluginsdk:** Add conformance testing support for Plugin interface ([#118](https://github.com/rshade/pulumicost-spec/issues/118)) ([8df49c0](https://github.com/rshade/pulumicost-spec/commit/8df49c041adc671843aa0fa6bda987c50a5bcc7a)), closes [#98](https://github.com/rshade/pulumicost-spec/issues/98)
- **pluginsdk:** add Prometheus metrics instrumentation for plugins ([#119](https://github.com/rshade/pulumicost-spec/issues/119)) ([9365aef](https://github.com/rshade/pulumicost-spec/commit/9365aef0636144f1bcf0695db68681823a889fe0)), closes [#80](https://github.com/rshade/pulumicost-spec/issues/80)
- **sdk/go/currency:** extract ISO 4217 validation as reusable package (T101) ([#116](https://github.com/rshade/pulumicost-spec/issues/116)) ([97e34f5](https://github.com/rshade/pulumicost-spec/commit/97e34f5fba0f1a635ab9651ed2bc80510898b962)), closes [#101](https://github.com/rshade/pulumicost-spec/issues/101)

## [0.4.2](https://github.com/rshade/pulumicost-spec/compare/v0.4.1...v0.4.2) (2025-11-30)

### Added

- **ci:** add performance regression testing workflow ([8944316](https://github.com/rshade/pulumicost-spec/commit/8944316a5337a12652efaa700999b7fd400517de))
- run concurrent benchmark for EstimateCost ([#113](https://github.com/rshade/pulumicost-spec/issues/113)) ([0ffcdc4](https://github.com/rshade/pulumicost-spec/commit/0ffcdc48e132b1eece31a1c51280cd250c608c23))
- **testing:** add distributed tracing example for EstimateCost (T042) ([#112](https://github.com/rshade/pulumicost-spec/issues/112)) ([b14dd3c](https://github.com/rshade/pulumicost-spec/commit/b14dd3c2eface44181d13d5add225eff57f53198)), closes [#85](https://github.com/rshade/pulumicost-spec/issues/85)
- **testing:** add metrics tracking example for EstimateCost (T041) ([#111](https://github.com/rshade/pulumicost-spec/issues/111)) ([944e078](https://github.com/rshade/pulumicost-spec/commit/944e0789a67be8e204d92d0d7a35ba181f9dd854)), closes [#84](https://github.com/rshade/pulumicost-spec/issues/84)
- **testing:** implement Plugin Conformance Test Suite ([#109](https://github.com/rshade/pulumicost-spec/issues/109)) ([03116ce](https://github.com/rshade/pulumicost-spec/commit/03116cef17567bdea85ba59e87aca322d2c42efb))

### Documentation

- **006-estimate-cost:** update data-model.md with actual decimal type (T054) ([#114](https://github.com/rshade/pulumicost-spec/issues/114)) ([45f4b2e](https://github.com/rshade/pulumicost-spec/commit/45f4b2e7c2a37c9414aada68343731a2f0e7913c)), closes [#89](https://github.com/rshade/pulumicost-spec/issues/89)

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
