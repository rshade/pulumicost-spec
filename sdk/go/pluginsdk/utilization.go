package pluginsdk

import (
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

const (
	// DefaultUtilization is the assumed utilization level when none is provided (50%).
	DefaultUtilization = 0.5
)

// GetUtilization extracts the utilization percentage from a GetProjectedCostRequest
// following standard precedence rules:
//  1. Resource-level override (ResourceDescriptor.utilization_percentage)
//  2. Global request default (GetProjectedCostRequest.utilization_percentage)
//  3. SDK baseline default (0.5 or 50%)
//
// Values are automatically clamped to the [0.0, 1.0] range.
func GetUtilization(req *pbc.GetProjectedCostRequest) float64 {
	if req == nil {
		return DefaultUtilization
	}

	utilization := req.GetUtilizationPercentage()
	resourceProvided := false

	// Check resource-level override (top priority)
	if req.GetResource() != nil && req.GetResource().UtilizationPercentage != nil {
		utilization = req.GetResource().GetUtilizationPercentage()
		resourceProvided = true
	}

	// If neither was explicitly provided (heuristic for global 0.0), use default 0.5.
	// Since global double is not optional, we treat 0.0 as "unset" ONLY IF
	// the resource level is also not provided.
	if !resourceProvided && utilization == 0 {
		utilization = DefaultUtilization
	}

	return ClampUtilization(utilization)
}

// ClampUtilization ensures a utilization value is within the valid [0.0, 1.0] range.
func ClampUtilization(u float64) float64 {
	if u < 0 {
		return 0.0
	}
	if u > 1.0 {
		return 1.0
	}
	return u
}
