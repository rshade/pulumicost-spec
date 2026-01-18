# Implementation Plan: TypeScript Client SDK

**Branch**: `028-typescript-sdk` | **Date**: 2026-01-16 | **Spec**: [specs/028-typescript-sdk/spec.md](specs/028-typescript-sdk/spec.md)
**Input**: Feature specification from `specs/028-typescript-sdk/spec.md`

## Summary

Implement the official TypeScript SDK (`finfocus-client`) for the FinFocus ecosystem. This SDK will serve as the
primary browser-side integration point, leveraging the Connect protocol (JSON/HTTP) to communicate with plugins. It
includes a core client, builder patterns for complex objects, validations, a REST API wrapper, and integration plugins
for major Node.js frameworks (Express, Fastify, NestJS).

## Technical Context

**Language/Version**: TypeScript 5.0+ (Target: ES2018)
**Primary Dependencies**:

- Core: `@connectrpc/connect`, `@connectrpc/connect-web`, `@bufbuild/protobuf`
- Validation: `zod` (internal usage for runtime checks if needed, or manual)
- Testing: `vitest`, `msw` (Mock Service Worker)
**Storage**: N/A
**Target Platform**: Browser (Chrome 60+, Firefox 55+, Safari 12+) and Node.js 18+
**Project Type**: TypeScript Monorepo (pnpm/npm workspaces) located in `sdk/typescript`
**Performance Goals**: Core bundle < 40KB (minified + gzipped)
**Constraints**: Zero external polyfills for target browsers, strictly typed APIs.
**Scale/Scope**: ~11 RPCs (Core), ~22 RPCs total (Admin/Registry included).

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First**: Implementation based on existing `proto/finfocus/v1` definitions.
- [x] **Spec Consumes**: SDK consumes APIs; includes no pricing logic.
- [x] **Multi-Provider**: SDK is agnostic; generic `ResourceDescriptor` used.
- [x] **FinFocus Alignment**: Uses "FinFocus" naming and FOCUS 1.2/1.3 standards.

## Project Structure

### Documentation (this feature)

```text
specs/028-typescript-sdk/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (OpenAPI for REST wrapper)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
sdk/typescript/
├── packages/
│   ├── client/              # Core Client (Browser compatible)
│   │   ├── src/
│   │   │   ├── client.ts    # Main entry
│   │   │   ├── generated/   # Buf generated code
│   │   │   └── builders/    # Fluent builders
│   │   └── test/
│   ├── middleware/          # REST API Wrapper (Node.js)
│   │   ├── src/
│   │   └── test/
│   └── framework-plugins/   # Express/Fastify/NestJS adapters
│       ├── src/
│       └── test/
├── package.json             # Workspace root
└── tsconfig.base.json
```

**Structure Decision**: Monorepo structure within `sdk/typescript` to separate the browser-light core client from
Node.js-specific middleware and framework plugins, ensuring the <40KB bundle size requirement is met for the core.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
| --------- | ---------- | ------------------------------------ |
| Monorepo  | Separation of browser/node concerns | Single package would bloat browser bundle with Node deps (Express, etc.) |
