// Package testing provides test utilities for PulumiCost plugin development.
// It includes mock plugin implementations, test harnesses, and conformance testing
// utilities for validating plugin behavior against the CostSource gRPC service spec.
package testing

import (
	"context"
	"fmt"
	"math"
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
	}
}

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
		Results: results,
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
