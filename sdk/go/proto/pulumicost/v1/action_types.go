// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
