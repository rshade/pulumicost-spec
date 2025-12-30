# Tasks: Multi-Protocol Plugin Access

**Feature**: `029-grpc-web-support`
**Status**: Complete
**Priority**: High
**Updated**: 2025-12-29 (Implemented with connect-go)

## Phase 1: Setup

*Goal: Configure buf for connect-go and add dependencies.*

- [x] T001 Add `buf.build/connectrpc/go` plugin to `buf.gen.yaml`
- [x] T002 Add `connectrpc.com/connect` dependency to `go.mod`
- [x] T003 Add `connectrpc.com/grpchealth` dependency to `go.mod`
- [x] T004 Add `golang.org/x/net/http2/h2c` dependency to `go.mod`
- [x] T005 Run `make generate` to create connect handlers
- [x] T006 Run `go build ./...` to verify generation succeeded

## Phase 2: Foundational (Server Core)

*Goal: Implement connect-go server with multi-protocol support.*

- [x] T007 Create `WebConfig` struct in `sdk/go/pluginsdk/sdk.go`
- [x] T008 Add `Web WebConfig` field to `ServeConfig` struct in `sdk/go/pluginsdk/sdk.go`
- [x] T009 [P] Implement `/healthz` HTTP handler in `sdk/go/pluginsdk/sdk.go`
- [x] T010 Create connect handler adapter in `sdk/go/pluginsdk/connect.go`
- [x] T011 Update `Serve()` to use connect handler with h2c wrapper
- [x] T012 Implement CORS support using custom middleware
- [x] T013 [P] Write unit tests for `WebConfig` in `sdk/go/pluginsdk/connect_test.go`
- [x] T014 Write integration test for connect server in `sdk/go/pluginsdk/connect_test.go`

## Phase 3: Protocol Verification

*Goal: Verify all three protocols work correctly.*

- [x] T015 [US1] Test gRPC protocol with existing grpc-go client
- [x] T016 [US1] Test Connect protocol with curl (JSON) - verified in tests
- [x] T017 [US1] Test gRPC-Web protocol with browser fetch - verified in tests
- [x] T018 [US1] Verify CORS headers for cross-origin requests - TestServeConnect_CORS
- [x] T019 [US5] Verify `/healthz` endpoint returns 200 OK - TestServeConnect_WithHealthEndpoint
- [ ] T020 [US5] Verify `grpc.health.v1.Health/Check` via all protocols - Deferred (not required for MVP)

## Phase 4: Go Client SDK

*Goal: Provide a robust Go client SDK.*

- [x] T021 [US6] Create client in `sdk/go/pluginsdk/client.go`
- [x] T022 [US6] Implement `NewClient` constructor with protocol support
- [x] T023 [US6] Implement typed methods for standard RPCs
- [ ] T024 [US2] Implement `BatchEstimateCost` with errgroup concurrency - Deferred for future iteration
- [ ] T025 [US2] Write integration test for batch operations - Deferred for future iteration
- [x] T026 [US6] Write unit tests for client in `sdk/go/pluginsdk/client_test.go`

## Phase 5: Multi-Tenant Isolation

*Goal: Ensure plugin SDK respects environment isolation.*

- [x] T027 [US3] Verify SDK does not cache credentials globally (code audit) - No global caching
- [ ] T028 [US3] Add multi-tenant deployment documentation - Deferred for future iteration

## Phase 6: Examples and Documentation

*Goal: Provide comprehensive examples and documentation.*

- [x] T029 Create browser example in `examples/plugins/web-client/index.html`
- [x] T030 Create curl examples in `examples/plugins/web-client/README.md`
- [x] T031 Update `sdk/go/pluginsdk/README.md` with connect-go usage
- [x] T032 Update quickstart.md with new API examples

## Phase 7: Polish and Validation

*Goal: Finalize code quality and run full validation.*

- [x] T033 Run `golangci-lint` and fix any issues - 0 issues
- [x] T034 Run `make test` ensuring no regressions - All tests pass
- [x] T035 Run markdown linting on all documentation - 0 errors
- [x] T036 Verify backward compatibility with existing tests - All existing tests pass

## Dependencies

1. Phase 1 (Setup) blocks all other phases
2. Phase 2 (Foundational) blocks Phases 3-6
3. Phase 3 (Protocol Verification) and Phase 4 (Client) can run in parallel
4. Phase 5 (Isolation) is independent after Phase 2
5. Phase 6 (Docs) can start after Phase 3
6. Phase 7 (Polish) runs last

## Implementation Strategy

1. **MVP**: Complete Phases 1-3 to prove multi-protocol connectivity
2. **Client SDK**: Complete Phase 4 for Go backend service integration
3. **Production Readiness**: Complete Phases 5-7

## Connect-go vs grpc-web Decision

**Rationale for connect-go**:

- Multi-protocol support (gRPC + gRPC-Web + Connect) on same handler
- JSON support enables curl/fetch debugging without special tooling
- Actively maintained by Buf team (same as our buf CLI)
- Existing gRPC clients work unchanged (wire-compatible)
- Built-in CORS support, no external middleware needed

See [research.md](./research.md) for full analysis.
