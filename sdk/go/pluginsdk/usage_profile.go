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
// AllUsageProfiles returns the pre-allocated slice of all valid UsageProfile values
// (UNSPECIFIED, PROD, DEV, BURST) for zero-allocation access. The returned slice must
// not be modified.
func AllUsageProfiles() []pbc.UsageProfile {
	return allUsageProfiles
}

// IsValidUsageProfile checks if the given profile is a known UsageProfile value.
// Returns true for UNSPECIFIED, PROD, DEV, and BURST.
// Returns false for unknown/future profile values.
//
// IsValidUsageProfile reports whether profile is one of the known UsageProfile values
// (UNSPECIFIED, PROD, DEV, BURST). It performs a membership check against a preallocated
// package-level slice for efficient, zero-allocation validation.
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
// ParseUsageProfile converts s to the corresponding pbc.UsageProfile.
// If s is empty it returns USAGE_PROFILE_UNSPECIFIED.
// Matching is case-insensitive for known profile names.
// If s is unrecognized it returns USAGE_PROFILE_UNSPECIFIED and an error.
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
// UsageProfileString returns the lowercase string representation of the given UsageProfile.
// For known profiles this yields "unspecified", "prod", "dev", or "burst". For unknown values it
// returns "unknown(<numeric>)" where <numeric> is the profile's integer value.
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
// NormalizeUsageProfile returns the provided profile if it is a known UsageProfile; otherwise it logs a warning with the numeric value and returns UsageProfile_USAGE_PROFILE_UNSPECIFIED.
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
// DefaultMonthlyHours returns the default monthly hours assumed for the given UsageProfile.
// For PROD it returns HoursProd (730), DEV returns HoursDev (160), and BURST returns HoursBurst (200).
// UNSPECIFIED and any unknown or future profiles default to HoursProd (production assumptions).
// Plugins may use different values based on their resource types.
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