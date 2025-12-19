package testing_test

import (
	"context"
	"strings"
	"testing"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
	"google.golang.org/protobuf/proto"
)

func TestUtilizationPrecedence(t *testing.T) {
	mock := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(mock)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	tests := []struct {
		name     string
		global   float64
		override *float64
		expected string // Part of BillingDetail in mock
	}{
		{
			name:     "Neither provided uses default 0.5",
			global:   0.0,
			override: nil,
			expected: "util:0.50",
		},
		{
			name:     "Global provided uses global",
			global:   0.75,
			override: nil,
			expected: "util:0.75",
		},
		{
			name:     "Override provided uses override",
			global:   0.75,
			override: proto.Float64(0.25),
			expected: "util:0.25",
		},
		{
			name:   "Global 0.0 with no resource override uses default (proto3 limitation)",
			global: 0.0,
			// Protobuf3 uses 0.0 as default for double, so we can't distinguish
			// "explicitly 0.0" from "not set" at the global level.
			override: nil,
			expected: "util:0.50",
		},
		{
			name:     "Explicit 0.0 at resource level is honored",
			global:   0.75,
			override: proto.Float64(0.0), // Explicit zero via pointer
			expected: "util:0.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := proto.Clone(resource).(*pbc.ResourceDescriptor)
			if tt.override != nil {
				res.UtilizationPercentage = tt.override
			}

			resp, err := harness.Client().GetProjectedCost(ctx, &pbc.GetProjectedCostRequest{
				Resource:              res,
				UtilizationPercentage: tt.global,
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
