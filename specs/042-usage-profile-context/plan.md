# Implementation Plan: Usage Profile Context

**Branch**: `042-usage-profile-context` | **Date**: 2026-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/042-usage-profile-context/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command.
See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add a `UsageProfile` enum to `GetProjectedCostRequest` and `GetRecommendationsRequest` messages, enabling
the Core to signal workload intent (DEV, PROD, BURST) to plugins for context-aware cost estimation and
recommendations. The SDK will provide profile-aware builder methods following existing patterns.

## Technical Context

**Language/Version**: Go 1.25.6 (per go.mod)
**Primary Dependencies**: Protocol Buffers v3, buf v1.32.1, google.golang.org/protobuf, google.golang.org/grpc, zerolog
**Storage**: N/A (stateless proto definitions and SDK helpers)
**Testing**: go test, conformance tests in sdk/go/testing/
**Target Platform**: gRPC plugins (Linux, macOS, Windows)
**Project Type**: Protocol specification + SDK
**Performance Goals**: SDK helpers <15 ns/op, 0 allocs/op (matching existing patterns)
**Constraints**: Backward compatibility required - existing plugins must work unchanged
**Scale/Scope**: 2 request messages extended, 1 enum added, SDK builder methods added

## Constitution Check (Pre-Design)

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First (I)**: Proto definitions updated before SDK implementation
- [x] **Spec Consumes (III)**: No complex pricing logic - plugins decide profile meaning
- [x] **Multi-Provider (II)**: Profile enum is provider-agnostic (DEV/PROD/BURST applies universally)
- [x] **FinFocus Alignment (VII)**: Follows existing naming conventions (UsageProfile)
- [x] **SDK Synchronization (XIII)**: TypeScript SDK update required after proto changes
- [x] **Test-First (V)**: Conformance tests required for profile handling
- [x] **Backward Compatibility (VI)**: UNSPECIFIED default preserves existing behavior
- [x] **Performance (VIII)**: Builder methods will use zero-allocation patterns

## Project Structure

### Documentation (this feature)

```text
specs/042-usage-profile-context/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
├── costsource.proto     # Extend GetProjectedCostRequest and GetRecommendationsRequest
└── enums.proto          # Add UsageProfile enum (per research.md decision)

sdk/go/
├── proto/finfocus/v1/   # Generated code (regenerated via make generate)
├── pluginsdk/
│   ├── usage_profile.go # Profile-aware builder helpers (new file)
│   └── usage_profile_test.go # Tests for profile helpers
└── testing/
    └── usage_profile_conformance_test.go # Conformance tests

sdk/typescript/          # TypeScript SDK updates (Constitution XIII)
├── src/client/          # Client wrapper updates
└── src/proto/           # Regenerated proto bindings
```

**Structure Decision**: Single-project protocol specification repository with Go and TypeScript SDKs.
Proto definitions in `proto/`, Go SDK in `sdk/go/`, TypeScript SDK in `sdk/typescript/`.

## Complexity Tracking

> No constitution violations requiring justification. This is a straightforward enum addition
> following established patterns (similar to GrowthType, RecommendationCategory).

## Constitution Check (Post-Design)

_Re-evaluated after Phase 1 design completion._

- [x] **Contract First (I)**: ✅ Proto contract defined in `contracts/usage_profile.proto`
- [x] **Spec Consumes (III)**: ✅ No pricing logic - profile is context hint only
- [x] **Multi-Provider (II)**: ✅ Enum values (DEV/PROD/BURST) are provider-agnostic
- [x] **FinFocus Alignment (VII)**: ✅ Naming follows existing patterns (UsageProfile, USAGE_PROFILE_*)
- [x] **SDK Synchronization (XIII)**: ✅ TypeScript SDK update planned in project structure
- [x] **Test-First (V)**: ✅ Conformance tests defined in quickstart.md examples
- [x] **Backward Compatibility (VI)**: ✅ UNSPECIFIED=0 default, unknown values handled gracefully
- [x] **Performance (VIII)**: ✅ Zero-allocation validation pattern documented in research.md
- [x] **Documentation (VII)**: ✅ Quickstart guide with complete examples created

**Post-Design Status**: All constitution checks pass. Ready for task generation via `/speckit.tasks`.

## Generated Artifacts

| Artifact | Path | Status |
|----------|------|--------|
| Implementation Plan | `specs/042-usage-profile-context/plan.md` | ✅ Complete |
| Research | `specs/042-usage-profile-context/research.md` | ✅ Complete |
| Data Model | `specs/042-usage-profile-context/data-model.md` | ✅ Complete |
| Proto Contract | `specs/042-usage-profile-context/contracts/usage_profile.proto` | ✅ Complete |
| Quickstart Guide | `specs/042-usage-profile-context/quickstart.md` | ✅ Complete |

## Next Steps

1. Run `/speckit.tasks` to generate implementation tasks
2. Implement proto changes in `enums.proto` and `costsource.proto`
3. Regenerate Go SDK with `make generate`
4. Implement SDK helpers in `sdk/go/pluginsdk/usage_profile.go`
5. Add conformance tests
6. Update TypeScript SDK
