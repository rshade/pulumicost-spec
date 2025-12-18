package pbc

// allRecommendationActionTypes is a package-level slice containing all valid RecommendationActionType values.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allRecommendationActionTypes = []RecommendationActionType{
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_ADJUST_REQUESTS,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_DELETE_UNUSED,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MIGRATE,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_CONSOLIDATE,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_SCHEDULE,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_REFACTOR,
	RecommendationActionType_RECOMMENDATION_ACTION_TYPE_OTHER,
}

// AllRecommendationActionTypes returns a slice of all supported recommendation action types.
func AllRecommendationActionTypes() []RecommendationActionType {
	return allRecommendationActionTypes
}

// IsValidRecommendationActionType checks if the given RecommendationActionType is a valid, defined value.
// This implementation provides zero-allocation validation.
func IsValidRecommendationActionType(at RecommendationActionType) bool {
	for _, valid := range allRecommendationActionTypes {
		if at == valid {
			return true
		}
	}
	return false
}
