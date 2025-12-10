# Feature Specification: GetBudgets RPC for Plugin-Provided Budget Information

**Feature Branch**: `001-get-budgets-rpc`
**Created**: 2025-12-09
**Status**: Draft
**Input**: User description: "Add a GetBudgets RPC to the CostSource service that allows
plugins to provide budget information from cloud cost management services (AWS Budgets, GCP
Billing Budgets, Azure Cost Management, Kubecost)."

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Unified Budget Visibility Across Providers (Priority: P1)

As a FinOps engineer managing cloud costs across multiple providers, I want to see all my
budgets (AWS, GCP, Azure, Kubecost) in a single pulumicost output so that I have a unified
view of my spending limits and can make informed decisions about resource allocation.

**Why this priority**: This is the core value proposition - consolidating budget information
from disparate sources into one view, which is essential for multi-cloud cost management.

**Independent Test**: Can be fully tested by configuring plugins for different providers and
verifying that budget information from each appears in the unified output, delivering
immediate value for budget monitoring.

**Acceptance Scenarios**:

1. **Given** AWS and GCP plugins are configured with budget data, **When** I run pulumicost
   commands, **Then** I see budgets from both providers in the output
2. **Given** a budget exceeds its threshold in AWS, **When** I view budget status, **Then**
   the exceeded budget is clearly marked and highlighted
3. **Given** multiple budgets with different currencies, **When** I view the summary, **Then**
   all amounts are displayed in their original currencies

---

### User Story 2 - Kubernetes Budget Tracking (Priority: P2)

As a Kubernetes operator responsible for cost governance, I want Kubecost budgets to be
exposed through pulumicost so that I can track namespace spending limits alongside
infrastructure costs.

**Why this priority**: Enables comprehensive cost visibility including container orchestration
budgets, which is critical for organizations running Kubernetes workloads.

**Independent Test**: Can be fully tested by configuring Kubecost plugin and verifying
namespace budget information appears in pulumicost output, providing standalone value for
Kubernetes cost management.

**Acceptance Scenarios**:

1. **Given** Kubecost is configured with namespace budgets, **When** I query budgets, **Then**
   I see namespace-level spending limits and current usage
2. **Given** a namespace budget is approaching its limit, **When** I check budget status,
   **Then** I receive appropriate warning indicators
3. **Given** multiple namespaces have budgets, **When** I view the summary, **Then** I see
   aggregated budget health across all namespaces

---

### User Story 3 - Multi-Cloud Budget Health Overview (Priority: P3)

As a multi-cloud user managing budgets across different cloud providers, I want aggregated
budget status so that I can understand my total budget health at a glance.

**Why this priority**: Provides high-level visibility into overall budget compliance across
the entire cloud estate, supporting executive-level decision making.

**Independent Test**: Can be fully tested by aggregating budget data from multiple providers
and displaying summary statistics, offering standalone value for high-level budget monitoring.

**Acceptance Scenarios**:

1. **Given** budgets across multiple providers, **When** I request budget summary, **Then**
   I see total number of budgets and breakdown by health status (OK, warning, exceeded)
2. **Given** some budgets are in warning state, **When** I view aggregated status, **Then**
   the overall health reflects the most critical budget status
3. **Given** forecasted spending data is available, **When** I check budget status, **Then**
   I see both current and forecasted budget utilization

---

### Edge Cases

- What happens when a plugin doesn't support budgets (should gracefully handle unimplemented
  functionality)?
- How does the system handle budgets with different currencies (should display in original
  currency)?
- What happens when budget data is temporarily unavailable (should show last known status or
  indicate unavailability)?
- How does the system handle very large numbers of budgets (should provide filtering and
  pagination)?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST allow plugins to provide budget information from cloud cost
  management services
- **FR-002**: System MUST display budget information in a unified format regardless of the
  source provider
- **FR-003**: System MUST show current spending status and utilization percentages for each
  budget
- **FR-004**: System MUST support filtering budgets by provider, region, resource type, and
  tags
- **FR-005**: System MUST indicate budget health status (OK, warning, critical, exceeded)
  based on threshold rules
- **FR-006**: System MUST provide aggregated budget summary showing total budgets and
  breakdown by health status
- **FR-007**: System MUST handle optional budget functionality gracefully when plugins don't
  support it
- **FR-008**: System MUST display budget amounts in their original currencies
- **FR-009**: System MUST show budget period information (daily, weekly, monthly, etc.)
- **FR-010**: System MUST support both actual and forecasted budget thresholds

### Key Entities _(include if feature involves data)_

- **Budget**: Represents a spending limit with alert thresholds, containing identification,
  amount limits, time periods, and status information
- **Budget Amount**: Specifies the monetary limit and currency for a budget
- **Budget Period**: Defines the time interval for budget calculations (daily, monthly,
  annually, etc.)
- **Budget Filter**: Allows narrowing down budgets by provider, region, resource type, or tags
- **Budget Threshold**: Defines alert points with percentages and types (actual vs forecasted)
- **Budget Status**: Shows current spending, utilization percentage, and health assessment
- **Budget Summary**: Provides aggregated statistics across multiple budgets

## Clarifications

### Session 2025-12-09

- Q: What are the expected scale limits for budgets? → A: Medium scale (100-1000 budgets)
  for departments
- Q: What authentication/authorization is needed for budget access? → A: Plugin developers
  deal with auth, not the spec
- Q: What logging/metrics are required for budget operations? → A: Use existing logging specs
- Q: How should plugin failures be handled when retrieving budgets? → A: Use existing specs
  for guidance

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Users can view budgets from all configured cloud providers in under 5 seconds
- **SC-002**: Budget information displays utilization percentages accurately within 1% of
  source data
- **SC-003**: System successfully retrieves budget data from at least 95% of configured plugins
- **SC-004**: Users report 80% satisfaction with budget visibility and unified display
- **SC-005**: Budget threshold alerts trigger within 1 minute of threshold being crossed
  (alerts handled by pulumicost-core consuming applications)
- **SC-006**: Multi-cloud users can assess overall budget health in under 30 seconds
- **SC-007**: System supports 100-1000 budgets per user/department with acceptable performance
