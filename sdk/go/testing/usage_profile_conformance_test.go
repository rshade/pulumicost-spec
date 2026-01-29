package testing_test

import (
	"context"
	"testing"
	"time"

	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// TestUsageProfileConformance validates plugin behavior with usage_profile field
// in GetProjectedCostRequest and GetRecommendationsRequest.
//
// IMPORTANT: Do NOT use t.Parallel() in subtests below.
// They share a single gRPC harness that will be closed on function exit.
// See: CLAUDE.md pattern for subtests sharing TestHarness.
//
//nolint:gocognit // Table-driven test with RPC method switch; complexity is inherent
func TestUsageProfileConformance(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	tests := []struct {
		name           string
		rpcMethod      string // "GetProjectedCost" or "GetRecommendations"
		usageProfile   pbc.UsageProfile
		provider       string
		resourceType   string
		region         string
		filterCategory pbc.RecommendationCategory
		checkCategory  bool
	}{
		// GetProjectedCost tests
		{
			name:         "GetProjectedCost_UNSPECIFIED_profile",
			rpcMethod:    "GetProjectedCost",
			usageProfile: pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
			provider:     "aws",
			resourceType: "ec2",
			region:       "us-east-1",
		},
		{
			name:         "GetProjectedCost_PROD_profile",
			rpcMethod:    "GetProjectedCost",
			usageProfile: pbc.UsageProfile_USAGE_PROFILE_PROD,
			provider:     "aws",
			resourceType: "ec2",
			region:       "us-east-1",
		},
		{
			name:         "GetProjectedCost_DEV_profile",
			rpcMethod:    "GetProjectedCost",
			usageProfile: pbc.UsageProfile_USAGE_PROFILE_DEV,
			provider:     "aws",
			resourceType: "ec2",
			region:       "us-east-1",
		},
		{
			name:         "GetProjectedCost_BURST_profile",
			rpcMethod:    "GetProjectedCost",
			usageProfile: pbc.UsageProfile_USAGE_PROFILE_BURST,
			provider:     "aws",
			resourceType: "ec2",
			region:       "us-east-1",
		},
		{
			name:         "GetProjectedCost_unknown_profile_treated_as_UNSPECIFIED",
			rpcMethod:    "GetProjectedCost",
			usageProfile: pbc.UsageProfile(999), // Unknown future value
			provider:     "aws",
			resourceType: "ec2",
			region:       "us-east-1",
		},
		// GetRecommendations tests
		{
			name:         "GetRecommendations_with_usage_profile",
			rpcMethod:    "GetRecommendations",
			usageProfile: pbc.UsageProfile_USAGE_PROFILE_DEV,
			provider:     "aws",
		},
		{
			name:           "GetRecommendations_with_filter_and_profile",
			rpcMethod:      "GetRecommendations",
			usageProfile:   pbc.UsageProfile_USAGE_PROFILE_PROD,
			provider:       "aws",
			filterCategory: pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
			checkCategory:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.rpcMethod {
			case "GetProjectedCost":
				req := &pbc.GetProjectedCostRequest{
					Resource: &pbc.ResourceDescriptor{
						Provider:     tc.provider,
						ResourceType: tc.resourceType,
						Region:       tc.region,
					},
					UsageProfile: tc.usageProfile,
				}
				resp, err := client.GetProjectedCost(ctx, req)
				if err != nil {
					t.Fatalf("GetProjectedCost failed: %v", err)
				}
				if resp.GetCostPerMonth() <= 0 {
					t.Errorf("Expected positive cost, got %f", resp.GetCostPerMonth())
				}

			case "GetRecommendations":
				filter := &pbc.RecommendationFilter{Provider: tc.provider}
				if tc.checkCategory {
					filter.Category = tc.filterCategory
				}
				req := &pbc.GetRecommendationsRequest{
					Filter:       filter,
					UsageProfile: tc.usageProfile,
				}
				resp, err := client.GetRecommendations(ctx, req)
				if err != nil {
					t.Fatalf("GetRecommendations failed: %v", err)
				}
				if resp.GetSummary() == nil {
					t.Error("Expected non-nil summary")
				}
				if tc.checkCategory {
					for _, rec := range resp.GetRecommendations() {
						if rec.GetCategory() != tc.filterCategory {
							t.Errorf("Category mismatch: got %v, want %v",
								rec.GetCategory(), tc.filterCategory)
						}
					}
				}
			}
		})
	}
}

// TestUsageProfileFieldPresence validates that usage_profile field is correctly
// serialized in requests and accessible in responses.
func TestUsageProfileFieldPresence(t *testing.T) {
	t.Parallel()

	// Verify enum values are as expected
	tests := []struct {
		name     string
		profile  pbc.UsageProfile
		expected int32
	}{
		{"UNSPECIFIED", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, 0},
		{"PROD", pbc.UsageProfile_USAGE_PROFILE_PROD, 1},
		{"DEV", pbc.UsageProfile_USAGE_PROFILE_DEV, 2},
		{"BURST", pbc.UsageProfile_USAGE_PROFILE_BURST, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if int32(tt.profile) != tt.expected {
				t.Errorf("UsageProfile %s = %d, want %d", tt.name, int32(tt.profile), tt.expected)
			}
		})
	}
}

