// Package testing provides test utilities for PulumiCost plugin development.
// It includes mock plugin implementations, test harnesses, and conformance testing
// utilities for validating plugin behavior against the CostSource gRPC service spec.
package testing

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

const (
	computeEngineResourceType  = "compute_engine"
	ec2ResourceType            = "ec2"
	blobStorageResourceType    = "blob_storage"
	cloudFunctionsResourceType = "cloud_functions"
	namespaceResourceType      = "namespace"
	cloudStorageResourceType   = "cloud_storage"
	lambdaResourceType         = "lambda"

	// Time and performance constants.
	defaultDataPoints    = 24   // 24 hours of hourly data
	defaultBaseRate      = 0.05 // Default hourly rate
	nameDelayMs          = 100  // Name RPC delay in milliseconds
	supportsDelayMs      = 200  // Supports RPC delay in milliseconds
	actualCostDelayMs    = 500  // ActualCost RPC delay in milliseconds
	projectedCostDelayMs = 300  // ProjectedCost RPC delay in milliseconds
	pricingSpecDelayMs   = 250  // PricingSpec RPC delay in milliseconds
	estimateCostDelayMs  = 150  // EstimateCost RPC delay in milliseconds
	costVariationMod     = 10   // Modulo for cost variation
	daysPerMonth         = 30   // Days per month for cost calculation

	// Cost variation constants.
	costVariationBase  = 0.8 // Base cost variation multiplier (80%)
	costVariationRange = 0.4 // Cost variation range (40% additional)

	// Rate multipliers for different resource types.
	computeRateMultiplier    = 2.0   // Compute resources (EC2, VM, Compute Engine)
	storageRateMultiplier    = 0.1   // Storage resources (S3, Blob, Cloud Storage)
	serverlessRateMultiplier = 0.001 // Serverless (Lambda, Cloud Functions)
	namespaceRateMultiplier  = 1.5   // Kubernetes namespace
	databaseRateMultiplier   = 3.0   // Database resources (SQL Database)

	// Recommendation generation constants.
	confidenceVariations = 4     // Number of confidence score variations (0.70, 0.775, 0.85, 0.925)
	confidenceBase       = 0.7   // Base confidence score
	confidenceStep       = 0.075 // Confidence increment per variation
	savingsBase          = 50.0  // Base savings amount
	savingsIncrement     = 25.0  // Savings increment per recommendation
)

// MockPlugin provides a configurable mock implementation of CostSourceServiceServer.
type MockPlugin struct {
	pbc.UnimplementedCostSourceServiceServer

	// Configuration
	PluginName         string
	SupportedProviders []string
	SupportedResources map[string][]string // provider -> resource types

	// Behavior configuration
	ShouldErrorOnName          bool
	ShouldErrorOnSupports      bool
	ShouldErrorOnActualCost    bool
	ShouldErrorOnProjectedCost bool
	ShouldErrorOnPricingSpec   bool
	ShouldErrorOnEstimateCost  bool

	// Response delays for testing timeouts
	NameDelay          time.Duration
	SupportsDelay      time.Duration
	ActualCostDelay    time.Duration
	ProjectedCostDelay time.Duration
	PricingSpecDelay   time.Duration
	EstimateCostDelay  time.Duration

	// Data generation configuration
	ActualCostDataPoints int
	BaseHourlyRate       float64
	Currency             string

	// Recommendations configuration
	RecommendationsConfig RecommendationsConfig

	// Budgets configuration
	ShouldErrorOnBudgets bool
	MockBudgets          []*pbc.Budget

	// FallbackHint configuration for GetActualCost responses.
	// Thread Safety: This field must be set before the plugin begins serving
	// requests. Use SetFallbackHint() for configuration, which documents the
	// thread safety constraints.
	FallbackHint pbc.FallbackHint
}

// NewMockPlugin creates a new mock plugin with default configuration.
func NewMockPlugin() *MockPlugin {
	return &MockPlugin{
		PluginName:         "mock-test-plugin",
		SupportedProviders: []string{"aws", "azure", "gcp", "kubernetes"},
		SupportedResources: map[string][]string{
			"aws":        {ec2ResourceType, "s3", lambdaResourceType, "rds"},
			"azure":      {"vm", blobStorageResourceType, "sql_database", "compute"},
			"gcp":        {computeEngineResourceType, cloudStorageResourceType, cloudFunctionsResourceType, "compute"},
			"kubernetes": {namespaceResourceType, "pod", "service"},
		},
		ActualCostDataPoints: defaultDataPoints,
		BaseHourlyRate:       defaultBaseRate,
		Currency:             "USD",
		// Pre-populate with sample recommendations for filtering tests
		RecommendationsConfig: RecommendationsConfig{
			Recommendations: GenerateSampleRecommendations(defaultRecommendationCount),
		},
	}
}

