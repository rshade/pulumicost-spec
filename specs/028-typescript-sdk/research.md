# Research & Technical Decisions

**Feature**: TypeScript Client SDK (`028-typescript-sdk`)
**Status**: Phase 0 Complete

## 1. Protocol Implementation

**Decision**: Use `@connectrpc/connect` ecosystem.

- Browser: `@connectrpc/connect-web` with `createConnectTransport` (uses `fetch`).
- Node: `@connectrpc/connect-node` with `createGrpcTransport` or `createConnectTransport` (uses `http2` or `http`).
- Protobuf: `@bufbuild/protobuf` for message manipulation.

**Rationale**:

- Connect is the modern standard, replacing `grpc-web`.
- Supports JSON/HTTP natively (critical for "browser compatibility" req).
- Fully type-safe with generated code.

## 2. Monorepo Structure

**Decision**: Use `npm` workspaces (or `pnpm` if available in environment) with a `packages/` directory.

- `packages/client`: Core logic, interfaces, builders. Browser compatible. Dependencies: `connect`, `protobuf`.
- `packages/middleware`: Node.js specific. Depends on `client`. Provides the "REST Gateway" functionality.
- `packages/frameworks`: Specific adapters for Express/Fastify/NestJS. Depends on `middleware` or `client`.

**Rationale**:

- Ensures strict separation of dependencies.
- `client` remains lightweight (<40KB) by excluding server-side frameworks.
- Allows shared development and versioning.

## 3. REST Wrapper (Gateway) Design

**Decision**: The REST Wrapper will act as a **Gateway/Proxy**.

- It accepts standard JSON HTTP requests.
- It translates them into strongly-typed RPC calls using the `CostSourceClient`.
- It returns the RPC response as JSON.
- **Mapping**:
  - Defaults to RPC-style `POST /api/v1/CostSource/GetActualCost`.
  - Can optionally support RESTful mapping (e.g., `GET /api/v1/costs/actual`) if configured, but RPC-style is
    native to Connect and easiest to maintain.
  - Given Spec "FR-006: converts 22 HTTP endpoints", we will support the standard Connect HTTP binding paths.

**Rationale**:

- Simplifies implementation.
- Connect protocol *is* HTTP/JSON friendly. The wrapper essentially forwards the request but handles the
  client-side authentication/connection logic to the remote plugin.

## 4. Framework Integration

**Decision**:

- **Express/Fastify**: Middleware that mounts the Gateway.
- **NestJS**: A `FinFocusModule` that registers the `CostSourceClient` as a provider and optionally mounts the Controller/Gateway.

**Rationale**: Native patterns for each framework ensure seamless adoption.

## 5. Build & Generate

**Decision**:

- Use `buf generate` to generate TS code into `packages/client/src/generated`.
- Check in generated code (per Spec FR-009) to avoid build-time dependency on `buf` for consumers.

**Rationale**: Compliance with FR-009.
