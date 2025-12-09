package pluginsdk

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/rshade/pulumicost-spec/sdk/go/currency"
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
func (cc *CostCalculator) CreateActualCostResponse(
	results []*pbc.ActualCostResult,
) *pbc.GetActualCostResponse {
	return &pbc.GetActualCostResponse{
		Results: results,
	}
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

// =============================================================================
// GetRecommendations Validation Helpers
// =============================================================================

// ValidateRecommendation validates a recommendation has all required fields.
// Returns an error if any required field is missing or invalid.
func ValidateRecommendation(rec *pbc.Recommendation) error {
	if rec == nil {
		return errors.New("recommendation cannot be nil")
	}
	if rec.GetId() == "" {
		return errors.New("recommendation.id is required")
	}
	if rec.GetCategory() == pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_UNSPECIFIED {
		return errors.New("recommendation.category must be specified")
	}
	if rec.GetActionType() == pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED {
		return errors.New("recommendation.action_type must be specified")
	}
	if rec.GetResource() == nil {
		return errors.New("recommendation.resource is required")
	}
	if err := ValidateResourceRecommendationInfo(rec.GetResource()); err != nil {
		return fmt.Errorf("recommendation.resource: %w", err)
	}
	if rec.GetImpact() == nil {
		return errors.New("recommendation.impact is required")
	}
	if err := ValidateRecommendationImpact(rec.GetImpact()); err != nil {
		return fmt.Errorf("recommendation.impact: %w", err)
	}
	if rec.ConfidenceScore != nil {
		if err := ValidateConfidenceScore(rec.ConfidenceScore); err != nil {
			return fmt.Errorf("recommendation.confidence_score: %w", err)
		}
	}
	return nil
}

// ValidateResourceRecommendationInfo validates resource information fields.
func ValidateResourceRecommendationInfo(res *pbc.ResourceRecommendationInfo) error {
	if res == nil {
		return errors.New("resource info cannot be nil")
	}
	if res.GetId() == "" {
		return errors.New("resource.id is required")
	}
	if res.GetProvider() == "" {
		return errors.New("resource.provider is required")
	}
	return nil
}

// ValidateRecommendationImpact validates impact information fields.
// Validates currency against ISO 4217 standard.
func ValidateRecommendationImpact(impact *pbc.RecommendationImpact) error {
	if impact == nil {
		return errors.New("impact cannot be nil")
	}
	if impact.GetCurrency() == "" {
		return errors.New("impact.currency is required")
	}
	if !currency.IsValid(impact.GetCurrency()) {
		return fmt.Errorf("impact.currency %q is not a valid ISO 4217 currency code", impact.GetCurrency())
	}
	if impact.GetEstimatedSavings() < 0 {
		return errors.New("impact.estimated_savings cannot be negative")
	}
	return nil
}

// ValidateConfidenceScore validates that a confidence score is within the valid range.
// Confidence scores must be between 0.0 and 1.0 inclusive.
func ValidateConfidenceScore(score *float64) error {
	if score == nil {
		return nil // nil is valid (confidence not available)
	}
	if *score < 0.0 || *score > 1.0 {
		return fmt.Errorf("confidence_score must be between 0.0 and 1.0, got %f", *score)
	}
	return nil
}

// ValidateRecommendationSummary validates summary information fields.
func ValidateRecommendationSummary(summary *pbc.RecommendationSummary) error {
	if summary == nil {
		return nil // Summary is optional
	}
	if summary.GetCurrency() != "" && !currency.IsValid(summary.GetCurrency()) {
		return fmt.Errorf("summary.currency %q is not a valid ISO 4217 currency code", summary.GetCurrency())
	}
	return nil
}

// =============================================================================
// GetRecommendations Filter Helpers
// =============================================================================

// ApplyRecommendationFilter filters recommendations based on the provided filter criteria.
// ApplyRecommendationFilter returns recommendations that match ALL specified filter criteria.
// Empty filter values are ignored (match all).
func ApplyRecommendationFilter(
	recommendations []*pbc.Recommendation,
	filter *pbc.RecommendationFilter,
) []*pbc.Recommendation {
	if filter == nil {
		return recommendations
	}

	result := make([]*pbc.Recommendation, 0, len(recommendations))
	for _, rec := range recommendations {
		if matchesFilter(rec, filter) {
			result = append(result, rec)
		}
	}
	return result
}

// matchesFilter checks if a recommendation matches all filter criteria.
func matchesFilter(rec *pbc.Recommendation, filter *pbc.RecommendationFilter) bool {
	// Filter by provider
	if filter.GetProvider() != "" {
		if rec.GetResource() == nil || rec.GetResource().GetProvider() != filter.GetProvider() {
			return false
		}
	}

	// Filter by region
	if filter.GetRegion() != "" {
		if rec.GetResource() == nil || rec.GetResource().GetRegion() != filter.GetRegion() {
			return false
		}
	}

	// Filter by resource type
	if filter.GetResourceType() != "" {
		if rec.GetResource() == nil || rec.GetResource().GetResourceType() != filter.GetResourceType() {
			return false
		}
	}

	// Filter by category
	if filter.GetCategory() != pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_UNSPECIFIED {
		if rec.GetCategory() != filter.GetCategory() {
			return false
		}
	}

	// Filter by action type
	if filter.GetActionType() != pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED {
		if rec.GetActionType() != filter.GetActionType() {
			return false
		}
	}

	return true
}

// =============================================================================
// GetRecommendations Pagination Helpers
// =============================================================================

// DefaultPageSize is the default number of recommendations per page.
const DefaultPageSize = 50

// MaxPageSize is the maximum allowed page size.
const MaxPageSize = 1000

// PaginateRecommendations applies pagination to a slice of recommendations.
// PaginateRecommendations returns the page of recommendations and the next page token (empty if last page).
func PaginateRecommendations(
	recommendations []*pbc.Recommendation,
	pageSize int32,
	pageToken string,
) ([]*pbc.Recommendation, string, error) {
	// Determine effective page size
	effectivePageSize := int(pageSize)
	if effectivePageSize <= 0 {
		effectivePageSize = DefaultPageSize
	}
	if effectivePageSize > MaxPageSize {
		effectivePageSize = MaxPageSize
	}

	// Decode offset from page token
	offset := 0
	if pageToken != "" {
		var err error
		offset, err = DecodePageToken(pageToken)
		if err != nil {
			return nil, "", fmt.Errorf("invalid page_token: %w", err)
		}
	}

	// Handle out-of-bounds offset
	total := len(recommendations)
	if offset >= total {
		return []*pbc.Recommendation{}, "", nil
	}

	// Calculate page boundaries
	end := offset + effectivePageSize
	if end > total {
		end = total
	}

	// Extract page
	page := recommendations[offset:end]

	// Generate next page token
	nextToken := ""
	if end < total {
		nextToken = EncodePageToken(end)
	}

	return page, nextToken, nil
}

// EncodePageToken creates an opaque page token from an offset.
func EncodePageToken(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}

// DecodePageToken decodes an opaque page token to an offset.
func DecodePageToken(token string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, errors.New("malformed page token")
	}
	offset, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0, errors.New("invalid page token value")
	}
	if offset < 0 {
		return 0, errors.New("page token offset cannot be negative")
	}
	return offset, nil
}

