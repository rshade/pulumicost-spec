// Package semver provides shared semantic versioning validation for the FinFocus SDK.
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

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

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

// Sentinel errors for semantic version validation.
// These allow callers to programmatically identify specific validation failures.
var (
	// ErrEmptyVersion is returned when the version string is empty.
	ErrEmptyVersion = errors.New("version is empty")

	// ErrMissingVPrefix is returned when the version doesn't start with 'v'.
	ErrMissingVPrefix = errors.New("version must start with 'v'")

	// ErrInvalidFormat is returned when the version doesn't match vMAJOR.MINOR.PATCH format.
	ErrInvalidFormat = errors.New("version does not match semantic versioning format vMAJOR.MINOR.PATCH")
)

// Validate returns an error with a contextual message if the version is invalid.
// This function provides more detailed error messages than IsValid, making it
// suitable for user-facing error reporting.
//
// Returned errors wrap sentinel errors for programmatic handling:
//   - [ErrEmptyVersion]: version string is empty
//   - [ErrMissingVPrefix]: version doesn't start with 'v' (wraps with version value)
//   - [ErrInvalidFormat]: version doesn't match vMAJOR.MINOR.PATCH (wraps with version value)
//
// Example:
//
//	if err := semver.Validate("1.0.0"); err != nil {
//	    if errors.Is(err, semver.ErrMissingVPrefix) {
//	        // Handle missing 'v' prefix specifically
//	    }
//	    // err.Error() = "version must start with 'v': \"1.0.0\""
//	}
func Validate(version string) error {
	if version == "" {
		return ErrEmptyVersion
	}
	if !strings.HasPrefix(version, "v") {
		return fmt.Errorf("%w: %q", ErrMissingVPrefix, version)
	}
	if !Regex.MatchString(version) {
		return fmt.Errorf("%w: %q", ErrInvalidFormat, version)
	}
	return nil
}
