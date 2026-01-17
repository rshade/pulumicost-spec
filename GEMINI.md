# GEMINI.md - Guiding Principles for finfocus-spec

This document outlines the core principles, architectural guidelines, and development philosophy
for the `finfocus-spec` repository, based on an analysis of existing documentation and direct
feedback. It serves as a constitution to ensure all future contributions are aligned with the
project's vision.

## 1. Project Vision

The `finfocus-spec` repository aims to define the **universal, open-source standard for cloud
cost observability**. It provides the foundational contracts (schemas, Protobufs) and developer
tools to create a robust ecosystem of cost-estimation plugins.

## 2. Core Principles

- **Performance is Paramount:** Code, especially within the Go SDK, must be highly performant and
  memory-efficient. A "zero-allocation" goal for common operations is the standard.
  All new code must be benchmarked.
- **Contracts are Sacred:** The Protobuf definitions and JSON schemas are the source of truth.
  They must be stable, well-documented, and evolved carefully through the established design
  spec process.
- **Developer Experience (DX) for Plugin Creators:** The primary audience for the SDK is the
  plugin developer. The SDK should provide simple, consistent, and powerful building blocks
  that make creating high-quality plugins as easy as possible.
- **Strict Separation of Concerns:** This repository defines the _specification_ and foundational
  tooling. It is not a monolithic application.
  - `finfocus-spec`: Defines the interfaces and data schemas. Provides SDKs for implementation.
  - `finfocus-core`: (Separate repo) Contains higher-level application logic, such as the
    public-facing Plugin Registry service.
  - `finfocus-plugins-*`: (Separate repos) Individual plugins that implement the spec.

## 3. Architectural & Development Guidelines

- **The Spec Consumes, It Does Not Calculate:** The `finfocus-spec` and the plugins that directly
  implement it are not responsible for complex pricing logic (e.g., tiered pricing, committed-use
  discounts). This logic belongs to upstream data providers (like Kubecost, Vantage, Flexera, etc.).
  The spec's role is to consume the final, _adjusted_ cost from these services and provide a
  standardized model for it.
- **Observability is for Maintainers:** Features like metrics (Prometheus) are intended for plugin
  maintainers to diagnose performance and efficiency. They are not primarily for end-users of the
  cost data. Therefore, such features should be implemented as optional, distinct components (e.g.,
  a separate gRPC interceptor) rather than being deeply integrated into core logic like logging.
- **Logging and Metrics are Separate:** `zerolog` is for structured, event-based logging. Prometheus
  is for aggregated, time-series metrics. These serve different purposes and should remain separate
  concerns in the SDK. The existing logging pattern is the standard.
- **Follow Established Patterns:** New contributions must adhere to existing, documented patterns,
  such as the "Standard Domain Enum Pattern" used in the Go SDK for high-performance,
  zero-allocation validation.
- **Changes Require Design Docs:** Significant changes or new features must be proposed and
  documented in a design specification under the `specs/` directory before implementation.

## Active Technologies

- TypeScript 5.0+ (Target: ES2018) (028-typescript-sdk)

- Go 1.25.5+, Protobuf 3 + `google.golang.org/grpc`, `google.golang.org/protobuf` (029-plugin-info-rpc)

- Go 1.25.5+, Protobuf 3 + `google.golang.org/protobuf`, `google.golang.org/grpc` (018-proto-add-arn)
- N/A (Proto definition) (018-proto-add-arn)

- Go 1.25.5+ + `google.golang.org/grpc`, `google.golang.org/protobuf` (existing in project) (017-pluginsdk-validation)

- Go 1.25.5+ (SDK), Protobuf 3 + `google.golang.org/grpc`, `google.golang.org/protobuf` (001-fallback-hint)
- N/A (API Specification) (001-fallback-hint)

## Recent Changes

- 001-fallback-hint: Added Go 1.25.5+ (SDK), Protobuf 3 + `google.golang.org/grpc`,
  `google.golang.org/protobuf`

## Common Issues & Solutions

- Issue: `make lint` and `make validate` may time out on this project.
  Solution: Run `golangci-lint run` directly for faster Go linting results, or `make test` for unit tests.

## Workflow Optimizations

- For CodeRabbit fixes: Always verify `git log` and file content first; reviews may reference older
  commits that have already been fixed by subsequent pushes.

## Project-Specific Patterns

- `pluginsdk.Serve`: Tests dealing with `Serve` should prefer injecting a `net.Listener` (via
  `ServeConfig.Listener`) rather than relying on `Port` and `listenOnLoopback` to avoid race
  conditions and ensure predictable port binding.
