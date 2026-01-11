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

package testing_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	plugintesting "github.com/rshade/pulumicost-spec/sdk/go/testing"
)

// TestUtilizationPassthrough verifies that out-of-range utilization values pass through
// without modification. The SDK follows a "strict validation" approach where:
//
//  1. ValidateProjectedCostRequest rejects out-of-range values with errors
//  2. GetUtilization returns values as-is (no clamping) for valid requests
//
// This test demonstrates that if validation is bypassed, values pass through unchanged.
// In production, always validate requests before processing.
func TestUtilizationPassthrough(t *testing.T) {
	mock := plugintesting.NewMockPlugin()
	harness := plugintesting.NewTestHarness(mock)
	harness.Start(t)
	defer harness.Stop()

	ctx := context.Background()
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	tests := []struct {
		name     string
		val      float64
		expected string // Expected value in billing detail (no clamping)
	}{
		{
			name:     "Negative passes through unchanged",
			val:      -0.5,
			expected: "util:-0.50",
		},
		{
			name:     "Above 1 passes through unchanged",
			val:      1.5,
			expected: "util:1.50",
		},
		{
			name:     "Valid value passes through",
			val:      0.75,
			expected: "util:0.75",
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

// TestValidationRejectsOutOfRange verifies that ValidateProjectedCostRequest
// properly rejects out-of-range utilization values with clear error messages.
func TestValidationRejectsOutOfRange(t *testing.T) {
	resource := plugintesting.CreateResourceDescriptor("aws", "ec2", "t3.micro", "us-east-1")

	tests := []struct {
		name    string
		val     float64
		wantErr bool
	}{
		{"negative is rejected", -0.1, true},
		{"above 1 is rejected", 1.1, true},
		{"0 is valid", 0.0, false},
		{"1 is valid", 1.0, false},
		{"0.5 is valid", 0.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pbc.GetProjectedCostRequest{
				Resource:              resource,
				UtilizationPercentage: tt.val,
			}
			err := pluginsdk.ValidateProjectedCostRequest(req)
			if tt.wantErr {
				if !errors.Is(err, pluginsdk.ErrUtilizationOutOfRange) {
					t.Errorf("expected ErrUtilizationOutOfRange, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
