// multi_provider_test.go provides tests for the multi-provider extraction
// functions demonstrating proper testing patterns for provider-specific logic.
package main

import (
	"testing"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
)

// TestExtractInstanceFamily tests AWS instance family extraction.
func TestExtractInstanceFamily(t *testing.T) {
	tests := []struct {
		name         string
		instanceType string
		expected     string
	}{
		{"t3.micro", "t3.micro", "t3"},
		{"t3.medium", "t3.medium", "t3"},
		{"m5.large", "m5.large", "m5"},
		{"m5.xlarge", "m5.xlarge", "m5"},
		{"c5n.18xlarge", "c5n.18xlarge", "c5n"},
		{"r6g.metal", "r6g.metal", "r6g"},
		{"no dot returns full string", "nodot", "nodot"},
		{"empty string", "", ""},
		{"just family", "t3", "t3"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractInstanceFamily(tc.instanceType)
			if result != tc.expected {
				t.Errorf("extractInstanceFamily(%q) = %q, want %q", tc.instanceType, result, tc.expected)
			}
		})
	}
}

// TestExtractAzureFamily tests Azure VM family extraction.
func TestExtractAzureFamily(t *testing.T) {
	tests := []struct {
		name     string
		vmSize   string
		expected string
	}{
		{"Standard_D2s_v3", "Standard_D2s_v3", "D"},
		{"Standard_D4s_v3", "Standard_D4s_v3", "D"},
		{"Standard_E8_v4", "Standard_E8_v4", "E"},
		{"Standard_E16_v4", "Standard_E16_v4", "E"},
		{"Standard_B2s", "Standard_B2s", "B"},
		{"Standard_NC24", "Standard_NC24", "NC"},
		{"Standard_NV12s_v3", "Standard_NV12s_v3", "NV"},
		{"short string returns as-is", "Short", "Short"},
		{"exactly 9 chars returns as-is", "Standard_", "Standard_"},
		{"empty string", "", ""},
		{"no digits returns rest", "Standard_DC", "DC"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractAzureFamily(tc.vmSize)
			if result != tc.expected {
				t.Errorf("extractAzureFamily(%q) = %q, want %q", tc.vmSize, result, tc.expected)
			}
		})
	}
}

// TestExtractGCPFamily tests GCP machine family extraction.
func TestExtractGCPFamily(t *testing.T) {
	tests := []struct {
		name        string
		machineType string
		expected    string
	}{
		{"n1-standard-4", "n1-standard-4", "n1"},
		{"n1-standard-8", "n1-standard-8", "n1"},
		{"n2-standard-4", "n2-standard-4", "n2"},
		{"e2-medium", "e2-medium", "e2"},
		{"e2-small", "e2-small", "e2"},
		{"c2-standard-60", "c2-standard-60", "c2"},
		{"m2-ultramem-416", "m2-ultramem-416", "m2"},
		{"no dash returns full string", "nodash", "nodash"},
		{"empty string", "", ""},
		{"just family", "n1", "n1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractGCPFamily(tc.machineType)
			if result != tc.expected {
				t.Errorf("extractGCPFamily(%q) = %q, want %q", tc.machineType, result, tc.expected)
			}
		})
	}
}

// TestExtractComputeDetails_AWS tests AWS property extraction.
func TestExtractComputeDetails_AWS(t *testing.T) {
	tests := []struct {
		name     string
		props    map[string]string
		expected ComputeDetails
	}{
		{
			name: "EC2 on-demand instance",
			props: map[string]string{
				"instanceType":      "t3.medium",
				"availabilityZone":  "us-east-1a",
				"instanceLifecycle": "on-demand",
			},
			expected: ComputeDetails{
				Provider:       "aws",
				SKU:            "t3.medium",
				Region:         "us-east-1",
				InstanceFamily: "t3",
				IsSpot:         false,
				IsPreemptible:  false,
			},
		},
		{
			name: "EC2 spot instance",
			props: map[string]string{
				"instanceType":      "m5.xlarge",
				"region":            "us-west-2",
				"instanceLifecycle": "spot",
			},
			expected: ComputeDetails{
				Provider:       "aws",
				SKU:            "m5.xlarge",
				Region:         "us-west-2",
				InstanceFamily: "m5",
				IsSpot:         true,
				IsPreemptible:  false,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractComputeDetails("aws", tc.props)
			if result.Provider != tc.expected.Provider {
				t.Errorf("Provider = %q, want %q", result.Provider, tc.expected.Provider)
			}
			if result.SKU != tc.expected.SKU {
				t.Errorf("SKU = %q, want %q", result.SKU, tc.expected.SKU)
			}
			if result.InstanceFamily != tc.expected.InstanceFamily {
				t.Errorf("InstanceFamily = %q, want %q", result.InstanceFamily, tc.expected.InstanceFamily)
			}
			if result.IsSpot != tc.expected.IsSpot {
				t.Errorf("IsSpot = %v, want %v", result.IsSpot, tc.expected.IsSpot)
			}
		})
	}
}

