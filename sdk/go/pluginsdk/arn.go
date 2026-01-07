package pluginsdk

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	AWSARNPrefix   = "arn:aws:"
	AzureARNPrefix = "/subscriptions/"
	GCPARNPrefix   = "//"
	// KubernetesFormat describes the expected format for Kubernetes resource IDs.
	// Unlike the prefix constants above, this is a format template used for documentation.
	KubernetesFormat = "{cluster}/{namespace}/{kind}/{name}"
)

// k8sSegmentRegex matches valid Kubernetes segment names (lowercase alphanumeric, hyphens, dots).
// Must start/end with alphanumeric.
var k8sSegmentRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$`)

// DetectARNProvider returns the cloud provider inferred from ARN format.
// DetectARNProvider infers the cloud/provider from an ARN-like string.
// It returns "aws", "azure", "gcp", "kubernetes", or an empty string if the format is unrecognized.
// Kubernetes detection is heuristic: the string must contain at least three '/' separators, must not start with '/'
// or "arn:", and every segment must match the Kubernetes segment validation pattern.
func DetectARNProvider(arn string) string {
	switch {
	case strings.HasPrefix(arn, AWSARNPrefix):
		return "aws"
	case strings.HasPrefix(arn, AzureARNPrefix):
		return "azure"
	case strings.HasPrefix(arn, GCPARNPrefix):
		return "gcp"
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
					return ""
				}
			}
			return "kubernetes"
		}
		return ""
	}
}

// ValidateARNConsistency checks if the ARN format matches the expected provider.
// ValidateARNConsistency verifies that the provider inferred from arn matches expectedProvider.
// It returns nil if the ARN is unrecognized or its detected provider equals expectedProvider; otherwise it returns an error describing the mismatch.
func ValidateARNConsistency(arn, expectedProvider string) error {
	detected := DetectARNProvider(arn)
	if detected == "" {
		return nil
	}
	if detected != expectedProvider {
		return fmt.Errorf("ARN format %q detected as %q but expected %q", arn, detected, expectedProvider)
	}
	return nil
}