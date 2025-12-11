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
// Validation order (fail-fast):
//  1. Request nil check
//  2. Resource nil check
//  3. Provider empty check
//  4. ResourceType empty check
//  5. SKU empty check (with mapping helper guidance)
//  6. Region empty check (with mapping helper guidance)
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

	return nil
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
