# Implementation Plan: SDK Polish v0.4.15

**Branch**: `034-sdk-polish` | **Date**: 2026-01-10 | **Spec**: spec.md
**Input**: Feature specification from `/specs/034-sdk-polish/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See
`.specify/templates/commands/plan.md` for the execution workflow.

## Summary

SDK Polish v0.4.15 enhances the PulumiCost Go SDK with configurable per-client timeouts, user-friendly error
messages for GetPluginInfo RPC calls, and performance conformance testing. Technical approach involves
extending ClientConfig with timeout support, improving error message formatting in Server.GetPluginInfo,
and adding performance validation to the conformance test framework.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: gRPC, protobuf, buf v1.32.1
**Storage**: N/A (SDK does not manage persistent storage)
**Testing**: Go testing framework, conformance tests (Basic/Standard/Advanced levels), performance benchmarks
**Target Platform**: Linux/Any (cross-platform gRPC SDK)
**Project Type**: Single project (SDK/library)
**Performance Goals**: GetPluginInfo < 100ms per call, timeout handling < 30s default,
10 iteration conformance test < 2s total
**Constraints**: Backward compatible, no breaking proto changes, existing tests must pass
**Scale/Scope**: SDK enhancement affecting sdk/go/pluginsdk/client.go, sdk.go, and sdk/go/testing/conformance_test.go

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle                           | Status             | Notes                                                                 |
| ----------------------------------- | ------------------ | --------------------------------------------------------------------- |
| I. gRPC Proto Specification-First   | ✅ PASS            | No protobuf changes required - SDK-level enhancement only             |
| II. Multi-Provider gRPC Consistency | ✅ PASS            | No provider-specific changes introduced                               |
| III. Test-First Protocol            | ✅ PASS            | Performance conformance tests will be written before implementation   |
| IV. Protobuf Backward Compatibility | ✅ PASS            | No proto changes - backward compatible SDK enhancement                |
| V. Comprehensive Documentation      | ⚠️ ACTION REQUIRED | README.md and docs/ must be updated in same PR                        |
| VI. Performance as gRPC Requirement | ✅ PASS            | GetPluginInfo performance conformance test included (100ms threshold) |
| VII. Validation at Multiple Levels  | ✅ PASS            | Must pass buf lint, schema validation, conformance tests, unit tests  |

**Gate Status**: ✅ PASS (with documentation action required in Phase 1)

## Project Structure

### Documentation (this feature)

```text
specs/034-sdk-polish/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
sdk/go/
├── pluginsdk/
│   ├── client.go        # Modify: ClientConfig.Timeout implementation
│   └── sdk.go           # Modify: GetPluginInfo error message formatting
├── testing/
│   └── conformance_test.go  # Modify: Add GetPluginInfoPerformance test

tests/
├── integration/
└── unit/

docs/ (to be updated)
└── (new/conformance testing documentation)
README.md (to be updated)
```

**Structure Decision**: Single project SDK enhancement modifying existing files in sdk/go/pluginsdk/
and sdk/go/testing/. No new directories required.
