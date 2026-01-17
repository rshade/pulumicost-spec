# Implementation Plan - Dual-Layer Capability Discovery

**Status**: Planned
**Feature**: Dual-Layer Capability Discovery
**Branch**: 037-capability-discovery
**Spec**: [specs/037-capability-discovery/spec.md](spec.md)

## Technical Context

This feature implements a "Dual-Layer" capability discovery mechanism to optimize Host-Plugin communication:

1. **Global Discovery**: Host asks "What major features (Budgets, Recommendations) do you support
   overall?" (at startup).
2. **Granular Discovery**: Host asks "Do you support Feature X for *this specific resource*?" (at
   runtime).

The goal is to prevent the Host from calling RPCs that are not supported, either globally or for
specific legacy resources.

**Key Technical Decisions**:

- **Protocol**: gRPC with Protobuf.
- **Granular Input**: `ResourceDescriptor` (Provider, Service, Type, Region) is passed to `Supports`
  to determine support.
- **Granular Output**: A strict `repeated PluginCapability capabilities_enum` list.
- **SDK Behavior**: "Opt-Out" auto-discovery. The SDK reflects on implemented Go interfaces to
  populate capabilities by default.
- **Fallback**: "Inherit Global". If granular list is empty, Host assumes global capabilities apply.
- **Enforcement**: Plugin rejects unsupported calls with `Unimplemented`.

**Changes Required**:

- **Protobuf**: Update `SupportsRequest` (to accept `ResourceDescriptor`) and `SupportsResponse`
  (to add `capabilities_enum`).
- **Go SDK**:
  - Implement reflection-based auto-discovery in `Serve` or `Plugin` initialization.
  - Update `Supports` method signature and logic.
  - Add `ResourceDescriptor` handling.

**Dependencies**:

- `google.golang.org/protobuf/proto`
- Existing `ResourceDescriptor` message definition.
- Existing `PluginCapability` enum.

## Constitution Check

| Principle | Status | Notes |
| :--- | :--- | :--- |
| **I. Proto-First** | ✅ | Plan starts with Protobuf changes (`SupportsRequest`, `SupportsResponse`). |
| **II. Multi-Provider** | ✅ | `ResourceDescriptor` is provider-agnostic. |
| **III. Spec Consumes** | ✅ | Discovery only; no pricing logic involved. |
| **IV. Separation of Concerns** | ✅ | Defines the *protocol* for discovery; Host logic is outside scope. |
| **V. Test-First** | ✅ | Tests will define `Supports` behavior before implementation. |
| **VI. Backward Compatibility** | ✅ | Legacy string `capabilities` map preserved. "Inherit Global" fallback ensures old plugins work. |
| **VII. Documentation** | ✅ | Docs will be updated. |
| **VIII. Performance** | ✅ | `ResourceDescriptor` is lightweight. Enum is efficient. |
| **IX. Observability** | N/A | No specific observability changes. |
| **X. Patterns** | ✅ | Uses Standard Domain Enum pattern (`PluginCapability`). |
| **XI. Copyright** | ✅ | Standard headers will be applied. |
| **XII. Auto-Discovery** | ✅ | SDK implements interface-driven discovery (Opt-Out). |

## Phase 0: Research & Validation

### 0.1 Research Questions

- [x] **Protobuf Compatibility**: Verify if modifying `SupportsRequest` to take `ResourceDescriptor`
  (instead of just `string resource_id` if that was the case, or adding fields) breaks existing
  clients. *Constraint: Must be backward compatible.*
- [x] **SDK Reflection Performance**: Assess the performance impact of reflecting on interfaces
  during `Serve` initialization (should be negligible as it's one-time).
- [x] **ResourceDescriptor Source**: Confirm where the Host gets the `ResourceDescriptor` to pass
  to `Supports`. (Likely from the `Resource` object it's processing).

### 0.2 Design Decisions

- **Decision 1**: `SupportsRequest` will likely need a new field or we need to check if existing
  fields can carry `ResourceDescriptor`. *Assume we add a new field `resource_descriptor`.*
- **Decision 2**: SDK will use `reflect.Type` to check for method implementations (e.g.,
  `GetRecommendations`) on the struct passed to `Serve`.

## Phase 1: Specifications & Contracts

### 1.1 Data Model (`data-model.md`)

- Define `SupportsRequest` structure (incorporating `ResourceDescriptor`).
- Define `SupportsResponse` structure (adding `capabilities_enum`).
- Define the `ResourceDescriptor` fields relevant for discovery.

### 1.2 API Contracts (`contracts/`)

- **File**: `proto/finfocus/v1/costsource.proto` (or where `Supports` is defined).
- **Update**:
  - `message SupportsRequest { ... ResourceDescriptor resource = X; ... }`
  - `message SupportsResponse { ... repeated PluginCapability capabilities_enum = Y; ... }`

### 1.3 Documentation

- Update `quickstart.md` to explain how to opt-out of capabilities if needed.

## Phase 2: Implementation Plan

### 2.1 Proto & Spec

1. **Modify Proto**: Add `resource` to `SupportsRequest` and `capabilities_enum` to `SupportsResponse`
   in `proto/finfocus/v1/*.proto`.
2. **Generate Go Code**: Run `buf generate`.
3. **Validate**: Ensure no breaking changes reported by `buf`.

### 2.2 Go SDK Implementation

1. **Refactor Plugin Registration**: In `server.go` (or equivalent), add logic to inspect the user's struct.
    - Check for `BudgetsProvider`, `RecommendationsProvider` interfaces.
    - Populate a default `[]PluginCapability` list.
2. **Implement `Supports` Handler**:
    - Update the default `Supports` implementation in the SDK to return the auto-discovered list.
    - Allow user override (if the user implements `Supports` explicitly, they control the logic,
      but we need a helper to check "global vs granular"). *Actually, `Supports` is usually
      implemented by the user for granular logic. The SDK handles GLOBAL discovery via `GetPluginInfo`.*
    - **Correction**: The SDK needs to handle `GetPluginInfo` auto-discovery (Global). For
      `Supports` (Granular), the *user* implements the logic, but the SDK should provide types/helpers.
    - **Refinement**:
        - `GetPluginInfo` (Global): SDK auto-discovers based on interfaces.
        - `Supports` (Granular): User implements this. If they don't, SDK default implementation
          returns "Supported=True, Capabilities=Global List".
        - The feature is about *granular* discovery. So we need to enable the user to say "Yes/No" per
          resource.

### 2.3 Host Integration (Mock/Test)

1. Create a test case mimicking the Host.
2. Call `GetPluginInfo` -> Assert global caps.
3. Call `Supports(resource)` -> Assert granular caps.

## Phase 3: Verification

### 3.1 Tests

- **Conformance Tests**: Update suite to verify `Supports` returns correct enum values.
- **Unit Tests**: Test SDK reflection logic with various struct combinations.
- **Backward Compat**: Verify legacy string map is still populated.
