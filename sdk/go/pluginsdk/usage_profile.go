// Package pluginsdk provides helpers for handling UsageProfile enum values.
//
// UsageProfile represents the intended workload context for cost estimation.
// The Core signals this intent via GetProjectedCostRequest and GetRecommendationsRequest,
// and plugins use it to apply profile-appropriate defaults.
//
// Usage:
//
//	profile := req.GetUsageProfile()
//	if !pluginsdk.IsValidUsageProfile(profile) {
//	    profile = pluginsdk.NormalizeUsageProfile(profile)
//	}
//	switch profile {
//	case pbc.UsageProfile_USAGE_PROFILE_DEV:
//	    // Apply development defaults (160hr/month)
//	case pbc.UsageProfile_USAGE_PROFILE_PROD:
//	    // Apply production defaults (730hr/month)
//	case pbc.UsageProfile_USAGE_PROFILE_BURST:
//	    // Apply burst defaults (high data transfer)
//	default:
//	    // UNSPECIFIED - apply plugin defaults
//	}
package pluginsdk

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

// allUsageProfiles is a package-level slice containing all valid UsageProfile values.
// This is allocated once at package initialization for zero-allocation validation.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation validation
var allUsageProfiles = []pbc.UsageProfile{
	pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
	pbc.UsageProfile_USAGE_PROFILE_PROD,
	pbc.UsageProfile_USAGE_PROFILE_DEV,
	pbc.UsageProfile_USAGE_PROFILE_BURST,
}

// usageProfileStringMap provides O(1) lookup for enum-to-string conversion.
// Map access is faster than switch statements for string conversion.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation lookup
var usageProfileStringMap = map[pbc.UsageProfile]string{
	pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED: "unspecified",
	pbc.UsageProfile_USAGE_PROFILE_PROD:        "prod",
	pbc.UsageProfile_USAGE_PROFILE_DEV:         "dev",
	pbc.UsageProfile_USAGE_PROFILE_BURST:       "burst",
}

// usageProfileParseMap provides O(1) lookup for string-to-enum parsing.
// Supports lowercase, uppercase, and common variants.
//
//nolint:gochecknoglobals // Intentional optimization for zero-allocation lookup
var usageProfileParseMap = map[string]pbc.UsageProfile{
	// Lowercase
	"unspecified": pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
	"prod":        pbc.UsageProfile_USAGE_PROFILE_PROD,
	"production":  pbc.UsageProfile_USAGE_PROFILE_PROD,
	"dev":         pbc.UsageProfile_USAGE_PROFILE_DEV,
	"development": pbc.UsageProfile_USAGE_PROFILE_DEV,
	"burst":       pbc.UsageProfile_USAGE_PROFILE_BURST,
}

// AllUsageProfiles returns all valid UsageProfile enum values.
// This returns the pre-allocated package-level slice for zero-allocation access.
func AllUsageProfiles() []pbc.UsageProfile {
	return allUsageProfiles
}

// IsValidUsageProfile checks if the given profile is a known UsageProfile value.
// Returns true for UNSPECIFIED, PROD, DEV, and BURST.
// Returns false for unknown/future profile values.
//
// Performance: <15 ns/op, 0 allocs/op using package-level slice iteration.
func IsValidUsageProfile(profile pbc.UsageProfile) bool {
	for _, valid := range allUsageProfiles {
		if profile == valid {
			return true
		}
	}
	return false
}

// ParseUsageProfile converts a string to the corresponding UsageProfile enum value.
// Supports case-insensitive matching and common variants:
//   - "dev", "development" → USAGE_PROFILE_DEV
//   - "prod", "production" → USAGE_PROFILE_PROD
//   - "burst" → USAGE_PROFILE_BURST
//   - "unspecified", "" → USAGE_PROFILE_UNSPECIFIED
//
// Returns an error for unrecognized strings.
func ParseUsageProfile(s string) (pbc.UsageProfile, error) {
	if s == "" {
		return pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED, nil
	}

	lower := strings.ToLower(s)
	if profile, ok := usageProfileParseMap[lower]; ok {
		return profile, nil
	}

	return pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED,
		fmt.Errorf("unknown usage profile: %q", s)
}

// UsageProfileString returns a lowercase string representation of the profile.
// For known profiles: "unspecified", "prod", "dev", "burst".
// For unknown profiles: "unknown(<value>)".
//
// This is useful for logging and human-readable output.
func UsageProfileString(profile pbc.UsageProfile) string {
	if str, ok := usageProfileStringMap[profile]; ok {
		return str
	}
	return fmt.Sprintf("unknown(%d)", int32(profile))
}

// NormalizeUsageProfile returns the profile if it's a known value, or UNSPECIFIED
// if the profile is unknown/future. This enables forward compatibility when
// plugins receive profile values from newer spec versions.
//
// Logs a warning at INFO level when normalizing an unknown value.
//
// Example:
//
//	profile := pluginsdk.NormalizeUsageProfile(req.GetUsageProfile())
//	// profile is guaranteed to be a known value
func NormalizeUsageProfile(profile pbc.UsageProfile) pbc.UsageProfile {
	if IsValidUsageProfile(profile) {
		return profile
	}

	log.Warn().
		Int32("usage_profile", int32(profile)).
		Msg("Unknown usage profile, treating as UNSPECIFIED")

	return pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED
}

// Constants for default monthly hours by profile.
const (
	// HoursProd represents 24/7 operation (730 hours = 365 days / 12 months * 24 hours).
	HoursProd = 730
	// HoursDev represents ~8 hours/day, 5 days/week, ~20 business days.
	HoursDev = 160
	// HoursBurst represents plugin discretion for batch/load-test scenarios.
	HoursBurst = 200
)

// DefaultMonthlyHours returns the default monthly usage hours for a profile.
// This is a convenience helper showing typical hour assumptions:
//   - PROD: 730 hours (24/7 operation, 365 days / 12 months * 24 hours)
//   - DEV: 160 hours (~8 hours/day, 5 days/week, ~20 business days)
//   - BURST: 200 hours (plugin discretion for batch/load-test scenarios)
//   - UNSPECIFIED: 730 hours (defaults to production assumptions)
//
// Plugins have discretion to use different values based on their resource types.
func DefaultMonthlyHours(profile pbc.UsageProfile) float64 {
	switch profile {
	case pbc.UsageProfile_USAGE_PROFILE_PROD:
		return HoursProd
	case pbc.UsageProfile_USAGE_PROFILE_DEV:
		return HoursDev
	case pbc.UsageProfile_USAGE_PROFILE_BURST:
		return HoursBurst
	case pbc.UsageProfile_USAGE_PROFILE_UNSPECIFIED:
		return HoursProd
	default:
		// Unknown profiles - use production assumptions
		return HoursProd
	}
}
