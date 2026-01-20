# Proto Contract Changes: Anomaly Detection

**Feature**: 040-anomaly-detection-recommendations
**File**: `proto/finfocus/v1/costsource.proto`

## Change Summary

Two enum value additions to existing enums. No new messages, fields, or RPCs.

## RecommendationCategory Enum Addition

**Location**: Lines ~867-873

```diff
 // RecommendationCategory classifies the type of optimization recommendation.
 enum RecommendationCategory {
   RECOMMENDATION_CATEGORY_UNSPECIFIED = 0;
   RECOMMENDATION_CATEGORY_COST = 1;
   RECOMMENDATION_CATEGORY_PERFORMANCE = 2;
   RECOMMENDATION_CATEGORY_SECURITY = 3;
   RECOMMENDATION_CATEGORY_RELIABILITY = 4;
+  // Cost anomaly requiring investigation.
+  // Used for unusual spending patterns detected by cost management services
+  // (AWS Cost Anomaly Detection, Azure Cost Management Anomalies, etc.).
+  // Anomaly recommendations typically use RECOMMENDATION_ACTION_TYPE_INVESTIGATE.
+  RECOMMENDATION_CATEGORY_ANOMALY = 5;
 }
```

## RecommendationActionType Enum Addition

**Location**: Lines ~875-904

```diff
 // RecommendationActionType specifies the type of action recommended.
 enum RecommendationActionType {
   RECOMMENDATION_ACTION_TYPE_UNSPECIFIED = 0;
   RECOMMENDATION_ACTION_TYPE_RIGHTSIZE = 1;
   RECOMMENDATION_ACTION_TYPE_TERMINATE = 2;
   RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT = 3;
   RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS = 4;
   RECOMMENDATION_ACTION_TYPE_MODIFY = 5;
   RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED = 6;
   // ... existing values 7-10 ...
   RECOMMENDATION_ACTION_TYPE_OTHER = 11;
+  // Investigate anomaly or issue requiring human analysis.
+  // No automated remediation is appropriate; the recommendation indicates
+  // something unusual that requires manual investigation.
+  // Commonly paired with RECOMMENDATION_CATEGORY_ANOMALY for cost anomalies.
+  RECOMMENDATION_ACTION_TYPE_INVESTIGATE = 12;
 }
```

## Breaking Change Analysis

| Check | Result | Notes |
|-------|--------|-------|
| buf breaking | PASS | Enum additions are backward compatible |
| Field removal | N/A | No fields removed |
| Field renaming | N/A | No fields renamed |
| Type changes | N/A | No type changes |
| Reserved fields | N/A | No reserved changes needed |

## Wire Format Impact

- **Existing messages**: Unchanged
- **New enum values**: Encoded as varint (5 and 12)
- **Unknown enum handling**: Proto3 preserves unknown enum values during round-trip

## SDK Regeneration Required

After applying these changes:

```bash
# Go SDK
make generate

# TypeScript SDK (if applicable)
cd sdk/typescript && npm run generate
```

## Example Proto Message

```protobuf
// Example anomaly recommendation (for documentation)
message ExampleAnomalyRecommendation {
  // Returns a recommendation like:
  // {
  //   "id": "aws-anomaly-12345",
  //   "category": "RECOMMENDATION_CATEGORY_ANOMALY",      // 5
  //   "action_type": "RECOMMENDATION_ACTION_TYPE_INVESTIGATE",  // 12
  //   "resource": {
  //     "provider": "aws",
  //     "resource_type": "ec2",
  //     "region": "us-east-1"
  //   },
  //   "impact": {
  //     "estimated_savings": -1500.00,  // Negative = overspend
  //     "currency": "USD"
  //   },
  //   "confidence_score": 0.85,
  //   "description": "Unusual EC2 spending: 150% above 30-day baseline",
  //   "metadata": {
  //     "baseline_amount": "1000.00",
  //     "deviation_percent": "150",
  //     "detection_time": "2026-01-19T10:30:00Z"
  //   }
  // }
}
```