const defaultRecommendationCount = 12 // 12 samples to cover all categories and action types

// ConfigurableErrorMockPlugin creates a mock plugin that can be configured to return errors.
func ConfigurableErrorMockPlugin() *MockPlugin {
	plugin := NewMockPlugin()
	plugin.PluginName = "error-test-plugin"
	return plugin
}

// SlowMockPlugin creates a mock plugin with artificial delays for timeout testing.
func SlowMockPlugin() *MockPlugin {
	plugin := NewMockPlugin()
	plugin.PluginName = "slow-test-plugin"
	plugin.NameDelay = nameDelayMs * time.Millisecond
	plugin.SupportsDelay = supportsDelayMs * time.Millisecond
	plugin.ActualCostDelay = actualCostDelayMs * time.Millisecond
	plugin.ProjectedCostDelay = projectedCostDelayMs * time.Millisecond
	plugin.PricingSpecDelay = pricingSpecDelayMs * time.Millisecond
	plugin.EstimateCostDelay = estimateCostDelayMs * time.Millisecond
	return plugin
}

// SetFallbackHint sets the fallback hint to be returned in GetActualCost responses.
//
// Thread Safety: This method is NOT safe for concurrent use. All calls to
// SetFallbackHint must complete before the plugin begins serving requests.
// Typical usage is to configure the mock during test setup, before calling
// harness.Start() or benchmark.ResetTimer().
func (m *MockPlugin) SetFallbackHint(hint pbc.FallbackHint) {
	m.FallbackHint = hint
}

// Name returns the plugin name.
func (m *MockPlugin) Name(_ context.Context, _ *pbc.NameRequest) (*pbc.NameResponse, error) {
	if m.NameDelay > 0 {
		time.Sleep(m.NameDelay)
	}

	if m.ShouldErrorOnName {
		return nil, status.Error(codes.Internal, "mock error: name operation failed")
	}

	return &pbc.NameResponse{
		Name: m.PluginName,
	}, nil
}

// Supports checks if a resource type is supported by this mock plugin.
func (m *MockPlugin) Supports(_ context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
	if m.SupportsDelay > 0 {
		time.Sleep(m.SupportsDelay)
	}

	if m.ShouldErrorOnSupports {
		return nil, status.Error(codes.InvalidArgument, "mock error: supports operation failed")
	}

	resource := req.GetResource()
	if resource == nil {
		return &pbc.SupportsResponse{
			Supported: false,
			Reason:    "resource descriptor is required",
		}, nil
	}

	provider := resource.GetProvider()
	resourceType := resource.GetResourceType()

	// Check if provider is supported
	providerSupported := false
	for _, supportedProvider := range m.SupportedProviders {
		if provider == supportedProvider {
			providerSupported = true
			break
		}
	}

	if !providerSupported {
		return &pbc.SupportsResponse{
			Supported: false,
			Reason:    fmt.Sprintf("provider %s is not supported", provider),
		}, nil
	}

	// Check if resource type is supported for this provider
	supportedResources, exists := m.SupportedResources[provider]
	if !exists {
		return &pbc.SupportsResponse{
			Supported: false,
			Reason:    fmt.Sprintf("no resource types configured for provider %s", provider),
		}, nil
	}

	for _, supportedResource := range supportedResources {
		if resourceType == supportedResource {
			return &pbc.SupportsResponse{
				Supported: true,
				Reason:    "",
			}, nil
		}
	}

	return &pbc.SupportsResponse{
		Supported: false,
		Reason:    fmt.Sprintf("resource type %s is not supported for provider %s", resourceType, provider),
	}, nil
}

