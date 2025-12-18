# Research: Extend RecommendationActionType Enum

**Feature**: 019-recommendation-action-types
**Date**: 2025-12-17

## Research Questions

### Q1: What action types do major FinOps platforms provide?

**Decision**: Add MIGRATE, CONSOLIDATE, SCHEDULE, REFACTOR, OTHER

**Rationale**: Comprehensive analysis of FinOps platforms identified gaps in current enum:

| Action Type | AWS | Azure Advisor | GCP Recommender | Kubecost |
|-------------|-----|---------------|-----------------|----------|
| MIGRATE | - | ✓ Move to reserved | ✓ Change machine type | - |
| CONSOLIDATE | - | ✓ Combine resources | - | ✓ Node consolidation |
| SCHEDULE | ✓ Instance Scheduler | ✓ Automation | - | - |
| REFACTOR | - | - | ✓ Use serverless | - |
| OTHER | ✓ Catch-all | ✓ Catch-all | ✓ Catch-all | ✓ Catch-all |

**Alternatives Considered**:

- **Expand MODIFY to cover all**: Rejected - loses semantic specificity for filtering/grouping
- **Add only OTHER as catch-all**: Rejected - platforms provide specific categories we should map

### Q2: Are proto3 enum extensions backward compatible?

**Decision**: Yes, additive enum extensions are fully backward compatible

**Rationale**: Proto3 specification guarantees:

1. Wire format unchanged - enums are varints on the wire
2. Unknown values preserved as numeric representation by receivers
3. Existing code continues to work without recompilation
4. New values can be added at any position (sequential preferred for readability)

**Alternatives Considered**:

- **Use reserved field numbers**: Not applicable - only needed for removed fields
- **Add deprecation markers to existing values**: Not needed - no values being deprecated

### Q3: What enum value numbers should be assigned?

**Decision**: Use sequential numbers 7-11

**Rationale**:

- Current highest value: DELETE_UNUSED = 6
- Sequential numbering (7, 8, 9, 10, 11) maintains consistency
- Proto3 best practice: reserve 1-15 for frequently used fields (irrelevant for enums)
- No gaps in numbering improves readability and tooling compatibility

**Alternatives Considered**:

- **Start at 100 for new values**: Rejected - creates confusing gap, no technical benefit
- **Group by category (10s, 20s)**: Rejected - existing values already sequential 0-6

### Q4: What documentation is required for each new value?

**Decision**: Add proto comments with:

1. Brief description (1 line)
2. Use case example
3. Source platforms where this action type originates

**Rationale**: Constitution requires inline proto comments for generated documentation.
Plugin developers need clear guidance on when to use each action type.

**Alternatives Considered**:

- **External documentation only**: Rejected - violates proto-first documentation principle
- **Minimal comments**: Rejected - insufficient for plugin developer guidance

### Q5: Are there any action-specific message types needed?

**Decision**: No new message types required

**Rationale**:

- MIGRATE, CONSOLIDATE, SCHEDULE, REFACTOR can use existing `ModifyAction` message
- `ModifyAction.modification_type` field captures the specific action
- `ModifyAction.current_config` and `recommended_config` maps handle any parameters
- OTHER requires no specific message - uses metadata map on Recommendation

**Alternatives Considered**:

- **Add MigrateAction, ConsolidateAction, etc.**: Rejected - over-engineering, ModifyAction
  suffices with modification_type differentiation
- **Add ScheduleAction with cron fields**: Rejected - schedule details vary by provider,
  better captured in metadata or recommended_config map

## Summary

This is a straightforward, additive enum extension with:

- **5 new values**: MIGRATE (7), CONSOLIDATE (8), SCHEDULE (9), REFACTOR (10), OTHER (11)
- **Zero new messages**: Existing ModifyAction and metadata fields handle all cases
- **Full backward compatibility**: Proto3 guarantees wire format stability
- **Clear documentation**: Proto comments required for each new value

No NEEDS CLARIFICATION items remain. Research phase complete.
