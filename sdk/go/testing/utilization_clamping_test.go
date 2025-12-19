package testing_test

import (
	"context"
	"strings"
	"testing"

	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func TestUtilizationClamping(t *testing.T) {
	mock := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(mock)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	tests := []struct {
		name     string
		val      float64
		expected string
	}{
		{
			name:     "Negative clamped to 0",
			val:      -0.5,
			expected: "util:0.00",
		},
		{
			name:     "Above 1 clamped to 1",
			val:      1.5,
			expected: "util:1.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := harness.Client().GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
				Resource:              resource,
				UtilizationPercentage: tt.val,
			})
			if err != nil {
				t.Fatalf("GetProjectedCost failed: %v", err)
			}

			if !strings.Contains(resp.GetBillingDetail(), tt.expected) {
				t.Errorf("Expected billing detail to contain %s, got %s", tt.expected, resp.GetBillingDetail())
			}
		})
	}
}
