# Implementation Plan: Documentation Drift Audit Remediation

**Branch**: `043-docs-drift-audit` | **Date**: 2026-01-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/043-docs-drift-audit/spec.md` + GitHub Issue #348

## Summary

Fix documentation drift across README.md and SDK documentation files. **Version consistency is
already resolved** (all v0.5.4). Remaining issues: correct RPC counts (claims 8, actually 11),
accurate example counts (claims 10/8, actually 9), compilable code examples in testing README,
and complete package documentation (mapping/ README).

## Technical Context

**Language/Version**: Markdown documentation (no code changes)
**Primary Dependencies**: N/A (documentation only)
**Storage**: N/A
**Testing**: Markdown linting (markdownlint-cli2), manual verification
**Target Platform**: GitHub repository documentation
**Project Type**: Documentation maintenance
**Performance Goals**: N/A
**Constraints**: Must pass `make lint-markdown`
**Scale/Scope**: 4 markdown files to update, 1 new README to create

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Contract First**: N/A - No proto changes (documentation only)
- [x] **Spec Consumes**: N/A - No pricing logic changes
- [x] **Multi-Provider**: N/A - No new provider examples
- [x] **FinFocus Alignment**: Yes - Updates use FinFocus naming consistently
- [x] **SDK Synchronization**: N/A - No proto changes requiring SDK updates
- [x] **Documentation Integrity**: **YES** - This feature directly implements Constitution Section XIV
  (Documentation Integrity) by fixing README ↔ Code sync and ensuring Package Completeness

**Note**: This feature is a direct remediation of Constitution XIV violations. No new violations introduced.

## Project Structure

### Documentation (this feature)

```text
specs/043-docs-drift-audit/
├── plan.md              # This file
├── research.md          # Phase 0 output - verification of current state
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Specification validation checklist
```

### Files to Modify (repository root)

```text
# P1 Critical Fixes
README.md                           # Version, RPC count, example count fixes
sdk/go/testing/README.md            # Function signature fixes

# P2 Package Completeness
sdk/go/pluginsdk/mapping/README.md  # NEW FILE - package documentation

# P3 SDK Helper Documentation (location TBD in research)
README.md or sdk/go/pluginsdk/README.md  # FallbackHint, NewActualCostResponse docs
```

**Structure Decision**: Documentation-only changes to existing markdown files plus one new
README.md for the mapping package.

## Complexity Tracking

No constitution violations to track - this feature remediates existing violations.

## Implementation Phases

### Phase 1: P1 Critical Fixes

**Files**: `README.md`, `sdk/go/testing/README.md`

1. **Version Consistency** (FR-001) ✅ ALREADY RESOLVED
   - All version references already show v0.5.4 (lines 1, 17, 832, 858)
   - No action needed

2. **RPC Count Accuracy** (FR-002)
   - Line 80: "8 RPC methods" → "11 RPC methods"
   - Line 699: "8 RPC methods" → "11 RPC methods"
   - The 11 CostSourceService RPCs: Name, Supports, GetActualCost, GetProjectedCost,
     GetPricingSpec, EstimateCost, GetRecommendations, DismissRecommendation, GetBudgets,
     GetPluginInfo, DryRun

3. **Example Count Accuracy** (FR-003)
   - Line 65: "10 comprehensive pricing examples" → "9 comprehensive pricing examples"
   - Line 660: "8 comprehensive examples" → "9 comprehensive examples"

4. **Function Signature Fixes** (FR-004, FR-005, FR-008)
   - `sdk/go/testing/README.md` lines 104, 171-174
   - Change: `result := RunBasicConformance(plugin)` →
     `result, err := RunBasicConformance(plugin)`
   - Add error handling to all conformance examples
   - Verify all documented function signatures match exports

### Phase 2: P2 Package Completeness

**Files**: `sdk/go/pluginsdk/mapping/README.md` (NEW)

1. **Create mapping/ README** (FR-006)
   - Document package purpose: cross-cloud property extraction
   - Document files: aws.go, azure.go, gcp.go, common.go, keys.go
   - Add usage examples for SKU/region extraction
   - Follow existing SDK README patterns

### Phase 3: P3 SDK Helper Documentation

**Files**: `README.md` or `sdk/go/pluginsdk/README.md`

1. **Document NewActualCostResponse** (FR-007, SC-006)
   - Add code example showing functional options pattern
   - Show WithFallbackHint, WithResults usage

2. **Document FallbackHint Enum** (FR-007, SC-007)
   - Document enum values: UNSPECIFIED, NONE, RECOMMENDED, REQUIRED
   - Explain orchestration semantics

3. **Document Validation Helpers** (FR-007)
   - Reference ValidateActualCostResponse, ValidateRecommendation
   - Show validation-before-return pattern

4. **Update Root README Code Examples** (from Issue #348)
   - Lines 259-279: Update manual struct construction to use SDK helpers
   - Show `pluginsdk.NewActualCostResponse()` pattern instead of manual `&pbc.GetActualCostResponse{}`

### Phase 4: P4 Godoc Coverage (from Issue #348)

**Files**: `sdk/go/pluginsdk/*.go`, `sdk/go/pricing/*.go`

1. **Add Godoc to Undocumented pluginsdk Functions** (~25 functions)
   - Response builders: NewActualCostResponse, NewProjectedCostResponse, NewDryRunResponse, etc.
   - Validation helpers: ValidateActualCostResponse, ValidateProjectedCostResponse, etc.
   - Field mapping utilities: FocusFieldNames, NewFieldMapping, AllFieldsWithStatus, etc.

2. **Add Godoc to Undocumented pricing Functions** (~10 functions)
   - ApplyGrowth, ValidateGrowthParams, CheckGrowthWarnings
   - ValidatePricingSpec, ValidateBillingMode

3. **Create registry/ Package README** (from Issue #348)
   - Document 8 enum types with validators
   - Located at `sdk/go/registry/`

### Phase 5: Verification

1. Run `make lint-markdown` to verify formatting
2. Verify all version references match v0.5.4
3. Count JSON files in examples/specs/ matches documentation
4. Manually verify code examples compile (future: CI automation)
5. Check godoc coverage for pluginsdk and pricing packages

## Artifacts Generated

- [x] `plan.md` - This file
- [x] `research.md` - Verification of current state and exact locations
- N/A `data-model.md` - Not needed (documentation feature)
- N/A `contracts/` - Not needed (no API changes)
- N/A `quickstart.md` - Not needed (documentation feature)

## Key Findings from Research

### From Issue #347 (Original Audit)

1. **Version consistency is RESOLVED** - All 4 version references now show v0.5.4
2. **RPC count needs 2 fixes** - Lines 80 and 699 claim "8", actual is 11
3. **Example count needs 2 fixes** - Lines 65 (claims 10) and 660 (claims 8), actual is 9
4. **Testing README needs 4 code block updates** - All conformance function examples
5. **Mapping README needs creation** - New ~100 line file
6. **Helper docs partially exist** - FallbackHint documented in CLAUDE.md, needs user-facing docs

### From Issue #348 (Additional Scope)

1. **Root README code examples outdated** - Lines 259-279 use manual struct instead of SDK helpers
2. **Godoc coverage gaps in pluginsdk** - ~25 exported functions undocumented
3. **Godoc coverage gaps in pricing** - ~10 exported functions undocumented
4. **Registry package needs README** - 8 enum types need package documentation

## Next Steps

Proceed to `/speckit.tasks` for task breakdown
