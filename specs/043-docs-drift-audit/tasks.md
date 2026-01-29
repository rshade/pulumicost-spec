# Tasks: Documentation Drift Audit Remediation

**Input**: Design documents from `/specs/043-docs-drift-audit/`
**Prerequisites**: plan.md, spec.md, research.md
**Source Issues**: GitHub #347, #348

**Tests**: Not applicable - documentation-only changes

**Organization**: Tasks grouped by user story priority (P1 → P2 → P3 → P4)

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Documentation-only feature - all paths relative to repository root:

- `README.md` - Root documentation
- `sdk/go/testing/README.md` - Testing package docs
- `sdk/go/pluginsdk/README.md` - Plugin SDK docs
- `sdk/go/pluginsdk/mapping/README.md` - Mapping package docs (NEW)
- `sdk/go/registry/README.md` - Registry package docs (NEW)
- `sdk/go/pluginsdk/*.go` - Go source files for godoc
- `sdk/go/pricing/*.go` - Go source files for godoc

---

## Phase 1: Setup

**Purpose**: Verify current state and prepare for documentation updates

- [x] T001 Verify current version references in README.md (expected: all v0.5.4)
- [x] T002 [P] Count RPC methods in proto/finfocus/v1/costsource.proto (expected: 11 CostSourceService + 3 Observability)
- [x] T003 [P] Count JSON example files in examples/specs/ (expected: 9 files)
- [x] T004 [P] Verify conformance function signatures in sdk/go/testing/conformance.go

**Checkpoint**: State verified - ready for documentation updates

---

## Phase 2: Foundational

**Purpose**: No blocking prerequisites for documentation changes

**Note**: This feature has no foundational phase - all user stories can begin immediately after setup verification.

---

## Phase 3: User Story 1 - Accurate Version Information (Priority: P1) ✅ RESOLVED

**Goal**: All version references in README.md display consistent v0.5.4

**Independent Test**: Search README.md for "v0." and verify all matches show v0.5.4

**Status**: Already resolved per research.md - all 4 version references (lines 1, 17, 832, 858) show v0.5.4

- [x] T005 [US1] Confirm version consistency in README.md lines 1, 17, 832, 858 (verification only)

**Checkpoint**: US1 complete - version consistency verified

---

## Phase 4: User Story 2 - Accurate RPC Documentation (Priority: P1)

**Goal**: RPC count in README matches actual proto definition (11 methods)

**Independent Test**: Compare README RPC count claims to proto file - counts must match

### Implementation for User Story 2

- [x] T006 [US2] Update RPC count from "8" to "11" in README.md (search: "CostSourceService with 8 RPC")
- [x] T007 [P] [US2] Update RPC count from "8" to "11" at README.md line 699
- [x] T008 [US2] Verify the 11 RPCs are listed in README.md

**Checkpoint**: US2 complete - RPC documentation accurate

---

## Phase 5: User Story 3 - Correct SDK Code Examples (Priority: P1)

**Goal**: Code examples in testing/README.md compile without errors

**Independent Test**: Copy code examples to a Go file and verify they compile

### Implementation for User Story 3

- [x] T009 [US3] Fix conformance example at sdk/go/testing/README.md line 104 to handle error return
- [x] T010 [P] [US3] Fix conformance example at sdk/go/testing/README.md line 171: add error handling for RunBasicConformance
- [x] T011 [P] [US3] Fix conformance example at sdk/go/testing/README.md line 172: add error handling for RunStandardConformance
- [x] T012 [P] [US3] Fix conformance example at sdk/go/testing/README.md line 173: add error handling for RunAdvancedConformance
- [x] T013 [US3] Verify all conformance functions documented in sdk/go/testing/README.md
  with correct signatures (RunBasic/Standard/AdvancedConformance, NewTestHarness, etc.)

**Checkpoint**: US3 complete - all code examples compile

---

## Phase 6: User Story 4 - Complete Package Documentation (Priority: P2)