// GetActualCost returns mock historical cost data.
func (m *MockPlugin) GetActualCost(
	_ context.Context,
	req *pbc.GetActualCostRequest,
) (*pbc.GetActualCostResponse, error) {
	if m.ActualCostDelay > 0 {
		time.Sleep(m.ActualCostDelay)
	}

	if m.ShouldErrorOnActualCost {
		return nil, status.Error(codes.NotFound, "mock error: actual cost data not available")
	}

	start := req.GetStart()
	end := req.GetEnd()

	if start == nil || end == nil {
		return nil, status.Error(codes.InvalidArgument, "start and end timestamps are required")
	}

	startTime := start.AsTime()
	endTime := end.AsTime()

	if endTime.Before(startTime) {
		return nil, status.Error(codes.InvalidArgument, "end time must be after start time")
	}

	// Generate mock cost data points
	duration := endTime.Sub(startTime)
	dataPoints := int(math.Min(float64(m.ActualCostDataPoints), duration.Hours()+1))

	results := make([]*pbc.ActualCostResult, 0, dataPoints)

	for i := range dataPoints {
		timestamp := startTime.Add(time.Duration(i) * time.Hour)
		if timestamp.After(endTime) {
			break
		}

		// Generate some variation in cost data
		costVariation := costVariationBase + (costVariationRange * float64(i%costVariationMod) / costVariationMod) // Varies between 0.8x and 1.2x base rate
		cost := m.BaseHourlyRate * costVariation
		usageAmount := 1.0 * costVariation // Assume 1 unit of usage

		result := &pbc.ActualCostResult{
			Timestamp:   timestamppb.New(timestamp),
			Cost:        cost,
			UsageAmount: usageAmount,
			UsageUnit:   "hour",
			Source:      m.PluginName,
		}
		results = append(results, result)
	}

	return &pbc.GetActualCostResponse{
		Results:      results,
		FallbackHint: m.FallbackHint,
	}, nil
}

// GetProjectedCost returns mock projected cost data.
func (m *MockPlugin) GetProjectedCost(
	_ context.Context,
	req *pbc.GetProjectedCostRequest,
) (*pbc.GetProjectedCostResponse, error) {
	if m.ProjectedCostDelay > 0 {
		time.Sleep(m.ProjectedCostDelay)
	}

	if m.ShouldErrorOnProjectedCost {
		return nil, status.Error(codes.Unavailable, "mock error: projected cost service unavailable")
	}

	resource := req.GetResource()
	if resource == nil {
		return nil, status.Error(codes.InvalidArgument, "resource descriptor is required")
	}

	// Calculate cost based on resource type (keep in sync with GetPricingSpec).
	multiplier := getRateMultiplier(resource.GetResourceType())

	unitPrice := m.BaseHourlyRate * multiplier
	costPerMonth := unitPrice * HoursPerDay * daysPerMonth

	billingDetail := fmt.Sprintf("mock-%s-rate", strings.ToLower(resource.GetProvider()))

	return &pbc.GetProjectedCostResponse{
		UnitPrice:     unitPrice,
		Currency:      m.Currency,
		CostPerMonth:  costPerMonth,
		BillingDetail: billingDetail,
	}, nil
}

// getBillingModeAndUnit returns billing mode and unit for a resource type.
func getBillingModeAndUnit(resourceType string) (string, string) {
	switch resourceType {
	case "s3", blobStorageResourceType, cloudStorageResourceType:
		return "per_gb_month", "GB-month"
	case lambdaResourceType, cloudFunctionsResourceType:
		return "per_invocation", "request"
	case namespaceResourceType:
		return "per_cpu_hour", "hour"
	case "sql_database":
		return "per_dtu", "DTU"
	default:
		return "per_hour", "hour"
	}
}

// isKnownResourceType returns true if the resource type is known/supported.
func isKnownResourceType(resourceType string) bool {
	knownTypes := []string{
		ec2ResourceType, "s3", lambdaResourceType, "rds",
		"vm", blobStorageResourceType, "sql_database",
		computeEngineResourceType, cloudStorageResourceType, cloudFunctionsResourceType,
		namespaceResourceType, "pod", "service",
	}
	for _, known := range knownTypes {
		if resourceType == known {
			return true
		}
	}
	return false
}

// getRateMultiplier returns the rate multiplier for a resource type.
func getRateMultiplier(resourceType string) float64 {
	switch resourceType {
	case ec2ResourceType, "vm", computeEngineResourceType, "compute":
		return computeRateMultiplier
	case "s3", blobStorageResourceType, cloudStorageResourceType:
		return storageRateMultiplier
	case lambdaResourceType, cloudFunctionsResourceType:
		return serverlessRateMultiplier
	case namespaceResourceType:
		return namespaceRateMultiplier
	case "sql_database":
		return databaseRateMultiplier
	default:
		return 1.0
	}
}

