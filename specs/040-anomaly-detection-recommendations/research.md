# Research: Anomaly Detection via Recommendations

**Feature**: 040-anomaly-detection-recommendations
**Date**: 2026-01-19
**Status**: Complete (no unknowns identified)

## Technical Context Resolution

This feature has no NEEDS CLARIFICATION items. All technical decisions are straightforward based on
existing patterns in the codebase.

### Decision 1: Enum Value Assignments

**Decision**: Use `RECOMMENDATION_CATEGORY_ANOMALY = 5` and `RECOMMENDATION_ACTION_TYPE_INVESTIGATE = 12`

**Rationale**:

- Current `RecommendationCategory` enum: 0-4 (UNSPECIFIED, COST, PERFORMANCE, SECURITY, RELIABILITY)
- Current `RecommendationActionType` enum: 0-11 (UNSPECIFIED through OTHER)
- Sequential assignment follows protobuf best practices for enum evolution
- No gaps needed since these are additive, non-breaking changes

**Alternatives Considered**:

- Using higher numbers (e.g., 100) for "extension" values: Rejected - unnecessary complexity
- Creating separate enums: Rejected - violates "unified recommendations" design goal

### Decision 2: Backward Compatibility Approach

**Decision**: Pure enum addition with no structural changes to Recommendation message

**Rationale**:

- Adding enum values to proto3 enums is fully backward compatible
- Old clients receiving ANOMALY/INVESTIGATE will see them as unknown enum values (proto3 default)
- Old servers will never send ANOMALY/INVESTIGATE, so old clients are unaffected
- buf breaking check will pass since no fields are removed/modified

**Alternatives Considered**:

- Creating new `AnomalyRecommendation` message type: Rejected per issue #315 - duplicates 80%+ of fields
- Adding `GetAnomalies` RPC: Rejected per issue #315 - fragments the API surface

### Decision 3: SDK Update Scope

**Decision**: Regenerate both Go and TypeScript SDKs; add documentation but no new helper functions

**Rationale**:

- Go SDK: `make generate` regenerates proto bindings automatically
- TypeScript SDK: buf generation regenerates bindings automatically
- No new SDK helpers needed - existing `Recommendation` creation patterns apply
- Documentation updates explain semantic mapping of existing fields to anomaly use cases

**Alternatives Considered**:

- Adding `NewAnomalyRecommendation()` helper: Deferred - may be added later if patterns emerge
- Adding validation to enforce ANOMALY + INVESTIGATE pairing: Rejected - spec allows flexibility

## Provider API Mapping

Research confirms provider anomaly APIs map cleanly to Recommendation fields:

| Provider | Anomaly API | confidence_score | estimated_savings | metadata keys |
|----------|-------------|------------------|-------------------|---------------|
| AWS | GetAnomalies | AnomalyScore (0-100 → 0.0-1.0) | Impact.TotalActualSpend | baseline, rootCauses |
| Azure | Cost Anomalies | ConfidenceScore | AnomalyAmount | threshold, contributors |
| GCP | Budget Alerts | N/A (threshold-based) | AlertSpendBasis | budgetName, thresholdRules |

**Note**: GCP lacks ML-based anomaly detection; plugins may use threshold triggers with lower confidence scores.

## Testing Strategy

**Decision**: Add conformance test for anomaly recommendation validation

**Test Cases**:

1. Plugin returns ANOMALY category → category filter works correctly
2. Plugin returns INVESTIGATE action → action_type filter works correctly
3. Negative estimated_savings is accepted (overspend anomaly)
4. Missing confidence_score is allowed (optional field)
5. ANOMALY + non-INVESTIGATE action combination is valid (per spec flexibility)

## Documentation Strategy

**Decision**: Add anomaly usage guide to SDK documentation

**Sections**:

1. Field mapping table (Recommendation field → anomaly semantic)
2. Example anomaly recommendation JSON
3. Provider-specific notes (AWS, Azure, GCP confidence score mapping)
4. Filtering best practices (category=ANOMALY, min_confidence_score)

## Conclusion

No research blockers identified. Implementation can proceed with:

1. Proto enum additions (2 values)
2. SDK regeneration (Go + TypeScript)
3. Conformance test addition
4. Documentation updates
