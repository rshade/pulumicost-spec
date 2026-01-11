// Package quickstart_test validates that all code samples in quickstart.md compile correctly.
package quickstart_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/rshade/finfocus-spec/sdk/go/pricing"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// TestQuickstartLinearGrowth validates the linear growth example compiles.
func TestQuickstartLinearGrowth(t *testing.T) {
	// Create a resource with 10% monthly linear growth
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.medium",
		Region:       "us-east-1",
		GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
		GrowthRate:   proto.Float64(0.10), // 10% growth
	}

	// Request projected costs
	req := &pbc.GetProjectedCostRequest{
		Resource: resource,
	}

	// Verify fields are set correctly
	if resource.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_LINEAR {
		t.Errorf("Expected LINEAR growth type")
	}
	if resource.GetGrowthRate() != 0.10 {
		t.Errorf("Expected 0.10 growth rate, got %v", resource.GetGrowthRate())
	}
	if req.GetResource() != resource {
		t.Errorf("Request resource not set")
	}
}

// TestQuickstartExponentialGrowth validates the exponential growth example compiles.
func TestQuickstartExponentialGrowth(t *testing.T) {
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "s3",
		Sku:          "standard",
		Region:       "us-east-1",
		GrowthType:   pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
		GrowthRate:   proto.Float64(0.05), // 5% compounding
	}

	if resource.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL {
		t.Errorf("Expected EXPONENTIAL growth type")
	}
}

// TestQuickstartDecliningUsage validates the declining usage example compiles.
func TestQuickstartDecliningUsage(t *testing.T) {
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.large",
		Region:       "us-east-1",
		GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
		GrowthRate:   proto.Float64(-0.10), // 10% decline
	}

	if resource.GetGrowthRate() != -0.10 {
		t.Errorf("Expected -0.10 growth rate")
	}
}

// TestQuickstartRequestOverride validates the request override example compiles.
func TestQuickstartRequestOverride(t *testing.T) {
	// Resource has 10% linear growth as default
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.medium",
		Region:       "us-east-1",
		GrowthType:   pbc.GrowthType_GROWTH_TYPE_LINEAR,
		GrowthRate:   proto.Float64(0.10),
	}

	// But request uses 5% exponential (override)
	req := &pbc.GetProjectedCostRequest{
		Resource:   resource,
		GrowthType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
		GrowthRate: proto.Float64(0.05),
	}

	if req.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL {
		t.Errorf("Expected EXPONENTIAL override")
	}
	if req.GetGrowthRate() != 0.05 {
		t.Errorf("Expected 0.05 rate override")
	}
}

// TestQuickstartValidation validates the validation example compiles.
func TestQuickstartValidation(t *testing.T) {
	// Validate before sending request
	err := pricing.ValidateGrowthParams(
		pbc.GrowthType_GROWTH_TYPE_LINEAR,
		proto.Float64(0.10),
	)
	if err != nil {
		t.Errorf("Validation should pass: %v", err)
	}
}

// TestQuickstartScalingScenarios validates the scaling scenarios example compiles.
func TestQuickstartScalingScenarios(t *testing.T) {
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.medium",
		Region:       "us-east-1",
	}

	// Aggressive growth scenario
	aggressive := &pbc.GetProjectedCostRequest{
		Resource:   resource,
		GrowthType: pbc.GrowthType_GROWTH_TYPE_EXPONENTIAL,
		GrowthRate: proto.Float64(0.20), // 20% monthly
	}

	// Conservative scenario
	conservative := &pbc.GetProjectedCostRequest{
		Resource:   resource,
		GrowthType: pbc.GrowthType_GROWTH_TYPE_LINEAR,
		GrowthRate: proto.Float64(0.05), // 5% monthly
	}

	if aggressive.GetGrowthRate() != 0.20 {
		t.Errorf("Aggressive rate should be 0.20")
	}
	if conservative.GetGrowthRate() != 0.05 {
		t.Errorf("Conservative rate should be 0.05")
	}
}

// TestQuickstartErrorHandling validates the error handling example compiles.
func TestQuickstartErrorHandling(t *testing.T) {
	// This just validates the code compiles - we don't have a real client
	err := status.Error(codes.InvalidArgument, "test error")

	// Create a logger for this test
	logger := zerolog.Nop()

	st, ok := status.FromError(err)
	if ok && st.Code() == codes.InvalidArgument {
		// Handle validation error (e.g., missing growth_rate)
		logger.Info().Str("message", st.Message()).Msg("Invalid request")
	}

	_ = context.Background() // Validate context import works
	t.Log("Error handling code compiles correctly")
}

// TestQuickstartBackwardCompatibility validates backward compatibility example compiles.
func TestQuickstartBackwardCompatibility(t *testing.T) {
	// This still works - no growth applied
	resource := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.medium",
		Region:       "us-east-1",
		// No growth fields set - defaults to GROWTH_TYPE_UNSPECIFIED (no growth)
	}

	if resource.GetGrowthType() != pbc.GrowthType_GROWTH_TYPE_UNSPECIFIED {
		t.Errorf("Expected UNSPECIFIED as default")
	}
	if resource.GrowthRate != nil {
		t.Errorf("Expected nil growth rate as default")
	}
}