// getMetricHints returns metric hints for a resource type.
func getMetricHints(resourceType string) []*pbc.UsageMetricHint {
	switch resourceType {
	case ec2ResourceType, "vm", computeEngineResourceType, "compute":
		return []*pbc.UsageMetricHint{
			{Metric: "vcpu_hours", Unit: "hour"},
			{Metric: "memory_gb_hours", Unit: "hour"},
		}
	case "s3", blobStorageResourceType, cloudStorageResourceType:
		return []*pbc.UsageMetricHint{
			{Metric: "storage_gb", Unit: "GB"},
			{Metric: "requests", Unit: "count"},
		}
	case lambdaResourceType, cloudFunctionsResourceType:
		return []*pbc.UsageMetricHint{
			{Metric: "invocations", Unit: "count"},
			{Metric: "duration_ms", Unit: "millisecond"},
		}
	case namespaceResourceType:
		return []*pbc.UsageMetricHint{
			{Metric: "cpu_cores", Unit: "hour"},
			{Metric: "memory_gb", Unit: "hour"},
		}
	default:
		return []*pbc.UsageMetricHint{}
	}
}

// GetPricingSpec returns mock pricing specification.
func (m *MockPlugin) GetPricingSpec(
	_ context.Context,
	req *pbc.GetPricingSpecRequest,
) (*pbc.GetPricingSpecResponse, error) {
	if m.PricingSpecDelay > 0 {
		time.Sleep(m.PricingSpecDelay)
	}

	if m.ShouldErrorOnPricingSpec {
		return nil, status.Error(codes.PermissionDenied, "mock error: pricing spec access denied")
	}

	resource := req.GetResource()
	if resource == nil {
		return nil, status.Error(codes.InvalidArgument, "resource descriptor is required")
	}

	// Validate required fields (FR-011)
	if resource.GetProvider() == "" {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}
	if resource.GetResourceType() == "" {
		return nil, status.Error(codes.InvalidArgument, "resource_type is required")
	}

	resourceType := resource.GetResourceType()

	// Handle unknown resource types with not_implemented (T039)
	if !isKnownResourceType(resourceType) {
		return &pbc.GetPricingSpecResponse{
			Spec: &pbc.PricingSpec{
				Provider:     resource.GetProvider(),
				ResourceType: resourceType,
				Sku:          resource.GetSku(),
				Region:       resource.GetRegion(),
				BillingMode:  "not_implemented",
				RatePerUnit:  0,
				Currency:     m.Currency,
				Unit:         "unknown",
				Description:  fmt.Sprintf("Pricing not implemented for %s", resourceType),
				Assumptions: []string{
					fmt.Sprintf("Resource type %s is not yet supported", resourceType),
					"Pricing data is unavailable for this resource type",
				},
				Source: m.PluginName,
			},
		}, nil
	}

	billingMode, unit := getBillingModeAndUnit(resourceType)
	ratePerUnit := m.BaseHourlyRate * getRateMultiplier(resourceType)
	metricHints := getMetricHints(resourceType)

	// Generate assumptions based on resource type
	assumptions := []string{
		fmt.Sprintf(
			"Pricing based on %s %s in %s",
			resource.GetProvider(),
			resource.GetResourceType(),
			resource.GetRegion(),
		),
		"On-demand pricing without reserved capacity discounts",
		"Standard tier without additional features",
	}

	spec := &pbc.PricingSpec{
		Provider:     resource.GetProvider(),
		ResourceType: resource.GetResourceType(),
		Sku:          resource.GetSku(),
		Region:       resource.GetRegion(),
		BillingMode:  billingMode,
		RatePerUnit:  ratePerUnit,
		Currency:     m.Currency,
		Description:  fmt.Sprintf("Mock pricing for %s %s", resource.GetProvider(), resource.GetResourceType()),
		MetricHints:  metricHints,
		PluginMetadata: map[string]string{
			"mock_plugin_version": "1.0.0",
			"test_mode":           "true",
			"data_source":         "synthetic",
		},
		Source:      m.PluginName,
		Unit:        unit,
		Assumptions: assumptions,
	}

	return &pbc.GetPricingSpecResponse{
		Spec: spec,
	}, nil
}

// parseResourceType parses a resource type string in format provider:module/resource:Type.
// Returns (provider, module, resource, typeName, error).
func parseResourceType(resourceType string) (string, string, string, string, error) {
	const (
		expectedColonParts = 3
		expectedSlashParts = 2
	)
	// Split into provider, module/resource, and type name
	parts := strings.SplitN(resourceType, ":", expectedColonParts)
	if len(parts) != expectedColonParts || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", "",
			fmt.Errorf(
				"invalid format: expected provider:module/resource:Type, got %s",
				resourceType,
			)
	}

	provider := parts[0]
	moduleResource := parts[1]
	typeName := parts[2]

	// Split module/resource
	moduleResourceParts := strings.SplitN(moduleResource, "/", expectedSlashParts)
	if len(moduleResourceParts) != expectedSlashParts || moduleResourceParts[0] == "" || moduleResourceParts[1] == "" {
		return "", "", "", "",
			fmt.Errorf(
				"invalid format: expected provider:module/resource:Type, got %s",
				resourceType,
			)
	}

	module := moduleResourceParts[0]
	resource := moduleResourceParts[1]

	return provider, module, resource, typeName, nil
}

