<!--
Sync Impact Report - Constitution v1.3.2
=======================================
Version Change: 1.3.1 → 1.3.2
Modified Principles:
  - Section XI (Mandatory Copyright Headers): Removed specific attribution text
    requirement ("Copyright 2026 PulumiCost/FinFocus Authors").
Removed Sections: None
Added Sections: None
Templates Requiring Updates: None
Follow-up TODOs: None
-->

# FinFocus Specification Constitution

## Core Principles

### I. gRPC Proto Specification-First Development (Contracts are Sacred)

Every change to the protocol MUST begin with protobuf specification updates before implementation:

- **Proto definitions are the source of truth**: All gRPC service definitions in
  `proto/finfocus/v1/` define the canonical contract.
- SDK code is generated from proto definitions; manual edits to generated code are FORBIDDEN.
- Proto changes require corresponding JSON schema updates for PricingSpec messages.
- All protobuf message types MUST have comprehensive validation.

**Rationale**: As a gRPC protobuf specification repository, the proto files define the wire protocol and
service contract. Stable, well-documented contracts are the foundation of the ecosystem.

### II. Multi-Provider gRPC Consistency

The gRPC service specification MUST maintain feature parity across all major cloud providers:

- New billing modes in PricingSpec messages MUST include cross-provider examples (AWS, Azure, GCP, Kubernetes).
- Examples MUST demonstrate real-world provider use cases as gRPC request/payloads.
- ResourceDescriptor message fields MUST be provider-agnostic.

**Rationale**: Plugin developers implement a universal interface. Provider-specific leakage fragments the
ecosystem and breaks interoperability.

### III. The Spec Consumes, It Does Not Calculate

The specification and its implementing plugins are NOT responsible for complex pricing logic (e.g., tiered
pricing, committed-use discounts):

- **Data is Pre-Calculated**: The spec's role is to consume final, _adjusted_ costs from upstream providers.
- **No Complex Math**: Avoid embedding discount calculators or tiered-pricing engines in the SDK or plugins.
- **Standardized Model**: Focus on providing a standardized model for the final cost data.

**Rationale**: Complex pricing logic belongs to specialized data providers (e.g., Kubecost, Vantage). The
spec's role is standardized observability, not financial calculation.

### IV. Strict Separation of Concerns

This repository defines the _specification_ and foundational tooling, not the end-user application:

- **Spec vs. Core**: `finfocus-spec` defines interfaces; `finfocus-core` contains application logic.
- **Plugin SDK focus**: The SDK is for plugin creators, not end-users.
- **Minimal Dependencies**: Core spec files MUST NOT depend on higher-level application components.

**Rationale**: Maintaining a clean separation ensures the specification remains a portable foundation for
diverse implementations.

### V. Test-First Protocol (NON-NEGOTIABLE)

TDD is mandatory for all gRPC specification changes:

1. Write conformance tests defining expected gRPC behavior (request → response).
2. Tests MUST fail against current proto/implementation.
3. Update proto definitions to make tests pass.
4. Regenerate SDK and validate all examples.

**Rationale**: gRPC protocol changes have high downstream impact. Tests define the RPC contract before
implementation to prevent breaking existing clients.

### VI. Protobuf Backward Compatibility

Breaking changes to protobuf definitions are strictly controlled:

- MAJOR version bump required for breaking proto changes (field removals, type changes, renaming).
- buf breaking change detection MUST pass in CI.
- Deprecated protobuf fields MUST remain for one MAJOR version.
- Use `reserved` keyword for removed fields to prevent field number reuse.

**Rationale**: Protobuf wire format compatibility is critical. Breaking changes cascade through all plugin
implementations and client applications.

### VII. Comprehensive Documentation & Identity Transition

Every gRPC specification element MUST be documented, and the project's identity must be preserved:

- **FinFocus Identity**: The project has been renamed to **FinFocus** to align with the
  industry-standard FinOps FOCUS specification.
- Proto messages and fields require inline comments for documentation generation.
- **Documentation Currency**: Documentation MUST be updated in the same PR as feature implementation. Stale
  docs are blocking for PR approval.
- Root `README.md`, `docs/`, and SDK `README.md` MUST stay in sync.

**Rationale**: gRPC specifications are only useful if understood. The rename ensures alignment with the
industry-standard FinOps FOCUS specification.

### VIII. Performance as a gRPC Requirement (Performance is Paramount)

Code, especially within the Go SDK, must be highly performant and memory-efficient:

- **Zero-Allocation Goal**: Common operations (validation, enum lookups) should aim for zero-allocation.
- Conformance tests include RPC response time requirements.
- Benchmarks are required for all new core SDK logic.

**Rationale**: Cloud cost data queries can involve large datasets. Inefficient designs degrade the entire
observability pipeline.

### IX. Observability & Validation (Maintainer Focused)

Observability features are for plugin maintainers, not necessarily end-users:

- **Logging and Metrics are Separate**: `zerolog` for structured events, Prometheus for time-series metrics.
- Metrics should be implemented as optional, distinct components (e.g., gRPC interceptors).
- **Validation Layers**: Multi-layer validation (Proto, Schema, Service, SDK, CI) ensures quality.

**Rationale**: Separating concerns allows maintainers to diagnose performance without cluttering core
business logic or over-complicating the end-user data model.

### X. Follow Established Patterns

New contributions MUST adhere to existing, documented patterns:

- **Standard Domain Enum Pattern**: Use the established pattern for high-performance validation.
- **Design Docs Required**: Significant changes MUST be proposed via a design spec in `specs/`.

**Rationale**: Consistency reduces cognitive load for maintainers and ensures the "zero-allocation" goals
are met across the codebase.

### XI. Mandatory Copyright Headers

Every source file (Go, Proto, Script, Schema) MUST include the standard Apache 2.0 copyright header:

- **License**: Explicitly state "Licensed under the Apache License, Version 2.0".
- **Persistence**: Headers must be maintained during the transition to FinFocus.

**Rationale**: Ensures license clarity even if files are separated from the repository and maintains a
professional enterprise standard aligned with Apache 2.0 recommendations.

## Governance

### Amendment Process

Constitution changes require:

1. Proposal documenting rationale and impact on gRPC development workflow.
2. Version bump per semantic versioning (MAJOR/MINOR/PATCH).
3. Update to all dependent template files.
4. Sync Impact Report prepended to constitution.

### Compliance Review

All PRs MUST verify:

- **gRPC proto-first approach**: Proto updated before implementation.
- **Spec Consumes**: No complex pricing logic added to the spec/SDK.
- **Test-first protocol**: Conformance tests written and passed.
- **Backward compatibility**: buf breaking check passes.
- **Documentation complete**: Updated root docs and examples.

### Runtime Development Guidance

For day-to-day development guidance, refer to `CLAUDE.md` or `GEMINI.md` in the repository root. The
constitution defines non-negotiable principles; these files provide practical workflow tips.

**Version**: 1.3.2 | **Ratified**: 2025-08-11 | **Last Amended**: 2026-01-11
