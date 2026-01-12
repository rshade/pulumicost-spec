# Implementation Plan: Resource ID and ARN Fields for ResourceDescriptor

**Branch**: `028-resource-id` | **Date**: 2025-12-26 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/028-resource-id/spec.md`

## Summary

Add two new fields to the `ResourceDescriptor` protobuf message:

- `id` (field 7): Client-specified opaque identifier for request/response correlation
  in batch operations
- `arn` (field 8): Canonical cloud resource identifier (AWS ARN, Azure Resource ID,
  GCP Full Resource Name, etc.) for exact resource matching

This is a backward-compatible protocol addition (minor version bump) enabling
finfocus-core to correlate batch recommendation responses and plugins to perform
precise resource lookups.

## Technical Context

**Language/Version**: Go 1.25.5 (per go.mod) + Protocol Buffers v3
**Primary Dependencies**: google.golang.org/protobuf, google.golang.org/grpc, buf v1.32.1
**Storage**: N/A (stateless proto definitions)
**Testing**: go test, conformance tests, buf lint/breaking
**Target Platform**: Cross-platform gRPC plugins
**Project Type**: Single (gRPC specification + Go SDK)
**Performance Goals**: Zero-allocation field access, O(1) correlation lookup
**Constraints**: Backward-compatible proto changes only, no breaking wire format
**Scale/Scope**: Used by all FinFocus plugins and core

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

| Principle | Status | Notes |
|-----------|--------|-------|
| I. gRPC Proto Specification-First | ✅ PASS | Proto changes first, SDK generation follows |
| II. Multi-Provider gRPC Consistency | ✅ PASS | ARN field supports AWS, Azure, GCP, K8s, Cloudflare |
| III. Test-First Protocol | ✅ PASS | Conformance tests will be written before proto change |
| IV. Protobuf Backward Compatibility | ✅ PASS | New optional fields only, minor version bump |
| V. Comprehensive Documentation | ✅ PASS | Proto comments + SDK godoc + examples planned |
| VI. Performance as gRPC Requirement | ✅ PASS | Zero-allocation field access design |
| VII. Validation at Multiple Levels | ✅ PASS | buf lint, conformance tests, integration tests |

**All gates pass. Proceeding to Phase 0.**

## Project Structure

### Documentation (this feature)

```text
specs/028-resource-id/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (proto diff)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
proto/finfocus/v1/
└── costsource.proto     # Add id (field 7) and arn (field 8) to ResourceDescriptor

sdk/go/
├── proto/               # Regenerated gRPC code (make generate)
│   └── finfocus/v1/
│       └── costsource.pb.go
├── pluginsdk/
│   └── helpers.go       # Add WithID/WithARN builder methods (if needed)
└── testing/
    └── resource_id_test.go  # New conformance tests for id/arn fields
```

**Structure Decision**: Single project structure - this is a protobuf specification
repository with generated Go SDK. Changes are localized to proto files and SDK
helpers.

## Complexity Tracking

> No violations. All changes are backward-compatible field additions.
