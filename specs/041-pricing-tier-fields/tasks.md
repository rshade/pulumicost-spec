# Implementation Tasks: Add Pricing Tier & Probability Fields

**Feature Branch**: `041-pricing-tier-fields`
**Input**: Feature plan from `/specs/041-pricing-tier-fields/plan.md`

## Phase 1: Setup

<!--
  Project initialization and infrastructure setup.
  No story dependencies.
-->

- [x] T001 Verify clean git state and branch is `041-pricing-tier-fields`
- [x] T002 Verify buf and go tools are installed and up to date

## Phase 2: Foundational

<!--
  Blocking prerequisites for all user stories.
  Contract updates must happen before any SDK generation or usage.
-->

- [x] T003 Update `proto/finfocus/v1/costsource.proto` to add `pricing_category` (field 3) and
      `spot_interruption_risk_score` (field 4) to `EstimateCostResponse` message
- [x] T004 Update `proto/finfocus/v1/costsource.proto` to add `pricing_category` (field 8) and
      `spot_interruption_risk_score` (field 9) to `GetProjectedCostResponse` message
- [x] T005 Run `make generate` to regenerate Go and TypeScript SDKs from updated protos

## Phase 3: Identify Spot Instance Risk (P1)

<!--
  User Story 1: Identify Spot Instance Risk
  Goal: Enable plugins to report Spot/Dynamic pricing and associated risk scores.
  Dependencies: Phase 2
-->

- [x] T006 [P] [US1] Create a new test file `sdk/go/spot_risk_test.go` to verify `EstimateCostResponse` fields
- [x] T007 [P] [US1] Create a new test file `sdk/go/projected_cost_test.go` to verify `GetProjectedCostResponse` fields
- [x] T008 [US1] Implement unit test in `sdk/go/spot_risk_test.go` that constructs `EstimateCostResponse`
      with `FOCUS_PRICING_CATEGORY_DYNAMIC` and risk score `0.8`
- [x] T009 [US1] Implement unit test in `sdk/go/projected_cost_test.go` that constructs `GetProjectedCostResponse`
      with `FOCUS_PRICING_CATEGORY_DYNAMIC` and risk score `0.8`
- [x] T010 [US1] Run `go test ./sdk/go/...` to verify the new fields can be set and retrieved correctly

## Phase 4: Explain Cost Basis (P2)

<!--
  User Story 2: Explain Cost Basis
  Goal: Enable plugins to report Committed/Reserved pricing.
  Dependencies: Phase 2
-->

- [x] T011 [US2] Add test case to `sdk/go/spot_risk_test.go` for `FOCUS_PRICING_CATEGORY_COMMITTED`
      scenario (Reserved Instance)
- [x] T012 [US2] Add test case to `sdk/go/projected_cost_test.go` for `FOCUS_PRICING_CATEGORY_COMMITTED`
      scenario (Savings Plan)
- [x] T013 [US2] Run `go test ./sdk/go/...` to verify committed pricing scenarios

## Final Phase: Polish

<!--
  Cross-cutting concerns, linting, and documentation.
-->

- [x] T014 Run `make buf-lint` to ensure proto style compliance
- [x] T015 Verify `sdk/typescript/` builds correctly with `npm run build` (or equivalent)
- [x] T016 Commit all changes

## Dependencies

- **Phase 1 -> Phase 2**: Setup required before proto modification.
- **Phase 2 -> Phase 3 & 4**: Proto updates and SDK regeneration are blocking prerequisites for all usage tests.
- **Phase 3 & 4**: Can be executed in parallel as they test different aspects of the same generated code.

## Implementation Strategy

1. **Contract First**: We modify the `.proto` definitions first. This is the source of truth.
2. **Generate**: We immediately regenerate SDKs to make the new fields available in code.
3. **Verify via Tests**: Since this is a spec/SDK change (not a backend service implementation), "implementation"
   means verifying the SDK consumes the new contract correctly. We write distinct tests for "Spot Risk" (US1) and
   "Committed Pricing" (US2) to validate the data model.
