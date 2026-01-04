# Tasks: JSON-LD / Schema.org Serialization

**Input**: Design documents from `/specs/032-jsonld-serialization/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are included as this is a library feature requiring conformance tests (per Constitution Check in plan.md)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Go package**: `sdk/go/jsonld/`
- **Tests**: `sdk/go/jsonld/*_test.go`
- **Examples**: `examples/jsonld/`
- **Contracts**: `specs/032-jsonld-serialization/contracts/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create sdk/go/jsonld/ package directory structure
- [x] T002 Create examples/jsonld/ directory for example outputs
- [x] T003 [P] Create doc.go with package documentation in sdk/go/jsonld/doc.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Implement Context type and Build() method in sdk/go/jsonld/context.go
- [x] T005 [P] Create vocabulary.go with FOCUS namespace constants in sdk/go/jsonld/vocabulary.go
- [x] T006 [P] Create schema_org.go with Schema.org type definitions in sdk/go/jsonld/schema_org.go
- [x] T007 Implement IDGenerator interface and SHA256-based generation in sdk/go/jsonld/id_generator.go
- [x] T008 [P] Create unit tests for Context.Build() in sdk/go/jsonld/context_test.go
- [x] T009 [P] Create unit tests for IDGenerator in sdk/go/jsonld/id_generator_test.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Basic FOCUS Record Serialization (Priority: P1) 🎯 MVP

**Goal**: Serialize individual FocusCostRecord protobuf messages to valid JSON-LD 1.1 format

**Independent Test**: Can be fully tested by serializing a single FocusCostRecord and validating the output is
well-formed JSON-LD that passes JSON-LD validation tools. Delivers value by enabling knowledge graph ingestion.

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T010 [P] [US1] Conformance test for required @context declaration in sdk/go/jsonld/serializer_test.go
- [x] T011 [P] [US1] Conformance test for @id generation (user-provided and fallback) in sdk/go/jsonld/serializer_test.go
- [x] T012 [P] [US1] Conformance test for empty value omission in sdk/go/jsonld/serializer_test.go
- [x] T013 [P] [US1] Conformance test for FOCUS 1.3 allocation fields serialization in sdk/go/jsonld/serializer_test.go
- [x] T014 [P] [US1] Conformance test for tags and extended_columns serialization in sdk/go/jsonld/serializer_test.go

### Implementation for User Story 1

- [x] T015 [US1] Create SerializerOptions struct in sdk/go/jsonld/serializer.go
- [x] T016 [US1] Implement NewSerializer() constructor with functional options in sdk/go/jsonld/serializer.go
- [x] T017 [US1] Implement Serialize() method for FocusCostRecord in sdk/go/jsonld/serializer.go (depends on T015, T016)
- [x] T018 [US1] Implement proto field to JSON-LD property mapping logic in sdk/go/jsonld/serializer.go (depends on T017)
- [x] T019 [US1] Implement timestamp to ISO 8601 string conversion in sdk/go/jsonld/serializer.go (depends on T017)
- [x] T020 [US1] Implement monetary amount serialization as Schema.org MonetaryAmount in
  sdk/go/jsonld/serializer.go (depends on T017)
- [x] T021 [US1] Implement enum serialization (string format) in sdk/go/jsonld/serializer.go (depends on T017)
- [x] T022 [US1] Implement map field serialization (tags, extended_columns) in sdk/go/jsonld/serializer.go (depends on T017)
- [x] T023 [US1] Implement deprecated field handling with annotations in sdk/go/jsonld/serializer.go (depends on T017)
- [x] T024 [US1] Add OmitEmptyFields logic (skip empty string, zero numeric, nil) in sdk/go/jsonld/serializer.go
  (depends on T017)
- [x] T024.1 [US1] Implement UTF-8 validation and sanitization for string fields in sdk/go/jsonld/serializer.go
  (depends on T017)
- [x] T025 [US1] Integrate IDGenerator into Serialize() method in sdk/go/jsonld/serializer.go (depends on T017, T007, T024.1)
- [x] T026 [US1] Integrate Context.Build() output into Serialize() method in sdk/go/jsonld/serializer.go
  (depends on T004, T017)

### Example Output for User Story 1

- [x] T027 [P] [US1] Create focus_cost_record.jsonld example in examples/jsonld/focus_cost_record.jsonld
- [x] T028 [P] [US1] Add README.md for examples in examples/jsonld/README.md

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Batch Serialization for Enterprise Scale (Priority: P2)

**Goal**: Serialize large volumes of FOCUS cost records efficiently with bounded memory usage

**Independent Test**: Can be tested by serializing 10,000 FocusCostRecords in a batch and measuring throughput.
Delivers value by enabling enterprise-scale data pipelines.

### Tests for User Story 2 ⚠️

- [x] T029 [P] [US2] Unit test for batch serialization with 10,000 records in sdk/go/jsonld/streaming_test.go
- [x] T030 [P] [US2] Unit test for streaming error handling (partial failures) in sdk/go/jsonld/streaming_test.go
- [x] T031 [P] [US2] Benchmark for single record serialization in sdk/go/jsonld/serializer_benchmark_test.go
- [x] T032 [P] [US2] Benchmark for batch serialization of 10,000 records in sdk/go/jsonld/serializer_benchmark_test.go
- [x] T033 [P] [US2] Memory usage benchmark for streaming in sdk/go/jsonld/serializer_benchmark_test.go

### Implementation for User Story 2

- [x] T034 [US2] Implement SerializeStream() method with io.Writer in sdk/go/jsonld/streaming.go
- [x] T035 [US2] Implement JSON array streaming pattern ([...] with separator logic) in sdk/go/jsonld/streaming.go
  (depends on T034)
- [x] T036 [US2] Implement error collection and reporting for batch operations in sdk/go/jsonld/streaming.go
  (depends on T034)
- [x] T037 [US2] Implement sync.Pool for byte buffer reuse in sdk/go/jsonld/streaming.go (depends on T034)
- [x] T038 [US2] Add streaming support for Serialize() method variants in sdk/go/jsonld/serializer.go (depends on T034)

### Example Output for User Story 2

- [x] T039 [P] [US2] Create batch_output.jsonld example in examples/jsonld/batch_output.jsonld

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Schema.org Vocabulary Mapping (Priority: P2)

**Goal**: Apply Schema.org vocabulary types to FOCUS fields where natural mappings exist (validated across all 66
mappable fields)

**Independent Test**: Can be tested by serializing a FocusCostRecord and validating that Schema.org types and
properties are correctly applied where mappings exist. Delivers value by enabling semantic web tool integration.

### Tests for User Story 3 ⚠️

- [x] T040 [P] [US3] Unit test for Schema.org MonetaryAmount type application in sdk/go/jsonld/schema_org_test.go
- [x] T041 [P] [US3] Unit test for Schema.org DateTime formatting (ISO 8601) in sdk/go/jsonld/schema_org_test.go
- [x] T042 [P] [US3] Unit test for custom FOCUS namespace fallback in sdk/go/jsonld/schema_org_test.go
- [x] T043 [P] [US3] Conformance test: validate output against JSON-LD 1.1 spec in sdk/go/jsonld/serializer_test.go
- [x] T043.1 [P] [US3] Conformance test: validate Schema.org mappings for all 66 mappable fields in sdk/go/jsonld/serializer_test.go

### Implementation for User Story 3

- [x] T044 [US3] Implement MonetaryAmount type embedding for cost fields (billed_cost, list_cost, effective_cost,
  contracted_cost) in sdk/go/jsonld/serializer.go (depends on T020)
- [x] T045 [US3] Implement DateTime type coercion for timestamp fields in sdk/go/jsonld/serializer.go (depends on T019)
- [x] T046 [US3] Implement PropertyValue type for tags/map serialization in sdk/go/jsonld/serializer.go (depends on T022)
- [x] T047 [US3] Add Schema.org @type annotations to mapped fields in sdk/go/jsonld/serializer.go
  (depends on T044, T045, T046)
- [x] T048 [US3] Implement custom FOCUS vocabulary namespace for unmapped fields in sdk/go/jsonld/serializer.go
  (depends on T005)

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently

---

## Phase 6: User Story 4 - Custom Context Configuration (Priority: P3)

**Goal**: Allow customization of JSON-LD context for enterprise ontology integration

**Independent Test**: Can be tested by providing a custom context configuration and validating the output uses the
specified property mappings. Delivers value by enabling integration with existing enterprise ontologies.

### Tests for User Story 4 ⚠️

- [x] T049 [P] [US4] Unit test for WithCustomMapping() in sdk/go/jsonld/context_test.go
- [x] T050 [P] [US4] Unit test for WithRemoteContext() in sdk/go/jsonld/context_test.go
- [x] T051 [P] [US4] Unit test for custom context serialization in sdk/go/jsonld/serializer_test.go
- [x] T052 [P] [US4] Unit test for invalid context configuration error handling in sdk/go/jsonld/context_test.go

### Implementation for User Story 4

- [x] T053 [US4] Implement WithCustomMapping() method on Context in sdk/go/jsonld/context.go
- [x] T054 [US4] Implement WithRemoteContext() method on Context in sdk/go/jsonld/context.go (depends on T053)
- [x] T055 [US4] Add custom mapping override logic in Context.Build() in sdk/go/jsonld/context.go (depends on T053)
- [x] T056 [US4] Add remote context URL handling in Context.Build() in sdk/go/jsonld/context.go (depends on T054)
- [x] T057 [US4] Implement context validation (valid URLs, IRI formats) in sdk/go/jsonld/context.go (depends on T053, T054)
- [x] T058 [US4] Add WithContext() functional option to NewSerializer() in sdk/go/jsonld/serializer.go
  (depends on T016, T053)

**Checkpoint**: All user stories should now be independently functional

---

## Phase 7: User Story 5 - Contract Commitment Dataset Serialization (Priority: P3)

**Goal**: Serialize ContractCommitment protobuf messages to JSON-LD with linked data references to cost records

**Independent Test**: Can be tested by serializing a ContractCommitment and validating it has proper linked data
references. Delivers value by enabling commitment analysis in knowledge graphs.

### Tests for User Story 5 ⚠️

- [x] T059 [P] [US5] Unit test for ContractCommitment serialization in sdk/go/jsonld/serializer_test.go
- [x] T060 [P] [US5] Unit test for commitment ID generation in sdk/go/jsonld/id_generator_test.go
- [x] T061 [P] [US5] Conformance test for commitment-cost record linking in sdk/go/jsonld/serializer_test.go

### Implementation for User Story 5

- [x] T062 [US5] Implement SerializeCommitment() method in sdk/go/jsonld/serializer.go
- [x] T063 [US5] Implement ContractCommitment field to JSON-LD mapping in sdk/go/jsonld/serializer.go (depends on T062)
- [x] T064 [US5] Implement GenerateCommitment() method in IDGenerator in sdk/go/jsonld/id_generator.go
- [x] T065 [US5] Add @id reference support for contract_applied field in sdk/go/jsonld/serializer.go (depends on T062, T017)

### Example Output for User Story 5

- [x] T066 [P] [US5] Create contract_commitment.jsonld example in examples/jsonld/contract_commitment.jsonld

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T067 [P] Add WithPrettyPrint() functional option in sdk/go/jsonld/serializer.go
- [x] T068 [P] Add WithOmitEmpty() functional option in sdk/go/jsonld/serializer.go
- [x] T069 [P] Add WithIRIEnums() functional option in sdk/go/jsonld/serializer.go
- [x] T070 [P] Add WithDeprecated() functional option in sdk/go/jsonld/serializer.go
- [x] T071 [P] Create comprehensive README.md in sdk/go/jsonld/README.md with usage examples
- [x] T072 [P] Add godoc comments to all exported types and methods in sdk/go/jsonld/
- [x] T073 Run all tests: go test ./sdk/go/jsonld/...
- [x] T074 Run benchmarks: go test -bench=. ./sdk/go/jsonld/...
- [x] T075 Validate quickstart.md examples compile and run
- [x] T076 Run golangci-lint on sdk/go/jsonld/ package
- [x] T077 Update examples/jsonld/README.md with all example descriptions

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Builds on US1 but independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Builds on US1 but independently testable
- **User Story 4 (P3)**: Can start after Foundational (Phase 2) - Integrates with US1/US2/US3 but independently testable
- **User Story 5 (P3)**: Can start after Foundational (Phase 2) - Builds on US1 serialization patterns

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Core infrastructure (Context, IDGenerator) before serialization logic
- Serialization logic before example outputs
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members
- Example file creation can run in parallel with implementation tasks

---

## Parallel Example: User Story 1

```bash
# Launch all conformance tests for User Story 1 together:
Task: "Conformance test for required @context declaration in sdk/go/jsonld/serializer_test.go"
Task: "Conformance test for @id generation in sdk/go/jsonld/serializer_test.go"
Task: "Conformance test for empty value omission in sdk/go/jsonld/serializer_test.go"
Task: "Conformance test for FOCUS 1.3 allocation fields serialization in sdk/go/jsonld/serializer_test.go"
Task: "Conformance test for tags and extended_columns serialization in sdk/go/jsonld/serializer_test.go"

# Launch example creation in parallel with implementation:
Task: "Create focus_cost_record.jsonld example in examples/jsonld/focus_cost_record.jsonld"
Task: "Add README.md for examples in examples/jsonld/README.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently with conformance tests
5. Run quickstart.md validation
6. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Add User Story 5 → Test independently → Deploy/Demo
7. Complete Polish phase → Final release
8. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (P1) - Critical path
   - Developer B: User Story 2 (P2) - Batch processing
   - Developer C: User Story 3 (P2) - Schema.org integration
3. Stories complete and integrate independently
4. Team works together on Polish phase

---

## Task Summary

- **Total Tasks**: 77
- **Setup Phase**: 3 tasks
- **Foundational Phase**: 6 tasks
- **User Story 1 (P1)**: 19 tasks (5 tests + 14 implementation)
- **User Story 2 (P2)**: 11 tasks (5 tests + 5 implementation + 1 example)
- **User Story 3 (P2)**: 9 tasks (4 tests + 5 implementation)
- **User Story 4 (P3)**: 10 tasks (4 tests + 6 implementation)
- **User Story 5 (P3)**: 8 tasks (3 tests + 4 implementation + 1 example)
- **Polish Phase**: 11 tasks

### Parallel Opportunities Identified

- **Phase 1**: 2 parallel tasks (T002, T003)
- **Phase 2**: 3 parallel tasks (T005, T006, T008, T009)
- **Phase 3**: 5 parallel test tasks, 2 parallel example tasks
- **Phase 4**: 4 parallel test tasks, 1 parallel example task
- **Phase 5**: 4 parallel test tasks
- **Phase 6**: 4 parallel test tasks
- **Phase 7**: 3 parallel test tasks, 1 parallel example task
- **Phase 8**: 6 parallel tasks (T067-T072)

### Suggested MVP Scope

**MVP = User Story 1 Only** (26 total tasks: 3 setup + 6 foundational + 19 US1 tasks)

This delivers:

- Single FocusCostRecord serialization to JSON-LD
- Valid @context with FOCUS vocabulary
- Proper @id generation (user-provided or SHA256 fallback)
- Empty value omission
- Tag and extended_columns serialization
- FOCUS 1.3 allocation field support
- Conformance tests ensuring JSON-LD 1.1 compliance
- Example output files

MVP enables organizations to start ingesting FOCUS cost data into knowledge graphs immediately.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (test-first approach)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Run `make validate` after completion to ensure linting and tests pass
- Focus on stdlib-only implementation (no external dependencies per research.md)
- Performance targets: <1ms single record, <5s for 10k records, bounded memory
- UTF-8 validation: String fields are validated for valid UTF-8; invalid sequences return
  \*jsonld.ValidationError with field name and sanitization applied
