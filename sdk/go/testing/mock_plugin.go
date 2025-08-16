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

// MockPlugin provides a configurable mock implementation of CostSourceServiceServer
type MockPlugin struct {
	pbc.UnimplementedCostSourceServiceServer

	// Configuration
	PluginName         string
	SupportedProviders []string
	SupportedResources map[string][]string // provider -> resource types
	
	// Behavior configuration
	ShouldErrorOnName            bool
	ShouldErrorOnSupports        bool
	ShouldErrorOnActualCost      bool
	ShouldErrorOnProjectedCost   bool
	ShouldErrorOnPricingSpec     bool
	
	// Response delays for testing timeouts
	NameDelay            time.Duration
	SupportsDelay        time.Duration
	ActualCostDelay      time.Duration
	ProjectedCostDelay   time.Duration
	PricingSpecDelay     time.Duration

	// Data generation configuration
	ActualCostDataPoints int
	BaseHourlyRate       float64
	Currency             string
}

// NewMockPlugin creates a new mock plugin with default configuration
func NewMockPlugin() *MockPlugin {
	return &MockPlugin{
		PluginName: "mock-test-plugin",
		SupportedProviders: []string{"aws", "azure", "gcp", "kubernetes"},
		SupportedResources: map[string][]string{
			"aws":        {"ec2", "s3", "lambda", "rds"},
			"azure":      {"vm", "blob_storage", "sql_database"},
			"gcp":        {"compute_engine", "cloud_storage", "cloud_functions"},
			"kubernetes": {"namespace", "pod", "service"},
		},
		ActualCostDataPoints: 24, // 24 hours of hourly data
		BaseHourlyRate:       0.05,
		Currency:             "USD",
	}
}

// ConfigurableErrorMockPlugin creates a mock plugin that can be configured to return errors
func ConfigurableErrorMockPlugin() *MockPlugin {
	plugin := NewMockPlugin()
	plugin.PluginName = "error-test-plugin"
	return plugin
}

// SlowMockPlugin creates a mock plugin with artificial delays for timeout testing
func SlowMockPlugin() *MockPlugin {
	plugin := NewMockPlugin()
	plugin.PluginName = "slow-test-plugin"
	plugin.NameDelay = 100 * time.Millisecond
	plugin.SupportsDelay = 200 * time.Millisecond
	plugin.ActualCostDelay = 500 * time.Millisecond
	plugin.ProjectedCostDelay = 300 * time.Millisecond
	plugin.PricingSpecDelay = 250 * time.Millisecond
	return plugin
}

// Name returns the plugin name
func (m *MockPlugin) Name(ctx context.Context, req *pbc.NameRequest) (*pbc.NameResponse, error) {
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

// Supports checks if a resource type is supported by this mock plugin
func (m *MockPlugin) Supports(ctx context.Context, req *pbc.SupportsRequest) (*pbc.SupportsResponse, error) {
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

// GetActualCost returns mock historical cost data
func (m *MockPlugin) GetActualCost(ctx context.Context, req *pbc.GetActualCostRequest) (*pbc.GetActualCostResponse, error) {
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
	
	for i := 0; i < dataPoints; i++ {
		timestamp := startTime.Add(time.Duration(i) * time.Hour)
		if timestamp.After(endTime) {
			break
		}

		// Generate some variation in cost data
		costVariation := 0.8 + (0.4 * float64(i%10) / 10) // Varies between 0.8x and 1.2x base rate
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

// GetProjectedCost returns mock projected cost data
func (m *MockPlugin) GetProjectedCost(ctx context.Context, req *pbc.GetProjectedCostRequest) (*pbc.GetProjectedCostResponse, error) {
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

	// Calculate cost based on resource type
	multiplier := 1.0
	switch resource.GetResourceType() {
	case "ec2", "vm", "compute_engine":
		multiplier = 2.0 // Compute is more expensive
	case "s3", "blob_storage", "cloud_storage":
		multiplier = 0.1 // Storage is cheaper
	case "lambda", "cloud_functions":
		multiplier = 0.001 // Serverless is very cheap per unit
	}

	unitPrice := m.BaseHourlyRate * multiplier
	costPerMonth := unitPrice * 24 * 30 // 30 days

	billingDetail := fmt.Sprintf("mock-%s-rate", strings.ToLower(resource.GetProvider()))

	return &pbc.GetProjectedCostResponse{
		UnitPrice:     unitPrice,
		Currency:      m.Currency,
		CostPerMonth:  costPerMonth,
		BillingDetail: billingDetail,
	}, nil
}

// GetPricingSpec returns mock pricing specification
func (m *MockPlugin) GetPricingSpec(ctx context.Context, req *pbc.GetPricingSpecRequest) (*pbc.GetPricingSpecResponse, error) {
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

	// Generate billing mode based on resource type
	billingMode := "per_hour"
	switch resource.GetResourceType() {
	case "s3", "blob_storage", "cloud_storage":
		billingMode = "per_gb_month"
	case "lambda", "cloud_functions":
		billingMode = "per_invocation"
	case "namespace":
		billingMode = "per_cpu_hour"
	case "sql_database":
		billingMode = "per_dtu"
	}

	// Calculate rate based on resource type
	multiplier := 1.0
	switch resource.GetResourceType() {
	case "ec2", "vm", "compute_engine":
		multiplier = 2.0
	case "s3", "blob_storage", "cloud_storage":
		multiplier = 0.1
	case "lambda", "cloud_functions":
		multiplier = 0.001
	case "namespace":
		multiplier = 1.5
	case "sql_database":
		multiplier = 3.0
	}

	ratePerUnit := m.BaseHourlyRate * multiplier

	// Generate metric hints based on resource type
	metricHints := []*pbc.UsageMetricHint{}
	switch resource.GetResourceType() {
	case "ec2", "vm", "compute_engine":
		metricHints = []*pbc.UsageMetricHint{
			{Metric: "vcpu_hours", Unit: "hour"},
			{Metric: "memory_gb_hours", Unit: "hour"},
		}
	case "s3", "blob_storage", "cloud_storage":
		metricHints = []*pbc.UsageMetricHint{
			{Metric: "storage_gb", Unit: "GB"},
			{Metric: "requests", Unit: "count"},
		}
	case "lambda", "cloud_functions":
		metricHints = []*pbc.UsageMetricHint{
			{Metric: "invocations", Unit: "count"},
			{Metric: "duration_ms", Unit: "millisecond"},
		}
	case "namespace":
		metricHints = []*pbc.UsageMetricHint{
			{Metric: "cpu_cores", Unit: "hour"},
			{Metric: "memory_gb", Unit: "hour"},
		}
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
			"test_mode":          "true",
			"data_source":        "synthetic",
		},
		Source: m.PluginName,
	}

	return &pbc.GetPricingSpecResponse{
		Spec: spec,
	}, nil
}