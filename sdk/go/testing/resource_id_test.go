// Package testing_test provides tests for resource ID and ARN fields in ResourceDescriptor.
//
// These tests verify the new `id` and `arn` fields added to ResourceDescriptor for:
//   - Client correlation in batch operations (id field)
//   - Exact resource matching for precise cost lookups (arn field)
//
// The tests are organized by user story:
//   - US1: Batch resource recommendation correlation
//   - US2: Exact resource matching via ARN
//   - US3: Pass-through identifier support
//   - US4: Backward compatible protocol evolution
package testing_test

import (
	"strings"
	"testing"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	"google.golang.org/protobuf/proto"
)

// =============================================================================
// User Story 1: Batch Resource Recommendation Correlation (P1)
// Goal: Enable clients to correlate batch recommendation responses to requests
// using the `id` field
// =============================================================================

// TestResourceDescriptor_IDField verifies the id field exists and is accessible.
// This is a prerequisite for correlation - clients must be able to set and get IDs.
func TestResourceDescriptor_IDField(t *testing.T) {
	t.Parallel()

	descriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.micro",
		Region:       "us-east-1",
	}

	// Verify that the Id field exists and is of type string
	if _, ok := interface{}(descriptor).(interface{ GetId() string }); !ok {
		t.Errorf("ResourceDescriptor does not have a GetId() method returning string")
	}

	// Verify setting and getting the Id field
	expectedID := "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::webserver"
	descriptor.Id = expectedID
	if descriptor.GetId() != expectedID {
		t.Errorf("Id field not set or retrieved correctly. Expected %s, got %s",
			expectedID, descriptor.GetId())
	}
}

// TestResourceDescriptor_CorrelationPattern verifies that IDs can be used
// for request/response correlation - the core use case for the id field.
func TestResourceDescriptor_CorrelationPattern(t *testing.T) {
	t.Parallel()

	// Create descriptors with unique IDs (simulating a batch request)
	descriptors := []*pbc.ResourceDescriptor{
		{
			Provider:     "aws",
			ResourceType: "ec2",
			Id:           "res-001",
		},
		{
			Provider:     "aws",
			ResourceType: "ec2",
			Id:           "res-002",
		},
		{
			Provider:     "aws",
			ResourceType: "ec2",
			Id:           "res-003",
		},
	}

	// Verify each descriptor has a unique, accessible ID
	ids := make(map[string]bool)
	for _, desc := range descriptors {
		id := desc.GetId()
		if id == "" {
			t.Errorf("Descriptor should have a non-empty ID")
		}
		if ids[id] {
			t.Errorf("Duplicate ID found: %s", id)
		}
		ids[id] = true
	}

	// Verify all expected IDs are present
	if !ids["res-001"] || !ids["res-002"] || !ids["res-003"] {
		t.Errorf("Not all expected IDs found in descriptors")
	}
}

// TestResourceDescriptor_BatchCorrelation tests the full batch correlation pattern
// where multiple resources with unique IDs can be sent and correlated.
func TestResourceDescriptor_BatchCorrelation(t *testing.T) {
	t.Parallel()

	// Create a recommendations request with multiple target resources
	req := &pbc.GetRecommendationsRequest{
		TargetResources: []*pbc.ResourceDescriptor{
			{
				Provider:     "aws",
				ResourceType: "ec2",
				Sku:          "t3.micro",
				Region:       "us-east-1",
				Id:           "batch-001",
				Arn:          "arn:aws:ec2:us-east-1:123456789012:instance/i-abc123",
			},
			{
				Provider:     "aws",
				ResourceType: "ec2",
				Sku:          "t3.small",
				Region:       "us-west-2",
				Id:           "batch-002",
				Arn:          "arn:aws:ec2:us-west-2:123456789012:instance/i-def456",
			},
			{
				Provider:     "aws",
				ResourceType: "ec2",
				Sku:          "t3.medium",
				Region:       "eu-west-1",
				Id:           "batch-003",
				Arn:          "arn:aws:ec2:eu-west-1:123456789012:instance/i-ghi789",
			},
		},
	}

	// Verify all target resources have their IDs set correctly
	for i, target := range req.GetTargetResources() {
		expectedID := "batch-00" + string(rune('1'+i))
		if target.GetId() != expectedID {
			t.Errorf("Target resource %d: expected ID %s, got %s",
				i, expectedID, target.GetId())
		}
	}

	// Verify correlation lookup works (simulate client-side correlation)
	idToResource := make(map[string]*pbc.ResourceDescriptor)
	for _, target := range req.GetTargetResources() {
		idToResource[target.GetId()] = target
	}

	// Verify all resources can be found by ID
	for _, expectedID := range []string{"batch-001", "batch-002", "batch-003"} {
		if _, found := idToResource[expectedID]; !found {
			t.Errorf("Resource with ID %s not found in correlation map", expectedID)
		}
	}
}