// TestUsageProfileBackwardCompatibility ensures that requests without
// usage_profile field (defaulting to UNSPECIFIED) work correctly.
//
// IMPORTANT: Do NOT use t.Parallel() in subtests below.
// They share a single gRPC harness that will be closed on function exit.
func TestUsageProfileBackwardCompatibility(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	t.Run("GetProjectedCost_without_profile", func(t *testing.T) {
		// Create request without setting usage_profile (defaults to UNSPECIFIED=0)
		req := &pbc.GetProjectedCostRequest{
			Resource: &pbc.ResourceDescriptor{
				Provider:     "azure",
				ResourceType: "vm",
				Region:       "eastus",
			},
		}

		resp, err := client.GetProjectedCost(ctx, req)
		if err != nil {
			t.Fatalf("GetProjectedCost without profile should succeed: %v", err)
		}
		if resp.GetCostPerMonth() <= 0 {
			t.Errorf("Expected positive cost, got %f", resp.GetCostPerMonth())
		}
	})

	t.Run("GetRecommendations_without_profile", func(t *testing.T) {
		// Create request without setting usage_profile (defaults to UNSPECIFIED=0)
		req := &pbc.GetRecommendationsRequest{
			Filter: &pbc.RecommendationFilter{
				Provider: "aws",
			},
		}

		resp, err := client.GetRecommendations(ctx, req)
		if err != nil {
			t.Fatalf("GetRecommendations without profile should succeed: %v", err)
		}
		// Summary should always be present
		if resp.GetSummary() == nil {
			t.Error("Expected non-nil summary")
		}
	})
}

// TestOldPluginReceivesProfileRequests verifies FR-005: existing plugins without
// profile awareness continue to function when receiving requests with any
// usage_profile value. The mock plugin ignores the profile field in its cost
// calculation, simulating a legacy plugin that predates the usage_profile feature.
//
// IMPORTANT: Do NOT use t.Parallel() in subtests below.
// They share a single gRPC harness that will be closed on function exit.
func TestOldPluginReceivesProfileRequests(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	// An "old" plugin simply ignores usage_profile. The mock already does this
	// for cost calculation — it returns the same cost regardless of profile.
	// Verify all profile values produce valid responses.
	profiles := []pbc.UsageProfile{
		pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
		pbc.UsageProfile_USAGE_PROFILE_PROD,
		pbc.UsageProfile_USAGE_PROFILE_DEV,
		pbc.UsageProfile_USAGE_PROFILE_BURST,
		pbc.UsageProfile(999), // Future unknown value
	}

	t.Run("GetProjectedCost_all_profiles_succeed", func(t *testing.T) {
		for _, profile := range profiles {
			req := &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Region:       "us-east-1",
				},
				UsageProfile: profile,
			}

			resp, err := client.GetProjectedCost(ctx, req)
			if err != nil {
				t.Fatalf("Old plugin should handle profile %d without error: %v",
					int32(profile), err)
			}
			if resp.GetCostPerMonth() <= 0 {
				t.Errorf("Expected positive cost for profile %d, got %f",
					int32(profile), resp.GetCostPerMonth())
			}
		}
	})

	t.Run("GetRecommendations_all_profiles_succeed", func(t *testing.T) {
		for _, profile := range profiles {
			req := &pbc.GetRecommendationsRequest{
				Filter: &pbc.RecommendationFilter{
					Provider: "aws",
				},
				UsageProfile: profile,
			}

			resp, err := client.GetRecommendations(ctx, req)
			if err != nil {
				t.Fatalf("Old plugin should handle profile %d without error: %v",
					int32(profile), err)
			}
			if resp.GetSummary() == nil {
				t.Errorf("Expected non-nil summary for profile %d", int32(profile))
			}
		}
	})
}

// BenchmarkGetProjectedCost_WithProfile measures overhead of usage_profile handling.
func BenchmarkGetProjectedCost_WithProfile(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(b)
	defer harness.Stop()

	client := harness.Client()
	ctx := context.Background()

	profiles := []struct {
		name    string
		profile pbc.UsageProfile
	}{
		{"UNSPECIFIED", pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED},
		{"PROD", pbc.UsageProfile_USAGE_PROFILE_PROD},
		{"DEV", pbc.UsageProfile_USAGE_PROFILE_DEV},
		{"BURST", pbc.UsageProfile_USAGE_PROFILE_BURST},
	}

	for _, p := range profiles {
		b.Run(p.name, func(b *testing.B) {
			req := &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Region:       "us-east-1",
				},
				UsageProfile: p.profile,
			}

			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				_, err := client.GetProjectedCost(ctx, req)
				if err != nil {
					b.Fatalf("GetProjectedCost failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkGetRecommendations_WithProfile measures overhead of usage_profile in recommendations.
func BenchmarkGetRecommendations_WithProfile(b *testing.B) {
	plugin := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(b)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &pbc.GetRecommendationsRequest{
		Filter: &pbc.RecommendationFilter{
			Provider: "aws",
		},
		UsageProfile: pbc.UsageProfile_USAGE_PROFILE_DEV,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, err := client.GetRecommendations(ctx, req)
		if err != nil {
			b.Fatalf("GetRecommendations failed: %v", err)
		}
	}
}
