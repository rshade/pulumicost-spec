# Feature Specification: TypeScript Client SDK for FinFocus Plugin Ecosystem

**Feature Branch**: `028-typescript-sdk`
**Created**: 2026-01-16
**Status**: Draft
**Input**: Comprehensive TypeScript SDK implementation plan for browser-first
Connect protocol support with REST API wrapper and framework plugins

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Browser Application Developer (Priority: P1)

A frontend developer building a cloud cost dashboard needs to integrate with a
FinFocus cost source plugin running in a browser environment. They need to fetch
actual costs, project future costs, and get optimization recommendations without
setting up backend infrastructure.

**Why this priority**: Core value proposition of the SDK - enabling browser-first
plugin integration. Unblocks entire client-side use case.

**Independent Test**: Core client package with all 11 CostSourceClient RPCs can
be fully tested in browser environment using MSW mock server, delivering complete
cost analysis workflow.

**Acceptance Scenarios**:

1. **Given** browser environment with plugin URL, **When** developer creates
   `CostSourceClient` and calls `getActualCost()`, **Then** returns valid cost
   data for specified resource
2. **Given** ResourceDescriptorBuilder instance, **When** developer chains methods
   (`.withProvider()`, `.withResourceType()`, `.build()`), **Then** produces valid
   ResourceDescriptor message
3. **Given** pagination response with `nextPageToken`, **When** developer uses
   `RecommendationsIterator`, **Then** automatically fetches next pages and
   aggregates results
4. **Given** invalid request (missing required fields), **When** validator runs,
   **Then** throws descriptive ValidationError before RPC call

---

### User Story 2 - REST API Integration (Priority: P2)

A service provider wants to expose FinFocus plugin capabilities via traditional
REST/JSON HTTP endpoints instead of Connect protocol. They need a middleware that
translates HTTP requests to Connect RPC calls and back.

**Why this priority**: Enables integration with non-browser clients and legacy
REST-only systems. Significantly broadens accessibility of plugins.

**Independent Test**: REST wrapper package with 22 HTTP endpoints can be tested
independently using HTTP client, transforming JSON payloads to protobuf and vice
versa.

**Acceptance Scenarios**:

1. **Given** Express middleware configured with plugin URL, **When** POST request
   to `/v1/cost/actual` with JSON body, **Then** response contains cost data as
   JSON
2. **Given** invalid request body, **When** middleware validates, **Then** returns
   HTTP 400 with descriptive error message
3. **Given** plugin connection failure, **When** middleware handles error, **Then**
   returns HTTP 503 with appropriate status code

---

### User Story 3 - Framework Integration (Priority: P3)

Application developers using popular Node.js frameworks need simplified
integration patterns. They want to use framework-specific decorators, modules, or
middleware without deep knowledge of Connect protocol details.

**Why this priority**: Reduces integration friction for framework users. Each
framework gets native integration pattern. Not blocking initial release but
valuable for ecosystem maturity.

**Independent Test**: Each framework plugin (Express, Fastify, NestJS) provides
working example application that demonstrates all 22 RPC methods accessible via
framework-native patterns.

**Acceptance Scenarios**:

1. **Given** Express app with finfocus middleware registered, **When** HTTP
   request arrives, **Then** processed through standard Express middleware chain
2. **Given** NestJS app with FinFocusModule imported, **When** service injects
   FinFocusService, **Then** has access to all client methods
3. **Given** framework plugin configured with plugin URL, **When** multiple
   concurrent requests arrive, **Then** handled correctly without state pollution

---

### Edge Cases

- What happens when plugin is unreachable? → Client throws descriptive error with
  timeout information
- How does SDK handle partial responses? → Validators catch incomplete messages
  and throw ValidationError
- What if browser doesn't support ES2018 features? → Bundle size grows ~10% with
  polyfills; document compatibility requirements
- How to handle CORS when plugin on different origin? → SDK uses Connect
  JSON-HTTP; document CORS header requirements in configuration
- What if recommendation filter references non-existent fields? → Validator
  validates against known fields enum; throws error with suggestions
- How to handle cost currency mismatches? → Validator validates ISO 4217 codes;
  clients must handle conversion or throw error

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: Core client MUST implement all 22 RPC methods across 3 services
  (11 CostSourceService, 3 ObservabilityService, 8 PluginRegistryService)
- **FR-002**: SDK MUST support Connect protocol with JSON/HTTP encoding for
  browser compatibility (no binary gRPC-Web)
- **FR-003**: SDK MUST provide builder pattern for ResourceDescriptor (10 methods),
  RecommendationFilter (16 filters), and FocusRecordBuilder (FOCUS 1.2/1.3)