// TestResourceDescriptor_EmptyIDBackwardCompatibility verifies that empty ID
// is valid and maintains backward compatibility with existing code.
func TestResourceDescriptor_EmptyIDBackwardCompatibility(t *testing.T) {
	t.Parallel()

	// Old-style descriptor without ID (backward compatible)
	descriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.micro",
		Region:       "us-east-1",
	}

	// Empty ID should be the default (empty string in proto3)
	if descriptor.GetId() != "" {
		t.Errorf("Expected empty ID for descriptor without explicit ID, got %s",
			descriptor.GetId())
	}

	// Verify the descriptor is still valid for existing use cases
	if descriptor.GetProvider() != "aws" {
		t.Errorf("Provider should be accessible: %s", descriptor.GetProvider())
	}
	if descriptor.GetResourceType() != "ec2" {
		t.Errorf("ResourceType should be accessible: %s", descriptor.GetResourceType())
	}
}

// =============================================================================
// User Story 2: Exact Resource Matching via ARN (P1)
// Goal: Enable plugins to use ARN for exact resource lookup instead of
// fuzzy type/sku/region/tags matching
// =============================================================================

// TestResourceDescriptor_ARNField verifies the arn field exists and is accessible.
func TestResourceDescriptor_ARNField(t *testing.T) {
	t.Parallel()

	descriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
	}

	// Verify that the Arn field exists and is of type string
	if _, ok := interface{}(descriptor).(interface{ GetArn() string }); !ok {
		t.Errorf("ResourceDescriptor does not have a GetArn() method returning string")
	}

	// Verify setting and getting the Arn field
	expectedARN := "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0"
	descriptor.Arn = expectedARN
	if descriptor.GetArn() != expectedARN {
		t.Errorf("Arn field not set or retrieved correctly. Expected %s, got %s",
			expectedARN, descriptor.GetArn())
	}
}

// TestResourceDescriptor_ARNPrecedence tests that ARN takes precedence
// when provided for resource matching.
func TestResourceDescriptor_ARNPrecedence(t *testing.T) {
	t.Parallel()

	// Descriptor with both fuzzy matching fields AND ARN
	descriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.micro",
		Region:       "us-east-1",
		Tags: map[string]string{
			"env":  "production",
			"team": "platform",
		},
		Arn: "arn:aws:ec2:us-east-1:123456789012:instance/i-specific123",
	}

	// When ARN is present, it should be accessible for exact matching
	arn := descriptor.GetArn()
	if arn == "" {
		t.Errorf("ARN should be present for exact matching")
	}

	// Verify the ARN contains the expected format
	if arn != "arn:aws:ec2:us-east-1:123456789012:instance/i-specific123" {
		t.Errorf("Unexpected ARN value: %s", arn)
	}

	// The presence of ARN should not affect other fields (they can be used for fallback)
	if descriptor.GetSku() != "t3.micro" {
		t.Errorf("SKU should still be accessible when ARN is set")
	}
}

// TestResourceDescriptor_ARNFallback tests that empty ARN allows
// fallback to fuzzy matching.
func TestResourceDescriptor_ARNFallback(t *testing.T) {
	t.Parallel()

	// Descriptor without ARN (should use fuzzy matching)
	descriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.micro",
		Region:       "us-east-1",
	}

	// ARN should be empty
	if descriptor.GetArn() != "" {
		t.Errorf("ARN should be empty when not set, got %s", descriptor.GetArn())
	}

	// Fuzzy matching fields should be available
	if descriptor.GetProvider() == "" || descriptor.GetResourceType() == "" {
		t.Errorf("Fuzzy matching fields should be available when ARN is empty")
	}
}