**Goal**: mapping/ package has README.md explaining purpose and usage

**Independent Test**: Verify sdk/go/pluginsdk/mapping/README.md exists and documents all files

### Implementation for User Story 4

- [x] T014 [US4] Create sdk/go/pluginsdk/mapping/README.md with package overview
- [x] T015 [P] [US4] Document aws.go functions in mapping/README.md (AWS property extraction)
- [x] T016 [P] [US4] Document azure.go functions in mapping/README.md (Azure property extraction)
- [x] T017 [P] [US4] Document gcp.go functions in mapping/README.md (GCP property extraction)
- [x] T018 [US4] Document common.go and keys.go in mapping/README.md (shared utilities and constants)
- [x] T019 [US4] Add usage examples showing SKU and region extraction patterns

**Checkpoint**: US4 complete - mapping package fully documented

---

## Phase 7: User Story 5 - Accurate Example Counts (Priority: P2)

**Goal**: Example count in README matches actual files (9 examples)

**Independent Test**: Compare README example count to `ls examples/specs/*.json | wc -l`

### Implementation for User Story 5

- [x] T020 [US5] Update example count from "10" to "9" in README.md (search: "10 comprehensive pricing")
- [x] T021 [P] [US5] Update example count from "8" to "9" at README.md line 660
- [x] T022 [US5] Verify example list matches actual files in examples/specs/

**Checkpoint**: US5 complete - example count accurate

---

## Phase 8: User Story 6 - Documented SDK Helpers (Priority: P3)

**Goal**: Key SDK helpers documented with usage examples

**Independent Test**: Verify NewActualCostResponse() and FallbackHint documented in user-facing README

### Implementation for User Story 6

- [x] T023 [US6] Add FallbackHint enum documentation to sdk/go/pluginsdk/README.md
- [x] T024 [P] [US6] Add NewActualCostResponse() example with functional options to sdk/go/pluginsdk/README.md
- [x] T025 [P] [US6] Add validation helper documentation (ValidateActualCostResponse, ValidateRecommendation) to sdk/go/pluginsdk/README.md
- [x] T026 [US6] Update README.md lines 259-279 to use pluginsdk.NewActualCostResponse() instead of manual struct construction
- [x] T027 [US6] Add cross-reference from root README.md to sdk/go/pluginsdk/README.md for SDK helpers

**Checkpoint**: US6 complete - SDK helpers documented

---

## Phase 9: Godoc Coverage (Priority: P4 - from Issue #348)

**Goal**: All exported functions have godoc comments

**Independent Test**: Run `go doc` on packages and verify no undocumented exports

### Implementation for Godoc Coverage

- [x] T028 [P] Add godoc comments to response builders in sdk/go/pluginsdk/helpers.go
- [x] T029 [P] Add godoc comments to validation helpers in sdk/go/pluginsdk/helpers.go
- [x] T030 [P] Add godoc comments to field mapping utilities in sdk/go/pluginsdk/dry_run.go
- [x] T031 [P] Add godoc comments to growth functions in sdk/go/pricing/growth.go
- [x] T032 [P] Add godoc comments to validation functions in sdk/go/pricing/validate.go
- [x] T033 Create sdk/go/registry/README.md documenting 8 enum types with validators

**Checkpoint**: Godoc coverage complete

---

## Phase 10: Polish & Verification

**Purpose**: Final validation and cleanup

- [x] T034 [P] Run `make lint-markdown` to verify all markdown formatting
- [x] T035 [P] Verify all version references still show v0.5.4
- [x] T036 [P] Verify RPC count matches proto (11 in CostSourceService)
- [x] T037 [P] Verify example count matches files (9 in examples/specs/)
- [x] T038 Manually verify code examples in sdk/go/testing/README.md compile
- [x] T039 Run `go doc ./sdk/go/pluginsdk` and verify key functions documented
- [x] T040 Run `go doc ./sdk/go/pricing` and verify key functions documented
- [x] T041 [P] Verify godoc coverage ≥80% for sdk/go/pluginsdk package (Constitution XIV)
- [x] T042 [P] Verify godoc coverage ≥80% for sdk/go/pricing package (Constitution XIV)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verification only
- **Foundational (Phase 2)**: N/A for documentation feature
- **User Stories (Phases 3-9)**: All independent - can run in parallel
- **Polish (Phase 10)**: Depends on all user stories complete

