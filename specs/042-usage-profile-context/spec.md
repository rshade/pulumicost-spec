# Feature Specification: Usage Profile Context

**Feature Branch**: `042-usage-profile-context`
**Created**: 2026-01-27
**Status**: Draft
**Input**: User description: "Add a UsageProfile enum to the CostEstimationRequest
message to allow the Core to signal intent (Dev vs Prod) to plugins."

## Clarifications

### Session 2026-01-27

- Q: Which proto messages should include the usage_profile field?
  A: GetProjectedCostRequest and GetRecommendationsRequest (not GetActualCostRequest,
  since actual costs reflect what happened, not intent).
- Q: Should plugins log or emit metrics when applying profile-specific behavior?
  A: SHOULD log at INFO level with structured logging including the profile value.
- Q: How should plugins handle unknown/future profile values?
  A: MUST treat unknown values as UNSPECIFIED (graceful degradation), SHOULD log a
  warning when doing so.
- Q: Should SDK provide helpers for profile handling?
  A: Yes, SDK provides profile-aware builder methods (e.g., `WithProfileDefaults(profile)`)
  matching existing builder pattern.
- Q: How should profile context influence recommendation priorities?
  A: Plugin developers determine profile-specific recommendation behavior; spec provides
  the context, plugins define the meaning.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Development Environment Cost Estimation (Priority: P1)

As a developer running cost estimates in a development environment, I want the Core to
signal "dev" intent to plugins so that I receive cost recommendations optimized for
development workloads (burstable instances, reduced usage assumptions, cost-efficient
defaults).

**Why this priority**: Development environments are the most common use case for cost
estimation during infrastructure prototyping. Developers need quick, realistic estimates
without production-grade assumptions that inflate costs unnecessarily.

**Independent Test**: Can be fully tested by running a cost estimate with `--profile=dev`
flag and verifying plugins return development-appropriate recommendations (e.g., t3/t4g
instances instead of m5, 160hr usage instead of 730hr).

**Acceptance Scenarios**:

1. **Given** a cost estimation request with `usage_profile=DEV`, **When** the request is
   sent to an AWS compute plugin, **Then** the plugin should recommend burstable instance
   types and assume 160 hours/month usage.
2. **Given** a cost estimation request with `usage_profile=DEV`, **When** the request is
   sent to a storage plugin, **Then** the plugin should warn about over-provisioned IOPS
   and suggest standard configurations.
3. **Given** no usage profile is specified, **When** the request is sent to any plugin,
   **Then** the plugin should treat it as `UNSPECIFIED` and apply plugin-default behavior.

---

### User Story 2 - Production Environment Cost Estimation (Priority: P1)

As a platform engineer planning production infrastructure, I want the Core to signal
"prod" intent to plugins so that I receive cost estimates with production-grade
assumptions (full-time utilization, high availability, appropriate redundancy).

**Why this priority**: Production cost estimates are critical for budgeting and capacity
planning. Accurate production assumptions prevent under-provisioning and budget surprises.

**Independent Test**: Can be fully tested by running a cost estimate with `--profile=prod`
flag and verifying plugins apply production assumptions (730hr usage, production instance
classes, appropriate retention periods).

**Acceptance Scenarios**:

1. **Given** a cost estimation request with `usage_profile=PROD`, **When** the request is
   sent to a compute plugin, **Then** the plugin should assume 730 hours/month (full
   utilization) and recommend production-grade instance types.
2. **Given** a cost estimation request with `usage_profile=PROD`, **When** the request is
   sent to a database plugin, **Then** the plugin should apply production defaults
   including higher retention periods and provisioned IOPS.

---

### User Story 3 - Burst Workload Cost Estimation (Priority: P2)

As a DevOps engineer planning for temporary high-load scenarios (batch processing, load
testing, seasonal peaks), I want to signal "burst" intent to plugins so that I receive
cost estimates appropriate for short-term, high-intensity workloads.

**Why this priority**: Burst scenarios have distinct cost characteristics (higher data
transfer, scale-out recommendations) that differ from both dev and prod profiles. This
enables accurate planning for batch jobs, load tests, and seasonal traffic spikes.

**Independent Test**: Can be fully tested by running a cost estimate with `--profile=burst`
flag and verifying plugins return burst-appropriate recommendations (high data transfer
estimates, scale-out architectures).

**Acceptance Scenarios**:

