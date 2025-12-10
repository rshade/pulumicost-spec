# Research: GetBudgets RPC Implementation

**Date**: 2025-12-09
**Feature**: GetBudgets RPC for Plugin-Provided Budget Information

## Research Tasks Completed

### Provider Budget API Analysis

#### AWS Budgets API

**Decision**: Use existing AWS Budgets API patterns with BudgetName, BudgetLimit, TimeUnit,
and Notification configurations.

**Rationale**: AWS Budgets API provides comprehensive budget management with alerts and
forecasting. The proto design maps cleanly to existing AWS structures.

**Alternatives considered**: Custom budget format - rejected because it would require
plugin-side transformation and lose AWS-native features.

#### GCP Cloud Billing Budget API

**Decision**: Map to GCP's budgetFilter, amount.specifiedAmount, and thresholdRules structures.

**Rationale**: GCP provides calendar-based budgeting with percentage thresholds. Proto fields
align well with GCP's hierarchical budget structure.

**Alternatives considered**: Simplify to basic amount/limit - rejected because GCP supports
complex filtering by projects/services.

#### Azure Cost Management Budgets

**Decision**: Support Azure's category-based budgets with timeGrain and notification settings.

**Rationale**: Azure budgets integrate with Cost Management and support resource group filtering.
Proto design accommodates Azure's category system.

**Alternatives considered**: Ignore Azure-specific categories - rejected because categories are
core to Azure budget organization.

#### Kubecost Budget API

**Decision**: Map Kubecost's name, spendLimit, interval, and actions[].percentage to proto fields.

**Rationale**: Kubecost provides namespace-level budgeting with percentage-based actions. Proto
supports the interval-based recurring budgets.

**Alternatives considered**: Generic budget format - rejected because Kubecost has unique
namespace scoping.

### Existing CostSource RPC Patterns

**Decision**: Follow existing optional RPC pattern with Unimplemented return for unsupported
plugins.

**Rationale**: CostSource service already has optional RPCs (like GetRecommendations). This
maintains consistency and allows gradual plugin adoption.

**Alternatives considered**: Required RPC - rejected because not all providers support budgets
(breaking change).

### Proto Message Structure Optimization

**Decision**: Use separate budget.proto file for budget messages, extend costsource.proto with
GetBudgets RPC.

**Rationale**: Budget messages are substantial and reusable. Separating concerns improves
maintainability while keeping RPC definitions together.

**Alternatives considered**: Inline budget messages in costsource.proto - rejected because it
would bloat the main service file.

### Cross-Provider Field Mapping

**Decision**: Define provider-agnostic proto fields with source field for provider identification.

**Rationale**: Allows unified client experience while preserving provider-specific information
in metadata maps.

**Alternatives considered**: Provider-specific message variants - rejected because it would
fragment the API and complicate clients.

### Error Handling Patterns

**Decision**: Use existing gRPC status codes and error patterns from CostSource service.

**Rationale**: Maintains consistency with existing RPC error handling. Plugins can return
appropriate status codes for auth failures, rate limits, etc.

**Alternatives considered**: Custom budget-specific errors - rejected because existing patterns
are sufficient and well-established.

### Performance Considerations

**Decision**: Include include_status flag to allow clients to request current spend data
separately from budget definitions.

**Rationale**: Budget definitions may be cached, but current spend requires live API calls.
This optimization prevents unnecessary API calls.

**Alternatives considered**: Always include status - rejected because it would impact
performance for clients that only need budget definitions.

### Testing Approach

**Decision**: Follow existing conformance test patterns with Basic/Standard/Advanced levels.

**Rationale**: Ensures comprehensive testing coverage. Basic tests validate RPC structure,
Standard tests validate cross-provider mapping, Advanced tests validate performance with
100+ budgets.

**Alternatives considered**: Minimal testing - rejected because gRPC protocol changes require
thorough validation.

## Key Findings

1. **Proto Design**: Budget messages need to accommodate provider-specific fields while
   maintaining API consistency
2. **Performance**: include_status flag provides necessary optimization for live vs cached data
3. **Compatibility**: Optional RPC approach allows gradual adoption without breaking existing
   plugins
4. **Validation**: Cross-provider examples required for all major cloud providers
5. **Error Handling**: Follow existing CostSource patterns for consistency

## Resolved Unknowns

- ✅ Provider API compatibility confirmed
- ✅ Proto structure optimized for maintainability
- ✅ Performance optimizations identified
- ✅ Testing approach aligned with existing patterns
- ✅ Error handling patterns established

**Status**: All research complete. Ready for Phase 1 design.