// =============================================================================
// GetRecommendations Support
// =============================================================================

// RecommendationsConfig holds configuration for mock recommendations.
type RecommendationsConfig struct {
	Recommendations []*pbc.Recommendation
	ShouldError     bool
	ErrorMessage    string
	Delay           time.Duration
}

// SetRecommendationsConfig configures the recommendations response.
func (m *MockPlugin) SetRecommendationsConfig(config RecommendationsConfig) {
	m.RecommendationsConfig = config
}

// GenerateSampleRecommendations creates realistic test recommendation data.
func GenerateSampleRecommendations(count int) []*pbc.Recommendation {
	categories := []pbc.RecommendationCategory{
		pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_COST,
		pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_PERFORMANCE,
		pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_SECURITY,
		pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_RELIABILITY,
	}

	actionTypes := []pbc.RecommendationActionType{
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_RIGHTSIZE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_TERMINATE,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_PURCHASE_COMMITMENT,
		pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_MODIFY,
	}

	priorities := []pbc.RecommendationPriority{
		pbc.RecommendationPriority_RECOMMENDATION_PRIORITY_HIGH,
		pbc.RecommendationPriority_RECOMMENDATION_PRIORITY_MEDIUM,
		pbc.RecommendationPriority_RECOMMENDATION_PRIORITY_LOW,
	}

	providers := []string{"aws", "azure", "gcp", "kubernetes"}
	regions := []string{"us-east-1", "us-west-2", "eu-west-1", "asia-pacific-1"}

	recs := make([]*pbc.Recommendation, count)
	for i := range count {
		confidence := confidenceBase + float64(i%confidenceVariations)*confidenceStep
		savings := savingsBase + float64(i)*savingsIncrement

		recs[i] = &pbc.Recommendation{
			Id:          fmt.Sprintf("rec-%d", i+1),
			Category:    categories[i%len(categories)],
			ActionType:  actionTypes[i%len(actionTypes)],
			Priority:    priorities[i%len(priorities)],
			Description: fmt.Sprintf("Recommendation %d: Optimize resource with cost savings", i+1),
			Resource: &pbc.ResourceRecommendationInfo{
				Id:           fmt.Sprintf("resource-%d", i+1),
				Provider:     providers[i%len(providers)],
				ResourceType: fmt.Sprintf("%s:compute/instance:Instance", providers[i%len(providers)]),
				Region:       regions[i%len(regions)],
				Name:         fmt.Sprintf("instance-%d", i+1),
			},
			Impact: &pbc.RecommendationImpact{
				EstimatedSavings: savings,
				Currency:         "USD",
				ProjectionPeriod: "12_months",
			},
			ConfidenceScore: &confidence,
		}
	}
	return recs
}

// GetRecommendations returns mock cost optimization recommendations.
func (m *MockPlugin) GetRecommendations(
	_ context.Context,
	req *pbc.GetRecommendationsRequest,
) (*pbc.GetRecommendationsResponse, error) {
	if m.RecommendationsConfig.Delay > 0 {
		time.Sleep(m.RecommendationsConfig.Delay)
	}

	if m.RecommendationsConfig.ShouldError {
		msg := m.RecommendationsConfig.ErrorMessage
		if msg == "" {
			msg = "mock error: recommendations unavailable"
		}
		return nil, status.Error(codes.Unavailable, msg)
	}

	// Validate request contract constraints (including MaxTargetResources limit)
	if err := ValidateGetRecommendationsRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "request validation failed: %v", err)
	}

	// Return configured recommendations or empty list
	recs := m.RecommendationsConfig.Recommendations
	if recs == nil {
		recs = []*pbc.Recommendation{}
	}

	// Apply target_resources filtering first (defines SCOPE).
	// When provided, only recommendations matching target resources are included.
	// When empty/nil, all recommendations are returned (backward compatible).
	if len(req.GetTargetResources()) > 0 {
		recs = filterByTargetResources(recs, req.GetTargetResources())
	}

	// Apply filter if provided (defines SELECTION CRITERIA within scope).
	// NOTE: This uses a local implementation rather than pluginsdk.ApplyRecommendationFilter
	// to avoid circular imports (pluginsdk imports testing for conformance functions).
	if req.GetFilter() != nil {
		recs = applyMockFilter(recs, req.GetFilter())
	}

	// Apply pagination if page_size is specified
	var nextToken string
	if req.GetPageSize() > 0 || req.GetPageToken() != "" {
		var paginationErr error
		recs, nextToken, paginationErr = paginateMockRecommendations(
			recs,
			req.GetPageSize(),
			req.GetPageToken(),
		)
		if paginationErr != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token: %v", paginationErr)
		}
	}

	// NOTE: This mock intentionally calculates summary from the paginated (current page)
	// results, not the total filtered set. This is by design for testing pagination
	// scenarios where clients need to verify per-page behavior. Production implementations
	// may calculate summary from the full filtered dataset before pagination.
	summary := CalculateMockSummary(recs, req.GetProjectionPeriod())

	return &pbc.GetRecommendationsResponse{
		Recommendations: recs,
		Summary:         summary,
		NextPageToken:   nextToken,
	}, nil
}

