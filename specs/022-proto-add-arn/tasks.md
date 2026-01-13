# Tasks: Add ARN to GetActualCostRequest

**Branch**: `018-proto-add-arn`
**Spec**: [spec.md](spec.md)
**Plan**: [plan.md](plan.md)

## Implementation Strategy

- **Phase 1**: Setup & Validation
- **Phase 2**: Protocol Definition (Proto)
- **Phase 3**: SDK Generation & Verification (User Story 1)
- **Phase 4**: Polish

## Dependencies

- US1 (Receive ARN) depends on Proto update.

## Phase 1: Setup

**Goal**: Ensure environment is ready for proto modification.

- [X] T001 Verify existence and content of proto/finfocus/v1/costsource.proto
- [X] T002 Verify buf CLI is installed and configured (run `buf --version`)

## Phase 2: Protocol Definition (Foundational)

**Goal**: Update the source of truth (protobuf) with the new field.

- [X] T003 [US1] Add `string arn = 5;` field to `GetActualCostRequest` message in proto/finfocus/v1/costsource.proto
- [X] T004 [US1] Add comments to the new `arn` field in proto/finfocus/v1/costsource.proto
  describing its purpose (Canonical Cloud Identifier)
- [X] T005 [US1] Run `buf lint` to ensure proto style compliance
- [X] T006 [US1] Run `buf breaking --against .git#branch=main` to ensure backward compatibility

## Phase 3: SDK Generation & Verification

**Goal**: Generate Go code and verify the new field is usable.

- [X] T007 [US1] Regenerate Go SDK using `make generate` (updates sdk/go/proto/costsource.pb.go)
- [X] T008 [US1] Create test file `sdk/go/testing/arn_test.go` to verify `GetActualCostRequest`
  struct has the `Arn` field
- [X] T009 [US1] Implement a round-trip test in `sdk/go/testing/arn_test.go` that sets and retrieves the `Arn` field
- [X] T010 [US1] Add a test case to `sdk/go/testing/arn_test.go` that constructs a
  `GetActualCostRequest` without the `arn` field, ensuring existing consumers can process it
  gracefully (backward compatibility verification).
- [X] T011 [US1] Run all SDK tests with `go test ./sdk/go/...` to ensure no regressions

## Phase 4: Polish

**Goal**: Final cleanups.

- [X] T012 Update `CHANGELOG.md` (if present) or `specs/018-proto-add-arn/spec.md` status
- [X] T013 Verify generated documentation matches proto comments
- [X] T014 Integrate content from `specs/018-proto-add-arn/quickstart.md` into relevant project
  documentation (e.g., `PLUGIN_DEVELOPER_GUIDE.md` or `README.md`)
