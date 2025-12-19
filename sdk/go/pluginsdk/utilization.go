package pluginsdk

import (
	"github.com/rshade/pulumicost-spec/sdk/go/internal/utilization"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// DefaultUtilization is the assumed utilization level when no utilization is
// explicitly provided at either request or resource level (50%).
const DefaultUtilization = utilization.DefaultUtilization

// GetUtilization extracts the utilization percentage from a GetProjectedCostRequest
// following standard precedence rules:
//
//  1. Resource-level override (ResourceDescriptor.utilization_percentage) - highest priority
//  2. Global request level (GetProjectedCostRequest.utilization_percentage)
//  3. SDK baseline default (0.5 or 50%) - lowest priority
//
// # Validation Requirement
//
// This function assumes the request has already been validated with
// [ValidateProjectedCostRequest]. Out-of-range values (< 0.0 or > 1.0) should be
// rejected at the validation layer, not silently corrected here.
//
// # Explicit Zero Handling
//
// A value of 0.0 is treated as an explicit request for 0% utilization, NOT as "unset".
// The default (0.5) is only applied when:
//   - Resource-level utilization_percentage is nil (not set), AND
//   - Global utilization_percentage is 0.0 (protobuf3 default for unset double)
//
// To request 0% utilization explicitly, set the resource-level field:
//
//	resource.UtilizationPercentage = proto.Float64(0.0)
//
// # Why Resource-Level for Explicit Zero?
//
// Protobuf3 uses 0.0 as the default value for double fields, making it impossible
// to distinguish between "explicitly set to 0.0" and "not set" at the global level.
// The resource-level field uses a pointer type (*float64), allowing nil to represent
// "not set" vs 0.0 representing "explicitly zero".
func GetUtilization(req *pbc.GetProjectedCostRequest) float64 {
	return utilization.Get(req)
}

// IsUtilizationValid checks if a utilization value is within the valid [0.0, 1.0] range.
// Use this for validation before processing. Returns false for NaN and Inf values.
func IsUtilizationValid(u float64) bool {
	return utilization.IsValid(u)
}
