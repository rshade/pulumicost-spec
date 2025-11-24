# Tasks: Plugin Registry Index JSON Schema

**Input**: Design documents from `/specs/003-plugin-registry-schema/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Test tasks included (schema validation is test-driven by nature).

**Organization**: Tasks grouped by user story to enable independent implementation and
testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Schema files**: `schemas/` at repository root
- **Examples**: `examples/` at repository root
- **Scripts**: `scripts/` at repository root (npm validation)
- **Documentation**: Repository root (README.md, package.json)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and verification of existing infrastructure

- [x] T001 Verify existing schema validation infrastructure in package.json
- [x] T002 [P] Review existing schemas for pattern consistency in schemas/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core schema file that all user stories depend on

**⚠️ CRITICAL**: No user story work can begin until the schema is created

- [x] T003 Create plugin_registry.schema.json in schemas/plugin_registry.schema.json
  - Include `$schema`, `$id`, `title`, `description`
  - Define root properties: `schema_version`, `plugins`
  - Define `$defs/RegistryEntry` with all required fields
  - Add all enums: `supported_providers`, `capabilities`, `security_level`
  - Add `dependentRequired` for deprecated/deprecation_message
  - Set `additionalProperties: false` at all levels

**Checkpoint**: Schema created - validation tests can now be written

---

## Phase 3: User Story 1 - Validate Registry Index File (Priority: P1) 🎯 MVP

**Goal**: Enable registry contributors to validate registry.json files against the schema

**Independent Test**: Run `npm run validate:registry` against example file and verify
valid entries pass, invalid entries fail with clear errors

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before full implementation**

- [x] T004 [US1] Create valid example registry in examples/registry.json
  - Include 2 plugins: kubecost (kubernetes) and aws-public (aws)
  - All required fields populated correctly
  - Optional fields demonstrated (license, homepage, capabilities, keywords)

- [x] T005 [US1] Add validate:registry npm script in package.json
  - Use AJV with `--strict=false` flag
  - Validate examples/registry.json against schemas/plugin_registry.schema.json

- [x] T006 [US1] Update validate npm script to include registry validation

### Implementation for User Story 1

- [x] T007 [US1] Verify schema validates example correctly by running npm run validate:registry
- [x] T008 [US1] Test invalid cases manually:
  - Missing required field (remove min_spec_version)
  - Invalid name pattern (uppercase letters)
  - Invalid provider enum (invalid value)
  - Deprecated without deprecation_message

**Checkpoint**: User Story 1 complete - registry validation works

---

## Phase 4: User Story 2 - Discover Plugin Metadata (Priority: P1)

**Goal**: Schema captures all metadata fields required by plugin installation system

**Independent Test**: Verify schema includes all fields from registry.proto PluginInfo,
PluginMetadata, and PluginRequirements messages

### Tests for User Story 2

- [x] T009 [US2] Verify proto alignment by comparing schema fields to registry.proto
  - PluginInfo: name, description, author, capabilities, security_level
  - PluginRequirements: min_spec_version, max_spec_version
  - PluginMetadata: homepage, repository, license, keywords

### Implementation for User Story 2

- [x] T010 [US2] Document proto alignment in schema field descriptions
- [x] T011 [US2] Add examples array to schema fields that have specific formats
  - repository: `"rshade/pulumicost-plugin-kubecost"`
  - schema_version: `"1.0.0"`
  - min_spec_version: `"0.1.0"`

**Checkpoint**: User Story 2 complete - schema aligned with proto

---

## Phase 5: User Story 3 - Contribute Plugin to Registry (Priority: P2)

**Goal**: Clear guidance for third-party plugin developers on registry contribution

**Independent Test**: New developer can create valid registry entry using only schema
documentation and quickstart guide

### Implementation for User Story 3

- [x] T012 [US3] Verify schema properties have clear descriptions in schemas/plugin_registry.schema.json
  - All 15+ properties should have description field
  - Descriptions explain purpose, not just type
  - Pattern fields include format guidance
- [x] T013 [US3] Update examples/README.md with registry.json documentation
  - Explain registry format
  - Link to quickstart.md for contribution guide
  - Document validation commands

**Checkpoint**: User Story 3 complete - contribution guidance available

---

## Phase 6: User Story 4 - Filter Plugins by Capability (Priority: P3)

**Goal**: Schema supports discoverability through capabilities, providers, and keywords

**Independent Test**: Schema validates arrays correctly, enforces uniqueItems

### Implementation for User Story 4

- [x] T014 [US4] Verify schema enum arrays support filtering use case
  - supported_providers: minItems 1, uniqueItems
  - capabilities: uniqueItems
  - keywords: maxItems 10, uniqueItems

- [x] T015 [US4] Verify example plugins demonstrate filtering diversity in examples/registry.json
  - kubecost: kubernetes provider, cost_retrieval/cost_projection/real_time_data
  - aws-public: aws provider, cost_projection/pricing_specs
  - Different capabilities enable filtering use case demonstration

**Checkpoint**: User Story 4 complete - filterability supported

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final validation

- [x] T016 [P] Run markdownlint on all new/modified markdown files
- [x] T017 [P] Run full validation suite: `make validate`
- [x] T018 Update CHANGELOG.md with new schema addition
- [x] T019 Run quickstart.md validation examples to verify documentation accuracy

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verify infrastructure
- **Foundational (Phase 2)**: Depends on Setup - creates core schema
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - US1 and US2 are both P1 but US1 must complete first (validation needed)
  - US3 and US4 can proceed after foundational
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends only on Foundational - core validation
- **User Story 2 (P1)**: Depends only on Foundational - proto alignment verification
- **User Story 3 (P2)**: Depends only on Foundational - documentation
- **User Story 4 (P3)**: Depends only on Foundational - enum array validation

All user stories are independently testable after Foundational phase.

### Within Each User Story

- Tests/verification before implementation
- Schema updates before documentation
- Validation before moving to next story

### Parallel Opportunities

- T001, T002 can run in parallel (Setup phase)
- T004, T005, T006 can run in parallel (US1 tests)
- T009 can run in parallel with T010, T011 (US2)
- T016, T017 can run in parallel (Polish phase)
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all setup tasks together:
Task: "Create valid example registry in examples/registry.json"
Task: "Add validate:registry npm script in package.json"
Task: "Update validate npm script to include registry validation"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (verify infrastructure)
2. Complete Phase 2: Foundational (create schema)
3. Complete Phase 3: User Story 1 (validation works)
4. **STOP and VALIDATE**: Run `npm run validate:registry`
5. Deploy/merge if ready

### Incremental Delivery

1. Complete Setup + Foundational → Schema exists
2. Add User Story 1 → Validation works → MVP!
3. Add User Story 2 → Proto alignment documented
4. Add User Story 3 → Contribution docs complete
5. Add User Story 4 → Filtering support verified
6. Each story adds value without breaking previous stories

### Single Developer Strategy

Recommended execution order:

1. T001 → T002 (Setup)
2. T003 (Foundational - main schema)
3. T004 → T005 → T006 → T007 → T008 (User Story 1)
4. T009 → T010 → T011 (User Story 2)
5. T012 → T013 (User Story 3)
6. T014 → T015 (User Story 4)
7. T016 → T017 → T018 → T019 (Polish)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- This feature has 19 total tasks across 7 phases
