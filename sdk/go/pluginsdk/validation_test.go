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

package pluginsdk_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateProjectedCostRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *pbc.GetProjectedCostRequest
		wantErr error
	}{
		{
			name:    "nil request returns error",
			req:     nil,
			wantErr: pluginsdk.ErrProjectedCostRequestNil,
		},
		{
			name:    "nil resource returns error",
			req:     &pbc.GetProjectedCostRequest{Resource: nil},
			wantErr: pluginsdk.ErrProjectedCostResourceNil,
		},
		{
			name:    "unset resource field returns error",
			req:     &pbc.GetProjectedCostRequest{}, // Resource field not set at all
			wantErr: pluginsdk.ErrProjectedCostResourceNil,
		},
		{
			name: "empty provider returns error",
			req: &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "",
					ResourceType: "ec2",
					Sku:          "t3.micro",
					Region:       "us-east-1",
				},
			},
			wantErr: pluginsdk.ErrProjectedCostProviderEmpty,
		},
		{
			name: "empty resource_type returns error",
			req: &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "",
					Sku:          "t3.micro",
					Region:       "us-east-1",
				},
			},
			wantErr: pluginsdk.ErrProjectedCostResourceTypeEmpty,
		},
		{
			name: "empty sku returns error with mapping guidance",
			req: &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Sku:          "",
					Region:       "us-east-1",
				},
			},
			wantErr: pluginsdk.ErrProjectedCostSkuEmpty,
		},
		{
			name: "empty region returns error with mapping guidance",
			req: &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Sku:          "t3.micro",
					Region:       "",
				},
			},
			wantErr: pluginsdk.ErrProjectedCostRegionEmpty,
		},
		{
			name: "valid request returns nil",
			req: &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Sku:          "t3.micro",
					Region:       "us-east-1",
				},
			},
			wantErr: nil,
		},
		{
			name: "utilization too high returns error",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: 1.1,
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Sku:          "t3.micro",
					Region:       "us-east-1",
				},
			},
			wantErr: pluginsdk.ErrUtilizationOutOfRange,
		},
		{
			name: "utilization too low returns error",
			req: &pbc.GetProjectedCostRequest{
				UtilizationPercentage: -0.1,
				Resource: &pbc.ResourceDescriptor{
					Provider:     "aws",
					ResourceType: "ec2",
					Sku:          "t3.micro",
					Region:       "us-east-1",
				},
			},
			wantErr: pluginsdk.ErrUtilizationOutOfRange,
		},
		{
			name: "resource utilization override too high returns error",
			req: &pbc.GetProjectedCostRequest{
				Resource: &pbc.ResourceDescriptor{
					Provider:              "aws",
					ResourceType:          "ec2",
					Sku:                   "t3.micro",
					Region:                "us-east-1",
					UtilizationPercentage: proto.Float64(1.1),
				},
			},
			wantErr: pluginsdk.ErrUtilizationOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pluginsdk.ValidateProjectedCostRequest(tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateProjectedCostRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateProjectedCostRequest_ErrorMessages(t *testing.T) {
	t.Run("sku error contains mapping guidance", func(t *testing.T) {
		req := &pbc.GetProjectedCostRequest{
			Resource: &pbc.ResourceDescriptor{
				Provider:     "aws",
				ResourceType: "ec2",
				Sku:          "",
				Region:       "us-east-1",
			},
		}
		err := pluginsdk.ValidateProjectedCostRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errMsg := err.Error()
		if !strings.Contains(errMsg, "mapping helpers") {
			t.Errorf("error message should contain mapping guidance, got: %s", errMsg)
		}
	})

	t.Run("region error contains mapping guidance", func(t *testing.T) {
		req := &pbc.GetProjectedCostRequest{
			Resource: &pbc.ResourceDescriptor{
				Provider:     "aws",
				ResourceType: "ec2",
				Sku:          "t3.micro",
				Region:       "",
			},
		}
		err := pluginsdk.ValidateProjectedCostRequest(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		errMsg := err.Error()
		if !strings.Contains(errMsg, "mapping helpers") {
			t.Errorf("error message should contain mapping guidance, got: %s", errMsg)
		}
	})
}

func TestValidateActualCostRequest(t *testing.T) {
	now := time.Now()
	startTime := timestamppb.New(now.Add(-time.Hour))
	endTime := timestamppb.New(now)

	tests := []struct {
		name    string
		req     *pbc.GetActualCostRequest
		wantErr error
	}{
		{
			name:    "nil request returns error",
			req:     nil,
			wantErr: pluginsdk.ErrActualCostRequestNil,
		},
		{
			name: "empty resource_id returns error",
			req: &pbc.GetActualCostRequest{
				ResourceId: "",
				Start:      startTime,
				End:        endTime,
			},
			wantErr: pluginsdk.ErrActualCostResourceIDEmpty,
		},
		{
			name: "nil start_time returns error",
			req: &pbc.GetActualCostRequest{
				ResourceId: "i-abc123",
				Start:      nil,
				End:        endTime,
			},
			wantErr: pluginsdk.ErrActualCostStartTimeNil,
		},
		{
			name: "nil end_time returns error",
			req: &pbc.GetActualCostRequest{
				ResourceId: "i-abc123",
				Start:      startTime,
				End:        nil,
			},
			wantErr: pluginsdk.ErrActualCostEndTimeNil,
		},
		{
			name: "end_time before start_time returns error",
			req: &pbc.GetActualCostRequest{
				ResourceId: "i-abc123",
				Start:      endTime,   // later time
				End:        startTime, // earlier time
			},
			wantErr: pluginsdk.ErrActualCostTimeRangeInvalid,
		},
		{
			name: "end_time equal to start_time returns error",
			req: &pbc.GetActualCostRequest{
				ResourceId: "i-abc123",
				Start:      startTime,
				End:        startTime, // same time
			},
			wantErr: pluginsdk.ErrActualCostTimeRangeInvalid,
		},
		{
			name: "valid request returns nil",
			req: &pbc.GetActualCostRequest{
				ResourceId: "i-abc123",
				Start:      startTime,
				End:        endTime,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pluginsdk.ValidateActualCostRequest(tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateActualCostRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSupportsResponse(t *testing.T) {
	tests := []struct {
		name    string
		res     *pbc.SupportsResponse
		wantErr error
	}{
		{
			name:    "nil response returns error",
			res:     nil,
			wantErr: pluginsdk.ErrSupportsResponseNil,
		},
		{
			name: "valid response with metrics returns nil",
			res: &pbc.SupportsResponse{
				Supported: true,
				SupportedMetrics: []pbc.MetricKind{
					pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT,
				},
			},
			wantErr: nil,
		},
		{
			name: "METRIC_KIND_UNSPECIFIED returns error",
			res: &pbc.SupportsResponse{
				Supported: true,
				SupportedMetrics: []pbc.MetricKind{
					pbc.MetricKind_METRIC_KIND_UNSPECIFIED,
				},
			},
			wantErr: pluginsdk.ErrMetricKindInvalid,
		},
		{
			name: "invalid metric kind returns error",
			res: &pbc.SupportsResponse{
				Supported: true,
				SupportedMetrics: []pbc.MetricKind{
					999,
				},
			},
			wantErr: pluginsdk.ErrMetricKindInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pluginsdk.ValidateSupportsResponse(tt.res)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateSupportsResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateActualCostRequest_TimeRangeEdgeCases(t *testing.T) {
	t.Run("one nanosecond difference is valid", func(t *testing.T) {
		now := time.Now()
		req := &pbc.GetActualCostRequest{
			ResourceId: "i-abc123",
			Start:      timestamppb.New(now),
			End:        timestamppb.New(now.Add(time.Nanosecond)),
		}
		err := pluginsdk.ValidateActualCostRequest(req)
		if err != nil {
			t.Errorf("expected nil error for 1ns difference, got: %v", err)
		}
	})

	t.Run("one nanosecond difference (end before start) returns error", func(t *testing.T) {
		now := time.Now()
		req := &pbc.GetActualCostRequest{
			ResourceId: "i-abc123",
			Start:      timestamppb.New(now),
			End:        timestamppb.New(now.Add(-time.Nanosecond)), // End time is 1ns before Start time
		}
		err := pluginsdk.ValidateActualCostRequest(req)
		if !errors.Is(err, pluginsdk.ErrActualCostTimeRangeInvalid) {
			t.Errorf("expected ErrActualCostTimeRangeInvalid for 1ns end before start, got: %v", err)
		}
	})

	t.Run("large time range is valid", func(t *testing.T) {
		now := time.Now()
		req := &pbc.GetActualCostRequest{
			ResourceId: "i-abc123",
			Start:      timestamppb.New(now.Add(-365 * 24 * time.Hour)), // 1 year ago
			End:        timestamppb.New(now),
		}
		err := pluginsdk.ValidateActualCostRequest(req)
		if err != nil {
			t.Errorf("expected nil error for large time range, got: %v", err)
		}
	})
}

// Benchmarks for validation functions.
func BenchmarkValidateProjectedCostRequest_Valid(b *testing.B) {
	req := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "aws",
			ResourceType: "ec2",
			Sku:          "t3.micro",
			Region:       "us-east-1",
		},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.ValidateProjectedCostRequest(req)
	}
}

func BenchmarkValidateProjectedCostRequest_Invalid_EmptyProvider(b *testing.B) {
	req := &pbc.GetProjectedCostRequest{
		Resource: &pbc.ResourceDescriptor{
			Provider:     "",
			ResourceType: "ec2",
			Sku:          "t3.micro",
			Region:       "us-east-1",
		},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.ValidateProjectedCostRequest(req)
	}
}

func BenchmarkValidateActualCostRequest_Valid(b *testing.B) {
	now := time.Now()
	req := &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      timestamppb.New(now.Add(-time.Hour)),
		End:        timestamppb.New(now),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.ValidateActualCostRequest(req)
	}
}

func BenchmarkValidateActualCostRequest_Invalid_EmptyResourceID(b *testing.B) {
	now := time.Now()
	req := &pbc.GetActualCostRequest{
		ResourceId: "",
		Start:      timestamppb.New(now.Add(-time.Hour)),
		End:        timestamppb.New(now),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.ValidateActualCostRequest(req)
	}
}

// BenchmarkIsValidMetricKind validates the zero-allocation claim for metric kind validation.
// This benchmark verifies that the validMetricKinds package-level slice optimization
// achieves the documented performance: ~5-12 ns/op with 0 allocs/op.
func BenchmarkIsValidMetricKind(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.IsValidMetricKind(pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT)
	}
}

// BenchmarkIsValidMetricKind_Invalid benchmarks validation of invalid metric kinds.
func BenchmarkIsValidMetricKind_Invalid(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.IsValidMetricKind(pbc.MetricKind_METRIC_KIND_UNSPECIFIED)
	}
}

// BenchmarkIsUtilizationValid validates the zero-allocation claim for utilization validation.
func BenchmarkIsUtilizationValid(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.IsUtilizationValid(0.75)
	}
}

// BenchmarkIsUtilizationValid_Invalid benchmarks validation of invalid utilization values.
func BenchmarkIsUtilizationValid_Invalid(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		_ = pluginsdk.IsUtilizationValid(1.5)
	}
}
