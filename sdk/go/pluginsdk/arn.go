package pluginsdk

import (
	"fmt"
	"regexp"
	"strings"
)

// Provider represents a cloud provider identifier for ARN detection.
type Provider string

// Provider constants for cloud provider identification.
const (
	ProviderAWS        Provider = "aws"
	ProviderAzure      Provider = "azure"
	ProviderGCP        Provider = "gcp"
	ProviderKubernetes Provider = "kubernetes"
	ProviderUnknown    Provider = ""
)

// allProviders is the list of recognized/supported cloud providers used for
// zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allProviders = []Provider{
	ProviderAWS,
	ProviderAzure,
	ProviderGCP,
	ProviderKubernetes,
}

// AllProviders returns all valid Provider values. The returned slice is shared and must not be modified.
func AllProviders() []Provider {
	return allProviders
}

// IsValidProvider returns true if the provider is a known valid value.
func IsValidProvider(p Provider) bool {
	for _, valid := range allProviders {
		if p == valid {
			return true
		}
	}
	return false
}

// String returns the string representation of the Provider.
func (p Provider) String() string {
	return string(p)
}

const (
	// AWSARNPrefix is the prefix for AWS ARN strings.
	AWSARNPrefix = "arn:aws:"
	// AzureARNPrefix is the prefix for Azure resource IDs.
	AzureARNPrefix = "/subscriptions/"
	// GCPResourcePrefix is the prefix for GCP full resource names.
	// GCP full resource names always start with "//" followed by a service endpoint.
	GCPResourcePrefix = "//"
	// GCPDomainSuffix is the domain suffix required for valid GCP resource names.
	GCPDomainSuffix = ".googleapis.com/"
	// KubernetesFormat describes the expected format for Kubernetes resource IDs.
	// Unlike the prefix constants above, this is a format template used for documentation.
	KubernetesFormat = "{cluster}/{namespace}/{kind}/{name}"
)

// k8sSegmentRegex matches valid Kubernetes segment names (lowercase alphanumeric, hyphens, dots).
// Must start/end with alphanumeric.
//
//nolint:gochecknoglobals // precompiled regex for DNS-1123 label validation to avoid repeated compilation
//nolint:nolintlint // gochecknoglobals may not trigger in all linter configurations; directive is kept for consistency
var k8sSegmentRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$`)

// IsGCPResourceName checks if the given string is a valid GCP full resource name.
// GCP full resource names follow the format: //{service}.googleapis.com/{resource-path}.
// For example: //compute.googleapis.com/projects/my-project/zones/us-central1-a/instances/my-instance.
func IsGCPResourceName(s string) bool {
	if !strings.HasPrefix(s, GCPResourcePrefix) {
		return false
	}
	// Check that the string contains .googleapis.com/ after the //
	remainder := s[len(GCPResourcePrefix):]
	return strings.Contains(remainder, GCPDomainSuffix[1:]) // Skip leading dot for Contains check
}

// DetectARNProvider infers the cloud/provider from an ARN-like string.
// It returns ProviderAWS, ProviderAzure, ProviderGCP, ProviderKubernetes, or ProviderUnknown (empty string)
// if the format is unrecognized.
// Kubernetes detection is heuristic: the string must contain at least three '/' separators, must not start with '/'
// or "arn:", and every segment must match the Kubernetes segment validation pattern.
// GCP detection requires the string to start with "//" and contain ".googleapis.com/" to avoid false positives.
func DetectARNProvider(arn string) Provider {
	switch {
	case strings.HasPrefix(arn, AWSARNPrefix):
		return ProviderAWS
	case strings.HasPrefix(arn, AzureARNPrefix):
		return ProviderAzure
	case IsGCPResourceName(arn):
		return ProviderGCP
	default:
		// Kubernetes detection is heuristic-based.
		// We expect a format like "{cluster}/{namespace}/{kind}/{name}" (at least 3 slashes).
		// We also reject strings starting with "/" to avoid confusion with file paths,
		// though Azure IDs (handled above) do start with "/".
		if strings.Count(arn, "/") >= 3 && !strings.HasPrefix(arn, "/") && !strings.HasPrefix(arn, "arn:") {
			// Validate segments to avoid false positives with generic paths
			parts := strings.Split(arn, "/")
			for _, part := range parts {
				if !k8sSegmentRegex.MatchString(part) {
					return ProviderUnknown
				}
			}
			return ProviderKubernetes
		}
		return ProviderUnknown
	}
}

// ValidateARNConsistency verifies that the provider inferred from arn matches expectedProvider.
// It returns nil if the ARN is unrecognized or its detected provider equals expectedProvider;
// otherwise it returns an error describing the mismatch.
func ValidateARNConsistency(arn string, expectedProvider Provider) error {
	detected := DetectARNProvider(arn)
	if detected == ProviderUnknown {
		return nil
	}
	if detected != expectedProvider {
		return fmt.Errorf("ARN format %q detected as %q but expected %q", arn, detected, expectedProvider)
	}
	return nil
}