### User Story Dependencies

| Story | Priority | Dependencies | Can Parallelize With |
| ----- | -------- | ------------ | -------------------- |
| US1   | P1       | None         | US2, US3, US4, US5, US6 |
| US2   | P1       | None         | US1, US3, US4, US5, US6 |
| US3   | P1       | None         | US1, US2, US4, US5, US6 |
| US4   | P2       | None         | US1, US2, US3, US5, US6 |
| US5   | P2       | None         | US1, US2, US3, US4, US6 |
| US6   | P3       | None         | US1, US2, US3, US4, US5 |
| P4    | P4       | None         | All user stories |

**Note**: All user stories are independent - documentation changes to different files have no dependencies.

### Parallel Opportunities

- **Phase 1**: T002, T003, T004 can run in parallel
- **Phase 4 (US2)**: T006, T007 can run in parallel
- **Phase 5 (US3)**: T010, T011, T012 can run in parallel
- **Phase 6 (US4)**: T015, T016, T017 can run in parallel
- **Phase 7 (US5)**: T020, T021 can run in parallel
- **Phase 8 (US6)**: T024, T025 can run in parallel
- **Phase 9 (P4)**: T028, T029, T030, T031, T032 can run in parallel
- **Phase 10**: T034, T035, T036, T037 can run in parallel

---

## Parallel Example: User Stories 2 and 3

```bash
# These can run simultaneously since they modify different files:

# User Story 2 (README.md RPC count):
Task: "Update RPC count from '8' to '11' at README.md line 80"
Task: "Update RPC count from '8' to '11' at README.md line 699"

# User Story 3 (testing/README.md examples):
Task: "Fix conformance example at sdk/go/testing/README.md line 104"
Task: "Fix conformance example at sdk/go/testing/README.md line 171"
```

---

## Implementation Strategy

### MVP First (User Stories 1-3 Only)

1. Complete Phase 1: Setup verification
2. Complete Phase 4: US2 - RPC count (2 line changes)
3. Complete Phase 5: US3 - Code examples (4 code block fixes)
4. **STOP and VALIDATE**: Run `make lint-markdown`, verify examples compile
5. This delivers critical P1 fixes

### Incremental Delivery

1. **MVP**: US1 + US2 + US3 → Critical accuracy fixes
2. **+US4**: mapping/ README → Package completeness
3. **+US5**: Example count → Minor accuracy fix
4. **+US6**: SDK helper docs → Improved discoverability
5. **+P4**: Godoc coverage → Comprehensive documentation

### Estimated Effort

| Phase | Tasks | Effort | Files Modified |
| ----- | ----- | ------ | -------------- |
| Setup | 4 | Low | 0 (verification) |
| US1 | 1 | Low | 0 (already resolved) |
| US2 | 3 | Low | 1 (README.md) |
| US3 | 5 | Medium | 1 (testing/README.md) |
| US4 | 6 | Medium | 1 (NEW mapping/README.md) |
| US5 | 3 | Low | 1 (README.md) |
| US6 | 5 | Medium | 2 (README.md, pluginsdk/README.md) |
| P4 | 6 | High | 4 (*.go files, registry/README.md) |
| Polish | 7 | Low | 0 (verification) |

**Total**: 42 tasks across 10 phases

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- US1 is already complete (version consistency resolved)
- All documentation changes must pass `make lint-markdown`
- Godoc tasks (P4) can be deferred if MVP scope is preferred
- Commit after each user story for clean git history
