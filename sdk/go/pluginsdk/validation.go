// Package pluginsdk provides a development SDK for PulumiCost plugins.
//
// # Request Validation
//
// This package provides lightweight request validation for plugin implementations.
// For comprehensive contract testing with detailed error reporting, see sdk/go/testing/contract.go.
//
// Key differences:
//   - pluginsdk: Simple errors, minimal dependencies, optimized for plugin defense-in-depth
//   - testing/contract: Rich error context, comprehensive rules, designed for test suites
//
// # Usage
//
// Validate requests at plugin entry points:
//
//	func (s *MyPlugin) GetProjectedCost(ctx context.Context, req *pb.GetProjectedCostRequest) (*pb.GetProjectedCostResponse, error) {
//	    if err := pluginsdk.ValidateProjectedCostRequest(req); err != nil {
//	        return nil, status.Error(codes.InvalidArgument, err.Error())
//	    }
//	    // Process valid request...
//	}
//
// # Performance
//
// Validation is optimized for performance:
//   - Zero allocations on the happy path (valid request returns nil)
//   - Error paths allocate only for the error message
//   - Target: <100ns execution time for valid requests
package pluginsdk

import (
	"errors"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// Validation error messages for GetProjectedCostRequest.
var (
	ErrProjectedCostRequestNil        = errors.New("request is required")
	ErrProjectedCostResourceNil       = errors.New("resource is required")
	ErrProjectedCostProviderEmpty     = errors.New("resource.provider is required")
	ErrProjectedCostResourceTypeEmpty = errors.New("resource.resource_type is required")
	ErrProjectedCostSkuEmpty          = errors.New("resource.sku is required (use mapping helpers)")
	ErrProjectedCostRegionEmpty       = errors.New("resource.region is required (use mapping helpers)")
	ErrUtilizationOutOfRange          = errors.New("utilization_percentage must be between 0.0 and 1.0")
	ErrMetricKindInvalid              = errors.New("invalid metric kind")
)

// Validation error messages for GetActualCostRequest.
var (
	ErrActualCostRequestNil       = errors.New("request is required")
	ErrActualCostResourceIDEmpty  = errors.New("resource_id is required")
	ErrActualCostStartTimeNil     = errors.New("start_time is required")
	ErrActualCostEndTimeNil       = errors.New("end_time is required")
	ErrActualCostTimeRangeInvalid = errors.New(
		"end_time must be strictly after start_time (equal timestamps not allowed)",
	)
)

// ValidateProjectedCostRequest validates a GetProjectedCostRequest for required fields.
// This function is designed for use in both:
//   - Core: Pre-flight validation before sending requests to plugins
//   - Plugins: Defense-in-depth validation upon receiving requests
//
// Validation order (fail-fast, structural before field values):
//  1. Request nil check
//  2. Resource nil check
//  3. Provider empty check
//  4. ResourceType empty check
//  5. SKU empty check (with mapping helper guidance)
//  6. Region empty check (with mapping helper guidance)
//  7. Global utilization range check (if non-zero)
//  8. Resource-level utilization range check (if provided)
//
// Performance: Zero allocations on the happy path (valid request returns nil).
// Error paths allocate for the error message.
//
// Returns nil if the request is valid, or an error describing the first validation failure.
func ValidateProjectedCostRequest(req *pbc.GetProjectedCostRequest) error {
	if req == nil {
		return ErrProjectedCostRequestNil
	}

	resource := req.GetResource()
	if resource == nil {
		return ErrProjectedCostResourceNil
	}

	if len(resource.GetProvider()) == 0 {
		return ErrProjectedCostProviderEmpty
	}

	if len(resource.GetResourceType()) == 0 {
		return ErrProjectedCostResourceTypeEmpty
	}

	if len(resource.GetSku()) == 0 {
		return ErrProjectedCostSkuEmpty
	}

	if len(resource.GetRegion()) == 0 {
		return ErrProjectedCostRegionEmpty
	}

	// Validate utilization values using centralized helper
	// Global utilization: non-zero values must be valid (protobuf3 default is 0.0)
	if u := req.GetUtilizationPercentage(); u != 0 && !IsUtilizationValid(u) {
		return ErrUtilizationOutOfRange
	}

	// Resource-level utilization: if explicitly set, must be valid
	if resource.UtilizationPercentage != nil && !IsUtilizationValid(resource.GetUtilizationPercentage()) {
		return ErrUtilizationOutOfRange
	}

	return nil
}

// ValidateSupportsResponse validates a SupportsResponse for correctness.
//
// Validation order:
//  1. Response nil check
//  2. Supported metrics validity check
//
// Returns nil if the response is valid, or an error describing the failure.
func ValidateSupportsResponse(res *pbc.SupportsResponse) error {
	if res == nil {
		return errors.New("response is required")
	}

	for _, kind := range res.GetSupportedMetrics() {
		if !IsValidMetricKind(kind) {
			return ErrMetricKindInvalid
		}
	}

	return nil
}

// validMetricKinds contains all valid sustainability metric kinds for zero-allocation validation.
// This follows the pattern established in sdk/go/registry for optimized enum validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var validMetricKinds = []pbc.MetricKind{
	pbc.MetricKind_METRIC_KIND_CARBON_FOOTPRINT,
	pbc.MetricKind_METRIC_KIND_ENERGY_CONSUMPTION,
	pbc.MetricKind_METRIC_KIND_WATER_USAGE,
}

// ValidMetricKinds returns all valid sustainability metric kinds.
// The returned slice is shared and must not be modified.
func ValidMetricKinds() []pbc.MetricKind {
	return validMetricKinds
}

// IsValidMetricKind returns true if the MetricKind is a recognized sustainability metric.
// METRIC_KIND_UNSPECIFIED is not in the valid list and returns false.
//
// Performance: Zero allocations, ~5-12 ns/op for small enum sets.
func IsValidMetricKind(kind pbc.MetricKind) bool {
	for _, valid := range validMetricKinds {
		if kind == valid {
			return true
		}
	}
	return false
}

// ValidateActualCostRequest validates a GetActualCostRequest for required fields.
// This function is designed for use in both:
//   - Core: Pre-flight validation before sending requests to plugins
//   - Plugins: Defense-in-depth validation upon receiving requests
//
// Validation order (fail-fast):
//  1. Request nil check
//  2. ResourceId empty check
//  3. StartTime nil check
//  4. EndTime nil check
//  5. TimeRange validation (EndTime must be after StartTime)
//
// Performance: Zero allocations on the happy path (valid request returns nil).
// Error paths allocate for the error message.
//
// Returns nil if the request is valid, or an error describing the first validation failure.
func ValidateActualCostRequest(req *pbc.GetActualCostRequest) error {
	if req == nil {
		return ErrActualCostRequestNil
	}

	if len(req.GetResourceId()) == 0 {
		return ErrActualCostResourceIDEmpty
	}

	startTime := req.GetStart()
	if startTime == nil {
		return ErrActualCostStartTimeNil
	}

	endTime := req.GetEnd()
	if endTime == nil {
		return ErrActualCostEndTimeNil
	}

	// Compare timestamps: end must be strictly after start
	// Using AsTime() for accurate comparison including nanoseconds
	if !endTime.AsTime().After(startTime.AsTime()) {
		return ErrActualCostTimeRangeInvalid
	}

	return nil
}
