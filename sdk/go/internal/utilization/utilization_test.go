package utilization_test

import (
	"math"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/internal/utilization"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/proto"
)

func TestGet_Precedence(t *testing.T) {
	tests := []struct {
		name     string
		req      *pbc.GetProjectedCostRequest
		expected float64
	}{
		{
			name:     "nil request returns default",
			req:      nil,
			expected: utilization.DefaultUtilization,
		},
		{
			name:     "empty request returns default",
			req:      &pbc.GetProjectedCostRequest{},
			expected: utilization.DefaultUtilization,
		},
		{
			name: "global 0.0 with no resource returns default (proto3 limitation)",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 0.0,
			},
			expected: utilization.DefaultUtilization,
		},
		{
			name: "non-zero global is used",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 0.75,
			},
			expected: 0.75,
		},
		{
			name: "resource-level overrides global",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 0.75,
				Resource: &pbc.ResourceDescriptor{
					UtilizationPercentage: proto.Float64(0.25),
				},
			},
			expected: 0.25,
		},
		{
			name: "resource-level explicit 0.0 is honored",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 0.75,
				Resource: &pbc.ResourceDescriptor{
					UtilizationPercentage: proto.Float64(0.0),
				},
			},
			expected: 0.0,
		},
		{
			name: "resource-level nil with non-zero global uses global",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 0.80,
				Resource: &pbc.ResourceDescriptor{
					Provider: "aws", // Resource exists but no utilization override
				},
			},
			expected: 0.80,
		},
		{
			name: "resource with nil utilization falls back to global",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 0.60,
				Resource: &pbc.ResourceDescriptor{
					UtilizationPercentage: nil, // Explicitly nil
				},
			},
			expected: 0.60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilization.Get(tt.req)
			if got != tt.expected {
				t.Errorf("Get() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected bool
	}{
		{"zero is valid", 0.0, true},
		{"one is valid", 1.0, true},
		{"0.5 is valid", 0.5, true},
		{"0.001 is valid", 0.001, true},
		{"0.999 is valid", 0.999, true},
		{"negative is invalid", -0.1, false},
		{"greater than 1 is invalid", 1.1, false},
		{"large negative is invalid", -100.0, false},
		{"large positive is invalid", 100.0, false},
		{"NaN is invalid", math.NaN(), false},
		{"positive infinity is invalid", math.Inf(1), false},
		{"negative infinity is invalid", math.Inf(-1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilization.IsValid(tt.value)
			if got != tt.expected {
				t.Errorf("IsValid(%v) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func BenchmarkGet_WithResourceOverride(b *testing.B) {
	req := &pbc.GetProjectedCostRequest{
		UtilizationPercentage: 0.75,
		Resource: &pbc.ResourceDescriptor{
			UtilizationPercentage: proto.Float64(0.25),
		},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		_ = utilization.Get(req)
	}
}

func BenchmarkGet_GlobalOnly(b *testing.B) {
	req := &pbc.GetProjectedCostRequest{
		UtilizationPercentage: 0.75,
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		_ = utilization.Get(req)
	}
}

func BenchmarkIsValid(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = utilization.IsValid(0.5)
	}
}