// GetBudgets implements the mock GetBudgets RPC method.
//
//nolint:gocognit // Budget filtering logic requires multiple conditions
func (m *MockPlugin) GetBudgets(
	_ context.Context,
	req *pbc.GetBudgetsRequest,
) (*pbc.GetBudgetsResponse, error) {
	if m.ShouldErrorOnBudgets {
		return nil, status.Error(codes.Internal, "mock error")
	}

	// Return configured budgets or empty list
	budgets := m.MockBudgets
	if budgets == nil {
		budgets = []*pbc.Budget{}
	}

	// Apply basic filtering if requested
	if req.GetFilter() != nil {
		filter := req.GetFilter()
		if len(filter.GetProviders()) > 0 {
			var filtered []*pbc.Budget
			for _, budget := range budgets {
				for _, provider := range filter.GetProviders() {
					if strings.Contains(budget.GetSource(), provider) {
						filtered = append(filtered, budget)
						break
					}
				}
			}
			budgets = filtered
		}
	}

	// Calculate summary
	totalBudgets := len(budgets)
	budgetsOk := 0
	budgetsWarning := 0
	budgetsCritical := 0
	budgetsExceeded := 0

	for _, budget := range budgets {
		if budget.GetStatus() != nil {
			switch budget.GetStatus().GetHealth() {
			case pbc.BudgetHealthStatus_BUDGET_HEALTH_STATUS_OK,
				pbc.BudgetHealthStatus_BUDGET_HEALTH_STATUS_UNSPECIFIED:
				budgetsOk++
			case pbc.BudgetHealthStatus_BUDGET_HEALTH_STATUS_WARNING:
				budgetsWarning++
			case pbc.BudgetHealthStatus_BUDGET_HEALTH_STATUS_CRITICAL:
				budgetsCritical++
			case pbc.BudgetHealthStatus_BUDGET_HEALTH_STATUS_EXCEEDED:
				budgetsExceeded++
			}
		} else {
			budgetsOk++ // No status means assume OK
		}
	}

	return &pbc.GetBudgetsResponse{
		Budgets: budgets,
		Summary: &pbc.BudgetSummary{
			TotalBudgets:    int32(totalBudgets),    //nolint:gosec // length will not exceed int32 max
			BudgetsOk:       int32(budgetsOk),       //nolint:gosec // count will not exceed int32 max
			BudgetsWarning:  int32(budgetsWarning),  //nolint:gosec // count will not exceed int32 max
			BudgetsCritical: int32(budgetsCritical), //nolint:gosec // count will not exceed int32 max
			BudgetsExceeded: int32(budgetsExceeded), //nolint:gosec // count will not exceed int32 max
		},
	}, nil
}

