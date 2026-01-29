# Feature Specification: Documentation Drift Audit Remediation

**Feature Branch**: `043-docs-drift-audit`
**Created**: 2026-01-29
**Status**: Draft
**Input**: GitHub Issues #347 and #348 - Comprehensive documentation drift audit findings

## Clarifications

### Session 2026-01-29

- Q: Which version number should all README references be updated to? → A: v0.5.4 (latest git tag)

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Accurate Version Information (Priority: P1)

A developer evaluating FinFocus for integration needs to quickly understand what specification
version they are working with. They visit the README and expect to see a single, consistent
version number that matches the actual release state.

**Why this priority**: Version confusion is critical - developers may integrate against the wrong
API version, leading to runtime failures. This is the most visible and impactful documentation issue.

**Independent Test**: Can be tested by searching the README for version strings and verifying
all occurrences match. Delivers clarity on current specification version.

**Acceptance Scenarios**:

1. **Given** a developer opens README.md, **When** they search for version numbers,
   **Then** all version references display the same, correct specification version
2. **Given** version numbers in the title, introduction, and footer sections,
   **When** compared against each other, **Then** they are identical
3. **Given** the README version number, **When** compared to the latest git tag,
   **Then** they match or the README clearly indicates "unreleased" status

---

### User Story 2 - Accurate RPC Documentation (Priority: P1)

A plugin developer needs to understand what RPCs are available in the CostSourceService.
They reference the README which claims "8 RPC methods" but discover there are actually
11 RPCs when reading the proto file, causing confusion about feature completeness.

**Why this priority**: Incorrect RPC counts lead developers to believe they've missed
functionality or that documentation is outdated, eroding trust in the specification.

**Independent Test**: Can be tested by counting RPCs in the proto file and comparing to
README claims. Delivers accurate service interface documentation.

**Acceptance Scenarios**:

1. **Given** a developer reads the CostSourceService description in README,
   **When** they count the listed RPC methods, **Then** the count matches the proto definition (11 RPCs)
2. **Given** the README lists RPC methods, **When** compared to proto file,
   **Then** all RPCs are documented: Name, Supports, GetActualCost, GetProjectedCost,
   GetPricingSpec, EstimateCost, GetRecommendations, DismissRecommendation, GetBudgets,
   GetPluginInfo, DryRun
3. **Given** observability RPCs exist (HealthCheck, GetMetrics, GetServiceLevelIndicators),
   **When** documented, **Then** they are clearly distinguished from core CostSourceService RPCs

---

### User Story 3 - Correct SDK Code Examples (Priority: P1)

A plugin developer follows the testing README examples to write conformance tests. They copy
the example `result := RunBasicConformance(plugin)` but get compile errors because the
actual function returns `(*ConformanceResult, error)`.

**Why this priority**: Broken code examples waste developer time and create friction.
Code that doesn't compile is the worst kind of documentation error.

**Independent Test**: Can be tested by compiling all code examples from READMEs.
Delivers working, copy-paste-ready code snippets.

**Acceptance Scenarios**:

1. **Given** code examples in testing/README.md, **When** copied into a Go file,
   **Then** they compile without modification
2. **Given** RunBasicConformance, RunStandardConformance, RunAdvancedConformance examples,
   **When** shown in README, **Then** they demonstrate proper error handling for the
   `(*ConformanceResult, error)` return type
3. **Given** any documented function, **When** its signature is shown,
   **Then** it matches the actual exported function signature

---

### User Story 4 - Complete Package Documentation (Priority: P2)

A developer discovers the `sdk/go/pluginsdk/mapping/` package while exploring the codebase.
They want to understand how to use the AWS, Azure, and GCP property extraction helpers
but find no README.md to explain the package purpose and usage patterns.

**Why this priority**: Undocumented packages reduce SDK discoverability and adoption.
The mapping package provides valuable cross-cloud helpers that developers may not find.

**Independent Test**: Can be tested by checking for README.md presence in mapping/ package.
Delivers discoverable documentation for GitHub browsing.

**Acceptance Scenarios**:

1. **Given** a developer browses sdk/go/pluginsdk/mapping/ on GitHub,
   **When** they look for documentation, **Then** a README.md exists explaining the package
2. **Given** the mapping package README, **When** read by a developer,
   **Then** it explains the purpose of aws.go, azure.go, gcp.go, common.go, and keys.go
3. **Given** property extraction helpers exist, **When** documented,
   **Then** usage examples show how to extract SKU, region, and other properties

---

### User Story 5 - Accurate Example Counts (Priority: P2)

A developer reads that there are "10 comprehensive pricing examples" but finds a different
number when listing the examples/specs/ directory, causing uncertainty about completeness.

