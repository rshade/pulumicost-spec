# Implementation Tasks: TypeScript Client SDK

**Feature**: `028-typescript-sdk`
**Status**: Pending
**Total Tasks**: 24

## Dependencies

1.  **Phase 1: Setup** (Project initialization)
2.  **Phase 2: Foundation** (Generated Code & Shared Types)
3.  **Phase 3: Browser Client (US1)** (Core Client, Builders, Validation)
4.  **Phase 4: REST API Integration (US2)** (Node Middleware)
5.  **Phase 5: Framework Integration (US3)** (Express, Fastify, NestJS)
6.  **Phase 6: Polish** (Docs, Bundle Size)

## Phase 1: Setup & Configuration
*Goal: Initialize the monorepo structure and build tooling.*

- [X] T001 Initialize npm workspaces configuration in `sdk/typescript/package.json`
- [X] T002 Create core client package structure in `sdk/typescript/packages/client/package.json`
- [X] T003 Configure TypeScript base config in `sdk/typescript/tsconfig.base.json`
- [X] T004 Install core dependencies (`@connectrpc/connect`, `@bufbuild/protobuf`) in `sdk/typescript/packages/client/`

## Phase 2: Foundation (Blocking)
*Goal: Generate Protobuf code and establish shared error handling.*

- [X] T005 Configure Buf generation for TypeScript in `sdk/typescript/buf.gen.yaml`
- [X] T006 Run buf generate to create proto artifacts in `sdk/typescript/packages/client/src/generated/`
- [X] T007 Implement `ValidationError` class in `sdk/typescript/packages/client/src/errors/validation-error.ts`
- [X] T008 [P] Export all enum types (BillingMode, Provider, etc.) in `sdk/typescript/packages/client/src/index.ts`

## Phase 3: User Story 1 - Browser Application Developer (P1)
*Goal: Enable browser-based cost fetching and optimization using the Core Client.*
*Independent Test: Verify `CostSourceClient` fetches data against MSW mock.*

- [X] T009 [P] [US1] Implement `ResourceDescriptorBuilder` with fluent API in `sdk/typescript/packages/client/src/builders/resource-descriptor.ts`
- [X] T010 [P] [US1] Implement `RecommendationFilterBuilder` in `sdk/typescript/packages/client/src/builders/recommendation-filter.ts`
- [X] T011 [P] [US1] Implement `FocusRecordBuilder` for FOCUS 1.2/1.3 in `sdk/typescript/packages/client/src/builders/focus-record.ts`
- [X] T012 [US1] Implement `CostSourceClient` wrapper with validation in `sdk/typescript/packages/client/src/clients/cost-source.ts`
- [X] T013 [US1] Implement `RecommendationsIterator` for async pagination in `sdk/typescript/packages/client/src/utils/pagination.ts`
- [X] T014 [US1] Implement `ObservabilityClient` and `RegistryClient` wrappers in `sdk/typescript/packages/client/src/clients/auxiliary.ts`
- [X] T015 [US1] Create MSW handlers for integration testing in `sdk/typescript/packages/client/test/mocks/handlers.ts`
- [X] T016 [US1] Write integration tests for full client workflow in `sdk/typescript/packages/client/test/integration.test.ts`

## Phase 4: User Story 2 - REST API Integration (P2)
*Goal: Expose plugin capabilities via standard HTTP endpoints using Node.js middleware.*
*Independent Test: Verify generic middleware translates JSON requests to RPC.*

- [X] T017 [US2] Initialize middleware package in `sdk/typescript/packages/middleware/package.json`
- [X] T018 [US2] Implement Node.js Connect transport factory in `sdk/typescript/packages/middleware/src/transport.ts`
- [X] T019 [US2] Implement generic `RESTGateway` logic (JSON->RPC proxy) in `sdk/typescript/packages/middleware/src/gateway.ts`

## Phase 5: User Story 3 - Framework Integration (P3)
*Goal: Provide native adapters for popular Node.js frameworks.*
*Independent Test: Verify Express app handles requests via the adapter.*

- [X] T020 [US3] Initialize framework plugins package in `sdk/typescript/packages/framework-plugins/package.json`
- [X] T021 [P] [US3] Implement Express adapter in `sdk/typescript/packages/framework-plugins/src/express/index.ts`
- [X] T022 [P] [US3] Implement Fastify adapter in `sdk/typescript/packages/framework-plugins/src/fastify/index.ts`
- [X] T023 [P] [US3] Implement NestJS module and service in `sdk/typescript/packages/framework-plugins/src/nestjs/index.ts`

## Phase 6: Polish & Cross-Cutting
*Goal: Finalize documentation and bundle size checks.*

- [X] T024 Ensure core client bundle size is under 40KB (minified) in `sdk/typescript/packages/client/`

## Implementation Strategy
- **MVP (Phase 1-3)**: Deliver the browser-compatible Core Client. This unblocks the primary use case (Dashboards).
- **Expansion (Phase 4-5)**: Add Node.js server-side support.
- **Parallelism**: Builders (T009-T011) and Framework adapters (T021-T023) can be implemented in parallel.