// CalculateMockSummary builds a RecommendationSummary from the given recommendations.
// NOTE: This duplicates pluginsdk.CalculateRecommendationSummary logic to avoid circular
// imports (pluginsdk imports testing for conformance functions).
func CalculateMockSummary(recs []*pbc.Recommendation, projectionPeriod string) *pbc.RecommendationSummary {
	summary := &pbc.RecommendationSummary{
		TotalRecommendations: int32(len(recs)), //nolint:gosec // length will not exceed int32 max
		CountByCategory:      make(map[string]int32),
		SavingsByCategory:    make(map[string]float64),
		CountByActionType:    make(map[string]int32),
		SavingsByActionType:  make(map[string]float64),
		ProjectionPeriod:     projectionPeriod,
	}

	var totalSavings float64
	var detectedCurrency string
	var currencyMismatch bool
	for _, rec := range recs {
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

// =============================================================================
// Target Resources Filtering (Feature 019-target-resources)
// =============================================================================

// filterByTargetResources filters recommendations to only include those matching
// at least one of the target resources. This implements the SCOPE filtering for
// stack-scoped recommendations, pre-deployment optimization, and batch analysis.
//
// Matching uses OR logic: a recommendation is included if it matches ANY target.
func filterByTargetResources(recs []*pbc.Recommendation, targets []*pbc.ResourceDescriptor) []*pbc.Recommendation {
	result := make([]*pbc.Recommendation, 0, len(recs))
	for _, rec := range recs {
		if matchesAnyTargetResource(rec, targets) {
			result = append(result, rec)
		}
	}
	return result
}

// matchesAnyTargetResource checks if a recommendation matches at least one target resource.
func matchesAnyTargetResource(rec *pbc.Recommendation, targets []*pbc.ResourceDescriptor) bool {
	resource := rec.GetResource()
	if resource == nil {
		return false
	}

	for _, target := range targets {
		if matchesResourceDescriptor(resource, target) {
			return true
		}
	}
	return false
}

// matchesResourceDescriptor checks if a ResourceRecommendationInfo matches a ResourceDescriptor.
// Matching rules (per spec FR-008):
//   - provider and resource_type must always match (required fields)
//   - sku, region, and tags are matched only when specified in the target
//   - If specified, optional fields must match exactly (strict matching)
func matchesResourceDescriptor(resource *pbc.ResourceRecommendationInfo, target *pbc.ResourceDescriptor) bool {
	if target == nil {
		return false
	}

	// Required fields must always match
	if resource.GetProvider() != target.GetProvider() {
		return false
	}
	if resource.GetResourceType() != target.GetResourceType() {
		return false
	}

	// Optional fields: only check if specified in target (strict matching)
	if target.GetSku() != "" && resource.GetSku() != target.GetSku() {
		return false
	}
	if target.GetRegion() != "" && resource.GetRegion() != target.GetRegion() {
		return false
	}

	// Tags: all specified tags must be present with exact values
	if len(target.GetTags()) > 0 {
		resourceTags := resource.GetTags()
		for key, targetValue := range target.GetTags() {
			if resourceTags[key] != targetValue {
				return false
			}
		}
	}

	return true
}

// applyMockFilter filters recommendations based on the filter criteria.
// NOTE: This duplicates pluginsdk.ApplyRecommendationFilter logic to avoid circular
// imports (pluginsdk imports testing for conformance functions).
func applyMockFilter(recs []*pbc.Recommendation, filter *pbc.RecommendationFilter) []*pbc.Recommendation {
	result := make([]*pbc.Recommendation, 0, len(recs))
	for _, rec := range recs {
		if matchesMockFilter(rec, filter) {
			result = append(result, rec)
		}
	}
	return result
}

// matchesMockFilter checks if a recommendation matches the filter criteria.
func matchesMockFilter(rec *pbc.Recommendation, filter *pbc.RecommendationFilter) bool {
	// Check provider filter
	if filter.GetProvider() != "" {
		if rec.GetResource() == nil || rec.GetResource().GetProvider() != filter.GetProvider() {
			return false
		}
	}

	// Check region filter
	if filter.GetRegion() != "" {
		if rec.GetResource() == nil || rec.GetResource().GetRegion() != filter.GetRegion() {
			return false
		}
	}

	// Check resource_type filter
	if filter.GetResourceType() != "" {
		if rec.GetResource() == nil || rec.GetResource().GetResourceType() != filter.GetResourceType() {
			return false
		}
	}

	// Check category filter
	if filter.GetCategory() != pbc.RecommendationCategory_RECOMMENDATION_CATEGORY_UNSPECIFIED {
		if rec.GetCategory() != filter.GetCategory() {
			return false
		}
	}

	// Check action_type filter
	if filter.GetActionType() != pbc.RecommendationActionType_RECOMMENDATION_ACTION_TYPE_UNSPECIFIED {
		if rec.GetActionType() != filter.GetActionType() {
			return false
		}
	}

	return true
}

// mockDefaultPageSize is the default page size for mock pagination.
// NOTE: This intentionally mirrors pluginsdk.DefaultPageSize (50) for consistency.
// We maintain a local constant to avoid circular imports (pluginsdk imports testing).
const mockDefaultPageSize int32 = 50

// paginateMockRecommendations applies pagination to recommendations.
// Returns the page of recommendations, next page token, and any error.
func paginateMockRecommendations(
	recs []*pbc.Recommendation,
	pageSize int32,
	pageToken string,
) ([]*pbc.Recommendation, string, error) {
	// Use default page size if not specified
	if pageSize <= 0 {
		pageSize = mockDefaultPageSize
	}

	// Decode offset from page token
	var offset int
	if pageToken != "" {
		decoded, err := decodeMockPageToken(pageToken)
		if err != nil {
			return nil, "", err
		}
		offset = decoded
	}

	// Handle offset beyond range - return empty
	if offset >= len(recs) {
		return []*pbc.Recommendation{}, "", nil
	}

	// Calculate end index
	end := offset + int(pageSize)
	if end > len(recs) {
		end = len(recs)
	}

	// Get page slice
	page := recs[offset:end]

	// Generate next token if more results exist
	var nextToken string
	if end < len(recs) {
		nextToken = encodeMockPageToken(end)
	}

	return page, nextToken, nil
}

// encodeMockPageToken encodes an offset as a base64 page token.
func encodeMockPageToken(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}

// decodeMockPageToken decodes a base64 page token to an offset.
func decodeMockPageToken(token string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, fmt.Errorf("invalid page token encoding: %w", err)
	}
	offset, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0, fmt.Errorf("invalid page token value: %w", err)
	}
	if offset < 0 {
		return 0, errors.New("invalid page token: negative offset")
	}
	return offset, nil
}