- **FR-004**: SDK MUST validate all requests against resource descriptor rules,
  response structure, FOCUS compliance, and ISO 4217 currency codes
- **FR-005**: SDK MUST include comprehensive helper utilities: cost calculations
  (hourly↔monthly, growth projections), 180+ ISO currency validation, 44+
  billing mode enums, async pagination iterator
- **FR-006**: SDK MUST provide REST API wrapper middleware converting 22 HTTP
  endpoints to Connect RPC calls with JSON↔protobuf transformation
- **FR-007**: SDK MUST provide Express, Fastify, and NestJS framework plugins
  with native integration patterns (middleware, modules, decorators)
- **FR-008**: SDK MUST compile to ES2018 target with ESM and CommonJS dual build
  outputs for maximum browser compatibility
- **FR-009**: SDK MUST include generated protobuf and Connect code (version-
  controlled in src/generated/) enabling fast npm installs without post-install
  build
- **FR-010**: SDK MUST export all client classes, builders, validators, and
  helpers from single public index.ts file for clean API surface
- **FR-011**: SDK MUST generate API documentation with JSDoc for all public
  methods including performance SLAs and usage examples
- **FR-012**: Core client testing MUST achieve 80%+ code coverage with Vitest
  unit tests and MSW mock server integration tests
- **FR-013**: Bundle size MUST remain under 40 KB (minified + gzipped) for core
  client only
- **FR-014**: SDK MUST support pagination for large result sets via async iterator
  with configurable page sizes and max page limits

### Key Entities

- **ResourceDescriptor**: Message representing a cloud resource with provider,
  resourceType, SKU, region, tags, ARN, utilization, growth parameters. Builder
  has 10 fluent methods for construction.
- **CostSourceClient**: Main client class wrapping CostSourceService with 11 RPC
  methods: Name, Supports, GetActualCost, GetProjectedCost, GetPricingSpec,
  EstimateCost, GetRecommendations, DismissRecommendation, GetBudgets,
  GetPluginInfo, DryRun.
- **ObservabilityClient**: Client for 3 RPC methods: HealthCheck, GetMetrics,
  GetServiceLevelIndicators supporting plugin observability and monitoring.
- **RegistryClient**: Client for 8 RPC methods: DiscoverPlugins, GetPluginManifest,
  ValidatePlugin, InstallPlugin, UpdatePlugin, RemovePlugin, ListInstalledPlugins,
  CheckPluginHealth supporting plugin lifecycle management.
- **RecommendationFilterBuilder**: Builder with 16 filter fields (provider, region,
  resourceType, category, actionType, priority, minEstimatedSavings, etc.) for
  filtering recommendations.
- **FocusRecordBuilder**: Builder for FOCUS 1.2/1.3 cost records with allocation
  fields (AllocatedMethodId, AllocatedResourceId), service provider fields, contract
  tracking.
- **RESTMiddleware**: Express/Fastify middleware factory that translates 22 HTTP
  endpoints to Connect RPC calls with JSON↔protobuf transformation and error
  mapping.
- **ValidationError**: Domain-specific error class extending Error with field name
  and error code for programmatic error handling.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: All 22 RPC methods have working implementations with request/response
  validation passing 100% of test cases
- **SC-002**: Core client bundle size ≤ 40 KB (minified + gzipped) with ES2018
  compilation
- **SC-003**: SDK works in Chrome 60+, Firefox 55+, Safari 12+ without external
  polyfills (or with minimal optional polyfills)
- **SC-004**: 5-line quick start example in README demonstrates functional cost
  retrieval workflow
- **SC-005**: All RPC methods have JSDoc documentation with parameter descriptions,
  return types, performance SLAs, and usage examples
- **SC-006**: Core client tests achieve ≥80% code coverage measured by Vitest with
  both unit and integration tests
- **SC-007**: REST wrapper converts all 22 RPCs to HTTP endpoints with correct
  request/response transformation
- **SC-008**: Three framework plugins (Express, Fastify, NestJS) each have working
  example applications demonstrating all RPC methods
- **SC-009**: Builder classes support method chaining with fluent API and throw
  descriptive errors on invalid configurations
- **SC-010**: Validator functions check all FOCUS compliance rules, ISO 4217
  currency codes, and resource descriptor constraints
- **SC-011**: Helper utilities provide cost calculations, pagination, filtering,
  and sorting with comprehensive examples
- **SC-012**: npm package publishes successfully as `finfocus-client` v0.1.0 with
  all dependencies properly specified
- **SC-013**: Documentation includes comprehensive guides for: installation, core
  concepts, API reference, advanced usage, REST integration, and framework-specific
  patterns
- **SC-014**: Project passes all CI/CD checks: linting, type checking, tests,
  bundle size validation, and documentation generation
