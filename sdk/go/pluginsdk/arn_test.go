package pluginsdk_test

import (
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func TestDetectARNProvider(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		provider string
	}{
		{
			name:     "AWS ARN returns aws",
			arn:      "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
			provider: "aws",
		},
		{
			name:     "Azure subscription ID returns azure",
			arn:      "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM",
			provider: "azure",
		},
		{
			name:     "GCP full resource name returns gcp",
			arn:      "//compute.googleapis.com/projects/my-project/zones/us-central1-a/instances/my-instance",
			provider: "gcp",
		},
		{
			name:     "Kubernetes format returns kubernetes",
			arn:      "my-cluster/my-namespace/pod/my-pod",
			provider: "kubernetes",
		},
		{
			name:     "absolute file path returns empty string",
			arn:      "/etc/passwd",
			provider: "",
		},
		{
			name:     "short relative path returns empty string",
			arn:      "foo/bar",
			provider: "",
		},
		{
			name:     "path with uppercase returns empty string",
			arn:      "File/Path/To/Resource",
			provider: "",
		},
		{
			name:     "path with underscores returns empty string",
			arn:      "file/path/to/resource_name",
			provider: "",
		},
		{
			name:     "unrecognized format returns empty string",
			arn:      "custom:identifier:format",
			provider: "",
		},
		{
			name:     "empty string returns empty string",
			arn:      "",
			provider: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pluginsdk.DetectARNProvider(tt.arn)
			if got != tt.provider {
				t.Errorf("DetectARNProvider(%q) = %q, want %q", tt.arn, got, tt.provider)
			}
		})
	}
}

func TestValidateARNConsistency(t *testing.T) {
	tests := []struct {
		name        string
		arn         string
		expected    string
		expectError bool
	}{
		{
			name:        "valid AWS ARN with AWS expected returns nil",
			arn:         "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
			expected:    "aws",
			expectError: false,
		},
		{
			name:        "AWS ARN with Azure expected returns error",
			arn:         "arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0",
			expected:    "azure",
			expectError: true,
		},
		{
			name:        "Azure subscription ID with Azure expected returns nil",
			arn:         "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup",
			expected:    "azure",
			expectError: false,
		},
		{
			name:        "Azure subscription ID with AWS expected returns error",
			arn:         "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup",
			expected:    "aws",
			expectError: true,
		},
		{
			name:        "GCP full resource name with GCP expected returns nil",
			arn:         "//compute.googleapis.com/projects/my-project/zones/us-central1-a/instances/my-instance",
			expected:    "gcp",
			expectError: false,
		},
		{
			name:        "Kubernetes format with Kubernetes expected returns nil",
			arn:         "my-cluster/my-namespace/pod/my-pod",
			expected:    "kubernetes",
			expectError: false,
		},
		{
			name:        "unrecognized ARN format returns nil",
			arn:         "custom:identifier:format",
			expected:    "aws",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pluginsdk.ValidateARNConsistency(tt.arn, tt.expected)
			if tt.expectError && err == nil {
				t.Errorf("ValidateARNConsistency(%q, %q) expected error, got nil", tt.arn, tt.expected)
			}
			if !tt.expectError && err != nil {
				t.Errorf("ValidateARNConsistency(%q, %q) unexpected error: %v", tt.arn, tt.expected, err)
			}
		})
	}
}