// EstimateCost returns mock cost estimate for a resource before deployment.
func (m *MockPlugin) EstimateCost(
	_ context.Context,
	req *pbc.EstimateCostRequest,
) (*pbc.EstimateCostResponse, error) {
	if m.EstimateCostDelay > 0 {
		time.Sleep(m.EstimateCostDelay)
	}

	if m.ShouldErrorOnEstimateCost {
		return nil, status.Error(codes.Unavailable, "mock error: pricing source unavailable")
	}

	resourceType := req.GetResourceType()

	// Validate resource type format (FR-003)
	if resourceType == "" {
		return nil, status.Error(codes.InvalidArgument, "resource_type is required")
	}

	provider, module, resource, _, err := parseResourceType(resourceType)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"resource_type must follow provider:module/resource:Type format, got: %s",
				resourceType,
			),
		)
	}

	// Check if provider is supported (FR-008)
	providerSupported := false
	for _, supportedProvider := range m.SupportedProviders {
		if provider == supportedProvider {
			providerSupported = true
			break
		}
	}

	if !providerSupported {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("resource type %s is not supported by this plugin", resourceType),
		)
	}

	// Check if resource type is supported for this provider
	supportedResources, exists := m.SupportedResources[provider]
	if !exists {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("resource type %s is not supported by this plugin", resourceType),
		)
	}

	// For mock purposes, try to match against module, resource, or module_resource combinations
	// This handles formats like "aws:ec2/instance:Instance" where "ec2" is the module
	simpleResourceType := module
	resourceSupported := false
	for _, supportedResource := range supportedResources {
		if module == supportedResource || resource == supportedResource || module+"_"+resource == supportedResource {
			resourceSupported = true
			simpleResourceType = supportedResource
			break
		}
	}

	if !resourceSupported {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("resource type %s is not supported by this plugin", resourceType),
		)
	}

	// Calculate monthly cost based on resource type
	multiplier := getRateMultiplier(simpleResourceType)
	hourlyRate := m.BaseHourlyRate * multiplier
	monthlyHours := float64(HoursPerDay * daysPerMonth)
	monthlyCost := hourlyRate * monthlyHours

	return &pbc.EstimateCostResponse{
		Currency:    m.Currency,
		CostMonthly: monthlyCost,
	}, nil
}

// =============================================================================
// FallbackHint Helper Functions for Benchmarks
// =============================================================================

// NewActualCostResponseWithHint creates a GetActualCostResponse with results and a fallback hint.
// This is a convenience function for benchmark tests.
func NewActualCostResponseWithHint(
	results []*pbc.ActualCostResult,
	hint pbc.FallbackHint,
) *pbc.GetActualCostResponse {
	return &pbc.GetActualCostResponse{
		Results:      results,
		FallbackHint: hint,
	}
}

// NOTE: ValidateActualCostResponseHelper was removed. Use pluginsdk.ValidateActualCostResponse
// instead, which is the canonical implementation for validating GetActualCostResponse messages.
// Import: "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