1. **Given** a cost estimation request with `usage_profile=BURST`, **When** the request
   is sent to a compute plugin, **Then** the plugin should provide scale-out
   recommendations and high data transfer estimates.
2. **Given** a cost estimation request with `usage_profile=BURST`, **When** the request
   is sent to a network plugin, **Then** the plugin should factor in elevated bandwidth
   costs typical of load testing scenarios.

---

### Edge Cases

- What happens when a plugin does not recognize the usage profile value? Plugin MUST
  treat unknown values as UNSPECIFIED and SHOULD log a warning (enables forward
  compatibility when new profiles are added in future spec versions).
- How does the system handle conflicting profile signals? The request-level profile takes
  precedence; plugins should not override the Core's stated intent.
- What if a resource type has no profile-specific behavior? Plugin applies the same
  estimation regardless of profile, returning consistent results.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST define a UsageProfile enumeration with values: UNSPECIFIED (0),
  PROD (1), DEV (2), BURST (3).
- **FR-002**: System MUST include the usage_profile field in GetProjectedCostRequest
  and GetRecommendationsRequest messages (not GetActualCostRequest, since actual costs
  reflect historical data, not intent).
- **FR-003**: Plugins MUST treat UNSPECIFIED as "no preference" and apply their default
  estimation behavior.
- **FR-004**: Plugins MUST apply profile-appropriate defaults when a recognized profile
  is specified (DEV: reduced usage assumptions, burstable recommendations; PROD: full
  utilization, production-grade defaults; BURST: high-intensity, scale-out assumptions).
- **FR-005**: System MUST preserve backward compatibility - existing plugins without
  profile awareness must continue to function when receiving requests with any profile
  value.
- **FR-006**: Core MUST map CLI flag values (--profile=dev, --profile=prod, --profile=burst)
  to corresponding enum values before sending requests to plugins.
  _(Implementation scope: finfocus-core repository, not finfocus-spec)_
- **FR-007**: Plugins SHOULD document their profile-specific behavior in their DryRun
  capability response or documentation.
- **FR-008**: Plugins SHOULD log at INFO level when applying profile-specific behavior,
  including the usage_profile value in structured log output (e.g., `"usage_profile": "DEV"`).
- **FR-009**: Plugins MUST treat unrecognized usage_profile values as UNSPECIFIED and
  SHOULD log a warning, enabling forward compatibility with future spec versions.
- **FR-010**: SDK SHOULD provide profile-aware builder methods (e.g., `WithProfileDefaults(profile)`)
  on existing builders to simplify profile-specific cost estimation logic.

### Key Entities

- **UsageProfile**: An enumeration representing the intended workload context. Determines
  default assumptions plugins apply to cost calculations. Values map to distinct cost
  estimation strategies.
- **GetProjectedCostRequest**: Extended to include usage_profile field, allowing Core
  to communicate workload intent for cost projections.
- **GetRecommendationsRequest**: Extended to include usage_profile field, enabling
  context-aware recommendations (e.g., cost-saving for DEV, reliability for PROD).

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Developers using DEV profile receive cost estimates that are lower than
  PROD profile estimates for equivalent infrastructure configurations (reflecting reduced
  usage assumptions).
- **SC-002**: All existing plugins continue to function without modification when
  receiving requests with any usage_profile value (backward compatibility).
- **SC-003**: Users can switch between profiles via a single CLI flag without modifying
  resource configurations. _(Validation scope: finfocus-core repository)_
- **SC-004**: Plugin developers can determine the usage profile from the request and
  apply profile-specific logic in under 5 lines of code.

## Assumptions

- The Core (finfocus-core) CLI will implement the `--profile` flag and handle mapping to
  enum values.
- Plugins are responsible for defining what each profile means for their specific resource
  types (applies to both cost estimation and recommendations).
- UNSPECIFIED is the safe default that preserves current behavior for all plugins.
- Usage hour assumptions are documented guidelines, not enforced values - plugins have discretion:
  - DEV: ~160 hours/month (business hours assumption)
  - PROD: 730 hours/month (24/7 operation)
  - BURST: Plugin discretion - focus is on intensity and scale-out patterns, not specific
    hour counts. Typical scenarios include batch jobs (hours to days), load tests (hours),
    and seasonal peaks (days to weeks).

## Out of Scope

- Profile-specific pricing overrides (plugins control their own pricing logic).
- User-defined custom profiles beyond the three specified.
- Profile persistence or user preference storage.
- Validation of profile appropriateness for specific resource types.