**Why this priority**: Incorrect counts create doubt about documentation accuracy and
whether examples are missing.

**Independent Test**: Can be tested by counting JSON files in examples/specs/ and comparing
to README claims. Delivers accurate example inventory.

**Acceptance Scenarios**:

1. **Given** the README states an example count, **When** compared to actual JSON files
   in examples/specs/, **Then** the count matches exactly (currently 9 JSON files)
2. **Given** example documentation, **When** listing covered billing models,
   **Then** the list matches actual example file contents

---

### User Story 6 - Documented SDK Helpers (Priority: P3)

A developer building a plugin wants to use SDK helper functions like `NewActualCostResponse()`
with functional options but cannot find documentation on available helpers, options, or
validation functions in the root README or SDK documentation.

**Why this priority**: Undocumented helpers reduce SDK usability. Developers may
implement features manually when helpers exist, leading to inconsistent implementations.

**Independent Test**: Can be tested by checking README for key helper function documentation.
Delivers discoverable API surface documentation.

**Acceptance Scenarios**:

1. **Given** NewActualCostResponse() exists in pluginsdk, **When** documented,
   **Then** the root README or SDK README shows usage with functional options
2. **Given** FallbackHint enum is used in plugin orchestration, **When** documented,
   **Then** users understand NONE, RECOMMENDED, REQUIRED hint semantics
3. **Given** validation helpers exist (ValidateActualCostResponse, ValidateRecommendation),
   **When** documented, **Then** developers know to use them before returning responses

---

### Edge Cases

- What happens when documentation references a file that has been renamed or moved?
  - Documentation MUST use relative links validated by CI
- How does the system handle version bumps?
  - A single source of truth for version (e.g., version constant) should propagate to README
- What if a new RPC is added to the proto?
  - Documentation update MUST be part of the proto change PR

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: All version numbers in README.md MUST be consistent (single version across title,
  introduction, and footer)
- **FR-002**: RPC method count in README MUST match the actual count in costsource.proto
  (currently 11 in CostSourceService + 3 in observability)
- **FR-003**: Example count in README MUST match the actual number of JSON files in
  examples/specs/ (currently 9)
- **FR-004**: All code examples in SDK READMEs MUST compile without modification
- **FR-005**: Function signatures shown in documentation MUST match actual exported signatures
- **FR-006**: The sdk/go/pluginsdk/mapping/ package MUST have a README.md for GitHub browsing
- **FR-007**: Key SDK helpers MUST be documented: NewActualCostResponse, FallbackHint enum,
  validation functions
- **FR-008**: Testing README MUST show proper error handling for conformance functions that
  return `(*ConformanceResult, error)`
- **FR-009**: All mandatory conformance test functions MUST be documented in testing/README.md

### Key Entities

- **Documentation File**: A markdown file (README.md, CLAUDE.md) containing user-facing content
  that must stay synchronized with source code
- **Code Example**: A snippet embedded in documentation that must compile and run correctly
- **Version Reference**: Any occurrence of a version string that must be consistent across documents
- **SDK Helper**: An exported function designed to simplify plugin development

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: 100% of version references in README.md display the same version number
- **SC-002**: RPC count in documentation matches proto file count exactly
- **SC-003**: Example count in documentation matches actual file count in examples/specs/
- **SC-004**: 100% of code examples in testing/README.md compile without errors
- **SC-005**: mapping/ package has README.md with package description and usage examples
- **SC-006**: At least one code example demonstrates NewActualCostResponse() with functional options
- **SC-007**: FallbackHint enum is documented in user-facing documentation (README or SDK docs)
- **SC-008**: All conformance functions with error returns have examples showing error handling

## Assumptions

- The current specification version is v0.5.4 (based on latest git tag)
- ~~Both v0.5.0 (title/intro) and v0.4.7 (footer) references in README are outdated~~
  **RESOLVED**: All 4 version references (lines 1, 17, 832, 858) already show v0.5.4 as of 2026-01-29
- Documentation changes require no breaking changes to code
- Markdown linting (markdownlint-cli2) will be run on all modified documentation
- No new SDK helpers need to be created; only documentation for existing helpers is needed

## Dependencies

- Access to proto/finfocus/v1/costsource.proto for accurate RPC counting
- Access to sdk/go/testing/ for accurate function signatures
- Access to sdk/go/pluginsdk/ for helper function inventory
- Access to examples/specs/ for example file counting

## Out of Scope

- Adding godoc coverage checking to CI (Phase 4 infrastructure - separate feature)
- Adding markdown link validation to CI (Phase 4 infrastructure - separate feature)
- Code example compilation testing in CI (Phase 4 infrastructure - separate feature)
- Restructuring SDK package organization
- Adding new SDK helper functions (only documenting existing ones)