// =============================================================================
// GetRecommendations Summary Calculation
// =============================================================================

// CalculateRecommendationSummary computes aggregated summary statistics for recommendations.
func CalculateRecommendationSummary(
	recommendations []*pbc.Recommendation,
	projectionPeriod string,
) *pbc.RecommendationSummary {
	summary := &pbc.RecommendationSummary{
		TotalRecommendations: int32(len(recommendations)), //nolint:gosec // length will not exceed int32 max
		ProjectionPeriod:     projectionPeriod,
		CountByCategory:      make(map[string]int32),
		SavingsByCategory:    make(map[string]float64),
		CountByActionType:    make(map[string]int32),
		SavingsByActionType:  make(map[string]float64),
	}

	var totalSavings float64
	var detectedCurrency string
	var currencyMismatch bool

	for _, rec := range recommendations {
		catName := rec.GetCategory().String()
		actionName := rec.GetActionType().String()

		summary.CountByCategory[catName]++
		summary.CountByActionType[actionName]++

		if impact := rec.GetImpact(); impact != nil {
			savings := impact.GetEstimatedSavings()
			totalSavings += savings
			summary.SavingsByCategory[catName] += savings
			summary.SavingsByActionType[actionName] += savings
			if c := impact.GetCurrency(); c != "" {
				if detectedCurrency == "" {
					detectedCurrency = c
				} else if detectedCurrency != c {
					currencyMismatch = true
				}
			}
		}
	}
	// Clear currency if recommendations have mixed currencies (sum is ambiguous)
	if currencyMismatch {
		detectedCurrency = ""
	}

	summary.TotalEstimatedSavings = totalSavings
	summary.Currency = detectedCurrency

	return summary
}