// TestExtractComputeDetails_Azure tests Azure property extraction.
func TestExtractComputeDetails_Azure(t *testing.T) {
	tests := []struct {
		name     string
		props    map[string]string
		expected ComputeDetails
	}{
		{
			name: "Azure VM regular priority",
			props: map[string]string{
				"vmSize":   "Standard_D2s_v3",
				"location": "eastus",
				"priority": "Regular",
			},
			expected: ComputeDetails{
				Provider:       "azure",
				SKU:            "Standard_D2s_v3",
				Region:         "eastus",
				InstanceFamily: "D",
				IsSpot:         false,
			},
		},
		{
			name: "Azure VM spot priority",
			props: map[string]string{
				"vmSize":   "Standard_E8_v4",
				"location": "westeurope",
				"priority": "Spot",
			},
			expected: ComputeDetails{
				Provider:       "azure",
				SKU:            "Standard_E8_v4",
				Region:         "westeurope",
				InstanceFamily: "E",
				IsSpot:         true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractComputeDetails("azure", tc.props)
			if result.Provider != tc.expected.Provider {
				t.Errorf("Provider = %q, want %q", result.Provider, tc.expected.Provider)
			}
			if result.SKU != tc.expected.SKU {
				t.Errorf("SKU = %q, want %q", result.SKU, tc.expected.SKU)
			}
			if result.InstanceFamily != tc.expected.InstanceFamily {
				t.Errorf("InstanceFamily = %q, want %q", result.InstanceFamily, tc.expected.InstanceFamily)
			}
			if result.IsSpot != tc.expected.IsSpot {
				t.Errorf("IsSpot = %v, want %v", result.IsSpot, tc.expected.IsSpot)
			}
		})
	}
}

// TestExtractComputeDetails_GCP tests GCP property extraction.
func TestExtractComputeDetails_GCP(t *testing.T) {
	tests := []struct {
		name     string
		props    map[string]string
		expected ComputeDetails
	}{
		{
			name: "GCP standard instance",
			props: map[string]string{
				"machineType":            "n1-standard-4",
				"zone":                   "us-central1-a",
				"scheduling.preemptible": "false",
			},
			expected: ComputeDetails{
				Provider:       "gcp",
				SKU:            "n1-standard-4",
				Region:         "us-central1",
				InstanceFamily: "n1",
				IsPreemptible:  false,
			},
		},
		{
			name: "GCP preemptible instance",
			props: map[string]string{
				"machineType":            "e2-medium",
				"zone":                   "europe-west1-b",
				"scheduling.preemptible": "true",
			},
			expected: ComputeDetails{
				Provider:       "gcp",
				SKU:            "e2-medium",
				Region:         "europe-west1",
				InstanceFamily: "e2",
				IsPreemptible:  true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractComputeDetails("gcp", tc.props)
			if result.Provider != tc.expected.Provider {
				t.Errorf("Provider = %q, want %q", result.Provider, tc.expected.Provider)
			}
			if result.SKU != tc.expected.SKU {
				t.Errorf("SKU = %q, want %q", result.SKU, tc.expected.SKU)
			}
			if result.InstanceFamily != tc.expected.InstanceFamily {
				t.Errorf("InstanceFamily = %q, want %q", result.InstanceFamily, tc.expected.InstanceFamily)
			}
			if result.IsPreemptible != tc.expected.IsPreemptible {
				t.Errorf("IsPreemptible = %v, want %v", result.IsPreemptible, tc.expected.IsPreemptible)
			}
		})
	}
}

// TestExtractComputeDetails_CustomProvider tests generic provider extraction.
func TestExtractComputeDetails_CustomProvider(t *testing.T) {
	props := map[string]string{
		"type":   "custom-large",
		"region": "custom-region",
	}

	result := extractComputeDetails("custom", props)

	if result.Provider != "custom" {
		t.Errorf("Provider = %q, want %q", result.Provider, "custom")
	}
	if result.SKU != "custom-large" {
		t.Errorf("SKU = %q, want %q", result.SKU, "custom-large")
	}
	if result.Region != "custom-region" {
		t.Errorf("Region = %q, want %q", result.Region, "custom-region")
	}
}

// TestMultiProviderMatcher tests the resource matcher configuration.
func TestMultiProviderMatcher(t *testing.T) {
	matcher := NewMultiProviderMatcher()

	tests := []struct {
		name         string
		provider     string
		resourceType string
		expected     bool
	}{
		// Supported resources
		{"AWS EC2 supported", "aws", "aws:ec2/instance:Instance", true},
		{"AWS RDS supported", "aws", "aws:rds/instance:Instance", true},
		{"AWS Lambda supported", "aws", "aws:lambda/function:Function", true},
		{"Azure VM supported", "azure", "azure:compute/virtualMachine:VirtualMachine", true},
		{"Azure SQL supported", "azure", "azure:sql/database:Database", true},
		{"GCP Compute supported", "gcp", "gcp:compute/instance:Instance", true},
		{"GCP SQL supported", "gcp", "gcp:sql/databaseInstance:DatabaseInstance", true},

		// Unsupported resources
		{"AWS S3 not supported", "aws", "aws:s3/bucket:Bucket", false},
		{"Azure Storage not supported", "azure", "azure:storage/account:Account", false},
		{"GCP BigTable not supported", "gcp", "gcp:bigtable/instance:Instance", false},
		{"Custom provider not supported", "custom", "custom:resource/type:Type", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			desc := pluginsdk.NewResourceDescriptor(tc.provider, tc.resourceType)
			result := matcher.matcher.Supports(desc)
			if result != tc.expected {
				t.Errorf("Supports(%q, %q) = %v, want %v",
					tc.provider, tc.resourceType, result, tc.expected)
			}
		})
	}
}

// TestAzureStandardPrefix verifies the constant is correct.
func TestAzureStandardPrefix(t *testing.T) {
	expectedPrefix := "Standard_"
	if azureStandardPrefix != expectedPrefix {
		t.Errorf("azureStandardPrefix = %q, want %q",
			azureStandardPrefix, expectedPrefix)
	}
}
