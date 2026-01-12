# Tasks: Add Supports() RPC Method to CostSourceService

**Input**: Design documents from `/specs/002-add-supports-rpc/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Status**: Feature already implemented in v0.1.0. Tasks below are administrative actions to
close issues and unblock dependent work.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different actions, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)

---

## Phase 1: Verification

**Purpose**: Confirm all components exist and pass validation

- [x] T001 Verify proto contains Supports RPC in proto/finfocus/v1/costsource.proto
- [x] T002 [P] Verify generated client method in sdk/go/proto/finfocus/v1/costsource_grpc.pb.go
- [x] T003 [P] Verify generated server interface in sdk/go/proto/finfocus/v1/costsource_grpc.pb.go
- [x] T004 [P] Verify service descriptor entry in sdk/go/proto/finfocus/v1/costsource_grpc.pb.go
- [x] T005 Run `make test` to confirm all tests pass
- [x] T006 Run `make lint` to confirm code quality
- [x] T006a [P] Verify CHANGELOG documents Supports RPC availability in CHANGELOG.md

**Checkpoint**: All components verified, tests and lint pass

---

## Phase 2: User Story 1 - Query Plugin Capabilities (Priority: P1)

**Goal**: Enable clients to query plugin capabilities via gRPC

**Independent Test**: Supports() RPC callable via gRPC test harness with proper response

### Verification for User Story 1

- [x] T007 [US1] Verify ValidateSupportsResponse() in sdk/go/testing/harness.go
- [x] T008 [US1] Verify MockPlugin.Supports() in sdk/go/testing/mock_plugin.go
- [x] T009 [US1] Verify integration tests include Supports in sdk/go/testing/integration_test.go
- [x] T010 [US1] Run integration test: `go test -v ./sdk/go/testing/ -run TestBasicPluginFunctionality`

**Checkpoint**: User Story 1 functionality verified

---

## Phase 3: User Story 2 - Graceful Capability Discovery (Priority: P2)

**Goal**: Plugin developers can implement Supports() via clear contract

**Independent Test**: Plugin can implement Supports() and respond via gRPC

### Verification for User Story 2

- [x] T011 [US2] Verify CostSourceServiceServer interface includes Supports in costsource_grpc.pb.go
- [x] T012 [US2] Verify UnimplementedCostSourceServiceServer.Supports() default in costsource_grpc.pb.go
- [x] T013 [US2] Verify error handling tests include Supports in sdk/go/testing/error_handling_test.go

**Checkpoint**: User Story 2 contract verified

---

## Phase 4: Issue Management

**Purpose**: Close finfocus-spec#64 and unblock finfocus-core#160

- [x] T014 Close GitHub Issue #64 with completion comment
- [x] T015 Comment on finfocus-core#160 to unblock it

### T014 - Close Issue #64 Comment

Post this comment and close the issue:

```text
## Resolution: Already Complete

The Supports() RPC method has been fully implemented in finfocus-spec since v0.1.0.

### Verified Components

**Proto Definition** (proto/finfocus/v1/costsource.proto):
- `rpc Supports(SupportsRequest) returns (SupportsResponse);` (line 17)
- `SupportsRequest` message (lines 38-42)
- `SupportsResponse` message (lines 44-50)

**Generated Go SDK** (sdk/go/proto/finfocus/v1/):
- `CostSourceServiceClient.Supports()` method
- `CostSourceServiceServer.Supports()` interface
- `UnimplementedCostSourceServiceServer.Supports()` default
- `_CostSourceService_Supports_Handler` gRPC handler
- Service descriptor entry in `CostSourceService_ServiceDesc.Methods`

**Testing Framework** (sdk/go/testing/):
- `ValidateSupportsResponse()` validator
- `MockPlugin.Supports()` implementation
- Integration tests for Supports RPC
- Error handling tests

### Next Steps

The actual issue is in **finfocus-core's pluginsdk**, which needs to expose the Supports
method to plugin implementations. See finfocus-core#160.
```

### T015 - Unblock finfocus-core#160 Comment

Post this comment on <https://github.com/rshade/finfocus-core/issues/160>:

```text
This issue is now unblocked. finfocus-spec#64 is already complete - the Supports()
RPC method has been in finfocus-spec since v0.1.0 release:

- Proto: `rpc Supports(SupportsRequest) returns (SupportsResponse);`
- Messages: SupportsRequest, SupportsResponse
- Generated code: Client/Server interfaces, handlers, service descriptor
- Testing: ValidateSupportsResponse, MockPlugin.Supports, integration tests

Update finfocus-core to use v0.1.0 or later and proceed with pluginsdk implementation.
```

**Checkpoint**: Issues properly documented and dependencies unblocked

---

## Dependencies & Execution Order

### Phase Dependencies

- **Verification (Phase 1)**: No dependencies - start immediately
- **User Story 1 (Phase 2)**: After Phase 1 verification passes
- **User Story 2 (Phase 3)**: After Phase 1 verification passes (parallel with Phase 2)
- **Issue Management (Phase 4)**: After all verification completes

### Parallel Opportunities

- T002, T003, T004 can run in parallel (different files)
- T007, T008, T009 can run in parallel (different files)
- T011, T012, T013 can run in parallel (different aspects)
- T014, T015 can run in parallel (different issues)

---

## Implementation Strategy

### Verification Flow

1. Complete Phase 1: Verify all proto and SDK components exist
2. Complete Phase 2: Verify User Story 1 testing components
3. Complete Phase 3: Verify User Story 2 contract components
4. Complete Phase 4: Close issue and unblock dependent work

### No Code Changes Required

This feature is **already implemented**. All tasks are verification and administrative actions.
No modifications to proto, SDK, or testing code are needed.

---

## Notes

- Feature fully implemented in v0.1.0 release
- Tasks verify existing implementation, not create new code
- Main action is issue management (close #64, unblock core#160)
- Dependent work continues in finfocus-core repository