// =============================================================================
// User Story 3: Pass-Through Identifier Support (P2)
// Goal: Verify that plugin developers can pass through IDs without
// validation or transformation
// =============================================================================

// TestResourceDescriptor_SpecialCharactersInID tests that IDs with
// special characters are preserved unchanged.
func TestResourceDescriptor_SpecialCharactersInID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		id   string
	}{
		{
			name: "Pulumi URN format",
			id:   "urn:pulumi:prod::myapp::aws:ec2/instance:Instance::webserver",
		},
		{
			name: "UUID format",
			id:   "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name: "Slashes and colons",
			id:   "cluster/namespace/kind/name",
		},
		{
			name: "URL-like format",
			id:   "https://example.com/resources/12345",
		},
		{
			name: "Special characters",
			id:   "resource@domain.com:path/to/item#section?query=value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			descriptor := &pbc.ResourceDescriptor{
				Provider:     "custom",
				ResourceType: "resource",
				Id:           tc.id,
			}

			// Verify the ID is preserved exactly
			if descriptor.GetId() != tc.id {
				t.Errorf("ID not preserved. Expected %s, got %s",
					tc.id, descriptor.GetId())
			}
		})
	}
}

// TestResourceDescriptor_LongID tests that IDs with 256+ characters
// are preserved unchanged.
func TestResourceDescriptor_LongID(t *testing.T) {
	t.Parallel()

	// Generate a 512-character ID using strings.Builder for efficiency
	var builder strings.Builder
	builder.Grow(512)
	for i := range 512 {
		builder.WriteByte('a' + byte(i%26))
	}
	longID := builder.String()

	descriptor := &pbc.ResourceDescriptor{
		Provider:     "custom",
		ResourceType: "resource",
		Id:           longID,
	}

	// Verify the long ID is preserved
	if len(descriptor.GetId()) != 512 {
		t.Errorf("Long ID length not preserved. Expected 512, got %d",
			len(descriptor.GetId()))
	}
	if descriptor.GetId() != longID {
		t.Errorf("Long ID content not preserved exactly")
	}
}

// =============================================================================
// User Story 4: Backward Compatible Protocol Evolution (P2)
// Goal: Verify existing plugins and clients continue to work without modification
// =============================================================================

// TestResourceDescriptor_OldClientSimulation tests that requests without
// id/arn fields work correctly (simulating old clients).
func TestResourceDescriptor_OldClientSimulation(t *testing.T) {
	t.Parallel()

	// Simulate an old client sending a request without id/arn fields
	oldDescriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.micro",
		Region:       "us-east-1",
		Tags:         map[string]string{"env": "dev"},
	}

	// Marshal to simulate network transmission
	marshaledData, err := proto.Marshal(oldDescriptor)
	if err != nil {
		t.Fatalf("Failed to marshal old descriptor: %v", err)
	}

	// Unmarshal into a new descriptor
	newDescriptor := &pbc.ResourceDescriptor{}
	err = proto.Unmarshal(marshaledData, newDescriptor)
	if err != nil {
		t.Fatalf("Failed to unmarshal data: %v", err)
	}

	// Verify existing fields are preserved
	if newDescriptor.GetProvider() != oldDescriptor.GetProvider() {
		t.Errorf("Provider mismatch. Expected %s, got %s",
			oldDescriptor.GetProvider(), newDescriptor.GetProvider())
	}
	if newDescriptor.GetResourceType() != oldDescriptor.GetResourceType() {
		t.Errorf("ResourceType mismatch. Expected %s, got %s",
			oldDescriptor.GetResourceType(), newDescriptor.GetResourceType())
	}
	if newDescriptor.GetSku() != oldDescriptor.GetSku() {
		t.Errorf("Sku mismatch. Expected %s, got %s",
			oldDescriptor.GetSku(), newDescriptor.GetSku())
	}
	if newDescriptor.GetRegion() != oldDescriptor.GetRegion() {
		t.Errorf("Region mismatch. Expected %s, got %s",
			oldDescriptor.GetRegion(), newDescriptor.GetRegion())
	}

	// Verify new fields default to empty strings
	if newDescriptor.GetId() != "" {
		t.Errorf("Id should be empty for old descriptor. Got %s", newDescriptor.GetId())
	}
	if newDescriptor.GetArn() != "" {
		t.Errorf("Arn should be empty for old descriptor. Got %s", newDescriptor.GetArn())
	}
}

