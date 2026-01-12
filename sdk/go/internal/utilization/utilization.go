// Package utilization provides shared utilization extraction logic for the PulumiCost SDK.
//
// This package exists to break circular dependencies between sdk/go/pluginsdk and
// sdk/go/testing. Both packages need utilization extraction logic, but pluginsdk
// imports testing for conformance test wrappers.
//
// # Usage
//
// This is an internal package. External consumers should use:
//   - [github.com/rshade/finfocus-spec/sdk/go/pluginsdk.GetUtilization]
package utilization

import (
	"math"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

const (
	// DefaultUtilization is the assumed utilization level when no utilization is
	// explicitly provided at either request or resource level (50%).
	DefaultUtilization = 0.5

	// MinUtilization is the minimum valid utilization percentage (0%).
	MinUtilization = 0.0

	// MaxUtilization is the maximum valid utilization percentage (100%).
	MaxUtilization = 1.0
)

// Get extracts the utilization percentage from a GetProjectedCostRequest
// following standard precedence rules:
//
//  1. Resource-level override (ResourceDescriptor.utilization_percentage) - highest priority
//  2. Global request level (GetProjectedCostRequest.utilization_percentage)
//  3. SDK baseline default (0.5 or 50%) - lowest priority
//
// # Validation Requirement
//
// This function assumes the request has already been validated.
// Out-of-range values (< 0.0 or > 1.0) should be rejected at the validation layer.
//
// # Explicit Zero Handling
//
// The default (0.5) is only applied when:
//   - Resource-level utilization_percentage is nil (not set), AND
//   - Global utilization_percentage is 0.0 (protobuf3 default for unset double)
//
// To request 0% utilization explicitly, set the resource-level field:
//
//	resource.UtilizationPercentage = proto.Float64(0.0)
func Get(req *pbc.GetProjectedCostRequest) float64 {
	if req == nil {
		return DefaultUtilization
	}

	// Check resource-level override first (highest priority, supports explicit 0.0)
	if req.GetResource() != nil && req.GetResource().UtilizationPercentage != nil {
		return req.GetResource().GetUtilizationPercentage()
	}

	// Use global request level if non-zero
	if req.GetUtilizationPercentage() != 0 {
		return req.GetUtilizationPercentage()
	}

	// Global is 0.0 and resource-level not provided - use SDK default
	return DefaultUtilization
}

// IsValid checks if a utilization value is within the valid [0.0, 1.0] range.
// Returns false for NaN and Inf values.
//
// The explicit NaN/Inf checks improve code clarity. While range comparisons
// naturally reject NaN (all NaN comparisons return false) and Inf falls outside
// [0.0, 1.0], explicit checks make the intent clear and prevent subtle bugs
// if the comparison logic ever changes.
func IsValid(u float64) bool {
	return !math.IsNaN(u) && !math.IsInf(u, 0) && u >= MinUtilization && u <= MaxUtilization
}
