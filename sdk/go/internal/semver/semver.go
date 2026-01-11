// Package semver provides shared semantic versioning validation for the PulumiCost SDK.
//
// This package exists to break circular dependencies between sdk/go/pluginsdk and
// sdk/go/testing. Both packages need semantic version validation but importing
// between them would create a circular dependency.
//
// # Usage
//
// This is an internal package. External consumers should use:
//   - [github.com/rshade/finfocus-spec/sdk/go/pluginsdk.ValidateSpecVersion]
//   - [github.com/rshade/finfocus-spec/sdk/go/testing.IsValidSemVer]
package semver

import "regexp"

// Regex is a precompiled regular expression for validating semantic versions.
// Matches: v1.2.3, v0.0.1, v10.20.30, etc.
// Does not match: 1.2.3 (no v prefix), v1.2 (missing patch), v1.2.3.4 (too many parts).
//
// The regex validates strict semantic versions in vMAJOR.MINOR.PATCH format only.
// Prerelease and build metadata (e.g., v1.0.0-alpha, v1.0.0+build) are intentionally
// excluded to ensure deterministic compatibility checking. This prevents ambiguity
// in version comparison (e.g., is v1.0.0-alpha < v1.0.0? Should v1.0.0-beta work
// with a core requiring v1.0.0?). By requiring exact vMAJOR.MINOR.PATCH format,
// version compatibility is unambiguous and consistently sortable.
//
// Additional validation rules:
//   - Required 'v' prefix
//   - No leading zeros (except for 0 itself)
//
// Example valid versions: "v0.4.11", "v1.0.0", "v2.15.3".
// Example invalid versions: "0.4.11" (no v prefix), "v1.2" (missing patch),
// "v01.2.3" (leading zero), "v1.0.0-alpha" (prerelease), "v1.0.0+build" (metadata).
var Regex = regexp.MustCompile(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)

// IsValid returns true if the version string is a valid semantic version.
// The version must be in the format vMAJOR.MINOR.PATCH where MAJOR, MINOR, and PATCH
// are non-negative integers without leading zeros (except for 0 itself).
func IsValid(version string) bool {
	if version == "" {
		return false
	}
	return Regex.MatchString(version)
}