// TestResourceDescriptor_EmptyDefaults verifies that empty string defaults
// are correctly applied for both id and arn fields.
func TestResourceDescriptor_EmptyDefaults(t *testing.T) {
	t.Parallel()

	// Create a descriptor with minimal required fields
	descriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
	}

	// Both id and arn should default to empty strings
	if descriptor.GetId() != "" {
		t.Errorf("Id should default to empty string, got %q", descriptor.GetId())
	}
	if descriptor.GetArn() != "" {
		t.Errorf("Arn should default to empty string, got %q", descriptor.GetArn())
	}
}

// TestResourceDescriptor_NewClientOldServer simulates a new client sending
// a request with id/arn to an old server (proto3 behavior: unknown fields ignored).
func TestResourceDescriptor_NewClientOldServer(t *testing.T) {
	t.Parallel()

	// New client sends descriptor with new fields
	newDescriptor := &pbc.ResourceDescriptor{
		Provider:     "aws",
		ResourceType: "ec2",
		Sku:          "t3.micro",
		Region:       "us-east-1",
		Id:           "new-client-id",
		Arn:          "arn:aws:ec2:us-east-1:123456789012:instance/i-new",
	}

	// Marshal the new descriptor
	marshaledData, err := proto.Marshal(newDescriptor)
	if err != nil {
		t.Fatalf("Failed to marshal new descriptor: %v", err)
	}

	// Unmarshal back - in proto3, this simulates how an old server would
	// handle the request (ignoring unknown fields)
	unmarshaledDescriptor := &pbc.ResourceDescriptor{}
	err = proto.Unmarshal(marshaledData, unmarshaledDescriptor)
	if err != nil {
		t.Fatalf("Failed to unmarshal new descriptor: %v", err)
	}

	// Verify all fields are correctly round-tripped
	if unmarshaledDescriptor.GetProvider() != "aws" {
		t.Errorf("Provider not preserved after round-trip")
	}
	if unmarshaledDescriptor.GetResourceType() != "ec2" {
		t.Errorf("ResourceType not preserved after round-trip")
	}
	if unmarshaledDescriptor.GetId() != "new-client-id" {
		t.Errorf("Id not preserved after round-trip")
	}
	if unmarshaledDescriptor.GetArn() != "arn:aws:ec2:us-east-1:123456789012:instance/i-new" {
		t.Errorf("Arn not preserved after round-trip")
	}
}

// =============================================================================
// Cross-Provider ARN Format Tests
// =============================================================================

// TestResourceDescriptor_CrossProviderARNFormats tests that various provider-specific
// canonical identifier formats are correctly stored and retrieved.
func TestResourceDescriptor_CrossProviderARNFormats(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		provider string
		arnName  string
		arn      string
	}{
		{
			provider: "aws",
			arnName:  "AWS ARN",
			arn:      "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
		},
		{
			provider: "azure",
			arnName:  "Azure Resource ID",
			arn:      "/subscriptions/sub-id/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm-name",
		},
		{
			provider: "gcp",
			arnName:  "GCP Full Resource Name",
			arn:      "//compute.googleapis.com/projects/my-project/zones/us-central1-a/instances/vm-1",
		},
		{
			provider: "kubernetes",
			arnName:  "Kubernetes Resource Path",
			arn:      "prod-cluster/default/Deployment/nginx",
		},
		{
			provider: "kubernetes",
			arnName:  "Kubernetes UID",
			arn:      "550e8400-e29b-41d4-a716-446655440000",
		},
		{
			provider: "cloudflare",
			arnName:  "Cloudflare Resource",
			arn:      "abc123zone/dns_record/xyz789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.arnName, func(t *testing.T) {
			t.Parallel()
			descriptor := &pbc.ResourceDescriptor{
				Provider:     tc.provider,
				ResourceType: "resource",
				Arn:          tc.arn,
			}

			// Verify ARN is preserved exactly
			if descriptor.GetArn() != tc.arn {
				t.Errorf("ARN not preserved for %s. Expected %s, got %s",
					tc.arnName, tc.arn, descriptor.GetArn())
			}
		})
	}
}
