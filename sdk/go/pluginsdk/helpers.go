package pluginsdk

import (
	"context"
	"errors"
	"fmt"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

// HoursPerMonth is the standard number of hours used for monthly cost calculations.
// This value (730) represents the average number of hours in a month (365.25 days / 12 months * 24 hours).
const HoursPerMonth = 730.0

// ResourceMatcher helps plugins determine if they support a resource.
//
// Thread Safety: ResourceMatcher is NOT safe for concurrent use. All calls to
// AddProvider and AddResourceType must complete before the plugin begins serving
// gRPC requests. Typical usage is to configure the matcher during plugin
// initialization, before calling Serve().
type ResourceMatcher struct {
	supportedProviders map[string]bool
	supportedTypes     map[string]bool
}

// NewResourceMatcher creates a ResourceMatcher with initialized empty maps for supported providers and supported resource types.
func NewResourceMatcher() *ResourceMatcher {
	return &ResourceMatcher{
		supportedProviders: make(map[string]bool),
		supportedTypes:     make(map[string]bool),
	}
}

// AddProvider adds a supported provider (e.g., "aws", "azure", "gcp").
// Empty strings are ignored.
func (rm *ResourceMatcher) AddProvider(provider string) {
	if provider == "" {
		return
	}
	rm.supportedProviders[provider] = true
}

// AddResourceType adds a supported resource type (e.g., "aws:ec2:Instance").
// Empty strings are ignored.
func (rm *ResourceMatcher) AddResourceType(resourceType string) {
	if resourceType == "" {
		return
	}
	rm.supportedTypes[resourceType] = true
}

// Supports checks if a resource is supported by this plugin.
func (rm *ResourceMatcher) Supports(resource *pbc.ResourceDescriptor) bool {
	if rm == nil || resource == nil {
		return false
	}

	if len(rm.supportedProviders) > 0 {
		if !rm.supportedProviders[resource.GetProvider()] {
			return false
		}
	}

	if len(rm.supportedTypes) > 0 {
		if !rm.supportedTypes[resource.GetResourceType()] {
			return false
		}
	}

	return true
}

// CostCalculator provides utilities for cost calculations.
type CostCalculator struct{}

// NewCostCalculator returns a new CostCalculator for performing cost conversions and creating cost responses.
func NewCostCalculator() *CostCalculator {
	return &CostCalculator{}
}

// HourlyToMonthly converts hourly cost to monthly cost (730 hours).
func (cc *CostCalculator) HourlyToMonthly(hourlyCost float64) float64 {
	return hourlyCost * HoursPerMonth
}

// MonthlyToHourly converts monthly cost to hourly cost (730 hours).
func (cc *CostCalculator) MonthlyToHourly(monthlyCost float64) float64 {
	return monthlyCost / HoursPerMonth
}

// CreateProjectedCostResponse creates a standard projected cost response.
// unitPrice is expected to be an hourly rate; CostPerMonth is derived using 730 hours.
func (cc *CostCalculator) CreateProjectedCostResponse(
	currency string,
	unitPrice float64,
	billingDetail string,
) *pbc.GetProjectedCostResponse {
	return &pbc.GetProjectedCostResponse{
		Currency:      currency,
		UnitPrice:     unitPrice,
		CostPerMonth:  cc.HourlyToMonthly(unitPrice),
		BillingDetail: billingDetail,
	}
}

// CreateActualCostResponse creates a standard actual cost response.
//
// The FallbackHint is not set, which means it defaults to FALLBACK_HINT_UNSPECIFIED (0).
// This default value is treated as "no fallback needed" by the core system, ensuring
// backwards compatibility with existing plugins that do not explicitly set a hint.
//
// For explicit control over the fallback hint, use NewActualCostResponse with
// the WithFallbackHint option instead.
func (cc *CostCalculator) CreateActualCostResponse(
	results []*pbc.ActualCostResult,
) *pbc.GetActualCostResponse {
	return &pbc.GetActualCostResponse{
		Results: results,
	}
}

// ActualCostResponseOption is a functional option for configuring GetActualCostResponse.
// Use these options with NewActualCostResponse to create responses with explicit
// fallback hints and results.
type ActualCostResponseOption func(*pbc.GetActualCostResponse)

// WithFallbackHint sets the fallback hint on the response.
//
// Use this option to explicitly signal to the core system whether it should
// attempt to query other plugins for the requested resource:
//
//   - FALLBACK_HINT_UNSPECIFIED (0): Default. Treated as "no fallback needed".
//   - FALLBACK_HINT_NONE (1): Plugin has data; do not attempt fallback.
//   - FALLBACK_HINT_RECOMMENDED (2): Plugin has no data; core SHOULD try other plugins.
//   - FALLBACK_HINT_REQUIRED (3): Plugin cannot handle request; core MUST try fallback.
//
// Example - Plugin has data, no fallback needed:
//
//	resp := pluginsdk.NewActualCostResponse(
//	    pluginsdk.WithResults(results),
//	    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
//	)
//
// Example - Plugin has no data, recommend fallback:
//
//	resp := pluginsdk.NewActualCostResponse(
//	    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
//	)
//
// Example - Plugin cannot handle this resource type:
//
//	resp := pluginsdk.NewActualCostResponse(
//	    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_REQUIRED),
//	)
//
// Important: For actual errors (API failures, network timeouts, invalid credentials),
// return a gRPC error instead of using a fallback hint. Hints are for "no data"
// scenarios, not system failures.
func WithFallbackHint(hint pbc.FallbackHint) ActualCostResponseOption {
	return func(resp *pbc.GetActualCostResponse) {
		resp.FallbackHint = hint
	}
}

// WithResults sets the cost results on the response.
//
// Use this option to include actual cost data in the response.
// If the results slice is non-empty, you typically want to use
// FALLBACK_HINT_NONE to indicate the plugin has authoritative data.
//
// Zero-cost results vs no data:
//   - Empty results [] with FALLBACK_HINT_RECOMMENDED = "no billing data found"
//   - Results with cost: 0.00 and FALLBACK_HINT_NONE = "resource is free tier"
func WithResults(results []*pbc.ActualCostResult) ActualCostResponseOption {
	return func(resp *pbc.GetActualCostResponse) {
		resp.Results = results
	}
}

// NewActualCostResponse creates a GetActualCostResponse using functional options.
//
// This is the preferred way to create responses when you need to explicitly
// set the fallback hint. The response starts with default values (no results,
// FALLBACK_HINT_UNSPECIFIED) and options are applied in order.
//
// Example - Plugin found cost data:
//
//	resp := pluginsdk.NewActualCostResponse(
//	    pluginsdk.WithResults(results),
//	    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_NONE),
//	)
//
// Example - Plugin has no data, recommend trying other plugins:
//
//	resp := pluginsdk.NewActualCostResponse(
//	    pluginsdk.WithFallbackHint(pbc.FallbackHint_FALLBACK_HINT_RECOMMENDED),
//	)
func NewActualCostResponse(opts ...ActualCostResponseOption) *pbc.GetActualCostResponse {
	resp := &pbc.GetActualCostResponse{}
	for _, opt := range opts {
		opt(resp)
	}
	return resp
}

// NotSupportedError returns an error indicating the specified resource type and provider are not supported.
// The formatted message includes the resource's ResourceType and Provider.
func NotSupportedError(resource *pbc.ResourceDescriptor) error {
	return fmt.Errorf(
		"resource type %s from provider %s is not supported",
		resource.GetResourceType(),
		resource.GetProvider(),
	)
}

// NoDataError returns a standard error when no cost data is available.
func NoDataError(resourceID string) error {
	return fmt.Errorf("no cost data available for resource %s", resourceID)
}

// BasePlugin provides common functionality for plugin implementations.
type BasePlugin struct {
	name    string
	matcher *ResourceMatcher
	calc    *CostCalculator
}

// NewBasePlugin creates a new BasePlugin with the given name and initializes its
// ResourceMatcher and CostCalculator for shared plugin functionality.
func NewBasePlugin(name string) *BasePlugin {
	return &BasePlugin{
		name:    name,
		matcher: NewResourceMatcher(),
		calc:    NewCostCalculator(),
	}
}

// Name returns the plugin name.
func (bp *BasePlugin) Name() string {
	return bp.name
}

// Matcher returns the resource matcher.
func (bp *BasePlugin) Matcher() *ResourceMatcher {
	return bp.matcher
}

// Calculator returns the cost calculator.
func (bp *BasePlugin) Calculator() *CostCalculator {
	return bp.calc
}

// GetProjectedCost provides a default implementation that returns not supported.
func (bp *BasePlugin) GetProjectedCost(
	_ context.Context,
	req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	if req == nil {
		return nil, errors.New("GetProjectedCostRequest cannot be nil")
	}
	resource := req.GetResource()
	if resource == nil {
		return nil, errors.New("resource cannot be nil")
	}
	return nil, NotSupportedError(resource)
}

// GetActualCost provides a default implementation that returns no data.
func (bp *BasePlugin) GetActualCost(
	_ context.Context,
	req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	if req == nil {
		return nil, errors.New("GetActualCostRequest cannot be nil")
	}
	return nil, NoDataError(req.GetResourceId())
}

// GetPricingSpec provides a default implementation that returns not implemented.
// Override this method in your plugin to return pricing specifications.
func (bp *BasePlugin) GetPricingSpec(
	_ context.Context,
	req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	if req == nil {
		return nil, errors.New("GetPricingSpecRequest cannot be nil")
	}
	return nil, errors.New("GetPricingSpec not implemented")
}

// EstimateCost provides a default implementation that returns not implemented.
// Override this method in your plugin to return cost estimates.
func (bp *BasePlugin) EstimateCost(
	_ context.Context,
	req *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	if req == nil {
		return nil, errors.New("EstimateCostRequest cannot be nil")
	}
	return nil, errors.New("EstimateCost not implemented")
}
