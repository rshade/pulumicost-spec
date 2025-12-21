# Feature Specification: GetRecommendations RPC for FinOps Optimization

**Feature Branch**: `013-recommendations-rpc`
**Created**: 2025-12-04
**Status**: Draft
**GitHub Issue**: [#122](https://github.com/rshade/pulumicost-spec/issues/122)
**Input**: Add GetRecommendations RPC to CostSourceService for FinOps recommendations

## Overview

This feature adds a `GetRecommendations` RPC to the `CostSourceService` that enables plugins to
surface cost optimization recommendations from various FinOps platforms (AWS Cost Explorer,
Kubecost, Azure Advisor, GCP Recommender). The interface provides a unified way to expose
actionable optimization opportunities to users.

## Clarifications

### Session 2025-12-04

- Q: How should plugins indicate they support recommendations? → A: Both auto-detection
  via optional interface AND deterministic capability declaration via metadata
- Q: What observability should GetRecommendations include? → A: Full observability with
  structured logging (zerolog), Prometheus metrics, and trace context propagation

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Retrieve All Recommendations (Priority: P1)

As a DevOps engineer using PulumiCost, I want to retrieve all cost optimization recommendations
from my connected cost sources so that I can identify opportunities to reduce cloud spending.

**Why this priority**: This is the core functionality - without the ability to retrieve
recommendations, no other features can work. Provides immediate value by surfacing
cost-saving opportunities.

**Independent Test**: Can be fully tested by calling GetRecommendations with an empty filter
and verifying recommendations are returned with required fields populated.

**Acceptance Scenarios**:

1. **Given** a plugin that supports recommendations, **When** GetRecommendations is called
   with no filter, **Then** all available recommendations are returned with valid category,
   action type, resource info, and estimated savings
2. **Given** a plugin that does not implement recommendations, **When** GetRecommendations
   is called, **Then** an empty list is returned without error
3. **Given** multiple recommendations exist, **When** GetRecommendations is called, **Then**
   a summary with total count and aggregated savings is included in the response

---

### User Story 2 - Filter Recommendations by Category (Priority: P2)

As a FinOps practitioner, I want to filter recommendations by category (cost, performance,
security, reliability) so that I can focus on specific optimization areas during my analysis.

**Why this priority**: Filtering enables targeted analysis, essential for large organizations
with many recommendations. Depends on P1 but provides significant value for practical use.

**Independent Test**: Can be tested by creating mock recommendations across categories and
verifying filters return only matching items.

**Acceptance Scenarios**:

1. **Given** recommendations exist in multiple categories, **When** filtering by COST
   category, **Then** only cost-related recommendations are returned
2. **Given** recommendations exist in multiple categories, **When** filtering by PERFORMANCE
   category, **Then** only performance recommendations are returned
3. **Given** a filter with multiple criteria (category + provider), **When** GetRecommendations
   is called, **Then** only recommendations matching ALL criteria are returned

---

### User Story 3 - Filter by Action Type (Priority: P2)

As a cloud architect, I want to filter recommendations by action type (rightsize, terminate,
purchase commitment) so that I can plan specific types of optimization initiatives.

**Why this priority**: Allows users to focus on actionable items matching their current
optimization initiative (e.g., rightsizing campaign).

**Independent Test**: Can be tested by creating mock recommendations with different action
types and verifying filter accuracy.

**Acceptance Scenarios**:

1. **Given** rightsizing and termination recommendations exist, **When** filtering by
   RIGHTSIZE action type, **Then** only rightsizing recommendations are returned
2. **Given** Kubernetes request adjustment recommendations exist, **When** filtering by
   ADJUST_REQUESTS action type, **Then** only container sizing recommendations are returned

---

### User Story 4 - Paginate Large Result Sets (Priority: P3)

As a system administrator managing multiple cloud accounts, I want to paginate through large
recommendation sets so that I can handle thousands of recommendations without overwhelming
system resources.

**Why this priority**: Essential for scalability but only relevant for users with large
environments. Core functionality works without pagination.

**Independent Test**: Can be tested by creating a large mock dataset and verifying page
tokens work correctly.

**Acceptance Scenarios**:

1. **Given** 500 recommendations exist and page_size is 100, **When** GetRecommendations is
   called, **Then** 100 recommendations and a next_page_token are returned
2. **Given** a valid page_token from a previous request, **When** GetRecommendations is called
   with that token, **Then** the next page of results is returned
3. **Given** the last page of results, **When** retrieved, **Then** next_page_token is empty

---

### User Story 5 - View Provider-Specific Recommendation Details (Priority: P3)

As a cloud engineer, I want to see provider-specific details for each recommendation (e.g.,
AWS instance type changes, Kubernetes container resource adjustments) so that I can
understand exactly what action to take.

**Why this priority**: Enriches the user experience but core identification of savings
opportunities works without detailed action breakdowns.

**Independent Test**: Can be tested by verifying action-specific fields are populated
correctly for each recommendation type.

**Acceptance Scenarios**:

1. **Given** an AWS rightsizing recommendation, **When** retrieved, **Then** current and
   recommended instance types are populated
2. **Given** a Kubernetes request sizing recommendation, **When** retrieved, **Then** cluster,
   namespace, controller, and container details are populated
3. **Given** a commitment purchase recommendation, **When** retrieved, **Then** commitment
   type, term, and recommended quantity are populated

---

### Edge Cases

- What happens when a plugin returns recommendations with invalid currency codes?
  - System validates currency against ISO 4217 and rejects invalid codes with appropriate error
- What happens when estimated savings is negative (cost increase)?
  - System accepts negative values to represent performance/reliability improvements that
    may increase cost
- What happens when pagination token is invalid or expired?
  - System returns an error indicating invalid token; client should restart from first page
- How does system handle recommendations without resource utilization data?
  - Utilization fields are optional; recommendations are valid without them
- What happens when confidence score is outside 0.0-1.0 range?
  - System validates confidence score bounds and rejects invalid values

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide a GetRecommendations operation that returns cost
  optimization recommendations
- **FR-002**: System MUST support filtering recommendations by provider, region, resource
  type, category, and action type
- **FR-003**: System MUST support pagination with configurable page size and continuation
  tokens
- **FR-004**: System MUST return a summary with total count, aggregated savings, and
  breakdowns by category and action type
- **FR-005**: System MUST support multiple recommendation categories: Cost, Performance,
  Security, Reliability
- **FR-006**: System MUST support multiple action types: Rightsize, Terminate, Purchase
  Commitment, Adjust Requests, Modify, Delete Unused
- **FR-007**: System MUST support multiple priority levels: Low, Medium, High, Critical
- **FR-008**: System MUST include resource information (ID, name, provider, type, region,
  tags, utilization) for each recommendation
- **FR-009**: System MUST include impact information (estimated savings, currency, projection
  period, current/projected cost) for each recommendation
- **FR-010**: System MUST validate currency codes against ISO 4217 standard using the
  existing currency package
- **FR-011**: System MUST support action-specific details via distinct action types
  (rightsize, terminate, commitment, kubernetes, modify)
- **FR-012**: System MUST return an empty list (not error) when a plugin does not support
  recommendations
- **FR-013**: System MUST include a projection_period parameter supporting "daily",
  "monthly" (default), and "annual" values
- **FR-014**: System MUST support optional confidence scores between 0.0 and 1.0 for each
  recommendation
- **FR-015**: System MUST include reasoning information explaining why each recommendation
  was generated
- **FR-016**: PluginSDK MUST provide an optional `RecommendationsProvider` interface that
  plugins can implement to opt-in to recommendations support
- **FR-017**: PluginSDK MUST auto-detect when a plugin implements the `RecommendationsProvider`
  interface and route GetRecommendations calls accordingly
- **FR-018**: Plugins MUST be able to declare recommendations capability via metadata in
  the `Supports` response for deterministic capability discovery
- **FR-019**: PluginSDK MUST maintain backward compatibility - existing plugins without
  recommendations support continue to work without modification
- **FR-020**: GetRecommendations MUST emit structured log entries using zerolog following
  existing logging patterns (request received, completion, errors)
- **FR-021**: GetRecommendations MUST expose Prometheus metrics for request count, latency
  histogram, error count, and recommendation count per response
- **FR-022**: GetRecommendations MUST propagate trace context for distributed tracing
  integration

### Key Entities

- **Recommendation**: A single optimization suggestion with category, action type, affected
  resource, impact assessment, priority, and source attribution
- **RecommendationFilter**: Criteria for narrowing recommendations by provider, region,
  resource type, category, or action type
- **RecommendationImpact**: Financial impact details including estimated savings, current
  cost, projected cost, and savings percentage
- **RecommendationSummary**: Aggregated view of recommendations with counts and savings
  totals broken down by category and action type
- **ResourceInfo**: Information about the cloud resource being recommended for optimization,
  including utilization metrics
- **ResourceUtilization**: Current utilization metrics (CPU, memory, storage, network) used
  to justify recommendations
- **Action Details**: Provider-specific details for each action type (RightsizeAction,
  TerminateAction, CommitmentAction, KubernetesAction, ModifyAction)
- **RecommendationsProvider**: Optional interface that plugins implement to provide
  recommendations; SDK auto-detects implementation for capability routing

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Plugin developers can implement recommendations support within one development
  sprint
- **SC-002**: All recommendation responses complete within 500ms for result sets under
  100 items
- **SC-003**: 100% of recommendations include valid category, action type, and estimated
  savings
- **SC-004**: Currency validation catches 100% of non-ISO 4217 currency codes
- **SC-005**: Pagination works correctly for result sets up to 10,000 recommendations
- **SC-006**: Filters accurately return only matching recommendations with zero false
  positives
- **SC-007**: Summary totals match the sum of individual recommendation impacts with no
  calculation errors
- **SC-008**: Conformance tests cover Basic (response validation), Standard (filtering,
  pagination), and Advanced (performance) levels
- **SC-009**: GetRecommendations emits metrics and logs consistent with other RPC methods
  in the PluginSDK

## Assumptions

- Recommendations are point-in-time snapshots; state tracking (applied/dismissed) is out of
  scope for this feature
- Each plugin is responsible for fetching recommendations from its respective backend service
- The projection period (daily/monthly/annual) affects how estimated savings are calculated
  but implementation is plugin-specific
- Default page size is 50 recommendations when not specified
- Empty filter returns all available recommendations (no filtering applied)
- Plugins opt-in to recommendations via optional interface implementation; SDK handles
  detection automatically and returns empty list for non-implementing plugins
- Plugins can also declare capability via `Supports` metadata for deterministic discovery
- Mock plugin will generate realistic sample data for all supported recommendation types

## Out of Scope

- CLI output formatting (follow-up in pulumicost-core)
- Aggregation logic for combining recommendations from multiple plugins
- Caching of recommendations
- State management for tracking applied/dismissed recommendations
- Automatic remediation or execution of recommendations
