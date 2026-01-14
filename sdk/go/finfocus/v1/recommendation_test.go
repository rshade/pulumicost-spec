package finfocus_test

import (
	"testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRecommendationReason_Serialization(t *testing.T) {
	// Create a recommendation with the new fields
	rec := &pbc.Recommendation{
		Id:          "rec-123",
		Description: "Downsize instance",
		// These fields will cause compilation error until proto is updated and sdk generated
		PrimaryReason: pbc.RecommendationReason_RECOMMENDATION_REASON_OVER_PROVISIONED,
		SecondaryReasons: []pbc.RecommendationReason{
			pbc.RecommendationReason_RECOMMENDATION_REASON_IDLE,
		},
	}

	// Verify fields are set
	assert.Equal(t, pbc.RecommendationReason_RECOMMENDATION_REASON_OVER_PROVISIONED, rec.GetPrimaryReason())
	assert.Len(t, rec.GetSecondaryReasons(), 1)
	assert.Equal(t, pbc.RecommendationReason_RECOMMENDATION_REASON_IDLE, rec.GetSecondaryReasons()[0])

	// Verify serialization
	bytes, err := proto.Marshal(rec)
	require.NoError(t, err)

	// Verify deserialization
	rec2 := &pbc.Recommendation{}
	err = proto.Unmarshal(bytes, rec2)
	require.NoError(t, err)

	assert.Equal(t, rec.GetPrimaryReason(), rec2.GetPrimaryReason())
	assert.Equal(t, rec.GetSecondaryReasons(), rec2.GetSecondaryReasons())
}

func TestRecommendationReason_SwitchConsumption(t *testing.T) {
	tests := []struct {
		name     string
		reason   pbc.RecommendationReason
		expected string
	}{
		{
			name:     "Over-provisioned",
			reason:   pbc.RecommendationReason_RECOMMENDATION_REASON_OVER_PROVISIONED,
			expected: "Downsize",
		},
		{
			name:     "Under-provisioned",
			reason:   pbc.RecommendationReason_RECOMMENDATION_REASON_UNDER_PROVISIONED,
			expected: "Upsize",
		},
		{
			name:     "Idle",
			reason:   pbc.RecommendationReason_RECOMMENDATION_REASON_IDLE,
			expected: "Terminate",
		},
		{
			name:     "Obsolete",
			reason:   pbc.RecommendationReason_RECOMMENDATION_REASON_OBSOLETE_GENERATION,
			expected: "Upgrade",
		},
		{
			name:     "Redundant",
			reason:   pbc.RecommendationReason_RECOMMENDATION_REASON_REDUNDANT,
			expected: "Consolidate",
		},
		{
			name:     "Unspecified",
			reason:   pbc.RecommendationReason_RECOMMENDATION_REASON_UNSPECIFIED,
			expected: "Review",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var action string
			switch tt.reason {
			case pbc.RecommendationReason_RECOMMENDATION_REASON_OVER_PROVISIONED:
				action = "Downsize"
			case pbc.RecommendationReason_RECOMMENDATION_REASON_UNDER_PROVISIONED:
				action = "Upsize"
			case pbc.RecommendationReason_RECOMMENDATION_REASON_IDLE:
				action = "Terminate"
			case pbc.RecommendationReason_RECOMMENDATION_REASON_OBSOLETE_GENERATION:
				action = "Upgrade"
			case pbc.RecommendationReason_RECOMMENDATION_REASON_REDUNDANT:
				action = "Consolidate"
			case pbc.RecommendationReason_RECOMMENDATION_REASON_UNSPECIFIED:
				action = "Review"
			default:
				action = "Review"
			}
			assert.Equal(t, tt.expected, action)
		})
	}
}
